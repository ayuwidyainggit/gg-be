package service

import (
	"context"
	"errors"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/times"
	"mobile/repository"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type (
	WeekService interface {
		List(dataFilter entity.WeekListQueryFilter) (data []entity.WeekListResponse, total int64, lastPage int, err error)
	}

	weekServiceImpl struct {
		WeekRepository           repository.WeekRepository
		PJPDistributorRepository repository.PjpDistributorRepository
		PJPPrincipalRepository   repository.PjpPrincipalRepository
	}
)

func NewWeekService(
	weekRepository repository.WeekRepository,
	pjpDistributorRepository repository.PjpDistributorRepository,
	pjpPrincipalRepository repository.PjpPrincipalRepository,

) WeekService {
	return &weekServiceImpl{
		WeekRepository:           weekRepository,
		PJPDistributorRepository: pjpDistributorRepository,
		PJPPrincipalRepository:   pjpPrincipalRepository,
	}
}

func (service *weekServiceImpl) List(dataFilter entity.WeekListQueryFilter) (data []entity.WeekListResponse, total int64, lastPage int, err error) {
	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "week_start:desc"
	}

	// Get weeks
	weeks, total, lastPage, err := service.WeekRepository.FindAll(dataFilter)
	if err != nil {
		log.Error("WeekService, List, FindAll, err:", err.Error())
		return data, 0, 0, errors.New("record not found")
	}

	// Map weeks to response
	data = make([]entity.WeekListResponse, 0)
	for _, week := range weeks {
		var workDays []model.WorkDayDetail
		var err error

		workDays, err = service.WeekRepository.FindWorkDaysByWeekId(week.WeekId, week.WeekStart.Format(time.DateOnly), week.WeekEnd.Format(time.DateOnly), dataFilter.ParentCustID, false)
		if err != nil {
			log.Error("WeekService, List, FindWorkDaysByWeekId, err:", err.Error())
			workDays = []model.WorkDayDetail{}
		}

		workDayData := make([]entity.WorkDayData, 0)
		for idx, wd := range workDays {
			var (
				routeCode           int64
				routeName           string
				numberOfOutlet      int
				numberOfDistributor int
				seqDay              int64 = int64(idx + 1)
			)

			// Call different repository method based on user type
			if dataFilter.IsDistributor {
				// Fetch dynamic route name
				rc, rn, err := service.PJPDistributorRepository.GetRouteNameByIndexDay(context.Background(), dataFilter.CustID, dataFilter.EmpID, seqDay)
				if err == nil && rn != "" {
					routeCode = rc
					routeName = rn
				} else {
					routeName = fmt.Sprintf("Route %d", idx+1)
				}

				totalOutlet, err := service.PJPDistributorRepository.CountOutletsByDate(context.Background(), dataFilter.EmpID, wd.WorkDate.Format(time.DateOnly))
				if err != nil {
					return nil, 0, 0, err
				}
				numberOfOutlet = totalOutlet

			} else {
				date, err := times.FormatDateWithZeroTime(wd.WorkDate.Format(time.DateOnly))
				if err != nil {
					return nil, 0, 0, errors.New("invalid date")
				}

				// Fetch dynamic route name
				rc, rn, err := service.PJPPrincipalRepository.GetRouteNameByIndexDay(context.Background(), dataFilter.CustID, dataFilter.EmpID, seqDay)
				if err == nil && rn != "" {
					routeCode = rc
					routeName = rn
				} else {
					routeName = fmt.Sprintf("Route %d", idx+1)
				}

				totalInfo, err := service.PJPPrincipalRepository.CountOutletsByDate(context.Background(), dataFilter.EmpID, date)
				if err != nil {
					return nil, 0, 0, err
				}
				numberOfOutlet = totalInfo.TotalOutlet
				numberOfDistributor = totalInfo.TotalDistributor
			}

			workDayData = append(workDayData, entity.WorkDayData{
				PerYear:             wd.PerYear,
				PerId:               wd.PerId,
				WeekId:              wd.WeekId,
				WorkDate:            wd.WorkDate,
				IsWork:              wd.IsWork,
				IsActive:            wd.IsActive,
				IsClosed:            wd.IsClosed,
				ClosedAt:            wd.ClosedAt,
				ClosedBy:            wd.ClosedBy,
				ClosedByName:        wd.ClosedByName,
				RouteCode:           routeCode,
				RouteName:           routeName,
				NumberOfOutlet:      numberOfOutlet,
				NumberOfDistributor: numberOfDistributor,
			})
		}

		closedBy := 0
		if week.ClosedBy > 0 {
			closedBy = week.ClosedBy
		}

		data = append(data, entity.WeekListResponse{
			PerYear:      week.PerYear,
			PerId:        week.PerId,
			WeekId:       week.WeekId,
			WeekStart:    week.WeekStart,
			WeekEnd:      week.WeekEnd,
			IsActive:     week.IsActive,
			IsClosed:     week.IsClosed,
			ClosedAt:     week.ClosedAt,
			ClosedBy:     closedBy,
			ClosedByName: week.ClosedByName,
			WorkDays:     workDayData,
		})
	}

	return data, total, lastPage, nil
}
