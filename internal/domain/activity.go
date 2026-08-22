package domain

import (
	"fmt"
	"strings"
)

type CaseFilter struct {
	Status   string
	Client   string
	Archived *bool
	Text     string
}
type FileFilter struct {
	CaseID     string
	Kind       string
	Archived   *bool
	Text       string
	MinVersion int
}
type Page struct {
	Number int
	Size   int
}
type PageInfo struct {
	Number int
	Size   int
	Total  int
	Pages  int
}

func (p Page) Normalize() Page {
	if p.Number < 1 {
		p.Number = 1
	}
	if p.Size < 1 {
		p.Size = 20
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p
}
func (p Page) Offset() int { n := p.Normalize(); return (n.Number - 1) * n.Size }
func (p Page) Info(total int) PageInfo {
	n := p.Normalize()
	pages := total / n.Size
	if total%n.Size != 0 {
		pages++
	}
	return PageInfo{Number: n.Number, Size: n.Size, Total: total, Pages: pages}
}
func MatchCase(c Case, f CaseFilter) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Client != "" && !strings.Contains(strings.ToLower(c.Client), strings.ToLower(f.Client)) {
		return false
	}
	if f.Archived != nil && c.Archived != *f.Archived {
		return false
	}
	if f.Text != "" && !strings.Contains(strings.ToLower(c.Title), strings.ToLower(f.Text)) {
		return false
	}
	return true
}
func MatchFile(f CaseFile, filter FileFilter) bool {
	if filter.CaseID != "" && f.CaseID != filter.CaseID {
		return false
	}
	if filter.Kind != "" && f.Kind != filter.Kind {
		return false
	}
	if filter.Archived != nil && f.Archived != *filter.Archived {
		return false
	}
	if filter.Text != "" && !strings.Contains(strings.ToLower(f.Name), strings.ToLower(filter.Text)) {
		return false
	}
	if filter.MinVersion > 0 && f.Version < filter.MinVersion {
		return false
	}
	return true
}
func CheckTransition(from, to string) error {
	allowed := map[string][]string{"open": {"closed"}, "closed": {"open"}}
	for _, v := range allowed[from] {
		if v == to {
			return nil
		}
	}
	return fmt.Errorf("transition %s -> %s not allowed", from, to)
}
func FileKindLabel(kind string) string {
	switch kind {
	case "contract":
		return "合同"
	case "evidence":
		return "证据"
	case "minutes":
		return "会议纪要"
	default:
		return "其他"
	}
}
func SortFiles(files []CaseFile) []CaseFile {
	out := append([]CaseFile(nil), files...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Name < out[i].Name {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func SummarizeCase(c Case) string {
	state := "进行中"
	if c.Archived {
		state = "已归档"
	}
	return fmt.Sprintf("%s (%s) - %s", c.Title, c.Client, state)
}
