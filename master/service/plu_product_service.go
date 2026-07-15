package service

import (
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type PluProductService interface {
	Detail(int, string) (entity.PluProductResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.PluProductResponse, total int, lastPage int, err error)
	Store(entity.CreatePluProductBody) (entity.PluProductResponse, error)
	Update(int, entity.UpdatePluProductRequest) error
	Delete(string, int, int64) error
}

func NewPluProductService(pluProductRepository repository.PluProductRepository) *pluProductServiceImpl {
	return &pluProductServiceImpl{
		PluProductRepository: pluProductRepository,
	}
}

type pluProductServiceImpl struct {
	PluProductRepository repository.PluProductRepository
}

func (service *pluProductServiceImpl) Detail(pluProId int, custId string) (response entity.PluProductResponse, err error) {
	pluProduct, err := service.PluProductRepository.FindOneByPluProductIdAndCustId(pluProId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(pluProduct, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *pluProductServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.PluProductResponse, total int, lastPage int, err error) {
	pluGroups, total, lastPage, err := service.PluProductRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(pluGroups) > 0 {
		for _, row := range pluGroups {
			var vResp entity.PluProductResponse
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *pluProductServiceImpl) Store(request entity.CreatePluProductBody) (response entity.PluProductResponse, err error) {
	// plu_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	// pluProduct, err := service.PluProductRepository.FindOneByPluProductGrpAndCustId(request.PluGrpId, request.CustId)
	// if err == nil {
	// 	return response, errors.New("plu_grp_id: " + strconv.Itoa(*pluProduct.PluGrpId) + " is already exists")
	// }

	timeNow := time.Now().In(time.UTC)
	pluProductData := model.PluProduct{
		CustId:    request.CustId,
		PluGrpId:  request.PluGrpId,
		ProId:     request.ProId,
		PluNo:     request.PluNo,
		CreatedAt: &timeNow,
		CreatedBy: &request.CreatedBy,
		UpdatedBy: &request.CreatedBy,
		UpdatedAt: &timeNow,
	}

	pluProductId, err := service.PluProductRepository.Store(pluProductData)
	if err != nil {
		return response, err
	}

	response.PluProId = pluProductId

	return response, err
}

func (service *pluProductServiceImpl) Update(pluProId int, request entity.UpdatePluProductRequest) (err error) {

	err = service.PluProductRepository.Update(pluProId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *pluProductServiceImpl) Delete(custId string, pluProId int, userId int64) (err error) {

	err = service.PluProductRepository.Delete(custId, pluProId, userId)
	if err != nil {
		return err
	}

	return err
}
