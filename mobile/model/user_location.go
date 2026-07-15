package model

import "time"

type UserLocation struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	CustId    string    `json:"cust_id" gorm:"column:cust_id"`
	EmpID     int64     `json:"emp_id" gorm:"column:emp_id"`
	Latitude  string    `json:"latitude" gorm:"column:latitude"`
	Longitude string    `json:"longitude" gorm:"column:longitude"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (UserLocation) TableName() string {
	return "sys.user_location"
}
