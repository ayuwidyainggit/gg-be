package service

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/repository"
	"time"
)

type CoreTaxVatExtractService interface {
	List(dataFilter entity.CoreTaxVatExtractQueryFilter) (datas []entity.CoreTaxVatExtractListResponse, total int64, lastPage int, err error)
	Extract(request entity.CoreTaxExtractReq) (resp entity.CoretaxVatExtractResp, err error)
	ExtractDownloadResult(vatExtractID int64, custID string, parentCustId string) (interface{}, error)
}

type CoreTaxVatExtractServiceImpl struct {
	CoreTaxVatExtractRepository repository.CoreTaxVatExtractRepository
	Transaction                 repository.Dbtransaction
}

func NewCoreTaxVatExtractService(repository repository.CoreTaxVatExtractRepository, transaction repository.Dbtransaction) *CoreTaxVatExtractServiceImpl {
	return &CoreTaxVatExtractServiceImpl{
		CoreTaxVatExtractRepository: repository,
		Transaction:                 transaction,
	}
}

func (service *CoreTaxVatExtractServiceImpl) List(dataFilter entity.CoreTaxVatExtractQueryFilter) (datas []entity.CoreTaxVatExtractListResponse, total int64, lastPage int, err error) {
	if dataFilter.InvoiceType == entity.CORETAX_TYPE_INVOICE {
		return service.ListInvoice(dataFilter)
	}

	return service.ListReturn(dataFilter)
}

func (service *CoreTaxVatExtractServiceImpl) ListInvoice(dataFilter entity.CoreTaxVatExtractQueryFilter) (datas []entity.CoreTaxVatExtractListResponse, total int64, lastPage int, err error) {
	orders, total, lastPage, err := service.CoreTaxVatExtractRepository.FindInvoiceListByCustId(dataFilter)
	if err != nil {
		return datas, total, lastPage, err
	}

	for _, order := range orders {
		var invoiceDate, taxExtractDate, invoiceNo string
		if order.InvoiceDate != nil {
			invoiceDate = order.InvoiceDate.Format("2006-01-02")
		}

		if order.TaxExtractDate != nil {
			taxExtractDate = order.TaxExtractDate.Format("2006-01-02")
		}

		if order.InvoiceNo != nil {
			invoiceNo = *order.InvoiceNo
		}
		datas = append(datas, entity.CoreTaxVatExtractListResponse{
			TransactionID:     order.RoNo,
			InvoiceNo:         invoiceNo,
			InvoiceDate:       invoiceDate,
			TaxNo:             order.TaxIdentifierNo,
			TaxType:           order.TaxIdentifierType,
			SalesId:           order.SalesId,
			SalesCode:         order.SalesCode,
			SalesName:         order.SalesName,
			OutletID:          order.OutletID,
			OutletCode:        order.OutletCode,
			OutletName:        order.OutletName,
			OutletAddress1:    order.OutletAddress1,
			OutletAddress2:    order.OutletAddress2,
			OutletTaxAddress1: order.OutletTaxAddress1,
			OutletTaxAddress2: order.OutletTaxAddress2,
			TaxExtractDate:    taxExtractDate,
			PPN:               order.Vat,
			PPNValue:          order.VatValue,
			PPNFinalValue:     order.VatValueFinal,
			DPP:               order.DPP,
			DPPFinal:          order.TotalFinal,
			PPNBM:             0,
			NITKU:             order.NITKU,
		})
	}
	return
}

func (service *CoreTaxVatExtractServiceImpl) ListReturn(dataFilter entity.CoreTaxVatExtractQueryFilter) (datas []entity.CoreTaxVatExtractListResponse, total int64, lastPage int, err error) {
	return
}

func (service *CoreTaxVatExtractServiceImpl) Extract(request entity.CoreTaxExtractReq) (resp entity.CoretaxVatExtractResp, err error) {
	c := context.Background()

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		now := time.Now()
		coretaxVatExtractModel := model.CoretaxVatExtract{
			InvoiceType:  request.InvoiceType,
			ExtractTotal: len(request.TransactionID),
			CreatedBy:    request.CreatedBy,
			CreatedAt:    &now,
		}

		err := service.CoreTaxVatExtractRepository.Store(txCtx, &coretaxVatExtractModel)
		if err != nil {
			return err
		}

		var coretaxVatExtractDetailModels []model.CoretaxVatExtractDetail

		for _, transactionID := range request.TransactionID {
			coretaxVatExtractDetailModels = append(coretaxVatExtractDetailModels, model.CoretaxVatExtractDetail{
				CustID:       request.CustID,
				VatExtractID: *coretaxVatExtractModel.CoretaxVatExtractID,
				ReferenceID:  transactionID,
			})
		}

		err = service.CoreTaxVatExtractRepository.StoreDetail(txCtx, coretaxVatExtractDetailModels)
		if err != nil {
			return err
		}

		resp.ID = *coretaxVatExtractModel.CoretaxVatExtractID
		return nil
	})
	return
}

