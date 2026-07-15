package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"strings"
	"time"
)

const (
	sapExportOrderType         = "ZOR"
	sapExportSalesOrganization = 5000
	sapExportDateTimeLayout    = "2006-01-02 15:04:05"
	sapExportDateLayout        = "2006-01-02"
)

func (service *replenishmentServiceImpl) SAPGetReplenishmentExport(query entity.SAPReplExportQuery) ([]entity.SAPReplExportItem, error) {
	ctx := context.Background()
	rows, err := service.ReplenishmentRepository.FindSAPExportRows(
		ctx,
		strings.TrimSpace(query.DistributorCode),
		query.DateFrom,
		query.DateTo,
	)
	if err != nil {
		return nil, err
	}
	return mapSAPReplExportRows(rows), nil
}

func mapSAPReplExportRows(rows []model.SAPReplExportRow) []entity.SAPReplExportItem {
	if len(rows) == 0 {
		return []entity.SAPReplExportItem{}
	}

	loc := constant.AsiaJakartaLocation
	ordered := make([]entity.SAPReplExportItem, 0)
	indexByNo := make(map[string]int)

	for _, row := range rows {
		key := row.ReplenishmentNo
		idx, ok := indexByNo[key]
		if !ok {
			item := entity.SAPReplExportItem{
				CustReference:       row.ReplenishmentNo,
				ShipToPartyCode:     strPtr(row.ShipToPartyCode),
				ReplenishmentType:   row.ReplenishmentType,
				IsAdditionFrom:      row.IsAdditionFrom,
				Status:              row.Status,
				OrderType:           sapExportOrderType,
				SalesOrganization:   sapExportSalesOrganization,
				DistributionChannel: strPtr(row.DistributionChannel),
				SalesOffice:         strPtr(row.SalesOffice),
				DeliveryDate:        formatSAPExportDate(row.DeliveryDate, loc),
				ShippingPoint:       strPtr(row.ShippingPoint),
				Plant:               strPtr(row.Plant),
				CreatedAt:           row.CreatedAt.In(loc).Format(sapExportDateTimeLayout),
				Details:             make([]entity.SAPReplExportDetail, 0),
			}
			ordered = append(ordered, item)
			indexByNo[key] = len(ordered) - 1
			idx = indexByNo[key]
		}

		ordered[idx].Details = append(ordered[idx].Details, entity.SAPReplExportDetail{
			Material:    strPtr(row.Material),
			Division:    strPtr(row.Division),
			OrderQty:    row.OrderQty,
			Uom:         strPtr(row.Uom),
			CustPoPrice: row.CustPoPrice,
		})
	}

	return ordered
}

func strPtr(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func formatSAPExportDate(t *time.Time, loc *time.Location) string {
	if t == nil {
		return ""
	}
	return t.In(loc).Format(sapExportDateLayout)
}
