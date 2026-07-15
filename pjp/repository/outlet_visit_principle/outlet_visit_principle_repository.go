package outlet_visit_principle

import (
	"context"
	"scyllax-pjp/data/entity"
	"scyllax-pjp/data/response"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

type OutletVisitPrincipleRepository interface {
	UpdateByPjpIDandDate(ctx context.Context, tx *gorm.DB, pjpID int, date string, data model.OutletVisitListPrinciple)
	UpdateOutletVisitListColumnAt(ctx context.Context, tx *gorm.DB, column string, currentTime *int64, date string, id int64)
	UpdateOutletVisitListWithFile(ctx context.Context, tx *gorm.DB, id int64, date string, fileInfo model.OutletVisitListPrinciple) error
	FindById(ctx context.Context, tx *gorm.DB, outletVisitId int) (data response.TodoListResponse, err error)
	GetVisitsByDateAndSalesman(ctx context.Context, tx *gorm.DB, dataFilter entity.SalesmanReportQueryFilter) ([]model.OutletVisitListPrinciple, error)
	CheckRouteOutletAdditional(ctx context.Context, tx *gorm.DB, destinationCode string, destinationId int) (bool, error)
	UpdateOutletVisitListSkipColumnAt(ctx context.Context, tx *gorm.DB, column string, currentTime int64, date string, id int64, skipReson string, inOutlet bool, fileInfo model.OutletVisitListPrinciple)
	FindDataByID(ctx context.Context, tx *gorm.DB, outletVisitId int) (data model.OutletVisitListPrinciple, err error)
	UpdateOutletVisitListByID(ctx context.Context, tx *gorm.DB, id int64, data model.OutletVisitListPrinciple) error
}
type outletVisitPrincipleRepository struct {
}

func NewOutletVisitPrincipleRepository() OutletVisitPrincipleRepository {
	return &outletVisitPrincipleRepository{}
}
