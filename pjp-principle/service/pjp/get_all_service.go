package pjp

import (
	"context"
	"math"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
)

func (service *pjpService) GetAll(
	ctx context.Context,
	limit int,
	page int,
	filters map[string]interface{},
	currentCustomerId string,
) ([]response.PjpResponse, response.Meta, error) {

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	result, totalData := service.pjpRepository.GetAll(ctx, tx, limit, page, filters, currentCustomerId)
	payload := make([]response.PjpResponse, 0, len(result))

	for _, value := range result {
		payload = append(payload, toPjpResponse(value))
	}

	meta := response.Meta{
		TotalData: int(totalData),
		Page:      page,
		Limit:     limit,
		TotalPage: int(math.Ceil(float64(totalData) / float64(limit))),
	}

	return payload, meta, nil
}
