package service

import (
	"context"
	"fmt"
	"mobile/entity"
	"mobile/model"
	"mobile/pkg/config/env"
	"mobile/pkg/conversion"
	"mobile/pkg/str"
	"mobile/pkg/structs"
	"mobile/repository"
	"slices"
	"time"
)

type OrderCanvasService interface {
	Store(request entity.CreateOrderBody) (err error)
	StoreNoOrder(request entity.CreateNoOrderBody) (err error)
	Detail(RoNo string, custID string) (response entity.OrderResponse, err error)
	List(dataFilter entity.OrderQueryFilter) (data []entity.OrderListResponse, total int64, lastPage int, err error)
	ListNoOrder(dataFilter entity.NoOrderQueryFilter) (data []entity.NoOrderListResponse, total int64, lastPage int, err error)
	Update(roNo string, request entity.UpdateOrderBody) (err error)
	UpdateFinal(roNo string, request entity.UpdateOrderDetailFinal) (err error)
	Delete(custId string, RoNo string, userId int64) (err error)
	Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error)
	BulkUpdateStatus(custId string, request entity.BulkUpdateStatusOrder) (err error)
	LookupSalesman(dataFilter entity.OrderQueryFilter) (data []entity.OutletBySalesmanResponse, total int64, lastPage int, err error)
	SummaryBySalesman(summaryFilter entity.SummaryTotalFilter) (response entity.ResponseSummaryTotal, err error)
}

func NewOrderCanvasService(
	config env.ConfigEnv,
	orderCanvasRepository repository.OrderCanvasRepository,
	discountRepository repository.DiscountRepository,
	transaction repository.Dbtransaction) *orderCanvasServiceImpl {
	return &orderCanvasServiceImpl{
		Config:                config,
		OrderCanvasRepository: orderCanvasRepository,
		DiscountRepository:    discountRepository,
		Transaction:           transaction,
	}
}

type orderCanvasServiceImpl struct {
	Config                env.ConfigEnv
	OrderCanvasRepository repository.OrderCanvasRepository
	DiscountRepository    repository.DiscountRepository
	Transaction           repository.Dbtransaction
}

// var detOrderCreatemapperToModel = func(det entity.CreateOrderDetBody, request entity.CreateOrderBody, roNo string) (model.OrderDetail, error) {
// 	var gdsDetModel model.OrderDetail

// 	// parse time format YYYY-mm-dd to Rfc3339
// 	if det.ExpDate != nil {
// 		expDate, err := str.DateStrToRfc3339String(*det.ExpDate)
// 		if err != nil {
// 			return gdsDetModel, err
// 		}
// 		det.ExpDate = &expDate
// 	}
// 	gdsDetModel.CustId = request.CustId
// 	gdsDetModel.RoNo = roNo
// 	err := structs.Automapper(det, &gdsDetModel)
// 	if err != nil {
// 		return gdsDetModel, err
// 	}
// 	return gdsDetModel, nil
// }

func generateRealOrderCanvasNumber(seq int, roDate *time.Time) string {
	// Format tanggal menjadi yy, mm, dan dd
	yy := roDate.Format("06") // 2 digit tahun
	mm := roDate.Format("01") // 2 digit bulan
	dd := roDate.Format("02") // 2 digit hari
	// Format urutan menjadi 4 digit
	seqFormatted := fmt.Sprintf("%04d", seq+1)

	// Gabungkan semuanya untuk membuat nomor faktur
	roNumber := fmt.Sprintf("SO%s%s%s%s", yy, mm, dd, seqFormatted)
	return roNumber
}

