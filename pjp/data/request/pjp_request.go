package request

type PjpRequest struct {
	ID            int    `json:"id,omitempty" validate:"omitempty,gte=1"`
	PjpCode       int    `validate:"required,min=1,max=9999" json:"pjp_code"`
	TeamSalesMan  string `validate:"required" json:"salesman_team"`
	OperationType string `validate:"required" json:"operation_type"`
	SalesManID    int    `validate:"gte=0" json:"salesman_id"`
	SalesmanName  string `validate:"required" json:"salesman_name"`
	SalesmanCode  string `validate:"required" json:"salesman_code"`
	WarehouseID   int    `validate:"gte=0" json:"warehouse_id"`
	WarehouseName string `validate:"required" json:"warehouse_name"`
	Status        string `validate:"required,oneof=true false" json:"is_active"`
}

type CreatePjpEnhanceRequest struct {
	PjpCode       int    `validate:"required,min=1,max=9999" json:"pjp_code"`
	TeamSalesMan  string `validate:"required"   json:"salesman_team" example:"Canvas"`
	OperationType string `validate:"required"   json:"operation_type" example:"GT"`
	SalesManID    int    `validate:"gte=0"   json:"salesman_id" example:"4"`
	SalesmanName  string `validate:"required"   json:"salesman_name" example:"Budi"`
	SalesmanCode  string `validate:"required"   json:"salesman_code" example:"3456"`
	WarehouseID   int    `validate:"gte=0"   json:"warehouse_id" example:"1"`
	WarehouseName string `validate:"required"   json:"warehouse_name" example:"Gudang"`
	// Status         string              `validate:"required,oneof=true false"   json:"is_active"`
	ApprovalStatus string              `validate:"required,oneof='Need Review' Draft" json:"approval_status" example:"Need Review"`
	Routes         []RoutesCreatePjp   `validate:"required,dive" json:"routes"`
	VisitDay       []VisitDayCreatePjp `json:"visit_day"`
}

type UpdateStatusPjpEnhanceRequest struct {
	Status        *bool   `json:"status" example:"true"`
	SalesmanName  *string `json:"salesman_name" example:"Charles Leclerc"`
	SalesmanCOde  *string `json:"salesman_code" example:"1234"`
	WarehouseID   *int    `json:"warehouse_id" example:"1"`
	WarehouseName *string `json:"warehouse_name" example:"Gudang"`
	OperationType *string `json:"operation_type" example:"Taking Order, Canvas"`
	TeamSalesMan  *string `json:"salesman_team" example:"MT"`
}

type RoutesCreatePjp struct {
	ID          *int          `json:"id,omitempty" example:"2"`           // digunakan untuk update/delete
	RouteID     *int          `json:"route_id,omitempty" example:"2"`     // alias from existing clients
	RouteCode   *int          `json:"route_code,omitempty" example:"234"` // optional untuk referensi
	RouteName   string        `validate:"required" json:"route_name" example:"Rute senin"`
	Destination []Destination `validate:"required,dive" json:"destinations,omitempty"`
}

type Destination struct {
	ID           int     `validate:"required" json:"id" example:"5"`
	Code         string  `validate:"required" json:"code" example:"OUT 5"`
	Name         string  `validate:"required" json:"name" example:"Outlet okok"`
	Longitude    string  `validate:"required" json:"longitude" example:"106.816666"`
	Latitude     string  `validate:"required" json:"latitude" example:"-6.200000"`
	Status       string  `json:"status" example:"1"`
	Address      string  `validate:"required" json:"address" example:"mampang"`
	AvgSalesWeek float64 `json:"avg_sales_week" example:"100"`
	Type         string  `validate:"required,oneof=outlet distributor" json:"type" example:"outlet"`
	IsAdditional bool    `json:"is_additional" example:"false"`
}

type VisitDayCreatePjp struct {
	ID                   int             `validate:"required" json:"id" example:"5"`
	Day                  string          `validate:"required" json:"day" example:"Senin"`
	IndexDay             int             `validate:"required" json:"indexDay" example:"1"`
	Week                 int             `validate:"required" json:"week" example:"1"`
	WorkingDayCalendarID *int64          `json:"working_day_calendar_id" example:"77"`
	StartWeek            string          `validate:"required,date" json:"startWeek" example:"2022-01-01"`
	Year                 int             `validate:"required" json:"year" example:"2022"`
	Date                 string          `validate:"required,date" json:"date" example:"2022-01-01"`
	IsInCurrentYear      bool            `validate:"required" json:"isInCurrentYear" example:"true"`
	Visit                RoutesCreatePjp `validate:"required" json:"visit"`
}
