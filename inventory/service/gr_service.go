package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/errmsg"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type GrService interface {
	Detail(grNo string, custID, parentCustId string, isAp bool) (response entity.GrWithDetailResponse, err error)
	DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.GrWithDetailResponse, err error)
	Store(request entity.CreateGrBody) (response entity.GrResponse, err error)
	List(dataFilter entity.GrQueryFilter, custId, parentCustId string) (data []entity.GrListResponse, total int64, lastPage int, err error)
	// Update(grNo string, request entity.UpdateGrRequest) (err error)
	// Delete(custId string, grNo string, userId int64) (err error)
	ListSupplier(dataFilter entity.GrSupplierQueryFilter, custId, parentCustId string) (data []entity.GrSupplierListResponse, total int64, lastPage int, err error)
	ListWarehouse(dataFilter entity.GrWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrWarehouseListResponse, total int64, lastPage int, err error)
	ListDistributor(custId, parentCustId string) (data []entity.GrDistributorListResponse, err error)
	ListLookupGrAp(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) (data []entity.GrLookupResponse, total int64, lastPage int, err error)
	Download(grNo string, custId, parentCustId string) (response entity.GrDownloadResponse, err error)
}

func NewGrService(
	grRepository repository.GrRepository,
	warehouseStockRepository repository.WarehouseStockRepository,
	stockRepository repository.StockRepository,
	replenishmentRepository repository.ReplenishmentRepository,
	transaction repository.Dbtransaction) *grServiceImpl {
	return &grServiceImpl{
		GrRepository:             grRepository,
		WarehouseStockRepository: warehouseStockRepository,
		StockRepository:          stockRepository,
		ReplenishmentRepository:  replenishmentRepository,
		Transaction:              transaction,
	}
}

type grServiceImpl struct {
	GrRepository             repository.GrRepository
	WarehouseStockRepository repository.WarehouseStockRepository
	StockRepository          repository.StockRepository
	ReplenishmentRepository  repository.ReplenishmentRepository
	Transaction              repository.Dbtransaction
}

