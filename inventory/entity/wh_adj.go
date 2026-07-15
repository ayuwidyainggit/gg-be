package entity

type CreateWhAdjBody struct {
	CustID    string               `json:"cust_id"`
	WhID      int64                `json:"wh_id"`
	Notes     string               `json:"notes"`
	CreatedBy int64                `json:"created_by"`
	Details   []CreateWhAdjDetBody `json:"details"`
}
type DetailWhAdjParams struct {
	AdjNo string `params:"adj_no" validate:"required"`
}
type UpdateWhAdjParams struct {
	AdjNo string `params:"adj_no" validate:"required"`
}
type WhAdjResponse struct {
	StockAdjNo    string             `json:"stock_adjustment_no"`
	StockAdjDate  string             `json:"stock_adjustment_date"`
	WhID          int64              `json:"wh_id"`
	WhCode        string             `json:"wh_code"`
	WhName        string             `json:"wh_name"`
	StockType     string             `json:"stock_type"`
	Notes         string             `json:"notes"`
	Status        int                `json:"status"`
	UpdatedAt     string             `json:"updated_at"`
	UpdatedByName string             `json:"updated_by_name"`
	IsClosed      bool               `json:"is_closed"`
	ClosedBy      int64              `json:"closed_by"`
	ClosedByName  string             `json:"closed_by_name"`
	ClosedAt      string             `json:"closed_at"`
	Details       []WhAdjDetresponse `json:"details"`
}
type WhAdjListResponse struct {
	StockAdjNo    string `json:"stock_adjustment_no"`
	StockAdjDate  string `json:"stock_adjustment_date"`
	WhID          int64  `json:"wh_id"`
	WhCode        string `json:"wh_code"`
	WhName        string `json:"wh_name"`
	Notes         string `json:"notes"`
	DataStatus    int64  `json:"status"`
	UpdatedAt     string `json:"updated_at"`
	IsClosed      bool   `json:"is_closed"`
	ClosedBy      int64  `json:"closed_by"`
	ClosedByName  string `json:"closed_by_name"`
	ClosedAt      string `json:"closed_at"`
	UpdatedByName string `json:"updated_by_name"`
}
type UpdateWhAdjStatusBody struct {
	CustID     string `json:"cust_id"`
	DataStatus int    `json:"status" validate:"oneof=2 9"`
	UpdatedBy  int64  `json:"updated_by"`
}

type WhAdjWarehouseQueryFilter struct {
	StartDate *int64 `query:"startDate"`
	EndDate   *int64 `query:"endDate"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}
type WarehouseAdjustment struct {
	WhId   *int64  `json:"wh_id"`
	WhCode *string `json:"wh_code"`
	WhName *string `json:"wh_name"`
}

type WhAdjQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	WhID         []int  `query:"warehouse_id"`
	Status       []int  `query:"status"`
	AdjusmentNo  string `query:"adjustment_no"`
	Sort         string `query:"sort"`
}
