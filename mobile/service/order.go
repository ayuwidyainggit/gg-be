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

// Operation Type Constants for Order Processing
const (
	// OprTypeOrder represents a standard Taking Order operation.
	OprTypeOrder = "O"

	// OprTypeCanvassing represents a Canvassing/Direct Sales operation.
	OprTypeCanvassing = "C"
)

type OrderService interface {
	Store(ctx context.Context, request entity.CreateOrderBody) (err error)
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
	StatisticReport(filters entity.StatisticReportFilter) (response entity.StatisticReportResponse, err error)
}

func NewOrderService(
	config env.ConfigEnv,
	orderRepository repository.OrderRepository,
	discountRepository repository.DiscountRepository,
	pjpPrincipleRepo repository.PjpPrincipalRepository,
	pjpDistributorRepo repository.PjpDistributorRepository,
	transaction repository.Dbtransaction,
) *orderServiceImpl {
	return &orderServiceImpl{
		Config:             config,
		OrderRepository:    orderRepository,
		PJPPrincipleRepo:   pjpPrincipleRepo,
		PJPDistributorRepo: pjpDistributorRepo,
		DiscountRepository: discountRepository,
		Transaction:        transaction,
	}
}

type orderServiceImpl struct {
	Config             env.ConfigEnv
	OrderRepository    repository.OrderRepository
	DiscountRepository repository.DiscountRepository
	PJPPrincipleRepo   repository.PjpPrincipalRepository
	PJPDistributorRepo repository.PjpDistributorRepository
	Transaction        repository.Dbtransaction
	Cache              *repository.Cache
}

var detOrderCreatemapperToModel = func(det entity.CreateOrderDetBody, request entity.CreateOrderBody, roNo string) (model.OrderDetail, error) {
	var gdsDetModel model.OrderDetail

	// parse time format YYYY-mm-dd to Rfc3339
	if det.ExpDate != nil {
		expDate, err := str.DateStrToRfc3339String(*det.ExpDate)
		if err != nil {
			return gdsDetModel, err
		}
		det.ExpDate = &expDate
	}
	gdsDetModel.CustId = request.CustId
	gdsDetModel.RoNo = roNo
	err := structs.Automapper(det, &gdsDetModel)
	if err != nil {
		return gdsDetModel, err
	}
	return gdsDetModel, nil
}

