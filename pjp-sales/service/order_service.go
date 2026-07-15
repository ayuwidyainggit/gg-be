package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/constant"
	"sales/pkg/conversion"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
)

type OrderService interface {
	Store(request entity.CreateOrderBody, validateOrderRequest entity.ValidateResponse) (data entity.CreateOrderResponse, err error)
	Detail(RoNo string, custID string) (response entity.OrderResponse, err error)
	DetailV2(RoNo string, custID string, parentCustID string) (response entity.OrderResponse, err error)
	List(dataFilter entity.OrderQueryFilter) (data []entity.OrderListResponse, total int64, lastPage int, err error)
	ProformaInvoiceList(dataFilter entity.ProformaInvoiceQueryFilter) (data []entity.ProformaInvoiceListResponse, total int64, lastPage int, err error)
	PrintProformaInvoice(ctx context.Context, request entity.PrintProformaInvoiceRequest, custId string, userId int64) (response entity.PrintProformaInvoiceResponse, err error)
	Update(roNo string, request entity.UpdateOrderBody, validateOrderRequest entity.ValidateResponse) (err error)
	UpdateFinal(roNo string, request entity.UpdateOrderDetailFinal, validateOrderRequest entity.ValidateResponse) (err error)
	Delete(custId string, RoNo string, userId int64) (err error)
	Conversion(conversionBody entity.CreateConversionBody, custID string, parentCustID string) (response entity.OrderConversionResponse, err error)
	LookupSalesman(dataFilter entity.OrderQueryFilter) (data []entity.OutletBySalesmanResponse, total int64, lastPage int, err error)
	BulkUpdateStatus(custId string, request entity.BulkUpdateStatusOrder) (err error)
	DetailDiscount(criteria entity.OrderDiscountQuery) (response entity.DiscountCriteria, err error)

	SetValidateOrderRequest(roNo string, validateOrderBody *entity.ValidateOrderBody) (err error)
	DetailNoCustID(RoNo string, custIDOrigin string, empID *int64) (response entity.OrderResponse, err error)
	GetMinimumPriceProduct(request entity.OrderMinimumPriceFilter) (entity.OrderMinimumPriceResp, error)
	UpdateEnhance(ctx context.Context, roNo string, request entity.EditOrderEnhanceBody) error
	ProcessEnhanceWithoutProductEdit(ctx context.Context, roNo string, custId string, updatedBy int64) error

	// Unified Calculation Helper Functions
	CalculateLineVAT(gross, disc, promo, vatRate float64) float64
	CalculateLineDiscount(custId, parentCustId string, outletId int, proId int, amount float64) (float64, string, error)
	RecalculateOrderTotals(details []model.OrderDetail) (subTotal, discTotal, vatTotal, total float64)

	ExportTemplate(format string) (*bytes.Buffer, string, string, error)
	ImportOrders(custId string, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error)
	ValidateImport(custId string, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportSummary, error)
}

func NewOrderService(orderRepository repository.OrderRepository, validateOrderRepository repository.ValidateOrderRepository, promotionRepository repository.PromotionRepository, promotionV2Repository repository.PromotionV2Repository, discountRepository repository.DiscountRepository, stockRepository repository.StockRepository, transaction repository.Dbtransaction) *orderServiceImpl {
	return &orderServiceImpl{
		OrderRepository:         orderRepository,
		ValidateOrderRepository: validateOrderRepository,
		PromotionRepository:     promotionRepository,
		PromotionV2Repository:   promotionV2Repository,
		DiscountRepository:      discountRepository,
		StockRepository:         stockRepository,
		Transaction:             transaction,
	}
}

type orderServiceImpl struct {
	OrderRepository         repository.OrderRepository
	ValidateOrderRepository repository.ValidateOrderRepository
	PromotionRepository     repository.PromotionRepository
	PromotionV2Repository   repository.PromotionV2Repository
	DiscountRepository      repository.DiscountRepository
	StockRepository         repository.StockRepository
	Transaction             repository.Dbtransaction
}

func (service *orderServiceImpl) validateCreateOrderStockWithRewards(request entity.CreateOrderBody, dataStatus *int64) error {
	// ponytail: import (data_source=3) bypasses stock validation per
	// SX-2434 product spec; revisit when import stock policy changes.
	if service.OrderRepository == nil || request.WhId == nil || !ShouldValidateStockOnCreate(request.OrderType) || !isProcessedDataStatus(dataStatus) || isImportedOrder(&request) {
		return nil
	}

	requestedByProduct := make(map[int64]float64)
	proIDSet := make(map[int64]struct{})
	appendDetail := func(detail entity.CreateOrderDetBody) error {
		proID := int64(detail.ProId)
		proIDSet[proID] = struct{}{}
		qtyUnit := &conversion.QtyUnit{
			Qty1:      int(getValueOrDefault(detail.Qty1, 0)),
			Qty2:      int(getValueOrDefault(detail.Qty2, 0)),
			Qty3:      int(getValueOrDefault(detail.Qty3, 0)),
			ConvUnit2: getValueOrDefaultInt(detail.ConvUnit2, 0),
			ConvUnit3: getValueOrDefaultInt(detail.ConvUnit3, 0),
		}
		totalQty, err := qtyUnit.ToTotalQuantity()
		if err != nil {
			return err
		}
		requestedByProduct[proID] += float64(totalQty)
		return nil
	}

	for _, detail := range request.Details.Normal {
		if err := appendDetail(detail); err != nil {
			return err
		}
	}
	for _, detail := range request.Details.Promo {
		if err := appendDetail(detail); err != nil {
			return err
		}
	}
	if len(proIDSet) == 0 {
		return nil
	}

	proIDs := make([]int64, 0, len(proIDSet))
	for proID := range proIDSet {
		proIDs = append(proIDs, proID)
	}
	availableByProduct, err := service.OrderRepository.FindWarehouseStockByWhIdAndProIds(request.CustId, *request.WhId, proIDs)
	if err != nil {
		return err
	}
	for proID, requestedQty := range requestedByProduct {
		if requestedQty > availableByProduct[proID] {
			return errors.New("Insufficient Stock")
		}
	}

	return nil
}

func getValueOrDefaultInt(ptr *int, defaultValue int) int {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func syncCreateOrderFinalQtyFromSalesQty(detail *model.OrderDetail) {
	detail.QtyFinal = detail.Qty
	detail.Qty1Final = detail.Qty1
	detail.Qty2Final = detail.Qty2
	detail.Qty3Final = detail.Qty3
	detail.Qty4Final = detail.Qty4
	detail.Qty5Final = detail.Qty5
}

// MapDataSourceToSource maps data_source (int64) to source (string)
// 1 = "web", 2 = "mobile", null = ""
func MapDataSourceToSource(dataSource *int64) *string {
	if dataSource == nil {
		emptyStr := ""
		return &emptyStr
	}
	var source string
	switch *dataSource {
	case 1:
		source = "web"
	case 2:
		source = "mobile"
	default:
		emptyStr := ""
		return &emptyStr
	}
	return &source
}

// Unified Calculation Helper Functions

// CalculateLineVAT calculates VAT value based on taxable amount
func (service *orderServiceImpl) CalculateLineVAT(gross, disc, promo, vatRate float64) float64 {
	return calculateVatValue(1, 0, 0, gross, 0, 0, promo, disc, vatRate)
}

// CalculateLineDiscount calculates discount based on product, outlet and amount (tier logic)
func (service *orderServiceImpl) CalculateLineDiscount(custId, parentCustId string, outletId int, proId int, amount float64) (float64, string, error) {
	product, err := service.OrderRepository.FindProductByID(proId)
	if err != nil {
		return 0, "", err
	}

	outlet, err := service.OrderRepository.FindOutletByID(outletId, custId, parentCustId)
	if err != nil {
		return 0, "", err
	}

	discount, err := service.OrderRepository.FindDiscountByProductAndOutlet(product, outlet)
	if err != nil {
		// No discount found is not a system error
		return 0, "", nil
	}

	// Find criteria based on amount (SubTotalPrincipal)
	criteria, err := service.OrderRepository.FindDiscountCriteriaBySubTotal(discount.DiscountId, int(amount))
	if err != nil {
		// No criteria met is not a system error
		return 0, "", nil
	}

	discValue := criteria.SlabReward
	if criteria.SlabRewardType == entity.DiscountRewardTypePercentage {
		discValue = math.Round((float64(discValue) * amount) / 100)
	}

	return float64(discValue), discount.DiscountId, nil
}

// RecalculateOrderTotals aggregates totals from order details
// Assumes details have correct Line Totals (Amount, DiscValue, VatValue)
func (service *orderServiceImpl) RecalculateOrderTotals(details []model.OrderDetail) (subTotal, discTotal, vatTotal, total float64) {
	for _, det := range details {
		nett := getValueOrDefault(det.AmountFinal, 0)
		vat := getValueOrDefault(det.VatValueFinal, 0)
		disc := getValueOrDefault(det.DiscValueFinal, 0)
		promo := getValueOrDefault(det.PromoValueFinal, 0)

		gross := nett - vat + disc + promo

		subTotal += gross
		discTotal += disc
		vatTotal += vat
		total += nett
	}
	return
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

func (service *orderServiceImpl) Store(request entity.CreateOrderBody, validateOrderRequest entity.ValidateResponse) (data entity.CreateOrderResponse, err error) {
	c := context.Background()

	promoSnapshotByProduct := map[int]promoAggregateRow{}
	promoSnapshotRemarks := []string{}

	if !isImportedOrder(&request) {
		var consultResp []entity.ConsultPromoResp
		consultResp, promoSnapshotByProduct, promoSnapshotRemarks, err = service.prepareCreateOrderPromoState(&request)
		if err != nil {
			return data, err
		}
		request.Rewards = buildCreateOrderRewardsFromPromoV2(request.CustId, consultResp)

		var consultDiscount entity.ConsultDiscountOrderBody
		if err = structs.Automapper(request, &consultDiscount); err != nil {
			return data, err
		}

		if err = service.ConsultDiscountBeforeStore(&consultDiscount); err != nil {
			return data, err
		}

		if err = structs.Automapper(consultDiscount, &request); err != nil {
			return data, err
		}
	} else {
		request.Rewards = nil
	}

	// parse time format YYYY-mm-dd to Rfc3339
	if request.RoDate != nil {
		roDate, err := str.DateStrToRfc3339String(*request.RoDate)
		if err != nil {
			return data, err
		}
		request.RoDate = &roDate
	}

	if request.ValDate != nil {
		valDate, err := str.DateStrToRfc3339String(*request.ValDate)
		if err != nil {
			return data, err
		}
		request.ValDate = &valDate
	}

	if request.DueDate != nil {
		DueDate, err := str.DateStrToRfc3339String(*request.DueDate)
		if err != nil {
			return data, err
		}
		request.DueDate = &DueDate
	}

	if request.DeliveryDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.DeliveryDate)
		if err != nil {
			return data, err
		}
		request.DeliveryDate = &deliveryDate
	}

	if request.InvoiceDate != nil {
		deliveryDate, err := str.DateStrToRfc3339String(*request.InvoiceDate)
		if err != nil {
			return data, err
		}
		request.InvoiceDate = &deliveryDate
	}

	total, err := service.OrderRepository.CountAllRoByCustId(request.CustId, *request.RoDate)
	if err != nil {
		return data, err
	}

	var orderModel model.Order
	err = structs.Automapper(request, &orderModel)
	if err != nil {
		return data, err
	}
	orderModel.PromoRemarksSo = model.JSONStringArray(append([]string{}, promoSnapshotRemarks...))
	orderModel.PromoRemarksFinal = model.JSONStringArray(append([]string{}, promoSnapshotRemarks...))

	// Fetch Outlet Data to get Address
	outletData, err := service.OrderRepository.FindOutletByID(int(request.OutletID), request.CustId, request.ParentCustId)
	if err != nil {
		return data, err
	}
	// Per docx: address1 of an imported order must be NULL, regardless of the
	// outlet master. Non-import flows keep the existing fallback.
	if isImportedOrder(&request) {
		orderModel.Address1 = nil
	} else {
		orderModel.Address1 = outletData.Address1
	}

	RoNumber := generateRealOrderNumber(total, orderModel.RoDate)
	// ponytail: allow import to provide its own ro_no (document_no);
	// fall back to the generator when the caller did not pre-set one.
	if strings.TrimSpace(orderModel.RoNo) == "" {
		orderModel.RoNo = RoNumber
	}
	log.Info("Ro Number : ", orderModel.RoNo)
	orderModel.PoNo = nil
	if !isImportedOrder(&request) {
		orderModel.InvoiceDate = nil
		orderModel.InvoiceNo = nil
	}

	statusDecision := determineSalesOrderStatus(validateOrderRequest, outletRulesFromOutletRead(outletData))
	if !isImportedOrder(&request) {
		if err = ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
			return data, err
		}
	}
	dataStatus := resolveCreateOrderDataStatus(request.OrderType, statusDecision)
	if isImportedOrder(&request) {
		dataStatus = int64(importDataStatus)
	}
	orderModel.DataStatus = &dataStatus
	if err = service.validateCreateOrderStockWithRewards(request, orderModel.DataStatus); err != nil {
		return data, err
	}
	applyValidationResultToOrderModel(&orderModel, validateOrderRequest)
	if isImportedOrder(&request) {
		orderModel.DataStatus = &dataStatus
		isSalesMapping := true
		orderModel.IsSalesMapping = &isSalesMapping
	}
	if IsTakingOrder(request.OrderType) {
		applyTakingOrderValidationSnapshot(&orderModel)
		if orderModel.OprType == nil || strings.TrimSpace(*orderModel.OprType) == "" {
			orderModel.OprType = stringPtr(orderTypeTakingOrder)
		}
	}

	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		err := service.OrderRepository.Store(txCtx, &orderModel)
		if err != nil {
			return err
		}

		var salesOrderStockUpdateEntities []*entity.SalesOrderStockUpdate

		for index, detail := range request.Details.Normal {
			var gdsDetModel model.OrderDetail

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

			gdsDetModel.Qty = float64(totalQty)
			if !isImportedOrder(&request) {
				gdsDetModel.QtyPo = float64(totalQty)
			}
			gdsDetModel.QtyFinal = float64(totalQty)

			err = structs.Automapper(detail, &gdsDetModel)
			if err != nil {
				return err
			}
			syncCreateOrderFinalQtyFromSalesQty(&gdsDetModel)
			gdsDetModel.ItemType = 1
			if IsTakingOrder(request.OrderType) {
				applyTakingOrderDetailFields(detail, &gdsDetModel, float64(totalQty))
				gdsDetModel.Qty = 0
				gdsDetModel.QtyFinal = 0
			}

			if detail.ProId == 0 {
				if isImportedOrder(&request) {
					gdsDetModel.SellPriceSystem1 = float64Ptr(0)
					gdsDetModel.SellPriceSystem2 = float64Ptr(0)
					gdsDetModel.SellPriceSystem3 = float64Ptr(0)
					gdsDetModel.SellPriceFinal1 = detail.SellPrice1
					gdsDetModel.SellPriceFinal2 = detail.SellPrice2
					gdsDetModel.SellPriceFinal3 = detail.SellPrice3
				} else {
					return errors.New("record not found")
				}
			} else {
				product, err := service.OrderRepository.FindProductByID(detail.ProId)
				if err != nil {
					return err
				}

				gdsDetModel.SellPriceSystem1 = &product.SellPrice1
				gdsDetModel.SellPriceSystem2 = &product.SellPrice1
				gdsDetModel.SellPriceSystem3 = &product.SellPrice1
				gdsDetModel.SellPriceFinal1 = detail.SellPrice1
				gdsDetModel.SellPriceFinal2 = detail.SellPrice2
				gdsDetModel.SellPriceFinal3 = detail.SellPrice3
			}

			if promoSnapshot, exists := promoSnapshotByProduct[index+1]; exists {
				gdsDetModel.PromoSo1 = float64Ptr(promoSnapshot.Promo1)
				gdsDetModel.PromoSo2 = float64Ptr(promoSnapshot.Promo2)
				gdsDetModel.PromoSo3 = float64Ptr(promoSnapshot.Promo3)
				gdsDetModel.PromoSo4 = float64Ptr(promoSnapshot.Promo4)
				gdsDetModel.PromoSo5 = float64Ptr(promoSnapshot.Promo5)
				gdsDetModel.PromoFinal1 = float64Ptr(promoSnapshot.Promo1)
				gdsDetModel.PromoFinal2 = float64Ptr(promoSnapshot.Promo2)
				gdsDetModel.PromoFinal3 = float64Ptr(promoSnapshot.Promo3)
				gdsDetModel.PromoFinal4 = float64Ptr(promoSnapshot.Promo4)
				gdsDetModel.PromoFinal5 = float64Ptr(promoSnapshot.Promo5)
				gdsDetModel.PromoRemarksSo = model.JSONStringArray(append([]string{}, promoSnapshot.Remarks...))
				gdsDetModel.PromoRemarksFinal = model.JSONStringArray(append([]string{}, promoSnapshot.Remarks...))
				gdsDetModel.IsProductPromotionSo = boolPtr(promoSnapshot.IsProductPromotion)
				gdsDetModel.IsProductPromotionFinal = boolPtr(promoSnapshot.IsProductPromotion)
				defaultFalse := false
				gdsDetModel.IsProductPromotionPo = &defaultFalse
			}
			if detail.IsProductPromotionSo != nil {
				gdsDetModel.IsProductPromotionSo = boolPtr(*detail.IsProductPromotionSo)
			}
			if detail.IsProductPromotionFinal != nil {
				gdsDetModel.IsProductPromotionFinal = boolPtr(*detail.IsProductPromotionFinal)
			}
			if gdsDetModel.IsProductPromotionPo == nil {
				defaultFalse := false
				gdsDetModel.IsProductPromotionPo = &defaultFalse
			}
			if IsTakingOrder(request.OrderType) {
				gdsDetModel.Qty1 = nil
				gdsDetModel.Qty2 = nil
				gdsDetModel.Qty3 = nil
				gdsDetModel.Qty1Final = nil
				gdsDetModel.Qty2Final = nil
				gdsDetModel.Qty3Final = nil
			}

			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}
			if isImportedOrder(&request) {
				if gdsDetModel.OrderDetailID == nil {
					return errors.New("stored order detail missing order_detail_id")
				}
				if err := service.OrderRepository.UpdateDetailPartial(txCtx, int64(*gdsDetModel.OrderDetailID), request.CustId, map[string]interface{}{"qty_po": nil}); err != nil {
					return err
				}
			}

			roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate)
			if err != nil {
				return err
			}

			if ShouldMutateInventoryOnCreate(request.OrderType) && isProcessedDataStatus(orderModel.DataStatus) {
				salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
					CustID:         request.CustId,
					WhID:           *request.WhId,
					ProID:          int64(detail.ProId),
					StockDate:      *roDate,
					TrCode:         orderModel.RoNo[0:2],
					TrNo:           orderModel.RoNo,
					QtyOrderBefore: nil,
					QtyOrder:       gdsDetModel.QtyFinal,
					UnitPrice:      *detail.SellPrice1,
					RefDetId:       int64(*gdsDetModel.OrderDetailID),
				}
				salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)
			}
		}

		for _, Detail := range request.Details.Promo {
			var gdsDetModel model.OrderDetail

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

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*Detail.Qty1),
				Qty2:      int(*Detail.Qty2),
				Qty3:      int(*Detail.Qty3),
				ConvUnit2: int(*Detail.ConvUnit2),
				ConvUnit3: int(*Detail.ConvUnit3),
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
			syncCreateOrderFinalQtyFromSalesQty(&gdsDetModel)
			applyOriginalTakingOrderQty(request.OrderType, &gdsDetModel)

			if Detail.ProId == 0 {
				if isImportedOrder(&request) {
					gdsDetModel.SellPriceSystem1 = float64Ptr(0)
					gdsDetModel.SellPriceSystem2 = float64Ptr(0)
					gdsDetModel.SellPriceSystem3 = float64Ptr(0)
					gdsDetModel.SellPriceFinal1 = Detail.SellPrice1
					gdsDetModel.SellPriceFinal2 = Detail.SellPrice2
					gdsDetModel.SellPriceFinal3 = Detail.SellPrice3
				} else {
					return errors.New("record not found")
				}
			} else {
				product, err := service.OrderRepository.FindProductByID(Detail.ProId)
				if err != nil {
					return err
				}

				gdsDetModel.SellPriceSystem1 = &product.SellPrice1
				gdsDetModel.SellPriceSystem2 = &product.SellPrice1
				gdsDetModel.SellPriceSystem3 = &product.SellPrice1
				gdsDetModel.SellPriceFinal1 = Detail.SellPrice1
				gdsDetModel.SellPriceFinal2 = Detail.SellPrice2
				gdsDetModel.SellPriceFinal3 = Detail.SellPrice3
			}
			if gdsDetModel.IsProductPromotionPo == nil {
				defaultFalse := false
				gdsDetModel.IsProductPromotionPo = &defaultFalse
			}

			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}
			if isImportedOrder(&request) {
				if gdsDetModel.OrderDetailID == nil {
					return errors.New("stored order detail missing order_detail_id")
				}
				if err := service.OrderRepository.UpdateDetailPartial(txCtx, int64(*gdsDetModel.OrderDetailID), request.CustId, map[string]interface{}{"qty_po": nil}); err != nil {
					return err
				}
			}

			roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate)
			if err != nil {
				return err
			}
			if ShouldMutateInventoryOnCreate(request.OrderType) && isProcessedDataStatus(orderModel.DataStatus) {
				salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
					CustID:         request.CustId,
					WhID:           *request.WhId,
					ProID:          int64(Detail.ProId),
					StockDate:      *roDate,
					TrCode:         orderModel.RoNo[0:2],
					TrNo:           orderModel.RoNo,
					QtyOrderBefore: nil,
					QtyOrder:       gdsDetModel.QtyFinal,
					UnitPrice:      *Detail.SellPrice1,
					RefDetId:       int64(*gdsDetModel.OrderDetailID),
				}
				salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)
			}
		}

		if len(salesOrderStockUpdateEntities) > 0 {
			err = service.StockRepository.SalesStockUpdates(txCtx, salesOrderStockUpdateEntities)
			if err != nil {
				return err
			}
		}

		log.Info("REWARD BEFORE STORE :", len(request.Rewards))
		for _, rewardRequest := range request.Rewards {
			var reward model.OrderReward

			if err = structs.Automapper(rewardRequest, &reward); err != nil {
				return err
			}

			reward.RoNo = orderModel.RoNo

			if err = service.OrderRepository.StoreReward(txCtx, &reward); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return data, err
	}

	data.RoNo = orderModel.RoNo
	return data, nil
}

func (service *orderServiceImpl) SetValidateOrderRequest(roNo string, validateOrderRequest *entity.ValidateOrderBody) (err error) {

	var ro model.OrderList
	var total float64
	if roNo != "" {
		ro, err = service.OrderRepository.FindByNo(roNo, validateOrderRequest.CustID)
		if err != nil {
			return err
		}
	}

	if validateOrderRequest.OutletID == 0 || validateOrderRequest.WhID == 0 {
		validateOrderRequest.OutletID = *ro.OutletID
		validateOrderRequest.WhID = int(*ro.WhId)
	}

	for _, proStok := range validateOrderRequest.ProStok {
		product, err := service.OrderRepository.FindProductByID(int(proStok.ProductId))
		if err != nil {
			return err
		}

		total += (product.SellPrice1 * float64(proStok.Qty1)) + (product.SellPrice2 * float64(proStok.Qty2)) + (product.SellPrice3 * float64(proStok.Qty3))
	}
	validateOrderRequest.Total = total
	return nil
}

func (service *orderServiceImpl) ConsultDiscountBeforeStore(request *entity.ConsultDiscountOrderBody) (err error) {
	// log.Info("Debugging ConsultDiscountBeforeStore1")
	outlet, err := service.OrderRepository.FindOutletByID(int(request.OutletID), request.CustId, request.ParentCustId)
	if err != nil {
		return err
	}

	// for index := range request.Details.Normal {
	// 	log.Info("promoValueFinal : ", *request.Details.Normal[index].PromoValueFinal)
	// 	log.Info("promoValue : ", *request.Details.Normal[index].PromoValue)
	// }

	SubTotalPrincipals := map[string]int{}
	TotalProductPerDiscount := map[string]int{}
	subTotal := 0.0
	for index := range request.Details.Normal {
		// log.Info("Debugging ConsultDiscountBeforeStore2")
		product, err := service.OrderRepository.FindProductByID(request.Details.Normal[index].ProId)
		if err != nil {
			return err
		}

		request.Details.Normal[index].Qty1Final = request.Details.Normal[index].Qty1
		request.Details.Normal[index].Qty2Final = request.Details.Normal[index].Qty2
		request.Details.Normal[index].Qty3Final = request.Details.Normal[index].Qty3
		request.Details.Normal[index].Qty4Final = request.Details.Normal[index].Qty4
		request.Details.Normal[index].Qty5Final = request.Details.Normal[index].Qty5

		// log.Info("Error Disini 0")
		amount := (*request.Details.Normal[index].Qty1 * *request.Details.Normal[index].SellPrice1) + (*request.Details.Normal[index].Qty2 * *request.Details.Normal[index].SellPrice2) + (*request.Details.Normal[index].Qty3 * *request.Details.Normal[index].SellPrice3)
		// log.Info("Error Disini 1")
		subTotal += amount
		amount -= *request.Details.Normal[index].PromoValueFinal
		// log.Info("Error Disini 2")
		request.Details.Normal[index].Amount = &amount
		request.Details.Normal[index].AmountFinal = &amount
		// log.Info("Debugging ConsultDiscountBeforeStore3")
		discountID := ""
		if discount, err := service.OrderRepository.FindDiscountByProductAndOutlet(product, outlet); err == nil {
			discountID = discount.DiscountId

			if subTotalPrincipals, isExist := SubTotalPrincipals[discountID]; isExist {
				SubTotalPrincipals[discountID] = subTotalPrincipals + int(*request.Details.Normal[index].Amount)
				TotalProductPerDiscount[discountID] += 1
			} else {
				SubTotalPrincipals[discountID] = int(*request.Details.Normal[index].Amount)
				TotalProductPerDiscount[discountID] = 1
			}
		}
		request.Details.Normal[index].DiscountID = &discountID
	}
	// log.Info("Debugging ConsultDiscountBeforeStore4")

	slabRewards := map[string]int{}
	decreaseSlabRewards := map[string]int{}
	discountCriterias := map[string]model.DiscountCriteria{}
	for discountID, SubTotalPrincipal := range SubTotalPrincipals {
		if discountCriteria, err := service.OrderRepository.FindDiscountCriteriaBySubTotal(discountID, SubTotalPrincipal); err == nil {
			slabReward := discountCriteria.SlabReward
			if discountCriteria.SlabRewardType == entity.DiscountRewardTypePercentage {
				slabReward = math.Round((float64(slabReward) * float64(SubTotalPrincipal)) / 100)
			}
			slabRewards[discountID] = int(slabReward)
			decreaseSlabRewards[discountID] = slabRewards[discountID]
			discountCriterias[discountID] = discountCriteria
		}
	}
	// log.Info("slabRewards Awal : ", slabRewards)
	// log.Info("decreaseSlabRewards Awal : ", decreaseSlabRewards)
	// log.Info("discountCriterias Awal : ", discountCriterias)
	// log.Info("SubTotalPrincipals Awal : ", SubTotalPrincipals)
	// log.Info("TotalProductPerDiscount Awal : ", TotalProductPerDiscount)

	discValues := 0.0
	vatValues := 0.0
	realAmounts := 0.0
	discountAppended := map[string]bool{}
	for index := range request.Details.Normal {
		realAmount := *request.Details.Normal[index].Amount
		log.Info("Price ", request.Details.Normal[index].ProId, " : ", realAmount)
		discValue := 0.0

		if *request.Details.Normal[index].DiscountID != "" {
			if slabRewards[*request.Details.Normal[index].DiscountID] == 0 {
				*request.Details.Normal[index].DiscountID = ""
			}

			rewardProduct := int(math.Round((float64(slabRewards[*request.Details.Normal[index].DiscountID]) * realAmount) / float64(SubTotalPrincipals[*request.Details.Normal[index].DiscountID])))
			TotalProductPerDiscount[*request.Details.Normal[index].DiscountID]--
			if TotalProductPerDiscount[*request.Details.Normal[index].DiscountID] <= 0 {
				rewardProduct = decreaseSlabRewards[*request.Details.Normal[index].DiscountID]
			}
			decreaseSlabRewards[*request.Details.Normal[index].DiscountID] -= rewardProduct

			log.Info("Reward Product ", request.Details.Normal[index].ProId, " : ", rewardProduct)

			discValue = float64(rewardProduct)

			if _, isExist := discountAppended[*request.Details.Normal[index].DiscountID]; !isExist && *request.Details.Normal[index].DiscountID != "" {
				orderReward := entity.CreateOrderRewardBody{
					CustId:       request.CustId,
					ReffId:       *request.Details.Normal[index].DiscountID,
					SlabDesc:     discountCriterias[*request.Details.Normal[index].DiscountID].SlabDesc,
					RewardTypeId: 2,
				}

				request.Rewards = append(request.Rewards, orderReward)

				discountAppended[*request.Details.Normal[index].DiscountID] = true
			}
		}
		realAmount -= discValue
		request.Details.Normal[index].DiscValue = &discValue
		request.Details.Normal[index].DiscValueFinal = &discValue
		request.Details.Normal[index].DiscPo = &discValue
		discValues += discValue

		log.Info("Discount ", request.Details.Normal[index].ProId, " : ", discValue)
		log.Info("Subtotal Discount ", request.Details.Normal[index].ProId, " : ", discValues)
		log.Info("Price ", request.Details.Normal[index].ProId, " After Discount : ", realAmount)

		vatValue := math.Round((realAmount * *request.Details.Normal[index].Vat) / 100.0)
		realAmount += vatValue
		request.Details.Normal[index].VatValue = &vatValue
		request.Details.Normal[index].VatValueFinal = &vatValue
		request.Details.Normal[index].VatValuePo = &vatValue
		vatValues += vatValue

		log.Info("PPN Product ", request.Details.Normal[index].ProId, " : ", vatValue)
		log.Info("Price ", request.Details.Normal[index].ProId, " After PPN : ", realAmount)

		request.Details.Normal[index].Amount = &realAmount
		request.Details.Normal[index].AmountFinal = &realAmount
		realAmounts += realAmount

		log.Info("Real Amounts ", request.Details.Normal[index].ProId, " After Added ", realAmount, " : ", realAmounts)
	}
	log.Info("Total Discount : ", discValues)
	log.Info("Total PPN : ", vatValues)
	log.Info("Total Price : ", realAmounts)
	request.SubTotal = &subTotal
	request.SubTotalFinal = &subTotal
	request.DiscValue = &discValues
	request.DiscValueFinal = &discValues
	request.VatValue = &vatValues
	request.VatValueFinal = &vatValues

	// tes
	// realAmounts -= *request.PromoBgValueFinal

	request.Total = &realAmounts
	request.TotalFinal = &realAmounts

	return nil
}

// calculateVatValue calculates VAT value based on the formula:
// ((qty1 * price1) + (qty2 * price2) + (qty3 * price3) - promo - disc) * vat / 100
// Returns the VAT value (not added to the base amount)
func calculateVatValue(qty1, qty2, qty3, price1, price2, price3, promo, disc, vat float64) float64 {
	subtotal := (qty1 * price1) + (qty2 * price2) + (qty3 * price3)
	taxableAmount := subtotal - promo - disc
	if taxableAmount < 0 {
		taxableAmount = 0
	}
	return math.Round(taxableAmount * vat / 100.0)
}

// Tambahkan fungsi ini di tempat yang sesuai
func getValueOrDefault(value *float64, defaultValue float64) float64 {
	if value == nil {
		return defaultValue
	}
	return *value
}

type promoAggregateRow struct {
	Promo1             float64
	Promo2             float64
	Promo3             float64
	Promo4             float64
	Promo5             float64
	PromoTotal         float64
	Remarks            []string
	IsProductPromotion bool
}

type promoSnapshotTab string

const (
	promoSnapshotTabSalesOrder  promoSnapshotTab = "sales_order"
	promoSnapshotTabFinalOrder  promoSnapshotTab = "final_order"
	promoSnapshotTabPurchase    promoSnapshotTab = "purchase_order"
	PROMO_SNAPSHOT_ROLLOUT_DATE                  = "2026-03-11T00:00:00Z"
	FLOAT_COMPARE_EPSILON                        = 1e-6
)

