package service

import (
	"database/sql"
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type DistributorReplenishmentSetupService interface {
	List(dataFilter entity.DistributorReplenishmentSetupQueryFilter) ([]entity.DistributorReplenishmentSetupListItem, int, int, error)
	ListSuppliersForPic(f entity.DistributorReplenishmentSupplierQueryFilter) ([]entity.DistributorReplenishmentSupplierItem, int, int, error)
	ListDistributorsForPic(f entity.DistributorReplenishmentDistributorQueryFilter) ([]entity.DistributorReplenishmentDistributorItem, int, int, error)
	ListByPicUserID(dataFilter entity.DistributorReplenishmentSetupQueryFilter, picUserID int) ([]entity.DistributorReplenishmentSetupListItem, int, int, error)
	Detail(id int, custID, parentCustID string) (entity.DistributorReplenishmentSetupDetailResponse, error)
	Create(payload entity.DistributorReplenishmentSetupCreatePayload, custID, parentCustID string, userID int64) (entity.DistributorReplenishmentSetupDetailResponse, error)
	Update(id int, payload entity.DistributorReplenishmentSetupCreatePayload, custID, parentCustID string, userID int64) (entity.DistributorReplenishmentSetupDetailResponse, error)
	Delete(id int, custID, parentCustID string, userID int64) error
}

func NewDistributorReplenishmentSetupService(repo repository.DistributorReplenishmentSetupRepository) DistributorReplenishmentSetupService {
	return &distributorReplenishmentSetupServiceImpl{repo: repo}
}

type distributorReplenishmentSetupServiceImpl struct {
	repo repository.DistributorReplenishmentSetupRepository
}

func (s *distributorReplenishmentSetupServiceImpl) ListSuppliersForPic(f entity.DistributorReplenishmentSupplierQueryFilter) ([]entity.DistributorReplenishmentSupplierItem, int, int, error) {
	rows, total, lastPage, err := s.repo.FindSuppliersForPic(f)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, ListSuppliersForPic, err:", err.Error())
		return nil, 0, 0, err
	}
	out := make([]entity.DistributorReplenishmentSupplierItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.DistributorReplenishmentSupplierItem{
			ID:      r.SupID,
			SupID:   r.SupID,
			SupCode: r.SupCode,
			SupName: r.SupName,
		})
	}
	return out, total, lastPage, nil
}

func (s *distributorReplenishmentSetupServiceImpl) ListDistributorsForPic(f entity.DistributorReplenishmentDistributorQueryFilter) ([]entity.DistributorReplenishmentDistributorItem, int, int, error) {
	rows, total, lastPage, err := s.repo.FindDistributorsForPic(f)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, ListDistributorsForPic, err:", err.Error())
		return nil, 0, 0, err
	}

	out := make([]entity.DistributorReplenishmentDistributorItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.DistributorReplenishmentDistributorItem{
			ID:              r.ID,
			DistributorID:   r.DistributorID,
			DistributorCode: r.DistributorCode,
			DistributorName: r.DistributorName,
		})
	}
	return out, total, lastPage, nil
}

func (s *distributorReplenishmentSetupServiceImpl) List(dataFilter entity.DistributorReplenishmentSetupQueryFilter) ([]entity.DistributorReplenishmentSetupListItem, int, int, error) {
	rows, total, lastPage, err := s.repo.FindAll(dataFilter)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, List, err:", err.Error())
		return nil, 0, 0, err
	}

	out := make([]entity.DistributorReplenishmentSetupListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapDistributorReplenishmentSetupToItem(r))
	}
	return out, total, lastPage, nil
}

func (s *distributorReplenishmentSetupServiceImpl) ListByPicUserID(dataFilter entity.DistributorReplenishmentSetupQueryFilter, picUserID int) ([]entity.DistributorReplenishmentSetupListItem, int, int, error) {
	rows, total, lastPage, err := s.repo.FindAllByPicUserID(dataFilter, picUserID)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, ListByPicUserID, err:", err.Error())
		return nil, 0, 0, err
	}

	out := make([]entity.DistributorReplenishmentSetupListItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapDistributorReplenishmentSetupToItem(r))
	}
	return out, total, lastPage, nil
}

func (s *distributorReplenishmentSetupServiceImpl) Detail(id int, custID, parentCustID string) (entity.DistributorReplenishmentSetupDetailResponse, error) {
	row, err := s.repo.FindByIDAndCustID(id, custID, parentCustID)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, Detail, FindByIDAndCustID:", err.Error())
		return entity.DistributorReplenishmentSetupDetailResponse{}, err
	}
	approvals, err := s.repo.FindApprovalsBySetupIDAndCustID(id, custID, parentCustID)
	if err != nil {
		log.Error("DistributorReplenishmentSetupService, Detail, FindApprovals:", err.Error())
		return entity.DistributorReplenishmentSetupDetailResponse{}, err
	}
	return mapToDetailResponse(row, approvals), nil
}

func (s *distributorReplenishmentSetupServiceImpl) Create(payload entity.DistributorReplenishmentSetupCreatePayload, custID, parentCustID string, userID int64) (entity.DistributorReplenishmentSetupDetailResponse, error) {
	if payload.IsApprovalRequired != nil && *payload.IsApprovalRequired && len(payload.ApprovalData) == 0 {
		return entity.DistributorReplenishmentSetupDetailResponse{}, errors.New("approval_data is required when is_approval_required is true")
	}

	setupID, err := s.repo.Create(custID, userID, payload)
	if err != nil {
		return entity.DistributorReplenishmentSetupDetailResponse{}, err
	}
	return s.Detail(setupID, custID, parentCustID)
}

