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

	// empCodes := strings.Split(request.SalesmanCode, ",") remove because get from auth salesman id

	if len(custId) <= 6 {

		pjpPrinciple := service.pjpPrincipleRepo.GetPjpsByEmpId(ctx, tx, request.EmpID, custId)
		if pjpPrinciple.ID == 0 {
			helper.ErrorPanic(errors.New("pjp data not found"))
		}

		if pjpPrinciple.DistributorID != 0 {
			helper.ErrorPanic(errors.New("pjp principle not found"))
		}
		data := model.OutletVisitListPrinciple{
			Start: &request.CurrentTime,
		}

		service.outletVisitPrincipleRepo.UpdateByPjpIDandDate(ctx, tx, pjpPrinciple.ID, request.Date, data)
	} else {
		pjpId := service.pjpRepo.GetPjpByEmpId(ctx, tx, request.EmpID, custId)
		if pjpId.ID == 0 {
			helper.ErrorPanic(errors.New("pjp data not found"))
		}
		data := model.OutletVisitList{
			Start: &request.CurrentTime,
		}

		service.outletVisitRepo.UpdateByPjpIDandDate(ctx, tx, pjpId.ID, request.Date, data)
	}
}
