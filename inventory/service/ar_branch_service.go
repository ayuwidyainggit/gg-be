package service

import (
	"context"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/str"
	"inventory/pkg/structs"
	"inventory/repository"
)

type ArBranchService interface {
	Detail(arBranchNo string, custID, parentCustId string, queryParam entity.ArBranchDetailQuery) (response entity.ArBranchWithDetailResponse, err error)
	// DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.ArBranchWithDetailResponse, err error)
	// Store(request entity.CreateArBranchBody) (response entity.ArBranchResponse, err error)
	List(dataFilter entity.ArBranchQueryFilter, custId, parentCustId string) (data []entity.ArBranchListResponse, total int64, lastPage int, err error)
	StoreArBranchPayment(request entity.CreateArBranchPaymentBody) (response entity.ArBranchPaymentResponse, err error)
	// Update(arBranchNo string, request entity.UpdateArBranchRequest) (err error)
	// Delete(custId string, arBranchNo string, userId int64) (err error)
	// ListWarehouse(dataFilter entity.ArBranchWarehouseQueryFilter, custId, parentCustId string) (data []entity.ArBranchWarehouseListResponse, total int64, lastPage int, err error)
	// OrderBookingDetail(orderBookingId int, custID string, parentCustId string) (responses []entity.ArBranchOrderBookingDetailResponse, err error)
	// OrderBookingList(dataFilter entity.ArBranchOrderBookingListQueryFilter, custId string, parentCustId string) (data []entity.ArBranchOrderBookingListResponse, total int64, lastPage int, err error)
	// BulkUpdateStatus(request entity.ArBranchBulkUpdateDataStatus, custId string, parentCustId string) (err error)
	// BulkPrint(request entity.ArBranchBulkPrint, custId string, userId int64) (err error)
	DistributorsFilter(dataFilter entity.ArBranchDistributorsFilterQueryFilter, custId, parentCustId string) (data []entity.ArBranchDistributorsFilterListResponse, total int64, lastPage int, err error)
	SuppliersFilter(dataFilter entity.ArBranchSuppliersFilterQueryFilter, custId, parentCustId string) (data []entity.ArBranchSuppliersFilterListResponse, total int64, lastPage int, err error)
}

func NewArBranchService(
	arBranchRepository repository.ArBranchRepository,
	warehouseStockRepository repository.WarehouseStockRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction) *arBranchServiceImpl {
	return &arBranchServiceImpl{
		ArBranchRepository:       arBranchRepository,
		WarehouseStockRepository: warehouseStockRepository,
		StockRepository:          stockRepository,
		Transaction:              transaction,
	}
}

type arBranchServiceImpl struct {
	ArBranchRepository       repository.ArBranchRepository
	WarehouseStockRepository repository.WarehouseStockRepository
	StockRepository          repository.StockRepository
	Transaction              repository.Dbtransaction
}

