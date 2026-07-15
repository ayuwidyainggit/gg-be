package pjp

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
)

func (service *pjpService) Create(ctx context.Context, req request.PjpRequest, currentCustomerId string) {
	if err := service.validate.Struct(req); err != nil {
		helper.ErrorPanic(err)
	}

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	pjp := mapRequestToModel(req, currentCustomerId)

	service.pjpRepository.Create(ctx, tx, pjp)
}
