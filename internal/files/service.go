package files

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/store"
	"strings"
	"time"
)

type Service struct {
	repo  *store.Repository
	clock func() time.Time
}

func New(repo *store.Repository) *Service {
	return &Service{repo: repo, clock: func() time.Time { return time.Unix(1700000000, 0).UTC() }}
}

func (s *Service) Upload(caseID, fileID, actor, name, kind string, content []byte) (domain.CaseFile, error) {
	f := domain.CaseFile{ID: fileID, CaseID: caseID, Name: name, Kind: kind, Content: append([]byte(nil), content...), Version: 1, UpdatedAt: s.clock()}
	if err := domain.ValidateFile(f); err != nil {
		return f, err
	}
	if _, ok, err := s.repo.FindCase(caseID); err != nil {
		return f, err
	} else if !ok {
		return f, fmt.Errorf("case %s missing", caseID)
	}
	if err := s.repo.SaveFile(f); err != nil {
		return f, err
	}
	v := domain.FileVersion{ID: fileID + ":1", FileID: fileID, Number: 1, Content: f.Content, CreatedBy: actor, CreatedAt: s.clock()}
	if err := s.repo.SaveVersion(v); err != nil {
		return f, err
	}
	return f, nil
}

func (s *Service) Edit(fileID, actor string, content []byte) (domain.CaseFile, error) {
	f, ok, err := s.repo.FindFile(fileID)
	if err != nil {
		return f, err
	}
	if !ok {
		return f, fmt.Errorf("file missing")
	}
	if f.Archived {
		return f, fmt.Errorf("file archived")
	}
	f.Version++
	f.Content = append([]byte(nil), content...)
	f.UpdatedAt = s.clock()
	if err := s.repo.SaveFile(f); err != nil {
		return f, err
	}
	v := domain.FileVersion{ID: fmt.Sprintf("%s:%d", fileID, f.Version), FileID: fileID, Number: f.Version, Content: f.Content, CreatedBy: actor, CreatedAt: s.clock()}
	return f, s.repo.SaveVersion(v)
}

func (s *Service) Archive(fileID string) error {
	f, ok, err := s.repo.FindFile(fileID)
	if err != nil || !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("file missing")
	}
	f.Archived = true
	return s.repo.SaveFile(f)
}

func (s *Service) Restore(fileID string) error {
	f, ok, err := s.repo.FindFile(fileID)
	if err != nil || !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("file missing")
	}
	f.Archived = false
	return s.repo.SaveFile(f)
}

func (s *Service) Preview(fileID string) (string, error) {
	f, ok, err := s.repo.FindFile(fileID)
	if err != nil || !ok {
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("file missing")
	}
	if len(f.Content) == 0 {
		return "", nil
	}
	return string(f.Content), nil
}

func (s *Service) Match(f domain.CaseFile, query string) bool {
	q := strings.ToLower(strings.TrimSpace(query))
	return q == "" || strings.Contains(strings.ToLower(f.Name), q) || strings.Contains(strings.ToLower(string(f.Content)), q)
}
