package share

import (
	"lawdrive/internal/store"
	"path/filepath"
	"testing"
	"time"
)

func TestShareLifecycle(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	svc := New(store.NewRepository(s))
	l, e := svc.Create("F", "u", time.Hour)
	if e != nil {
		t.Fatal(e)
	}
	if e = svc.Validate(l, time.Unix(1700000001, 0)); e != nil {
		t.Fatal(e)
	}
	if e = svc.Revoke(l); e != nil {
		t.Fatal(e)
	}
}
