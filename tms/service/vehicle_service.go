package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"math"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/repository"
)

type VehicleService interface {
	GetVehicle(ctx context.Context, dataFilter entity.VehicleQueryFilter) ([]entity.VehicleResponse, entity.Meta, error)
}

type VehicleServiceImpl struct {
	vehicleRepo repository.VehicleRepo
	validate    *validator.Validate
}

func NewVehicleServiceImpl(vehicleRepo repository.VehicleRepo, validate *validator.Validate) VehicleService {
	return &VehicleServiceImpl{
		vehicleRepo: vehicleRepo,
		validate:    validate,
	}
}

func (service *VehicleServiceImpl) GetVehicle(ctx context.Context, dataFilter entity.VehicleQueryFilter) ([]entity.VehicleResponse, entity.Meta, error) {
	result, total := service.vehicleRepo.GetVehicle(ctx, dataFilter)
	var data []entity.VehicleResponse

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	if dataFilter.Page == 0 {
		dataFilter.Page = 1
	}

	for _, value := range result {
		var res entity.VehicleResponse
		helper.Automapper(value, &res)

		data = append(data, res)
	}

	pagination := &entity.Meta{
		TotalData: total,
		Page:      dataFilter.Page,
		Limit:     dataFilter.Limit,
		TotalPage: int(math.Ceil(float64(total) / float64(dataFilter.Limit))),
	}

	return data, *pagination, nil
}
