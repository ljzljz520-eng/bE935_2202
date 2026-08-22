package store

import (
	"encoding/json"
	"fmt"
	"lawdrive/internal/domain"
)

type Snapshot struct {
	Cases  []domain.Case       `json:"cases"`
	Files  []domain.CaseFile   `json:"files"`
	Audits []domain.AuditEntry `json:"audits"`
}

func (r *Repository) Snapshot() (Snapshot, error) {
	cs, e := r.AllCases()
	if e != nil {
		return Snapshot{}, e
	}
	fs, e := r.AllFiles()
	if e != nil {
		return Snapshot{}, e
	}
	as, e := r.AllAudits()
	if e != nil {
		return Snapshot{}, e
	}
	return Snapshot{Cases: cs, Files: fs, Audits: as}, nil
}
func (r *Repository) ExportJSON() ([]byte, error) {
	snap, e := r.Snapshot()
	if e != nil {
		return nil, e
	}
	return json.MarshalIndent(snap, "", "  ")
}
func (r *Repository) ImportJSON(data []byte) error {
	var snap Snapshot
	if e := json.Unmarshal(data, &snap); e != nil {
		return e
	}
	if len(snap.Cases) == 0 && len(snap.Files) == 0 && len(snap.Audits) == 0 {
		return fmt.Errorf("empty snapshot")
	}
	for _, c := range snap.Cases {
		if e := r.SaveCase(c); e != nil {
			return e
		}
	}
	for _, f := range snap.Files {
		if e := r.SaveFile(f); e != nil {
			return e
		}
	}
	for _, a := range snap.Audits {
		if e := r.SaveAudit(a); e != nil {
			return e
		}
	}
	return nil
}
func (r *Repository) EnsureCase(id string) error {
	_, ok, e := r.FindCase(id)
	if e != nil {
		return e
	}
	if !ok {
		return fmt.Errorf("case %s missing", id)
	}
	return nil
}
func (r *Repository) EnsureFile(id string) error {
	_, ok, e := r.FindFile(id)
	if e != nil {
		return e
	}
	if !ok {
		return fmt.Errorf("file %s missing", id)
	}
	return nil
}
