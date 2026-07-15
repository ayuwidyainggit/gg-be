package route

import (
	"context"
	"math"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
)

func (service *routeService) GetAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) ([]response.ApprovalRouteMappingEnhanceResponse, response.Meta, error) {
	tx := service.db.Begin()
	if tx.Error != nil {
		return nil, response.Meta{}, tx.Error
	}
	defer helper.CommitOrRollback(tx)

	result, totalCount := service.destinationRepo.GetAllEnhance(ctx, tx, page, limit, filters, currentCustomerId)

	var routes []response.ApprovalRouteMappingEnhanceResponse

	for _, value := range result {
		var verifiedDate string
		if len(value.RouteOutlets) > 0 && value.RouteOutlets[0].VerifiedDate != nil {
			verifiedDate = value.RouteOutlets[0].VerifiedDate.Format("2006-01-02 15:04")
		} else {
			verifiedDate = ""
		}

		routes = append(routes, response.ApprovalRouteMappingEnhanceResponse{
			PjpID:        value.ID,
			PjpCode:      helper.FormatPjpCode(value.PjpCode),
			Status:       value.ApprovalStatus,
			SalesmanName: value.SalesmanName,
			SalesmanCode: value.SalesmanCode,
			Date:         value.CreatedAt.Format("2006-01-02 15:04"),
			VerifiedDate: verifiedDate,
		})
	}

	pagination := response.Meta{
		TotalData: int(totalCount),
		Page:      page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(totalCount) / float64(limit))),
	}

	return routes, pagination, nil
}