func generateRealOrderNumber(seq int, roDate *time.Time) string {
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

func (service *orderServiceImpl) Store(ctx context.Context, request entity.CreateOrderBody) (err error) {
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

	// salesman, err := service.OrderRepository.GetSalesman(request.CustId, request.SalesmanId)
	// if err != nil {
	// 	return err
	// }

	oprType := OprTypeOrder
	if request.OprType != nil {
		oprType = *request.OprType
	}

	if oprType == OprTypeCanvassing {
		// Validation
		for _, detail := range request.Details.Normal {
			if detail.Qty1 != nil && detail.Qty1Stok != nil && *detail.Qty1 > *detail.Qty1Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
			if detail.Qty2 != nil && detail.Qty2Stok != nil && *detail.Qty2 > *detail.Qty2Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
			if detail.Qty3 != nil && detail.Qty3Stok != nil && *detail.Qty3 > *detail.Qty3Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
		}

		for _, detail := range request.Details.Promo {
			if detail.Qty1 != nil && detail.Qty1Stok != nil && *detail.Qty1 > *detail.Qty1Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
			if detail.Qty2 != nil && detail.Qty2Stok != nil && *detail.Qty2 > *detail.Qty2Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
			if detail.Qty3 != nil && detail.Qty3Stok != nil && *detail.Qty3 > *detail.Qty3Stok {
				return fmt.Errorf("insufficient stock for product %d", detail.ProId)
			}
		}
	}

	// Snapshot the date once before the transaction starts.
	// This prevents a midnight-boundary edge case where dateStr inside the
	// transaction could differ from the one used for the invoice prefix.
	dateStr := time.Now().Format("060102") // YYMMDD

	// LOCK ORDER WARNING (Fix 2):
	// Inside the transaction, advisory locks are always acquired in this order:
	//   1. GetNextRoNumber  → pg_advisory_xact_lock(hash("rono_<custId>_<date>"))
	//   2. GetNextInvoiceNumber → pg_advisory_xact_lock(invoiceLockKey, dateInt)
	// Never reverse this order in future code paths — it will cause deadlocks
	// under concurrent canvassing requests.
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		var orderModel model.Order
		err = structs.Automapper(request, &orderModel)
		if err != nil {
			return err
		}

		// RO Number Generation inside transaction
		total, err := service.OrderRepository.GetNextRoNumber(txCtx, request.CustId, *request.RoDate)
		if err != nil {
			return err
		}
		RoNumber := generateRealOrderNumber(total, orderModel.RoDate)
		orderModel.RoNo = RoNumber
		orderModel.PoNo = nil
		orderModel.InvoiceDate = nil
		orderModel.InvoiceNo = nil
		orderModel.OprType = &oprType

		if oprType == OprTypeCanvassing {
			// Generate Invoice No (dateStr captured before transaction to avoid midnight-boundary skew)
			nextNum, err := service.OrderRepository.GetNextInvoiceNumber(txCtx, dateStr, request.CustId)
			if err != nil {
				return err
			}
			invoiceNo := fmt.Sprintf("INV%s%04d", dateStr, nextNum)
			orderModel.InvoiceNo = &invoiceNo

			now := time.Now()
			orderModel.InvoiceDate = &now

			msg := "Sufficient Stock"
			orderModel.ValidateStokMessage = &msg
		}
		if oprType == OprTypeOrder || oprType == OprTypeCanvassing {
			orderModel.ValidateStok = true
		}

		err = service.OrderRepository.Store(txCtx, &orderModel)
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

		productsModel, err := service.OrderRepository.FindProductByListID(txCtx, productIDs)
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
				Qty1:      detail.GetSafeQTY(1),
				Qty2:      detail.GetSafeQTY(2),
				Qty3:      detail.GetSafeQTY(3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			if oprType == OprTypeCanvassing {
				subtotalStock, err := service.OrderRepository.GetStockByWarehouseProductCustomer(txCtx, request.WhId, int64(detail.ProId), request.CustId)
				if err != nil {
					return err
				}

				if totalQty > int(subtotalStock) {
					return fmt.Errorf("insufficient stock for product %d", detail.ProId)
				}
			}

			gdsDetModel.Qty = float64(totalQty)
			gdsDetModel.QtyPo = float64(totalQty)
			// gdsDetModel.QtyFinal = float64(totalQty)
			gdsDetModel.QtyFinal = 0

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

			gdsDetModel.SellPriceSystem1 = detail.SellPrice1
			gdsDetModel.SellPriceSystem2 = detail.SellPrice2
			gdsDetModel.SellPriceSystem3 = detail.SellPrice3
			gdsDetModel.QtyPo1 = detail.Qty1
			gdsDetModel.QtyPo2 = detail.Qty2
			gdsDetModel.QtyPo3 = detail.Qty3
			gdsDetModel.SellPricePo1 = detail.SellPrice1
			gdsDetModel.SellPricePo2 = detail.SellPrice2
			gdsDetModel.SellPricePo3 = detail.SellPrice3
			gdsDetModel.SellPriceFinal1 = detail.SellPrice1
			gdsDetModel.SellPriceFinal2 = detail.SellPrice2
			gdsDetModel.SellPriceFinal3 = detail.SellPrice3
			gdsDetModel.DiscPo = detail.DiscValue
			gdsDetModel.DiscValuePo = detail.DiscValue
			gdsDetModel.DiscValueFinal = detail.DiscValue
			gdsDetModel.VatPo = detail.Vat
			gdsDetModel.VatValuePo = detail.VatValue

			gdsDetModel.ConvUnit2 = &productModel.ConvUnit2
			gdsDetModel.ConvUnit3 = &productModel.ConvUnit3
			gdsDetModel.ConvUnit5 = &productModel.ConvUnit4

			gdsDetModel.OriginalQtyPo1 = detail.Qty
			gdsDetModel.OriginalQtyPo2 = detail.Qty
			gdsDetModel.OriginalQtyPo3 = detail.Qty

			// set amount final and vat final (sx-1944)
			if detail.VatValue != nil {
				gdsDetModel.VatValueFinal = detail.VatValue

				amount := request.Amount
				if amount != nil {
					amountFinal := *amount + *detail.VatValue
					gdsDetModel.AmountFinal = &amountFinal
				} else {
					gdsDetModel.AmountFinal = detail.VatValue
				}
			}

			if request.Amount != nil && detail.VatValue == nil {
				gdsDetModel.AmountFinal = request.Amount
			}

			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

			var refDetID int64
			if gdsDetModel.OrderDetailID != nil {
				refDetID = int64(*gdsDetModel.OrderDetailID)
			}

			unitPrice := float64(0)
			if detail.SellPrice1 != nil {
				unitPrice = *detail.SellPrice1
			}

			if oprType == OprTypeCanvassing {
				stockSO := model.Stock{
					CustID:      request.CustId,
					StockDate:   time.Now().Format("2006-01-02"),
					TrCode:      "SO",
					TrNo:        orderModel.RoNo,
					WhID:        int64(request.WhId),
					ProID:       int64(detail.ProId),
					ItemCdn:     1,
					QtyIn:       0,
					QtyOut:      float64(totalQty),
					UnitPrice:   unitPrice,
					Cogs:        productModel.Cogs,
					RefDetID:    refDetID,
					CreatedAt:   time.Now().Unix(),
					QtyInOrder:  0,
					QtyOutOrder: 0,
				}

				if err := service.OrderRepository.StoreStock(txCtx, &stockSO); err != nil {
					return err
				}

				// takeout by request.
				// stockCO := stockSO
				// stockCO.StockID = 0 // use auto-increment from DB
				// stockCO.TrCode = "CO"
				// stockCO.TrNo = orderModel.RoNo + "-CO"
				// stockCO.QtyOut = 0
				// stockCO.QtyInOrder = float64(totalQty)

				// if err := service.OrderRepository.StoreStock(txCtx, &stockCO); err != nil {
				// 	return err
				// }

				if err := service.OrderRepository.UpdateWarehouseStock(txCtx, request.WhId, int64(detail.ProId), request.CustId, float64(totalQty)); err != nil {
					return err
				}

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
				Qty1:      Detail.GetSafeQTY(1),
				Qty2:      Detail.GetSafeQTY(2),
				Qty3:      Detail.GetSafeQTY(3),
				ConvUnit2: int(productModel.ConvUnit2),
				ConvUnit3: int(productModel.ConvUnit3),
			}

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			if oprType == OprTypeCanvassing {
				subtotalStock, err := service.OrderRepository.GetStockByWarehouseProductCustomer(txCtx, request.WhId, int64(Detail.ProId), request.CustId)
				if err != nil {
					return err
				}

				if totalQty > int(subtotalStock) {
					return fmt.Errorf("insufficient stock for product %d", Detail.ProId)
				}
			}

			gdsDetModel.Qty = float64(totalQty)
			gdsDetModel.QtyPo = float64(totalQty)
			gdsDetModel.QtyFinal = float64(totalQty)

			err = structs.Automapper(Detail, &gdsDetModel)
			if err != nil {
				return err
			}

			gdsDetModel.SellPriceSystem1 = Detail.SellPrice1
			gdsDetModel.SellPriceSystem2 = Detail.SellPrice2
			gdsDetModel.SellPriceSystem3 = Detail.SellPrice3
			gdsDetModel.SellPriceFinal1 = Detail.SellPrice1
			gdsDetModel.SellPriceFinal2 = Detail.SellPrice2
			gdsDetModel.SellPriceFinal3 = Detail.SellPrice3
			gdsDetModel.QtyPo1 = Detail.Qty1
			gdsDetModel.QtyPo2 = Detail.Qty2
			gdsDetModel.QtyPo3 = Detail.Qty3
			gdsDetModel.SellPricePo1 = Detail.SellPrice1
			gdsDetModel.SellPricePo2 = Detail.SellPrice2
			gdsDetModel.SellPricePo3 = Detail.SellPrice3
			gdsDetModel.DiscPo = Detail.DiscValue
			gdsDetModel.DiscValuePo = Detail.DiscValue
			gdsDetModel.DiscValueFinal = Detail.DiscValue
			gdsDetModel.VatPo = Detail.Vat
			gdsDetModel.VatValuePo = Detail.VatValue

			gdsDetModel.ConvUnit2 = &productModel.ConvUnit2
			gdsDetModel.ConvUnit3 = &productModel.ConvUnit3
			gdsDetModel.ConvUnit5 = &productModel.ConvUnit4
			gdsDetModel.OriginalQtyPo1 = Detail.Qty
			gdsDetModel.OriginalQtyPo2 = Detail.Qty
			gdsDetModel.OriginalQtyPo3 = Detail.Qty
			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

			var refDetID int64
			if gdsDetModel.OrderDetailID != nil {
				refDetID = int64(*gdsDetModel.OrderDetailID)
			}

			unitPrice := float64(0)
			if Detail.SellPrice1 != nil {
				unitPrice = *Detail.SellPrice1
			}

			if oprType == OprTypeCanvassing {
				stockSO := model.Stock{
					CustID:      request.CustId,
					StockDate:   time.Now().Format("2006-01-02"),
					TrCode:      "SO",
					TrNo:        orderModel.RoNo,
					WhID:        int64(request.WhId),
					ProID:       int64(Detail.ProId),
					ItemCdn:     1,
					QtyIn:       0,
					QtyOut:      float64(totalQty),
					UnitPrice:   unitPrice,
					Cogs:        productModel.Cogs,
					RefDetID:    refDetID,
					CreatedAt:   time.Now().Unix(),
					QtyInOrder:  0,
					QtyOutOrder: 0,
				}
				if err := service.OrderRepository.StoreStock(txCtx, &stockSO); err != nil {
					return err
				}

				// takeout by request.
				// stockCO := stockSO
				// stockCO.StockID = 0 // use auto-increment from DB
				// stockCO.TrCode = "CO"
				// stockCO.TrNo = orderModel.RoNo + "-CO"
				// stockCO.QtyOut = 0
				// stockCO.QtyOutOrder = float64(totalQty)

				// if err := service.OrderRepository.StoreStock(txCtx, &stockCO); err != nil {
				// 	return err
				// }

				if err := service.OrderRepository.UpdateWarehouseStock(txCtx, request.WhId, int64(Detail.ProId), request.CustId, float64(totalQty)); err != nil {
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

func (service *orderServiceImpl) StoreNoOrder(request entity.CreateNoOrderBody) (err error) {
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
		if err := service.OrderRepository.StoreNoOrder(txCtx, &noOrderModel); err != nil {
			return fmt.Errorf("failed to store NoOrder: %w", err)
		}
		return nil
	})

	return err
}

func (service *orderServiceImpl) Detail(RoNo string, custID string) (response entity.OrderResponse, err error) {
	ro, err := service.OrderRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	details, err := service.OrderRepository.FindDetail(RoNo, custID)
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

func getStartAndEndOfMonth(dateStr string) (string, string, error) {
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

func (service *orderServiceImpl) SummaryBySalesman(summaryFilter entity.SummaryTotalFilter) (response entity.ResponseSummaryTotal, err error) {

	currentDate := time.Now().Format("2006-01-02")

	if summaryFilter.Date != nil && *summaryFilter.Date != 0 {
		currentDate = str.UnixTimestampToUtcTime(*summaryFilter.Date).Format("2006-01-02") // Mengubah time.Time menjadi string
	}

	startOfMonth, endOfMonth, err := getStartAndEndOfMonth(currentDate)
	if err != nil {
		return response, err
	}

	var dataStatus = []int{1, 2}

	summaryToday, err := service.OrderRepository.SummaryTotalBySalesmanAndDate(summaryFilter.SalesmanId, currentDate, currentDate, dataStatus, summaryFilter.CustId)
	if err != nil {
		return response, err
	}
	dataStatus = []int{1, 2, 3, 4, 5, 6, 7}
	summaryMonth, err := service.OrderRepository.SummaryTotalBySalesmanAndDate(summaryFilter.SalesmanId, startOfMonth, endOfMonth, dataStatus, summaryFilter.CustId)
	if err != nil {
		return response, err
	}

	response.SummaryToday = summaryToday.TotalSummary
	response.SummaryMonth = summaryMonth.TotalSummary
	response.SalesmanId = summaryFilter.SalesmanId
	response.CreatedAt = currentDate

	return response, nil
}

func (service *orderServiceImpl) List(dataFilter entity.OrderQueryFilter) (data []entity.OrderListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderRepository.FindAllByCustId(dataFilter)
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

func (service *orderServiceImpl) ListNoOrder(dataFilter entity.NoOrderQueryFilter) (data []entity.NoOrderListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderRepository.FindAllNoOrderByCustId(dataFilter)
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

func (service *orderServiceImpl) Update(roNo string, request entity.UpdateOrderBody) (err error) {
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
		err = service.OrderRepository.Update(txCtx, roNo, Model)
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
			err := service.OrderRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
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
				err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
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
				err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			} else {
				roDetModel.CustId = ""
				err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
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

func (service *orderServiceImpl) UpdateFinal(roNo string, request entity.UpdateOrderDetailFinal) (err error) {
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
		// 	err := service.OrderRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
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

		productsModel, err := service.OrderRepository.FindProductByListID(txCtx, productIDs)
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
			err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
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
			err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
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

func (service *orderServiceImpl) Delete(custId string, RoNo string, userId int64) (err error) {
	c := context.Background()
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err = service.OrderRepository.Delete(txCtx, custId, RoNo, userId)
		if err != nil {
			return err
		}
		return nil
	})

	return err
}

func (service *orderServiceImpl) Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error) {
	product, err := service.OrderRepository.FindOneProductByProductIdAndCustId(conversionBody.ProductId, custID, parentCustID)
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

func (service *orderServiceImpl) LookupSalesman(dataFilter entity.OrderQueryFilter) (data []entity.OutletBySalesmanResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderRepository.FindAllOutletBySalesmanId(dataFilter)
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

func (service *orderServiceImpl) BulkUpdateStatus(custId string, request entity.BulkUpdateStatusOrder) (err error) {
	c := context.Background()

	for index := range request.Orders {
		// End parse time format YYYY-mm-dd to Rfc339
		var Model model.Order
		err = structs.Automapper(request.Orders[index], &Model)
		if err != nil {
			return err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			err = service.OrderRepository.Update(txCtx, request.Orders[index].RoNo, Model)
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

func (service *orderServiceImpl) StatisticReport(filters entity.StatisticReportFilter) (response entity.StatisticReportResponse, err error) {
	var isPrincipal bool
	var pjpID int64
	ctx := context.Background()
	if len(filters.CustID) == 6 && !filters.IsDistributor {
		isPrincipal = true
		pjp, errFind := service.PJPPrincipleRepo.FindOneBySalesmanAndCustID(ctx, int64(filters.EmpID), filters.CustID)
		if errFind != nil {
			return response, errFind
		}
		pjpID = pjp.ID
	} else {
		pjp, errFind := service.PJPDistributorRepo.FindOneBySalesmanAndCustID(ctx, int64(filters.EmpID), filters.CustID)
		if errFind != nil {
			return response, errFind
		}
		pjpID = pjp.ID
	}

	var startDate string
	var endDate string
	now := time.Now()
	switch filters.Type {
	case "day":
		startDate = now.Format("2006-01-02") + " 00:00:00"
		endDate = now.Format("2006-01-02") + " 23:59:59"
	case "month":
		firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		startDate = firstOfMonth.Format("2006-01-02") + " 00:00:00"
		endDate = lastOfMonth.Format("2006-01-02") + " 23:59:59"
	default:
		startDate = now.Format("2006-01-02") + " 00:00:00"
		endDate = now.Format("2006-01-02") + " 23:59:59"
	}

	summarySales, err := service.OrderRepository.GetSummarySales(filters.CustID, filters.EmpID, startDate, endDate)
	if err != nil {
		return response, err
	}

	newResp := entity.StatisticReportResponse{
		Sales: summarySales,
	}

	if isPrincipal {
		overviewList, err := service.PJPPrincipleRepo.GetVisitOverview(ctx, pjpID, int64(filters.EmpID), startDate, endDate)
		if err != nil {
			return response, err
		}
		newResp.TotalVisitList = overviewList.TotalVisitList
		newResp.TotalVisit = overviewList.TotalVisit
		newResp.TotalNotVisit = overviewList.TotalNotVisit
		newResp.VisitPlanned = overviewList.VisitPlanned
		newResp.VisitNotPlanned = overviewList.VisitNotPlanned
		newResp.NotVisitPlanned = overviewList.NotVisitPlanned
		newResp.NotVisitNotPlanned = overviewList.NotVisitNotPlanned
		newResp.TotalBuy = overviewList.TotalBuy
		newResp.TotalNotBuy = overviewList.TotalNotBuy

		totalNotBuy, err := service.PJPPrincipleRepo.GetNotBuyReasons(ctx, pjpID, int64(filters.EmpID), startDate, endDate)
		if err != nil {
			return response, err
		}
		var listNotBuy []entity.NotBuyReasonItem
		for _, v := range totalNotBuy {
			listNotBuy = append(listNotBuy, entity.NotBuyReasonItem{
				Reason: v.Reason,
				Total:  v.Total,
			})
		}
		if len(listNotBuy) > 0 {
			newResp.NotBuyReasonData = listNotBuy
		}

		totalSkipReason, err := service.PJPPrincipleRepo.GetSkipReasons(ctx, pjpID, startDate, endDate)
		if err != nil {
			return response, err
		}
		var listSkipReason []entity.SkipReasonItem
		for _, v := range totalSkipReason {
			listSkipReason = append(listSkipReason, entity.SkipReasonItem{
				SkipReason: v.SkipReason,
				Total:      v.Total,
			})
		}
		if len(listSkipReason) > 0 {
			newResp.SkipReasonData = listSkipReason
		}
	} else {
		overviewList, err := service.PJPDistributorRepo.GetVisitOverview(ctx, pjpID, int64(filters.EmpID), startDate, endDate)
		if err != nil {
			return response, err
		}
		newResp.TotalVisitList = overviewList.TotalVisitList
		newResp.TotalVisit = overviewList.TotalVisit
		newResp.TotalNotVisit = overviewList.TotalNotVisit
		newResp.VisitPlanned = overviewList.VisitPlanned
		newResp.VisitNotPlanned = overviewList.VisitNotPlanned
		newResp.NotVisitPlanned = overviewList.NotVisitPlanned
		newResp.NotVisitNotPlanned = overviewList.NotVisitNotPlanned
		newResp.TotalBuy = overviewList.TotalBuy
		newResp.TotalNotBuy = overviewList.TotalNotBuy

		totalNotBuy, err := service.PJPDistributorRepo.GetNotBuyReasons(ctx, pjpID, int64(filters.EmpID), startDate, endDate)
		if err != nil {
			return response, err
		}
		var listNotBuy []entity.NotBuyReasonItem
		for _, v := range totalNotBuy {
			listNotBuy = append(listNotBuy, entity.NotBuyReasonItem{
				Reason: v.Reason,
				Total:  v.Total,
			})
		}
		if len(listNotBuy) > 0 {
			newResp.NotBuyReasonData = listNotBuy
		}

		totalSkipReason, err := service.PJPDistributorRepo.GetSkipReasons(ctx, pjpID, startDate, endDate)
		if err != nil {
			return response, err
		}
		var listSkipReason []entity.SkipReasonItem
		for _, v := range totalSkipReason {
			listSkipReason = append(listSkipReason, entity.SkipReasonItem{
				SkipReason: v.SkipReason,
				Total:      v.Total,
			})
		}
		if len(listSkipReason) > 0 {
			newResp.SkipReasonData = listSkipReason
		}
	}

	return newResp, nil
}
