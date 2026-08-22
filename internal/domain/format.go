package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func FormatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}
func ParseTimestamp(v string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", v, time.UTC)
}
func CaseHeader(c Case) string {
	return fmt.Sprintf("%s | %s | %s", c.ID, c.Title, StatusLabel(c.Status))
}
func FileHeader(f CaseFile) string {
	return fmt.Sprintf("%s | %s | v%d | %s", f.ID, f.Name, f.Version, FileKindLabel(f.Kind))
}
func PermissionHeader(p Permission) string {
	return fmt.Sprintf("%s/%s | %s", p.CaseID, p.UserID, RoleLabel(p.Role))
}
func AuditHeader(a AuditEntry) string {
	return fmt.Sprintf("%s %s %s", FormatTimestamp(a.At), ActionLabel(a.Action), a.Actor)
}
func ShareHeader(s ShareLink) string {
	return fmt.Sprintf("%s -> %s", s.ID, FormatTimestamp(s.ExpiresAt))
}
func JoinLabels(values []string) string {
	clean := make([]string, 0, len(values))
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			clean = append(clean, strings.TrimSpace(v))
		}
	}
	return strings.Join(clean, " / ")
}
func UniqueStrings(values []string) []string {
	set := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, v := range values {
		if !set[v] {
			set[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
func CaseKinds(files []CaseFile) []string {
	out := make([]string, 0)
	for _, f := range files {
		out = append(out, f.Kind)
	}
	return UniqueStrings(out)
}
func CaseActors(audits []AuditEntry) []string {
	out := make([]string, 0)
	for _, a := range audits {
		out = append(out, a.Actor)
	}
	return UniqueStrings(out)
}
func CaseIDs(cases []Case) []string {
	out := make([]string, 0)
	for _, c := range cases {
		out = append(out, c.ID)
	}
	return UniqueStrings(out)
}
func FileIDs(files []CaseFile) []string {
	out := make([]string, 0)
	for _, f := range files {
		out = append(out, f.ID)
	}
	return UniqueStrings(out)
}
func IsRecent(t, now time.Time, window time.Duration) bool {
	return !t.IsZero() && !t.After(now) && now.Sub(t) <= window
}
func ClampVersion(v int) int {
	if v < 1 {
		return 1
	}
	if v > 10000 {
		return 10000
	}
	return v
}
