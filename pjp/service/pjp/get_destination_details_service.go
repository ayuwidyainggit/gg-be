package pjp

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"strings"
)

func (service *pjpService) GetDestinationDetails(
	ctx context.Context,
	req request.DestinationDetailsRequest,
	currentCustomerId string,
) (response.DestinationDetailsData, response.LiveMonitoringPaging, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}

	sortOrder := strings.ToLower(req.SortOrder)
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	tx := service.db.Begin()
	if tx.Error != nil {
		return response.DestinationDetailsData{}, response.LiveMonitoringPaging{}, tx.Error
	}
	defer helper.CommitOrRollback(tx)

	isPrincipal, err := service.pjpRepository.IsPrincipalCustomer(ctx, tx, currentCustomerId)
	if err != nil {
		return response.DestinationDetailsData{}, response.LiveMonitoringPaging{}, err
	}

	rows, total, err := service.pjpRepository.GetDestinationDetails(ctx, tx, req.PjpID, req.Date, limit, page, sortOrder, currentCustomerId, isPrincipal)
	if err != nil {
		return response.DestinationDetailsData{}, response.LiveMonitoringPaging{}, err
	}

	return buildDestinationDetailsData(rows), buildDestinationDetailsPaging(total, page, limit), nil
}

func buildDestinationDetailsData(rows []response.DestinationDetailRow) response.DestinationDetailsData {
	data := response.DestinationDetailsData{
		Outlets:      []response.DestinationDetailOutlet{},
		Distributors: []response.DestinationDetailDistributor{},
	}
	if len(rows) == 0 {
		return data
	}

	first := rows[0]
	data.RouteCode = first.RouteCode
	data.RouteName = first.RouteName
	data.Week = first.Week
	data.Year = first.Year
	data.Date = first.Date.UTC()

	for _, row := range rows {
		switch strings.ToLower(row.DestinationType) {
		case "distributor":
			data.Distributors = append(data.Distributors, response.DestinationDetailDistributor{
				DistributorID:      row.DestinationID,
				DistributorCode:    row.DestinationCode,
				DistributorName:    row.DestinationName,
				DistributorStatus:  row.DestinationStatus,
				DistributorAddress: row.DestinationAddress,
				Longitude:          row.Longitude,
				Latitude:           row.Latitude,
			})
		default:
			data.Outlets = append(data.Outlets, response.DestinationDetailOutlet{
				OutletID:      row.DestinationID,
				OutletCode:    row.DestinationCode,
				OutletName:    row.DestinationName,
				Longitude:     row.Longitude,
				Latitude:      row.Latitude,
				OutletStatus:  row.DestinationStatus,
				OutletAddress: row.DestinationAddress,
			})
		}
	}

	return data
}

func buildDestinationDetailsPaging(total int64, page, limit int) response.LiveMonitoringPaging {
	totalPage := int(total) / limit
	if int(total)%limit > 0 {
		totalPage++
	}

	return response.LiveMonitoringPaging{
		TotalRecord: int(total),
		PageCurrent: page,
		PageLimit:   limit,
		PageTotal:   totalPage,
	}
}
