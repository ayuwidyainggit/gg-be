package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/data/response"

	"gorm.io/gorm"
)

func (repo *outletVisitPrincipleRepository) FindById(ctx context.Context, tx *gorm.DB, outletVisitId int) (data response.TodoListResponse, err error) {
	result := tx.WithContext(ctx).
		Table("pjp_principles.outlet_visit_list").
		Select("arrive_at, leave_at", "on_hold", "resume_at", "skip_at", "skip_reason", "skip_in_outlet as in_outlet").
		Where("id = ?", outletVisitId).
		Scan(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}
