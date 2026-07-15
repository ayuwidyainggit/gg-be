package service

import (
	"context"
	"errors"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func (service *pjpServiceImpl) SubmitPjpDistributor(ctx context.Context, request entity.SubmitPjpDistributorRequest, custID string) error {
	now := time.Now()

	salesmanInfo, err := service.PjpDistributorRepository.GetSalesmanAndTeam(ctx, request.SalesmanId)
	if err != nil {
		log.Error("PjpService, SubmitPjpDistributor, GetSalesmanAndTeam, err:", err.Error())
		return err
	}

	existPJP, err := service.PjpDistributorRepository.FindOnePJPByCodeAndCustId(ctx, request.PjpCode, custID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("PjpService, SubmitPjpDistributor, FindOnePJPByCodeAndCustId, err:", err.Error())
		return err
	}
	if existPJP.ID != 0 {
		return errors.New("PJP already exist")
	}

	whName, err := service.PjpDistributorRepository.GetWarehouseName(ctx, request.WarehouseId, custID)
	if err != nil {
		return errors.New("warehouse not found")
	}

	pjpModel := model.PermanentJourneyPlan{
		PjpCode:        request.PjpCode,
		OperationType:  request.OperationType,
		SalesmanId:     request.SalesmanId, // from fe this is emp.id
		TeamSalesman:   salesmanInfo.SalesTeamName,
		SalesmanName:   salesmanInfo.SalesName,
		SalesmanCode:   salesmanInfo.EmpCode,
		WarehouseId:    request.WarehouseId,
		WarehouseName:  &whName,
		PjpMode:        "Manual",  // Hardcode as per requirement
		Status:         "Pending", // Hardcode as per requirement
		CustId:         custID,
		CreatedAt:      now,
		ApprovalStatus: &request.ApprovalStatus, // Use value from request
	}

	var pjpId int64
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := service.PjpDistributorRepository.StorePermanentJourneyPlan(txCtx, &pjpModel); err != nil {
			log.Error("PjpService, SubmitPjpDistributor, StorePermanentJourneyPlan, err:", err.Error())
			return err
		}
		pjpId = pjpModel.ID

		lastRouteCode, err := service.PjpDistributorRepository.GetLastRouteCode(txCtx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("PjpService, SubmitPjpDistributor, GetLastRouteCode, err:", err.Error())
			return err
		}

		for i := 0; i < len(request.VisitDay); i++ {
			lastRouteCode++
			visitDay := request.VisitDay[i]
			routeModel := model.Route{
				RouteCode: lastRouteCode,
				RouteName: visitDay.Visit.RouteName,
				Sequence:  int64(visitDay.IndexDay),
				IsAssign:  nil,
				CustId:    custID,
				CreatedAt: now,
				UpdatedAt: nil,
				PjpId:     pjpId,
				IsPjpOld:  nil,
			}

			if err := service.PjpDistributorRepository.StoreRoute(txCtx, &routeModel); err != nil {
				log.Error("PjpService, SubmitPjpDistributor, StoreRoute, err:", err.Error())
				return err
			}

			routeCodeStr := strconv.FormatInt(routeModel.RouteCode, 10)

			for _, outlet := range visitDay.Visit.OutletDestinations {
				outletStatus := outlet.OutletStatus
				if outletStatus == "" {
					outletStatus = "0"
				}

				avgSalesWeek := 0.0
				if outlet.AvgSalesWeek != nil {
					avgSalesWeek = *outlet.AvgSalesWeek
				}

				if outlet.Longitude == "" && outlet.Latitude == "" {
					return fmt.Errorf("longitude and latitude for outlet %s cannot be empty", outlet.OutletName)
				}

				routeOutletModel := model.RouteOutlet{
					RouteCode:     routeCodeStr,
					RouteName:     visitDay.Visit.RouteName,
					OutletId:      outlet.OutletId,
					OutletCode:    outlet.OutletCode,
					OutletName:    outlet.OutletName,
					Longitude:     outlet.Longitude,
					Latitude:      outlet.Latitude,
					OutletStatus:  outletStatus,
					OutletAddress: outlet.Address1,
					PjpId:         pjpId,
					PjpCode:       request.PjpCode,
					CustId:        custID,
					Status:        "pending",
					CreatedAt:     now,
					UpdatedAt:     nil,
					VerifiedDate:  nil,
					OldPjpId:      nil,
					OldPjpCode:    nil,
					OldRouteCode:  nil,
					OldRouteName:  nil,
					Photo:         nil,
					Signature:     nil,
					AvgSalesWeek:  &avgSalesWeek,
				}
				if err := service.PjpDistributorRepository.StoreRouteOutlet(txCtx, &routeOutletModel); err != nil {
					log.Error("PjpService, SubmitPjpDistributor, StoreRouteOutlet, err:", err.Error())
					return err
				}

				// Insert route_outlet_history
				routeOutletHistoryModel := model.RouteOutletHistory{
					RouteCode:       routeCodeStr,
					RouteName:       visitDay.Visit.RouteName,
					OutletId:        outlet.OutletId,
					OutletCode:      outlet.OutletCode,
					OutletName:      outlet.OutletName,
					Longitude:       outlet.Longitude,
					Latitude:        outlet.Latitude,
					OutletStatus:    outletStatus,
					OutletAddress:   outlet.Address1,
					PjpId:           pjpId,
					PjpCode:         request.PjpCode,
					CustId:          custID,
					Status:          "pending",
					CreatedAt:       now,
					UpdatedAt:       nil,
					VerifiedDate:    nil,
					OldPjpId:        nil,
					OldPjpCode:      nil,
					OldRouteCode:    nil,
					OldRouteName:    nil,
					Photo:           nil,
					Signature:       nil,
					AvgSalesWeek:    &avgSalesWeek,
					IndexDay:        visitDay.IndexDay,
					StartWeek:       visitDay.StartWeek.Time,
					IsInCurrentYear: visitDay.IsInCurrentYear,
					Week:            visitDay.Week,
					Year:            visitDay.Year,
					Date:            visitDay.Date.Time,
					IsAdditional:    nil,
				}
				if err := service.PjpDistributorRepository.StoreRouteOutletHistory(txCtx, &routeOutletHistoryModel); err != nil {
					log.Error("PjpService, SubmitPjpDistributor, StoreRouteOutletHistory, err:", err.Error())
					return err
				}
			}

			// Insert route_pop_permanent
			routePopPermanentModel := model.RoutePopPermanent{
				Year:      visitDay.Year,
				Week:      visitDay.Week,
				Date:      visitDay.Date.Time,
				Day:       visitDay.Day,
				RouteCode: routeCodeStr,
				PjpId:     pjpId,
				PjpCode:   request.PjpCode,
				CustId:    custID,
				CreatedAt: nil, // NULL as per requirement
				UpdatedAt: nil,
			}
			if err := service.PjpDistributorRepository.StoreRoutePopPermanent(txCtx, &routePopPermanentModel); err != nil {
				log.Error("PjpService, SubmitPjpDistributor, StoreRoutePopPermanent, err:", err.Error())
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *pjpServiceImpl) UpdatePJPDistributor(ctx context.Context, payload entity.UpdatePJPDistributorRequest) error {
	existPJP, err := service.PjpDistributorRepository.FindOnePJPByCodeAndSalesman(ctx, payload.PJPCode, payload.SalesmanID, payload.CustomerID)
	if err != nil {
		return err
	}

	if existPJP.PjpCode != payload.PJPCode {
		return errors.New("PJP not matched")
	}

	// find all route outlets by params (pjp_id, route_id, cust_id)
	var (
		now                     = time.Now()
		existingRouteOutletsMap = make(map[string]model.RouteOutlet)
		existingRouteOutlets    []model.RouteOutlet
		generateKey             = func(routeCode string, outletID int64) string {
			return fmt.Sprintf("%s:%d", routeCode, outletID)
		}
	)

	for _, upvd := range payload.VisitDay {
		routeOutlets, err := service.PjpDistributorRepository.FindAllRouteOutletsByParams(ctx,
			existPJP.ID,
			upvd.Visit.RouteCode,
			payload.CustomerID,
		)

		if err != nil {
			return err
		}

		// collect all route outlets into single slice
		existingRouteOutlets = append(existingRouteOutlets, routeOutlets...)
		for _, ro := range routeOutlets {
			keyOutlet := generateKey(ro.RouteCode, ro.OutletId)
			existingRouteOutletsMap[keyOutlet] = ro
		}
	}

	// find all route outlet history by params (pjp_id, route_id, cust_id)
	var existingRouteOutletHistories []model.RouteOutletHistory
	var existinRouteOutletHistoriesMap = make(map[string]model.RouteOutletHistory)
	for _, upvd := range payload.VisitDay {
		routeOutletHistories, err := service.PjpDistributorRepository.FindAllRouteOutletHistoriesByPJPID(ctx,
			existPJP.ID,
			upvd.Visit.RouteCode,
			payload.CustomerID,
			upvd.Date.Format(time.DateOnly),
		)
		if err != nil {
			return err
		}

		existingRouteOutletHistories = append(existingRouteOutletHistories, routeOutletHistories...)
		for _, roh := range routeOutletHistories {
			keyOutletHist := generateKey(roh.RouteCode, roh.OutletId)
			existinRouteOutletHistoriesMap[keyOutletHist] = roh
		}
	}

	var (
		newRouteOutlets            = make([]model.RouteOutlet, 0)
		updatedRouteOutlets        = make(map[int64]model.RouteOutlet)
		newRouteOutletHistories    = make([]model.RouteOutletHistory, 0)
		updateRouteOutletHistories = make(map[int64]model.RouteOutletHistory)
		newRoutePopPermanents      = make([]model.RoutePopPermanent, 0)
		newRoutePopPermanetMap     = make(map[string]model.RoutePopPermanent)
		weeks                      = make(map[int]int)
		years                      = make(map[int]int)
		routeCodes                 = make(map[int64]int64)
	)

	for _, upvd := range payload.VisitDay {
		weeks[upvd.Week] = upvd.Week
		years[upvd.Year] = upvd.Year
		routeCodes[upvd.Visit.RouteCode] = upvd.Visit.RouteCode
		for _, upod := range upvd.Visit.OutletDestinations {
			var (
				isNewRouteOutlet        = true
				isNewRouteOutletHistory = true
			)

			// modify / update route outlet
			keyOutlet := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), upod.OutletID)
			if existingRouteOutlet, ok := existingRouteOutletsMap[keyOutlet]; ok {
				existingRouteOutlet.RouteName = upvd.Visit.RouteName
				existingRouteOutlet.OutletName = upod.OutletName
				existingRouteOutlet.UpdatedAt = &now

				if upod.OutletCode != "" {
					existingRouteOutlet.OutletCode = upod.OutletCode
				}

				if upod.Address1 != "" {
					existingRouteOutlet.OutletAddress = upod.Address1
				}

				if upod.OutletStatus != "" {
					existingRouteOutlet.OutletStatus = upod.OutletStatus
				}

				if upod.Latitude != "" {
					existingRouteOutlet.Latitude = upod.Latitude
				}

				if upod.Longitude != "" {
					existingRouteOutlet.Longitude = upod.Longitude
				}

				if upod.AvgSalesWeek != nil {
					existingRouteOutlet.AvgSalesWeek = upod.AvgSalesWeek
				}

				updatedRouteOutlets[existingRouteOutlet.ID] = existingRouteOutlet
				isNewRouteOutlet = false
			}

			// new route outlet
			if isNewRouteOutlet {
				newRouteOutlets = append(newRouteOutlets, model.RouteOutlet{
					RouteCode:     fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:     upvd.Visit.RouteName,
					OutletId:      upod.OutletID,
					OutletCode:    upod.OutletCode,
					OutletName:    upod.OutletName,
					OutletAddress: upod.Address1,
					OutletStatus:  upod.OutletStatus,
					Latitude:      upod.Latitude,
					Longitude:     upod.Longitude,
					AvgSalesWeek:  upod.AvgSalesWeek,
					PjpId:         existPJP.ID,
					PjpCode:       payload.PJPCode,
					CustId:        payload.CustomerID,
					CreatedAt:     now,
					Status:        "pending",
				})
			}

			// route outlet history
			keyOutletHist := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), upod.OutletID)
			if existingRouteOutletHistory, ok := existinRouteOutletHistoriesMap[keyOutletHist]; ok {
				existingRouteOutletHistory.RouteCode = fmt.Sprintf("%d", upvd.Visit.RouteCode)
				existingRouteOutletHistory.RouteName = upvd.Visit.RouteName
				existingRouteOutletHistory.OutletId = upod.OutletID
				existingRouteOutletHistory.OutletName = upod.OutletName
				existingRouteOutletHistory.OutletAddress = upod.Address1
				existingRouteOutletHistory.OutletStatus = upod.OutletStatus
				existingRouteOutletHistory.Longitude = upod.Longitude
				existingRouteOutletHistory.Latitude = upod.Latitude
				existingRouteOutletHistory.PjpId = existPJP.ID
				existingRouteOutletHistory.PjpCode = payload.PJPCode
				existingRouteOutletHistory.CustId = payload.CustomerID
				existingRouteOutletHistory.AvgSalesWeek = upod.AvgSalesWeek
				existingRouteOutletHistory.IndexDay = upvd.IndexDay
				existingRouteOutletHistory.StartWeek = time.Time(upvd.StartWeek.Time)
				existingRouteOutletHistory.IsInCurrentYear = upvd.IsInCurrentYear
				existingRouteOutletHistory.Week = upvd.Week
				existingRouteOutletHistory.Year = upvd.Year
				existingRouteOutletHistory.Date = time.Time(upvd.Date.Time)
				existingRouteOutletHistory.UpdatedAt = &now

				if upod.OutletCode != "" {
					existingRouteOutletHistory.OutletCode = upod.OutletCode
				}

				updateRouteOutletHistories[existingRouteOutletHistory.ID] = existingRouteOutletHistory
				isNewRouteOutletHistory = false
			}

			if isNewRouteOutletHistory {
				newRouteOutletHistories = append(newRouteOutletHistories, model.RouteOutletHistory{
					RouteCode:       fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:       upvd.Visit.RouteName,
					OutletId:        upod.OutletID,
					OutletCode:      upod.OutletCode,
					OutletName:      upod.OutletName,
					Longitude:       upod.Longitude,
					Latitude:        upod.Latitude,
					OutletStatus:    upod.OutletStatus,
					OutletAddress:   upod.Address1,
					PjpId:           existPJP.ID,
					PjpCode:         payload.PJPCode,
					CustId:          payload.CustomerID,
					Status:          "pending",
					CreatedAt:       now,
					AvgSalesWeek:    upod.AvgSalesWeek,
					IndexDay:        upvd.IndexDay,
					StartWeek:       upvd.StartWeek.Time,
					IsInCurrentYear: upvd.IsInCurrentYear,
					Week:            upvd.Week,
					Year:            upvd.Year,
					Date:            upvd.Date.Time,
				})
			}

			// route permament, collect all data
			routePopPermanent := model.RoutePopPermanent{
				Year:      upvd.Year,
				Week:      upvd.Week,
				Date:      upvd.Date.Time,
				Day:       upvd.Day,
				RouteCode: fmt.Sprintf("%d", upvd.Visit.RouteCode),
				PjpId:     existPJP.ID,
				PjpCode:   payload.PJPCode,
				CustId:    payload.CustomerID,
			}

			keyRoutePopPermanent := fmt.Sprintf("%s:%s", upvd.Date.Format(time.DateOnly), routePopPermanent.RouteCode)
			if _, ok := newRoutePopPermanetMap[keyRoutePopPermanent]; !ok {
				newRoutePopPermanetMap[keyRoutePopPermanent] = routePopPermanent
			}
		}
	}

	routePOPPPermanents, err := service.PjpDistributorRepository.FindAllRoutePopPermanentsByParams(ctx, existPJP.ID, payload.CustomerID)
	if err != nil {
		return err
	}

	for _, rpp := range routePOPPPermanents {
		keyRoutePopPermanent := fmt.Sprintf("%s:%s", rpp.Date.Format(time.DateOnly), rpp.RouteCode)
		delete(newRoutePopPermanetMap, keyRoutePopPermanent)
	}

	for _, rpp := range newRoutePopPermanetMap {
		newRoutePopPermanents = append(newRoutePopPermanents, rpp)
	}

	var week []int
	var year []int
	var routes []int
	for _, w := range weeks {
		week = append(week, w)
	}
	for _, y := range years {
		year = append(year, y)
	}
	for _, r := range routeCodes {
		routes = append(routes, int(r))
	}

	mapCountUsedRouteOutlet, err := service.PjpDistributorRepository.CheckRouteOutlet(ctx, existPJP.ID, routes, week, year)
	if err != nil {
		return err
	}

	// find old route outlet and delete
	var deleteRouteOutlets = make([]int64, 0)
	for _, routeOutlet := range existingRouteOutlets {
		countUsed := 0
		if _, oks := mapCountUsedRouteOutlet[routeOutlet.RouteCode]; oks {
			countUsed = mapCountUsedRouteOutlet[routeOutlet.RouteCode]
		}
		_, ok := updatedRouteOutlets[routeOutlet.ID]
		if !ok && countUsed == 0 {
			deleteRouteOutlets = append(deleteRouteOutlets, routeOutlet.ID)
		}
	}

	// delete route outlet histories
	var deleteRouteOutletHistories = make([]int64, 0)
	for _, roh := range existingRouteOutletHistories {
		if _, ok := updateRouteOutletHistories[roh.ID]; !ok {
			deleteRouteOutletHistories = append(deleteRouteOutletHistories, roh.ID)
		}
	}

	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// update status pjp to need review
		err := service.PjpDistributorRepository.UpdatePJPStatus(txCtx, existPJP.ID, payload.ApprovalStatus)
		if err != nil {
			return err
		}

		for _, newRouteOutlet := range newRouteOutlets {
			err = service.PjpDistributorRepository.StoreRouteOutlet(txCtx, &newRouteOutlet)
			if err != nil {
				return err
			}
		}

		for _, updatedRouteOutlet := range updatedRouteOutlets {
			err = service.PjpDistributorRepository.UpdateRouteOutlet(txCtx, &updatedRouteOutlet)
			if err != nil {
				return err
			}
		}

		if len(deleteRouteOutlets) > 0 {
			err := service.PjpDistributorRepository.BulkDeleteRouteOutlets(txCtx, deleteRouteOutlets)
			if err != nil {
				return err
			}
		}

		for _, roh := range newRouteOutletHistories {
			err := service.PjpDistributorRepository.StoreRouteOutletHistory(txCtx, &roh)
			if err != nil {
				return err
			}
		}

		for _, roh := range updateRouteOutletHistories {
			err := service.PjpDistributorRepository.UpdateRouteOutletHistory(txCtx, &roh)
			if err != nil {
				return err
			}
		}

		if len(deleteRouteOutletHistories) > 0 {
			err = service.PjpDistributorRepository.BulkDeleteRouteOutletsHistories(txCtx, deleteRouteOutletHistories)
			if err != nil {
				return err
			}
		}

		for _, rpp := range newRoutePopPermanents {
			err = service.PjpDistributorRepository.StoreRoutePopPermanent(txCtx, &rpp)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
