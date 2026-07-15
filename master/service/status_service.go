package service

import (
	"master/entity"
	"master/pkg/structs"
	"master/repository"
)

type StatusService interface {
	Detail(string, int, string) (entity.StatusResponse, error)
	List(entity.StatusQueryFilter) (data []entity.StatusListResponse, total int, lastPage int, err error)
	LookupList(entity.StatusQueryFilter) (data []entity.StatusLookupResponse, total int, lastPage int, err error)
	// Store(entity.CreateStatusBody) (entity.StatusResponse, error)
	// Update(int, entity.UpdateStatusRequest) error
	// Delete(string, int, int64) error
}

func NewStatusService(statusRepository repository.StatusRepository) *statusServiceImpl {
	return &statusServiceImpl{
		StatusRepository: statusRepository,
	}
}

type statusServiceImpl struct {
	StatusRepository repository.StatusRepository
}

func (service *statusServiceImpl) Detail(statusId string, statusValue int, langId string) (response entity.StatusResponse, err error) {
	status, err := service.StatusRepository.FindOneByStatusIdAndStatusValue(statusId, statusValue, langId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(status, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *statusServiceImpl) List(dataFilter entity.StatusQueryFilter) (data []entity.StatusListResponse, total int, lastPage int, err error) {
	statuses, total, lastPage, err := service.StatusRepository.FindAllByStatusIdAndLangId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range statuses {
		var vResp entity.StatusListResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *statusServiceImpl) LookupList(dataFilter entity.StatusQueryFilter) (data []entity.StatusLookupResponse, total int, lastPage int, err error) {
	statuss, total, lastPage, err := service.StatusRepository.FindAllByStatusIdAndLangIdLookup(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range statuss {
		var vResp entity.StatusLookupResponse
		err = structs.Automapper(row, &vResp)
		if err != nil {
			return data, total, lastPage, err
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

// func (service *statusServiceImpl) Store(request entity.CreateStatusBody) (response entity.StatusResponse, err error) {

// 	timeNow := time.Now().In(time.UTC)
// 	var statusData model.Status
// 	err = structs.Automapper(request, &statusData)
// 	if err != nil {
// 		return response, err
// 	}
// 	statusData.CreatedAt = &timeNow
// 	statusData.CreatedBy = &request.CreatedBy
// 	statusData.UpdatedAt = &timeNow
// 	statusData.UpdatedBy = &request.CreatedBy

// 	statusId, err := service.StatusRepository.Store(statusData)
// 	if err != nil {
// 		return response, err
// 	}

// 	response.StatusId = statusId

// 	return response, err
// }

// func (service *statusServiceImpl) Update(statusId int, request entity.UpdateStatusRequest) (err error) {

// 	err = service.StatusRepository.Update(statusId, request)
// 	if err != nil {
// 		return err
// 	}

// 	return err
// }

// func (service *statusServiceImpl) Delete(custId string, statusId int, userId int64) (err error) {

// 	err = service.StatusRepository.Delete(custId, statusId, userId)
// 	if err != nil {
// 		return err
// 	}

// 	return err
// }
