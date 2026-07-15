package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/pkg/structs"
	"mobile/repository"
	"time"
)

type ActivitiesService interface {
	ActivitiesSummaryDaily(entity.SummaryDailyRequest) (entity.SummaryDailyResponse, error)
}

type ActivitiesServiceImpl struct {
	Config         env.ConfigEnv
	Transaction    repository.Dbtransaction
	ActivitiesRepo repository.ActivitiesRepository
}

func NewActivitiesService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
	activitiesRepo repository.ActivitiesRepository,
) *ActivitiesServiceImpl {
	return &ActivitiesServiceImpl{
		Config:         config,
		Transaction:    transaction,
		ActivitiesRepo: activitiesRepo,
	}
}

func (service *ActivitiesServiceImpl) ActivitiesSummaryDaily(request entity.SummaryDailyRequest) (response entity.SummaryDailyResponse, err error) {
	employee, err := service.ActivitiesRepo.FindSummaryActivity(request)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(employee, &response)
	if err != nil {
		return response, err
	}

	if employee.StartTime != nil && employee.EndTime != nil {
		totalStart, err := time.Parse("15:04:05", *employee.StartTime)
		if err != nil {
			return response, err
		}
		totalEnd, err := time.Parse("15:04:05", *employee.EndTime)
		if err != nil {
			return response, err
		}
		totalSpent := totalEnd.Sub(totalStart)
		response.EstTime = totalSpent.String()
	}

	return response, err
}
