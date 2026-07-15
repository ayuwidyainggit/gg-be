package entity

type MConfigQueryFilter struct {
	Page     int      `query:"page"`
	Limit    int      `query:"limit" validate:"required"`
	Query    string   `query:"q"`
	Mode     string   `query:"mode"`
	Sort     string   `query:"sort"`
	ConfigId []string `query:"config_id"`
}
type CreateMConfigBody struct {
	CustID      string  `json:"cust_id"`
	ConfigID    string  `json:"config_id"`
	ConfigValue *string `json:"config_value"`
	DataType    *string `json:"data_type"`
	ConfigDesc  *string `json:"config_desc"`
	Module      *string `json:"module" validate:"required,oneof='sys' 'inv' 'acf' 'mst' 'inv'"`
}

type UpdateMConfigBody struct {
	CustId      string  `json:"cust_id"`
	ConfigValue *string `json:"config_value"`
	DataType    *string `json:"data_type"`
	ConfigDesc  *string `json:"config_desc"`
	Module      *string `json:"module" validate:"required,oneof='sys' 'inv' 'acf' 'mst' 'inv'"`
}

type UpdateMConfigBodyParam struct {
	ConfigID string `params:"config_id" validate:"required"`
}

type MConfigResponse struct {
	ConfigID    string  `json:"config_id"`
	ConfigValue *string `json:"config_value"`
	DataType    *string `json:"data_type"`
	ConfigDesc  *string `json:"config_desc"`
	Module      *string `json:"module" validate:"required,oneof='sys' 'inv' 'acf' 'mst' 'inv'"`
}
type DetailConfigBodyParam struct {
	ConfigID string `params:"config_id" validate:"required"`
}
type DeleteConfigBodyParam struct {
	ConfigID string `params:"config_id" validate:"required"`
}
