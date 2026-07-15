package entity

import "encoding/json"

type ResponseShipmentInvoive struct {
	SalesmanID      int                              `json:"salesman_id"`
	SalesmanCode    string                           `json:"salesman_code"`
	SalesName       string                           `json:"sales_name"`
	WhID            int                              `json:"wh_id"`
	WhCode          string                           `json:"wh_code"`
	WhName          string                           `json:"wh_name"`
	WhLatitude      string                           `json:"wh_latitude"`
	WhLongitude     string                           `json:"wh_longitude"`
	OutletID        int                              `json:"outlet_id"`
	OutletCode      string                           `json:"outlet_code"`
	OutletName      string                           `json:"outlet_name"`
	OutletAddress   string                           `json:"outlet_address"`
	OutletLatitude  string                           `json:"outlet_latitude"`
	OutletLongitude string                           `json:"outlet_longitude"`
	DeliveryDate    string                           `json:"delivery_date"`
	OrderNo         string                           `json:"order_no"`
	PoNo            interface{}                      `json:"po_no"`
	VehicleNo       interface{}                      `json:"vehicle_no"`
	PayType         int                              `json:"pay_type"`
	PayTypeName     string                           `json:"pay_type_name"`
	ReffNo          interface{}                      `json:"reff_no"`
	MobileID        int                              `json:"mobile_id"`
	Disc            float64                          `json:"disc"`
	DiscValue       float64                          `json:"disc_value"`
	PromoValue      float64                          `json:"promo_value"`
	PromoBgValue    float64                          `json:"promo_bg_value"`
	SubTotal        float64                          `json:"sub_total"`
	Vat             float64                          `json:"vat"`
	VatValue        float64                          `json:"vat_value"`
	Total           float64                          `json:"total"`
	DataStatus      int                              `json:"data_status"`
	DataStatusName  string                           `json:"data_status_name"`
	DueDate         string                           `json:"due_date"`
	InvoiceNo       string                           `json:"invoice_no"`
	InvoiceDate     string                           `json:"invoice_date"`
	TotalWeight     float64                          `json:"total_weight"`
	TotalVolume     float64                          `json:"total_volume"`
	Details         []ResponseShipmentInvoiveDetails `json:"details"`
}

type CustomShipmentInvoice struct {
	SalesmanID      int                              `json:"salesman_id"`
	SalesmanCode    string                           `json:"salesman_code"`
	SalesName       string                           `json:"sales_name"`
	WhID            int                              `json:"wh_id"`
	WhCode          string                           `json:"wh_code"`
	WhName          string                           `json:"wh_name"`
	WhLatitude      string                           `json:"wh_latitude"`
	WhLongitude     string                           `json:"wh_longitude"`
	OutletID        int                              `json:"outlet_id"`
	OutletCode      string                           `json:"outlet_code"`
	OutletName      string                           `json:"outlet_name"`
	OutletAddress   string                           `json:"outlet_address"`
	OutletLatitude  string                           `json:"outlet_latitude"`
	OutletLongitude string                           `json:"outlet_longitude"`
	DeliveryDate    string                           `json:"delivery_date"`
	OrderNo         string                           `json:"order_no"`
	PoNo            interface{}                      `json:"po_no"`
	VehicleNo       interface{}                      `json:"vehicle_no"`
	PayType         int                              `json:"pay_type"`
	PayTypeName     string                           `json:"pay_type_name"`
	ReffNo          interface{}                      `json:"reff_no"`
	MobileID        int                              `json:"mobile_id"`
	Disc            float64                          `json:"disc"`
	DiscValue       float64                          `json:"disc_value"`
	TotalPromo      float64                          `json:"total_promo"` // Mapped from promo_value
	TotalBruto      float64                          `json:"total_bruto"` // Mapped from sub_total
	Vat             float64                          `json:"vat"`
	TotalPPN        float64                          `json:"total_ppn"`   // Mapped from vat_value
	TotalNetto      float64                          `json:"total_netto"` // Mapped from total
	DataStatus      int                              `json:"data_status"`
	DataStatusName  string                           `json:"data_status_name"`
	DueDate         string                           `json:"due_date"`
	InvoiceNo       string                           `json:"invoice_no"`
	InvoiceDate     string                           `json:"invoice_date"`
	TotalWeight     float64                          `json:"total_weight"`
	TotalVolume     float64                          `json:"total_volume"`
	Details         []ResponseShipmentInvoiveDetails `json:"details"`
}

type ResponseShipmentInvoiveDetails struct {
	OrderDetailID   int     `json:"order_detail_id"`
	SeqNo           int     `json:"seq_no"`
	ProID           int     `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	ItemType        int     `json:"item_type"`
	Qty1            int     `json:"qty1"`
	Qty2            int     `json:"qty2"`
	Qty3            int     `json:"qty3"`
	Stock1          int     `json:"stock1"`
	Stock2          int     `json:"stock2"`
	Stock3          int     `json:"stock3"`
	PurchPrice1     float32 `json:"purch_price1"`
	PurchPrice2     float32 `json:"purch_price2"`
	PurchPrice3     float32 `json:"purch_price3"`
	SellPrice1      float32 `json:"sell_price1"`
	SellPrice2      float32 `json:"sell_price2"`
	SellPrice3      float32 `json:"sell_price3"`
	ConvUnit1       int     `json:"conv_unit1"`
	ConvUnit2       int     `json:"conv_unit2"`
	ConvUnit3       int     `json:"conv_unit3"`
	UnitID1         string  `json:"unit_id1"`
	UnitID2         string  `json:"unit_id2"`
	UnitID3         string  `json:"unit_id3"`
	Amount          float64 `json:"amount"`
	DiscValue       int     `json:"disc_value"`
	BatchNo         string  `json:"batch_no"`
	ExpDate         string  `json:"exp_date"`
	PriceIncludePpn int     `json:"price_include_ppn"`
	PriceExcludePpn int     `json:"price_exclude_ppn"`
	NetValue        int     `json:"net_value"`
	Vat             float64 `json:"vat"`
	VatValue        float64 `json:"vat_value"`
	Volume1         float64 `json:"volume1"`
	Volume2         float64 `json:"volume2"`
	Volume3         float64 `json:"volume3"`
	Weight1         float64 `json:"weight1"`
	Weight2         float64 `json:"weight2"`
	Weight3         float64 `json:"weight3"`
}

func (r *ResponseShipmentInvoiveDetails) UnmarshalJSON(data []byte) error {
	// bikin struct alias biar gak infinite recursion
	type Alias ResponseShipmentInvoiveDetails
	aux := &struct {
		Volume1 float64 `json:"volume1"`
		Volume2 float64 `json:"volume2"`
		Volume3 float64 `json:"volume3"`
		Weight1 float64 `json:"weight1"`
		Weight2 float64 `json:"weight2"`
		Weight3 float64 `json:"weight3"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	// unmarshal dulu ke aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// mapping sesuai kebutuhan
	r.Volume1 = aux.Volume3
	r.Volume2 = aux.Volume2
	r.Volume3 = aux.Volume1

	r.Weight1 = aux.Weight3
	r.Weight2 = aux.Weight2
	r.Weight3 = aux.Weight1

	return nil
}
