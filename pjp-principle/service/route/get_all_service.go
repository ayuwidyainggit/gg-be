package route

import (
	"context"
	"math"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"time"
)

func (service *routeService) GetAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingResponse, response.Meta, error) {

	tx := service.db.Begin()
	if tx.Error != nil {
		return nil, response.Meta{}, tx.Error
	}
	defer helper.CommitOrRollback(tx)

	routeModels, totalData := service.destinationRepo.GetAll(ctx, tx, page, limit, filters, currentCustomerId)

	routes := make([]response.ApprovalRouteMappingResponse, 0, len(routeModels))
	for _, model := range routeModels {
		routes = append(routes, mapRouteToResponse(model))
	}

	pagination := response.Meta{
		TotalData: totalData,
		Page:      page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(totalData) / float64(limit))),
	}

	return routes, pagination, nil
}

func mapRouteToResponse(value model.Destination) response.ApprovalRouteMappingResponse {
	res := response.ApprovalRouteMappingResponse{
		ID:           value.ID,
		RouteCode:    value.OldRouteCode,
		RouteName:    value.OldRouteName,
		NewRouteName: value.RouteName,
		NewRouteCode: value.RouteCode,
		Status:       value.Status,
		Date:         value.CreatedAt.Format("2006-01-02 15:04"),
	}

	if value.VerifiedDate != nil {
		res.VerifiedDate = time.Now().Format("2006-01-02 15:04")
	}
	if value.PjpID != nil || value.PjpCode != nil {
		res.PjpID = value.PjpID
		res.PjpCode = value.PjpCode
	}
	if value.Pjp != nil {
		if value.Pjp.SalesmanCode != "" {
			res.SalesmanCode = &value.Pjp.SalesmanCode
		}
		if value.Pjp.SalesmanName != "" {
			res.SalesmanName = &value.Pjp.SalesmanName
		}
	}
	if value.PjpOld != nil {
		if value.PjpOld.SalesmanCode != "" {
			res.NewSalesmanCode = &value.PjpOld.SalesmanCode
		}
		if value.PjpOld.SalesmanName != "" {
			res.NewSalesmanName = &value.PjpOld.SalesmanName
		}
	}
	if value.DestinationCode != "" {
		res.Outlets = &response.DestinationResponse{
			DestinationID:      value.DestinationID,
			DestinationCode:    value.DestinationCode,
			DestinationName:    value.DestinationName,
			Longitude:          value.Longitude,
			Latitude:           value.Latitude,
			DestinationStatus:  value.DestinationStatus,
			DestinationAddress: value.DestinationAddress,
			AvgSalesWeek:       value.AvgSalesWeek,
			Status:             value.Status,
		}
	}
	return res
}
