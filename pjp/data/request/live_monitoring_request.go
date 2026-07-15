package request

import (
	"strconv"
	"strings"
)

// LiveMonitoringRequest is the query parameter request for location monitoring (Principal & Distributor)
type LiveMonitoringRequest struct {
	RegionID      int      `form:"region_id"`
	AreaID        int      `form:"area_id"`
	DistributorID int      `form:"distributor_id"`
	Date          int64    `form:"date" binding:"required"` // epoch timestamp
	EmpIDs        []int    `form:"emp_id[]"`
	LegacyEmpID   string   `form:"emp_id"`
	Status        []string `form:"status[]" binding:"required"`
	Page          int      `form:"page"`
	Limit         int      `form:"limit"`
}

// GetEmpIDs returns employee IDs from array binding and keeps backward compatibility with comma-separated emp_id.
func (r *LiveMonitoringRequest) GetEmpIDs() []int {
	if len(r.EmpIDs) > 0 {
		return r.EmpIDs
	}

	if r.LegacyEmpID == "" {
		return nil
	}

	parts := strings.Split(r.LegacyEmpID, ",")
	var ids []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if id, err := strconv.Atoi(p); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

// LiveMonitoringDetailRequest is the query parameter request for location monitoring detail
type LiveMonitoringDetailRequest struct {
	EmpID         int    `form:"emp_id" binding:"required"`
	DistributorID *int   `form:"distributor_id"`
	Date          string `form:"date" binding:"required"`
}

// UpdateLocationsRequest contains query parameters for an employee location timeline.
type UpdateLocationsRequest struct {
	EmpID int    `form:"emp_id" binding:"required"`
	Date  string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}