func deterministicTabSignature(items []entity.OrderDetResponse) string {
	rows := make([]string, 0, len(items))
	for _, item := range items {
		row := strconv.Itoa(item.ProId) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.Qty1, 0), 'f', 4, 64) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.Qty2, 0), 'f', 4, 64) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.Qty3, 0), 'f', 4, 64) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.SellPrice1, 0), 'f', 4, 64) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.SellPrice2, 0), 'f', 4, 64) + "|" +
			strconv.FormatFloat(getValueOrDefault(item.SellPrice3, 0), 'f', 4, 64)
		rows = append(rows, row)
	}

	sort.Strings(rows)
	raw := strings.Join(rows, ";")
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

type consultAggregateDetail struct {
	ProID      int
	Qty1       float64
	Qty2       float64
	Qty3       float64
	GrossValue float64
}

type promoComponent struct {
	Kind   string
	Promo1 float64
	Promo2 float64
	Promo3 float64
	Promo4 float64
	Promo5 float64
}

type promoRowDistributionInfo struct {
	DetailID  int
	ProductID int
	Gross     float64
	Qty       float64
}

func aggregateConsultDetailsByProduct(details []entity.ConPromoV2Det) []entity.ConPromoV2Det {
	aggregates := make(map[int]*consultAggregateDetail)
	orderedProductIDs := make([]int, 0)

	for _, detail := range details {
		aggregate, exists := aggregates[detail.ProID]
		if !exists {
			aggregate = &consultAggregateDetail{ProID: detail.ProID}
			aggregates[detail.ProID] = aggregate
			orderedProductIDs = append(orderedProductIDs, detail.ProID)
		}

		aggregate.Qty1 += detail.Qty1
		aggregate.Qty2 += detail.Qty2
		aggregate.Qty3 += detail.Qty3
		aggregate.GrossValue += detail.Total
	}

	result := make([]entity.ConPromoV2Det, 0, len(orderedProductIDs))
	for _, productID := range orderedProductIDs {
		aggregate := aggregates[productID]
		result = append(result, entity.ConPromoV2Det{
			ProID:      aggregate.ProID,
			Qty1:       aggregate.Qty1,
			Qty2:       aggregate.Qty2,
			Qty3:       aggregate.Qty3,
			Total:      aggregate.GrossValue,
			GrossValue: int(math.Round(aggregate.GrossValue)),
		})
	}

	return result
}

func buildConsultPayloadByTab(ro model.OrderList, items []entity.OrderDetResponse, custID string, parentCustID string) entity.ConsultPromoV2Req {
	orderDate := ""
	if ro.RoDate != nil {
		orderDate = ro.RoDate.Format("2006-01-02")
	}

	rawDetails := make([]entity.ConPromoV2Det, 0, len(items))
	for _, item := range items {
		qty1 := getValueOrDefault(item.Qty1, 0)
		qty2 := getValueOrDefault(item.Qty2, 0)
		qty3 := getValueOrDefault(item.Qty3, 0)
		price1 := getValueOrDefault(item.SellPrice1, 0)
		price2 := getValueOrDefault(item.SellPrice2, 0)
		price3 := getValueOrDefault(item.SellPrice3, 0)
		gross := (qty1 * price1) + (qty2 * price2) + (qty3 * price3)

		rawDetails = append(rawDetails, entity.ConPromoV2Det{
			ProID:      item.ProId,
			Qty1:       qty1,
			Qty2:       qty2,
			Qty3:       qty3,
			Total:      gross,
			GrossValue: int(math.Round(gross)),
		})
	}

	payload := entity.ConsultPromoV2Req{
		CustID:       custID,
		ParentCustID: parentCustID,
		OrderDate:    orderDate,
		Details:      aggregateConsultDetailsByProduct(rawDetails),
	}

	if ro.OutletID != nil {
		payload.OutletID = int(*ro.OutletID)
	}
	if ro.SalesmanId != nil {
		payload.SalesmanID = int(*ro.SalesmanId)
	}
	if ro.WhId != nil {
		payload.WhID = int(*ro.WhId)
	}

	return payload
}

func determinePromoComponentsPerProduct(consultResp []entity.ConsultPromoResp) map[int][]promoComponent {
	components := make(map[int][]promoComponent)

	appendComponent := func(productID int, component promoComponent) {
		components[productID] = append(components[productID], component)
	}

	for _, promo := range consultResp {
		perScope := strings.ToLower(strings.TrimSpace(promo.SlabPerScope))
		for _, reward := range promo.RewardPercentage {
			appendComponent(reward.ProID, promoComponent{
				Kind:   "reward_percentage",
				Promo1: reward.Promo1,
				Promo2: reward.Promo2,
				Promo3: reward.Promo3,
				Promo4: reward.Promo4,
				Promo5: reward.Promo5,
			})
		}
		for _, reward := range promo.RewardValue {
			componentKind := "reward_value_per_order"
			if perScope == strings.ToLower(string(model.PerScopeProduct)) {
				componentKind = "reward_value_per_product"
			}
			appendComponent(reward.ProID, promoComponent{
				Kind:   componentKind,
				Promo1: reward.Promo1,
				Promo2: reward.Promo2,
				Promo3: reward.Promo3,
				Promo4: reward.Promo4,
				Promo5: reward.Promo5,
			})
		}
		for _, reward := range promo.RewardProduct {
			appendComponent(reward.ProID, promoComponent{
				Kind:   "reward_product",
				Promo1: reward.Promo1,
				Promo2: reward.Promo2,
				Promo3: reward.Promo3,
				Promo4: reward.Promo4,
				Promo5: reward.Promo5,
			})
		}
	}

	return components
}

func distributePromoValueByWeight(total float64, rows []promoRowDistributionInfo, weightFn func(promoRowDistributionInfo) float64) map[int]float64 {
	allocations := make(map[int]float64, len(rows))
	if len(rows) == 0 || total == 0 {
		return allocations
	}

	weights := make([]float64, len(rows))
	totalWeight := 0.0
	for i, row := range rows {
		weight := weightFn(row)
		if weight < 0 {
			weight = 0
		}
		weights[i] = weight
		totalWeight += weight
	}

	if totalWeight == 0 {
		equalShare := 0.0
		if len(rows) > 0 {
			equalShare = math.Floor((total / float64(len(rows)) * 1000000)) / 1000000
		}
		remaining := total
		for i, row := range rows {
			allocation := equalShare
			if i == len(rows)-1 {
				allocation = remaining
			}
			allocations[row.DetailID] += allocation
			remaining -= allocation
		}
		return allocations
	}

	remaining := total
	for i, row := range rows {
		allocation := 0.0
		if i == len(rows)-1 {
			allocation = remaining
		} else {
			allocation = total * (weights[i] / totalWeight)
			allocation = math.Round(allocation*1000000) / 1000000
			remaining -= allocation
		}
		allocations[row.DetailID] += allocation
	}

	return allocations
}

func applyPromoComponentAllocation(target *promoAggregateRow, component promoComponent, allocation float64) {
	if allocation == 0 {
		return
	}

	totalComponentPromo := component.Promo1 + component.Promo2 + component.Promo3 + component.Promo4 + component.Promo5
	if totalComponentPromo == 0 {
		return
	}

	promoValues := []*float64{&target.Promo1, &target.Promo2, &target.Promo3, &target.Promo4, &target.Promo5}
	componentValues := []float64{component.Promo1, component.Promo2, component.Promo3, component.Promo4, component.Promo5}
	remaining := allocation

	for i, componentValue := range componentValues {
		if componentValue == 0 {
			continue
		}

		portion := 0.0
		if i == len(componentValues)-1 {
			portion = remaining
		} else {
			portion = allocation * (componentValue / totalComponentPromo)
			portion = math.Round(portion*1000000) / 1000000
			remaining -= portion
		}
		*promoValues[i] += portion
	}
}

func buildPromoRowsByProductFromDetails(details []model.OrderDetailRead, tab promoSnapshotTab) map[int][]promoRowDistributionInfo {
	productRows := make(map[int][]promoRowDistributionInfo)

	for _, detail := range details {
		if detail.ItemType == 2 || detail.OrderDetailID == nil {
			continue
		}

		var qty1, qty2, qty3, price1, price2, price3 float64
		switch tab {
		case promoSnapshotTabFinalOrder:
			qty1 = getValueOrDefault(detail.Qty1Final, 0)
			qty2 = getValueOrDefault(detail.Qty2Final, 0)
			qty3 = getValueOrDefault(detail.Qty3Final, 0)
			price1 = getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0))
			price2 = getValueOrDefault(detail.SellPriceFinal2, getValueOrDefault(detail.SellPrice2, 0))
			price3 = getValueOrDefault(detail.SellPriceFinal3, getValueOrDefault(detail.SellPrice3, 0))
		case promoSnapshotTabPurchase:
			qty1 = getValueOrDefault(detail.QtyPo1, 0)
			qty2 = getValueOrDefault(detail.QtyPo2, 0)
			qty3 = getValueOrDefault(detail.QtyPo3, 0)
			price1 = getValueOrDefault(detail.SellPricePo1, 0)
			price2 = getValueOrDefault(detail.SellPricePo2, 0)
			price3 = getValueOrDefault(detail.SellPricePo3, 0)
		default:
			qty1 = getValueOrDefault(detail.Qty1, 0)
			qty2 = getValueOrDefault(detail.Qty2, 0)
			qty3 = getValueOrDefault(detail.Qty3, 0)
			price1 = getValueOrDefault(detail.SellPrice1, 0)
			price2 = getValueOrDefault(detail.SellPrice2, 0)
			price3 = getValueOrDefault(detail.SellPrice3, 0)
		}

		gross := (qty1 * price1) + (qty2 * price2) + (qty3 * price3)
		qty := qty1 + qty2 + qty3
		productRows[detail.ProId] = append(productRows[detail.ProId], promoRowDistributionInfo{
			DetailID:  *detail.OrderDetailID,
			ProductID: detail.ProId,
			Gross:     gross,
			Qty:       qty,
		})
	}

	return productRows
}

func buildPromoRowsByProductFromCreateDetails(details []entity.CreateOrderDetBody) map[int][]promoRowDistributionInfo {
	productRows := make(map[int][]promoRowDistributionInfo)

	for i := range details {
		gross := (getValueOrDefault(details[i].Qty1, 0) * getValueOrDefault(details[i].SellPrice1, 0)) +
			(getValueOrDefault(details[i].Qty2, 0) * getValueOrDefault(details[i].SellPrice2, 0)) +
			(getValueOrDefault(details[i].Qty3, 0) * getValueOrDefault(details[i].SellPrice3, 0))
		qty := getValueOrDefault(details[i].Qty1, 0) + getValueOrDefault(details[i].Qty2, 0) + getValueOrDefault(details[i].Qty3, 0)
		productRows[details[i].ProId] = append(productRows[details[i].ProId], promoRowDistributionInfo{
			DetailID:  i + 1,
			ProductID: details[i].ProId,
			Gross:     gross,
			Qty:       qty,
		})
	}

	return productRows
}

func distributePromoToRowsByProduct(aggregate map[int]promoAggregateRow, productRows map[int][]promoRowDistributionInfo, consultResp []entity.ConsultPromoResp) map[int]promoAggregateRow {
	result := make(map[int]promoAggregateRow)
	componentsByProduct := determinePromoComponentsPerProduct(consultResp)
	normalAggregate := aggregatePromoByProductForNormalRows(consultResp)

	for productID, rows := range productRows {
		if len(rows) == 0 {
			continue
		}

		components := componentsByProduct[productID]
		hasRewardProductComponent := false
		for _, component := range components {
			if component.Kind == "reward_product" {
				hasRewardProductComponent = true
			}
			switch component.Kind {
			case "reward_percentage":
				allocations := distributePromoValueByWeight(component.Promo1+component.Promo2+component.Promo3+component.Promo4+component.Promo5, rows, func(row promoRowDistributionInfo) float64 {
					return row.Gross
				})
				for _, row := range rows {
					rowAggregate := result[row.DetailID]
					applyPromoComponentAllocation(&rowAggregate, component, allocations[row.DetailID])
					result[row.DetailID] = rowAggregate
				}
			case "reward_value_per_order":
				allocations := distributePromoValueByWeight(component.Promo1+component.Promo2+component.Promo3+component.Promo4+component.Promo5, rows, func(row promoRowDistributionInfo) float64 {
					return 1
				})
				for _, row := range rows {
					rowAggregate := result[row.DetailID]
					applyPromoComponentAllocation(&rowAggregate, component, allocations[row.DetailID])
					result[row.DetailID] = rowAggregate
				}
			case "reward_value_per_product":
				allocations := distributePromoValueByWeight(component.Promo1+component.Promo2+component.Promo3+component.Promo4+component.Promo5, rows, func(row promoRowDistributionInfo) float64 {
					return row.Qty
				})
				for _, row := range rows {
					rowAggregate := result[row.DetailID]
					applyPromoComponentAllocation(&rowAggregate, component, allocations[row.DetailID])
					result[row.DetailID] = rowAggregate
				}
			case "reward_product":
				continue
			}
		}

		normalRow := normalAggregate[productID]
		for _, row := range rows {
			rowAggregate := result[row.DetailID]
			rowAggregate.Remarks = mergeUniqueRemarks(rowAggregate.Remarks, normalRow.Remarks)
			if normalRow.IsProductPromotion {
				rowAggregate.IsProductPromotion = true
			} else if aggregateRow, ok := aggregate[productID]; ok && aggregateRow.IsProductPromotion && !hasRewardProductComponent {
				rowAggregate.IsProductPromotion = true
			}
			rowAggregate.PromoTotal = rowAggregate.Promo1 + rowAggregate.Promo2 + rowAggregate.Promo3 + rowAggregate.Promo4 + rowAggregate.Promo5
			result[row.DetailID] = rowAggregate
		}
	}

	for productID, aggregateRow := range aggregate {
		if _, exists := productRows[productID]; exists {
			continue
		}
		_ = aggregateRow
	}

	return result
}

func distributePromoToDetailRowsV2(aggregate map[int]promoAggregateRow, details []model.OrderDetailRead, tab promoSnapshotTab, consultResp []entity.ConsultPromoResp) map[int]promoAggregateRow {
	return distributePromoToRowsByProduct(aggregate, buildPromoRowsByProductFromDetails(details, tab), consultResp)
}

func orchestratePromoConsultByTabs(
	payloads map[string]entity.ConsultPromoV2Req,
	signatures map[string]string,
	consultFn func(entity.ConsultPromoV2Req) ([]entity.ConsultPromoResp, error),
) map[string][]entity.ConsultPromoResp {
	result := map[string][]entity.ConsultPromoResp{
		"normal":   {},
		"final":    {},
		"purchase": {},
	}

	tabs := []string{"normal", "final", "purchase"}
	normalSig := signatures["normal"]
	finalSig := signatures["final"]
	purchaseSig := signatures["purchase"]

	if normalSig == finalSig && finalSig == purchaseSig {
		resp, err := consultFn(payloads["normal"])
		if err == nil {
			result["normal"] = resp
			result["final"] = resp
			result["purchase"] = resp
			return result
		}
		// Fallback: hit tiap tab agar partial success tetap terisi
	}

	cacheBySignature := make(map[string][]entity.ConsultPromoResp)
	for _, tab := range tabs {
		sig := signatures[tab]
		if cached, ok := cacheBySignature[sig]; ok {
			result[tab] = cached
			continue
		}

		resp, err := consultFn(payloads[tab])
		if err != nil {
			result[tab] = []entity.ConsultPromoResp{}
			cacheBySignature[sig] = []entity.ConsultPromoResp{}
			continue
		}
		result[tab] = resp
		cacheBySignature[sig] = resp
	}

	return result
}

func aggregatePromoByProduct(consultResp []entity.ConsultPromoResp) map[int]promoAggregateRow {
	return aggregatePromoByProductWithEligibleScope(consultResp, false)
}

func aggregatePromoByProductForDetailSnapshot(consultResp []entity.ConsultPromoResp) map[int]promoAggregateRow {
	return aggregatePromoByProductWithEligibleScope(consultResp, true)
}

func aggregatePromoByProductForNormalRows(consultResp []entity.ConsultPromoResp) map[int]promoAggregateRow {
	result := make(map[int]promoAggregateRow)

	appendRemark := func(row promoAggregateRow, promoID string) promoAggregateRow {
		if promoID == "" {
			return row
		}
		if !slices.Contains(row.Remarks, promoID) {
			row.Remarks = append(row.Remarks, promoID)
			sort.Strings(row.Remarks)
		}
		return row
	}

	for _, promo := range consultResp {
		attachEligibleRemarks := len(promo.RewardProduct) == 0
		eligibleSet := make(map[int]struct{}, len(promo.ProductsEligible))
		for _, eligibleProID := range promo.ProductsEligible {
			eligibleSet[eligibleProID] = struct{}{}
			row := result[eligibleProID]
			if attachEligibleRemarks {
				row = appendRemark(row, promo.PromoID)
			}
			result[eligibleProID] = row
		}

		for _, reward := range promo.RewardPercentage {
			if _, ok := eligibleSet[reward.ProID]; !ok {
				continue
			}
			row := result[reward.ProID]
			row.Promo1 += reward.Promo1
			row.Promo2 += reward.Promo2
			row.Promo3 += reward.Promo3
			row.Promo4 += reward.Promo4
			row.Promo5 += reward.Promo5
			row.PromoTotal = row.Promo1 + row.Promo2 + row.Promo3 + row.Promo4 + row.Promo5
			result[reward.ProID] = row
		}

		for _, reward := range promo.RewardValue {
			if _, ok := eligibleSet[reward.ProID]; !ok {
				continue
			}
			row := result[reward.ProID]
			row.Promo1 += reward.Promo1
			row.Promo2 += reward.Promo2
			row.Promo3 += reward.Promo3
			row.Promo4 += reward.Promo4
			row.Promo5 += reward.Promo5
			row.PromoTotal = row.Promo1 + row.Promo2 + row.Promo3 + row.Promo4 + row.Promo5
			result[reward.ProID] = row
		}
	}

	return result
}

func aggregatePromoByProductWithEligibleScope(consultResp []entity.ConsultPromoResp, restrictRewardToEligible bool) map[int]promoAggregateRow {
	result := make(map[int]promoAggregateRow)

	appendRemark := func(row promoAggregateRow, promoID string) promoAggregateRow {
		if promoID == "" {
			return row
		}
		if !slices.Contains(row.Remarks, promoID) {
			row.Remarks = append(row.Remarks, promoID)
			sort.Strings(row.Remarks)
		}
		return row
	}

	for _, promo := range consultResp {
		attachEligibleRemarks := len(promo.RewardProduct) == 0
		eligibleSet := make(map[int]struct{}, len(promo.ProductsEligible))
		for _, eligibleProID := range promo.ProductsEligible {
			eligibleSet[eligibleProID] = struct{}{}
			row := result[eligibleProID]
			if attachEligibleRemarks {
				row = appendRemark(row, promo.PromoID)
			}
			result[eligibleProID] = row
		}

		for _, reward := range promo.RewardPercentage {
			if restrictRewardToEligible {
				if _, ok := eligibleSet[reward.ProID]; !ok {
					continue
				}
			}
			row := result[reward.ProID]
			row.Promo1 += reward.Promo1
			row.Promo2 += reward.Promo2
			row.Promo3 += reward.Promo3
			row.Promo4 += reward.Promo4
			row.Promo5 += reward.Promo5
			row.PromoTotal = row.Promo1 + row.Promo2 + row.Promo3 + row.Promo4 + row.Promo5
			result[reward.ProID] = row
		}

		for _, reward := range promo.RewardValue {
			if restrictRewardToEligible {
				if _, ok := eligibleSet[reward.ProID]; !ok {
					continue
				}
			}
			row := result[reward.ProID]
			row.Promo1 += reward.Promo1
			row.Promo2 += reward.Promo2
			row.Promo3 += reward.Promo3
			row.Promo4 += reward.Promo4
			row.Promo5 += reward.Promo5
			row.PromoTotal = row.Promo1 + row.Promo2 + row.Promo3 + row.Promo4 + row.Promo5
			result[reward.ProID] = row
		}

		for _, reward := range promo.RewardProduct {
			if restrictRewardToEligible {
				if _, ok := eligibleSet[reward.ProID]; !ok {
					continue
				}
			}
			row := result[reward.ProID]
			row.Promo1 += reward.Promo1
			row.Promo2 += reward.Promo2
			row.Promo3 += reward.Promo3
			row.Promo4 += reward.Promo4
			row.Promo5 += reward.Promo5
			row.PromoTotal = row.Promo1 + row.Promo2 + row.Promo3 + row.Promo4 + row.Promo5
			row.IsProductPromotion = true
			row = appendRemark(row, promo.PromoID)
			result[reward.ProID] = row
		}
	}

	return result
}

func applyPersistedPromoSnapshotToItems(items []entity.OrderDetResponse, tab promoSnapshotTab) []entity.OrderDetResponse {
	mapped := make([]entity.OrderDetResponse, len(items))
	copy(mapped, items)

	for i := range mapped {
		var (
			promo1             float64
			promo2             float64
			promo3             float64
			promo4             float64
			promo5             float64
			remarks            []string
			isProductPromotion bool
		)

		switch tab {
		case promoSnapshotTabSalesOrder:
			promo1 = mapped[i].PromoSo1
			promo2 = mapped[i].PromoSo2
			promo3 = mapped[i].PromoSo3
			promo4 = mapped[i].PromoSo4
			promo5 = mapped[i].PromoSo5
			remarks = append([]string{}, mapped[i].PromoRemarksSo...)
			isProductPromotion = mapped[i].IsProductPromotionSo
		case promoSnapshotTabFinalOrder:
			promo1 = mapped[i].PromoFinal1
			promo2 = mapped[i].PromoFinal2
			promo3 = mapped[i].PromoFinal3
			promo4 = mapped[i].PromoFinal4
			promo5 = mapped[i].PromoFinal5
			remarks = append([]string{}, mapped[i].PromoRemarksFinal...)
			isProductPromotion = mapped[i].IsProductPromotionFinal
		case promoSnapshotTabPurchase:
			promo1 = mapped[i].PromoPo1
			promo2 = mapped[i].PromoPo2
			promo3 = mapped[i].PromoPo3
			promo4 = mapped[i].PromoPo4
			promo5 = mapped[i].PromoPo5
			remarks = append([]string{}, mapped[i].PromoRemarksPo...)
			isProductPromotion = mapped[i].IsProductPromotionPo
		}

		sort.Strings(remarks)
		mapped[i].Promo1 = promo1
		mapped[i].Promo2 = promo2
		mapped[i].Promo3 = promo3
		mapped[i].Promo4 = promo4
		mapped[i].Promo5 = promo5
		mapped[i].PromoTotal = promo1 + promo2 + promo3 + promo4 + promo5
		mapped[i].Remarks = remarks
		mapped[i].IsProductPromotion = isProductPromotion
	}

	return mapped
}

func hasPersistedPromoSnapshot(headerRemarks []string, items []entity.OrderDetResponse, tab promoSnapshotTab) bool {
	if len(headerRemarks) > 0 {
		return true
	}

	for _, item := range items {
		switch tab {
		case promoSnapshotTabSalesOrder:
			if item.PromoSo1 != 0 || item.PromoSo2 != 0 || item.PromoSo3 != 0 || item.PromoSo4 != 0 || item.PromoSo5 != 0 || len(item.PromoRemarksSo) > 0 || item.IsProductPromotionSo {
				return true
			}
		case promoSnapshotTabFinalOrder:
			if item.PromoFinal1 != 0 || item.PromoFinal2 != 0 || item.PromoFinal3 != 0 || item.PromoFinal4 != 0 || item.PromoFinal5 != 0 || len(item.PromoRemarksFinal) > 0 || item.IsProductPromotionFinal {
				return true
			}
		case promoSnapshotTabPurchase:
			if item.PromoPo1 != 0 || item.PromoPo2 != 0 || item.PromoPo3 != 0 || item.PromoPo4 != 0 || item.PromoPo5 != 0 || len(item.PromoRemarksPo) > 0 || item.IsProductPromotionPo {
				return true
			}
		}
	}

	return false
}

func buildDetailPromoSnapshotUpdates(detail model.OrderDetailRead, aggregate map[int]promoAggregateRow, tab promoSnapshotTab) map[string]interface{} {
	row := promoAggregateRow{}
	if detail.OrderDetailID != nil {
		row = aggregate[*detail.OrderDetailID]
	}
	updates := map[string]interface{}{}

	switch tab {
	case promoSnapshotTabSalesOrder:
		updates["promo_so1"] = row.Promo1
		updates["promo_so2"] = row.Promo2
		updates["promo_so3"] = row.Promo3
		updates["promo_so4"] = row.Promo4
		updates["promo_so5"] = row.Promo5
		updates["promo_remarks_so"] = model.JSONStringArray(append([]string{}, row.Remarks...))
		updates["is_product_promotion_so"] = row.IsProductPromotion
	case promoSnapshotTabFinalOrder:
		updates["promo_final1"] = row.Promo1
		updates["promo_final2"] = row.Promo2
		updates["promo_final3"] = row.Promo3
		updates["promo_final4"] = row.Promo4
		updates["promo_final5"] = row.Promo5
		updates["promo_remarks_final"] = model.JSONStringArray(append([]string{}, row.Remarks...))
		updates["is_product_promotion_final"] = row.IsProductPromotion
	case promoSnapshotTabPurchase:
		updates["promo_po1"] = row.Promo1
		updates["promo_po2"] = row.Promo2
		updates["promo_po3"] = row.Promo3
		updates["promo_po4"] = row.Promo4
		updates["promo_po5"] = row.Promo5
		updates["promo_remarks_po"] = model.JSONStringArray(append([]string{}, row.Remarks...))
		updates["is_product_promotion_po"] = row.IsProductPromotion
	}

	return updates
}

func buildHeaderPromoSnapshotUpdate(remarks []string, tab promoSnapshotTab) map[string]interface{} {
	updates := map[string]interface{}{}
	copied := append([]string{}, remarks...)
	sort.Strings(copied)
	jsonRemarks := model.JSONStringArray(copied)

	switch tab {
	case promoSnapshotTabSalesOrder:
		updates["promo_remarks_so"] = jsonRemarks
	case promoSnapshotTabFinalOrder:
		updates["promo_remarks_final"] = jsonRemarks
	case promoSnapshotTabPurchase:
		updates["promo_remarks_po"] = jsonRemarks
	}

	return updates
}

func buildRewardProductStockDeltas(custID, roNo string, whID int64, roDate time.Time, existingRewards []model.OrderDetailRead, newRewards []entity.OrderRewardProductResponse) []*entity.SalesOrderStockUpdate {
	updates := make([]*entity.SalesOrderStockUpdate, 0, len(existingRewards)+len(newRewards))
	trCode := ""
	if len(roNo) >= 2 {
		trCode = roNo[0:2]
	}

	for _, detail := range existingRewards {
		if detail.ItemType != 2 || detail.OrderDetailID == nil {
			continue
		}

		qtyBefore := getValueOrDefault(detail.QtyFinal, 0)
		if qtyBefore == 0 {
			continue
		}

		updates = append(updates, &entity.SalesOrderStockUpdate{
			CustID:         custID,
			WhID:           whID,
			ProID:          int64(detail.ProId),
			StockDate:      roDate,
			TrCode:         trCode,
			TrNo:           roNo,
			QtyOrderBefore: float64Ptr(qtyBefore),
			QtyOrder:       0,
			UnitPrice:      getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0)),
			RefDetId:       int64(*detail.OrderDetailID),
		})
	}

	for _, reward := range newRewards {
		qtyOrder := reward.Qty1 + reward.Qty2 + reward.Qty3
		if qtyOrder == 0 {
			continue
		}

		updates = append(updates, &entity.SalesOrderStockUpdate{
			CustID:         custID,
			WhID:           whID,
			ProID:          int64(reward.ProID),
			StockDate:      roDate,
			TrCode:         trCode,
			TrNo:           roNo,
			QtyOrderBefore: nil,
			QtyOrder:       qtyOrder,
			UnitPrice:      reward.SellPrice1,
		})
	}

	sort.Slice(updates, func(i, j int) bool {
		if updates[i].ProID == updates[j].ProID {
			return updates[i].RefDetId < updates[j].RefDetId
		}
		return updates[i].ProID < updates[j].ProID
	})

	return updates
}

func injectPromoToOrderItems(items []entity.OrderDetResponse, promoMap map[int]promoAggregateRow) []entity.OrderDetResponse {
	for i := range items {
		items[i].Promo1 = 0
		items[i].Promo2 = 0
		items[i].Promo3 = 0
		items[i].Promo4 = 0
		items[i].Promo5 = 0
		items[i].PromoTotal = 0
		items[i].Remarks = []string{}
		items[i].IsProductPromotion = false

		if items[i].OrderDetId == 0 {
			continue
		}
		promoData, ok := promoMap[int(items[i].OrderDetId)]
		if !ok {
			continue
		}

		items[i].Promo1 = promoData.Promo1
		items[i].Promo2 = promoData.Promo2
		items[i].Promo3 = promoData.Promo3
		items[i].Promo4 = promoData.Promo4
		items[i].Promo5 = promoData.Promo5
		items[i].PromoTotal = promoData.PromoTotal
		items[i].Remarks = append([]string{}, promoData.Remarks...)
		items[i].IsProductPromotion = promoData.IsProductPromotion
		sort.Strings(items[i].Remarks)
	}

	return items
}

func buildFinalRemarks(consultResp []entity.ConsultPromoResp) []string {
	remarks := make([]string, 0)
	for _, promo := range consultResp {
		if promo.PromoID == "" {
			continue
		}
		if !slices.Contains(remarks, promo.PromoID) {
			remarks = append(remarks, promo.PromoID)
		}
	}
	sort.Strings(remarks)
	return remarks
}

func buildProductMetaMap(items []entity.OrderDetResponse) map[int]entity.OrderDetResponse {
	result := make(map[int]entity.OrderDetResponse)
	for _, item := range items {
		if _, exists := result[item.ProId]; !exists {
			result[item.ProId] = item
		}
	}
	return result
}

func buildProductMasterMap(products []model.Product) map[int]model.Product {
	result := make(map[int]model.Product, len(products))
	for _, product := range products {
		result[int(product.ProductId)] = product
	}
	return result
}

func collectPersistedRewardProductIDsForFallback(details []model.OrderDetailRead) []int64 {
	productIDs := make([]int64, 0)
	seen := make(map[int64]struct{})

	for _, detail := range details {
		if detail.ItemType != 2 || !isEmptyStringPtr(detail.UnitId1) {
			continue
		}

		productID := int64(detail.ProId)
		if _, exists := seen[productID]; exists {
			continue
		}

		seen[productID] = struct{}{}
		productIDs = append(productIDs, productID)
	}

	return productIDs
}

func collectPromoItemProductIDsForFallback(items []entity.OrderDetResponse) []int64 {
	productIDs := make([]int64, 0)
	seen := make(map[int64]struct{})

	for _, item := range items {
		if !isEmptyStringPtr(item.UnitId1) {
			continue
		}

		productID := int64(item.ProId)
		if _, exists := seen[productID]; exists {
			continue
		}

		seen[productID] = struct{}{}
		productIDs = append(productIDs, productID)
	}

	return productIDs
}

func buildCreateOrderPromoPayload(request entity.CreateOrderBody) entity.ConsultPromoV2Req {
	payload := entity.ConsultPromoV2Req{
		CustID:       request.CustId,
		ParentCustID: request.ParentCustId,
		OutletID:     int(request.OutletID),
		SalesmanID:   int(request.SalesmanId),
	}

	if request.WhId != nil {
		payload.WhID = int(*request.WhId)
	}
	if request.RoDate != nil {
		if roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate); err == nil {
			payload.OrderDate = roDate.Format("2006-01-02")
		} else if len(*request.RoDate) >= 10 {
			payload.OrderDate = (*request.RoDate)[:10]
		} else {
			payload.OrderDate = *request.RoDate
		}
	}

	rawDetails := make([]entity.ConPromoV2Det, 0, len(request.Details.Normal))
	for _, detail := range request.Details.Normal {
		qty1 := getValueOrDefault(detail.Qty1, 0)
		qty2 := getValueOrDefault(detail.Qty2, 0)
		qty3 := getValueOrDefault(detail.Qty3, 0)
		price1 := getValueOrDefault(detail.SellPrice1, 0)
		price2 := getValueOrDefault(detail.SellPrice2, 0)
		price3 := getValueOrDefault(detail.SellPrice3, 0)
		gross := (qty1 * price1) + (qty2 * price2) + (qty3 * price3)

		rawDetails = append(rawDetails, entity.ConPromoV2Det{
			ProID:      detail.ProId,
			Qty1:       qty1,
			Qty2:       qty2,
			Qty3:       qty3,
			Total:      gross,
			GrossValue: int(math.Round(gross)),
		})
	}
	payload.Details = aggregateConsultDetailsByProduct(rawDetails)

	return payload
}

