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
	"log"
	"math"
	"slices"
	"time"
)

type GrBranchService interface {
	Detail(grBranchNo string, custID, parentCustId string, queryParam entity.GrBranchDetailQuery) (response entity.GrBranchWithDetailResponse, err error)
	DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.GrBranchWithDetailResponse, err error)
	Store(request entity.CreateGrBranchBody) (response entity.GrBranchResponse, err error)
	List(dataFilter entity.GrBranchQueryFilter, custId, parentCustId string) (data []entity.GrBranchListResponse, total int64, lastPage int, err error)
	Update(grBranchNo string, request entity.UpdateGrBranchRequest) (err error)
	Delete(custId string, grBranchNo string, userId int64) (err error)
	ListSupplier(dataFilter entity.GrBranchSupplierQueryFilter, custId, parentCustId string) (data []entity.GrBranchSupplierListResponse, total int64, lastPage int, err error)
	ListDistributor(dataFilter entity.GrBranchDistributorQueryFilter, custId, parentCustId string) (data []entity.GrBranchDistributorListResponse, total int64, lastPage int, err error)
	ListWarehouse(dataFilter entity.GrBranchWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrBranchWarehouseListResponse, total int64, lastPage int, err error)
	ListPrintWarehouse(dataFilter entity.GrBranchPrintWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrBranchWarehouseListResponse, total int64, lastPage int, err error)
	OrderBookingDetail(orderBookingId int, custID string, parentCustId string) (responses []entity.GrBranchOrderBookingDetailResponse, err error)
	OrderBookingList(dataFilter entity.GrBranchOrderBookingListQueryFilter, custId string, parentCustId string) (data []entity.GrBranchOrderBookingListResponse, total int64, lastPage int, err error)
	BulkUpdateStatus(request entity.GrBranchBulkUpdateDataStatus, custId string, parentCustId string) (err error)
	BulkPrint(request entity.GrBranchBulkPrint, custId string, parentCustId string, userId int64) (err error)
}

func NewGrBranchService(
	orderBookingRepository repository.OrderBookingRepository,
	grBranchRepository repository.GrBranchRepository,
	warehouseStockRepository repository.WarehouseStockRepository,
	stockRepository repository.StockRepository,
	transaction repository.Dbtransaction) *grBranchServiceImpl {
	return &grBranchServiceImpl{
		OrderBookingRepository:   orderBookingRepository,
		GrBranchRepository:       grBranchRepository,
		WarehouseStockRepository: warehouseStockRepository,
		StockRepository:          stockRepository,
		Transaction:              transaction,
	}
}

type grBranchServiceImpl struct {
	OrderBookingRepository   repository.OrderBookingRepository
	GrBranchRepository       repository.GrBranchRepository
	WarehouseStockRepository repository.WarehouseStockRepository
	StockRepository          repository.StockRepository
	Transaction              repository.Dbtransaction
}

