package service

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type SupplierReturnService interface {
	Store(request *entity.CreateSupplierReturnBody) (response entity.SupplierReturnResponse, err error)
	List(dataFilter entity.SupplierReturnQueryFilter) (data []entity.SupplierReturnListResponse, total int64, lastPage int, err error)
	Detail(supplierReturnNo, custId, parentCustId string) (response entity.SupplierReturnGetResp, err error)
	ListSupplier(dataFilter entity.ReturnSupplierQueryFilter, custId, parentCustId string) (data []entity.SupplierReturnSupplierListResponse, total int64, lastPage int, err error)
	UpdateStatus(supplierReturnNo string, request entity.UpdateSupplierReturnStatusBody) (err error)
}

func NewSupplierReturnService(returnRepository repository.SupplierReturnRepository, transaction repository.Dbtransaction, stockRepository repository.StockRepository) *SupplierReturnServiceImpl {
	return &SupplierReturnServiceImpl{
		SupplierReturnRepository: returnRepository,
		Transaction:              transaction,
		StockRepository:          stockRepository,
	}
}

type SupplierReturnServiceImpl struct {
	SupplierReturnRepository repository.SupplierReturnRepository
	Transaction              repository.Dbtransaction
	StockRepository          repository.StockRepository
}

func (service *SupplierReturnServiceImpl) Store(request *entity.CreateSupplierReturnBody) (response entity.SupplierReturnResponse, err error) {
	c := context.Background()

	if len(request.Details) == 0 {
		return response, errors.New("item details is required")
	}

	// parse time format YYYY-mm-dd to Rfc3339

	// End parse time format YYYY-mm-dd to Rfc339
	var SupplierReturnModel *model.SupplierReturn
	err = structs.Automapper(request, &SupplierReturnModel)
	if err != nil {
		return response, err
	}
	returnDate := str.GetJakartaDate()
	SupplierReturnModel.SupplierReturnDate = &returnDate
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		_, err := service.SupplierReturnRepository.GetInvoiceFromAPReturn(request.CustID, request.InvoiceNo)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			apDet, err := service.SupplierReturnRepository.GetApProductList(request.InvoiceNo, request.CustID, detail.ProID)
			if err != nil {
				return err
			}
			detail.Disc = apDet.Disc
			detail.ConvUnit2 = apDet.ConvUnit2
			detail.ConvUnit3 = apDet.ConvUnit3
			detail.UnitPrice1 = apDet.UnitPrice1
			detail.UnitPrice2 = apDet.UnitPrice2
			detail.UnitPrice3 = apDet.UnitPrice3
			detail.QtyRemaining = apDet.QtyRemaining
			detail.Vat = apDet.Vat
			detail.VatLg = apDet.VatLg
			detail.VatBg = apDet.VatBg
			detail.Calculate()
		}

		request.Calculate()

		SupplierReturnModel.DiscountValue = &request.DiscountValue
		SupplierReturnModel.SubTotal = &request.Subtotal
		SupplierReturnModel.Total = &request.Total
		SupplierReturnModel.VatValue = &request.TotalVatValue
		SupplierReturnModel.VatLgValue = &request.TotalVatLgValue
		SupplierReturnModel.VatBgValue = &request.TotalVatBgValue

		err = service.SupplierReturnRepository.Store(txCtx, SupplierReturnModel)
		if err != nil {
			return err
		}
		response.SupplierReturnNo = SupplierReturnModel.SupplierReturnNo
		for index, detail := range request.Details {
			if detail.Qty > detail.QtyRemaining {
				return errors.New(fmt.Sprintf("invalid return qty product %v", detail.ProID))
			}

			qtyToUpdateBalance := detail.QtyRemaining - detail.Qty
			// process gr balance
			invoicebalance := model.ProductInvoiceBalances{
				CustID:    request.CustID,
				InvoiceNo: request.InvoiceNo,
				ProID:     detail.ProID,
				Qty:       int(qtyToUpdateBalance),
			}
			service.SupplierReturnRepository.SaveRemainingQty(txCtx, invoicebalance)

			// if qty 0, check gr can return
			if invoicebalance.Qty == 0 {
				remainingQtyProducts, err := service.SupplierReturnRepository.GetRemainingProductQtyByInvoiceNo(request.InvoiceNo)
				if err != nil {
					return err
				}

				var grCanReturn bool
				for _, remainingQtyProduct := range remainingQtyProducts {
					if remainingQtyProduct.QtyRemaining > 0 {
						grCanReturn = true
						break
					}
				}

				if !grCanReturn {
					err := service.SupplierReturnRepository.SetInvoiceNoIsCanReturn(txCtx, request.InvoiceNo, false)
					if err != nil {
						return err
					}
				}
			}

			var SupplierReturnDetailModel model.SupplierReturnDet
			seq := index + 1

			err = structs.Automapper(detail, &SupplierReturnDetailModel)
			if err != nil {
				return err
			}

			SupplierReturnDetailModel.SupplierReturnNo = SupplierReturnModel.SupplierReturnNo
			SupplierReturnDetailModel.SeqNo = seq
			SupplierReturnDetailModel.CustID = request.CustID
			SupplierReturnDetailModel.Qty = int(detail.Qty)
			SupplierReturnDetailModel.UnitPrice1 = detail.UnitPrice1
			SupplierReturnDetailModel.UnitPrice2 = detail.UnitPrice2
			SupplierReturnDetailModel.UnitPrice3 = detail.UnitPrice3
			SupplierReturnDetailModel.SubTotal = detail.Subtotal
			SupplierReturnDetailModel.Discount = detail.Disc
			SupplierReturnDetailModel.DiscountValue = detail.DiscValue
			SupplierReturnDetailModel.Vat = detail.Vat
			SupplierReturnDetailModel.VatValue = detail.VatValue
			SupplierReturnDetailModel.VatLg = detail.VatLg
			SupplierReturnDetailModel.VatLgValue = detail.VatLgValue
			SupplierReturnDetailModel.VatBg = detail.VatBg
			SupplierReturnDetailModel.VatBgValue = detail.VatBgValue
			SupplierReturnDetailModel.Total = detail.Total
			_, err = service.SupplierReturnRepository.StoreDetail(txCtx, &SupplierReturnDetailModel)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return response, err
}

