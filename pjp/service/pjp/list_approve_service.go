package pjp

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"strconv"
)

func (service *pjpService) ListPjpApprove(ctx context.Context, q string, custId string) []response.PjpResponse {
	pjpMap := make(map[int]*response.PjpResponse)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	rows := service.pjpRepository.ListPjpApprove(ctx, tx, q, custId)
	for _, row := range rows {
		pjpId := row.ID

		if _, exists := pjpMap[pjpId]; !exists {
			statusBool, _ := strconv.ParseBool(row.Status)

			pjpMap[pjpId] = &response.PjpResponse{
				ID:             row.ID,
				PjpCode:        helper.FormatPjpCode(row.PjpCode),
				OperationType:  row.OperationType,
				TeamSalesMan:   row.TeamSalesMan,
				SalesManID:     row.SalesManID,
				SalesmanName:   row.SalesmanName,
				SalesmanCode:   row.SalesmanCode,
				WarehouseID:    row.WarehouseID,
				WarehouseName:  row.WarehouseName,
				PjpMode:        row.PjpMode,
				Status:         statusBool,
				ApprovalStatus: row.ApprovalStatus,
				Route:          []response.RoutesEntity{},
			}
		}

		if row.RouteCode != 0 {
			pjpMap[pjpId].Route = append(pjpMap[pjpId].Route, response.RoutesEntity{
				RouteCode:   row.RouteCode,
				TotalOutlet: row.TotalOutlet,
			})
		}
	}

	// Ubah map jadi slice
	var payload []response.PjpResponse
	for _, pjp := range pjpMap {
		payload = append(payload, *pjp)
	}

	return payload
}
