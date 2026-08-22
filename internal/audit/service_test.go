package audit

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestRecordList(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	a := New(store.NewRepository(s))
	a.Record(domain.NewEvent("upload", "A", "F", "u", "ok"))
	rows, e := a.List("A")
	if e != nil || len(rows) != 1 {
		t.Fatalf("%d %v", len(rows), e)
	}
}
