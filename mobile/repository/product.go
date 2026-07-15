package repository

import (
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"strconv"

	"gorm.io/gorm"
)

type (
	RepositoryProductImpl struct {
		*gorm.DB
	}
)

type ProductRepository interface {
	GetSalesmanTakingOrder(custID string, empID int64) (salesman model.Salesman, err error)
	GetSalesmanCanvas(custID string, empID int64) (salesman model.SalesmanCanvas, err error)
	FindAllByCustId(dataFilter entity.ProductsQueryFilter, custId, parentCustId string, WhId int64, IsTakingOrder bool) ([]model.Product, int64, int, error)
	FindOneByProductIdAndCustId(params entity.DetailProductParams) (model.Product, error)
	// FindByEmpIdCanvas(EmpId int64, custId string) (stockWarehouse model.StockWareHouseList, err error)
}

func NewProductRepository(db *gorm.DB) *RepositoryProductImpl {
	return &RepositoryProductImpl{db}
}

func (repository *RepositoryProductImpl) GetSalesmanTakingOrder(custID string, empID int64) (salesman model.Salesman, err error) {
	err = repository.
		Select("*").
		Where("cust_id = ? AND emp_id = ?", custID, empID).
		Take(&salesman).Error

	return salesman, err
}

func (repository *RepositoryProductImpl) GetSalesmanCanvas(custID string, empID int64) (salesmanCanvas model.SalesmanCanvas, err error) {
	err = repository.
		Select("*").
		Where("cust_id = ? AND emp_id = ?", custID, empID).
		Take(&salesmanCanvas).Error

	return salesmanCanvas, err
}

func (repository *RepositoryProductImpl) FindAllByCustId(dataFilter entity.ProductsQueryFilter, custId, parentCustId string, WhId int64, IsTakingOrder bool) ([]model.Product, int64, int, error) {
	var products []model.Product
	var total int64
	var limit int
	if dataFilter.Limit == 0 {
		limit = 10
	} else {
		limit = dataFilter.Limit
	}

	var queryFilter string
	if dataFilter.Query != "" {
		queryFilter = fmt.Sprintf("and LOWER(m_product.pro_name) ilike '%%%s%%'", dataFilter.Query)
	}

	stock := "> 0"
	if IsTakingOrder && dataFilter.Mode != "canvas" {
		stock = ">= 0"
	}

	innerQuery := fmt.Sprintf(`
		SELECT 
			COALESCE(stk.wh_id, %d) as wh_id,
			m_product.cust_id, m_product.pro_id, m_product.is_active, 
			m_product.pro_code, m_product.pro_name, m_product.sup_id, 
			sbr.brand_id, m_product.sbrand1_id, 
			m_product.unit_id1, m_product.unit_id2, m_product.unit_id3, m_product.unit_id4, m_product.unit_id5,
			m_product.conv_unit2, m_product.conv_unit3, m_product.conv_unit4, m_product.conv_unit5, m_product.image_url,
			s.sup_code, s.sup_name,
			br.brand_code, br.brand_name,
			sbr.sbrand1_code, sbr.sbrand1_name, 
			mmp.limit_action, mmp.price1_minimum, mmp.price2_minimum, mmp.price3_minimum, mmp.price4_minimum, mmp.price5_minimum,
			mmp.status_manage_minimum_price as status_manage_minimum_price_name, mmp.limit_action,
			m_product.sell_price1, m_product.sell_price2, m_product.sell_price3, m_product.sell_price4, m_product.sell_price5,m_product.vat,m_product.vat_bg,m_product.vat_lg_purch,m_product.vat_lg_sell,
			COALESCE(SUM(stk.qty_in), 0)-COALESCE(SUM(stk.qty_out), 0) AS qty 
		FROM "mst"."m_product" 
		LEFT JOIN mst.m_supplier s ON s.sup_id = m_product.sup_id
		LEFT JOIN mst.m_sub_brand1 sbr ON sbr.sbrand1_id = m_product.sbrand1_id
		LEFT JOIN mst.m_brand br ON br.brand_id = sbr.brand_id
		LEFT JOIN mst.manage_minimum_price mmp on mmp.pro_id = mst.m_product.pro_id AND mmp.is_del = false and mmp.status_manage_minimum_price = 2 AND mmp.cust_id = '%s'
		LEFT JOIN inv.stock as stk ON stk.pro_id = m_product.pro_id AND stk.wh_id = %d
		WHERE m_product.cust_id = '%s' AND m_product.is_active = true AND m_product.is_del = false %s
		GROUP BY stk.wh_id,m_product.cust_id, m_product.pro_id, m_product.is_active, 
			m_product.pro_code, m_product.pro_name, m_product.sup_id, 
			sbr.brand_id, m_product.sbrand1_id, 
			m_product.unit_id1, m_product.unit_id2, m_product.unit_id3, m_product.unit_id4, m_product.unit_id5,
			m_product.conv_unit2, m_product.conv_unit3, m_product.conv_unit4, m_product.conv_unit5, m_product.image_url,
			s.sup_code, s.sup_name,
			br.brand_code, br.brand_name,
			sbr.sbrand1_code, sbr.sbrand1_name, 
			mmp.limit_action, mmp.price1_minimum, mmp.price2_minimum, mmp.price3_minimum, mmp.price4_minimum, mmp.price5_minimum,
			mmp.status_manage_minimum_price, mmp.limit_action,
			m_product.sell_price1, m_product.sell_price2, m_product.sell_price3, m_product.sell_price4, m_product.sell_price5,m_product.vat,m_product.vat_bg,m_product.vat_lg_purch,m_product.vat_lg_sell
		ORDER BY m_product.pro_code ASC
	`, WhId, custId, WhId, custId, queryFilter)

	finalQuery := fmt.Sprintf("SELECT * FROM (%s) as stock WHERE stock.qty "+stock+" LIMIT %d OFFSET %d", innerQuery, limit, (dataFilter.Page-1)*limit)

	err := repository.Raw(finalQuery).Scan(&products).Error
	if err != nil {
		return products, total, 0, err
	}

	// Count total for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as stock WHERE stock.qty "+stock+"", innerQuery)
	var count int64
	err = repository.Raw(countQuery).Scan(&count).Error
	if err != nil {
		return products, total, 0, err
	}
	total = count

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return products, total, lastPage, nil

}

