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
	RepositoryInvoiceImpl struct {
		*gorm.DB
	}
)
type InvoiceRepository interface {
	FindByNo(RoNo string, custId string) (realInvoice model.InvoiceList, err error)
	FindDetail(RoNo string, custId string) (details []model.InvoiceDetRead, err error)
	FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error)
	FindAllByInvoiceNombersAndCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, error)
	GenerateInvoiceNo(c context.Context, custId string, invoiceDate time.Time) (string, error)

	Update(c context.Context, RoNo, custID string, data model.Invoice) error
	UpdateOutletStatusFromPreDormantIfSet(c context.Context, custID string, outletID int64, updatedBy int64) error

	Print(c context.Context, custId string, invoiceNo string, printedBy int64) error
}

func NewInvoiceRepo(db *gorm.DB) *RepositoryInvoiceImpl {
	return &RepositoryInvoiceImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *RepositoryInvoiceImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryInvoiceImpl) FindByNo(roNo string, custId string) (realInvoice model.InvoiceList, err error) {
	err = repository.
		Select(`sls.order.*, 
			us.user_fullname AS updated_by_name,
			ot.outlet_code, ot.outlet_name, ot.address1 AS outlet_address, 
			ot.latitude AS outlet_latitude, ot.longitude AS outlet_longitude,
			sales.sales_name, ot.payment_type, ot.top,
			
			sls.order.is_printed,
			sls.order.printed_by,
			sls.order.printed_at,
			printer.user_fullname AS printed_by_name,

			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join sys.m_user printer on printer.user_id = sls.order.printed_by").
		Joins("left join mst.m_salesman sales on sales.emp_id = sls.order.salesman_id AND sales.cust_id = ?", custId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", custId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", custId).
		Where("sls.order.ro_no = ? AND sls.order.cust_id=?", roNo, custId).
		Take(&realInvoice).Error
	return realInvoice, err
}

func (repository *RepositoryInvoiceImpl) FindDetail(roNo string, custId string) (details []model.InvoiceDetRead, err error) {
	err = repository.Select(`sls.order_detail.*, 
		p.pro_code, p.pro_name, p.weight, p.length, p.width, p.height, p.volume,p.volume1,p.volume2,p.volume3,p.weight1,p.weight2,p.weight3,
		p.conv_unit2, p.conv_unit3, p.unit_id1, p.unit_id2, p.unit_id3`).
		Joins("LEFT JOIN mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Where("sls.order_detail.ro_no = ? AND sls.order_detail.cust_id=?", roNo, custId).
		Find(&details).Error
	return details, err
}

func (repository *RepositoryInvoiceImpl) FindAllByCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, int64, int, error) {
	var invoice []model.InvoiceList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	queryCount := repository.Select("ro_no")
	query := repository.Select(`
			sls.order.cust_id, 
			sls.order.ro_date, 
			sls.order.val_date, 
			sls.order.due_date, 
			sls.order.salesman_id, 
			sls.order.wh_id, 
			wh.latitude as wh_latitude, 
			wh.longitude as wh_longitude, 
			sls.order.outlet_id, 
			sls.order.delivery_date, 
			sls.order.ro_no as order_no, 
			sls.order.po_no, 
			sls.order.vehicle_no, 
			sls.order.pay_type, 
			sls.order.reff_no, 
			sls.order.mobile_id, 
			sls.order.sub_total, 
			sls.order.disc, 
			sls.order.disc_value, 
			sls.order.promo_value, 
			sls.order.promo_value_final,
			sls.order.promo_bg_value, 
			sls.order.promo_bg_value_final,
			sls.order.cash_disc_value, 
			sls.order.tot_disc1, 
			sls.order.tot_disc2, 
			sls.order.vat, 
			sls.order.vat_value, 
			sls.order.total, 
			sls.order.data_status, 
			sls.order.invoice_no, 
			sls.order.invoice_date, 
			sls.order.data_source, 
			sls.order.deleted_at, 
			emp.emp_code AS salesman_code,
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, 
			ot.outlet_name, 
			ot.address1 AS outlet_address, 
			ot.latitude AS outlet_latitude, 
			ot.longitude AS outlet_longitude,
			sales.sales_name,
			wh.wh_code, 
			wh.wh_name
			`).
		Joins("LEFT JOIN sys.m_user us ON us.user_id = sls.order.updated_by").
		Joins("LEFT JOIN mst.m_salesman sales ON sales.emp_id = sls.order.salesman_id AND sales.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_employee emp ON emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_warehouse wh ON wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_outlet ot ON ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	queryCount.Where("sls.order.cust_id = ?", dataFilter.CustId)
	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
		queryCount.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		queryCount.Where("sls.order.ro_no=?", dataFilter.Query)
		query.Where("sls.order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		queryCount.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		queryCount.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.Status) > 0 {
		queryCount.Where("sls.order.data_status in ?", dataFilter.Status)
		query.Where("sls.order.data_status in ?", dataFilter.Status)
	}

	if dataFilter.IsInvoice != nil && *dataFilter.IsInvoice {
		queryCount.Where("sls.order.invoice_no IS NOT NULL")
		query.Where("sls.order.invoice_no IS NOT NULL")
	}

	// queryCount.Where("sls.order.data_status=4")
	// query.Where("sls.order.data_status=4")

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
		query.Order("sls.order.ro_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&invoice).Error
	if err != nil {
		return invoice, total, 0, err
	}
	err = queryCount.Model(&invoice).Count(&total).Error
	if err != nil {
		return invoice, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return invoice, total, lastPage, nil
}

func (repository *RepositoryInvoiceImpl) GenerateInvoiceNo(c context.Context, custId string, invoiceDate time.Time) (string, error) {
	var invoiceNo string
	err := repository.model(c).
		Raw("SELECT sls.generate_invoice_no(?, ?) AS invoice_no", custId, invoiceDate).
		Row().
		Scan(&invoiceNo)
	if err != nil {
		return "", err
	}
	return invoiceNo, nil
}

func (repository *RepositoryInvoiceImpl) FindAllByInvoiceNombersAndCustId(dataFilter entity.InvoiceQueryFilter) ([]model.InvoiceList, error) {
	var invoice []model.InvoiceList

	query := repository.Select(
		`sls.order.*,
			sls.order.ro_no as order_no,
			emp.emp_code AS salesman_code,
			us.user_fullname AS updated_by_name, 
			ot.outlet_code, ot.outlet_name, 
			sales.sales_name,
			wh.wh_code, wh.wh_name`).
		Joins("left join sys.m_user us on us.user_id = sls.order.updated_by").
		Joins("left join mst.m_salesman sales on sales.emp_id = sls.order.salesman_id AND sales.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_employee emp on emp.emp_id = sls.order.salesman_id AND emp.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_warehouse wh on wh.wh_id = sls.order.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("left join mst.m_outlet ot on ot.outlet_id = sls.order.outlet_id AND ot.cust_id = ?", dataFilter.CustId)

	query.Where("sls.order.cust_id=?", dataFilter.CustId)

	if dataFilter.From != nil && dataFilter.To != nil {
		query.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.Query != "" {
		query.Where("sls.order.ro_no=?", dataFilter.Query)
	}

	if len(dataFilter.SalesmanId) > 0 {
		query.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		query.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if len(dataFilter.InvoiceNo) > 0 {
		query.Where("sls.order.invoice_no in ?", dataFilter.InvoiceNo)
	}

	query.Where("sls.order.data_status=6")
	// log.Info("InvoiceRepository, FindAllByCustId")
	// sortBy := ``
	// if dataFilter.Sort != "" {
	// 	mSortBy := strings.Split(dataFilter.Sort, ",")
	// 	for _, row := range mSortBy {
	// 		colSort := strings.Split(row, ":")
	// 		if len(colSort) > 1 {
	// 			sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
	// 		}
	// 	}
	// 	sortBy = strings.TrimSuffix(sortBy, ", ")
	// 	query.Order(sortBy)
	// } else {
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
		query.Order("sls.order.ro_no DESC")
	}

	err := query.Find(&invoice).Error
	if err != nil {
		return invoice, err
	}

	return invoice, nil
}

func (repository *RepositoryInvoiceImpl) Update(c context.Context, RoNo, custID string, data model.Invoice) error {
	result := repository.model(c).Model(&data).Where("ro_no=? AND cust_id = ?", RoNo, custID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	// if result.RowsAffected == 0 {
	// 	return errors.New("no rows affected")
	// }
	return nil
}

func (repository *RepositoryInvoiceImpl) UpdateOutletStatusFromPreDormantIfSet(c context.Context, custID string, outletID int64, updatedBy int64) error {
	if outletID == 0 {
		return nil
	}
	return repository.model(c).Exec(`
		UPDATE mst.m_outlet
		SET outlet_status = CASE
		      WHEN pre_dormant_status IS NOT NULL AND pre_dormant_status <> 0
		      THEN pre_dormant_status
		      ELSE outlet_status
		    END,
		    last_trans_date = CURRENT_DATE,
		    updated_at = NOW(),
		    updated_by = ?
		WHERE cust_id = ?
		  AND outlet_id = ?`,
		updatedBy, custID, outletID,
	).Error
}

func (repository *RepositoryInvoiceImpl) Print(c context.Context, custId string, invoiceNo string, printedBy int64) error {
	var data model.Invoice
	result := repository.model(c).Model(&data).Where("invoice_no=? AND cust_id = ?", invoiceNo, custId).
		Updates(map[string]interface{}{"is_printed": true, "printed_by": printedBy, "printed_at": time.Now()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
