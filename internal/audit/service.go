package audit

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"sort"
)

type Service struct {
	repo     *store.Repository
	sequence int
}

func New(repo *store.Repository) *Service { return &Service{repo: repo} }
func (s *Service) Record(event domain.Event) (domain.AuditEntry, error) {
	if !domain.EventNeedsAudit(event.Name) {
		return domain.AuditEntry{}, nil
	}
	s.sequence++
	a := domain.AuditEntry{ID: fmt.Sprintf("audit-%06d", s.sequence), CaseID: event.CaseID, FileID: event.FileID, Actor: event.Actor, Action: event.Name, Detail: event.Data}
	return a, s.repo.SaveAudit(a)
}
func (s *Service) List(caseID string) ([]domain.AuditEntry, error) {
	rows, e := s.repo.AllAudits()
	if e != nil {
		return nil, e
	}
	filtered := make([]domain.AuditEntry, 0)
	for _, a := range rows {
		if caseID == "" || a.CaseID == caseID {
			filtered = append(filtered, a)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })
	return filtered, nil
}
func (s *Service) Count(caseID string) (int, error) { rows, e := s.List(caseID); return len(rows), e }
