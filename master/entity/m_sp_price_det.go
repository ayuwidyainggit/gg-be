package entity

type MSpPriceDetGroup struct {
	SalesTeam   []MSpPriceDet `json:"sales_team,omitempty"`
	OutletType  []MSpPriceDet `json:"outlet_type"`
	OutletGroup []MSpPriceDet `json:"outlet_group"`
	Outlet      []MSpPriceDet `json:"outlet"`
	Output      []MSpPriceDet `json:"output,omitempty"`
}

type MSpPriceDet struct {
	CustID        string   `json:"cust_id,omitempty"`
	SpPriceDetID  int      `json:"sp_price_det_id,omitempty"`
	RefType       *int     `json:"ref_type,omitempty"`
	SpPriceID     int      `json:"sp_price_id,omitempty"`
	RefID         *int64   `json:"ref_id"`
	NewSellPrice1 *float64 `json:"new_sell_price1"`
	NewSellPrice2 *float64 `json:"new_sell_price2"`
	NewSellPrice3 *float64 `json:"new_sell_price3"`
	CreatedBy     string   `json:"created_by,omitempty"`
}

type MSpPriceDetRespGroup struct {
	SalesTeam   []MSpPriceDetResp `json:"sales_team,omitempty"`
	OutletType  []MSpPriceDetResp `json:"outlet_type"`
	OutletGroup []MSpPriceDetResp `json:"outlet_group"`
	Outlet      []MSpPriceDetResp `json:"outlet"`
}

type MSpPriceDetResp struct {
	RefType       int     `json:"ref_type,omitempty"`
	RefID         int64   `json:"ref_id"`
	RefName       string  `json:"ref_name,omitempty"`
	SellPrice1    float64 `json:"sell_price1"`
	SellPrice2    float64 `json:"sell_price2"`
	SellPrice3    float64 `json:"sell_price3"`
	NewSellPrice1 float64 `json:"new_sell_price1"`
	NewSellPrice2 float64 `json:"new_sell_price2"`
	NewSellPrice3 float64 `json:"new_sell_price3"`
}

type MSpPriceDetUpdateGroup struct {
	SalesTeam   []MSpPriceDetUpdate `json:"sales_team,omitempty"`
	OutletType  []MSpPriceDetUpdate `json:"outlet_type"`
	OutletGroup []MSpPriceDetUpdate `json:"outlet_group"`
	Outlet      []MSpPriceDetUpdate `json:"outlet"`
}

type MSpPriceDetUpdate struct {
	CustID        string   `json:"cust_id"`
	SpPriceDetID  *string  `json:"sp_price_det_id"`
	RefType       *int64   `json:"ref_type"`
	SpPriceID     string   `json:"sp_price_id"`
	RefID         *int64   `json:"ref_id"`
	NewSellPrice1 *float64 `json:"new_sell_price1"`
	NewSellPrice2 *float64 `json:"new_sell_price2"`
	NewSellPrice3 *float64 `json:"new_sell_price3"`
	CreatedBy     *int64   `json:"created_by"`
}
