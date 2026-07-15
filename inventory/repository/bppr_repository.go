package repository

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/str"
	"log"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	RepositoryBpprImpl struct {
		*gorm.DB
	}
)

type BpprRepository interface {
	Store(c context.Context, data *model.Bppr) error
	FindByNo(bpprNo, custId, parentCustId string) (bppr model.BpprList, err error)
	FindActiveWorkDay(custId string) (workDay model.WorkDay, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.BpprList, int64, int, error)
	Update(c context.Context, bpprNo string, data model.Bppr) (string, error)
	Delete(custId string, bpprNo string, deletedBy int64) error
	FindBpprDetails(bpprno string, custId string) (Details []model.BpprDetRead, err error)
	CreateBpprDetail(c context.Context, detail *model.BpprDet) (*model.BpprDet, error)
	UpdateBpprDetail(c context.Context, Detail *model.BpprDet) error
	DeleteBpprDetailByBpprNo(c context.Context, bpprNo string) error
}

func NewBpprRepo(db *gorm.DB) *RepositoryBpprImpl {
	return &RepositoryBpprImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryBpprImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryBpprImpl) DeleteBpprDetailByBpprNo(c context.Context, bpprNo string) error {
	var bpprDetails model.BpprDet
	err := repository.model(c).Where("bppr_no = ?", bpprNo).Delete(&bpprDetails).Error
	return err
}

func (repository *RepositoryBpprImpl) FindByNo(bpprNo, custId, parentCustId string) (bppr model.BpprList, err error) {
	err = repository.
		Select("bppr.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = bppr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = bppr.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = bppr.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_supplier sup on sup.sup_id = bppr.sup_id AND sup.cust_id = ?", parentCustId).
		Where("bppr.bppr_no = ? AND bppr.cust_id=?", bpprNo, custId).
		Take(&bppr).Error
	return bppr, err
}

func (repository *RepositoryBpprImpl) FindActiveWorkDay(custId string) (workDay model.WorkDay, err error) {
	err = repository.
		Select("per_year, per_id, week_id, work_date, is_active, is_closed").
		Where("is_active = true AND cust_id = ?", custId).
		Take(&workDay).Error
	return workDay, err
}

func (repository *RepositoryBpprImpl) Store(c context.Context, data *model.Bppr) error {
	err := repository.model(c).Create(data)
	if err != nil {
		return err.Error
	}
	return nil
}

func (repository *RepositoryBpprImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.BpprList, int64, int, error) {
	var bpprs []model.BpprList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("bppr_no")
	query := repository.Select("bppr.*,us.user_fullname AS updated_by_name,us2.user_fullname AS closed_by_name, sup.sup_code, sup.sup_name, wh.wh_code, wh.wh_name").
		Joins("left join sys.m_user us on us.user_id = bppr.updated_by").
		Joins("left join sys.m_user us2 on us2.user_id = bppr.closed_by").
		Joins("left join mst.m_warehouse wh on wh.wh_id = bppr.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier sup on sup.sup_id = bppr.sup_id AND sup.cust_id = ?", dataFilter.ParentCustId)

	queryCount.Where("bppr.cust_id=?", dataFilter.CustId)
	query.Where("bppr.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("bppr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("bppr.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("bppr.bppr_no=?", dataFilter.Query)
		queryCount.Where("bppr.bppr_no=?", dataFilter.Query)
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
		query.Order("bppr_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&bpprs).Error
	if err != nil {
		return bpprs, total, 0, err
	}

	err = queryCount.Model(&bpprs).Count(&total).Error
	if err != nil {
		log.Println("queryCount, err:", err.Error())
		return bpprs, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))

	return bpprs, total, lastPage, nil

}

func (repository *RepositoryBpprImpl) Update(c context.Context, bpprNo string, data model.Bppr) (string, error) {
	dataUpdate := data
	dataUpdate.CustID = nil

	result := repository.model(c).
		Model(&data).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "return_no"}}}).
		Where("bppr_no=? AND cust_id=?", bpprNo, data.CustID).Updates(dataUpdate)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("no rows affected")
	}
	// log.Println("result.RowsAffected:", structs.StructToJson(result.RowsAffected))
	// log.Println("data:", structs.StructToJson(data))
	// log.Println("dataUpdate:", structs.StructToJson(dataUpdate))
	return data.ReturnNo, nil
}

func (repository *RepositoryBpprImpl) Delete(custId string, bpprNo string, deletedBy int64) error {
	var data model.Bppr
	result := repository.Model(&data).Where("bppr_no=? AND cust_id = ? AND is_del= ? ", bpprNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryBpprImpl) FindBpprDetails(bpprno string, custId string) (Details []model.BpprDetRead, err error) {
	err = repository.Select("bppr_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = bppr_det.pro_id").
		Where("bppr_no = ? AND bppr_det.cust_id=?", bpprno, custId).Order("seq_no ASC").
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryBpprImpl) CreateBpprDetail(c context.Context, detail *model.BpprDet) (*model.BpprDet, error) {
	result := repository.model(c).Create(detail)
	if result.Error != nil {
		return detail, result.Error
	}
	if result.RowsAffected == 0 {
		return detail, errors.New("no rows affected")
	}

	return detail, nil
}
func (repository *RepositoryBpprImpl) UpdateBpprDetail(c context.Context, Detail *model.BpprDet) error {
	result := repository.model(c).Updates(&Detail)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
