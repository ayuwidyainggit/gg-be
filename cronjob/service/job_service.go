package service

import (
	"context"
	"cronjob/entity"
	"cronjob/model"
	"cronjob/pkg/config/env"
	"cronjob/pkg/structs"
	"cronjob/repository"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

type JobService interface {
	Store(request entity.CreateJobBody) (err error)
	Update(jobID uuid.UUID, request entity.UpdateJobBody) (err error)
	Inactive(jobID uuid.UUID) (err error)
	Detail(jobID uuid.UUID) (response entity.JobList, err error)
	Delete(jobID uuid.UUID) (err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.JobList, total int64, lastPage int, err error)
}

type jobServiceImpl struct {
	JobRepository repository.JobRepository
	Transaction   repository.Dbtransaction
	Config        env.ConfigEnv
}

func NewJobService(
	config env.ConfigEnv,
	jobRepository repository.JobRepository,
	transaction repository.Dbtransaction,
) *jobServiceImpl {
	return &jobServiceImpl{
		Config:        config,
		JobRepository: jobRepository,
		Transaction:   transaction,
	}
}

func (service *jobServiceImpl) Store(request entity.CreateJobBody) (err error) {
	c := context.Background()
	var jobModel model.Job

	if err = structs.Automapper(request, &jobModel); err != nil {
		log.Errorf("%v:", err)
		return err
	}

	if len(request.RunAt) > 0 {
		runAtTz, err := time.Parse(time.RFC3339, request.RunAt)
		if err != nil {
			log.Errorf("%v:", err)
			return err
		}
		log.Info("runAtTz:", runAtTz.Format(time.RFC3339))
		runAtUTC := runAtTz.In(time.UTC).Format(time.RFC3339)
		log.Info("runAtUTC:", runAtUTC)
		jobModel.RunAt = &runAtUTC
	}

	if len(request.TimeOfDay) > 0 {
		timeOfDayTz, err := time.Parse(time.TimeOnly, request.TimeOfDay)
		if err != nil {
			log.Errorf("%v:", err)
			return err
		}
		// log.Info("timeOfDayTz:", timeOfDayTz.Format(time.RFC3339))
		timeOfDayUTC := timeOfDayTz.Format(time.TimeOnly)
		// log.Info("timeOfDayUTC:", timeOfDayUTC)
		jobModel.TimeOfDay = &timeOfDayUTC
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// log.Infof("jobModel: %+v", jobModel)
		err := service.JobRepository.Store(txCtx, &jobModel)
		if err != nil {
			log.Errorf("%+v:", err)
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *jobServiceImpl) Detail(jobID uuid.UUID) (response entity.JobList, err error) {
	job, err := service.JobRepository.FindOneByJobID(jobID)
	if err != nil {
		log.Errorf("%+v:", err)
		return response, err
	}
	if err = structs.Automapper(job, &response); err != nil {
		log.Errorf("%+v:", err)
		return response, err
	}
	return response, nil
}

func (service *jobServiceImpl) Delete(jobID uuid.UUID) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.JobRepository.Delete(txCtx, jobID)
		if err != nil {
			log.Errorf("%+v:", err)
			return err
		}
		return nil
	})

	return err
}

func (service *jobServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.JobList, total int64, lastPage int, err error) {
	jobs, total, lastPage, err := service.JobRepository.FindAll(dataFilter)
	if err != nil {
		log.Errorf("%+v:", err)
		return data, total, lastPage, err
	}

	for _, row := range jobs {
		var vResp entity.JobList
		if err = structs.Automapper(row, &vResp); err != nil {
			log.Errorf("%+v:", err)
			return data, total, lastPage, err
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *jobServiceImpl) Update(jobID uuid.UUID, request entity.UpdateJobBody) (err error) {
	c := context.Background()

	var jobModel model.Job
	_, err = service.JobRepository.FindOneByJobID(jobID)
	if err != nil {
		log.Errorf("%+v:", err)
		return err
	}
	if err = structs.Automapper(request, &jobModel); err != nil {
		log.Errorf("%+v:", err)
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if err := service.JobRepository.Update(txCtx, jobID, jobModel); err != nil {
			log.Errorf("%+v:", err)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *jobServiceImpl) Inactive(jobID uuid.UUID) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.JobRepository.Inactive(txCtx, jobID)
		if err != nil {
			log.Errorf("%+v:", err)
			return err
		}
		return nil
	})

	return err
}
