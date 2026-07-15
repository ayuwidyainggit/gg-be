package service

import (
	"errors"
	"master/entity"
	"master/model"
	"master/repository"
	"strings"
	"time"
)

type OutletCodeService interface {
	List(filter entity.OutletCodeListFilter, custId string) ([]entity.OutletCodeItem, int, int, error)
	Create(body entity.CreateOutletCodeBody, custId, createdBy string) error
	Update(id string, body entity.UpdateOutletCodeBody, custId, updatedBy string) error
	UpdateStatus(id string, body entity.UpdateOutletCodeStatusBody, custId, updatedBy string) error
	SetupOutletCheck(custId, parentCustId string, year int, status []string, createdBy string) (*entity.SetupOutletCheckData, error)
}

func NewOutletCodeService(repo repository.OutletCodeRepository) OutletCodeService {
	return &outletCodeServiceImpl{OutletCodeRepository: repo}
}

type outletCodeServiceImpl struct {
	OutletCodeRepository repository.OutletCodeRepository
}

func (s *outletCodeServiceImpl) List(filter entity.OutletCodeListFilter, custId string) ([]entity.OutletCodeItem, int, int, error) {
	custIDs := []string{}
	if custId != "" {
		custIDs = append(custIDs, custId)
	}
	rows, total, lastPage, err := s.OutletCodeRepository.List(filter, custIDs)
	if err != nil {
		return nil, 0, 0, err
	}
	out := make([]entity.OutletCodeItem, 0, len(rows))
	for _, r := range rows {
		item := entity.OutletCodeItem{
			Id:             r.Id,
			CustId:         r.CustId,
			SerialCode:     r.SerialCode,
			YearCode:       r.YearCode,
			LastSequenceNo: r.LastSequenceNo,
			Status:         mapStatus(r.Status),
			CreatedAt:      timePtrToStr(r.CreatedAt),
			CreatedBy:      strPtrToStr(r.CreatedBy),
			UpdatedAt:      timePtrToStr(r.UpdatedAt),
			UpdatedBy:      strPtrToStr(r.UpdatedBy),
		}
		out = append(out, item)
	}
	return out, total, lastPage, nil
}

func (s *outletCodeServiceImpl) Create(body entity.CreateOutletCodeBody, custId, createdBy string) error {
	exists, err := s.OutletCodeRepository.ExistsBySerialCodeAndYearCode(custId, strings.TrimSpace(body.SerialCode), body.YearCode)
	if err != nil {
		return err
	}
	if exists {
		return entity.ErrOutletCodeDuplicate
	}
	row := model.OutletCode{
		CustId:         custId,
		SerialCode:     strings.TrimSpace(body.SerialCode),
		YearCode:       body.YearCode,
		LastSequenceNo: strings.TrimSpace(body.LastSequenceNo),
		Status:         "Active",
		CreatedBy:      strPtr(createdBy),
		UpdatedBy:      strPtr(createdBy),
	}
	if err := s.OutletCodeRepository.Store(row); err != nil {
		return err
	}
	return nil
}

func (s *outletCodeServiceImpl) Update(id string, body entity.UpdateOutletCodeBody, custId, updatedBy string) error {
	existing, err := s.OutletCodeRepository.GetByID(id)
	if err != nil {
		return err
	}

	effectiveCustId := existing.CustId
	if strings.TrimSpace(custId) != "" {
		effectiveCustId = strings.TrimSpace(custId)
	}

	serial := strings.TrimSpace(body.SerialCode)
	if serial == "" {
		return errors.New("serial_code is required")
	}

	exists, err := s.OutletCodeRepository.ExistsBySerialCodeAndYearCodeExceptID(effectiveCustId, serial, existing.YearCode, existing.Id)
	if err != nil {
		return err
	}
	if exists {
		return entity.ErrOutletCodeDuplicate
	}

	if err := s.OutletCodeRepository.UpdateSerialCode(existing.Id, serial, strPtr(updatedBy)); err != nil {
		return err
	}
	return nil
}

func (s *outletCodeServiceImpl) UpdateStatus(id string, body entity.UpdateOutletCodeStatusBody, custId, updatedBy string) error {
	existing, err := s.OutletCodeRepository.GetByID(id)
	if err != nil {
		return err
	}

	raw := strings.TrimSpace(body.Status)
	if raw == "" {
		return errors.New("status is required")
	}
	var normalized string
	switch strings.ToLower(raw) {
	case "active":
		normalized = "Active"
	case "deactivate":
		normalized = "Deactivate"
	case "non active", "non_active":
		normalized = "Non Active"
	default:
		return errors.New("invalid status")
	}

	_ = existing

	if err := s.OutletCodeRepository.UpdateStatus(id, normalized, strPtr(updatedBy)); err != nil {
		return err
	}
	return nil
}

func (s *outletCodeServiceImpl) setupOutletCheckFindOne(custId string, year int, status []string, createdBy string) (*model.OutletCode, error) {
	row, err := s.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(custId, year, status, createdBy)
	if err != nil {
		return nil, err
	}
	if row != nil {
		return row, nil
	}
	row, err = s.OutletCodeRepository.FindOneByCustIdYearAndStatus(custId, year, status)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (s *outletCodeServiceImpl) SetupOutletCheck(custId, parentCustId string, year int, status []string, createdBy string) (*entity.SetupOutletCheckData, error) {
	row, err := s.setupOutletCheckFindOne(custId, year, status, createdBy)
	if err != nil {
		return nil, err
	}
	parentCustId = strings.TrimSpace(parentCustId)
	if row == nil && parentCustId != "" && parentCustId != custId {
		row, err = s.setupOutletCheckFindOne(parentCustId, year, status, createdBy)
		if err != nil {
			return nil, err
		}
	}
	if row == nil {
		return nil, nil
	}
	return &entity.SetupOutletCheckData{
		Id:             row.Id,
		CustId:         row.CustId,
		SerialCode:     row.SerialCode,
		YearCode:       row.YearCode,
		LastSequenceNo: row.LastSequenceNo,
	}, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func mapStatus(s string) string {
	if s == "" {
		return ""
	}
	switch strings.ToLower(s) {
	case "deactivate":
		return "Deactivate"
	case "non active", "non_active":
		return "Non Active"
	case "active":
		return "Active"
	default:
		return s
	}
}

func timePtrToStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func strPtrToStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
