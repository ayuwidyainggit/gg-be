package repository

import (
	"system/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type (
	RepositorySmcMCustomerImpl struct {
		*gorm.DB
	}
)
type SmcMCustomerRepository interface {
	FindOneByDomain(domain string) (model.SmcMCustomer, error)
	FindOneByCustId(custId string) (model.SmcMCustomer, error)
}

func NewSmcMCustomerRepository(db *gorm.DB) *RepositorySmcMCustomerImpl {
	return &RepositorySmcMCustomerImpl{db}
}

// func (repo *RepositorySmcMCustomerImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositorySmcMCustomerImpl) FindOneByDomain(domain string) (model.SmcMCustomer, error) {
	user := model.SmcMCustomer{}

	err := repository.
		Where("is_del = ? AND domain=?", false, domain).
		Take(&user).Error

	if err != nil {
		log.Error("err.Error():", err.Error())
		return user, err
	}

	return user, nil
}

func (repository *RepositorySmcMCustomerImpl) FindOneByCustId(custId string) (model.SmcMCustomer, error) {
	user := model.SmcMCustomer{}

	err := repository.
		Where("is_del = ? AND cust_id = ?", false, custId).
		Take(&user).Error

	if err != nil {
		log.Error("err.Error():", err.Error())
		return user, err
	}

	return user, nil
}
