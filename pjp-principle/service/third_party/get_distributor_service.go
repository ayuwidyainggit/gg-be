package thirdparty

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"
)

func (service *thirdPartyService) GetDistributor(ctx context.Context, dataFilter model.DistributorQueryFilter, custId string) ([]model.DistributorDms, response.Meta) {
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data
	disrtributorDms := service.distributorRepo.GetDistributorDms(ctx, tx, dataFilter, custId)

	// Hitung total
	total := service.distributorRepo.CountDistributorDms(ctx, tx, dataFilter, custId)

	// Mapping
	var disrtributors []model.DistributorDms
	for _, dms := range disrtributorDms {
		disrtributors = append(disrtributors, model.DistributorDms{
			DistributorID:      dms.DistributorID,
			DistributorCode:    dms.DistributorCode,
			DistributorName:    dms.DistributorName,
			DistributorAddress: dms.DistributorAddress,
			Longitude:          dms.Longitude,
			Latitude:           dms.Latitude,
			DistributorStatus:  dms.DistributorStatus,
			AvgSalesWeek:       dms.AvgSalesWeek,
			CustID:             dms.CustID,
		})
	}

	// Pagination meta
	page := 1
	limit := int(total)
	if p, err := strconv.Atoi(dataFilter.Page); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(dataFilter.Limit); err == nil && l > 0 {
		limit = l
	}

	totalPage := 1
	if limit > 0 {
		totalPage = int((total + int64(limit) - 1) / int64(limit))
	}

	meta := response.Meta{
		TotalData: int(total),
		Page:      page,
		Limit:     limit,
		TotalPage: totalPage,
	}

	return disrtributors, meta
}
