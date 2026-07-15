package service

import (
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"
)

type TprLimitService interface {
	Detail(int, string) (entity.TprLimitListResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.TprLimitListResponse, total int, lastPage int, err error)
	Store(entity.CreatedTprLimitBody) (entity.TprLimitResponse, error)
	Update(int, entity.UpdateTprLimitRequest) error
	Delete(string, int, int64) error
}

func NewTprLimitService(tprLimitRepository repository.TprLimitRepository) *tprLimitServiceImpl {
	return &tprLimitServiceImpl{
		TprLimitRepository: tprLimitRepository,
	}
}

type tprLimitServiceImpl struct {
	TprLimitRepository repository.TprLimitRepository
}

func (service *tprLimitServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.TprLimitListResponse, total int, lastPage int, err error) {
	tprLimits, total, lastPage, err := service.TprLimitRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range tprLimits {
		var vResp entity.TprLimitListResponse
		structs.Automapper(row, &vResp)

		if row.DateStart != nil {
			dateStart := row.DateStart.Format("2006-01-02")
			vResp.DateStart = dateStart
		}

		if row.DateEnd != nil {
			dateEnd := row.DateEnd.Format("2006-01-02")
			vResp.DateEnd = dateEnd
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *tprLimitServiceImpl) Detail(outletId int, tprLimitID string) (response entity.TprLimitListResponse, err error) {
	tprLimit, err := service.TprLimitRepository.FindOneByTprLimitIdAndCustId(outletId, tprLimitID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(tprLimit, &response)
	if err != nil {
		return response, err
	}

	if tprLimit.DateStart != nil {
		dateStart := tprLimit.DateStart.Format("2006-01-02")
		response.DateStart = dateStart
	}

	if tprLimit.DateEnd != nil {
		dateEnd := tprLimit.DateEnd.Format("2006-01-02")
		response.DateEnd = dateEnd
	}

	return response, err
}

func (service *tprLimitServiceImpl) Store(request entity.CreatedTprLimitBody) (response entity.TprLimitResponse, err error) {

	// conv_grp_code & cust id validation, if err == nil, this means that code & cust id already exists
	// convGroup, err := service.tprLimitRepository.FindOneByConvGroupCodeAndCustId(request.ConvGroupCode, request.CustId)
	// if err == nil {
	// 	return response, errors.New("conv_grp_code: " + convGroup.ConvGroupCode + " is already exists")
	// }

	var tprLimit model.MTprLimit
	if request.DateStart != nil {
		if *request.DateStart == "" {
			request.DateStart = nil
		}
	}

	if request.DateEnd != nil {
		if *request.DateEnd == "" {
			request.DateEnd = nil
		}
	}
	timeNow := time.Now().In(time.UTC)

	structs.Automapper(request, &tprLimit)

	tprLimit.CreatedAt = &timeNow
	tprLimit.CreatedBy = &request.CreatedBy
	tprLimit.UpdatedBy = &request.UpdatedBy
	tprLimit.UpdatedAt = &timeNow

	TprLimitId, err := service.TprLimitRepository.Store(tprLimit)
	if err != nil {
		return response, err
	}

	response.TprLimitId = TprLimitId

	return response, err
}

func (service *tprLimitServiceImpl) Update(proId int, request entity.UpdateTprLimitRequest) (err error) {

	// conv_grp_code & cust id validation, if err == nil and params convGroupId != convGroup.Id, this means that code & cust id already exists
	// convGroup, err := service.ConvGroupDetRepository.FindOneByPluProIdAndCustId(request.ProId, request.CustId)
	// if err == nil && convGroup.ProId != proId {
	// 	return errors.New("pro_id: " + strconv.Itoa(convGroup.ProId) + " is already exists")
	// }

	err = service.TprLimitRepository.Update(proId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *tprLimitServiceImpl) Delete(custId string, tprLimitId int, userId int64) (err error) {

	err = service.TprLimitRepository.Delete(custId, tprLimitId, userId)
	if err != nil {
		return err
	}

	return err
}
