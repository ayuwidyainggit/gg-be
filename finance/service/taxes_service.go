package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"time"
)

var (
	TAXES_GENERATE_STATUS_ACTIVE = 1
	TAXES_GENERATE_STATUS_DELETE = 0
)

type TaxesService interface {
	Generate(request entity.TaxesGenerateReq) (err error)
	ListReport(dataFilter entity.TaxesQueryFilter) (data []entity.TaxesResponse, total int64, lastPage int, err error)
	ListGenerate(dataFilter entity.TaxesGenerateQueryFilter) (data []entity.TaxesGenerateResponse, total int64, lastPage int, err error)
	Delete(custId string, taxesID int64, userId int64) (err error)
	DeleteBulk(custId string, taxesIDs []int64, userId int64) (err error)
}

type TaxesServiceImpl struct {
	Transaction      repository.Dbtransaction
	TaxesRepository  repository.TaxesRepository
	MTaxesRepository repository.MTaxesRepository
}

func NewTaxesService(taxesRepository repository.TaxesRepository, repository repository.MTaxesRepository, transaction repository.Dbtransaction) *TaxesServiceImpl {
	return &TaxesServiceImpl{
		TaxesRepository:  taxesRepository,
		MTaxesRepository: repository,
		Transaction:      transaction,
	}
}

func (service *TaxesServiceImpl) Generate(request entity.TaxesGenerateReq) (err error) {
	c := context.Background()
	year := time.Now().Year()

	invoices, err := service.TaxesRepository.GetInvoiceInfo(request.CustID, request.Invoices)
	if err != nil {
		return err
	}

	var mapOutletInvoices = entity.MapOutletInvoice{}
	for _, invoice := range invoices {
		mapOutletInvoices.MapOutletInvoice(invoice.OutletID, invoice.TaxInvoiceForm, invoice.InvoiceNo)
	}

	mTaxes, err := service.MTaxesRepository.GetNewestSerialByStatus(request.CustID, entity.STATUS_TAXES_ACTIVE, year)
	if err != nil {
		mTaxes, err = service.MTaxesRepository.GetNewestSerialByStatus(request.CustID, entity.STATUS_TAXES_RESERVED, year)
		if err != nil {
			return err
		}
	}

	taxGenerated := &entity.TaxesObj{
		TransactionStatusCode: mTaxes.TransactionStatusCode,
		SerialCode:            mTaxes.SerialCode,
		RemainingQty:          mTaxes.RemainingQty,
		SerialFrom:            mTaxes.SerialFrom,
		SerialTo:              mTaxes.SerialTo,
	}

	for _, mapOutletInvoice := range mapOutletInvoices {
		taxGenerated.OutletInvoices = append(taxGenerated.OutletInvoices, *mapOutletInvoice)
	}

	taxGenerated.GenerateInvoice()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		var taxModels []model.Taxes
		var lastGeneratedTax string
		for index, InvoiceGenerated := range taxGenerated.InvoiceGenerated {
			taxModels = append(taxModels, model.Taxes{
				CustID:    request.CustID,
				MTaxId:    *mTaxes.MTaxID,
				TaxNo:     InvoiceGenerated.Tax,
				InvoiceNo: InvoiceGenerated.Invoice,
				CreatedBy: &request.CreatedBy,
				UpdatedBy: &request.CreatedBy,
				Status:    TAXES_GENERATE_STATUS_ACTIVE,
			})
			if index == len(taxGenerated.InvoiceGenerated)-1 {
				lastGeneratedTax = InvoiceGenerated.Tax
			}
		}

		err := service.TaxesRepository.Store(txCtx, taxModels)
		if err != nil {
			return err
		}

		err = service.MTaxesRepository.UpdateAfterGenerate(txCtx, *mTaxes.MTaxID, taxGenerated.RemainingQty, taxGenerated.Status, lastGeneratedTax)
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

func (service *TaxesServiceImpl) ListReport(dataFilter entity.TaxesQueryFilter) (data []entity.TaxesResponse, total int64, lastPage int, err error) {
	Taxess, total, lastPage, err := service.TaxesRepository.FindAllReportTaxes(dataFilter, dataFilter.CustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Taxess {
		var vResp entity.TaxesResponse
		structs.Automapper(row, &vResp)

		if row.InvoiceDate != nil {
			InvoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = InvoiceDate
		}

		if dataFilter.Type == "invoice" {
			if row.InvoiceDate != nil {
				TaxesDate := row.InvoiceDate.Format("2006-01-02")
				vResp.TaxDate = TaxesDate
			}
			vResp.Type = "invoice"
		} else {
			if row.ReturnDate != nil {
				TaxesDate := row.ReturnDate.Format("2006-01-02")
				vResp.TaxDate = TaxesDate
			}
			vResp.Type = "return"
		}

		// if row.StatusCheque != nil {
		// 	statusText := entity.ConvStatus(entity.StatusGiro, *row.StatusCheque)
		// 	vResp.StatusClearing = &statusText
		// }

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *TaxesServiceImpl) ListGenerate(dataFilter entity.TaxesGenerateQueryFilter) (data []entity.TaxesGenerateResponse, total int64, lastPage int, err error) {
	Taxess, total, lastPage, err := service.TaxesRepository.TaxGenerateList(dataFilter, dataFilter.CustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range Taxess {
		var vResp entity.TaxesGenerateResponse
		structs.Automapper(row, &vResp)

		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = invoiceDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *TaxesServiceImpl) Delete(custId string, taxesID int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.TaxesRepository.Delete(txCtx, custId, taxesID, userId)
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

func (service *TaxesServiceImpl) DeleteBulk(custId string, taxesIDs []int64, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.TaxesRepository.DeleteBulk(txCtx, custId, taxesIDs, userId)
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