func (service *orderCanvasServiceImpl) Store(request entity.CreateOrderBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339

	if request.RoDate != nil {
		roDate, err := str.DateStrToRfc3339String(*request.RoDate)
		if err != nil {
			return err
		}
		request.RoDate = &roDate
	}

	if request.ValDate != nil {
		valDate, err := str.DateStrToRfc3339String(*request.ValDate)
		if err != nil {
			return err
		}
		request.ValDate = &valDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	if request.InvoiceDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
		if err != nil {
			return err
		}
		request.InvoiceDate = &deliveryDate
	}

	total, err := service.OrderCanvasRepository.CountAllRoByCustId(request.CustId, *request.RoDate)
	if err != nil {
		return err
	}

	var orderModel model.Order
	err = structs.Automapper(request, &orderModel)
	if err != nil {
		return err
	}

	RoNumber := generateRealOrderCanvasNumber(total, orderModel.RoDate)
	orderModel.RoNo = RoNumber
	fmt.Println("Ro Number : ", orderModel.RoNo)
	orderModel.PoNo = nil
	orderModel.InvoiceDate = nil
	orderModel.InvoiceNo = nil
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.OrderCanvasRepository.Store(txCtx, &orderModel)
		if err != nil {
			return err
		}

		var productIDs []int64
		for _, detail := range request.Details.Normal {
			productIDs = append(productIDs, int64(detail.ProId))
		}
		for _, detail := range request.Details.Promo {
			if !slices.Contains(productIDs, int64(detail.ProId)) {
				productIDs = append(productIDs, int64(detail.ProId))
			}
		}

		productsModel, err := service.OrderCanvasRepository.FindProductByListID(productIDs)
		if err != nil {
			return err
		}

		var productMap = model.MapProduct{}

		for _, productModel := range productsModel {
			productMap.SetProduct(productModel.ProductId, productModel)
		}

		for _, detail := range request.Details.Normal {
			var gdsDetModel model.OrderDetail

			// parse time format YYYY-mm-dd to Rfc3339
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}
			gdsDetModel.CustId = request.CustId
			gdsDetModel.RoNo = orderModel.RoNo
			gdsDetModel.ItemType = 1

			productModel, err := productMap.GetByID(int64(detail.ProId))
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*detail.Qty1),
				Qty2:      int(*detail.Qty2),
				Qty3:      int(*detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			gdsDetModel.Qty = float64(totalQty)
			gdsDetModel.QtyPo = float64(totalQty)
			gdsDetModel.QtyFinal = float64(totalQty)

			if detail.Qty1Stok == nil {
				detail.Qty1Stok = new(float64)
			}
			if detail.Qty2Stok == nil {
				detail.Qty2Stok = new(float64)
			}
			if detail.Qty3Stok == nil {
				detail.Qty3Stok = new(float64)
			}

			err = structs.Automapper(detail, &gdsDetModel)
			if err != nil {
				return err
			}

			err = service.OrderCanvasRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

		}
		for _, Detail := range request.Details.Promo {
			var gdsDetModel model.OrderDetail

			// parse time format YYYY-mm-dd to Rfc3339
			if Detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*Detail.ExpDate)
				if err != nil {
					return err
				}
				Detail.ExpDate = &expDate
			}
			gdsDetModel.CustId = request.CustId
			gdsDetModel.RoNo = orderModel.RoNo
			gdsDetModel.ItemType = 2

			productModel, err := productMap.GetByID(int64(Detail.ProId))
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*Detail.Qty1),
				Qty2:      int(*Detail.Qty2),
				Qty3:      int(*Detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			gdsDetModel.Qty = float64(totalQty)
			gdsDetModel.QtyPo = float64(totalQty)
			gdsDetModel.QtyFinal = float64(totalQty)

			err = structs.Automapper(Detail, &gdsDetModel)
			if err != nil {
				return err
			}
			err = service.OrderCanvasRepository.StoreDetail(txCtx, &gdsDetModel)
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

func (service *orderCanvasServiceImpl) StoreNoOrder(request entity.CreateNoOrderBody) (err error) {
	c := context.Background()

	// Parse NoOrderDate to RFC3339 format if provided
	if request.NoOrderDate != nil {
		noOrderDate, err := str.DateStrToRfc3339String(*request.NoOrderDate)
		if err != nil {
			return fmt.Errorf("invalid NoOrderDate format: %w", err)
		}
		request.NoOrderDate = &noOrderDate
	}

	// Set CreatedAt to the current time if not provided
	if request.CreatedAt == nil {
		currentTime := time.Now().Format(time.RFC3339)
		request.CreatedAt = &currentTime
	}

	// Map the request entity to the model
	var noOrderModel model.NoOrder
	if err = structs.Automapper(request, &noOrderModel); err != nil {
		return fmt.Errorf("failed to map request to model: %w", err)
	}

	// Store the NoOrder in a transaction
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if err := service.OrderCanvasRepository.StoreNoOrder(txCtx, &noOrderModel); err != nil {
			return fmt.Errorf("failed to store NoOrder: %w", err)
		}
		return nil
	})

	return err
}

func (service *orderCanvasServiceImpl) Detail(RoNo string, custID string) (response entity.OrderResponse, err error) {
	ro, err := service.OrderCanvasRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.OrderCanvasRepository.FindDetail(RoNo, custID)
	if err != nil {
		return response, err
	}
	for _, detail := range details {
		var detailData entity.OrderDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}
		if detailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, detailData)
		} else {
			response.Details.Promo = append(response.Details.Promo, detailData)
		}
	}
	if ro.RoDate != nil {
		roDate := ro.RoDate.Format("2006-01-02")
		response.RoDate = &roDate
	}
	if ro.ValDate != nil {
		ValDate := ro.ValDate.Format("2006-01-02")
		response.ValDate = &ValDate
	}
	if ro.DeliveryDate != nil {
		DelivDate := ro.DeliveryDate.Format("2006-01-02")
		response.DeliveryDate = &DelivDate
	}
	if ro.InvoiceDate != nil {
		invoiceDate := ro.InvoiceDate.Format("2006-01-02")
		response.InvoiceDate = &invoiceDate
	}

	statusName := response.GenerateDataStatusName()
	response.DataStatusName = statusName

	payTypeName := response.GeneratePayTypeName()
	response.PayTypeName = payTypeName

	return response, nil
}

