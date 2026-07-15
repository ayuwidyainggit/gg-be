package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/repository"
	"time"
)

type EventsService interface {
	Events(entity.EventsRequest) ([]entity.EventsResponse, error)
}

type EventsServiceImpl struct {
	Config env.ConfigEnv
	// MCustomerRepository repository.MCustomerRepository,
	Transaction repository.Dbtransaction
}

func NewEventsService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
) *EventsServiceImpl {
	return &EventsServiceImpl{
		Config:      config,
		Transaction: transaction,
	}
}

func (service *EventsServiceImpl) Events(request entity.EventsRequest) (response []entity.EventsResponse, err error) {

	for i := 0; i < 2; i++ {
		rowResp := entity.EventDet{
			Name: "Google meet",
			Link: "link",
			Icon: "Icon",
		}

		rowEvent := entity.EventsResponse{
			Id:        "id",
			Type:      "type",
			Name:      "name",
			StartTime: time.Now().Format("2006-01-02 15:04:05"),
			EndTime:   time.Now().Format("2006-01-02 15:04:05"),
			Location:  rowResp,
		}
		response = append(response, rowEvent)

	}

	return response, err
}
