package model

import "time"

type Pjp struct {
	ID            int    `gorm:"type:int;primary_key" json:"id"`
	PjpCode       int    `gorm:"column:pjp_code;type:int;uniqueIndex;not null" json:"pjp_code"`
	OperationType string `gorm:"column:operation_type;type:varchar(125);not null" json:"operation_type"`
	TeamSalesMan  string `gorm:"column:team_salesman;type:varchar(125)" json:"team_salesman"`
	SalesManID    int    `gorm:"column:salesman_id" json:"salesman_id"`
	SalesmanName  string `gorm:"column:salesman_name;type:varchar(125)" json:"salesman_name"`
	WarehouseID   int    `gorm:"column:warehouse_id" json:"warehouse_id"`
	WarehouseName string `gorm:"column:warehouse_name;type:varchar(125)" json:"warehouse_name"`
	SalesmanCode  string `gorm:"column:salesman_code" json:"salesman_code"`
	PjpMode       string `gorm:"column:pjp_mode;type:varchar(125);null" json:"pjp_mode"`
	// EmpCode       string `gorm:"column:emp_code;type:varchar(125);null" json:"emp_code"`
	EmpCode           string    `gorm:"->" json:"emp_code"`
	Status            string    `gorm:"column:status;type:varchar(125)" json:"status"`
	ApprovalStatus    string    `gorm:"type:varchar(32);not null;default:'Draft'" validate:"required,oneof='In Review' Draft Approved 'Approved With Propose' Reject" json:"approval_status"`
	CustID            string    `gorm:"column:cust_id;type:varchar(125);null" json:"cust_id"`
	TotalOutlet       int       `gorm:"->" json:"total_outlet"`
	TotalDestinations int       `gorm:"->" json:"total_destinations"`
	TotalRoute        int       `gorm:"->" json:"total_route"`
	RouteCode         int       `gorm:"->" json:"route_code"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	RouteOutlets []RouteOutlet `gorm:"foreignKey:pjp_id;references:ID"`
}

type PjpWithEmpCode struct {
	Pjp
	EmployeeCode string `json:"emp_code"`
}

type Tabler interface {
	TableName() string
}

func (Pjp) TableName() string {
	return "pjp.permanent_journey_plans"
}
