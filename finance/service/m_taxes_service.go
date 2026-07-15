package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"
)

const (
	MAX_RANGE = 1000
)

type MTaxesService interface {
	Store(request entity.MTaxesCreateReq) (err error)
	Detail(mTaxID int64, custID string) (response entity.MTaxesResp, err error)
	List(dataFilter entity.MTaxQueryFilter) (data []entity.MTaxesResp, total int64, lastPage int, err error)
	Delete(custId string, mtaxID int64, userId int64) (err error)
	Update(mTaxID int64, request entity.MTaxesUpdateReq) (err error)
}

type MTaxesServiceImpl struct {
	Repository      repository.MTaxesRepository
	Transaction     repository.Dbtransaction
	TaxesRepository repository.TaxesRepository
}

func NewMTaxesService(repository repository.MTaxesRepository, taxesRepository repository.TaxesRepository, transaction repository.Dbtransaction) *MTaxesServiceImpl {
	return &MTaxesServiceImpl{
		Repository:      repository,
		Transaction:     transaction,
		TaxesRepository: taxesRepository,
	}
}
func (service *MTaxesServiceImpl) Store(request entity.MTaxesCreateReq) (err error) {
	c := context.Background()

	var mTaxes model.MTaxes
	err = structs.Automapper(request, &mTaxes)
	if err != nil {
		return err
	}

	if request.To-request.From > MAX_RANGE {
		return errors.New(fmt.Sprintf("max range from and to can't be larger than %v", MAX_RANGE))
	}

	taxFrom, err := service.Repository.GetSeriesFromNoRange(request.CustID, request.From, request.To, request.Year)
	if err == nil {
		if taxFrom.MTaxID != nil {
			return errors.New(fmt.Sprintf("from %v AND to %v already at range tax ID %v", request.From, request.To, *taxFrom.MTaxID))
		}
	}

	taxTo, err := service.Repository.GetSeriesToNoRange(request.CustID, request.From, request.To, request.Year)
	if err == nil {
		if taxTo.MTaxID != nil {
			return errors.New(fmt.Sprintf("from %v AND to %v already at range tax ID %v", request.From, request.To, *taxTo.MTaxID))
		}
	}

	var sequence int
	lasttaxByYear, err := service.Repository.GetLastSequenceByYear(request.CustID, request.Year)
	if err != nil {
		sequence = 1
	} else {
		sequence = lasttaxByYear.Sequence + 1
	}

	status := entity.STATUS_TAXES_RESERVED
	mTaxes.Sequence = sequence
	mTaxes.Status = &status
	mTaxes.SerialFrom = request.From
	mTaxes.SerialTo = request.To
	mTaxes.RemainingQty = mTaxes.SerialTo - (mTaxes.SerialFrom - 1)
	mTaxes.TotalTaxNo = (request.To + 1) - request.From

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.Store(txCtx, &mTaxes)
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

func (service *MTaxesServiceImpl) Detail(mTaxID int64, custID string) (response entity.MTaxesResp, err error) {
	tax, err := service.Repository.GetByID(custID, mTaxID)
	if err != nil {
		return response, err
	}

	statusUsedCount, err := service.TaxesRepository.CountTaxesByStatusAndMTax(custID, mTaxID, 1)
	if err != nil {
		return response, err
	}

	statusDeletedCount, err := service.TaxesRepository.CountTaxesByStatusAndMTax(custID, mTaxID, 0)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(tax, &response)
	if err != nil {
		return response, err
	}

	response.From = tax.SerialFrom
	response.To = tax.SerialTo
	response.UsedTotal = statusUsedCount
	response.DeletedTotal = statusDeletedCount
	return response, err
}

func (service *MTaxesServiceImpl) List(dataFilter entity.MTaxQueryFilter) (data []entity.MTaxesResp, total int64, lastPage int, err error) {
	taxes, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range taxes {
		var vResp entity.MTaxesResp
		structs.Automapper(row, &vResp)

		vResp.From = row.SerialFrom
		vResp.To = row.SerialTo

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *MTaxesServiceImpl) Delete(custId string, mtaxID int64, userId int64) (err error) {
	c := context.Background()

	tax, err := service.Repository.GetByID(custId, mtaxID)
	if err == nil {
		if *tax.Status != entity.STATUS_TAXES_RESERVED {
			return errors.New("delete only tax status with status reserved")
		}
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, mtaxID, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *MTaxesServiceImpl) Update(mTaxID int64, request entity.MTaxesUpdateReq) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.MTaxes
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	_, err = service.Repository.GetByID(request.CustID, mTaxID)
	if err != nil {
		return err
	}

	taxFrom, err := service.Repository.GetSeriesFromNoRange(request.CustID, request.From, request.To, request.Year)
	if err == nil {
		if taxFrom.MTaxID != nil && taxFrom.MTaxID == &mTaxID {
			return errors.New(fmt.Sprintf("from %v AND to %v already at range tax ID %v", request.From, request.To, *taxFrom.MTaxID))
		}
	}

	taxTo, err := service.Repository.GetSeriesToNoRange(request.CustID, request.From, request.To, request.Year)
	if err == nil {
		if taxTo.MTaxID != nil && taxTo.MTaxID == &mTaxID {
			return errors.New(fmt.Sprintf("from %v AND to %v already at range tax ID %v", request.From, request.To, *taxTo.MTaxID))
		}
	}

	Model.CustID = ""
	Model.TotalTaxNo = (Model.SerialTo + 1) - Model.SerialFrom
	Model.SerialFrom = request.From
	Model.SerialTo = request.To
	Model.RemainingQty = Model.SerialTo - (Model.SerialFrom - 1)

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Update(txCtx, mTaxID, Model)
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
