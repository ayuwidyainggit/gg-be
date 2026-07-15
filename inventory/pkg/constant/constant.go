package constant

import "time"

const (
	HEADER_ACCEPT_LANG = "Accept-Language"

	STOCK_ADJUSTMENT_DET_TYPE_ADD_STOCK    = 1
	STOCK_ADJUSTMENT_DET_TYPE_REMOVE_STOCK = 2

	STOCK_ADJUTMENT_STATUS_APPROVED = 2
	STOCK_ADJUTMENT_STATUS_REJECT   = 9

	DEFAULT_PAGE_LIMIT     = 10
	MAX_FILE_SIZE_BYTES    = 10 * 1024 * 1024
	TR_CODE_STOCK_DISPOSAL = "SD"
	DATE_FORMAT_DISPLAY    = "2006-01-02"
	DATE_FORMAT_DETAIL     = "02/01/2006"

	// Response Messages
	DATA_NOT_FOUND              = "Data not found"
	RECORD_NOT_FOUND            = "Record not found"
	SUCCESS                     = "Success"
	DATA_SAVED_SUCCESSFULLY     = "Data successfully saved"
	DATA_DISPLAYED_SUCCESSFULLY = "Data successfully displayed"
	DATA_UPDATED_SUCCESSFULLY   = "Data successfully updated"
	DATA_CREATED_SUCCESSFULLY   = "Created Successfully"
	NO_DATA                     = "No Data"
	ERR_VALIDATION              = "Validation error"
	INVALID_JSON_BODY           = "invalid JSON body"

	// Stock Opname Messages
	STOCK_OPNAME_TEMPLATE_DOWNLOAD_SUCCESS = "Template downloaded successfully."
	STOCK_OPNAME_TEMPLATE_DOWNLOAD_FAILED  = "Failed to download template."
	STOCK_OPNAME_FILTER_REQUIRED           = "At least one of principal_id, pl_lane, brand_id, or sbrand1_id must be provided"
	STOCK_OPNAME_FILTER_REQUIRED_V2        = "At least one of principal_id, pl_id, brand_id, or sbrand1_id must be provided"
	DATA_STARTED_SUCCESSFULLY              = "Data successfully started"

	// Stock Disposal
	STOCK_DISPOSAL_PRODUCT_LOOKUP_SUCCESS = "Products loaded successfully."
	STOCK_DISPOSAL_PRODUCT_LOOKUP_FAILED  = "Failed to load products."
)

var AsiaJakartaLocation *time.Location

func init() {
	AsiaJakartaLocation, _ = time.LoadLocation("Asia/Jakarta")
}

var MapTransactionType = map[string]string{
	"SA": "Stock Adjustment",
	"SD": "Stock Disposal",
	"ST": "Stock Transfer",
}
