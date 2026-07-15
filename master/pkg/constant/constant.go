package constant

const (
	HEADER_ACCEPT_LANG            = "Accept-Language"
	CUST_ID                       = "Cust_id"
	SUCCESSFULLY_ADDED            = "Successfully Added"
	SUCCESSFULLY_UPDATED          = "Successfully Updated"
	SUCCESSFULLY_CANCELLED        = "Successfully Cancelled"
	SUCCESSFULLY_SAVED            = "Successfully Saved"
	SUCCESSFULLY_PREVIEWED        = "Successfully Previewed"
	RMQ_DEFAULT_QUEUE_TYPE        = "quorum"
	RMQ_DEFAULT_EXCHANGE          = "events"
	RMQ_DEFAULT_DELAY_SUFFIX      = ".delay"
	RMQ_MANAGE_PRICE_CREATE_EVENT = "manage-price.events.create"
	RMQ_OUTLET_PRICE_START_EVENT  = "outlet-price.events.start"
	RMQ_OUTLET_PRICE_END_EVENT    = "outlet-price.events.end"
	DATE_LAYOUT_YYYY_MM_DD        = "2006-01-02"
	JOB_TYPE_ONE_TIME             = "one_time"
	JOB_TASK_HTTP_REQ             = "http_request"
	SP_PRICE_PUBLISH_UNPUBLISH    = "/v1/outlet-prices/scheduler/publish-unpublish"
	DIST_PRICE_PUBLISH_UNPUBLISH  = "/v1/distributor-prices/scheduler/publish-unpublish"
	SALESMAN_ISACTIVE             = "/v1/salesman/scheduler/isactive"
	SALESMAN_DEACTIVE             = "/v1/salesman/scheduler/deactive"
	NO_DATA                       = "No Data Found"
	SUCCESS_NO_DATA               = "Success"
	SUCCESS_NO_DATA_DISPLAYED     = "Data successfully displayed"

	SUCCESS_GET_DISTRIBUTOR_REPLENISHMENT_SETUP        = "Success get distributor replenishment setup"
	SUCCESS_GET_DETAIL_DISTRIBUTOR_REPLENISHMENT_SETUP = "Success get detail distributor replenishment setup"
	RECORD_NOT_FOUND                                   = "record not found"

	ProductMappingDuplicateUOMErrorMsg = "Duplicate unit of measure is not allowed within the same product"
)
