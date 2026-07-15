package repository

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"
)

type (
	RepositoryCndnImpl struct {
		*gorm.DB
	}
)

type CndnRepository interface {
	Store(c context.Context, data *model.Cndn) error
	FindByNo(CndnNo string, custId string, parentCustId string) (whAdj model.CndnGetDetil, err error)
	FindAllByCustId(dataFilter entity.CndnQueryFilter) ([]model.CndnGet, int64, int, error)
	Update(c context.Context, CndnNo string, custId string, data model.Cndn) error
	Delete(c context.Context, custId string, CndnNo string, deletedBy int64) error
}

func NewCndnRepo(db *gorm.DB) *RepositoryCndnImpl {
	return &RepositoryCndnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryCndnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryCndnImpl) Store(c context.Context, data *model.Cndn) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryCndnImpl) FindByNo(CndnNo string, custId string, parentCustId string) (whAdj model.CndnGetDetil, err error) {
	err = repository.Select(
		`acf.cndn.cust_id,
		acf.cndn.cndn_no,
		acf.cndn.cndn_date,
		acf.cndn.owner_id,
		acf.cndn.cndn_jenis,
		acf.cndn.amount,
		acf.cndn.last_transaction_date,
		acf.cndn.notes,
		acf.cndn.created_by,
		acf.cndn.created_at,
		acf.cndn.updated_by,
		acf.cndn.updated_at,
		us.user_fullname AS updated_by_name, 
		case when acf.cndn.owner_id=1 then o.outlet_id  else ms.sup_id end as outlet_id, 
		case when acf.cndn.owner_id=1 then o.outlet_name  else ms.sup_name end as outlet_name, 
		case when acf.cndn.owner_id=1 then o.outlet_code else ms.sup_code end as outlet_code,  
		mcndn.cndn_name, mcndn.cndn_id, mcndn.cndn_code,
		appo.payment_amount as used_amount, ab.payment as used_amount_outlet`).
		Joins("left join sys.m_user us on us.user_id = acf.cndn.updated_by").
		Joins("left join mst.m_cndn mcndn on mcndn.cndn_id = acf.cndn.cndn_type").
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cndn.outlet_id AND o.cust_id = ?", custId).
		Joins("left join mst.m_supplier ms on ms.sup_id = acf.cndn.outlet_id and ms.cust_id = ?", custId).
		Joins("left join ("+
			"select "+
			"acf.deposit_payment.document_no, "+
			"coalesce(SUM(acf.deposit_payment.payment_amount), 0) as payment "+
			"from "+
			"acf.deposit_payment "+
			"where acf.deposit_payment.pay_type = 5 AND acf.deposit_payment.cust_id = '"+custId+"' "+
			"group by "+
			"document_no"+
			") ab on ab.document_no = acf.cndn.cndn_no").
		Joins("left join acf.account_payable_payment_options appo on appo.document_no = acf.cndn.cndn_no and appo.cust_id = ?", custId).
		Where("acf.cndn.cndn_no = ? AND acf.cndn.cust_id=?", CndnNo, custId).
		Where("acf.cndn.is_del=false").
		Take(&whAdj).Error
	return whAdj, err
}

func (repository *RepositoryCndnImpl) FindAllByCustId(dataFilter entity.CndnQueryFilter) ([]model.CndnGet, int64, int, error) {
	var Cndn []model.CndnGet
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("cndn_no")

	query := repository.Select(`acf.cndn.*,us.user_fullname AS updated_by_name,case when acf.cndn.owner_id=1 then o.outlet_id  else ms.sup_id end as outlet_id, 
	case when acf.cndn.owner_id=1 then o.outlet_name  else ms.sup_name end as outlet_name, 
	case when acf.cndn.owner_id=1 then o.outlet_code else ms.sup_code end as outlet_code,  mcndn.cndn_name as cndn_type`).
		Joins("left join sys.m_user us on us.user_id = acf.cndn.updated_by").
		Joins("left join mst.m_cndn mcndn on mcndn.cndn_id = acf.cndn.cndn_type").
		Joins("left join mst.m_outlet o on o.outlet_id = acf.cndn.outlet_id AND o.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_supplier ms on ms.sup_id = acf.cndn.outlet_id and ms.cust_id = ?", dataFilter.CustId)
	queryCount.Where("acf.cndn.cust_id=?", dataFilter.CustId)
	query.Where("acf.cndn.cust_id=?", dataFilter.CustId)

	queryCount.Where("acf.cndn.is_del=false")
	query.Where("acf.cndn.is_del=false")

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("acf.cndn.cndn_date between ? AND ?", *dataFilter.From, *dataFilter.To)
		queryCount.Where("acf.cndn.cndn_date between ? AND ?", *dataFilter.From, *dataFilter.To)
	}

	if dataFilter.DocumentNo != "" {
		query.Where("acf.cndn.cndn_no like ?", "%"+dataFilter.DocumentNo+"%")
		queryCount.Where("acf.cndn.cndn_no like ?", "%"+dataFilter.DocumentNo+"%")
	}

	if dataFilter.OwnerId != 0 {
		query.Where("acf.cndn.owner_id=?", dataFilter.OwnerId)
		queryCount.Where("acf.cndn.owner_id=?", dataFilter.OwnerId)
	}

	if dataFilter.OwnerId == 1 {

		if dataFilter.OutletId != nil && *dataFilter.OutletId != 0 {
			query.Where("acf.cndn.outlet_id=?", *dataFilter.OutletId)
			queryCount.Where("acf.cndn.outlet_id=?", *dataFilter.OutletId)
		}
	}

	if dataFilter.OwnerId == 2 {
		if dataFilter.SuptId != nil && *dataFilter.SuptId != 0 {
			query.Where("acf.cndn.outlet_id=?", *dataFilter.SuptId)
			queryCount.Where("acf.cndn.outlet_id=?", *dataFilter.SuptId)
		}
	}

	if dataFilter.CndnJenis != nil && *dataFilter.CndnJenis != "" {
		query.Where("acf.cndn.cndn_jenis=?", *dataFilter.CndnJenis)
		queryCount.Where("acf.cndn.cndn_jenis=?", *dataFilter.CndnJenis)
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
		query.Order("cndn_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&Cndn).Error
	if err != nil {
		return Cndn, total, 0, err
	}
	err = queryCount.Model(&Cndn).Count(&total).Error
	if err != nil {
		return Cndn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return Cndn, total, lastPage, nil
}

func (repository *RepositoryCndnImpl) Delete(c context.Context, custId string, CndnNo string, deletedBy int64) error {
	var data model.Cndn
	result := repository.model(c).Model(&data).Where("cndn_no=? AND cust_id = ? AND is_del= ? ", CndnNo, custId, false).
		Updates(map[string]interface{}{"is_del": true, "deleted_by": deletedBy, "deleted_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryCndnImpl) Update(c context.Context, CndnNo string, custId string, data model.Cndn) error {

	result := repository.model(c).Model(&data).Where("cndn_no=? AND cust_id = ?", CndnNo, custId).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }

	return nil
}
