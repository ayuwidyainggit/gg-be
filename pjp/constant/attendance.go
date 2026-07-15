package constant

// AttendanceMessage contains all response messages for attendance check API
var AttendanceMessage = struct {
	CheckInAvailable   string
	CheckInUnavailable string
}{
	CheckInAvailable:   "Check-in Available",
	CheckInUnavailable: "Check-in Unavailable",
}

// AttendanceDescription contains all description messages for attendance check API
var AttendanceDescription = struct {
	NoPlanPrincipal       string
	NoPlanDistributor     string
	NoStock               string
	NoPlanAndNoStock      string
	ErrorFetchingSalesman string
}{
	NoPlanPrincipal:       "Check-in cannot be completed as there is no scheduled route plan. Please reach out to your administrator for further assistance.",
	NoPlanDistributor:     "Check-in cannot be completed because there is no scheduled route plan. Please contact your administrator.",
	NoStock:               "Check-in cannot be completed because the stock is unavailable. Please contact your administrator.",
	NoPlanAndNoStock:      "Check-in cannot be completed because there is no scheduled route plan and the stock is unavailable. Please contact your administrator.",
	ErrorFetchingSalesman: "Error fetching salesman info",
}

// SalesmanOperationType defines operation types for salesman
const (
	OprTypeTakingOrder = "O" // Taking Order salesman
	OprTypeCanvas      = "C" // Canvas salesman
)

// DateFormat for database queries
const DateFormat = "2006-01-02"
