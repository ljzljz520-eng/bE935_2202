package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
	"sort"
	"strings"
)

type CaseReport struct {
	Case       domain.Case
	Files      []domain.CaseFile
	TotalBytes int
	Kinds      map[string]int
	Complete   bool
}

func (s *Service) Report(caseID string) (CaseReport, error) {
	c, ok, e := s.repo.FindCase(caseID)
	if e != nil {
		return CaseReport{}, e
	}
	if !ok {
		return CaseReport{}, fmt.Errorf("case missing")
	}
	files, e := s.repo.FilesForCase(caseID)
	if e != nil {
		return CaseReport{}, e
	}
	return CaseReport{Case: c, Files: files, TotalBytes: filesSize(files), Kinds: kindTotals(files), Complete: IsArchiveComplete(c, files)}, nil
}
func filesSize(fs []domain.CaseFile) int {
	n := 0
	for _, f := range fs {
		n += len(f.Content)
	}
	return n
}
func kindTotals(fs []domain.CaseFile) map[string]int {
	m := map[string]int{}
	for _, f := range fs {
		m[f.Kind]++
	}
	return m
}
func (s *Service) CasesByClient(client string) ([]domain.Case, error) {
	all, e := s.repo.AllCases()
	if e != nil {
		return nil, e
	}
	out := make([]domain.Case, 0)
	for _, c := range all {
		if strings.EqualFold(c.Client, client) {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
func (s *Service) CloseCase(caseID string) error {
	c, ok, e := s.repo.FindCase(caseID)
	if e != nil || !ok {
		if e != nil {
			return e
		}
		return fmt.Errorf("case missing")
	}
	if e = domain.CheckTransition(c.Status, "closed"); e != nil {
		return e
	}
	c.Status = "closed"
	return s.repo.SaveCase(c)
}
func (s *Service) OpenCase(caseID string) error {
	c, ok, e := s.repo.FindCase(caseID)
	if e != nil || !ok {
		if e != nil {
			return e
		}
		return fmt.Errorf("case missing")
	}
	if e = domain.CheckTransition(c.Status, "open"); e != nil {
		return e
	}
	c.Status = "open"
	return s.repo.SaveCase(c)
}
func (s *Service) ValidateFileSet(caseID string) error {
	fs, e := s.repo.FilesForCase(caseID)
	if e != nil {
		return e
	}
	if len(fs) == 0 {
		return fmt.Errorf("case has no files")
	}
	for _, f := range fs {
		if f.Version < 1 {
			return fmt.Errorf("file %s has no version", f.ID)
		}
	}
	return nil
}
func (s *Service) ArchiveReady(caseID string) bool {
	c, ok, _ := s.repo.FindCase(caseID)
	if !ok {
		return false
	}
	fs, _ := s.repo.FilesForCase(caseID)
	return len(fs) > 0 && !c.Archived
}
func (s *Service) FormatReport(r CaseReport) string {
	return fmt.Sprintf("%s: %d files, %d bytes", r.Case.ID, len(r.Files), r.TotalBytes)
}
