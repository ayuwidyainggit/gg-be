package service

import (
	"context"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/structs"
	"mobile/repository"
)

type MConfigService interface {
	Store(request entity.CreateMConfigBody) (err error)
	Detail(configID string, custID string) (response entity.MConfigResponse, err error)
	List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MConfigResponse, total int64, lastPage int, err error)
	ListDetail(dataFilter entity.MConfigQueryFilter, custId string) (data []entity.MConfigResponse, err error)
	Delete(custId string, configID string) (err error)
	Update(ConfigId string, request entity.UpdateMConfigBody) (err error)
}

type mconfigServiceImpl struct {
	MConfigRepository repository.MConfigRepository
	Transaction       repository.Dbtransaction
}

func NewMConfigService(mconfigRepository repository.MConfigRepository, transaction repository.Dbtransaction) *mconfigServiceImpl {
	return &mconfigServiceImpl{
		MConfigRepository: mconfigRepository,
		Transaction:       transaction,
	}
}

func (service *mconfigServiceImpl) Store(request entity.CreateMConfigBody) (err error) {
	c := context.Background()
	var model model.MConfig

	err = structs.Automapper(request, &model)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.MConfigRepository.Store(txCtx, &model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *mconfigServiceImpl) Detail(configID string, custID string) (response entity.MConfigResponse, err error) {
	user, err := service.MConfigRepository.FindDetail(configID, custID)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(user, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (service *mconfigServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MConfigResponse, total int64, lastPage int, err error) {
	configs, total, lastPage, err := service.MConfigRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range configs {
		var vResp entity.MConfigResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mconfigServiceImpl) ListDetail(dataFilter entity.MConfigQueryFilter, custId string) (data []entity.MConfigResponse, err error) {
	configs, err := service.MConfigRepository.FindAllByCustIdDetails(dataFilter, custId)
	if err != nil {
		return data, err
	}

	for _, configId := range dataFilter.ConfigId {
		for _, row := range configs {
			if configId == row.ConfigID {
				var vResp entity.MConfigResponse
				structs.Automapper(row, &vResp)
				data = append(data, vResp)
			}
		}
	}

	return data, err
}

func (service *mconfigServiceImpl) Delete(custId string, configID string) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MConfigRepository.Delete(txCtx, custId, configID)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *mconfigServiceImpl) Update(ConfigId string, request entity.UpdateMConfigBody) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MConfig
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.MConfigRepository.Update(txCtx, ConfigId, Model)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
