package route

import (
	"context"
	"log"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strconv"
	"time"
)

func (service *routeService) UpdateStatusEnhance(ctx context.Context, request request.UpdateStatusEnhanceRequest, custId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	for i := range request.PjpCode {
		pjpCode, _ := strconv.Atoi(request.PjpCode[i])
		route := service.destinationRepo.FindAllOutletsByPjpCode(ctx, tx, pjpCode, custId)

		pjp := model.Pjp{}
		pjp.ID = *route[0].PjpID

		if request.Status == "Approved" {
			pjp.ApprovalStatus = "Approved"
			pjp.Status = "true"
			service.pjpRepo.Update(ctx, tx, pjp)

			routes := service.routePopPermanentRepo.FindByPjpID(ctx, tx, *route[0].PjpID, custId)
			if len(routes) != 0 {
				for _, route := range routes {
					dataset := model.RoutePopDaily{
						RouteCode: route.RouteCode,
						Week:      route.Week,
						Day:       route.Day,
						Date:      route.Date,
						PjpCode:   route.PjpCode,
						PjpID:     route.PjpID,
						Year:      route.Year,
						CustID:    route.CustID,
						Status:    "permanent",
					}
					service.routePopDailyRepo.Create(ctx, tx, dataset)
				}
			}

		} else {
			pjp.ApprovalStatus = "Rejected"
			service.pjpRepo.Update(ctx, tx, pjp)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Println("Gagal load lokasi waktu, fallback ke WIB manual:", err)
			loc = time.FixedZone("WIB", 7*3600) // fallback jika LoadLocation gagal
		}

		now := time.Now().In(loc)

		for _, data := range route {
			routes := model.Destination{
				ID:           data.ID,
				Status:       request.Status,
				VerifiedDate: &now,
				// RouteCode:    route.RouteCode,
				// DestinationCode:   request.DestinationCode[i],
			}
			service.destinationRepo.UpdatePivot(ctx, tx, routes)
		}
	}
}
