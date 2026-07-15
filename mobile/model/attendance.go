package model

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	CustID       string    `gorm:"column:cust_id" json:"cust_id"`
	AttendanceId *int64    `gorm:"column:attendance_id;primaryKey" json:"attendance_id"`
	EmpCode      *string   `gorm:"column:emp_code" json:"emp_code"`
	Latitude     *string   `gorm:"column:latitude" json:"latitude"`
	Longitude    *string   `gorm:"column:longitude" json:"longitude"`
	LeaveID      *int64    `gorm:"column:leave_id" json:"leave_id"`
	Type         *int      `gorm:"column:type" json:"type"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (m *Attendance) BeforeCreate(trx *gorm.DB) (err error) {
	return nil
}
func (Attendance) TableName() string {
	return "mobile.attendances"
}

type AttendanceRead struct {
	AttendanceId *int64    `gorm:"column:attendance_id;primaryKey" json:"attendance_id"`
	EmpCode      *string   `gorm:"column:emp_code" json:"emp_code"`
	Latitude     *string   `gorm:"column:latitude" json:"latitude"`
	Longitude    *string   `gorm:"column:longitude" json:"longitude"`
	Type         *int      `gorm:"column:type" json:"type"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (AttendanceRead) TableName() string {
	return "mobile.attendances"
}
