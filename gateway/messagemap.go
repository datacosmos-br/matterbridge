package gateway

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"

	"github.com/matterbridge-org/matterbridge/gateway/sqlitedb"
)

// messageCacheSize is the capacity of the in-memory message map.
const messageCacheSize = 5000

// messageMap maps a canonical message ID ("<protocol> <ID>" of the source
// message) to the IDs of the messages relayed to each destination.
type messageMap interface {
	// get returns the destinations of a canonical message and whether the
	// canonical message is known.
	get(canonical string) ([]*BrMsgID, bool, error)
	// add stores the destinations of a canonical message, replacing any
	// previous entry.
	add(canonical string, ids []*BrMsgID) error
	// findCanonical returns the canonical ID of a canonical or downstream
	// message ID, or "" when unknown.
	findCanonical(id string) (string, error)
	close() error
}

// key uniquely identifies a destination entry.
func (b *BrMsgID) key() string {
	return strings.Join([]string{b.Protocol, b.DestName, b.ChannelID, b.ID}, " ")
}

// lruMessageMap is the default in-memory message map.
type lruMessageMap struct {
	cache *lru.Cache[string, []*BrMsgID]
}

func newLRUMessageMap() (*lruMessageMap, error) {
	cache, err := lru.New[string, []*BrMsgID](messageCacheSize)
	if err != nil {
		return nil, fmt.Errorf("create message cache: %w", err)
	}

	return &lruMessageMap{cache: cache}, nil
}

func (m *lruMessageMap) get(canonical string) ([]*BrMsgID, bool, error) {
	ids, ok := m.cache.Get(canonical)
	return ids, ok, nil
}

func (m *lruMessageMap) add(canonical string, ids []*BrMsgID) error {
	m.cache.Add(canonical, ids)
	return nil
}

func (m *lruMessageMap) findCanonical(id string) (string, error) {
	if m.cache.Contains(id) {
		return id, nil
	}
	// If not keyed, iterate through cache for downstream, and infer upstream.
	for _, canonical := range m.cache.Keys() {
		ids, _ := m.cache.Peek(canonical)
		for _, downstream := range ids {
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
	destinations TEXT NOT NULL
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
// messages holds the JSON-encoded destinations of each canonical message;
// downstream indexes every destination message ID to its canonical message.
type sqliteMessageMap struct {
	db *sql.DB

	selectMsg, selectByID, upsertMsg, deleteDownstream, insertDownstream *sql.Stmt
}

// messageMapPath returns the database file of gateway name below dir.
func messageMapPath(dir, name string) string {
	return filepath.Join(dir, url.PathEscape(name)+".sqlite")
}

func openSQLiteMessageMap(path string) (*sqliteMessageMap, error) {
	db, err := sqlitedb.Open(path, messageMapSchema)
	if err != nil {
		return nil, fmt.Errorf("persistent message store: %w", err)
	}

	m := &sqliteMessageMap{db: db}

	err = m.prepareAll()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("persistent message store %s: %w", path, err), m.close())
	}

	return m, nil
}

// prepareAll prepares the statements of the message map; close() closes them.
func (m *sqliteMessageMap) prepareAll() error {
	var err error

	m.selectMsg, err = sqlitedb.Prepare(m.db, `SELECT destinations FROM messages WHERE canonical = ?`)
	if err != nil {
		return err
	}

	m.selectByID, err = sqlitedb.Prepare(m.db, `SELECT canonical FROM downstream WHERE id = ? LIMIT 1`)
	if err != nil {
		return err
	}

	m.upsertMsg, err = sqlitedb.Prepare(m.db, `INSERT INTO messages (canonical, destinations) VALUES (?, ?)
		ON CONFLICT (canonical) DO UPDATE SET destinations = excluded.destinations`)
	if err != nil {
		return err
	}

	m.deleteDownstream, err = sqlitedb.Prepare(m.db, `DELETE FROM downstream WHERE canonical = ?`)
	if err != nil {
		return err
	}

	m.insertDownstream, err = sqlitedb.Prepare(m.db, `INSERT OR REPLACE INTO downstream (entry_key, id, canonical) VALUES (?, ?, ?)`)

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

func (m *sqliteMessageMap) add(canonical string, ids []*BrMsgID) error {
	encoded, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("encode destinations of %s: %w", canonical, err)
	}

	ctx := context.Background()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store message %s: %w", canonical, err)
	}

	err = m.addTx(ctx, tx, canonical, string(encoded), ids)
	if err != nil {
		return errors.Join(fmt.Errorf("store message %s: %w", canonical, err), tx.Rollback())
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("store message %s: %w", canonical, err)
	}

	return nil
}

func (m *sqliteMessageMap) addTx(ctx context.Context, tx *sql.Tx, canonical, encoded string, ids []*BrMsgID) error {
	err := execTx(ctx, tx, m.upsertMsg, canonical, encoded)
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
	return sqlitedb.CloseAll(m.db, m.selectMsg, m.selectByID, m.upsertMsg, m.deleteDownstream, m.insertDownstream)
}
