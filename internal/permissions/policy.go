package permissions

import "lawdrive/internal/domain"

func DefaultForRole(role string) domain.Permission {
	role = domain.NormalizeRole(role)
	return domain.Permission{Role: role, CanDownload: domain.RoleAllows(role, "download"), CanShare: domain.RoleAllows(role, "share"), CanEdit: domain.RoleAllows(role, "edit")}
}

func AllowedKinds(role string) []string {
	if domain.RoleAllows(role, "edit") {
		return []string{"contract", "evidence", "minutes"}
	}
	if domain.RoleAllows(role, "download") {
		return []string{"contract", "evidence"}
	}
	return []string{"contract", "evidence", "minutes"}
}

func CanManage(role string) bool { return domain.RoleAllows(role, "archive") }
