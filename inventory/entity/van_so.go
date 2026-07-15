package entity

type CreateVanSoBody struct {
	CustID     string               `json:"cust_id"`
	VanSoNo    string               `json:"van_so_no"`
	VanSoDate  *string              `json:"van_so_date"`
	TrCode     *string              `json:"tr_code" validate:"required,len=3"`
	EmpID      *int64               `json:"emp_id"`
	Notes      *string              `json:"notes"`
	DataStatus *int64               `json:"data_status"`
	CreatedBy  *int64               `json:"created_by"`
	Details    []CreateVanSoDetBody `json:"details"`
}
type DetailVanSoParams struct {
	VanSoNo string `params:"van_so_no" validate:"required"`
}
type UpdateVanSoParams struct {
	VanSoNo string `params:"van_so_no" validate:"required"`
}
type VanSoResponse struct {
	VanSoNo       string             `json:"van_so_no"`
	VanSoDate     *string            `json:"van_so_date"`
	TrCode        *string            `json:"tr_code"`
	EmpID         *int64             `json:"emp_id"`
	EmpCode       string             `json:"emp_code"`
	EmpName       string             `json:"emp_name"`
	Notes         *string            `json:"notes"`
	DataStatus    *int64             `json:"data_status"`
	UpdatedAt     string             `json:"updated_at"`
	UpdatedByName string             `json:"updated_by_name"`
	IsClosed      bool               `json:"is_closed"`
	ClosedBy      int64              `json:"closed_by"`
	ClosedByName  string             `json:"closed_by_name"`
	ClosedAt      string             `json:"closed_at"`
	Details       []VanSoDetResponse `json:"details"`
}

type VanSoListResponse struct {
	VanSoNo       string  `json:"van_so_no"`
	VanSoDate     *string `json:"van_so_date"`
	TrCode        *string `json:"tr_code"`
	EmpID         *int64  `json:"emp_id"`
	EmpCode       string  `json:"emp_code"`
	EmpName       string  `json:"emp_name"`
	Notes         *string `json:"notes"`
	DataStatus    *int64  `json:"data_status"`
	UpdatedAt     string  `json:"updated_at"`
	UpdatedByName string  `json:"updated_by_name"`
	IsClosed      bool    `json:"is_closed"`
	ClosedBy      int64   `json:"closed_by"`
	ClosedByName  string  `json:"closed_by_name"`
	ClosedAt      string  `json:"closed_at"`
}
type UpdateVanSoBody struct {
	CustID     string               `json:"cust_id"`
	VanSoNo    string               `json:"van_so_no"`
	VanSoDate  *string              `json:"van_so_date"`
	TrCode     *string              `json:"tr_code"`
	EmpID      *int64               `json:"emp_id"`
	Notes      *string              `json:"notes"`
	DataStatus *int64               `json:"data_status"`
	CreatedBy  *int64               `json:"created_by"`
	UpdatedBy  int64                `json:"updated_by"`
	Details    []UpdateVanSoDetBody `json:"details"`
}
