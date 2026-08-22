package store

import (
	"fmt"
	"lawdrive/internal/domain"
)

type Repository struct{ store *Store }

func NewRepository(s *Store) *Repository { return &Repository{store: s} }

func (r *Repository) SaveCase(c domain.Case) error { return r.store.Put("cases", c.ID, c) }
func (r *Repository) FindCase(id string) (domain.Case, bool, error) {
	var c domain.Case
	ok, err := r.store.Get("cases", id, &c)
	return c, ok, err
}
func (r *Repository) SaveFile(f domain.CaseFile) error { return r.store.Put("files", f.ID, f) }
func (r *Repository) FindFile(id string) (domain.CaseFile, bool, error) {
	var f domain.CaseFile
	ok, err := r.store.Get("files", id, &f)
	return f, ok, err
}
func (r *Repository) SavePermission(p domain.Permission) error {
	return r.store.Put("permissions", p.CaseID+"/"+p.UserID, p)
}
func (r *Repository) FindPermission(caseID, userID string) (domain.Permission, bool, error) {
	var p domain.Permission
	ok, err := r.store.Get("permissions", caseID+"/"+userID, &p)
	return p, ok, err
}
func (r *Repository) SaveVersion(v domain.FileVersion) error { return r.store.Put("versions", v.ID, v) }
func (r *Repository) SaveAudit(a domain.AuditEntry) error {
	if a.ID == "" {
		return fmt.Errorf("audit id required")
	}
	return r.store.Put("audits", a.ID, a)
}
func (r *Repository) SaveShare(s domain.ShareLink) error { return r.store.Put("shares", s.ID, s) }

func (r *Repository) AllFiles() ([]domain.CaseFile, error) {
	rows, err := r.store.List("files")
	if err != nil {
		return nil, err
	}
	result := make([]domain.CaseFile, 0, len(rows))
	for _, row := range rows {
		var f domain.CaseFile
		if e := decode(row, &f); e != nil {
			return nil, e
		}
		result = append(result, f)
	}
	return result, nil
}

func (r *Repository) AllCases() ([]domain.Case, error) {
	rows, err := r.store.List("cases")
	if err != nil {
		return nil, err
	}
	result := make([]domain.Case, 0, len(rows))
	for _, row := range rows {
		var c domain.Case
		if e := decode(row, &c); e != nil {
			return nil, e
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *Repository) AllAudits() ([]domain.AuditEntry, error) {
	rows, err := r.store.List("audits")
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEntry, 0, len(rows))
	for _, row := range rows {
		var a domain.AuditEntry
		if e := decode(row, &a); e != nil {
			return nil, e
		}
		result = append(result, a)
	}
	return result, nil
}
