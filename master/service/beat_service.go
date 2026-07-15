package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type BeatService interface {
	Detail(int, string) (entity.BeatResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.BeatResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.BeatLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateBeatBody) (entity.BeatResponse, error)
	Update(int, entity.UpdateBeatRequest) error
	Delete(string, int, int64) error
}

func NewBeatService(beatRepository repository.BeatRepository) *beatServiceImpl {
	return &beatServiceImpl{
		BeatRepository: beatRepository,
	}
}

type beatServiceImpl struct {
	BeatRepository repository.BeatRepository
}

func (service *beatServiceImpl) Detail(beatId int, custId string) (response entity.BeatResponse, err error) {
	beat, err := service.BeatRepository.FindOneByBeatIdAndCustId(beatId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(beat, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *beatServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.BeatResponse, total int, lastPage int, err error) {
	beats, total, lastPage, err := service.BeatRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range beats {
		var vResp entity.BeatResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *beatServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.BeatLookupResponse, total int, lastPage int, err error) {
	beats, total, lastPage, err := service.BeatRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range beats {
		var vResp entity.BeatLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *beatServiceImpl) Store(request entity.CreateBeatBody) (response entity.BeatResponse, err error) {

	// beat_code & cust id validation, if err == nil, this means that code & cust id already exists
	beat, err := service.BeatRepository.FindOneByBeatCodeAndCustId(request.BeatCode, request.CustId)
	if err == nil {
		return response, errors.New("beat_code: " + beat.BeatCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	beatData := model.Beat{
		CustId:     request.CustId,
		BeatCode:   request.BeatCode,
		BeatName:   request.BeatName,
		DistrictId: &request.DistrictId,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	beatId, err := service.BeatRepository.Store(beatData)
	if err != nil {
		return response, err
	}

	response.BeatId = beatId

	return response, err
}

func (service *beatServiceImpl) Update(beatId int, request entity.UpdateBeatRequest) (err error) {

	// beat_code & cust id validation, if err == nil and params beatId != beat.Id, this means that code & cust id already exists
	beat, err := service.BeatRepository.FindOneByBeatCodeAndCustId(request.BeatCode, request.CustId)
	if err == nil && beat.BeatId != beatId {
		return errors.New("beat_code: " + beat.BeatCode + " is already exists")
	}

	err = service.BeatRepository.Update(beatId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *beatServiceImpl) Delete(custId string, beatId int, userId int64) (err error) {

	err = service.BeatRepository.Delete(custId, beatId, userId)
	if err != nil {
		return err
	}

	return err
}
