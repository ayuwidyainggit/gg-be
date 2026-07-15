package service

import (
	"context"
	"errors"
	"mobile/entity"
	"mobile/model"
	"mobile/repository"

	"github.com/gofiber/fiber/v2/log"
)

type (
	DistributorService interface {
		MobileDistributorList(ctx context.Context, dataFilter entity.MobileDistributorListQueryFilter, custId string) (data []entity.MobileDistributorListResponse, total int64, lastPage int, err error)
		GetPrincipalInfo(ctx context.Context, custId string) (*model.PrincipalInfo, error)
	}

	distributorServiceImpl struct {
		DistributorRepository repository.DistributorRepository
	}
)

func NewDistributorService(
	distributorRepository repository.DistributorRepository,
) DistributorService {
	return &distributorServiceImpl{
		DistributorRepository: distributorRepository,
	}
}

func (service *distributorServiceImpl) MobileDistributorList(ctx context.Context, dataFilter entity.MobileDistributorListQueryFilter, custId string) (data []entity.MobileDistributorListResponse, total int64, lastPage int, err error) {
	if dataFilter.IsActive == nil {
		isActive := 1 // tue
		dataFilter.IsActive = &isActive
	}

	if dataFilter.Page <= 0 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "distributor_code:asc"
	}

	distributors, total, lastPage, err := service.DistributorRepository.FindMobileDistributorList(dataFilter, custId)
	if err != nil {
		log.Error("DistributorService, MobileDistributorList, FindMobileDistributorList, err:", err.Error())
		return data, 0, 0, errors.New("record not found")
	}

	data = make([]entity.MobileDistributorListResponse, 0)
	for _, dist := range distributors {
		data = append(data, entity.MobileDistributorListResponse{
			DistributorId:   dist.DistributorId,
			DistributorCode: dist.DistributorCode,
			DistributorName: dist.DistributorName,
			Address:         dist.Address,
			Latitude:        dist.Latitude,
			Longitude:       dist.Longitude,
			RegionID:        dist.RegionID,
			AreaID:          dist.AreaID,
		})
	}

	return data, total, lastPage, nil
}

func (service *distributorServiceImpl) GetPrincipalInfo(ctx context.Context, custId string) (*model.PrincipalInfo, error) {
	principal, err := service.DistributorRepository.GetPrincipalInfo(ctx, custId)
	if err != nil {
		return nil, err
	}

	return principal, nil
}
