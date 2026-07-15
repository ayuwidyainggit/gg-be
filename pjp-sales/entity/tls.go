package entity

type CreateTlsBody struct {
	CustId     string             `json:"cust_id"`
	TlsDate    *string            `json:"tls_date"`
	SalesmanId *int64             `json:"salesman_id"`
	Notes      *string            `json:"notes"`
	DataStatus *int64             `json:"data_status"`
	CreatedBy  *int64             `json:"created_by"`
	Details    []CreateTlsDetBody `json:"details"`
}

type TlsResponse struct {
	TlsId         int64            `json:"tls_id"`
	TlsDate       *string          `json:"tls_date"`
	SalesmanId    *int64           `json:"salesman_id"`
	SalesmanCode  *string          `json:"salesman_code"`
	SalesmanName  *string          `json:"salesman_name"`
	Notes         *string          `json:"notes"`
	DataStatus    *int64           `json:"data_status"`
	UpdatedAt     string           `json:"updated_at"`
	UpdatedByName string           `json:"updated_by_name"`
	Details       []TlsDetResponse `json:"details"`
}

type DetailTlsParams struct {
	TlsId int `params:"tls_id" validate:"required"`
}
type DeleteTlsParams struct {
	TlsId int `params:"tls_id" validate:"required"`
}

type UpdateTlsParams struct {
	TlsId int `params:"tls_id" validate:"required"`
}

type TlsListResponse struct {
	TlsId         int64   `json:"tls_id"`
	TlsDate       *string `json:"tls_date"`
	SalesmanId    *int64  `json:"salesman_id"`
	SalesmanCode  string  `json:"salesman_code"`
	SalesmanName  string  `json:"salesman_name"`
	Notes         *string `json:"notes"`
	DataStatus    *int64  `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
}

type UpdateTlsBody struct {
	CustId     string             `json:"cust_id"`
	TlsDate    *string            `json:"tls_date"`
	SalesmanId *int64             `json:"salesman_id"`
	Notes      *string            `json:"notes"`
	DataStatus *int64             `json:"data_status"`
	CreatedBy  *int64             `json:"created_by"`
	CreatedAt  *string            `json:"created_at"`
	UpdatedBy  int64              `json:"updated_by"`
	Details    []UpdateTlsDetBody `json:"details"`
}
