package outletvisit

import (
	"context"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type OutletVisitRepository interface {
	UpdateByPjpIDandDate(ctx context.Context, tx *gorm.DB, pjpID int, date string, data model.OutletVisitList)
}
type outletVisitRepository struct {
}

func NewOutletVisitRepository() OutletVisitRepository {
	return &outletVisitRepository{}
}