func (service *SupplierReturnServiceImpl) List(dataFilter entity.SupplierReturnQueryFilter) (data []entity.SupplierReturnListResponse, total int64, lastPage int, err error) {
	supplierReturns, total, lastPage, err := service.SupplierReturnRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}
	for _, row := range supplierReturns {
		var vResp entity.SupplierReturnListResponse
		structs.Automapper(row, &vResp)
		supplierReturnDate := row.SupplierReturnDate.Format("2006-01-02")

		if row.InvoiceDate != nil {
			vResp.InvoiceDate = row.InvoiceDate.Format("2006-01-02")
		}

		if row.TaxInvoiceDate != nil {
			vResp.TaxInvoiceDate = row.TaxInvoiceDate.Format("2006-01-02")
		}

		if row.DueDate != nil {
			vResp.DueDate = row.DueDate.Format("2006-01-02")
		}

		vResp.SupplierReturnDate = supplierReturnDate
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *SupplierReturnServiceImpl) Detail(supplierReturnNo, custId, parentCustId string) (response entity.SupplierReturnGetResp, err error) {
	supplierReturn, err := service.SupplierReturnRepository.FindByNo(supplierReturnNo, custId, parentCustId)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(supplierReturn, &response)
	if err != nil {
		return response, err
	}
	supplierReturnDate := supplierReturn.SupplierReturnDate.Format("2006-01-02")
	response.SupplierReturnDate = supplierReturnDate

	if supplierReturn.InvoiceDate != nil {
		response.InvoiceDate = supplierReturn.InvoiceDate.Format("2006-01-02")
	}

	if supplierReturn.TaxInvoiceDate != nil {
		response.TaxInvoiceDate = supplierReturn.TaxInvoiceDate.Format("2006-01-02")
	}

	if supplierReturn.DueDate != nil {
		response.DueDate = supplierReturn.DueDate.Format("2006-01-02")
	}

	response.TotalSkuPrice = response.SubTotal - response.DiscountValue
	response.SubTotal = response.SubTotal - response.DiscountValue
	supplierReturnDets, err := service.SupplierReturnRepository.FindDetBySupplierReturnNo(supplierReturnNo, custId)
	if err != nil {
		return response, err
	}

	var supplierReturnDetailsData []entity.SupplierReturnGetDetResp
	for _, supplierReturnDet := range supplierReturnDets {
		var supplierReturnDetailData entity.SupplierReturnGetDetResp
		err = structs.Automapper(supplierReturnDet, &supplierReturnDetailData)
		if err != nil {
			return response, err
		}

		qtyInvoice := &conversion.Qty{
			Qty:       int(supplierReturnDet.InvoiceQty),
			ConvUnit2: int(supplierReturnDet.ConvUnit2),
			ConvUnit3: int(supplierReturnDet.ConvUnit3),
		}
		qtyInvoiceConversion := qtyInvoice.ConvToQtyConversion()

		qtyReturn := &conversion.Qty{
			Qty:       int(supplierReturnDet.Qty),
			ConvUnit2: int(supplierReturnDet.ConvUnit2),
			ConvUnit3: int(supplierReturnDet.ConvUnit3),
		}
		qtyReturnConversion := qtyReturn.ConvToQtyConversion()

		qtyRemaining := &conversion.Qty{
			Qty:       int(supplierReturnDet.RemainingQty),
			ConvUnit2: int(supplierReturnDet.ConvUnit2),
			ConvUnit3: int(supplierReturnDet.ConvUnit3),
		}
		qtyRemainingConversion := qtyRemaining.ConvToQtyConversion()

		qtyWarehouse := &conversion.Qty{
			Qty:       int(supplierReturnDet.WhQty),
			ConvUnit2: int(supplierReturnDet.ConvUnit2),
			ConvUnit3: int(supplierReturnDet.ConvUnit3),
		}
		qtyWarehouseConversion := qtyWarehouse.ConvToQtyConversion()

		remainingQty1 := qtyRemainingConversion.Qty1
		remainingQty2 := qtyRemainingConversion.Qty2
		remainingQty3 := qtyRemainingConversion.Qty3

		supplierReturnDetailData.Qty1 = qtyReturnConversion.Qty1
		supplierReturnDetailData.Qty2 = qtyReturnConversion.Qty2
		supplierReturnDetailData.Qty3 = qtyReturnConversion.Qty3
		supplierReturnDetailData.RemainingQty1 = remainingQty1
		supplierReturnDetailData.RemainingQty2 = remainingQty2
		supplierReturnDetailData.RemainingQty3 = remainingQty3
		supplierReturnDetailData.InvoiceQty1 = qtyInvoiceConversion.Qty1
		supplierReturnDetailData.InvoiceQty2 = qtyInvoiceConversion.Qty2
		supplierReturnDetailData.InvoiceQty3 = qtyInvoiceConversion.Qty3
		supplierReturnDetailData.ConvUnit1 = supplierReturnDetailData.ConvUnit2 * supplierReturnDetailData.ConvUnit3
		supplierReturnDetailData.SubTotal = supplierReturnDet.SubTotal
		supplierReturnDetailData.Nett = supplierReturnDet.SubTotal - supplierReturnDetailData.DiscountValue
		supplierReturnDetailData.WhQty1 = qtyWarehouseConversion.Qty1
		supplierReturnDetailData.WhQty2 = qtyWarehouseConversion.Qty2
		supplierReturnDetailData.WhQty3 = qtyWarehouseConversion.Qty3

		supplierReturnDetailsData = append(supplierReturnDetailsData, supplierReturnDetailData)

	}
	response.Details = supplierReturnDetailsData
	return response, nil
}

