package store

import (
	"fmt"
	"lawdrive/internal/domain"
)

type Batch struct {
	repo   *Repository
	cases  []domain.Case
	files  []domain.CaseFile
	audits []domain.AuditEntry
}

func NewBatch(repo *Repository) *Batch               { return &Batch{repo: repo} }
func (b *Batch) AddCase(c domain.Case) *Batch        { b.cases = append(b.cases, c); return b }
func (b *Batch) AddFile(f domain.CaseFile) *Batch    { b.files = append(b.files, f); return b }
func (b *Batch) AddAudit(a domain.AuditEntry) *Batch { b.audits = append(b.audits, a); return b }
func (b *Batch) Commit() error {
	for _, c := range b.cases {
		if e := b.repo.SaveCase(c); e != nil {
			return e
		}
	}
	for _, f := range b.files {
		if e := b.repo.SaveFile(f); e != nil {
			return e
		}
	}
	for _, a := range b.audits {
		if e := b.repo.SaveAudit(a); e != nil {
			return e
		}
	}
	return nil
}
func (b *Batch) Validate() error {
	if len(b.cases) == 0 && len(b.files) == 0 && len(b.audits) == 0 {
		return fmt.Errorf("empty batch")
	}
	for _, c := range b.cases {
		if e := domain.ValidateCase(c); e != nil {
			return e
		}
	}
	for _, f := range b.files {
		if e := domain.ValidateFile(f); e != nil {
			return e
		}
	}
	return nil
}
func (r *Repository) SaveCaseWithFiles(c domain.Case, files []domain.CaseFile) error {
	b := NewBatch(r).AddCase(c)
	for _, f := range files {
		b.AddFile(f)
	}
	if e := b.Validate(); e != nil {
		return e
	}
	return b.Commit()
}
func (r *Repository) FilesForCase(caseID string) ([]domain.CaseFile, error) {
	all, e := r.AllFiles()
	if e != nil {
		return nil, e
	}
	out := make([]domain.CaseFile, 0)
	for _, f := range all {
		if f.CaseID == caseID {
			out = append(out, f)
		}
	}
	return domain.SortFiles(out), nil
}
