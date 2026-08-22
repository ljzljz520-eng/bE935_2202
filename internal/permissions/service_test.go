package permissions

import (
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestGrantAndCheck(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	p := New(r)
	if _, err := p.Grant("A", "u", "lawyer", true, true, true); err != nil {
		t.Fatal(err)
	}
	ok, err := p.Check("A", "u", "share")
	if err != nil || !ok {
		t.Fatalf("share denied: %v", err)
	}
}
