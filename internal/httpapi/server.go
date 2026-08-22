package httpapi

import (
	"encoding/json"
	"lawdrive/internal/query"
	"lawdrive/internal/workflow"
	"net/http"
)

type Server struct {
	workflow *workflow.Service
	query    *query.Service
}

func New(w *workflow.Service, q *query.Service) *Server { return &Server{workflow: w, query: q} }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(http.StatusOK); rw.Write([]byte("ok")) })
	mux.HandleFunc("/cases", s.cases)
	return mux
}
func (s *Server) cases(rw http.ResponseWriter, r *http.Request) {
	rows, e := s.query.Cases(r.URL.Query().Get("status"))
	if e != nil {
		http.Error(rw, e.Error(), 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(rows)
}
