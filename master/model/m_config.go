package model

import (
	"time"
)

type MConfig struct {
	ConfigId    string     `db:"config_id" json:"config_id"`
	ConfigValue string     `db:"config_value" json:"config_value"`
	DataType    string     `db:"data_type" json:"data_type"`
	ConfigDesc  string     `db:"config_desc" json:"config_desc"`
	Module      string     `db:"module" json:"module"`
	CreatedDate *time.Time `db:"created_date" json:"created_date"`
	ConfigGroup *string    `db:"config_group" json:"config_group"`
}
