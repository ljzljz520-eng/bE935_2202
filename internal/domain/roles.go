package domain

import "fmt"

var roleOrder = map[string]int{"partner": 4, "lawyer": 3, "assistant": 2, "auditor": 1}

func NormalizeRole(role string) string {
	if _, ok := roleOrder[role]; ok {
		return role
	}
	return "auditor"
}

func RoleAllows(role, action string) bool {
	level, ok := roleOrder[NormalizeRole(role)]
	if !ok {
		return false
	}
	switch action {
	case "read":
		return level >= 1
	case "download":
		return level >= 2
	case "share":
		return level >= 3
	case "edit":
		return level >= 3
	case "archive":
		return level >= 4
	default:
		return false
	}
}

func ValidateCase(c Case) error {
	if c.ID == "" {
		return fmt.Errorf("case id required")
	}
	if c.Title == "" {
		return fmt.Errorf("case title required")
	}
	if c.Status != "open" && c.Status != "closed" {
		return fmt.Errorf("invalid status")
	}
	return nil
}

func ValidateFile(f CaseFile) error {
	if f.ID == "" || f.CaseID == "" {
		return fmt.Errorf("file identity required")
	}
	if f.Name == "" {
		return fmt.Errorf("file name required")
	}
	if f.Kind != "contract" && f.Kind != "evidence" && f.Kind != "minutes" {
		return fmt.Errorf("unsupported file kind")
	}
	return nil
}
