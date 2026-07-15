package pjpenhance

import (
	"context"
	"scyllax-pjp/data/request"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

func (service *pjpEnhanceService) UpdatePjp(ctx context.Context, id int, request request.CreatePjpEnhanceRequest, currentCustomerId string) {
	err := service.validate.Struct(request)
	helper.ErrorPanic(err)

	tx := service.db.Begin()
	if tx.Error != nil {
		helper.ErrorPanic(tx.Error)
	}
	defer helper.CommitOrRollback(tx)

	findPjp := service.pjpRepository.GetPjpById(ctx, tx, id, currentCustomerId)
	existingRoutes := service.routeRepository.FindByPjpID(ctx, tx, findPjp.ID, currentCustomerId)

	pjp := buildPjpModel(request, currentCustomerId)
	pjp.ID = findPjp.ID
	service.pjpRepository.Update(ctx, tx, pjp)

	service.routeRepository.DeleteByPjpId(ctx, tx, pjp.ID, currentCustomerId)

	savedRoutes := service.createRoutes(ctx, tx, pjp, request.Routes, currentCustomerId, existingRoutes)
	editableVisitDays := filterEditableVisitDays(request.VisitDay, currentJakartaTime())
	service.deleteEditableVisitHistory(ctx, tx, pjp, savedRoutes, editableVisitDays, currentCustomerId)
	routePopPermanents := service.createVisitHistory(ctx, tx, pjp, savedRoutes, editableVisitDays, currentCustomerId)

	if routePopPermanents != nil {
		service.routePopRepository.CreateBulk(ctx, tx, routePopPermanents)
	}
}

func filterEditableVisitDays(visitDays []request.VisitDayCreatePjp, now time.Time) []request.VisitDayCreatePjp {
	var editable []request.VisitDayCreatePjp
	for _, visit := range visitDays {
		if visit.Date == "-" {
			editable = append(editable, visit)
			continue
		}
		dateTime, err := parseVisitDate(visit.Date)
		if err != nil {
			helper.ErrorPanic(err)
		}
		if !helper.IsBeforeCurrentWeek(dateTime, now) {
			editable = append(editable, visit)
		}
	}
	return editable
}

func (service *pjpEnhanceService) deleteEditableVisitHistory(ctx context.Context, tx *gorm.DB, savedPjp model.Pjp, savedRoutes []model.Route, visitDays []request.VisitDayCreatePjp, currentCustomerId string) {
	for _, visit := range visitDays {
		if visit.Date == "-" {
			continue
		}
		dateTime, err := parseVisitDate(visit.Date)
		helper.ErrorPanic(err)

		routeCode := getRouteCodeByName(savedRoutes, visit.Visit.RouteName)
		service.routeOutletHistoryRepository.DeleteByVisitDay(ctx, tx, model.RouteOutletHistory{
			PjpID:     &savedPjp.ID,
			CustID:    currentCustomerId,
			RouteCode: routeCode,
			Week:      visit.Week,
			Year:      visit.Year,
			Date:      dateTime,
		})
		service.routePopRepository.DeleteByVisitDay(ctx, tx, model.RoutePopPermanent{
			PjpID:     &savedPjp.ID,
			CustID:    currentCustomerId,
			RouteCode: &routeCode,
			Week:      visit.Week,
			Year:      visit.Year,
			Date:      dateTime,
		})
	}
}

func parseVisitDate(date string) (time.Time, error) {
	if date == "-" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", strings.Split(date, "T")[0])
}

func currentJakartaTime() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	return time.Now().In(loc)
}
