package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
)

type OutletConfigService interface {
	List(filter entity.OutletConfigListFilter, custID, parentCustID string) ([]entity.OutletConfigListResponse, int, int, error)
	Detail(outletConfigID int, custID, parentCustID string) (entity.OutletConfigDetailResponse, error)
	Create(body entity.CreateOutletConfigBody, custID string, createdBy int64) error
	Update(outletConfigID int, custID, parentCustID string, body entity.CreateOutletConfigBody, updatedBy int64) error
	Delete(outletConfigID int, custID, parentCustID string, updatedBy int64) error
}

func NewOutletConfigService(repo repository.OutletConfigRepository) *outletConfigServiceImpl {
	return &outletConfigServiceImpl{OutletConfigRepository: repo}
}

type outletConfigServiceImpl struct {
	OutletConfigRepository repository.OutletConfigRepository
}

func (s *outletConfigServiceImpl) visibleCustIDs(custID, parentCustID string) ([]string, error) {
	if custID == "" {
		return nil, nil
	}
	if parentCustID != "" && parentCustID != custID {
		return []string{custID, parentCustID}, nil
	}
	children, err := s.OutletConfigRepository.GetCustIdsByParentCustId(custID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, 1+len(children))
	ids = append(ids, custID)
	ids = append(ids, children...)
	return ids, nil
}

func (s *outletConfigServiceImpl) canEditOrDelete(callerCustID, parentCustID, configOwnerCustID string) bool {
	if configOwnerCustID == "" {
		return false
	}
	if parentCustID != "" && callerCustID != parentCustID {
		return configOwnerCustID == callerCustID
	}
	if parentCustID == "" || callerCustID == parentCustID {
		if configOwnerCustID == callerCustID {
			return true
		}
		children, _ := s.OutletConfigRepository.GetCustIdsByParentCustId(callerCustID)
		for _, c := range children {
			if c == configOwnerCustID {
				return true
			}
		}
	}
	return false
}

func (s *outletConfigServiceImpl) List(filter entity.OutletConfigListFilter, custID, parentCustID string) ([]entity.OutletConfigListResponse, int, int, error) {
	visible, err := s.visibleCustIDs(custID, parentCustID)
	if err != nil {
		return nil, 0, 0, err
	}
	if len(visible) == 0 {
		return []entity.OutletConfigListResponse{}, 0, 0, nil
	}
	rows, total, lastPage, err := s.OutletConfigRepository.FindAll(filter, visible)
	if err != nil {
		return nil, 0, 0, err
	}
	list := make([]entity.OutletConfigListResponse, 0, len(rows))
	for _, r := range rows {
		canEdit := s.canEditOrDelete(custID, parentCustID, r.CustId)
		list = append(list, mapOutletConfigToResponse(r, canEdit, canEdit))
	}
	return list, total, lastPage, nil
}

func mapOutletConfigToResponse(r model.OutletConfig, canEdit, canDelete bool) entity.OutletConfigListResponse {
	resp := entity.OutletConfigListResponse{
		CustId:             r.CustId,
		OutletConfigId:     r.OutletConfigId,
		VerificationStatus: r.VerificationStatus,
		RulesType:          r.RulesType,
		Status:             r.Status,
		IsActive:           r.Status == 1,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
		CanEdit:            canEdit,
		CanDelete:          canDelete,
	}
	if r.CreatedByName != nil {
		resp.CreatedBy = *r.CreatedByName
	}
	if r.UpdatedByName != nil {
		resp.UpdatedBy = *r.UpdatedByName
	}
	return resp
}

