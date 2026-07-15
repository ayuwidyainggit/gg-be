package repository

import (
	"context"
	"cronjob/entity"
	"cronjob/model"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	RepositoryJobImpl struct {
		*gorm.DB
	}
)
type JobRepository interface {
	Store(c context.Context, data *model.Job) error
	FindOneByJobID(jobID uuid.UUID) (model.Job, error)
	Update(c context.Context, jobID uuid.UUID, data model.Job) error
	Delete(c context.Context, jobID uuid.UUID) error
	Inactive(c context.Context, jobID uuid.UUID) error
	FindAll(dataFilter entity.GeneralQueryFilter) ([]model.Job, int64, int, error)
}

func NewJobRepository(db *gorm.DB) *RepositoryJobImpl {
	return &RepositoryJobImpl{db}
}

func (repo *RepositoryJobImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryJobImpl) FindOneByJobID(jobID uuid.UUID) (model.Job, error) {
	job := model.Job{}

	err := repository.
		Where("job_id = ?", jobID).
		Take(&job).Error
	if err != nil {
		log.Errorf("%+v:", err)
		return job, err
	}

	return job, nil
}

func (repository *RepositoryJobImpl) Store(c context.Context, data *model.Job) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		log.Errorf("%+v:", err)
		return err
	}
	return nil
}

func (repository *RepositoryJobImpl) Update(c context.Context, jobID uuid.UUID, data model.Job) error {
	result := repository.model(c).Model(&data).Where("job_id = ?", jobID).Updates(&data)
	if result.Error != nil {
		log.Errorf("%+v:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryJobImpl) Delete(c context.Context, jobID uuid.UUID) error {
	var data model.Job
	result := repository.model(c).Model(&data).Delete("job_id = ?", jobID)
	if result.Error != nil {
		log.Errorf("%+v:", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryJobImpl) FindAll(dataFilter entity.GeneralQueryFilter) ([]model.Job, int64, int, error) {
	var data []model.Job
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("job_id")
	query := repository.Select("*")

	if dataFilter.Active == 1 {
		queryCount = queryCount.Where("active = ?", true)
		query = query.Where("active = ?", true)
	} else if dataFilter.Active == 2 {
		queryCount = queryCount.Where("active = ?", false)
		query = query.Where("active = ?", false)
	}

	if dataFilter.Query != "" {

	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("job_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		log.Errorf("%+v:", err)
		return data, total, 0, err
	}

	if dataFilter.Count {
		err = queryCount.Model(&data).Count(&total).Error
		if err != nil {
			log.Errorf("%+v:", err)
			return data, total, 0, err
		}
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return data, total, lastPage, nil
}

func (repository *RepositoryJobImpl) Inactive(c context.Context, jobID uuid.UUID) error {
	return repository.Model(&model.Job{}).
		Where("job_id = ?", jobID).
		Update("active", false).Error
}
