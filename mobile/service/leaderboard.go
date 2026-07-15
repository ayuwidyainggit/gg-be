package service

import (
	"mobile/entity"
	"mobile/pkg/config/env"
	"mobile/repository"
)

type LeaderboardsService interface {
	Leaderboards(entity.LeaderboardsRequest) ([]entity.LeaderboardsResponse, error)
}

type LeaderboardsServiceImpl struct {
	Config env.ConfigEnv
	// MCustomerRepository repository.MCustomerRepository,
	Transaction repository.Dbtransaction
}

func NewLeaderboardsService(
	config env.ConfigEnv,
	transaction repository.Dbtransaction,
) *LeaderboardsServiceImpl {
	return &LeaderboardsServiceImpl{
		Config:      config,
		Transaction: transaction,
	}
}

func (service *LeaderboardsServiceImpl) Leaderboards(request entity.LeaderboardsRequest) (response []entity.LeaderboardsResponse, err error) {

	for i := 1; i < 3; i++ {

		rowLeaderboards := entity.LeaderboardsResponse{
			Rank:     i,
			Name:     "name",
			Currency: "IDR",
			Amount:   1234,
			Unit:     "M",
			Image:    "image",
			IsMe:     false,
		}
		response = append(response, rowLeaderboards)

	}

	return response, err
}
