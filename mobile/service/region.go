package service

import (
	"context"
	"mobile/model"
	"mobile/repository"
)

type RegionService interface {
	List(ctx context.Context, custID string) ([]model.Region, error)
}

type regionService struct {
	Repo repository.RegionRepository
}

func NewRegionService(repo repository.RegionRepository) RegionService {
	return &regionService{
		Repo: repo,
	}
}

func (s *regionService) List(ctx context.Context, custID string) ([]model.Region, error) {
	records, err := s.Repo.List(ctx, custID)
	if err != nil {
		return nil, err
	}

	return records, nil
}
