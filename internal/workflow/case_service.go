package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
	"lawdrive/internal/files"
	"lawdrive/internal/permissions"
	"lawdrive/internal/store"
	"time"
)

type Service struct {
	repo  *store.Repository
	files *files.Service
	perms *permissions.Service
	now   func() time.Time
}

func New(repo *store.Repository) *Service {
	return &Service{repo: repo, files: files.New(repo), perms: permissions.New(repo), now: func() time.Time { return time.Unix(1700000000, 0).UTC() }}
}

func (s *Service) CreateCase(id, title, client string) (domain.Case, error) {
	c := domain.Case{ID: id, Title: title, Client: client, Status: "open", CreatedAt: s.now()}
	if err := domain.ValidateCase(c); err != nil {
		return c, err
	}
	return c, s.repo.SaveCase(c)
}
func (s *Service) GrantAccess(caseID, userID, role string, download, share, edit bool) (domain.Permission, error) {
	return s.perms.Grant(caseID, userID, role, download, share, edit)
}
func (s *Service) UploadAttachment(caseID, fileID, actor, name, kind string, content []byte) (domain.CaseFile, error) {
	return s.files.Upload(caseID, fileID, actor, name, kind, content)
}
func (s *Service) EditAttachment(fileID, actor string, content []byte) (domain.CaseFile, error) {
	return s.files.Edit(fileID, actor, content)
}
func (s *Service) Download(caseID, fileID, userID string) ([]byte, error) {
	ok, err := s.perms.Check(caseID, userID, "download")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("download denied")
	}
	f, found, err := s.repo.FindFile(fileID)
	if err != nil || !found {
		return nil, fmt.Errorf("file missing")
	}
	return append([]byte(nil), f.Content...), nil
}
func (s *Service) Share(caseID, fileID, userID string) (domain.ShareLink, error) {
	ok, err := s.perms.Check(caseID, userID, "share")
	if err != nil || !ok {
		if err != nil {
			return domain.ShareLink{}, err
		}
		return domain.ShareLink{}, fmt.Errorf("share denied")
	}
	link := domain.ShareLink{ID: caseID + "/" + fileID + "/" + userID, FileID: fileID, CreatedBy: userID, ExpiresAt: s.now().Add(24 * time.Hour)}
	return link, s.repo.SaveShare(link)
}

func (s *Service) ArchiveCase(req domain.ArchiveRequest) (records []domain.CaseFile, err error) {
	c, ok, err := s.repo.FindCase(req.CaseID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("case missing")
	}
	if c.Archived {
		return nil, fmt.Errorf("case already archived")
	}
	all, err := s.repo.AllFiles()
	if err != nil {
		return nil, err
	}
	records = make([]domain.CaseFile, 0)
	for i := range all {
		f := all[i]
		if f.CaseID != req.CaseID {
			continue
		}
		if f.Archived {
			continue
		}
		defer func(file domain.CaseFile) { records = append(records, file) }(f)
		f.Archived = true
		if err := s.repo.SaveFile(f); err != nil {
			return nil, err
		}
	}
	c.Archived = true
	c.Status = "closed"
	if err := s.repo.SaveCase(c); err != nil {
		return nil, err
	}
	return records, nil
}

func (s *Service) ReopenCase(caseID string) error {
	c, ok, err := s.repo.FindCase(caseID)
	if err != nil || !ok {
		if err != nil {
			return err
		}
		return fmt.Errorf("case missing")
	}
	if !c.Archived {
		return nil
	}
	c.Archived = false
	c.Status = "open"
	return s.repo.SaveCase(c)
}

func (s *Service) FileService() *files.Service             { return s.files }
func (s *Service) PermissionService() *permissions.Service { return s.perms }