func buildRewardProducts(
	consultResp []entity.ConsultPromoResp,
	productMeta map[int]entity.OrderDetResponse,
	findProductFn func(int) (model.ProductRead, error),
) []entity.OrderRewardProductResponse {
	rewardProducts := make([]entity.OrderRewardProductResponse, 0)
	for _, promo := range consultResp {
		for _, rewardProduct := range promo.RewardProduct {
			meta, ok := productMeta[rewardProduct.ProID]
			if !ok && findProductFn != nil {
				if product, err := findProductFn(rewardProduct.ProID); err == nil {
					meta = entity.OrderDetResponse{
						ProId:      product.ProId,
						ProCode:    product.ProCode,
						ProName:    product.ProName,
						UnitId1:    stringPtr(product.UnitId1),
						UnitId2:    stringPtr(product.UnitId2),
						UnitId3:    stringPtr(product.UnitId3),
						UnitId4:    stringPtr(product.UnitId4),
						UnitId5:    stringPtr(product.UnitId5),
						SellPrice1: &product.SellPrice1,
						SellPrice2: &product.SellPrice2,
						SellPrice3: &product.SellPrice3,
					}
				}
			}

			rewardProducts = append(rewardProducts, entity.OrderRewardProductResponse{
				ProID:      rewardProduct.ProID,
				ProCode:    meta.ProCode,
				ProName:    meta.ProName,
				UnitId1:    meta.UnitId1,
				UnitId2:    meta.UnitId2,
				UnitId3:    meta.UnitId3,
				UnitId4:    meta.UnitId4,
				UnitId5:    meta.UnitId5,
				SellPrice1: getValueOrDefault(meta.SellPrice1, 0),
				SellPrice2: getValueOrDefault(meta.SellPrice2, 0),
				SellPrice3: getValueOrDefault(meta.SellPrice3, 0),
				Qty1:       rewardProduct.Qty1,
				Qty2:       rewardProduct.Qty2,
				Qty3:       rewardProduct.Qty3,
				GrossValue: rewardProduct.GrossValue,
				Promo1:     rewardProduct.Promo1,
				Promo2:     rewardProduct.Promo2,
				Promo3:     rewardProduct.Promo3,
				Promo4:     rewardProduct.Promo4,
				Promo5:     rewardProduct.Promo5,
			})
		}
	}

	return rewardProducts
}

func buildCreateOrderRewardsFromPromoV2(custID string, consultResp []entity.ConsultPromoResp) []entity.CreateOrderRewardBody {
	rewards := make([]entity.CreateOrderRewardBody, 0, len(consultResp))
	for _, promo := range consultResp {
		if promo.PromoID == "" {
			continue
		}
		rewards = append(rewards, entity.CreateOrderRewardBody{
			CustId:       custID,
			ReffId:       promo.PromoID,
			SlabId:       0,
			SlabDesc:     promo.SlabDesc,
			RewardTypeId: 1,
		})
	}
	return rewards
}

func buildCreateOrderRewardDetails(consultResp []entity.ConsultPromoResp, findProductFn func(int) (model.ProductRead, error)) ([]entity.CreateOrderDetBody, float64, error) {
	details := make([]entity.CreateOrderDetBody, 0)
	promoBgTotal := 0.0

	for _, promo := range consultResp {
		for _, rewardProduct := range promo.RewardProduct {
			product, err := findProductFn(rewardProduct.ProID)
			if err != nil {
				return nil, 0, err
			}

			qty1 := rewardProduct.Qty1
			qty2 := rewardProduct.Qty2
			qty3 := rewardProduct.Qty3
			qty4 := 0.0
			qty5 := 0.0
			amount := 0.0
			discValue := 0.0
			vatValue := 0.0
			promoValue := rewardProduct.GrossValue
			promoSo1 := rewardProduct.Promo1
			promoSo2 := rewardProduct.Promo2
			promoSo3 := rewardProduct.Promo3
			promoSo4 := rewardProduct.Promo4
			promoSo5 := rewardProduct.Promo5
			promoFinal1 := rewardProduct.Promo1
			promoFinal2 := rewardProduct.Promo2
			promoFinal3 := rewardProduct.Promo3
			promoFinal4 := rewardProduct.Promo4
			promoFinal5 := rewardProduct.Promo5
			promoRemarksSo := []string{promo.PromoID}
			promoRemarksFinal := []string{promo.PromoID}
			isProductPromotionSo := true
			isProductPromotionFinal := true
			promoID := promo.PromoID
			itemType := 2
			convUnit2 := int(product.ConvUnit2)
			convUnit3 := int(product.ConvUnit3)

			details = append(details, entity.CreateOrderDetBody{
				ProId:                   rewardProduct.ProID,
				ItemType:                itemType,
				PromoID:                 &promoID,
				Qty1:                    &qty1,
				Qty2:                    &qty2,
				Qty3:                    &qty3,
				Qty4:                    &qty4,
				Qty5:                    &qty5,
				Qty1Final:               &qty1,
				Qty2Final:               &qty2,
				Qty3Final:               &qty3,
				SellPrice1:              &product.SellPrice1,
				SellPrice2:              &product.SellPrice2,
				SellPrice3:              &product.SellPrice3,
				PromoValue:              &promoValue,
				PromoValueFinal:         &promoValue,
				PromoSo1:                &promoSo1,
				PromoSo2:                &promoSo2,
				PromoSo3:                &promoSo3,
				PromoSo4:                &promoSo4,
				PromoSo5:                &promoSo5,
				PromoFinal1:             &promoFinal1,
				PromoFinal2:             &promoFinal2,
				PromoFinal3:             &promoFinal3,
				PromoFinal4:             &promoFinal4,
				PromoFinal5:             &promoFinal5,
				PromoRemarksSo:          promoRemarksSo,
				PromoRemarksFinal:       promoRemarksFinal,
				IsProductPromotionSo:    &isProductPromotionSo,
				IsProductPromotionFinal: &isProductPromotionFinal,
				DiscValue:               &discValue,
				DiscValueFinal:          &discValue,
				Vat:                     &vatValue,
				VatValue:                &vatValue,
				VatValueFinal:           &vatValue,
				Amount:                  &amount,
				AmountFinal:             &amount,
				ConvUnit2:               &convUnit2,
				ConvUnit3:               &convUnit3,
			})
			promoBgTotal += rewardProduct.GrossValue
		}
	}

	return details, promoBgTotal, nil
}

func buildRewardOrderDetailModels(custID string, roNo string, rewardDetails []entity.CreateOrderDetBody) ([]model.OrderDetail, error) {
	rewardModels := make([]model.OrderDetail, 0, len(rewardDetails))

	for _, rewardDetail := range rewardDetails {
		var rewardModel model.OrderDetail
		rewardModel.CustId = custID
		rewardModel.RoNo = roNo
		rewardModel.ItemType = 2

		qtyUnit := &conversion.QtyUnit{
			Qty1:      int(getValueOrDefault(rewardDetail.Qty1, 0)),
			Qty2:      int(getValueOrDefault(rewardDetail.Qty2, 0)),
			Qty3:      int(getValueOrDefault(rewardDetail.Qty3, 0)),
			ConvUnit2: int(*rewardDetail.ConvUnit2),
			ConvUnit3: int(*rewardDetail.ConvUnit3),
		}

		totalQty, err := qtyUnit.ToTotalQuantity()
		if err != nil {
			return nil, err
		}

		rewardModel.Qty = float64(totalQty)
		rewardModel.QtyPo = float64(totalQty)
		rewardModel.QtyFinal = float64(totalQty)

		if err := structs.Automapper(rewardDetail, &rewardModel); err != nil {
			return nil, err
		}

		rewardModel.CustId = custID
		rewardModel.RoNo = roNo
		rewardModel.ItemType = 2
		rewardModels = append(rewardModels, rewardModel)
	}

	return rewardModels, nil
}

func buildRewardProductsFromPersistedDetails(details []model.OrderDetailRead, tab promoSnapshotTab, productMap map[int]model.Product) []entity.OrderRewardProductResponse {
	rewardProducts := make([]entity.OrderRewardProductResponse, 0)

	for _, detail := range details {
		if detail.ItemType != 2 {
			continue
		}

		reward := entity.OrderRewardProductResponse{
			ProID:   detail.ProId,
			ProCode: detail.ProCode,
			ProName: detail.ProName,
			UnitId1: detail.UnitId1,
			UnitId2: detail.UnitId2,
			UnitId3: detail.UnitId3,
			UnitId4: detail.UnitId4,
			UnitId5: detail.UnitId5,
		}

		// Fallback to product master if persisted unit fields are empty
		if isEmptyStringPtr(reward.UnitId1) {
			if product, exists := productMap[detail.ProId]; exists {
				reward.UnitId1 = stringPtr(product.UnitId1)
				reward.UnitId2 = stringPtr(product.UnitId2)
				reward.UnitId3 = stringPtr(product.UnitId3)
				reward.UnitId4 = product.UnitId4
				reward.UnitId5 = product.UnitId5
			}
		}

		switch tab {
		case promoSnapshotTabFinalOrder:
			reward.Qty1 = getValueOrDefault(detail.Qty1Final, 0)
			reward.Qty2 = getValueOrDefault(detail.Qty2Final, 0)
			reward.Qty3 = getValueOrDefault(detail.Qty3Final, 0)
			reward.SellPrice1 = getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0))
			reward.SellPrice2 = getValueOrDefault(detail.SellPriceFinal2, getValueOrDefault(detail.SellPrice2, 0))
			reward.SellPrice3 = getValueOrDefault(detail.SellPriceFinal3, getValueOrDefault(detail.SellPrice3, 0))
			reward.Promo1 = getValueOrDefault(detail.PromoFinal1, 0)
			reward.Promo2 = getValueOrDefault(detail.PromoFinal2, 0)
			reward.Promo3 = getValueOrDefault(detail.PromoFinal3, 0)
			reward.Promo4 = getValueOrDefault(detail.PromoFinal4, 0)
			reward.Promo5 = getValueOrDefault(detail.PromoFinal5, 0)
		case promoSnapshotTabPurchase:
			reward.Qty1 = getValueOrDefault(detail.QtyPo1, 0)
			reward.Qty2 = getValueOrDefault(detail.QtyPo2, 0)
			reward.Qty3 = getValueOrDefault(detail.QtyPo3, 0)
			reward.SellPrice1 = getValueOrDefault(detail.SellPricePo1, 0)
			reward.SellPrice2 = getValueOrDefault(detail.SellPricePo2, 0)
			reward.SellPrice3 = getValueOrDefault(detail.SellPricePo3, 0)
			reward.Promo1 = getValueOrDefault(detail.PromoPo1, 0)
			reward.Promo2 = getValueOrDefault(detail.PromoPo2, 0)
			reward.Promo3 = getValueOrDefault(detail.PromoPo3, 0)
			reward.Promo4 = getValueOrDefault(detail.PromoPo4, 0)
			reward.Promo5 = getValueOrDefault(detail.PromoPo5, 0)
		default:
			reward.Qty1 = getValueOrDefault(detail.Qty1, 0)
			reward.Qty2 = getValueOrDefault(detail.Qty2, 0)
			reward.Qty3 = getValueOrDefault(detail.Qty3, 0)
			reward.SellPrice1 = getValueOrDefault(detail.SellPrice1, 0)
			reward.SellPrice2 = getValueOrDefault(detail.SellPrice2, 0)
			reward.SellPrice3 = getValueOrDefault(detail.SellPrice3, 0)
			reward.Promo1 = getValueOrDefault(detail.PromoSo1, 0)
			reward.Promo2 = getValueOrDefault(detail.PromoSo2, 0)
			reward.Promo3 = getValueOrDefault(detail.PromoSo3, 0)
			reward.Promo4 = getValueOrDefault(detail.PromoSo4, 0)
			reward.Promo5 = getValueOrDefault(detail.PromoSo5, 0)
		}

		reward.GrossValue = (reward.Qty1 * reward.SellPrice1) + (reward.Qty2 * reward.SellPrice2) + (reward.Qty3 * reward.SellPrice3)
		rewardProducts = append(rewardProducts, reward)
	}

	return rewardProducts
}

func buildRewardProductStockDeltasFromModels(custID, roNo string, whID int64, roDate time.Time, existingRewards []model.OrderDetailRead, newRewardDetails []model.OrderDetail) []*entity.SalesOrderStockUpdate {
	updates := make([]*entity.SalesOrderStockUpdate, 0, len(existingRewards)+len(newRewardDetails))
	trCode := ""
	if len(roNo) >= 2 {
		trCode = roNo[0:2]
	}

	for _, detail := range existingRewards {
		if detail.ItemType != 2 || detail.OrderDetailID == nil {
			continue
		}

		qtyBefore := getValueOrDefault(detail.QtyFinal, 0)
		if qtyBefore == 0 {
			continue
		}

		updates = append(updates, &entity.SalesOrderStockUpdate{
			CustID:         custID,
			WhID:           whID,
			ProID:          int64(detail.ProId),
			StockDate:      roDate,
			TrCode:         trCode,
			TrNo:           roNo,
			QtyOrderBefore: float64Ptr(qtyBefore),
			QtyOrder:       0,
			UnitPrice:      getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0)),
			RefDetId:       int64(*detail.OrderDetailID),
		})
	}

	for _, detail := range newRewardDetails {
		if detail.ItemType != 2 {
			continue
		}

		qtyOrder := detail.QtyFinal
		if qtyOrder == 0 {
			continue
		}

		refDetID := int64(0)
		if detail.OrderDetailID != nil {
			refDetID = int64(*detail.OrderDetailID)
		}

		updates = append(updates, &entity.SalesOrderStockUpdate{
			CustID:         custID,
			WhID:           whID,
			ProID:          int64(detail.ProId),
			StockDate:      roDate,
			TrCode:         trCode,
			TrNo:           roNo,
			QtyOrderBefore: nil,
			QtyOrder:       qtyOrder,
			UnitPrice:      getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0)),
			RefDetId:       refDetID,
		})
	}

	sort.Slice(updates, func(i, j int) bool {
		if updates[i].ProID == updates[j].ProID {
			return updates[i].RefDetId < updates[j].RefDetId
		}
		return updates[i].ProID < updates[j].ProID
	})

	return updates
}

func (service *orderServiceImpl) syncRewardProductState(ctx context.Context, ro model.OrderList, roNo string, custID string, parentCustID string, details []model.OrderDetailRead, tab promoSnapshotTab) ([]*entity.SalesOrderStockUpdate, error) {
	if service.PromotionRepository == nil || service.PromotionV2Repository == nil || ro.WhId == nil || ro.RoDate == nil {
		return nil, nil
	}

	promoService := NewPromotionService(service.PromotionRepository, service.PromotionV2Repository, service.Transaction)
	consultResp, err := promoService.ConsultV2(buildConsultPayloadByTab(ro, buildOrderItemsForPromoTab(details, tab), custID, parentCustID))
	if err != nil {
		return nil, err
	}

	existingRewards := make([]model.OrderDetailRead, 0)
	for _, detail := range details {
		if detail.ItemType == 2 {
			existingRewards = append(existingRewards, detail)
		}
	}

	if err := service.OrderRepository.DeletePromoDetails(ctx, roNo, custID); err != nil {
		return nil, err
	}

	rewardDetails, _, err := buildCreateOrderRewardDetails(consultResp, service.OrderRepository.FindProductByID)
	if err != nil {
		return nil, err
	}

	rewardModels, err := buildRewardOrderDetailModels(custID, roNo, rewardDetails)
	if err != nil {
		return nil, err
	}

	for i := range rewardModels {
		defaultFalse := false
		rewardModels[i].IsProductPromotionPo = &defaultFalse
		if err := service.OrderRepository.StoreDetail(ctx, &rewardModels[i]); err != nil {
			return nil, err
		}
	}

	return buildRewardProductStockDeltasFromModels(custID, roNo, int64(*ro.WhId), *ro.RoDate, existingRewards, rewardModels), nil
}

func (service *orderServiceImpl) prepareCreateOrderPromoState(request *entity.CreateOrderBody) ([]entity.ConsultPromoResp, map[int]promoAggregateRow, []string, error) {
	if service.PromotionRepository == nil || service.PromotionV2Repository == nil {
		return nil, map[int]promoAggregateRow{}, nil, nil
	}

	promoService := NewPromotionService(service.PromotionRepository, service.PromotionV2Repository, service.Transaction)
	consultResp, err := promoService.ConsultV2(buildCreateOrderPromoPayload(*request))
	if err != nil {
		return nil, nil, nil, err
	}

	productAggregate := aggregatePromoByProductForNormalRows(consultResp)
	remarks := buildFinalRemarks(consultResp)
	promoTotal := 0.0
	rowAggregate := distributePromoToRowsByProduct(productAggregate, buildPromoRowsByProductFromCreateDetails(request.Details.Normal), consultResp)

	for i := range request.Details.Normal {
		promoValue := rowAggregate[i+1].PromoTotal
		request.Details.Normal[i].PromoValue = float64Ptr(promoValue)
		request.Details.Normal[i].PromoValueFinal = float64Ptr(promoValue)
		promoTotal += promoValue
	}
	request.PromoValue = float64Ptr(promoTotal)
	request.PromoValueFinal = float64Ptr(promoTotal)

	rewardDetails, promoBgTotal, err := buildCreateOrderRewardDetails(consultResp, service.OrderRepository.FindProductByID)
	if err != nil {
		return nil, nil, nil, err
	}
	request.Details.Promo = rewardDetails
	request.PromoBgValue = float64Ptr(promoBgTotal)
	request.PromoBgValueFinal = float64Ptr(promoBgTotal)

	return consultResp, rowAggregate, remarks, nil
}

func activeQtyForTab(detail model.OrderDetailRead, tab promoSnapshotTab) float64 {
	switch tab {
	case promoSnapshotTabFinalOrder:
		return getValueOrDefault(detail.Qty1Final, 0) + getValueOrDefault(detail.Qty2Final, 0) + getValueOrDefault(detail.Qty3Final, 0)
	case promoSnapshotTabPurchase:
		return getValueOrDefault(detail.QtyPo1, 0) + getValueOrDefault(detail.QtyPo2, 0) + getValueOrDefault(detail.QtyPo3, 0)
	default:
		return getValueOrDefault(detail.Qty1, 0) + getValueOrDefault(detail.Qty2, 0) + getValueOrDefault(detail.Qty3, 0)
	}
}

func isActiveDetailForTab(detail model.OrderDetailRead, tab promoSnapshotTab) bool {
	if detail.ItemType == 2 {
		return false
	}
	return activeQtyForTab(detail, tab) > 0
}

func hasPurchaseDisplayQty(detail model.OrderDetailRead) bool {
	return getValueOrDefault(detail.QtyPo1, 0) > 0 ||
		getValueOrDefault(detail.QtyPo2, 0) > 0 ||
		getValueOrDefault(detail.QtyPo3, 0) > 0 ||
		getValueOrDefault(detail.OriginalQtyPo1, 0) > 0 ||
		getValueOrDefault(detail.OriginalQtyPo2, 0) > 0 ||
		getValueOrDefault(detail.OriginalQtyPo3, 0) > 0
}

func shouldIncludePurchaseDetailRow(detail model.OrderDetailRead) bool {
	if detail.ItemType == 2 {
		return false
	}
	return hasPurchaseDisplayQty(detail) || isActiveDetailForTab(detail, promoSnapshotTabSalesOrder)
}

func filterActiveNormalDetailsForTab(details []model.OrderDetailRead, tab promoSnapshotTab) []model.OrderDetailRead {
	active := make([]model.OrderDetailRead, 0, len(details))
	for _, detail := range details {
		if !isActiveDetailForTab(detail, tab) {
			continue
		}
		active = append(active, detail)
	}
	return active
}

func buildOrderItemsForPromoTab(details []model.OrderDetailRead, tab promoSnapshotTab) []entity.OrderDetResponse {
	activeDetails := filterActiveNormalDetailsForTab(details, tab)
	items := make([]entity.OrderDetResponse, 0, len(activeDetails))
	for _, detail := range activeDetails {
		item := entity.OrderDetResponse{ProId: detail.ProId}
		switch tab {
		case promoSnapshotTabFinalOrder:
			item.Qty1 = float64Ptr(getValueOrDefault(detail.Qty1Final, 0))
			item.Qty2 = float64Ptr(getValueOrDefault(detail.Qty2Final, 0))
			item.Qty3 = float64Ptr(getValueOrDefault(detail.Qty3Final, 0))
			item.SellPrice1 = float64Ptr(getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0)))
			item.SellPrice2 = float64Ptr(getValueOrDefault(detail.SellPriceFinal2, getValueOrDefault(detail.SellPrice2, 0)))
			item.SellPrice3 = float64Ptr(getValueOrDefault(detail.SellPriceFinal3, getValueOrDefault(detail.SellPrice3, 0)))
		case promoSnapshotTabPurchase:
			item.Qty1 = float64Ptr(getValueOrDefault(detail.QtyPo1, 0))
			item.Qty2 = float64Ptr(getValueOrDefault(detail.QtyPo2, 0))
			item.Qty3 = float64Ptr(getValueOrDefault(detail.QtyPo3, 0))
			item.SellPrice1 = float64Ptr(getValueOrDefault(detail.SellPricePo1, 0))
			item.SellPrice2 = float64Ptr(getValueOrDefault(detail.SellPricePo2, 0))
			item.SellPrice3 = float64Ptr(getValueOrDefault(detail.SellPricePo3, 0))
		default:
			item.Qty1 = float64Ptr(getValueOrDefault(detail.Qty1, 0))
			item.Qty2 = float64Ptr(getValueOrDefault(detail.Qty2, 0))
			item.Qty3 = float64Ptr(getValueOrDefault(detail.Qty3, 0))
			item.SellPrice1 = float64Ptr(getValueOrDefault(detail.SellPrice1, 0))
			item.SellPrice2 = float64Ptr(getValueOrDefault(detail.SellPrice2, 0))
			item.SellPrice3 = float64Ptr(getValueOrDefault(detail.SellPrice3, 0))
		}
		items = append(items, item)
	}
	return items
}

type promoFlagOverride struct {
	SalesOrder    *bool
	FinalOrder    *bool
	PurchaseOrder *bool
}

func applyExplicitPromoFlagOverride(updates map[string]interface{}, tab promoSnapshotTab, override promoFlagOverride) {
	switch tab {
	case promoSnapshotTabSalesOrder:
		if override.SalesOrder != nil {
			updates["is_product_promotion_so"] = *override.SalesOrder
		}
	case promoSnapshotTabFinalOrder:
		if override.FinalOrder != nil {
			updates["is_product_promotion_final"] = *override.FinalOrder
		}
	case promoSnapshotTabPurchase:
		if override.PurchaseOrder != nil {
			updates["is_product_promotion_po"] = *override.PurchaseOrder
		}
	}
}

func (service *orderServiceImpl) recomputePromoStateForTab(ctx context.Context, ro model.OrderList, custID string, parentCustID string, details []model.OrderDetailRead, tab promoSnapshotTab, explicitPromoOverrides map[int64]promoFlagOverride) (model.Order, error) {
	headerUpdate := model.Order{}
	activeDetails := filterActiveNormalDetailsForTab(details, tab)
	items := buildOrderItemsForPromoTab(activeDetails, tab)
	aggregate := map[int]promoAggregateRow{}
	remarks := []string{}

	if service.PromotionRepository != nil && service.PromotionV2Repository != nil {
		promoService := NewPromotionService(service.PromotionRepository, service.PromotionV2Repository, service.Transaction)
		consultResp, err := promoService.ConsultV2(buildConsultPayloadByTab(ro, items, custID, parentCustID))
		if err != nil {
			return headerUpdate, err
		}
		aggregate = distributePromoToDetailRowsV2(aggregatePromoByProductForDetailSnapshot(consultResp), activeDetails, tab, consultResp)
		remarks = buildFinalRemarks(consultResp)
	}

	subTotal := 0.0
	discTotal := 0.0
	vatTotal := 0.0
	total := 0.0
	promoTotal := 0.0

	for _, detail := range activeDetails {
		if detail.OrderDetailID == nil {
			continue
		}

		var qty1, qty2, qty3, price1, price2, price3 float64
		switch tab {
		case promoSnapshotTabFinalOrder:
			qty1 = getValueOrDefault(detail.Qty1Final, 0)
			qty2 = getValueOrDefault(detail.Qty2Final, 0)
			qty3 = getValueOrDefault(detail.Qty3Final, 0)
			price1 = getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0))
			price2 = getValueOrDefault(detail.SellPriceFinal2, getValueOrDefault(detail.SellPrice2, 0))
			price3 = getValueOrDefault(detail.SellPriceFinal3, getValueOrDefault(detail.SellPrice3, 0))
		case promoSnapshotTabPurchase:
			qty1 = getValueOrDefault(detail.QtyPo1, 0)
			qty2 = getValueOrDefault(detail.QtyPo2, 0)
			qty3 = getValueOrDefault(detail.QtyPo3, 0)
			price1 = getValueOrDefault(detail.SellPricePo1, 0)
			price2 = getValueOrDefault(detail.SellPricePo2, 0)
			price3 = getValueOrDefault(detail.SellPricePo3, 0)
		default:
			qty1 = getValueOrDefault(detail.Qty1, 0)
			qty2 = getValueOrDefault(detail.Qty2, 0)
			qty3 = getValueOrDefault(detail.Qty3, 0)
			price1 = getValueOrDefault(detail.SellPrice1, 0)
			price2 = getValueOrDefault(detail.SellPrice2, 0)
			price3 = getValueOrDefault(detail.SellPrice3, 0)
		}

		gross := (qty1 * price1) + (qty2 * price2) + (qty3 * price3)
		promoValue := aggregate[*detail.OrderDetailID].PromoTotal
		discValue := 0.0
		if ro.OutletID != nil {
			discValue, _, _ = service.CalculateLineDiscount(custID, parentCustID, int(*ro.OutletID), detail.ProId, gross-promoValue)
		}
		vatPercent := getValueOrDefault(detail.Vat, 0)
		if detail.Vat == nil {
			if ro.Vat != nil {
				vatPercent = *ro.Vat
			} else if product, productErr := service.OrderRepository.FindProductByID(detail.ProId); productErr == nil {
				vatPercent = product.Vat
			}
		}
		vatValue := calculateVatValue(qty1, qty2, qty3, price1, price2, price3, promoValue, discValue, vatPercent)
		amount := gross - promoValue - discValue + vatValue

		updates := buildDetailPromoSnapshotUpdates(detail, aggregate, tab)
		if override, ok := explicitPromoOverrides[int64(*detail.OrderDetailID)]; ok {
			applyExplicitPromoFlagOverride(updates, tab, override)
		}
		switch tab {
		case promoSnapshotTabFinalOrder:
			updates["promo_value_final"] = promoValue
			updates["disc_value_final"] = discValue
			updates["vat_value_final"] = vatValue
			updates["amount_final"] = amount
		case promoSnapshotTabPurchase:
			updates["disc_po"] = discValue
			updates["vat_value_po"] = vatValue
		default:
			updates["promo_value"] = promoValue
			updates["disc_value"] = discValue
			updates["vat_value"] = vatValue
			updates["amount"] = amount
		}

		if err := service.OrderRepository.UpdateDetailPartial(ctx, int64(*detail.OrderDetailID), custID, updates); err != nil {
			return headerUpdate, err
		}

		subTotal += gross
		discTotal += discValue
		vatTotal += vatValue
		total += amount
		promoTotal += promoValue
	}

	switch tab {
	case promoSnapshotTabFinalOrder:
		headerUpdate.SubTotalFinal = float64Ptr(subTotal)
		headerUpdate.DiscValueFinal = float64Ptr(discTotal)
		headerUpdate.VatValueFinal = float64Ptr(vatTotal)
		headerUpdate.TotalFinal = float64Ptr(total)
		headerUpdate.PromoValueFinal = float64Ptr(promoTotal)
		headerUpdate.PromoRemarksFinal = model.JSONStringArray(append([]string{}, remarks...))
	case promoSnapshotTabPurchase:
		headerUpdate.PromoRemarksPo = model.JSONStringArray(append([]string{}, remarks...))
	default:
		headerUpdate.SubTotal = float64Ptr(subTotal)
		headerUpdate.DiscValue = float64Ptr(discTotal)
		headerUpdate.VatValue = float64Ptr(vatTotal)
		headerUpdate.Total = float64Ptr(total)
		headerUpdate.PromoValue = float64Ptr(promoTotal)
		headerUpdate.PromoRemarksSo = model.JSONStringArray(append([]string{}, remarks...))
	}

	return headerUpdate, nil
}

func buildEnhanceProjectionItems(currentDetails []model.OrderDetailRead, targetDetailID int64, tab promoSnapshotTab, qty1 float64, qty2 float64, qty3 float64, price1 float64, price2 float64, price3 float64) []entity.OrderDetResponse {
	items := make([]entity.OrderDetResponse, 0, len(currentDetails))

	for _, detail := range currentDetails {
		if detail.ItemType == 2 {
			continue
		}

		item := entity.OrderDetResponse{
			ProId:           detail.ProId,
			Qty1:            float64Ptr(getValueOrDefault(detail.Qty1, 0)),
			Qty2:            float64Ptr(getValueOrDefault(detail.Qty2, 0)),
			Qty3:            float64Ptr(getValueOrDefault(detail.Qty3, 0)),
			Qty1Final:       float64Ptr(getValueOrDefault(detail.Qty1Final, 0)),
			Qty2Final:       float64Ptr(getValueOrDefault(detail.Qty2Final, 0)),
			Qty3Final:       float64Ptr(getValueOrDefault(detail.Qty3Final, 0)),
			QtyPo1:          float64Ptr(getValueOrDefault(detail.QtyPo1, 0)),
			QtyPo2:          float64Ptr(getValueOrDefault(detail.QtyPo2, 0)),
			QtyPo3:          float64Ptr(getValueOrDefault(detail.QtyPo3, 0)),
			SellPrice1:      float64Ptr(getValueOrDefault(detail.SellPrice1, 0)),
			SellPrice2:      float64Ptr(getValueOrDefault(detail.SellPrice2, 0)),
			SellPrice3:      float64Ptr(getValueOrDefault(detail.SellPrice3, 0)),
			SellPricePo1:    float64Ptr(getValueOrDefault(detail.SellPricePo1, 0)),
			SellPricePo2:    float64Ptr(getValueOrDefault(detail.SellPricePo2, 0)),
			SellPricePo3:    float64Ptr(getValueOrDefault(detail.SellPricePo3, 0)),
			SellPriceFinal1: float64Ptr(getValueOrDefault(detail.SellPriceFinal1, getValueOrDefault(detail.SellPrice1, 0))),
			SellPriceFinal2: float64Ptr(getValueOrDefault(detail.SellPriceFinal2, getValueOrDefault(detail.SellPrice2, 0))),
			SellPriceFinal3: float64Ptr(getValueOrDefault(detail.SellPriceFinal3, getValueOrDefault(detail.SellPrice3, 0))),
		}

		if detail.OrderDetailID != nil && int64(*detail.OrderDetailID) == targetDetailID {
			switch tab {
			case promoSnapshotTabSalesOrder:
				item.Qty1 = float64Ptr(qty1)
				item.Qty2 = float64Ptr(qty2)
				item.Qty3 = float64Ptr(qty3)
				item.SellPrice1 = float64Ptr(price1)
				item.SellPrice2 = float64Ptr(price2)
				item.SellPrice3 = float64Ptr(price3)
				item.Qty1Final = float64Ptr(qty1)
				item.Qty2Final = float64Ptr(qty2)
				item.Qty3Final = float64Ptr(qty3)
				item.SellPriceFinal1 = float64Ptr(price1)
				item.SellPriceFinal2 = float64Ptr(price2)
				item.SellPriceFinal3 = float64Ptr(price3)
			case promoSnapshotTabFinalOrder:
				item.Qty1Final = float64Ptr(qty1)
				item.Qty2Final = float64Ptr(qty2)
				item.Qty3Final = float64Ptr(qty3)
				item.SellPriceFinal1 = float64Ptr(price1)
				item.SellPriceFinal2 = float64Ptr(price2)
				item.SellPriceFinal3 = float64Ptr(price3)
			case promoSnapshotTabPurchase:
				item.QtyPo1 = float64Ptr(qty1)
				item.QtyPo2 = float64Ptr(qty2)
				item.QtyPo3 = float64Ptr(qty3)
				item.SellPricePo1 = float64Ptr(price1)
				item.SellPricePo2 = float64Ptr(price2)
				item.SellPricePo3 = float64Ptr(price3)
			}
		}

		switch tab {
		case promoSnapshotTabSalesOrder:
			item.SellPrice1 = float64Ptr(getValueOrDefault(item.SellPrice1, 0))
			item.SellPrice2 = float64Ptr(getValueOrDefault(item.SellPrice2, 0))
			item.SellPrice3 = float64Ptr(getValueOrDefault(item.SellPrice3, 0))
		case promoSnapshotTabFinalOrder:
			item.Qty1 = float64Ptr(getValueOrDefault(item.Qty1Final, 0))
			item.Qty2 = float64Ptr(getValueOrDefault(item.Qty2Final, 0))
			item.Qty3 = float64Ptr(getValueOrDefault(item.Qty3Final, 0))
			item.SellPrice1 = float64Ptr(getValueOrDefault(item.SellPriceFinal1, 0))
			item.SellPrice2 = float64Ptr(getValueOrDefault(item.SellPriceFinal2, 0))
			item.SellPrice3 = float64Ptr(getValueOrDefault(item.SellPriceFinal3, 0))
		case promoSnapshotTabPurchase:
			item.Qty1 = float64Ptr(getValueOrDefault(item.QtyPo1, 0))
			item.Qty2 = float64Ptr(getValueOrDefault(item.QtyPo2, 0))
			item.Qty3 = float64Ptr(getValueOrDefault(item.QtyPo3, 0))
			item.SellPrice1 = float64Ptr(getValueOrDefault(item.SellPricePo1, 0))
			item.SellPrice2 = float64Ptr(getValueOrDefault(item.SellPricePo2, 0))
			item.SellPrice3 = float64Ptr(getValueOrDefault(item.SellPricePo3, 0))
		}

		items = append(items, item)
	}

	return items
}

