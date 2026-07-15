package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type BankRepository interface {
	FindByCustIDAndOutletID(ctx context.Context, custID string, outletID int64) (*model.OutletBank, error)
	GetNewDocNoBank(ctx context.Context, custID string) (string, error)
}

type bankRepositoryImpl struct {
	*gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &bankRepositoryImpl{db}
}

func (repo *bankRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.DB.WithContext(ctx)
}

func (repo *bankRepositoryImpl) FindByCustIDAndOutletID(ctx context.Context, custID string, outletID int64) (*model.OutletBank, error) {
	var outletBank model.OutletBank
	tx := repo.model(ctx).Raw(`
	SELECT
		b.bank_name,
		b.bank_code,
		ob.account_name,
		ob.account_no,
		ob.bank_id,
		ob.outlet_id,
		ob.outlet_bank_id
	FROM mst.m_outlet_bank ob
	LEFT JOIN mst.m_bank b ON b.bank_id = ob.bank_id
		AND b.cust_id = ?
	WHERE ob.cust_id = ?
		AND ob.outlet_id = ?
		AND (b.is_del = false OR b.is_del IS NULL)
	LIMIT 1
	`, custID, custID, outletID)
	if err := tx.Scan(&outletBank).Error; err != nil {
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &outletBank, nil
}

func (repo *bankRepositoryImpl) GetNewDocNoBank(ctx context.Context, custID string) (string, error) {
	var docNo string
	err := repo.model(ctx).Raw(`
	SELECT
		'TF' || TO_CHAR(CURRENT_DATE, 'YYMMDD') || LPAD(
			(COALESCE(
				CAST(SUBSTRING(MAX(doc_no_bank) FROM 10 FOR 4) AS INTEGER),
				0
			) + 1)::TEXT,
			4,
			'0'
		) AS doc_no_bank
	FROM acf.bank_transfer
	WHERE cust_id = ?
	AND is_del = false
	AND doc_no_bank LIKE 'TF' || TO_CHAR(CURRENT_DATE, 'YYMMDD') || '%'
	`, custID).Scan(&docNo).Error
	if err != nil {
		return "", err
	}
	return docNo, nil
}