func (service *grServiceImpl) Store(request entity.CreateGrBody) (response entity.GrResponse, err error) {
	c := context.Background()

	// check invoice
	if request.InvoiceNo != "" {
		grExisting, _ := service.GrRepository.GetByInvoiceNo(request.InvoiceNo, request.CustID, request.ParentCustID)
		if grExisting.GrNo != "" {
			return response, errors.New(errmsg.ERROR_INVOICE_ALREADY_USED)
		}
	}

	// return response, errors.New("test")
	if len(request.Details.Normal) == 0 {
		return response, errors.New("item details is required")
	}

	// parse time format YYYY-mm-dd to Rfc3339
	deliveryDate, err := str.DateStrToRfc3339String(request.DeliveryDate)
	if err != nil {
		return response, err
	}
	request.DeliveryDate = deliveryDate

	if request.InvoiceDate != nil {
		if *request.InvoiceDate != "" {
			invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
			if err != nil {
				return response, err
			}
			request.InvoiceDate = &invoiceDate
		} else {
			request.InvoiceDate = nil
		}

	}
	// End parse time format YYYY-mm-dd to Rfc339

	var grModel model.Gr
	err = structs.Automapper(request, &grModel)
	if err != nil {
		return response, err
	}
	grDate := str.GetJakartaDate()

	grModel.GrDate = &grDate

	// Set good_receipt_type: only set if provided in request (optional field)
	if request.GoodReceiptType != nil && *request.GoodReceiptType != "" {
		grModel.GoodReceiptType = request.GoodReceiptType
	}
	// If not provided, leave it as nil (will be NULL in database)

	if request.WithReference != nil {
		grModel.WithReference = request.WithReference
	}

	if request.DeliveryFee != nil {
		grModel.DeliveryFee = request.DeliveryFee
	}

	if request.SoNo != nil && *request.SoNo != "" {
		grModel.SoNo = request.SoNo
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		err := service.GrRepository.Store(txCtx, &grModel)
		if err != nil {
			return err
		}
		response.GrNo = grModel.GrNo

		var productIDs []int64
		for _, detail := range request.Details.Normal {
			productIDs = append(productIDs, detail.ProID)
		}
		for _, detail := range request.Details.Promo {
			if !slices.Contains(productIDs, detail.ProID) {
				productIDs = append(productIDs, detail.ProID)
			}
		}

		productsModel, err := service.GrRepository.FindProductByListID(request.CustID, request.ParentCustID, request.DistributorID, productIDs)
		if err != nil {
			return err
		}

		var productMap = model.MapProduct{}

		for _, productModel := range productsModel {
			productMap.SetProduct(productModel.ProductId, productModel)
		}

		var stockUpdateEntities []*entity.StockUpdate

		for index, detail := range request.Details.Normal {
			productModel, err := productMap.GetByID(detail.ProID)
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      detail.Qty1,
				Qty2:      detail.Qty2,
				Qty3:      detail.Qty3,
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			QtyShipUnit := &conversion.QtyUnit{
				Qty1:      detail.QtyShip1,
				Qty2:      detail.QtyShip2,
				Qty3:      detail.QtyShip3,
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQtyShip, err := QtyShipUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			var grDetailModel model.GrDetCreate
			// seq := index + 1
			if detail.Qty1 == 0 && detail.Qty2 == 0 && detail.Qty3 == 0 {
				return errors.New(fmt.Sprintf("please input qty on id_product %v", detail.ProID))
			}

			if detail.QtyShip1 == 0 && detail.QtyShip2 == 0 && detail.QtyShip3 == 0 {
				return errors.New(fmt.Sprintf("please input qty_ship on id_product %v", detail.ProID))
			}

			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}

			// Unit price: use payload if provided (>0), otherwise fallback to master product price
			unitPrice1 := productModel.PurchPrice1
			if detail.UnitPrice1 != nil && *detail.UnitPrice1 > 0 {
				unitPrice1 = *detail.UnitPrice1
			}
			unitPrice2 := productModel.PurchPrice2
			if detail.UnitPrice2 != nil && *detail.UnitPrice2 > 0 {
				unitPrice2 = *detail.UnitPrice2
			}
			unitPrice3 := productModel.PurchPrice3
			if detail.UnitPrice3 != nil && *detail.UnitPrice3 > 0 {
				unitPrice3 = *detail.UnitPrice3
			}

			grDetailModel.CustID = request.CustID
			grDetailModel.GrNo = grModel.GrNo
			grDetailModel.ItemType = model.ITEM_TYPE_NORMAL
			grDetailModel.SeqNo = index + 1
			grDetailModel.UnitPrice1 = unitPrice1
			grDetailModel.UnitPrice2 = unitPrice2
			grDetailModel.UnitPrice3 = unitPrice3
			grDetailModel.Qty1 = detail.Qty1
			grDetailModel.Qty2 = detail.Qty2
			grDetailModel.Qty3 = detail.Qty3
			grDetailModel.QtyShip1 = detail.QtyShip1
			grDetailModel.QtyShip2 = detail.QtyShip2
			grDetailModel.QtyShip3 = detail.QtyShip3
			grDetailModel.Qty = totalQty
			grDetailModel.QtyShip = totalQtyShip

			grDet, err := service.GrRepository.CreateGrDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			stockUpdateEntity := entity.StockUpdate{
				CustID:    grDetailModel.CustID,
				WhID:      *grModel.WhID,
				ProID:     detail.ProID,
				StockDate: *grModel.GrDate,
				TrCode:    grModel.GrNo[0:2],
				TrNo:      grModel.GrNo,
				QtyIn:     float64(grDet.Qty),
				UnitPrice: grDetailModel.UnitPrice1,
				RefDetId:  grDet.ID,
			}

			stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)

		}

		for index, detail := range request.Details.Promo {
			productModel, err := productMap.GetByID(detail.ProID)
			if err != nil {
				return err
			}
			QtyUnit := &conversion.QtyUnit{
				Qty1:      detail.Qty1,
				Qty2:      detail.Qty2,
				Qty3:      detail.Qty3,
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			var grDetailModel model.GrDetCreate

			if detail.Qty1 == 0 && detail.Qty2 == 0 && detail.Qty3 == 0 {
				return errors.New(fmt.Sprintf("please input qty on id_product %v", detail.ProID))
			}

			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}
			// grDetailModel.SeqNo = seq
			grDetailModel.CustID = request.CustID
			grDetailModel.GrNo = grModel.GrNo
			grDetailModel.ItemType = model.ITEM_TYPE_PROMO
			grDetailModel.SeqNo = index + 1
			// Unit price: use payload if provided (>0), otherwise fallback to master product price
			unitPrice1 := productModel.PurchPrice1
			if detail.UnitPrice1 != nil && *detail.UnitPrice1 > 0 {
				unitPrice1 = *detail.UnitPrice1
			}
			unitPrice2 := productModel.PurchPrice2
			if detail.UnitPrice2 != nil && *detail.UnitPrice2 > 0 {
				unitPrice2 = *detail.UnitPrice2
			}
			unitPrice3 := productModel.PurchPrice3
			if detail.UnitPrice3 != nil && *detail.UnitPrice3 > 0 {
				unitPrice3 = *detail.UnitPrice3
			}
			grDetailModel.UnitPrice1 = unitPrice1
			grDetailModel.UnitPrice2 = unitPrice2
			grDetailModel.UnitPrice3 = unitPrice3
			grDetailModel.QtyShip1 = detail.QtyShip1
			grDetailModel.QtyShip2 = detail.QtyShip2
			grDetailModel.QtyShip3 = detail.QtyShip3
			grDetailModel.Qty = totalQty
			grDet, err := service.GrRepository.CreateGrDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			stockUpdateEntity := entity.StockUpdate{
				CustID:    grDetailModel.CustID,
				WhID:      *grModel.WhID,
				ProID:     detail.ProID,
				StockDate: *grModel.GrDate,
				TrCode:    grModel.GrNo[0:2],
				TrNo:      grModel.GrNo,
				QtyIn:     float64(grDet.Qty),
				UnitPrice: grDetailModel.UnitPrice1,
				RefDetId:  grDet.ID,
			}

			stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
		}

		err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
		if err != nil {
			return err
		}

		// Update replenishment status to Completed if with_reference is true and good_receipt_type is Replenishment or Replenishment Event
		if request.WithReference != nil && *request.WithReference &&
			request.GoodReceiptType != nil && (*request.GoodReceiptType == "Replenishment" || *request.GoodReceiptType == "Replenishment Event") &&
			request.PoNo != "" {
			err = service.ReplenishmentRepository.UpdateStatusByReplenishmentNo(txCtx, request.PoNo, request.CustID, StatusCompleted, request.UpdatedBy)
			if err != nil {
				return fmt.Errorf("failed to update replenishment status: %w", err)
			}
		}

		return nil
	})

	return response, err
}

