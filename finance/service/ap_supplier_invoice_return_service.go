package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/conversion"
	"finance/pkg/str"
	"finance/pkg/structs"
	"finance/repository"
	"fmt"
	"strings"
)

type ApSupplierInvoiceReturnService interface {
	Store(request *entity.ApSupplierInvoiceReturnCreate) (err error)
	Detail(accountPayableID uint, custID, parentCustId string) (response entity.ApSupplierInvoiceReturnRespone, err error)
	List(dataFilter entity.ApSupplierInoviceReturnQueryFilter) (data []entity.ApSupplierInvoiceReturnResponeList, total int64, lastPage int, err error)
	Delete(custId string, accountPayableID uint, userId int64) (err error)
	Update(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) (err error)
}

type ApSupplierInvoiceReturnServiceImpl struct {
	Repository  repository.ApSupplierInvoiceReturnRepository
	Transaction repository.Dbtransaction
}

func NewApSupplierInvoiceReturnService(repository repository.ApSupplierInvoiceReturnRepository, transaction repository.Dbtransaction) *ApSupplierInvoiceReturnServiceImpl {
	return &ApSupplierInvoiceReturnServiceImpl{
		Repository:  repository,
		Transaction: transaction,
	}
}

func (service *ApSupplierInvoiceReturnServiceImpl) Store(request *entity.ApSupplierInvoiceReturnCreate) (err error) {
	// parse time format YYYY-mm-dd to Rfc3339
	if request.AccountPayableDate != nil {
		if *request.AccountPayableDate != "" {
			apDate, err := str.DateStrToRfc3339String(*request.AccountPayableDate)
			if err != nil {
				return err
			}
			request.AccountPayableDate = &apDate
		} else {
			request.AccountPayableDate = nil
		}
	}

	if request.InvoiceDate != nil {
		if *request.InvoiceDate != "" {
			invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
			if err != nil {
				return err
			}
			request.InvoiceDate = &invoiceDate
		} else {
			request.InvoiceDate = nil
		}
	}

	if request.TaxInvoiceDate != nil {
		if *request.TaxInvoiceDate != "" {
			taxDate, err := str.DateStrToRfc3339String(*request.TaxInvoiceDate)
			if err != nil {
				return err
			}
			request.TaxInvoiceDate = &taxDate
		} else {
			request.TaxInvoiceDate = nil
		}
	}

	if request.DueDate != nil {
		if *request.DueDate != "" {
			dueDate, err := str.DateStrToRfc3339String(*request.DueDate)
			if err != nil {
				return err
			}
			request.DueDate = &dueDate
		} else {
			request.DueDate = nil
		}
	}

	apWhereDocument, err := service.Repository.FindByDocumentNoAndType(request.DocumentNo, request.ApType, request.CustId, request.ParentCustID)
	if err == nil {
		if apWhereDocument.ID != 0 {
			return errors.New(fmt.Sprintf("Document No %v already created on AP with id %v", request.DocumentNo, apWhereDocument.ID))
		}
	}

	if strings.TrimSpace(request.InvoiceNo) != "" {
		apWhereInvoice, err := service.Repository.FindByInvoiceNo(request.InvoiceNo, request.CustId)
		if err == nil && apWhereInvoice.ID != 0 {
			return errors.New("Invoice number must be unique.")
		}
	}

	if request.ApType == entity.AP_TYPE_INVOICE {
		err := service.StoreApInvoice(request)
		if err != nil {
			return err
		}
	} else {
		err := service.StoreApReturn(request)
		if err != nil {
			return err
		}
	}

	return nil
}
func (service *ApSupplierInvoiceReturnServiceImpl) StoreApReturn(request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()

	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		_, err := service.Repository.GetReturnSupplierByDocumentNo(request.DocumentNo, request.CustId)
		if err != nil {
			return err
		}

		returnProducts, err := service.Repository.FindDetBySupplierReturnNo(request.DocumentNo, request.CustId)
		if err != nil {
			return err
		}

		for _, returnProduct := range returnProducts {
			product := &entity.ProductListCreate{
				ProId:       returnProduct.ProID,
				Qty:         int(returnProduct.Qty),
				UnitPrice1:  returnProduct.UnitPrice1,
				UnitPrice2:  returnProduct.UnitPrice2,
				UnitPrice3:  returnProduct.UnitPrice3,
				ConvUnit2:   returnProduct.ConvUnit2,
				ConvUnit3:   returnProduct.ConvUnit3,
				Vat:         returnProduct.Vat,
				VatLg:       returnProduct.VatLg,
				Disc:        returnProduct.Discount,
				InvoiceDisc: request.DiscountPercent,
				Type:        entity.AP_DET_TYPE_NORMAL,
			}

			product.Calculate() // calculate per product
			request.ProductLists = append(request.ProductLists, product)
		}
		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturn
		err = structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Store(txCtx, &Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {

			detModel := model.AccountPayableProduct{
				AccountPayableID:  Apmodel.ID,
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
				Type:              ProductList.Type,
			}
			err = service.Repository.StoreDetail(txCtx, &detModel)
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
func (service *ApSupplierInvoiceReturnServiceImpl) StoreApInvoice(request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()

	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		// check GR
		// Ambil 3 karakter pertama
		prefix := request.DocumentNo[:3]

		fmt.Println(prefix) // Output: GRB

		if prefix == "GRB" {
			_, err := service.Repository.GetGrbByInvoiceNo(request.DocumentNo, request.CustId, request.ParentCustID, request.CustIdParam)
			if err != nil {
				return err
			}

			GrProducts, err := service.Repository.FindGrbdetail(request.DocumentNo, request.CustId, request.CustIdParam)
			if err != nil {
				return err
			}

			for _, GrProduct := range GrProducts {
				if GrProduct.ItemType == 1 { // normal
					normalProduct := &entity.ProductListCreate{
						ProId:       GrProduct.ProID,
						Qty:         int(GrProduct.Qty),
						UnitPrice1:  GrProduct.UnitPrice1,
						UnitPrice2:  GrProduct.UnitPrice2,
						UnitPrice3:  GrProduct.UnitPrice3,
						ConvUnit2:   GrProduct.ConvUnit2,
						ConvUnit3:   GrProduct.ConvUnit3,
						Vat:         GrProduct.Vat,
						VatLg:       GrProduct.VatLgPurch,
						Disc:        GrProduct.Discount,
						InvoiceDisc: request.DiscountPercent,
						Type:        entity.AP_DET_TYPE_NORMAL,
					}
					normalProduct.Calculate() // calculate per product

					request.ProductLists = append(request.ProductLists, normalProduct)
				} else {
					promoProduct := &entity.ProductListCreate{
						ProId:     GrProduct.ProID,
						Qty:       int(GrProduct.Qty),
						ConvUnit2: GrProduct.ConvUnit2,
						ConvUnit3: GrProduct.ConvUnit3,
						Type:      entity.AP_DET_TYPE_PROMO,
					}

					request.ProductLists = append(request.ProductLists, promoProduct)
				}
			}
		} else {
			_, err := service.Repository.GetGrByInvoiceNo(request.DocumentNo, request.CustId, request.ParentCustID)
			if err != nil {
				return err
			}

			GrProducts, err := service.Repository.FindGrdetail(request.DocumentNo, request.CustId)
			if err != nil {
				return err
			}

			for _, GrProduct := range GrProducts {
				if GrProduct.ItemType == 1 { // normal
					normalProduct := &entity.ProductListCreate{
						ProId:       GrProduct.ProID,
						Qty:         int(GrProduct.Qty),
						UnitPrice1:  GrProduct.UnitPrice1,
						UnitPrice2:  GrProduct.UnitPrice2,
						UnitPrice3:  GrProduct.UnitPrice3,
						ConvUnit2:   GrProduct.ConvUnit2,
						ConvUnit3:   GrProduct.ConvUnit3,
						Vat:         GrProduct.Vat,
						VatLg:       GrProduct.VatLgPurch,
						Disc:        GrProduct.Discount,
						InvoiceDisc: request.DiscountPercent,
						Type:        entity.AP_DET_TYPE_NORMAL,
					}
					normalProduct.Calculate() // calculate per product

					request.ProductLists = append(request.ProductLists, normalProduct)
				} else {
					promoProduct := &entity.ProductListCreate{
						ProId:     GrProduct.ProID,
						Qty:       int(GrProduct.Qty),
						ConvUnit2: GrProduct.ConvUnit2,
						ConvUnit3: GrProduct.ConvUnit3,
						Type:      entity.AP_DET_TYPE_PROMO,
					}

					request.ProductLists = append(request.ProductLists, promoProduct)
				}
			}
		}

		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturn
		err := structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Store(txCtx, &Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {

			detModel := model.AccountPayableProduct{
				AccountPayableID:  Apmodel.ID,
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
				Type:              ProductList.Type,
			}
			err = service.Repository.StoreDetail(txCtx, &detModel)
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

func (service *ApSupplierInvoiceReturnServiceImpl) Detail(accountPayableID uint, custID, parentCustId string) (response entity.ApSupplierInvoiceReturnRespone, err error) {
	ap, err := service.Repository.FindByID(accountPayableID, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ap, &response)
	if err != nil {
		return response, err
	}

	ProductLists, err := service.Repository.FindProductList(accountPayableID, custID)
	if err != nil {
		return response, err
	}
	var ProductListDatas []entity.ProductListRespone
	var ProductPromoDatas []entity.ProductPromoRespone

	for _, detail := range ProductLists {
		if detail.ItemType == entity.AP_DET_TYPE_NORMAL {

			var detailData entity.ProductListRespone
			err = structs.Automapper(detail, &detailData)

			qty := &conversion.Qty{
				Qty:       int(detailData.Qty),
				ConvUnit2: int(detailData.ConvUnit2),
				ConvUnit3: int(detailData.ConvUnit3),
			}

			qtyConversion := qty.ConvToQtyConversion()
			// fmt.Println("====>", qtyConversion)

			detailData.Qty1 = int(qtyConversion.Qty1)
			detailData.Qty2 = int(qtyConversion.Qty2)
			detailData.Qty3 = int(qtyConversion.Qty3)

			qtyRemaining := &conversion.Qty{
				Qty:       int(detail.QtyRemaining),
				ConvUnit2: int(detail.ConvUnit2),
				ConvUnit3: int(detail.ConvUnit3),
			}
			qtyRemainingConversion := qtyRemaining.ConvToQtyConversion()

			detailData.QtyRemaining1 = qtyRemainingConversion.Qty1
			detailData.QtyRemaining2 = qtyRemainingConversion.Qty2
			detailData.QtyRemaining3 = qtyRemainingConversion.Qty3

			if err != nil {
				return response, err
			}
			if ap.ApType == entity.AP_TYPE_INVOICE {

				prefix := ap.DocumentNo[:3]

				if prefix == "GRB" {
					whQty, err := service.Repository.GetWarehouseStockFromGrb(custID, ap.DocumentNo, detail.ProId)
					if err != nil {
						return response, err
					}
					qtyWh := &conversion.Qty{
						Qty:       int(whQty.Qty),
						ConvUnit2: int(detail.ConvUnit2),
						ConvUnit3: int(detail.ConvUnit3),
					}
					qtyWhConversion := qtyWh.ConvToQtyConversion()
					detailData.WhQty1 = qtyWhConversion.Qty1
					detailData.WhQty2 = qtyWhConversion.Qty2
					detailData.WhQty3 = qtyWhConversion.Qty3

				} else {
					whQty, err := service.Repository.GetWarehouseStockFromGr(custID, ap.DocumentNo, detail.ProId)
					if err != nil {
						return response, err
					}
					qtyWh := &conversion.Qty{
						Qty:       int(whQty.Qty),
						ConvUnit2: int(detail.ConvUnit2),
						ConvUnit3: int(detail.ConvUnit3),
					}
					qtyWhConversion := qtyWh.ConvToQtyConversion()
					detailData.WhQty1 = qtyWhConversion.Qty1
					detailData.WhQty2 = qtyWhConversion.Qty2
					detailData.WhQty3 = qtyWhConversion.Qty3
				}
			} else {
				whQty, err := service.Repository.GetWarehouseStockFromReturn(custID, ap.DocumentNo, detail.ProId)
				if err != nil {
					return response, err
				}
				qtyWh := &conversion.Qty{
					Qty:       int(whQty.Qty),
					ConvUnit2: int(detail.ConvUnit2),
					ConvUnit3: int(detail.ConvUnit3),
				}
				qtyWhConversion := qtyWh.ConvToQtyConversion()
				detailData.WhQty1 = qtyWhConversion.Qty1
				detailData.WhQty2 = qtyWhConversion.Qty2
				detailData.WhQty3 = qtyWhConversion.Qty3
			}

			ProductListDatas = append(ProductListDatas, detailData)
		} else {
			var detailData entity.ProductPromoRespone
			err = structs.Automapper(detail, &detailData)

			qty := &conversion.Qty{
				Qty:       int(detailData.Qty),
				ConvUnit2: int(detailData.ConvUnit2),
				ConvUnit3: int(detailData.ConvUnit3),
			}

			qtyConversion := qty.ConvToQtyConversion()

			detailData.Qty1 = int(qtyConversion.Qty1)
			detailData.Qty2 = int(qtyConversion.Qty2)
			detailData.Qty3 = int(qtyConversion.Qty3)

			if err != nil {
				return response, err
			}

			ProductPromoDatas = append(ProductPromoDatas, detailData)
		}
	}

	if ap.AccountPayableDate != nil {
		apDate := ap.AccountPayableDate.Format("2006-01-02")
		response.AccountPayableDate = apDate
	}

	if ap.InvoiceDate != nil {
		invDate := ap.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = invDate
	}

	if ap.TaxInvoiceDate != nil {
		taxInvDate := ap.TaxInvoiceDate.Format("2006-01-02")
		response.TaxInvoiceDate = taxInvDate
	}

	if ap.DueDate != nil {
		dueDate := ap.DueDate.Format("2006-01-02")
		response.DueDate = dueDate
	}

	if response.ApType == "I" {
		response.ApType = "Invoice"
	} else {
		response.ApType = "Return"
	}

	response.ProductList = ProductListDatas
	response.ProductPromo = ProductPromoDatas
	return response, nil
}

func (service *ApSupplierInvoiceReturnServiceImpl) List(dataFilter entity.ApSupplierInoviceReturnQueryFilter) (data []entity.ApSupplierInvoiceReturnResponeList, total int64, lastPage int, err error) {
	aps, total, lastPage, err := service.Repository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range aps {
		var vResp entity.ApSupplierInvoiceReturnResponeList
		structs.Automapper(row, &vResp)
		if row.AccountPayableDate != nil {
			accountPayableDate := row.AccountPayableDate.Format("2006-01-02")
			vResp.AccountPayableDate = accountPayableDate
		}

		if row.InvoiceDate != nil {
			invDate := row.InvoiceDate.Format("2006-01-02")
			vResp.AccountPayableDate = invDate
		}

		if row.TaxInvoiceDate != nil {
			taxInvDate := row.TaxInvoiceDate.Format("2006-01-02")
			vResp.TaxInvoiceDate = taxInvDate
		}

		if row.DueDate != nil {
			dueDate := row.DueDate.Format("2006-01-02")
			vResp.DueDate = dueDate
		}

		if row.ApType == "I" {
			vResp.ApType = "Invoice"
		} else {
			vResp.ApType = "Return"
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *ApSupplierInvoiceReturnServiceImpl) Delete(custId string, accountPayableID uint, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.Repository.Delete(txCtx, custId, accountPayableID, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *ApSupplierInvoiceReturnServiceImpl) Update(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) (err error) {

	// parse time format YYYY-mm-dd to Rfc3339
	if request.AccountPayableDate != nil {
		if *request.AccountPayableDate != "" {
			apDate, err := str.DateStrToRfc3339String(*request.AccountPayableDate)
			if err != nil {
				return err
			}
			request.AccountPayableDate = &apDate
		} else {
			request.AccountPayableDate = nil
		}
	}

	if request.InvoiceDate != nil {
		if *request.InvoiceDate != "" {
			invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
			if err != nil {
				return err
			}
			request.InvoiceDate = &invoiceDate
		} else {
			request.InvoiceDate = nil
		}
	}

	if request.TaxInvoiceDate != nil {
		if *request.TaxInvoiceDate != "" {
			taxDate, err := str.DateStrToRfc3339String(*request.TaxInvoiceDate)
			if err != nil {
				return err
			}
			request.TaxInvoiceDate = &taxDate
		} else {
			request.TaxInvoiceDate = nil
		}
	}

	if request.DueDate != nil {
		if *request.DueDate != "" {
			dueDate, err := str.DateStrToRfc3339String(*request.DueDate)
			if err != nil {
				return err
			}
			request.DueDate = &dueDate
		} else {
			request.DueDate = nil
		}
	}

	if request.ReturnDate != nil {
		if *request.ReturnDate != "" {
			returnDate, err := str.DateStrToRfc3339String(*request.ReturnDate)
			if err != nil {
				return err
			}
			request.ReturnDate = &returnDate
		} else {
			request.ReturnDate = nil
		}
	}

	apWhereDocument, err := service.Repository.FindByDocumentNoAndType(request.DocumentNo, request.ApType, request.CustId, request.ParentCustID)
	if err == nil {
		if apWhereDocument.ID != 0 && apWhereDocument.ID != accountPayableID {
			return errors.New(fmt.Sprintf("Document No %v already created on AP with id %v", request.DocumentNo, apWhereDocument.ID))
		}
	}

	if strings.TrimSpace(request.InvoiceNo) != "" {
		apWhereInvoice, err := service.Repository.FindByInvoiceNo(request.InvoiceNo, request.CustId)
		if err == nil && apWhereInvoice.ID != 0 && apWhereInvoice.ID != accountPayableID {
			return errors.New("Invoice number must be unique.")
		}
	}

	if request.ApType == entity.AP_TYPE_INVOICE {
		err := service.EditApInvoice(accountPayableID, request)
		if err != nil {
			return err
		}
	} else {
		err := service.EditApReturn(accountPayableID, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ApSupplierInvoiceReturnServiceImpl) EditApReturn(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	ap, err := service.Repository.FindByID(accountPayableID, request.CustId, request.ParentCustID)
	if err != nil {
		return err
	}

	supplierReturn, err := service.Repository.GetReturnSupplierByDocumentNo(request.DocumentNo, request.CustId)
	if err == nil {
		if supplierReturn.GrNO != "" {
			return errors.New(fmt.Sprintf("can't update account payable. invoice already processed on return supplier by document no %v", supplierReturn.SupplierReturnNo))
		}
	}

	apPay, err := service.Repository.GetArPayByInvoiceNo(ap.InvoiceNo, request.CustId)
	if err == nil {
		if apPay.ApPayNo != "" {
			return errors.New(fmt.Sprintf("can't update account payable. invoice already processed on account payable pay by document no %v", apPay.ApPayNo))
		}
	}

	if ap.DocumentNo != request.DocumentNo {
		err := service.EditApReturnNotEqualDocumentNo(accountPayableID, request)
		if err != nil {
			return err
		}
	} else {
		err := service.EditApReturnEqualDocument(accountPayableID, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *ApSupplierInvoiceReturnServiceImpl) EditApReturnNotEqualDocumentNo(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()
	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.DeleteDetail(txCtx, request.CustId, accountPayableID, request.UpdatedBy)
		if err != nil {
			return err
		}

		returnProducts, err := service.Repository.FindDetBySupplierReturnNo(request.DocumentNo, request.CustId)
		if err != nil {
			return err
		}

		for _, returnProduct := range returnProducts {
			product := &entity.ProductListCreate{
				ProId:       returnProduct.ProID,
				Qty:         int(returnProduct.Qty),
				UnitPrice1:  returnProduct.UnitPrice1,
				UnitPrice2:  returnProduct.UnitPrice2,
				UnitPrice3:  returnProduct.UnitPrice3,
				ConvUnit2:   returnProduct.ConvUnit2,
				ConvUnit3:   returnProduct.ConvUnit3,
				Vat:         returnProduct.Vat,
				VatLg:       returnProduct.VatLg,
				Disc:        returnProduct.Discount,
				InvoiceDisc: request.DiscountPercent,
				Type:        entity.AP_DET_TYPE_NORMAL,
			}

			product.Calculate() // calculate per product
			request.ProductLists = append(request.ProductLists, product)
		}
		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturnupdate
		err = structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Update(txCtx, accountPayableID, Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {
			detModel := model.AccountPayableProduct{
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
			}
			err = service.Repository.UpdateProductNormal(txCtx, *ProductList.ID, &detModel)
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

func (service *ApSupplierInvoiceReturnServiceImpl) EditApReturnEqualDocument(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()
	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		ProductNormalLists, err := service.Repository.FindProductList(accountPayableID, request.CustId)
		if err != nil {
			return err
		}

		for _, ProductNormal := range ProductNormalLists {
			normalProduct := &entity.ProductListCreate{
				ID:          &ProductNormal.AccountPayableDetailID,
				ProId:       ProductNormal.ProId,
				Qty:         int(ProductNormal.Qty),
				UnitPrice1:  ProductNormal.UnitPrice1,
				UnitPrice2:  ProductNormal.UnitPrice2,
				UnitPrice3:  ProductNormal.UnitPrice3,
				ConvUnit2:   ProductNormal.ConvUnit2,
				ConvUnit3:   ProductNormal.ConvUnit3,
				Vat:         ProductNormal.Vat,
				VatLg:       ProductNormal.VatLg,
				Disc:        ProductNormal.Disc,
				InvoiceDisc: request.DiscountPercent,
			}

			normalProduct.Calculate() // calculate per product

			request.ProductLists = append(request.ProductLists, normalProduct)
		}

		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturnupdate
		err = structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Update(txCtx, accountPayableID, Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {
			detModel := model.AccountPayableProduct{
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
			}
			err = service.Repository.UpdateProductNormal(txCtx, *ProductList.ID, &detModel)
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

func (service *ApSupplierInvoiceReturnServiceImpl) EditApInvoice(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	ap, err := service.Repository.FindByID(accountPayableID, request.CustId, request.ParentCustID)
	if err != nil {
		return err
	}

	supplierReturn, err := service.Repository.GetReturnSupplierByDocumentNo(ap.DocumentNo, request.CustId)
	if err == nil {
		if supplierReturn.GrNO != "" {
			return errors.New(fmt.Sprintf("can't update account payable. invoice already processed on return supplier by document no %v", supplierReturn.SupplierReturnNo))
		}
	}

	apPay, err := service.Repository.GetArPayByInvoiceNo(ap.InvoiceNo, request.CustId)
	if err == nil {
		if apPay.ApPayNo != "" {
			return errors.New(fmt.Sprintf("can't update account payable. invoice already processed on account payable pay by document no %v", apPay.ApPayNo))
		}
	}

	if ap.DocumentNo != request.DocumentNo {
		err := service.EditApInvoiceNotEqualInvoice(accountPayableID, request)
		if err != nil {
			return err
		}
	} else {
		err := service.EditApInvoiceEqualInvoice(accountPayableID, request)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ApSupplierInvoiceReturnServiceImpl) EditApInvoiceNotEqualInvoice(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()
	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.Repository.DeleteDetail(txCtx, request.CustId, accountPayableID, request.UpdatedBy)
		if err != nil {
			return err
		}

		// check GR
		_, err = service.Repository.GetGrByInvoiceNo(request.DocumentNo, request.CustId, request.ParentCustID)
		if err != nil {
			return err
		}

		GrProducts, err := service.Repository.FindGrdetail(request.DocumentNo, request.CustId)
		if err != nil {
			return err
		}

		for _, GrProduct := range GrProducts {
			if GrProduct.ItemType == 1 { // normal
				normalProduct := &entity.ProductListCreate{
					ProId:       GrProduct.ProID,
					Qty:         int(GrProduct.Qty),
					UnitPrice1:  GrProduct.UnitPrice1,
					UnitPrice2:  GrProduct.UnitPrice2,
					UnitPrice3:  GrProduct.UnitPrice3,
					ConvUnit2:   GrProduct.ConvUnit2,
					ConvUnit3:   GrProduct.ConvUnit3,
					Vat:         GrProduct.Vat,
					VatLg:       GrProduct.VatLgPurch,
					Disc:        GrProduct.Discount,
					InvoiceDisc: request.DiscountPercent,
					Type:        entity.AP_DET_TYPE_NORMAL,
				}
				normalProduct.Calculate() // calculate per product

				request.ProductLists = append(request.ProductLists, normalProduct)
			} else {
				promoProduct := &entity.ProductListCreate{
					ProId:     GrProduct.ProID,
					Qty:       int(GrProduct.Qty),
					ConvUnit2: GrProduct.ConvUnit2,
					ConvUnit3: GrProduct.ConvUnit3,
					Type:      entity.AP_DET_TYPE_PROMO,
				}

				request.ProductLists = append(request.ProductLists, promoProduct)
			}
		}

		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturnupdate
		err = structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Update(txCtx, accountPayableID, Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {

			detModel := model.AccountPayableProduct{
				AccountPayableID:  &accountPayableID,
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
				Type:              ProductList.Type,
			}
			err = service.Repository.StoreDetail(txCtx, &detModel)
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

func (service *ApSupplierInvoiceReturnServiceImpl) EditApInvoiceEqualInvoice(accountPayableID uint, request *entity.ApSupplierInvoiceReturnCreate) error {
	c := context.Background()
	err := service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		ProductNormalLists, err := service.Repository.FindProductList(accountPayableID, request.CustId)
		if err != nil {
			return err
		}

		for _, ProductNormal := range ProductNormalLists {
			normalProduct := &entity.ProductListCreate{
				ID:          &ProductNormal.AccountPayableDetailID,
				ProId:       ProductNormal.ProId,
				Qty:         int(ProductNormal.Qty),
				UnitPrice1:  ProductNormal.UnitPrice1,
				UnitPrice2:  ProductNormal.UnitPrice2,
				UnitPrice3:  ProductNormal.UnitPrice3,
				ConvUnit2:   ProductNormal.ConvUnit2,
				ConvUnit3:   ProductNormal.ConvUnit3,
				Vat:         ProductNormal.Vat,
				VatLg:       ProductNormal.VatLg,
				Disc:        ProductNormal.Disc,
				InvoiceDisc: request.DiscountPercent,
			}

			normalProduct.Calculate() // calculate per product

			request.ProductLists = append(request.ProductLists, normalProduct)
		}

		request.Calculate() // calculate all

		var Apmodel model.ApSupplierInvoiceReturnupdate
		err = structs.Automapper(request, &Apmodel)
		if err != nil {
			return err
		}
		Apmodel.DiscountPercent = &request.DiscountPercent
		Apmodel.DiscountRp = &request.DiscountValue
		Apmodel.Amount = request.TotalSkuprice
		Apmodel.SubTotal = request.Subtotal
		Apmodel.VatValue = request.TotalVatValue
		Apmodel.VatLgValue = request.TotalVatLgValue
		Apmodel.Total = request.Total

		err = service.Repository.Update(txCtx, accountPayableID, Apmodel)
		if err != nil {
			return err
		}

		for _, ProductList := range request.ProductLists {
			detModel := model.AccountPayableProduct{
				CustId:            request.CustId,
				ProId:             ProductList.ProId,
				UnitPrice1:        ProductList.UnitPrice1,
				UnitPrice2:        ProductList.UnitPrice2,
				UnitPrice3:        ProductList.UnitPrice3,
				ConvUnit2:         ProductList.ConvUnit2,
				ConvUnit3:         ProductList.ConvUnit3,
				SubTotal:          ProductList.Gross,
				Disc:              ProductList.Disc,
				DiscValue:         ProductList.DiscValue,
				SubTotalBeforePpn: ProductList.NetAmount,
				Vat:               ProductList.Vat,
				VatValue:          ProductList.VatValue,
				VatLg:             ProductList.VatLg,
				VatLgValue:        ProductList.VatLgValue,
				Total:             ProductList.NetAmountAfterInvoiceDiscount,
				Qty:               ProductList.Qty,
			}
			err = service.Repository.UpdateProductNormal(txCtx, *ProductList.ID, &detModel)
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
