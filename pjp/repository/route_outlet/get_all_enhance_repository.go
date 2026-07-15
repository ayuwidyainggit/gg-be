package routeoutlet

import (
	"context"
	"fmt"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (repo *routeOutletRepository) GetAllEnhance(ctx context.Context, tx *gorm.DB, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]model.Pjp, int) {
	var pjp []model.Pjp
	var totalCount int64

	query := tx.Model(&pjp).Preload("RouteOutlets").Where("approval_status != ?", "Draft")

	for field, value := range filters {
		switch field {
		case "salesman_name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("salesman_name = ?", v)
			}
		case "salesman_code":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}
				query = query.Where("salesman_code IN ?", codes)
			}
		case "status":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}
				query = query.Where("approval_status IN ?", codes)
			}
		case "start_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				query = query.Where("created_at >= ?", v)
			}
		case "end_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				v = v.AddDate(0, 0, 1)
				query = query.Where("created_at <= ?", v)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					codes := strings.Split(v, ",")
					for i := range codes {
						codes[i] = strings.TrimSpace(codes[i])
					}
					query = query.Where(fmt.Sprintf("%s IN ?", field), codes)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("%s = ?", field), v)
				}
			}
		}
	}

	// Tambahkan filter cust_id
	query = query.Where("cust_id = ?", currentCustomerId)

	// Hitung total sebelum apply limit-offset
	query.Count(&totalCount)

	// Ambil data dengan pagination
	result := query.Scopes(response.Scopes(page, limit)).WithContext(ctx).Find(&pjp)
	helper.ErrorPanic(result.Error)

	return pjp, int(totalCount)
}
