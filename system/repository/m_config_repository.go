package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"system/entity"
	"system/model"

	"gorm.io/gorm"
)

type (
	RepositoryMConfigImpl struct {
		*gorm.DB
	}
)
type MConfigRepository interface {
	Store(c context.Context, data *model.MConfig) error
	FindDetail(configID string, custId string) (Details model.MConfig, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MConfig, int64, int, error)
	FindAllByCustIdDetails(dataFilter entity.MConfigQueryFilter, custId string) ([]model.MConfig, error)
	Delete(c context.Context, custId string, configID string) error
	Update(c context.Context, ConfigId string, data model.MConfig) error
}

func NewMConfigRepository(db *gorm.DB) *RepositoryMConfigImpl {
	return &RepositoryMConfigImpl{db}
}

func (repo *RepositoryMConfigImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMConfigImpl) Store(c context.Context, data *model.MConfig) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryMConfigImpl) FindDetail(configID string, custId string) (Details model.MConfig, err error) {
	err = repository.
		Where("config_id = ? AND cust_id=?", configID, custId).
		Take(&Details).Error
	return Details, err
}

func (repository *RepositoryMConfigImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MConfig, int64, int, error) {
	var data []model.MConfig
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("config_id")
	query := repository.Select("*")

	if dataFilter.Query != "" {

	}
	if dataFilter.Sort != "" {

	} else {
		query.Order("config_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return data, total, 0, err
	}
	err = queryCount.Model(&data).Count(&total).Error
	if err != nil {
		return data, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return data, total, lastPage, nil
}

func (repository *RepositoryMConfigImpl) FindAllByCustIdDetails(dataFilter entity.MConfigQueryFilter, custId string) ([]model.MConfig, error) {
	var data []model.MConfig
	qWhere := `cust_id = '` + custId + `' AND `
	tempConfigid := []interface{}{}
	for _, xx := range dataFilter.ConfigId {
		tempConfigid = append(tempConfigid, "'"+xx+"'")
	}
	varConfigid := make([]string, len(tempConfigid))
	for i, v := range tempConfigid {
		varConfigid[i] = fmt.Sprint(v)
	}
	var Configid = strings.Join(varConfigid, ",")
	// fmt.Println("1. FILTER >>>>", Configid)

	if len(dataFilter.ConfigId) > 0 {
		qWhere += ` config_id in (` + Configid + `) `
	}

	query := repository.Select("*").Where(qWhere)
	if dataFilter.Query != "" {

	}

	if dataFilter.Sort != "" {

	} else {
		// query.Order("config_id DESC")
	}

	err := query.Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (repository *RepositoryMConfigImpl) Delete(c context.Context, custId string, configID string) error {
	var data model.MConfig
	result := repository.model(c).Delete(&data, "config_id=? AND cust_id = ?", configID, custId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryMConfigImpl) Update(c context.Context, ConfigId string, data model.MConfig) error {
	result := repository.model(c).Model(&data).Where("config_id=?", ConfigId).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
