package service

import (
	"errors"
	"strings"

	"sales/model"
	"sales/repository"
)

var (
	ErrOpenAPIUnauthorized     = errors.New("open api unauthorized")
	ErrOpenAPIForbidden        = errors.New("open api forbidden")
	ErrOpenAPIEndpointNotFound = errors.New("open api endpoint not found")
)

type OpenAPIAuthContext struct {
	ConfigID          int64
	SystemIntegration string
	CustID            string
	APICode           string
}

type OpenAPIService interface {
	Authenticate(clientID, clientSecret, method, path, custID string) (OpenAPIAuthContext, error)
}

type openAPIServiceImpl struct {
	repo repository.OpenAPIRepository
}

func NewOpenAPIService(repo repository.OpenAPIRepository) OpenAPIService {
	return &openAPIServiceImpl{repo: repo}
}

func (s *openAPIServiceImpl) Authenticate(clientID, clientSecret, method, path, custID string) (OpenAPIAuthContext, error) {
	cfg, err := s.repo.FindActiveConfigByClientID(clientID)
	if err != nil {
		return OpenAPIAuthContext{}, ErrOpenAPIUnauthorized
	}

	if !validateClientSecret(cfg, clientSecret) {
		return OpenAPIAuthContext{}, ErrOpenAPIUnauthorized
	}

	custID = strings.TrimSpace(custID)
	if custID == "" {
		return OpenAPIAuthContext{}, ErrOpenAPIForbidden
	}

	ok, err := s.repo.IsCustomerWhitelisted(cfg.ID, custID)
	if err != nil || !ok {
		return OpenAPIAuthContext{}, ErrOpenAPIForbidden
	}

	ep, err := s.repo.FindActiveInboundEndpoint(cfg.ID, method, path)
	if err != nil {
		return OpenAPIAuthContext{}, ErrOpenAPIEndpointNotFound
	}

	return OpenAPIAuthContext{
		ConfigID:          cfg.ID,
		SystemIntegration: cfg.SystemIntegration,
		CustID:            custID,
		APICode:           ep.APICode,
	}, nil
}

func validateClientSecret(cfg model.OpenAPIConfig, secret string) bool {
	if cfg.ClientSecret == nil {
		return false
	}
	return strings.TrimSpace(*cfg.ClientSecret) == strings.TrimSpace(secret)
}
