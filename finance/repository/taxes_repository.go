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
	RepositoryTaxesImpl struct {
		*gorm.DB
	}
)

type TaxesRepository interface {
	GetInvoiceInfo(custId string, invoices []string) (orders []model.Order, err error)
	FindAllReportTaxes(dataFilter entity.TaxesQueryFilter, custId string) ([]model.TaxesList, int64, int, error)
	Store(c context.Context, data []model.Taxes) error
	TaxGenerateList(dataFilter entity.TaxesGenerateQueryFilter, custId string) ([]model.TaxesGenerateRead, int64, int, error)
	Delete(c context.Context, custId string, id int64, deletedBy int64) error
	CountTaxesByStatusAndMTax(custId string, mtaxID int64, status int) (int, error)
	DeleteBulk(c context.Context, custId string, ids []int64, deletedBy int64) error
}

func NewTaxesRepo(db *gorm.DB) *RepositoryTaxesImpl {
	return &RepositoryTaxesImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryTaxesImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryTaxesImpl) GetInvoiceInfo(custId string, invoices []string) (orders []model.Order, err error) {
	err = repository.Select("sls.order.outlet_id, sls.order.invoice_no, m.tax_invoice_form").
		Joins("JOIN mst.m_outlet m on m.outlet_id = sls.order.outlet_id").
		Where("sls.order.cust_id = ? AND sls.order.invoice_no in ?", custId, invoices).Find(&orders).
		Error

	return orders, err
}

