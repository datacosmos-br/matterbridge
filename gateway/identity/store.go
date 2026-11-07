package identity

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Store maps WhatsApp user IDs to corporate users.
type Store struct {
	path string
	mu   sync.Mutex
	data map[string]string
}

func New(path string) (*Store, error) {
	s := &Store{path: path, data: map[string]string{}}
	b, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(b, &s.data)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

func (s *Store) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

func (s *Store) Link(whatsappID, corporateUser string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[whatsappID] = corporateUser
	return s.save()
}

func (s *Store) Resolve(whatsappID string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.data[whatsappID]
	if !ok {
		return "", false, nil
	}
	return u, true, nil
}

func (s *Store) Close() error { return nil }

func (s *Store) String() string { return fmt.Sprintf("identity.Store") }
