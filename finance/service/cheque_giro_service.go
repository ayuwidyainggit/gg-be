package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type ChequeGiroService interface {
	Store(request entity.CreateChequeGiroBody) (err error)
	Detail(ChequeGiroNo int, custID string) (response entity.ChequeGiroResponse, err error)
	List(dataFilter entity.CheckGiroQueryFilter) (data []entity.ChequeGiroResponse, total int64, lastPage int, err error)
	Delete(custId string, ChequeGiroNo int, userId int64) (err error)
	Update(ChequeGiroNo int, request entity.UpdateChequeGiroBody) (err error)

	LookupBank(dataFilter entity.CheckGiroQueryFilter) (data []entity.BankLookup, total int64, lastPage int, err error)
	LookupBankAccount(dataFilter entity.CheckGiroQueryFilter) (data []entity.BankAccountLookup, total int64, lastPage int, err error)
}

type ChequeGiroServiceImpl struct {
	ChequeGiroRepository repository.ChequeGiroRepository
	Transaction          repository.Dbtransaction
}

func NewChequeGiroService(repository repository.ChequeGiroRepository, transaction repository.Dbtransaction) *ChequeGiroServiceImpl {
	return &ChequeGiroServiceImpl{
		ChequeGiroRepository: repository,
		Transaction:          transaction,
	}
}

func (service *ChequeGiroServiceImpl) Store(request entity.CreateChequeGiroBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DocDateCheque != nil {
		ChqDate, err := str.DateStrToRfc3339String(*request.DocDateCheque)
		if err != nil {
			return err
		}
		request.DocDateCheque = &ChqDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	var ChequeGiromodel model.ChequeGiro
	err = structs.Automapper(request, &ChequeGiromodel)
	if err != nil {
		return err
	}

	if ChequeGiromodel.OwnerID == 1 {
		ChequeGiromodel.SupplierID = nil
	} else {
		ChequeGiromodel.OutletID = nil
		ChequeGiromodel.SalesmanID = nil
		ChequeGiromodel.OutletBankID = nil
	}

	ChequeGiromodel.StatusCheque = 2

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.ChequeGiroRepository.Store(txCtx, &ChequeGiromodel)
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

func (service *ChequeGiroServiceImpl) Detail(ChequeGiroNo int, custID string) (response entity.ChequeGiroResponse, err error) {
	ChequeGiro, err := service.ChequeGiroRepository.FindByNo(ChequeGiroNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ChequeGiro, &response)
	if err != nil {
		return response, err
	}

	ChqDate := ChequeGiro.DocDateCheque.Format("2006-01-02")
	response.DocDateCheque = &ChqDate

	ChqDueDate := ChequeGiro.DueDate.Format("2006-01-02")
	response.DueDate = &ChqDueDate

	ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	response.OwnerName = ownerName

	statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	if response.OwnerID == 1 {
		response.UsedAmount = ChequeGiro.UsedAmountOutlet
	}
	response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *ChequeGiroServiceImpl) List(dataFilter entity.CheckGiroQueryFilter) (data []entity.ChequeGiroResponse, total int64, lastPage int, err error) {
	ChequeGiros, total, lastPage, err := service.ChequeGiroRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ChequeGiros {
		var vResp entity.ChequeGiroResponse
		structs.Automapper(row, &vResp)
		if row.DocDateCheque != nil {
			DocDateCheque := row.DocDateCheque.Format("2006-01-02")
			vResp.DocDateCheque = &DocDateCheque
		}
		if row.DueDate != nil {
			DueDate := row.DueDate.Format("2006-01-02")
			vResp.DueDate = &DueDate
		}

		ownerName := entity.ConvStatus(entity.OwnerGiro, row.OwnerID)
		vResp.OwnerName = ownerName

		statusText := entity.ConvStatus(entity.StatusGiro, row.StatusCheque)
		vResp.StatusChequeText = &statusText

		vResp.UsedAmount = float64(0)
		vResp.RemainingAmount = row.Amount - vResp.UsedAmount

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ChequeGiroServiceImpl) Delete(custId string, ChequeGiroNo int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeGiroRepository.Delete(txCtx, custId, ChequeGiroNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ChequeGiroServiceImpl) Update(ChequeGiroNo int, request entity.UpdateChequeGiroBody) (err error) {
	c := context.Background()

	if request.DocDateCheque != nil {
		ChqDate, err := str.DateStrToRfc3339String(*request.DocDateCheque)
		if err != nil {
			return err
		}
		request.DocDateCheque = &ChqDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ChequeGiro
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeGiroRepository.Update(txCtx, ChequeGiroNo, Model)
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

func (service *ChequeGiroServiceImpl) LookupBank(dataFilter entity.CheckGiroQueryFilter) (data []entity.BankLookup, total int64, lastPage int, err error) {
	ChequeGiros, total, lastPage, err := service.ChequeGiroRepository.FindAllBankByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ChequeGiros {
		var vResp entity.BankLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ChequeGiroServiceImpl) LookupBankAccount(dataFilter entity.CheckGiroQueryFilter) (data []entity.BankAccountLookup, total int64, lastPage int, err error) {
	ChequeGiros, total, lastPage, err := service.ChequeGiroRepository.FindAllBankAccountByCustId(dataFilter, dataFilter.BankID)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ChequeGiros {
		var vResp entity.BankAccountLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
