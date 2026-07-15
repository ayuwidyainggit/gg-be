package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/repository"
	"time"
)

type AnnouncementsService interface {
	Announcements(entity.AnnouncementsRequest) ([]entity.AnnouncementsResponse, error)
}

type AnnouncementsServiceImpl struct {
	Config env.ConfigEnv
	// MCustomerRepository repository.MCustomerRepository,
	Transaction repository.Dbtransaction
}

func NewAnnouncementsService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
) *AnnouncementsServiceImpl {
	return &AnnouncementsServiceImpl{
		Config:      config,
		Transaction: transaction,
	}
}

func (service *AnnouncementsServiceImpl) Announcements(request entity.AnnouncementsRequest) (response []entity.AnnouncementsResponse, err error) {

	for i := 1; i < 3; i++ {

		rowAnnouncements := entity.AnnouncementsResponse{
			Id:        i,
			Name:      "name",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			IsNew:     false,
			Image:     "image",
		}
		response = append(response, rowAnnouncements)

	}

	return response, err
}
