package entity

// EditOrderEnhanceBody represents the enhanced edit order request
// FE should send only one of: purchase_order, sales_order, or final_order
// For adding new products, use add_purchase_order, add_sales_order, or add_final_order
type EditOrderEnhanceBody struct {
	CustId          string                    `json:"-"`
	ParentCustId    string                    `json:"-"`
	UpdatedBy       int64                     `json:"-"`
	RoNo            string                    `json:"ro_no,omitempty"`
	PurchaseOrder   []EditPurchaseOrderDetail `json:"purchase_order,omitempty"`
	PurchaseDetails []EditPurchaseOrderDetail `json:"purchase_details,omitempty"`
	SalesOrder      []EditSalesOrderDetail    `json:"sales_order,omitempty"`
	FinalOrder      []EditFinalOrderDetail    `json:"final_order,omitempty"`

	// Backward/alternate request aliases used by some clients
	SalesOrderDetails []EditSalesOrderDetail `json:"sales_order_details,omitempty"`
	FinalOrderDetails []EditFinalOrderDetail `json:"final_order_details,omitempty"`

	// Add new product support
	AddPurchaseOrder   []AddPurchaseOrderDetail `json:"add_purchase_order,omitempty"`
	AddPurchaseDetails []AddPurchaseOrderDetail `json:"add_purchase_details,omitempty"`
	AddSalesOrder      []AddSalesOrderDetail    `json:"add_sales_order,omitempty"`
	AddFinalOrder      []AddFinalOrderDetail    `json:"add_final_order,omitempty"`
}

// EditPurchaseOrderDetail for Case 1 - Purchase Order tab
// Updates qty_po1/2/3 and sell_price_po1/2/3 fields in sls.order_detail
type EditPurchaseOrderDetail struct {
	OrderDetailId int64    `json:"order_detail_id" validate:"required"`
	QtyPo1        *float64 `json:"qty_po1,omitempty"`
	QtyPo2        *float64 `json:"qty_po2,omitempty"`
	QtyPo3        *float64 `json:"qty_po3,omitempty"`
	SellPricePo1  *float64 `json:"sell_price_po1,omitempty"`
	SellPricePo2  *float64 `json:"sell_price_po2,omitempty"`
	SellPricePo3  *float64 `json:"sell_price_po3,omitempty"`
	DiscPo        *float64 `json:"disc_po,omitempty"`
	VatValuePo    *float64 `json:"vat_value_po,omitempty"`
}

// EditSalesOrderDetail for Case 2 - Sales Order tab
// Updates qty1/2/3 and sell_price1/2/3 fields in sls.order_detail
type EditSalesOrderDetail struct {
	OrderDetailId        int64    `json:"order_detail_id" validate:"required"`
	Qty1                 *float64 `json:"qty1,omitempty"`
	Qty2                 *float64 `json:"qty2,omitempty"`
	Qty3                 *float64 `json:"qty3,omitempty"`
	SellPrice1           *float64 `json:"sell_price1,omitempty"`
	SellPrice2           *float64 `json:"sell_price2,omitempty"`
	SellPrice3           *float64 `json:"sell_price3,omitempty"`
	IsProductPromotion   *bool    `json:"is_product_promotion,omitempty"`
	IsProductPromotionSo *bool    `json:"is_product_promotion_so,omitempty"`
}

// EditFinalOrderDetail for Case 3 - Final Order tab
// Updates qty1_final/2_final/3_final and sell_price_final1/2/3 fields in sls.order_detail
type EditFinalOrderDetail struct {
	OrderDetailId           int64    `json:"order_detail_id" validate:"required"`
	Qty1Final               *float64 `json:"qty1_final,omitempty"`
	Qty2Final               *float64 `json:"qty2_final,omitempty"`
	Qty3Final               *float64 `json:"qty3_final,omitempty"`
	SellPriceFinal1         *float64 `json:"sell_price_final1,omitempty"`
	SellPriceFinal2         *float64 `json:"sell_price_final2,omitempty"`
	SellPriceFinal3         *float64 `json:"sell_price_final3,omitempty"`
	IsProductPromotion      *bool    `json:"is_product_promotion,omitempty"`
	IsProductPromotionFinal *bool    `json:"is_product_promotion_final,omitempty"`
}

// AddPurchaseOrderDetail for adding new product from Purchase Order tab
// Cascades to PO, SO, and Final fields in sls.order_detail
type AddPurchaseOrderDetail struct {
	ProId                int64   `json:"pro_id" validate:"required"`
	QtyPo1               float64 `json:"qty_po1"`
	QtyPo2               float64 `json:"qty_po2"`
	QtyPo3               float64 `json:"qty_po3"`
	SellPriceSystem1     float64 `json:"sell_price_system1"`
	SellPriceSystem2     float64 `json:"sell_price_system2"`
	SellPriceSystem3     float64 `json:"sell_price_system3"`
	SellPricePo1         float64 `json:"sell_price_po1"`
	SellPricePo2         float64 `json:"sell_price_po2"`
	SellPricePo3         float64 `json:"sell_price_po3"`
	DiscPo               float64 `json:"disc_po"`
	VatValuePo           float64 `json:"vat_value_po"`
	UnitId1              string  `json:"unit_id1,omitempty"`
	UnitId2              string  `json:"unit_id2,omitempty"`
	UnitId3              string  `json:"unit_id3,omitempty"`
	IsProductPromotionPo *bool   `json:"is_product_promotion_po,omitempty"`
}

// AddSalesOrderDetail for adding new product from Sales Order tab
// Cascades to SO and Final fields in sls.order_detail
type AddSalesOrderDetail struct {
	ProId                int64   `json:"pro_id" validate:"required"`
	Qty1                 float64 `json:"qty1"`
	Qty2                 float64 `json:"qty2"`
	Qty3                 float64 `json:"qty3"`
	SellPriceSystem1     float64 `json:"sell_price_system1"`
	SellPriceSystem2     float64 `json:"sell_price_system2"`
	SellPriceSystem3     float64 `json:"sell_price_system3"`
	SellPrice1           float64 `json:"sell_price1"`
	SellPrice2           float64 `json:"sell_price2"`
	SellPrice3           float64 `json:"sell_price3"`
	UnitId1              string  `json:"unit_id1"`
	UnitId2              string  `json:"unit_id2"`
	UnitId3              string  `json:"unit_id3"`
	Qty1Stock            float64 `json:"qty1_stock"`
	Qty2Stock            float64 `json:"qty2_stock"`
	Qty3Stock            float64 `json:"qty3_stock"`
	IsProductPromotionSo *bool   `json:"is_product_promotion_so,omitempty"`
}

// AddFinalOrderDetail for adding new product from Final Order tab
// Only updates Final fields in sls.order_detail
type AddFinalOrderDetail struct {
	ProId                   int64   `json:"pro_id" validate:"required"`
	Qty1Final               float64 `json:"qty1_final"`
	Qty2Final               float64 `json:"qty2_final"`
	Qty3Final               float64 `json:"qty3_final"`
	SellPriceSystem1        float64 `json:"sell_price_system1"`
	SellPriceSystem2        float64 `json:"sell_price_system2"`
	SellPriceSystem3        float64 `json:"sell_price_system3"`
	SellPriceFinal1         float64 `json:"sell_price_final1"`
	SellPriceFinal2         float64 `json:"sell_price_final2"`
	SellPriceFinal3         float64 `json:"sell_price_final3"`
	IsProductPromotionFinal *bool   `json:"is_product_promotion_final,omitempty"`
}