func (service *grServiceImpl) Detail(grNo string, custID, parentCustId string, isAp bool) (response entity.GrWithDetailResponse, err error) {
	gr, err := service.GrRepository.FindByNo(grNo, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gr, &response)
	if err != nil {
		return response, err
	}

	err = service.GetGrDetail(grNo, custID, &response, isAp)
	if err != nil {
		return response, err
	}

	grData := gr.GrDate.Format("2006-01-02")
	deliveryDate := gr.DeliveryDate.Format("2006-01-02")
	if gr.InvoiceDate != nil {
		invoiceDate := gr.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = invoiceDate
	}

	response.GrDate = grData
	response.DeliveryDate = deliveryDate

	if gr.GoodReceiptType != nil {
		response.GoodReceiptType = *gr.GoodReceiptType
	}
	if gr.SoNo != nil {
		response.SoNo = *gr.SoNo
	}
	if gr.DeliveryFee != nil {
		response.DeliveryFee = *gr.DeliveryFee
	}

	return response, nil
}

func (service *grServiceImpl) DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.GrWithDetailResponse, err error) {
	gr, err := service.GrRepository.GetByInvoiceNo(invoice, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gr, &response)
	if err != nil {
		return response, err
	}

	err = service.GetGrDetail(invoice, custID, &response, isAp)
	if err != nil {
		return response, err
	}

	grData := gr.GrDate.Format("2006-01-02")
	deliveryDate := gr.DeliveryDate.Format("2006-01-02")
	invoiceDate := gr.InvoiceDate.Format("2006-01-02")

	response.GrDate = grData
	response.DeliveryDate = deliveryDate
	response.InvoiceDate = invoiceDate
	return response, nil
}

func (service *grServiceImpl) GetGrDetail(grNo string, custID string, gr *entity.GrWithDetailResponse, isAp bool) (err error) {
	var grDetails []model.GrDetList
	var discountValueTotal float64
	if !isAp {
		grDetails, err = service.GrRepository.FindGrdetail(grNo, custID)
		if err != nil {
			return err
		}
	} else {
		grDetails, err = service.GrRepository.FindGrdetailWithDiscount(grNo, custID)
		if err != nil {
			return err
		}
	}

	gr.Details.Promo = []entity.GrDetList{}
	gr.Details.Normal = []entity.GrDetList{}

	for _, grDetail := range grDetails {
		var grDetailData entity.GrDetList
		err = structs.Automapper(grDetail, &grDetailData)
		if err != nil {
			return err
		}

		qty := &conversion.Qty{
			Qty:       int(grDetailData.Qty),
			ConvUnit2: int(grDetailData.ConvUnit2),
			ConvUnit3: int(grDetailData.ConvUnit3),
		}

		qtyConversion := qty.ConvToQtyConversion()

		qtyRemaining := &conversion.Qty{
			Qty:       int(grDetail.QtyRemaining),
			ConvUnit2: int(grDetailData.ConvUnit2),
			ConvUnit3: int(grDetailData.ConvUnit3),
		}
		qtyRemainingConversion := qtyRemaining.ConvToQtyConversion()

		qtyWarehouse := &conversion.Qty{
			Qty:       int(grDetail.WhQty),
			ConvUnit2: int(grDetailData.ConvUnit2),
			ConvUnit3: int(grDetailData.ConvUnit3),
		}
		qtyWarehouseConversion := qtyWarehouse.ConvToQtyConversion()

		grDetailData.Qty1 = qtyConversion.Qty1
		grDetailData.Qty2 = qtyConversion.Qty2
		grDetailData.Qty3 = qtyConversion.Qty3

		// Populate qty_ship per unit from stored columns; fallback to conversion if nil
		if grDetail.QtyShip1 != nil {
			grDetailData.QtyShip1 = int(*grDetail.QtyShip1)
		} else {
			grDetailData.QtyShip1 = 0
		}
		if grDetail.QtyShip2 != nil {
			grDetailData.QtyShip2 = int(*grDetail.QtyShip2)
		} else {
			grDetailData.QtyShip2 = 0
		}
		if grDetail.QtyShip3 != nil {
			grDetailData.QtyShip3 = int(*grDetail.QtyShip3)
		} else {
			grDetailData.QtyShip3 = 0
		}

		grDetailData.QtyRemaining1 = qtyRemainingConversion.Qty1
		grDetailData.QtyRemaining2 = qtyRemainingConversion.Qty2
		grDetailData.QtyRemaining3 = qtyRemainingConversion.Qty3

		grDetailData.WhQty1 = qtyWarehouseConversion.Qty1
		grDetailData.WhQty2 = qtyWarehouseConversion.Qty2
		grDetailData.WhQty3 = qtyWarehouseConversion.Qty3

		if grDetail.ItemType == 1 {
			var discountValue, discount float64
			Subtotal := (grDetail.UnitPrice1 * float64(grDetailData.Qty1)) + (grDetail.UnitPrice2 * float64(grDetailData.Qty2)) + (grDetail.UnitPrice3 * float64(grDetailData.Qty3))
			if isAp {
				if grDetail.Discount != nil {
					discountValue = (*grDetail.Discount / 100) * Subtotal
					discount = *grDetail.Discount
				}
				grDetailData.DiscountValue = &discountValue
				grDetailData.Discount = &discount

				discountValueTotal += discountValue
			}

			ppn := (Subtotal - discountValue) * grDetail.Vat / 100
			pbn := (Subtotal - discountValue) * grDetail.VatLgPurch / 100
			//ppnDp := (Subtotal - discountValue) * grDetail.VatBg / 100
			// total := (Subtotal - discountValue) + ppn + pbn
			total := Subtotal + ppn + gr.DeliveryFee

			gr.SubTotal += Subtotal - discountValue
			gr.Total += total
			gr.TotalVat += ppn
			gr.TotalVatLgPurch += pbn
			gr.TotalSkuPrice += Subtotal - discountValue
			grDetailData.Nett = Subtotal - discountValue
			grDetailData.SubTotal = Subtotal
			grDetailData.Total = total
			grDetailData.VatValue = ppn
			grDetailData.VatLgPurchValue = pbn

			grDetailData.ConvUnit1 = grDetail.ConvUnit2 * grDetail.ConvUnit3
			gr.Details.Normal = append(gr.Details.Normal, grDetailData)
		} else {
			promoPrice := (grDetail.UnitPrice1 * float64(grDetailData.Qty1)) +
				(grDetail.UnitPrice2 * float64(grDetailData.Qty2)) +
				(grDetail.UnitPrice3 * float64(grDetailData.Qty3))
			grDetailData.PromoPrice = &promoPrice
			gr.Details.Promo = append(gr.Details.Promo, grDetailData)
		}
	}

	if gr.Details.Promo == nil {
		gr.Details.Promo = []entity.GrDetList{}
	}

	if isAp {
		gr.DiscountValue = &discountValueTotal
	}
	return
}

