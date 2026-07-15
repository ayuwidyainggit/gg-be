package entity

import "time"

type ManageMinimumPriceQueryFilter struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit" validate:"required"`
	Query    string `query:"q"`
	Mode     string `query:"mode"`
	Sort     string `query:"sort"`
	IsActive *int   `query:"is_active"`
	Status   []int  `query:"status"`
}

type ManageMinimumPrice struct {
	ManageMinimumPriceId     int      `json:"manage_minimum_price_id"`
	BasePrice                *int     `json:"base_price"`
	LimitAction              *int     `json:"limit_action"`
	Threshold                *float64 `json:"threshold"`
	StatusManageMinimumPrice *int     `json:"status_manage_minimum_price"`
	ProId                    *int     `json:"pro_id"`
	Price1                   *float64 `json:"price1"`
	Price2                   *float64 `json:"price2"`
	Price3                   *float64 `json:"price3"`
	Price4                   *float64 `json:"price4"`
	Price5                   *float64 `json:"price5"`
	PriceMinimum1            *float64 `json:"price1_minimum"`
	PriceMinimum2            *float64 `json:"price2_minimum"`
	PriceMinimum3            *float64 `json:"price3_minimum"`
	PriceMinimum4            *float64 `json:"price4_minimum"`
	PriceMinimum5            *float64 `json:"price5_minimum"`
	UnitId1                  *string  `json:"unit_id1"`
	UnitId2                  *string  `json:"unit_id2"`
	UnitId3                  *string  `json:"unit_id3"`
	UnitId4                  *string  `json:"unit_id4"`
	UnitId5                  *string  `json:"unit_id5"`
	ConvUnit2                *int     `json:"conv_unit2"`
	ConvUnit3                *int     `json:"conv_unit3"`
	ConvUnit4                *int     `json:"conv_unit4"`
	ConvUnit5                *int     `json:"conv_unit5"`
}

type ManageMinimumPriceRead struct {
	ManageMinimumPrice
	BasePriceName                string `json:"base_price_name"`
	LimitActionName              string `json:"limit_action_name"`
	StatusManageMinimumPriceName string `json:"status_manage_minimum_price_name"`
	ProductName                  string `json:"pro_name"`
	ProductCode                  string `json:"pro_code"`
}

type CreateManageMinimumPrice struct {
	BasePrice                *int     `json:"base_price"`
	LimitAction              *int     `json:"limit_action"`
	Threshold                *float64 `json:"threshold"`
	StatusManageMinimumPrice *int     `json:"status_manage_minimum_price"`
	ProId                    *int     `json:"pro_id"`
	Price1                   *float64 `json:"price1"`
	Price2                   *float64 `json:"price2"`
	Price3                   *float64 `json:"price3"`
	Price4                   *float64 `json:"price4"`
	Price5                   *float64 `json:"price5"`
	PriceMinimum1            *float64 `json:"price1_minimum"`
	PriceMinimum2            *float64 `json:"price2_minimum"`
	PriceMinimum3            *float64 `json:"price3_minimum"`
	PriceMinimum4            *float64 `json:"price4_minimum"`
	PriceMinimum5            *float64 `json:"price5_minimum"`
	UnitId1                  *string  `json:"unit_id1"`
	UnitId2                  *string  `json:"unit_id2"`
	UnitId3                  *string  `json:"unit_id3"`
	UnitId4                  *string  `json:"unit_id4"`
	UnitId5                  *string  `json:"unit_id5"`
	ConvUnit2                *int     `json:"conv_unit2"`
	ConvUnit3                *int     `json:"conv_unit3"`
	ConvUnit4                *int     `json:"conv_unit4"`
	ConvUnit5                *int     `json:"conv_unit5"`
}

type BodyCreateManageMinimumPrice struct {
	CustId       string
	ParentCustId string
	CreatedBy    int64                      `json:"created_by"`
	CreatedAt    time.Time                  `json:"created_at"`
	UpdatedBy    int64                      `json:"updated_by"`
	UpdatedAt    time.Time                  `json:"updated_at"`
	Body         []CreateManageMinimumPrice `json:"data"`
}

type UpdateManageMinimumPrice struct {
	CustId        string
	ParentCustId  string
	BasePrice     *int     `json:"base_price"`
	LimitAction   *int     `json:"limit_action"`
	Threshold     *float64 `json:"threshold"`
	PriceMinimum1 *float64 `json:"price1_minimum"`
	PriceMinimum2 *float64 `json:"price2_minimum"`
	PriceMinimum3 *float64 `json:"price3_minimum"`
	PriceMinimum4 *float64 `json:"price4_minimum"`
	PriceMinimum5 *float64 `json:"price5_minimum"`
	UpdatedBy     int64    `json:"updated_by"`
}

type DetailManageMinimumPriceParams struct {
	CustId               string
	ParentCustId         string
	ManageMinimumPriceId int64 `params:"manage_minimum_price_id" validate:"required"`
}

type UpdateManageMinimumPriceParams struct {
	ManageMinimumPriceId int64 `params:"manage_minimum_price_id" validate:"required"`
}

type DeleteManageMinimumPriceParams struct {
	ManageMinimumPriceId int64 `params:"manage_minimum_price_id" validate:"required"`
}

type UpdateStatusMinimumPrice struct {
	CustId string `json:"cust_id"`
	UserId int64  `json:"user_id"`
	Status int64  `json:"status"`
}

type BasePriceLookupResponse struct {
	BasePrice     int    `json:"base_price"`
	BasePriceName string `json:"base_price_name"`
}

var BasePrice = []BasePriceLookupResponse{
	{BasePrice: 1, BasePriceName: "Purchase Price"},
	{BasePrice: 2, BasePriceName: "COGS"},
}

type LimitActionLookupResponse struct {
	LimitAction     int    `json:"limit_action"`
	LimitActionName string `json:"limit_action_name"`
}

var LimitAction = []LimitActionLookupResponse{
	{LimitAction: 1, LimitActionName: "Warning"},
	{LimitAction: 2, LimitActionName: "Restricted"},
}

func GetBasePriceName(id int) string {
	for _, item := range BasePrice {
		if item.BasePrice == id {
			return item.BasePriceName
		}
	}
	return "-"
}

func GetLimitActionName(id int) string {
	for _, item := range LimitAction {
		if item.LimitAction == id {
			return item.LimitActionName
		}
	}
	return "-"
}

var StatusManageMinimumPrice = map[int]string{
	0: "Non Active",
	1: "Submit",
	2: "Active",
}

func ConvStatusManageMinimumPrice(param int) string {
	result, ok := StatusManageMinimumPrice[param] // langsung gunakan int sebagai key
	if !ok {
		result = "Unknown"
	}
	return result
}
