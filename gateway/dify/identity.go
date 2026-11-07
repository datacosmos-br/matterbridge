package dify

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/matterbridge-org/matterbridge/gateway/sqlitedb"
)

const identitySchema = `
CREATE TABLE IF NOT EXISTS identities (
	account   TEXT NOT NULL,
	user_id   TEXT NOT NULL,
	corp_user TEXT NOT NULL,
	PRIMARY KEY (account, user_id)
) WITHOUT ROWID;
`

// IdentityStore maps the user ID of a chat user on an account (for example a
// WhatsApp phone number) to their corporate user (IdentityStorePath).
type IdentityStore struct {
	db *sql.DB

	selectUser, upsertUser *sql.Stmt
}

// OpenIdentityStore opens (creating if needed) the identity database at path.
func OpenIdentityStore(path string) (*IdentityStore, error) {
	db, err := sqlitedb.Open(path, identitySchema)
	if err != nil {
		return nil, fmt.Errorf("identity store: %w", err)
	}

	s := &IdentityStore{db: db}

	s.selectUser, err = sqlitedb.Prepare(db, `SELECT corp_user FROM identities WHERE account = ? AND user_id = ?`)
	if err == nil {
		s.upsertUser, err = sqlitedb.Prepare(db, `INSERT INTO identities (account, user_id, corp_user) VALUES (?, ?, ?)
			ON CONFLICT (account, user_id) DO UPDATE SET corp_user = excluded.corp_user`)
	}

	if err != nil {
		return nil, errors.Join(fmt.Errorf("identity store %s: %w", path, err), s.Close())
	}

	return s, nil
}

// Resolve returns the corporate user linked to userID on account, and false
// when the user is not linked.
func (s *IdentityStore) Resolve(account, userID string) (string, bool, error) {
	var corpUser string

	err := s.selectUser.QueryRowContext(context.Background(), account, userID).Scan(&corpUser)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}

	if err != nil {
		return "", false, fmt.Errorf("resolve identity of %s on %s: %w", userID, account, err)
	}

	return corpUser, true, nil
}

// Link links userID on account to corpUser, replacing any previous link.
func (s *IdentityStore) Link(l Link) error {
	_, err := s.upsertUser.ExecContext(context.Background(), l.Account, l.UserID, l.CorpUser)
	if err != nil {
		return fmt.Errorf("link %s on %s to %s: %w", l.UserID, l.Account, l.CorpUser, err)
	}

	return nil
}

// Close closes the database.
func (s *IdentityStore) Close() error {
	return sqlitedb.CloseAll(s.db, s.selectUser, s.upsertUser)
}

// ErrInvalidLink reports a malformed identity link.
var ErrInvalidLink = errors.New("invalid identity link")

// Link is one identity link, written as "account:user_id=corp_user".
type Link struct {
	Account  string
	UserID   string
	CorpUser string
}

// ParseLink parses "account:user_id=corp_user". The account is the
// "protocol.name" of a bridge and cannot contain ':'; the user ID is
// everything between the first ':' and the last '='.
func ParseLink(s string) (Link, error) {
	ref, corpUser, ok := cutLast(s, "=")
	if !ok {
		return Link{}, fmt.Errorf("%w %q: missing '=corp_user'", ErrInvalidLink, s)
	}

	account, userID, ok := strings.Cut(ref, ":")
	if !ok {
		return Link{}, fmt.Errorf("%w %q: missing 'account:'", ErrInvalidLink, s)
	}

	if account == "" || userID == "" || corpUser == "" {
		return Link{}, fmt.Errorf("%w %q: account, user_id and corp_user must not be empty", ErrInvalidLink, s)
	}

	return Link{Account: account, UserID: userID, CorpUser: corpUser}, nil
}

func cutLast(s, sep string) (string, string, bool) {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return s, "", false
	}

	return s[:i], s[i+len(sep):], true
}