func (service *grServiceImpl) List(dataFilter entity.GrQueryFilter, custId, parentCustId string) (data []entity.GrListResponse, total int64, lastPage int, err error) {
	grs, total, lastPage, err := service.GrRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grs {
		var vResp entity.GrListResponse
		structs.Automapper(row, &vResp)
		grData := row.GrDate.Format("2006-01-02")
		deliveryDate := row.DeliveryDate.Format("2006-01-02")
		if row.WithReference != nil {
			vResp.WithReference = *row.WithReference
		}
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = invoiceDate
		}

		vResp.GrDate = grData
		vResp.DeliveryDate = deliveryDate
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *grServiceImpl) ListSupplier(dataFilter entity.GrSupplierQueryFilter, custId, parentCustId string) (data []entity.GrSupplierListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.GrRepository.FindSupplierGr(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.GrSupplierListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *grServiceImpl) ListWarehouse(dataFilter entity.GrWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrWarehouseListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.GrRepository.FindWarehouseGr(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.GrWarehouseListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

// func (service *grServiceImpl) Update(grNo string, request entity.UpdateGrRequest) (err error) {
// 	c := context.Background()

// 	if request.DeliveryDate != nil {
// 		// parse time format YYYY-mm-dd to Rfc3339
// 		if *request.DeliveryDate != "" {
// 			deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
// 			if err != nil {
// 				return err
// 			}
// 			request.DeliveryDate = &deliveryDate
// 		}
// 	}

// 	var grModel model.Gr
// 	err = structs.Automapper(request, &grModel)
// 	if err != nil {
// 		return err
// 	}
// 	grModel.CustID = ""
// 	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
// 		gr, err := service.GrRepository.FindByNo(grNo, request.CustID, request.ParentCustID)
// 		if err != nil {
// 			return err
// 		}

// 		err = service.GrRepository.Update(txCtx, gr.GrNo, grModel)
// 		if err != nil {
// 			return err
// 		}

// 		// grDetails, err := service.GrRepository.FindGrdetail(grNo, request.CustID)
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		// for _, grDet := range grDetails {
// 		// 	err = service.WhStockRepository.UpdateOldWhStock(txCtx, request.CustID, *request.WhID, grDet.ProID, grDet.Qty)
// 		// 	if err != nil {
// 		// 		return err
// 		// 	}
// 		// }

// 		err = service.GrRepository.DeleteGrDetailByGrNo(txCtx, grNo)
// 		if err != nil {
// 			return err
// 		}

// 		newGrDetIds := make([]int64, 0)
// 		for _, detail := range request.Details.Normal {
// 			/* if detail.GrDetId == nil || *detail.GrDetId == 0 { */
// 			detail.GrDetId = nil
// 			var grDetailModel model.GrDetCreate
// 			err = structs.Automapper(detail, &grDetailModel)
// 			if err != nil {
// 				return err
// 			}
// 			// grDetailModel.SeqNo = sequence
// 			grDetailModel.CustID = request.CustID
// 			grDetailModel.GrNo = grNo
// 			// grDet, err := service.GrRepository.CreateGrDetail(txCtx, &grDetailModel)
// 			grDet, err := service.GrRepository.CreateGrDetail(txCtx, &grDetailModel)
// 			if err != nil {
// 				return err
// 			}

// 			newGrDetIds = append(newGrDetIds, grDet.ID)
// 		}

// 		for _, detail := range request.Details.Promo {
// 			detail.GrDetId = nil
// 			var grDetailModel model.GrDetCreate
// 			err = structs.Automapper(detail, &grDetailModel)
// 			if err != nil {
// 				return err
// 			}
// 			grDetailModel.CustID = request.CustID
// 			grDetailModel.GrNo = grNo
// 			grDet, err := service.GrRepository.CreateGrDetail(txCtx, &grDetailModel)
// 			if err != nil {
// 				return err
// 			}

// 			newGrDetIds = append(newGrDetIds, grDet.ID)
// 		}

// 		err = service.GrRepository.DeleteStockNotInRefIds(txCtx, request.CustID, gr.GrNo, newGrDetIds)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (service *grServiceImpl) Delete(custId string, grNo string, userId int64) (err error) {
// 	c := context.Background()
// 	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
// 		grDetails, err := service.GrRepository.FindGrdetailForUpdateWhStock(grNo, custId)
// 		if err != nil {
// 			return err
// 		}
// 		log.Println("grDetails:", structs.StructToJson(grDetails))

// 		oldGrDetId := make([]int64, 0)
// 		for _, grDet := range grDetails {
// 			err = service.GrRepository.UpdateOldWhStock(txCtx, custId, grDet.WhID, grDet.ProID, *grDet.Qty)
// 			if err != nil {
// 				return err
// 			}
// 			oldGrDetId = append(oldGrDetId, grDet.ID)
// 		}

// 		log.Println("oldGrDetId:", structs.StructToJson(oldGrDetId))
// 		log.Println("oldGrDetId:", structs.StructToJson(oldGrDetId))
// 		if len(oldGrDetId) > 0 {
// 			err = service.GrRepository.DeleteStockInRefIds(txCtx, custId, grNo, oldGrDetId)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		err = service.GrRepository.Delete(txCtx, custId, grNo, userId)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})

// 	return err
// }

func (service *grServiceImpl) ListDistributor(custId, parentCustId string) (data []entity.GrDistributorListResponse, err error) {
	grDistributor, err := service.GrRepository.FindDistributorGr(custId)
	if err != nil {
		return data, err
	}

	for _, row := range grDistributor {
		var vResp entity.GrDistributorListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, err
}

func (service *grServiceImpl) ListLookupGrAp(dataFilter entity.GrLookupQueryFilter, custId, parentCustId string) (data []entity.GrLookupResponse, total int64, lastPage int, err error) {

	gr, total, lastPage, err := service.GrRepository.FindAllGrByCustIdSupId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	grb, total, lastPage, err := service.GrRepository.FindAllGrBranchByCustIdSupId(dataFilter, dataFilter.CustIdParam, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(custId) < 10 {

		for _, row := range gr {
			var vResp entity.GrLookupResponse
			structs.Automapper(row, &vResp)

			vResp.GrNo = row.GrNo
			// fmt.Println("gr", vResp.GrNo)
			data = append(data, vResp)
		}

		for _, row := range grb {
			var vResp entity.GrLookupResponse
			structs.Automapper(row, &vResp)

			vResp.GrNo = row.GrBranchNo
			// fmt.Println("grb", vResp.GrNo)
			data = append(data, vResp)
		}

	} else {
		for _, row := range gr {
			var vResp entity.GrLookupResponse
			structs.Automapper(row, &vResp)

			vResp.GrNo = row.GrNo
			data = append(data, vResp)
		}

	}

	return data, total, lastPage, err
}

func (service *grServiceImpl) Download(grNo string, custId, parentCustId string) (response entity.GrDownloadResponse, err error) {
	ctx := context.Background()

	// Check if there's a report in progress
	isInProgress, err := service.GrRepository.CheckReportInProgress(ctx, "DownloadGoodReceipt")
	if err != nil {
		return response, err
	}
	if isInProgress {
		return response, errors.New("Processing time may vary by file size. Please check Download History to access the file")
	}

	// Fetch GR data
	gr, grDetails, err := service.GrRepository.FindGrDetailForDownload(grNo, custId, parentCustId)
	if err != nil {
		return response, err
	}

	// Generate Excel file
	file, err := service.generateExcelFile(gr, grDetails)
	if err != nil {
		return response, err
	}

	// Convert to base64
	var buf bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		return response, err
	}
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Generate report name: DownloadGoodReceipt-DDMMYY-3digitRunningNumber
	now := time.Now()
	dateStr := now.Format("020106") // DDMMYY
	sequenceNumber, err := getNextSequenceNumber(dateStr)
	if err != nil {
		return response, fmt.Errorf("failed to get sequence number: %w", err)
	}
	reportName := fmt.Sprintf("DownloadGoodReceipt-%s-%03d", dateStr, sequenceNumber)

	response.ReportName = reportName
	response.FileBase64 = base64Str

	return response, nil
}

func (service *grServiceImpl) generateExcelFile(gr model.GrList, grDetails []model.GrDetList) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set sheet name
	sheetName := "Good Receipt"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Create styles
	boldStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return nil, err
	}

	boldCenterStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	boldBorderCenterStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	numberBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
		NumFmt:    3, // #,##0 format
	})
	if err != nil {
		return nil, err
	}

	textBorderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	// Goods Receipt Detail section (Row 1-6, two columns layout)
	f.SetCellValue(sheetName, "A1", "Goods Receipt Detail")
	f.SetCellStyle(sheetName, "A1", "A1", boldStyle)

	// Left column (A-B)
	f.SetCellValue(sheetName, "A2", "Goods Receipt Type")
	if gr.GoodReceiptType != nil {
		f.SetCellValue(sheetName, "B2", *gr.GoodReceiptType)
	}

	f.SetCellValue(sheetName, "A3", "PO No.")
	if gr.PoNo != nil {
		f.SetCellValue(sheetName, "B3", *gr.PoNo)
	}

	f.SetCellValue(sheetName, "A4", "Supplier")
	if gr.SupName != nil {
		f.SetCellValue(sheetName, "B4", *gr.SupName)
	}

	f.SetCellValue(sheetName, "A5", "Warehouse")
	if gr.WhName != nil {
		f.SetCellValue(sheetName, "B5", *gr.WhName)
	}

	f.SetCellValue(sheetName, "A6", "Note")
	if gr.Notes != nil {
		notes := *gr.Notes
		if notes == "" {
			notes = "-"
		}
		f.SetCellValue(sheetName, "B6", notes)
	} else {
		f.SetCellValue(sheetName, "B6", "-")
	}

	// Right column (E-F)
	f.SetCellValue(sheetName, "E2", "Delivery No.")
	if gr.DeliveryNo != nil {
		f.SetCellValue(sheetName, "F2", *gr.DeliveryNo)
	}

	f.SetCellValue(sheetName, "E3", "Delivery Date")
	if gr.DeliveryDate != nil {
		// Format: DD/MM/YYYY
		f.SetCellValue(sheetName, "F3", gr.DeliveryDate.Format("02/01/2006"))
	}

	f.SetCellValue(sheetName, "E4", "Vehicle No.")
	if gr.VehicleNo != nil {
		f.SetCellValue(sheetName, "F4", *gr.VehicleNo)
	}

	f.SetCellValue(sheetName, "E5", "SO No.")
	if gr.SoNo != nil {
		f.SetCellValue(sheetName, "F5", *gr.SoNo)
	}

	row := 8
	// Header Row 8: "Stock" (A-N) and "Free Goods" (O-U)
	f.SetCellValue(sheetName, "A8", "Stock")
	f.MergeCell(sheetName, "A8", "N8")
	f.SetCellStyle(sheetName, "A8", "N8", boldCenterStyle)
	f.SetCellValue(sheetName, "O8", "Free Goods")
	f.MergeCell(sheetName, "O8", "U8")
	f.SetCellStyle(sheetName, "O8", "U8", boldCenterStyle)

	// Header Row 9: Main headers
	row = 9
	boldBorderLeftStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	if err != nil {
		return nil, err
	}

	stockHeaders := []string{"Product Code", "Product Name", "UOM", "Qty Shipment", "", "", "Qty Received", "", "", "Price", "", "", "Sub Total", "PPn"}
	stockCols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N"}
	for i, header := range stockHeaders {
		if header != "" {
			cell := fmt.Sprintf("%s%d", stockCols[i], row)
			f.SetCellValue(sheetName, cell, header)
			if i < 3 {
				f.SetCellStyle(sheetName, cell, cell, boldBorderLeftStyle)
			} else {
				f.SetCellStyle(sheetName, cell, cell, boldBorderCenterStyle)
			}
		}
	}

	// Merge cells for Qty Shipment (D9-F9)
	f.MergeCell(sheetName, "D9", "F9")
	f.SetCellStyle(sheetName, "D9", "F9", boldBorderCenterStyle)
	// Merge cells for Qty Received (G9-I9)
	f.MergeCell(sheetName, "G9", "I9")
	f.SetCellStyle(sheetName, "G9", "I9", boldBorderCenterStyle)
	// Merge cells for Price (J9-L9)
	f.MergeCell(sheetName, "J9", "L9")
	f.SetCellStyle(sheetName, "J9", "L9", boldBorderCenterStyle)
	f.SetCellStyle(sheetName, "N9", "N9", boldBorderCenterStyle)

	// Free Goods headers (O-U)
	f.SetCellValue(sheetName, "O9", "Qty Shipment")
	f.SetCellValue(sheetName, "R9", "Qty Received")
	f.SetCellValue(sheetName, "U9", "Sub Total")

	// Merge cells for Free Goods Qty Shipment (O9-Q9)
	f.MergeCell(sheetName, "O9", "Q9")
	f.SetCellStyle(sheetName, "O9", "Q9", boldBorderCenterStyle)
	// Merge cells for Free Goods Qty Received (R9-T9)
	f.MergeCell(sheetName, "R9", "T9")
	f.SetCellStyle(sheetName, "R9", "T9", boldBorderCenterStyle)
	f.SetCellStyle(sheetName, "U9", "U9", boldBorderCenterStyle)

	// Header Row 10: Sub headers (Largest, Middle, Smallest)
	row = 10
	stockSubHeaders := []string{"", "", "", "Largest", "Middle", "Smallest", "Largest", "Middle", "Smallest", "Largest", "Middle", "Smallest", "", ""}
	for i, header := range stockSubHeaders {
		if header != "" {
			cell := fmt.Sprintf("%s%d", stockCols[i], row)
			f.SetCellValue(sheetName, cell, header)
			f.SetCellStyle(sheetName, cell, cell, boldBorderCenterStyle)
		}
	}

	freeGoodsSubHeaders := []string{"Largest", "Middle", "Smallest", "Largest", "Middle", "Smallest", ""}
	freeGoodsCols := []string{"O", "P", "Q", "R", "S", "T", "U"}
	for i, header := range freeGoodsSubHeaders {
		if header != "" {
			cell := fmt.Sprintf("%s%d", freeGoodsCols[i], row)
			f.SetCellValue(sheetName, cell, header)
			f.SetCellStyle(sheetName, cell, cell, boldBorderCenterStyle)
		}
	}

	// Group data by ProID (to combine stock and promoted for same product)
	productMap := make(map[int64]map[int]model.GrDetList)
	for _, detail := range grDetails {
		if productMap[detail.ProID] == nil {
			productMap[detail.ProID] = make(map[int]model.GrDetList)
		}
		productMap[detail.ProID][detail.ItemType] = detail
	}

	// Convert to sorted list (by sequence)
	var sortedProducts []model.GrDetList
	for _, detail := range grDetails {
		if detail.ItemType == 1 {
			sortedProducts = append(sortedProducts, detail)
		}
	}

	var sumStockSubTotal, sumPpn float64

	// Write data rows (Row 11+)
	row = 11
	for _, stockDetail := range sortedProducts {
		proID := stockDetail.ProID
		stockData := stockDetail
		promoData, hasPromo := productMap[proID][2]

		// UOM
		uom := ""
		if stockData.UnitId3 != nil {
			uom += *stockData.UnitId3
		}
		if stockData.UnitId2 != nil {
			if uom != "" {
				uom += "/"
			}
			uom += *stockData.UnitId2
		}
		if stockData.UnitId1 != nil {
			if uom != "" {
				uom += "/"
			}
			uom += *stockData.UnitId1
		}

		// Stock Qty Shipment
		qtyShip3 := 0.0
		if stockData.QtyShip3 != nil {
			qtyShip3 = *stockData.QtyShip3
		}
		qtyShip2 := 0.0
		if stockData.QtyShip2 != nil {
			qtyShip2 = *stockData.QtyShip2
		}
		qtyShip1 := 0.0
		if stockData.QtyShip1 != nil {
			qtyShip1 = *stockData.QtyShip1
		}

		// Stock Qty Received
		qty3 := stockData.Qty3
		qty2 := stockData.Qty2
		qty1 := stockData.Qty1

		// Stock Price
		price3 := stockData.UnitPrice3
		price2 := stockData.UnitPrice2
		price1 := stockData.UnitPrice1

		stockSubTotal := (qty3 * price3) + (qty2 * price2) + (qty1 * price1)
		ppnRow := stockSubTotal * (stockData.Vat / 100)
		sumStockSubTotal += stockSubTotal
		sumPpn += ppnRow

		// Write Stock data (A-N)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), stockData.ProCode)
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), textBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), stockData.ProName)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), textBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), uom)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), textBorderStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), qtyShip3)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), qtyShip2)
		f.SetCellStyle(sheetName, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), qtyShip1)
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), numberBorderStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), qty3)
		f.SetCellStyle(sheetName, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), qty2)
		f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), qty1)
		f.SetCellStyle(sheetName, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), numberBorderStyle)

		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), price3)
		f.SetCellStyle(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), price2)
		f.SetCellStyle(sheetName, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), price1)
		f.SetCellStyle(sheetName, fmt.Sprintf("L%d", row), fmt.Sprintf("L%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), stockSubTotal)
		f.SetCellStyle(sheetName, fmt.Sprintf("M%d", row), fmt.Sprintf("M%d", row), numberBorderStyle)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), ppnRow)
		f.SetCellStyle(sheetName, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), numberBorderStyle)

		// Write Free Goods data (O-U)
		if hasPromo {
			promoQtyShip3 := 0.0
			if promoData.QtyShip3 != nil {
				promoQtyShip3 = *promoData.QtyShip3
			}
			promoQtyShip2 := 0.0
			if promoData.QtyShip2 != nil {
				promoQtyShip2 = *promoData.QtyShip2
			}
			promoQtyShip1 := 0.0
			if promoData.QtyShip1 != nil {
				promoQtyShip1 = *promoData.QtyShip1
			}

			promoQty3 := promoData.Qty3
			promoQty2 := promoData.Qty2
			promoQty1 := promoData.Qty1

			promoPrice3 := promoData.UnitPrice3
			promoPrice2 := promoData.UnitPrice2
			promoPrice1 := promoData.UnitPrice1

			promoSubTotal := (promoQty3 * promoPrice3) + (promoQty2 * promoPrice2) + (promoQty1 * promoPrice1)

			f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), promoQtyShip3)
			f.SetCellStyle(sheetName, fmt.Sprintf("O%d", row), fmt.Sprintf("O%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), promoQtyShip2)
			f.SetCellStyle(sheetName, fmt.Sprintf("P%d", row), fmt.Sprintf("P%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), promoQtyShip1)
			f.SetCellStyle(sheetName, fmt.Sprintf("Q%d", row), fmt.Sprintf("Q%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), promoQty3)
			f.SetCellStyle(sheetName, fmt.Sprintf("R%d", row), fmt.Sprintf("R%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("S%d", row), promoQty2)
			f.SetCellStyle(sheetName, fmt.Sprintf("S%d", row), fmt.Sprintf("S%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("T%d", row), promoQty1)
			f.SetCellStyle(sheetName, fmt.Sprintf("T%d", row), fmt.Sprintf("T%d", row), numberBorderStyle)
			f.SetCellValue(sheetName, fmt.Sprintf("U%d", row), promoSubTotal)
			f.SetCellStyle(sheetName, fmt.Sprintf("U%d", row), fmt.Sprintf("U%d", row), numberBorderStyle)
		} else {
			// Empty Free Goods columns - fill with zeros
			for _, col := range []string{"O", "P", "Q", "R", "S", "T", "U"} {
				f.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), 0)
				f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), numberBorderStyle)
			}
		}

		row++
	}

	// Summary section (Subtotal, PPn, Delivery fee, Total) - bottom right like image (S/T columns)
	summaryRow := row + 2
	deliveryFee := 0.0
	if gr.DeliveryFee != nil {
		deliveryFee = *gr.DeliveryFee
	}
	total := sumStockSubTotal + sumPpn + deliveryFee

	f.SetCellValue(sheetName, fmt.Sprintf("S%d", summaryRow), "Subtotal")
	f.SetCellStyle(sheetName, fmt.Sprintf("S%d", summaryRow), fmt.Sprintf("S%d", summaryRow), boldStyle)
	f.SetCellValue(sheetName, fmt.Sprintf("T%d", summaryRow), sumStockSubTotal)
	f.SetCellStyle(sheetName, fmt.Sprintf("T%d", summaryRow), fmt.Sprintf("T%d", summaryRow), numberBorderStyle)

	f.SetCellValue(sheetName, fmt.Sprintf("S%d", summaryRow+1), "PPn")
	f.SetCellStyle(sheetName, fmt.Sprintf("S%d", summaryRow+1), fmt.Sprintf("S%d", summaryRow+1), boldStyle)
	f.SetCellValue(sheetName, fmt.Sprintf("T%d", summaryRow+1), sumPpn)
	f.SetCellStyle(sheetName, fmt.Sprintf("T%d", summaryRow+1), fmt.Sprintf("T%d", summaryRow+1), numberBorderStyle)

	f.SetCellValue(sheetName, fmt.Sprintf("S%d", summaryRow+2), "Delivery fee")
	f.SetCellStyle(sheetName, fmt.Sprintf("S%d", summaryRow+2), fmt.Sprintf("S%d", summaryRow+2), boldStyle)
	f.SetCellValue(sheetName, fmt.Sprintf("T%d", summaryRow+2), deliveryFee)
	f.SetCellStyle(sheetName, fmt.Sprintf("T%d", summaryRow+2), fmt.Sprintf("T%d", summaryRow+2), numberBorderStyle)

	f.SetCellValue(sheetName, fmt.Sprintf("S%d", summaryRow+3), "Total")
	f.SetCellStyle(sheetName, fmt.Sprintf("S%d", summaryRow+3), fmt.Sprintf("S%d", summaryRow+3), boldStyle)
	f.SetCellValue(sheetName, fmt.Sprintf("T%d", summaryRow+3), total)
	f.SetCellStyle(sheetName, fmt.Sprintf("T%d", summaryRow+3), fmt.Sprintf("T%d", summaryRow+3), numberBorderStyle)

	return f, nil
}

var (
	sequenceMutex sync.Mutex
	sequenceFile  = ".download_good_receipt_sequence.json"
)

type SequenceStorage struct {
	Sequences map[string]int `json:"sequences"` // map[DDMMYY]sequenceNumber
}

func getNextSequenceNumber(dateStr string) (int, error) {
	sequenceMutex.Lock()
	defer sequenceMutex.Unlock()

	storage, err := readSequenceStorage()
	if err != nil {
		return 0, err
	}

	if storage.Sequences == nil {
		storage.Sequences = make(map[string]int)
	}

	currentSeq := storage.Sequences[dateStr]

	nextSeq := currentSeq + 1
	storage.Sequences[dateStr] = nextSeq

	err = writeSequenceStorage(storage)
	if err != nil {
		return 0, err
	}

	return nextSeq, nil
}

func readSequenceStorage() (*SequenceStorage, error) {
	storage := &SequenceStorage{
		Sequences: make(map[string]int),
	}

	if _, err := os.Stat(sequenceFile); os.IsNotExist(err) {
		return storage, nil
	}

	data, err := os.ReadFile(sequenceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read sequence file: %w", err)
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, storage)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sequence file: %w", err)
		}
	}

	return storage, nil
}

func writeSequenceStorage(storage *SequenceStorage) error {
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sequence storage: %w", err)
	}

	err = os.WriteFile(sequenceFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write sequence file: %w", err)
	}

	return nil
}
