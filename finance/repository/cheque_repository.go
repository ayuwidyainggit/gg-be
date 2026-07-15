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
	RepositoryChequeImpl struct {
		*gorm.DB
	}
)

type ChequeRepository interface {
	Store(c context.Context, data *model.Cheque) error
	FindByNo(ChequeNo int, custId string) (whAdj model.ChequeList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ChequeList, int64, int, error)
	Update(c context.Context, ChequeNo int, custId string, data model.Cheque) error
	Delete(c context.Context, custId string, ChequeNo int, deletedBy int64) error
}

func NewChequeRepo(db *gorm.DB) *RepositoryChequeImpl {
	return &RepositoryChequeImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryChequeImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryChequeImpl) Store(c context.Context, data *model.Cheque) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryChequeImpl) FindByNo(ChequeNo int, custId string) (whAdj model.ChequeList, err error) {
	err = repository.Select("acf.cheque.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name,b.bank_code,b.bank_name,o.outlet_code,o.outlet_name").
		Joins("left join sys.m_user us on us.user_id = acf.cheque.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.cheque.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_bank b on b.bank_id = acf.cheque.bank_id AND b.cust_id = ?", custId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cheque.outlet_id AND o.cust_id = ?", custId).
		Where("acf.cheque.chq_no = ? AND acf.cheque.cust_id=?", ChequeNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryChequeImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ChequeList, int64, int, error) {
	var Cheque []model.ChequeList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("chq_no")
	query := repository.Select("acf.cheque.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name,b.bank_code,b.bank_name,o.outlet_code,o.outlet_name").
		Joins("left join sys.m_user us on us.user_id = acf.cheque.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = acf.cheque.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_bank b on b.bank_id = acf.cheque.bank_id AND b.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cheque.outlet_id AND o.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.cheque.cust_id=?", dataFilter.CustId)
	query.Where("acf.cheque.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.cheque.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.cheque.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("acf.cheque.chq_tr_no=?", dataFilter.Query)
		queryCount.Where("acf.cheque.chq_tr_no=?", dataFilter.Query)
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
		query.Order("chq_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&Cheque).Error
	if err != nil {
		return Cheque, total, 0, err
	}
	err = queryCount.Model(&Cheque).Count(&total).Error
	if err != nil {
		return Cheque, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return Cheque, total, lastPage, nil
}

func (repository *RepositoryChequeImpl) Delete(c context.Context, custId string, ChequeNo int, deletedBy int64) error {
	var data model.Cheque
	result := repository.model(c).Model(&data).Where("chq_no=? AND cust_id = ? AND is_del= ?", ChequeNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryChequeImpl) Update(c context.Context, ChequeNo int, custId string, data model.Cheque) error {

	result := repository.model(c).Model(&data).Where("chq_no=? AND cust_id = ?", ChequeNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}
