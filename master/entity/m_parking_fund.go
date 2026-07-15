package entity

import "time"

type MParkingFundResponse struct {
	ParkingFundId int        `json:"parking_fund_id"`
	OutletId      int64      `json:"outlet_id"`
	OutletCode    string     `json:"outlet_code"`
	OutletName    string     `json:"outlet_name"`
	ProId         int64      `json:"pro_id"`
	ProCode       string     `json:"pro_code"`
	ProName       string     `json:"pro_name"`
	PDisc         float64    `json:"p_disc"`
	UpdatedAt     *time.Time `json:"updated_at"`
	UpdatedByName *string    `json:"updated_by_name"`
}

type CreateMParkingFundBody struct {
	CustId    string  `json:"cust_id"`
	OutletId  int     `json:"outlet_id"`
	ProId     int     `json:"pro_id"`
	PDisc     float64 `json:"p_disc"`
	CreatedBy int64   `json:"created_by"`
	UpdatedBy int64   `json:"updated_by"`
}

type DetailMParkingFundParams struct {
	ParkingFundId int `params:"parking_fund_id" validate:"required"`
}

type UpdateMParkingFundParams struct {
	ParkingFundId int `params:"parking_fund_id" validate:"required"`
}

type DeleteMParkingFundParams struct {
	ParkingFundId int `params:"parking_fund_id" validate:"required"`
}

type UpdateMParkingFundRequest struct {
	CustId    string     `json:"cust_id"`
	OutletId  int64      `json:"outlet_id"`
	ProId     int64      `json:"pro_id"`
	PDisc     float64    `json:"p_disc"`
	UpdatedBy int64      `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}
