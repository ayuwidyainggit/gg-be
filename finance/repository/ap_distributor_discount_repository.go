package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryApDistributorDiscountImpl struct {
		*gorm.DB
	}
)

type ApDistributorDiscountRepository interface {
	Store(c context.Context, data *model.ApDistributorDiscount) error
	FindByID(ApDistributorDiscountId int64, custId string, parentCustID string) (ApCndn model.ApDistributorDiscountList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApDistributorDiscountList, int64, int, error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter) ([]model.ApDistributorDiscountList, int64, int, error)
	Delete(c context.Context, custId string, ApDistributorDiscountId int64, deletedBy int64) error
	Update(c context.Context, ApDistributorDiscountId int64, data model.ApDistributorDiscount, custId string) error
}

func NewApDistributorDiscountRepo(db *gorm.DB) *RepositoryApDistributorDiscountImpl {
	return &RepositoryApDistributorDiscountImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApDistributorDiscountImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApDistributorDiscountImpl) Store(c context.Context, data *model.ApDistributorDiscount) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryApDistributorDiscountImpl) FindByID(ApDistributorDiscountId int64, custId string, parentCustID string) (ApDistributorDiscount model.ApDistributorDiscountList, err error) {
	err = repository.Select("acf.account_payable_discounts.*, us.user_fullname AS updated_by_name, c.pro_code, c.pro_name ").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable_discounts.updated_by").
		Joins("left join mst.m_product c on c.pro_id = acf.account_payable_discounts.pro_id AND c.cust_id = ?", parentCustID).
		Where("acf.account_payable_discounts.distributor_discount_id = ? AND acf.account_payable_discounts.cust_id=?", ApDistributorDiscountId, custId).
		Take(&ApDistributorDiscount).Error
	return ApDistributorDiscount, err
}

func (repository *RepositoryApDistributorDiscountImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApDistributorDiscountList, int64, int, error) {
	var ApDistributorDiscount []model.ApDistributorDiscountList
	var total int64
	var limit, page, offset int

	queryCount := repository.Select("distributor_discount_id")
	query := repository.Select("acf.account_payable_discounts.*, us.user_fullname AS updated_by_name, c.pro_code, c.pro_name").
		Joins("left join sys.m_user us on us.user_id = acf.account_payable_discounts.updated_by").
		Joins("left join mst.m_product c on c.pro_id = acf.account_payable_discounts.pro_id AND c.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.account_payable_discounts.cust_id=?", dataFilter.CustId)
	query.Where("acf.account_payable_discounts.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable_discounts.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.account_payable_discounts.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {

		query.Where("c.pro_code LIKE ?", "%"+dataFilter.Query+"%")
		query.Or("c.pro_name LIKE ?", "%"+dataFilter.Query+"%")

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
		query.Order("distributor_discount_id DESC")
	}

	if dataFilter.Mode != "lookup" {
		if dataFilter.Limit == 0 {
			limit = 10
		} else {
			limit = dataFilter.Limit
		}
		page = dataFilter.Page
		if page-1 < 1 {
			page = 1
		}
		offset = (page - 1) * dataFilter.Limit

		query.Limit(limit).Offset(offset)
	}

	err := query.Find(&ApDistributorDiscount).Error
	if err != nil {
		return ApDistributorDiscount, total, 0, err
	}
	err = queryCount.Model(&ApDistributorDiscount).Count(&total).Error
	if err != nil {
		return ApDistributorDiscount, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ApDistributorDiscount, total, lastPage, nil
}

func (repository *RepositoryApDistributorDiscountImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter) ([]model.ApDistributorDiscountList, int64, int, error) {
	var ApDistributorDiscount []model.ApDistributorDiscountList

	query := repository.Select("acf.account_payable_discounts.*, c.pro_code, c.pro_name").
		Joins("left join mst.m_product c on c.pro_id = acf.account_payable_discounts.pro_id AND c.cust_id = ?", dataFilter.ParentCustId)

	query.Where("acf.account_payable_discounts.cust_id = ?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.account_payable_discounts.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("c.pro_code LIKE ?", "%"+dataFilter.Query+"%")
		query.Or("c.pro_name LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("distributor_discount_id DESC")
	}

	err := query.Find(&ApDistributorDiscount).Error
	if err != nil {
		return ApDistributorDiscount, 0, 0, err
	}
	return ApDistributorDiscount, 0, 0, nil
}

func (repository *RepositoryApDistributorDiscountImpl) Delete(c context.Context, custId string, ApDistributorDiscountId int64, deletedBy int64) error {
	var data model.ApDistributorDiscount
	result := repository.model(c).Model(&data).Where("distributor_discount_id =? AND cust_id = ? AND is_del= ? ", ApDistributorDiscountId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApDistributorDiscountImpl) Update(c context.Context, ApDistributorDiscountId int64, data model.ApDistributorDiscount, custId string) error {
	result := repository.model(c).Model(&data).Where("cust_id = ? AND distributor_discount_id = ?", custId, ApDistributorDiscountId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
