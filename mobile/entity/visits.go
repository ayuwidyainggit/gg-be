package entity

import (
	"mime/multipart"
	"mobile/model"
)

var (
	TYPE_START_ID   int = 1
	TYPE_SKIP_ID    int = 5
	TYPE_ARRIVE_ID  int = 10
	TYPE_ON_HOLD_ID int = 15
	TYPE_RESUME_ID  int = 20
	TYPE_LEAVE_ID   int = 25
	TYPE_END_ID     int = 100

	ERROR_ARRIVE_MUST_BE_START                   = "do START before update type to ARRIVE"
	ERROR_ARRIVE_MUST_BE_FIRST_STATUS            = "ARRIVE must be first status. Last Status is %v"
	ERROR_ARRIVE_ALREADY_EXIST                   = "ARRIVE status already exist"
	ERROR_ARRIVE_DATE_MUST_GREAT_START           = "ARRIVE date must be greater than START time"
	ERROR_SKIP_MUST_BE_START                     = "do START before update type to SKIP"
	ERROR_SKIP_ALREADY_EXIST                     = "status SKIP on outlet code %v status already exist"
	ERROR_ON_HOLD_MUST_BE_ARRIVE                 = "do ARRIVE before update type to ON HOLD"
	ERROR_ON_HOLD_LAST_MUST_BE_ARRIVE            = "last status must be ARRIVE"
	ERROR_ON_HOLD_DATE_MUST_GREAT_ARRIVE         = "ON HOLD date must be greater than ARRIVE time"
	ERROR_RESUME_MUST_BE_ON_HOLD                 = "do ON HOLD before update type to RESUME"
	ERROR_RESUME_LAST_MUST_BE_ON_HOLD            = "last status must be ON HOLD"
	ERROR_RESUME_DATE_MUST_GREAT_ON_HOLD         = "RESUME date must be greater than ON HOLD time"
	ERROR_LEAVE_MUST_BE_ARRIVE_OR_RESUME         = "do ARRIVE or RESUME before update type to LEAVE"
	ERROR_LEAVE_LAST_MUST_BE_ARRIVE_OR_RESUME    = "last status must be ARRIVE or RESUME"
	ERROR_LEAVE_DATE_MUST_GREAT_ARRIVE_OR_RESUME = "LEAVE date must be greater than ARRIVE or RESUME time"
)

func GetTypeNameByID(typeID int) string {
	switch typeID {
	case TYPE_START_ID:
		return "START"

	case TYPE_SKIP_ID:
		return "SKIP"

	case TYPE_ARRIVE_ID:
		return "ARRIVE"

	case TYPE_ON_HOLD_ID:
		return "ON HOLD"

	case TYPE_END_ID:
		return "END"

	case TYPE_LEAVE_ID:
		return "LEAVE"

	case TYPE_LEAVE_ID:
		return "LEAVE"

	}
	return ""
}

type VisitsRequest struct {
	CurrentTime int64 `json:"current_time" validate:"required"`
	Email       string
	CustID      string
}

type VisitsListRequest struct {
	CustID        string
	OutletID      *int64  `json:"outlet_id" query:"outlet_id"`
	EmpCode       *string `json:"-"`
	EmpID         int64   `json:"-"`
	IsDistributor bool    `json:"-"`
	PJPID         int64   `json:"-" query:"-"`
}

type VisitsSummariesRequest struct {
	Email  string
	CustID string
}

type VisitsResponse struct {
	Sequence   int    `json:"sequence"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
	OutletImg  string `json:"outlet_img"`
	Status     string `json:"status"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	ArrivedAt  string `json:"arrived_at"`
	Barcode    string `json:"barcode"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	ZipCode    string `json:"zip_code"`
	PhoneNo    string `json:"phone_no"`
	WaNo       string `json:"wa_no"`
	Email      string `json:"email"`
}

type SummariesResponse struct {
	CurrentTime     string  `json:"current_time"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	Plan            float64 `json:"plan"`
	Finished        float64 `json:"finished"`
	Skipped         float64 `json:"skipped"`
	InProgress      float64 `json:"in_progress"`
	OnHold          float64 `json:"on_hold"`
	CurrentProgress float64 `json:"current_progress"`

	Totals            TotalSummaryDeposit        `json:"totals"`
	InvoiceList       []model.InvoiceListItem    `json:"invoice_list"`
	ExpenseList       []model.ExpenseListItem    `json:"expense_list"`
	CollectionSummary []model.CollectionListItem `json:"collections"`
}

