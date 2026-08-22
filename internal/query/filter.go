package query

import (
	"lawdrive/internal/domain"
	"sort"
)

func FilterCases(cases []domain.Case, filter domain.CaseFilter) []domain.Case {
	out := make([]domain.Case, 0)
	for _, c := range cases {
		if domain.MatchCase(c, filter) {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func PageCases(cases []domain.Case, page domain.Page) ([]domain.Case, domain.PageInfo) {
	p := page.Normalize()
	info := p.Info(len(cases))
	start := p.Offset()
	if start >= len(cases) {
		return []domain.Case{}, info
	}
	end := start + p.Size
	if end > len(cases) {
		end = len(cases)
	}
	return cases[start:end], info
}
func GroupFiles(files []domain.CaseFile) map[string][]domain.CaseFile {
	out := map[string][]domain.CaseFile{}
	for _, f := range files {
		out[f.CaseID] = append(out[f.CaseID], f)
	}
	for k, v := range out {
		out[k] = domain.SortFiles(v)
	}
	return out
}
func LatestOnly(files []domain.CaseFile) []domain.CaseFile {
	group := GroupFiles(files)
	out := make([]domain.CaseFile, 0, len(group))
	for _, fs := range group {
		if len(fs) == 0 {
			continue
		}
		latest := fs[0]
		for _, f := range fs[1:] {
			if f.UpdatedAt.After(latest.UpdatedAt) {
				latest = f
			}
		}
		out = append(out, latest)
	}
	return domain.SortFiles(out)
}
