package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
	"strings"
)

type WorkflowMetrics struct {
	Uploaded   int
	Edited     int
	Shared     int
	Downloaded int
	Archived   int
	Denied     int
}

func (m WorkflowMetrics) Total() int {
	return m.Uploaded + m.Edited + m.Shared + m.Downloaded + m.Archived + m.Denied
}
func (m WorkflowMetrics) Successes() int {
	return m.Uploaded + m.Edited + m.Shared + m.Downloaded + m.Archived
}
func (m WorkflowMetrics) FailureRate() float64 {
	if m.Total() == 0 {
		return 0
	}
	return float64(m.Denied) / float64(m.Total())
}
func (m *WorkflowMetrics) Record(action string, ok bool) {
	if !ok {
		m.Denied++
		return
	}
	switch action {
	case "upload":
		m.Uploaded++
	case "edit":
		m.Edited++
	case "share":
		m.Shared++
	case "download":
		m.Downloaded++
	case "archive":
		m.Archived++
	}
}
func (m WorkflowMetrics) Summary() string {
	return fmt.Sprintf("success=%d denied=%d", m.Successes(), m.Denied)
}
func ActionSequence(events []domain.Event) []string {
	out := make([]string, 0, len(events))
	for _, e := range events {
		out = append(out, e.Name)
	}
	return out
}
func HasAction(events []domain.Event, action string) bool {
	for _, e := range events {
		if e.Name == action {
			return true
		}
	}
	return false
}
func SequenceComplete(events []domain.Event) bool {
	needed := []string{"upload", "edit", "share"}
	for _, n := range needed {
		if !HasAction(events, n) {
			return false
		}
	}
	return true
}
func SequenceLabel(events []domain.Event) string { return strings.Join(ActionSequence(events), " -> ") }
func MergeMetrics(a, b WorkflowMetrics) WorkflowMetrics {
	return WorkflowMetrics{Uploaded: a.Uploaded + b.Uploaded, Edited: a.Edited + b.Edited, Shared: a.Shared + b.Shared, Downloaded: a.Downloaded + b.Downloaded, Archived: a.Archived + b.Archived, Denied: a.Denied + b.Denied}
}
func (s *Service) CaseMetrics(caseID string) (WorkflowMetrics, error) {
	audits, e := s.repo.AllAudits()
	if e != nil {
		return WorkflowMetrics{}, e
	}
	m := WorkflowMetrics{}
	for _, a := range audits {
		if a.CaseID == caseID {
			m.Record(a.Action, true)
		}
	}
	return m, nil
}
func IsHealthy(m WorkflowMetrics) bool     { return m.Denied == 0 && m.Successes() > 0 }
func NormalizeAction(action string) string { return strings.ToLower(strings.TrimSpace(action)) }
func AcceptableAction(action string) bool {
	switch NormalizeAction(action) {
	case "upload", "edit", "share", "download", "archive":
		return true
	default:
		return false
	}
}

func (m WorkflowMetrics) NeedsReview() bool { return m.Denied > 0 || m.FailureRate() > 0.1 }
