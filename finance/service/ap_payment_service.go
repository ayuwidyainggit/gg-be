package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type ApPaymentService interface {
	Store(request entity.CreateAccountPayablePaymentBody) (response entity.CreateAccountPayablePaymentBody, err error)
	Detail(AccountPayablePaymentNo string, custID string, ParentCustId string) (response entity.AccountPayablePaymentDetailResponse, err error)
	List(dataFilter entity.ApPaymentQueryFilter) (data []entity.AccountPayablePaymentList, total int64, lastPage int, err error)
	Delete(custId string, AccountPayablePaymentNo string, userId int64) (err error)
	Update(AccountPayablePaymentNo string, request entity.UpdateAccountPayablePaymentBody) (err error)

	ListBalancePaymentDepositByCustId(dataFilter entity.GeneralQueryFilter) (data []entity.DepositPaymentLookup, total int64, lastPage int, err error)
	ListInvoiceNo(dataFilter entity.ApLookupSupplierInoviceReturnQueryFilter) (data []entity.ApLookupSupplierInvoiceReturnResponeList, total int64, lastPage int, err error)
}

type ApPaymentServiceImpl struct {
	ApPaymentRepository repository.ApPaymentRepository
	Transaction         repository.Dbtransaction
}

func NewApPaymentService(repository repository.ApPaymentRepository, transaction repository.Dbtransaction) *ApPaymentServiceImpl {
	return &ApPaymentServiceImpl{
		ApPaymentRepository: repository,
		Transaction:         transaction,
	}
}