func getStartAndEndOfMonthService(dateStr string) (string, string, error) {
	// Parsing string tanggal ke objek time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", "", err
	}

	// Menghitung tanggal awal bulan
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())

	// Menghitung tanggal akhir bulan
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Mengembalikan tanggal awal dan akhir bulan dalam format yyyy-mm-dd
	return startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"), nil
}

func (service *orderCanvasServiceImpl) SummaryBySalesman(summaryFilter entity.SummaryTotalFilter) (response entity.ResponseSummaryTotal, err error) {

	currentDate := time.Now().Format("2006-01-02")

	if summaryFilter.Date != nil && *summaryFilter.Date != 0 {
		currentDate = str.UnixTimestampToUtcTime(*summaryFilter.Date).Format("2006-01-02") // Mengubah time.Time menjadi string
	}

	startOfMonth, endOfMonth, err := getStartAndEndOfMonthService(currentDate)
	if err != nil {
		return response, err
	}

	var dataStatus = []int{1, 2}

	summaryToday, err := service.OrderCanvasRepository.SummaryTotalBySalesmanAndDate(summaryFilter.SalesmanId, currentDate, currentDate, dataStatus, summaryFilter.CustId)
	if err != nil {
		return response, err
	}
	dataStatus = []int{1, 2, 3, 4, 5, 6, 7}
	summaryMonth, err := service.OrderCanvasRepository.SummaryTotalBySalesmanAndDate(summaryFilter.SalesmanId, startOfMonth, endOfMonth, dataStatus, summaryFilter.CustId)
	if err != nil {
		return response, err
	}

	response.SummaryToday = summaryToday.TotalSummary
	response.SummaryMonth = summaryMonth.TotalSummary
	response.SalesmanId = summaryFilter.SalesmanId
	response.CreatedAt = currentDate

	return response, nil
}

