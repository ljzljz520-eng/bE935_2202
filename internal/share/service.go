package share

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"time"
)

type Service struct {
	repo *store.Repository
	now  func() time.Time
}

func New(repo *store.Repository) *Service {
	return &Service{repo: repo, now: func() time.Time { return time.Unix(1700000000, 0).UTC() }}
}
func (s *Service) Create(fileID, user string, ttl time.Duration) (domain.ShareLink, error) {
	if fileID == "" || user == "" {
		return domain.ShareLink{}, fmt.Errorf("file and user required")
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	link := domain.ShareLink{ID: fileID + "/" + user, FileID: fileID, CreatedBy: user, ExpiresAt: s.now().Add(ttl)}
	return link, s.repo.SaveShare(link)
}
func (s *Service) Validate(link domain.ShareLink, at time.Time) error {
	if link.Revoked {
		return fmt.Errorf("link revoked")
	}
	if !at.Before(link.ExpiresAt) {
		return fmt.Errorf("link expired")
	}
	return nil
}
func (s *Service) Revoke(link domain.ShareLink) error {
	link.Revoked = true
	return s.repo.SaveShare(link)
}
