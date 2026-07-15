package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryDistributorImpl struct {
		*gorm.DB
	}
)

type DistributorRepository interface {
	FindMobileDistributorList(dataFilter entity.MobileDistributorListQueryFilter, custId string) (distributors []model.MobileDistributorList, total int64, lastPage int, err error)
	GetDistributorByID(ctx context.Context, id int64) (*model.Distributor, error)
	GetPrincipalInfo(ctx context.Context, custId string) (*model.PrincipalInfo, error)
}

func NewDistributorRepository(db *gorm.DB) *RepositoryDistributorImpl {
	return &RepositoryDistributorImpl{db}
}

func (repo *RepositoryDistributorImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryDistributorImpl) GetDistributorByID(ctx context.Context, id int64) (*model.Distributor, error) {
	var result model.Distributor

	err := repository.WithContext(ctx).
		Table("mst.m_distributor").
		Where("distributor_id = ? AND is_del = false", id).
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (repository *RepositoryDistributorImpl) GetPrincipalInfo(ctx context.Context, custId string) (*model.PrincipalInfo, error) {
	var result model.PrincipalInfo

	err := repository.model(ctx).
		Table("smc.m_customer").
		Select(`cust_id, cust_name, COALESCE(distributor_id, 0) as distributor_id`).
		Where("cust_id = ? AND is_active = true", custId).
		Limit(1).
		Take(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (repository *RepositoryDistributorImpl) FindMobileDistributorList(dataFilter entity.MobileDistributorListQueryFilter, custId string) ([]model.MobileDistributorList, int64, int, error) {
	distributors := []model.MobileDistributorList{}

	query := repository.DB.
		Select(`
			md.distributor_id,
			COALESCE(md.distributor_code, '') AS distributor_code,
			COALESCE(md.distributor_name, '') AS distributor_name,
			COALESCE(md.address, '') AS address,
			COALESCE(md.latitude, '') AS latitude,
			COALESCE(md.longitude, '') AS longitude,
			COALESCE(md.region_id, 0) AS region_id,
			COALESCE(md.area_id, 0) AS area_id
		`).
		Table("mst.m_distributor AS md").
		Where("md.is_del=false").
		Where("md.parent_cust_id LIKE ?", custId+"%")

	if dataFilter.Query != "" {
		query = query.Where("(md.distributor_code ILIKE ? OR md.distributor_name ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			query = query.Where("md.is_active = true")
		} else if *dataFilter.IsActive == 0 {
			query = query.Where("md.is_active = false")
		}
	}

	if len(dataFilter.RegionID) > 0 {
		query = query.Where("md.region_id IN (?)", dataFilter.RegionID)
	}

	if len(dataFilter.AreaID) > 0 {
		query = query.Where("md.area_id IN (?)", dataFilter.AreaID)
	}

	// Count total
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return distributors, 0, 0, err
	}

	// Sorting
	sortBy := "md.distributor_code ASC" // default sort
	if dataFilter.Sort != "" {
		sortParts := strings.Split(dataFilter.Sort, ":")
		if len(sortParts) == 2 {
			sortField := sortParts[0]
			sortOrder := strings.ToUpper(sortParts[1])
			if sortOrder == "ASC" || sortOrder == "DESC" {
				if sortField == "created_date" {
					sortField = "md.created_date"
				} else {
					sortField = fmt.Sprintf("md.%s", sortField)
				}
				sortBy = fmt.Sprintf("%s %s", sortField, sortOrder)
			}
		}
	}
	query = query.Order(sortBy)

	// Pagination
	if dataFilter.Limit <= 0 {
		dataFilter.Limit = 5
	}
	page := dataFilter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))
	if lastPage <= 0 {
		lastPage = 1
	}

	err = query.Offset(offset).Limit(dataFilter.Limit).Find(&distributors).Error
	if err != nil {
		return distributors, 0, 0, err
	}

	return distributors, total, lastPage, nil
}
