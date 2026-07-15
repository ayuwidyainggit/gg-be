package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"strconv"
)

func (service *pjpEnhanceService) UpdateStatusPjp(ctx context.Context, id int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	findPjp := service.pjpRepository.GetPjpById(ctx, tx, id, currentCustomerId)

	findPjp.Status = strconv.FormatBool(*request.Status)
	service.pjpRepository.Update(ctx, tx, findPjp)
}
