package thirdparty

import (
	"context"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"
)

func (service *thirdPartyService) GetOutlet(ctx context.Context, dataFilter model.OutletQueryFilter, custId string) ([]model.OutletDms, response.Meta) {
	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	// Ambil data
	outletDms := service.outletRepo.GetOutletDms(ctx, tx, dataFilter, custId)

	// Hitung total
	total := service.outletRepo.CountOutletDms(ctx, tx, dataFilter, custId)

	// Mapping
	var outlets []model.OutletDms
	for _, dms := range outletDms {
		outlets = append(outlets, model.OutletDms{
			OutletID:      dms.OutletID,
			OutletCode:    dms.OutletCode,
			OutletName:    dms.OutletName,
			OutletAddress: dms.OutletAddress,
			Longitude:     dms.Longitude,
			Latitude:      dms.Latitude,
			OutletStatus:  dms.OutletStatus,
			AvgSalesWeek:  dms.AvgSalesWeek,
			CustID:        dms.CustID,
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

	return outlets, meta
}
