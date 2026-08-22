package httpapi

import (
	"encoding/json"
	"lawdrive/internal/domain"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func decodeJSON(r *http.Request, value any) error { return json.NewDecoder(r.Body).Decode(value) }
func method(w http.ResponseWriter, r *http.Request, allowed string) bool {
	if r.Method != allowed {
		w.Header().Set("Allow", allowed)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func caseSummary(c domain.Case) map[string]any {
	return map[string]any{"id": c.ID, "title": c.Title, "client": c.Client, "status": c.Status, "archived": c.Archived}
}
func fileSummary(f domain.CaseFile) map[string]any {
	return map[string]any{"id": f.ID, "case_id": f.CaseID, "name": f.Name, "kind": f.Kind, "version": f.Version, "archived": f.Archived, "size": len(f.Content)}
}
func (s *Server) routeCaseList(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	rows, e := s.query.Cases(r.URL.Query().Get("status"))
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]map[string]any, 0, len(rows))
	for _, c := range rows {
		out = append(out, caseSummary(c))
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) routeHealth(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (s *Server) routeNotFound(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}
func parseArchiveRequest(r *http.Request) (domain.ArchiveRequest, error) {
	var req domain.ArchiveRequest
	if e := decodeJSON(r, &req); e != nil {
		return req, e
	}
	return req, nil
}
