package httpapi

import (
	"lawdrive/internal/query"
	"lawdrive/internal/store"
	"lawdrive/internal/workflow"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealth(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := store.NewRepository(s)
	w := workflow.New(r)
	q := query.New(r)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	New(w, q).Handler().ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatal(res.Code)
	}
}
