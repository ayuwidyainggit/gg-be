package visit

import (
	"context"
	"errors"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
)

func (service *visitService) StartVisit(ctx context.Context, request request.StartVisitRequest, custId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	if request.SalesmanCode == "" {
		helper.ErrorPanic(errors.New("salesman code is empty"))
	}

	pjpId := service.pjpRepo.GetPjpIdByEmpCode(ctx, tx, request.SalesmanCode, custId)

	data := model.OutletVisitList{
		Start: &request.CurrentTime,
	}

	service.outletVisitRepo.UpdateByPjpIDandDate(ctx, tx, pjpId.ID, request.Date, data)
}
