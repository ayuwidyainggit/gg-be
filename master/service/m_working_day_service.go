package service

import (
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"strconv"
)

type MWorkingDayService interface {
	Active(string, string) (entity.MWorkingDayActiveResponse, error)
	Detail(int, int, int, string, string, string) (entity.MWorkingDayResponse, error)
	List(entity.MWorkingDayQueryFilter, string, string) (data []entity.MWorkingDayResponse, total int, lastPage int, err error)
	Store(entity.CreateMWorkingDayBody) (entity.MWorkingDayResponse, error)
	Update(int, int, int, string, entity.UpdateMWorkingDayRequest) error
	Delete(string, int, int, int, string, int64, string) error
}

func NewMWorkingDayService(mWorkingDayRepository repository.MWorkingDayRepository) *mWorkingDayServiceImpl {
	return &mWorkingDayServiceImpl{
		MWorkingDayRepository: mWorkingDayRepository,
	}
}

type mWorkingDayServiceImpl struct {
	MWorkingDayRepository repository.MWorkingDayRepository
}

func (service *mWorkingDayServiceImpl) Active(custId string, parentCustId string) (response entity.MWorkingDayActiveResponse, err error) {
	scopeCustId, generatedOnly, err := service.resolveReadScope(entity.MWorkingDayQueryFilter{}, custId, parentCustId)
	if err != nil {
		return response, err
	}

	mWorkDay, err := service.MWorkingDayRepository.FindOneActive(scopeCustId, generatedOnly)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mWorkDay, &response)
	if err != nil {
		return response, err
	}

	if mWorkDay.WorkDate != nil {
		workDate := mWorkDay.WorkDate.Format("2006-01-02")
		response.WorkDate = &workDate
	}

	return response, err
}

func (service *mWorkingDayServiceImpl) Detail(PerYear int, PerId int, WeekId int, WorkDate string, custId string, parentCustId string) (response entity.MWorkingDayResponse, err error) {
	scopeCustId, generatedOnly, err := service.resolveReadScope(entity.MWorkingDayQueryFilter{PerYear: strconv.Itoa(PerYear)}, custId, parentCustId)
	if err != nil {
		return response, err
	}

	mWeeks, err := service.MWorkingDayRepository.FindOneByPerYearAndPerIdAndWorkingDayIdAndCustId(PerYear, PerId, WeekId, WorkDate, scopeCustId, generatedOnly)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(mWeeks, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mWorkingDayServiceImpl) List(dataFilter entity.MWorkingDayQueryFilter, custId string, parentCustId string) (data []entity.MWorkingDayResponse, total int, lastPage int, err error) {
	scopeCustId, generatedOnly, err := service.resolveReadScope(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	MWorkingDays, total, lastPage, err := service.MWorkingDayRepository.FindAllByCustId(dataFilter, scopeCustId, generatedOnly)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range MWorkingDays {
		var vResp entity.MWorkingDayResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mWorkingDayServiceImpl) resolveReadScope(dataFilter entity.MWorkingDayQueryFilter, custId string, parentCustId string) (string, bool, error) {
	ownerCustId := parentCustId
	if ownerCustId == "" {
		ownerCustId = custId
	}

	if custId != ownerCustId {
		hasLegacyRows, err := service.MWorkingDayRepository.HasLegacyRows(custId, dataFilter.PerYear)
		if err != nil {
			return "", false, err
		}
		if hasLegacyRows {
			return custId, false, nil
		}
	}

	hasGeneratedRows, err := service.MWorkingDayRepository.HasGeneratedCalendarRows(ownerCustId, dataFilter.PerYear)
	if err != nil {
		return "", false, err
	}
	if hasGeneratedRows {
		return ownerCustId, true, nil
	}

	return custId, false, nil
}

func (service *mWorkingDayServiceImpl) Store(request entity.CreateMWorkingDayBody) (response entity.MWorkingDayResponse, err error) {
	// MWorkingDays, err := service.MWorkingDayRepository.FindOneByPerYearAndPerIdAndWeekIdAndCustId(request.PerYear, request.PerId, request.WeekId, request.CustId)
	// if err == nil {
	// 	return response, errors.New("week_id: " + strconv.Itoa(MWorkingDays.WeekId) + " is already exists")
	// }

	var mWorkingDayData model.MWorkingDay
	if request.WorkDate != nil {
		if *request.WorkDate == "" {
			request.WorkDate = nil
		}
	}
	structs.Automapper(request, &mWorkingDayData)
	PerId, err := service.MWorkingDayRepository.Store(mWorkingDayData)
	if err != nil {
		return response, err
	}

	response.PerId = PerId

	return response, err
}

func (service *mWorkingDayServiceImpl) Update(PerYear int, PerId int, WeekId int, WorkDate string, request entity.UpdateMWorkingDayRequest) (err error) {
	// MWorkingDays_code & cust id validation, if err == nil and params MWorkingDaysId != MWorkingDays.Id, this means that code & cust id already exists
	// MWorkingDays, err := service.MWorkingDaysRepository.FindOneByMWorkingDaysCodeAndCustId(request.RemPromoCode, request.CustId)
	// if err == nil && MWorkingDays.RemPromoId != MWorkingDaysId {
	// 	return errors.New("remark_promo_code: " + MWorkingDays.RemPromoCode + " is already exists")
	// }
	err = service.MWorkingDayRepository.Update(PerYear, PerId, WeekId, WorkDate, request)
	if err != nil {
		return err
	}

	return err
}

func (service *mWorkingDayServiceImpl) Delete(custId string, PerYear int, PerId int, WeekId int, WorkDate string, closedBy int64, closedByName string) (err error) {

	err = service.MWorkingDayRepository.Delete(custId, PerYear, PerId, WeekId, WorkDate, closedBy, closedByName)
	if err != nil {
		return err
	}
	return err
}
