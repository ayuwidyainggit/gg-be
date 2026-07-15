package pjp

import (
	"context"
	"fmt"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"
	"strconv"

	"gorm.io/gorm"
)

type PjpRepository interface {
	Create(ctx context.Context, tx *gorm.DB, pjp model.Pjp) model.Pjp
	GetAll(ctx context.Context, tx *gorm.DB, limit int, page int, filters map[string]interface{}, currentCustomerId string) ([]model.Pjp, int64)
	GetPjpWithRoute(ctx context.Context, tx *gorm.DB, q string, custId string) []response.PjpWithRouteRow
	Update(ctx context.Context, tx *gorm.DB, pjp model.Pjp)
	DeleteByPjpId(ctx context.Context, tx *gorm.DB, pjpId int, custID string)
	GetPjpById(ctx context.Context, tx *gorm.DB, pjpId int, currentCustomerId string) model.Pjp
	GetPjpIdByEmpId(ctx context.Context, tx *gorm.DB, salesmanId int, currentCustomerId string) model.Pjp
	GetPjpIdByEmpCode(ctx context.Context, tx *gorm.DB, salesmanCode string, currentCustomerId string) model.Pjp
	ListPjpApprove(ctx context.Context, tx *gorm.DB, q string, custId string) []response.PjpWithRouteRow
	Patch(ctx context.Context, tx *gorm.DB, pjpId int, pjpMode, custId string)
}

type pjpRepository struct{}

func NewPjpRepository() PjpRepository {
	return &pjpRepository{}
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
