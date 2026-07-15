package entity

type ProductAssignmentQueryFilter struct {
	CustId         string
	ParentCustId   string
	DistributorId  int64
	Page           int      `query:"page"`
	Limit          int      `query:"limit" validate:"required"`
	Query          string   `query:"q"`
	Sort           string   `query:"sort"`
	AssignmentType []string `query:"assignment_type"`
	QDistributorId []int64  `query:"distributor_id"`
	ProId          []int64  `query:"pro_id"`
}

type ProductAssignmentResponse struct {
	ID              int64  `json:"id"`
	CustID          string `json:"cust_id"`
	ProID           int64  `json:"pro_id,omitempty"`
	ProCode         string `json:"pro_code"`
	ProName         string `json:"pro_name,omitempty"`
	DistributorID   int64  `json:"distributor_id,omitempty"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name,omitempty"`
	AssignmentType  string `json:"assignment_type,omitempty"`
	CreatedBy       int64  `json:"created_by,omitempty"`
	CreatedByName   string `json:"created_by_name,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
}

type ProductAssignmentImportRequest struct {
	FileURL string `json:"file_url" validate:"required"`
}

type ProductAssignmentImportResponse struct {
	FileURL       string   `json:"file_url"`
	TotalRow      int      `json:"total_row"`
	SuccessRow    int      `json:"success_row"`
	FailedRow     int      `json:"failed_row"`
	FailedReasons []string `json:"failed_reasons"`
	ProcessedAt   string   `json:"processed_at"`
}

type ProductAssignmentDownloadTemplateResponse struct {
	Message   string                        `json:"message"`
	Data      ProductAssignmentTemplateData `json:"data"`
	RequestID string                        `json:"request_id"`
}

type ProductAssignmentTemplateData struct {
	Name           string `json:"name"`
	DownloadBase64 string `json:"download_base64"`
}

type ProductAssignmentImportRow struct {
	ProCode         string `db:"pro_code"`
	DistributorCode string `db:"distributor_code"`
}
