package permissions

import (
	"fmt"
	"lawdrive/internal/domain"
)

type Matrix struct{ entries map[string]domain.Permission }

func NewMatrix() *Matrix { return &Matrix{entries: map[string]domain.Permission{}} }
func (m *Matrix) Set(p domain.Permission) error {
	if p.CaseID == "" || p.UserID == "" {
		return fmt.Errorf("identity required")
	}
	m.entries[p.CaseID+"/"+p.UserID] = p
	return nil
}
func (m *Matrix) Get(caseID, userID string) (domain.Permission, bool) {
	p, ok := m.entries[caseID+"/"+userID]
	return p, ok
}
func (m *Matrix) Remove(caseID, userID string) { delete(m.entries, caseID+"/"+userID) }
func (m *Matrix) Users(caseID string) []string {
	out := []string{}
	prefix := caseID + "/"
	for k := range m.entries {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, k[len(prefix):])
		}
	}
	return out
}
func (m *Matrix) Allows(caseID, userID, action string) bool {
	p, ok := m.Get(caseID, userID)
	if !ok {
		return false
	}
	switch action {
	case "download":
		return p.CanDownload
	case "share":
		return p.CanShare
	case "edit":
		return p.CanEdit
	default:
		return domain.RoleAllows(p.Role, action)
	}
}
func Merge(a, b *Matrix) *Matrix {
	m := NewMatrix()
	for k, v := range a.entries {
		m.entries[k] = v
	}
	for k, v := range b.entries {
		m.entries[k] = v
	}
	return m
}
func Restrict(p domain.Permission) domain.Permission {
	p.CanDownload = p.CanDownload && domain.RoleAllows(p.Role, "download")
	p.CanShare = p.CanShare && domain.RoleAllows(p.Role, "share")
	p.CanEdit = p.CanEdit && domain.RoleAllows(p.Role, "edit")
	return p
}