func (repository *RepositoryProductImpl) FindOneByProductIdAndCustId(params entity.DetailProductParams) (model.Product, error) {
	var product model.Product

	err := repository.Table("mst.m_product p").
		Select(`m_product.cust_id, m_product.pro_id, m_product.pro_code, m_product.bar_code, m_product.pro_name,
				m_product.vat, m_product.vat_bg, m_product.vat_lg_purch, m_product.vat_lg_sell, m_product.cogs, m_product.pro_status,
				m_product.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				m_product.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				m_product.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				m_product.flavor_id, fv.flavor_code, fv.flavor_name, 
				m_product.ptype_id, pt.ptype_code, pt.ptype_name,
				m_product.psize_id, ps.psize_code, ps.psize_name,
				m_product.sup_id, su.sup_code, su.sup_name,
				m_product.principal_id, pr.principal_code, pr.principal_name,
				m_product.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				m_product.is_main_pro, m_product.sort_no, m_product.item_no, m_product.unit_id1, m_product.unit_id2, m_product.unit_id3, m_product.unit_id4, m_product.unit_id5, 
				m_product.conv_unit2, m_product.conv_unit3, m_product.conv_unit4, m_product.conv_unit5, 
				m_product.is_batch, m_product.is_exp_date, 
				m_product.weight,m_product.length, m_product.width, m_product.height, m_product.volume,
				m_product.purch_price1, m_product.sell_price1,
				m_product.purch_price2, m_product.sell_price2,
				m_product.purch_price3, m_product.sell_price3,
				m_product.purch_price4, m_product.sell_price4,
				m_product.purch_price5, m_product.sell_price5,
				m_product.weight1, m_product.length1, m_product.width1, m_product.height1, m_product.volume1, 
				m_product.weight2, m_product.length2, m_product.width2, m_product.height2, m_product.volume2, 
				m_product.weight3, m_product.length3, m_product.width3, m_product.height3, m_product.volume3, 
				m_product.weight4, m_product.length4, m_product.width4, m_product.height4, m_product.volume4, 
				m_product.weight5, m_product.length5, m_product.width5, m_product.height5, m_product.volume5,  
				m_product.parent_pro_id, 
				m_product.excise_rate, m_product.excise_tax, 
				m_product.is_active, m_product.created_by, m_product.created_at, 
				m_product.updated_by, m_product.updated_at, m_product.is_del, m_product.deleted_by, m_product.deleted_at, m_product.image_url,
				m_product.saf_stock_unit_id, m_product.saf_stock_qty, m_product.min_stock_unit_id, m_product.min_stock_qty,
				saf_unit.unit_name AS saf_stock_unit_name,
				min_unit.unit_name AS min_stock_unit_name,
				u.user_fullname AS updated_by_name,
				parent.pro_code AS parent_pro_code, parent.pro_name AS parent_pro_name `).
		Joins("LEFT JOIN sys.m_user u ON u.user_id = m_product.updated_by").
		Joins("LEFT JOIN mst.m_unit saf_unit ON saf_unit.unit_id = m_product.saf_stock_unit_id AND saf_unit.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_unit min_unit ON min_unit.unit_id = m_product.min_stock_unit_id AND min_unit.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = m_product.sbrand1_id AND sb1.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = m_product.pcat_id AND pc.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = m_product.sbrand2_id AND sb2.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_pack_size ps ON ps.psize_id = m_product.psize_id AND ps.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = m_product.ptype_id AND pt.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_flavor fv ON fv.flavor_id = m_product.flavor_id AND fv.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_principal pr ON pr.principal_id = m_product.principal_id AND pr.cust_id = ?", params.ParentCustID).
		Joins("LEFT JOIN mst.m_supplier su ON su.sup_id = m_product.sup_id").
		Joins("LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = m_product.c_pro_id").
		Joins("LEFT JOIN mst.m_product parent ON parent.pro_id = m_product.parent_pro_id").
		Where("m_product.pro_id = ? AND m_product.cust_id = ?", params.ProductId, params.ParentCustID).
		First(&product).Error

	if err != nil {
		log.Println("productRepository, FindOneByProductTypeAndCustId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
