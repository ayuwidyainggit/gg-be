package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	CheckIsPrinciple(ctx context.Context, custId string) (bool, error)
	GetCustIdByEmpCode(ctx context.Context, empCode string) (string, error)
	GetCustIdByEmpId(ctx context.Context, empId string) (string, error)
}
type CustomerRepositoryImpl struct {
	Db *gorm.DB
}

func NewCustomerRepositoryImpl(Db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{Db: Db}
}

func (c CustomerRepositoryImpl) CheckIsPrinciple(ctx context.Context, custId string) (bool, error) {
	var distributorID *string
	var isPrincipal bool

	err := c.Db.WithContext(ctx).Table("smc.m_customer").
		Select("distributor_id").
		Where("cust_id = ?", custId).
		Scan(&distributorID).Error

	if err != nil {
		return false, fmt.Errorf("failed to check customer status: %v", err)
	}
	isPrincipal = len(custId) == 6 && distributorID == nil
	return isPrincipal, nil
}

func (c CustomerRepositoryImpl) GetCustIdByEmpCode(ctx context.Context, empCode string) (string, error) {
	var customerID string
	err := c.Db.WithContext(ctx).Table("mst.m_employee").
		Select("cust_id").
		Where("emp_code = ?", empCode).
		Scan(&customerID).Error

	if err != nil {
		return "", fmt.Errorf("failed to get cust id by emp code: %v", err)
	}

	return customerID, nil
}

func (c CustomerRepositoryImpl) GetCustIdByEmpId(ctx context.Context, empId string) (string, error) {
	var customerID string
	err := c.Db.WithContext(ctx).Table("mst.m_employee").
		Select("cust_id").
		Where("emp_id = ?", empId).
		Scan(&customerID).Error

	if err != nil {
		return "", fmt.Errorf("failed to get cust id by emp id: %v", err)
	}

	return customerID, nil
}
