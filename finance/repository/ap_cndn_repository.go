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
	RepositoryApCndnImpl struct {
		*gorm.DB
	}
)

type ApCndnRepository interface {
	Store(c context.Context, data *model.ApCndn) error
	FindByNo(ApCndnNo string, custId string) (ApCndn model.ApCndnList, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApCndnList, int64, int, error)
	Delete(c context.Context, custId string, ApCndnNo string, deletedBy int64) error
	Update(c context.Context, ApCndnNo string, data model.ApCndn) error
}

func NewApCndnRepo(db *gorm.DB) *RepositoryApCndnImpl {
	return &RepositoryApCndnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryApCndnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryApCndnImpl) Store(c context.Context, data *model.ApCndn) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryApCndnImpl) FindByNo(ApCndnNo string, custId string) (ApCndn model.ApCndnList, err error) {
	err = repository.Select("acf.ap_cndn.*, us.user_fullname AS updated_by_name,c.cndn_code,c.cndn_name, sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.ap_cndn.updated_by").
		Joins("left join mst.m_cndn c on c.cndn_id = acf.ap_cndn.cndn_id AND c.cust_id = ?", custId).
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.ap_cndn.sup_id AND sup.cust_id = ?", custId).
		Where("acf.ap_cndn.ap_cndn_no = ? AND acf.ap_cndn.cust_id=?", ApCndnNo, custId).
		Take(&ApCndn).Error
	return ApCndn, err
}
func (repository *RepositoryApCndnImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ApCndnList, int64, int, error) {
	var ApCndn []model.ApCndnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ap_cndn_no")
	query := repository.Select("acf.ap_cndn.*, us.user_fullname AS updated_by_name,c.cndn_code,c.cndn_name, sup.sup_code,sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = acf.ap_cndn.updated_by").
		Joins("left join mst.m_cndn c on c.cndn_id = acf.ap_cndn.cndn_id AND c.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sup on sup.sup_id = acf.ap_cndn.sup_id AND sup.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.ap_cndn.cust_id=?", dataFilter.CustId)
	query.Where("acf.ap_cndn.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.ap_cndn.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("acf.ap_cndn.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("acf.ap_cndn.ap_cndn_no=?", dataFilter.Query)
		query.Where("acf.ap_cndn.ap_cndn_no=?", dataFilter.Query)

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
		query.Order("ap_cndn_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ApCndn).Error
	if err != nil {
		return ApCndn, total, 0, err
	}
	err = queryCount.Model(&ApCndn).Count(&total).Error
	if err != nil {
		return ApCndn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ApCndn, total, lastPage, nil
}
func (repository *RepositoryApCndnImpl) Delete(c context.Context, custId string, ApCndnNo string, deletedBy int64) error {
	var data model.ApCndn
	result := repository.model(c).Model(&data).Where("ap_cndn_no=? AND cust_id = ? AND is_del= ? ", ApCndnNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryApCndnImpl) Update(c context.Context, ApCndnNo string, data model.ApCndn) error {
	result := repository.model(c).Model(&data).Where("ap_cndn_no=?", ApCndnNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
