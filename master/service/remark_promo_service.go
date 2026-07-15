package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type RemarkPromoService interface {
	Detail(int, string) (entity.RemarkPromoResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.RemarkPromoResponse, total int, lastPage int, err error)
	Store(entity.CreateRemarkPromoBody) (entity.RemarkPromoResponse, error)
	Update(int, entity.UpdateRemarkPromoRequest) error
	Delete(string, int, int64) error
}

func NewRemarkPromoService(remarkPromoRepository repository.RemarkPromoRepository) *remarkPromoServiceImpl {
	return &remarkPromoServiceImpl{
		RemarkPromoRepository: remarkPromoRepository,
	}
}

type remarkPromoServiceImpl struct {
	RemarkPromoRepository repository.RemarkPromoRepository
}

func (service *remarkPromoServiceImpl) Detail(remPromoId int, custId string) (response entity.RemarkPromoResponse, err error) {
	remarkPromo, err := service.RemarkPromoRepository.FindOneByRemarkPromoIdAndCustId(remPromoId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(remarkPromo, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *remarkPromoServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.RemarkPromoResponse, total int, lastPage int, err error) {
	remarkPromos, total, lastPage, err := service.RemarkPromoRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(remarkPromos) > 0 {
		for _, row := range remarkPromos {
			var vResp entity.RemarkPromoResponse
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *remarkPromoServiceImpl) Store(request entity.CreateRemarkPromoBody) (response entity.RemarkPromoResponse, err error) {
	remarkPromo, err := service.RemarkPromoRepository.FindOneByRemarkPromoCodeAndCustId(request.RemPromoCode, request.CustId)
	if err == nil {
		return response, errors.New("remark_promo_code: " + remarkPromo.RemPromoCode + " is already exists")
	}

	var remarkPromoData model.RemarkPromo

	timeNow := time.Now().In(time.UTC)
	structs.Automapper(request, &remarkPromoData)
	remarkPromoData.CreatedAt = &timeNow
	remarkPromoData.CreatedBy = &request.CreatedBy
	remarkPromoData.UpdatedBy = &request.UpdatedBy
	remarkPromoData.UpdatedAt = &timeNow
	RemPromoId, err := service.RemarkPromoRepository.Store(remarkPromoData)
	if err != nil {
		return response, err
	}

	response.RemPromoId = RemPromoId

	return response, err
}

func (service *remarkPromoServiceImpl) Update(remarkPromoId int, request entity.UpdateRemarkPromoRequest) (err error) {
	// remarkPromo_code & cust id validation, if err == nil and params remarkPromoId != remarkPromo.Id, this means that code & cust id already exists
	remarkPromo, err := service.RemarkPromoRepository.FindOneByRemarkPromoCodeAndCustId(request.RemPromoCode, request.CustId)
	if err == nil && remarkPromo.RemPromoId != remarkPromoId {
		return errors.New("remark_promo_code: " + remarkPromo.RemPromoCode + " is already exists")
	}
	err = service.RemarkPromoRepository.Update(remarkPromoId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *remarkPromoServiceImpl) Delete(custId string, remarkPromoId int, userId int64) (err error) {

	err = service.RemarkPromoRepository.Delete(custId, remarkPromoId, userId)
	if err != nil {
		return err
	}
	return err
}
