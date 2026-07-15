package service

import (
	"context"
	"errors"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/times"
	"mobile/repository"
	"strconv"
	"time"
)

type ExtraCallService interface {
	Create(request entity.CreateExtraCallRequest) error
}

type ExtraCallServiceImpl struct {
	Transaction        repository.Dbtransaction
	PjpDistributorRepo repository.PjpDistributorRepository
	PjpPrincipalRepo   repository.PjpPrincipalRepository
	MOutletRepo        repository.MOutletRepository
}

func NewExtraCallService(
	transaction repository.Dbtransaction,
	pjpDistributorRepo repository.PjpDistributorRepository,
	pjpPrincipalRepo repository.PjpPrincipalRepository,
	mOutletRepo repository.MOutletRepository,
) *ExtraCallServiceImpl {
	return &ExtraCallServiceImpl{
		Transaction:        transaction,
		PjpDistributorRepo: pjpDistributorRepo,
		PjpPrincipalRepo:   pjpPrincipalRepo,
		MOutletRepo:        mOutletRepo,
	}
}

func (service *ExtraCallServiceImpl) Create(request entity.CreateExtraCallRequest) error {
	ctx := context.Background()
	var err error
	var routeCode int
	route := model.Route{}

	isPrincipal := len(request.CustID) == 6 && !request.IsDistributor
	startDefault := time.Now().UnixMilli()
	if !isPrincipal {
		pjpD, err := service.PjpDistributorRepo.FindPJPByEmpIDAndCustID(ctx, int64(request.EmpID), request.CustID)
		if err != nil {
			return err
		}

		if request.RouteCode == "" {
			request.RouteCode = pjpD.RouteCode
		}

		if request.PJPCode == "" {
			request.PJPCode = pjpD.PJPCode
		}

		if request.PJPID == 0 {
			request.PJPID = pjpD.PJPID
		}
		routeCode, err = strconv.Atoi(request.RouteCode)
		if err != nil {
			return fmt.Errorf("invalid route code")
		}

		route, err = service.PjpDistributorRepo.GetRouteByRouteCode(ctx, routeCode)
		if err != nil {
			return err
		}

	} else {
		pjpD, err := service.PjpPrincipalRepo.FindPJPByEmpIDAndCustID(ctx, int64(request.EmpID), request.CustID)
		if err != nil {
			return err
		}

		if request.RouteCode == "" {
			request.RouteCode = pjpD.RouteCode
		}

		if request.PJPCode == "" {
			request.PJPCode = pjpD.PJPCode
		}

		if request.PJPID == 0 {
			request.PJPID = pjpD.PJPID
		}
		routeCode, err = strconv.Atoi(request.RouteCode)
		if err != nil {
			return fmt.Errorf("invalid route code")
		}
		route, err = service.PjpPrincipalRepo.GetRouteByRouteCode(ctx, routeCode)
		if err != nil {
			return err
		}
	}

	outlets, err := service.MOutletRepo.FindByDestinations(ctx, request.DestinationType, request.DestinationIDs)
	if err != nil {
		return err
	}

	if len(outlets) == 0 {
		return errors.New("outlet not found")
	}

	trx, err := service.MOutletRepo.TrxBegin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			trx.TrxRollback()
		}
	}()

	if route.PjpId != request.PJPID {
		return errors.New("route not found")
	}

	pjpCode, err := strconv.Atoi(request.PJPCode)
	if err != nil {
		return fmt.Errorf("invalid pjp code")
	}
	pjpID := int(request.PJPID)
	starts := int(request.Start)
	if starts == 0 {
		starts = int(startDefault)
	}

	timeNow, errTime := times.GetCurrentTime()
	if errTime != nil {
		return errTime
	}
	var (
		isAdditional         = true
		isCurrentYear        = true
		isExtraCall          = true
		startWeek            = timeNow.Format(time.DateOnly)
		startOfYear          = time.Date(timeNow.Year(), 1, 1, 0, 0, 0, 0, timeNow.Location())
		daysSinceStartOfYear = int(timeNow.Sub(startOfYear).Hours() / 24)
		week                 = (daysSinceStartOfYear+int(startOfYear.Weekday()))/7 + 1
		yearNow              = timeNow.Year()
	)

	for _, rows := range outlets {
		outletData := model.MOutletCreadFromList{
			RouteCode:       &routeCode,
			RouteName:       &route.RouteName,
			OutletId:        rows.OutletID,
			OutletCode:      rows.OutletCode,
			OutletName:      rows.OutletName,
			Longitude:       rows.Longitude,
			Latitude:        rows.Latitude,
			OutletAddress:   rows.Address,
			PjpId:           &pjpID,
			PjpCode:         &pjpCode,
			CustId:          request.CustID,
			Status:          "pending",
			CreatedAt:       &timeNow,
			UpdatedAt:       &timeNow,
			IsAdditional:    &isAdditional,
			StartWeek:       &startWeek,
			IsInCurrentYear: &isCurrentYear,
			Week:            &week,
			Year:            &yearNow,
			Date:            &timeNow,
			IsExtraCall:     &isExtraCall,
		}
		if rows.OutletStatus != nil {
			outletData.OutletStatus = *rows.OutletStatus
		}

		var (
			onlyDateVisitList   = time.Date(outletData.Date.Year(), outletData.Date.Month(), outletData.Date.Day(), 0, 0, 0, 0, outletData.Date.Location())
			isPlannedVisitList  = false
			outletDataVisitList = model.MOutletCreadFromListOutletVisitList{
				Week:            outletData.Week,
				Year:            outletData.Year,
				Date:            &onlyDateVisitList,
				Day:             outletData.Date.Format("Mon"),
				RouteCode:       outletData.RouteCode,
				OutletCode:      outletData.OutletCode,
				PjpId:           outletData.PjpId,
				PjpCode:         outletData.PjpCode,
				CreatedAt:       outletData.CreatedAt,
				OutletId:        outletData.OutletId,
				IsPlanned:       &isPlannedVisitList,
				IsExtraCall:     &isExtraCall,
				Latitude:        outletData.Latitude,
				Longitude:       outletData.Longitude,
				Start:           &starts,
				DestinationType: &request.DestinationType,
			}
		)

		if request.IsDistributor {
			err = trx.StoreFromList(ctx, &outletData)
			if err != nil {
				return err
			}

			err = trx.StoreFromListOutletVisitList(ctx, &outletDataVisitList)
			if err != nil {
				return err
			}
		} else {
			err = trx.StoreFromListPrinciple(ctx, &outletData)
			if err != nil {
				return err
			}

			err = trx.StoreFromListOutletVisitListPrinciple(ctx, &outletDataVisitList)
			if err != nil {
				return err
			}
		}
	}

	err = trx.TrxCommit()
	if err != nil {
		return err
	}

	return nil
}
