package files

import (
	"fmt"
	"lawdrive/internal/domain"
	"sort"
	"time"
)

type RetentionRule struct {
	Kind         string
	Days         int
	KeepVersions int
}
type RetentionDecision struct {
	FileID string
	Delete bool
	Reason string
}

func DefaultRules() []RetentionRule {
	return []RetentionRule{{Kind: "contract", Days: 3650, KeepVersions: 20}, {Kind: "evidence", Days: 3650, KeepVersions: 50}, {Kind: "minutes", Days: 1825, KeepVersions: 10}}
}
func RuleFor(rules []RetentionRule, kind string) RetentionRule {
	for _, r := range rules {
		if r.Kind == kind {
			return r
		}
	}
	return RetentionRule{Kind: kind, Days: 365, KeepVersions: 5}
}
func ShouldRetain(f domain.CaseFile, now time.Time, rules []RetentionRule) bool {
	r := RuleFor(rules, f.Kind)
	if f.UpdatedAt.IsZero() {
		return true
	}
	return now.Sub(f.UpdatedAt) <= time.Duration(r.Days)*24*time.Hour
}
func PlanRetention(files []domain.CaseFile, now time.Time, rules []RetentionRule) []RetentionDecision {
	out := make([]RetentionDecision, 0, len(files))
	for _, f := range files {
		keep := ShouldRetain(f, now, rules)
		reason := "within retention"
		if !keep {
			reason = "retention expired"
		}
		out = append(out, RetentionDecision{FileID: f.ID, Delete: !keep, Reason: reason})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FileID < out[j].FileID })
	return out
}
func ValidateRules(rules []RetentionRule) error {
	if len(rules) == 0 {
		return fmt.Errorf("retention rules required")
	}
	seen := map[string]bool{}
	for _, r := range rules {
		if r.Kind == "" || seen[r.Kind] {
			return fmt.Errorf("duplicate retention kind")
		}
		if r.Days < 1 || r.KeepVersions < 1 {
			return fmt.Errorf("invalid retention values")
		}
		seen[r.Kind] = true
	}
	return nil
}
func (s *Service) RetentionPlan(now time.Time) ([]RetentionDecision, error) {
	fs, e := s.repo.AllFiles()
	if e != nil {
		return nil, e
	}
	rules := DefaultRules()
	return PlanRetention(fs, now, rules), nil
}
func (s *Service) PurgeExpired(now time.Time) (int, error) {
	plan, e := s.RetentionPlan(now)
	if e != nil {
		return 0, e
	}
	n := 0
	for _, d := range plan {
		if d.Delete {
			if e := s.Archive(d.FileID); e != nil {
				return n, e
			}
			n++
		}
	}
	return n, nil
}
func VersionLimit(f domain.CaseFile, rules []RetentionRule) int {
	return RuleFor(rules, f.Kind).KeepVersions
}
func SortDecisions(ds []RetentionDecision) []RetentionDecision {
	out := append([]RetentionDecision(nil), ds...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Delete == out[j].Delete {
			return out[i].FileID < out[j].FileID
		}
		return out[i].Delete
	})
	return out
}