func (service *CoreTaxVatExtractServiceImpl) ExtractDownloadResult(vatExtractID int64, custID string, parentCustId string) (interface{}, error) {
	coretaxVatExtract, err := service.CoreTaxVatExtractRepository.FindCoretaxVatById(vatExtractID)
	if err != nil {
		return nil, err
	}

	if coretaxVatExtract.InvoiceType == entity.CORETAX_TYPE_INVOICE {
		orders, err := service.ExtractDownloadResultInvoice(vatExtractID, custID, parentCustId)
		if err != nil {
			return nil, err
		}
		return orders, err
	}
	return nil, nil
}

func (service *CoreTaxVatExtractServiceImpl) ExtractDownloadResultInvoice(vatExtractID int64, custID string, parentCustId string) (data entity.CoretaxVatExtractDownload, err error) {
	orders, err := service.CoreTaxVatExtractRepository.FindInvoiceExtractByID(vatExtractID, custID)
	if err != nil {
		return data, err
	}

	for index, order := range orders {
		var invoiceDate, invoiceNo string
		var invoiceDetailCoretaxExtract []entity.CoretaxVatExtractResultDetail

		if order.InvoiceDate != nil {
			invoiceDate = order.InvoiceDate.Format("2006-01-02")
		}

		if order.InvoiceNo != nil {
			invoiceNo = *order.InvoiceNo
		}

		orderDetails, err := service.CoreTaxVatExtractRepository.FindDetailsByRoNoAndItemType(order.RoNo, 1, custID, parentCustId)
		if err != nil {
			return data, err
		}

		orderDetailsPromo, err := service.CoreTaxVatExtractRepository.FindDetailsByRoNoAndItemType(order.RoNo, 2, custID, parentCustId)
		if err != nil {
			return data, err
		}

		var coretaxInvoiceOrderDetailReadMap = model.CoretaxInvoiceOrderDetailReadMap{}

		for _, orderDetailsPromo := range orderDetailsPromo {
			coretaxInvoiceOrderDetailReadMap.SetTempEmployeeValidationMap(orderDetailsPromo.ProId, orderDetailsPromo)
		}

		data.NPWPSeller = ""
		for _, orderDetail := range orderDetails {
			var promoTotalPrice1, promoTotalPrice2, promoTotalPrice3 float64
			var promoQty1Final, promoQty2Final, promoQty3Final float64
			orderDetailPromo, _ := coretaxInvoiceOrderDetailReadMap.GetByID(orderDetail.ProId)
			if orderDetailPromo != nil {
				orderDetailPromo.Extracted = true

				if orderDetailPromo.Qty1Final > 0 {
					promoTotalPrice1 = orderDetailPromo.SellPrice1 * orderDetailPromo.Qty1Final
					promoQty1Final = orderDetailPromo.Qty1Final
				}

				if orderDetailPromo.Qty2Final > 0 {
					promoTotalPrice2 = orderDetailPromo.SellPrice2 * orderDetailPromo.Qty2Final
					promoQty2Final = orderDetailPromo.Qty2Final

				}
				if orderDetailPromo.Qty3Final > 0 {
					promoTotalPrice3 = orderDetailPromo.SellPrice3 * orderDetailPromo.Qty3Final
					promoQty3Final = orderDetailPromo.Qty3Final

				}
			}

			grossTotal := (orderDetail.SellPrice1 * orderDetail.Qty1Final) + (orderDetail.SellPrice2 * orderDetail.Qty2Final) + (orderDetail.SellPrice3 * orderDetail.Qty3Final)
			if orderDetail.Qty1Final > 0 {
				totalPrice := (orderDetail.SellPrice1 * orderDetail.Qty1Final)

				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice1 * orderDetail.Qty1Final / grossTotal) * order.PromoValueFinal
				}

				if orderDetail.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice1 * orderDetail.Qty1Final / order.SubTotalFinal) * order.DiscValueFinal

				}
				promoDiscount := promo + discount + promoTotalPrice1

				DPP := totalPrice - promoDiscount

				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)

				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty1Final + promoQty1Final,
					Price:         orderDetail.SellPrice1,
					UnitId:        orderDetail.UnitId1,
					UnitIdCoretax: orderDetail.UnitIdCoreTax1,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}

			if orderDetail.Qty2Final > 0 {
				totalPrice := (orderDetail.SellPrice2 * orderDetail.Qty2Final)

				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice2 * orderDetail.Qty2Final / grossTotal) * order.PromoValueFinal
				}

				if order.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice2 * orderDetail.Qty2Final / order.SubTotalFinal) * order.DiscValueFinal
				}
				promoDiscount := promo + discount + promoTotalPrice2
				DPP := totalPrice - promoDiscount
				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)

				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty2Final + promoQty2Final,
					Price:         orderDetail.SellPrice2,
					UnitId:        orderDetail.UnitId2,
					UnitIdCoretax: orderDetail.UnitIdCoreTax2,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}

			if orderDetail.Qty3Final > 0 {
				totalPrice := (orderDetail.SellPrice3 * orderDetail.Qty3Final)

				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice3 * orderDetail.Qty3Final / grossTotal) * order.PromoValueFinal
				}
				if orderDetail.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice3 * orderDetail.Qty3Final / order.SubTotalFinal) * order.DiscValueFinal

				}
				promoDiscount := promo + discount + promoTotalPrice3

				DPP := totalPrice - promoDiscount
				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)
				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty3Final + promoQty3Final,
					Price:         orderDetail.SellPrice3,
					UnitId:        orderDetail.UnitId3,
					UnitIdCoretax: orderDetail.UnitIdCoreTax3,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}
		}

		for _, orderDetail := range orderDetailsPromo {
			orderDetailPromo, err := coretaxInvoiceOrderDetailReadMap.GetByID(orderDetail.ProId)
			if err != nil {
				return data, err
			}

			if orderDetailPromo != nil {
				if orderDetailPromo.Extracted {
					continue
				}
			}

			grossTotal := (orderDetail.SellPrice1 * orderDetail.Qty1Final) + (orderDetail.SellPrice2 * orderDetail.Qty2Final) + (orderDetail.SellPrice3 * orderDetail.Qty3Final)
			if orderDetail.Qty1Final > 0 {
				totalPrice := orderDetail.SellPrice1 * orderDetail.Qty1Final
				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice1 * orderDetail.Qty1Final / grossTotal) * orderDetail.PromoValueFinal
				}

				if orderDetail.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice1 * orderDetail.Qty1Final / order.SubTotalFinal) * order.DiscValueFinal

				}
				promoDiscount := promo + discount

				DPP := totalPrice - promoDiscount

				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)

				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty1Final,
					Price:         orderDetail.SellPrice1,
					UnitId:        orderDetail.UnitId1,
					UnitIdCoretax: orderDetail.UnitIdCoreTax1,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}

			if orderDetail.Qty2Final > 0 {
				totalPrice := orderDetail.SellPrice2 * orderDetail.Qty2Final
				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice2 * orderDetail.Qty2Final / grossTotal) * orderDetail.PromoValueFinal
				}

				if order.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice2 * orderDetail.Qty2Final / order.SubTotalFinal) * order.DiscValueFinal
				}
				promoDiscount := promo + discount
				DPP := totalPrice - promoDiscount
				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)

				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty2Final,
					Price:         orderDetail.SellPrice2,
					UnitId:        orderDetail.UnitId2,
					UnitIdCoretax: orderDetail.UnitIdCoreTax2,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}

			if orderDetail.Qty3Final > 0 {

				totalPrice := orderDetail.SellPrice3 * orderDetail.Qty3Final
				var promo, discount float64
				if orderDetail.PromoValueFinal > 0 {
					promo = (orderDetail.SellPrice3 * orderDetail.Qty3Final / grossTotal) * orderDetail.PromoValueFinal
				}

				if orderDetail.DiscValueFinal > 0 {
					discount = (orderDetail.SellPrice3 * orderDetail.Qty3Final / order.SubTotalFinal) * order.DiscValueFinal

				}
				promoDiscount := promo + discount

				DPP := totalPrice - discount
				DppOther := DPP * (11.0 / 12.0)
				PPNValue := DPP * (orderDetail.Vat / 100)
				invoiceDetailCoretaxExtract = append(invoiceDetailCoretaxExtract, entity.CoretaxVatExtractResultDetail{
					Item:          "A",
					ItemCode:      orderDetail.ProCodeCoretax,
					ItemName:      orderDetail.ProCode + " - " + orderDetail.ProName,
					Qty:           orderDetail.Qty3Final,
					Price:         orderDetail.SellPrice3,
					UnitId:        orderDetail.UnitId3,
					UnitIdCoretax: orderDetail.UnitIdCoreTax3,
					Discount:      promoDiscount,
					DPP:           DPP,
					DPPOther:      DppOther,
					PPN:           orderDetail.Vat,
					PPNValue:      PPNValue,
				})
			}
		}

		var buyerIDTku, buyerNPWP string
		if order.TaxIdentifierType == entity.TAX_IDENTIFIER_TYPE_TIN {
			buyerIDTku = order.TaxIdentifierNo
			buyerNPWP = order.TaxIdentifierNo + order.NITKU
			//taxIdentifierType = "-"
		} else {
			buyerIDTku = order.NITKU
			buyerNPWP = "0000000000000000"
			//taxIdentifierType = order.TaxIdentifierType
		}

		data.ExtractResults = append(data.ExtractResults, entity.CoretaxVatExtractResult{
			Row:                   index + 1,
			FakturDate:            invoiceDate,
			FakturType:            "Normal",
			TransactionCode:       "04",
			AdditionalDescription: "",
			DocumentSupport:       "",
			Reference:             invoiceNo,
			FacilityMark:          "",
			SellerTKUId:           "", // tanya mas agus
			BuyerNPWPorNIK:        buyerNPWP,
			BuyerTypeID:           order.TaxIdentifierType,
			BuyerCountry:          "IDN",
			BuyerDocumentNo:       order.IdentityNo,
			BuyerName:             order.TaxName,
			BuyerAddress:          order.AddressTax,
			BuyerEmail:            "",
			BuyerIDTKU:            buyerIDTku,
			Lists:                 invoiceDetailCoretaxExtract,
		})

	}

	return data, nil
}
