package files

import (
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"sort"
	"strings"
)

type Searcher struct{ repo *store.Repository }

func NewSearcher(repo *store.Repository) *Searcher { return &Searcher{repo: repo} }

func (s *Searcher) Search(query string, includeArchived bool) ([]domain.SearchResult, error) {
	files, err := s.repo.AllFiles()
	if err != nil {
		return nil, err
	}
	results := make([]domain.SearchResult, 0)
	for _, f := range files {
		if f.Archived && !includeArchived {
			continue
		}
		if !matches(f, query) {
			continue
		}
		c, _, _ := s.repo.FindCase(f.CaseID)
		score := scoreFile(f, query)
		results = append(results, domain.SearchResult{File: f, CaseTitle: c.Title, Score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].File.ID < results[j].File.ID
		}
		return results[i].Score > results[j].Score
	})
	return results, nil
}

func matches(f domain.CaseFile, query string) bool {
	q := strings.ToLower(strings.TrimSpace(query))
	return q == "" || strings.Contains(strings.ToLower(f.Name), q) || strings.Contains(strings.ToLower(string(f.Content)), q)
}
func scoreFile(f domain.CaseFile, query string) int {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return 1
	}
	if strings.EqualFold(f.Name, q) {
		return 100
	}
	if strings.Contains(strings.ToLower(f.Name), q) {
		return 50
	}
	return 10
}