func (service *ApPaymentServiceImpl) Store(request entity.CreateAccountPayablePaymentBody) (response entity.CreateAccountPayablePaymentBody, err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if *request.AccountPayablePaymentDate != "" {
		ApPaymentDate, err := str.DateStrToRfc3339String(*request.AccountPayablePaymentDate)
		if err != nil {
			return response, err
		}
		request.AccountPayablePaymentDate = &ApPaymentDate
	}

	if *request.Details[0].TotalPayment > *request.Details[0].RemainingAmount {
		return response, errors.New("Total Payment tidak boleh lebih besar dari Remaining Amount")
	}

	var apPaymentModel model.AccountPayablePayment
	err = structs.Automapper(request, &apPaymentModel)
	if err != nil {
		return response, err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		err := service.ApPaymentRepository.Store(txCtx, &apPaymentModel)

		if err != nil {
			return err
		}

		response.AccountPayablePaymentNo = apPaymentModel.AccountPayablePaymentNo

		for _, detail := range request.Details {
			var apPaymentDetailModel model.AccountPayablePaymentDetail

			if *detail.InvoiceDate != "" {
				invDate, err := str.DateStrToRfc3339String(*detail.InvoiceDate)
				if err != nil {
					return err
				}
				detail.InvoiceDate = &invDate
			}

			err := structs.Automapper(detail, &apPaymentDetailModel)
			if err != nil {
				return err
			}

			apPaymentDetailModel.AccountPayablePaymentNo = apPaymentModel.AccountPayablePaymentNo
			apPaymentDetailModel.CustId = request.CustId
			_, err = service.ApPaymentRepository.StoreApPaymentDetail(txCtx, &apPaymentDetailModel)
			if err != nil {
				return err
			}

			for _, payment := range detail.Payment {
				var apPaymentOptionsModel model.AccountPayablePaymentOptions

				err := structs.Automapper(payment, &apPaymentOptionsModel)
				if err != nil {
					return err
				}
				apPaymentOptionsModel.AccountPayablePaymentNo = apPaymentModel.AccountPayablePaymentNo
				apPaymentOptionsModel.CustId = request.CustId

				_, err = service.ApPaymentRepository.StoreApPaymentOptions(txCtx, &apPaymentOptionsModel)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	return response, nil
}

func (service *ApPaymentServiceImpl) Detail(AccountPayablePaymentNo string, custID string, ParentCustId string) (response entity.AccountPayablePaymentDetailResponse, err error) {
	ApPayment, err := service.ApPaymentRepository.FindByNo(AccountPayablePaymentNo, custID, ParentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ApPayment, &response)
	if err != nil {
		return response, err
	}

	details, err := service.ApPaymentRepository.FindDetailByNo(AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Details = make([]entity.AccountPayablePaymentDetailRespone, len(details))
	for i, detail := range details {
		var detailData entity.AccountPayablePaymentDetailRespone

		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		var totalPayment float64
		if detail.TotalPayment != nil {
			totalPayment = *detail.TotalPayment
		}

		detailsPayment, err := service.ApPaymentRepository.FindDetailPaymentByNo(AccountPayablePaymentNo, *detailData.InvoiceNo, custID, totalPayment)
		if err != nil {
			return response, err
		}
		// Fallback for multi payment-method lines (options amounts sum to total, none equals total alone)
		if len(detailsPayment) == 0 && totalPayment > 0 {
			detailsPayment, err = service.ApPaymentRepository.FindDetailPaymentByNo(AccountPayablePaymentNo, *detailData.InvoiceNo, custID, 0)
			if err != nil {
				return response, err
			}
		}

		detailData.Payment = make([]entity.AccountPayablePaymentOptionsRespone, len(detailsPayment))
		for i, detailPaymentOptions := range detailsPayment {
			var detailDataPayment entity.AccountPayablePaymentOptionsRespone
			err = structs.Automapper(detailPaymentOptions, &detailDataPayment)
			if err != nil {
				return response, err
			}
			// 1: Cash; 2: Check; 3: Transfer; 4: Retrun; 5:Return
			if *detailPaymentOptions.PayType == 1 {
				detailDataPayment.PayTypeName = "Cash"
			}
			if *detailPaymentOptions.PayType == 2 {
				detailDataPayment.PayTypeName = "Cheque"
			}
			if *detailPaymentOptions.PayType == 3 {
				detailDataPayment.PayTypeName = "Transfer"
			}
			if *detailPaymentOptions.PayType == 4 {
				detailDataPayment.PayTypeName = "Return"
			}
			if *detailPaymentOptions.PayType == 5 {
				detailDataPayment.PayTypeName = "Credit/Debit"
			}

			InvDate := detailPaymentOptions.InvDate.Format("2006-01-02")
			detailDataPayment.InvoiceDate = &InvDate

			detailData.Payment[i] = detailDataPayment
		}

		if detail.InvDate != nil {
			ApPaymentDate := detail.InvDate.Format("2006-01-02")
			detailData.InvoiceDate = &ApPaymentDate
		}

		response.Details[i] = detailData
	}

	if ApPayment.AccountPayablPaymenteDate != nil {
		ApPaymentDate := ApPayment.AccountPayablPaymenteDate.Format("2006-01-02")
		response.AccountPayablePaymentDate = &ApPaymentDate
	}

	cashs, err := service.ApPaymentRepository.FindDetailApPaymentOptionsByNo(1, AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Cash = make([]entity.AccountPayablePaymentOptionsRespone, len(cashs))
	var tmpTotalCash, tmpTotalCek, tmpTotalTransfer, tmpTotalReturn, tmpTotalCndn float64
	for i, cash := range cashs {
		var detailData entity.AccountPayablePaymentOptionsRespone
		err = structs.Automapper(cash, &detailData)
		if err != nil {
			return response, err
		}
		if *detailData.PayType == 1 {
			detailData.PayTypeName = "Cash"
		}

		InvDate := cash.InvDate.Format("2006-01-02")
		detailData.InvoiceDate = &InvDate
		tmpTotalCash += detailData.PaymentAmount

		response.Cash[i] = detailData
	}
	response.TotalCash = tmpTotalCash

	cek, err := service.ApPaymentRepository.FindDetailApPaymentOptionsByNo(2, AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Cek = make([]entity.AccountPayablePaymentOptionsRespone, len(cek))

	for i, ceks := range cek {
		var detailData entity.AccountPayablePaymentOptionsRespone
		err = structs.Automapper(ceks, &detailData)
		if err != nil {
			return response, err
		}
		if *detailData.PayType == 2 {
			detailData.PayTypeName = "Cheque"
		}
		InvDate := ceks.InvDate.Format("2006-01-02")
		detailData.InvoiceDate = &InvDate
		tmpTotalCek += detailData.PaymentAmount

		response.Cek[i] = detailData
	}
	response.TotalCek = tmpTotalCek

	transfers, err := service.ApPaymentRepository.FindDetailApPaymentOptionsByNo(3, AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Trasfer = make([]entity.AccountPayablePaymentOptionsRespone, len(transfers))

	for i, transfer := range transfers {
		var detailData entity.AccountPayablePaymentOptionsRespone
		err = structs.Automapper(transfer, &detailData)
		if err != nil {
			return response, err
		}

		if *detailData.PayType == 3 {
			detailData.PayTypeName = "Transfer"
		}

		InvDate := transfer.InvDate.Format("2006-01-02")
		detailData.InvoiceDate = &InvDate

		tmpTotalTransfer += detailData.PaymentAmount

		response.Trasfer[i] = detailData
	}
	response.TotalTransfer = tmpTotalTransfer

	returns, err := service.ApPaymentRepository.FindDetailApPaymentOptionsByNo(4, AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Return = make([]entity.AccountPayablePaymentOptionsRespone, len(returns))

	for i, returna := range returns {
		var detailData entity.AccountPayablePaymentOptionsRespone
		err = structs.Automapper(returna, &detailData)
		if err != nil {
			return response, err
		}

		if *detailData.PayType == 4 {
			detailData.PayTypeName = "Return"
		}

		InvDate := returna.InvDate.Format("2006-01-02")
		detailData.InvoiceDate = &InvDate

		tmpTotalReturn += detailData.PaymentAmount

		response.Return[i] = detailData
	}
	response.TotalRetrun = tmpTotalReturn

	cndns, err := service.ApPaymentRepository.FindDetailApPaymentOptionsByNo(5, AccountPayablePaymentNo, custID)
	if err != nil {
		return response, err
	}

	response.Cndn = make([]entity.AccountPayablePaymentOptionsRespone, len(cndns))

	for i, cndn := range cndns {
		var detailData entity.AccountPayablePaymentOptionsRespone
		err = structs.Automapper(cndn, &detailData)
		if err != nil {
			return response, err
		}

		if *detailData.PayType == 5 {
			detailData.PayTypeName = "Credit/Debit"
		}

		InvDate := cndn.InvDate.Format("2006-01-02")
		detailData.InvoiceDate = &InvDate

		tmpTotalCndn += detailData.PaymentAmount

		response.Cndn[i] = detailData
	}
	response.TotalCndn = tmpTotalCndn

	return response, nil
}

func (service *ApPaymentServiceImpl) List(dataFilter entity.ApPaymentQueryFilter) (data []entity.AccountPayablePaymentList, total int64, lastPage int, err error) {
	ApPayment, total, lastPage, err := service.ApPaymentRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range ApPayment {
		var vResp entity.AccountPayablePaymentList
		structs.Automapper(row, &vResp)

		if row.AccountPayablPaymenteDate != nil {
			apPaymentDate := row.AccountPayablPaymenteDate.Format("2006-01-02")
			vResp.AccountPayablePaymentDate = &apPaymentDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ApPaymentServiceImpl) Delete(custId string, AccountPayablePaymentNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ApPaymentRepository.Delete(txCtx, custId, AccountPayablePaymentNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ApPaymentServiceImpl) Update(AccountPayablePaymentNo string, request entity.UpdateAccountPayablePaymentBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if *request.AccountPayablePaymentDate != "" {
		ApPaymentDate, err := str.DateStrToRfc3339String(*request.AccountPayablePaymentDate)
		if err != nil {
			return err
		}
		request.AccountPayablePaymentDate = &ApPaymentDate
	}

	var apPaymentModel model.AccountPayablePayment
	err = structs.Automapper(request, &apPaymentModel)
	if err != nil {
		return err
	}
	apPaymentModel.AccountPayablePaymentNo = AccountPayablePaymentNo

	apPaymentModel.CustId = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ApPaymentRepository.DeleteAllDetailByApPayment(txCtx, AccountPayablePaymentNo, request.CustId)
		if err != nil {
			return err
		}

		err = service.ApPaymentRepository.DeleteAllDetailPaymentByApPayment(txCtx, AccountPayablePaymentNo, request.CustId)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var apPaymentDetailModel model.AccountPayablePaymentDetail

			if *detail.InvoiceDate != "" {
				invDate, err := str.DateStrToRfc3339String(*detail.InvoiceDate)
				if err != nil {
					return err
				}
				detail.InvoiceDate = &invDate
			}

			err := structs.Automapper(detail, &apPaymentDetailModel)
			if err != nil {
				return err
			}
			apPaymentDetailModel.AccountPayablePaymentNo = AccountPayablePaymentNo

			apPaymentDetailModel.CustId = request.CustId
			_, err = service.ApPaymentRepository.StoreApPaymentDetail(txCtx, &apPaymentDetailModel)
			if err != nil {
				return err
			}

			for _, payment := range detail.Payment {
				var apPaymentOptionsModel model.AccountPayablePaymentOptions

				err := structs.Automapper(payment, &apPaymentOptionsModel)
				if err != nil {
					return err
				}
				apPaymentOptionsModel.AccountPayablePaymentNo = AccountPayablePaymentNo
				apPaymentOptionsModel.CustId = request.CustId

				_, err = service.ApPaymentRepository.StoreApPaymentOptions(txCtx, &apPaymentOptionsModel)
				if err != nil {
					return err
				}
			}
		}

		err = service.ApPaymentRepository.Update(txCtx, AccountPayablePaymentNo, request.CustId, apPaymentModel)
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

func (service *ApPaymentServiceImpl) ListBalancePaymentDepositByCustId(dataFilter entity.GeneralQueryFilter) (data []entity.DepositPaymentLookup, total int64, lastPage int, err error) {
	DepositLookups, total, lastPage, err := service.ApPaymentRepository.FindAllBalancePaymentDepositByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range DepositLookups {
		var vResp entity.DepositPaymentLookup
		structs.Automapper(row, &vResp)

		if row.Balance != 0 && row.Balance > 0 {
			data = append(data, vResp)
		}
	}

	return data, total, lastPage, err
}

func (service *ApPaymentServiceImpl) ListInvoiceNo(dataFilter entity.ApLookupSupplierInoviceReturnQueryFilter) (data []entity.ApLookupSupplierInvoiceReturnResponeList, total int64, lastPage int, err error) {
	aps, total, lastPage, err := service.ApPaymentRepository.FindAllInvoiceNo(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range aps {
		var vResp entity.ApLookupSupplierInvoiceReturnResponeList
		structs.Automapper(row, &vResp)

		if row.ApType == "I" {
			vResp.ApType = "Invoice"
		} else {
			vResp.ApType = "Return"
		}

		if row.Amount != row.PaidAmount {
			data = append(data, vResp)
		}

	}

	return data, total, lastPage, err
}
