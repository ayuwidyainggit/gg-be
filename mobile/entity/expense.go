package entity

// ExpenseQueryFilter represents query parameters for expense list
type ExpenseQueryFilter struct {
	CustId    string
	EmpID     int64
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	Sort      string `query:"sort"`
	StartDate *int64 `query:"start_date"`
	EndDate   *int64 `query:"end_date"`
}

// ExpenseListResponse represents expense data in list response
type ExpenseListResponse struct {
	ExpenseID   string  `json:"expense_id"`
	DocNo       string  `json:"doc_no"`
	Date        string  `json:"date"`
	ExpenseName string  `json:"expense_name"`
	Amount      float64 `json:"amount"`
	Reason      string  `json:"reason"`
}

// ExpenseDetailResponse represents expense detail response
type ExpenseDetailResponse struct {
	ExpenseID   string                 `json:"expense_id"`
	DocNo       string                 `json:"doc_no"`
	Date        string                 `json:"date"`
	ExpenseName string                 `json:"expense_name"`
	Amount      float64                `json:"amount"`
	Reason      string                 `json:"reason"`
	IsClockOut  *bool                  `json:"is_clock_out"`
	Visits      []ExpenseVisitResponse `json:"visits"`
	File        []ExpenseFileResponse  `json:"file"`
}

// ExpenseVisitResponse represents outlet visit in expense detail
type ExpenseVisitResponse struct {
	OutletID       string `json:"outlet_id"`
	OutletCode     string `json:"outlet_code"`
	OutletName     string `json:"outlet_name"`
	OutletAddress1 string `json:"outlet_address1"`
}

// ExpenseFileResponse represents file attachment in expense detail
type ExpenseFileResponse struct {
	FileID        string `json:"file_id"`
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	FileKey       string `json:"file_key"`
	MediaCategory string `json:"media_category"`
	FileURL       string `json:"file_url"`
	FileSize      int64  `json:"file_size"`
}

// ExpenseTypeLookupResponse represents expense type in lookup response
type ExpenseTypeLookupResponse struct {
	ExpenseTypeID   int    `json:"expense_type_id"`
	ExpenseTypeCode string `json:"expense_type_code"`
	ExpenseTypeName string `json:"expense_type_name"`
}

// ExpenseTypeQueryFilter represents query parameters for expense type lookup
type ExpenseTypeQueryFilter struct {
	IsActive *int   `query:"is_active"`
	Mode     string `query:"mode" validate:"required,eq=lookup"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

// OutletLookupResponse represents outlet in lookup response
type OutletLookupResponse struct {
	OutletID      string  `json:"outlet_id"`
	OutletCode    string  `json:"outlet_code"`
	OutletName    string  `json:"outlet_name"`
	Address       string  `json:"address"`
	Latitude      *string `json:"latitude"`
	Longitude     *string `json:"longitude"`
	OutletStatus  *int    `json:"outlet_status"`
	DistributorID *int    `db:"distributor_id" json:"distributor_id"`
	RegionID      *int    `db:"region_id" json:"region_id"`
	AreaID        *int    `db:"area_id" json:"area_id"`
}

// OutletLookupQueryFilter represents query parameters for outlet lookup
type OutletLookupQueryFilter struct {
	IsActive        *int   `query:"is_active"`
	Statuses        []int  `query:"status"`
	Search          string `query:"q"`
	Sort            string `query:"sort"`
	Page            int    `query:"page"`
	Limit           int    `query:"limit"`
	RegionID        []int  `query:"region_id"`
	AreaID          []int  `query:"area_id"`
	DistributorID   []int  `query:"distributor_id"`
	OutletPrinciple bool   `query:"outlet_principle"`
}

// CreateExpenseBody represents request body for creating expense (multipart/form-data)
type CreateExpenseBody struct {
	CustID        string  `form:"-"`
	EmpID         int64   `form:"-"`
	ExpenseTypeID int     `form:"expense_type_id" validate:"required"`
	OutletID      []int   `form:"outlet_id[]" validate:"omitempty"`
	Amount        float64 `form:"amount" validate:"required,gt=0"`
	Note          string  `form:"note" validate:"omitempty"`
	Folder        string  `form:"folder" validate:"required"`
}

// UpdateExpenseBody represents request body for updating expense (multipart/form-data)
type UpdateExpenseBody struct {
	CustID        string  `form:"-"`
	EmpID         int64   `form:"-"`
	ExpenseTypeID int     `form:"expense_type_id" validate:"required"`
	OutletID      []int   `form:"outlet_id[]" validate:"omitempty"`
	Amount        float64 `form:"amount" validate:"required,gt=0"`
	Note          string  `form:"note" validate:"omitempty"`
	Folder        string  `form:"folder" validate:"required"`
	DeleteFileIDs []int64 `form:"delete_file_ids" validate:"omitempty"`
}

// DetailExpenseParams represents path parameters for expense detail
type DetailExpenseParams struct {
	ExpenseID string `params:"expense_id" validate:"required"`
}

// UpdateExpenseParams represents path parameters for expense update
type UpdateExpenseParams struct {
	ExpenseID string `params:"expense_id" validate:"required"`
}

// DeleteExpenseParams represents path parameters for expense delete
type DeleteExpenseParams struct {
	ExpenseID string `params:"expense_id" validate:"required"`
}

// ExpenseListDataResponse represents the data wrapper for expense list
type ExpenseListDataResponse struct {
	IsClockOut  bool                  `json:"is_clock_out"`
	ExpenseData []ExpenseListResponse `json:"expense_data"`
}
