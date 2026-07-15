package service

import (
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type ManageMinimumPriceService interface {
	LookupBasePrice() (data []entity.BasePriceLookupResponse, err error)
	LookupLimitAction() (data []entity.LimitActionLookupResponse, err error)
	List(dataFilter entity.ManageMinimumPriceQueryFilter, custId, parentCustId string) (data []entity.ManageMinimumPriceRead, total int, lastPage int, err error)
	Detail(params entity.DetailManageMinimumPriceParams) (response entity.ManageMinimumPriceRead, err error)
	Store(request entity.BodyCreateManageMinimumPrice) (response entity.ManageMinimumPriceRead, err error)
	Update(manageMinimumPriceId int64, request entity.UpdateManageMinimumPrice) (err error)
	Delete(custId string, manageMinimumPriceId int64, userId int64) (err error)
	UpdateStatus(custId string, manageMinimumPriceId int64, status int, userId int64) (err error)
}

func NewManageMinimumPriceService(manageMinimumPriceRepository repository.ManageMinimumPriceRepository) *manageMinimumPriceServiceImpl {
	return &manageMinimumPriceServiceImpl{
		ManageMinimumPriceRepository: manageMinimumPriceRepository,
	}
}

type manageMinimumPriceServiceImpl struct {
	ManageMinimumPriceRepository repository.ManageMinimumPriceRepository
	// MProductRepository repository.MProductRepository
}

func (service *manageMinimumPriceServiceImpl) LookupBasePrice() (data []entity.BasePriceLookupResponse, err error) {

	data = entity.BasePrice

	return data, err
}

func (service *manageMinimumPriceServiceImpl) LookupLimitAction() (data []entity.LimitActionLookupResponse, err error) {

	data = entity.LimitAction

	return data, err
}

func (service *manageMinimumPriceServiceImpl) List(dataFilter entity.ManageMinimumPriceQueryFilter, custId, parentCustId string) (data []entity.ManageMinimumPriceRead, total int, lastPage int, err error) {
	manageMinimumPrice, total, lastPage, err := service.ManageMinimumPriceRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range manageMinimumPrice {
		var vResp entity.ManageMinimumPriceRead
		structs.Automapper(row, &vResp)

		vResp.BasePriceName = entity.GetBasePriceName(*vResp.BasePrice)
		vResp.LimitActionName = entity.GetLimitActionName(*vResp.LimitAction)
		vResp.StatusManageMinimumPriceName = entity.ConvStatusManageMinimumPrice(*vResp.StatusManageMinimumPrice)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *manageMinimumPriceServiceImpl) Detail(params entity.DetailManageMinimumPriceParams) (response entity.ManageMinimumPriceRead, err error) {
	manageMinimumPrice, err := service.ManageMinimumPriceRepository.FindDetailById(params.ManageMinimumPriceId, params.CustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(manageMinimumPrice, &response)
	if err != nil {
		return response, err
	}

	response.BasePriceName = entity.GetBasePriceName(*response.BasePrice)
	response.LimitActionName = entity.GetLimitActionName(*response.LimitAction)
	response.StatusManageMinimumPriceName = entity.ConvStatusManageMinimumPrice(*response.StatusManageMinimumPrice)

	return response, err
}

func (service *manageMinimumPriceServiceImpl) Store(request entity.BodyCreateManageMinimumPrice) (response entity.ManageMinimumPriceRead, err error) {
	service.ManageMinimumPriceRepository.TrxBegin()

	timeNow := time.Now().In(time.UTC)

	defer func() {
		if p := recover(); p != nil {
			service.ManageMinimumPriceRepository.TrxRollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	for _, detail := range request.Body {

		checkProId, _ := service.ManageMinimumPriceRepository.FindDetailByProId(int64(*detail.ProId), request.CustId)
		if checkProId.ManageMinimumPriceId != nil && *checkProId.ManageMinimumPriceId != 0 {
			service.ManageMinimumPriceRepository.Delete(request.CustId, int64(*checkProId.ManageMinimumPriceId), request.CreatedBy)
		}

		modelDetail := model.ManageMinimumPrice{
			CustId:    request.CustId,
			CreatedAt: timeNow,
			CreatedBy: &request.CreatedBy,
			UpdatedAt: timeNow,
			UpdatedBy: &request.CreatedBy,
		}

		err = structs.Automapper(detail, &modelDetail)
		if err != nil {
			return response, err
		}

		status := 1
		modelDetail.StatusManageMinimumPrice = &status

		_, err := service.ManageMinimumPriceRepository.Store(modelDetail)
		if err != nil {
			service.ManageMinimumPriceRepository.TrxRollback()
			return response, err
		}
	}

	err = service.ManageMinimumPriceRepository.TrxCommit()
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *manageMinimumPriceServiceImpl) Update(manageMinimumPriceId int64, request entity.UpdateManageMinimumPrice) (err error) {

	service.ManageMinimumPriceRepository.TrxBegin()

	defer func() {
		if p := recover(); p != nil {
			service.ManageMinimumPriceRepository.TrxRollback()
			panic(p) // re-throw panic after Rollback
		}
	}()

	err = service.ManageMinimumPriceRepository.Update(manageMinimumPriceId, request)
	if err != nil {
		service.ManageMinimumPriceRepository.TrxRollback()
		return err
	}

	service.ManageMinimumPriceRepository.TrxCommit()

	return err
}

func (service *manageMinimumPriceServiceImpl) Delete(custId string, manageMinimumPriceId int64, userId int64) (err error) {

	err = service.ManageMinimumPriceRepository.Delete(custId, manageMinimumPriceId, userId)
	if err != nil {
		return err
	}

	return err
}

func (service *manageMinimumPriceServiceImpl) UpdateStatus(custId string, manageMinimumPriceId int64, status int, userId int64) (err error) {

	err = service.ManageMinimumPriceRepository.UpdateStatus(manageMinimumPriceId, status, custId, userId)
	if err != nil {
		return err
	}

	return err
}
