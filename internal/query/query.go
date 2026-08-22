package query

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"sort"
	"strings"
)

type Service struct{ repo *store.Repository }

func New(repo *store.Repository) *Service { return &Service{repo: repo} }
func (s *Service) Cases(status string) ([]domain.Case, error) {
	rows, e := s.repo.AllCases()
	if e != nil {
		return nil, e
	}
	out := make([]domain.Case, 0)
	for _, c := range rows {
		if status == "" || c.Status == status {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *Service) FindFiles(caseID, name string) ([]domain.CaseFile, error) {
	all, e := s.repo.AllFiles()
	if e != nil {
		return nil, e
	}
	out := make([]domain.CaseFile, 0)
	for _, f := range all {
		if caseID != "" && f.CaseID != caseID {
			continue
		}
		if name != "" && !strings.Contains(strings.ToLower(f.Name), strings.ToLower(name)) {
			continue
		}
		out = append(out, f)
	}
	return out, nil
}
