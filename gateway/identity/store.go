package identity

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store maps WhatsApp user IDs to corporate users.
type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS whatsapp_user_map (
		whatsapp_id TEXT PRIMARY KEY,
		corporate_user TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func (s *Store) Link(whatsappID, corporateUser string) error {
	_, err := s.db.Exec(`INSERT INTO whatsapp_user_map (whatsapp_id, corporate_user) VALUES (?, ?)
		ON CONFLICT(whatsapp_id) DO UPDATE SET corporate_user=excluded.corporate_user`, whatsappID, corporateUser)
	return err
}

func (s *Store) Resolve(whatsappID string) (string, bool, error) {
	var u string
	err := s.db.QueryRow(`SELECT corporate_user FROM whatsapp_user_map WHERE whatsapp_id=?`, whatsappID).Scan(&u)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return u, true, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) String() string { return fmt.Sprintf("identity.Store") }
