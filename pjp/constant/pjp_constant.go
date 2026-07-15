package constant

// Error messages
const (
	ErrPjpCodeExists = "PJP Code already exists"

	// OutletChangeRequestSourceMobile represents mobile source for outlet change request
	OutletChangeRequestSourceMobile = 2
	// OutletChangeRequestStatusPending represents pending status for outlet change request
	OutletChangeRequestStatusPending = 1
	ErrCustomerIDNotFound            = "customer_id not found"
	ErrInvalidRequestParams          = "Invalid request parameters: "
)

// Response messages
const (
	MsgSuccess      = "Success"
	MsgNoData       = "No Data"
	MsgUnauthorized = "Unauthorized"
	MsgOK           = "OK"
)

// User levels
const (
	LevelPrincipal   = "Principal"
	LevelDistributor = "Distributor"
)

// Approval status
const (
	ApprovalStatusApproved = "Approved"
)

// Default values
const (
	DefaultCompanyCode = "-"
)