func ternary(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func (repository *RepositoryTaxesImpl) FindAllReportTaxes(dataFilter entity.TaxesQueryFilter, custId string) ([]model.TaxesList, int64, int, error) {
	var taxes []model.TaxesList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	tableFrom := ternary(dataFilter.Type == "invoice", "sls.order", "sls.return")
	fieldDate := ternary(dataFilter.Type == "invoice", "ro_date", "return_date")
	fieldNo := ternary(dataFilter.Type == "invoice", "ro_no", "return_no")

	queryCount := repository.Select(tableFrom + "." + fieldNo + "")
	query := repository.Select(`
			`+tableFrom+`.cust_id, 
			`+tableFrom+`.`+fieldNo+`, 
			`+tableFrom+`.`+fieldDate+`,
			`+tableFrom+`.salesman_id, 
			`+tableFrom+`.outlet_id, 
			`+tableFrom+`.pay_type, 
			`+tableFrom+`.mobile_id, 
			`+tableFrom+`.sub_total, 
			`+tableFrom+`.vat, 
			`+tableFrom+`.vat_value, 
			`+tableFrom+`.total, 
			`+tableFrom+`.data_status, 
			`+tableFrom+`.invoice_no, 
			`+tableFrom+`.invoice_date, 
			`+tableFrom+`.data_source, 
			`+tableFrom+`.deleted_at, 
			emp.emp_code AS salesman_code,
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, 
			ot.outlet_name, 
			ot.address1 AS outlet_address, 
			ot.latitude AS outlet_latitude, 
			ot.longitude AS outlet_longitude,
			sales.sales_name,
			tx.taxes_id,
			tx.tax_no,
			(SELECT otx.tax_no 
				FROM mst.m_outlet_tax otx 
				WHERE otx.outlet_id = ot.outlet_id AND otx.tax_no IS NOT NULL 
				ORDER BY otx.outlet_tax_id DESC 
				LIMIT 1) AS npwp
			`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = "+tableFrom+".updated_by").
		Joins("LEFT JOIN mst.m_salesman sales ON sales.emp_id = "+tableFrom+".salesman_id AND sales.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = "+tableFrom+".salesman_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = "+tableFrom+".outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where(tableFrom+".cust_id = ?", dataFilter.CustId)
	query.Where(tableFrom+".cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where(tableFrom+".created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where(tableFrom+".created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount = queryCount.Where(tableFrom+".invoice_no ILIKE ? OR tx.tax_no ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
		query = query.Where(tableFrom+".invoice_no ILIKE ? OR tx.tax_no ILIKE ?)", "%"+dataFilter.Query+"%", "%"+dataFilter.Query+"%")
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where(tableFrom+".salesman_id in ?", dataFilter.SalesmanId)
		query.Where(tableFrom+".salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where(tableFrom+".outlet_id in ?", dataFilter.OutletID)
		query.Where(tableFrom+".outlet_id in ?", dataFilter.OutletID)
	}

	if !dataFilter.Taxes {
		queryCount.Where("tx.invoice_no IS NULL").
			Joins("LEFT JOIN acf.taxes tx ON tx.invoice_no = "+tableFrom+".invoice_no AND tx.cust_id = ?", dataFilter.CustId)
		query.Where("tx.invoice_no IS NULL").
			Joins("LEFT JOIN acf.taxes tx ON tx.invoice_no = "+tableFrom+".invoice_no AND tx.cust_id = ?", dataFilter.CustId)
	} else {
		queryCount.Where("tx.status = 1").
			Joins("RIGHT JOIN acf.taxes tx ON tx.invoice_no = "+tableFrom+".invoice_no AND tx.cust_id = ?", dataFilter.CustId)
		query.Where("tx.status = 1").
			Joins("RIGHT JOIN acf.taxes tx ON tx.invoice_no = "+tableFrom+".invoice_no AND tx.cust_id = ?", dataFilter.CustId)
	}

	queryCount.Where(tableFrom + ".invoice_no IS NOT NULL")
	query.Where(tableFrom + ".invoice_no IS NOT NULL")

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
		query.Order(tableFrom + "." + fieldNo + " DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Table(tableFrom).Limit(limit).Offset(offset).Find(&taxes).Error
	if err != nil {
		return taxes, total, 0, err
	}
	err = queryCount.Table(tableFrom).Model(&taxes).Count(&total).Error
	if err != nil {
		return taxes, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return taxes, total, lastPage, nil
}

func (repository *RepositoryTaxesImpl) Store(c context.Context, data []model.Taxes) error {
	err := repository.model(c).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryTaxesImpl) TaxGenerateList(dataFilter entity.TaxesGenerateQueryFilter, custId string) ([]model.TaxesGenerateRead, int64, int, error) {
	var generateTaxes []model.TaxesGenerateRead
	var total int64
	var limit int

	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("taxes_id")
	query := repository.Select("acf.taxes.*, ord.invoice_date").
		Joins("left join sls.order ord on acf.taxes.invoice_no = ord.invoice_no AND ord.cust_id = ?", dataFilter.CustId)

	queryCount.Where("acf.taxes.cust_id=?", dataFilter.CustId)
	query.Where("acf.taxes.cust_id=?", dataFilter.CustId)

	if dataFilter.MTaxID != 0 {
		queryCount.Where("acf.taxes.m_tax_id=?", dataFilter.MTaxID)
		query.Where("acf.taxes.m_tax_id=?", dataFilter.MTaxID)
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
		query.Order("m_tax_id DESC")
	}
	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&generateTaxes).Error
	if err != nil {
		return generateTaxes, total, 0, err
	}
	err = queryCount.Model(&generateTaxes).Count(&total).Error
	if err != nil {
		return generateTaxes, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return generateTaxes, total, lastPage, nil
}

func (repository *RepositoryTaxesImpl) Delete(c context.Context, custId string, id int64, deletedBy int64) error {
	var data model.Taxes
	result := repository.model(c).Model(&data).Where("cust_id = ? AND taxes_id = ? ", custId, id).
		Updates(map[string]interface{}{"status": 0, "updated_by": deletedBy, "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *RepositoryTaxesImpl) CountTaxesByStatusAndMTax(custId string, mtaxID int64, status int) (int, error) {
	var total int64
	var data model.Taxes

	err := repository.Select("taxes_id").Where("m_tax_id = ? AND status = ?", mtaxID, status).Model(&data).Count(&total).Error

	if err != nil {
		return int(total), err
	}

	return int(total), nil
}

func (repository *RepositoryTaxesImpl) DeleteBulk(c context.Context, custId string, ids []int64, deletedBy int64) error {
	var data model.Taxes
	result := repository.model(c).Model(&data).Where("cust_id = ? AND taxes_id in ? ", custId, ids).
		Updates(map[string]interface{}{"status": 0, "updated_by": deletedBy, "updated_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
