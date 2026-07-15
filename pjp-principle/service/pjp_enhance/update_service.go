package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
)

func (service *pjpEnhanceService) UpdatePjp(ctx context.Context, id int, request request.CreatePjpEnhanceRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	findPjp := service.pjpRepository.GetPjpById(ctx, tx, id, currentCustomerId)

	pjp := buildPjpModel(request, currentCustomerId)
	pjp.ID = findPjp.ID
	service.pjpRepository.Update(ctx, tx, pjp)

	service.routeRepository.DeleteByPjpId(ctx, tx, pjp.ID, currentCustomerId)
	service.routeOutletHistoryRepository.DeleteByPjpId(ctx, tx, pjp.ID, currentCustomerId)

	savedRoutes := service.createRoutes(ctx, tx, pjp, request.Routes, currentCustomerId)
	routePopPermanents := service.createVisitHistory(ctx, tx, pjp, savedRoutes, request.VisitDay, currentCustomerId)

	if routePopPermanents != nil {
		service.routePopRepository.CreateBulk(ctx, tx, routePopPermanents)
	}
}