func buildEnhancePromoPayload(ro model.OrderList, custID string, parentCustID string, currentDetails []model.OrderDetailRead, targetDetailID int64, tab promoSnapshotTab, qty1 float64, qty2 float64, qty3 float64, price1 float64, price2 float64, price3 float64) entity.ConsultPromoV2Req {
	items := buildEnhanceProjectionItems(currentDetails, targetDetailID, tab, qty1, qty2, qty3, price1, price2, price3)
	return buildConsultPayloadByTab(ro, items, custID, parentCustID)
}

func mergeUniqueRemarks(base []string, incoming []string) []string {
	merged := append([]string{}, base...)
	for _, remark := range incoming {
		if remark == "" {
			continue
		}
		if !slices.Contains(merged, remark) {
			merged = append(merged, remark)
		}
	}
	sort.Strings(merged)
	return merged
}

func backfillPromoItemUnits(items []entity.OrderDetResponse, productMap map[int]model.Product) {
	for i := range items {
		item := &items[i]
		if !isEmptyStringPtr(item.UnitId1) {
			continue
		}

		product, exists := productMap[item.ProId]
		if !exists {
			continue
		}

		item.UnitId1 = stringPtr(product.UnitId1)
		item.UnitId2 = stringPtr(product.UnitId2)
		item.UnitId3 = stringPtr(product.UnitId3)
		item.UnitId4 = product.UnitId4
		item.UnitId5 = product.UnitId5
	}
}

func movePromoDetailsToNormal(normal []entity.OrderDetResponse, promo []entity.OrderDetResponse) ([]entity.OrderDetResponse, []entity.OrderDetResponse) {
	moved := append([]entity.OrderDetResponse{}, normal...)
	for _, promoItem := range promo {
		promoItem.IsProductPromotion = true
		moved = append(moved, promoItem)
	}
	return moved, []entity.OrderDetResponse{}
}

func stockDisplayQtyByPriority(detail model.OrderDetailRead) (qty1 float64, qty2 float64, qty3 float64) {
	if detail.Qty1Final != nil || detail.Qty2Final != nil || detail.Qty3Final != nil {
		return getValueOrDefault(detail.Qty1Final, 0), getValueOrDefault(detail.Qty2Final, 0), getValueOrDefault(detail.Qty3Final, 0)
	}

	if detail.Qty1 != nil || detail.Qty2 != nil || detail.Qty3 != nil {
		return getValueOrDefault(detail.Qty1, 0), getValueOrDefault(detail.Qty2, 0), getValueOrDefault(detail.Qty3, 0)
	}

	return getValueOrDefault(detail.QtyPo1, 0), getValueOrDefault(detail.QtyPo2, 0), getValueOrDefault(detail.QtyPo3, 0)
}

func getValueOld(service *orderServiceImpl, details []model.OrderDetailRead, proID int64, qty1 float64, qty2 float64, qty3 float64) (float64, error) {
	var productIDs []int64
	for _, detail := range details {
		productIDs = append(productIDs, int64(detail.ProId))
	}

	productsModel, err := service.OrderRepository.FindProductByListID(productIDs)
	if err != nil {
		return 0.0, err
	}

	var productMap = model.MapProduct{}

	for _, productModel := range productsModel {
		productMap.SetProduct(productModel.ProductId, productModel)
	}

	productModel, err := productMap.GetByID(int64(proID))
	if err != nil {
		return 0.0, err
	}

	QtyUnit := &conversion.QtyUnit{
		Qty1:      int(qty1),
		Qty2:      int(qty2),
		Qty3:      int(qty3),
		ConvUnit2: int(productModel.ConvUnit2),
		ConvUnit3: int(productModel.ConvUnit3),
	}

	totalQty, err := QtyUnit.ToTotalQuantity()
	if err != nil {
		return 0.0, err
	}

	return float64(totalQty), nil
}

func normalizeMobileUnitID(unitID *string) string {
	if unitID == nil {
		return ""
	}
	return strings.TrimSpace(strings.ToUpper(*unitID))
}

func getMobileQtyValue(primary *float64, fallback *float64) float64 {
	if primary != nil {
		return *primary
	}
	if fallback != nil {
		return *fallback
	}
	return 0
}

func getConvUnitValue(conv *int) int {
	if conv != nil && *conv > 0 {
		return *conv
	}
	return 1
}

func getConvUnitValueWithFallback(primary *int, fallback *int) int {
	if primary != nil && *primary > 0 {
		return *primary
	}
	if fallback != nil && *fallback > 0 {
		return *fallback
	}
	return 1
}

func floatEquals(a float64, b float64) bool {
	return math.Abs(a-b) <= FLOAT_COMPARE_EPSILON
}

func toSmallestUnitQty(qty1 float64, qty2 float64, qty3 float64, conv2 int, conv3 int) (float64, error) {
	qtyUnit := &conversion.QtyUnit{
		Qty1:      int(math.Round(qty1)),
		Qty2:      int(math.Round(qty2)),
		Qty3:      int(math.Round(qty3)),
		ConvUnit2: conv2,
		ConvUnit3: conv3,
	}
	totalQty, err := qtyUnit.ToTotalQuantity()
	if err != nil {
		return 0, err
	}
	return float64(totalQty), nil
}

func mobileDetailIdentity(orderDetailID int64, proID int, unit1 *string, unit2 *string, unit3 *string) string {
	if orderDetailID > 0 {
		return fmt.Sprintf("id:%d", orderDetailID)
	}
	return fmt.Sprintf("key:%d|%s|%s|%s", proID, normalizeMobileUnitID(unit1), normalizeMobileUnitID(unit2), normalizeMobileUnitID(unit3))
}

func (service *orderServiceImpl) isMobileProcessNoMeaningfulDetailChange(roNo string, ro model.OrderList, request entity.UpdateOrderBody) (bool, error) {
	if ro.DataSource == nil || *ro.DataSource != 2 {
		return false, nil
	}

	existingDetails, err := service.OrderRepository.FindDetail(roNo, request.CustId)
	if err != nil {
		return false, err
	}

	existingNormal := make(map[string]model.OrderDetailRead)
	for _, detail := range existingDetails {
		if detail.ItemType != 1 {
			continue
		}
		orderDetailID := int64(0)
		if detail.OrderDetailID != nil {
			orderDetailID = int64(*detail.OrderDetailID)
		}
		identity := mobileDetailIdentity(orderDetailID, detail.ProId, detail.UnitId1, detail.UnitId2, detail.UnitId3)
		existingNormal[identity] = detail
	}

	if len(request.Details.Normal) == 0 {
		if len(existingNormal) == 0 {
			return true, nil
		}
		return false, nil
	}

	requestNormal := make(map[string]entity.UpdateOrderDetBody)
	for _, detail := range request.Details.Normal {
		orderDetailID := int64(0)
		if detail.OrderDetId != nil {
			orderDetailID = *detail.OrderDetId
		}
		identity := mobileDetailIdentity(orderDetailID, detail.ProId, detail.UnitId1, detail.UnitId2, detail.UnitId3)
		requestNormal[identity] = detail
	}

	if len(requestNormal) != len(existingNormal) {
		return false, nil
	}

	for identity, requestDetail := range requestNormal {
		existingDetail, exists := existingNormal[identity]
		if !exists {
			return false, nil
		}

		if requestDetail.ProId != existingDetail.ProId {
			return false, nil
		}

		requestHasOrderDetailID := requestDetail.OrderDetId != nil && *requestDetail.OrderDetId > 0
		existingHasOrderDetailID := existingDetail.OrderDetailID != nil && *existingDetail.OrderDetailID > 0
		if !(requestHasOrderDetailID && existingHasOrderDetailID) {
			if normalizeMobileUnitID(requestDetail.UnitId1) != normalizeMobileUnitID(existingDetail.UnitId1) ||
				normalizeMobileUnitID(requestDetail.UnitId2) != normalizeMobileUnitID(existingDetail.UnitId2) ||
				normalizeMobileUnitID(requestDetail.UnitId3) != normalizeMobileUnitID(existingDetail.UnitId3) {
				return false, nil
			}
		}

		requestQtyPo1 := getMobileQtyValue(requestDetail.QtyPo1, requestDetail.Qty1)
		requestQtyPo2 := getMobileQtyValue(requestDetail.QtyPo2, requestDetail.Qty2)
		requestQtyPo3 := getMobileQtyValue(requestDetail.QtyPo3, requestDetail.Qty3)

		existingQtyPo1 := getMobileQtyValue(existingDetail.QtyPo1, existingDetail.Qty1)
		existingQtyPo2 := getMobileQtyValue(existingDetail.QtyPo2, existingDetail.Qty2)
		existingQtyPo3 := getMobileQtyValue(existingDetail.QtyPo3, existingDetail.Qty3)

		if !floatEquals(requestQtyPo1, existingQtyPo1) ||
			!floatEquals(requestQtyPo2, existingQtyPo2) ||
			!floatEquals(requestQtyPo3, existingQtyPo3) {
			return false, nil
		}

		requestSellPrice1 := getValueOrDefault(requestDetail.SellPrice1, getValueOrDefault(existingDetail.SellPrice1, 0))
		requestSellPrice2 := getValueOrDefault(requestDetail.SellPrice2, getValueOrDefault(existingDetail.SellPrice2, 0))
		requestSellPrice3 := getValueOrDefault(requestDetail.SellPrice3, getValueOrDefault(existingDetail.SellPrice3, 0))
		existingSellPrice1 := getValueOrDefault(existingDetail.SellPrice1, 0)
		existingSellPrice2 := getValueOrDefault(existingDetail.SellPrice2, 0)
		existingSellPrice3 := getValueOrDefault(existingDetail.SellPrice3, 0)

		if !floatEquals(requestSellPrice1, existingSellPrice1) ||
			!floatEquals(requestSellPrice2, existingSellPrice2) ||
			!floatEquals(requestSellPrice3, existingSellPrice3) {
			return false, nil
		}

		requestSmallest, err := toSmallestUnitQty(
			requestQtyPo1,
			requestQtyPo2,
			requestQtyPo3,
			getConvUnitValueWithFallback(requestDetail.ConvUnit2, existingDetail.MpConvUnit2),
			getConvUnitValueWithFallback(requestDetail.ConvUnit3, existingDetail.MpConvUnit3),
		)
		if err != nil {
			return false, err
		}
		existingSmallest, err := toSmallestUnitQty(
			existingQtyPo1,
			existingQtyPo2,
			existingQtyPo3,
			getConvUnitValueWithFallback(existingDetail.MpConvUnit2, requestDetail.ConvUnit2),
			getConvUnitValueWithFallback(existingDetail.MpConvUnit3, requestDetail.ConvUnit3),
		)
		if err != nil {
			return false, err
		}

		if !floatEquals(requestSmallest, existingSmallest) {
			return false, nil
		}
	}

	return true, nil
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

	promoDetails := make([]entity.OrderDetResponse, 0)
	for index, detail := range details {
		log.Info("INDEX : ", index)
		var detailData entity.OrderDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}

		if detailData.Qty == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.Qty = &qtyOld
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

		qtyOrder := int(getValueOrDefault(detailData.Qty, 0))
		qtyFinal := int(getValueOrDefault(detailData.QtyFinal, 0))

		if qtyFinal < qtyOrder {
			detailData.OrderStatus = "Partial Reject"
		}
		if qtyFinal < 1 {
			detailData.OrderStatus = "Rejected"
		}

		if detailData.ItemType == 1 {
			response.Details.Normal = append(response.Details.Normal, detailData)
		} else {
			promoDetails = append(promoDetails, detailData)
		}
	}

	log.Info("MASUK PROMO DETAIL")
	response.Details.Promo = append(response.Details.Promo, promoDetails...)
	//final
	log.Info("MASUK FINAL")
	promoFinalDetails := make([]entity.OrderDetResponse, 0)
	for index, detail := range details {
		log.Info("INDEX FINAL : ", index)
		var detailData entity.OrderDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}
		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}

		if detailData.QtyFinal == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.QtyFinal = &qtyOld
		}

		qty := &conversion.Qty{
			Qty:       int(getValueOrDefault(detailData.QtyFinal, 0)),
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

		if detailData.ItemType == 1 {
			response.DetailsFinal.Normal = append(response.DetailsFinal.Normal, detailData)
		} else {
			promoFinalDetails = append(promoFinalDetails, detailData)
		}
	}

	response.DetailsFinal.Promo = append(response.DetailsFinal.Promo, promoFinalDetails...)

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

	rewards, err := service.OrderRepository.FindReward(RoNo, custID)
	if err != nil {
		return response, err
	}

	for _, reward := range rewards {
		var remark entity.OrderRewardResponse

		if err = structs.Automapper(reward, &remark); err != nil {
			return response, err
		}

		remark.RewardTypeName = remark.GenerateRewardTypeName()

		response.Remarks = append(response.Remarks, remark)
	}

	return response, nil
}

