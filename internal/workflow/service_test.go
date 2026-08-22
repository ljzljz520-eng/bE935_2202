package workflow

import (
	"lawdrive/internal/files"
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowUploadShareDownload(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	w := New(r)
	if _, e := w.CreateCase("A", "Matter", "Client"); e != nil {
		t.Fatal(e)
	}
	if _, e := w.GrantAccess("A", "lawyer", "lawyer", true, true, true); e != nil {
		t.Fatal(e)
	}
	if _, e := w.UploadAttachment("A", "F", "lawyer", "meeting.txt", "minutes", []byte("notes")); e != nil {
		t.Fatal(e)
	}
	if _, e := w.Share("A", "F", "lawyer"); e != nil {
		t.Fatal(e)
	}
	b, e := w.Download("A", "F", "lawyer")
	if e != nil || string(b) != "notes" {
		t.Fatalf("download %q %v", b, e)
	}
}

func TestWorkflowThreeEditAndSearch(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	w := New(r)
	w.CreateCase("A", "Search Matter", "C")
	w.GrantAccess("A", "p", "partner", true, true, true)
	w.UploadAttachment("A", "F", "p", "contract.docx", "contract", []byte("alpha"))
	w.EditAttachment("F", "p", []byte("beta"))
	rows, e := files.NewSearcher(r).Search("beta", false)
	if e != nil || len(rows) != 1 {
		t.Fatalf("search %d %v", len(rows), e)
	}
}
