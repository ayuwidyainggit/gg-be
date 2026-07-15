package entity

type WorkingDayCalendarQueryFilter struct {
	CustID       string
	ParentCustID string
	Page         int    `query:"page" validate:"required,min=1"`
	Limit        int    `query:"limit" validate:"required,min=1,max=100"`
	Query        string `query:"q"`
	Sort         string `query:"sort"`
}

type CreateWorkingDayCalendarBody struct {
	CustID          string
	ParentCustID    string
	CreatedBy       int64
	Title           string `json:"title" validate:"required,max=100"`
	StartDate       string `json:"start_date" validate:"required"`
	NumberOfWeeks   int    `json:"number_of_weeks" validate:"required,min=1,max=99"`
	DefaultHolidays []int  `json:"default_holidays"`
}

type WorkingDayCalendarIDParams struct {
	WorkingDayCalendarID int64 `params:"working_day_calendar_id" validate:"required,min=1"`
}

type WorkingDayCalendarViewFilter struct {
	View  string `query:"view"`
	Month int    `query:"month"`
	Year  int    `query:"year"`
}

type WorkingDayCalendarImportHolidayRequest struct {
	FileURL string `json:"file_url" validate:"required,url"`
}

type WorkingDayCalendarImportHolidayResponse struct {
	FileURL       string   `json:"file_url"`
	FileName      string   `json:"file_name"`
	ProcessedAt   string   `json:"processed_at"`
	TotalRow      int      `json:"total_row"`
	SuccessRow    int      `json:"success_row"`
	FailedRow     int      `json:"failed_row"`
	FailedReasons []string `json:"failed_reasons"`
}

type WorkingDayCalendarListItem struct {
	WorkingDayCalendarID int64  `json:"working_day_calendar_id"`
	Title                string `json:"title"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	NumberOfWeeks        int    `json:"number_of_weeks"`
	DefaultHolidays      []int  `json:"default_holidays"`
}

type WorkingDayCalendarDetailResponse struct {
	WorkingDayCalendarID int64                       `json:"working_day_calendar_id"`
	Title                string                      `json:"title"`
	StartDate            string                      `json:"start_date"`
	EndDate              string                      `json:"end_date"`
	NumberOfWeeks        int                         `json:"number_of_weeks"`
	DefaultHolidays      []int                       `json:"default_holidays"`
	ImportedHolidays     []WorkingDayCalendarHoliday `json:"imported_holidays"`
}

type WorkingDayCalendarHoliday struct {
	DistributorCustID *string `json:"distributor_cust_id,omitempty"`
	Date              string  `json:"date"`
	Notes             string  `json:"notes"`
}

type WorkingDayCalendarDateItem struct {
	Date              string  `json:"date"`
	WeekID            int     `json:"week_id"`
	CalendarWeekNo    int     `json:"calendar_week_no"`
	WeekLabel         string  `json:"week_label"`
	IsWork            bool    `json:"is_work"`
	IsDefaultHoliday  bool    `json:"is_default_holiday"`
	IsImportedHoliday bool    `json:"is_imported_holiday"`
	Notes             *string `json:"notes"`
}

type WorkingDayCalendarViewResponse struct {
	WorkingDayCalendarID int64                        `json:"working_day_calendar_id"`
	Title                string                       `json:"title"`
	View                 string                       `json:"view"`
	Month                *int                         `json:"month,omitempty"`
	Year                 int                          `json:"year"`
	Dates                []WorkingDayCalendarDateItem `json:"dates"`
}