func (service *orderServiceImpl) DetailV2(RoNo string, custID string, parentCustID string) (response entity.OrderResponse, err error) {
	ro, err := service.OrderRepository.FindByNo(RoNo, custID)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	// Map data_source to source
	response.Source = MapDataSourceToSource(ro.DataSource)

	// Set is_proforma_inv (default false if null)
	if ro.IsProformaInv != nil {
		response.IsProformaInv = ro.IsProformaInv
	} else {
		falseVal := false
		response.IsProformaInv = &falseVal
	}

	details, err := service.OrderRepository.FindDetail(RoNo, custID)
	if err != nil {
		return response, err
	}

	var whId int64
	if ro.WhId != nil {
		whId = *ro.WhId
	}

	var proIds []int64
	for _, detail := range details {
		proIds = append(proIds, int64(detail.ProId))
	}

	warehouseStockMap, err := service.OrderRepository.FindWarehouseStockByWhIdAndProIds(custID, whId, proIds)
	if err != nil {
		warehouseStockMap = make(map[int64]float64)
	}

	useWarehouseCurrentOnly := ro.DataStatus != nil && *ro.DataStatus == entity.CANCELLED

	promoDetails := make([]entity.OrderDetResponse, 0)
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

		// Backward compatible: fallback SellPriceFinal to SellPrice if nil or 0
		if detailData.SellPriceFinal1 == nil || *detailData.SellPriceFinal1 == 0 {
			detailData.SellPriceFinal1 = detailData.SellPrice1
		}
		if detailData.SellPriceFinal2 == nil || *detailData.SellPriceFinal2 == 0 {
			detailData.SellPriceFinal2 = detailData.SellPrice2
		}
		if detailData.SellPriceFinal3 == nil || *detailData.SellPriceFinal3 == 0 {
			detailData.SellPriceFinal3 = detailData.SellPrice3
		}

		if detailData.Qty == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.Qty = &qtyOld
		}

		isSalesMapping := ro.IsSalesMapping != nil && *ro.IsSalesMapping
		if !isSalesMapping {
			qty := &conversion.Qty{
				Qty:       int(getValueOrDefault(detailData.Qty, 0)),
				ConvUnit2: int(*detail.MpConvUnit2),
				ConvUnit3: int(*detail.MpConvUnit3),
			}

			qtyConversion := qty.ConvToQtyConversion()
			detailDataQty1 := float64(qtyConversion.Qty3)
			detailDataQty2 := float64(qtyConversion.Qty2)
			detailDataQty3 := float64(qtyConversion.Qty1)

			detailData.Qty1 = &detailDataQty1
			detailData.Qty2 = &detailDataQty2
			detailData.Qty3 = &detailDataQty3
		}

		whStockQty := warehouseStockMap[int64(detail.ProId)]
		convUnit2 := 1
		convUnit3 := 1
		if detail.MpConvUnit2 != nil && *detail.MpConvUnit2 > 0 {
			convUnit2 = *detail.MpConvUnit2
		}
		if detail.MpConvUnit3 != nil && *detail.MpConvUnit3 > 0 {
			convUnit3 = *detail.MpConvUnit3
		}
		dispQty1, dispQty2, dispQty3 := stockDisplayQtyByPriority(detail)
		stockBreakdown := computeDisplayedAvailableStockBreakdown(
			int(whStockQty),
			dispQty3,
			dispQty2,
			dispQty1,
			!useWarehouseCurrentOnly,
			convUnit2,
			convUnit3,
		)
		applyStockBreakdownToPointers(&detailData.Qty1Stok, &detailData.Qty2Stok, &detailData.Qty3Stok, stockBreakdown)

		qtyOrder := int(getValueOrDefault(detailData.Qty, 0))
		qtyFinal := int(getValueOrDefault(detailData.QtyFinal, 0))

		if qtyFinal < qtyOrder {
			detailData.OrderStatus = "Partial Reject"
		}
		if qtyFinal < 1 {
			detailData.OrderStatus = "Rejected"
		}

		if detailData.ItemType == 1 {
			if isActiveDetailForTab(detail, promoSnapshotTabSalesOrder) {
				response.Details.Normal = append(response.Details.Normal, detailData)
			}
		} else {
			promoDetails = append(promoDetails, detailData)
		}
	}

	response.Details.Promo = append(response.Details.Promo, promoDetails...)
	promoProductMasterMap := make(map[int]model.Product)
	promoProductIDs := collectPromoItemProductIDsForFallback(response.Details.Promo)
	if len(promoProductIDs) > 0 {
		productsModel, productErr := service.OrderRepository.FindProductByListID(promoProductIDs)
		if productErr != nil {
			return response, fmt.Errorf("failed to fetch product masters for promo items: %w", productErr)
		}

		promoProductMasterMap = buildProductMasterMap(productsModel)
	}
	backfillPromoItemUnits(response.Details.Promo, promoProductMasterMap)
	response.Details.Normal, response.Details.Promo = movePromoDetailsToNormal(response.Details.Normal, response.Details.Promo)

	// Recalculate VatValue for Details.Normal based on formula:
	// ((qty1 * sell_price1) + (qty2 * sell_price2) + (qty3 * sell_price3) - promo_value_final - disc_value_final) * vat%
	for i := range response.Details.Normal {
		item := &response.Details.Normal[i]
		qty1 := getValueOrDefault(item.Qty1, 0)
		qty2 := getValueOrDefault(item.Qty2, 0)
		qty3 := getValueOrDefault(item.Qty3, 0)
		price1 := getValueOrDefault(item.SellPrice1, 0)
		price2 := getValueOrDefault(item.SellPrice2, 0)
		price3 := getValueOrDefault(item.SellPrice3, 0)
		promo := getValueOrDefault(item.PromoValueFinal, 0)
		disc := getValueOrDefault(item.DiscValueFinal, 0)
		vat := getValueOrDefault(item.Vat, 0)
		if item.Vat == nil {
			if response.Vat != nil {
				vat = *response.Vat
			} else if product, productErr := service.OrderRepository.FindProductByID(item.ProId); productErr == nil {
				vat = product.Vat
			}
		}

		vatValue := calculateVatValue(qty1, qty2, qty3, price1, price2, price3, promo, disc, vat)
		item.VatValue = &vatValue
	}

	// Build PurchaseDetails from persisted purchase rows instead of copying sales details.
	promoPurchaseDetails := make([]entity.OrderDetResponse, 0)
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

		purchaseQtyPoTotal := getValueOrDefault(detailData.QtyPo1, 0) + getValueOrDefault(detailData.QtyPo2, 0) + getValueOrDefault(detailData.QtyPo3, 0)
		if purchaseQtyPoTotal > 0 {
			detailData.Qty1 = detailData.QtyPo1
			detailData.Qty2 = detailData.QtyPo2
			detailData.Qty3 = detailData.QtyPo3
		}
		if detailData.SellPricePo1 != nil {
			detailData.SellPrice1 = detailData.SellPricePo1
		}
		if detailData.SellPricePo2 != nil {
			detailData.SellPrice2 = detailData.SellPricePo2
		}
		if detailData.SellPricePo3 != nil {
			detailData.SellPrice3 = detailData.SellPricePo3
		}

		whStockQtyPurchase := warehouseStockMap[int64(detail.ProId)]
		convUnit2Purchase := 1
		convUnit3Purchase := 1
		if detail.MpConvUnit2 != nil && *detail.MpConvUnit2 > 0 {
			convUnit2Purchase = *detail.MpConvUnit2
		}
		if detail.MpConvUnit3 != nil && *detail.MpConvUnit3 > 0 {
			convUnit3Purchase = *detail.MpConvUnit3
		}
		dispQty1Purchase, dispQty2Purchase, dispQty3Purchase := stockDisplayQtyByPriority(detail)
		stockBreakdownPurchase := computeDisplayedAvailableStockBreakdown(
			int(whStockQtyPurchase),
			dispQty3Purchase,
			dispQty2Purchase,
			dispQty1Purchase,
			!useWarehouseCurrentOnly,
			convUnit2Purchase,
			convUnit3Purchase,
		)
		applyStockBreakdownToPointers(&detailData.Qty1Stok, &detailData.Qty2Stok, &detailData.Qty3Stok, stockBreakdownPurchase)

		if detailData.ItemType == 1 {
			if shouldIncludePurchaseDetailRow(detail) {
				response.PurchaseDetails.Normal = append(response.PurchaseDetails.Normal, detailData)
			}
		} else {
			promoPurchaseDetails = append(promoPurchaseDetails, detailData)
		}
	}

	response.PurchaseDetails.Promo = append(response.PurchaseDetails.Promo, promoPurchaseDetails...)
	promoPurchaseProductMasterMap := make(map[int]model.Product)
	promoPurchaseProductIDs := collectPromoItemProductIDsForFallback(response.PurchaseDetails.Promo)
	if len(promoPurchaseProductIDs) > 0 {
		productsModel, productErr := service.OrderRepository.FindProductByListID(promoPurchaseProductIDs)
		if productErr != nil {
			return response, fmt.Errorf("failed to fetch product masters for purchase promo items: %w", productErr)
		}

		promoPurchaseProductMasterMap = buildProductMasterMap(productsModel)
	}
	backfillPromoItemUnits(response.PurchaseDetails.Promo, promoPurchaseProductMasterMap)
	response.PurchaseDetails.Normal, response.PurchaseDetails.Promo = movePromoDetailsToNormal(response.PurchaseDetails.Normal, response.PurchaseDetails.Promo)

	// Set OrderStatus in PurchaseDetails based on header data_status (sls.order.data_status)
	// Also map sell_price_po to sell_price for purchase_details
	var headerStatus string
	if response.DataStatus != nil {
		headerStatus = fmt.Sprintf("%d", *response.DataStatus)
	}
	for i := range response.PurchaseDetails.Normal {
		response.PurchaseDetails.Normal[i].OrderStatus = headerStatus
		// Map sell_price_po to sell_price for purchase_details
		if response.PurchaseDetails.Normal[i].SellPricePo1 != nil {
			response.PurchaseDetails.Normal[i].SellPrice1 = response.PurchaseDetails.Normal[i].SellPricePo1
		}
		if response.PurchaseDetails.Normal[i].SellPricePo2 != nil {
			response.PurchaseDetails.Normal[i].SellPrice2 = response.PurchaseDetails.Normal[i].SellPricePo2
		}
		if response.PurchaseDetails.Normal[i].SellPricePo3 != nil {
			response.PurchaseDetails.Normal[i].SellPrice3 = response.PurchaseDetails.Normal[i].SellPricePo3
		}
	}
	for i := range response.PurchaseDetails.Promo {
		response.PurchaseDetails.Promo[i].OrderStatus = headerStatus
		// Map sell_price_po to sell_price for purchase_details
		if response.PurchaseDetails.Promo[i].SellPricePo1 != nil {
			response.PurchaseDetails.Promo[i].SellPrice1 = response.PurchaseDetails.Promo[i].SellPricePo1
		}
		if response.PurchaseDetails.Promo[i].SellPricePo2 != nil {
			response.PurchaseDetails.Promo[i].SellPrice2 = response.PurchaseDetails.Promo[i].SellPricePo2
		}
		if response.PurchaseDetails.Promo[i].SellPricePo3 != nil {
			response.PurchaseDetails.Promo[i].SellPrice3 = response.PurchaseDetails.Promo[i].SellPricePo3
		}
	}

	for i := range response.PurchaseDetails.Normal {
		item := &response.PurchaseDetails.Normal[i]

		if item.SellPricePo1 == nil || *item.SellPricePo1 == 0 {
			item.SellPricePo1 = item.SellPrice1
		}
		if item.SellPricePo2 == nil || *item.SellPricePo2 == 0 {
			item.SellPricePo2 = item.SellPrice2
		}
		if item.SellPricePo3 == nil || *item.SellPricePo3 == 0 {
			item.SellPricePo3 = item.SellPrice3
		}

		qtyPo1 := getValueOrDefault(item.QtyPo1, 0)
		qtyPo2 := getValueOrDefault(item.QtyPo2, 0)
		qtyPo3 := getValueOrDefault(item.QtyPo3, 0)
		pricePo1 := getValueOrDefault(item.SellPricePo1, 0)
		pricePo2 := getValueOrDefault(item.SellPricePo2, 0)
		pricePo3 := getValueOrDefault(item.SellPricePo3, 0)
		promo := getValueOrDefault(item.PromoValueFinal, 0)
		disc := getValueOrDefault(item.DiscValueFinal, 0)
		vat := getValueOrDefault(item.Vat, 0)
		if item.Vat == nil {
			if response.Vat != nil {
				vat = *response.Vat
			} else if product, productErr := service.OrderRepository.FindProductByID(item.ProId); productErr == nil {
				vat = product.Vat
			}
		}

		vatValue := calculateVatValue(qtyPo1, qtyPo2, qtyPo3, pricePo1, pricePo2, pricePo3, promo, disc, vat)
		item.VatValue = &vatValue

		vatValueFinal := calculateVatValue(qtyPo1, qtyPo2, qtyPo3, pricePo1, pricePo2, pricePo3, promo, disc, vat)
		item.VatValueFinal = &vatValueFinal
	}

	//final
	promoFinalDetails := make([]entity.OrderDetResponse, 0)
	for _, detail := range details {
		var detailData entity.OrderDetResponse
		err = structs.Automapper(detail, &detailData)
		if err != nil {
			return response, err
		}

		// Backward compatible: fallback SellPriceFinal to SellPrice if nil or 0
		if detailData.SellPriceFinal1 == nil || *detailData.SellPriceFinal1 == 0 {
			detailData.SellPriceFinal1 = detailData.SellPrice1
		}
		if detailData.SellPriceFinal2 == nil || *detailData.SellPriceFinal2 == 0 {
			detailData.SellPriceFinal2 = detailData.SellPrice2
		}
		if detailData.SellPriceFinal3 == nil || *detailData.SellPriceFinal3 == 0 {
			detailData.SellPriceFinal3 = detailData.SellPrice3
		}

		if detail.ExpDate != nil {
			expDate := detail.ExpDate.Format("2006-01-02")
			detailData.ExpDate = &expDate
		}

		if detailData.QtyFinal == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.QtyFinal = &qtyOld
		}

		qty := &conversion.Qty{
			Qty:       int(getValueOrDefault(detailData.QtyFinal, 0)),
			ConvUnit2: int(*detail.MpConvUnit2),
			ConvUnit3: int(*detail.MpConvUnit3),
		}

		qtyConversion := qty.ConvToQtyConversion()

		detailDataQty1 := float64(qtyConversion.Qty1)
		detailDataQty2 := float64(qtyConversion.Qty2)
		detailDataQty3 := float64(qtyConversion.Qty3)

		detailData.Qty1 = &detailDataQty1
		detailData.Qty2 = &detailDataQty2
		detailData.Qty3 = &detailDataQty3

		whStockQtyFinal := warehouseStockMap[int64(detail.ProId)]
		convUnit2Final := 1
		convUnit3Final := 1
		if detail.MpConvUnit2 != nil && *detail.MpConvUnit2 > 0 {
			convUnit2Final = *detail.MpConvUnit2
		}
		if detail.MpConvUnit3 != nil && *detail.MpConvUnit3 > 0 {
			convUnit3Final = *detail.MpConvUnit3
		}
		dispQty1Final, dispQty2Final, dispQty3Final := stockDisplayQtyByPriority(detail)
		stockBreakdownFinal := computeDisplayedAvailableStockBreakdown(
			int(whStockQtyFinal),
			dispQty3Final,
			dispQty2Final,
			dispQty1Final,
			!useWarehouseCurrentOnly,
			convUnit2Final,
			convUnit3Final,
		)
		applyStockBreakdownToPointers(&detailData.Qty1Stok, &detailData.Qty2Stok, &detailData.Qty3Stok, stockBreakdownFinal)

		if detailData.ItemType == 1 {
			if isActiveDetailForTab(detail, promoSnapshotTabFinalOrder) {
				response.DetailsFinal.Normal = append(response.DetailsFinal.Normal, detailData)
			}
		} else {
			promoFinalDetails = append(promoFinalDetails, detailData)
		}
	}

	response.DetailsFinal.Promo = append(response.DetailsFinal.Promo, promoFinalDetails...)
	promoFinalProductMasterMap := make(map[int]model.Product)
	promoFinalProductIDs := collectPromoItemProductIDsForFallback(response.DetailsFinal.Promo)
	if len(promoFinalProductIDs) > 0 {
		productsModel, productErr := service.OrderRepository.FindProductByListID(promoFinalProductIDs)
		if productErr != nil {
			return response, fmt.Errorf("failed to fetch product masters for final promo items: %w", productErr)
		}

		promoFinalProductMasterMap = buildProductMasterMap(productsModel)
	}
	backfillPromoItemUnits(response.DetailsFinal.Promo, promoFinalProductMasterMap)
	response.DetailsFinal.Normal, response.DetailsFinal.Promo = movePromoDetailsToNormal(response.DetailsFinal.Normal, response.DetailsFinal.Promo)

	for i := range response.DetailsFinal.Normal {
		item := &response.DetailsFinal.Normal[i]
		qty1Final := getValueOrDefault(item.Qty1Final, 0)
		qty2Final := getValueOrDefault(item.Qty2Final, 0)
		qty3Final := getValueOrDefault(item.Qty3Final, 0)
		priceFinal1 := getValueOrDefault(item.SellPriceFinal1, 0)
		priceFinal2 := getValueOrDefault(item.SellPriceFinal2, 0)
		priceFinal3 := getValueOrDefault(item.SellPriceFinal3, 0)
		promo := getValueOrDefault(item.PromoValueFinal, 0)
		disc := getValueOrDefault(item.DiscValueFinal, 0)
		vat := getValueOrDefault(item.Vat, 0)
		if item.Vat == nil {
			if response.Vat != nil {
				vat = *response.Vat
			} else if product, productErr := service.OrderRepository.FindProductByID(item.ProId); productErr == nil {
				vat = product.Vat
			}
		}

		vatValueFinal := calculateVatValue(qty1Final, qty2Final, qty3Final, priceFinal1, priceFinal2, priceFinal3, promo, disc, vat)
		item.VatValueFinal = &vatValueFinal
	}

	if parentCustID == "" {
		parentCustID = ro.CustID
	}
	if parentCustID == "" {
		parentCustID = custID
	}

	normalHasSnapshot := hasPersistedPromoSnapshot(ro.PromoRemarksSo, response.Details.Normal, promoSnapshotTabSalesOrder)
	finalHasSnapshot := hasPersistedPromoSnapshot(ro.PromoRemarksFinal, response.DetailsFinal.Normal, promoSnapshotTabFinalOrder)
	purchaseHasSnapshot := hasPersistedPromoSnapshot(ro.PromoRemarksPo, response.PurchaseDetails.Normal, promoSnapshotTabPurchase)
	rewardProductMasterMap := make(map[int]model.Product)

	if normalHasSnapshot || finalHasSnapshot || purchaseHasSnapshot {
		rewardProductIDs := collectPersistedRewardProductIDsForFallback(details)
		if len(rewardProductIDs) > 0 {
			productsModel, productErr := service.OrderRepository.FindProductByListID(rewardProductIDs)
			if productErr != nil {
				return response, fmt.Errorf("failed to fetch product masters for persisted reward products: %w", productErr)
			}

			rewardProductMasterMap = buildProductMasterMap(productsModel)
		}
	}

	response.Details.PromoRemarksSo = append([]string{}, ro.PromoRemarksSo...)
	response.Details.FinalRemarks = append([]string{}, ro.PromoRemarksSo...)
	response.DetailsFinal.PromoRemarksFinal = append([]string{}, ro.PromoRemarksFinal...)
	response.DetailsFinal.FinalRemarks = append([]string{}, ro.PromoRemarksFinal...)
	response.PurchaseDetails.PromoRemarksPo = append([]string{}, ro.PromoRemarksPo...)
	response.PurchaseDetails.FinalRemarks = append([]string{}, ro.PromoRemarksPo...)

	if normalHasSnapshot {
		response.Details.Normal = applyPersistedPromoSnapshotToItems(response.Details.Normal, promoSnapshotTabSalesOrder)
		response.Details.RewardProducts = buildRewardProductsFromPersistedDetails(details, promoSnapshotTabSalesOrder, rewardProductMasterMap)
	}
	if finalHasSnapshot {
		response.DetailsFinal.Normal = applyPersistedPromoSnapshotToItems(response.DetailsFinal.Normal, promoSnapshotTabFinalOrder)
		response.DetailsFinal.RewardProducts = buildRewardProductsFromPersistedDetails(details, promoSnapshotTabFinalOrder, rewardProductMasterMap)
	}
	if purchaseHasSnapshot {
		response.PurchaseDetails.Normal = applyPersistedPromoSnapshotToItems(response.PurchaseDetails.Normal, promoSnapshotTabPurchase)
		response.PurchaseDetails.RewardProducts = buildRewardProductsFromPersistedDetails(details, promoSnapshotTabPurchase, rewardProductMasterMap)
	}

	if (!normalHasSnapshot || !finalHasSnapshot || !purchaseHasSnapshot) && service.PromotionRepository != nil && service.PromotionV2Repository != nil {
		promoService := NewPromotionService(service.PromotionRepository, service.PromotionV2Repository, service.Transaction)

		payloads := map[string]entity.ConsultPromoV2Req{
			"normal":   buildConsultPayloadByTab(ro, response.Details.Normal, custID, parentCustID),
			"final":    buildConsultPayloadByTab(ro, response.DetailsFinal.Normal, custID, parentCustID),
			"purchase": buildConsultPayloadByTab(ro, response.PurchaseDetails.Normal, custID, parentCustID),
		}
		signatures := map[string]string{
			"normal":   deterministicTabSignature(response.Details.Normal),
			"final":    deterministicTabSignature(response.DetailsFinal.Normal),
			"purchase": deterministicTabSignature(response.PurchaseDetails.Normal),
		}

		consultByTab := orchestratePromoConsultByTabs(payloads, signatures, promoService.ConsultV2)

		if !normalHasSnapshot {
			normalConsult := consultByTab["normal"]
			normalAggregate := distributePromoToDetailRowsV2(aggregatePromoByProductForDetailSnapshot(normalConsult), details, promoSnapshotTabSalesOrder, normalConsult)
			response.Details.Normal = injectPromoToOrderItems(response.Details.Normal, normalAggregate)
			response.Details.RewardProducts = buildRewardProducts(normalConsult, buildProductMetaMap(response.Details.Normal), service.OrderRepository.FindProductByID)
			response.Details.PromoRemarksSo = buildFinalRemarks(normalConsult)
			response.Details.FinalRemarks = append([]string{}, response.Details.PromoRemarksSo...)
		}
		if !finalHasSnapshot {
			finalConsult := consultByTab["final"]
			finalAggregate := distributePromoToDetailRowsV2(aggregatePromoByProductForDetailSnapshot(finalConsult), details, promoSnapshotTabFinalOrder, finalConsult)
			response.DetailsFinal.Normal = injectPromoToOrderItems(response.DetailsFinal.Normal, finalAggregate)
			response.DetailsFinal.RewardProducts = buildRewardProducts(finalConsult, buildProductMetaMap(response.DetailsFinal.Normal), service.OrderRepository.FindProductByID)
			response.DetailsFinal.PromoRemarksFinal = buildFinalRemarks(finalConsult)
			response.DetailsFinal.FinalRemarks = append([]string{}, response.DetailsFinal.PromoRemarksFinal...)
		}
		if !purchaseHasSnapshot {
			purchaseConsult := consultByTab["purchase"]
			purchaseAggregate := distributePromoToDetailRowsV2(aggregatePromoByProductForDetailSnapshot(purchaseConsult), details, promoSnapshotTabPurchase, purchaseConsult)
			response.PurchaseDetails.Normal = injectPromoToOrderItems(response.PurchaseDetails.Normal, purchaseAggregate)
			response.PurchaseDetails.RewardProducts = buildRewardProducts(purchaseConsult, buildProductMetaMap(response.PurchaseDetails.Normal), service.OrderRepository.FindProductByID)
			response.PurchaseDetails.PromoRemarksPo = buildFinalRemarks(purchaseConsult)
			response.PurchaseDetails.FinalRemarks = append([]string{}, response.PurchaseDetails.PromoRemarksPo...)
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

	rewards, err := service.OrderRepository.FindReward(RoNo, custID)
	if err != nil {
		return response, err
	}

	for _, reward := range rewards {
		var remark entity.OrderRewardResponse

		if err = structs.Automapper(reward, &remark); err != nil {
			return response, err
		}

		remark.RewardTypeName = remark.GenerateRewardTypeName()

		response.Remarks = append(response.Remarks, remark)
	}

	return response, nil
}

func (service *orderServiceImpl) DetailDiscount(criteria entity.OrderDiscountQuery) (response entity.DiscountCriteria, err error) {
	ro, err := service.OrderRepository.FindDiscountCriteria(criteria.ProID, criteria.OutletID, criteria.OrderDate, criteria.GrossValue)

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

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

		// Map outlet_address1, fallback to outlet_address2 if address1 is nil
		if row.OutletAddress1 != nil {
			vResp.OutletAddress1 = row.OutletAddress1
		} else if row.OutletAddress2 != nil {
			vResp.OutletAddress1 = row.OutletAddress2
		}

		// Map data_source to source
		vResp.Source = MapDataSourceToSource(row.DataSource)

		statusName := vResp.GenerateDataStatusName()
		vResp.DataStatusName = statusName

		payTypeName := vResp.GeneratePayTypeName()
		vResp.PayTypeName = payTypeName

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderServiceImpl) ProformaInvoiceList(dataFilter entity.ProformaInvoiceQueryFilter) (data []entity.ProformaInvoiceListResponse, total int64, lastPage int, err error) {
	realOrders, total, lastPage, err := service.OrderRepository.FindProformaInvoiceList(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range realOrders {
		var vResp entity.ProformaInvoiceListResponse

		// Map basic fields
		vResp.RoNo = row.RoNo
		// Use order_no as fallback if po_no is empty
		if row.PoNo != nil && *row.PoNo != "" {
			vResp.PoNo = row.PoNo
		} else if row.OrderNo != nil {
			vResp.PoNo = row.OrderNo
		}
		if row.OutletID != nil {
			vResp.OutletID = row.OutletID
		}
		if row.OutletCode != nil {
			vResp.OutletCode = *row.OutletCode
		}
		if row.OutletName != nil {
			vResp.OutletName = *row.OutletName
		}
		if row.OutletAddress1 != nil {
			vResp.OutletAddress = *row.OutletAddress1
		}
		if row.SalesmanId != nil {
			vResp.SalesmanId = row.SalesmanId
		}
		if row.SalesmanCode != nil {
			vResp.SalesmanCode = *row.SalesmanCode
		}
		if row.SalesName != nil {
			vResp.SalesmanName = *row.SalesName
		}
		if row.Total != nil {
			vResp.TotalValue = row.Total
		}
		if row.DataStatus != nil {
			vResp.DataStatus = row.DataStatus
		}
		vResp.DataStatusName = "Processed"

		// Map proforma invoice fields
		if row.IsProformaInv != nil {
			vResp.IsProformaInv = row.IsProformaInv
		}
		// ProformaInvNo will be null for now, will be filled when generate proforma invoice is implemented
		vResp.ProformaInvNo = nil

		// Format date to DD/MM/YYYY
		if row.RoDate != nil {
			roDate := row.RoDate.Format("02/01/2006")
			vResp.RoDate = &roDate
		}
		if row.FirstIssueDate != nil {
			firstIssueDate := row.FirstIssueDate.Format("02/01/2006")
			vResp.FirstIssueDate = &firstIssueDate
		}

		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *orderServiceImpl) PrintProformaInvoice(ctx context.Context, request entity.PrintProformaInvoiceRequest, custId string, userId int64) (response entity.PrintProformaInvoiceResponse, err error) {
	// Validasi ro_no tidak kosong
	if len(request.RoNo) == 0 {
		return response, fmt.Errorf("ro_no is required and must contain at least one value")
	}

	var orders []model.OrderList
	var orderDetails []model.OrderDetailRead
	var isCetakUlang bool

	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Validasi dan fetch order data
		var fetchErr error
		orders, fetchErr = service.OrderRepository.FindOrdersByRoNos(txCtx, request.RoNo, custId)
		if fetchErr != nil {
			return fetchErr
		}

		// Validasi semua ro_no ditemukan (bisa ada duplicate, jadi cek unique ro_no)
		foundRoNos := make(map[string]bool)
		for _, order := range orders {
			if order.RoNo != "" {
				foundRoNos[order.RoNo] = true
			}
		}

		var missingRoNos []string
		for _, roNo := range request.RoNo {
			if !foundRoNos[roNo] {
				missingRoNos = append(missingRoNos, roNo)
			}
		}

		if len(missingRoNos) > 0 {
			return fmt.Errorf("some ro_no not found: %v", missingRoNos)
		}

		// Cek apakah sudah pernah generate (untuk determine cetak_ulang)
		for _, order := range orders {
			if order.IsProformaInv != nil && *order.IsProformaInv {
				isCetakUlang = true
				break
			}
		}

		// Fetch order details
		orderDetails, fetchErr = service.OrderRepository.FindOrderDetailsForProforma(txCtx, request.RoNo, custId)
		if fetchErr != nil {
			return fetchErr
		}

		// Update proforma invoice flags
		fetchErr = service.OrderRepository.UpdateProformaInvoiceFlags(txCtx, request.RoNo, custId, userId)
		if fetchErr != nil {
			// Check if it's orders not found error and translate it
			if errors.Is(fetchErr, repository.ErrOrdersNotFound) {
				return fmt.Errorf("orders not found for ro_no: %v", request.RoNo)
			}
			return fetchErr
		}

		return nil
	})

	if err != nil {
		return response, err
	}

	// Mapping response
	response.RoNo = request.RoNo

	// Map header fields dari order pertama
	if len(orders) > 0 {
		service.mapOrderToResponseHeader(&response, orders[0], isCetakUlang)
	}

	// Map products
	response.Products = service.mapOrderDetailsToProducts(orderDetails)

	// Map footer fields (aggregate dari semua order)
	service.calculateFooterValues(&response, orders)

	return response, nil
}

// mapOrderToResponseHeader maps order data to response header fields
func (service *orderServiceImpl) mapOrderToResponseHeader(response *entity.PrintProformaInvoiceResponse, order model.OrderList, isCetakUlang bool) {
	// Use ro_no for no_so (RO is Sales Order number)
	roNo := order.RoNo
	response.NoSo = &roNo
	// Use order_no for no_po (order_no contains PO number)
	if order.OrderNo != nil {
		response.NoPo = order.OrderNo
	}
	response.SalesmanId = order.SalesmanId
	if order.SalesName != nil {
		response.SalesmanName = *order.SalesName
	}
	if order.Notes != nil {
		response.Notes = order.Notes
	}
	if order.InvoiceNo != nil {
		response.NoInvoice = order.InvoiceNo
	}
	if order.InvoiceDate != nil {
		tglInvoice := order.InvoiceDate.Format(constant.DATE_FORMAT_DD_MM_YYYY)
		response.TglInvoice = &tglInvoice
	}
	if order.DueDate != nil {
		tglJatuhTempo := order.DueDate.Format(constant.DATE_FORMAT_DD_MM_YYYY)
		response.TglJatuhTempo = &tglJatuhTempo
	}

	// Generate pay type name
	if order.PayType != nil {
		payTypeName := entity.PayTypeName(*order.PayType)
		response.TypeBayar = payTypeName
	}

	// Map data_source to source, same as DetailV2
	response.Source = MapDataSourceToSource(order.DataSource)

	// Map outlet fields
	response.OutletCode = order.OutletCode
	response.OutletName = order.OutletName
	response.Address1 = order.OutletAddress1
	response.ZipCode = order.ZipCode

	// Set is_proforma_inv dan cetak_ulang
	response.IsProformaInv = true
	response.CetakUlang = isCetakUlang
}

// mapOrderDetailsToProducts maps order details to product list
func (service *orderServiceImpl) mapOrderDetailsToProducts(orderDetails []model.OrderDetailRead) []entity.PrintProformaInvoiceProduct {
	var products []entity.PrintProformaInvoiceProduct

	for _, detail := range orderDetails {
		product := entity.PrintProformaInvoiceProduct{
			ProductCode: detail.ProCode,
			ProductName: detail.ProName,
			Qty1:        detail.Qty1,
			Qty2:        detail.Qty2,
			Qty3:        detail.Qty3,
			UnitId1:     detail.UnitId1,
			UnitId2:     detail.UnitId2,
			UnitId3:     detail.UnitId3,
			SellPrice1:  detail.SellPrice1,
			SellPrice2:  detail.SellPrice2,
			SellPrice3:  detail.SellPrice3,
			DiscValue:   detail.DiscValueFinal,
			Remarks:     detail.Notes,
		}

		// Calculate total and nett_value
		total := service.calculateProductTotal(detail)
		product.Total = &total

		nettValue := service.calculateProductNettValue(total, detail.DiscValueFinal)
		product.NettValue = &nettValue

		// Set default promo values
		product.Promo1 = constant.DEFAULT_PROMO_VALUE
		product.Promo2 = constant.DEFAULT_PROMO_VALUE
		product.Promo3 = constant.DEFAULT_PROMO_VALUE
		product.Promo4 = constant.DEFAULT_PROMO_VALUE
		product.Promo5 = constant.DEFAULT_PROMO_VALUE

		products = append(products, product)
	}

	return products
}

// calculateProductTotal calculates total price for a product (qty * sell_price)
func (service *orderServiceImpl) calculateProductTotal(detail model.OrderDetailRead) float64 {
	var total float64
	if detail.Qty1 != nil && detail.SellPrice1 != nil {
		total += *detail.Qty1 * *detail.SellPrice1
	}
	if detail.Qty2 != nil && detail.SellPrice2 != nil {
		total += *detail.Qty2 * *detail.SellPrice2
	}
	if detail.Qty3 != nil && detail.SellPrice3 != nil {
		total += *detail.Qty3 * *detail.SellPrice3
	}
	return total
}

// calculateProductNettValue calculates nett value (total - disc_value)
func (service *orderServiceImpl) calculateProductNettValue(total float64, discValue *float64) float64 {
	nettValue := total
	if discValue != nil {
		nettValue -= *discValue
	}
	return nettValue
}

// calculateFooterValues calculates and sets footer values by aggregating all orders
func (service *orderServiceImpl) calculateFooterValues(response *entity.PrintProformaInvoiceResponse, orders []model.OrderList) {
	var gross, promotionMoney, discount, vatValue, fakturAmount float64

	for _, order := range orders {
		if order.SubTotal != nil {
			gross += *order.SubTotal
		}
		if order.PromoValue != nil {
			promotionMoney += *order.PromoValue
		}
		if order.DiscValue != nil {
			discount += *order.DiscValue
		}
		if order.VatValue != nil {
			vatValue += *order.VatValue
		}
		if order.TotalFinal != nil {
			fakturAmount += *order.TotalFinal
		}
	}

	response.Gross = &gross
	response.PromotionMoney = &promotionMoney
	response.Discount = &discount
	response.VatValue = &vatValue
	response.FakturAmount = &fakturAmount
	response.PromotionProduct = constant.DEFAULT_PROMO_VALUE
	response.Remark = ""
}

/*
	func (service *orderServiceImpl) UpdateOld(roNo string, request entity.UpdateOrderBody) (err error) {
		c := context.Background()

		// var consultDiscount entity.ConsultDiscountOrderBody
		// if err = structs.Automapper(request, &consultDiscount); err != nil {
		// 	return err
		// }

		// if err = service.ConsultDiscountBeforeStore(&consultDiscount); err != nil {
		// 	return err
		// }

		// if err = structs.Automapper(consultDiscount, &request); err != nil {
		// 	return err
		// }

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
			// DetailIds := []int64{}

			// for _, detail := range request.Details.Normal {
			// 	if detail.OrderDetId != nil {
			// 		DetailIds = append(DetailIds, *detail.OrderDetId)
			// 	}
			// }
			// for _, detail := range request.Details.Promo {
			// 	if detail.OrderDetId != nil {
			// 		DetailIds = append(DetailIds, *detail.OrderDetId)
			// 	}
			// }
			// log.Info(len(DetailIds))
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

			productsModel, err := service.OrderRepository.FindProductByListID(productIDs)
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

				totalQty, err := QtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				roDetModel.Qty = float64(totalQty)

				if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
					roDetModel.OrderDetailID = nil
					err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
					if err != nil {
						return err
					}
				} else {
					// roDetModel.CustId = ""
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

				roDetModel.Qty = float64(totalQty)

				if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
					roDetModel.OrderDetailID = nil
					err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
					if err != nil {
						return err
					}
				} else {
					// roDetModel.CustId = ""
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
*/
func (service *orderServiceImpl) Update(roNo string, request entity.UpdateOrderBody, validateOrderRequest entity.ValidateResponse) (err error) {
	c := context.Background()

	if request.DataStatus != nil && *request.DataStatus == entity.PROCESSED && isValidationResponseEmpty(validateOrderRequest) && service.ValidateOrderRepository != nil {
		statusDecision, decisionErr := service.determineStatusForExistingOrder(roNo, request.CustId)
		if decisionErr != nil {
			return decisionErr
		}
		if guardErr := ensureSalesOrderStatusDecisionAllowed(statusDecision); guardErr != nil {
			return guardErr
		}
		validateOrderRequest = validationResultFromOrderList(model.OrderList{
			ValidateStok:        statusDecision.DataStatus == int64(entity.PROCESSED),
			ValidateCreditLimit: statusDecision.DataStatus == int64(entity.PROCESSED),
			ValidateOverdue:     statusDecision.DataStatus == int64(entity.PROCESSED),
			ValidateOutstanding: statusDecision.DataStatus == int64(entity.PROCESSED),
			ValidateSummary:     statusDecision.DataStatus == int64(entity.PROCESSED),
		})
	}

	var consultDiscount entity.ConsultDiscountOrderBody
	if err = structs.Automapper(request, &consultDiscount); err != nil {
		return err
	}

	if err = service.ConsultDiscountBeforeStore(&consultDiscount); err != nil {
		return err
	}

	if err = structs.Automapper(consultDiscount, &request); err != nil {
		return err
	}

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

	if err = structs.Automapper(request, &Model); err != nil {
		return err
	}

	if (request.DataStatus == nil || *request.DataStatus != entity.CANCELLED) && !isValidationResponseEmpty(validateOrderRequest) {
		// Fetch Outlet Data to get Address if OutletID is present
		var statusDecision salesOrderStatusDecision
		if request.OutletID != nil {
			outletData, err := service.OrderRepository.FindOutletByID(int(*request.OutletID), request.CustId, request.ParentCustId)
			if err != nil {
				return err
			}
			Model.Address1 = outletData.Address1
			statusDecision = determineSalesOrderStatus(validateOrderRequest, outletRulesFromOutletRead(outletData))
		} else {
			ro, err := service.OrderRepository.FindByNo(roNo, request.CustId)
			if err != nil {
				return err
			}
			statusDecision = determineSalesOrderStatus(validateOrderRequest, outletRulesFromOrderList(ro))
		}
		if err = ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
			return err
		}
		Model.DataStatus = &statusDecision.DataStatus
		applyValidationResultToOrderModel(&Model, validateOrderRequest)
	}

	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		if request.DataStatus != nil {
			if *request.DataStatus == entity.CANCELLED { // jika canceled, maka tidak perlu lanjut proses update
				ro, err := service.OrderRepository.FindByNo(roNo, request.CustId)
				if err != nil {
					return err
				}

				details, err := service.OrderRepository.FindDetail(roNo, request.CustId)
				if err != nil {
					return err
				}

				var salesDetailCanceledUpdates []*entity.SalesOrderStockUpdate

				for _, detail := range details {
					salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
						CustID:         request.CustId,
						WhID:           *request.WhId,
						ProID:          int64(detail.ProId),
						StockDate:      *ro.RoDate,
						TrCode:         request.RoNo[0:2],
						TrNo:           request.RoNo,
						QtyOrderBefore: detail.QtyFinal,
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

				return nil
			}
		}

		err = service.OrderRepository.Update(txCtx, roNo, request.CustId, Model)
		if err != nil {
			return err
		}
		//panic("sdsd")

		DetailIds := []int64{}

		log.Info("Get Detail ID")
		for _, detail := range request.Details.Normal {
			if detail.OrderDetId != nil {
				DetailIds = append(DetailIds, *detail.OrderDetId)
			}
		}
		// for _, detail := range request.Details.Promo {
		// 	if detail.OrderDetId != nil {
		// 		DetailIds = append(DetailIds, *detail.OrderDetId)
		// 	}
		// }
		log.Info(len(DetailIds))

		log.Info("Get Master Order")
		ro, err := service.OrderRepository.FindByNo(roNo, request.CustId)
		if err != nil {
			return err
		}

		isMobileNoChange, err := service.isMobileProcessNoMeaningfulDetailChange(roNo, ro, request)
		if err != nil {
			return err
		}
		if isMobileNoChange {
			return nil
		}

		var salesDetailDeletedUpdates []*entity.SalesOrderStockUpdate

		log.Info("Get Deleted Details")
		deletedDetails, err := service.OrderRepository.FindDetailByNotInDetailIDs(DetailIds, roNo, request.CustId)
		if err != nil {
			return err
		}

		for _, deletedDetail := range deletedDetails {
			salesDetailDeletedUpdate := entity.SalesOrderStockUpdate{
				CustID:         request.CustId,
				WhID:           *request.WhId,
				ProID:          int64(deletedDetail.ProId),
				StockDate:      *ro.RoDate,
				TrCode:         request.RoNo[0:2],
				TrNo:           request.RoNo,
				QtyOrderBefore: deletedDetail.QtyFinal,
				QtyOrder:       0,
				UnitPrice:      *deletedDetail.SellPrice1,
				RefDetId:       int64(*deletedDetail.OrderDetailID),
			}
			salesDetailDeletedUpdates = append(salesDetailDeletedUpdates, &salesDetailDeletedUpdate)
		}

		if len(salesDetailDeletedUpdates) > 0 {
			log.Info("Update Stock Deleted Details")
			err = service.StockRepository.SalesStockUpdates(txCtx, salesDetailDeletedUpdates)
			if err != nil {
				return err
			}
		}

		log.Info("Delete Deleted Details")
		err = service.OrderRepository.DeleteDetailNotInIDs(txCtx, roNo, request.CustId, DetailIds)
		if err != nil {
			return err
		}

		var salesOrderStockUpdateEntities []*entity.SalesOrderStockUpdate
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

			if isProcessedDataStatus(Model.DataStatus) { // jika status process, qty final disamakan dengan qty
				detail.Qty1Final = detail.Qty1
				detail.Qty2Final = detail.Qty2
				detail.Qty3Final = detail.Qty3
			}

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo

			// Map based on data_source
			if ro.DataSource != nil && *ro.DataSource == 2 {
				// Mobile (Purchase Order): map qty1/2/3 to qty_po1/2/3 and sell_price1/2/3 to sell_price_po1/2/3
				roDetModel.QtyPo1 = detail.Qty1
				roDetModel.QtyPo2 = detail.Qty2
				roDetModel.QtyPo3 = detail.Qty3
				roDetModel.SellPricePo1 = detail.SellPrice1
				roDetModel.SellPricePo2 = detail.SellPrice2
				roDetModel.SellPricePo3 = detail.SellPrice3
			} else {
				// Web (Sales Order) or null: map to normal qty1/2/3 and sell_price1/2/3
				// Automapper already handles this, but we ensure it's set explicitly
				roDetModel.Qty1 = detail.Qty1
				roDetModel.Qty2 = detail.Qty2
				roDetModel.Qty3 = detail.Qty3
				roDetModel.SellPrice1 = detail.SellPrice1
				roDetModel.SellPrice2 = detail.SellPrice2
				roDetModel.SellPrice3 = detail.SellPrice3
			}

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

			roDetModel.Qty = float64(totalQty)
			roDetModel.QtyFinal = float64(totalQty)

			// Fix Task 4: Refresh Stock Snapshot
			// Fetch current warehouse stock
			if request.WhId != nil {
				currentStock, err := service.StockRepository.GetCurrentStock(txCtx, request.CustId, int64(*request.WhId), int64(detail.ProId))
				if err == nil {
					// Get product for conversion units (if not already available, fetch it)
					product, err := service.OrderRepository.FindProductByID(int(detail.ProId))
					if err == nil {
						stockBreakdown := canonicalAPIStockBreakdown(int(currentStock), int(product.ConvUnit2), int(product.ConvUnit3))
						applyStockBreakdownToPointers(&roDetModel.Qty1Stok, &roDetModel.Qty2Stok, &roDetModel.Qty3Stok, stockBreakdown)
					}
				}
			}

			if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
				roDetModel.OrderDetailID = nil
				err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
				roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate)
				if err != nil {

					return err
				}
				salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
					CustID:         request.CustId,
					WhID:           *request.WhId,
					ProID:          int64(detail.ProId),
					StockDate:      *roDate,
					TrCode:         request.RoNo[0:2],
					TrNo:           request.RoNo,
					QtyOrderBefore: nil,
					QtyOrder:       roDetModel.QtyFinal,
					UnitPrice:      *detail.SellPrice1,
					RefDetId:       int64(*roDetModel.OrderDetailID),
				}
				salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)

				// Sync final fields to ensure consistency
				if err := service.OrderRepository.SyncFinalOrderFields(txCtx, int64(*roDetModel.OrderDetailID)); err != nil {
					return err
				}

			} else {
				roDetModelExist, err := service.OrderRepository.FindDetailByDetailID(*detail.OrderDetId, roNo, request.CustId)
				if err != nil {
					return err
				}

				// roDetModel.CustId = ""
				err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}

				// Sync final fields to ensure consistency
				if err := service.OrderRepository.SyncFinalOrderFields(txCtx, int64(*detail.OrderDetId)); err != nil {
					return err
				}

				roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate)
				if err != nil {
					return err
				}

				salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
					CustID:         request.CustId,
					WhID:           *request.WhId,
					ProID:          int64(detail.ProId),
					StockDate:      *roDate,
					TrCode:         roDetModelExist.RoNo[0:2],
					TrNo:           roDetModelExist.RoNo,
					QtyOrderBefore: roDetModelExist.QtyFinal,
					QtyOrder:       roDetModel.QtyFinal,
					UnitPrice:      *roDetModelExist.SellPrice1,
					RefDetId:       int64(*roDetModelExist.OrderDetailID),
				}
				salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)

			}
		}

		/*
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

				roDetModel.Qty = float64(totalQty)
				roDetModel.QtyFinal = float64(totalQty)

				if detail.OrderDetId == nil || *detail.OrderDetId == 0 {
					roDetModel.OrderDetailID = nil
					err = service.OrderRepository.StoreDetail(txCtx, &roDetModel)
					if err != nil {
						return err
					}
				} else {
					// roDetModel.CustId = ""
					err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
					if err != nil {
						return err
					}

				}
			}
		*/
		if err := service.OrderRepository.DeletePromoDetails(txCtx, roNo, request.CustId); err != nil {
			return err
		}

		details, err := service.OrderRepository.FindDetail(roNo, request.CustId)
		if err != nil {
			return err
		}

		// kembalikan stok, karena proses diatas dihapus semua detail yang promo
		var salesDetailPromoDeleted []*entity.SalesOrderStockUpdate
		for _, detail := range details {
			if detail.ItemType == 2 { // jika promo
				salesDetailCanceledUpdate := entity.SalesOrderStockUpdate{
					CustID:         request.CustId,
					WhID:           *request.WhId,
					ProID:          int64(detail.ProId),
					StockDate:      *ro.RoDate,
					TrCode:         request.RoNo[0:2],
					TrNo:           request.RoNo,
					QtyOrderBefore: detail.Qty,
					QtyOrder:       0,
					UnitPrice:      *detail.SellPrice1,
					RefDetId:       int64(*detail.OrderDetailID),
				}
				salesDetailPromoDeleted = append(salesDetailPromoDeleted, &salesDetailCanceledUpdate)
			}
		}

		if len(salesDetailPromoDeleted) > 0 {
			log.Info("Update Stock Deleted Details")
			err = service.StockRepository.SalesStockUpdates(txCtx, salesDetailPromoDeleted)
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
			gdsDetModel.RoNo = roNo
			gdsDetModel.ItemType = 2

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*Detail.Qty1),
				Qty2:      int(*Detail.Qty2),
				Qty3:      int(*Detail.Qty3),
				ConvUnit2: int(*Detail.ConvUnit2),
				ConvUnit3: int(*Detail.ConvUnit3),
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
			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

			roDate, err := str.ConvertStringTimeToTimeObject(*request.RoDate)
			if err != nil {
				return err
			}

			salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
				CustID:         request.CustId,
				WhID:           *request.WhId,
				ProID:          int64(Detail.ProId),
				StockDate:      *roDate,
				TrCode:         request.RoNo[0:2],
				TrNo:           request.RoNo,
				QtyOrderBefore: nil,
				QtyOrder:       gdsDetModel.QtyFinal,
				UnitPrice:      *Detail.SellPrice1,
				RefDetId:       int64(*gdsDetModel.OrderDetailID),
			}
			salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)

		}

		err = service.StockRepository.SalesStockUpdates(txCtx, salesOrderStockUpdateEntities)
		if err != nil {
			return err
		}

		if err := service.OrderRepository.DeleteRewards(txCtx, roNo, request.CustId); err != nil {
			return err
		}

		for _, rewardRequest := range request.Rewards {
			var reward model.OrderReward

			if err = structs.Automapper(rewardRequest, &reward); err != nil {
				return err
			}

			reward.RoNo = roNo

			if err = service.OrderRepository.StoreReward(txCtx, &reward); err != nil {
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

func (service *orderServiceImpl) UpdateFinal(roNo string, request entity.UpdateOrderDetailFinal, validateOrderRequest entity.ValidateResponse) (err error) {
	c := context.Background()

	ro, err := service.OrderRepository.FindByNo(roNo, request.CustId)
	if err != nil {
		return err
	}

	request.OutletID = ro.OutletID

	for index, detail := range request.Details.Normal {
		detailOrder, err := service.OrderRepository.FindDetailByDetailID(*detail.OrderDetId, roNo, request.CustId)
		if err != nil {
			return err
		}

		request.Details.Normal[index].SellPrice1 = detailOrder.SellPrice1
		request.Details.Normal[index].SellPrice2 = detailOrder.SellPrice2
		request.Details.Normal[index].SellPrice3 = detailOrder.SellPrice3
		request.Details.Normal[index].ConvUnit2 = detailOrder.ConvUnit2
		request.Details.Normal[index].ConvUnit3 = detailOrder.ConvUnit3
		request.Details.Normal[index].Vat = detailOrder.Vat
	}
	var consultDiscount entity.ConsultDiscountOrderBody
	if err = structs.Automapper(request, &consultDiscount); err != nil {
		return err
	}
	if err = structs.Automapper(request.Details, &consultDiscount.Details); err != nil {
		return err
	}
	if err = service.ConsultDiscountBeforeStore(&consultDiscount); err != nil {
		return err
	}
	if err = structs.Automapper(consultDiscount, &request); err != nil {
		return err
	}
	if err = structs.Automapper(consultDiscount.Details, &request.Details); err != nil {
		return err
	}

	// End parse time format YYYY-mm-dd to Rfc339
	var Model model.Order

	if err = structs.Automapper(request, &Model); err != nil {
		return err
	}

	// Fetch Outlet Data to get Address if OutletID is present
	var statusDecision salesOrderStatusDecision
	if request.OutletID != nil {
		outletData, err := service.OrderRepository.FindOutletByID(int(*request.OutletID), request.CustId, request.ParentCustId)
		if err != nil {
			return err
		}
		Model.Address1 = outletData.Address1
		statusDecision = determineSalesOrderStatus(validateOrderRequest, outletRulesFromOutletRead(outletData))
	} else {
		statusDecision = determineSalesOrderStatus(validateOrderRequest, outletRulesFromOrderList(ro))
	}
	if err = ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
		return err
	}
	Model.DataStatus = &statusDecision.DataStatus
	applyValidationResultToOrderModel(&Model, validateOrderRequest)

	Model.CustID = ""
	err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
		basisRows, err := service.StockRepository.GetCancelStockBasis(txCtx, request.CustId, roNo)
		if err != nil {
			return err
		}
		if err := validateFinalOrderStockBasis(basisRows); err != nil {
			return err
		}

		err = service.OrderRepository.Update(txCtx, roNo, request.CustId, Model)
		if err != nil {
			return err
		}
		// DetailIds := []int64{}

		// for _, detail := range request.Details.Normal {
		// 	if detail.OrderDetId != nil {
		// 		DetailIds = append(DetailIds, *detail.OrderDetId)
		// 	}
		// }
		// for _, detail := range request.Details.Promo {
		// 	if detail.OrderDetId != nil {
		// 		DetailIds = append(DetailIds, *detail.OrderDetId)
		// 	}
		// }
		// log.Info(len(DetailIds))
		// if len(DetailIds) > 0 {
		// 	err := service.OrderRepository.DeleteDetailNotInIDs(txCtx, roNo, DetailIds)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		var salesOrderStockUpdateEntities []*entity.SalesOrderStockUpdate

		for _, detail := range request.Details.Normal {
			// parse time format YYYY-mm-dd to Rfc3339

			var roDetModel model.OrderDetail

			err = structs.Automapper(detail, &roDetModel)
			if err != nil {
				return err
			}
			roDetModel.CustId = request.CustId
			roDetModel.RoNo = roNo

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*detail.Qty1Final),
				Qty2:      int(*detail.Qty2Final),
				Qty3:      int(*detail.Qty3Final),
				ConvUnit2: int(*detail.ConvUnit2),
				ConvUnit3: int(*detail.ConvUnit3),
			}

			QtyUnit.DoConversion()

			totalQty, err := QtyUnit.ToTotalQuantity()
			if err != nil {
				return err
			}

			roDetModel.Qty1 = nil
			roDetModel.Qty2 = nil
			roDetModel.Qty3 = nil
			roDetModel.Qty4 = nil
			roDetModel.Qty5 = nil
			roDetModel.QtyFinal = float64(totalQty)

			roDetModelExist, err := service.OrderRepository.FindDetailByDetailID(*detail.OrderDetId, roNo, request.CustId)
			if err != nil {
				return err
			}

			err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
			if err != nil {
				return err
			}

			salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
				CustID:         request.CustId,
				WhID:           *ro.WhId,
				ProID:          int64(detail.ProId),
				StockDate:      *ro.RoDate,
				TrCode:         roDetModelExist.RoNo[0:2],
				TrNo:           roDetModelExist.RoNo,
				QtyOrderBefore: roDetModelExist.QtyFinal,
				QtyOrder:       roDetModel.QtyFinal,
				UnitPrice:      *roDetModelExist.SellPrice1,
				RefDetId:       int64(*roDetModelExist.OrderDetailID),
			}
			salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)
		}

		/*
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
				QtyUnit := &conversion.QtyUnit{
					Qty1:      int(*detail.Qty1Final),
					Qty2:      int(*detail.Qty2Final),
					Qty3:      int(*detail.Qty3Final),
					ConvUnit2: int(*detail.ConvUnit2),
					ConvUnit3: int(*detail.ConvUnit3),
				}

				QtyUnit.DoConversion()

				totalQty, err := QtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				roDetModel.Qty1 = nil
				roDetModel.Qty2 = nil
				roDetModel.Qty3 = nil
				roDetModel.Qty4 = nil
				roDetModel.Qty5 = nil
				roDetModel.QtyFinal = float64(totalQty)

				err = service.OrderRepository.UpdateDetail(txCtx, &roDetModel)
				if err != nil {
					return err
				}
			}
		*/
		if err := service.OrderRepository.DeletePromoDetails(txCtx, roNo, request.CustId); err != nil {
			return err
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
			gdsDetModel.RoNo = roNo
			gdsDetModel.ItemType = 2

			QtyUnit := &conversion.QtyUnit{
				Qty1:      int(*Detail.Qty1),
				Qty2:      int(*Detail.Qty2),
				Qty3:      int(*Detail.Qty3),
				ConvUnit2: int(*Detail.ConvUnit2),
				ConvUnit3: int(*Detail.ConvUnit3),
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
			err = service.OrderRepository.StoreDetail(txCtx, &gdsDetModel)
			if err != nil {
				return err
			}

			roDate, err := str.ConvertStringTimeToTimeObject(ro.RoDate.String())
			if err != nil {
				return err
			}

			salesOrderStockUpdateEntity := entity.SalesOrderStockUpdate{
				CustID:         request.CustId,
				WhID:           *ro.WhId,
				ProID:          int64(Detail.ProId),
				StockDate:      *roDate,
				TrCode:         request.RoNo[0:2],
				TrNo:           request.RoNo,
				QtyOrderBefore: nil,
				QtyOrder:       gdsDetModel.QtyFinal,
				UnitPrice:      *Detail.SellPrice1,
				RefDetId:       int64(*gdsDetModel.OrderDetailID),
			}
			salesOrderStockUpdateEntities = append(salesOrderStockUpdateEntities, &salesOrderStockUpdateEntity)
		}

		err = service.StockRepository.SalesStockUpdates(txCtx, salesOrderStockUpdateEntities)
		if err != nil {
			return err
		}

		if err := service.OrderRepository.DeleteRewards(txCtx, roNo, request.CustId); err != nil {
			return err
		}

		for _, rewardRequest := range request.Rewards {
			var reward model.OrderReward

			if err = structs.Automapper(rewardRequest, &reward); err != nil {
				return err
			}

			reward.RoNo = roNo

			if err = service.OrderRepository.StoreReward(txCtx, &reward); err != nil {
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
	// Find full product details to get Price and VAT info
	product, err := service.OrderRepository.FindProductByID(int(conversionBody.ProductId))
	if err != nil {
		return response, err
	}

	// Detect input format: prioritize qty_final if exists, otherwise use qty
	useFinalFormat := conversionBody.Qty1Final != nil || conversionBody.Qty2Final != nil || conversionBody.Qty3Final != nil

	var qty1, qty2, qty3 int64

	if useFinalFormat {
		// Use qty_final values
		if conversionBody.Qty1Final != nil {
			qty1 = *conversionBody.Qty1Final
		}
		if conversionBody.Qty2Final != nil {
			qty2 = *conversionBody.Qty2Final
		}
		if conversionBody.Qty3Final != nil {
			qty3 = *conversionBody.Qty3Final
		}
	} else {
		// Use qty values (backward compatibility)
		qty1 = conversionBody.Qty1
		qty2 = conversionBody.Qty2
		qty3 = conversionBody.Qty3
	}

	// Normalize using QtyUnit.DoConversion()
	qtyUnit := &conversion.QtyUnit{
		Qty1:      int(qty1),
		Qty2:      int(qty2),
		Qty3:      int(qty3),
		ConvUnit2: int(product.ConvUnit2),
		ConvUnit3: int(product.ConvUnit3),
	}

	qtyUnit.DoConversion()

	// Set response based on input format
	if useFinalFormat {
		qty1Final := int64(qtyUnit.Qty1)
		qty2Final := int64(qtyUnit.Qty2)
		qty3Final := int64(qtyUnit.Qty3)
		response.Qty1Final = &qty1Final
		response.Qty2Final = &qty2Final
		response.Qty3Final = &qty3Final
		response.Qty1 = nil
		response.Qty2 = nil
		response.Qty3 = nil
	} else {
		qty1Result := int64(qtyUnit.Qty1)
		qty2Result := int64(qtyUnit.Qty2)
		qty3Result := int64(qtyUnit.Qty3)
		response.Qty1 = &qty1Result
		response.Qty2 = &qty2Result
		response.Qty3 = &qty3Result
		response.Qty1Final = nil
		response.Qty2Final = nil
		response.Qty3Final = nil
	}

	// Calculate total quantity
	totalQty := (int64(product.ConvUnit2)*int64(product.ConvUnit3))*int64(qtyUnit.Qty3) + (int64(product.ConvUnit2) * int64(qtyUnit.Qty2)) + int64(qtyUnit.Qty1)
	response.TotalQty = totalQty

	// Calculate Financials (Preview)
	// Use SellPrice1 as base price
	price := product.SellPrice1
	grossAmount := float64(totalQty) * price

	response.Price = &price

	var discValue float64 = 0
	var vatValue float64 = 0

	// Calculate Discount if OutletID is provided
	if conversionBody.OutletID != nil {
		disc, _, _ := service.CalculateLineDiscount(custID, parentCustID, int(*conversionBody.OutletID), int(product.ProId), grossAmount)
		discValue = disc
	}

	vatRate := product.Vat
	vatValue = service.CalculateLineVAT(grossAmount, discValue, 0, vatRate)

	response.DiscValue = &discValue
	response.VatValue = &vatValue

	total := grossAmount - discValue + vatValue
	response.Total = &total

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

type cancelOrderStockBasis struct {
	CustID         string
	RoNo           string
	WhID           int64
	ProID          int64
	RefDetID       int64
	QtyOutSmallest float64
	UnitPrice      float64
	StockDate      time.Time
}

func validateCancelTransition(currentStatus int64) bool {
	switch currentStatus {
	case entity.NEED_REVIEW, entity.PROCESSED:
		return true
	default:
		return false
	}
}

func buildCancelStockWriteCommands(basis []cancelOrderStockBasis) []entity.CancelStockWrite {
	commands := make([]entity.CancelStockWrite, 0, len(basis))
	for _, row := range basis {
		if row.QtyOutSmallest <= 0 {
			continue
		}

		commands = append(commands, entity.CancelStockWrite{
			CustID:      row.CustID,
			RoNo:        row.RoNo,
			WhID:        row.WhID,
			ProID:       row.ProID,
			RefDetID:    row.RefDetID,
			QtySmallest: row.QtyOutSmallest,
			UnitPrice:   row.UnitPrice,
		})
	}
	return commands
}

func validateCancelStockBasis(rows []entity.CancelStockBasis) error {
	if len(rows) == 0 {
		return errors.New("order cannot be cancelled because final detail and stock ledger are inconsistent")
	}

	for _, row := range rows {
		if row.IsAmbiguous {
			return errors.New("order cannot be cancelled because final detail and stock ledger are inconsistent")
		}
		if row.IsMissingSource && row.QtyOutSmallest <= 0 {
			return errors.New("order cannot be cancelled because final detail and stock ledger are inconsistent")
		}
		if row.QtyOutstanding < 0 || row.QtyOutSmallest <= 0 {
			return errors.New("order cannot be cancelled because stock basis is invalid")
		}
	}

	return nil
}

func validateFinalOrderStockBasis(rows []entity.CancelStockBasis) error {
	if len(rows) == 0 {
		return errors.New("final order cannot be saved because stock basis is not synchronized")
	}

	for _, row := range rows {
		if row.IsAmbiguous || (row.IsMissingSource && row.QtyOutSmallest == 0) || row.QtyOutstanding < 0 {
			return errors.New("final order cannot be saved because stock basis is not synchronized")
		}
	}

	return nil
}

func (service *orderServiceImpl) BulkUpdateStatus(custId string, request entity.BulkUpdateStatusOrder) (err error) {
	c := context.Background()

	for index := range request.Orders {
		if request.Orders[index].DataStatus == nil {
			return errors.New("data_status is required")
		}

		requestedStatus := *request.Orders[index].DataStatus
		if requestedStatus == entity.PROCESSED {
			statusDecision, err := service.determineStatusForExistingOrder(request.Orders[index].RoNo, custId)
			if err != nil {
				return err
			}
			if err := ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
				return err
			}
			request.Orders[index].DataStatus = &statusDecision.DataStatus
		}

		var orderModel model.Order
		err = structs.Automapper(request.Orders[index], &orderModel)
		if err != nil {
			return err
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			if *request.Orders[index].DataStatus == entity.CANCELLED {
				orderData, err := service.OrderRepository.FindByNo(request.Orders[index].RoNo, custId)
				if err != nil {
					return err
				}

				if orderData.DataStatus == nil {
					return errors.New("current order status is empty")
				}

				if *orderData.DataStatus == entity.CANCELLED {
					return nil
				}

				if !validateCancelTransition(*orderData.DataStatus) {
					return fmt.Errorf("invalid status transition from %d to %d", *orderData.DataStatus, entity.CANCELLED)
				}

				basisRows, err := service.StockRepository.GetCancelStockBasis(txCtx, custId, request.Orders[index].RoNo)
				if err != nil {
					return err
				}

				skipCancelStockWrite := false
				if *orderData.DataStatus == entity.NEED_REVIEW {
					hasOutstandingStock := false
					for _, row := range basisRows {
						if row.QtyOutSmallest > 0 {
							hasOutstandingStock = true
							break
						}
					}
					if !hasOutstandingStock {
						skipCancelStockWrite = true
						log.Warn("BulkUpdateStatus cancel skip stock reversal due empty/no-outstanding basis for need review order: ", request.Orders[index].RoNo)
					} else if err := validateCancelStockBasis(basisRows); err != nil {
						return err
					}
				} else if err := validateCancelStockBasis(basisRows); err != nil {
					return err
				}

				if !skipCancelStockWrite {
					basis := make([]cancelOrderStockBasis, 0, len(basisRows))
					for _, row := range basisRows {
						basis = append(basis, cancelOrderStockBasis{
							CustID:         row.CustID,
							RoNo:           request.Orders[index].RoNo,
							WhID:           row.WhID,
							ProID:          row.ProID,
							RefDetID:       row.RefDetID,
							QtyOutSmallest: row.QtyOutSmallest,
							UnitPrice:      row.UnitPrice,
							StockDate:      row.StockDate,
						})
					}

					commands := buildCancelStockWriteCommands(basis)
					stockDate := time.Now()
					if orderData.RoDate != nil {
						stockDate = *orderData.RoDate
					}

					err = service.StockRepository.CancelSalesStockUpdates(txCtx, request.Orders[index].RoNo, stockDate, commands)
					if err != nil {
						return err
					}
				}
			}

			err := service.OrderRepository.Update(txCtx, request.Orders[index].RoNo, custId, orderModel)
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

func (service *orderServiceImpl) DetailNoCustID(RoNo string, custIDOrigin string, empID *int64) (response entity.OrderResponse, err error) {
	ro, err := service.OrderRepository.FindByNoNoCustID(RoNo, custIDOrigin)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(ro, &response)
	if err != nil {
		return response, err
	}

	if empID != nil {
		reqApprovalDetail, err := service.OrderRepository.FindOrderApprovalRequestDetailByRoAndEmp(RoNo, *empID)
		if err == nil {
			response.OrderApprovalRequestEmpApprovalStatus = reqApprovalDetail.Status
		}
	}

	details, err := service.OrderRepository.FindDetailNoCustID(RoNo, custIDOrigin)
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

		if detailData.Qty == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.Qty = &qtyOld
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

		qtyOrder := int(getValueOrDefault(detailData.Qty, 0))
		qtyFinal := int(getValueOrDefault(detailData.QtyFinal, 0))

		if qtyFinal < qtyOrder {
			detailData.OrderStatus = "Partial Reject"
		}
		if qtyFinal < 1 {
			detailData.OrderStatus = "Rejected"
		}

		if detailData.ItemType == 1 {
			if isActiveDetailForTab(detail, promoSnapshotTabSalesOrder) {
				response.Details.Normal = append(response.Details.Normal, detailData)
			}
		} else {
			response.Details.Promo = append(response.Details.Promo, detailData)
		}
	}
	//final
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

		if detailData.QtyFinal == nil {
			qty1 := getValueOrDefault(detailData.Qty1, 0)
			qty2 := getValueOrDefault(detailData.Qty2, 0)
			qty3 := getValueOrDefault(detailData.Qty3, 0)
			qtyOld, _ := getValueOld(service, details, int64(detailData.ProId), qty1, qty2, qty3)
			detailData.QtyFinal = &qtyOld
		}

		qty := &conversion.Qty{
			Qty:       int(getValueOrDefault(detailData.QtyFinal, 0)),
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

		if detailData.ItemType == 1 {
			if isActiveDetailForTab(detail, promoSnapshotTabFinalOrder) {
				response.DetailsFinal.Normal = append(response.DetailsFinal.Normal, detailData)
			}
		} else {
			response.DetailsFinal.Promo = append(response.DetailsFinal.Promo, detailData)
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

	rewards, err := service.OrderRepository.FindRewardNoCustID(RoNo, custIDOrigin)
	if err != nil {
		return response, err
	}

	for _, reward := range rewards {
		var remark entity.OrderRewardResponse

		if err = structs.Automapper(reward, &remark); err != nil {
			return response, err
		}

		remark.RewardTypeName = remark.GenerateRewardTypeName()

		response.Remarks = append(response.Remarks, remark)
	}

	return response, nil
}

func (service *orderServiceImpl) GetMinimumPriceProduct(request entity.OrderMinimumPriceFilter) (entity.OrderMinimumPriceResp, error) {
	var response entity.OrderMinimumPriceResp
	var err error
	salesman, err := service.OrderRepository.FindSalesman(request.SalesmanId, request.CustId)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(salesman, &response)
	if err != nil {
		return response, err
	}

	minimumPrice, err := service.OrderRepository.FindMinimumPriceActiveByProID(request.ProID, request.CustId)
	if err != nil {
		return response, nil
	}

	var minimumPriceEntity entity.OrderMinimumPriceSettingResp
	err = structs.Automapper(minimumPrice, &minimumPriceEntity)
	if err != nil {
		return response, err
	}
	response.MinimumPrice = &minimumPriceEntity
	return response, nil
}

// UpdateEnhance handles enhanced edit order for 3 different cases:
// Case 1: Purchase Order tab - updates qty_po1/2/3 and sell_price_po1/2/3
// Case 2: Sales Order tab - updates qty1/2/3 and sell_price1/2/3
// Case 3: Final Order tab - updates qty1_final/2_final/3_final and sell_price_final1/2/3
// Also supports adding new products via add_purchase_order, add_sales_order, add_final_order
func (service *orderServiceImpl) UpdateEnhance(ctx context.Context, roNo string, request entity.EditOrderEnhanceBody) error {
	if err := normalizeEnhancePromoFlags(&request); err != nil {
		return err
	}

	// Determine which tab(s) are being used
	hasPurchaseOrder := len(request.PurchaseOrder) > 0 || len(request.AddPurchaseOrder) > 0
	hasSalesOrder := len(request.SalesOrder) > 0 || len(request.AddSalesOrder) > 0
	hasFinalOrder := len(request.FinalOrder) > 0 || len(request.AddFinalOrder) > 0

	// Validate that at least one tab is being edited
	if !hasPurchaseOrder && !hasSalesOrder && !hasFinalOrder {
		return errors.New("at least one of purchase_order, sales_order, final_order or their add_ variants must be provided")
	}

	return service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		ro, err := service.OrderRepository.FindByNo(roNo, request.CustId)
		if err != nil {
			return fmt.Errorf("failed to find order %s: %w", roNo, err)
		}

		shouldApplyProcessedEffects := true

		needSalesRecompute := false
		needFinalRecompute := false
		needPurchaseRecompute := false
		var stockUpdates []*entity.SalesOrderStockUpdate
		explicitPromoOverrides := map[int64]promoFlagOverride{}

		if hasFinalOrder && shouldApplyProcessedEffects {
			basisRows, basisErr := service.StockRepository.GetCancelStockBasis(txCtx, request.CustId, roNo)
			if basisErr != nil {
				return basisErr
			}
			if basisValidationErr := validateFinalOrderStockBasis(basisRows); basisValidationErr != nil {
				return basisValidationErr
			}
		}

		// Case 1: Purchase Order tab
		for index, detail := range request.PurchaseOrder {
			existingDetail, err := service.OrderRepository.FindOrderDetailByDetailID(detail.OrderDetailId, request.CustId)
			if err != nil {
				return fmt.Errorf("failed to find order detail at index %d (PurchaseOrder): %w", index, err)
			}

			updates := make(map[string]interface{})
			sellPrice1 := getValueOrDefault(existingDetail.SellPrice1, 0)
			// Cascade price updates to SO and Final
			if detail.SellPricePo1 != nil {
				sellPrice1 = *detail.SellPricePo1
				updates["sell_price_po1"] = *detail.SellPricePo1
				updates["sell_price1"] = *detail.SellPricePo1
				updates["sell_price_final1"] = *detail.SellPricePo1
			}
			if detail.SellPricePo2 != nil {
				updates["sell_price_po2"] = *detail.SellPricePo2
				updates["sell_price2"] = *detail.SellPricePo2
				updates["sell_price_final2"] = *detail.SellPricePo2
			}
			if detail.SellPricePo3 != nil {
				updates["sell_price_po3"] = *detail.SellPricePo3
				updates["sell_price3"] = *detail.SellPricePo3
				updates["sell_price_final3"] = *detail.SellPricePo3
			}
			if detail.DiscPo != nil {
				updates["disc_po"] = *detail.DiscPo
				updates["disc_value"] = *detail.DiscPo
				updates["disc_value_final"] = *detail.DiscPo
			}
			if detail.VatValuePo != nil {
				updates["vat_value_po"] = *detail.VatValuePo
				updates["vat_value"] = *detail.VatValuePo
				updates["vat_value_final"] = *detail.VatValuePo
			}

			if detail.QtyPo1 != nil || detail.QtyPo2 != nil || detail.QtyPo3 != nil {
				totalQty, qtyUnit, err := calculateNormalizedQty(
					getQtyValue(detail.QtyPo1, existingDetail.QtyPo1),
					getQtyValue(detail.QtyPo2, existingDetail.QtyPo2),
					getQtyValue(detail.QtyPo3, existingDetail.QtyPo3),
					*existingDetail.MpConvUnit2, *existingDetail.MpConvUnit3,
				)
				if err != nil {
					return fmt.Errorf("failed to calculate normalized qty at index %d (PurchaseOrder): %w", index, err)
				}
				// Cascade qty updates to SO and Final
				updates["qty_po"] = float64(totalQty)
				updates["qty_po1"] = float64(qtyUnit.Qty1)
				updates["qty_po2"] = float64(qtyUnit.Qty2)
				updates["qty_po3"] = float64(qtyUnit.Qty3)

				updates["qty"] = float64(totalQty)
				updates["qty1"] = float64(qtyUnit.Qty1)
				updates["qty2"] = float64(qtyUnit.Qty2)
				updates["qty3"] = float64(qtyUnit.Qty3)

				updates["qty_final"] = float64(totalQty)
				updates["qty1_final"] = float64(qtyUnit.Qty1)
				updates["qty2_final"] = float64(qtyUnit.Qty2)
				updates["qty3_final"] = float64(qtyUnit.Qty3)

				if shouldApplyProcessedEffects && ro.WhId != nil && ro.RoDate != nil && existingDetail.OrderDetailID != nil {
					var trCode string
					if len(roNo) > 2 {
						trCode = roNo[0:2]
					}
					stockUpdate := &entity.SalesOrderStockUpdate{
						CustID:         request.CustId,
						WhID:           int64(*ro.WhId),
						ProID:          int64(existingDetail.ProId),
						StockDate:      *ro.RoDate,
						TrCode:         trCode,
						TrNo:           roNo,
						QtyOrderBefore: existingDetail.QtyFinal,
						QtyOrder:       float64(totalQty),
						UnitPrice:      sellPrice1,
						RefDetId:       int64(*existingDetail.OrderDetailID),
					}
					stockUpdates = append(stockUpdates, stockUpdate)
				}
			}

			if len(updates) > 0 {
				if err := service.OrderRepository.UpdateDetailPartial(txCtx, detail.OrderDetailId, request.CustId, updates); err != nil {
					return fmt.Errorf("failed to update detail partial at index %d (PurchaseOrder): %w", index, err)
				}
				needPurchaseRecompute = true
				needSalesRecompute = true
				needFinalRecompute = true
			}
		}

		// Case 2: Sales Order tab
		for index, detail := range request.SalesOrder {
			existingDetail, err := service.OrderRepository.FindOrderDetailByDetailID(detail.OrderDetailId, request.CustId)
			if err != nil {
				return fmt.Errorf("failed to find order detail at index %d (SalesOrder): %w", index, err)
			}

			updates := make(map[string]interface{})

			sellPrice1 := getValueOrDefault(existingDetail.SellPrice1, 0)
			sellPrice2 := getValueOrDefault(existingDetail.SellPrice2, 0)
			sellPrice3 := getValueOrDefault(existingDetail.SellPrice3, 0)
			if detail.SellPrice1 != nil {
				sellPrice1 = *detail.SellPrice1
				updates["sell_price1"] = sellPrice1
				updates["sell_price_final1"] = sellPrice1
			}
			if detail.SellPrice2 != nil {
				sellPrice2 = *detail.SellPrice2
				updates["sell_price2"] = sellPrice2
				updates["sell_price_final2"] = sellPrice2
			}
			if detail.SellPrice3 != nil {
				sellPrice3 = *detail.SellPrice3
				updates["sell_price3"] = sellPrice3
				updates["sell_price_final3"] = sellPrice3
			}

			qty1 := getQtyValue(detail.Qty1, existingDetail.Qty1)
			qty2 := getQtyValue(detail.Qty2, existingDetail.Qty2)
			qty3 := getQtyValue(detail.Qty3, existingDetail.Qty3)
			var totalQty int

			if detail.Qty1 != nil || detail.Qty2 != nil || detail.Qty3 != nil {
				var qtyUnit *conversion.QtyUnit
				totalQty, qtyUnit, err = calculateNormalizedQty(
					qty1, qty2, qty3,
					*existingDetail.MpConvUnit2, *existingDetail.MpConvUnit3,
				)
				if err != nil {
					return fmt.Errorf("failed to calculate normalized qty at index %d (SalesOrder): %w", index, err)
				}

				updates["qty"] = float64(totalQty)
				updates["qty1"] = float64(qtyUnit.Qty1)
				updates["qty2"] = float64(qtyUnit.Qty2)
				updates["qty3"] = float64(qtyUnit.Qty3)

				updates["qty_final"] = float64(totalQty)
				updates["qty1_final"] = float64(qtyUnit.Qty1)
				updates["qty2_final"] = float64(qtyUnit.Qty2)
				updates["qty3_final"] = float64(qtyUnit.Qty3)

				if shouldApplyProcessedEffects && ro.WhId != nil && ro.RoDate != nil {
					qtyFinalBefore := existingDetail.QtyFinal
					var trCode string
					if len(roNo) > 2 {
						trCode = roNo[0:2]
					}
					stockUpdate := &entity.SalesOrderStockUpdate{
						CustID:         request.CustId,
						WhID:           int64(*ro.WhId),
						ProID:          int64(existingDetail.ProId),
						StockDate:      *ro.RoDate,
						TrCode:         trCode,
						TrNo:           roNo,
						QtyOrderBefore: qtyFinalBefore,
						QtyOrder:       float64(totalQty),
						UnitPrice:      sellPrice1,
						RefDetId:       int64(*existingDetail.OrderDetailID),
					}
					stockUpdates = append(stockUpdates, stockUpdate)
				}
			} else {
				totalQty = int(getValueOrDefault(existingDetail.Qty, 0))
			}

			amount := (float64(qty1) * sellPrice1) + (float64(qty2) * sellPrice2) + (float64(qty3) * sellPrice3)
			promoValueSales := getValueOrDefault(existingDetail.PromoValue, 0)
			promoValueFinal := getValueOrDefault(existingDetail.PromoValueFinal, 0)
			updates["promo_value"] = promoValueSales
			updates["promo_value_final"] = promoValueFinal
			amountAfterPromo := amount - promoValueSales

			// Refresh Stock Snapshot
			if ro.WhId != nil {
				currentStock, err := service.StockRepository.GetCurrentStock(txCtx, request.CustId, int64(*ro.WhId), int64(existingDetail.ProId))
				if err == nil {
					productForStock, err := service.OrderRepository.FindProductByID(existingDetail.ProId)
					if err == nil {
						stockBreakdown := canonicalAPIStockBreakdown(int(currentStock), int(productForStock.ConvUnit2), int(productForStock.ConvUnit3))

						updates["qty1_stok"] = stockBreakdown.Qty1
						updates["qty2_stok"] = stockBreakdown.Qty2
						updates["qty3_stok"] = stockBreakdown.Qty3
					}
				}
			}

			// Calculate Discount
			discValue := 0.0
			if ro.OutletID != nil {
				dVal, _, _ := service.CalculateLineDiscount(request.CustId, request.ParentCustId, int(*ro.OutletID), existingDetail.ProId, amountAfterPromo)
				discValue = dVal
			}
			updates["disc_value"] = discValue
			updates["disc_value_final"] = discValue

			// Calculate VAT
			vatPercent := getValueOrDefault(existingDetail.Vat, 0)
			amountAfterDisc := amountAfterPromo - discValue
			vatValue := service.CalculateLineVAT(amountAfterDisc, 0, 0, vatPercent)

			updates["vat_value"] = vatValue
			updates["vat_value_final"] = vatValue

			finalAmount := amountAfterDisc + vatValue
			updates["amount"] = finalAmount
			updates["amount_final"] = finalAmount

			if detail.IsProductPromotionSo != nil {
				updates["is_product_promotion_so"] = *detail.IsProductPromotionSo
				override := explicitPromoOverrides[detail.OrderDetailId]
				override.SalesOrder = boolPtr(*detail.IsProductPromotionSo)
				explicitPromoOverrides[detail.OrderDetailId] = override
			}

			if len(updates) > 0 {
				if err := service.OrderRepository.UpdateDetailPartial(txCtx, detail.OrderDetailId, request.CustId, updates); err != nil {
					return fmt.Errorf("failed to update detail partial at index %d (SalesOrder): %w", index, err)
				}
				needSalesRecompute = true
				needFinalRecompute = true
			}
		}

		// Case 3: Final Order tab
		finalOrderDetailIDs := make([]int64, 0, len(request.FinalOrder))
		for _, detail := range request.FinalOrder {
			finalOrderDetailIDs = append(finalOrderDetailIDs, detail.OrderDetailId)
		}

		finalOrderDetails, err := service.OrderRepository.FindOrderDetailsByIDs(finalOrderDetailIDs, request.CustId)
		if err != nil {
			return fmt.Errorf("failed to fetch order details (FinalOrder): %w", err)
		}

		finalOrderDetailMap := make(map[int64]model.OrderDetailRead, len(finalOrderDetails))
		for _, existingDetail := range finalOrderDetails {
			if existingDetail.OrderDetailID == nil {
				continue
			}

			finalOrderDetailMap[int64(*existingDetail.OrderDetailID)] = existingDetail
		}

		for index, detail := range request.FinalOrder {
			existingDetail, ok := finalOrderDetailMap[detail.OrderDetailId]
			if !ok {
				return fmt.Errorf("order detail not found at index %d (FinalOrder)", index)
			}

			updates := make(map[string]interface{})
			applyUpdateIfNotNull(updates, "sell_price_final1", detail.SellPriceFinal1)
			applyUpdateIfNotNull(updates, "sell_price_final2", detail.SellPriceFinal2)
			applyUpdateIfNotNull(updates, "sell_price_final3", detail.SellPriceFinal3)

			finalQty1 := getQtyValue(detail.Qty1Final, existingDetail.Qty1Final)
			finalQty2 := getQtyValue(detail.Qty2Final, existingDetail.Qty2Final)
			finalQty3 := getQtyValue(detail.Qty3Final, existingDetail.Qty3Final)
			finalPrice1 := getValueOrDefault(existingDetail.SellPriceFinal1, getValueOrDefault(existingDetail.SellPrice1, 0))
			finalPrice2 := getValueOrDefault(existingDetail.SellPriceFinal2, getValueOrDefault(existingDetail.SellPrice2, 0))
			finalPrice3 := getValueOrDefault(existingDetail.SellPriceFinal3, getValueOrDefault(existingDetail.SellPrice3, 0))
			if detail.SellPriceFinal1 != nil {
				finalPrice1 = *detail.SellPriceFinal1
			}
			if detail.SellPriceFinal2 != nil {
				finalPrice2 = *detail.SellPriceFinal2
			}
			if detail.SellPriceFinal3 != nil {
				finalPrice3 = *detail.SellPriceFinal3
			}

			if detail.Qty1Final != nil || detail.Qty2Final != nil || detail.Qty3Final != nil {
				totalQty, qtyUnit, err := calculateNormalizedQty(
					finalQty1,
					finalQty2,
					finalQty3,
					*existingDetail.MpConvUnit2, *existingDetail.MpConvUnit3,
				)
				if err != nil {
					return fmt.Errorf("failed to calculate normalized qty at index %d (FinalOrder): %w", index, err)
				}
				updates["qty_final"] = float64(totalQty)
				updates["qty1_final"] = float64(qtyUnit.Qty1)
				updates["qty2_final"] = float64(qtyUnit.Qty2)
				updates["qty3_final"] = float64(qtyUnit.Qty3)

				if shouldApplyProcessedEffects && ro.WhId != nil && ro.RoDate != nil && existingDetail.OrderDetailID != nil {
					var trCode string
					if len(roNo) > 2 {
						trCode = roNo[0:2]
					}
					stockUpdate := &entity.SalesOrderStockUpdate{
						CustID:         request.CustId,
						WhID:           int64(*ro.WhId),
						ProID:          int64(existingDetail.ProId),
						StockDate:      *ro.RoDate,
						TrCode:         trCode,
						TrNo:           roNo,
						QtyOrderBefore: existingDetail.QtyFinal,
						QtyOrder:       float64(totalQty),
						UnitPrice:      getValueOrDefault(existingDetail.SellPrice1, 0),
						RefDetId:       int64(*existingDetail.OrderDetailID),
					}
					stockUpdates = append(stockUpdates, stockUpdate)
				}
			}

			finalAmount := (float64(finalQty1) * finalPrice1) + (float64(finalQty2) * finalPrice2) + (float64(finalQty3) * finalPrice3)
			promoValueFinal := getValueOrDefault(existingDetail.PromoValueFinal, 0)
			updates["promo_value_final"] = promoValueFinal
			amountAfterPromo := finalAmount - promoValueFinal
			discValue := 0.0
			if ro.OutletID != nil {
				dVal, _, _ := service.CalculateLineDiscount(request.CustId, request.ParentCustId, int(*ro.OutletID), existingDetail.ProId, amountAfterPromo)
				discValue = dVal
			}
			updates["disc_value_final"] = discValue
			vatPercent := getValueOrDefault(existingDetail.Vat, 0)
			vatValue := service.CalculateLineVAT(amountAfterPromo-discValue, 0, 0, vatPercent)
			updates["vat_value_final"] = vatValue
			updates["amount_final"] = amountAfterPromo - discValue + vatValue

			if detail.IsProductPromotionFinal != nil {
				updates["is_product_promotion_final"] = *detail.IsProductPromotionFinal
				override := explicitPromoOverrides[detail.OrderDetailId]
				override.FinalOrder = boolPtr(*detail.IsProductPromotionFinal)
				explicitPromoOverrides[detail.OrderDetailId] = override
			}

			if len(updates) > 0 {
				if err := service.OrderRepository.UpdateDetailPartial(txCtx, detail.OrderDetailId, request.CustId, updates); err != nil {
					return fmt.Errorf("failed to update detail partial at index %d (FinalOrder): %w", index, err)
				}
				needFinalRecompute = true
			}
		}

		// Add new products - Case 1: Purchase Order tab (cascades to PO, SO, and Final)
		for _, addDetail := range request.AddPurchaseOrder {
			if ro.WhId == nil {
				return fmt.Errorf("warehouse is required for add_purchase_order")
			}
			if ro.RoDate == nil {
				return fmt.Errorf("ro_date is required for add_purchase_order")
			}
			stockUpdate, err := service.createOrderDetailFromPurchaseOrder(txCtx, roNo, request.CustId, int64(*ro.WhId), *ro.RoDate, addDetail)
			if err != nil {
				return err
			}
			if shouldApplyProcessedEffects {
				stockUpdates = append(stockUpdates, stockUpdate)
			}
			needPurchaseRecompute = true
			needSalesRecompute = true
			needFinalRecompute = true
		}

		// Add new products - Case 2: Sales Order tab (cascades to SO and Final)
		for _, addDetail := range request.AddSalesOrder {
			stockUpdate, err := service.createOrderDetailFromSalesOrder(txCtx, roNo, request.CustId, int64(*ro.WhId), *ro.RoDate, addDetail)
			if err != nil {
				return err
			}
			if shouldApplyProcessedEffects {
				stockUpdates = append(stockUpdates, stockUpdate)
			}
			if addDetail.IsProductPromotionSo != nil {
				override := explicitPromoOverrides[stockUpdate.RefDetId]
				override.SalesOrder = boolPtr(*addDetail.IsProductPromotionSo)
				explicitPromoOverrides[stockUpdate.RefDetId] = override
			}
			needSalesRecompute = true
			needFinalRecompute = true
		}

		// Add new products - Case 3: Final Order tab (only Final fields)
		for _, addDetail := range request.AddFinalOrder {
			stockUpdate, err := service.createOrderDetailFromFinalOrder(txCtx, roNo, request.CustId, int64(*ro.WhId), *ro.RoDate, addDetail)
			if err != nil {
				return err
			}
			if shouldApplyProcessedEffects {
				stockUpdates = append(stockUpdates, stockUpdate)
			}
			if addDetail.IsProductPromotionFinal != nil {
				override := explicitPromoOverrides[stockUpdate.RefDetId]
				override.FinalOrder = boolPtr(*addDetail.IsProductPromotionFinal)
				explicitPromoOverrides[stockUpdate.RefDetId] = override
			}
			needFinalRecompute = true
		}

		updatedDetails, err := service.OrderRepository.FindOrderDetailsForProforma(txCtx, []string{roNo}, request.CustId)
		if err != nil {
			return fmt.Errorf("failed to fetch details for header recalculation: %w", err)
		}

		rewardPromoTab := promoSnapshotTab("")
		switch {
		case needFinalRecompute:
			rewardPromoTab = promoSnapshotTabFinalOrder
		case needSalesRecompute:
			rewardPromoTab = promoSnapshotTabSalesOrder
		case needPurchaseRecompute:
			rewardPromoTab = promoSnapshotTabPurchase
		}
		if rewardPromoTab != "" {
			rewardStockUpdates, rewardErr := service.syncRewardProductState(txCtx, ro, roNo, request.CustId, request.ParentCustId, updatedDetails, rewardPromoTab)
			if rewardErr != nil {
				return fmt.Errorf("failed to sync reward product state: %w", rewardErr)
			}
			if shouldApplyProcessedEffects {
				stockUpdates = append(stockUpdates, rewardStockUpdates...)
			}
			updatedDetails, err = service.OrderRepository.FindOrderDetailsForProforma(txCtx, []string{roNo}, request.CustId)
			if err != nil {
				return fmt.Errorf("failed to refresh details after reward sync: %w", err)
			}
		}

		headerUpdate := model.Order{}
		if needSalesRecompute {
			salesHeaderUpdate, recomputeErr := service.recomputePromoStateForTab(txCtx, ro, request.CustId, request.ParentCustId, updatedDetails, promoSnapshotTabSalesOrder, explicitPromoOverrides)
			if recomputeErr != nil {
				return fmt.Errorf("failed to recompute sales promo snapshot: %w", recomputeErr)
			}
			headerUpdate.SubTotal = salesHeaderUpdate.SubTotal
			headerUpdate.DiscValue = salesHeaderUpdate.DiscValue
			headerUpdate.VatValue = salesHeaderUpdate.VatValue
			headerUpdate.Total = salesHeaderUpdate.Total
			headerUpdate.PromoValue = salesHeaderUpdate.PromoValue
			headerUpdate.PromoRemarksSo = salesHeaderUpdate.PromoRemarksSo
		}
		if needFinalRecompute {
			finalHeaderUpdate, recomputeErr := service.recomputePromoStateForTab(txCtx, ro, request.CustId, request.ParentCustId, updatedDetails, promoSnapshotTabFinalOrder, explicitPromoOverrides)
			if recomputeErr != nil {
				return fmt.Errorf("failed to recompute final promo snapshot: %w", recomputeErr)
			}
			headerUpdate.SubTotalFinal = finalHeaderUpdate.SubTotalFinal
			headerUpdate.DiscValueFinal = finalHeaderUpdate.DiscValueFinal
			headerUpdate.VatValueFinal = finalHeaderUpdate.VatValueFinal
			headerUpdate.TotalFinal = finalHeaderUpdate.TotalFinal
			headerUpdate.PromoValueFinal = finalHeaderUpdate.PromoValueFinal
			headerUpdate.PromoRemarksFinal = finalHeaderUpdate.PromoRemarksFinal
		}
		if needPurchaseRecompute {
			purchaseHeaderUpdate, recomputeErr := service.recomputePromoStateForTab(txCtx, ro, request.CustId, request.ParentCustId, updatedDetails, promoSnapshotTabPurchase, explicitPromoOverrides)
			if recomputeErr != nil {
				return fmt.Errorf("failed to recompute purchase promo snapshot: %w", recomputeErr)
			}
			headerUpdate.PromoRemarksPo = purchaseHeaderUpdate.PromoRemarksPo
		}

		if !needSalesRecompute && !needFinalRecompute {
			var (
				subTotalFinal  float64
				discValueFinal float64
				vatValueFinal  float64
				vatPercent     float64
			)

			for i, detail := range updatedDetails {
				if i == 0 && detail.Vat != nil {
					vatPercent = *detail.Vat
				}

				q1 := getValueOrDefault(detail.Qty1Final, 0)
				p1 := getValueOrDefault(detail.SellPriceFinal1, 0)
				q2 := getValueOrDefault(detail.Qty2Final, 0)
				p2 := getValueOrDefault(detail.SellPriceFinal2, 0)
				q3 := getValueOrDefault(detail.Qty3Final, 0)
				p3 := getValueOrDefault(detail.SellPriceFinal3, 0)

				subTotalFinal += (q1 * p1) + (q2 * p2) + (q3 * p3)
				discValueFinal += getValueOrDefault(detail.DiscValueFinal, 0)
				vatValueFinal += getValueOrDefault(detail.VatValueFinal, 0)
			}

			totalFinal := subTotalFinal - discValueFinal + vatValueFinal
			headerUpdate.SubTotal = &subTotalFinal
			headerUpdate.SubTotalFinal = &subTotalFinal
			headerUpdate.DiscValue = &discValueFinal
			headerUpdate.DiscValueFinal = &discValueFinal
			headerUpdate.Vat = &vatPercent
			headerUpdate.VatValue = &vatValueFinal
			headerUpdate.VatValueFinal = &vatValueFinal
			headerUpdate.Total = &totalFinal
			headerUpdate.TotalFinal = &totalFinal
		}

		statusDecision, err := service.determineStatusForOrderProjection(ro, request.CustId, updatedDetails, headerUpdate)
		if err != nil {
			return err
		}
		if err := ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
			return err
		}

		headerUpdate.DataStatus = &statusDecision.DataStatus
		headerUpdate.UpdatedBy = &request.UpdatedBy

		if err := service.OrderRepository.Update(txCtx, roNo, request.CustId, headerUpdate); err != nil {
			return fmt.Errorf("failed to update order header: %w", err)
		}

		if statusDecision.DataStatus == int64(entity.PROCESSED) && len(stockUpdates) > 0 {
			if err := service.StockRepository.SalesStockUpdates(txCtx, stockUpdates); err != nil {
				return fmt.Errorf("failed to update warehouse stock: %w", err)
			}
		}

		return nil
	})
}

func (service *orderServiceImpl) ProcessEnhanceWithoutProductEdit(ctx context.Context, roNo string, custId string, updatedBy int64) error {
	return service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		ro, err := service.OrderRepository.FindByNo(roNo, custId)
		if err != nil {
			return fmt.Errorf("failed to find order %s: %w", roNo, err)
		}

		statusDecision := determineSalesOrderStatus(validationResultFromOrderList(ro), outletRulesFromOrderList(ro))
		if err := ensureSalesOrderStatusDecisionAllowed(statusDecision); err != nil {
			return err
		}

		headerUpdate := model.Order{
			DataStatus: &statusDecision.DataStatus,
			UpdatedBy:  &updatedBy,
		}

		if err := service.OrderRepository.Update(txCtx, roNo, custId, headerUpdate); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	})
}

func normalizeEnhancePromoFlags(request *entity.EditOrderEnhanceBody) error {
	if len(request.PurchaseDetails) > 0 {
		request.PurchaseOrder = append(request.PurchaseOrder, request.PurchaseDetails...)
	}
	if len(request.SalesOrderDetails) > 0 {
		request.SalesOrder = append(request.SalesOrder, request.SalesOrderDetails...)
	}
	if len(request.FinalOrderDetails) > 0 {
		request.FinalOrder = append(request.FinalOrder, request.FinalOrderDetails...)
	}
	if len(request.AddPurchaseDetails) > 0 {
		request.AddPurchaseOrder = append(request.AddPurchaseOrder, request.AddPurchaseDetails...)
	}

	for i := range request.SalesOrder {
		detail := &request.SalesOrder[i]
		if detail.IsProductPromotion == nil {
			continue
		}
		if detail.IsProductPromotionSo != nil && *detail.IsProductPromotionSo != *detail.IsProductPromotion {
			return fmt.Errorf("sales_order[%d]: is_product_promotion conflicts with is_product_promotion_so", i)
		}
		if detail.IsProductPromotionSo == nil {
			detail.IsProductPromotionSo = boolPtr(*detail.IsProductPromotion)
		}
	}

	for i := range request.FinalOrder {
		detail := &request.FinalOrder[i]
		if detail.IsProductPromotion == nil {
			continue
		}
		if detail.IsProductPromotionFinal != nil && *detail.IsProductPromotionFinal != *detail.IsProductPromotion {
			return fmt.Errorf("final_order[%d]: is_product_promotion conflicts with is_product_promotion_final", i)
		}
		if detail.IsProductPromotionFinal == nil {
			detail.IsProductPromotionFinal = boolPtr(*detail.IsProductPromotion)
		}
	}

	return nil
}

// calculateNormalizedQty normalizes qty1/2/3 and returns total quantity
func calculateNormalizedQty(qty1, qty2, qty3, convUnit2, convUnit3 int) (int, *conversion.QtyUnit, error) {
	qtyUnit := &conversion.QtyUnit{
		Qty1:      qty1,
		Qty2:      qty2,
		Qty3:      qty3,
		ConvUnit2: convUnit2,
		ConvUnit3: convUnit3,
	}
	totalQty, err := qtyUnit.ToTotalQuantity()
	if err != nil {
		return 0, nil, err
	}
	return totalQty, qtyUnit, nil
}

// getQtyValue returns new value if provided, otherwise existing value
func getQtyValue(newVal, existingVal *float64) int {
	if newVal != nil {
		return int(*newVal)
	}
	return int(getValueOrDefault(existingVal, 0))
}

// applyUpdateIfNotNull is a helper to build partial update map
func applyUpdateIfNotNull(updates map[string]interface{}, key string, val *float64) {
	if val != nil {
		updates[key] = *val
	}
}

// createOrderDetailFromPurchaseOrder creates a new order detail from Purchase Order tab
// Cascades values to PO, SO, and Final fields
func (service *orderServiceImpl) createOrderDetailFromPurchaseOrder(ctx context.Context, roNo string, custId string, whId int64, roDate time.Time, addDetail entity.AddPurchaseOrderDetail) (*entity.SalesOrderStockUpdate, error) {
	product, err := service.OrderRepository.FindProductByID(int(addDetail.ProId))
	if err != nil {
		return nil, fmt.Errorf("product with id %d not found: %w", addDetail.ProId, err)
	}

	requestedQtyPo1 := addDetail.QtyPo1
	requestedQtyPo2 := addDetail.QtyPo2
	requestedQtyPo3 := addDetail.QtyPo3

	qtyPo1 := requestedQtyPo1
	qtyPo2 := requestedQtyPo2
	qtyPo3 := requestedQtyPo3
	qty1Stok := 0.0
	qty2Stok := 0.0
	qty3Stok := 0.0

	if service.StockRepository != nil {
		currentStock, stockErr := service.StockRepository.GetCurrentStock(ctx, custId, whId, addDetail.ProId)
		if stockErr != nil {
			log.Warnf("createOrderDetailFromPurchaseOrder stock lookup failed ro=%s cust=%s pro_id=%d wh_id=%d: %v", roNo, custId, addDetail.ProId, whId, stockErr)
		} else {
			stockBreakdown := canonicalAPIStockBreakdown(int(currentStock), int(product.ConvUnit2), int(product.ConvUnit3))
			qty1Stok = stockBreakdown.Qty1
			qty2Stok = stockBreakdown.Qty2
			qty3Stok = stockBreakdown.Qty3
			qtyPo1 = math.Min(requestedQtyPo1, stockBreakdown.Qty1)
			qtyPo2 = math.Min(requestedQtyPo2, stockBreakdown.Qty2)
			qtyPo3 = math.Min(requestedQtyPo3, stockBreakdown.Qty3)
		}
	}

	unitId1 := product.UnitId1
	if strings.TrimSpace(addDetail.UnitId1) != "" {
		unitId1 = addDetail.UnitId1
	}
	unitId2 := product.UnitId2
	if strings.TrimSpace(addDetail.UnitId2) != "" {
		unitId2 = addDetail.UnitId2
	}
	unitId3 := product.UnitId3
	if strings.TrimSpace(addDetail.UnitId3) != "" {
		unitId3 = addDetail.UnitId3
	}

	totalQty, _, err := calculateNormalizedQty(int(qtyPo1), int(qtyPo2), int(qtyPo3), int(product.ConvUnit2), int(product.ConvUnit3))
	if err != nil {
		return nil, err
	}

	isProductPromotionPo := false
	if addDetail.IsProductPromotionPo != nil {
		isProductPromotionPo = *addDetail.IsProductPromotionPo
	}

	orderDetail := &model.OrderDetail{
		CustId:               custId,
		RoNo:                 roNo,
		ProId:                int(addDetail.ProId),
		ItemType:             1,
		QtyPo:                float64(totalQty),
		QtyPo1:               &qtyPo1,
		QtyPo2:               &qtyPo2,
		QtyPo3:               &qtyPo3,
		OriginalQtyPo1:       &requestedQtyPo1,
		OriginalQtyPo2:       &requestedQtyPo2,
		OriginalQtyPo3:       &requestedQtyPo3,
		SellPricePo1:         &addDetail.SellPricePo1,
		SellPricePo2:         &addDetail.SellPricePo2,
		SellPricePo3:         &addDetail.SellPricePo3,
		Qty:                  float64(totalQty),
		Qty1:                 &qtyPo1,
		Qty2:                 &qtyPo2,
		Qty3:                 &qtyPo3,
		SellPrice1:           &addDetail.SellPricePo1,
		SellPrice2:           &addDetail.SellPricePo2,
		SellPrice3:           &addDetail.SellPricePo3,
		QtyFinal:             float64(totalQty),
		Qty1Final:            &qtyPo1,
		Qty2Final:            &qtyPo2,
		Qty3Final:            &qtyPo3,
		SellPriceFinal1:      &addDetail.SellPricePo1,
		SellPriceFinal2:      &addDetail.SellPricePo2,
		SellPriceFinal3:      &addDetail.SellPricePo3,
		SellPriceSystem1:     &addDetail.SellPriceSystem1,
		SellPriceSystem2:     &addDetail.SellPriceSystem2,
		SellPriceSystem3:     &addDetail.SellPriceSystem3,
		Vat:                  float64Ptr(product.Vat),
		ConvUnit2:            intPtr(int(product.ConvUnit2)),
		ConvUnit3:            intPtr(int(product.ConvUnit3)),
		UnitId1:              stringPtr(unitId1),
		UnitId2:              stringPtr(unitId2),
		UnitId3:              stringPtr(unitId3),
		Qty1Stok:             &qty1Stok,
		Qty2Stok:             &qty2Stok,
		Qty3Stok:             &qty3Stok,
		IsProductPromotionPo: &isProductPromotionPo,
	}

	if err := service.OrderRepository.StoreDetail(ctx, orderDetail); err != nil {
		return nil, err
	}

	var trCode string
	if len(roNo) > 2 {
		trCode = roNo[0:2]
	}

	stockUpdate := &entity.SalesOrderStockUpdate{
		CustID:         custId,
		WhID:           whId,
		ProID:          int64(addDetail.ProId),
		StockDate:      roDate,
		TrCode:         trCode,
		TrNo:           roNo,
		QtyOrder:       float64(totalQty),
		QtyOrderBefore: nil,
		UnitPrice:      addDetail.SellPricePo1,
		RefDetId:       int64(*orderDetail.OrderDetailID),
	}

	return stockUpdate, nil
}

// createOrderDetailFromSalesOrder creates a new order detail from Sales Order tab
// Cascades to SO and Final fields
func (service *orderServiceImpl) createOrderDetailFromSalesOrder(ctx context.Context, roNo string, custId string, whId int64, roDate time.Time, addDetail entity.AddSalesOrderDetail) (*entity.SalesOrderStockUpdate, error) {
	// Get product info for conversion units
	product, err := service.OrderRepository.FindProductByID(int(addDetail.ProId))
	if err != nil {
		return nil, fmt.Errorf("product with id %d not found: %w", addDetail.ProId, err)
	}

	qty1 := addDetail.Qty1
	qty2 := addDetail.Qty2
	qty3 := addDetail.Qty3
	totalQty, _, err := calculateNormalizedQty(int(qty1), int(qty2), int(qty3), int(product.ConvUnit2), int(product.ConvUnit3))
	if err != nil {
		return nil, err
	}

	unitId1 := addDetail.UnitId1
	unitId2 := addDetail.UnitId2
	unitId3 := addDetail.UnitId3
	qty1Stok := addDetail.Qty1Stock
	qty2Stok := addDetail.Qty2Stock
	qty3Stok := addDetail.Qty3Stock

	defaultFalse := false
	finalPromotionFlag := false
	orderDetail := &model.OrderDetail{
		CustId:   custId,
		RoNo:     roNo,
		ProId:    int(addDetail.ProId),
		ItemType: 1,
		// PO fields 0
		QtyPo:        0,
		QtyPo1:       float64Ptr(0),
		QtyPo2:       float64Ptr(0),
		QtyPo3:       float64Ptr(0),
		SellPricePo1: float64Ptr(0),
		SellPricePo2: float64Ptr(0),
		SellPricePo3: float64Ptr(0),
		// Sales fields
		Qty:        float64(totalQty),
		Qty1:       &qty1,
		Qty2:       &qty2,
		Qty3:       &qty3,
		SellPrice1: &addDetail.SellPrice1,
		SellPrice2: &addDetail.SellPrice2,
		SellPrice3: &addDetail.SellPrice3,
		// Final fields
		QtyFinal:                float64(totalQty),
		Qty1Final:               &qty1,
		Qty2Final:               &qty2,
		Qty3Final:               &qty3,
		SellPriceFinal1:         &addDetail.SellPrice1,
		SellPriceFinal2:         &addDetail.SellPrice2,
		SellPriceFinal3:         &addDetail.SellPrice3,
		SellPriceSystem1:        &addDetail.SellPriceSystem1,
		SellPriceSystem2:        &addDetail.SellPriceSystem2,
		SellPriceSystem3:        &addDetail.SellPriceSystem3,
		Vat:                     float64Ptr(product.Vat),
		IsProductPromotionSo:    addDetail.IsProductPromotionSo,
		IsProductPromotionFinal: &finalPromotionFlag,
		IsProductPromotionPo:    &defaultFalse,
		ConvUnit2:               intPtr(int(product.ConvUnit2)),
		ConvUnit3:               intPtr(int(product.ConvUnit3)),
		UnitId1:                 &unitId1,
		UnitId2:                 &unitId2,
		UnitId3:                 &unitId3,
		Qty1Stok:                &qty1Stok,
		Qty2Stok:                &qty2Stok,
		Qty3Stok:                &qty3Stok,
	}

	if err := service.OrderRepository.StoreDetail(ctx, orderDetail); err != nil {
		return nil, err
	}

	var trCode string
	if len(roNo) > 2 {
		trCode = roNo[0:2]
	}

	stockUpdate := &entity.SalesOrderStockUpdate{
		CustID:         custId,
		WhID:           whId,
		ProID:          int64(addDetail.ProId),
		StockDate:      roDate,
		TrCode:         trCode,
		TrNo:           roNo,
		QtyOrder:       float64(totalQty),
		QtyOrderBefore: nil,
		UnitPrice:      addDetail.SellPrice1,
		RefDetId:       int64(*orderDetail.OrderDetailID),
	}

	return stockUpdate, nil
}

// createOrderDetailFromFinalOrder creates a new order detail from Final Order tab
// Only updates Final fields
func (service *orderServiceImpl) createOrderDetailFromFinalOrder(ctx context.Context, roNo string, custId string, whId int64, roDate time.Time, addDetail entity.AddFinalOrderDetail) (*entity.SalesOrderStockUpdate, error) {
	// Get product info for conversion units
	product, err := service.OrderRepository.FindProductByID(int(addDetail.ProId))
	if err != nil {
		return nil, fmt.Errorf("product with id %d not found: %w", addDetail.ProId, err)
	}

	// Calculate total quantity
	qty1Final := addDetail.Qty1Final
	qty2Final := addDetail.Qty2Final
	qty3Final := addDetail.Qty3Final
	totalQty, _, err := calculateNormalizedQty(int(qty1Final), int(qty2Final), int(qty3Final), int(product.ConvUnit2), int(product.ConvUnit3))
	if err != nil {
		return nil, err
	}

	// Build order detail with Final fields only
	orderDetail := &model.OrderDetail{
		CustId:   custId,
		RoNo:     roNo,
		ProId:    int(addDetail.ProId),
		ItemType: 1, // Normal item
		// Purchase Order fields (Initialize to 0)
		QtyPo:        0,
		QtyPo1:       float64Ptr(0),
		QtyPo2:       float64Ptr(0),
		QtyPo3:       float64Ptr(0),
		SellPricePo1: float64Ptr(0),
		SellPricePo2: float64Ptr(0),
		SellPricePo3: float64Ptr(0),
		// Sales Order fields (Initialize to 0)
		Qty:        0,
		Qty1:       float64Ptr(0),
		Qty2:       float64Ptr(0),
		Qty3:       float64Ptr(0),
		SellPrice1: float64Ptr(0),
		SellPrice2: float64Ptr(0),
		SellPrice3: float64Ptr(0),
		// Final Order fields only
		QtyFinal:        float64(totalQty),
		Qty1Final:       &qty1Final,
		Qty2Final:       &qty2Final,
		Qty3Final:       &qty3Final,
		SellPriceFinal1: &addDetail.SellPriceFinal1,
		SellPriceFinal2: &addDetail.SellPriceFinal2,
		SellPriceFinal3: &addDetail.SellPriceFinal3,
		// System prices
		SellPriceSystem1:        &addDetail.SellPriceSystem1,
		SellPriceSystem2:        &addDetail.SellPriceSystem2,
		SellPriceSystem3:        &addDetail.SellPriceSystem3,
		Vat:                     float64Ptr(product.Vat),
		IsProductPromotionFinal: addDetail.IsProductPromotionFinal,
		// Conversion units from product
		ConvUnit2: intPtr(int(product.ConvUnit2)),
		ConvUnit3: intPtr(int(product.ConvUnit3)),
		UnitId1:   stringPtr(product.UnitId1),
		UnitId2:   stringPtr(product.UnitId2),
		UnitId3:   stringPtr(product.UnitId3),
	}

	if err := service.OrderRepository.StoreDetail(ctx, orderDetail); err != nil {
		return nil, err
	}

	trCode := ""
	if len(roNo) >= 2 {
		trCode = roNo[0:2]
	}

	return &entity.SalesOrderStockUpdate{
		CustID:         custId,
		WhID:           whId,
		ProID:          int64(addDetail.ProId),
		StockDate:      roDate,
		TrCode:         trCode,
		TrNo:           roNo,
		QtyOrderBefore: nil,
		QtyOrder:       float64(totalQty),
		UnitPrice:      addDetail.SellPriceFinal1,
		RefDetId:       int64(*orderDetail.OrderDetailID),
	}, nil
}

// float64Ptr is a helper to get pointer to float64
func float64Ptr(v float64) *float64 {
	return &v
}

func applyOriginalTakingOrderQty(orderType *string, detail *model.OrderDetail) {
	if orderType == nil || *orderType != "O" || detail == nil {
		return
	}

	detail.OriginalQtyPo1 = detail.QtyPo1
	detail.OriginalQtyPo2 = detail.QtyPo2
	detail.OriginalQtyPo3 = detail.QtyPo3
}

// boolPtr is a helper to get pointer to bool
func boolPtr(v bool) *bool {
	return &v
}

// intPtr is a helper to get pointer to int
func intPtr(v int) *int {
	return &v
}

// stringPtr is a helper to get pointer to string
func stringPtr(v string) *string {
	return &v
}

func stringFromPtr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// isEmptyStringPtr returns true if the string pointer is nil or points to an empty string
func isEmptyStringPtr(s *string) bool {
	return s == nil || *s == ""
}

// orderImportHeaders mirrors the QA-provided template (SX-2470 follow-up).
// The template is intentionally flat: only fields the FE has data for are
// present. Fields like warehouse, pay_type, and unit conversion are derived
// from lookup tables or defaulted so that downstream Store() still gets a
// fully populated CreateOrderBody.
var orderImportHeaders = []string{
	"DocumentNo", "DocumentDate", "OutletCode", "OutletName",
	"SalesmanCode", "Salesman Name", "ProCode", "ProName",
	"Price", "Unit", "QTY", "GrossSales", "Promo", "Discount",
	"PPN", "NetSalesIncPPN",
}

func (service *orderServiceImpl) ExportTemplate(format string) (*bytes.Buffer, string, string, error) {
	f := excelize.NewFile()
	// excelize.NewFile() always creates a default "Sheet1". To make the
	// downloaded template a true single-sheet workbook (the importable header
	// row only), we rename that default sheet to "Order Template" and reuse it
	// instead of creating a new one; excelize forbids deleting the last sheet.
	templateSheet := f.GetSheetName(0)
	if templateSheet == "" {
		templateSheet = "Order Template"
		if idx, err := f.NewSheet(templateSheet); err == nil {
			f.SetActiveSheet(idx)
		}
	} else {
		if err := f.SetSheetName(templateSheet, "Order Template"); err == nil {
			templateSheet = "Order Template"
		}
		idx, _ := f.GetSheetIndex(templateSheet)
		f.SetActiveSheet(idx)
	}
	for i, header := range orderImportHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(templateSheet, cell, header)
	}
	// Sample row (B2) with date format so Excel renders it as a date.
	styleID, _ := f.NewStyle(&excelize.Style{NumFmt: 14})
	f.SetCellValue(templateSheet, "B2", "09/07/2026")
	f.SetCellStyle(templateSheet, "B2", "B2", styleID)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}
	// Only xlsx is supported; excelize does not generate legacy XLS.
	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "order_import_template.xlsx", nil
}

