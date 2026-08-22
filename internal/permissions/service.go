package permissions

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
)

type Service struct{ repo *store.Repository }

func New(repo *store.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Grant(caseID, userID, role string, download, share, edit bool) (domain.Permission, error) {
	role = domain.NormalizeRole(role)
	if caseID == "" || userID == "" {
		return domain.Permission{}, fmt.Errorf("case and user required")
	}
	p := domain.Permission{CaseID: caseID, UserID: userID, Role: role, CanDownload: download && domain.RoleAllows(role, "download"), CanShare: share && domain.RoleAllows(role, "share"), CanEdit: edit && domain.RoleAllows(role, "edit")}
	return p, s.repo.SavePermission(p)
}

func (s *Service) Check(caseID, userID, action string) (bool, error) {
	p, ok, err := s.repo.FindPermission(caseID, userID)
	if err != nil || !ok {
		return false, err
	}
	switch action {
	case "download":
		return p.CanDownload, nil
	case "share":
		return p.CanShare, nil
	case "edit":
		return p.CanEdit, nil
	case "read":
		return domain.RoleAllows(p.Role, "read"), nil
	default:
		return false, nil
	}
}

func (s *Service) Describe(caseID, userID string) (domain.Permission, error) {
	p, _, err := s.repo.FindPermission(caseID, userID)
	return p, err
}
