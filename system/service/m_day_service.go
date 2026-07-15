package service

import (
	"system/entity"
	"system/pkg/structs"
	"system/repository"
)

type MDayService interface {
	List(dataFilter entity.GeneralQueryFilter, LangId string) (data []entity.MDayListResponse, total int64, lastPage int, err error)
	Detail(dayId int64, langId string) (response entity.MDayListResponse, err error)
}

type mDayServiceImpl struct {
	MDayRepository repository.MDayRepository
}

func NewMDayService(mDayRepository repository.MDayRepository) *mDayServiceImpl {
	return &mDayServiceImpl{
		MDayRepository: mDayRepository,
	}
}

func (service *mDayServiceImpl) List(dataFilter entity.GeneralQueryFilter, LangId string) (data []entity.MDayListResponse, total int64, lastPage int, err error) {
	whAdjs, total, lastPage, err := service.MDayRepository.FindAllByLangId(dataFilter, LangId)
	if err != nil {
		return data, total, lastPage, err
	}
	if len(whAdjs) > 0 {
		for _, row := range whAdjs {
			var vResp entity.MDayListResponse
			structs.Automapper(row, &vResp)

			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *mDayServiceImpl) Detail(dayId int64, langId string) (response entity.MDayListResponse, err error) {
	user, err := service.MDayRepository.FindDetail(dayId, langId)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(user, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
