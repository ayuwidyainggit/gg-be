package pjp

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
)

func (service *pjpService) Update(ctx context.Context, request request.PjpRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	pjp := mapRequestToModel(request, currentCustomerId)

	service.pjpRepository.Update(ctx, tx, pjp)
}
