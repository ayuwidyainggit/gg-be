package repository

import (
	"sales/model"
	"strings"

	"gorm.io/gorm"
)

type OpenAPIRepository interface {
	FindActiveConfigByClientID(clientID string) (model.OpenAPIConfig, error)
	IsCustomerWhitelisted(configID int64, custID string) (bool, error)
	FindActiveInboundEndpoint(configID int64, method, endpointURL string) (model.OpenAPIEndpoint, error)
}

type RepositoryOpenAPIImpl struct {
	*gorm.DB
}

func NewOpenAPIRepo(db *gorm.DB) *RepositoryOpenAPIImpl {
	return &RepositoryOpenAPIImpl{DB: db}
}

func (r *RepositoryOpenAPIImpl) FindActiveConfigByClientID(clientID string) (model.OpenAPIConfig, error) {
	var cfg model.OpenAPIConfig
	err := r.DB.
		Where("client_id = ? AND status = ?", strings.TrimSpace(clientID), model.OpenAPIStatusActive).
		First(&cfg).Error
	return cfg, err
}

func (r *RepositoryOpenAPIImpl) IsCustomerWhitelisted(configID int64, custID string) (bool, error) {
	var count int64
	err := r.DB.Model(&model.OpenAPIConfigCustomer{}).
		Where("open_api_config_id = ? AND cust_id = ? AND status = ?", configID, strings.TrimSpace(custID), model.OpenAPIStatusActive).
		Count(&count).Error
	return count > 0, err
}

func (r *RepositoryOpenAPIImpl) FindActiveInboundEndpoint(configID int64, method, endpointURL string) (model.OpenAPIEndpoint, error) {
	var ep model.OpenAPIEndpoint
	err := r.DB.
		Where(
			"open_api_config_id = ? AND UPPER(method) = ? AND endpoint_url = ? AND UPPER(api_type) = ? AND is_active = TRUE",
			configID,
			strings.ToUpper(strings.TrimSpace(method)),
			normalizeEndpointPath(endpointURL),
			model.OpenAPITypeInbound,
		).
		First(&ep).Error
	return ep, err
}

func normalizeEndpointPath(path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return p
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}
