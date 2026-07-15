package constant

const (
	HEADER_ACCEPT_LANG                        = "Accept-Language"
	YYYY_MM_DD                                = "2006-01-02"
	DD_MM_YYYY                                = "02-01-2006"
	DATE_FORMAT_DD_MM_YYYY                    = "02/01/2006"
	DEFAULT_PROMO_VALUE                       = 0.0
	AttrProduct                               = "PRO"
	AttrOutletClass                           = "OCL"
	AttrOutletType                            = "OTY"
	AttrOutletGroup                           = "OTG"
	AttrSalesType                             = "STY"
	AttrSalesTeam                             = "STE"
	RMQ_DEFAULT_QUEUE_TYPE                    = "quorum"
	RMQ_DEFAULT_EXCHANGE                      = "events"
	RMQ_DEFAULT_DELAY_SUFFIX                  = ".delay"
	RMQ_SECONDARY_SALES_EXPORT                = "secondary-sales.events.export"
	RMQ_SALESMAN_ACTIVITY_REPORT_SALES_EXPORT = "salesman-activity.events.export"
	MsgReturnCreateProductListSuccess        = "Successfully retrieved return create product list"
	MsgActivityReportSalesListSuccess        = "Successfully retrieved salesman activity report"
	MsgActivityReportSalesListFailed         = "Failed to retrieve salesman activity report"
	MsgActivityReportSalesExportSuccess      = "Successfully submitted salesman activity report export"
	MsgActivityReportSalesExportFailed       = "Failed to submit salesman activity report export"

	// Open API headers
	HeaderOpenAPIClientID     = "X-Client-Id"
	HeaderOpenAPIClientSecret = "X-Client-Secret"
	HeaderOpenAPICustID       = "X-Cust-Id"

	// Open API messages
	MsgOpenAPIMissingCredentials     = "Missing Open API credentials"
	MsgOpenAPIUnauthorized           = "Unauthorized"
	MsgOpenAPIForbidden              = "Forbidden"
	MsgOpenAPIEndpointNotAllowed     = "Endpoint not allowed"
	MsgOpenAPICreatePromotionSuccess = "Successfully added"
)

var uomName = map[int]string{
	0: "", 1: "Smallest", 2: "Middle", 3: "Largest",
}

func GetUomName(uomInt int) string {
	return uomName[uomInt]
}

var promoScopeLevelName = map[int]string{
	0: "", 1: "Distributor", 2: "Salesman", 3: "Outlet", 4: "Area",
}

func GetPromoScopeLevelName(scopeLevelInt int) string {
	return promoScopeLevelName[scopeLevelInt]
}

var includeExcludeDisplayName = map[string]string{
	"I": "Include", "E": "Exclude",
}

func GetIncludeExcludeDisplayName(key string) string {
	return includeExcludeDisplayName[key]
}

var qtyAmountPercentDisplayName = map[int]string{
	0: "", 1: "Quantity", 2: "Amount", 3: "Percentage",
}

func GetQtyAmountPercentDisplayName(key int) string {
	return qtyAmountPercentDisplayName[key]
}

var promoAttributeDisplayName = map[string]string{
	"PRO": "Product", "OCL": "Outlet Class",
	"OTY": "Outlet Type",
	"OTG": "Outlet Group", "STY": "Sales Type",
	"STE": "Sales Team",
}

func GetPromoAttributeDisplayName(key string) string {
	return promoAttributeDisplayName[key]
}