var jakartaLoc = time.FixedZone("Asia/Jakarta", 7*60*60)

func validateImportDate(documentDate, today time.Time) string {
	today = today.In(jakartaLoc).Truncate(24 * time.Hour)
	date := documentDate.In(jakartaLoc).Truncate(24 * time.Hour)
	if date.After(today) {
		return "Transaction Date cannot be later than the current date."
	}
	if date.Before(today.AddDate(0, 0, -7)) {
		return "Transaction Date cannot be more than 7 days before the current date."
	}
	return ""
}

func (service *orderServiceImpl) importSecondarySales(ctx context.Context, custId string, parentCustId string, userId int64, parsed []entity.CreateOrderBody) error {
	scopeSet := make(map[string]time.Time)
	for _, order := range parsed {
		if order.RoDate == nil {
			continue
		}
		date, err := time.ParseInLocation("2006-01-02", *order.RoDate, jakartaLoc)
		if err != nil {
			return err
		}
		scopeSet[date.Format("2006-01-02")] = date
	}
	scope := make([]time.Time, 0, len(scopeSet))
	for _, date := range scopeSet {
		scope = append(scope, date)
	}
	return service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := service.OrderRepository.LockOrderByScope(txCtx, custId, scope); err != nil {
			return err
		}
		if _, err := service.OrderRepository.DeleteOrderDetailByScope(txCtx, custId, scope); err != nil {
			return err
		}
		if _, err := service.OrderRepository.DeleteOrderByScope(txCtx, custId, scope); err != nil {
			return err
		}
		for _, order := range parsed {
			orderModel := model.Order{
				CustID:            order.CustId,
				RoNo:              order.RoNo,
				IsSalesMapping:    boolPtr(true),
				SalesmanId:        &order.SalesmanId,
				WhId:              order.WhId,
				OutletID:          &order.OutletID,
				PayType:           order.PayType,
				SubTotal:          order.SubTotal,
				SubTotalFinal:     order.SubTotalFinal,
				Disc:              order.Disc,
				DiscValue:         order.DiscValue,
				DiscValueFinal:    order.DiscValueFinal,
				PromoValue:        order.PromoValue,
				PromoValueFinal:   order.PromoValueFinal,
				PromoBgValue:      order.PromoBgValue,
				PromoBgValueFinal: order.PromoBgValueFinal,
				CashDiscValue:     order.CashDiscValue,
				TotDisc1:          order.TotDisc1,
				TotDisc2:          order.TotDisc2,
				Vat:               order.Vat,
				VatValue:          order.VatValue,
				VatValueFinal:     order.VatValueFinal,
				Total:             order.Total,
				TotalFinal:        order.TotalFinal,
				DataStatus:        int64Ptr(int64(importDataStatus)),
				DataSource:        int64Ptr(int64(importDataSource)),
				CreatedBy:         &userId,
				InvoiceNo:         order.InvoiceNo,
				Notes:             order.Notes,
				Address1:          nil,
			}
			if order.RoDate != nil {
				roDate, err := time.Parse("2006-01-02", *order.RoDate)
				if err != nil {
					return err
				}
				orderModel.RoDate = &roDate
			}
			if order.DeliveryDate != nil {
				dd, err := time.Parse("2006-01-02", *order.DeliveryDate)
				if err != nil {
					return err
				}
				orderModel.DeliveryDate = &dd
			}
			if order.DueDate != nil {
				dd, err := time.Parse("2006-01-02", *order.DueDate)
				if err != nil {
					return err
				}
				orderModel.DueDate = &dd
			}
			if order.InvoiceDate != nil {
				id, err := time.Parse("2006-01-02", *order.InvoiceDate)
				if err != nil {
					return err
				}
				orderModel.InvoiceDate = &id
			}
			if err := service.OrderRepository.Store(txCtx, &orderModel); err != nil {
				return err
			}
			for _, det := range order.Details.Normal {
				detail := model.OrderDetail{
					CustId:                  order.CustId,
					RoNo:                    order.RoNo,
					SeqNo:                   det.SeqNo,
					ProId:                   det.ProId,
					ItemType:                1,
					Qty:                     getValueOrDefault(det.Qty, 0),
					QtyFinal:                getValueOrDefault(det.Qty, 0),
					Qty1:                    det.Qty1,
					Qty2:                    det.Qty2,
					Qty3:                    det.Qty3,
					Qty1Final:               det.Qty1,
					Qty2Final:               det.Qty2,
					Qty3Final:               det.Qty3,
					PurchPrice1:             det.PurchPrice1,
					PurchPrice2:             det.PurchPrice2,
					PurchPrice3:             det.PurchPrice3,
					SellPrice1:              det.SellPrice1,
					SellPrice2:              det.SellPrice2,
					SellPrice3:              det.SellPrice3,
					SellPriceSystem1:        det.SellPriceSystem1,
					SellPriceSystem2:        det.SellPriceSystem2,
					SellPriceSystem3:        det.SellPriceSystem3,
					Amount:                  det.Amount,
					AmountFinal:             det.AmountFinal,
					DiscValue:               det.DiscValue,
					DiscValueFinal:          det.DiscValueFinal,
					PromoValue:              det.PromoValue,
					PromoValueFinal:         det.PromoValueFinal,
					PromoSo1:                det.PromoSo1,
					PromoFinal1:             det.PromoFinal1,
					Vat:                     det.Vat,
					VatValue:                det.VatValue,
					VatValueFinal:           det.VatValueFinal,
					UnitId1:                 det.UnitId1,
					UnitId2:                 det.UnitId2,
					UnitId3:                 det.UnitId3,
					ConvUnit2:               det.ConvUnit2,
					ConvUnit3:               det.ConvUnit3,
					IsProductPromotionSo:    det.IsProductPromotionSo,
					IsProductPromotionFinal: det.IsProductPromotionFinal,
				}
				if err := service.OrderRepository.StoreDetail(txCtx, &detail); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (service *orderServiceImpl) ImportOrders(custId string, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportResult, []entity.OrderImportError, error) {
	parsed, errs, summary, err := service.parseImportOrders(custId, parentCustId, userId, file, filename)
	if err != nil {
		return entity.OrderImportResult{}, nil, err
	}
	if len(errs) > 0 {
		failedReasons := make([]string, 0, len(errs))
		for _, e := range errs {
			failedReasons = append(failedReasons, formatImportFailedReason(e))
		}
		return entity.OrderImportResult{
			StartDate:       summary.StartDate,
			EndDate:         summary.EndDate,
			NumberOfInvoice: summary.NumberOfInvoice,
			NumberOfOutlet:  summary.NumberOfOutlet,
			Amount:          summary.Amount,
			CreatedRoNos:    []string{},
		}, nil, &entity.ImportFailedError{FailedReasons: failedReasons}
	}
	result := entity.OrderImportResult{
		StartDate:       summary.StartDate,
		EndDate:         summary.EndDate,
		NumberOfInvoice: summary.NumberOfInvoice,
		NumberOfOutlet:  summary.NumberOfOutlet,
		Amount:          summary.Amount,
		CreatedRoNos:    []string{},
	}
	importValidation := entity.ValidateResponse{
		Validate1Success: true,
		Validate1:        "Sufficient Stock",
		Validate2Success: true,
		Validate2:        "Within Limit",
		Validate3Success: true,
		Validate3:        "Allowed",
		Validate4Success: true,
		Validate4:        "Allowed",
	}
	_ = importValidation
	if err := service.importSecondarySales(context.Background(), custId, parentCustId, userId, parsed); err != nil {
		return entity.OrderImportResult{}, nil, err
	}
	for _, req := range parsed {
		result.CreatedRoNos = append(result.CreatedRoNos, req.RoNo)
	}
	return result, nil, nil
}

func (service *orderServiceImpl) ValidateImport(custId string, parentCustId string, userId int64, file io.Reader, filename string) (entity.OrderImportSummary, error) {
	parsed, errs, summary, err := service.parseImportOrders(custId, parentCustId, userId, file, filename)
	if err != nil {
		return entity.OrderImportSummary{}, err
	}
	_ = parsed
	if len(errs) > 0 {
		failedReasons := make([]string, 0, len(errs))
		for _, e := range errs {
			failedReasons = append(failedReasons, formatImportFailedReason(e))
		}
		summary.FailedReasons = failedReasons
	}
	return summary, nil
}

func (service *orderServiceImpl) parseImportOrders(custId string, parentCustId string, userId int64, file io.Reader, filename string) ([]entity.CreateOrderBody, []entity.OrderImportError, entity.OrderImportSummary, error) {
	_ = filename
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, nil, entity.OrderImportSummary{}, err
	}
	defer f.Close()
	var matchedSheet string
	for _, sheetName := range f.GetSheetList() {
		sheetRows, _ := f.GetRows(sheetName)
		if len(sheetRows) == 0 {
			continue
		}
		if len(sheetRows[0]) < len(orderImportHeaders) {
			continue
		}
		match := true
		for i, h := range orderImportHeaders {
			if strings.TrimSpace(sheetRows[0][i]) != h {
				match = false
				break
			}
		}
		if match {
			if matchedSheet != "" {
				return nil, []entity.OrderImportError{{Row: 1, Field: "file", Message: "multiple template sheets found"}}, entity.OrderImportSummary{}, nil
			}
			matchedSheet = sheetName
		}
	}
	if matchedSheet == "" {
		return nil, []entity.OrderImportError{{Row: 1, Field: "file", Message: "template sheet is empty"}}, entity.OrderImportSummary{}, nil
	}
	rows, err := f.GetRows(matchedSheet)
	if err != nil || len(rows) < 2 {
		return nil, []entity.OrderImportError{{Row: 1, Field: "file", Message: "template sheet is empty"}}, entity.OrderImportSummary{}, nil
	}

	type parsedImportLine struct {
		row                int
		documentNo         string
		documentDate       string
		outlet             model.OutletRead
		salesman           model.SalesmanDetail
		product            model.ProductRead
		parentProId        int
		groupingProductKey string
		proCode            string
		outletName         string
		salesmanName       string
		proName            string
		unit               string
		qty                float64
		price              float64
		grossSales         float64
		promo              float64
		discount           float64
		ppn                float64
		netSales           float64
		detail             entity.CreateOrderDetBody
	}

	parseRequired := func(field string, line int, raw string) (string, bool, []entity.OrderImportError) {
		if strings.TrimSpace(raw) == "" {
			return "", false, []entity.OrderImportError{{Row: line, Field: field, Message: "required"}}
		}
		return strings.TrimSpace(raw), true, nil
	}

	parseRequiredNumber := func(field string, line int, raw string) (float64, bool, []entity.OrderImportError) {
		v, ok, errs := parseRequired(field, line, raw)
		if !ok {
			return 0, false, errs
		}
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false, []entity.OrderImportError{{Row: line, Field: field, Message: "invalid number"}}
		}
		return n, true, nil
	}

	parseOptional := func(field string, line int, raw string) (string, bool, []entity.OrderImportError) {
		return strings.TrimSpace(raw), true, nil
	}

	parseOptionalNumber := func(field string, line int, raw string) (float64, bool, []entity.OrderImportError) {
		v := strings.TrimSpace(raw)
		if v == "" {
			return 0, true, nil
		}
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false, []entity.OrderImportError{{Row: line, Field: field, Message: "invalid number"}}
		}
		return n, true, nil
	}

	matchUnitSlot := func(product model.ProductRead, unit string) (qty1, qty2, qty3 float64, sell1, sell2, sell3 *float64, conv2, conv3 *int, err error) {
		unit = strings.TrimSpace(strings.ToUpper(unit))
		u1 := strings.TrimSpace(strings.ToUpper(product.UnitId1))
		u2 := strings.TrimSpace(strings.ToUpper(product.UnitId2))
		u3 := strings.TrimSpace(strings.ToUpper(product.UnitId3))
		u4 := strings.TrimSpace(strings.ToUpper(product.UnitId4))
		u5 := strings.TrimSpace(strings.ToUpper(product.UnitId5))
		switch unit {
		case u1:
			return 0, 0, 0, float64Ptr(product.SellPrice1), float64Ptr(product.SellPrice2), float64Ptr(product.SellPrice3), intPtrFromFloat32(product.ConvUnit2), intPtrFromFloat32(product.ConvUnit3), nil
		case u2:
			return 0, 0, 0, float64Ptr(product.SellPrice1), float64Ptr(product.SellPrice2), float64Ptr(product.SellPrice3), intPtrFromFloat32(product.ConvUnit2), intPtrFromFloat32(product.ConvUnit3), nil
		case u3:
			return 0, 0, 0, float64Ptr(product.SellPrice1), float64Ptr(product.SellPrice2), float64Ptr(product.SellPrice3), intPtrFromFloat32(product.ConvUnit2), intPtrFromFloat32(product.ConvUnit3), nil
		case u4:
			return 0, 0, 0, float64Ptr(product.SellPrice1), float64Ptr(product.SellPrice2), float64Ptr(product.SellPrice3), intPtrFromFloat32(product.ConvUnit2), intPtrFromFloat32(product.ConvUnit3), nil
		case u5:
			return 0, 0, 0, float64Ptr(product.SellPrice1), float64Ptr(product.SellPrice2), float64Ptr(product.SellPrice3), intPtrFromFloat32(product.ConvUnit2), intPtrFromFloat32(product.ConvUnit3), nil
		default:
			return 0, 0, 0, nil, nil, nil, nil, nil, fmt.Errorf("unit not mapped")
		}
	}

	// mapQtyAndPriceToSlot materializes a CreateOrderDetBody for an import row.
	// Per Import_Sales_Order_BE.docx, VatValue/VatValueFinal must come from the
	// uploaded row's PPN; the percentage Vat field is left at 0 because the
	// store path does not derive VAT percentage from the imported row.
	mapQtyAndPriceToSlot := func(parentProduct model.ProductRead, distributorProduct model.ProductRead, unit string, qty float64, price float64, ppnValue float64) (entity.CreateOrderDetBody, error) {
		unit = strings.TrimSpace(strings.ToUpper(unit))
		detail := entity.CreateOrderDetBody{
			ProId:            parentProduct.ProId,
			Qty1:             float64Ptr(0),
			Qty2:             float64Ptr(0),
			Qty3:             float64Ptr(0),
			ConvUnit2:        intPtrFromFloat32(parentProduct.ConvUnit2),
			ConvUnit3:        intPtrFromFloat32(parentProduct.ConvUnit3),
			PurchPrice1:      float64Ptr(parentProduct.PurchPrice1),
			PurchPrice2:      float64Ptr(parentProduct.PurchPrice2),
			PurchPrice3:      float64Ptr(parentProduct.PurchPrice3),
			SellPrice1:       float64Ptr(0),
			SellPrice2:       float64Ptr(0),
			SellPrice3:       float64Ptr(0),
			SellPriceSystem1: float64Ptr(parentProduct.SellPrice1),
			SellPriceSystem2: float64Ptr(parentProduct.SellPrice3),
			SellPriceSystem3: float64Ptr(parentProduct.SellPrice3),
			Vat:              float64Ptr(0),
			PromoSo1:         float64Ptr(0),
			PromoFinal1:      float64Ptr(0),
			PromoValue:       float64Ptr(0),
			PromoValueFinal:  float64Ptr(0),
			QtyPo:            nil,
			DiscValue:        float64Ptr(0),
			DiscValueFinal:   float64Ptr(0),
			Amount:           float64Ptr(0),
			AmountFinal:      float64Ptr(0),
			VatValue:         float64Ptr(ppnValue),
			VatValueFinal:    float64Ptr(ppnValue),
			UnitId1:          strPtrIfNotEmpty(parentProduct.UnitId1),
			UnitId2:          strPtrIfNotEmpty(parentProduct.UnitId2),
			UnitId3:          strPtrIfNotEmpty(parentProduct.UnitId3),
		}
		switch unit {
		case strings.TrimSpace(strings.ToUpper(distributorProduct.UnitId1)):
			detail.Qty1 = float64Ptr(qty)
			detail.SellPrice1 = float64Ptr(price)
		case strings.TrimSpace(strings.ToUpper(distributorProduct.UnitId2)):
			detail.Qty2 = float64Ptr(qty)
			detail.SellPrice2 = float64Ptr(price)
		case strings.TrimSpace(strings.ToUpper(distributorProduct.UnitId3)):
			detail.Qty3 = float64Ptr(qty)
			detail.SellPrice3 = float64Ptr(price)
		default:
			return detail, fmt.Errorf("unit not mapped")
		}
		return detail, nil
	}

	var errs []entity.OrderImportError
	linesByDocument := make(map[string][]parsedImportLine)
	seenDocumentUnit := make(map[string]map[string]struct{})

	for idx, row := range rows[1:] {
		line := idx + 2
		get := func(i int) string {
			if i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}

		// Skip stray non-data rows (e.g. trailing formula totals or
		// the validation rule table that may live in the same sheet).
		if get(0) == "" {
			continue
		}
		documentNo, ok, rowErrs := parseRequired("DocumentNo", line, get(0))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		documentDate, ok, rowErrs := parseRequired("DocumentDate", line, get(1))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		parsedDocumentDate, err := parseImportDate(documentDate)
		if err != nil {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "DocumentDate", Message: "invalid date format"})
			continue
		}
		if message := validateImportDate(parsedDocumentDate, time.Now().In(jakartaLoc)); message != "" {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "DocumentDate", Message: message})
			continue
		}
		documentDate = parsedDocumentDate.Format("2006-01-02")
		outletCode, ok, rowErrs := parseRequired("OutletCode", line, get(2))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		outletName, ok, rowErrs := parseRequired("OutletName", line, get(3))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		salesmanCode, ok, rowErrs := parseRequired("SalesmanCode", line, get(4))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		salesmanName, ok, rowErrs := parseRequired("Salesman Name", line, get(5))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		proCode, _, _ := parseOptional("ProCode", line, get(6))
		proName, ok, rowErrs := parseRequired("ProName", line, get(7))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		price, ok, rowErrs := parseRequiredNumber("Price", line, get(8))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		unit, ok, rowErrs := parseRequired("Unit", line, get(9))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		qty, ok, rowErrs := parseRequiredNumber("QTY", line, get(10))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		if qty <= 0 {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "QTY", Message: "quantity must be > 0"})
			continue
		}
		grossSales, ok, rowErrs := parseRequiredNumber("GrossSales", line, get(11))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		promo, _, rowErrs := parseOptionalNumber("Promo", line, get(12))
		if rowErrs != nil {
			errs = append(errs, rowErrs...)
			continue
		}
		discount, _, rowErrs := parseOptionalNumber("Discount", line, get(13))
		if rowErrs != nil {
			errs = append(errs, rowErrs...)
			continue
		}
		ppn, _, rowErrs := parseOptionalNumber("PPN", line, get(14))
		if rowErrs != nil {
			errs = append(errs, rowErrs...)
			continue
		}
		netSales, ok, rowErrs := parseRequiredNumber("NetSalesIncPPN", line, get(15))
		if !ok {
			errs = append(errs, rowErrs...)
			continue
		}
		if grossSales != price*qty {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "GrossSales", Message: "must equal Price * QTY"})
			continue
		}
		if netSales != grossSales-promo-discount+ppn {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "NetSalesIncPPN", Message: "must equal GrossSales - Promo - Discount + PPN"})
			continue
		}

		outlet, err := service.OrderRepository.FindOutletByCode(outletCode, custId)
		if err != nil {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "OutletCode", Message: "not found"})
			continue
		}
		if strings.TrimSpace(outletName) == "" || !strings.EqualFold(strings.TrimSpace(outlet.OutletName), strings.TrimSpace(outletName)) {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "OutletName", Message: "does not match master outlet"})
			continue
		}
		salesman, err := service.OrderRepository.FindSalesmanByCode(salesmanCode, custId)
		if err != nil {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "SalesmanCode", Message: "not found"})
			continue
		}
		if strings.TrimSpace(salesmanName) == "" || !strings.EqualFold(strings.TrimSpace(salesman.SalesmanName), strings.TrimSpace(salesmanName)) {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "Salesman Name", Message: "does not match master salesman"})
			continue
		}
		product := model.ProductRead{}
		parentProduct := model.ProductRead{}
		if strings.TrimSpace(proCode) != "" {
			product, err = service.OrderRepository.FindProductByCode(proCode, custId)
			if err != nil {
				errs = append(errs, entity.OrderImportError{Row: line, Field: "ProCode", Message: "not found"})
				continue
			}
			if strings.TrimSpace(proName) == "" || !strings.EqualFold(strings.TrimSpace(product.ProName), strings.TrimSpace(proName)) {
				errs = append(errs, entity.OrderImportError{Row: line, Field: "ProName", Message: "does not match product mapping"})
				continue
			}
		} else {
			if strings.TrimSpace(proName) == "" {
				errs = append(errs, entity.OrderImportError{Row: line, Field: "ProName", Message: "is required when ProCode is blank"})
				continue
			}
			product, err = service.OrderRepository.FindProductByName(proName, custId)
			if err != nil {
				errs = append(errs, entity.OrderImportError{Row: line, Field: "ProName", Message: "not found"})
				continue
			}
		}
		// Per docx: pro_id, purch_price*, sell_price_system*, conv_unit* and
		// the unit_id* values stored in sls.order_detail all come from the
		// parent product. The distributor product only determines which
		// unit slot (1/2/3) the row maps to.
		parentProduct = product
		if product.ParentProId != nil && *product.ParentProId > 0 {
			parent, err := service.OrderRepository.FindProductByID(*product.ParentProId)
			if err != nil {
				errs = append(errs, entity.OrderImportError{Row: line, Field: "ProCode", Message: "parent product not found"})
				continue
			}
			parentProduct = parent
		}
		if _, _, _, _, _, _, _, _, err := matchUnitSlot(product, unit); err != nil {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "Unit", Message: "is not mapped for this product"})
			continue
		}
		detail, err := mapQtyAndPriceToSlot(parentProduct, product, unit, qty, price, ppn)
		if err != nil {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "Unit", Message: "is not mapped for this product"})
			continue
		}
		detail.PromoSo1 = float64Ptr(promo)
		detail.PromoFinal1 = float64Ptr(promo)
		detail.PromoValue = float64Ptr(promo)
		detail.PromoValueFinal = float64Ptr(promo)
		detail.DiscValue = float64Ptr(discount)
		detail.DiscValueFinal = float64Ptr(discount)
		detail.VatValue = float64Ptr(ppn)
		detail.VatValueFinal = float64Ptr(ppn)
		detail.Amount = float64Ptr(netSales)
		detail.AmountFinal = float64Ptr(netSales)

		groupingProductKey := strings.ToUpper(strings.TrimSpace(proCode))
		if groupingProductKey == "" {
			groupingProductKey = strings.ToUpper(strings.TrimSpace(product.ProCode))
		}
		if groupingProductKey == "" {
			groupingProductKey = fmt.Sprintf("PARENT:%d", detail.ProId)
		}
		if _, ok := seenDocumentUnit[documentNo]; !ok {
			seenDocumentUnit[documentNo] = map[string]struct{}{}
		}
		dupKey := fmt.Sprintf("%s|%s", groupingProductKey, strings.ToUpper(strings.TrimSpace(unit)))
		if _, exists := seenDocumentUnit[documentNo][dupKey]; exists {
			errs = append(errs, entity.OrderImportError{Row: line, Field: "Unit", Message: "duplicate Product and Unit found in Document No"})
			continue
		}
		seenDocumentUnit[documentNo][dupKey] = struct{}{}

		linesByDocument[documentNo] = append(linesByDocument[documentNo], parsedImportLine{
			row:                line,
			documentNo:         documentNo,
			documentDate:       documentDate,
			outlet:             outlet,
			salesman:           salesman,
			product:            product,
			parentProId:        detail.ProId,
			groupingProductKey: groupingProductKey,
			proCode:            proCode,
			outletName:         outletName,
			salesmanName:       salesmanName,
			proName:            proName,
			unit:               unit,
			qty:                qty,
			price:              price,
			grossSales:         grossSales,
			promo:              promo,
			discount:           discount,
			ppn:                ppn,
			netSales:           netSales,
			detail:             detail,
		})
	}

	if len(errs) > 0 {
		return nil, errs, entity.OrderImportSummary{}, nil
	}

	parsed := make([]entity.CreateOrderBody, 0, len(linesByDocument))
	for documentNo, lines := range linesByDocument {
		if len(lines) == 0 {
			continue
		}
		first := lines[0]
		subTotal := 0.0
		promoTotal := 0.0
		discTotal := 0.0
		vatTotal := 0.0
		netTotal := 0.0
		groupIndex := map[string]int{}
		groups := make([]parsedImportLine, 0, len(lines))
		for _, line := range lines {
			subTotal += line.grossSales
			promoTotal += line.promo
			discTotal += line.discount
			vatTotal += line.ppn
			netTotal += line.netSales
			if idx, exists := groupIndex[line.groupingProductKey]; exists {
				agg := &groups[idx]
				if agg.detail.Qty1 != nil && *agg.detail.Qty1 == 0 && line.detail.Qty1 != nil {
					agg.detail.Qty1 = line.detail.Qty1
					agg.detail.SellPrice1 = line.detail.SellPrice1
				}
				if agg.detail.Qty2 != nil && *agg.detail.Qty2 == 0 && line.detail.Qty2 != nil {
					agg.detail.Qty2 = line.detail.Qty2
					agg.detail.SellPrice2 = line.detail.SellPrice2
				}
				if agg.detail.Qty3 != nil && *agg.detail.Qty3 == 0 && line.detail.Qty3 != nil {
					agg.detail.Qty3 = line.detail.Qty3
					agg.detail.SellPrice3 = line.detail.SellPrice3
				}
				agg.grossSales += line.grossSales
				agg.promo += line.promo
				agg.discount += line.discount
				agg.ppn += line.ppn
				agg.netSales += line.netSales
				continue
			}
			groupIndex[line.groupingProductKey] = len(groups)
			groups = append(groups, line)
		}
		details := make([]entity.CreateOrderDetBody, 0, len(groups))
		for idx := range groups {
			g := &groups[idx]
			g.detail.SeqNo = idx + 1
			g.detail.PromoSo1 = float64Ptr(g.promo)
			g.detail.PromoFinal1 = float64Ptr(g.promo)
			g.detail.PromoValue = float64Ptr(g.promo)
			g.detail.PromoValueFinal = float64Ptr(g.promo)
			g.detail.DiscValue = float64Ptr(g.discount)
			g.detail.DiscValueFinal = float64Ptr(g.discount)
			g.detail.VatValue = float64Ptr(g.ppn)
			g.detail.VatValueFinal = float64Ptr(g.ppn)
			g.detail.Amount = float64Ptr(g.netSales)
			g.detail.AmountFinal = float64Ptr(g.netSales)
			details = append(details, g.detail)
		}
		payType := int64(1)
		dataStatus := int(importDataStatus)
		dataSource := int64(importDataSource)
		isSalesMapping := true
		parsed = append(parsed, entity.CreateOrderBody{
			RoNo:              documentNo,
			CustId:            custId,
			ParentCustId:      parentCustId,
			RoDate:            &first.documentDate,
			ValDate:           nil,
			DueDate:           &first.documentDate,
			SalesmanId:        int64(first.salesman.SalesmanId),
			WhId:              int64PtrForImport(int64(first.salesman.WhId)),
			OutletID:          int64(first.outlet.OutletId),
			OutletAddress1:    nil,
			DeliveryDate:      &first.documentDate,
			OrderNo:           nil,
			PoNo:              nil,
			VehicleNo:         nil,
			PayType:           &payType,
			ReffNo:            nil,
			MobileID:          nil,
			SubTotal:          float64Ptr(subTotal),
			SubTotalFinal:     float64Ptr(subTotal),
			Disc:              float64Ptr(0),
			DiscValue:         float64Ptr(discTotal),
			DiscValueFinal:    float64Ptr(discTotal),
			PromoValue:        float64Ptr(promoTotal),
			PromoValueFinal:   float64Ptr(promoTotal),
			PromoBgValue:      nil,
			PromoBgValueFinal: nil,
			CashDiscValue:     float64Ptr(0),
			TotDisc1:          float64Ptr(0),
			TotDisc2:          float64Ptr(0),
			Vat:               float64Ptr(0),
			VatValue:          float64Ptr(vatTotal),
			VatValueFinal:     float64Ptr(vatTotal),
			Total:             float64Ptr(netTotal),
			TotalFinal:        float64Ptr(netTotal),
			DataStatus:        &dataStatus,
			CreatedBy:         &userId,
			DataSource:        &dataSource,
			Details:           entity.OrderDetWithGroup{Normal: details},
			OrderType:         nil,
			IsClosed:          false,
			Notes:             nil,
			InvoiceNo:         strPtrIfNotEmpty(documentNo),
			InvoiceDate:       &first.documentDate,
			IsSalesMapping:    &isSalesMapping,
		})
	}
	startDate := ""
	endDate := ""
	amount := 0.0
	outletSet := make(map[int64]struct{})
	for _, order := range parsed {
		if order.RoDate != nil && *order.RoDate != "" {
			if startDate == "" || *order.RoDate < startDate {
				startDate = *order.RoDate
			}
			if endDate == "" || *order.RoDate > endDate {
				endDate = *order.RoDate
			}
		}
		if order.Total != nil {
			amount += *order.Total
		}
		outletSet[order.OutletID] = struct{}{}
	}
	return parsed, nil, entity.OrderImportSummary{
		StartDate:       startDate,
		EndDate:         endDate,
		NumberOfInvoice: len(parsed),
		NumberOfOutlet:  len(outletSet),
		Amount:          amount,
		FailedReasons:   []string{},
	}, nil
}

