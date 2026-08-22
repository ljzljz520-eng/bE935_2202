package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
	"strings"
)

type ChecklistItem struct {
	Code     string
	Label    string
	Required bool
	Complete bool
}
type Checklist struct {
	CaseID string
	Items  []ChecklistItem
}

func DefaultChecklist(caseID string) Checklist {
	return Checklist{CaseID: caseID, Items: []ChecklistItem{{Code: "client", Label: "client identity", Required: true}, {Code: "contract", Label: "signed contract", Required: true}, {Code: "evidence", Label: "evidence bundle", Required: true}, {Code: "minutes", Label: "meeting minutes", Required: false}}}
}
func (c *Checklist) Mark(code string, complete bool) error {
	for i := range c.Items {
		if c.Items[i].Code == code {
			c.Items[i].Complete = complete
			return nil
		}
	}
	return fmt.Errorf("checklist item %s missing", code)
}
func (c Checklist) Ready() bool {
	for _, i := range c.Items {
		if i.Required && !i.Complete {
			return false
		}
	}
	return true
}
func (c Checklist) Missing() []string {
	out := []string{}
	for _, i := range c.Items {
		if i.Required && !i.Complete {
			out = append(out, i.Code)
		}
	}
	return out
}
func (c Checklist) Summary() string {
	if c.Ready() {
		return "ready"
	}
	return "missing: " + strings.Join(c.Missing(), ", ")
}
func ValidateChecklist(c Checklist) error {
	if c.CaseID == "" {
		return fmt.Errorf("case id required")
	}
	seen := map[string]bool{}
	for _, i := range c.Items {
		if i.Code == "" || seen[i.Code] {
			return fmt.Errorf("invalid checklist")
		}
		seen[i.Code] = true
	}
	return nil
}
func BuildChecklistFromFiles(caseID string, files []domain.CaseFile) Checklist {
	c := DefaultChecklist(caseID)
	for _, f := range files {
		for i := range c.Items {
			if c.Items[i].Code == f.Kind {
				c.Items[i].Complete = true
			}
		}
	}
	return c
}
func (s *Service) Checklist(caseID string) (Checklist, error) {
	fs, e := s.repo.FilesForCase(caseID)
	if e != nil {
		return Checklist{}, e
	}
	c := BuildChecklistFromFiles(caseID, fs)
	return c, nil
}
func (s *Service) RequireReady(caseID string) error {
	c, e := s.Checklist(caseID)
	if e != nil {
		return e
	}
	if !c.Ready() {
		return fmt.Errorf("case not ready: %s", c.Summary())
	}
	return nil
}
func CompleteRequired(c Checklist) Checklist {
	for i := range c.Items {
		if c.Items[i].Required {
			c.Items[i].Complete = true
		}
	}
	return c
}
func ChecklistCodes(c Checklist) []string {
	out := make([]string, 0, len(c.Items))
	for _, i := range c.Items {
		out = append(out, i.Code)
	}
	return out
}
