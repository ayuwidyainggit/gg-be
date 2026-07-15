package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"log"
)

type ApDistributorDiscountService interface {
	Store(request entity.CreateApDistributorDiscountBody) (err error)
	Detail(ApDistributorDiscountId int64, custID string, parentCustID string) (response entity.ApDistributorDiscountResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ApDistributorDiscountListResponse, total int64, lastPage int, err error)
	LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ApDistributorDiscountLookupListResponse, total int64, lastPage int, err error)
	Delete(custId string, ApDistributorDiscountId int64, userId int64) (err error)
	Update(ApDistributorDiscountId int64, request entity.UpdateApDistributorDiscountBody, custId string) (err error)
}

type ApDistributorDiscountServiceImpl struct {
	ApDistributorDiscountRepository repository.ApDistributorDiscountRepository
	Transaction                     repository.Dbtransaction
}

func NewApDistributorDiscountService(repository repository.ApDistributorDiscountRepository, transaction repository.Dbtransaction) *ApDistributorDiscountServiceImpl {
	return &ApDistributorDiscountServiceImpl{
		ApDistributorDiscountRepository: repository,
		Transaction:                     transaction,
	}
}

func (service *ApDistributorDiscountServiceImpl) Store(request entity.CreateApDistributorDiscountBody) (err error) {
	c := context.Background()

	var ApDistributorDiscountModel model.ApDistributorDiscount
	err = structs.Automapper(request, &ApDistributorDiscountModel)
	if err != nil {
		return err
	}
	log.Println("ApDistributorDiscountModel:", structs.StructToJson(ApDistributorDiscountModel))
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.ApDistributorDiscountRepository.Store(txCtx, &ApDistributorDiscountModel)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *ApDistributorDiscountServiceImpl) Detail(ApDistributorDiscountId int64, custID string, parentCustID string) (response entity.ApDistributorDiscountResponse, err error) {
	ApDistributorDiscount, err := service.ApDistributorDiscountRepository.FindByID(ApDistributorDiscountId, custID, parentCustID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ApDistributorDiscount, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (service *ApDistributorDiscountServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ApDistributorDiscountListResponse, total int64, lastPage int, err error) {
	ApDistributorDiscounts, total, lastPage, err := service.ApDistributorDiscountRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ApDistributorDiscounts {
		var vResp entity.ApDistributorDiscountListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ApDistributorDiscountServiceImpl) LookupList(dataFilter entity.GeneralQueryFilter) (data []entity.ApDistributorDiscountLookupListResponse, total int64, lastPage int, err error) {
	ApDistributorDiscounts, total, lastPage, err := service.ApDistributorDiscountRepository.FindAllByCustIdLookup(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ApDistributorDiscounts {
		var vResp entity.ApDistributorDiscountLookupListResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ApDistributorDiscountServiceImpl) Delete(custId string, ApDistributorDiscountId int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ApDistributorDiscountRepository.Delete(txCtx, custId, ApDistributorDiscountId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ApDistributorDiscountServiceImpl) Update(ApDistributorDiscountId int64, request entity.UpdateApDistributorDiscountBody, custId string) (err error) {
	c := context.Background()

	// log.Println("request:", structs.StructToJson(request))
	// End parse time format YYYY-mm-dd to Rfc339
	var ApDistributorDiscount model.ApDistributorDiscount
	err = structs.Automapper(request, &ApDistributorDiscount)
	if err != nil {
		return err
	}
	ApDistributorDiscount.CustID = ""

	// log.Println("model.ApDistributorDiscount:", structs.StructToJson(ApDistributorDiscount))
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ApDistributorDiscountRepository.Update(txCtx, ApDistributorDiscountId, ApDistributorDiscount, custId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