func int64PtrForImport(v int64) *int64 { return &v }

func int64Ptr(v int64) *int64 { return &v }

func float64PtrValueOrDefault(p *float64, def float64) *float64 {
	if p == nil {
		v := def
		return &v
	}
	return p
}

// import-related constants per SX-2434 product spec.
const (
	importDataStatus     = 6 // Invoicing
	importDataSource     = 3 // Import
	defaultImportVat     = 11
	defaultImportVatRate = 11
)

func intPtrFromFloat32(v float32) *int {
	i := int(v)
	return &i
}

func isImportedOrder(req *entity.CreateOrderBody) bool {
	return req != nil && req.DataSource != nil && *req.DataSource == importDataSource
}

func formatImportFailedReason(e entity.OrderImportError) string {
	field := strings.TrimSpace(e.Field)
	if field == "" {
		return fmt.Sprintf("row %d: %s", e.Row, e.Message)
	}
	return fmt.Sprintf("row %d: %s %s", e.Row, field, e.Message)
}

// parseImportDate accepts the date formats documented in the import spec:
// DD/MM/YYYY is the primary user-facing format; YYYY-MM-DD is kept as a
// fallback for compatibility with automated exports. Excel date cells are
// also accepted: when a cell stores a date as a numeric serial we convert
// it to time using the standard 1900-based epoch (matching excelize).
func parseImportDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	// Strip legacy apostrophe prefix used by some Excel exports.
	dateStr = strings.TrimPrefix(dateStr, "'")
	if dateStr != "" {
		if n, err := strconv.ParseFloat(dateStr, 64); err == nil {
			return excelSerialToTime(n), nil
		}
	}
	for _, layout := range []string{"02/01/2006", "2006-01-02", "01-02-06"} {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %q", dateStr)
}

