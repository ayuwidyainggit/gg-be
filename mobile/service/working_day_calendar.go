package service

import (
	"errors"
	"mobile/entity"
	"mobile/repository"

	"github.com/gofiber/fiber/v2/log"
)

type (
	WorkingDayCalendarService interface {
		List(dataFilter entity.WorkingDayCalendarQueryFilter) (data []entity.WorkingDayCalendarResponse, err error)
		ListMonths(dataFilter entity.WorkingDayCalendarMonthQueryFilter) (data []entity.WorkingDayCalendarMonthResponse, err error)
	}

	workingDayCalendarServiceImpl struct {
		WorkingDayCalendarRepository repository.WorkingDayCalendarRepository
	}
)

func NewWorkingDayCalendarService(
	workingDayCalendarRepository repository.WorkingDayCalendarRepository,
) WorkingDayCalendarService {
	return &workingDayCalendarServiceImpl{
		WorkingDayCalendarRepository: workingDayCalendarRepository,
	}
}

func (service *workingDayCalendarServiceImpl) List(dataFilter entity.WorkingDayCalendarQueryFilter) (data []entity.WorkingDayCalendarResponse, err error) {
	calendars, err := service.WorkingDayCalendarRepository.FindAll(dataFilter)
	if err != nil {
		log.Error("WorkingDayCalendarService, List, FindAll, err:", err.Error())
		return data, errors.New("record not found")
	}

	data = make([]entity.WorkingDayCalendarResponse, 0)
	for _, cal := range calendars {
		data = append(data, entity.WorkingDayCalendarResponse{
			WorkingDayCalendarID: cal.WorkingDayCalendarID,
			CustID:               cal.CustID,
			Title:                cal.Title,
			StartDate:            cal.StartDate,
			NumberOfWeeks:        cal.NumberOfWeeks,
			EndDate:              cal.EndDate,
			DefaultHolidays:      cal.DefaultHolidays,
			IsClosed:             cal.IsClosed,
			IsActive:             cal.IsActive,
		})
	}

	return data, nil
}

func (service *workingDayCalendarServiceImpl) ListMonths(dataFilter entity.WorkingDayCalendarMonthQueryFilter) (data []entity.WorkingDayCalendarMonthResponse, err error) {
	months, err := service.WorkingDayCalendarRepository.FindMonthsByWDCID(dataFilter.WDCID)
	if err != nil {
		log.Error("WorkingDayCalendarService, ListMonths, FindMonthsByWDCID, err:", err.Error())
		return data, errors.New("record not found")
	}

	data = make([]entity.WorkingDayCalendarMonthResponse, 0)
	for _, m := range months {
		data = append(data, entity.WorkingDayCalendarMonthResponse{
			IsActive:  m.IsActive,
			Month:     m.Month,
			Year:      m.Year,
			TextMonth: m.TextMonth,
		})
	}

	return data, nil
}
