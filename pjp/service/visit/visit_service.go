package visit

import (
	"context"
	"scyllax-pjp/data/request"
	outletvisit "scyllax-pjp/repository/outlet_visit"
	"scyllax-pjp/repository/outlet_visit_principle"
	"scyllax-pjp/repository/pjp"
	"scyllax-pjp/repository/pjp_principle"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type VisitService interface {
	StartVisit(ctx context.Context, request request.StartVisitRequest, custId string)
}
type visitService struct {
	validate                 *validator.Validate
	db                       *gorm.DB
	pjpRepo                  pjp.PjpRepository
	pjpPrincipleRepo         pjp_principle.PjpPrincipleRepository
	outletVisitRepo          outletvisit.OutletVisitRepository
	outletVisitPrincipleRepo outlet_visit_principle.OutletVisitPrincipleRepository
}

func NewVisitService(pjpRepo pjp.PjpRepository, pjpPrincipleRepo pjp_principle.PjpPrincipleRepository, outletVisitRepo outletvisit.OutletVisitRepository, outletVisitPrincipleRepo outlet_visit_principle.OutletVisitPrincipleRepository, validate *validator.Validate, db *gorm.DB) VisitService {
	return &visitService{
		pjpRepo:                  pjpRepo,
		outletVisitRepo:          outletVisitRepo,
		validate:                 validate,
		db:                       db,
		pjpPrincipleRepo:         pjpPrincipleRepo,
		outletVisitPrincipleRepo: outletVisitPrincipleRepo,
	}
}
