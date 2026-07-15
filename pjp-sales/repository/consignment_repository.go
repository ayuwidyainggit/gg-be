package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sales/entity"
	"sales/model"
	"sales/pkg/str"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryConsignmentImpl struct {
		*gorm.DB
	}
)
type ConsignmentRepository interface {
	Store(c context.Context, data *model.Consignment) error
	StoreDetail(c context.Context, data *model.ConsignmentDet) error
	FindByNo(consNo string, custId string) (consg model.ConsignmentList, err error)
	FindDetail(consNo string, custId string) (Details []model.ConsignmentDetRead, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ConsignmentList, int64, int, error)
	Delete(c context.Context, custId string, consNo string, deletedBy int64) error
	Update(c context.Context, consNo string, data model.Consignment) error
	DeleteDetailNotInIDs(c context.Context, consNo string, IDs []int64) error
	UpdateDetail(c context.Context, Details *model.ConsignmentDet) error
}

func NewConsignmentRepo(db *gorm.DB) *RepositoryConsignmentImpl {
	return &RepositoryConsignmentImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryConsignmentImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryConsignmentImpl) Store(c context.Context, data *model.Consignment) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryConsignmentImpl) StoreDetail(c context.Context, data *model.ConsignmentDet) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryConsignmentImpl) FindByNo(consNo string, custId string) (consg model.ConsignmentList, err error) {
	err = repository.Select("sls.consign.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name,ot.outlet_code, ot.outlet_name").
		Joins("left join sys.m_user us on us.user_id = sls.consign.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.consign.salesman_id AND sls.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.consign.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.consign.cust_id=? AND sls.consign.cons_no = ? ", custId, consNo).
		Take(&consg).Error
	return consg, err
}

func (repository *RepositoryConsignmentImpl) FindDetail(consNo string, custId string) (Details []model.ConsignmentDetRead, err error) {
	err = repository.Select("sls.consign_det.*, p.pro_code, p.pro_name").
		Joins("left join mst.m_product p on p.pro_id = sls.consign_det.pro_id").
		Where("sls.consign_det.cust_id=? AND cons_no = ?", custId, consNo).
		Find(&Details).Error
	return Details, err
}
func (repository *RepositoryConsignmentImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter) ([]model.ConsignmentList, int64, int, error) {
	var consignment []model.ConsignmentList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("cons_no")
	query := repository.Select("sls.consign.*, us.user_fullname AS updated_by_name,sls.emp_id as salesman_code, sls.sales_name as salesman_name,ot.outlet_code, ot.outlet_name").
		Joins("left join sys.m_user us on us.user_id = sls.consign.updated_by").
		Joins("left join mst.m_salesman sls on sls.emp_id = sls.consign.salesman_id AND sls.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.consign.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.consign.cust_id=?", dataFilter.CustId)
	query.Where("sls.consign.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.consign.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.consign.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.consign.cons_no=?", dataFilter.Query)
		query.Where("sls.consign.cons_no=?", dataFilter.Query)
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
		query.Order("cons_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&consignment).Error
	if err != nil {
		return consignment, total, 0, err
	}
	err = queryCount.Model(&consignment).Count(&total).Error
	if err != nil {
		return consignment, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return consignment, total, lastPage, nil
}
func (repository *RepositoryConsignmentImpl) Delete(c context.Context, custId string, consNo string, deletedBy int64) error {
	var data model.Consignment
	result := repository.model(c).Model(&data).Where("cust_id = ? AND cons_no=? AND is_del= ? ", custId, consNo, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
func (repository *RepositoryConsignmentImpl) Update(c context.Context, consNo string, data model.Consignment) error {
	result := repository.model(c).Model(&data).Where("cons_no=?", consNo).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
func (repository *RepositoryConsignmentImpl) DeleteDetailNotInIDs(c context.Context, consNo string, IDs []int64) error {
	var Details model.ConsignmentDet
	err := repository.model(c).Where("cons_no=? AND cons_det_id not in (?) ", consNo, IDs).Delete(&Details).Error
	return err
}
func (repository *RepositoryConsignmentImpl) UpdateDetail(c context.Context, Details *model.ConsignmentDet) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}
