package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryOutletBankImpl struct {
		*gorm.DB
	}
)

type OutletBankRepository interface {
	FindByOutletID(ctx context.Context, custID string, outletID int64) ([]model.OutletBankInfo, error)
	FindFirstByOutletID(ctx context.Context, custID string, outletID int64) (*model.OutletBankInfo, error)
}

func NewOutletBankRepository(db *gorm.DB) *RepositoryOutletBankImpl {
	return &RepositoryOutletBankImpl{db}
}

func (repo *RepositoryOutletBankImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

// FindByOutletID returns all bank accounts linked to the given outlet.
func (repo *RepositoryOutletBankImpl) FindByOutletID(ctx context.Context, custID string, outletID int64) ([]model.OutletBankInfo, error) {
	var results []model.OutletBankInfo

	err := repo.model(ctx).
		Table("mst.m_outlet_bank ob").
		Select(`
			ob.outlet_bank_id,
			ob.outlet_id,
			ob.bank_id,
			ob.account_no,
			ob.account_name,
			b.bank_code,
			b.bank_name
		`).
		Joins("LEFT JOIN mst.m_bank b ON b.bank_id = ob.bank_id AND b.cust_id = ?", custID).
		Where("ob.cust_id = ?", custID).
		Where("ob.outlet_id = ?", outletID).
		Where("b.is_del = ? OR b.is_del IS NULL", false).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindFirstByOutletID returns the first bank account linked to the given outlet, or nil if none exists.
func (repo *RepositoryOutletBankImpl) FindFirstByOutletID(ctx context.Context, custID string, outletID int64) (*model.OutletBankInfo, error) {
	var result model.OutletBankInfo

	err := repo.model(ctx).
		Table("mst.m_outlet_bank ob").
		Select(`
			ob.outlet_bank_id,
			ob.outlet_id,
			ob.bank_id,
			ob.account_no,
			ob.account_name,
			b.bank_code,
			b.bank_name
		`).
		Joins("LEFT JOIN mst.m_bank b ON b.bank_id = ob.bank_id AND b.cust_id = ?", custID).
		Where("ob.cust_id = ?", custID).
		Where("ob.outlet_id = ?", outletID).
		Where("b.is_del = ? OR b.is_del IS NULL", false).
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// result will be zero-value if no row found (Scan doesn't return ErrRecordNotFound)
	if result.OutletBankID == 0 {
		return nil, nil
	}

	return &result, nil
}
