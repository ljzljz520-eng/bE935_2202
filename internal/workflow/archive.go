package workflow

import (
	"fmt"
	"lawdrive/internal/domain"
)

func ValidateArchiveRequest(req domain.ArchiveRequest) error {
	if req.CaseID == "" {
		return fmt.Errorf("case id required")
	}
	if req.Actor == "" {
		return fmt.Errorf("actor required")
	}
	return nil
}
func ArchiveSummary(rows []domain.CaseFile) (int, int) {
	total := len(rows)
	bytes := 0
	for _, f := range rows {
		bytes += len(f.Content)
	}
	return total, bytes
}
func IsArchiveComplete(c domain.Case, rows []domain.CaseFile) bool {
	if !c.Archived {
		return false
	}
	for _, f := range rows {
		if !f.Archived {
			return false
		}
	}
	return true
}
