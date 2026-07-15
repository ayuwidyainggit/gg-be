package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
)

type ChequeGiroClearingService interface {
	ListChequeGiroClearing(dataFilter entity.CheckGiroClearingQueryFilter) (data []entity.ChequeGiroClearingResponse, total int64, lastPage int, err error)
	Detail(ChequeGiroNo int, custID string) (response entity.ChequeGiroClearingResponse, err error)
	Update(ChequeGiroNo int, request entity.UpdateChequeGiroClearingBody) (err error)
	UpdateChange(ChequeGiroNo int, request entity.UpdateChequeGiroClearingChangeBody) (err error)
}

type ChequeGiroClearingServiceImpl struct {
	ChequeGiroClearingRepository repository.ChequeGiroClearingRepository
	Transaction                  repository.Dbtransaction
}

func NewChequeGiroClearingService(repository repository.ChequeGiroClearingRepository, transaction repository.Dbtransaction) *ChequeGiroClearingServiceImpl {
	return &ChequeGiroClearingServiceImpl{
		ChequeGiroClearingRepository: repository,
		Transaction:                  transaction,
	}
}

func (service *ChequeGiroClearingServiceImpl) ListChequeGiroClearing(dataFilter entity.CheckGiroClearingQueryFilter) (data []entity.ChequeGiroClearingResponse, total int64, lastPage int, err error) {
	Deposits, total, lastPage, err := service.ChequeGiroClearingRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Deposits {
		var vResp entity.ChequeGiroClearingResponse
		structs.Automapper(row, &vResp)

		if row.StatusCheque != 0 {
			statusText := entity.ConvStatus(entity.StatusGiro, row.StatusCheque)
			vResp.StatusClearing = &statusText
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ChequeGiroClearingServiceImpl) Detail(ChequeGiroNo int, custID string) (response entity.ChequeGiroClearingResponse, err error) {
	ChequeGiro, err := service.ChequeGiroClearingRepository.FindByNo(ChequeGiroNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ChequeGiro, &response)
	if err != nil {
		return response, err
	}

	response.PayingBankName = ChequeGiro.BankName

	ChqDate := ChequeGiro.DocDateCheque.Format("2006-01-02")
	response.DocDateCheque = &ChqDate

	ChqDueDate := ChequeGiro.DueDate.Format("2006-01-02")
	response.DueDate = &ChqDueDate

	ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	response.OwnerName = ownerName

	statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	response.StatusClearing = &statusText

	// response.UsedAmount = float64(0)
	response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *ChequeGiroClearingServiceImpl) Update(ChequeGiroNo int, request entity.UpdateChequeGiroClearingBody) (err error) {
	c := context.Background()

	if request.ClearingDate != "" {
		ChqDate, err := str.DateStrToRfc3339String(request.ClearingDate)
		if err != nil {
			return err
		}
		request.ClearingDate = ChqDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ChequeGiroClearingList
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	// Model.CustID = ""

	Model.StatusCheque = 3

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeGiroClearingRepository.UpdateClearing(txCtx, ChequeGiroNo, request.CustID, Model)
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

func (service *ChequeGiroClearingServiceImpl) UpdateChange(ChequeGiroNo int, request entity.UpdateChequeGiroClearingChangeBody) (err error) {
	c := context.Background()

	if request.ClearingDate != "" {
		ChqDate, err := str.DateStrToRfc3339String(request.ClearingDate)
		if err != nil {
			return err
		}
		request.ClearingDate = ChqDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.ChequeGiroClearingList
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}

	// Model.CustID = ""

	Model.StatusCheque = 1

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.ChequeGiroClearingRepository.UpdateClearing(txCtx, ChequeGiroNo, request.CustID, Model)
		if err != nil {
			return err
		}

		var replacedID []int

		Cekpayment, err := service.ChequeGiroClearingRepository.FindDetailPaymentDepositByGiroNo(request.DocNoCheque, request.CustID)
		if err != nil {
			return err
		}

		request.CashRemaining = request.CashAmount
		request.ChequeRemaining = request.ChequeAmount
		request.TransferRemaining = request.TransferAmount

		for _, payment := range Cekpayment {
			replacedID = append(replacedID, payment.DepositPaymentID)

			depositPayment := model.DepositPayment{
				DepositNo: payment.DepositNo,
				InvoiceNo: payment.InvoiceNo,
			}

			remainingAmount := payment.PaymentAmount

			diskon := payment.Discount
			materai := payment.Materai

			// Process cash
			if request.CashRemaining > 0 {
				if remainingAmount <= request.CashRemaining {
					depositPayment.PayType = 1
					depositPayment.DocumentNo = request.CashNo
					depositPayment.Balance = request.CashRemaining
					depositPayment.PaymentAmount = remainingAmount

					request.CashRemaining -= remainingAmount
					remainingAmount = 0
				} else {
					depositPayment.PayType = 1
					depositPayment.DocumentNo = request.CashNo
					depositPayment.Balance = request.CashRemaining
					depositPayment.PaymentAmount = request.CashRemaining

					remainingAmount -= request.CashRemaining
					request.CashRemaining = 0
				}

				cashPayment, err := service.ChequeGiroClearingRepository.FindDetailPaymentCashByDepositInvoiceNo(1, payment.DepositNo, payment.InvoiceNo, request.CustID)
				if err != nil {
					return err
				}

				if len(cashPayment) > 0 {
					depositPaymentCash := model.DepositPayment{
						DocumentNo:    depositPayment.DocumentNo,
						Balance:       cashPayment[0].Balance + depositPayment.PaymentAmount,
						PaymentAmount: cashPayment[0].Balance + depositPayment.PaymentAmount,
					}

					if materai != nil && *materai > 0 {
						depositPaymentCash.Materai = materai
						materai = nil
					}
					if diskon != nil && *diskon > 0 {
						depositPaymentCash.Discount = diskon
						diskon = nil
					}

					if err := service.ChequeGiroClearingRepository.UpdateCashAmount(txCtx, payment.DepositPaymentID, depositPaymentCash); err != nil {
						return err
					}
				} else {

					if materai != nil && *materai > 0 {
						depositPayment.Materai = materai
						materai = nil
					}
					if diskon != nil && *diskon > 0 {
						depositPayment.Discount = diskon
						diskon = nil
					}

					if _, err := service.ChequeGiroClearingRepository.StorePayment(txCtx, &depositPayment); err != nil {
						return err
					}
				}

				if remainingAmount == 0 {
					continue
				}
			}

			// Process transfer
			if request.TransferRemaining > 0 {
				if remainingAmount <= request.TransferRemaining {
					depositPayment.PayType = 3
					depositPayment.DocumentNo = request.TransferNo
					depositPayment.Balance = request.TransferRemaining
					depositPayment.PaymentAmount = remainingAmount

					request.TransferRemaining -= remainingAmount
					remainingAmount = 0
				} else {
					depositPayment.PayType = 3
					depositPayment.DocumentNo = request.TransferNo
					depositPayment.Balance = request.TransferRemaining
					depositPayment.PaymentAmount = request.TransferRemaining

					remainingAmount -= request.TransferRemaining
					request.TransferRemaining = 0
				}

				if materai != nil && *materai > 0 {
					depositPayment.Materai = materai
					materai = nil
				}
				if diskon != nil && *diskon > 0 {
					depositPayment.Discount = diskon
					diskon = nil
				}

				if _, err := service.ChequeGiroClearingRepository.StorePayment(txCtx, &depositPayment); err != nil {
					return err
				}
				if remainingAmount == 0 {
					continue
				}
			}

			// Process cheque
			if request.ChequeRemaining > 0 {
				if remainingAmount <= request.ChequeRemaining {
					depositPayment.PayType = 2
					depositPayment.DocumentNo = request.ChequeNo
					depositPayment.Balance = request.ChequeRemaining
					depositPayment.PaymentAmount = remainingAmount

					request.ChequeRemaining -= remainingAmount
					remainingAmount = 0
				} else {
					depositPayment.PayType = 2
					depositPayment.DocumentNo = request.ChequeNo
					depositPayment.Balance = request.ChequeRemaining
					depositPayment.PaymentAmount = request.ChequeRemaining

					remainingAmount -= request.ChequeRemaining
					request.ChequeRemaining = 0
				}

				if materai != nil && *materai > 0 {
					depositPayment.Materai = materai
					materai = nil
				}
				if diskon != nil && *diskon > 0 {
					depositPayment.Discount = diskon
					diskon = nil
				}

				if _, err := service.ChequeGiroClearingRepository.StorePayment(txCtx, &depositPayment); err != nil {
					return err
				}
				if remainingAmount == 0 {
					continue
				}
			}
		}

		if len(replacedID) > 0 {
			err = service.ChequeGiroClearingRepository.DeleteAllDetailPaymentByDepositPaymentID(txCtx, replacedID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
