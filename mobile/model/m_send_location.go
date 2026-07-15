package model

type MSendLocation struct {
	ID       int64 `json:"id" gorm:"primary_key"`
	Duration int64 `json:"duration" gorm:"column:duration"`
	Distance int64 `json:"distance" gorm:"column:distance"`
}

func (MSendLocation) TableName() string {
	return "sys.m_send_location"
}
