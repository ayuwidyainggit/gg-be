package repository

import (
	"mobile/entity"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryMTakingOrderImpl struct {
		*gorm.DB
	}
)
type MTakingOrderRepository interface {
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MTakingOrder, error)
}

func NewMTakingOrderRepository(db *gorm.DB) *RepositoryMTakingOrderImpl {
	return &RepositoryMTakingOrderImpl{db}
}

// func (repo *RepositoryMTakingOrderImpl) model(ctx context.Context) *gorm.DB {
// 	tx := extractTx(ctx)
// 	if tx != nil {
// 		return tx.WithContext(ctx)
// 	}
// 	return repo.WithContext(ctx)
// }

func (repository *RepositoryMTakingOrderImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.MTakingOrder, error) {
	var products []model.MTakingOrder

	query := repository.Select("*").Where("is_active = true AND cust_id = ?", dataFilter.CustId).Order("taking_order_id DESC")

	err := query.Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil

}
