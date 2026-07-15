package service

import (
	"context"
	"encoding/json"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/structs"
	"mobile/repository"
	"strings"
	"time"
)

type InvoicesService interface {
	GetInvoices(req entity.InvoicesListReq) (resp []entity.InvoicesGetResp, err error)
	CreatePayment(ctx context.Context, req entity.InvoicesPaymentCreate) (err error)
	DetailInvoice(invoiceNo string, custID string) (response entity.DetilPaymentDepositInvoiceDetailResponse, err error)
}

func NewInvoicesService(
	config env.ConfigEnv,
	invoicesRepository repository.InvoicesRepository,
	paymentRepository repository.PaymentRepository,
	outletBankRepository repository.OutletBankRepository,
	transaction repository.Dbtransaction) *InvoicesServiceImpl {
	return &InvoicesServiceImpl{
		Config:               config,
		InvoicesRepository:   invoicesRepository,
		PaymentRepository:    paymentRepository,
		OutletBankRepository: outletBankRepository,
		Transaction:          transaction,
	}
}

type InvoicesServiceImpl struct {
	Config               env.ConfigEnv
	InvoicesRepository   repository.InvoicesRepository
	PaymentRepository    repository.PaymentRepository
	OutletBankRepository repository.OutletBankRepository
	Transaction          repository.Dbtransaction
}

func (service *InvoicesServiceImpl) GetInvoices(req entity.InvoicesListReq) (resp []entity.InvoicesGetResp, err error) {
	resp = []entity.InvoicesGetResp{
		entity.InvoicesGetResp{
			InvoiceNumber:  "INV001",
			OrderNumber:    "100000001",
			InvoiceDate:    "2024-02-01",
			DueDate:        "2024-02-01",
			PaymentOption:  "",
			SettlementDate: "",
			InvoiceAmount:  100000,
			PaymentDetail: []entity.PaymentDetail{
				entity.PaymentDetail{
					Amount:        1000000,
					PaymentMethod: "cash",
				},
			},
			Images:         []string{"1.jpg", "2.jpg"},
			CollectHistory: []any{},
		}, entity.InvoicesGetResp{
			InvoiceNumber:  "INV002",
			OrderNumber:    "100000002",
			InvoiceDate:    "2024-02-02",
			DueDate:        "2024-02-02",
			PaymentOption:  "",
			SettlementDate: "",
			InvoiceAmount:  100000,
			PaymentDetail:  []entity.PaymentDetail{},
			Images:         []string{"1.jpg", "2.jpg"},
			CollectHistory: []any{},
		},
	}
	return resp, nil
}

