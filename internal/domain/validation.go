package domain

import (
	"fmt"
	"strings"
)

func ValidatePermission(p Permission) error {
	if p.CaseID == "" || p.UserID == "" {
		return fmt.Errorf("permission identity required")
	}
	if NormalizeRole(p.Role) != p.Role {
		return fmt.Errorf("role must be normalized")
	}
	if p.CanEdit && !RoleAllows(p.Role, "edit") {
		return fmt.Errorf("role cannot edit")
	}
	if p.CanShare && !RoleAllows(p.Role, "share") {
		return fmt.Errorf("role cannot share")
	}
	return nil
}
func ValidateVersion(v FileVersion) error {
	if v.FileID == "" || v.ID == "" {
		return fmt.Errorf("version identity required")
	}
	if v.Number < 1 {
		return fmt.Errorf("version must be positive")
	}
	if strings.TrimSpace(v.CreatedBy) == "" {
		return fmt.Errorf("version creator required")
	}
	return nil
}
func ValidateShare(s ShareLink) error {
	if s.ID == "" || s.FileID == "" {
		return fmt.Errorf("share identity required")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("share creator required")
	}
	if s.ExpiresAt.IsZero() {
		return fmt.Errorf("share expiry required")
	}
	return nil
}
func ValidateAudit(a AuditEntry) error {
	if a.ID == "" || a.Actor == "" {
		return fmt.Errorf("audit identity required")
	}
	if a.Action == "" {
		return fmt.Errorf("audit action required")
	}
	return nil
}
func NormalizeCaseTitle(title string) string { return strings.Join(strings.Fields(title), " ") }
func NormalizeClient(name string) string     { return strings.TrimSpace(name) }
func IsTextKind(kind string) bool            { return kind == "minutes" || kind == "contract" }
func IsBinaryKind(kind string) bool          { return kind == "evidence" }
func IsKnownKind(kind string) bool           { return IsTextKind(kind) || IsBinaryKind(kind) }
func StatusLabel(status string) string {
	switch status {
	case "open":
		return "进行中"
	case "closed":
		return "已结案"
	default:
		return "未知"
	}
}
func RoleLabel(role string) string {
	switch NormalizeRole(role) {
	case "partner":
		return "合伙人"
	case "lawyer":
		return "律师"
	case "assistant":
		return "助理"
	default:
		return "审计员"
	}
}
func ActionLabel(action string) string {
	switch action {
	case "upload":
		return "上传"
	case "edit":
		return "编辑"
	case "download":
		return "下载"
	case "share":
		return "分享"
	case "archive":
		return "归档"
	default:
		return action
	}
}
