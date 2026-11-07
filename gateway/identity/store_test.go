package identity

import (
	"os"
	"testing"
)

func TestStoreLinkResolve(t *testing.T) {
	f, _ := os.CreateTemp("", "id-*.db")
	f.Close()
	os.Remove(f.Name())
	s, err := New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Link("wa1", "corp1"); err != nil {
		t.Fatal(err)
	}
	u, ok, err := s.Resolve("wa1")
	if err != nil || !ok || u != "corp1" {
		t.Fatalf("resolve failed: %v %v %s", err, ok, u)
	}
	_, ok, _ = s.Resolve("missing")
	if ok {
		t.Fatal("expected not found")
	}
}
