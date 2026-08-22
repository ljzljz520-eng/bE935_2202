package store

import (
	"lawdrive/internal/domain"
	"path/filepath"
	"testing"
)

func TestStoreRoundTripAndBuckets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cases.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Put("cases", "A-1", domain.Case{ID: "A-1", Title: "Matter", Status: "open"}); err != nil {
		t.Fatal(err)
	}
	var got domain.Case
	found, err := s.Get("cases", "A-1", &got)
	if err != nil || !found || got.ID != "A-1" {
		t.Fatalf("case lookup found=%v err=%v value=%+v", found, err, got)
	}
	count, err := s.Count("cases")
	if err != nil || count != 1 {
		t.Fatalf("case count=%d err=%v", count, err)
	}
	if err = s.Delete("cases", "A-1"); err != nil {
		t.Fatal(err)
	}
	found, err = s.Get("cases", "A-1", &got)
	if err != nil || found {
		t.Fatalf("deleted case found=%v err=%v", found, err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoryBatchAndSnapshot(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "cases.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r := NewRepository(s)
	caseRow := domain.Case{ID: "A-2", Title: "Batch Matter", Client: "Client", Status: "open"}
	fileRow := domain.CaseFile{ID: "F-2", CaseID: "A-2", Name: "evidence.pdf", Kind: "evidence", Content: []byte("proof"), Version: 1}
	if err = r.SaveCaseWithFiles(caseRow, []domain.CaseFile{fileRow}); err != nil {
		t.Fatal(err)
	}
	if err = r.SaveAudit(domain.AuditEntry{ID: "audit-1", CaseID: "A-2", FileID: "F-2", Action: "upload"}); err != nil {
		t.Fatal(err)
	}
	snapshot, err := r.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Cases) != 1 || len(snapshot.Files) != 1 || len(snapshot.Audits) != 1 {
		t.Fatalf("snapshot=%+v", snapshot)
	}
	files, err := r.FilesForCase("A-2")
	if err != nil || len(files) != 1 || files[0].ID != "F-2" {
		t.Fatalf("files=%+v err=%v", files, err)
	}
}
