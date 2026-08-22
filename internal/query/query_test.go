package query

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestCasesAndFiles(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	r.SaveCase(domain.Case{ID: "A", Title: "A", Status: "open"})
	r.SaveFile(domain.CaseFile{ID: "F", CaseID: "A", Name: "contract.docx"})
	q := New(r)
	cases, e := q.Cases("open")
	if e != nil || len(cases) != 1 {
		t.Fatal(e)
	}
	fs, e := q.FindFiles("A", "contract")
	if e != nil || len(fs) != 1 {
		t.Fatal(e)
	}
}
