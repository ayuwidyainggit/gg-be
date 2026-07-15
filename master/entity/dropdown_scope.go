package entity

import "strings"

// EmployeeDropdownScope keeps normalized scope for dropdown endpoints.
type EmployeeDropdownScope struct {
	RegionScope      string
	AreaScope        string
	DistributorScope string
}

func NormalizeDropdownScope(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "specific", "spesific", "selected":
		return "specific"
	default:
		return "all"
	}
}
