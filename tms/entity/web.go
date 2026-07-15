package entity

import "gorm.io/gorm"

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	TraceID string      `json:"trace_id"`
}

type Error struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Errors  interface{} `json:"errors,omitempty"`
	TraceID string      `json:"trace_id"`
}

type Paging struct {
	TotalRecord int `json:"total_record"`
	PageCurrent int `json:"page_current"`
	PageLimit   int `json:"page_limit"`
	PageTotal   int `json:"page_total"`
}

type Meta struct {
	Limit     int `json:"limit"`
	Page      int `json:"page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func Scopes(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func QueryByDriverID(db *gorm.DB, driverId int) *gorm.DB {
	return db.Where("driver_id = ?", driverId)
}

func QueryByCustomerID(db *gorm.DB, customerId int) *gorm.DB {
	return db.Where("cust_id = ?", customerId)
}

type GeneralQueryFilter struct {
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:"required"`
	Query        string `query:"q"`
	Mode         string `query:"mode"`
	Sort         string `query:"sort"`
	IsActive     int    `query:"is_active"`
	OutletId     string `query:"outlet_id"`
	DocumentType string `query:"document_type"`
	EmpId        string `query:"emp_id"`
}
