package pjp_principle

import (
	"context"
	"fmt"
	"scyllax-pjp/model"
	"strconv"

	"gorm.io/gorm"
)

type PjpPrincipleRepository interface {
	GetPjpsByEmpCodes(ctx context.Context, tx *gorm.DB, salesmanCode []string, currentCustomerId string) model.PjpPrinciple
	GetPjpsByEmpId(ctx context.Context, tx *gorm.DB, empID int, currentCustomerId string) model.PjpPrinciple
}

type pjpPrincipleRepository struct{}

func NewPjpPrincipleRepository() PjpPrincipleRepository {
	return &pjpPrincipleRepository{}
}

func applyFilters(db *gorm.DB, query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for field, value := range filters {
		switch field {
		case "status":
			if v, ok := value.(string); ok && v != "" {
				switch v {
				case "1":
					query = query.Where("pjp_principles.status = ?", "true")
				case "2":
					query = query.Where("pjp_principles.status = ?", "false")
				default:
					query = query.Where("pjp_principles.status = ?", v)
				}
			}
		case "q":
			if v, ok := value.(string); ok && v != "" {
				if intValue, err := strconv.Atoi(v); err == nil {
					query = query.Where("pjp_principles.pjp_code = ?", intValue)
				} else {
					query = query.Where(
						db.Where("pjp_principles.salesman_name ILIKE ?", "%"+v+"%").
							Or("pjp_principles.operation_type ILIKE ?", "%"+v+"%").
							Or("pjp_principles.team_salesman ILIKE ?", "%"+v+"%"),
					)
				}
			}
		case "pjp_code":
			if codes, ok := value.([]string); ok && len(codes) > 0 {
				query = query.Where("pjp_principles.pjp_code IN ?", codes)
			}
		case "team_salesman":
			if salesmen, ok := value.([]string); ok && len(salesmen) > 0 {
				query = query.Where("pjp_principles.team_salesman IN ?", salesmen)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					query = query.Where(fmt.Sprintf("pjp_principles.%s = ?", field), v)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("pjp_principles.%s = ?", field), v)
				}
			}
		}
	}
	return query
}
