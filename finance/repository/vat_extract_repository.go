package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	RepositoryVatExtractImpl struct {
		*gorm.DB
	}
)

type VatExtractRepository interface {
	Store(c context.Context, data *model.VatExtract) error
	StoreDetail(c context.Context, data []model.VatExtractDetail) error
	FindExtractResult(vatExtractID int64, custID string, parentCustId string) ([]model.VatExtractDetailList, error)
	FindAllVatExtractByCustId(dataFilter entity.VatExtractResultQueryFilter) ([]model.VatExtractList, int64, int, error)
}

func NewVatExtractRepo(db *gorm.DB) *RepositoryVatExtractImpl {
	return &RepositoryVatExtractImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryVatExtractImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}
func (repository *RepositoryVatExtractImpl) Store(c context.Context, data *model.VatExtract) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVatExtractImpl) StoreDetail(c context.Context, data []model.VatExtractDetail) error {
	err := repository.model(c).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryVatExtractImpl) FindExtractResult(vatExtractID int64, custID string, parentCustId string) ([]model.VatExtractDetailList, error) {
	var details []model.VatExtractDetailList

	err := repository.Select(`ap.*,us.user_fullname AS updated_by_name, us.user_fullname AS created_by_name ,
			sup.sup_code,sup.sup_name,sup.address1 AS address,
			sup.tax_no AS npwp, CASE WHEN ved.vat_extract_id IS NULL THEN 'not extracted' ELSE 'extracted' END AS extract_status, 
			acf.vat_extracts.created_at AS extracted_at, acf.vat_extracts.vat_extract_type, acf.vat_extracts.invoice_type`).
		Joins("left join acf.vat_extract_details ved on acf.vat_extracts.vat_extract_id = ved.vat_extract_id").
		Joins("left join acf.account_payable ap on ved.reference_id = ap.account_payable_id").
		Joins("LEFT JOIN mst.m_supplier sup ON sup.sup_id = ap.sup_id AND sup.cust_id = ?", parentCustId).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = ap.updated_by").
		Where("acf.vat_extracts.cust_id = ? AND acf.vat_extracts.vat_extract_id=?", custID, vatExtractID).Find(&details).Error
	return details, err

}

func (repository *RepositoryVatExtractImpl) FindAllVatExtractByCustId(dataFilter entity.VatExtractResultQueryFilter) ([]model.VatExtractList, int64, int, error) {
	var vatExtracts []model.VatExtractList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("vat_extract_id").Where("cust_id=?", dataFilter.CustId)
	query := repository.Select("*").Where("cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.VatTypes) > 0 {
		query.Where("vat_extract_type in ?", dataFilter.VatTypes)
		queryCount.Where("vat_extract_type in ?", dataFilter.VatTypes)

	}

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
		query.Order("vat_extract_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&vatExtracts).Error
	if err != nil {
		return vatExtracts, total, 0, err
	}
	err = queryCount.Model(&vatExtracts).Count(&total).Error
	if err != nil {
		return vatExtracts, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return vatExtracts, total, lastPage, nil
}
