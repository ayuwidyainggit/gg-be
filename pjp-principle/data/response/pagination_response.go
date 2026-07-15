package response

import "gorm.io/gorm"

type Meta struct {
	Limit     int `json:"limit"`
	Page      int `json:"page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type Pagination struct {
	Code    int         `json:"code"`   
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Meta    Meta        `json:"meta"`
	TraceID string      `json:"trace_id"`
}

func Scopes(page int, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}


