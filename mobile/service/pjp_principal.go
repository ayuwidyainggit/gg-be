package service

import (
	"context"
	"errors"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type PjpPrincipalService interface {
	SubmitPjpPrincipal(ctx context.Context, request entity.SubmitPjpPrincipalRequest, custId string) error
	UpdatePjpPrincipal(ctx context.Context, payload entity.UpdatePjpPrincipalRequest) error
}

type pjpPrincipalServiceImpl struct {
	PjpPrincipalRepository repository.PjpPrincipalRepository
	MCustomerRepository    repository.MCustomerRepository
	Transaction            repository.Dbtransaction
}

func NewPjpPrincipalService(
	pjpPrincipalRepository repository.PjpPrincipalRepository,
	mCustomerRepository repository.MCustomerRepository,
	transaction repository.Dbtransaction,
) PjpPrincipalService {
	return &pjpPrincipalServiceImpl{
		PjpPrincipalRepository: pjpPrincipalRepository,
		MCustomerRepository:    mCustomerRepository,
		Transaction:            transaction,
	}
}

// SubmitPjpPrincipal submits PJP principal data
func (service *pjpPrincipalServiceImpl) SubmitPjpPrincipal(ctx context.Context, request entity.SubmitPjpPrincipalRequest, custId string) error {
	now := time.Now()

	customer, err := service.MCustomerRepository.FindOneByCustId(custId)
	if err != nil {
		log.Error("PjpPrincipalService, SubmitPjpPrincipal, FindOneByCustId, err:", err.Error())
		return err
	}
	parentCustId := customer.ParentCustID

	var warehouseName *string
	whName, err := service.PjpPrincipalRepository.GetWarehouseName(ctx, request.WarehouseId, parentCustId)
	if err == nil && whName != "" {
		warehouseName = &whName
	}

	salesmanInfo, err := service.PjpPrincipalRepository.GetSalesmanAndTeam(ctx, request.SalesmanId)
	if err != nil {
		log.Error("PjpPrincipalService, SubmitPjpPrincipal, GetSalesmanAndTeam, err:", err.Error())
		return err
	}

	existPJP, err := service.PjpPrincipalRepository.FindOnePJPByCodeAndCustId(ctx, request.PjpCode, custId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("PjpPrincipalService, SubmitPjpPrincipal, FindOnePJPByCodeAndCustId, err:", err.Error())
		return err
	}
	if existPJP.ID != 0 {
		return errors.New("PJP already exist")
	}

	// Create PermanentJourneyPlanPrincipal
	pjpModel := model.PermanentJourneyPlanPrincipal{
		PjpCode:        request.PjpCode,
		OperationType:  request.OperationType,
		SalesmanId:     request.SalesmanId,
		TeamSalesman:   salesmanInfo.SalesTeamName,
		SalesmanCode:   salesmanInfo.EmpCode,
		SalesmanName:   salesmanInfo.SalesName,
		WarehouseId:    request.WarehouseId,
		WarehouseName:  warehouseName, // Get from database
		PjpMode:        "Manual",      // Hardcode as per requirement
		Status:         "Pending",     // Hardcode as per requirement
		CustId:         custId,
		CreatedAt:      now,
		ApprovalStatus: &request.ApprovalStatus, // Use value from request
	}

	var pjpId int64

	// Transaction: insert all related data
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := service.PjpPrincipalRepository.StorePermanentJourneyPlan(txCtx, &pjpModel); err != nil {
			log.Error("PjpPrincipalService, SubmitPjpPrincipal, StorePermanentJourneyPlan, err:", err.Error())
			return err
		}
		pjpId = pjpModel.ID

		lastRouteCode, err := service.PjpPrincipalRepository.GetLastRouteCode(ctx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("PjpPrincipalService, SubmitPjpPrincipal, GetLastRouteCode, err:", err.Error())
			return err
		}

		for i := 0; i < len(request.VisitDay); i++ {
			lastRouteCode++
			visitDay := request.VisitDay[i]

			// Insert route
			routeModel := model.RoutePrincipal{
				RouteCode: lastRouteCode,
				RouteName: visitDay.Visit.RouteName,
				Sequence:  int64(visitDay.IndexDay),
				CustId:    custId,
				CreatedAt: now,
				UpdatedAt: nil,
				PjpId:     pjpId,
			}
			if err := service.PjpPrincipalRepository.StoreRoute(txCtx, &routeModel); err != nil {
				log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreRoute, err:", err.Error())
				return err
			}

			routeCodeStr := strconv.FormatInt(routeModel.RouteCode, 10)

			for _, outlet := range visitDay.Visit.OutletDestinations {
				destinationStatus := outlet.OutletStatus
				if destinationStatus == "" {
					destinationStatus = "0"
				}

				avgSalesWeek := 0.0
				if outlet.AvgSalesWeek != nil {
					avgSalesWeek = *outlet.AvgSalesWeek
				}

				if outlet.Longitude == "" && outlet.Latitude == "" {
					return fmt.Errorf("longitude and latitude for outlet %s cannot be empty", outlet.OutletName)
				}

				// Insert destination (outlet type)
				destinationModel := model.DestinationPrincipal{
					RouteCode:          routeCodeStr,
					RouteName:          visitDay.Visit.RouteName,
					Status:             "pending",
					VerifiedDate:       nil,
					DestinationId:      outlet.OutletId,
					DestinationCode:    outlet.OutletCode,
					DestinationStatus:  destinationStatus,
					DestinationName:    outlet.OutletName,
					DestinationAddress: outlet.Address1,
					DestinationType:    "outlet",
					Longitude:          outlet.Longitude,
					Latitude:           outlet.Latitude,
					PjpId:              pjpId,
					PjpCode:            request.PjpCode,
					OldPjpId:           nil,
					OldPjpCode:         nil,
					OldRouteCode:       nil,
					OldRouteName:       nil,
					Photo:              nil,
					Signature:          nil,
					AvgSalesWeek:       &avgSalesWeek,
					CustId:             custId,
					CreatedAt:          now,
					UpdatedAt:          nil,
				}
				if err := service.PjpPrincipalRepository.StoreDestination(txCtx, &destinationModel); err != nil {
					log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreDestination (outlet), err:", err.Error())
					return err
				}

				// Insert destination_history (outlet type)
				destinationHistoryModel := model.DestinationHistoryPrincipal{
					RouteCode:          routeCodeStr,
					RouteName:          visitDay.Visit.RouteName,
					VerifiedDate:       nil,
					Date:               visitDay.Date.Time,
					Week:               visitDay.Week,
					Year:               visitDay.Year,
					IndexDay:           visitDay.IndexDay,
					StartWeek:          visitDay.StartWeek.Time,
					IsInCurrentYear:    visitDay.IsInCurrentYear,
					IsAdditional:       nil,
					DestinationId:      outlet.OutletId,
					DestinationCode:    outlet.OutletCode,
					DestinationStatus:  destinationStatus,
					DestinationName:    outlet.OutletName,
					DestinationAddress: outlet.Address1,
					DestinationType:    "outlet",
					Longitude:          outlet.Longitude,
					Latitude:           outlet.Latitude,
					PjpId:              pjpId,
					PjpCode:            request.PjpCode,
					OldPjpId:           nil,
					OldPjpCode:         nil,
					OldRouteCode:       nil,
					OldRouteName:       nil,
					Photo:              nil,
					Signature:          nil,
					AvgSalesWeek:       &avgSalesWeek,
					CustId:             custId,
					CreatedAt:          now,
					UpdatedAt:          nil,
				}
				if err := service.PjpPrincipalRepository.StoreDestinationHistory(txCtx, &destinationHistoryModel); err != nil {
					log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreDestinationHistory (outlet), err:", err.Error())
					return err
				}
			}

			// Process distributor destinations
			for _, distributor := range visitDay.Visit.DistributorDestinations {

				if distributor.Longitude == "" && distributor.Latitude == "" {
					return fmt.Errorf("longitude and latitude for distributor %s cannot be empty", distributor.DistributorName)
				}

				destinationModel := model.DestinationPrincipal{
					RouteCode:          routeCodeStr,
					RouteName:          visitDay.Visit.RouteName,
					Status:             "pending",
					VerifiedDate:       nil,
					DestinationId:      distributor.DistributorId,
					DestinationCode:    distributor.DistributorCode,
					DestinationStatus:  "", // Empty for distributor
					DestinationName:    distributor.DistributorName,
					DestinationAddress: distributor.Address,
					DestinationType:    "distributor",
					Longitude:          distributor.Longitude,
					Latitude:           distributor.Latitude,
					PjpId:              pjpId,
					PjpCode:            request.PjpCode,
					OldPjpId:           nil,
					OldPjpCode:         nil,
					OldRouteCode:       nil,
					OldRouteName:       nil,
					Photo:              nil,
					Signature:          nil,
					AvgSalesWeek:       nil, // No avg_sales_week for distributor
					CustId:             custId,
					CreatedAt:          now,
					UpdatedAt:          nil,
				}
				if err := service.PjpPrincipalRepository.StoreDestination(txCtx, &destinationModel); err != nil {
					log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreDestination (distributor), err:", err.Error())
					return err
				}

				// Insert destination_history (distributor type)
				destinationHistoryModel := model.DestinationHistoryPrincipal{
					RouteCode:          routeCodeStr,
					RouteName:          visitDay.Visit.RouteName,
					VerifiedDate:       nil,
					Date:               visitDay.Date.Time,
					Week:               visitDay.Week,
					Year:               visitDay.Year,
					IndexDay:           visitDay.IndexDay,
					StartWeek:          visitDay.StartWeek.Time,
					IsInCurrentYear:    visitDay.IsInCurrentYear,
					IsAdditional:       nil,
					DestinationId:      distributor.DistributorId,
					DestinationCode:    distributor.DistributorCode,
					DestinationStatus:  "", // Empty for distributor
					DestinationName:    distributor.DistributorName,
					DestinationAddress: distributor.Address,
					DestinationType:    "distributor",
					Longitude:          distributor.Longitude,
					Latitude:           distributor.Latitude,
					PjpId:              pjpId,
					PjpCode:            request.PjpCode,
					OldPjpId:           nil,
					OldPjpCode:         nil,
					OldRouteCode:       nil,
					OldRouteName:       nil,
					Photo:              nil,
					Signature:          nil,
					AvgSalesWeek:       nil, // No avg_sales_week for distributor
					CustId:             custId,
					CreatedAt:          now,
					UpdatedAt:          nil,
				}
				if err := service.PjpPrincipalRepository.StoreDestinationHistory(txCtx, &destinationHistoryModel); err != nil {
					log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreDestinationHistory (distributor), err:", err.Error())
					return err
				}
			}

			// Insert route_pop_permanent
			routePopPermanentModel := model.RoutePopPermanentPrincipal{
				Year:      visitDay.Year,
				Week:      visitDay.Week,
				Date:      visitDay.Date.Time,
				Day:       visitDay.Day,
				RouteCode: routeModel.RouteCode,
				PjpId:     pjpId,
				PjpCode:   request.PjpCode,
				CustId:    custId,
				CreatedAt: nil, // NULL as per requirement
				UpdatedAt: nil,
			}
			if err := service.PjpPrincipalRepository.StoreRoutePopPermanent(txCtx, &routePopPermanentModel); err != nil {
				log.Error("PjpPrincipalService, SubmitPjpPrincipal, StoreRoutePopPermanent, err:", err.Error())
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

// UpdatePjpPrincipal updates PJP principal data
func (service *pjpPrincipalServiceImpl) UpdatePjpPrincipal(ctx context.Context, payload entity.UpdatePjpPrincipalRequest) error {

	var (
		now         = time.Now()
		generateKey = func(routeCode, destinationType string, destID int64) string {
			return fmt.Sprintf("%s:%s:%d", routeCode, destinationType, destID)
		}
	)

	existPJP, err := service.PjpPrincipalRepository.FindOnePJPByCodeAndSalesmanID(ctx, payload.PJPCode, payload.SalesmanID, payload.CustomerID)
	if err != nil {
		return err
	}

	if existPJP.PjpCode != payload.PJPCode {
		return errors.New("PJP not matched")
	}

	// find all destinations
	var existingDestinations []model.DestinationPrincipal
	var existingDestinationsMap = make(map[string]model.DestinationPrincipal)
	for _, upvd := range payload.VisitDay {
		destinations, err := service.PjpPrincipalRepository.FindAllDestinationsByParams(ctx,
			existPJP.ID,
			upvd.Visit.RouteCode,
			payload.CustomerID,
		)

		if err != nil {
			return err
		}

		existingDestinations = append(existingDestinations, destinations...)
		for _, dest := range destinations {
			keyDest := generateKey(dest.RouteCode, dest.DestinationType, dest.DestinationId)
			existingDestinationsMap[keyDest] = dest
		}
	}

	// find all destination histories
	var existingDestinationHistories []model.DestinationHistoryPrincipal
	var existingDestinationHistoriesMap = make(map[string]model.DestinationHistoryPrincipal)
	for _, upvd := range payload.VisitDay {
		histories, err := service.PjpPrincipalRepository.FindAllDestinationHistoriesByPJPID(ctx,
			existPJP.ID,
			upvd.Visit.RouteCode,
			payload.CustomerID,
			upvd.Date.Format(time.DateOnly),
		)
		if err != nil {
			return err
		}

		existingDestinationHistories = append(existingDestinationHistories, histories...)
		for _, h := range histories {
			keyDestHist := generateKey(h.RouteCode, h.DestinationType, h.DestinationId)
			existingDestinationHistoriesMap[keyDestHist] = h
		}
	}

	var (
		newDestinations                = make([]model.DestinationPrincipal, 0)
		updatedDestinationsMap         = make(map[string]model.DestinationPrincipal)
		newDestinationHistories        = make([]model.DestinationHistoryPrincipal, 0)
		updatedDestinationHistoriesMap = make(map[int64]model.DestinationHistoryPrincipal)
		newRoutePopPermanents          = make([]model.RoutePopPermanentPrincipal, 0)
		newRoutePopPermanetMap         = make(map[string]model.RoutePopPermanentPrincipal)
		weeks                          = make(map[int]int)
		years                          = make(map[int]int)
		routeCodes                     = make(map[int64]int64)
	)

	for _, upvd := range payload.VisitDay {
		weeks[upvd.Week] = upvd.Week
		years[upvd.Year] = upvd.Year
		routeCodes[upvd.Visit.RouteCode] = upvd.Visit.RouteCode
		// Outlets
		for _, outletDest := range upvd.Visit.OutletDestinations {
			// Modify Destination
			keyOutletDest := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), "outlet", outletDest.OutletID)
			if existingDest, ok := existingDestinationsMap[keyOutletDest]; ok {
				existingDest.RouteName = upvd.Visit.RouteName
				existingDest.DestinationName = outletDest.OutletName
				existingDest.UpdatedAt = &now
				existingDest.DestinationAddress = outletDest.Address1
				existingDest.AvgSalesWeek = outletDest.AvgSalesWeek

				if outletDest.OutletStatus != "" {
					existingDest.DestinationStatus = outletDest.OutletStatus
				}

				if outletDest.Latitude != "" {
					existingDest.Latitude = outletDest.Latitude
				}

				if outletDest.Longitude != "" {
					existingDest.Longitude = outletDest.Longitude
				}

				if outletDest.OutletCode != "" {
					existingDest.DestinationCode = outletDest.OutletCode
				}

				updatedDestinationsMap[keyOutletDest] = existingDest
			} else {
				// New Destination
				newDestinations = append(newDestinations, model.DestinationPrincipal{
					RouteCode:          fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:          upvd.Visit.RouteName,
					DestinationId:      outletDest.OutletID,
					DestinationCode:    outletDest.OutletCode,
					DestinationName:    outletDest.OutletName,
					DestinationAddress: outletDest.Address1,
					DestinationStatus:  outletDest.OutletStatus,
					Latitude:           outletDest.Latitude,
					Longitude:          outletDest.Longitude,
					AvgSalesWeek:       outletDest.AvgSalesWeek,
					DestinationType:    "outlet",
					PjpId:              existPJP.ID,
					PjpCode:            payload.PJPCode,
					CustId:             payload.CustomerID,
					CreatedAt:          now,
					Status:             "pending",
				})
			}

			// Modify History
			keyDestHist := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), "outlet", outletDest.OutletID)
			if existingHist, ok := existingDestinationHistoriesMap[keyDestHist]; ok {
				existingHist.RouteName = upvd.Visit.RouteName
				existingHist.DestinationName = outletDest.OutletName
				existingHist.DestinationAddress = outletDest.Address1
				existingHist.DestinationStatus = outletDest.OutletStatus
				existingHist.Longitude = outletDest.Longitude
				existingHist.Latitude = outletDest.Latitude
				existingHist.AvgSalesWeek = outletDest.AvgSalesWeek
				existingHist.IndexDay = upvd.IndexDay
				existingHist.StartWeek = time.Time(upvd.StartWeek.Time)
				existingHist.IsInCurrentYear = upvd.IsInCurrentYear
				existingHist.Week = upvd.Week
				existingHist.Year = upvd.Year
				existingHist.Date = time.Time(upvd.Date.Time)
				existingHist.UpdatedAt = &now

				if outletDest.OutletCode != "" {
					existingHist.DestinationCode = outletDest.OutletCode
				}

				updatedDestinationHistoriesMap[existingHist.ID] = existingHist
			} else {
				newDestinationHistories = append(newDestinationHistories, model.DestinationHistoryPrincipal{
					RouteCode:          fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:          upvd.Visit.RouteName,
					DestinationId:      outletDest.OutletID,
					DestinationCode:    outletDest.OutletCode,
					DestinationName:    outletDest.OutletName,
					Longitude:          outletDest.Longitude,
					Latitude:           outletDest.Latitude,
					DestinationStatus:  outletDest.OutletStatus,
					DestinationAddress: outletDest.Address1,
					DestinationType:    "outlet",
					PjpId:              existPJP.ID,
					PjpCode:            payload.PJPCode,
					CustId:             payload.CustomerID,
					CreatedAt:          now,
					AvgSalesWeek:       outletDest.AvgSalesWeek,
					IndexDay:           upvd.IndexDay,
					StartWeek:          upvd.StartWeek.Time,
					IsInCurrentYear:    upvd.IsInCurrentYear,
					Week:               upvd.Week,
					Year:               upvd.Year,
					Date:               upvd.Date.Time,
				})
			}
		}

		// Distributors
		for _, dist := range upvd.Visit.DistributorDestinations {
			keyDest := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), "distributor", dist.DistributorID)
			if existingDest, ok := existingDestinationsMap[keyDest]; ok {
				existingDest.RouteName = upvd.Visit.RouteName
				existingDest.UpdatedAt = &now

				updatedDestinationsMap[keyDest] = existingDest
			} else {
				newDestinations = append(newDestinations, model.DestinationPrincipal{
					RouteCode:          fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:          upvd.Visit.RouteName,
					DestinationId:      dist.DistributorID,
					DestinationCode:    dist.DistributorCode,
					DestinationName:    dist.DistributorName,
					DestinationAddress: dist.Address,
					Latitude:           dist.Latitude,
					Longitude:          dist.Longitude,
					DestinationType:    "distributor",
					PjpId:              existPJP.ID,
					PjpCode:            payload.PJPCode,
					CustId:             payload.CustomerID,
					CreatedAt:          now,
					Status:             "pending",
				})
			}

			keyDestHist := generateKey(fmt.Sprintf("%d", upvd.Visit.RouteCode), "distributor", dist.DistributorID)
			if existingHist, ok := existingDestinationHistoriesMap[keyDestHist]; ok {
				existingHist.RouteName = upvd.Visit.RouteName
				existingHist.UpdatedAt = &now
				updatedDestinationHistoriesMap[existingHist.ID] = existingHist
			} else {
				newDestinationHistories = append(newDestinationHistories, model.DestinationHistoryPrincipal{
					RouteCode:          fmt.Sprintf("%d", upvd.Visit.RouteCode),
					RouteName:          upvd.Visit.RouteName,
					DestinationId:      dist.DistributorID,
					DestinationCode:    dist.DistributorCode,
					DestinationName:    dist.DistributorName,
					Longitude:          dist.Longitude,
					Latitude:           dist.Latitude,
					DestinationAddress: dist.Address,
					DestinationType:    "distributor",
					PjpId:              existPJP.ID,
					PjpCode:            payload.PJPCode,
					CustId:             payload.CustomerID,
					CreatedAt:          now,
					IndexDay:           upvd.IndexDay,
					StartWeek:          upvd.StartWeek.Time,
					IsInCurrentYear:    upvd.IsInCurrentYear,
					Week:               upvd.Week,
					Year:               upvd.Year,
					Date:               upvd.Date.Time,
				})
			}
		}

		// route permament, collect all data
		routePopPermanent := model.RoutePopPermanentPrincipal{
			Year:      upvd.Year,
			Week:      upvd.Week,
			Date:      upvd.Date.Time,
			Day:       upvd.Day,
			RouteCode: upvd.Visit.RouteCode,
			PjpId:     existPJP.ID,
			PjpCode:   payload.PJPCode,
			CustId:    payload.CustomerID,
		}

		keyRoutePopPermanent := fmt.Sprintf("%s:%d", upvd.Date.Format(time.DateOnly), routePopPermanent.RouteCode)
		if _, ok := newRoutePopPermanetMap[keyRoutePopPermanent]; !ok {
			newRoutePopPermanetMap[keyRoutePopPermanent] = routePopPermanent
		}
	}

	routePOPPPermanents, err := service.PjpPrincipalRepository.FindAllRoutePopPermanentsByParams(ctx, existPJP.ID, payload.CustomerID)
	if err != nil {
		return err
	}

	for _, rpp := range routePOPPPermanents {
		keyRoutePopPermanent := fmt.Sprintf("%s:%d", rpp.Date.Format(time.DateOnly), rpp.RouteCode)
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

	mapCountUsedRouteOutlet, err := service.PjpPrincipalRepository.CheckRouteOutlet(ctx, existPJP.ID, routes, week, year)
	if err != nil {
		return err
	}

	// Delete logic
	var deleteDestinations = make([]int64, 0)
	for _, dest := range existingDestinations {
		keyDest := generateKey(dest.RouteCode, dest.DestinationType, dest.DestinationId)
		countUsed := 0
		if _, ok := mapCountUsedRouteOutlet[dest.RouteCode]; ok {
			countUsed = mapCountUsedRouteOutlet[dest.RouteCode]
		}
		if _, ok := updatedDestinationsMap[keyDest]; !ok && countUsed == 0 {
			deleteDestinations = append(deleteDestinations, dest.ID)
		}
	}

	var deleteDestinationHistories = make([]int64, 0)
	for _, h := range existingDestinationHistories {
		if _, ok := updatedDestinationHistoriesMap[h.ID]; !ok {
			deleteDestinationHistories = append(deleteDestinationHistories, h.ID)
		}
	}

	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		err := service.PjpPrincipalRepository.UpdatePJPStatus(txCtx, existPJP.ID, payload.ApprovalStatus)
		if err != nil {
			return err
		}

		for _, newDest := range newDestinations {
			err = service.PjpPrincipalRepository.StoreDestination(txCtx, &newDest)
			if err != nil {
				return err
			}
		}

		for _, updatedDest := range updatedDestinationsMap {
			err = service.PjpPrincipalRepository.UpdateDestination(txCtx, &updatedDest)
			if err != nil {
				return err
			}
		}

		if len(deleteDestinations) > 0 {
			err = service.PjpPrincipalRepository.BulkDeleteDestinations(txCtx, deleteDestinations)
			if err != nil {
				return err
			}
		}

		for _, newHist := range newDestinationHistories {
			err = service.PjpPrincipalRepository.StoreDestinationHistory(txCtx, &newHist)
			if err != nil {
				return err
			}
		}

		for _, updatedHist := range updatedDestinationHistoriesMap {
			err = service.PjpPrincipalRepository.UpdateDestinationHistory(txCtx, &updatedHist)
			if err != nil {
				return err
			}
		}

		if len(deleteDestinationHistories) > 0 {
			err = service.PjpPrincipalRepository.BulkDeleteDestinationHistories(txCtx, deleteDestinationHistories)
			if err != nil {
				return err
			}
		}

		for _, rpp := range newRoutePopPermanents {
			err = service.PjpPrincipalRepository.StoreRoutePopPermanent(txCtx, &rpp)
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
