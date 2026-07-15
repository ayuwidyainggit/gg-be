package service

import (
	"context"
	"mobile/model"
	"mobile/repository"
)

type AreaService interface {
	List(ctx context.Context, custID string) ([]model.Area, error)
}

type areaService struct {
	Repo repository.AreaRepository
}

func NewAreaService(repo repository.AreaRepository) AreaService {
	return &areaService{
		Repo: repo,
	}
}

func (s *areaService) List(ctx context.Context, custID string) ([]model.Area, error) {
	records, err := s.Repo.List(ctx, custID)
	if err != nil {
		return nil, err
	}

	return records, nil
}
