package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"time"
)

type VatExtractService interface {
	ListVatIn(dataFilter entity.VatExtractQueryFilter) (data []entity.VatExtractListResponse, total int64, lastPage int, err error)
	Extract(request entity.VatExtractReq) (resp entity.VatExtractResp, err error)
	ExtractResult(vatExtractID int64, custID string, parentCustId string) (resp []entity.VatExtractListResponse, err error)
	ListVatExtract(dataFilter entity.VatExtractResultQueryFilter) (datas []entity.VatExtractList, total int64, lastPage int, err error)
	ExtractDownloadResult(vatExtractID int64, custID string, parentCustId string) (interface{}, error)
}

type VatExtractServiceImpl struct {
	Repository   repository.VatExtractRepository
	ApRepository repository.ApSupplierInvoiceReturnRepository
	Transaction  repository.Dbtransaction
}

func NewVatExtractService(repository repository.VatExtractRepository, apRepository repository.ApSupplierInvoiceReturnRepository, transaction repository.Dbtransaction) *VatExtractServiceImpl {
	return &VatExtractServiceImpl{
		Repository:   repository,
		ApRepository: apRepository,
		Transaction:  transaction,
	}
}

func (service *VatExtractServiceImpl) ListVatIn(dataFilter entity.VatExtractQueryFilter) (datas []entity.VatExtractListResponse, total int64, lastPage int, err error) {
	if dataFilter.InvoiceType == entity.AP_TYPE_INVOICE {
		return service.ListVatInInvoice(dataFilter)
	} else {
		return service.ListVatInReturn(dataFilter)
	}
}

func (service *VatExtractServiceImpl) ListVatInInvoice(dataFilter entity.VatExtractQueryFilter) (datas []entity.VatExtractListResponse, total int64, lastPage int, err error) {
	vatIns, total, lastPage, err := service.ApRepository.FindAllVatInByCustId(dataFilter)
	if err != nil {
		return
	}

	for _, row := range vatIns {
		var invoiceDate, taxInvoiceDate, extractedAt string
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}

		if row.TaxInvoiceDate != nil {
			taxInvoiceDate = row.TaxInvoiceDate.Format("2006-01-02")
		}

		if row.ExtractedAt != nil {
			extractedAt = row.ExtractedAt.Format("2006-01-02")
		}

		datas = append(datas, entity.VatExtractListResponse{
			TransactionID:  row.ID,
			InvoiceNo:      row.InvoiceNo,
			InvoiceDate:    invoiceDate,
			InvoiceType:    row.ApType,
			NPWP:           row.Npwp,
			SupplierCode:   row.SupCode,
			SupplierName:   row.SupName,
			Address:        row.Address,
			DPP:            *row.SubTotal,
			PPN:            *row.VatValue,
			PPNBM:          *row.VatLgValue,
			TaxNo:          *row.TaxInvoiceNo,
			TaxDate:        taxInvoiceDate,
			TaxExtractDate: extractedAt,
		})
	}

	return
}

func (service *VatExtractServiceImpl) ListVatInReturn(dataFilter entity.VatExtractQueryFilter) (datas []entity.VatExtractListResponse, total int64, lastPage int, err error) {
	vatIns, total, lastPage, err := service.ApRepository.FindAllVatInByCustId(dataFilter)
	if err != nil {
		return
	}
	var documents []string

	for _, row := range vatIns {
		var invoiceDate, taxInvoiceDate, extractedAt string
		if row.InvoiceDate != nil {
			invoiceDate = row.InvoiceDate.Format("2006-01-02")
		}

		if row.TaxInvoiceDate != nil {
			taxInvoiceDate = row.TaxInvoiceDate.Format("2006-01-02")
		}

		if row.ExtractedAt != nil {
			extractedAt = row.ExtractedAt.Format("2006-01-02")
		}

		datas = append(datas, entity.VatExtractListResponse{
			TransactionID:    row.ID,
			InvoiceNo:        row.InvoiceNo,
			InvoiceDate:      invoiceDate,
			InvoiceType:      row.ApType,
			NPWP:             row.Npwp,
			SupplierCode:     row.SupCode,
			SupplierName:     row.SupName,
			Address:          row.Address,
			DPP:              *row.SubTotal,
			PPN:              *row.VatValue,
			PPNBM:            *row.VatLgValue,
			TaxNo:            *row.TaxInvoiceNo,
			TaxDate:          taxInvoiceDate,
			TaxExtractDate:   extractedAt,
			ReturnDocumentNo: row.DocumentNo,
		})

		documents = append(documents, row.DocumentNo)
	}

	apReturns, err := service.ApRepository.FindByDocumentsNo(documents, dataFilter.CustId, dataFilter.ParentCustId)
	if err != nil {
		return
	}

	apReturnsMap := make(map[string]model.ApSuppilerInvoiceReturnList)
	for _, apReturn := range apReturns {
		apReturnsMap[apReturn.DocumentNo] = apReturn
	}

	for index, _ := range datas {
		datas[index].ReturnDate = apReturnsMap[datas[index].ReturnDocumentNo].InvoiceDate.Format("2006-01-02")
	}

	return
}

