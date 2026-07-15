package repository

import (
	"context"
	"finance/entity"
	"finance/model"
	"finance/pkg/str"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

type (
	CoreTaxVatExtractRepositoryImpl struct {
		*gorm.DB
	}
)

type CoreTaxVatExtractRepository interface {
	FindInvoiceListByCustId(dataFilter entity.CoreTaxVatExtractQueryFilter) ([]model.CoretaxInvoiceOrderList, int64, int, error)
	Store(c context.Context, data *model.CoretaxVatExtract) error
	StoreDetail(c context.Context, data []model.CoretaxVatExtractDetail) error
	FindInvoiceExtractByID(coretaxVatExtraxtID int64, custID string) ([]model.CoretaxInvoiceOrderList, error)
	FindCoretaxVatById(coretaxVatExtraxtID int64) (coretaxVat model.CoretaxVatExtract, err error)
	FindDetailsByRoNo(roNo string, custId string, parentCustId string) (details []model.CoretaxInvoiceOrderDetailRead, err error)
	FindDetailsByRoNoAndItemType(roNo string, itemType int, custId string, parentCustId string) (details []*model.CoretaxInvoiceOrderDetailRead, err error)
}

func NewCoreTaxVatExtractRepository(db *gorm.DB) *CoreTaxVatExtractRepositoryImpl {
	return &CoreTaxVatExtractRepositoryImpl{db}
}

// model returns query model with context with or without transaction extracted from context
func (repo *CoreTaxVatExtractRepositoryImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *CoreTaxVatExtractRepositoryImpl) Store(c context.Context, data *model.CoretaxVatExtract) error {
	err := repository.model(c).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *CoreTaxVatExtractRepositoryImpl) StoreDetail(c context.Context, data []model.CoretaxVatExtractDetail) error {
	err := repository.model(c).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *CoreTaxVatExtractRepositoryImpl) FindCoretaxVatById(coretaxVatExtraxtID int64) (coretaxVat model.CoretaxVatExtract, err error) {
	err = repository.Select("*").
		Where("coretax_vat_extract_id=?", coretaxVatExtraxtID).
		Take(&coretaxVat).Error
	return coretaxVat, err
}

func (repository *CoreTaxVatExtractRepositoryImpl) FindInvoiceExtractByID(coretaxVatExtraxtID int64, custID string) ([]model.CoretaxInvoiceOrderList, error) {
	var ro []model.CoretaxInvoiceOrderList
	dataStatus := []int{6, 7}
	subquery := repository.Table("sls.order").
		Select(`cv.created_at AS tax_extract_date,
		wh.wh_id,
		wh.wh_code,
		wh.wh_name,
		s.emp_id AS sales_id,
		s.sales_name AS sales_name,
		em.emp_code AS sales_code,
		mo.outlet_id,
		mo.outlet_code,
		mo.outlet_name, 
		mo.tax_no,
		mo.address1 as outlet_address1,
		mo.address2 as outlet_address2,
		mo.tax_addr1 as outlet_tax_address1,
		mo.tax_addr2 as outlet_tax_address2,
		CASE WHEN cvd.coretax_vat_extract_id IS NULL THEN 'not extracted' ELSE 'extracted' END AS extract_status,
		motx.nitku,
		motx.tax_identifier_type,
		motx.tax_identifier_no,
		motx.tax_name,
		motx.address_tax,
		moc.identity_no,
		sls.order.*`).
		Joins("LEFT JOIN acf.coretax_vat_extracts_details cvd ON cvd.reference_id = sls.order.ro_no AND cvd.cust_id = ?", custID).
		Joins("LEFT JOIN acf.coretax_vat_extracts cv ON cvd.coretax_vat_extract_id = cv.coretax_vat_extract_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON sls.order.wh_id = wh.wh_id AND wh.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_salesman s ON sls.order.salesman_id = s.emp_id AND s.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_employee em on em.emp_id = s.emp_id AND em.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_outlet_tax motx on motx.outlet_id = mo.outlet_id AND motx.cust_id = ?", custID).
		Joins("LEFT JOIN mst.m_outlet_contact moc on moc.outlet_id = mo.outlet_id AND moc.cust_id = ?", custID)

		// Joins("LEFT JOIN mst.m_distributor md on md.cust_id = sls.order.cust_id").
	// Joins("LEFT JOIN mst.m_distributor_tax mdtx on mdtx.distributor_id = md.distributor_id")

	subquery.Where("sls.order.cust_id=? AND cvd.coretax_vat_extract_id = ? AND sls.order.data_status in ?", custID, coretaxVatExtraxtID, dataStatus)

	query := repository.Table("(?) AS subquery", subquery)
	err := query.Find(&ro).Error
	if err != nil {
		return ro, err
	}

	return ro, nil
}

func (repository *CoreTaxVatExtractRepositoryImpl) FindDetailsByRoNo(roNo string, custId string, parentCustId string) (details []model.CoretaxInvoiceOrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.pro_code_coretax, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3,unc1.unit_id_coretax as unit_id_coretax1, unc2.unit_id_coretax as unit_id_coretax2, unc3.unit_id_coretax as unit_id_coretax3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Joins("LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un1.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un1.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax ").
		Joins("LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax").
		Joins("LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax").
		Where("ro_no = ? AND sls.order_detail.cust_id = ?", roNo, custId).
		Find(&details).Error

	return details, err
}

func (repository *CoreTaxVatExtractRepositoryImpl) FindDetailsByRoNoAndItemType(roNo string, itemType int, custId string, parentCustId string) (details []*model.CoretaxInvoiceOrderDetailRead, err error) {
	err = repository.Select("sls.order_detail.*, p.pro_code, p.pro_name, p.pro_code_coretax, p.conv_unit2 as mconv_unit2,p.conv_unit3 as mconv_unit3,unc1.unit_id_coretax as unit_id_coretax1, unc2.unit_id_coretax as unit_id_coretax2, unc3.unit_id_coretax as unit_id_coretax3").
		Joins("left join mst.m_product p on p.pro_id = sls.order_detail.pro_id").
		Joins("LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = ?", parentCustId).
		Joins("LEFT JOIN mst.m_unit_coretax unc1 ON un1.unit_id_coretax = unc1.unit_id_coretax AND unc1.cust_id =?", parentCustId).
		Joins("LEFT JOIN mst.m_unit_coretax unc2 ON un2.unit_id_coretax = unc2.unit_id_coretax AND unc2.cust_id =?", parentCustId).
		Joins("LEFT JOIN mst.m_unit_coretax unc3 ON un3.unit_id_coretax = unc3.unit_id_coretax AND unc3.cust_id =?", parentCustId).
		Where("ro_no = ? AND sls.order_detail.item_type = ? AND sls.order_detail.cust_id = ?", roNo, itemType, custId).
		Order("sls.order_detail.order_detail_id ASC").
		Find(&details).Error

	return details, err
}

func (repository *CoreTaxVatExtractRepositoryImpl) FindInvoiceListByCustId(dataFilter entity.CoreTaxVatExtractQueryFilter) ([]model.CoretaxInvoiceOrderList, int64, int, error) {
	var ro []model.CoretaxInvoiceOrderList
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}
	dataStatus := []int{6, 7}

	subquery := repository.Table("sls.order").
		Select(`DISTINCT ON (sls.order.ro_no)  (SELECT MAX(cv.created_at) 
         FROM acf.coretax_vat_extracts cv 
         JOIN acf.coretax_vat_extracts_details cvd 
         ON cvd.coretax_vat_extract_id = cv.coretax_vat_extract_id 
         WHERE cvd.reference_id = sls.order.ro_no 
         AND cvd.cust_id = 'C220010001') AS tax_extract_date,
		wh.wh_id,
		wh.wh_code,
		wh.wh_name,
		s.emp_id AS sales_id,
		s.sales_name AS sales_name,
		em.emp_code AS sales_code,
		mo.outlet_id,
		mo.outlet_code,
		mo.outlet_name, 
		mo.tax_no,
		mo.address1 as outlet_address1,
		mo.address2 as outlet_address2,
		mo.tax_addr1 as outlet_tax_address1,
		mo.tax_addr2 as outlet_tax_address2,
		CASE WHEN cvd.coretax_vat_extract_id IS NULL THEN 'not extracted' ELSE 'extracted' END AS extract_status,
		motx.nitku,
		motx.tax_identifier_type,
		motx.tax_identifier_no,
		(sls.order.sub_total_final - sls.order.disc_value_final - sls.order.promo_value_final) AS dpp,
		sls.order.*`).
		Joins("LEFT JOIN acf.coretax_vat_extracts_details cvd ON cvd.reference_id = sls.order.ro_no AND cvd.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN acf.coretax_vat_extracts cv ON cvd.coretax_vat_extract_id = cv.coretax_vat_extract_id").
		Joins("LEFT JOIN mst.m_warehouse wh ON sls.order.wh_id = wh.wh_id AND wh.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_salesman s ON sls.order.salesman_id = s.emp_id AND s.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_employee em on em.emp_id = s.emp_id AND em.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_outlet mo on mo.outlet_id = sls.order.outlet_id AND mo.cust_id = ?", dataFilter.CustId).
		Joins("LEFT JOIN mst.m_outlet_tax motx on motx.outlet_id = mo.outlet_id AND motx.cust_id = ?", dataFilter.CustId)

	subquery.Where("sls.order.cust_id=? AND sls.order.data_status in ?", dataFilter.CustId, dataStatus)

	if dataFilter.From != nil && dataFilter.To != nil {
		subquery.Where("sls.order.created_at between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.From), str.UnixTimestampToUtcTime(*dataFilter.To))
	}

	if dataFilter.InvoiceFrom != nil && dataFilter.InvoiceTo != nil {
		subquery.Where("sls.order.invoice_date between ? AND ?", str.UnixTimestampToUtcTime(*dataFilter.InvoiceFrom), str.UnixTimestampToUtcTime(*dataFilter.InvoiceTo))
	}

	if len(dataFilter.SalesmanId) > 0 {
		subquery.Where("sls.order.salesman_id in ?", dataFilter.SalesmanId)
	}

	if len(dataFilter.OutletID) > 0 {
		subquery.Where("sls.order.outlet_id in ?", dataFilter.OutletID)
	}

	if dataFilter.ExtractionStatus != "" {
		if dataFilter.ExtractionStatus == "E" {
			subquery.Where("cvd.coretax_vat_extract_id is not null")
		} else if dataFilter.ExtractionStatus == "NE" {
			subquery.Where("cvd.coretax_vat_extract_id is null")
		}
	}

	query := repository.Table("(?) AS subquery", subquery)
	queryCount := repository.Table("(?) AS subquery", subquery)

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
		query.Order("ro_no DESC")
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	err := query.Limit(limit).Offset(offset).Find(&ro).Error
	if err != nil {
		return ro, total, 0, err
	}
	err = queryCount.Model(&ro).Count(&total).Error
	if err != nil {
		return ro, total, 0, err
	}

	lastPage := int(math.Ceil(float64(float64(total) / float64(limit))))
	return ro, total, lastPage, nil
}