// excelSerialToTime converts an Excel date serial to a time.Time. Excel
// uses 1900-01-01 as serial 1 (with the well-known 1900 leap year bug,
// which is small enough to ignore for our use case since the values we
// handle are well past 1900-03-01).
func excelSerialToTime(serial float64) time.Time {
	const secondsPerDay = 24 * 60 * 60
	epoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	return epoch.Add(time.Duration(serial * secondsPerDay * float64(time.Second)))
}
func strPtrIfNotEmpty(v string) *string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return &v
}
func intPtrFromString(v string) *int {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	n, _ := strconv.Atoi(v)
	return &n
}
func float64PtrFromString(v string) *float64 {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	f, _ := strconv.ParseFloat(v, 64)
	return &f
}
func parseRequiredFloat(v string) (float64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}
func parseOptionalFloat(v string) (float64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}
func intPtrFromStringStrict(v string) (*int, bool) {
	if strings.TrimSpace(v) == "" {
		return nil, true
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil, false
	}
	return &n, true
}
func float64PtrFromStringStrict(v string) (*float64, bool) {
	if strings.TrimSpace(v) == "" {
		return nil, true
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, false
	}
	return &f, true
}
func payTypeFromName(v string) *int64 {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "cash on delivery":
		v := int64(entity.PAY_TYPE_CASH_ON_DELIVERY)
		return &v
	case "cash before delivery":
		v := int64(entity.PAY_TYPE_CASH_BEFORE_DELIVERY)
		return &v
	case "credit":
		v := int64(entity.PAY_TYPE_CREDIT)
		return &v
	}
	return nil
}
