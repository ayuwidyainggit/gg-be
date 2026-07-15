package service

import (
	"context"
	"encoding/json"
	"mobile/entity"
	"mobile/model"
	"time"

	"mobile/repository"
)

type CollectionPayService interface {
	Store(ctx context.Context, request entity.CreateCollectionPayRequest) (data []entity.StoreCollectionPayResponse, err error)
	List(ctx context.Context, dataFilter entity.CollectionPayQueryFilter) (data []entity.CollectionPayResponse, total int64, lastPage int64, err error)
	StoreNoPayment(ctx context.Context, request entity.CreateNoPaymentRequest) error
}

func NewCollectionPayService(
	paymentRepo repository.PaymentRepository,
	bankRepo repository.BankRepository,
	collectionRepo repository.CollectionRepository,
	transaction repository.Dbtransaction) *collectionPayServiceImpl {
	return &collectionPayServiceImpl{
		PaymentRepo:    paymentRepo,
		BankRepo:       bankRepo,
		CollectionRepo: collectionRepo,
		Transaction:    transaction,
	}
}

type collectionPayServiceImpl struct {
	PaymentRepo    repository.PaymentRepository
	BankRepo       repository.BankRepository
	CollectionRepo repository.CollectionRepository
	Transaction    repository.Dbtransaction
}

func (service *collectionPayServiceImpl) Store(ctx context.Context, request entity.CreateCollectionPayRequest) (data []entity.StoreCollectionPayResponse, err error) {
	timeNow := time.Now()
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, detail := range request.Detail {
			documentNo, err := service.PaymentRepo.GetNewDocumentNo(txCtx, request.CustID)
			if err != nil {
				return err
			}
			var imgArray string
			if len(detail.Image) > 0 {
				imgJson, err := json.Marshal(detail.Image)
				if err != nil {
					return err
				}
				imgArray = string(imgJson)
			}

			paymentTrx := &model.PaymentTrx{
				CustID:           request.CustID,
				InvoiceNo:        detail.InvoiceNo,
				CollectionNo:     request.CollectionNo,
				DocumentNo:       documentNo,
				OutletID:         request.OutletID,
				PONumber:         detail.OrderNo,
				TrxSource:        "L", // hardcode by request SA to avoid wrong data input
				EmpID:            request.EmpID,
				TotalTransaction: detail.InvoiceAmount,
				PaymentAmount:    detail.PaidAmount,
				RemainingAmount:  detail.NewRemainingAmount,
				Notes:            &detail.Notes,
				Date:             timeNow,
				CreatedBy:        request.UserID,
				CreatedAt:        timeNow,
				UpdatedAt:        &timeNow,
			}
			if imgArray != "" {
				paymentTrx.Files = imgArray
			}
			err = service.PaymentRepo.StorePaymentTrx(txCtx, paymentTrx)
			if err != nil {
				return err
			}
			for _, payment := range detail.Payment {
				respData := entity.StoreCollectionPayResponse{
					InvoiceNo: detail.InvoiceNo,
				}
				paymentTrxDet := &model.PaymentTrxDet{
					CustID:       request.CustID,
					PaymentTrxID: paymentTrx.PaymentTrxID,
					Amount:       payment.PaymentAmount,
					CreatedBy:    request.UserID,
					CreatedAt:    timeNow,
					UpdatedAt:    &timeNow,
					PayType:      int16(payment.PayType),
				}
				if payment.PayType == 3 {
					bank, err := service.BankRepo.FindByCustIDAndOutletID(txCtx, request.CustID, int64(request.OutletID))
					if err != nil {
						return err
					}
					docNoBank, err := service.BankRepo.GetNewDocNoBank(txCtx, request.CustID)
					if err != nil {
						return err
					}
					bankTransfer := &model.BankTransfer{
						CustID:             request.CustID,
						DocNoBank:          docNoBank,
						OwnerID:            1,
						SalesmanID:         &request.EmpID,
						OutletID:           &request.OutletID,
						BankID:             &bank.BankID,
						AccountNo:          &bank.AccountNo,
						TransferDate:       &timeNow,
						Amount:             payment.PaymentAmount,
						StatusBankTransfer: 2,
						CreatedBy:          request.UserID,
						CreatedAt:          timeNow,
						UpdatedAt:          &timeNow,
						OutletBankID:       &bank.OutletBankID,
						PaidAmount:         detail.PaidAmount,
						RemainingAmount:    detail.NewRemainingAmount,
						AccountName:        bank.AccountName,
					}
					err = service.PaymentRepo.StoreBankTransfer(txCtx, bankTransfer)
					if err != nil {
						return err
					}

					respData.DocNoBank = docNoBank
					paymentTrxDet.BankTransferNo = &bankTransfer.BankTransferNo
					paymentTrxDet.PayType = 3
				}
				err = service.PaymentRepo.StorePaymentTrxDetail(txCtx, paymentTrxDet)
				if err != nil {
					return err
				}
				data = append(data, respData)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (service *collectionPayServiceImpl) List(ctx context.Context, dataFilter entity.CollectionPayQueryFilter) (data []entity.CollectionPayResponse, total int64, lastPage int64, err error) {
	arpays, total, lastPage, err := service.PaymentRepo.GetByFilter(ctx, dataFilter)
	if err != nil {
		return nil, total, lastPage, err
	}

	return arpays, total, lastPage, err
}

func (service *collectionPayServiceImpl) StoreNoPayment(ctx context.Context, request entity.CreateNoPaymentRequest) error {
	invoices, err := service.CollectionRepo.GetInvoiceList(ctx, request.CustID, int(request.OutletID))
	if err != nil {
		return err
	}
	if len(invoices) == 0 {
		return nil
	}

	paymentDate, errParse := time.Parse(time.DateOnly, request.PaymentDate)
	if errParse != nil {
		return errParse
	}
	userID := int(request.UserID)
	timeNow := time.Now()

	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		for _, invoice := range invoices {
			insert := &model.CollectionNoPayment{
				CustID:                 request.CustID,
				CollectionNo:           invoice.CollectionNo,
				InvoiceNo:              invoice.InvoiceNo,
				SalesmanID:             &request.EmpID,
				MissedPaymentReasonsID: &request.ReasonID,
				Reason:                 &request.Reason,
				PaymentDate:            &paymentDate,
				CreatedBy:              &userID,
				CreatedAt:              &timeNow,
			}

			err := service.CollectionRepo.StoreCollectionNoPayment(txCtx, insert)
			if err != nil {
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