func (service *VatExtractServiceImpl) Extract(request entity.VatExtractReq) (resp entity.VatExtractResp, err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		now := time.Now()
		vatExtractModel := model.VatExtract{
			CustID:         request.CustID,
			VatExtractType: entity.VAT_TYPE_IN,
			ExtractTotal:   len(request.TransactionID),
			CreatedBy:      request.CreatedBy,
			UpdatedBy:      request.CreatedBy,
			InvoiceType:    request.InvoiceType,
			CreatedAt:      &now,
		}

		err := service.Repository.Store(txCtx, &vatExtractModel)
		if err != nil {
			return err
		}

		var vatExtractDetailModel []model.VatExtractDetail

		for _, transactionID := range request.TransactionID {
			vatExtractDetailModel = append(vatExtractDetailModel, model.VatExtractDetail{
				VatExtractID: *vatExtractModel.VatExtractID,
				ReferenceID:  transactionID,
			})
		}

		err = service.Repository.StoreDetail(txCtx, vatExtractDetailModel)
		if err != nil {
			return err
		}

		resp.ID = *vatExtractModel.VatExtractID
		return nil
	})
	return
}

func (service *VatExtractServiceImpl) ExtractResult(vatExtractID int64, custID string, parentCustId string) (resp []entity.VatExtractListResponse, err error) {
	vatExtractDetails, err := service.Repository.FindExtractResult(vatExtractID, custID, parentCustId)
	if err != nil {
		return
	}

	if len(vatExtractDetails) > 0 {
		if vatExtractDetails[0].VatExtractType == entity.VAT_TYPE_IN {
			if vatExtractDetails[0].InvoiceType == "I" {
				return service.ExtractResultVatInInvoice(vatExtractDetails)
			} else {
				return service.ExtractResultVatInReturn(vatExtractDetails, custID, parentCustId)
			}
		}
	}
	return
}

func (service *VatExtractServiceImpl) ExtractResultVatInInvoice(vatExtractDetails []model.VatExtractDetailList) (resp []entity.VatExtractListResponse, err error) {
	for _, vatExtractDetail := range vatExtractDetails {
		var invoiceDate, taxInvoiceDate, extractedAt string
		if vatExtractDetail.InvoiceDate != nil {
			invoiceDate = vatExtractDetail.InvoiceDate.Format("2006-01-02")
		}

		if vatExtractDetail.TaxInvoiceDate != nil {
			taxInvoiceDate = vatExtractDetail.TaxInvoiceDate.Format("2006-01-02")
		}

		if vatExtractDetail.ExtractedAt != nil {
			extractedAt = vatExtractDetail.ExtractedAt.Format("2006-01-02")
		}

		resp = append(resp, entity.VatExtractListResponse{
			TransactionID:  vatExtractDetail.ID,
			InvoiceNo:      vatExtractDetail.InvoiceNo,
			InvoiceDate:    invoiceDate,
			InvoiceType:    vatExtractDetail.ApType,
			NPWP:           vatExtractDetail.Npwp,
			SupplierCode:   vatExtractDetail.SupCode,
			SupplierName:   vatExtractDetail.SupName,
			Address:        vatExtractDetail.Address,
			DPP:            *vatExtractDetail.SubTotal,
			PPN:            *vatExtractDetail.VatValue,
			PPNBM:          *vatExtractDetail.VatLgValue,
			TaxNo:          *vatExtractDetail.TaxInvoiceNo,
			TaxDate:        taxInvoiceDate,
			TaxExtractDate: extractedAt,
		})
	}

	return
}

