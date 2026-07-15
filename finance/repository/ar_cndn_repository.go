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
	RepositoryArCndnImpl struct {
		*gorm.DB
	}
)

type ArCndnRepository interface {
	Store(c context.Context, data *model.ArCndn) error
	FindById(ArCndnId int, custId, parentCustId string) (whAdj model.ArCndnList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ArCndnList, int64, int, error)
	Update(c context.Context, ArCndnId int, data model.ArCndn) error
	Delete(c context.Context, custId string, ArCndnId int, deletedBy int64) error
}

func NewArCndnRepo(db *gorm.DB) *RepositoryArCndnImpl {
	return &RepositoryArCndnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryArCndnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryArCndnImpl) Store(c context.Context, data *model.ArCndn) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryArCndnImpl) FindById(ArCndnId int, custId, parentCustId string) (whAdj model.ArCndnList, err error) {
	err = repository.Select("acf.ar_cndn.*, us.user_fullname AS updated_by_name,c.cndn_code,c.cndn_name,o.outlet_code,o.outlet_name").
		Joins("left join sys.m_user us on us.user_id = acf.ar_cndn.updated_by").
		Joins("left join mst.m_outlet o on o.outlet_id = acf.ar_cndn.outlet_id AND o.cust_id = ?", custId).
		Joins("left join mst.m_cndn c on c.cndn_id = acf.ar_cndn.cndn_id AND c.cust_id = ?", parentCustId).
		Where("ar_cndn_id = ? AND acf.ar_cndn.cust_id=?", ArCndnId, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryArCndnImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ArCndnList, int64, int, error) {
	var ArCndn []model.ArCndnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ar_cndn_no")
	query := repository.Select("acf.ar_cndn.*, us.user_fullname AS updated_by_name,c.cndn_code,c.cndn_name,o.outlet_code,o.outlet_name").
		Joins("left join sys.m_user us on us.user_id = acf.ar_cndn.updated_by").
		Joins("left join mst.m_outlet o on o.outlet_id = acf.ar_cndn.outlet_id AND o.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_cndn c on c.cndn_id = acf.ar_cndn.cndn_id AND c.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("acf.ar_cndn.cust_id=?", dataFilter.CustId)
	query.Where("acf.ar_cndn.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.ar_cndn.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.ar_cndn.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.ar_cndn.ar_cndn_no=?", dataFilter.Query)
		query.Where("acf.ar_cndn.ar_cndn_no=?", dataFilter.Query)

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
		query.Order("ar_cndn_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ArCndn).Error
	if err != nil {
		return ArCndn, total, 0, err
	}
	err = queryCount.Model(&ArCndn).Count(&total).Error
	if err != nil {
		return ArCndn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ArCndn, total, lastPage, nil
}

func (repository *RepositoryArCndnImpl) Delete(c context.Context, custId string, ArCndnId int, deletedBy int64) error {
	var data model.ArCndn
	result := repository.model(c).Model(&data).Where("ar_cndn_id=? AND cust_id = ? AND is_del= ? ", ArCndnId, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryArCndnImpl) Update(c context.Context, ArCndnId int, data model.ArCndn) error {

	result := repository.model(c).Model(&data).Where("ar_cndn_id=?", ArCndnId).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
