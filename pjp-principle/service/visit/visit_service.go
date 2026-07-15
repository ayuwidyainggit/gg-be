package visit

import (
	"context"
	"scyllax-pjp/data/request"
	outletvisit "scyllax-pjp/repository/outlet_visit"
	"scyllax-pjp/repository/pjp"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type VisitService interface {
	StartVisit(ctx context.Context, request request.StartVisitRequest, custId string)
}
type visitService struct {
	validate        *validator.Validate
	db              *gorm.DB
	pjpRepo         pjp.PjpRepository
	outletVisitRepo outletvisit.OutletVisitRepository
}

func NewVisitService(pjpRepo pjp.PjpRepository, outletVisitRepo outletvisit.OutletVisitRepository, validate *validator.Validate, db *gorm.DB) VisitService {
	return &visitService{
		pjpRepo:         pjpRepo,
		outletVisitRepo: outletVisitRepo,
		validate:        validate,
		db:              db,
	}
}
