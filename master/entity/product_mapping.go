package entity

type ProductMappingListQueryFilter struct {
	CustId       string
	ParentCustId string
	Page         int    `query:"page"`
	Limit        int    `query:"limit"`
	Search       string `query:"search"`
	SortBy       string `query:"sort_by"`
	SortOrder    string `query:"sort_order"`
}

type ProductMappingListItem struct {
	DistributorID   int64  `json:"distributor_id"`
	DistributorCode string `json:"distributor_code"`
	DistributorName string `json:"distributor_name"`
	TotalProduct    int    `json:"total_product"`
	CreatedBy       int64  `json:"created_by,omitempty"`
	CreatedByName   string `json:"created_by_name,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedBy       int64  `json:"updated_by,omitempty"`
	UpdatedByName   string `json:"updated_by_name,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

type ProductMappingDetailQueryFilter struct {
	CustId         string
	ParentCustId   string
	DistributorId  int64
	Page           int    `query:"page"`
	Limit          int    `query:"limit"`
}

type ProductMappingDetailItem struct {
	ProID         int64   `json:"pro_id"`
	ParentProID   int64   `json:"parent_pro_id"`
	ParentProCode string  `json:"parent_pro_code"`
	ParentProName string  `json:"parent_pro_name"`
	ProCode       string  `json:"pro_code"`
	ProName       string  `json:"pro_name"`
	LargestUOM    string  `json:"largest_uom"`
	MiddleUOM     *string `json:"middle_uom"`
	SmallestUOM   *string `json:"smallest_uom"`
}

type ProductMappingDetailResponse struct {
	DistributorCode string                     `json:"distributor_code"`
	DistributorName string                     `json:"distributor_name"`
	TotalProduct    int                        `json:"total_product"`
	Items           []ProductMappingDetailItem `json:"items"`
}

type ProductMappingUpdateRequest struct {
	ProCode  string `json:"pro_code"`
	ProName  string `json:"pro_name" validate:"required"`
	UnitID1  string `json:"unit_id1" validate:"required"`
	UnitID2  string `json:"unit_id2"`
	UnitID3  string `json:"unit_id3"`
}

type ProductMappingImportRequest struct {
	URL      string `json:"url"`
	FileURL  string `json:"file_url"`
	Validate bool   `json:"validate"`
}

type ProductMappingImportResponse struct {
	URL           string   `json:"url,omitempty"`
	TotalRow      int      `json:"total_rows"`
	SuccessRow    int      `json:"success_rows"`
	FailedRow     int      `json:"failed_rows"`
	FailedReasons []string `json:"failed_reasons,omitempty"`
	ProcessedAt   string   `json:"processed_at,omitempty"`
}