func (s *outletConfigServiceImpl) Detail(outletConfigID int, custID, parentCustID string) (entity.OutletConfigDetailResponse, error) {
	header, err := s.OutletConfigRepository.FindHeaderByID(outletConfigID)
	if err != nil {
		return entity.OutletConfigDetailResponse{}, err
	}
	visible, err := s.visibleCustIDs(custID, parentCustID)
	if err != nil {
		return entity.OutletConfigDetailResponse{}, err
	}
	visibleMap := make(map[string]bool)
	for _, c := range visible {
		visibleMap[c] = true
	}
	if !visibleMap[header.CustId] {
		return entity.OutletConfigDetailResponse{}, errors.New("sql: no rows in result set")
	}
	_, details, err := s.OutletConfigRepository.FindByOutletConfigID(outletConfigID, header.CustId)
	if err != nil {
		return entity.OutletConfigDetailResponse{}, err
	}
	canEdit := s.canEditOrDelete(custID, parentCustID, header.CustId)
	resp := entity.OutletConfigDetailResponse{
		CustId:             header.CustId,
		OutletConfigId:     header.OutletConfigId,
		VerificationStatus: header.VerificationStatus,
		RulesType:          header.RulesType,
		OutletStatus:       make([]entity.OutletConfigDetItem, 0, len(details)),
		CanEdit:            canEdit,
		CanDelete:          canEdit,
	}
	for _, d := range details {
		desc := "UNKNOWN"
		if d.StatusDescription != nil && *d.StatusDescription != "" {
			desc = *d.StatusDescription
		}
		resp.OutletStatus = append(resp.OutletStatus, entity.OutletConfigDetItem{
			OutletConfigDetId: d.OutletConfigDetId,
			Status:            d.Status,
			StatusDesc:        desc,
			ValidateTrx:       d.ValidateTrx,
			CountingPeriod:    d.CountingPeriod,
		})
	}
	return resp, nil
}

func (s *outletConfigServiceImpl) Create(body entity.CreateOutletConfigBody, custID string, createdBy int64) error {
	return s.OutletConfigRepository.Store(custID, body, createdBy)
}

func (s *outletConfigServiceImpl) Update(outletConfigID int, custID, parentCustID string, body entity.CreateOutletConfigBody, updatedBy int64) error {
	header, err := s.OutletConfigRepository.FindHeaderByID(outletConfigID)
	if err != nil {
		return err
	}
	if !s.canEditOrDelete(custID, parentCustID, header.CustId) {
		return errors.New("not allowed to edit this outlet config")
	}
	return s.OutletConfigRepository.Update(outletConfigID, header.CustId, body, updatedBy)
}

func (s *outletConfigServiceImpl) Delete(outletConfigID int, custID, parentCustID string, updatedBy int64) error {
	header, err := s.OutletConfigRepository.FindHeaderByID(outletConfigID)
	if err != nil {
		return err
	}
	if !s.canEditOrDelete(custID, parentCustID, header.CustId) {
		return errors.New("not allowed to delete this outlet config")
	}
	return s.OutletConfigRepository.Delete(outletConfigID, header.CustId, updatedBy)
}

type OutletConfigStatusService interface {
	List(filter entity.OutletConfigStatusListFilter) ([]entity.OutletConfigStatusListResponse, int, int, error)
}

func NewOutletConfigStatusService(repo repository.OutletConfigStatusRepository) *outletConfigStatusServiceImpl {
	return &outletConfigStatusServiceImpl{OutletConfigStatusRepository: repo}
}

type outletConfigStatusServiceImpl struct {
	OutletConfigStatusRepository repository.OutletConfigStatusRepository
}

func (s *outletConfigStatusServiceImpl) List(filter entity.OutletConfigStatusListFilter) ([]entity.OutletConfigStatusListResponse, int, int, error) {
	rows, total, lastPage, err := s.OutletConfigStatusRepository.FindAll(filter)
	if err != nil {
		return nil, 0, 0, err
	}
	list := make([]entity.OutletConfigStatusListResponse, 0, len(rows))
	for _, r := range rows {
		list = append(list, mapOutletConfigStatusToResponse(r))
	}
	return list, total, lastPage, nil
}

func mapOutletConfigStatusToResponse(r model.OutletConfigStatus) entity.OutletConfigStatusListResponse {
	return entity.OutletConfigStatusListResponse{
		OutletConfigStatusId: r.OutletConfigStatusId,
		StatusCode:           r.StatusCode,
		StatusDescription:    r.StatusDescription,
		IsTrx:                r.IsTrx,
		SortOrder:            r.SortOrder,
		IsActive:             r.IsActive,
	}
}
