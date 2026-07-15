package model

type MDay struct {
	DayId   *int64  `gorm:"column:day_id;primaryKey" json:"day_id"`
	DayName *string `gorm:"column:day_name" json:"day_name"`
	LangId  *string `gorm:"column:lang_id" json:"lang_id"`
}

func (MDay) TableName() string {
	return "sys.m_day"
}
