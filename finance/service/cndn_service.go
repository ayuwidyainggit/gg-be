package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type CndnService interface {
	Store(request entity.CreateCndnBody) (err error)
	Detail(CndnNo string, custID string, parentCustId string) (response entity.CndnListDetailResponse, err error)
	List(dataFilter entity.CndnQueryFilter) (data []entity.CndnListResponse, total int64, lastPage int, err error)
	Delete(custId string, CndnNo string, userId int64) (err error)
	Update(CndnNo string, request entity.UpdateCndnBody) (err error)
}

type CndnServiceImpl struct {
	CndnRepository repository.CndnRepository
	Transaction    repository.Dbtransaction
}

func NewCndnService(repository repository.CndnRepository, transaction repository.Dbtransaction) *CndnServiceImpl {
	return &CndnServiceImpl{
		CndnRepository: repository,
		Transaction:    transaction,
	}
}

func (service *CndnServiceImpl) Store(request entity.CreateCndnBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.CndnDate != nil {
		CndnDate, err := str.DateStrToRfc3339String(*request.CndnDate)
		if err != nil {
			return err
		}
		request.CndnDate = &CndnDate
	}

	if request.LastTransactionDate != nil {
		LastTransactionDate, err := str.DateStrToRfc3339String(*request.LastTransactionDate)
		if err != nil {
			return err
		}
		request.LastTransactionDate = &LastTransactionDate
	}

	var Cndnmodel model.Cndn
	err = structs.Automapper(request, &Cndnmodel)
	if err != nil {
		return err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.CndnRepository.Store(txCtx, &Cndnmodel)
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

func (service *CndnServiceImpl) Detail(CndnNo string, custID string, parentCustId string) (response entity.CndnListDetailResponse, err error) {
	Cndn, err := service.CndnRepository.FindByNo(CndnNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Cndn, &response)
	if err != nil {
		return response, err
	}

	if Cndn.CndnDate != nil {
		CndnDate := Cndn.CndnDate.Format("2006-01-02")
		response.CndnDate = &CndnDate
	}

	if Cndn.LastTransactionDate != nil {
		LastTransactionDate := Cndn.LastTransactionDate.Format("2006-01-02")
		response.LastTransactionDate = &LastTransactionDate
	}

	ownerName := entity.ConvStatusOwnerId(entity.OwnerId, response.OwnerId)
	response.OwnerName = ownerName

	cndnJenisName := entity.ConvCndnJenis(entity.CndnJenis, response.CndnJenis)
	response.CndnJenis = cndnJenisName

	// response.UsedAmount = float64(0)
	if response.OwnerId == 1 {
		response.UsedAmount = Cndn.UsedAmountOutlet
	}
	response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *CndnServiceImpl) List(dataFilter entity.CndnQueryFilter) (data []entity.CndnListResponse, total int64, lastPage int, err error) {
	Cndns, total, lastPage, err := service.CndnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Cndns {
		var vResp entity.CndnListResponse
		structs.Automapper(row, &vResp)
		if row.CndnDate != nil {
			CndnDates := row.CndnDate.Format("2006-01-02")
			vResp.CndnDate = &CndnDates
		}
		if row.LastTransactionDate != nil {
			LastTransactionDates := row.LastTransactionDate.Format("2006-01-02")
			vResp.LastTransactionDate = &LastTransactionDates
		}

		ownerName := entity.ConvStatusOwnerId(entity.OwnerId, row.OwnerId)
		vResp.OwnerName = ownerName

		cndnJenisName := entity.ConvCndnJenis(entity.CndnJenis, *row.CndnJenis)
		vResp.CndnJenis = cndnJenisName

		vResp.UsedAmount = float64(0)
		vResp.RemainingAmount = *row.Amount - vResp.UsedAmount

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *CndnServiceImpl) Delete(custId string, CndnNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.CndnRepository.Delete(txCtx, custId, CndnNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *CndnServiceImpl) Update(CndnNo string, request entity.UpdateCndnBody) (err error) {
	c := context.Background()

	if request.CndnDate != nil {
		CndnDate, err := str.DateStrToRfc3339String(*request.CndnDate)
		if err != nil {
			return err
		}
		request.CndnDate = &CndnDate
	}

	if request.LastTransactionDate != nil {
		LastTransactionDate, err := str.DateStrToRfc3339String(*request.LastTransactionDate)
		if err != nil {
			return err
		}
		request.LastTransactionDate = &LastTransactionDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Cndn
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.CndnRepository.Update(txCtx, CndnNo, request.CustID, Model)
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