func (service *VatExtractServiceImpl) ExtractResultVatInReturn(vatExtractDetails []model.VatExtractDetailList, custID string, parentCustId string) (resp []entity.VatExtractListResponse, err error) {
	var documents []string

	for _, vatExtractDetail := range vatExtractDetails {
		documents = append(documents, vatExtractDetail.DocumentNo)
	}

	apReturns, err := service.ApRepository.FindByDocumentsNo(documents, custID, parentCustId)
	if err != nil {
		return
	}

	apReturnsMap := make(map[string]model.ApSuppilerInvoiceReturnList)
	for _, apReturn := range apReturns {
		apReturnsMap[apReturn.DocumentNo] = apReturn
	}

	for _, vatExtractDetail := range vatExtractDetails {
		var invoiceDate, taxInvoiceDate, extractedAt, returnDate string
		if vatExtractDetail.InvoiceDate != nil {
			invoiceDate = vatExtractDetail.InvoiceDate.Format("2006-01-02")
		}

		if vatExtractDetail.TaxInvoiceDate != nil {
			taxInvoiceDate = vatExtractDetail.TaxInvoiceDate.Format("2006-01-02")
		}

		if vatExtractDetail.ExtractedAt != nil {
			extractedAt = vatExtractDetail.ExtractedAt.Format("2006-01-02")
		}

		if apReturnsMap[vatExtractDetail.DocumentNo].InvoiceDate != nil {
			returnDate = apReturnsMap[vatExtractDetail.DocumentNo].InvoiceDate.Format("2006-01-02")
		}

		resp = append(resp, entity.VatExtractListResponse{
			TransactionID:    vatExtractDetail.ID,
			InvoiceNo:        vatExtractDetail.InvoiceNo,
			InvoiceDate:      invoiceDate,
			InvoiceType:      vatExtractDetail.ApType,
			NPWP:             vatExtractDetail.Npwp,
			SupplierCode:     vatExtractDetail.SupCode,
			SupplierName:     vatExtractDetail.SupName,
			Address:          vatExtractDetail.Address,
			DPP:              *vatExtractDetail.SubTotal,
			PPN:              *vatExtractDetail.VatValue,
			PPNBM:            *vatExtractDetail.VatLgValue,
			TaxNo:            *vatExtractDetail.TaxInvoiceNo,
			TaxDate:          taxInvoiceDate,
			TaxExtractDate:   extractedAt,
			ReturnDocumentNo: vatExtractDetail.DocumentNo,
			ReturnDate:       returnDate,
		})
	}
	return
}

func (service *VatExtractServiceImpl) ListVatExtract(dataFilter entity.VatExtractResultQueryFilter) (datas []entity.VatExtractList, total int64, lastPage int, err error) {
	vatExtracts, total, lastPage, err := service.Repository.FindAllVatExtractByCustId(dataFilter)
	if err != nil {
		return
	}

	for _, row := range vatExtracts {
		var vResp entity.VatExtractList
		structs.Automapper(row, &vResp)

		datas = append(datas, vResp)
	}

	return
}

func (service *VatExtractServiceImpl) ExtractDownloadResult(vatExtractID int64, custID string, parentCustId string) (interface{}, error) {
	vatExtractDetails, err := service.Repository.FindExtractResult(vatExtractID, custID, parentCustId)
	if err != nil {
		return nil, err
	}
	if len(vatExtractDetails) > 0 {
		if vatExtractDetails[0].VatExtractType == entity.VAT_TYPE_IN {
			if vatExtractDetails[0].InvoiceType == "I" {
				vatInInvoice, err := service.ExtractDownloadResultVatInInvoice(vatExtractDetails)
				if err != nil {
					return nil, err
				}
				return vatInInvoice, nil
			} else {
				vatInReturn, err := service.ExtractDownloadResultVatInReturn(vatExtractDetails, custID, parentCustId)
				if err != nil {
					return nil, err
				}
				return vatInReturn, nil
			}
		} else {

		}
	}
	return nil, err
}

