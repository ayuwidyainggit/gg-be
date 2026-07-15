package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"strconv"
)

func (service *pjpEnhanceService) UpdateStatusByEmpId(ctx context.Context, empId int, request request.UpdateStatusPjpEnhanceRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	findPjp := service.pjpRepository.GetPjpIdByEmpId(ctx, tx, empId, currentCustomerId)

	findPjp.Status = strconv.FormatBool(request.Status)
	service.pjpRepository.Update(ctx, tx, findPjp)
}
