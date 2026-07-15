package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type SubBeatService interface {
	Detail(int, string) (entity.SubBeatResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.SubBeatResponse, total int, lastPage int, err error)
	LookupList(entity.GeneralQueryFilter, string) (data []entity.SubBeatLookupResponse, total int, lastPage int, err error)
	Store(entity.CreateSubBeatBody) (entity.SubBeatResponse, error)
	Update(int, entity.UpdateSubBeatRequest) error
	Delete(string, int, int64) error
}

func NewSubBeatService(subBeatRepository repository.SubBeatRepository) *subBeatServiceImpl {
	return &subBeatServiceImpl{
		SubBeatRepository: subBeatRepository,
	}
}

type subBeatServiceImpl struct {
	SubBeatRepository repository.SubBeatRepository
}

func (service *subBeatServiceImpl) Detail(subBeatId int, custId string) (response entity.SubBeatResponse, err error) {
	subBeat, err := service.SubBeatRepository.FindOneBySubBeatIdAndCustId(subBeatId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(subBeat, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *subBeatServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SubBeatResponse, total int, lastPage int, err error) {
	subBeats, total, lastPage, err := service.SubBeatRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBeats {
		var vResp entity.SubBeatResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *subBeatServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.SubBeatLookupResponse, total int, lastPage int, err error) {
	subBeats, total, lastPage, err := service.SubBeatRepository.FindAllByCustIdLookupMode(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range subBeats {
		var vResp entity.SubBeatLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *subBeatServiceImpl) Store(request entity.CreateSubBeatBody) (response entity.SubBeatResponse, err error) {

	// subBeat_code & cust id validation, if err == nil, this means that code & cust id already exists
	subBeat, err := service.SubBeatRepository.FindOneBySubBeatCodeAndCustId(request.SbeatCode, request.CustId)
	if err == nil {
		return response, errors.New("sbeat_code: " + subBeat.SbeatCode + " is already exists")
	}

	timeNow := time.Now().In(time.UTC)
	subBeatData := model.SubBeat{
		CustId:     request.CustId,
		SbeatCode:  request.SbeatCode,
		SbeatName:  request.SbeatName,
		BeatId:     request.BeatId,
		DistrictId: request.DistrictId,
		IsActive:   request.IsActive,
		CreatedAt:  &timeNow,
		CreatedBy:  &request.CreatedBy,
		UpdatedAt:  &timeNow,
		UpdatedBy:  &request.CreatedBy,
	}

	subBeatId, err := service.SubBeatRepository.Store(subBeatData)
	if err != nil {
		return response, err
	}

	response.SbeatId = subBeatId

	return response, err
}

func (service *subBeatServiceImpl) Update(subBeatId int, request entity.UpdateSubBeatRequest) (err error) {

	// subBeat_code & cust id validation, if err == nil and params subBeatId != subBeat.Id, this means that code & cust id already exists
	subBeat, err := service.SubBeatRepository.FindOneBySubBeatCodeAndCustId(request.SbeatCode, request.CustId)
	if err == nil && subBeat.SbeatId != subBeatId {
		return errors.New("sbeat_code: " + subBeat.SbeatCode + " is already exists")
	}

	err = service.SubBeatRepository.Update(subBeatId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *subBeatServiceImpl) Delete(custId string, subBeatId int, userId int64) (err error) {

	err = service.SubBeatRepository.Delete(custId, subBeatId, userId)
	if err != nil {
		return err
	}

	return err
}
