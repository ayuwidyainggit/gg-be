package entity

type ProductRipeningQueryFilter struct {
	CustId           string
	ParentCustId     string
	JwtDistributorId int64
	Page             int     `query:"page" validate:"required,min=1"`
	Limit            int     `query:"limit" validate:"required,min=1,max=100"`
	Query            string  `query:"q"`
	Sort             string  `query:"sort"`
	Format           string  `query:"format"`
	DistributorID    []int64 `query:"distributor_id"`
	WeekID           *int    `query:"week_id"`
	PerYear          *int    `query:"per_year"`
	PerID            *int    `query:"per_id"`
}

type ProductRipeningDetailParams struct {
	DistributorID int64 `params:"distributor_id" validate:"required,min=1"`
	PerYear       int   `params:"per_year" validate:"required,min=1"`
	WeekID        int   `params:"week_id" validate:"required,min=1"`
}

type ProductRipeningImportRequest struct {
	FileURL string `json:"file_url" validate:"required,url"`
}

type ProductRipeningImportResponse struct {
	FileURL       string   `json:"file_url"`
	FileName      string   `json:"file_name"`
	ProcessedAt   string   `json:"processed_at"`
	TotalRow      int      `json:"total_row"`
	SuccessRow    int      `json:"success_row"`
	FailedRow     int      `json:"failed_row"`
	FailedReasons []string `json:"failed_reasons"`
}

type ProductRipeningListItem struct {
	ID              int64   `json:"id"`
	DistributorID   int64   `json:"distributor_id"`
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	PerYear         int     `json:"per_year"`
	PerID           int     `json:"per_id"`
	WeekID          int     `json:"week_id"`
	WeekStart       string  `json:"week_start"`
	WeekEnd         string  `json:"week_end"`
	WeekLabel       string  `json:"week_label"`
	IsActive        bool    `json:"is_active"`
	CanEdit         bool    `json:"can_edit"`
	TotalProduct    int     `json:"total_product"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedBy       *int64  `json:"updated_by"`
	UpdatedByName   *string `json:"updated_by_name"`
	UpdatedAt       *string `json:"updated_at"`
}

type ProductRipeningEditableDays struct {
	Sunday    bool `json:"sunday"`
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
}

type ProductRipeningDetailResponse struct {
	DistributorID   int64                       `json:"distributor_id"`
	DistributorCode string                      `json:"distributor_code"`
	DistributorName string                      `json:"distributor_name"`
	PerYear         int                         `json:"per_year"`
	PerID           int                         `json:"per_id"`
	WeekID          int                         `json:"week_id"`
	WeekStart       string                      `json:"week_start"`
	WeekEnd         string                      `json:"week_end"`
	WeekLabel       string                      `json:"week_label"`
	IsActive        bool                        `json:"is_active"`
	CanEdit         bool                        `json:"can_edit"`
	EditableDays    ProductRipeningEditableDays `json:"editable_days"`
	Rows            []ProductRipeningDetailRow  `json:"rows"`
}

type ProductRipeningDetailRow struct {
	ID            int64   `json:"id"`
	ProID         int64   `json:"pro_id"`
	ProductCode   string  `json:"product_code"`
	ProductName   string  `json:"product_name"`
	SundayQty     int     `json:"sunday_qty"`
	MondayQty     int     `json:"monday_qty"`
	TuesdayQty    int     `json:"tuesday_qty"`
	WednesdayQty  int     `json:"wednesday_qty"`
	ThursdayQty   int     `json:"thursday_qty"`
	FridayQty     int     `json:"friday_qty"`
	SaturdayQty   int     `json:"saturday_qty"`
	CreatedBy     int64   `json:"created_by"`
	CreatedByName string  `json:"created_by_name"`
	CreatedAt     string  `json:"created_at"`
	UpdatedBy     *int64  `json:"updated_by"`
	UpdatedByName *string `json:"updated_by_name"`
	UpdatedAt     *string `json:"updated_at"`
}

type ProductRipeningUpdateRow struct {
	ID           int64 `json:"id" validate:"required,min=1"`
	SundayQty    int   `json:"sunday_qty" validate:"min=0"`
	MondayQty    int   `json:"monday_qty" validate:"min=0"`
	TuesdayQty   int   `json:"tuesday_qty" validate:"min=0"`
	WednesdayQty int   `json:"wednesday_qty" validate:"min=0"`
	ThursdayQty  int   `json:"thursday_qty" validate:"min=0"`
	FridayQty    int   `json:"friday_qty" validate:"min=0"`
	SaturdayQty  int   `json:"saturday_qty" validate:"min=0"`
}

type ProductRipeningUpdateRequest struct {
	Rows []ProductRipeningUpdateRow `json:"rows" validate:"required,min=1"`
}

type ProductRipeningExportRow struct {
	ID              int64   `json:"id"`
	DistributorCode string  `json:"distributor_code"`
	DistributorName string  `json:"distributor_name"`
	PerYear         int     `json:"per_year"`
	PerID           int     `json:"per_id"`
	WeekID          int     `json:"week_id"`
	WeekStart       string  `json:"week_start"`
	WeekEnd         string  `json:"week_end"`
	TotalProduct    int     `json:"total_product"`
	CreatedBy       int64   `json:"created_by"`
	CreatedByName   string  `json:"created_by_name"`
	CreatedAt       string  `json:"created_at"`
	UpdatedBy       *int64  `json:"updated_by"`
	UpdatedByName   *string `json:"updated_by_name"`
	UpdatedAt       *string `json:"updated_at"`
}