func (s *distributorReplenishmentSetupServiceImpl) Update(id int, payload entity.DistributorReplenishmentSetupCreatePayload, custID, parentCustID string, userID int64) (entity.DistributorReplenishmentSetupDetailResponse, error) {
	if payload.IsApprovalRequired != nil && *payload.IsApprovalRequired && len(payload.ApprovalData) == 0 {
		return entity.DistributorReplenishmentSetupDetailResponse{}, errors.New("approval_data is required when is_approval_required is true")
	}

	if _, err := s.repo.FindByIDAndCustID(id, custID, parentCustID); err != nil {
		return entity.DistributorReplenishmentSetupDetailResponse{}, err
	}

	if err := s.repo.Update(id, custID, userID, payload); err != nil {
		return entity.DistributorReplenishmentSetupDetailResponse{}, err
	}
	return s.Detail(id, custID, parentCustID)
}

func (s *distributorReplenishmentSetupServiceImpl) Delete(id int, custID, parentCustID string, userID int64) error {
	if _, err := s.repo.FindByIDAndCustID(id, custID, parentCustID); err != nil {
		return err
	}
	return s.repo.Delete(id, custID, userID)
}

func mapToDetailResponse(r model.DistributorReplenishmentSetup, approvals []model.DistributorReplenishmentApproval) entity.DistributorReplenishmentSetupDetailResponse {
	li := mapDistributorReplenishmentSetupToItem(r)
	ad := make([]entity.DistributorReplenishmentApprovalItem, 0, len(approvals))
	for _, a := range approvals {
		item := entity.DistributorReplenishmentApprovalItem{
			ID:                       a.ID,
			DistReplenishmentSetupID: a.DistReplenishmentSetupID,
			Level:                    a.Level,
			Sequence:                 a.Sequence,
			BusinessUnit:             a.BusinessUnit,
			Pic:                      a.Pic,
			IsActive:                 a.IsActive,
		}
		if a.BusinessUnitName.Valid {
			item.BusinessUnitName = a.BusinessUnitName.String
		}
		if a.PicName.Valid {
			item.PicName = a.PicName.String
		}
		ad = append(ad, item)
	}
	return entity.DistributorReplenishmentSetupDetailResponse{
		ID:                 li.ID,
		SupID:              li.SupID,
		SupCode:            li.SupCode,
		SupName:            li.SupName,
		DistributorID:      li.DistributorID,
		DistributorCode:    li.DistributorCode,
		DistributorName:    li.DistributorName,
		DistributorType:    li.DistributorType,
		WhLimitAction:      li.WhLimitAction,
		WhCapacity:         li.WhCapacity,
		WhVolume:           li.WhVolume,
		CreditLimitAction:  li.CreditLimitAction,
		PlafonCredit:       li.PlafonCredit,
		LeadTimeDays:       li.LeadTimeDays,
		IsApprovalRequired: li.IsApprovalRequired,
		CreatedBy:          li.CreatedBy,
		CreatedByName:      li.CreatedByName,
		CreatedAt:          li.CreatedAt,
		UpdatedBy:          li.UpdatedBy,
		UpdatedByName:      li.UpdatedByName,
		UpdatedAt:          li.UpdatedAt,
		ApprovalData:       ad,
	}
}

func mapDistributorReplenishmentSetupToItem(r model.DistributorReplenishmentSetup) entity.DistributorReplenishmentSetupListItem {
	item := entity.DistributorReplenishmentSetupListItem{
		ID:                 r.ID,
		SupID:              r.SupID,
		SupCode:            r.SupCode,
		SupName:            r.SupName,
		DistributorID:      r.DistributorID,
		DistributorCode:    r.DistributorCode,
		DistributorName:    r.DistributorName,
		DistributorType:    r.DistributorType,
		WhLimitAction:      emptyStringToNilPtr(r.WhLimitAction),
		WhCapacity:         nullInt64ToPtr(r.WhCapacity),
		WhVolume:           nullInt64ToPtr(r.WhVolume),
		CreditLimitAction:  creditLimitActionLabel(r.CreditLimitAction),
		PlafonCredit:       nullInt64ToPtr(r.PlafonCredit),
		LeadTimeDays:       r.LeadTimeDays,
		IsApprovalRequired: r.IsApprovalRequired,
	}
	if r.CreatedBy.Valid {
		item.CreatedBy = r.CreatedBy.Int64
	}
	if r.CreatedAt.Valid {
		item.CreatedAt = r.CreatedAt.Time.Format(time.RFC3339)
	}
	if r.CreatedByName.Valid {
		item.CreatedByName = r.CreatedByName.String
	}
	if r.UpdatedBy.Valid {
		item.UpdatedBy = r.UpdatedBy.Int64
	}
	if r.UpdatedByName.Valid {
		item.UpdatedByName = r.UpdatedByName.String
	}
	if r.UpdatedAt.Valid {
		item.UpdatedAt = r.UpdatedAt.Time.Format(time.RFC3339)
	}
	return item
}

func creditLimitActionLabel(n int) string {
	switch n {
	case 1:
		return "Restricted"
	case 2:
		return "Unrestricted"
	default:
		return ""
	}
}

func emptyStringToNilPtr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func nullInt64ToPtr(n sql.NullInt64) *int {
	if !n.Valid {
		return nil
	}
	v := int(n.Int64)
	return &v
}