func (service *VatExtractServiceImpl) ExtractDownloadResultVatInInvoice(vatExtractDetails []model.VatExtractDetailList) (datas []entity.VatExtractDowloadVatInInvoice, err error) {
	for _, vatExtractDetail := range vatExtractDetails {
		var month time.Month
		faktur := entity.GenerateFakturComponents(*vatExtractDetail.TaxInvoiceNo)
		var taxInvoiceDate string

		if vatExtractDetail.TaxInvoiceDate != nil {
			taxInvoiceDate = vatExtractDetail.TaxInvoiceDate.Format("2006-01-02")
			month = vatExtractDetail.TaxInvoiceDate.Month()
		}
		datas = append(datas, entity.VatExtractDowloadVatInInvoice{
			KdJenisTransaksi: faktur.KdJjenisTansaksi,
			FgPengganti:      faktur.FgPengganti,
			NomorFaktur:      faktur.NomorFaktur,
			MasaPajak:        int(month),
			TahunPajak:       faktur.TahunPajak,
			TanggalFaktur:    taxInvoiceDate,
			NPWP:             vatExtractDetail.Npwp,
			Nama:             vatExtractDetail.SupName,
			AlamatLengkap:    vatExtractDetail.Address,
			JumlahDPP:        *vatExtractDetail.SubTotal,
			JumlahPPN:        *vatExtractDetail.VatValue,
			JumlahPPNBM:      *vatExtractDetail.VatLgValue,
			IsCreditable:     1,
		})
	}

	return
}

func (service *VatExtractServiceImpl) ExtractDownloadResultVatInReturn(vatExtractDetails []model.VatExtractDetailList, custID string, parentCustId string) (datas []entity.VatExtractDowloadVatInReturn, err error) {
	var documents []string

	for _, vatExtractDetail := range vatExtractDetails {
		documents = append(documents, vatExtractDetail.DocumentNo)
	}

	apReturns, err := service.ApRepository.FindByDocumentsNo(documents, custID, parentCustId)
	if err != nil {
		return
	}

	apReturnsMap := make(map[string]model.ApSuppilerInvoiceReturnList)
	for _, apReturn := range apReturns {
		apReturnsMap[apReturn.DocumentNo] = apReturn
	}

	for _, vatExtractDetail := range vatExtractDetails {
		faktur := entity.GenerateFakturComponents(*vatExtractDetail.TaxInvoiceNo)
		var taxInvoiceDate, returnDate string
		var monthReturn time.Month
		if vatExtractDetail.TaxInvoiceDate != nil {
			taxInvoiceDate = vatExtractDetail.TaxInvoiceDate.Format("2006-01-02")
		}

		if apReturnsMap[vatExtractDetail.DocumentNo].InvoiceDate != nil {
			returnDate = apReturnsMap[vatExtractDetail.DocumentNo].InvoiceDate.Format("2006-01-02")
			monthReturn = apReturnsMap[vatExtractDetail.DocumentNo].InvoiceDate.Month()
		}

		datas = append(datas, entity.VatExtractDowloadVatInReturn{
			Npwp:              vatExtractDetail.Npwp,
			Nama:              vatExtractDetail.SupName,
			KdJenisTransaksi:  faktur.KdJjenisTansaksi,
			FgPengganti:       faktur.FgPengganti,
			NomorFaktur:       faktur.NomorFaktur,
			TanggalFaktur:     taxInvoiceDate,
			IsCreditable:      1,
			NomorDokumenRetur: vatExtractDetail.DocumentNo,
			TanggalRetur:      returnDate,
			TahunPajak:        faktur.TahunPajak,
			MasaPajakRetur:    int(monthReturn),
			NilaiReturDPP:     *vatExtractDetail.SubTotal,
			NilaiReturPPN:     *vatExtractDetail.VatValue,
			NilaiReturPPNBM:   *vatExtractDetail.VatLgValue,
		})
	}

	return
}