func (service *InvoicesServiceImpl) CreatePayment(ctx context.Context, request entity.InvoicesPaymentCreate) (err error) {
	var (
		defaultOwnerID            = 1
		defaultCreditCNDN         = "credit"
		defaultTitipBayarCNDNType = int64(72)
	)
	/*
	   semua type baik Canvas/TO akan masuk ke transaction dan transcation detail.
	   khusus kondisi berikut perlu insert ke table terkait
	   TO :
	   cash = cndn

	   C/TO:
	   transfer = bank_transfer
	*/
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		detail := request.Detail

		var filesBytes []byte
		if len(request.Files) > 0 {
			filesBytes, err = json.Marshal(request.Files)
			if err != nil {
				return err
			}
		}

		remainingAmount := request.TotalPaymentBalance - request.TotalPayment
		paymentTrx := &model.PaymentTrx{
			CustID:           request.CustID,
			OutletID:         request.OutletID,
			PONumber:         detail.InvoiceNo,
			TrxSource:        request.OprType.String(),
			EmpID:            request.EmpID,
			TotalTransaction: request.TotalPaymentBalance,
			CreatedBy:        request.CreatedBy,
			PaymentAmount:    request.TotalPayment,
			RemainingAmount:  remainingAmount,
			Date:             time.Now(),
			Notes:            request.Notes,
			Files:            filesBytes,
		}

		err = service.PaymentRepository.StorePaymentTrx(txCtx, paymentTrx)
		if err != nil {
			return err
		}
		for _, payment := range detail.Payment {
			var bankTransferNo *int
			// Check payment type in DB by ID
			pt, err := service.PaymentRepository.FindPaymentTypeByID(ctx, int(payment.PayType))
			if err != nil {
				return err
			}

			if request.OprType.IsTakingOrder() && strings.ToLower(pt.PaymentTypeCode) == "cash" {
				nowDate := time.Now()
				cndn := &model.Cndn{
					CustID:    request.CustID,
					CndnDate:  &nowDate,
					OwnerId:   defaultOwnerID,
					CndnJenis: &defaultCreditCNDN,
					CndnType:  defaultTitipBayarCNDNType,
					Amount:    &payment.PaymentAmount,
					CreatedBy: &request.CreatedBy,
					CreatedAt: time.Now(),
					OutletId:  &request.OutletID,
					Notes:     request.Notes,
				}

				err = service.InvoicesRepository.StorePaymentCndn(txCtx, cndn)
				if err != nil {
					return err
				}
			}

			if payment.PayType == 3 {
				transferDate := time.Now().UTC()

				// Fetch bank info from the outlet's registered bank account
				outletBank, err := service.OutletBankRepository.FindFirstByOutletID(txCtx, request.CustID, request.OutletID)
				if err != nil {
					return err
				}

				bankTfr := &model.BankTransfer{
					CustID:          request.CustID,
					SalesmanID:      &request.SalesmanID,
					OutletID:        &request.OutletID,
					Amount:          payment.PaymentAmount,
					CreatedBy:       request.CreatedBy,
					TransferDate:    &transferDate,
					RemainingAmount: remainingAmount,
				}

				if outletBank != nil {
					bankTfr.BankID = &outletBank.BankID
					bankTfr.AccountNo = outletBank.AccountNo
					bankTfr.OutletBankID = &outletBank.OutletBankID

					if outletBank.AccountName != nil {
						bankTfr.AccountName = *outletBank.AccountName
					}
				}
				err = service.PaymentRepository.StoreBankTransfer(txCtx, bankTfr)
				if err != nil {
					return err
				}

				bankTransferNo = &bankTfr.BankTransferNo
			}

			paymentTrxDetail := &model.PaymentTrxDet{
				PaymentTrxID:   paymentTrx.PaymentTrxID,
				CustID:         request.CustID,
				PayType:        payment.PayType.Value(),
				BankTransferNo: bankTransferNo,
				Amount:         payment.PaymentAmount,
				CreatedBy:      request.CreatedBy,
			}

			err = service.PaymentRepository.StorePaymentTrxDetail(txCtx, paymentTrxDetail)
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

func (service *InvoicesServiceImpl) DetailInvoice(invoiceNo string, custID string) (response entity.DetilPaymentDepositInvoiceDetailResponse, err error) {
	cashs, err := service.InvoicesRepository.FindDetailPaymentByInvoice(1, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Cash = make([]entity.DetilDepositPaymentInvoice, len(cashs))

	for i, cash := range cashs {
		var detailData entity.DetilDepositPaymentInvoice
		err = structs.Automapper(cash, &detailData)
		if err != nil {
			return response, err
		}

		// Calculate TotalPayment based on the fields: PaymentAmount, Materai, and Discount
		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance

		paymentImages, err := service.InvoicesRepository.FindPaymentImagesByNo(cash.DepositNo, cash.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositInvoicePaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositInvoicePaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Cash[i] = detailData
	}

	cek, err := service.InvoicesRepository.FindDetailPaymentByInvoice(2, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Cek = make([]entity.DetilDepositPaymentInvoice, len(cek))

	for i, ceks := range cek {
		var detailData entity.DetilDepositPaymentInvoice
		err = structs.Automapper(ceks, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.InvoicesRepository.FindPaymentImagesByNo(ceks.DepositNo, ceks.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositInvoicePaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositInvoicePaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Cek[i] = detailData
	}

	transfers, err := service.InvoicesRepository.FindDetailPaymentByInvoice(3, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Trasfer = make([]entity.DetilDepositPaymentInvoice, len(transfers))

	for i, transfer := range transfers {
		var detailData entity.DetilDepositPaymentInvoice
		err = structs.Automapper(transfer, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.InvoicesRepository.FindPaymentImagesByNo(transfer.DepositNo, transfer.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositInvoicePaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositInvoicePaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Trasfer[i] = detailData
	}

	returns, err := service.InvoicesRepository.FindDetailPaymentByInvoice(4, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.Return = make([]entity.DetilDepositPaymentInvoice, len(returns))

	for i, returna := range returns {
		var detailData entity.DetilDepositPaymentInvoice
		err = structs.Automapper(returna, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.InvoicesRepository.FindPaymentImagesByNo(returna.DepositNo, returna.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositInvoicePaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositInvoicePaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.Return[i] = detailData
	}

	cndns, err := service.InvoicesRepository.FindDetailPaymentByInvoice(5, invoiceNo, custID)
	if err != nil {
		return response, err
	}

	response.CNDN = make([]entity.DetilDepositPaymentInvoice, len(cndns))

	for i, cndn := range cndns {
		var detailData entity.DetilDepositPaymentInvoice
		err = structs.Automapper(cndn, &detailData)
		if err != nil {
			return response, err
		}

		detailData.TotalPayment = detailData.PaymentAmount + detailData.Materai + detailData.Discount
		// detailData.PaymentBalance = Deposit.TotalPaymentBalance
		paymentImages, err := service.InvoicesRepository.FindPaymentImagesByNo(cndn.DepositNo, cndn.InvoiceNo)
		if err != nil {
			return response, err
		}

		detailData.Images = make([]entity.DepositInvoicePaymentImage, len(paymentImages))
		for j, img := range paymentImages {
			var imageData entity.DepositInvoicePaymentImage
			err = structs.Automapper(img, &imageData)
			if err != nil {
				return response, err
			}
			detailData.Images[j] = imageData
		}

		response.CNDN[i] = detailData
	}

	// ownerName := entity.ConvStatus(entity.OwnerGiro, response.OwnerID)
	// response.OwnerName = ownerName

	// statusText := entity.ConvStatus(entity.StatusGiro, response.StatusCheque)
	// response.StatusChequeText = &statusText

	// response.UsedAmount = float64(0)
	// response.RemainingAmount = response.Amount - response.UsedAmount

	return response, nil
}