func (service *orderCanvasServiceImpl) List(dataFilter entity.OrderQueryFilter) (data []entity.OrderListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderCanvasRepository.FindAllByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.OrderListResponse
		structs.Automapper(row, &vResp)
		if row.RoDate != nil {
			roDate := row.RoDate.Format("2006-01-02")
			vResp.RoDate = &roDate
		}
		if row.ValDate != nil {
			ValDate := row.ValDate.Format("2006-01-02")
			vResp.ValDate = &ValDate
		}
		if row.DeliveryDate != nil {
			DelivDate := row.DeliveryDate.Format("2006-01-02")
			vResp.DeliveryDate = &DelivDate
		}
		if row.InvoiceDate != nil {
			invoiceDate := row.InvoiceDate.Format("2006-01-02")
			vResp.InvoiceDate = &invoiceDate
		}

		statusName := vResp.GenerateDataStatusName()
		vResp.DataStatusName = statusName

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderCanvasServiceImpl) ListNoOrder(dataFilter entity.NoOrderQueryFilter) (data []entity.NoOrderListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderCanvasRepository.FindAllNoOrderByCustId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.NoOrderListResponse
		structs.Automapper(row, &vResp)
		if row.NoOrderDate != nil {
			roDate := row.NoOrderDate.Format("2006-01-02")
			vResp.NoOrderDate = &roDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderCanvasServiceImpl) Update(roNo string, request entity.UpdateOrderBody) (err error) {
	c := context.Background()

	// parse time format YYYY-mm-dd to Rfc3339
	if request.RoDate != nil {
		RoDate, err := str.DateStrToRfc3339String(*request.RoDate)
		if err != nil {
			return err
		}
		request.RoDate = &RoDate
	}

	if request.ValDate != nil {
		ValDate, err := str.DateStrToRfc3339String(*request.ValDate)
		if err != nil {
			return err
		}
		request.ValDate = &ValDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return err
		}
		request.DueDate = &DueDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return err
		}
		request.DeliveryDate = &deliveryDate
	}

	if request.InvoiceDate != nil {
		invoiceDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
		if err != nil {
			return err
		}
		request.InvoiceDate = &invoiceDate
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Order
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderCanvasRepository.Update(txCtx, roNo, Model)
		if err != nil {
			return err
		}
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.OrderDetId != nil {
				DetailIds = append(DetailIds, *detail.OrderDetId)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.OrderDetId != nil {
				DetailIds = append(DetailIds, *detail.OrderDetId)
			}
		}
		if len(DetailIds) > 0 {
			err := service.OrderCanvasRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.OrderDetail
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo
			if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
				roDetModel.OrderDetailID = nil
				err = service.OrderCanvasRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.OrderCanvasRepository.UpdateDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}

			}
		}
		for _, detail := range request.Details.Promo {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.OrderDetail
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo
			if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
				roDetModel.OrderDetailID = nil
				err = service.OrderCanvasRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.OrderCanvasRepository.UpdateDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *orderCanvasServiceImpl) UpdateFinal(roNo string, request entity.UpdateOrderDetailFinal) (err error) {
	c := context.Background()

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Order
	err = structs.Automapper(request, &Model)
	if err != nil {
		return err
	}
	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		DetailIds := []int64{}

		for _, detail := range request.Details.Normal {
			if detail.OrderDetId != nil {
				DetailIds = append(DetailIds, *detail.OrderDetId)
			}
		}
		for _, detail := range request.Details.Promo {
			if detail.OrderDetId != nil {
				DetailIds = append(DetailIds, *detail.OrderDetId)
			}
		}
		// fmt.Println(len(DetailIds))
		// if len(DetailIds) > 0 {
		// 	err := service.OrderCanvasRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		var productIDs []int64
		for _, detail := range request.Details.Normal {
			productIDs = append(productIDs, int64(detail.ProId))
		}
		for _, detail := range request.Details.Promo {
			if !slices.Contains(productIDs, int64(detail.ProId)) {
				productIDs = append(productIDs, int64(detail.ProId))
			}
		}

		productsModel, err := service.OrderCanvasRepository.FindProductByListID(productIDs)
		if err != nil {
			return err
		}

		var productMap = model.MapProduct{}

		for _, productModel := range productsModel {
			productMap.SetProduct(productModel.ProductId, productModel)
		}

		for _, detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.OrderDetail

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo

			productModel, err := productMap.GetByID(int64(detail.ProId))
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*detail.Qty1),
				Qty2:      int(*detail.Qty2),
				Qty3:      int(*detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			QtyUnit.DoConversion()

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			roDetModel.QtyFinal = float64(totalQty)
			err = service.OrderCanvasRepository.UpdateDetail(txCtx, &roDetModel)
			if err != nil {
				return err
			}
		}

		for _, detail := range request.Details.Promo {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.OrderDetail
			if detail.ExpDate != nil {
				expDate, err := str.DateStrToRfc3339String(*detail.ExpDate)
				if err != nil {
					return err
				}
				detail.ExpDate = &expDate
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo

			productModel, err := productMap.GetByID(int64(detail.ProId))
			if err != nil {
				return err
			}

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*detail.Qty1),
				Qty2:      int(*detail.Qty2),
				Qty3:      int(*detail.Qty3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			QtyUnit.DoConversion()

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			roDetModel.QtyFinal = float64(totalQty)
			err = service.OrderCanvasRepository.UpdateDetail(txCtx, &roDetModel)
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

func (service *orderCanvasServiceImpl) Delete(custId string, RoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderCanvasRepository.Delete(txCtx, custId, RoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *orderCanvasServiceImpl) Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error) {
	product, err := service.OrderCanvasRepository.FindOneProductByProductIdAndCustId(conversionBody.ProductId, custID, parentCustID)
	if err != nil {
		return response, err
	}

	qty1 := conversionBody.Qty1
	qty2 := conversionBody.Qty2
	qty3 := conversionBody.Qty3

	rQty2 := qty1 / int64(product.ConvUnit2)
	if rQty2 > 0 {
		qty1 = qty1 % int64(product.ConvUnit2)
		qty2 += rQty2
	}

	rQty3 := qty2 / int64(product.ConvUnit3)
	if rQty3 > 0 {
		qty2 = qty2 % int64(product.ConvUnit3)
		qty3 += rQty3
	}

	response.Qty1 = qty1
	response.Qty2 = qty2
	response.Qty3 = qty3

	response.TotalQty = (int64(product.ConvUnit2)*int64(product.ConvUnit3))*qty3 + (int64(product.ConvUnit2) * qty2) + qty1

	return response, err
}

func (service *orderCanvasServiceImpl) LookupSalesman(dataFilter entity.OrderQueryFilter) (data []entity.OutletBySalesmanResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderCanvasRepository.FindAllOutletBySalesmanId(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.OutletBySalesmanResponse
		structs.Automapper(row, &vResp)

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderCanvasServiceImpl) BulkUpdateStatus(custId string, request entity.BulkUpdateStatusOrder) (err error) {
	c := context.Background()

	for index := range request.Orders {
		// End parse time format YYYY-mm-dd to Rfc339
		var Model model.Order
		err = structs.Automapper(request.Orders[index], &Model)
		if err != nil {
			return err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err = service.OrderCanvasRepository.Update(txCtx, request.Orders[index].RoNo, Model)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}
