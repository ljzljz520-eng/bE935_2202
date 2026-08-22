package share

import (
	"fmt"
	"lawdrive/internal/domain"
	"time"
)

func DefaultTTL(role string) time.Duration {
	switch domain.NormalizeRole(role) {
	case "partner":
		return 7 * 24 * time.Hour
	case "lawyer":
		return 72 * time.Hour
	default:
		return 24 * time.Hour
	}
}
func CanCreate(role string) bool { return domain.RoleAllows(role, "share") }
func ValidateRequest(fileID, user, role string, ttl time.Duration) error {
	if fileID == "" || user == "" {
		return fmt.Errorf("file and user required")
	}
	if !CanCreate(role) {
		return fmt.Errorf("share permission required")
	}
	if ttl <= 0 {
		return fmt.Errorf("ttl must be positive")
	}
	return nil
}
func (s *Service) CreateForRole(fileID, user, role string) (domain.ShareLink, error) {
	ttl := DefaultTTL(role)
	if e := ValidateRequest(fileID, user, role, ttl); e != nil {
		return domain.ShareLink{}, e
	}
	return s.Create(fileID, user, ttl)
}
func IsUsable(link domain.ShareLink, now time.Time) bool {
	return !link.Revoked && now.Before(link.ExpiresAt)
}
