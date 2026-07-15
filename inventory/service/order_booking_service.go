package service

import (
	"context"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/conversion"
	"inventory/pkg/structs"
	"inventory/repository"
	"time"
)

type OrderService interface {
	Store(request entity.CreateOrderBody) (data entity.CreateOrderBookingResponse, err error)
	Detail(OrderBookingId int, custID string, parentCustId string) (response entity.OrderBookingResponse, err error)
	List(dataFilter entity.OrderQueryFilter) (data []entity.OrderBookingListResponse, total int64, lastPage int, err error)
	Delete(custId string, OrderBookingId int, userId int64) (err error)
	LookupStatus() (data []entity.StatusList, err error)
	UpdateReject(OrderBookingId int, custID string, parentCustId string) (err error)
	UpdateAppove(OrderBookingId int, request entity.OrderBookingResponse, parentCustId string) (err error)
}

func NewOrderBookingService(orderRepository repository.OrderBookingRepository, transaction repository.Dbtransaction) *orderServiceImpl {
	return &orderServiceImpl{
		OrderRepository: orderRepository,
		Transaction:     transaction,
	}
}

type orderServiceImpl struct {
	OrderRepository repository.OrderBookingRepository
	Transaction     repository.Dbtransaction
}

// Tambahkan fungsi ini di tempat yang sesuai
func getValueOrDefault(value *float64, defaultValue float64) float64 {
	if value == nil {
		return defaultValue
	}
	return *value
}

func (service *orderServiceImpl) LookupStatus() (data []entity.StatusList, err error) {
	var statusList = entity.DataStatus

	for id, name := range statusList {
		// fmt.Println(id, name)
		data = append(data, entity.StatusList{
			StatusOrderBooking:     int(id),
			StatusOrderBookingName: name,
		})
	}

	return data, err
}

