package pjp

import (
	"context"
	"scyllax-pjp/helper"
)

func (service *pjpService) Delete(ctx context.Context, pjpId int, custId string) {
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	pjp := service.pjpRepository.GetPjpById(ctx, tx, pjpId, custId)
	routeOutlets := service.routeOutletRepository.FindAllOutletsByPjpIds(ctx, tx, []int{pjp.ID}, custId)

	// Ambil route_code unik
	routeCodeMap := make(map[int]struct{})
	for _, ro := range routeOutlets {
		routeCodeMap[ro.RouteCode] = struct{}{}
	}

	var routeCodes []int
	for code := range routeCodeMap {
		routeCodes = append(routeCodes, code)
	}

	service.routeRepository.DeleteByRouteCodes(ctx, tx, routeCodes, custId)
	service.pjpRepository.DeleteByPjpId(ctx, tx, pjpId, custId)
}
