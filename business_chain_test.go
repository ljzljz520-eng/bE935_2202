package lawdrive

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"lawdrive/internal/workflow"
	"path/filepath"
	"testing"
)

func TestBusinessChain06(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	w := workflow.New(r)
	if _, e := w.CreateCase("A-52", "Archive Matter", "Client"); e != nil {
		t.Fatal(e)
	}
	for i := 1; i <= 6; i++ {
		id := fmt.Sprintf("A52-F-%d", i)
		if _, e := w.UploadAttachment("A-52", id, "partner", fmt.Sprintf("attachment-%d.pdf", i), "evidence", []byte("record")); e != nil {
			t.Fatal(e)
		}
	}
	rows, e := w.ArchiveCase(domain.ArchiveRequest{CaseID: "A-52", Actor: "partner", Reason: "closed"})
	if e != nil {
		t.Fatal(e)
	}
	if len(rows) != 6 {
		t.Fatalf("archive returned %d records, want 6", len(rows))
	}
	for _, f := range rows {
		if !f.Archived {
			t.Fatalf("file %s status disagrees with archived case", f.ID)
		}
	}
}
