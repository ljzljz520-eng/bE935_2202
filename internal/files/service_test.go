package files

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
)

func TestUploadEditPreview(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	r.SaveCase(domain.Case{ID: "A", Title: "Matter", Status: "open"})
	f := New(r)
	if _, err := f.Upload("A", "F", "u", "contract.docx", "contract", []byte("draft")); err != nil {
		t.Fatal(err)
	}
	got, err := f.Edit("F", "u", []byte("final"))
	if err != nil || got.Version != 2 {
		t.Fatal(err)
	}
	text, err := f.Preview("F")
	if err != nil || text != "final" {
		t.Fatalf("preview %q", text)
	}
}
