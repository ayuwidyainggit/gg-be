package model

import "time"

type MParkingFund struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	ParkingFundId int        `json:"parking_fund_id" db:"parking_fund_id"`
	OutletId      *int       `json:"outlet_id" db:"outlet_id"`
	ProId         *int       `json:"pro_id" db:"pro_id"`
	PDisc         *float64   `json:"p_disc" db:"p_disc"`
	CreatedBy     *int64     `json:"created_by" db:"created_by"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

type MParkingFundList struct {
	CustId        string     `json:"cust_id" db:"cust_id"`
	ParkingFundId int        `json:"parking_fund_id" db:"parking_fund_id"`
	OutletId      *int       `json:"outlet_id" db:"outlet_id"`
	OutletCode    *string    `json:"outlet_code" db:"outlet_code"`
	OutletName    *string    `json:"outlet_name" db:"outlet_name"`
	ProId         *int       `json:"pro_id" db:"pro_id"`
	ProCode       *string    `json:"pro_code" db:"pro_code"`
	ProName       *string    `json:"pro_name" db:"pro_name"`
	PDisc         *float64   `json:"p_disc" db:"p_disc"`
	CreatedBy     *int64     `json:"created_by" db:"created_by"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedBy     *int64     `json:"updated_by" db:"updated_by"`
	UpdatedByName *string    `json:"updated_by_name" db:"updated_by_name"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
	IsDel         bool       `json:"is_del" db:"is_del"`
	DeletedBy     *int64     `json:"deleted_by" db:"deleted_by"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

type MParkingFundUpdate struct {
	OutletId  *int64     `json:"outlet_id" sql:"outlet_id"`
	ProId     *int64     `json:"pro_id" sql:"pro_id"`
	PDisc     *float64   `json:"p_disc" sql:"p_disc"`
	UpdatedBy *int64     `json:"updated_by" sql:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at" sql:"updated_at"`
}
