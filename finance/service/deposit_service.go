package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/shopspring/decimal"
)

type DepositService interface {
	StoreCollection(request entity.CreateDepositBodyByCollection) (err error)
	StoreInvoice(request entity.CreateDepositBodyByInvoice) (err error)
	Detail(depositNo string, custID string) (response entity.DepositDetailResponse, err error)
	DetailReport(depositNo string, custID string) (response entity.DepositDetailReportResponse, err error)
	List(dataFilter entity.DepositQueryFilter) (data []entity.DepositResponse, total int64, lastPage int, err error)
	ListDepositNumber(dataFilter entity.DepositNumberListQueryFilter) (data []entity.DepositNumberListItemResponse, total int64, lastPage int, err error)
	Delete(custId string, depositNo string, userId int64) (err error)
	UpdateCollection(depositNo string, request entity.UpdateDepositBodyCollection) (err error)
	UpdateInvoice(depositNo string, request entity.UpdateDepositBodyInvoice) (err error)
	ProofOfPayment(depositNo string, q string, typ string, custID string) (items []map[string]interface{}, err error)
}

type DepositServiceImpl struct {
	DepositRepository repository.DepositRepository
	Transaction       repository.Dbtransaction
}

func NewDepositService(repository repository.DepositRepository, transaction repository.Dbtransaction) *DepositServiceImpl {
	return &DepositServiceImpl{
		DepositRepository: repository,
		Transaction:       transaction,
	}
}

