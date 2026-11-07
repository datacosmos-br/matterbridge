package gateway

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"

	_ "modernc.org/sqlite" // database/sql driver "sqlite"
)

// messageCacheSize is the capacity of the in-memory message map.
const messageCacheSize = 5000

// messageMap maps a canonical message ID ("<protocol> <ID>" of the source
// message) to the IDs of the messages relayed to each destination.
//
// The source message itself is kept apart from the destinations: it is a
// delete target (a delete made on a destination reaches it) but never an edit
// target, since a bridge cannot edit a message another user posted.
type messageMap interface {
	// get returns the destinations of a canonical message and whether the
	// canonical message is known.
	get(canonical string) ([]*BrMsgID, bool, error)
	// source returns the source message of a canonical message and whether
	// it is recorded.
	source(canonical string) (*BrMsgID, bool, error)
	// add stores the source and destinations of a canonical message,
	// replacing any previous entry.
	add(canonical string, source *BrMsgID, ids []*BrMsgID) error
	// findCanonical returns the canonical ID of a canonical or downstream
	// message ID, or "" when unknown.
	findCanonical(id string) (string, error)
	close() error
}

// key uniquely identifies a destination entry.
func (b *BrMsgID) key() string {
	return strings.Join([]string{b.Protocol, b.DestName, b.ChannelID, b.ID}, " ")
}

// messageEntry is the in-memory record of a canonical message.
type messageEntry struct {
	source *BrMsgID
	ids    []*BrMsgID
}

// lruMessageMap is the default in-memory message map.
type lruMessageMap struct {
	cache *lru.Cache[string, messageEntry]
}

func newLRUMessageMap() (*lruMessageMap, error) {
	cache, err := lru.New[string, messageEntry](messageCacheSize)
	if err != nil {
		return nil, fmt.Errorf("create message cache: %w", err)
	}

	return &lruMessageMap{cache: cache}, nil
}

func (m *lruMessageMap) get(canonical string) ([]*BrMsgID, bool, error) {
	entry, ok := m.cache.Get(canonical)
	return entry.ids, ok, nil
}

func (m *lruMessageMap) source(canonical string) (*BrMsgID, bool, error) {
	entry, _ := m.cache.Get(canonical)
	return entry.source, entry.source != nil, nil
}

func (m *lruMessageMap) add(canonical string, source *BrMsgID, ids []*BrMsgID) error {
	m.cache.Add(canonical, messageEntry{source: source, ids: ids})
	return nil
}

func (m *lruMessageMap) findCanonical(id string) (string, error) {
	if m.cache.Contains(id) {
		return id, nil
	}
	// If not keyed, iterate through cache for downstream, and infer upstream.
	for _, canonical := range m.cache.Keys() {
		entry, _ := m.cache.Peek(canonical)
		for _, downstream := range entry.ids {
			if id == downstream.ID {
				return canonical, nil
			}
		}
	}

	return "", nil
}

func (m *lruMessageMap) close() error {
	return nil
}

const messageMapSchema = `
CREATE TABLE IF NOT EXISTS messages (
	canonical    TEXT PRIMARY KEY,
	destinations TEXT NOT NULL,
	source       TEXT NOT NULL DEFAULT ''
) WITHOUT ROWID;
CREATE TABLE IF NOT EXISTS downstream (
	entry_key TEXT PRIMARY KEY,
	id        TEXT NOT NULL,
	canonical TEXT NOT NULL REFERENCES messages(canonical) ON DELETE CASCADE
) WITHOUT ROWID;
CREATE INDEX IF NOT EXISTS downstream_id ON downstream(id);
CREATE INDEX IF NOT EXISTS downstream_canonical ON downstream(canonical);
`

// sqliteMessageMap persists the message map in one sqlite database per
// gateway (PersistentMessageStorePath).
//
// messages holds the JSON-encoded destinations and source (empty when not
// recorded, as in databases written before the source column) of each
// canonical message;
// downstream indexes every destination message ID to its canonical message.
type sqliteMessageMap struct {
	db *sql.DB

	selectMsg, selectSource, selectByID, upsertMsg, deleteDownstream, insertDownstream *sql.Stmt
}

// openSQLite opens (creating if needed) the sqlite database at path in WAL
// mode and applies schema.
func openSQLite(path, schema string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		return nil, fmt.Errorf("create directory of %s: %w", path, err)
	}

	dsn := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create schema in %s: %w", path, err), db.Close())
	}

	return db, nil
}

// prepare returns a prepared statement owned by the caller.
func prepare(db *sql.DB, query string) (*sql.Stmt, error) {
	stmt, err := db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("prepare %q: %w", query, err)
	}

	return stmt, nil
}

// closeAll closes the statements that were prepared and then the database.
func closeAll(db *sql.DB, stmts ...*sql.Stmt) error {
	var errs []error

	for _, stmt := range stmts {
		if stmt != nil {
			errs = append(errs, stmt.Close())
		}
	}

	errs = append(errs, db.Close())

	return errors.Join(errs...)
}

// messageMapPath returns the database file of gateway name below dir.
func messageMapPath(dir, name string) string {
	return filepath.Join(dir, url.PathEscape(name)+".sqlite")
}

func openSQLiteMessageMap(path string) (*sqliteMessageMap, error) {
	db, err := openSQLite(path, messageMapSchema)
	if err != nil {
		return nil, fmt.Errorf("persistent message store: %w", err)
	}

	m := &sqliteMessageMap{db: db}

	err = addSourceColumn(db)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("persistent message store %s: %w", path, err), m.close())
	}

	err = m.prepareAll()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("persistent message store %s: %w", path, err), m.close())
	}

	return m, nil
}

