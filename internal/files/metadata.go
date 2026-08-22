package files

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"path"
	"strings"
)

type Metadata struct {
	Extension  string
	Size       int
	Label      string
	Searchable bool
}

func ExtractMetadata(f domain.CaseFile) Metadata {
	ext := strings.ToLower(path.Ext(f.Name))
	return Metadata{Extension: ext, Size: len(f.Content), Label: domain.FileKindLabel(f.Kind), Searchable: len(f.Content) > 0}
}
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name required")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("path separators not allowed")
	}
	return nil
}
func NormalizeName(name string) string { return strings.TrimSpace(strings.ReplaceAll(name, "  ", " ")) }
func (s *Service) Metadata(fileID string) (Metadata, error) {
	f, ok, e := s.repo.FindFile(fileID)
	if e != nil {
		return Metadata{}, e
	}
	if !ok {
		return Metadata{}, fmt.Errorf("file missing")
	}
	return ExtractMetadata(f), nil
}
func (s *Service) FilesForCase(caseID string) ([]domain.CaseFile, error) {
	return s.repo.FilesForCase(caseID)
}
func (s *Service) Filter(files []domain.CaseFile, filter domain.FileFilter) []domain.CaseFile {
	out := make([]domain.CaseFile, 0)
	for _, f := range files {
		if domain.MatchFile(f, filter) {
			out = append(out, f)
		}
	}
	return domain.SortFiles(out)
}
func BuildManifest(files []domain.CaseFile) map[string]Metadata {
	m := make(map[string]Metadata, len(files))
	for _, f := range files {
		m[f.ID] = ExtractMetadata(f)
	}
	return m
}
func KindCounts(files []domain.CaseFile) map[string]int {
	m := map[string]int{}
	for _, f := range files {
		m[f.Kind]++
	}
	return m
}
func TotalSize(files []domain.CaseFile) int {
	total := 0
	for _, f := range files {
		total += len(f.Content)
	}
	return total
}
func EnsureRepository(r *store.Repository) error {
	if r == nil {
		return fmt.Errorf("repository required")
	}
	return nil
}
