package audit

import (
	"encoding/json"
	"fmt"
	"lawdrive/internal/domain"
	"sort"
	"strings"
)

func Encode(entries []domain.AuditEntry) ([]byte, error) {
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	return json.Marshal(entries)
}
func Actions(entries []domain.AuditEntry) map[string]int {
	m := map[string]int{}
	for _, e := range entries {
		m[e.Action]++
	}
	return m
}
func Actors(entries []domain.AuditEntry) []string {
	set := map[string]bool{}
	for _, e := range entries {
		set[e.Actor] = true
	}
	out := make([]string, 0, len(set))
	for a := range set {
		out = append(out, a)
	}
	sort.Strings(out)
	return out
}
func Filter(entries []domain.AuditEntry, action, actor string) []domain.AuditEntry {
	out := make([]domain.AuditEntry, 0)
	for _, e := range entries {
		if action != "" && !strings.EqualFold(action, e.Action) {
			continue
		}
		if actor != "" && !strings.EqualFold(actor, e.Actor) {
			continue
		}
		out = append(out, e)
	}
	return out
}
func Summary(entries []domain.AuditEntry) string {
	if len(entries) == 0 {
		return "no activity"
	}
	return fmt.Sprintf("%d actions by %d actors", len(entries), len(Actors(entries)))
}
