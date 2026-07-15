package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryGdsImpl struct {
		*gorm.DB
	}
)

func NewGdsRepo(db *gorm.DB) *RepositoryGdsImpl {
	return &RepositoryGdsImpl{db}
}

type GdsRepository interface {
	Store(c context.Context, data *model.Gds) error
	StoreDetail(c context.Context, data *model.GdsDet) error
	FindByNo(gdNo, custId, parentCustId string) (whAdj model.GdsList, err error)
	FindDetail(gdNo string, custId string, langId string) (Details []model.GdsDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.GdsList, int64, int, error)
	Delete(c context.Context, custId string, gdsNo string, deletedBy int64) error
	Update(c context.Context, gdsNo string, data model.Gds) error
	DeleteDetailNotInIDs(c context.Context, gdsNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.GdsDet) error
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryGdsImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryGdsImpl) Store(c context.Context, data *model.Gds) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryGdsImpl) StoreDetail(c context.Context, data *model.GdsDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryGdsImpl) FindByNo(gdsNo, custId, parentCustId string) (whAdj model.GdsList, err error) {
	err = repository.Select("inv.gds.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, sup.sup_code, sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = inv.gds.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gds.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.gds.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_supplier sup on sup.sup_id = inv.gds.sup_id AND sup.cust_id = ?", parentCustId).
		Where("inv.gds.gds_no = ? AND inv.gds.cust_id=?", gdsNo, custId).
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryGdsImpl) FindDetail(gdsNo string, custId string, langId string) (Details []model.GdsDetRead, err error) {
	err = repository.Select("gds_det.*, p.pro_code, p.pro_name, st.status_name AS item_cnd_name").
		Joins("left join mst.m_product p on p.pro_id = gds_det.pro_id").
		Joins("left join mst.m_status st on st.status_id = 'item_cdn' AND st.status_value = inv.gds_det.item_cnd AND st.lang_id = '"+langId+"'").
		Where("inv.gds_det.gds_no = ? AND inv.gds_det.cust_id=?", gdsNo, custId).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryGdsImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.GdsList, int64, int, error) {
	var gdSo []model.GdsList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("gds_no")
	query := repository.Select("inv.gds.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, wh.wh_code, wh.wh_name, sup.sup_code, sup.sup_name").
		Joins("left join sys.m_user us on us.user_id = inv.gds.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = gds.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = inv.gds.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sup on sup.sup_id = inv.gds.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("inv.gds.cust_id=?", dataFilter.CustId)
	query.Where("inv.gds.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("inv.gds.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("inv.gds.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("inv.gds.gds_no=?", dataFilter.Query)
		queryCount.Where("inv.gds.gds_no=?", dataFilter.Query)
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
		query.Order("gds_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&gdSo).Error
	if err != nil {
		return gdSo, total, 0, err
	}
	err = queryCount.Model(&gdSo).Count(&total).Error
	if err != nil {
		return gdSo, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return gdSo, total, lastPage, nil
}

func (repository *RepositoryGdsImpl) Delete(c context.Context, custId string, gdsNo string, deletedBy int64) error {
	var data model.Gds
	result := repository.model(c).Model(&data).Where("gds_no=? AND cust_id = ? AND is_del= ? ", gdsNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryGdsImpl) Update(c context.Context, gdsNo string, data model.Gds) error {
	result := repository.model(c).Model(&data).Where("gds_no=?", gdsNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryGdsImpl) DeleteDetailNotInIDs(c context.Context, gdsNo string, IDs []int64) error {
	var Details model.GdsDet
	err := repository.model(c).Where("gds_no=? AND gds_det_id not in (?) ", gdsNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryGdsImpl) UpdateDetail(c context.Context, Details *model.GdsDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
