package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryMEmployeeImpl struct {
		*gorm.DB
	}
)
type MEmployeeRepository interface {
	FindOneByEmail(email string) (employee model.MEmployee, err error)
	FindOneByEmailCustID(email string, custID string) (employee model.MEmployee, err error)
	// Update(c context.Context, ConfigId string, data model.MEmployee) error
	FindOneByDistributor(distId int) (distributor model.DistributorDetail, err error)
}

func NewMEmployeeRepository(db *gorm.DB) *RepositoryMEmployeeImpl {
	return &RepositoryMEmployeeImpl{db}
}

func (repo *RepositoryMEmployeeImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMEmployeeImpl) FindOneByEmail(email string) (employee model.MEmployee, err error) {
	err = repository.
		Where("emp_type_id IN ('D', 'S') AND is_active = true AND email = ?", email).
		Take(&employee).Error
	return employee, err
}

func (repository *RepositoryMEmployeeImpl) FindOneByEmailCustID(email string, custID string) (employee model.MEmployee, err error) {
	err = repository.
		Where("cust_id = ? AND emp_type_id IN ('D', 'S') AND is_active = true AND email = ?", custID, email).
		Take(&employee).Error
	return employee, err
}

func (repository *RepositoryMEmployeeImpl) FindOneByDistributor(distId int) (distributor model.DistributorDetail, err error) {
	err = repository.Select(`
		distributor_id,
		distributor_code,
		distributor_name,
		address
	`).
		Table("mst.m_distributor").
		Where("distributor_id = ?", distId).
		Take(&distributor).Error
	return distributor, err
}

// func (repository *RepositoryMEmployeeImpl) Update(c context.Context, ConfigId string, data model.MEmployee) error {
// 	result := repository.model(c).Model(&data).Where("config_id=?", ConfigId).Updates(&data)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return errors.New("no rows affected")
// 	}
// 	return nil
// }
