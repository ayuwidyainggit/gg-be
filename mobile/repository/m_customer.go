package repository

import (
	"log"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryMCustomerImpl struct {
		*gorm.DB
	}
)
type MCustomerRepository interface {
	FindOneByCustId(custId string) (model.MCustomer, error)
}

func NewMCustomerRepository(db *gorm.DB) *RepositoryMCustomerImpl {
	return &RepositoryMCustomerImpl{db}
}

// func (repo *RepositoryMCustomerImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositoryMCustomerImpl) FindOneByCustId(custId string) (model.MCustomer, error) {
	customer := model.MCustomer{}

	err := repository.
		Select("smc.m_customer.*, mst.m_distributor.distributor_code, mst.m_distributor.distributor_name, mst.m_distributor.address as distributor_address").
		Joins("LEFT JOIN mst.m_distributor ON mst.m_distributor.distributor_id = smc.m_customer.distributor_id").
		Where("smc.m_customer.is_del = ? AND smc.m_customer.cust_id = ?", false, custId).
		Take(&customer).Error

	if err != nil {
		log.Println("err.Error():", err.Error())
		return customer, err
	}

	return customer, nil
}