func (service *SupplierReturnServiceImpl) ListSupplier(dataFilter entity.ReturnSupplierQueryFilter, custId, parentCustId string) (data []entity.SupplierReturnSupplierListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.SupplierReturnRepository.FindSupplierReturn(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.SupplierReturnSupplierListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *SupplierReturnServiceImpl) UpdateStatus(supplierReturnNo string, request entity.UpdateSupplierReturnStatusBody) (err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.SupplierReturnRepository.UpdateStatus(txCtx, supplierReturnNo, request.CustID, request.DataStatus)
		if err != nil {
			return err
		}

		supplierReturn, err := service.SupplierReturnRepository.FindByNo(supplierReturnNo, request.CustID, request.ParentCustID)
		if err != nil {
			return err
		}
		supplierReturnDets, err := service.SupplierReturnRepository.FindDetBySupplierReturnNo(supplierReturnNo, request.CustID)
		if err != nil {
			return err
		}

		if request.DataStatus == 2 {
			var stockUpdateEntities []*entity.StockUpdate
			for _, supplierReturnDet := range supplierReturnDets {
				stockUpdateEntity := entity.StockUpdate{
					CustID:    request.CustID,
					WhID:      *supplierReturn.WhID,
					ProID:     int64(supplierReturnDet.ProID),
					StockDate: *supplierReturn.SupplierReturnDate,
					TrCode:    supplierReturn.SupplierReturnNo[0:2],
					TrNo:      supplierReturn.SupplierReturnNo,
					QtyIn:     0,
					QtyOut:    float64(supplierReturnDet.Qty),
					UnitPrice: supplierReturnDet.UnitPrice1,
					RefDetId:  supplierReturnDet.ID,
				}

				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)

			}

			err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			if err != nil {
				return err
			}

		} else {
			for _, supplierReturnDet := range supplierReturnDets {
				remainingQty, err := service.SupplierReturnRepository.GetRemainingQtyInvoice(*supplierReturn.InvoiceNo, request.CustID, int64(supplierReturnDet.ProID))
				if err != nil {
					return err
				}

				// process gr balance
				grbalance := model.ProductInvoiceBalances{
					CustID:    request.CustID,
					InvoiceNo: *supplierReturn.InvoiceNo,
					ProID:     int64(supplierReturnDet.ProID),
					Qty:       int(remainingQty.RemainingQty) + int(supplierReturnDet.Qty),
				}
				service.SupplierReturnRepository.SaveRemainingQty(txCtx, grbalance)
			}

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
