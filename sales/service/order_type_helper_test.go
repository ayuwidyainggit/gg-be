package service

import (
	"context"
	"testing"

	"sales/entity"
	"sales/model"
)

func TestResolveCreateOrderDataStatusForTakingOrder(t *testing.T) {
	processedDecision := salesOrderStatusDecision{DataStatus: int64(entity.PROCESSED)}
	needReviewDecision := salesOrderStatusDecision{DataStatus: int64(entity.NEED_REVIEW)}
	orderTypeTaking := "O"
	orderTypeSales := "SO"

	tests := []struct {
		name      string
		orderType *string
		decision  salesOrderStatusDecision
		expected  int64
	}{
		{name: "taking order forces need review from processed decision", orderType: &orderTypeTaking, decision: processedDecision, expected: int64(entity.NEED_REVIEW)},
		{name: "taking order keeps need review decision", orderType: &orderTypeTaking, decision: needReviewDecision, expected: int64(entity.NEED_REVIEW)},
		{name: "sales order keeps existing decision", orderType: &orderTypeSales, decision: processedDecision, expected: int64(entity.PROCESSED)},
		{name: "nil order type keeps existing decision", orderType: nil, decision: processedDecision, expected: int64(entity.PROCESSED)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveCreateOrderDataStatus(tt.orderType, tt.decision); got != tt.expected {
				t.Fatalf("resolveCreateOrderDataStatus() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestStore_SX2184TakingOrderSkipsStockMutationAndPersistsOriginalQty(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-06-08"
	orderType := "O"
	zero := 0.0
	qty1 := 2.0
	qty2 := 0.0
	qty3 := 1.0
	qtyPo1 := 7.0
	qtyPo2 := 1.0
	qtyPo3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 1000.0
	vat := 11.0

	var storedOrder model.Order
	var storedDetails []model.OrderDetail
	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error {
			storedOrder = *data
			return nil
		},
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := len(storedDetails) + 1
			data.OrderDetailID = &id
			storedDetails = append(storedDetails, *data)
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		OrderType:    &orderType,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			QtyPo1:          &qtyPo1,
			QtyPo2:          &qtyPo2,
			QtyPo3:          &qtyPo3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	request.DataStatus = intPtrForTest(int64(entity.PROCESSED))
	response, err := service.Store(request, validateResponse)
	if err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if response.RoNo == "" {
		t.Fatalf("ro number must be generated")
	}
	if stockWriteCalled {
		t.Fatal("expected no stock mutation for taking order create")
	}
	if len(storedDetails) != 1 {
		t.Fatalf("expected 1 stored detail, got %d", len(storedDetails))
	}
	if storedOrder.OrderType == nil || *storedOrder.OrderType != orderType {
		t.Fatalf("expected stored order_type O, got %+v", storedOrder.OrderType)
	}
	if storedOrder.OprType == nil || *storedOrder.OprType != orderType {
		t.Fatalf("expected stored opr_type O, got %+v", storedOrder.OprType)
	}
	if storedOrder.DataStatus == nil || *storedOrder.DataStatus != int64(entity.NEED_REVIEW) {
		t.Fatalf("expected taking order data_status Need Review even when payload asks Processed, got %+v", storedOrder.DataStatus)
	}
	if storedOrder.ValidateStok == nil || *storedOrder.ValidateStok {
		t.Fatalf("expected validate_stok false for taking order, got %+v", storedOrder.ValidateStok)
	}
	if storedOrder.ValidateStokMessage != nil {
		t.Fatalf("expected nil validate_stok_message for taking order, got %+v", storedOrder.ValidateStokMessage)
	}

	detail := storedDetails[0]
	if detail.Qty != 0 {
		t.Fatalf("expected sales qty to remain zero for taking order, got %v", detail.Qty)
	}
	if detail.QtyFinal != 0 {
		t.Fatalf("expected final qty to remain zero for taking order, got %v", detail.QtyFinal)
	}
	if detail.Qty1 != nil || detail.Qty2 != nil || detail.Qty3 != nil {
		t.Fatalf("expected sales qty tiers nil for taking order, got %+v %+v %+v", detail.Qty1, detail.Qty2, detail.Qty3)
	}
	if got := getValueOrDefault(detail.QtyPo1, 0); got != qty1 {
		t.Fatalf("expected qty_po1 to prefer qty1 payload when both qty and qty_po are present, got %v want %v", got, qty1)
	}
	if got := getValueOrDefault(detail.QtyPo2, 0); got != qty2 {
		t.Fatalf("expected qty_po2 to prefer qty2 payload when both qty and qty_po are present, got %v want %v", got, qty2)
	}
	if got := getValueOrDefault(detail.QtyPo3, 0); got != qty3 {
		t.Fatalf("expected qty_po3 to prefer qty3 payload when both qty and qty_po are present, got %v want %v", got, qty3)
	}
	if got := getValueOrDefault(detail.OriginalQtyPo1, 0); got != qty1 {
		t.Fatalf("expected original_qty_po1 to preserve original qty1 payload, got %v want %v", got, qty1)
	}
	if got := getValueOrDefault(detail.OriginalQtyPo2, 0); got != qty2 {
		t.Fatalf("expected original_qty_po2 to preserve original qty2 payload, got %v want %v", got, qty2)
	}
	if got := getValueOrDefault(detail.OriginalQtyPo3, 0); got != qty3 {
		t.Fatalf("expected original_qty_po3 to preserve original qty3 payload, got %v want %v", got, qty3)
	}
	if detail.QtyPo != 52 {
		t.Fatalf("expected qty_po total from qty1/2/3 conversion, got %v want 52", detail.QtyPo)
	}
}

func TestStore_SX2184NilOrderTypeStillMutatesInventory(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-06-08"
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0

	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			if len(stockUpdates) != 1 {
				t.Fatalf("expected 1 stock update, got %d", len(stockUpdates))
			}
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if !stockWriteCalled {
		t.Fatal("expected stock mutation for nil order_type create")
	}
}

func TestStore_SX2184EmptyOrderTypeStillMutatesInventory(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-06-08"
	orderType := ""
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0

	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			if len(stockUpdates) != 1 {
				t.Fatalf("expected 1 stock update, got %d", len(stockUpdates))
			}
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		OrderType:    &orderType,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if !stockWriteCalled {
		t.Fatal("expected stock mutation for empty order_type create")
	}
}

func TestStore_SX2184OrderTypeCStillMutatesInventory(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-06-08"
	orderType := "C"
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0

	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			if len(stockUpdates) != 1 {
				t.Fatalf("expected 1 stock update, got %d", len(stockUpdates))
			}
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		OrderType:    &orderType,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if !stockWriteCalled {
		t.Fatal("expected stock mutation for order_type C create")
	}
}

func TestStore_SX2184SalesOrderStillMutatesInventory(t *testing.T) {
	custID := "C220010001"
	parentCustID := custID
	whID := int64(301)
	userID := int64(99)
	orderDate := "2026-06-08"
	orderType := "SO"
	zero := 0.0
	qty1 := 1.0
	qty2 := 0.0
	qty3 := 0.0
	sellPrice1 := 100.0
	sellPrice2 := 0.0
	sellPrice3 := 0.0
	vat := 0.0

	stockWriteCalled := false

	orderRepo := &mockOrderRepositoryStore{
		countAllRoByCustIdFn: func(custId string, roDate string) (int, error) { return 0, nil },
		findOutletByIDFn: func(outletID int, custId string, parentCustId string) (model.OutletRead, error) {
			return model.OutletRead{Address1: stringPtr("Outlet Address")}, nil
		},
		findProductByIDFn: func(productID int) (model.ProductRead, error) {
			return model.ProductRead{ProId: productID, SellPrice1: sellPrice1, ConvUnit2: 10, ConvUnit3: 5}, nil
		},
		findDiscountByProductAndOutletFn: func(product model.ProductRead, outlet model.OutletRead) (model.DiscountRead, error) {
			return model.DiscountRead{}, nil
		},
		findDiscountCriteriaBySubTotalFn: func(discountID string, subTotal int) (model.DiscountCriteria, error) {
			return model.DiscountCriteria{}, nil
		},
		storeFn: func(c context.Context, data *model.Order) error { return nil },
		storeDetailFn: func(c context.Context, data *model.OrderDetail) error {
			id := 1
			data.OrderDetailID = &id
			return nil
		},
		storeRewardFn: func(c context.Context, data *model.OrderReward) error { return nil },
	}

	service := &orderServiceImpl{
		OrderRepository: orderRepo,
		StockRepository: &mockStockRepository{salesStockUpdatesFn: func(c context.Context, stockUpdates []*entity.SalesOrderStockUpdate) error {
			stockWriteCalled = true
			if len(stockUpdates) != 1 {
				t.Fatalf("expected 1 stock update, got %d", len(stockUpdates))
			}
			return nil
		}},
		Transaction: &mockDbtransaction{},
	}

	request := entity.CreateOrderBody{
		CustId:       custID,
		ParentCustId: parentCustID,
		RoDate:       &orderDate,
		SalesmanId:   11,
		WhId:         &whID,
		OutletID:     21,
		CreatedBy:    &userID,
		OrderType:    &orderType,
		Details: entity.OrderDetWithGroup{Normal: []entity.CreateOrderDetBody{{
			ProId:           748,
			Qty1:            &qty1,
			Qty2:            &qty2,
			Qty3:            &qty3,
			ConvUnit2:       intPtr(10),
			ConvUnit3:       intPtr(5),
			SellPrice1:      &sellPrice1,
			SellPrice2:      &sellPrice2,
			SellPrice3:      &sellPrice3,
			PromoValue:      &zero,
			PromoValueFinal: &zero,
			DiscValue:       &zero,
			Vat:             &vat,
			VatValue:        &zero,
			Amount:          &zero,
			AmountFinal:     &zero,
		}}},
	}

	validateResponse := entity.ValidateResponse{Validate1Success: true, Validate2Success: true, Validate3Success: true, Validate4Success: true, IsSuccessValidate: true}

	if _, err := service.Store(request, validateResponse); err != nil {
		t.Fatalf("Store returned error: %v", err)
	}
	if !stockWriteCalled {
		t.Fatal("expected stock mutation for sales order create")
	}
}
