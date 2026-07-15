package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type TprService interface {
	Detail(int64, string) (entity.TprResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.TprListResponse, total int, lastPage int, err error)
	Store(request entity.CreateTprBody) (response entity.TprResponse, err error)
	Update(tprID int64, request entity.UpdateTprRequest) (err error)
	Delete(custId string, tprId int64, userId int64) (err error)
}

func NewTprService(tprRepository repository.TprRepository) *tprServiceImpl {
	return &tprServiceImpl{
		TprRepository: tprRepository,
	}
}

type tprServiceImpl struct {
	TprRepository repository.TprRepository
}

func (service *tprServiceImpl) Detail(tprID int64, custID string) (response entity.TprResponse, err error) {
	tpr, err := service.TprRepository.FindOneByTprIdAndCustId(tprID, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(tpr, &response)
	if err != nil {
		return response, err
	}

	if tpr.DateStart != nil {
		dateStart := tpr.DateStart.Format("2006-01-02")
		response.DateStart = &dateStart
	}

	if tpr.DateEnd != nil {
		dateEnd := tpr.DateEnd.Format("2006-01-02")
		response.DateEnd = &dateEnd
	}

	return response, err
}

func (service *tprServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.TprListResponse, total int, lastPage int, err error) {
	tprs, total, lastPage, err := service.TprRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(tprs) > 0 {
		for _, row := range tprs {
			var vResp entity.TprListResponse
			structs.Automapper(row, &vResp)

			if row.DateStart != nil {
				dateStart := row.DateStart.Format("2006-01-02")
				vResp.DateStart = &dateStart
			}
			if row.DateEnd != nil {
				dateEnd := row.DateEnd.Format("2006-01-02")
				vResp.DateEnd = &dateEnd
			}
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *tprServiceImpl) Store(request entity.CreateTprBody) (response entity.TprResponse, err error) {
	tpr, err := service.TprRepository.FindOneByTprCodeAndCustId(request.TprCode, request.CustID)
	if err == nil {
		return response, errors.New("tpr_id: " + tpr.TprCode + " is already exists")
	}
	timeNow := time.Now().In(time.UTC)

	if request.DateEnd != nil {
		if *request.DateEnd != "" {
			dateEnd, err := str.DateStrToRfc3339String(*request.DateEnd)
			if err != nil {
				return response, err
			}
			request.DateEnd = &dateEnd
		} else {
			request.DateEnd = nil
		}
	}

	if request.DateStart != nil {
		if *request.DateStart != "" {

			dateStart, err := str.DateStrToRfc3339String(*request.DateStart)
			if err != nil {
				return response, err
			}
			request.DateStart = &dateStart
		} else {
			request.DateStart = nil
		}
	}
	var tprData model.MTpr

	structs.Automapper(request, &tprData)

	tprData.CreatedAt = &timeNow
	tprData.CreatedBy = &request.CreatedBy
	tprData.UpdatedAt = &timeNow

	tprID, err := service.TprRepository.Store(tprData)
	if err != nil {
		return response, err
	}

	response.TprID = tprID
	return response, err
}

func (service *tprServiceImpl) Update(tprID int64, request entity.UpdateTprRequest) (err error) {

	// vehicle_code & cust id validation, if err == nil and params vehicleId != vehicle.Id, this means that code & cust id already exists
	tpr, err := service.TprRepository.FindOneByTprCodeAndCustId(request.TprCode, request.CustId)
	if err == nil && tpr.TprID != tprID {
		return errors.New("tpr_code: " + tpr.TprCode + " is already exists")
	}

	err = service.TprRepository.Update(tprID, request)
	if err != nil {
		return err
	}

	return err
}
func (service *tprServiceImpl) Delete(custId string, tprId int64, userId int64) (err error) {
	err = service.TprRepository.Delete(custId, tprId, userId)
	if err != nil {
		return err
	}
	return err
}