func (s *SummariesResponse) PlanIncrement() {
	s.Plan++
}
func (s *SummariesResponse) FinishIncrement() {
	s.Finished++
}
func (s *SummariesResponse) SkippedIncrement() {
	s.Skipped++
}
func (s *SummariesResponse) InProgressIncrement() {
	s.InProgress++
}
func (s *SummariesResponse) OnHoldIncrement() {
	s.OnHold++
}

type VisitListResponse struct {
	Sequence   int    `json:"sequence"`
	OutletCode string `json:"outlet_code"`
	OutletName string `json:"outlet_name"`
	OutletImg  string `json:"outlet_img"`
	Status     string `json:"status"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	Barcode    string `json:"barcode"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	ZipCode    string `json:"zip_code"`
	PhoneNo    string `json:"phone_no"`
	WaNo       string `json:"wa_no"`
	Email      string `json:"email"`
}

type SkipReasonsResponse struct {
	// IsInOutlet bool   `json:"is_in_outlet"`
	Reason string `json:"reason"`
}

type StartRequest struct {
	CurrentTime int64 `json:"current_time"`
	Email       string
	CustID      string
}

type SkipRequest struct {
	CurrentTime   int64  `json:"current_time"`
	OutletCode    string `json:"outlet_code"`
	Latitude      string `json:"latitude" validate:"required"`
	Longitude     string `json:"longitude" validate:"required"`
	Reason        string `json:"reason" validate:"required"`
	IsInOutlet    bool   `json:"is_in_outlet"`
	UpcomingVisit string `json:"upcoming_visit"`
	Email         string
	CustID        string
}

type ArriveRequest struct {
	CurrentTime      int64                 `json:"current_time" form:"current_time" validate:"required"`
	OutletCode       string                `json:"outlet_code" form:"outlet_code" validate:"required"`
	Latitude         string                `json:"latitude" form:"latitude" validate:"required"`
	Longitude        string                `json:"longitude" form:"longitude" validate:"required"`
	Folder           string                `json:"folder" form:"folder" validate:"required"`
	IsUpdateLocation bool                  `json:"is_update_location" form:"is_update_location"`
	File             *multipart.FileHeader `json:"-"`
	Email            string
	CustID           string
}

type ArriveResponse struct {
	PhotoURL          string `json:"photo_url"`
	IsUpdateLocation  bool   `json:"is_update_location"`
	OutletVisitListID int64  `json:"outlet_visit_list_id"`
}

type HoldRequest struct {
	CurrentTime int64  `json:"current_time"`
	OutletCode  string `json:"outlet_code"`
	Email       string
	CustID      string
}

type ResumeRequest struct {
	CurrentTime int64  `json:"current_time"`
	OutletCode  string `json:"outlet_code"`
	Email       string
	CustID      string
}

type LeaveRequest struct {
	CurrentTime int64  `json:"current_time"`
	OutletCode  string `json:"outlet_code"`
	Email       string
	CustID      string
}

type EndRequest struct {
	CurrentTime int64 `json:"current_time"`
	Email       string
	CustID      string
}

type VisitQueryFilter struct {
	CustId       string
	ParentCustId string
	Sort         string `query:"sort"`
}

type SkipReasonsQueryFilter struct {
	CustId       string
	ParentCustId string
	Sort         string `query:"sort"`
}

type VisitListDetailResponse struct {
	CheckinAt  string    `json:"checkin_at"`
	OutletCode string    `json:"outlet_code"`
	OutletName string    `json:"outlet_name"`
	Longitude  string    `json:"longitude"`
	Latitude   string    `json:"latitude"`
	File       *FileInfo `json:"file"`
}

type FileInfo struct {
	FileName      string `json:"file_name"`
	FileURL       string `json:"file_url"`
	FileKey       string `json:"file_key"`
	MediaCategory string `json:"media_category"`
	FileSize      int64  `json:"file_size"`
}

type SummariesRequest struct {
	EmpID  int64
	CustID string
	Date   string `query:"date"`
}
