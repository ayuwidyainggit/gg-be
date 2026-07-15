package entity

type GamificationQueryFilter struct {
	CustId       string
	ParentCustId string
	Status       []int  `query:"status"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type GamificationListResponse struct {
	GamificationNo   string  `json:"gamification_no"`
	GamificationDate *string `json:"gamification_date"`
	Title            *string `json:"title"`
	// Description      *string `json:"description"`
	StartDate              *string `json:"start_date"`
	EndDate                *string `json:"end_date"`
	FinishedDate           *string `json:"finished_date"`
	GamificationStatus     *int64  `json:"gamification_status"`
	GamificationStatusName *string `json:"gamification_status_name"`
	// CreatedBy        *int64  `json:"created_by"`
	// CreatedByName    *string `json:"created_by_name"`
	// CreatedAt        *string `json:"created_at"`
}

var dataGamificationStatusName = map[int64]string{
	1: "Not Started",
	2: "Active",
	3: "Finished",
	4: "Stopped",
	5: "Canceled",
}

func (gamification GamificationListResponse) GenerateGamificationStatusName() string {
	if gamification.GamificationStatus != nil {
		return dataGamificationStatusName[*gamification.GamificationStatus]
	}
	return ""
}

type GamificationStatusesLookupResponse struct {
	GamificationStatus     *int64  `json:"gamification_status"`
	GamificationStatusName *string `json:"gamification_status_name"`
}

func (gamificationStatus GamificationStatusesLookupResponse) GenerateDataGamificationStatusName() string {
	if gamificationStatus.GamificationStatus != nil {
		return dataGamificationStatusName[*gamificationStatus.GamificationStatus]
	}
	return ""
}

type CreateGamificationBody struct {
	CustID             string   `json:"cust_id"`
	GamificationNo     string   `json:"gamification_no"`
	GamificationDate   string   `json:"gamification_date"`
	Title              string   `json:"title" validate:"required"`
	Description        *string  `json:"description"`
	EmpGrpID           int64    `json:"emp_grp_id" validate:"required"`
	ProductFilter      string   `json:"product_filter"`
	OutletFilter       string   `json:"outlet_filter"`
	SourceId           int64    `json:"source_id"`
	StartDate          string   `json:"start_date"`
	EndDate            string   `json:"end_date"`
	MeasurementId      int64    `json:"measurement_id"`
	AnnouncementId     int64    `json:"announcement_id"`
	AnnouncementTarget float64  `json:"announcement_target"`
	SubAnnouncementId  int64    `json:"subannouncement_id"`
	CreatedBy          int64    `json:"created_by"`
	Participants       []int64  `json:"participants"`
	Products           []int64  `json:"products"`
	Outlets            []int64  `json:"outlets"`
	Announcements      []string `json:"announcements"`
	SalesTeams         []int64  `json:"sales_teams"`
	OperationTypes     []string `json:"operation_types"`
	Warehouses         []int64  `json:"warehouses"`
	ProductCategories  []int64  `json:"product_categories"`
	ProductLines       []int64  `json:"product_lines"`
	Brands             []int64  `json:"brands"`
	SubBrands1         []int64  `json:"sub_brands1"`
	SubBrands2         []int64  `json:"sub_brands2"`
	OutletLocations    []int64  `json:"outlet_locations"`
	OutletGroups       []int64  `json:"outlet_groups"`
	OutletClasses      []int64  `json:"outlet_classes"`
	Markets            []int64  `json:"markets"`
	Districts          []int64  `json:"districts"`
	Industries         []int64  `json:"industries"`
}

type DetailGamificationParams struct {
	GamificationNo string `params:"gamification_no" validate:"required"`
}

type GamificationResponse struct {
	CustID                 string                              `json:"cust_id"`
	GamificationNo         string                              `json:"gamification_no"`
	GamificationDate       string                              `json:"gamification_date"`
	Title                  string                              `json:"title"`
	Description            *string                             `json:"description"`
	EmpGrpID               int64                               `json:"emp_grp_id"`
	EmpGrpCode             *string                             `json:"emp_grp_code"`
	EmpGrpName             *string                             `json:"emp_grp_name"`
	ProductFilter          string                              `json:"product_filter"`
	OutletFilter           string                              `json:"outlet_filter"`
	SourceId               int64                               `json:"source_id"`
	SourceName             string                              `json:"source_name"`
	StartDate              string                              `json:"start_date"`
	EndDate                string                              `json:"end_date"`
	FinishedDate           *string                             `json:"finished_date"`
	MeasurementId          int64                               `json:"measurement_id"`
	MeasurementName        string                              `json:"measurement_name"`
	AnnouncementId         int64                               `json:"announcement_id"`
	AnnouncementName       string                              `json:"announcement_name"`
	AnnouncementTarget     float64                             `json:"announcement_target"`
	SubAnnouncementId      *int64                              `json:"subannouncement_id"`
	SubAnnouncementName    *string                             `json:"subannouncement_name"`
	GamificationStatus     int64                               `json:"gamification_status"`
	GamificationStatusName string                              `json:"gamification_status_name"`
	Participants           []GamificationParticipantsResponse  `json:"participants"`
	Products               []GamificationProductsResponse      `json:"products"`
	Outlets                []GamificationOutletsResponse       `json:"outlets"`
	Announcements          []GamificationAnnouncementsResponse `json:"announcements"`
	Rankings               []GamificationRankingsResponse      `json:"rankings"`
	SalesTeams             []int64                             `json:"sales_teams"`
	OperationTypes         []string                            `json:"operation_types"`
	Warehouses             []int64                             `json:"warehouses"`
	ProductCategories      []int64                             `json:"product_categories"`
	ProductLines           []int64                             `json:"product_lines"`
	Brands                 []int64                             `json:"brands"`
	SubBrands1             []int64                             `json:"sub_brands1"`
	SubBrands2             []int64                             `json:"sub_brands2"`
	OutletLocations        []int64                             `json:"outlet_locations"`
	OutletGroups           []int64                             `json:"outlet_groups"`
	OutletClasses          []int64                             `json:"outlet_classes"`
	Markets                []int64                             `json:"markets"`
	Districts              []int64                             `json:"districts"`
	Industries             []int64                             `json:"industries"`
}

type EmployeeGroupsLookupResponse struct {
	EmpGroupId   int    `json:"emp_grp_id"`
	EmpGroupCode string `json:"emp_grp_code"`
	EmpGroupName string `json:"emp_grp_name"`
}

type SalesTeamsLookupResponse struct {
	SalesTeamId   int    `json:"sales_team_id"`
	SalesTeamCode string `json:"sales_team_code"`
	SalesTeamName string `json:"sales_team_name"`
}

type OperationTypesLookupResponse struct {
	OperationTypeCode string `json:"operation_type_code"`
	OperationTypeName string `json:"operation_type_name"`
}

type SalesmanQueryFilter struct {
	CustId       string
	ParentCustId string
	EmpGrpId     []int    `query:"emp_grp_id"`
	SalesTeamId  []int    `query:"sales_team_id"`
	OprType      []string `query:"opr_type"`
	WhId         []int    `query:"wh_id"`
	From         *int64   `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64   `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int      `query:"page"`
	Limit        int      `query:"limit" validate:""`
	Query        string   `query:"q"`
	Mode         string   `query:"mode"`
	Sort         string   `query:"sort"`
	IsActive     *int     `query:"is_active"`
}

type ProductQueryFilter struct {
	CustId       string
	ParentCustId string
	PCatId       []int `query:"pcat_id"`
	// PLId         []int  `query:"pl_id"`
	// BrandId      []int  `query:"brand_id"`
	SBrand1Id []int  `query:"sbrand1_id"`
	SBrand2Id []int  `query:"sbrand2_id"`
	From      *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To        *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page      int    `query:"page"`
	Query     string `query:"q"`
	Mode      string `query:"mode"`
	Sort      string `query:"sort"`
	IsActive  *int   `query:"is_active"`
}

type ProductCategoriesLookupResponse struct {
	PCatId   int    `json:"pcat_id"`
	PCatCode string `json:"pcat_code"`
	PCatName string `json:"pcat_name"`
}

type ProductLinesLookupResponse struct {
	PLId   int    `json:"pl_id"`
	PLCode string `json:"pl_code"`
	PLName string `json:"pl_name"`
}

type BrandQueryFilter struct {
	CustId       string
	ParentCustId string
	PLId         []int  `query:"pl_id"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type BrandsLookupResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type SubBrand1QueryFilter struct {
	CustId       string
	ParentCustId string
	BrandId      []int  `query:"brand_id"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type SubBrands1LookupResponse struct {
	SBrand1Id   int    `json:"sbrand1_id"`
	SBrand1Code string `json:"sbrand1_code"`
	SBrand1Name string `json:"sbrand1_name"`
}

type SubBrands2LookupResponse struct {
	SBrand2Id   int    `json:"sbrand2_id"`
	SBrand2Code string `json:"sbrand2_code"`
	SBrand2Name string `json:"sbrand2_name"`
}

type ProductsLookupResponse struct {
	ProductID   *int64  `json:"pro_id"`
	ProductCode *string `json:"pro_code"`
	ProductName *string `json:"pro_name"`
}

type OutletQueryFilter struct {
	CustId       string
	ParentCustId string
	OtLocId      []int  `query:"ot_loc_id"`
	OtGrpId      []int  `query:"ot_grp_id"`
	OtClassId    []int  `query:"ot_class_id"`
	MarketId     []int  `query:"market_id"`
	DistrictId   []int  `query:"district_id"`
	IndustryId   []int  `query:"industry_id"`
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     *int   `query:"is_active"`
}

type OutletLocationsLookupResponse struct {
	OtLocId   int    `json:"ot_loc_id"`
	OtLocCode string `json:"ot_loc_code"`
	OtLocName string `json:"ot_loc_name"`
}

type OutletGroupsLookupResponse struct {
	OtGrpId   int    `json:"ot_grp_id"`
	OtGrpCode string `json:"ot_grp_code"`
	OtGrpName string `json:"ot_grp_name"`
}

type OutletClassesLookupResponse struct {
	OtClassId   int    `json:"ot_class_id"`
	OtClassCode string `json:"ot_class_code"`
	OtClassName string `json:"ot_class_name"`
}

type MarketsLookupResponse struct {
	MarketId   int    `json:"market_id"`
	MarketCode string `json:"market_code"`
	MarketName string `json:"market_name"`
}

type DistrictsLookupResponse struct {
	DistrictId   int    `json:"district_id"`
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

type IndustriesLookupResponse struct {
	IndustryId   int    `json:"industry_id"`
	IndustryCode string `json:"industry_code"`
	IndustryName string `json:"industry_name"`
}

var dataGamificationSourceName = map[int64]string{
	1: "Purchase Order",
	2: "Invoice",
}

var dataGamificationMeasurementName = map[int64]string{
	1: "Revenue",
	2: "Quantity",
	3: "Total Purchase Order",
	4: "Total Invoice",
}

var dataGamificationAnnouncementName = map[int64]string{
	1: "One Time",
	2: "Recurring",
}

var dataGamificationSubAnnouncementName = map[int64]string{
	1: "By Period",
	2: "By Target",
	3: "Weekly",
	4: "Biweekly",
	5: "Monthly",
	6: "Quarterly",
}

var dataGamificationFilterName = map[string]string{
	"sales_team": "Sales Team",
	"op_type":    "Operation Type",
	"warehouse":  "Warehouse",
	"pcat":       "Product Category",
	"pline":      "Product Line",
	"brand":      "Brand",
	"sbrand1":    "Sub Brand 1",
	"sbrand2":    "Sub Brand 2",
	"ot_loc":     "Outlet Location",
	"ot_grp":     "Outlet Group",
	"ot_class":   "Outlet Class",
	"market":     "Market",
	"district":   "District",
	"industry":   "Industry",
}

type SourcesLookupResponse struct {
	SourceId   int64  `json:"source_id"`
	SourceName string `json:"source_name"`
}

type MeasurementsLookupResponse struct {
	MeasurementId   int64  `json:"measurement_id"`
	MeasurementName string `json:"measurement_name"`
}

type MeasurementQueryFilter struct {
	CustId       string
	ParentCustId string
	SourceId     int `query:"source_id"`
	Page         int `query:"page"`
	Limit        int `query:"limit"`
}

type AnnouncementsLookupResponse struct {
	AnnouncementId   int64  `json:"announcement_id"`
	AnnouncementName string `json:"announcement_name"`
}

type SubAnnouncementsLookupResponse struct {
	SubAnnouncementId   int64  `json:"subannouncement_id"`
	SubAnnouncementName string `json:"subannouncement_name"`
}

func (source SourcesLookupResponse) GetDataSource() map[int64]string {
	return dataGamificationSourceName
}

func (measurement MeasurementsLookupResponse) GetDataMeasurement() map[int64]string {
	return dataGamificationMeasurementName
}

func (announcement AnnouncementsLookupResponse) GetDataAnnouncement() map[int64]string {
	return dataGamificationAnnouncementName
}

func (subAnnouncement SubAnnouncementsLookupResponse) GetDataSubAnnouncement() map[int64]string {
	return dataGamificationSubAnnouncementName
}

type SubAnnouncementQueryFilter struct {
	CustId         string
	ParentCustId   string
	AnnouncemnetId int `query:"announcement_id"`
	Page           int `query:"page"`
	Limit          int `query:"limit"`
}

func (gamification GamificationResponse) GenerateGamificationStatusName() string {
	return dataGamificationStatusName[gamification.GamificationStatus]
}

func (gamification GamificationResponse) GenerateGamificationSourceName() string {
	return dataGamificationSourceName[gamification.SourceId]
}

func (gamification GamificationResponse) GenerateGamificationMeasurementName() string {
	return dataGamificationMeasurementName[gamification.MeasurementId]
}

func (gamification GamificationResponse) GenerateGamificationAnnouncementName() string {
	return dataGamificationAnnouncementName[gamification.AnnouncementId]
}

func (gamification GamificationResponse) GenerateGamificationSubAnnouncementName() string {
	if gamification.SubAnnouncementId != nil {
		return dataGamificationSubAnnouncementName[*gamification.SubAnnouncementId]
	}
	return ""
}

type UpdateGamificationParams struct {
	GamificationNo string `params:"gamification_no" validate:"required"`
}

type UpdateGamificationBody struct {
	CustID             string   `json:"cust_id"`
	Title              string   `json:"title" validate:"required"`
	Description        *string  `json:"description"`
	EmpGrpID           int64    `json:"emp_grp_id" validate:"required"`
	ProductFilter      string   `json:"product_filter"`
	OutletFilter       string   `json:"outlet_filter"`
	SourceId           int64    `json:"source_id"`
	StartDate          string   `json:"start_date"`
	EndDate            string   `json:"end_date"`
	MeasurementId      int64    `json:"measurement_id"`
	AnnouncementId     int64    `json:"announcement_id"`
	AnnouncementTarget float64  `json:"announcement_target"`
	SubAnnouncementId  int64    `json:"subannouncement_id"`
	UpdatedBy          int64    `json:"updated_by"`
	Participants       []int64  `json:"participants"`
	Products           []int64  `json:"products"`
	Outlets            []int64  `json:"outlets"`
	Announcements      []string `json:"announcements"`
	SalesTeams         []int64  `json:"sales_teams"`
	OperationTypes     []string `json:"operation_types"`
	Warehouses         []int64  `json:"warehouses"`
	ProductCategories  []int64  `json:"product_categories"`
	ProductLines       []int64  `json:"product_lines"`
	Brands             []int64  `json:"brands"`
	SubBrands1         []int64  `json:"sub_brands1"`
	SubBrands2         []int64  `json:"sub_brands2"`
	OutletLocations    []int64  `json:"outlet_locations"`
	OutletGroups       []int64  `json:"outlet_groups"`
	OutletClasses      []int64  `json:"outlet_classes"`
	Markets            []int64  `json:"markets"`
	Districts          []int64  `json:"districts"`
	Industries         []int64  `json:"industries"`
}

type CancelGamificationParams struct {
	GamificationNo string `params:"gamification_no" validate:"required"`
}

type CancelGamificationBody struct {
	CustID    string `json:"cust_id"`
	UpdatedBy int64  `json:"updated_by"`
}

type StopGamificationParams struct {
	GamificationNo string `params:"gamification_no" validate:"required"`
}

type StopGamificationBody struct {
	CustID    string `json:"cust_id"`
	UpdatedBy int64  `json:"updated_by"`
}
