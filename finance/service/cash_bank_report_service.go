package service

import (
	"finance/entity"
	"finance/pkg/structs"
	"finance/repository"
)

type CashBankReportService interface {
	ListReport(dataFilter entity.DepositQueryFilter) (data []entity.DepositReportResponse, total int64, lastPage int, err error)
	LookupFilterReportDepositNo(dataFilter entity.GeneralQueryFilter) (data []entity.DepositNoLookup, total int64, lastPage int, err error)
	LookupFilterReportPayType(dataFilter entity.GeneralQueryFilter) (data []entity.DepositPayTypeLookup, total int64, lastPage int, err error)
}

type CashBankReportServiceImpl struct {
	CashBankReportRepository repository.CashBankReportRepository
	Transaction              repository.Dbtransaction
}

func NewCashBankReportService(repository repository.CashBankReportRepository, transaction repository.Dbtransaction) *CashBankReportServiceImpl {
	return &CashBankReportServiceImpl{
		CashBankReportRepository: repository,
		Transaction:              transaction,
	}
}

func (service *CashBankReportServiceImpl) ListReport(dataFilter entity.DepositQueryFilter) (data []entity.DepositReportResponse, total int64, lastPage int, err error) {
	Deposits, total, lastPage, err := service.CashBankReportRepository.FindAllReportDeposit(dataFilter, dataFilter.CustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Deposits {
		var vResp entity.DepositReportResponse
		structs.Automapper(row, &vResp)

		// if row.DepositDate != nil {
		// 	DepositDate := row.DepositDate.Format("2006-01-02")
		// 	vResp.DepositDate = &DepositDate
		// }

		if row.StatusCheque != nil {
			statusText := entity.ConvStatus(entity.StatusGiro, *row.StatusCheque)
			vResp.StatusClearing = &statusText
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *CashBankReportServiceImpl) LookupFilterReportDepositNo(dataFilter entity.GeneralQueryFilter) (data []entity.DepositNoLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.CashBankReportRepository.FindAllDepositNoReportFilterByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.DepositNoLookup
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *CashBankReportServiceImpl) LookupFilterReportPayType(dataFilter entity.GeneralQueryFilter) (data []entity.DepositPayTypeLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.CashBankReportRepository.FindAllDepositPaymentTypeReportFilterByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.DepositPayTypeLookup
		structs.Automapper(row, &vResp)
		PayTypeName := entity.ConvStatusBankTransfer(entity.PayType, row.PayType)
		vResp.PayTypeName = PayTypeName
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}