// addSourceColumn adds the source column to a messages table created before
// it existed (CREATE TABLE IF NOT EXISTS leaves such a table unchanged).
func addSourceColumn(db *sql.DB) error {
	ctx := context.Background()

	var count int

	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('messages') WHERE name = 'source'`).Scan(&count)
	if err != nil {
		return fmt.Errorf("inspect messages table: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, err = db.ExecContext(ctx, `ALTER TABLE messages ADD COLUMN source TEXT NOT NULL DEFAULT ''`)
	if err != nil {
		return fmt.Errorf("add source column: %w", err)
	}

	return nil
}

// prepareAll prepares the statements of the message map; close() closes them.
func (m *sqliteMessageMap) prepareAll() error {
	var err error

	m.selectMsg, err = prepare(m.db, `SELECT destinations FROM messages WHERE canonical = ?`)
	if err != nil {
		return err
	}

	m.selectSource, err = prepare(m.db, `SELECT source FROM messages WHERE canonical = ?`)
	if err != nil {
		return err
	}

	m.selectByID, err = prepare(m.db, `SELECT canonical FROM downstream WHERE id = ? LIMIT 1`)
	if err != nil {
		return err
	}

	m.upsertMsg, err = prepare(m.db, `INSERT INTO messages (canonical, destinations, source) VALUES (?, ?, ?)
		ON CONFLICT (canonical) DO UPDATE SET destinations = excluded.destinations, source = excluded.source`)
	if err != nil {
		return err
	}

	m.deleteDownstream, err = prepare(m.db, `DELETE FROM downstream WHERE canonical = ?`)
	if err != nil {
		return err
	}

	m.insertDownstream, err = prepare(m.db, `INSERT OR REPLACE INTO downstream (entry_key, id, canonical) VALUES (?, ?, ?)`)

	return err
}

// queryString returns the single string column selected by stmt, and whether
// a row matched.
func queryString(stmt *sql.Stmt, arg string) (string, bool, error) {
	var value string

	err := stmt.QueryRowContext(context.Background(), arg).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	return value, true, nil
}

func (m *sqliteMessageMap) get(canonical string) ([]*BrMsgID, bool, error) {
	encoded, found, err := queryString(m.selectMsg, canonical)
	if err != nil {
		return nil, false, fmt.Errorf("read message %s: %w", canonical, err)
	}

	if !found {
		return nil, false, nil
	}

	var ids []*BrMsgID

	err = json.Unmarshal([]byte(encoded), &ids)
	if err != nil {
		return nil, false, fmt.Errorf("decode destinations of %s: %w", canonical, err)
	}

	return ids, true, nil
}

func (m *sqliteMessageMap) source(canonical string) (*BrMsgID, bool, error) {
	encoded, _, err := queryString(m.selectSource, canonical)
	if err != nil {
		return nil, false, fmt.Errorf("read source of %s: %w", canonical, err)
	}

	if encoded == "" {
		return nil, false, nil
	}

	var source BrMsgID

	err = json.Unmarshal([]byte(encoded), &source)
	if err != nil {
		return nil, false, fmt.Errorf("decode source of %s: %w", canonical, err)
	}

	return &source, true, nil
}

func (m *sqliteMessageMap) add(canonical string, source *BrMsgID, ids []*BrMsgID) error {
	encoded, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("encode destinations of %s: %w", canonical, err)
	}

	var encodedSource []byte
	if source != nil {
		encodedSource, err = json.Marshal(source)
		if err != nil {
			return fmt.Errorf("encode source of %s: %w", canonical, err)
		}
	}

	ctx := context.Background()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store message %s: %w", canonical, err)
	}

	err = m.addTx(ctx, tx, canonical, string(encoded), string(encodedSource), ids)
	if err != nil {
		return errors.Join(fmt.Errorf("store message %s: %w", canonical, err), tx.Rollback())
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("store message %s: %w", canonical, err)
	}

	return nil
}

func (m *sqliteMessageMap) addTx(ctx context.Context, tx *sql.Tx, canonical, encoded, encodedSource string, ids []*BrMsgID) error {
	err := execTx(ctx, tx, m.upsertMsg, canonical, encoded, encodedSource)
	if err != nil {
		return err
	}

	err = execTx(ctx, tx, m.deleteDownstream, canonical)
	if err != nil {
		return err
	}

	for _, id := range ids {
		err = execTx(ctx, tx, m.insertDownstream, id.key(), id.ID, canonical)
		if err != nil {
			return err
		}
	}

	return nil
}

// execTx runs a prepared statement inside tx.
func execTx(ctx context.Context, tx *sql.Tx, stmt *sql.Stmt, args ...any) (err error) {
	txStmt := tx.StmtContext(ctx, stmt)
	defer func() { err = errors.Join(err, txStmt.Close()) }()

	_, err = txStmt.ExecContext(ctx, args...)

	return err
}

func (m *sqliteMessageMap) findCanonical(id string) (string, error) {
	_, found, err := queryString(m.selectMsg, id)
	if err != nil {
		return "", fmt.Errorf("read message %s: %w", id, err)
	}

	if found {
		return id, nil
	}

	canonical, _, err := queryString(m.selectByID, id)
	if err != nil {
		return "", fmt.Errorf("look up canonical message of %s: %w", id, err)
	}

	return canonical, nil
}

func (m *sqliteMessageMap) close() error {
	return closeAll(m.db, m.selectMsg, m.selectSource, m.selectByID, m.upsertMsg, m.deleteDownstream, m.insertDownstream)
}
