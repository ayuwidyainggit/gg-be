package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"time"
)

type ChequeService interface {
	Store(request entity.CreateChequeBody) (err error)
	Detail(ChequeNo int, custID string) (response entity.ChequeResponse, err error)
	List(dataFilter entity.GeneralQueryFilter) (data []entity.ChequeListResponse, total int64, lastPage int, err error)
	Delete(custId string, ChequeNo int, userId int64) (err error)
	Update(ChequeNo int, request entity.UpdateChequeBody) (err error)
}

type ChequeServiceImpl struct {
	ChequeRepository repository.ChequeRepository
	Transaction      repository.Dbtransaction
}

func NewChequeService(repository repository.ChequeRepository, transaction repository.Dbtransaction) *ChequeServiceImpl {
	return &ChequeServiceImpl{
		ChequeRepository: repository,
		Transaction:      transaction,
	}
}

func (service *ChequeServiceImpl) Store(request entity.CreateChequeBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.ChqDate != nil {
		ChqDate, err := str.DateStrToRfc3339String(*request.ChqDate)
		if err != nil {
			return err
		}
		request.ChqDate = &ChqDate
	}

	if request.ChqDueDate != nil {
		ChqDueDate, err := str.DateStrToRfc3339String(*request.ChqDueDate)
		if err != nil {
			return err
		}
		request.ChqDueDate = &ChqDueDate
	}

	if request.ClearingDate != nil {
		ClearingDate, err := str.DateStrToRfc3339String(*request.ClearingDate)
		if err != nil {
			return err
		}
		request.ClearingDate = &ClearingDate
	}

	if request.StatusDate != nil {
		StatusDate, err := str.DateStrToRfc3339String(*request.StatusDate)
		if err != nil {
			return err
		}
		request.StatusDate = &StatusDate
	}

	var Chequemodel model.Cheque
	err = structs.Automapper(request, &Chequemodel)
	if err != nil {
		return err
	}
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.ChequeRepository.Store(txCtx, &Chequemodel)
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

func (service *ChequeServiceImpl) Detail(ChequeNo int, custID string) (response entity.ChequeResponse, err error) {
	Cheque, err := service.ChequeRepository.FindByNo(ChequeNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Cheque, &response)
	if err != nil {
		return response, err
	}

	ChqDate := Cheque.ChqDate.Format("2006-01-02")
	response.ChqDate = &ChqDate

	ChqDueDate := Cheque.ChqDueDate.Format("2006-01-02")
	response.ChqDueDate = &ChqDueDate

	ClearingDate := Cheque.ClearingDate.Format("2006-01-02")
	response.ClearingDate = &ClearingDate

	StatusDate := Cheque.StatusDate.Format("2006-01-02")
	response.StatusDate = &StatusDate

	return response, nil
}

func (service *ChequeServiceImpl) List(dataFilter entity.GeneralQueryFilter) (data []entity.ChequeListResponse, total int64, lastPage int, err error) {
	Cheques, total, lastPage, err := service.ChequeRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Cheques {
		var vResp entity.ChequeListResponse
		structs.Automapper(row, &vResp)
		if row.ChqDate != nil {
			ChqDate := row.ChqDate.Format("2006-01-02")
			vResp.ChqDate = &ChqDate
		}
		if row.ChqDueDate != nil {
			ChqDueDate := row.ChqDueDate.Format("2006-01-02")
			vResp.ChqDueDate = &ChqDueDate
		}
		if row.ClearingDate != nil {
			ClearingDate := row.ClearingDate.Format("2006-01-02")
			vResp.ClearingDate = &ClearingDate
		}
		if row.StatusDate != nil {
			StatusDate := row.StatusDate.Format("2006-01-02")
			vResp.StatusDate = &StatusDate
		}
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ChequeServiceImpl) Delete(custId string, ChequeNo int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeRepository.Delete(txCtx, custId, ChequeNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ChequeServiceImpl) Update(ChequeNo int, request entity.UpdateChequeBody) (err error) {
	c := context.Background()

	if request.ChqDate != nil {
		ChqDate, err := str.DateStrToRfc3339String(*request.ChqDate)
		if err != nil {
			return err
		}
		request.ChqDate = &ChqDate
	}

	if request.ChqDueDate != nil {
		ChqDueDate, err := str.DateStrToRfc3339String(*request.ChqDueDate)
		if err != nil {
			return err
		}
		request.ChqDueDate = &ChqDueDate
	}

	if request.ClearingDate != nil {
		ClearingDate, err := str.DateStrToRfc3339String(*request.ClearingDate)
		if err != nil {
			return err
		}
		request.ClearingDate = &ClearingDate
	}

	if request.StatusDate != nil {
		StatusDate, err := str.DateStrToRfc3339String(*request.StatusDate)
		if err != nil {
			return err
		}
		request.StatusDate = &StatusDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Cheque
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	if Model.IsPosted != nil {
		if *Model.IsPosted {
			now := time.Now()
			Model.PostedAt = &now
		}
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeRepository.Update(txCtx, ChequeNo, request.CustID, Model)
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
