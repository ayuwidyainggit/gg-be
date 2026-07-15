package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ProductDistRepository interface {
	FindOneByProductIdAndCustId(productId int64, custId, parentCustId string) (model.ProductDist, error)
	FindOneByProductCodeAndCustId(productCode string, custId, parentCustId string) (model.ProductDist, error)
	FindAllByCustId(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (consPro []model.ProductDist, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (consPro []model.ProductDist, total int, lastPage int, err error)
	FindAllByCustIdSearch(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) (consPro []model.ProductDist, total int, lastPage int, err error)
	Store(product model.ProductDist) (int64, error)
	Update(productId int64, request entity.UpdateProductDistRequest) error
	Delete(custId string, productId int64, deletedBy int64) error
}

func NewProductDistRepository(db *sqlx.DB) ProductDistRepository {
	return &productDistRepositoryImpl{db}
}

type productDistRepositoryImpl struct {
	*sqlx.DB
}

func (repository *productDistRepositoryImpl) FindOneByProductIdAndCustId(productId int64, custId, parentCustId string) (model.ProductDist, error) {
	product := model.ProductDist{}
	query := `SELECT 
				pd.cust_id, pd.pro_id, pd.is_alloc, pd.min_stock, pd.min_stock_str, 
				pd.safety_stock, pd.safety_stock_str, pd.po_formula, pd.is_active, 
				pd.is_new_pro, pd.s_mweek1, pd.s_mweek2, pd.cogs,
				pd.vat, pd.vat_bg, pd.vat_lg_purch, pd.vat_lg_sell, 
				p.pro_code, p.bar_code, p.pro_name, 
				p.pcat_id, pc.pcat_code, pc.pcat_name,
				br.brand_id, br.brand_code, br.brand_name, 
				br.pl_id, pl.pl_code, pl.pl_name,
				p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
				p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
				p.flavor_id, fv.flavor_code, fv.flavor_name, 
				p.ptype_id, pt.ptype_code, pt.ptype_name,
				p.psize_id, ps.psize_code, ps.psize_name,
				p.sup_id, su.sup_code, su.sup_name,
				p.principal_id, pr.principal_code, pr.principal_name,
				p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
				p.is_main_pro, p.sort_no, p.item_no, 
				p.weight, p.is_batch, p.is_exp_date, p.length, p.width, p.height, p.volume, 
				p.parent_pro_id, p.excise_rate, p.excise_tax, p.image_url,
				p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
				un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5, 
				p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5, 
				pd.updated_by, pd.updated_at, 
				u.user_fullname AS updated_by_name
			FROM mst.m_product_dist pd
			LEFT JOIN mst.m_product p ON p.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `' 
			LEFT JOIN sys.m_user u ON u.user_id = pd.updated_by AND u.cust_id = '` + custId + `' 
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustId + `' 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un4 ON un4.unit_id = p.unit_id4 AND un4.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un5 ON un5.unit_id = p.unit_id5 AND un5.cust_id = '` + parentCustId + `'
			WHERE pd.pro_id = $1 
			AND pd.cust_id = $2`
	err := repository.Get(&product, query, productId, custId)
	if err != nil {
		log.Println("productDistRepository, FindOneByProductIdAndCustId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productDistRepositoryImpl) FindOneByProductCodeAndCustId(productCode string, custId, parentCustId string) (model.ProductDist, error) {
	product := model.ProductDist{}
	query := `SELECT 
				p.cust_id, p.pro_id, p.pro_code, p.bar_code, p.pro_name, 
				p.pcat_id, p.sbrand1_id, p.sbrand2_id, p.flavor_id, p.ptype_id, 
				p.psize_id, p.sup_id, p.principal_id, p.c_pro_id, p.is_main_pro, p.is_alloc, 
				p.s_mweek1, p.s_mweek2, p.sort_no, p.item_no, 
				dp.unit_id1, dp.unit_id2, dp.unit_id3, dp.unit_id4, dp.unit_id5, 
				un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5, 
				dp.conv_unit2, dp.conv_unit3, dp.conv_unit4, dp.conv_unit5, 
				dp.purch_price1, dp.purch_price2, dp.purch_price3, dp.purch_price4, dp.purch_price5, 
				dp.sell_price1, dp.sell_price2, dp.sell_price3, dp.sell_price4, dp.sell_price5, 
				p.weight, p.is_batch, p.is_exp_date,  
				p.length, p.width, p.height, p.volume, 
				pd.min_stock, pd.min_stock_str, pd.safety_stock, pd.safety_stock_str, 
				pd.po_formula, p.parent_pro_id, pd.is_new_pro, 
				pd.vat, pd.vat_bg, pd.vat_lg_purch, pd.vat_lg_sell, p.excise_rate, p.excise_tax, 
				p.is_active, p.created_by, p.created_at, 
				p.updated_by, p.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
				u.user_fullname AS updated_by_name
			FROM mst.m_product_dist pd
			LEFT JOIN mst.m_product p ON p.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `' 
			LEFT JOIN mst.m_dist_price dp ON dp.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `' 
			LEFT JOIN sys.m_user u ON u.user_id = pd.updated_by AND u.cust_id = '` + custId + `' 
			LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustId + `' 
			LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un1 ON un1.unit_id = dp.unit_id1 AND un1.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un2 ON un2.unit_id = dp.unit_id2 AND un2.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un3 ON un3.unit_id = dp.unit_id3 AND un3.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un4 ON un4.unit_id = dp.unit_id4 AND un4.cust_id = '` + parentCustId + `'
			LEFT JOIN mst.m_unit un5 ON un5.unit_id = dp.unit_id5 AND un5.cust_id = '` + parentCustId + `'
			WHERE p.pro_code = $1 AND pd.cust_id = $2`
	err := repository.Get(&product, query, productCode, custId)
	if err != nil {
		log.Println("productDistRepository, FindOneByProductCodeAndCustId, err:", err.Error())
		return product, err
	}

	return product, nil
}

func (repository *productDistRepositoryImpl) FindAllByCustId(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) ([]model.ProductDist, int, int, error) {

	products := []model.ProductDist{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pd.cust_id, pd.pro_id, pd.is_alloc, pd.min_stock, pd.min_stock_str, 
					pd.safety_stock, pd.safety_stock_str, pd.po_formula, 
					pd.is_new_pro, pd.s_mweek1, pd.s_mweek2, pd.cogs,
					pd.vat, pd.vat_bg, pd.vat_lg_purch, pd.vat_lg_sell, 
					p.parent_pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro, p.sort_no, p.item_no, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
					un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5,  
					p.weight,
					p.is_batch, p.is_exp_date, 
					p.length, p.width, p.height, p.volume, 
					p.excise_rate, p.excise_tax, 
					pd.is_active, p.created_by, p.created_at, 
					pd.updated_by, pd.updated_at, p.is_del, p.deleted_by, p.deleted_at, p.image_url,
					u.user_fullname AS updated_by_name `
	qWhere := ` LEFT JOIN mst.m_product p ON p.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `' 
				LEFT JOIN sys.m_user u ON u.user_id = pd.updated_by AND u.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un4 ON un4.unit_id = p.unit_id4 AND un4.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un5 ON un5.unit_id = p.unit_id5 AND un5.cust_id = '` + parentCustId + `'
				WHERE pd.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND pd.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND pd.is_active = false `
		}
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product_dist pd `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Println("productDistRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productDistRepository, count total, err:", err.Error())
		return products, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `p.pro_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("productDistRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Println("productDistRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productDistRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) ([]model.ProductDist, int, int, error) {

	products := []model.ProductDist{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pd.cust_id, pd.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.pcat_id, pc.pcat_code, pc.pcat_name,
					br.brand_id, br.brand_code, br.brand_name, 
					br.pl_id, pl.pl_code, pl.pl_name,
					p.sbrand1_id, sb1.sbrand1_code, sb1.sbrand1_name,
					p.sbrand2_id, sb2.sbrand2_code, sb2.sbrand2_name,
					p.flavor_id, fv.flavor_code, fv.flavor_name, 
					p.ptype_id, pt.ptype_code, pt.ptype_name,
					p.psize_id, ps.psize_code, ps.psize_name,
					p.sup_id, su.sup_code, su.sup_name,
					p.principal_id, pr.principal_code, pr.principal_name,
					p.c_pro_id, cp.c_pro_code, cp.c_pro_name,
					p.is_main_pro, pd.is_alloc, p.image_url,
					p.weight,p.is_batch, p.is_exp_date, p.sort_no, p.item_no, 
					p.length, p.width, p.height, p.volume, 
					p.parent_pro_id, p.excise_rate, p.excise_tax, 
					pd.is_active, pd.s_mweek1, pd.s_mweek2, 
					pd.vat, pd.vat_bg, pd.vat_lg_purch, pd.vat_lg_sell, 
					pd.min_stock, pd.min_stock_str, pd.safety_stock, pd.safety_stock_str, 
					pd.po_formula, pd.is_new_pro, pd.updated_by, pd.updated_at, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
					un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5,  
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5,  
					u.user_fullname AS updated_by_name `
	qWhere := ` LEFT JOIN mst.m_product p ON p.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `'
				LEFT JOIN sys.m_user u ON u.user_id = pd.updated_by AND u.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un4 ON un4.unit_id = p.unit_id4 AND un4.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un5 ON un5.unit_id = p.unit_id5 AND un5.cust_id = '` + parentCustId + `'
				WHERE pd.is_active = true 
				AND pd.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	// if dataFilter.DistPriceGroupId > 0 {
	// 	qWhere += ` AND dp.dist_price_group_id = ` + strconv.Itoa(dataFilter.DistPriceGroupId)
	// }

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product_dist pd `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Println("productDistRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productDistRepository, count total, err:", err.Error())
		return products, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `pd.pro_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("productDistRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Println("productDistRepository, FindAllByCustId, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productDistRepositoryImpl) FindAllByCustIdSearch(dataFilter entity.ProductDistQueryFilter, custId, parentCustId string) ([]model.ProductDist, int, int, error) {

	products := []model.ProductDist{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` pd.cust_id, pd.pro_id, p.pro_code, p.bar_code, p.pro_name, 
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5, 
					un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5, 
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5 `
	qWhere := ` LEFT JOIN mst.m_product p ON p.pro_id = pd.pro_id AND p.cust_id = '` + parentCustId + `' 
				LEFT JOIN sys.m_user u ON u.user_id = pd.updated_by AND u.cust_id = '` + custId + `' 
				LEFT JOIN mst.m_sub_brand1 sb1 ON sb1.sbrand1_id = p.sbrand1_id AND sb1.cust_id = '` + parentCustId + `' 
				LEFT JOIN mst.m_brand br ON br.brand_id = sb1.brand_id AND br.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_line pl ON pl.pl_id = br.pl_id AND pl.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_cat pc ON pc.pcat_id = p.pcat_id AND pc.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_sub_brand2 sb2 ON sb2.sbrand2_id = p.sbrand2_id AND sb2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_size ps ON ps.psize_id = p.psize_id AND ps.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_pack_type pt ON pt.ptype_id = p.ptype_id AND pt.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_flavor fv ON fv.flavor_id = p.flavor_id AND fv.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_principal pr ON pr.principal_id = p.principal_id AND pr.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_supplier su ON su.sup_id = p.sup_id AND su.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_cons_product cp ON cp.c_pro_id = p.c_pro_id AND cp.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un1 ON un1.unit_id = p.unit_id1 AND un1.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un2 ON un2.unit_id = p.unit_id2 AND un2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un3 ON un3.unit_id = p.unit_id3 AND un3.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un4 ON un4.unit_id = p.unit_id4 AND un4.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un5 ON un5.unit_id = p.unit_id5 AND un5.cust_id = '` + parentCustId + `'
				WHERE pd.is_active = true 
				AND pd.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.SupCode != "" {
		qWhere += ` AND su.sup_code = '` + dataFilter.SupCode + `' `
	}

	tempPCatId := []interface{}{}
	for _, xx := range dataFilter.PcatId {
		tempPCatId = append(tempPCatId, ""+strconv.Itoa(xx)+"")
	}
	varPCatId := make([]string, len(tempPCatId))
	for i, v := range tempPCatId {
		varPCatId[i] = fmt.Sprint(v)
	}
	var PCatId = strings.Join(varPCatId, ",")

	if len(dataFilter.PcatId) > 0 {
		qWhere += ` AND p.pcat_id in (` + PCatId + `) `
	}

	queryFrom := ` FROM mst.m_product_dist pd `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Println("productDistRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("productDistRepository, FindAllByCustIdSearch, err:", err.Error())
		return products, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `pd.pro_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	// log.Println("productDistRepository, querySelect:", querySelect)
	err = repository.Select(&products, querySelect)
	if err != nil {
		log.Println("productDistRepository, FindAllByCustIdSearch, err:", err.Error())
		return products, total, lastPage, err
	}

	return products, total, lastPage, nil
}

func (repository *productDistRepositoryImpl) Store(product model.ProductDist) (int64, error) {
	query :=
		`INSERT INTO mst.m_product_dist(
			cust_id, pro_code, bar_code, 
			pro_name, pcat_id, 
			sbrand1_id, sbrand2_id, flavor_id, 
			ptype_id, psize_id, sup_id, principal_id, 
			c_pro_id, is_main_pro, is_alloc, 
			s_mweek1, s_mweek2, sort_no, item_no, 
			unit_id1, unit_id2, unit_id3, unit_id4, unit_id5, 
			conv_unit2, conv_unit3, conv_unit4, conv_unit5, 
			purch_price, 
			sell_price1, 
			sell_price2, sell_price3, 
			sell_price4, sell_price5, 
			margin,weight, 
			is_batch, is_exp_date, 
			length, width, height, volume, 
			min_stock, min_stock_str, 
			safety_stock, safety_stock_str, 
			po_formula, parent_pro_id, is_new_pro, 
			vat, vat_bg, vat_lg_purch, vat_lg_sell, 
			excise_rate, excise_tax, is_active, 
			created_by, created_at, 
			updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			:cust_id, :pro_code, :bar_code, 
			:pro_name, :pcat_id, 
			:sbrand1_id, :sbrand2_id, :flavor_id, 
			:ptype_id, :psize_id, :sup_id, :principal_id, 
			:c_pro_id, :is_main_pro, :is_alloc, 
			:s_mweek1, :s_mweek2, :sort_no, :item_no, 
			:unit_id1, :unit_id2, :unit_id3, :unit_id4, :unit_id5, 
			:conv_unit2, :conv_unit3, :conv_unit4, :conv_unit5, 
			:purch_price, 
			:sell_price1, 
			:sell_price2, :sell_price3, 
			:sell_price4, :sell_price5, 
			:margin, :weight, 
			:is_batch, :is_exp_date, 
			:length, :width, :height, :volume, 
			:min_stock, :min_stock_str, 
			:safety_stock, :safety_stock_str, 
			:po_formula, :parent_pro_id, :is_new_pro, 
			:vat, :vat_bg, :vat_lg_purch, :vat_lg_sell, 
			:excise_rate, :excise_tax, :is_active, 
			:created_by, :created_at, 
			:updated_by, :updated_at, 
			:is_del, :deleted_by, :deleted_at
		) RETURNING pro_id;`
	// lastInsertId := product.ProductId
	_, err := repository.NamedExec(query, product) // .Scan(&lastInsertId)
	if err != nil {
		log.Println("productDistRepository, Store, err:", err.Error())
		return product.ProductId, err
	}
	return product.ProductId, nil
}

func (repository *productDistRepositoryImpl) Update(productId int64, request entity.UpdateProductDistRequest) error {
	var (
		r            model.ProductDistUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("productDistRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_product_dist
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id 
			  AND pro_id = :pro_id_old;`
	// log.Println("productDistRepository, Update, query:", query)

	sqlPatch.Args["pro_id_old"] = productId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("productDistRepository, Update, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *productDistRepositoryImpl) Delete(custId string, productId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_product_dist
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND pro_id = :pro_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"pro_id":     productId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ProductDistRepository, Delete, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
