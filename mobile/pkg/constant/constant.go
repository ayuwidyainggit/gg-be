package constant

const (
	CUST_ID              = "Cust_id"
	STATUS_OK            = "OK"
	HEADER_ACCEPT_LANG   = "Accept-Language"
	STATUS_DB_NOT_FOUND  = "data not found"
	NOT_FOUND            = "not found"
	RECORD_NOT_FOUND     = "record not found"
	SUCCESSFULLY_UPDATED = "Successfully Updated"
	SUCCESSFULLY_ADDED   = "Successfully Added"
	SUCCESSFULLY_SAVED   = "Successfully Saved"
	YYYY_MM_DD           = "2006-01-02"
	AttrProduct          = "PRO"
	AttrOutletClass      = "OCL"
	AttrOutletType       = "OTY"
	AttrOutletGroup      = "OTG"
	AttrSalesType        = "STY"
	AttrSalesTeam        = "STE"
	PaymentTypeCash      = "Cash"
	PaymentTypeTransfer  = "Transfer"
	CHECKIN_AVAILABLE                      = "Check-in Available"
	CHECKIN_UNAVAILABLE                    = "Check-in Unavailable"
	CHECKIN_UNAVAILABLE_DESCRIPTION        = "Check-in cannot be completed as there is no scheduled route plan and/or available stock. Please reach out to your administrator for further assistance."
	CHECKIN_AVAILABLE_DESCRIPTION          = "Check-in can be completed as there is scheduled route plan and available stock."
	CHECKIN_UNAVAILABLE_NO_PLAN            = "Check-in cannot be completed as there is no available route plan. Please reach out to your administrator for further assistance."
	CHECKIN_UNAVAILABLE_NO_STOCK           = "Check-in cannot be completed as there is no available stock. Please reach out to your administrator for further assistance."
	CHECKIN_UNAVAILABLE_NO_PLAN_AND_STOCK  = "Check-in cannot be completed as there is no route plan and available stock. Please reach out to your administrator for further assistance."

	// Source constants (common for inserts/updates)
	SourceMobile = "MOBILE"
	SourceWeb    = "WEB"
	SourceOther  = "OTHER"
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
