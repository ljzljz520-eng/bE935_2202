package query

import (
	"lawdrive/internal/domain"
	"strings"
)

type Index struct {
	byName map[string][]string
	byCase map[string][]string
}

func NewIndex() *Index { return &Index{byName: map[string][]string{}, byCase: map[string][]string{}} }
func (i *Index) Add(f domain.CaseFile) {
	key := strings.ToLower(f.Name)
	i.byName[key] = append(i.byName[key], f.ID)
	i.byCase[f.CaseID] = append(i.byCase[f.CaseID], f.ID)
}
func (i *Index) Remove(f domain.CaseFile) {
	key := strings.ToLower(f.Name)
	i.byName[key] = removeID(i.byName[key], f.ID)
	i.byCase[f.CaseID] = removeID(i.byCase[f.CaseID], f.ID)
}
func removeID(ids []string, id string) []string {
	out := ids[:0]
	for _, v := range ids {
		if v != id {
			out = append(out, v)
		}
	}
	return out
}
func (i *Index) FindName(name string) []string {
	return append([]string(nil), i.byName[strings.ToLower(name)]...)
}
func (i *Index) FindCase(caseID string) []string { return append([]string(nil), i.byCase[caseID]...) }
func (i *Index) Count() int {
	n := 0
	for _, ids := range i.byName {
		n += len(ids)
	}
	return n
}
func (i *Index) Rebuild(files []domain.CaseFile) {
	i.byName = map[string][]string{}
	i.byCase = map[string][]string{}
	for _, f := range files {
		i.Add(f)
	}
}
func SortByVersion(files []domain.CaseFile) []domain.CaseFile {
	out := append([]domain.CaseFile(nil), files...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Version > out[i].Version {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
