package model

import "time"

const (
	OpenAPIStatusActive   = "A"
	OpenAPIStatusInactive = "I"

	OpenAPITypeInbound  = "INBOUND"
	OpenAPITypeOutbound = "OUTBOUND"
)

type OpenAPIConfig struct {
	ID                  int64      `gorm:"column:id;primaryKey" json:"id"`
	SystemIntegration   string     `gorm:"column:system_integration" json:"system_integration"`
	ClientID            *string    `gorm:"column:client_id" json:"client_id"`
	ClientSecret        *string    `gorm:"column:client_secret" json:"-"`
	Environment         string     `gorm:"column:environment" json:"environment"`
	BaseURL             *string    `gorm:"column:base_url" json:"base_url"`
	SignatureAlgorithm  *string    `gorm:"column:signature_algorithm" json:"signature_algorithm"`
	Status              string     `gorm:"column:status" json:"status"`
	CreatedBy           string     `gorm:"column:created_by" json:"created_by"`
	CreatedDate         time.Time  `gorm:"column:created_date" json:"created_date"`
	UpdatedBy           *string    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedDate         *time.Time `gorm:"column:updated_date" json:"updated_date"`
}

func (OpenAPIConfig) TableName() string { return "mst.open_api_config" }

type OpenAPIConfigIP struct {
	ID               int64  `gorm:"column:id;primaryKey" json:"id"`
	OpenAPIConfigID  int64  `gorm:"column:open_api_config_id" json:"open_api_config_id"`
	IPAddress        string `gorm:"column:ip_address" json:"ip_address"`
	Status           string `gorm:"column:status" json:"status"`
}

func (OpenAPIConfigIP) TableName() string { return "mst.open_api_config_ip" }

type OpenAPIConfigCustomer struct {
	ID              int64  `gorm:"column:id;primaryKey" json:"id"`
	OpenAPIConfigID int64  `gorm:"column:open_api_config_id" json:"open_api_config_id"`
	CustID          string `gorm:"column:cust_id" json:"cust_id"`
	Status          string `gorm:"column:status" json:"status"`
}

func (OpenAPIConfigCustomer) TableName() string { return "mst.open_api_config_customer" }

type OpenAPIEndpoint struct {
	ID              int64   `gorm:"column:id;primaryKey" json:"id"`
	OpenAPIConfigID int64   `gorm:"column:open_api_config_id" json:"open_api_config_id"`
	APICode         string  `gorm:"column:api_code" json:"api_code"`
	APIName         string  `gorm:"column:api_name" json:"api_name"`
	EndpointURL     string  `gorm:"column:endpoint_url" json:"endpoint_url"`
	Method          string  `gorm:"column:method" json:"method"`
	APIType         string  `gorm:"column:api_type" json:"api_type"`
	TimeoutSecond   *int    `gorm:"column:timeout_second" json:"timeout_second"`
	RetryCount      *int    `gorm:"column:retry_count" json:"retry_count"`
	IsActive        bool    `gorm:"column:is_active" json:"is_active"`
	Description     *string `gorm:"column:description" json:"description"`
}

func (OpenAPIEndpoint) TableName() string { return "mst.open_api_endpoint" }
