package repository

import (
	"context"
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
	RepositoryStockReturnImpl struct {
		*gorm.DB
	}
)

type StockReturnRepository interface {
	FindAllByCustId(dataFilter entity.StockReturnQueryFilter) ([]model.StockReturnList, int64, int, error)
	FindOneByReturnNo(returnNo string, custId string, parentCustId string) (rtn model.StockReturnRead, err error)
	FindReturnDetail(returnNo string, custId string, parentCustId string) (Details []model.StockReturnDetailRead, err error)
	CountReturnedProductQty(invoiceNo string, productId int64, custId string) (qtySummary model.StockReturnedDetailRead, err error)
	UpdateDetail(c context.Context, Details *model.StockReturnDetail) error
	UpdateStatus(c context.Context, custId string, returnNo string, status int) error
	FindReturnDetailByListNo(returnNo []string, custId string, parentCustId string) (rtns []model.StockReturnRead, err error)
	UpdatebatchStatus(c context.Context, custId string, returnNo []string, status int) error
	UpdatebatchClosedAt(c context.Context, custId string, returnNo []string, closedAt time.Time) error
	UpdateClosedAt(c context.Context, custId string, returnNo string, closedAt time.Time) error
}

func NewStockReturnRepo(db *gorm.DB) *RepositoryStockReturnImpl {
	return &RepositoryStockReturnImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryStockReturnImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryStockReturnImpl) FindAllByCustId(dataFilter entity.StockReturnQueryFilter) ([]model.StockReturnList, int64, int, error) {
	var rtn []model.StockReturnList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("return_no")
	query := repository.Select(`
			sls.return.refference_no, 
			sls.return.return_no, 
			sls.return.invoice_no, 
			sls.return.invoice_date, 
			sls.return.salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name, 
			sls.return.outlet_id, 
			ot.outlet_code, 
			ot.outlet_name,
			sls.return.data_status, 
			sls.return.created_by, 
			creator.user_fullname AS created_by_name,
			sls.return.created_at, 
			sls.return.reviewed_by, 
			reviewer.user_fullname AS reviewed_by_name,
			sls.return.reviewed_at
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.return.salesman_id AND employee.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.return.outlet_id AND ot.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user creator on creator.user_id = sls.return.created_by AND creator.cust_id = ?", dataFilter.CustId).
		Joins("left join sys.m_user reviewer on reviewer.user_id = sls.return.reviewed_by AND reviewer.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.return.cust_id=? AND data_status in (5,6,9)", dataFilter.CustId)
	query.Where("sls.return.cust_id=? AND data_status in (5,6,9)", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.return.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.return.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.return.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.return.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.return.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.return.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.return.data_status in ?", dataFilter.Status)
		query.Where("sls.return.data_status in ?", dataFilter.Status)
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.return.return_no LIKE ?", "%"+dataFilter.Query+"%")
		query.Where("sls.return.return_no LIKE ?", "%"+dataFilter.Query+"%")
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
		query.Order("return_no DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&rtn).Error
	if err != nil {
		return rtn, total, 0, err
	}
	err = queryCount.Model(&rtn).Count(&total).Error
	if err != nil {
		return rtn, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return rtn, total, lastPage, nil
}

func (repository *RepositoryStockReturnImpl) FindOneByReturnNo(returnNo string, custId string, parentCustId string) (rtn model.StockReturnRead, err error) {
	err = repository.Select(`
			sls.return.cust_id, 
			sls.return.refference_no, 
			sls.return.return_no, 
			sls.return.return_date, 
			sls.return.invoice_no, 
			sls.return.invoice_date, 
			sls.return.salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name, 
			sls.return.outlet_id, 
			ot.outlet_code, 
			ot.outlet_name,
			sls.return.tpr_cash_value, 
			sls.return.tpr_item_value, 
			sls.return.discount, 
			sls.return.disc_value, 
			sls.return.vat, 
			sls.return.vat_value, 
			sls.return.sub_total, 
			sls.return.total, 
			sls.return.data_status 
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.return.salesman_id AND employee.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.return.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.return.cust_id=? AND sls.return.return_no = ? ", custId, returnNo).
		Take(&rtn).Error

	return rtn, err
}

func (repository *RepositoryStockReturnImpl) FindReturnDetail(returnNo string, custId string, parentCustId string) (Details []model.StockReturnDetailRead, err error) {
	err = repository.Select(`
			sls.return_det.return_detail_id, 
			sls.return_det.order_detail_id, 
			sls.return_det.return_no, 
			sls.return_det.product_id, 
			COALESCE(product_parent.pro_code, product_cust.pro_code) as product_code, 
			COALESCE(product_parent.pro_name, product_cust.pro_name) as product_name, 
			sls.return_det.item_cnd,
			sls.return_det.qty1, 
			sls.return_det.qty2, 
			sls.return_det.qty3, 
			invoice_detail.qty1 as invoice_qty1, 
			invoice_detail.qty2 as invoice_qty2, 
			invoice_detail.qty3 as invoice_qty3,
			sls.return_det.sell_price1, 
			sls.return_det.sell_price2, 
			sls.return_det.sell_price3, 
			sls.return_det.unit_id1, 
			sls.return_det.unit_id2,
			sls.return_det.unit_id3,
			unit1.unit_name as unit_name1, 
			unit2.unit_name as unit_name2, 
			unit3.unit_name as unit_name3, 
			sls.return_det.conv_unit2, 
			sls.return_det.conv_unit3, 
			sls.return_det.vat, 
			sls.return_det.vat_value, 
			sls.return_det.sub_total, 
			sls.return_det.total, 
			sls.return_det.return_reason_id, 
			return_reason.return_reason_code, 
			return_reason.return_reason_name,
			sls.return_det.wh_id,
			wh.wh_code,
			wh.wh_name 
		`).
		Joins("left join sls.order_detail invoice_detail on invoice_detail.order_detail_id = sls.return_det.order_detail_id AND invoice_detail.cust_id = ?", custId).
		Joins("left join mst.m_product product_parent on product_parent.pro_id = sls.return_det.product_id AND product_parent.cust_id = ?", parentCustId).
		Joins("left join mst.m_product product_cust on product_cust.pro_id = sls.return_det.product_id AND product_cust.cust_id = ?", custId).
		Joins("left join mst.m_return_reason return_reason on return_reason.return_reason_id = sls.return_det.return_reason_id AND return_reason.cust_id = ?", parentCustId).
		Joins("left join mst.m_unit unit1 on unit1.unit_id = sls.return_det.unit_id1 AND unit1.cust_id = ?", parentCustId).
		Joins("left join mst.m_unit unit2 on unit2.unit_id = sls.return_det.unit_id2 AND unit2.cust_id = ?", parentCustId).
		Joins("left join mst.m_unit unit3 on unit3.unit_id = sls.return_det.unit_id3 AND unit3.cust_id = ?", parentCustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.return_det.wh_id AND wh.cust_id = ?", parentCustId).
		Where("sls.return_det.cust_id=? AND sls.return_det.return_no = ? ", custId, returnNo).
		Find(&Details).Error

	return Details, err
}

func (repository *RepositoryStockReturnImpl) CountReturnedProductQty(invoiceNo string, productId int64, custId string) (qtySummary model.StockReturnedDetailRead, err error) {
	err = repository.Select(`
			coalesce(sum(sls.return_det.qty1), 0) as remaining_qty1, 
			coalesce(sum(sls.return_det.qty2), 0) as remaining_qty2, 
			coalesce(sum(sls.return_det.qty3), 0) as remaining_qty3
		`).
		Where("sls.return_det.return_no in (select sls.return.return_no from sls.return where sls.return.invoice_no = ? and sls.return.data_status = 5)", invoiceNo).
		Where("sls.return_det.product_id = ?", productId).
		Find(&qtySummary).Error

	return qtySummary, err
}

func (repository *RepositoryStockReturnImpl) UpdateStatus(c context.Context, custId string, returnNo string, status int) error {
	repository.Table("sls.return").Where("cust_id = ? AND return_no = ?", custId, returnNo).Update("data_status", status)

	return nil
}

func (repository *RepositoryStockReturnImpl) UpdateDetail(c context.Context, Details *model.StockReturnDetail) error {
	result := repository.model(c).Updates(&Details)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository *RepositoryStockReturnImpl) FindReturnDetailByListNo(returnNo []string, custId string, parentCustId string) (rtns []model.StockReturnRead, err error) {
	err = repository.Select(`
			sls.return.cust_id, 
			sls.return.refference_no, 
			sls.return.return_no, 
			sls.return.return_date, 
			sls.return.invoice_no, 
			sls.return.invoice_date, 
			sls.return.salesman_id, 
			employee.emp_code as salesman_code, 
			employee.emp_name as salesman_name, 
			sls.return.outlet_id, 
			ot.outlet_code, 
			ot.outlet_name,
			sls.return.tpr_cash_value, 
			sls.return.tpr_item_value, 
			sls.return.discount, 
			sls.return.disc_value, 
			sls.return.vat, 
			sls.return.vat_value, 
			sls.return.sub_total, 
			sls.return.total, 
			sls.return.data_status 
		`).
		Joins("left join mst.m_employee employee on employee.emp_id = sls.return.salesman_id AND employee.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.return.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.return.cust_id=? AND sls.return.return_no in ? ", custId, returnNo).
		Find(&rtns).Error

	return rtns, err
}

func (repository *RepositoryStockReturnImpl) UpdatebatchStatus(c context.Context, custId string, returnNo []string, status int) error {
	repository.Table("sls.return").Where("cust_id = ? AND return_no in ?", custId, returnNo).Update("data_status", status)

	return nil
}

func (repository *RepositoryStockReturnImpl) UpdateClosedAt(c context.Context, custId string, returnNo string, closedAt time.Time) error {
	repository.Table("sls.return").Where("cust_id = ? AND return_no = ?", custId, returnNo).Update("closed_at", closedAt)

	return nil
}

func (repository *RepositoryStockReturnImpl) UpdatebatchClosedAt(c context.Context, custId string, returnNo []string, closedAt time.Time) error {
	repository.Table("sls.return").Where("cust_id = ? AND return_no in ?", custId, returnNo).Update("closed_at", closedAt)

	return nil
}
