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

type MParkingFundService interface {
	Detail(int, string) (entity.MParkingFundResponse, error)
	List(entity.GeneralQueryFilter, string) (data []entity.MParkingFundResponse, total int, lastPage int, err error)
	Store(entity.CreateMParkingFundBody) (entity.MParkingFundResponse, error)
	Update(int, entity.UpdateMParkingFundRequest) error
	Delete(string, int, int64) error
}

func NewMParkingFundService(mParkingFundRepository repository.MParkingFundRepository) *mParkingFundServiceImpl {
	return &mParkingFundServiceImpl{
		MParkingFundRepository: mParkingFundRepository,
	}
}

type mParkingFundServiceImpl struct {
	MParkingFundRepository repository.MParkingFundRepository
}

func (service *mParkingFundServiceImpl) Detail(ParkingFundId int, custId string) (response entity.MParkingFundResponse, err error) {
	remarkPromo, err := service.MParkingFundRepository.FindOneBymParkingFundIdAndCustId(ParkingFundId, custId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(remarkPromo, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

func (service *mParkingFundServiceImpl) List(dataFilter entity.GeneralQueryFilter, custId string) (data []entity.MParkingFundResponse, total int, lastPage int, err error) {
	remarkPromos, total, lastPage, err := service.MParkingFundRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range remarkPromos {
		var vResp entity.MParkingFundResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *mParkingFundServiceImpl) Store(request entity.CreateMParkingFundBody) (response entity.MParkingFundResponse, err error) {
	mParkingFunds, err := service.MParkingFundRepository.FindOneByOutletAndProIdIdAndCustId(request.OutletId, request.ProId, request.CustId)
	if err == nil {
		return response, errors.New("outlet_id: " + strconv.Itoa(*mParkingFunds.OutletId) + " is already exists")
	}

	var mParkingFund model.MParkingFund

	timeNow := time.Now().In(time.UTC)
	structs.Automapper(request, &mParkingFund)
	mParkingFund.CreatedAt = &timeNow
	mParkingFund.CreatedBy = &request.CreatedBy
	mParkingFund.UpdatedBy = &request.UpdatedBy
	mParkingFund.UpdatedAt = &timeNow
	ParkingFundId, err := service.MParkingFundRepository.Store(mParkingFund)
	if err != nil {
		return response, err
	}

	response.ParkingFundId = ParkingFundId

	return response, err
}

func (service *mParkingFundServiceImpl) Update(parkingFundId int, request entity.UpdateMParkingFundRequest) (err error) {
	// remarkPromo_code & cust id validation, if err == nil and params parkingFundId != remarkPromo.Id, this means that code & cust id already exists
	// remarkPromo, err := service.MParkingFundRepository.FindOneByRemarkPromoCodeAndCustId(request.RemPromoCode, request.CustId)
	// if err == nil && remarkPromo.RemPromoId != parkingFundId {
	// 	return errors.New("remark_promo_code: " + remarkPromo.RemPromoCode + " is already exists")
	// }
	err = service.MParkingFundRepository.Update(parkingFundId, request)
	if err != nil {
		return err
	}

	return err
}

func (service *mParkingFundServiceImpl) Delete(custId string, parkingFundId int, userId int64) (err error) {

	err = service.MParkingFundRepository.Delete(custId, parkingFundId, userId)
	if err != nil {
		return err
	}
	return err
}
