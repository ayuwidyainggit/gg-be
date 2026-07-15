package repository

import (
	"context"
	"fmt"
	"math"
	"mobile/entity"
	"mobile/model"
	"strings"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	FindAllPaymentTypeLookup(ctx context.Context) ([]model.PaymentType, error)
	StorePaymentTrx(ctx context.Context, data *model.PaymentTrx) error
	StorePaymentTrxDetail(ctx context.Context, data *model.PaymentTrxDet) error
	FindOnePaymentTrx(ctx context.Context, custId string, paymentTrxId int64) (model.PaymentTrx, error)
	FindDetailsByPaymentTrxId(ctx context.Context, custId string, paymentTrxId int64) ([]model.PaymentTrxDet, error)
	StoreBankTransfer(ctx context.Context, data *model.BankTransfer) error
	FindPaymentTypeByID(ctx context.Context, id int) (model.PaymentType, error)
	GetNewDocumentNo(ctx context.Context, custID string) (string, error)
	GetByFilter(ctx context.Context, dataFilter entity.CollectionPayQueryFilter) ([]entity.CollectionPayResponse, int64, int64, error)
}

type paymentRepositoryImpl struct {
	*gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{db}
}

func (repo *paymentRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.DB.WithContext(ctx)
}

func (repo *paymentRepositoryImpl) FindAllPaymentTypeLookup(ctx context.Context) ([]model.PaymentType, error) {
	var paymentTypes []model.PaymentType
	err := repo.model(ctx).
		Where("is_del = false").
		Order("payment_type_id ASC").
		Find(&paymentTypes).Error
	return paymentTypes, err
}

func (repo *paymentRepositoryImpl) StorePaymentTrx(ctx context.Context, data *model.PaymentTrx) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *paymentRepositoryImpl) StorePaymentTrxDetail(ctx context.Context, data *model.PaymentTrxDet) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *paymentRepositoryImpl) FindOnePaymentTrx(ctx context.Context, custId string, paymentTrxId int64) (model.PaymentTrx, error) {
	var paymentTrx model.PaymentTrx
	err := repo.model(ctx).
		Where("cust_id = ? AND payment_trx_id = ? AND is_del = false", custId, paymentTrxId).
		First(&paymentTrx).Error
	return paymentTrx, err
}

func (repo *paymentRepositoryImpl) FindDetailsByPaymentTrxId(ctx context.Context, custId string, paymentTrxId int64) ([]model.PaymentTrxDet, error) {
	var details []model.PaymentTrxDet
	err := repo.model(ctx).
		Where("cust_id = ? AND payment_trx_id = ? AND is_del = false", custId, paymentTrxId).
		Find(&details).Error
	return details, err
}

func (repo *paymentRepositoryImpl) StoreBankTransfer(ctx context.Context, data *model.BankTransfer) error {
	return repo.model(ctx).Create(data).Error
}

func (repo *paymentRepositoryImpl) FindPaymentTypeByID(ctx context.Context, id int) (model.PaymentType, error) {
	var paymentType model.PaymentType
	err := repo.model(ctx).
		Where("payment_type_id = ? AND is_del = false", id).
		First(&paymentType).Error
	return paymentType, err
}

func (repo *paymentRepositoryImpl) GetNewDocumentNo(ctx context.Context, custID string) (string, error) {
	var documentNo string
	err := repo.model(ctx).Raw(`
		SELECT
			'TRX' || TO_CHAR(CURRENT_DATE, 'YYMMDD') || LPAD(
				(COALESCE(
					CAST(SUBSTRING(MAX(document_no) FROM 10 FOR 4) AS INTEGER),
					0
				) + 1)::TEXT,
				4,
				'0'
			) AS document_no
		FROM acf.payment_trx
		WHERE cust_id = ?
		AND is_del = false
		AND document_no LIKE 'TRX' || TO_CHAR(CURRENT_DATE, 'YYMMDD') || '%'
	`, custID).Scan(&documentNo).Error
	return documentNo, err
}

func (repo *paymentRepositoryImpl) GetByFilter(ctx context.Context, dataFilter entity.CollectionPayQueryFilter) ([]entity.CollectionPayResponse, int64, int64, error) {
	var data []entity.CollectionPayResponse
	var counting int64
	limit := 10

	if dataFilter.Limit > 0 {
		limit = dataFilter.Limit
	}

	query := repo.WithContext(ctx).Table("acf.payment_trx pt").
		Select("pt.payment_trx_id, pt.outlet_id, pt.emp_id, pt.po_number, pt.document_no, pt.total_transaction, pt.payment_amount, pt.remaining_amount, ptd.pay_type, ptd.amount, TO_CHAR(pt.created_at, 'YYYY-MM-DD') as payment_date").
		Joins("JOIN acf.payment_trx_detail ptd on pt.payment_trx_id = ptd.payment_trx_id").
		Where("pt.invoice_no = ? AND pt.outlet_id = ? AND pt.trx_source = 'L' AND pt.is_del = false", dataFilter.InvoiceNo, dataFilter.OutletID)

	queryCount := repo.WithContext(ctx).Table("acf.payment_trx pt").
		Select("COUNT(pt.payment_trx_id)").
		Joins("JOIN acf.payment_trx_detail ptd on pt.payment_trx_id = ptd.payment_trx_id").
		Where("pt.invoice_no = ? AND pt.outlet_id = ? AND pt.trx_source = 'L' AND pt.is_del = false", dataFilter.InvoiceNo, dataFilter.OutletID)

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		query.Order(sortBy)
	} else {
		query.Order("pt.created_at DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		return nil, 0, 0, err
	}
	err = queryCount.Scan(&counting).Error
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int64(math.Ceil(float64(float64(counting) / float64(limit))))
	return data, int64(counting), lastPage, nil
}