func (service *grBranchServiceImpl) Store(request entity.CreateGrBranchBody) (response entity.GrBranchResponse, err error) {
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

	var grBranchModel model.GrBranch
	if err = structs.Automapper(request, &grBranchModel); err != nil {
		return response, err
	}

	grDate := str.GetJakartaDate()
	grBranchModel.GrBranchDate = &grDate

	orderBooking, err := service.GrBranchRepository.FindGrBranchOrderBooking(request.PoNo, request.CustID, request.ParentCustID)
	if err != nil {
		return response, err
	}

	dataStatusGrBranch := entity.GR_BRANCH_SUBMITTED
	if *orderBooking.TypeApproval == 1 {
		dataStatusGrBranch = entity.GR_BRANCH_PROCESSED
	}
	grBranchModel.DataStatus = &dataStatusGrBranch

	var subTotal float64
	var vatValue float64
	var total float64
	var productIDs []int64
	var grBranchDetailModels []*model.GrBranchDetailCreate
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

		var grBranchDetailModel model.GrBranchDetailCreate

		if detail.QtyShip1 == 0 && detail.QtyShip2 == 0 && detail.QtyShip3 == 0 {
			return response, errors.New(fmt.Sprintf("please input qty_ship on id_product %v", detail.ProID))
		}

		if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
			return response, errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
		}

		err = structs.Automapper(detail, &grBranchDetailModel)
		if err != nil {
			return response, err
		}

		grBranchDetailModel.CustID = request.CustID
		// grBranchDetailModel.GrBranchNo = grBranchModel.GrBranchNo
		grBranchDetailModel.ItemType = model.ITEM_TYPE_NORMAL
		grBranchDetailModel.SeqNo = index + 1
		grBranchDetailModel.UnitPrice1 = detail.UnitPrice1
		grBranchDetailModel.UnitPrice2 = detail.UnitPrice2
		grBranchDetailModel.UnitPrice3 = detail.UnitPrice3
		grBranchDetailModel.QtyShip = totalQtyShip
		grBranchDetailModel.QtyReceived = totalReceivedQty

		amount := (float64(detail.QtyReceived1) * detail.UnitPrice1) + (float64(detail.QtyReceived2) * detail.UnitPrice2) + (float64(detail.QtyReceived3) * detail.UnitPrice3)
		subTotal += amount
		grBranchDetailModel.Amount = amount

		vatVal := math.Round((amount * *detail.Vat) / 100.0)
		grBranchDetailModel.VatValue = vatVal

		amount += vatVal
		vatValue += vatVal

		grBranchDetailModels = append(grBranchDetailModels, &grBranchDetailModel)
	}

	grBranchModel.SubTotal = &subTotal
	grBranchModel.VatValue = &vatValue
	total = subTotal + vatValue + *orderBooking.DeliveryFee
	grBranchModel.Total = &total
	grBranchModel.DeliveryFee = orderBooking.DeliveryFee

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

		var grBranchDetailModel model.GrBranchDetailCreate

		if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
			return response, errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
		}

		if err = structs.Automapper(detail, &grBranchDetailModel); err != nil {
			return response, err
		}

		grBranchDetailModel.CustID = request.CustID
		// grBranchDetailModel.GrBranchNo = grBranchModel.GrBranchNo
		grBranchDetailModel.ItemType = model.ITEM_TYPE_PROMO
		grBranchDetailModel.SeqNo = index + 1
		grBranchDetailModel.UnitPrice1 = detail.UnitPrice1
		grBranchDetailModel.UnitPrice2 = detail.UnitPrice2
		grBranchDetailModel.UnitPrice3 = detail.UnitPrice3
		grBranchDetailModel.QtyShip = totalQtyShip
		grBranchDetailModel.QtyReceived = totalReceivedQty

		grBranchDetailModels = append(grBranchDetailModels, &grBranchDetailModel)
	}

	if _, err := service.GrBranchRepository.FindProductByListID(productIDs); err != nil {
		return response, err
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		// var stockUpdateEntities []*entity.StockUpdate

		if err = service.GrBranchRepository.Store(txCtx, &grBranchModel); err != nil {
			return err
		}
		response.GrBranchNo = grBranchModel.GrBranchNo

		for _, grBranchDetailModel := range grBranchDetailModels {
			var grBranchDetail model.GrBranchDetailCreate
			if err := structs.Automapper(*grBranchDetailModel, &grBranchDetail); err != nil {
				return err
			}
			// grBranchDetail := *grBranchDetailModel

			grBranchDetail.GrBranchNo = grBranchModel.GrBranchNo
			grDet, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grBranchDetail)
			if err != nil {
				return err
			}

			fmt.Println("Gr Branch Detail ID grBranchDetailModel : ", grBranchDetailModel.GrBranchDetId)
			fmt.Println("Item Type grBranchDetailModel : ", grBranchDetailModel.ItemType)
			fmt.Println("Gr Branch Detail ID grBranchDetail : ", grDet.GrBranchDetId)
			fmt.Println("Item Type grBranchDetail : ", grDet.ItemType)
			// stockUpdateEntity := entity.StockUpdate{
			// 	CustID:    grBranchModel.CustID,
			// 	WhID:      *grBranchModel.WhID,
			// 	ProID:     grDet.ProID,
			// 	StockDate: *grBranchModel.GrBranchDate,
			// 	TrCode:    grBranchModel.GrBranchNo[0:3],
			// 	TrNo:      grBranchModel.GrBranchNo,
			// 	QtyIn:     float64(grDet.QtyReceived),
			// 	UnitPrice: grDet.UnitPrice1,
			// 	RefDetId:  grDet.GrBranchDetId,
			// }

			// stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
		}

		// if *grBranchModel.DataStatus == entity.GR_BRANCH_PROCESSED {
		// 	var grBranchDetails []model.GrBranchDetailList
		// 	var stockUpdateEntities []*entity.StockUpdate

		// 	grBranchDetails, err = service.GrBranchRepository.FindGrBranchdetailWithDiscount(grBranchModel.GrBranchNo, grBranchModel.CustID)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	for _, grBranchDetail := range grBranchDetails {
		// 		ReceivedQtyUnit := &conversion.QtyUnit{
		// 			Qty1:      int(*grBranchDetail.QtyReceived1),
		// 			Qty2:      int(*grBranchDetail.QtyReceived2),
		// 			Qty3:      int(*grBranchDetail.QtyReceived3),
		// 			ConvUnit2: int(grBranchDetail.ConvUnit2),
		// 			ConvUnit3: int(grBranchDetail.ConvUnit3),
		// 		}

		// 		totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
		// 		if err != nil {
		// 			return err
		// 		}

		// 		stockUpdateEntity := entity.StockUpdate{
		// 			CustID:    grBranchModel.CustID,
		// 			WhID:      *grBranchModel.WhID,
		// 			ProID:     grBranchDetail.ProID,
		// 			StockDate: *grBranchModel.GrBranchDate,
		// 			TrCode:    grBranchModel.GrBranchNo[0:3],
		// 			TrNo:      grBranchModel.GrBranchNo,
		// 			QtyIn:     float64(totalReceivedQty),
		// 			UnitPrice: grBranchDetail.UnitPrice1,
		// 			RefDetId:  int64(grBranchDetail.GrBranchDetId),
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

	if *grBranchModel.DataStatus == entity.GR_BRANCH_PROCESSED {
		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			var obModel model.OrderBookingDetailStatus

			obModel.CustID = grBranchModel.CustID
			statusOrderBooking := int64(3)
			obModel.StatusOrderBooking = &statusOrderBooking
			obModel.UpdatedBy = grBranchModel.CreatedBy
			updatedAt := time.Now()
			obModel.UpdatedAt = &updatedAt
			if err = service.OrderBookingRepository.UpdateCompleted(txCtx, *orderBooking.PoNo, obModel); err != nil {
				return err
			}

			var grBranchDetails []model.GrBranchDetailList
			var stockUpdateEntities []*entity.StockUpdate

			grBranchDetails, err = service.GrBranchRepository.FindGrBranchdetailWithDiscount(grBranchModel.GrBranchNo, grBranchModel.CustID)
			if err != nil {
				return err
			}

			for _, grBranchDetail := range grBranchDetails {
				ReceivedQtyUnit := &conversion.QtyUnit{
					Qty1:      int(*grBranchDetail.QtyReceived1),
					Qty2:      int(*grBranchDetail.QtyReceived2),
					Qty3:      int(*grBranchDetail.QtyReceived3),
					ConvUnit2: int(grBranchDetail.ConvUnit2),
					ConvUnit3: int(grBranchDetail.ConvUnit3),
				}

				totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    grBranchModel.CustID,
					WhID:      *grBranchModel.WhID,
					ProID:     grBranchDetail.ProID,
					StockDate: *grBranchModel.GrBranchDate,
					TrCode:    grBranchModel.GrBranchNo[0:3],
					TrNo:      grBranchModel.GrBranchNo,
					QtyIn:     float64(totalReceivedQty),
					UnitPrice: grBranchDetail.UnitPrice1,
					RefDetId:  int64(grBranchDetail.GrBranchDetId),
				}

				fmt.Println("CustID : ", grBranchModel.CustID)
				fmt.Println("WhID : ", *grBranchModel.WhID)
				fmt.Println("ProID : ", grBranchDetail.ProID)
				fmt.Println("StockDate : ", *grBranchModel.GrBranchDate)
				fmt.Println("TrCode : ", grBranchModel.GrBranchNo[0:3])
				fmt.Println("TrNo : ", grBranchModel.GrBranchNo)
				fmt.Println("QtyIn : ", float64(totalReceivedQty))
				fmt.Println("UnitPrice : ", grBranchDetail.UnitPrice1)
				fmt.Println("RefDetId : ", int64(grBranchDetail.GrBranchDetId))

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

/*
	func (service *grBranchServiceImpl) StoreOld(request entity.CreateGrBranchBody) (response entity.GrBranchResponse, err error) {
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

		var grBranchModel model.GrBranch
		if err = structs.Automapper(request, &grBranchModel); err != nil {
			return response, err
		}

		grDate := str.GetJakartaDate()
		grBranchModel.GrBranchDate = &grDate

		orderBooking, err := service.GrBranchRepository.FindGrBranchOrderBooking(request.PoNo, request.CustID, request.ParentCustID)
		if err != nil {
			return response, err
		}

		dataStatusGrBranch := entity.GR_BRANCH_PROCESSED
		if *orderBooking.TypeApproval == 2 {
			dataStatusGrBranch = entity.GR_BRANCH_SUBMITTED
		}
		grBranchModel.DataStatus = &dataStatusGrBranch

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

			// err := service.GrBranchRepository.Store(txCtx, &grBranchModel)
			// if err != nil {
			// 	return err
			// }
			// response.GrBranchNo = grBranchModel.GrBranchNo
			var subTotal int64

			var productIDs []int64
			for _, detail := range request.Details.Normal {
				productIDs = append(productIDs, detail.ProID)
			}
			for _, detail := range request.Details.Promo {
				if !slices.Contains(productIDs, detail.ProID) {
					productIDs = append(productIDs, detail.ProID)
				}
			}

			productsModel, err := service.GrBranchRepository.FindProductByListID(productIDs)
			if err != nil {
				return err
			}

			var productMap = model.MapProduct{}

			for _, productModel := range productsModel {
				productMap.SetProduct(productModel.ProductId, productModel)
			}

			var stockUpdateEntities []*entity.StockUpdate

			if err = service.GrBranchRepository.Store(txCtx, &grBranchModel); err != nil {
				return err
			}
			response.GrBranchNo = grBranchModel.GrBranchNo

			for index, detail := range request.Details.Normal {
				productModel, err := productMap.GetByID(detail.ProID)
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

				ReceivedQtyUnit := &conversion.QtyUnit{
					Qty1:      detail.QtyReceived1,
					Qty2:      detail.QtyReceived2,
					Qty3:      detail.QtyReceived3,
					ConvUnit2: int(productModel.ConvUnit2),
					ConvUnit3: int(productModel.ConvUnit3),
				}

				totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				var grDetailModel model.GrBranchDetailCreate
				// seq := index + 1

				if detail.QtyShip1 == 0 && detail.QtyShip2 == 0 && detail.QtyShip3 == 0 {
					return errors.New(fmt.Sprintf("please input qty_ship on id_product %v", detail.ProID))
				}

				if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
					return errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
				}

				err = structs.Automapper(detail, &grDetailModel)
				if err != nil {
					return err
				}

				grDetailModel.CustID = request.CustID
				grDetailModel.GrBranchNo = grBranchModel.GrBranchNo
				grDetailModel.ItemType = model.ITEM_TYPE_NORMAL
				grDetailModel.SeqNo = index + 1
				grDetailModel.UnitPrice1 = productModel.PurchPrice1
				grDetailModel.UnitPrice2 = productModel.PurchPrice2
				grDetailModel.UnitPrice3 = productModel.PurchPrice3
				grDetailModel.QtyShip = totalQtyShip
				grDetailModel.QtyReceived = totalReceivedQty

				grDet, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grDetailModel)
				if err != nil {
					return err
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    grDetailModel.CustID,
					WhID:      *grBranchModel.WhID,
					ProID:     detail.ProID,
					StockDate: *grBranchModel.GrBranchDate,
					TrCode:    grBranchModel.GrBranchNo[0:3],
					TrNo:      grBranchModel.GrBranchNo,
					QtyIn:     float64(grDet.Qty),
					UnitPrice: grDetailModel.UnitPrice1,
					// RefDetId:  grDet.GrBranchDetId,
				}

				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)

			}

			for index, detail := range request.Details.Promo {
				productModel, err := productMap.GetByID(detail.ProID)
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

				ReceivedQtyUnit := &conversion.QtyUnit{
					Qty1:      detail.QtyReceived1,
					Qty2:      detail.QtyReceived2,
					Qty3:      detail.QtyReceived3,
					ConvUnit2: int(productModel.ConvUnit2),
					ConvUnit3: int(productModel.ConvUnit3),
				}

				totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				var grDetailModel model.GrBranchDetailCreate

				if detail.QtyReceived1 == 0 && detail.QtyReceived2 == 0 && detail.QtyReceived3 == 0 {
					return errors.New(fmt.Sprintf("please input qty_received on id_product %v", detail.ProID))
				}

				err = structs.Automapper(detail, &grDetailModel)
				if err != nil {
					return err
				}
				// grDetailModel.SeqNo = seq
				grDetailModel.CustID = request.CustID
				grDetailModel.GrBranchNo = grBranchModel.GrBranchNo
				grDetailModel.ItemType = model.ITEM_TYPE_PROMO
				grDetailModel.SeqNo = index + 1
				grDetailModel.UnitPrice1 = productModel.PurchPrice1
				grDetailModel.UnitPrice2 = productModel.PurchPrice2
				grDetailModel.UnitPrice3 = productModel.PurchPrice3
				grDetailModel.QtyShip = totalQtyShip
				grDetailModel.QtyReceived = totalReceivedQty
				grDet, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grDetailModel)
				if err != nil {
					return err
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    grDetailModel.CustID,
					WhID:      *grBranchModel.WhID,
					ProID:     detail.ProID,
					StockDate: *grBranchModel.GrBranchDate,
					TrCode:    grBranchModel.GrBranchNo[0:3],
					TrNo:      grBranchModel.GrBranchNo,
					QtyIn:     float64(grDet.Qty),
					UnitPrice: grDetailModel.UnitPrice1,
					// RefDetId:  grDet.GrBranchDetId,
				}

				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)

			}

			// err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			// if err != nil {
			// 	return err
			// }

			return nil
		})

		return response, err
	}
*/
func (service *grBranchServiceImpl) Detail(grBranchNo string, custID, parentCustId string, queryParam entity.GrBranchDetailQuery) (response entity.GrBranchWithDetailResponse, err error) {
	gr, err := service.GrBranchRepository.FindByNo(grBranchNo, queryParam.GrBranchCustId, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gr, &response)
	if err != nil {
		return response, err
	}

	err = service.GetGrBranchDetail(grBranchNo, queryParam.GrBranchCustId, &response)
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

func (service *grBranchServiceImpl) DetailByInvoice(invoice string, custID, parentCustId string, isAp bool) (response entity.GrBranchWithDetailResponse, err error) {
	gr, err := service.GrBranchRepository.GetByInvoiceNo(invoice, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(gr, &response)
	if err != nil {
		return response, err
	}

	err = service.GetGrBranchDetail(invoice, custID, &response)
	if err != nil {
		return response, err
	}

	grData := gr.GrBranchDate.Format("2006-01-02")
	// deliveryDate := gr.DeliveryDate.Format("2006-01-02")
	// invoiceDate := gr.InvoiceDate.Format("2006-01-02")

	response.GrBranchDate = grData
	// response.DeliveryDate = deliveryDate
	// response.InvoiceDate = invoiceDate
	return response, nil
}

func (service *grBranchServiceImpl) GetGrBranchDetail(grBranchNo string, custID string, grBranch *entity.GrBranchWithDetailResponse) (err error) {
	var grBranchDetails []model.GrBranchDetailList
	// var discountValueTotal float64

	grBranchDetails, err = service.GrBranchRepository.FindGrBranchdetailWithDiscount(grBranchNo, custID)
	if err != nil {
		return err
	}

	grBranch.Details.Promo = []entity.GrBranchDetailList{}
	grBranch.Details.Normal = []entity.GrBranchDetailList{}

	for _, grBranchDetail := range grBranchDetails {
		var grBranchDetailData entity.GrBranchDetailList
		err = structs.Automapper(grBranchDetail, &grBranchDetailData)
		if err != nil {
			return err
		}

		qtyShip := &conversion.Qty{
			Qty:       int(grBranchDetailData.QtyShip),
			ConvUnit2: int(grBranchDetailData.ConvUnit2),
			ConvUnit3: int(grBranchDetailData.ConvUnit3),
		}

		qtyShipConversion := qtyShip.ConvToQtyConversion()

		qtyReceived := &conversion.Qty{
			Qty:       int(grBranchDetailData.QtyReceived),
			ConvUnit2: int(grBranchDetailData.ConvUnit2),
			ConvUnit3: int(grBranchDetailData.ConvUnit3),
		}
		qtyReceivedConversion := qtyReceived.ConvToQtyConversion()

		grBranchDetailData.QtyShip1 = qtyShipConversion.Qty1
		grBranchDetailData.QtyShip2 = qtyShipConversion.Qty2
		grBranchDetailData.QtyShip3 = qtyShipConversion.Qty3

		grBranchDetailData.QtyReceived1 = qtyReceivedConversion.Qty1
		grBranchDetailData.QtyReceived2 = qtyReceivedConversion.Qty2
		grBranchDetailData.QtyReceived3 = qtyReceivedConversion.Qty3

		if grBranchDetail.ItemType == 1 {
			// var discountValue, discount float64
			// Subtotal := (grBranchDetail.UnitPrice1 * float64(grBranchDetailData.QtyReceived1)) + (grBranchDetail.UnitPrice2 * float64(grBranchDetailData.QtyReceived2)) + (grBranchDetail.UnitPrice3 * float64(grBranchDetailData.QtyReceived3))
			// if isAp {
			// 	if grBranchDetail.Discount != nil {
			// 		discountValue = (*grBranchDetail.Discount / 100) * Subtotal
			// 		discount = *grBranchDetail.Discount
			// 	}
			// 	grBranchDetailData.DiscountValue = &discountValue
			// 	grBranchDetailData.Discount = &discount

			// 	discountValueTotal += discountValue
			// }

			// ppn := (Subtotal - discountValue) * grBranchDetail.Vat / 100
			// pbn := (Subtotal - discountValue) * grBranchDetail.VatLgPurch / 100
			//ppnDp := (Subtotal - discountValue) * grBranchDetail.VatBg / 100
			// total := (Subtotal - discountValue) + ppn + pbn

			// grBranch.SubTotal += Subtotal
			// grBranch.SubTotal += Subtotal - discountValue
			// grBranch.Total += total
			// grBranch.TotalVat += ppn
			// grBranch.TotalVatLgPurch += pbn
			// grBranch.TotalSkuPrice += Subtotal - discountValue
			// grBranchDetailData.Nett = Subtotal - discountValue
			// grBranchDetailData.SubTotal = Subtotal
			// grBranchDetailData.Total = total
			// grBranchDetailData.VatValue = ppn
			// grBranchDetailData.VatLgPurchValue = pbn

			// grBranchDetailData.ConvUnit1 = grBranchDetail.ConvUnit2 * grBranchDetail.ConvUnit3
			grBranch.Details.Normal = append(grBranch.Details.Normal, grBranchDetailData)
		} else {
			grBranch.Details.Promo = append(grBranch.Details.Promo, grBranchDetailData)
		}
	}

	if grBranch.Details.Promo == nil {
		grBranch.Details.Promo = []entity.GrBranchDetailList{}
	}

	// if isAp {
	// 	grBranch.DiscountValue = &discountValueTotal
	// }
	return
}

func (service *grBranchServiceImpl) List(dataFilter entity.GrBranchQueryFilter, custId, parentCustId string) (data []entity.GrBranchListResponse, total int64, lastPage int, err error) {
	grs, total, lastPage, err := service.GrBranchRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grs {
		var vResp entity.GrBranchListResponse
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
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *grBranchServiceImpl) ListSupplier(dataFilter entity.GrBranchSupplierQueryFilter, custId, parentCustId string) (data []entity.GrBranchSupplierListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.GrBranchRepository.FindSupplierGrBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.GrBranchSupplierListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *grBranchServiceImpl) ListDistributor(dataFilter entity.GrBranchDistributorQueryFilter, custId, parentCustId string) (data []entity.GrBranchDistributorListResponse, total int64, lastPage int, err error) {
	grDistributor, total, lastPage, err := service.GrBranchRepository.FindDistributorGrBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grDistributor {
		var vResp entity.GrBranchDistributorListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *grBranchServiceImpl) ListWarehouse(dataFilter entity.GrBranchWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrBranchWarehouseListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.GrBranchRepository.FindWarehouseGrBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.GrBranchWarehouseListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *grBranchServiceImpl) Update(grNo string, request entity.UpdateGrBranchRequest) (err error) {
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

	var grModel model.GrBranch
	err = structs.Automapper(request, &grModel)
	if err != nil {
		return err
	}
	grModel.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		gr, err := service.GrBranchRepository.FindByNo(grNo, request.CustID, request.ParentCustID)
		if err != nil {
			return err
		}

		err = service.GrBranchRepository.Update(txCtx, gr.GrBranchNo, &grModel)
		if err != nil {
			return err
		}

		// grDetails, err := service.GrBranchRepository.FindGrBranchdetail(grNo, request.CustID)
		// if err != nil {
		// 	return err
		// }

		// for _, grDet := range grDetails {
		// 	err = service.WhStockRepository.UpdateOldWhStock(txCtx, request.CustID, *request.WhID, grDet.ProID, grDet.Qty)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		err = service.GrBranchRepository.DeleteGrBranchDetailByGrBranchNo(txCtx, grNo, gr.CustID)
		if err != nil {
			return err
		}

		newGrBranchDetIds := make([]int64, 0)
		for _, detail := range request.Details.Normal {
			/* if detail.GrBranchDetId == nil || *detail.GrBranchDetId == 0 { */
			detail.GrBranchDetId = nil
			var grDetailModel model.GrBranchDetailCreate
			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}
			// grDetailModel.SeqNo = sequence
			grDetailModel.CustID = request.CustID
			grDetailModel.GrBranchNo = grNo
			// grDet, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grDetailModel)
			_, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			// newGrBranchDetIds = append(newGrBranchDetIds, grDet.GrBranchDetId)
			newGrBranchDetIds = append(newGrBranchDetIds, grDetailModel.ProID)
		}

		for _, detail := range request.Details.Promo {
			detail.GrBranchDetId = nil
			var grDetailModel model.GrBranchDetailCreate
			err = structs.Automapper(detail, &grDetailModel)
			if err != nil {
				return err
			}
			grDetailModel.CustID = request.CustID
			grDetailModel.GrBranchNo = grNo
			_, err := service.GrBranchRepository.CreateGrBranchDetail(txCtx, &grDetailModel)
			if err != nil {
				return err
			}

			// newGrBranchDetIds = append(newGrBranchDetIds, grDet.GrBranchDetId)
			newGrBranchDetIds = append(newGrBranchDetIds, grDetailModel.ProID)
		}

		err = service.GrBranchRepository.DeleteStockNotInRefIds(txCtx, request.CustID, gr.GrBranchNo, newGrBranchDetIds)
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

func (service *grBranchServiceImpl) Delete(custId string, grNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		grDetails, err := service.GrBranchRepository.FindGrBranchdetailForUpdateWhStock(grNo, custId)
		if err != nil {
			return err
		}
		log.Println("grDetails:", structs.StructToJson(grDetails))

		oldGrBranchDetId := make([]int64, 0)
		for _, grDet := range grDetails {
			err = service.GrBranchRepository.UpdateOldWhStock(txCtx, custId, grDet.WhID, grDet.ProID, *grDet.Qty)
			if err != nil {
				return err
			}
			oldGrBranchDetId = append(oldGrBranchDetId, grDet.GrBranchDetId)
		}

		log.Println("oldGrBranchDetId:", structs.StructToJson(oldGrBranchDetId))
		log.Println("oldGrBranchDetId:", structs.StructToJson(oldGrBranchDetId))
		if len(oldGrBranchDetId) > 0 {
			err = service.GrBranchRepository.DeleteStockInRefIds(txCtx, custId, grNo, oldGrBranchDetId)
			if err != nil {
				return err
			}
		}

		err = service.GrBranchRepository.Delete(txCtx, custId, grNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *grBranchServiceImpl) OrderBookingDetail(orderBookingId int, custID string, parentCustId string) (responses []entity.GrBranchOrderBookingDetailResponse, err error) {
	orderBookingDetails, err := service.GrBranchRepository.FindGrBranchOrderBookingDetails(orderBookingId, custID, parentCustId)
	if err != nil {
		return responses, err
	}

	for _, orderBookingDetail := range orderBookingDetails {
		var response entity.GrBranchOrderBookingDetailResponse
		err = structs.Automapper(orderBookingDetail, &response)
		if err != nil {
			return responses, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (service *grBranchServiceImpl) OrderBookingList(dataFilter entity.GrBranchOrderBookingListQueryFilter, custId string, parentCustId string) (data []entity.GrBranchOrderBookingListResponse, total int64, lastPage int, err error) {
	orderBookingList, total, lastPage, err := service.GrBranchRepository.FindGrBranchOrderBookingList(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range orderBookingList {
		var vResp entity.GrBranchOrderBookingListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}

func (service *grBranchServiceImpl) BulkUpdateStatus(request entity.GrBranchBulkUpdateDataStatus, custId string, parentCustId string) (err error) {
	c := context.Background()

	for index := range request.GrBranches {
		// End parse time format YYYY-mm-dd to Rfc339
		var Model model.GrBranch
		err = structs.Automapper(request.GrBranches[index], &Model)
		if err != nil {
			return err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

			grBranch, err := service.GrBranchRepository.FindByNo(request.GrBranches[index].GrBranchNo, request.GrBranches[index].CustId, parentCustId)
			if err != nil {
				return err
			}

			Model.SupID = grBranch.SupId
			err = service.GrBranchRepository.Update(txCtx, request.GrBranches[index].GrBranchNo, &Model)
			if err != nil {
				return err
			}

			if *request.GrBranches[index].DataStatus == entity.GR_BRANCH_PROCESSED { // jika processed, trigger stock
				var obModel model.OrderBookingDetailStatus

				obModel.CustID = grBranch.CustID
				statusOrderBooking := int64(3)
				obModel.StatusOrderBooking = &statusOrderBooking
				updatedBy := request.GrBranches[index].UpdatedBy
				obModel.UpdatedBy = &updatedBy
				updatedAt := time.Now()
				obModel.UpdatedAt = &updatedAt
				if err = service.OrderBookingRepository.UpdateCompleted(txCtx, *grBranch.PoNo, obModel); err != nil {
					return err
				}
				// details, err := service.GrBranchRepository.FindGrBranchdetailWithDiscount(request.GrBranches[index].GrBranchNo, custId)
				// if err != nil {
				// 	return err
				// }

				// var salesDetailCanceledUpdates []*entity.SalesOrderStockUpdate

				// for _, detail := range details {
				// 	salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
				// 		CustID:         custId,
				// 		WhID:           *grBranch.WhId,
				// 		ProID:          int64(detail.ProId),
				// 		StockDate:      *grBranch.RoDate,
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

				var grBranchDetails []model.GrBranchDetailList
				var stockUpdateEntities []*entity.StockUpdate

				grBranchDetails, err = service.GrBranchRepository.FindGrBranchdetailWithDiscount(request.GrBranches[index].GrBranchNo, request.GrBranches[index].CustId)
				if err != nil {
					return err
				}

				for _, grBranchDetail := range grBranchDetails {
					ReceivedQtyUnit := &conversion.QtyUnit{
						Qty1:      int(*grBranchDetail.QtyReceived1),
						Qty2:      int(*grBranchDetail.QtyReceived2),
						Qty3:      int(*grBranchDetail.QtyReceived3),
						ConvUnit2: int(grBranchDetail.ConvUnit2),
						ConvUnit3: int(grBranchDetail.ConvUnit3),
					}

					totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
					if err != nil {
						return err
					}

					stockUpdateEntity := entity.StockUpdate{
						CustID:    custId,
						WhID:      *grBranch.WhId,
						ProID:     grBranchDetail.ProID,
						StockDate: *grBranch.GrBranchDate,
						TrCode:    grBranch.GrBranchNo[0:3],
						TrNo:      grBranch.GrBranchNo,
						QtyIn:     float64(totalReceivedQty),
						UnitPrice: grBranchDetail.UnitPrice1,
						RefDetId:  int64(grBranchDetail.GrBranchDetId),
					}

					// fmt.Println("CustID : ", grBranchModel.CustID)
					// fmt.Println("WhID : ", *grBranchModel.WhID)
					// fmt.Println("ProID : ", grBranchDetail.ProID)
					// fmt.Println("StockDate : ", *grBranchModel.GrBranchDate)
					// fmt.Println("TrCode : ", grBranchModel.GrBranchNo[0:3])
					// fmt.Println("TrNo : ", grBranchModel.GrBranchNo)
					// fmt.Println("QtyIn : ", float64(totalReceivedQty))
					// fmt.Println("UnitPrice : ", grBranchDetail.UnitPrice1)
					// fmt.Println("RefDetId : ", int64(grBranchDetail.GrBranchDetId))

					stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
				}

				err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
				if err != nil {
					return err
				}
			}

			/*
				if *request.GrBranches[index].DataStatus == entity.GR_BRANCH_REJECTED { // jika cancel, trigger stock
					grBranch, err := service.GrBranchRepository.FindByNo(request.GrBranches[index].GrBranchNo, custId, parentCustId)
					if err != nil {
						return err
					}

					details, err := service.GrBranchRepository.FindGrBranchdetailWithDiscount(request.GrBranches[index].GrBranchNo, custId)
					if err != nil {
						return err
					}


					var salesDetailCanceledUpdates []*entity.SalesOrderStockUpdate

					for _, detail := range details {
						salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
							CustID:         custId,
							WhID:           *ro.WhId,
							ProID:          int64(detail.ProId),
							StockDate:      *ro.RoDate,
							TrCode:         request.Orders[index].RoNo[0:2],
							TrNo:           request.Orders[index].RoNo,
							QtyOrderBefore: detail.Qty,
							QtyOrder:       0,
							UnitPrice:      *detail.SellPrice1,
							RefDetId:       int64(*detail.OrderDetailID),
						}
						salesDetailCanceledUpdates = append(salesDetailCanceledUpdates, &salesDetailCanceledUpdate)

					}

					if len(salesDetailCanceledUpdates) > 0 {
						log.Info("Update Stock Deleted Details")
						err = service.StockRepository.SalesStockUpdates(txCtx, salesDetailCanceledUpdates)
						if err != nil {
							return err
						}
					}
				}
			*/
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *grBranchServiceImpl) BulkPrint(request entity.GrBranchBulkPrint, custId string, parentCustId string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		fmt.Println("GrBranchService Print")
		err = service.GrBranchRepository.PrintGrBranch(txCtx, custId, request.GrBranches, userId)
		if err != nil {
			return err
		}

		for index := range request.GrBranches {
			grBranch, err := service.GrBranchRepository.FindByNo(request.GrBranches[index].GrBranchNo, request.GrBranches[index].CustId, parentCustId)
			if err != nil {
				return err
			}

			var grBranchDetails []model.GrBranchDetailList
			var stockUpdateEntities []*entity.StockUpdate

			grBranchDetails, err = service.GrBranchRepository.FindGrBranchdetailWithDiscount(request.GrBranches[index].GrBranchNo, request.GrBranches[index].CustId)
			if err != nil {
				return err
			}

			for _, grBranchDetail := range grBranchDetails {
				ReceivedQtyUnit := &conversion.QtyUnit{
					Qty1:      int(*grBranchDetail.QtyReceived1),
					Qty2:      int(*grBranchDetail.QtyReceived2),
					Qty3:      int(*grBranchDetail.QtyReceived3),
					ConvUnit2: int(grBranchDetail.ConvUnit2),
					ConvUnit3: int(grBranchDetail.ConvUnit3),
				}

				totalReceivedQty, err := ReceivedQtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				WhIdPrincipal := request.GrBranches[index].WhId
				if *grBranch.TypeApproval == entity.ORDER_BOOKING_TYPE_APPROVAL_EKSTERNAL {
					WhIdPrincipal = *grBranch.WhId
				}

				stockUpdateEntity := entity.StockUpdate{
					CustID:    custId,
					WhID:      WhIdPrincipal,
					ProID:     grBranchDetail.ProID,
					StockDate: *grBranch.GrBranchDate,
					TrCode:    grBranch.GrBranchNo[0:3],
					TrNo:      grBranch.GrBranchNo,
					QtyOut:    float64(totalReceivedQty),
					UnitPrice: grBranchDetail.UnitPrice1,
					RefDetId:  int64(grBranchDetail.GrBranchDetId),
				}
				stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)

				if *grBranch.TypeApproval == entity.ORDER_BOOKING_TYPE_APPROVAL_EKSTERNAL {
					stockUpdateEntity := entity.StockUpdate{
						CustID:    grBranch.CustID,
						WhID:      request.GrBranches[index].WhId,
						ProID:     grBranchDetail.ProID,
						StockDate: *grBranch.GrBranchDate,
						TrCode:    grBranch.GrBranchNo[0:3],
						TrNo:      grBranch.GrBranchNo,
						QtyIn:     float64(totalReceivedQty),
						UnitPrice: grBranchDetail.UnitPrice1,
						RefDetId:  int64(grBranchDetail.GrBranchDetId),
					}
					stockUpdateEntities = append(stockUpdateEntities, &stockUpdateEntity)
				}
			}

			err = service.StockRepository.StockUpdates(txCtx, stockUpdateEntities)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (service *grBranchServiceImpl) ListPrintWarehouse(dataFilter entity.GrBranchPrintWarehouseQueryFilter, custId, parentCustId string) (data []entity.GrBranchWarehouseListResponse, total int64, lastPage int, err error) {
	grsupplier, total, lastPage, err := service.GrBranchRepository.FindPrintWarehouseGrBranch(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range grsupplier {
		var vResp entity.GrBranchWarehouseListResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}
	return data, total, lastPage, err
}