func (service *orderServiceImpl) List(dataFilter entity.OrderQueryFilter) (data []entity.OrderBookingListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.OrderBookingListResponse
		structs.Automapper(row, &vResp)
		if !row.CreatedAt.IsZero() {
			roDate := row.CreatedAt.Format("2006-01-02")
			vResp.CreatedAt = roDate
		}

		statusName := vResp.GenerateDataStatusName()

		vResp.StatusName = statusName

		vResp.TotalTotal = vResp.Total + getValueOrDefault(vResp.TotalAlloc, 0)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderServiceImpl) Detail(OrderBookingId int, custID string, parentCustId string) (response entity.OrderBookingResponse, err error) {
	ro, err := service.OrderRepository.FindByNo(OrderBookingId, custID, parentCustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.OrderRepository.FindDetail(OrderBookingId, custID)
	fmt.Println(details)
	if err != nil {
		return response, err
	}

	defFloat := float64(0)

	for _, detail := range details {
		var detailData entity.OrderBookingDetailResponse
		var detailFinalData entity.OrderBookingDetailFinalResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		err = structs.Automapper(detail, &detailFinalData)
		if err != nil {
			return response, err
		}

		qty := &conversion.Qty{
			Qty:       int(getValueOrDefault(detailData.Qty, 0)),
			ConvUnit2: int(*detail.MpConvUnit2),
			ConvUnit3: int(*detail.MpConvUnit3),
		}

		qtyConversion := qty.ConvToQtyConversion()
		detailDataQty1 := float64(qtyConversion.Qty1)
		detailDataQty2 := float64(qtyConversion.Qty2)
		detailDataQty3 := float64(qtyConversion.Qty3)

		detailData.Qty1 = &detailDataQty1 // Tambahkan pointer
		detailData.Qty2 = &detailDataQty2 // Tambahkan pointer
		detailData.Qty3 = &detailDataQty3 // Tambahkan pointer

		detailData.Qty1Total = *detailData.Qty1 + getValueOrDefault(detailData.Qty1Alloc, 0)
		detailData.Qty2Total = *detailData.Qty2 + getValueOrDefault(detailData.Qty2Alloc, 0)
		detailData.Qty3Total = *detailData.Qty3 + getValueOrDefault(detailData.Qty3Alloc, 0)
		// detailData.Qty4Total = *detailData.Qty4 + getValueOrDefault(detailData.Qty4Alloc, 0)
		// detailData.Qty5Total = *detailData.Qty5 + getValueOrDefault(detailData.Qty5Alloc, 0)

		response.Details = append(response.Details, detailData)

		if (detailData.Qty1Total + detailData.Qty2Total + detailData.Qty3Total) > 0 {
			response.OrderBookingApproval = append(response.OrderBookingApproval, detailData)
		}

		if ro.GrBranchNo == nil {
			response.OrderBookingFinal = append(response.OrderBookingFinal, detailFinalData)
		}
	}

	if !ro.CreatedAt.IsZero() {
		roDate := ro.CreatedAt.Format("2006-01-02")
		response.CreatedAt = roDate
	}

	statusName := response.GenerateDataStatusName()
	response.StatusName = statusName

	response.SubTotalTotal = response.SubTotal + getValueOrDefault(response.SubTotalAlloc, 0)
	response.VatValueTotal = response.VatValue + getValueOrDefault(response.VatValueAlloc, 0)
	response.TotalTotal = response.Total + getValueOrDefault(response.TotalAlloc, 0)

	response.OrderBooking = response.Details

	if ro.GrBranchNo != nil {
		response.TotalTotalFinal = &defFloat
		detailsFinal, err := service.OrderRepository.FindDetailFinal(*ro.GrBranchNo, custID, OrderBookingId)
		fmt.Println(detailsFinal)
		if err != nil {
			return response, err
		}
		var detailFinalData []entity.OrderBookingDetailFinalResponse
		err = structs.Automapper(detailsFinal, &detailFinalData)
		if err != nil {
			return response, err
		}
		response.SubTotalFinal = ro.SubTotalFinal
		response.VatValueFinal = ro.VatValueFinal
		response.DeliveryFeeFinal = ro.DeliveryFeeFinal
		response.TotalTotalFinal = ro.TotalFinal

		response.OrderBookingFinal = detailFinalData
	}

	return response, nil
}

func (service *orderServiceImpl) Store(request entity.CreateOrderBody) (data entity.CreateOrderBookingResponse, err error) {
	c := context.Background()

	count, err := service.OrderRepository.CountAllByCustId(request.CustId)

	obDate := time.Now()

	// Format tanggal menjadi yy, mm, dan dd
	yy := obDate.Format("2006") // 2 digit tahun
	mm := obDate.Format("01")   // 2 digit bulan
	dd := obDate.Format("02")   // 2 digit hari
	// Format urutan menjadi 4 digit
	seqFormatted := fmt.Sprintf("%03d", count+1)

	// Gabungkan semuanya untuk membuat nomor faktur
	obNumber := fmt.Sprintf("OB%s%s%s%s", dd, mm, yy, seqFormatted)

	var orderModel model.OrderBooking
	err = structs.Automapper(request, &orderModel)
	if err != nil {
		return data, err
	}
	orderModel.OrderBookingId = nil // pastikan ID kosong
	orderModel.StatusOrderBooking = 1
	orderModel.PoNo = &obNumber

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.OrderRepository.Store(txCtx, &orderModel)
		if err != nil {
			return err
		}

		for _, detail := range request.Details {
			var gdsDetModel model.OrderBookingDetail

			gdsDetModel.CustID = request.CustId

			gdsDetModel.ItemType = 1

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*detail.Qty1),
				Qty2:      int(*detail.Qty2),
				Qty3:      int(*detail.Qty3),
				ConvUnit2: int(*detail.ConvUnit2),
				ConvUnit3: int(*detail.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			gdsDetModel.QtyBo = float64(totalQty)
			// gdsDetModel.QtyAlloc = float64(totalQty)
			gdsDetModel.ConvUnit2 = detail.ConvUnit2
			gdsDetModel.ConvUnit3 = detail.ConvUnit3

			err = structs.Automapper(detail, &gdsDetModel)
			if err != nil {
				return err
			}

			gdsDetModel.OrderBookingId = *orderModel.OrderBookingId

			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return data, err
	}

	data.OrderBookingId = *orderModel.OrderBookingId
	return data, nil
}

func (service *orderServiceImpl) Delete(custId string, OrderBookingId int, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderRepository.Delete(txCtx, custId, OrderBookingId, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *orderServiceImpl) UpdateReject(OrderBookingId int, custID string, parentCustId string) (err error) {
	c := context.Background()

	_, err = service.OrderRepository.FindByNo(OrderBookingId, custID, parentCustId)
	if err != nil {
		return err
	}

	var orderModel model.OrderBookingDetailStatus

	status := int64(0)
	orderModel.StatusOrderBooking = &status
	orderModel.CustID = ""

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderRepository.UpdateApproval(txCtx, OrderBookingId, orderModel)
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

func derefInt(ptr *float64) int {
	if ptr != nil {
		return int(*ptr)
	}
	return 0
}

func (service *orderServiceImpl) UpdateAppove(OrderBookingId int, request entity.OrderBookingResponse, parentCustId string) (err error) {
	c := context.Background()

	_, err = service.OrderRepository.FindByNo(OrderBookingId, request.CustId, parentCustId)
	if err != nil {
		return err
	}

	status := int64(2)
	var orderModel model.OrderBookingDetailStatusApproval

	orderModel.StatusOrderBooking = status
	orderModel.CustID = ""
	orderModel.SubTotalAlloc = *request.SubTotalAlloc
	orderModel.VatValueAlloc = *request.VatValueAlloc
	orderModel.TotalAlloc = *request.TotalAlloc
	orderModel.DeliveryFee = *request.DeliveryFee
	orderModel.TypeApproval = int64(*request.TypeApproval)

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {

		DetailIds := []int64{}

		for _, detail := range request.Details {
			if detail.OrderBookingDetailId != nil {
				DetailIds = append(DetailIds, int64(*detail.OrderBookingDetailId))
			}
		}

		deletedDetails, err := service.OrderRepository.FindDetailByNotInDetailIDs(DetailIds, OrderBookingId, request.CustId)
		if err != nil {
			return err
		}

		for _, deletedDetails := range deletedDetails {

			qtyAllocFloat := -float64(deletedDetails.QtyBo)
			qty1Float := -float64(*deletedDetails.Qty1)
			qty2Float := -float64(*deletedDetails.Qty2)
			qty3Float := -float64(*deletedDetails.Qty3)
			amountAllocFloat := float64(0)
			vatValueAllocFloat := float64(0)
			salesDetailDeletedUpdate := model.OrderBookingDetailApproval{
				CustID:        "",
				QtyAlloc:      &qtyAllocFloat,
				Qty1Alloc:     &qty1Float,
				Qty2Alloc:     &qty2Float,
				Qty3Alloc:     &qty3Float,
				AmountAlloc:   &amountAllocFloat,
				VatValueAlloc: &vatValueAllocFloat,
			}

			err = service.OrderRepository.UpdateOrderDetail(txCtx, int(*deletedDetails.OrderBookingDetailId), salesDetailDeletedUpdate)
			if err != nil {
				return err
			}

		}

		for _, detail := range request.Details {

			if detail.OrderBookingDetailId != nil {

				QtyUnit := &conversion.QtyUnit{
					Qty1:      int(derefInt(detail.Qty1Alloc)),
					Qty2:      int(derefInt(detail.Qty2Alloc)),
					Qty3:      int(derefInt(detail.Qty3Alloc)),
					ConvUnit2: int(*detail.ConvUnit2),
					ConvUnit3: int(*detail.ConvUnit3),
				}

				totalQty, err := QtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				totalQtyFloat := float64(totalQty)

				salesDetailUpdate := model.OrderBookingDetailApproval{
					CustID:        "",
					QtyAlloc:      &totalQtyFloat,
					Qty1Alloc:     detail.Qty1Alloc,
					Qty2Alloc:     detail.Qty2Alloc,
					Qty3Alloc:     detail.Qty3Alloc,
					AmountAlloc:   detail.AmountAlloc,
					VatValueAlloc: detail.VatValueAlloc,
				}

				err = service.OrderRepository.UpdateOrderDetail(txCtx, *detail.OrderBookingDetailId, salesDetailUpdate)
				if err != nil {
					return err
				}

			} else {
				var gdsDetModel model.OrderBookingDetail

				gdsDetModel.CustID = request.CustId

				gdsDetModel.ItemType = 1

				QtyUnit := &conversion.QtyUnit{
					Qty1:      int(*detail.Qty1Alloc),
					Qty2:      int(*detail.Qty2Alloc),
					Qty3:      int(*detail.Qty3Alloc),
					ConvUnit2: int(*detail.ConvUnit2),
					ConvUnit3: int(*detail.ConvUnit3),
				}

				totalQty, err := QtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				err = structs.Automapper(detail, &gdsDetModel)
				if err != nil {
					return err
				}

				qty1Float := float64(0)
				qty2Float := float64(0)
				qty3Float := float64(0)
				amountFloat := float64(0)
				vatValueFloat := float64(0)

				gdsDetModel.QtyBo = float64(0)
				gdsDetModel.Qty1 = &qty1Float
				gdsDetModel.Qty2 = &qty2Float
				gdsDetModel.Qty3 = &qty3Float
				gdsDetModel.Amount = &amountFloat
				gdsDetModel.VatValue = &vatValueFloat
				gdsDetModel.QtyAlloc = float64(totalQty)

				gdsDetModel.OrderBookingId = OrderBookingId

				err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
				if err != nil {
					return err
				}

			}

			err = service.OrderRepository.UpdateOrder(txCtx, OrderBookingId, orderModel)
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