/*
	func (service *arBranchServiceImpl) Store(request entity.CreateArBranchBody) (response entity.ArBranchResponse, err error) {
		c := context.Background()

		if len(request.Details.Normal) == 0 {
			return response, errors.New("item details is required")
		}

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

		var arBranchModel model.ArBranch
		if err = structs.Automapper(request, &arBranchModel); err != nil {
			return response, err
		}

		grDate := str.GetJakartaDate()
		arBranchModel.ArBranchDate = &grDate

		orderBooking, err := service.ArBranchRepository.FindArBranchOrderBooking(request.PoNo, request.CustID, request.ParentCustID)
		if err != nil {
			return response, err
		}

		dataStatusArBranch := entity.GR_BRANCH_SUBMITTED
		if *orderBooking.TypeApproval == 1 {
			dataStatusArBranch = entity.GR_BRANCH_PROCESSED
		}
		arBranchModel.DataStatus = &dataStatusArBranch

		var subTotal float64
		var vatValue float64
		var total float64
		var productIDs []int64
		var arBranchDetailModels []*model.ArBranchDetailCreate
		for index, detail := range request.Details.Normal {
			productIDs = append(productIDs, detail.ProID)

			QtyShipUnit := &conversion.QtyUnit{
				Qty1:      detail.QtyShip1,
				Qty2:      detail.QtyShip2,
				Qty3:      detail.QtyShip3,
				ConvUnit2: int(detail.ConvUnit2),
				ConvUnit3: int(detail.ConvUnit3),
			}

			totalQtyShip, err := QtyShipUnit.ToTotalQuantity()
			if err != nil {
				return response, err
			}

			ReceivedQtyUnit := &conversion.QtyUnit{
				Qty1:      detail.QtyReceived1,
				Qty2:      detail.QtyReceived2,
				Qty3:      detail.QtyReceived3,
				ConvUnit2: int(detail.ConvUnit2),
				ConvUnit3: int(detail.ConvUnit3),
			}

			totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
			if err != nil {
				return response, err
			}

			var arBranchDetailModel model.ArBranchDetailCreate

			if detail.QtyShip1 == 0 && detail.QtyShip2 == 0 && detail.QtyShip3 == 0 {
				return response, errors.New(fmt.Sprintf("please input qty_ship on id_product %v", detail.ProID))
			}

			if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
				return response, errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
			}

			err = structs.Automapper(detail, &arBranchDetailModel)
			if err != nil {
				return response, err
			}

			arBranchDetailModel.CustID = request.CustID
			// arBranchDetailModel.ArBranchNo = arBranchModel.ArBranchNo
			arBranchDetailModel.ItemType = model.ITEM_TYPE_NORMAL
			arBranchDetailModel.SeqNo = index + 1
			arBranchDetailModel.UnitPrice1 = detail.UnitPrice1
			arBranchDetailModel.UnitPrice2 = detail.UnitPrice2
			arBranchDetailModel.UnitPrice3 = detail.UnitPrice3
			arBranchDetailModel.QtyShip = totalQtyShip
			arBranchDetailModel.QtyReceived = totalReceivedQty

			amount := (float64(detail.QtyReceived1) * detail.UnitPrice1) + (float64(detail.QtyReceived2) * detail.UnitPrice2) + (float64(detail.QtyReceived3) * detail.UnitPrice3)
			subTotal += amount
			arBranchDetailModel.Amount = amount

			vatVal := math.Round((amount * *detail.Vat) / 100.0)
			arBranchDetailModel.VatValue = vatVal

			amount += vatVal
			vatValue += vatVal

			arBranchDetailModels = append(arBranchDetailModels, &arBranchDetailModel)
		}

		arBranchModel.SubTotal = &subTotal
		arBranchModel.VatValue = &vatValue
		total = subTotal + vatValue + *orderBooking.DeliveryFee
		arBranchModel.Total = &total
		arBranchModel.DeliveryFee = orderBooking.DeliveryFee

		for index, detail := range request.Details.Promo {
			if !slices.Contains(productIDs, detail.ProID) {
				productIDs = append(productIDs, detail.ProID)
			}

			QtyShipUnit := &conversion.QtyUnit{
				Qty1:      detail.QtyShip1,
				Qty2:      detail.QtyShip2,
				Qty3:      detail.QtyShip3,
				ConvUnit2: int(detail.ConvUnit2),
				ConvUnit3: int(detail.ConvUnit3),
			}

			totalQtyShip, err := QtyShipUnit.ToTotalQuantity()
			if err != nil {
				return response, err
			}

			ReceivedQtyUnit := &conversion.QtyUnit{
				Qty1:      detail.QtyReceived1,
				Qty2:      detail.QtyReceived2,
				Qty3:      detail.QtyReceived3,
				ConvUnit2: int(detail.ConvUnit2),
				ConvUnit3: int(detail.ConvUnit3),
			}

			totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
			if err != nil {
				return response, err
			}

			var arBranchDetailModel model.ArBranchDetailCreate

			if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
				return response, errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
			}

			if err = structs.Automapper(detail, &arBranchDetailModel); err != nil {
				return response, err
			}

			arBranchDetailModel.CustID = request.CustID
			// arBranchDetailModel.ArBranchNo = arBranchModel.ArBranchNo
			arBranchDetailModel.ItemType = model.ITEM_TYPE_PROMO
			arBranchDetailModel.SeqNo = index + 1
			arBranchDetailModel.UnitPrice1 = detail.UnitPrice1
			arBranchDetailModel.UnitPrice2 = detail.UnitPrice2
			arBranchDetailModel.UnitPrice3 = detail.UnitPrice3
			arBranchDetailModel.QtyShip = totalQtyShip
			arBranchDetailModel.QtyReceived = totalReceivedQty

			arBranchDetailModels = append(arBranchDetailModels, &arBranchDetailModel)
		}

		if _, err := service.ArBranchRepository.FindProductByListID(productIDs); err != nil {
			return response, err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

			// var stockUpdateEntities []*entity.StockUpdate

			if err = service.ArBranchRepository.Store(txCtx, &arBranchModel); err != nil {
				return err
			}
			response.ArBranchNo = arBranchModel.ArBranchNo

			for _, arBranchDetailModel := range arBranchDetailModels {
				var arBranchDetail model.ArBranchDetailCreate
				if err := structs.Automapper(*arBranchDetailModel, &arBranchDetail); err != nil {
					return err
				}
				// arBranchDetail := *arBranchDetailModel

				arBranchDetail.ArBranchNo = arBranchModel.ArBranchNo
				grDet, err := service.ArBranchRepository.CreateArBranchDetail(txCtx, &arBranchDetail)
				if err != nil {
					return err
				}

				fmt.Println("Gr Branch Detail ID arBranchDetailModel : ", arBranchDetailModel.ArBranchDetId)
				fmt.Println("Item Type arBranchDetailModel : ", arBranchDetailModel.ItemType)
				fmt.Println("Gr Branch Detail ID arBranchDetail : ", grDet.ArBranchDetId)
				fmt.Println("Item Type arBranchDetail : ", grDet.ItemType)
				// stockUpdateEntity := entity.StockUpdate{
				// 	CustID:    arBranchModel.CustID,
				// 	WhID:      *arBranchModel.WhID,
				// 	ProID:     grDet.ProID,
				// 	StockDate: *arBranchModel.ArBranchDate,
				// 	TrCode:    arBranchModel.ArBranchNo[0:3],
				// 	TrNo:      arBranchModel.ArBranchNo,
				// 	QtyIn:     float64(grDet.QtyReceived),
				// 	UnitPrice: grDet.UnitPrice1,
				// 	RefDetId:  grDet.ArBranchDetId,
				// }

				// stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
			}

			// if *arBranchModel.DataStatus == entity.GR_BRANCH_PROCESSED {
			// 	var arBranchDetails []model.ArBranchDetailList
			// 	var stockUpdateEntities []*entity.StockUpdate

			// 	arBranchDetails, err = service.ArBranchRepository.FindArBranchdetailWithDiscount(arBranchModel.ArBranchNo, arBranchModel.CustID)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	for _, arBranchDetail := range arBranchDetails {
			// 		ReceivedQtyUnit := &conversion.QtyUnit{
			// 			Qty1:      int(*arBranchDetail.QtyReceived1),
			// 			Qty2:      int(*arBranchDetail.QtyReceived2),
			// 			Qty3:      int(*arBranchDetail.QtyReceived3),
			// 			ConvUnit2: int(arBranchDetail.ConvUnit2),
			// 			ConvUnit3: int(arBranchDetail.ConvUnit3),
			// 		}

			// 		totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
			// 		if err != nil {
			// 			return err
			// 		}

			// 		stockUpdateEntity := entity.StockUpdate{
			// 			CustID:    arBranchModel.CustID,
			// 			WhID:      *arBranchModel.WhID,
			// 			ProID:     arBranchDetail.ProID,
			// 			StockDate: *arBranchModel.ArBranchDate,
			// 			TrCode:    arBranchModel.ArBranchNo[0:3],
			// 			TrNo:      arBranchModel.ArBranchNo,
			// 			QtyIn:     float64(totalReceivedQty),
			// 			UnitPrice: arBranchDetail.UnitPrice1,
			// 			RefDetId:  int64(arBranchDetail.ArBranchDetId),
			// 		}

			// 		stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
			// 	}

			// 	err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			return nil
		})

		if *arBranchModel.DataStatus == entity.GR_BRANCH_PROCESSED {
			err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
				var arBranchDetails []model.ArBranchDetailList
				var stockUpdateEntities []*entity.StockUpdate

				arBranchDetails, err = service.ArBranchRepository.FindArBranchdetailWithDiscount(arBranchModel.ArBranchNo, arBranchModel.CustID)
				if err != nil {
					return err
				}

				for _, arBranchDetail := range arBranchDetails {
					ReceivedQtyUnit := &conversion.QtyUnit{
						Qty1:      int(*arBranchDetail.QtyReceived1),
						Qty2:      int(*arBranchDetail.QtyReceived2),
						Qty3:      int(*arBranchDetail.QtyReceived3),
						ConvUnit2: int(arBranchDetail.ConvUnit2),
						ConvUnit3: int(arBranchDetail.ConvUnit3),
					}

					totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
					if err != nil {
						return err
					}

					stockUpdateEntity := entity.StockUpdate{
						CustID:    arBranchModel.CustID,
						WhID:      *arBranchModel.WhID,
						ProID:     arBranchDetail.ProID,
						StockDate: *arBranchModel.ArBranchDate,
						TrCode:    arBranchModel.ArBranchNo[0:3],
						TrNo:      arBranchModel.ArBranchNo,
						QtyIn:     float64(totalReceivedQty),
						UnitPrice: arBranchDetail.UnitPrice1,
						RefDetId:  int64(arBranchDetail.ArBranchDetId),
					}

					// fmt.Println("CustID : ", arBranchModel.CustID)
					// fmt.Println("WhID : ", *arBranchModel.WhID)
					// fmt.Println("ProID : ", arBranchDetail.ProID)
					// fmt.Println("StockDate : ", *arBranchModel.ArBranchDate)
					// fmt.Println("TrCode : ", arBranchModel.ArBranchNo[0:3])
					// fmt.Println("TrNo : ", arBranchModel.ArBranchNo)
					// fmt.Println("QtyIn : ", float64(totalReceivedQty))
					// fmt.Println("UnitPrice : ", arBranchDetail.UnitPrice1)
					// fmt.Println("RefDetId : ", int64(arBranchDetail.ArBranchDetId))

					stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
				}

				err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
				if err != nil {
					return err
				}

				return nil
			})
		}

		return response, err
	}
*/
func (service *arBranchServiceImpl) StoreArBranchPayment(request entity.CreateArBranchPaymentBody) (response entity.ArBranchPaymentResponse, err error) {
	c := context.Background()
	var arBranchPaymentModel model.ArBranchPaymentCreate
	if err = structs.Automapper(request, &arBranchPaymentModel); err != nil {
		return response, err
	}

	arBranch, err := service.ArBranchRepository.FindByNo(*request.GrBranchNo, request.CustID, *request.ParentCustID)
	if err != nil {
		return response, err
	}

	paymentDate := str.GetJakartaDate()
	arBranchPaymentModel.DepositDate = &paymentDate

	// orderBooking, err := service.ArBranchRepository.FindArBranchOrderBooking(request.PoNo, request.CustID, request.ParentCustID)
	// if err != nil {
	// 	return response, err
	// }
	dataPaymentType := entity.AR_BRANCH_PAYMENT_TYPE_CASH
	arBranchPaymentModel.PaymentType = &dataPaymentType

	dataPaymentOption := entity.AR_BRANCH_PAYMENT_OPTION_PARTIAL
	if arBranchPaymentModel.PaymentAmount != nil && arBranchPaymentModel.Discount != nil && arBranchPaymentModel.PaymentBalance != nil {
		if (*arBranchPaymentModel.PaymentAmount + *arBranchPaymentModel.Discount + *arBranchPaymentModel.PaymentBalance) >= *arBranch.Total {
			dataPaymentOption = entity.AR_BRANCH_PAYMENT_OPTION_FULL
		}
	}
	arBranchPaymentModel.PaymentOption = &dataPaymentOption

	dataVerificationStatus := entity.AR_BRANCH_VERIFICATION_STATUS_NEED_REVIEW
	arBranchPaymentModel.VerificationStatus = &dataVerificationStatus

	// dataPaymentAmount := *arBranchPaymentModel.TotalPayment + *arBranchPaymentModel.Discount - *arBranchPaymentModel.PaymentBalance
	// arBranchPaymentModel.PaymentAmount = &dataPaymentAmount

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		// var stockUpdateEntities []*entity.StockUpdate

		if err = service.ArBranchRepository.StoreArBranchPayment(txCtx, &arBranchPaymentModel); err != nil {
			return err
		}
		response.DepositNo = *arBranchPaymentModel.DepositNo

		return nil
	})

	return response, err
}

