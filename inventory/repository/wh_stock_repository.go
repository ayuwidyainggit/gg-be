package repository

import (
	"context"
	"database/sql"
	"errors"
	"inventory/entity"
	"inventory/model"
	"log"

	"gorm.io/gorm"
)

type (
	RepositoryWhStockImpl struct {
		*gorm.DB
	}
)

type WhStockRepository interface {
	FindByWhIdAndProId(wsQuery entity.WhStockQuery) (whStock model.WhStock, err error)
	StoreWhStock(c context.Context, data *model.WhStock) (err error)
	UpdateWhStockByWhIdAndProId(c context.Context, whStock model.WhStock) error
	UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error
}

func NewWhStockRepo(db *gorm.DB) *RepositoryWhStockImpl {
	return &RepositoryWhStockImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryWhStockImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryWhStockImpl) FindByWhIdAndProId(wsQuery entity.WhStockQuery) (whStock model.WhStock, err error) {
	err = repository.
		Select("wh_stock.*").
		Where(`wh_stock.cust_id = ? 
		AND wh_stock.wh_id = ? 
		AND wh_stock.pro_id = ? `, wsQuery.CustID, wsQuery.WhId, wsQuery.ProId).
		Take(&whStock).Error
	isErrRecordNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	if isErrRecordNotFound {
		return whStock, nil
	}
	return whStock, err
}

func (repository *RepositoryWhStockImpl) UpdateOldWhStock(c context.Context, custId string, whId, proId int64, qty float64) error {
	err := repository.model(c).Exec(
		`UPDATE inv.wh_stock 
		SET qty = qty-@qty 
		WHERE cust_id = @cust_id AND pro_id = @pro_id AND wh_id = @wh_id;`,
		sql.Named("cust_id", custId),
		sql.Named("wh_id", whId),
		sql.Named("pro_id", proId),
		sql.Named("qty", qty)).Error
	if err != nil {
		log.Println("UpdateOldStock, error:", err.Error())
		return err
	}
	return nil
}

func (repository *RepositoryWhStockImpl) UpdateWhStockByWhIdAndProId(c context.Context, data model.WhStock) error {
	// dataModel :=
	dataUpdate := data
	dataUpdate.CustID = ""
	dataUpdate.WhID = nil
	dataUpdate.ProID = nil
	result := repository.model(c).Model(&data).
		Where(
			`cust_id = ? 
			AND wh_id = ? 
			AND pro_id = ?`, data.CustID, data.WhID, data.ProID).Updates(dataUpdate)
	if result.Error != nil {
		log.Println("UpdateWhStockByWhIdAndProId, error:", result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("UpdateWhStockByWhIdAndProId, no rows affected, pro_id:", data.ProID)
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryWhStockImpl) StoreWhStock(c context.Context, data *model.WhStock) error {
	err := repository.model(c).Exec(
		`INSERT INTO inv.wh_stock (
			cust_id, wh_id, pro_id, qty, stock_id
		) VALUES (
			@cust_id, @wh_id, @pro_id, @qty, @stock_id
		) ON CONFLICT ON CONSTRAINT wh_stock_pkey 
		DO UPDATE SET qty = inv.wh_stock.qty + EXCLUDED.qty, stock_id = @stock_id;`,
		sql.Named("cust_id", data.CustID),
		sql.Named("wh_id", data.WhID),
		sql.Named("pro_id", data.ProID),
		sql.Named("qty", data.Qty),
		sql.Named("stock_id", data.StockId)).Error
	if err != nil {
		log.Println("StoreWhStock, error:", err.Error())
		return err
	}
	return nil
}
