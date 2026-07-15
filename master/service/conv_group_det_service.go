package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
	"time"
)

type ConvGroupDetService interface {
	Detail(int, int, string) (entity.ConvGroupDetResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.ConvGroupDetResponse, total int, lastPage int, err error)
	Store(entity.CreateConvGroupDetBody) (entity.ConvGroupDetResponse, error)
	Update(int, int, entity.UpdateConvGroupDetRequest) error
	Delete(string, int, int, int64) error
}

func NewConvGroupDetService(convGroupDetRepository repository.ConvGroupDetRepository) *convGroupDetServiceImpl {
	return &convGroupDetServiceImpl{
		ConvGroupDetRepository: convGroupDetRepository,
	}
}

type convGroupDetServiceImpl struct {
	ConvGroupDetRepository repository.ConvGroupDetRepository
}

func (service *convGroupDetServiceImpl) Detail(convGrpId int, proId int, custId string) (response entity.ConvGroupDetResponse, err error) {
	convGroup, err := service.ConvGroupDetRepository.FindOneByConvGroupIdAndCustId(convGrpId, custId, proId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(convGroup, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *convGroupDetServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.ConvGroupDetResponse, total int, lastPage int, err error) {
	convGroupsDet, total, lastPage, err := service.ConvGroupDetRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range convGroupsDet {
		var vResp entity.ConvGroupDetResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *convGroupDetServiceImpl) Store(request entity.CreateConvGroupDetBody) (response entity.ConvGroupDetResponse, err error) {

	proId, err := service.ConvGroupDetRepository.FindOneByProIdAndCustId(request.ConvGrpId, request.ProId, request.CustId)
	if err == nil {
		return response, errors.New("pro_id: " + strconv.Itoa(proId.ProId) + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)

	var convGroupDet model.ConvGroupDet
	structs.Automapper(request, &convGroupDet)

	convGroupDet.CreatedAt = &timeNow
	convGroupDet.CreatedBy = &request.CreatedBy
	convGroupDet.UpdatedBy = &request.UpdatedBy
	convGroupDet.UpdatedAt = &timeNow

	ConvGrpId, err := service.ConvGroupDetRepository.Store(convGroupDet)
	if err != nil {
		return response, err
	}

	response.ConvGrpId = ConvGrpId

	return response, err
}

func (service *convGroupDetServiceImpl) Update(convGrpId int, proId int, request entity.UpdateConvGroupDetRequest) (err error) {

	// conv_grp_code & cust id validation, if err == nil and params convGroupId != convGroup.Id, this means that code & cust id already exists
	convGroup, err := service.ConvGroupDetRepository.FindOneByProIdAndCustId(request.ConvGrpId, request.ProId, request.CustId)
	if err == nil && convGroup.ProId != proId {
		return errors.New("pro_id: " + strconv.Itoa(convGroup.ProId) + " is already exists")
	}

	err = service.ConvGroupDetRepository.Update(convGrpId, proId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *convGroupDetServiceImpl) Delete(custId string, convGrpId int, proId int, userId int64) (err error) {

	err = service.ConvGroupDetRepository.Delete(custId, convGrpId, proId, userId)
	if err != nil {
		return err
	}

	return err
}
