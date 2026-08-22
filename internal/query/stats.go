package query

import (
	"lawdrive/internal/domain"
	"sort"
)

type Stats struct {
	Cases         int
	OpenCases     int
	ArchivedCases int
	Files         int
	Bytes         int
	Kinds         map[string]int
}

func BuildStats(cases []domain.Case, files []domain.CaseFile) Stats {
	s := Stats{Cases: len(cases), Files: len(files), Kinds: map[string]int{}}
	for _, c := range cases {
		if c.Archived {
			s.ArchivedCases++
		}
		if c.Status == "open" {
			s.OpenCases++
		}
	}
	for _, f := range files {
		s.Bytes += len(f.Content)
		s.Kinds[f.Kind]++
	}
	return s
}
func (s Stats) Completion() float64 {
	if s.Files == 0 {
		return 0
	}
	return float64(s.Bytes) / float64(s.Files)
}
func (s Stats) HasData() bool { return s.Cases > 0 || s.Files > 0 }
func (s Stats) Labels() []string {
	out := make([]string, 0, len(s.Kinds))
	for k := range s.Kinds {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func CountByStatus(cases []domain.Case) map[string]int {
	m := map[string]int{}
	for _, c := range cases {
		m[c.Status]++
	}
	return m
}
func CountByClient(cases []domain.Case) map[string]int {
	m := map[string]int{}
	for _, c := range cases {
		m[c.Client]++
	}
	return m
}
func CountByKind(files []domain.CaseFile) map[string]int {
	m := map[string]int{}
	for _, f := range files {
		m[f.Kind]++
	}
	return m
}
func LargestFile(files []domain.CaseFile) (domain.CaseFile, bool) {
	if len(files) == 0 {
		return domain.CaseFile{}, false
	}
	out := files[0]
	for _, f := range files[1:] {
		if len(f.Content) > len(out.Content) {
			out = f
		}
	}
	return out, true
}
func NewestFile(files []domain.CaseFile) (domain.CaseFile, bool) {
	if len(files) == 0 {
		return domain.CaseFile{}, false
	}
	out := files[0]
	for _, f := range files[1:] {
		if f.UpdatedAt.After(out.UpdatedAt) {
			out = f
		}
	}
	return out, true
}
func ArchivedRatio(cases []domain.Case) float64 {
	if len(cases) == 0 {
		return 0
	}
	n := 0
	for _, c := range cases {
		if c.Archived {
			n++
		}
	}
	return float64(n) / float64(len(cases))
}
func VersionAverage(files []domain.CaseFile) float64 {
	if len(files) == 0 {
		return 0
	}
	n := 0
	for _, f := range files {
		n += f.Version
	}
	return float64(n) / float64(len(files))
}
func SortCasesByTitle(cases []domain.Case) []domain.Case {
	out := append([]domain.Case(nil), cases...)
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}
func SortFilesBySize(files []domain.CaseFile) []domain.CaseFile {
	out := append([]domain.CaseFile(nil), files...)
	sort.Slice(out, func(i, j int) bool { return len(out[i].Content) > len(out[j].Content) })
	return out
}
