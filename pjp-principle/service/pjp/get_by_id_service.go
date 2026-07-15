package pjp

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"strconv"
)

func (service *pjpService) GetById(ctx context.Context, pjpId int, currentCustomerId string) response.PjpResponse {
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	data := service.pjpRepository.GetPjpById(ctx, tx, pjpId, currentCustomerId)

	statusBool, err := strconv.ParseBool(data.Status)
	if err != nil {
		statusBool = false
	}

	var res response.PjpResponse
	res.PjpCode = helper.FormatPjpCode(data.PjpCode)
	res.Status = statusBool
	helper.Automapper(data, &res)

	return res
}