func (service *DepositServiceImpl) StoreCollection(request entity.CreateDepositBodyByCollection) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DepositDate != "" {
		DepositDate, err := str.DateStrToRfc3339String(request.DepositDate)
		if err != nil {
			return err
		}
		request.DepositDate = DepositDate
	}

	var depositModel model.Deposit
	err = structs.Automapper(request, &depositModel)
	if err != nil {
		return err
	}

	total, err := service.DepositRepository.CountAllByCustId(request.CustID, request.DepositDate)
	if err != nil {
		return err
	}

	depositNo := entity.GenerateNumber("DP", total, depositModel.DepositDate)
	depositModel.DepositNo = depositNo

	fmt.Println(depositNo)

	if depositModel.CollectionNo != nil {
		depositModel.SalesmanID = nil
		depositModel.InvoiceDateFrom = nil
		depositModel.InvoiceDateTo = nil
		depositModel.DueDateTo = nil
		depositModel.DueDateTo = nil
	}

	depositModel.DepositStatus = 1

	var remainingAmount = 0.0

	var auditor int64
	if request.CreatedBy != nil {
		auditor = *request.CreatedBy
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		depositModel.RemainingAmount = 0
		err = service.DepositRepository.Store(txCtx, &depositModel)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var depositDetails model.DepositDetail

			err := structs.Automapper(detail, &depositDetails)
			if err != nil {
				return err
			}

			depositDetails.DepositNo = depositNo
			depositDetails.RemainingPayment = detail.RemainingPayment
			remainingAmount += (depositDetails.RemainingPayment - depositDetails.TotalPayment)
			depositDetails.CustID = request.CustID
			_, err = service.DepositRepository.StoreDetail(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			err = service.DepositRepository.CalcCollectionPaidByInvoice(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			if depositDetails.TotalPayment > 0 {
				indexGiro := -1
				for i, payment := range detail.Payment {
					if payment.PayType == 2 {
						indexGiro = i
						break
					}
				}

				for index, payment := range detail.Payment {
					var depositPayment model.DepositPayment

					err := structs.Automapper(payment, &depositPayment)
					if err != nil {
						return err
					}
					depositPayment.DepositNo = depositNo
					depositPayment.CustID = request.CustID

					if index == 0 {
						depositPayment.Discount = &detail.Discount
					}

					if indexGiro > -1 && indexGiro == index {
						depositPayment.Materai = &detail.Materai
					} else {
						zeroFloat := float64(0)
						depositPayment.Materai = &zeroFloat
					}

					switch depositPayment.PayType {
					case 2: // cheque_giro
						err := service.DepositRepository.UpdateAmountProgressionCheque(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 3: // bank_transfer
						err := service.DepositRepository.UpdateAmountProgressionTransfer(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 4: // return
						err := service.DepositRepository.UpdateAmountProgressionReturn(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 5: // cndn
						err := service.DepositRepository.UpdateAmountProgressionCNDN(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					}

					_, err = service.DepositRepository.StorePayment(txCtx, &depositPayment)
					if err != nil {
						return err
					}
				}
			}

		}

		// Process expenses at top level
		for _, expense := range request.Expense {
			var depositExpense model.DepositExpense

			err := structs.Automapper(expense, &depositExpense)
			if err != nil {
				return err
			}

			depositExpense.DepositNo = depositNo
			depositExpense.CustID = request.CustID
			depositExpense.PaymentAmount = decimal.NewFromFloat(expense.Amount)
			if request.CreatedBy != nil {
				depositExpense.CreatedBy = *request.CreatedBy
			}

			_, err = service.DepositRepository.StoreExpense(txCtx, &depositExpense)
			if err != nil {
				return err
			}

			if err := service.DepositRepository.DeductExpense(txCtx, &depositExpense); err != nil {
				return err
			}
		}

		depositModel.RemainingAmount = remainingAmount
		err := service.DepositRepository.Update(txCtx, depositNo, request.CustID, depositModel)
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

func (service *DepositServiceImpl) StoreInvoice(request entity.CreateDepositBodyByInvoice) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DepositDate != "" {
		DepositDate, err := str.DateStrToRfc3339String(request.DepositDate)
		if err != nil {
			return err
		}
		request.DepositDate = DepositDate
	}

	if request.InvoiceDateFrom != "" {
		InvoiceDateFrom, err := str.DateStrToRfc3339String(request.InvoiceDateFrom)
		if err != nil {
			return err
		}
		request.InvoiceDateFrom = InvoiceDateFrom
	}

	if request.InvoiceDateTo != "" {
		InvoiceDateTo, err := str.DateStrToRfc3339String(request.InvoiceDateTo)
		if err != nil {
			return err
		}
		request.InvoiceDateTo = InvoiceDateTo
	}

	if request.DueDateFrom != "" {
		DueDateFrom, err := str.DateStrToRfc3339String(request.DueDateFrom)
		if err != nil {
			return err
		}
		request.DueDateFrom = DueDateFrom
	}

	if request.DueDateTo != "" {
		DueDateTo, err := str.DateStrToRfc3339String(request.DueDateTo)
		if err != nil {
			return err
		}
		request.DueDateTo = DueDateTo
	}

	var depositModel model.Deposit
	err = structs.Automapper(request, &depositModel)
	if err != nil {
		return err
	}

	total, err := service.DepositRepository.CountAllByCustId(request.CustID, request.DepositDate)
	if err != nil {
		return err
	}

	depositNo := entity.GenerateNumber("DP", total, depositModel.DepositDate)
	depositModel.DepositNo = depositNo

	if depositModel.CollectionNo != nil {
		depositModel.SalesmanID = nil
		depositModel.InvoiceDateFrom = nil
		depositModel.InvoiceDateTo = nil
		depositModel.DueDateTo = nil
		depositModel.DueDateTo = nil
	}

	depositModel.DepositStatus = 1

	var remainingAmount = 0.0

	var auditor int64
	if request.CreatedBy != nil {
		auditor = *request.CreatedBy
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		depositModel.RemainingAmount = 0
		err = service.DepositRepository.Store(txCtx, &depositModel)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var depositDetails model.DepositDetail

			err := structs.Automapper(detail, &depositDetails)
			if err != nil {
				return err
			}

			depositDetails.DepositNo = depositNo
			remainingAmountByInv, errR := service.DepositRepository.CountRemainingAmountByInvoice(txCtx, detail.InvoiceNo, request.CustID)
			if errR != nil {
				return errR
			}

			depositDetails.RemainingPayment = depositDetails.InvoiceAmount - remainingAmountByInv
			// remainingAmount += depositDetails.RemainingPayment
			remainingAmount += (depositDetails.RemainingPayment - depositDetails.TotalPayment)
			depositDetails.CustID = request.CustID
			_, err = service.DepositRepository.StoreDetail(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			err = service.DepositRepository.CalcCollectionPaidByInvoice(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			if depositDetails.TotalPayment > 0 {
				indexGiro := -1
				for i, payment := range detail.Payment {
					if payment.PayType == 2 {
						indexGiro = i
						break
					}
				}
				indexDiscount := -1
				indexDiscountPayType := 0
				for i, payment := range detail.Payment {
					if payment.PayType == 1 {
						indexDiscount = i
						indexDiscountPayType = 1
						break
					} else if payment.PayType == 2 {
						indexDiscount = i
						indexDiscountPayType = 2
						continue
					} else if payment.PayType == 3 {
						if indexDiscountPayType > 2 || indexDiscountPayType == 0 {
							indexDiscount = i
							indexDiscountPayType = 3
						}
					} else {
						if indexDiscountPayType > 3 || indexDiscountPayType == 0 {
							indexDiscount = 0
						}
					}
				}

				for index, payment := range detail.Payment {
					var depositPayment model.DepositPayment

					err := structs.Automapper(payment, &depositPayment)
					if err != nil {
						return err
					}
					depositPayment.DepositNo = depositNo
					depositPayment.CustID = request.CustID

					// if index == 0 {
					// 	depositPayment.Discount = &detail.Discount
					// }

					if index == indexDiscount {
						depositPayment.Discount = &detail.Discount
					}

					if indexGiro > -1 && indexGiro == index {
						depositPayment.Materai = &detail.Materai
					} else {
						zeroFloat := float64(0)
						depositPayment.Materai = &zeroFloat
					}

					switch depositPayment.PayType {
					case 2: // cheque_giro
						err := service.DepositRepository.UpdateAmountProgressionCheque(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 3: // bank_transfer
						err := service.DepositRepository.UpdateAmountProgressionTransfer(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 4: // return
						err := service.DepositRepository.UpdateAmountProgressionReturn(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 5: // cndn
						err := service.DepositRepository.UpdateAmountProgressionCNDN(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					}

					_, err = service.DepositRepository.StorePayment(txCtx, &depositPayment)
					if err != nil {
						return err
					}
				}
			}
		}

		for _, expense := range request.Expense {
			var depositExpense model.DepositExpense

			err := structs.Automapper(expense, &depositExpense)
			if err != nil {
				return err
			}

			depositExpense.DepositNo = depositNo
			depositExpense.CustID = request.CustID
			depositExpense.PaymentAmount = decimal.NewFromFloat(expense.Amount)
			if request.CreatedBy != nil {
				depositExpense.CreatedBy = *request.CreatedBy
			}

			_, err = service.DepositRepository.StoreExpense(txCtx, &depositExpense)
			if err != nil {
				return err
			}

			if err := service.DepositRepository.DeductExpense(txCtx, &depositExpense); err != nil {
				return err
			}
		}

		depositModel.RemainingAmount = remainingAmount
		err := service.DepositRepository.Update(txCtx, depositNo, request.CustID, depositModel)

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

func (service *DepositServiceImpl) ProofOfPayment(depositNo string, q string, typ string, custID string) (items []map[string]interface{}, err error) {
	return service.DepositRepository.FindProofOfPayment(depositNo, q, typ, custID)
}

func (service *DepositServiceImpl) Detail(depositNo string, custID string) (response entity.DepositDetailResponse, err error) {
	Deposit, err := service.DepositRepository.FindByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Deposit, &response)
	if err != nil {
		return response, err
	}

	statusText := entity.ConvStatus(entity.StatusDeposit, response.DepositStatus)
	response.DepositStatusName = statusText

	details, err := service.DepositRepository.FindDetailByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	finalRemainingAmount := 0.0
	response.Details = make([]entity.DepositDetail, len(details))
	for i, detail := range details {
		var detailData entity.DepositDetail
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		detailData.RemainingPayment = detail.RemainingAmount
		finalRemainingAmount += detailData.RemainingPayment

		detailsPayment, err := service.DepositRepository.FindDetailPaymentByNo(depositNo, detailData.InvoiceNo, custID)
		if err != nil {
			return response, err
		}

		detailData.Payment = make([]entity.DepositPayment, len(detailsPayment))
		for i, detail := range detailsPayment {
			var detailDataPayment entity.DepositPayment
			err = structs.Automapper(detail, &detailDataPayment)
			if err != nil {
				return response, err
			}
			detailData.Payment[i] = detailDataPayment
		}

		response.Details[i] = detailData
	}

	response.RemainingAmount = finalRemainingAmount

	if Deposit.DepositDate != nil {
		DepositDate := Deposit.DepositDate.Format("2006-01-02")
		response.DepositDate = &DepositDate
	}

	if Deposit.CollectionDate != nil {
		CollectionDate := Deposit.CollectionDate.Format("2006-01-02")
		response.CollectionDate = &CollectionDate
	}

	if Deposit.InvoiceDateFrom != nil {
		InvoiceDateFrom := Deposit.InvoiceDateFrom.Format("2006-01-02")
		response.InvoiceDateFrom = &InvoiceDateFrom
	}

	if Deposit.InvoiceDateTo != nil {
		InvoiceDateTo := Deposit.InvoiceDateTo.Format("2006-01-02")
		response.InvoiceDateTo = &InvoiceDateTo
	}

	if Deposit.DueDateFrom != nil {
		DueDateFrom := Deposit.DueDateFrom.Format("2006-01-02")
		response.DueDateFrom = &DueDateFrom
	}

	if Deposit.DueDateTo != nil {
		DueDateTo := Deposit.DueDateTo.Format("2006-01-02")
		response.DueDateTo = &DueDateTo
	}

	cashs, err := service.DepositRepository.FindDetailPaymentInvoiceByNo(1, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Cash = make([]entity.DepositPaymentInvoice, len(cashs))

	for i, cash := range cashs {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cash, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Cash[i] = detailData
	}

	cek, err := service.DepositRepository.FindDetailPaymentInvoiceByNo(2, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Cek = make([]entity.DepositPaymentInvoice, len(cek))

	for i, ceks := range cek {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(ceks, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Cek[i] = detailData
	}

	transfers, err := service.DepositRepository.FindDetailPaymentInvoiceByNo(3, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Trasfer = make([]entity.DepositPaymentInvoice, len(transfers))

	for i, transfer := range transfers {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(transfer, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Trasfer[i] = detailData
	}

	returns, err := service.DepositRepository.FindDetailPaymentInvoiceByNo(4, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Return = make([]entity.DepositPaymentInvoice, len(returns))

	for i, returna := range returns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(returna, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.Return[i] = detailData
	}

	cndns, err := service.DepositRepository.FindDetailPaymentInvoiceByNo(5, depositNo, custID)
	if err != nil {
		return response, err
	}

	response.CNDN = make([]entity.DepositPaymentInvoice, len(cndns))

	for i, cndn := range cndns {
		var detailData entity.DepositPaymentInvoice
		err = structs.Automapper(cndn, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		detailData.PaymentBalance = Deposit.TotalPaymentBalance

		response.CNDN[i] = detailData
	}

	expenses, err := service.DepositRepository.FindExpenseByDepositNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Expense = make([]entity.DepositExpense, len(expenses))
	expenseTotal := 0.0
	for i, expense := range expenses {
		var expenseData entity.DepositExpense

		_ = structs.Automapper(expense, &expenseData)

		paymentAmount, _ := expense.PaymentAmount.Float64()
		expenseData.PaymentAmount = paymentAmount
		expenseTotal += paymentAmount

		response.Expense[i] = expenseData
	}
	response.ExpenseTotal = expenseTotal

	// ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	// response.OwnerName = ownerName

	// statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	// response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	// response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *DepositServiceImpl) DetailReport(depositNo string, custID string) (response entity.DepositDetailReportResponse, err error) {
	Deposit, err := service.DepositRepository.FindByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(Deposit, &response)
	if err != nil {
		return response, err
	}

	details, err := service.DepositRepository.FindDetailByNo(depositNo, custID)
	if err != nil {
		return response, err
	}

	response.Details = make([]entity.DepositDetailReport, len(details))
	for i, detail := range details {
		var detailData entity.DepositDetailReport
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		nilaiTunai := 0.0
		nilaiTransfer := 0.0
		nilaiCek := 0.0
		nilaiRetur := 0.0
		nilaiCndn := 0.0

		detailsPayment, err := service.DepositRepository.FindDetailPaymentByNo(depositNo, detailData.InvoiceNo, custID)
		if err != nil {
			return response, err
		}

		for _, detailp := range detailsPayment {
			switch detailp.PayType {
			case 1:
				nilaiTunai += detailp.PaymentAmount
				detailData.NoTunai = &detailp.DocumentNo
			case 2:
				nilaiCek += detailp.PaymentAmount
				detailData.NoCek = &detailp.DocumentNo
			case 3:
				nilaiTransfer += detailp.PaymentAmount
				detailData.NoTransfer = &detailp.DocumentNo
			case 4:
				nilaiCndn += detailp.PaymentAmount
				detailData.NoCndn = &detailp.DocumentNo
			case 5:
				nilaiRetur += detailp.PaymentAmount
				detailData.NoRetur = &detailp.DocumentNo
			}
		}

		detailData.NilaiTunai = &nilaiTunai
		detailData.NilaiTransfer = &nilaiTransfer
		detailData.NilaiCek = &nilaiCek
		detailData.NilaiCndn = &nilaiCndn
		detailData.NilaiRetur = &nilaiRetur

		response.Details[i] = detailData
	}

	if Deposit.DepositDate != nil {
		DepositDate := Deposit.DepositDate.Format("2006-01-02")
		response.DepositDate = &DepositDate
	}

	if Deposit.CollectionDate != nil {
		CollectionDate := Deposit.CollectionDate.Format("2006-01-02")
		response.CollectionDate = &CollectionDate
	}

	if Deposit.InvoiceDateFrom != nil {
		InvoiceDateFrom := Deposit.InvoiceDateFrom.Format("2006-01-02")
		response.InvoiceDateFrom = &InvoiceDateFrom
	}

	if Deposit.InvoiceDateTo != nil {
		InvoiceDateTo := Deposit.InvoiceDateTo.Format("2006-01-02")
		response.InvoiceDateTo = &InvoiceDateTo
	}

	if Deposit.DueDateFrom != nil {
		DueDateFrom := Deposit.DueDateFrom.Format("2006-01-02")
		response.DueDateFrom = &DueDateFrom
	}

	if Deposit.DueDateTo != nil {
		DueDateTo := Deposit.DueDateTo.Format("2006-01-02")
		response.DueDateTo = &DueDateTo
	}

	// ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	// response.OwnerName = ownerName

	// statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	// response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	// response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}

func (service *DepositServiceImpl) List(dataFilter entity.DepositQueryFilter) (data []entity.DepositResponse, total int64, lastPage int, err error) {
	Deposits, total, lastPage, err := service.DepositRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Deposits {
		var vResp entity.DepositResponse
		structs.Automapper(row, &vResp)

		if row.DepositDate != nil {
			DepositDate := row.DepositDate.Format("2006-01-02")
			vResp.DepositDate = &DepositDate
		}

		if row.CollectionDate != nil {
			CollectionDate := row.CollectionDate.Format("2006-01-02")
			vResp.CollectionDate = &CollectionDate
		}

		if row.InvoiceDateFrom != nil {
			InvoiceDateFrom := row.InvoiceDateFrom.Format("2006-01-02")
			vResp.InvoiceDateFrom = &InvoiceDateFrom
		}

		if row.InvoiceDateTo != nil {
			InvoiceDateTo := row.InvoiceDateTo.Format("2006-01-02")
			vResp.InvoiceDateTo = &InvoiceDateTo
		}

		if row.DueDateFrom != nil {
			DueDateFrom := row.DueDateFrom.Format("2006-01-02")
			vResp.DueDateFrom = &DueDateFrom
		}

		if row.DueDateTo != nil {
			DueDateTo := row.DueDateTo.Format("2006-01-02")
			vResp.DueDateTo = &DueDateTo
		}

		statusText := entity.ConvStatus(entity.StatusDeposit, row.DepositStatus)
		vResp.DepositStatusName = statusText

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *DepositServiceImpl) ListDepositNumber(dataFilter entity.DepositNumberListQueryFilter) (data []entity.DepositNumberListItemResponse, total int64, lastPage int, err error) {
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit < 1 {
		dataFilter.Limit = 20
	}
	if dataFilter.Limit > 9999 {
		dataFilter.Limit = 9999
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	deposits, total, lastPage, err := service.DepositRepository.FindDepositNumberListByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range deposits {
		item := entity.DepositNumberListItemResponse{
			DepositNo:   row.DepositNo,
			CollectorID: row.CollectorID,
		}

		if row.DepositDate != nil {
			item.DepositDate = row.DepositDate.Format("2006-01-02T15:04:05Z")
		}

		data = append(data, item)
	}

	return data, total, lastPage, nil
}

func (service *DepositServiceImpl) Delete(custId string, depositNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.DepositRepository.Delete(txCtx, custId, depositNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *DepositServiceImpl) UpdateCollection(depositNo string, request entity.UpdateDepositBodyCollection) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DepositDate != "" {
		DepositDate, err := str.DateStrToRfc3339String(request.DepositDate)
		if err != nil {
			return err
		}
		request.DepositDate = DepositDate
	}

	var depositModel model.Deposit
	err = structs.Automapper(request, &depositModel)
	if err != nil {
		return err
	}

	depositModel.DepositNo = depositNo

	if depositModel.CollectionNo != nil {
		depositModel.SalesmanID = nil
		depositModel.InvoiceDateFrom = nil
		depositModel.InvoiceDateTo = nil
		depositModel.DueDateTo = nil
		depositModel.DueDateTo = nil
	}

	depositModel.DepositStatus = 1

	remainingAmount := 0.0

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		// restore expense balances by adding back any deposit_expense.payment_amount
		if err := service.DepositRepository.RestoreExpensesByDeposit(txCtx, depositNo, request.CustID); err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllDetailByDeposit(txCtx, depositNo)
		if err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllDetailPaymentByDeposit(txCtx, depositNo)
		if err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllExpenseByDeposit(txCtx, depositNo, request.CustID)
		if err != nil {
			return err
		}

		var auditor int64
		if request.UpdatedBy != nil {
			auditor = *request.UpdatedBy
		}

		for _, detail := range request.Detail {
			var depositDetails model.DepositDetail

			err := structs.Automapper(detail, &depositDetails)
			if err != nil {
				return err
			}

			depositDetails.DepositNo = depositNo
			// for collection update, RemainingPayment comes from request.detail
			depositDetails.RemainingPayment = detail.RemainingPayment

			remainingAmount += (depositDetails.RemainingPayment - depositDetails.TotalPayment)

			depositDetails.CustID = request.CustID
			_, err = service.DepositRepository.StoreDetail(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			err = service.DepositRepository.CalcCollectionPaidByInvoice(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			if depositDetails.TotalPayment > 0 {
				indexGiro := -1
				for i, payment := range detail.Payment {
					if payment.PayType == 2 {
						indexGiro = i
						break
					}
				}

				for index, payment := range detail.Payment {
					var depositPayment model.DepositPayment

					err := structs.Automapper(payment, &depositPayment)
					if err != nil {
						return err
					}
					depositPayment.DepositNo = depositNo
					depositPayment.CustID = request.CustID

					if index == 0 {
						depositPayment.Discount = &detail.Discount
					}

					if indexGiro > -1 && indexGiro == index {
						depositPayment.Materai = &detail.Materai
					} else {
						zeroFloat := float64(0)
						depositPayment.Materai = &zeroFloat
					}

					switch depositPayment.PayType {
					case 2:
						err := service.DepositRepository.UpdateAmountProgressionCheque(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 3:
						err := service.DepositRepository.UpdateAmountProgressionTransfer(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 4:
						err := service.DepositRepository.UpdateAmountProgressionReturn(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 5:
						err := service.DepositRepository.UpdateAmountProgressionCNDN(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					}

					_, err = service.DepositRepository.StorePayment(txCtx, &depositPayment)
					if err != nil {
						return err
					}
				}
			}

		}

		// Process expenses at top level
		for _, expense := range request.Expense {
			var depositExpense model.DepositExpense

			err := structs.Automapper(expense, &depositExpense)
			if err != nil {
				return err
			}

			depositExpense.DepositNo = depositNo
			depositExpense.CustID = request.CustID
			depositExpense.PaymentAmount = decimal.NewFromFloat(expense.Amount)
			if request.UpdatedBy != nil {
				depositExpense.CreatedBy = *request.UpdatedBy
			}

			_, err = service.DepositRepository.StoreExpense(txCtx, &depositExpense)
			if err != nil {
				return err
			}

			if err := service.DepositRepository.DeductExpense(txCtx, &depositExpense); err != nil {
				return err
			}
		}

		depositModel.RemainingAmount = remainingAmount
		err = service.DepositRepository.Update(txCtx, depositNo, request.CustID, depositModel)
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

func (service *DepositServiceImpl) UpdateInvoice(depositNo string, request entity.UpdateDepositBodyInvoice) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.DepositDate != "" {
		DepositDate, err := str.DateStrToRfc3339String(request.DepositDate)
		if err != nil {
			return err
		}
		request.DepositDate = DepositDate
	}

	if request.InvoiceDateFrom != "" {
		InvoiceDateFrom, err := str.DateStrToRfc3339String(request.InvoiceDateFrom)
		if err != nil {
			return err
		}
		request.InvoiceDateFrom = InvoiceDateFrom
	}

	if request.InvoiceDateTo != "" {
		InvoiceDateTo, err := str.DateStrToRfc3339String(request.InvoiceDateTo)
		if err != nil {
			return err
		}
		request.InvoiceDateTo = InvoiceDateTo
	}

	if request.DueDateFrom != "" {
		DueDateFrom, err := str.DateStrToRfc3339String(request.DueDateFrom)
		if err != nil {
			return err
		}
		request.DueDateFrom = DueDateFrom
	}

	if request.DueDateTo != "" {
		DueDateTo, err := str.DateStrToRfc3339String(request.DueDateTo)
		if err != nil {
			return err
		}
		request.DueDateTo = DueDateTo
	}

	var depositModel model.Deposit
	err = structs.Automapper(request, &depositModel)
	if err != nil {
		log.Error(err)
		// return err
	}

	depositModel.DepositNo = depositNo

	if depositModel.CollectionNo != nil {
		depositModel.SalesmanID = nil
		depositModel.InvoiceDateFrom = nil
		depositModel.InvoiceDateTo = nil
		depositModel.DueDateTo = nil
		depositModel.DueDateTo = nil
	}

	remainingAmount := 0.0

	// Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		// restore expense balances by adding back any deposit_expense.payment_amount
		if err := service.DepositRepository.RestoreExpensesByDeposit(txCtx, depositNo, request.CustID); err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllDetailByDeposit(txCtx, depositNo)
		if err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllDetailPaymentByDeposit(txCtx, depositNo)
		if err != nil {
			return err
		}

		err = service.DepositRepository.DeleteAllExpenseByDeposit(txCtx, depositNo, request.CustID)
		if err != nil {
			return err
		}

		var auditor int64
		if request.UpdatedBy != nil {
			auditor = *request.UpdatedBy
		}

		for _, detail := range request.Detail {
			var depositDetails model.DepositDetail

			err := structs.Automapper(detail, &depositDetails)
			if err != nil {
				return err
			}

			depositDetails.DepositNo = depositNo
			remainingAmountByInv, errR := service.DepositRepository.CountRemainingAmountByInvoice(txCtx, detail.InvoiceNo, request.CustID)
			if errR != nil {
				return errR
			}

			depositDetails.RemainingPayment = depositDetails.InvoiceAmount - remainingAmountByInv

			remainingAmount += (depositDetails.RemainingPayment - depositDetails.TotalPayment)

			depositDetails.CustID = request.CustID
			_, err = service.DepositRepository.StoreDetail(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			err = service.DepositRepository.CalcCollectionPaidByInvoice(txCtx, &depositDetails)
			if err != nil {
				return err
			}

			if depositDetails.TotalPayment > 0 {
				indexGiro := -1
				for i, payment := range detail.Payment {
					if payment.PayType == 2 {
						indexGiro = i
						break
					}
				}
				indexDiscount := -1
				indexDiscountPayType := 0
				for i, payment := range detail.Payment {
					if payment.PayType == 1 {
						indexDiscount = i
						indexDiscountPayType = 1
						break
					} else if payment.PayType == 2 {
						indexDiscount = i
						indexDiscountPayType = 2
						continue
					} else if payment.PayType == 3 {
						if indexDiscountPayType > 2 || indexDiscountPayType == 0 {
							indexDiscount = i
							indexDiscountPayType = 3
						}
					} else {
						if indexDiscountPayType > 3 || indexDiscountPayType == 0 {
							indexDiscount = 0
						}
					}
				}

				for index, payment := range detail.Payment {
					var depositPayment model.DepositPayment

					err := structs.Automapper(payment, &depositPayment)
					if err != nil {
						return err
					}
					depositPayment.DepositNo = depositNo
					depositPayment.CustID = request.CustID

					if index == indexDiscount {
						depositPayment.Discount = &detail.Discount
					}

					if indexGiro > -1 && indexGiro == index {
						depositPayment.Materai = &detail.Materai
					} else {
						zeroFloat := float64(0)
						depositPayment.Materai = &zeroFloat
					}

					switch depositPayment.PayType {
					case 2: // cheque_giro
						err := service.DepositRepository.UpdateAmountProgressionCheque(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 3: // bank_transfer
						err := service.DepositRepository.UpdateAmountProgressionTransfer(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 4: // return
						err := service.DepositRepository.UpdateAmountProgressionReturn(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					case 5: // cndn
						err := service.DepositRepository.UpdateAmountProgressionCNDN(txCtx, &depositPayment, auditor)
						if err != nil {
							return err
						}
					}

					_, err = service.DepositRepository.StorePayment(txCtx, &depositPayment)
					if err != nil {
						return err
					}
				}
			}
		}

		// handle expenses if present
		for _, expense := range request.Expense {
			var depositExpense model.DepositExpense

			err := structs.Automapper(expense, &depositExpense)
			if err != nil {
				return err
			}

			depositExpense.DepositNo = depositNo
			depositExpense.CustID = request.CustID
			depositExpense.PaymentAmount = decimal.NewFromFloat(expense.Amount)
			if request.UpdatedBy != nil {
				depositExpense.CreatedBy = *request.UpdatedBy
			}

			_, err = service.DepositRepository.StoreExpense(txCtx, &depositExpense)
			if err != nil {
				return err
			}

			if err := service.DepositRepository.DeductExpense(txCtx, &depositExpense); err != nil {
				return err
			}
		}

		depositModel.RemainingAmount = remainingAmount
		err = service.DepositRepository.Update(txCtx, depositNo, request.CustID, depositModel)
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