func (service *arBranchServiceImpl) Detail(arBranchNo string, custID, parentCustId string, queryParam entity.ArBranchDetailQuery) (response entity.ArBranchWithDetailResponse, err error) {
	gr, err := service.ArBranchRepository.FindByNo(arBranchNo, queryParam.ArBranchCustId, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gr, &response)
	if err != nil {
		return response, err
	}

	err = service.GetArBranchDetail(arBranchNo, queryParam.ArBranchCustId, &response)
	if err != nil {
		return response, err
	}

	grDate := gr.GrBranchDate.Format("2006-01-02")
	response.GrBranchDate = grDate

	deliveryDate := gr.DeliveryDate.Format("2006-01-02")
	response.DeliveryDate = deliveryDate

	if gr.InvoiceDate != nil {
		invoiceDate := gr.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	dataStatusName := response.GenerateDataStatusName()
	response.DataStatusName = dataStatusName

	typeApprovalName := response.GenerateOrderBookingTypeApprovalName()
	response.TypeApprovalName = &typeApprovalName

	invoiceStatus := int64(entity.AR_BRANCH_INVOICE_STATUS_OUTSTANDING)
	if *response.RemainingAmount <= 0 {
		invoiceStatus = entity.AR_BRANCH_INVOICE_STATUS_PAID
	}
	response.InvoiceStatus = &invoiceStatus

	invoiceStatusName := response.GenerateInvoiceStatusName()
	response.InvoiceStatusName = &invoiceStatusName

	if gr.InvoiceDateBranch != nil {
		invoiceDateBranch := gr.InvoiceDateBranch.Format("2006-01-02")
		response.InvoiceDateBranch = &invoiceDateBranch
	}

	if gr.InvoiceDueDateBranch != nil {
		invoiceDueDateBranch := gr.InvoiceDueDateBranch.Format("2006-01-02")
		response.InvoiceDueDateBranch = &invoiceDueDateBranch
	}

	return response, nil
}

/*
	func (service *arBranchServiceImpl) DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.ArBranchWithDetailResponse, err error) {
		gr, err := service.ArBranchRepository.GetByInvoiceNo(invoice, custID, parentCustId)
		if err != nil {
			return response, err
		}

		err = structs.Automapper(gr, &response)
		if err != nil {
			return response, err
		}

		err = service.GetArBranchDetail(invoice, custID, &response)
		if err != nil {
			return response, err
		}

		grData := gr.ArBranchDate.Format("2006-01-02")
		// deliveryDate := gr.DeliveryDate.Format("2006-01-02")
		// invoiceDate := gr.InvoiceDate.Format("2006-01-02")

		response.ArBranchDate = grData
		// response.DeliveryDate = deliveryDate
		// response.InvoiceDate = invoiceDate
		return response, nil
	}
*/
func (service *arBranchServiceImpl) GetArBranchDetail(arBranchNo string, custID string, arBranch *entity.ArBranchWithDetailResponse) (err error) {
	var arBranchDetails []model.ArBranchDetailList
	// var discountValueTotal float64

	arBranchDetails, err = service.ArBranchRepository.FindArBranchDetailWithDiscount(arBranchNo, custID)
	if err != nil {
		return err
	}

	arBranch.Details.Promo = []entity.ArBranchDetailList{}
	arBranch.Details.Normal = []entity.ArBranchDetailList{}

	for _, arBranchDetail := range arBranchDetails {
		var arBranchDetailData entity.ArBranchDetailList
		err = structs.Automapper(arBranchDetail, &arBranchDetailData)
		if err != nil {
			return err
		}

		qtyShip := &conversion.Qty{
			Qty:       int(arBranchDetailData.QtyShip),
			ConvUnit2: int(arBranchDetailData.ConvUnit2),
			ConvUnit3: int(arBranchDetailData.ConvUnit3),
		}

		qtyShipConversion := qtyShip.ConvToQtyConversion()

		qtyReceived := &conversion.Qty{
			Qty:       int(arBranchDetailData.QtyReceived),
			ConvUnit2: int(arBranchDetailData.ConvUnit2),
			ConvUnit3: int(arBranchDetailData.ConvUnit3),
		}
		qtyReceivedConversion := qtyReceived.ConvToQtyConversion()

		arBranchDetailData.QtyShip1 = qtyShipConversion.Qty1
		arBranchDetailData.QtyShip2 = qtyShipConversion.Qty2
		arBranchDetailData.QtyShip3 = qtyShipConversion.Qty3

		arBranchDetailData.QtyReceived1 = qtyReceivedConversion.Qty1
		arBranchDetailData.QtyReceived2 = qtyReceivedConversion.Qty2
		arBranchDetailData.QtyReceived3 = qtyReceivedConversion.Qty3

		if arBranchDetail.ItemType == 1 {
			// var discountValue, discount float64
			// Subtotal := (arBranchDetail.UnitPrice1 * float64(arBranchDetailData.QtyReceived1)) + (arBranchDetail.UnitPrice2 * float64(arBranchDetailData.QtyReceived2)) + (arBranchDetail.UnitPrice3 * float64(arBranchDetailData.QtyReceived3))
			// if isAp {
			// 	if arBranchDetail.Discount != nil {
			// 		discountValue = (*arBranchDetail.Discount / 100) * Subtotal
			// 		discount = *arBranchDetail.Discount
			// 	}
			// 	arBranchDetailData.DiscountValue = &discountValue
			// 	arBranchDetailData.Discount = &discount

			// 	discountValueTotal += discountValue
			// }

			// ppn := (Subtotal - discountValue) * arBranchDetail.Vat / 100
			// pbn := (Subtotal - discountValue) * arBranchDetail.VatLgPurch / 100
			//ppnDp := (Subtotal - discountValue) * arBranchDetail.VatBg / 100
			// total := (Subtotal - discountValue) + ppn + pbn

			// arBranch.SubTotal += Subtotal
			// arBranch.SubTotal += Subtotal - discountValue
			// arBranch.Total += total
			// arBranch.TotalVat += ppn
			// arBranch.TotalVatLgPurch += pbn
			// arBranch.TotalSkuPrice += Subtotal - discountValue
			// arBranchDetailData.Nett = Subtotal - discountValue
			// arBranchDetailData.SubTotal = Subtotal
			// arBranchDetailData.Total = total
			// arBranchDetailData.VatValue = ppn
			// arBranchDetailData.VatLgPurchValue = pbn

			// arBranchDetailData.ConvUnit1 = arBranchDetail.ConvUnit2 * arBranchDetail.ConvUnit3
			arBranch.Details.Normal = append(arBranch.Details.Normal, arBranchDetailData)
		} else {
			arBranch.Details.Promo = append(arBranch.Details.Promo, arBranchDetailData)
		}
	}

	if arBranch.Details.Promo == nil {
		arBranch.Details.Promo = []entity.ArBranchDetailList{}
	}

	var arBranchPayments []model.ArBranchPaymentList
	// var discountValueTotal float64

	arBranchPayments, err = service.ArBranchRepository.FindArBranchPayments(*arBranch.InvoiceNoBranch, custID)
	if err != nil {
		return err
	}

	// arBranch.Payments = []entity.ArBranchPaymentList{}
	for _, arBranchPayment := range arBranchPayments {
		var arBranchPaymentData entity.ArBranchPaymentResponse

		if err = structs.Automapper(arBranchPayment, &arBranchPaymentData); err != nil {
			return err
		}

		if arBranchPaymentData.DepositDate != nil {
			depositDate := arBranchPayment.DepositDate.Format("2006-01-02")
			arBranchPaymentData.DepositDate = &depositDate
		}

		paymentOptionName := arBranchPaymentData.GeneratePaymentOptionName()
		arBranchPaymentData.PaymentOptionName = &paymentOptionName

		paymentTypeName := arBranchPaymentData.GeneratePaymentTypeName()
		arBranchPaymentData.PaymentTypeName = &paymentTypeName

		verificationStatusName := arBranchPaymentData.GenerateVerificationStatusName()
		arBranchPaymentData.VerificationStatusName = &verificationStatusName

		arBranch.Payments = append(arBranch.Payments, arBranchPaymentData)
	}

	if arBranch.Payments == nil {
		arBranch.Payments = []entity.ArBranchPaymentResponse{}
	}

	return
}

func (service *arBranchServiceImpl) List(dataFilter entity.ArBranchQueryFilter, custId, parentCustId string) (data []entity.ArBranchListResponse, total int64, lastPage int, err error) {
	grs, total, lastPage, err := service.ArBranchRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grs {
		var vResp entity.ArBranchListResponse
		structs.Automapper(row, &vResp)

		grDate := row.GrBranchDate.Format("2006-01-02")
		vResp.GrBranchDate = grDate

		if row.DeliveryDate != nil {
			deliveryDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &deliveryDate
		}

		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}

		if row.InvoiceDateBranch != nil {
			invoiceDateBranch := row.InvoiceDateBranch.Format("2006-01-02")
			vResp.InvoiceDateBranch = &invoiceDateBranch
		}

		vResp.DataStatusName = vResp.GenerateDataStatusName()

		invoiceStatus := int64(entity.AR_BRANCH_INVOICE_STATUS_OUTSTANDING)
		if *vResp.RemainingAmount <= 0 {
			invoiceStatus = entity.AR_BRANCH_INVOICE_STATUS_PAID
		}
		vResp.InvoiceStatus = &invoiceStatus

		invoiceStatusName := vResp.GenerateInvoiceStatusName()
		vResp.InvoiceStatusName = &invoiceStatusName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

/*
func (service *arBranchServiceImpl) ListSupplier(dataFilter entity.ArBranchSupplierQueryFilter, custId, parentCustId string) (data []entity.ArBranchSupplierListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.ArBranchRepository.FindSupplierArBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.ArBranchSupplierListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *arBranchServiceImpl) ListWarehouse(dataFilter entity.ArBranchWarehouseQueryFilter, custId, parentCustId string) (data []entity.ArBranchWarehouseListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.ArBranchRepository.FindWarehouseArBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.ArBranchWarehouseListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *arBranchServiceImpl) Update(grNo string, request entity.UpdateArBranchRequest) (err error) {
	c := context.Background()

	if request.DeliveryDate != nil {
		// parse time format YYYY-mm-dd to Rfc3339
		if *request.DeliveryDate != "" {
			deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
			if err != nil {
				return err
			}
			request.DeliveryDate = &deliveryDate
		}
	}

	var grModel model.ArBranch
	err = structs.Automapper(request, &grModel)
	if err != nil {
		return err
	}
	grModel.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		gr, err := service.ArBranchRepository.FindByNo(grNo, request.CustID, request.ParentCustID)
		if err != nil {
			return err
		}

		err = service.ArBranchRepository.Update(txCtx, gr.ArBranchNo, &grModel)
		if err != nil {
			return err
		}

		// grDetails, err := service.ArBranchRepository.FindArBranchdetail(grNo, request.CustID)
		// if err != nil {
		// 	return err
		// }

		// for _, grDet := range grDetails {
		// 	err = service.WhStockRepository.UpdateOldWhStock(txCtx, request.CustID, *request.WhID, grDet.ProID, grDet.Qty)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		err = service.ArBranchRepository.DeleteArBranchDetailByArBranchNo(txCtx, grNo)
		if err != nil {
			return err
		}

		newArBranchDetIds := make([]int64, 0)
		for _, detail := range request.Details.Normal {
			//  if detail.ArBranchDetId == nil || *detail.ArBranchDetId == 0 {
			detail.ArBranchDetId = nil
			var grDetailModel model.ArBranchDetailCreate
			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}
			// grDetailModel.SeqNo = sequence
			grDetailModel.CustID = request.CustID
			grDetailModel.ArBranchNo = grNo
			// grDet, err := service.ArBranchRepository.CreateArBranchDetail(txCtx, &grDetailModel)
			_, err := service.ArBranchRepository.CreateArBranchDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			// newArBranchDetIds = append(newArBranchDetIds, grDet.ArBranchDetId)
			newArBranchDetIds = append(newArBranchDetIds, grDetailModel.ProID)
		}

		for _, detail := range request.Details.Promo {
			detail.ArBranchDetId = nil
			var grDetailModel model.ArBranchDetailCreate
			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}
			grDetailModel.CustID = request.CustID
			grDetailModel.ArBranchNo = grNo
			_, err := service.ArBranchRepository.CreateArBranchDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			// newArBranchDetIds = append(newArBranchDetIds, grDet.ArBranchDetId)
			newArBranchDetIds = append(newArBranchDetIds, grDetailModel.ProID)
		}

		err = service.ArBranchRepository.DeleteStockNotInRefIds(txCtx, request.CustID, gr.ArBranchNo, newArBranchDetIds)
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

func (service *arBranchServiceImpl) Delete(custId string, grNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		grDetails, err := service.ArBranchRepository.FindArBranchdetailForUpdateWhStock(grNo, custId)
		if err != nil {
			return err
		}
		log.Println("grDetails:", structs.StructToJson(grDetails))

		oldArBranchDetId := make([]int64, 0)
		for _, grDet := range grDetails {
			err = service.ArBranchRepository.UpdateOldWhStock(txCtx, custId, grDet.WhID, grDet.ProID, *grDet.Qty)
			if err != nil {
				return err
			}
			oldArBranchDetId = append(oldArBranchDetId, grDet.ArBranchDetId)
		}

		log.Println("oldArBranchDetId:", structs.StructToJson(oldArBranchDetId))
		log.Println("oldArBranchDetId:", structs.StructToJson(oldArBranchDetId))
		if len(oldArBranchDetId) > 0 {
			err = service.ArBranchRepository.DeleteStockInRefIds(txCtx, custId, grNo, oldArBranchDetId)
			if err != nil {
				return err
			}
		}

		err = service.ArBranchRepository.Delete(txCtx, custId, grNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *arBranchServiceImpl) OrderBookingDetail(orderBookingId int, custID string, parentCustId string) (responses []entity.ArBranchOrderBookingDetailResponse, err error) {
	orderBookingDetails, err := service.ArBranchRepository.FindArBranchOrderBookingDetails(orderBookingId, custID, parentCustId)
	if err != nil {
		return responses, err
	}

	for _, orderBookingDetail := range orderBookingDetails {
		var response entity.ArBranchOrderBookingDetailResponse
		err = structs.Automapper(orderBookingDetail, &response)
		if err != nil {
			return responses, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *arBranchServiceImpl) OrderBookingList(dataFilter entity.ArBranchOrderBookingListQueryFilter, custId string, parentCustId string) (data []entity.ArBranchOrderBookingListResponse, total int64, lastPage int, err error) {
	orderBookingList, total, lastPage, err := service.ArBranchRepository.FindArBranchOrderBookingList(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range orderBookingList {
		var vResp entity.ArBranchOrderBookingListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *arBranchServiceImpl) BulkUpdateStatus(request entity.ArBranchBulkUpdateDataStatus, custId string, parentCustId string) (err error) {
	c := context.Background()

	for index := range request.ArBranches {
		// End parse time format YYYY-mm-dd to Rfc339
		var Model model.ArBranch
		err = structs.Automapper(request.ArBranches[index], &Model)
		if err != nil {
			return err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

			err = service.ArBranchRepository.Update(txCtx, request.ArBranches[index].ArBranchNo, &Model)
			if err != nil {
				return err
			}

			if *request.ArBranches[index].DataStatus == entity.GR_BRANCH_PROCESSED { // jika processed, trigger stock
				arBranch, err := service.ArBranchRepository.FindByNo(request.ArBranches[index].ArBranchNo, request.ArBranches[index].CustId, parentCustId)
				if err != nil {
					return err
				}

				// details, err := service.ArBranchRepository.FindArBranchdetailWithDiscount(request.ArBranches[index].ArBranchNo, custId)
				// if err != nil {
				// 	return err
				// }

				// var salesDetailCanceledUpdates []*entity.SalesOrderStockUpdate

				// for _, detail := range details {
				// 	salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
				// 		CustID:         custId,
				// 		WhID:           *arBranch.WhId,
				// 		ProID:          int64(detail.ProId),
				// 		StockDate:      *arBranch.RoDate,
				// 		TrCode:         request.Orders[index].RoNo[0:2],
				// 		TrNo:           request.Orders[index].RoNo,
				// 		QtyOrderBefore: detail.Qty,
				// 		QtyOrder:       0,
				// 		UnitPrice:      *detail.SellPrice1,
				// 		RefDetId:       int64(*detail.OrderDetailID),
				// 	}
				// 	salesDetailCanceledUpdates = append(salesDetailCanceledUpdates, &salesDetailCanceledUpdate)

				// }

				// if len(salesDetailCanceledUpdates) > 0 {
				// 	log.Info("Update Stock Deleted Details")
				// 	err = service.StockRepository.SalesStockUpdates(txCtx, salesDetailCanceledUpdates)
				// 	if err != nil {
				// 		return err
				// 	}
				// }

				var arBranchDetails []model.ArBranchDetailList
				var stockUpdateEntities []*entity.StockUpdate

				arBranchDetails, err = service.ArBranchRepository.FindArBranchdetailWithDiscount(request.ArBranches[index].ArBranchNo, request.ArBranches[index].CustId)
				if err != nil {
					return err
				}

				for _, arBranchDetail := range arBranchDetails {
					ReceivedQtyUnit := &conversion.QtyUnit{
						Qty1:      int(*arBranchDetail.QtyReceived1),
						Qty2:      int(*arBranchDetail.QtyReceived2),
						Qty3:      int(*arBranchDetail.QtyReceived3),
						ConvUnit2: int(arBranchDetail.ConvUnit2),
						ConvUnit3: int(arBranchDetail.ConvUnit3),
					}

					totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
					if err != nil {
						return err
					}

					stockUpdateEntity := entity.StockUpdate{
						CustID:    custId,
						WhID:      *arBranch.WhId,
						ProID:     arBranchDetail.ProID,
						StockDate: *arBranch.ArBranchDate,
						TrCode:    arBranch.ArBranchNo[0:3],
						TrNo:      arBranch.ArBranchNo,
						QtyIn:     float64(totalReceivedQty),
						UnitPrice: arBranchDetail.UnitPrice1,
						RefDetId:  int64(arBranchDetail.ArBranchDetId),
					}

					// fmt.Println("CustID : ", arBranchModel.CustID)
					// fmt.Println("WhID : ", *arBranchModel.WhID)
					// fmt.Println("ProID : ", arBranchDetail.ProID)
					// fmt.Println("StockDate : ", *arBranchModel.ArBranchDate)
					// fmt.Println("TrCode : ", arBranchModel.ArBranchNo[0:3])
					// fmt.Println("TrNo : ", arBranchModel.ArBranchNo)
					// fmt.Println("QtyIn : ", float64(totalReceivedQty))
					// fmt.Println("UnitPrice : ", arBranchDetail.UnitPrice1)
					// fmt.Println("RefDetId : ", int64(arBranchDetail.ArBranchDetId))

					stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
				}

				err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
				if err != nil {
					return err
				}
			}


				// if *request.ArBranches[index].DataStatus == entity.GR_BRANCH_REJECTED { // jika cancel, trigger stock
				// 	arBranch, err := service.ArBranchRepository.FindByNo(request.ArBranches[index].ArBranchNo, custId, parentCustId)
				// 	if err != nil {
				// 		return err
				// 	}

				// 	details, err := service.ArBranchRepository.FindArBranchdetailWithDiscount(request.ArBranches[index].ArBranchNo, custId)
				// 	if err != nil {
				// 		return err
				// 	}


				// 	var salesDetailCanceledUpdates []*entity.SalesOrderStockUpdate

				// 	for _, detail := range details {
				// 		salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
				// 			CustID:         custId,
				// 			WhID:           *ro.WhId,
				// 			ProID:          int64(detail.ProId),
				// 			StockDate:      *ro.RoDate,
				// 			TrCode:         request.Orders[index].RoNo[0:2],
				// 			TrNo:           request.Orders[index].RoNo,
				// 			QtyOrderBefore: detail.Qty,
				// 			QtyOrder:       0,
				// 			UnitPrice:      *detail.SellPrice1,
				// 			RefDetId:       int64(*detail.OrderDetailID),
				// 		}
				// 		salesDetailCanceledUpdates = append(salesDetailCanceledUpdates, &salesDetailCanceledUpdate)

				// 	}

				// 	if len(salesDetailCanceledUpdates) > 0 {
				// 		log.Info("Update Stock Deleted Details")
				// 		err = service.StockRepository.SalesStockUpdates(txCtx, salesDetailCanceledUpdates)
				// 		if err != nil {
				// 			return err
				// 		}
				// 	}
				// }

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *arBranchServiceImpl) BulkPrint(request entity.ArBranchBulkPrint, custId string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		fmt.Println("ArBranchService Print")
		err = service.ArBranchRepository.PrintArBranch(txCtx, custId, request.ArBranches, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}
*/

func (service *arBranchServiceImpl) DistributorsFilter(dataFilter entity.ArBranchDistributorsFilterQueryFilter, custId, parentCustId string) (data []entity.ArBranchDistributorsFilterListResponse, total int64, lastPage int, err error) {
	distributors, total, lastPage, err := service.ArBranchRepository.FindDistributorsArBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range distributors {
		var vResp entity.ArBranchDistributorsFilterListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *arBranchServiceImpl) SuppliersFilter(dataFilter entity.ArBranchSuppliersFilterQueryFilter, custId, parentCustId string) (data []entity.ArBranchSuppliersFilterListResponse, total int64, lastPage int, err error) {
	suppliers, total, lastPage, err := service.ArBranchRepository.FindSuppliersArBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range suppliers {
		var vResp entity.ArBranchSuppliersFilterListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}
