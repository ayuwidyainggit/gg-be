package repository

import (
	"context"
	"math"
	"system/entity"
	"system/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type (
	RepositoryMDayImpl struct {
		*gorm.DB
	}
)

type MDayRepository interface {
	FindAllByLangId(dataFilter entity.GeneralQueryFilter, LangId string) ([]model.MDay, int64, int, error)
	FindDetail(dayId int64, langId string) (Details model.MDay, err error)
}

func NewMDayRepository(db *gorm.DB) *RepositoryMDayImpl {
	return &RepositoryMDayImpl{db}
}

func (repo *RepositoryMDayImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMDayImpl) FindAllByLangId(dataFilter entity.GeneralQueryFilter, LangId string) ([]model.MDay, int64, int, error) {
	var mdays []model.MDay
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("lang_id")
	query := repository.Select("*")

	queryCount.Where("lang_id=?", LangId)
	query.Where("lang_id=?", LangId)

	if dataFilter.Query != "" {
		query.Where("day_name=?", dataFilter.Query)
		queryCount.Where("day_name=?", dataFilter.Query)
	}

	if dataFilter.Sort != "" {

	} else {
		// query.Order("bppr_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&mdays).Error
	if err != nil {
		return mdays, total, 0, err
	}

	err = queryCount.Model(&mdays).Count(&total).Error
	if err != nil {
		log.Error("queryCount, err:", err.Error())
		return mdays, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))

	return mdays, total, lastPage, nil
}

func (repository *RepositoryMDayImpl) FindDetail(dayId int64, langId string) (Details model.MDay, err error) {
	err = repository.
		Where("day_id = ? and lang_id = ?", dayId, langId).
		Take(&Details).Error
	return Details, err
}
