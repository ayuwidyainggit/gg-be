package repository

import (
	"context"
	"errors"
	"fmt"
	"sales/model"

	"gorm.io/gorm"
)

type (
	RepositoryOrderApprovalRequestImpl struct {
		*gorm.DB
	}
)
type OrderApprovalRequestRepository interface {
	Store(c context.Context, data *model.OrderApprovalRequest) error
	StoreDetail(c context.Context, data *model.OrderApprovalRequestDetail) error
	FindOneByID(orderApprovalRequestID int64, custID string) (model.OrderApprovalRequestRead, error)
	FindDetail(orderApprovalRequestID int64) (details []model.OrderApprovalRequestDetailRead, err error)
	FindOneByRoNo(roNo string, custID string) (model.OrderApprovalRequestRead, error)
	FindApprovalProcessedByRoNo(roNo string, custID string) (model.OrderApprovalRequestRead, error)
}

func NewOrderApprovalRequestRepo(db *gorm.DB) *RepositoryOrderApprovalRequestImpl {
	return &RepositoryOrderApprovalRequestImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryOrderApprovalRequestImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryOrderApprovalRequestImpl) Store(c context.Context, data *model.OrderApprovalRequest) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderApprovalRequestImpl) StoreDetail(c context.Context, data *model.OrderApprovalRequestDetail) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryOrderApprovalRequestImpl) FindOneByID(orderApprovalRequestID int64, custID string) (model.OrderApprovalRequestRead, error) {
	orderApprovalRequest := model.OrderApprovalRequestRead{}
	err := repository.Select("*").Where("order_approval_request_id = ? AND cust_id = ?", orderApprovalRequestID, custID).Take(&orderApprovalRequest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return orderApprovalRequest, errors.New(fmt.Sprintf("order approval request id : %v not found on cust_id %v", orderApprovalRequestID, custID))
		}

		return orderApprovalRequest, err
	}
	return orderApprovalRequest, nil
}

func (repository *RepositoryOrderApprovalRequestImpl) FindOneByRoNo(roNo string, custID string) (model.OrderApprovalRequestRead, error) {
	orderApprovalRequest := model.OrderApprovalRequestRead{}
	err := repository.Select("*").Where("ro_no = ? AND cust_id = ?", roNo, custID).Take(&orderApprovalRequest).Error
	if err != nil {

		return orderApprovalRequest, err
	}
	return orderApprovalRequest, nil
}

func (repository *RepositoryOrderApprovalRequestImpl) FindApprovalProcessedByRoNo(roNo string, custID string) (model.OrderApprovalRequestRead, error) {
	orderApprovalRequest := model.OrderApprovalRequestRead{}
	err := repository.Select("*").Where("ro_no = ? AND cust_id = ? AND finished_at is null", roNo, custID).Take(&orderApprovalRequest).Error
	if err != nil {

		return orderApprovalRequest, err
	}
	return orderApprovalRequest, nil
}

func (repository *RepositoryOrderApprovalRequestImpl) FindDetail(orderApprovalRequestID int64) (details []model.OrderApprovalRequestDetailRead, err error) {
	err = repository.Select(`sls.order_approval_requests_details.*, emp.emp_name, emp.emp_code, emp.image_url`).
		Joins("LEFT JOIN mst.m_employee emp ON sls.order_approval_requests_details.emp_id = emp.emp_id").
		Where("sls.order_approval_requests_details.order_approval_request_id=?", orderApprovalRequestID).
		Order("level ASC, seq ASC").
		Find(&details).Error
	return details, err
}
