package lawdrive

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cases.db")
	s, e := store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	r := store.NewRepository(s)
	if e = r.SaveCase(domain.Case{ID: "A-52", Title: "Archive", Status: "open"}); e != nil {
		t.Fatal(e)
	}
	if e = r.SaveFile(domain.CaseFile{ID: "F-1", CaseID: "A-52", Name: "evidence.pdf", Kind: "evidence", Content: []byte("proof"), Version: 1}); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r = store.NewRepository(s)
	c, ok, e := r.FindCase("A-52")
	if e != nil || !ok || c.Title != "Archive" {
		t.Fatalf("%+v %v %v", c, ok, e)
	}
	f, ok, e := r.FindFile("F-1")
	if e != nil || !ok || string(f.Content) != "proof" {
		t.Fatalf("%+v %v %v", f, ok, e)
	}
}
