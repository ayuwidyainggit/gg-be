package entity

type CreateApDistributorDiscountBody struct {
	CustID                string   `json:"cust_id"`
	DistributorDiscountId *int64   `json:"distributor_discount_id"`
	ProId                 *int64   `json:"pro_id"`
	PurchPrice            *float64 `json:"purch_price"`
	Discount              *float64 `json:"discount"`
	NetPrice              *float64 `json:"net_price"`
	IsActive              bool     `json:"is_active"`
	CreatedBy             int64    `json:"created_by"`
	UpdatedBy             int64    `json:"updated_by"`
}

type UpdateApDistributorDiscountBody struct {
	CustID     string   `json:"cust_id"`
	ProId      *int64   `json:"pro_id"`
	PurchPrice *float64 `json:"purch_price"`
	Discount   *float64 `json:"discount"`
	NetPrice   *float64 `json:"net_price"`
	IsActive   *bool    `json:"is_active"`
	UpdatedBy  int64    `json:"updated_by"`
}
type ApDistributorDiscountResponse struct {
	DistributorDiscountId int      `json:"distributor_discount_id"`
	ProId                 *int64   `json:"pro_id"`
	ProCode               *string  `json:"pro_code"`
	ProName               *string  `json:"pro_name"`
	PurchPrice            *float64 `json:"purch_price"`
	Discount              *float64 `json:"discount"`
	NetPrice              *float64 `json:"net_price"`
	IsActive              bool     `json:"is_active"`
	UpdatedAt             *string  `json:"updated_at"`
	UpdatedByName         string   `json:"updated_by_name"`
}
type ApDistributorDiscountListResponse struct {
	DistributorDiscountId int64    `json:"distributor_discount_id"`
	ProId                 *int64   `json:"pro_id"`
	ProCode               *string  `json:"pro_code"`
	ProName               *string  `json:"pro_name"`
	PurchPrice            *float64 `json:"purch_price"`
	Discount              *float64 `json:"discount"`
	NetPrice              *float64 `json:"net_price"`
	IsActive              bool     `json:"is_active"`
	UpdatedByName         string   `json:"updated_by_name"`
	UpdatedAt             *string  `json:"updated_at"`
}

type ApDistributorDiscountLookupListResponse struct {
	DistributorDiscountId int64    `json:"distributor_discount_id"`
	ProId                 *int64   `json:"pro_id"`
	ProCode               *string  `json:"pro_code"`
	ProName               *string  `json:"pro_name"`
	PurchPrice            *float64 `json:"purch_price"`
	Discount              *float64 `json:"discount"`
	NetPrice              *float64 `json:"net_price"`
	IsActive              bool     `json:"is_active"`
}

type DetailApDistributorDiscountParams struct {
	DistributorDiscountId int64 `params:"distributor_discount_id" validate:"required"`
}

type UpdateApDistributorDiscountParams struct {
	DistributorDiscountId int64 `params:"distributor_discount_id" validate:"required"`
}
type DeleteApDistributorDiscountParams struct {
	DistributorDiscountId int64 `params:"distributor_discount_id" validate:"required"`
}
