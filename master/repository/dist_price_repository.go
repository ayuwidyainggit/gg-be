package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/jmoiron/sqlx"
)

type DistPriceRepository interface {
	// TrxBegin() (*distPriceTransaction, error)
	FindOneByDistPriceIdAndCustId(distPriceId int64, custId string) (model.DistPriceDetail, error)
	// FindOneByDiscCodeAndCustId(distPriceCode, custId string) (model.DistPrice, error)
	FindAllByCustId(dataFilter entity.DistPriceQueryFilter, custId string) (distPrice []model.DistPriceList, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.DistPriceQueryFilter, parentCustId, custId string) (distPrice []model.DistPriceLookup, total int, lastPage int, err error)
	FindAllByCustIdLookupProduct(dataFilter entity.DistPriceQueryFilter, custId string) (distPrice []model.DistPriceLookupProduct, total int, lastPage int, err error)
	Store(distPrice model.DistPrice) (int64, error)
	Update(distPriceId int64, request entity.UpdateDistPriceRequest) error
	UpdateEndDateNullByDistPriceId(distPriceId int64, request entity.UpdateDistPriceRequest) error
	Delete(custId string, distPriceId int64, deletedBy int64) error
	TrxBegin()
	TrxCommit() error
	TrxRollback() error
	CountByProIDAndCustID(proID int64, custID string) int64
	UpdateStatusByRMQ(request entity.PublishUnpublishDistPriceReq) error
}

func NewDistPriceRepository(db *sqlx.DB) DistPriceRepository {
	return &distPriceRepositoryImpl{db: db}
}

type distPriceRepositoryImpl struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func (repository *distPriceRepositoryImpl) TrxBegin() {
	repository.tx = repository.db.MustBegin()
}
func (repo *distPriceRepositoryImpl) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *distPriceRepositoryImpl) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *distPriceRepositoryImpl) FindOneByDistPriceIdAndCustId(distPriceId int64, custId string) (model.DistPriceDetail, error) {
	distPrice := model.DistPriceDetail{}
	query := `SELECT 
				d.cust_id, d.dist_price_id, d.dist_price_group_id,
				d.start_date, d.end_date, d.pro_id,
				d.unit_id1, d.unit_id2, d.unit_id3, d.unit_id4, d.unit_id5,
				d.conv_unit2, d.conv_unit3, d.conv_unit4, d.conv_unit5,
				d.purch_price1, d.purch_price2, d.purch_price3, d.purch_price4, d.purch_price5,
				d.sell_price1, d.sell_price2, d.sell_price3, d.sell_price4, d.sell_price5, 
				d.status,

				mp.pro_code, mp.pro_name, 
				un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, 
				un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, 
				un5.unit_name AS unit_name5, 
				dpg.dist_price_grp_code, dpg.dist_price_grp_name,
				d.updated_at, u.user_fullname AS updated_by_name,
				d.dist_price_id_old,
				old_d.dist_price_group_id AS dist_price_group_id_old,
				old_d.start_date AS start_date_old, old_d.end_date AS end_date_old, old_d.pro_id AS pro_id_old,
				old_d.unit_id1 AS unit_id1_old, old_d.unit_id2 AS unit_id2_old, old_d.unit_id3 AS unit_id3_old, 
				old_d.unit_id4 AS unit_id4_old, old_d.unit_id5 AS unit_id5_old,
				old_d.conv_unit2 AS conv_unit2_old, old_d.conv_unit3 AS conv_unit3_old, 
				old_d.conv_unit4 AS conv_unit4_old, old_d.conv_unit5 AS conv_unit5_old,

				old_d.purch_price1 AS purch_price1_old, old_d.purch_price2 AS purch_price2_old, 
				old_d.purch_price3 AS purch_price3_old, old_d.purch_price4 AS purch_price4_old, 
				old_d.purch_price5 AS purch_price5_old,

				old_d.sell_price1 AS sell_price1_old, old_d.sell_price2 AS sell_price2_old, 
				old_d.sell_price3 AS sell_price3_old, old_d.sell_price4 AS sell_price4_old, 
				old_d.sell_price5 AS sell_price5_old, 

				old_un1.unit_name AS unit_name1_old, old_un2.unit_name AS unit_name2_old, 
				old_un3.unit_name AS unit_name3_old, old_un4.unit_name AS unit_name4_old, 
				old_un5.unit_name AS unit_name5_old, 
				old_dpg.dist_price_grp_code AS dist_price_grp_code_old, old_dpg.dist_price_grp_name AS dist_price_grp_name_old,
				
				d.updated_at, u.user_fullname AS updated_by_name,
				d.dist_price_id_old
			  FROM mst.m_dist_price d
			  LEFT JOIN sys.m_user u ON u.user_id = d.updated_by 
			  LEFT JOIN mst.m_product mp ON mp.pro_id = d.pro_id AND mp.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = d.dist_price_group_id AND dpg.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit un1 ON un1.unit_id = d.unit_id1 AND un1.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit un2 ON un2.unit_id = d.unit_id2 AND un2.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit un3 ON un3.unit_id = d.unit_id3 AND un3.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit un4 ON un4.unit_id = d.unit_id4 AND un4.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit un5 ON un5.unit_id = d.unit_id5 AND un5.cust_id = '` + custId + `'
			  
			  LEFT JOIN mst.m_dist_price old_d ON old_d.dist_price_id = d.dist_price_id_old
			  LEFT JOIN mst.m_dist_price_group old_dpg ON old_dpg.dist_price_grp_id = old_d.dist_price_group_id AND old_dpg.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit old_un1 ON old_un1.unit_id = old_d.unit_id1 AND old_un1.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit old_un2 ON old_un2.unit_id = old_d.unit_id2 AND old_un2.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit old_un3 ON old_un3.unit_id = old_d.unit_id3 AND old_un3.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit old_un4 ON old_un4.unit_id = old_d.unit_id4 AND old_un4.cust_id = '` + custId + `'
			  LEFT JOIN mst.m_unit old_un5 ON old_un5.unit_id = old_d.unit_id5 AND old_un5.cust_id = '` + custId + `'
			  WHERE d.is_del = false 
			  AND d.dist_price_id = $1 
			  AND d.cust_id = $2`
	err := repository.db.Get(&distPrice, query, distPriceId, custId)
	if err != nil {
		log.Info("distPriceRepository, FindOneByDistPriceIdAndCustId, err:", err.Error())
		return distPrice, err
	}

	return distPrice, nil
}

func (repository *distPriceRepositoryImpl) FindAllByCustId(dataFilter entity.DistPriceQueryFilter, custId string) ([]model.DistPriceList, int, int, error) {

	distPrices := []model.DistPriceList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` d.cust_id, d.dist_price_id, d.dist_price_group_id,
					d.start_date,
					d.end_date,
					d.pro_id,
					d.unit_id1,
					d.unit_id2,
					d.unit_id3,
					d.unit_id4,
					d.unit_id5,
					d.conv_unit2,
					d.conv_unit3,
					d.conv_unit4,
					d.conv_unit5,
					d.purch_price1,
					d.purch_price2,
					d.purch_price3,
					d.purch_price4,
					d.purch_price5,
					d.sell_price1,
					d.sell_price2,
					d.sell_price3,
					d.sell_price4,
					d.sell_price5, 
					d.updated_at, 
					d.dist_price_id_old,
					p.pro_code, p.pro_name,
					u.user_fullname AS updated_by_name,
					dpg.dist_price_grp_code, dpg.dist_price_grp_name,
					d.status `
	qWhere := ` WHERE d.is_del = false 
				AND d.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND d.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND d.is_active = false `
		}
	}

	if dataFilter.DistPriceGroupId != nil {
		distPriceGrpIdStr := fmt.Sprintf("%d", *dataFilter.DistPriceGroupId)
		distPriceGrpIdInt, err := strconv.Atoi(distPriceGrpIdStr)
		if err != nil {
			return distPrices, 0, 0, err
		}
		if distPriceGrpIdInt != 0 {
			qWhere += ` AND d.dist_price_group_id = ` + fmt.Sprintf("%d", *dataFilter.DistPriceGroupId) + ` `
		}
	}

	qFrom := ` 	FROM mst.m_dist_price d
			   	LEFT JOIN mst.m_product p ON p.pro_id = d.pro_id 
				LEFT JOIN sys.m_user u ON u.user_id = d.updated_by 
				LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = d.dist_price_group_id `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Info("distPriceRepository, queryCount:", queryCount)
	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Info("distPriceRepository, count total, err:", err.Error())
		return distPrices, 0, 0, err
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
		sortBy := `dist_price_id`
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

	// log.Info("distPriceRepository, querySelect:", querySelect)
	err = repository.db.Select(&distPrices, querySelect)
	if err != nil {
		log.Info("distPriceRepository, FindAllByCustId, err:", err.Error())
		return distPrices, total, lastPage, err
	}

	return distPrices, total, lastPage, nil
}

func (repository *distPriceRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.DistPriceQueryFilter, parentCustId, custId string) ([]model.DistPriceLookup, int, int, error) {
	var strFilterByWhId string
	if dataFilter.WhId != nil {
		if *dataFilter.WhId > 0 {
			strFilterByWhId = `AND wh.wh_id = ` + strconv.Itoa(*dataFilter.WhId)
		}
	}

	distPrices := []model.DistPriceLookup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` d.dist_price_id, d.dist_price_group_id,
					d.start_date, d.end_date, d.pro_id, 
					p.pro_code, p.pro_name,
					p.unit_id1, p.unit_id2, p.unit_id3, p.unit_id4, p.unit_id5,
					p.conv_unit2, p.conv_unit3, p.conv_unit4, p.conv_unit5,
					pd.vat, pd.vat_bg, pd.vat_lg_purch, pd.vat_lg_sell, p.excise_rate, p.excise_tax,
					d.purch_price1, d.purch_price2, d.purch_price3, d.purch_price4, d.purch_price5,
					d.sell_price1, d.sell_price2, d.sell_price3, d.sell_price4, d.sell_price5, 
					d.dist_price_id_old,
					p.sup_id, sup.sup_code, sup.sup_name,
					un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, un4.unit_name AS unit_name4, un5.unit_name AS unit_name5, 
					dpg.dist_price_grp_code, dpg.dist_price_grp_name,
					
					COALESCE((SELECT SUM(COALESCE(wh.qty, 0)) 
					FROM inv.wh_stock wh
						WHERE wh.pro_id = d.pro_id
							` + strFilterByWhId + `
							AND wh.cust_id = '` + custId + `'),0) AS qty,

					COALESCE((SELECT SUM(COALESCE(wh.qty_on_order, 0)) 
						FROM inv.wh_stock wh
						WHERE wh.pro_id= d.pro_id
							` + strFilterByWhId + `
							AND wh.cust_id = '` + custId + `'),0) AS qty_on_order,

					COALESCE((SELECT SUM(COALESCE(wh.qty_on_shipping, 0)) 
						FROM inv.wh_stock wh
						WHERE wh.pro_id= d.pro_id
							` + strFilterByWhId + `
							AND wh.cust_id = '` + custId + `'),0) AS qty_on_shipping,

					COALESCE((SELECT SUM(COALESCE(wh.qty_bs, 0)) 
						FROM inv.wh_stock wh
						WHERE wh.pro_id= d.pro_id
							` + strFilterByWhId + `
							AND wh.cust_id = '` + custId + `'),0) AS qty_bs,

					COALESCE((SELECT SUM(COALESCE(wh.qty_exp, 0)) 
						FROM inv.wh_stock wh
						WHERE wh.pro_id= d.pro_id
							` + strFilterByWhId + `
							AND wh.cust_id = '` + custId + `'),0) AS qty_exp
				
				`
	qWhere := ` WHERE d.is_del = false 
				AND pd.is_active = true 
				AND d.start_date <= '` + dataFilter.WhDate + `'
				AND COALESCE(d.end_date, CURRENT_DATE) >= '` + dataFilter.WhDate + `'
				AND d.cust_id = '` + parentCustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (upper(p.pro_code) ILIKE '%` + strings.ToUpper(dataFilter.Query) + `%' 
					OR upper(p.pro_name) ILIKE '%` + strings.ToUpper(dataFilter.Query) + `%' )`
	}

	if dataFilter.DistPriceGroupId != nil {
		qWhere += ` AND d.dist_price_group_id = ` + fmt.Sprintf("%d", *dataFilter.DistPriceGroupId) + ` `
	}

	if dataFilter.SupId != nil {
		qWhere += ` AND p.sup_id = ` + fmt.Sprintf("%d", *dataFilter.SupId) + ` `
	}

	if dataFilter.ProId != nil {
		if *dataFilter.ProId > 0 {
			qWhere += ` AND d.pro_id  = ` + fmt.Sprintf("%d", *dataFilter.ProId) + ` `
		}
	}

	qFrom := ` 	FROM mst.m_dist_price d
				LEFT JOIN mst.m_product p ON p.pro_id = d.pro_id AND p.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_product_dist pd ON pd.pro_id = p.pro_id AND pd.cust_id = '` + custId + `'
				LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = d.dist_price_group_id AND dpg.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_supplier sup ON sup.sup_id = p.sup_id AND sup.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un1 ON un1.unit_id = d.unit_id1 AND un1.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un2 ON un2.unit_id = d.unit_id2 AND un2.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un3 ON un3.unit_id = d.unit_id3 AND un3.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un4 ON un4.unit_id = d.unit_id4 AND un4.cust_id = '` + parentCustId + `'
				LEFT JOIN mst.m_unit un5 ON un5.unit_id = d.unit_id5 AND un5.cust_id = '` + parentCustId + `' `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Info("distPriceRepository, count total, err:", err.Error())
		return distPrices, 0, 0, err
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
		sortBy := `d.dist_price_id`
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

	err = repository.db.Select(&distPrices, querySelect)
	if err != nil {
		log.Info("distPriceRepository, FindAllByCustId, err:", err.Error())
		return distPrices, total, lastPage, err
	}

	return distPrices, total, lastPage, nil
}

func (repository *distPriceRepositoryImpl) FindAllByCustIdLookupProduct(dataFilter entity.DistPriceQueryFilter, custId string) ([]model.DistPriceLookupProduct, int, int, error) {

	distPrices := []model.DistPriceLookupProduct{}
	selectCount := ` COUNT(*) AS total `
	selectField := `dp.dist_price_id, dp.pro_id, dp.dist_price_id_old, dp.dist_price_group_id,
					dp.unit_id1, dp.unit_id2, dp.unit_id3, dp.unit_id4, dp.unit_id5,
					dp.conv_unit2, dp.conv_unit3, dp.conv_unit4, dp.conv_unit5, 
					dp.purch_price1, dp.purch_price2, dp.purch_price3, dp.purch_price4, dp.purch_price5,
					dp.sell_price1, dp.sell_price2, dp.sell_price3, dp.sell_price4, dp.sell_price5,
					p.pro_code, p.pro_name,
					dpg.dist_price_grp_code, dpg.dist_price_grp_name,
					un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3, 
					un4.unit_name AS unit_name4, un5.unit_name AS unit_name5 `
	qWhere := ` WHERE dp.is_del = false AND dp.end_date IS NULL
				AND dp.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.DistPriceGroupId != nil && *dataFilter.DistPriceGroupId > 0 {
		qWhere += ` AND dp.dist_price_group_id = ` + strconv.Itoa(*dataFilter.DistPriceGroupId) + ` `
	}

	qFrom := ` 	FROM mst.m_dist_price dp
				LEFT JOIN mst.m_product p ON p.pro_id = dp.pro_id 
				LEFT JOIN mst.m_dist_price_group dpg ON dpg.dist_price_grp_id = dp.dist_price_group_id AND dpg.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit un1 ON un1.unit_id = dp.unit_id1 AND un1.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit un2 ON un2.unit_id = dp.unit_id2 AND un2.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit un3 ON un3.unit_id = dp.unit_id3 AND un3.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit un4 ON un4.unit_id = dp.unit_id4 AND un4.cust_id = '` + custId + `'
				LEFT JOIN mst.m_unit un5 ON un5.unit_id = dp.unit_id5 AND un5.cust_id = '` + custId + `' `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Info("distPriceRepository, count total, err:", err.Error())
		return distPrices, 0, 0, err
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
		sortBy := `dp.dist_price_id`
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

	err = repository.db.Select(&distPrices, querySelect)
	if err != nil {
		log.Info("distPriceRepository, FindAllByCustIdLookupProduct, err:", err.Error())
		return distPrices, total, lastPage, err
	}

	return distPrices, total, lastPage, nil
}

func (repository *distPriceRepositoryImpl) Store(distPrice model.DistPrice) (int64, error) {
	query :=
		`INSERT INTO mst.m_dist_price(
			cust_id, dist_price_group_id, start_date, end_date,
			pro_id, unit_id1, unit_id2, unit_id3,
			unit_id4, unit_id5, conv_unit2, conv_unit3,
			conv_unit4, conv_unit5, purch_price1, purch_price2,
			purch_price3, purch_price4, purch_price5, sell_price1,
			sell_price2, sell_price3, sell_price4, sell_price5,
			created_by, created_at, updated_by, updated_at, dist_price_id_old)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28, $29
		) RETURNING dist_price_id;`
	// lastInsertId := distPrice.DistPriceId
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		distPrice.CustId, distPrice.DistPriceGroupId, distPrice.StartDate, distPrice.EndDate,
		distPrice.ProId, distPrice.UnitId1, distPrice.UnitId2, distPrice.UnitId3,
		distPrice.UnitId4, distPrice.UnitId5, distPrice.ConvUnit2, distPrice.ConvUnit3,
		distPrice.ConvUnit4, distPrice.ConvUnit5, distPrice.PurchPrice1, distPrice.PurchPrice2,
		distPrice.PurchPrice3, distPrice.PurchPrice4, distPrice.PurchPrice5, distPrice.SellPrice1,
		distPrice.SellPrice2, distPrice.SellPrice3, distPrice.SellPrice4, distPrice.SellPrice5,
		distPrice.CreatedBy, distPrice.CreatedAt, distPrice.UpdatedBy, distPrice.UpdatedAt, distPrice.DistPriceIdOld).Scan(&lastInsertId)
	if err != nil {
		log.Info("DistPriceRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *distPriceRepositoryImpl) Update(distPriceId int64, request entity.UpdateDistPriceRequest) error {
	var (
		r            model.DistPriceUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("distPriceRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_dist_price
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND dist_price_id = :dist_price_id;`

	// log.Info("distPriceRepository, Update, query:", query)

	sqlPatch.Args["dist_price_id"] = distPriceId
	sqlPatch.Args["cust_id"] = request.CustId
	// log.Info("repository, Update, sqlPatch:", structs.StructToJson(sqlPatch))

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Info("distPriceRepository, Update, err:", err.Error())
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

func (repository *distPriceRepositoryImpl) UpdateEndDateNullByDistPriceId(distPriceId int64, request entity.UpdateDistPriceRequest) error {
	var (
		r     model.DistPriceUpdate
		nRows int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	query := `UPDATE mst.m_dist_price
			  SET 
			  	end_date = NULL,
			  	updated_by = :updated_by,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND dist_price_id = :dist_price_id;`

	// log.Info("distPriceRepository, Update, query:", query)

	args := map[string]interface{}{
		"updated_by":    request.UpdatedBy,
		"dist_price_id": distPriceId,
		"cust_id":       request.CustId,
	}
	// log.Info("repository, UpdateEndDateNullByDistPriceId, sqlPatch:", structs.StructToJson(sqlPatch))

	result, err := repository.tx.NamedExec(query, args)
	if err != nil {
		log.Info("distPriceRepository, UpdateEndDateNullByDistPriceId, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return err
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *distPriceRepositoryImpl) Delete(custId string, distPriceId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_dist_price
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND dist_price_id = :dist_price_id;`

	wMap := map[string]interface{}{
		"cust_id":       custId,
		"dist_price_id": distPriceId,
		"deleted_by":    deletedBy,
	}

	result, err := repository.tx.NamedExec(query, wMap)
	if err != nil {
		log.Info("DistPriceRepository, Delete, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return err
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *distPriceRepositoryImpl) CountByProIDAndCustID(proID int64, custID string) (count int64) {
	query := `SELECT count(*) AS total FROM mst.m_product WHERE is_del = false AND pro_id = $1 AND cust_id = $2`
	err := repository.db.QueryRow(query, proID, custID).Scan(&count)
	if err != nil {
		log.Info("distPriceRepository, CountByProIDAndCustID, err:", err.Error())
		return count
	}
	return count
}

func (repository *distPriceRepositoryImpl) UpdateStatusByRMQ(request entity.PublishUnpublishDistPriceReq) error {
	var (
		r            model.DistPricePublish
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_dist_price ` +
		`SET ` + sqlSetFields + `, ` +
		`updated_at = CURRENT_TIMESTAMP ` +
		`WHERE cust_id = :cust_id ` +
		`AND dist_price_id = :dist_price_id;`

	sqlPatch.Args["dist_price_id"] = request.DistPriceID
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.db.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error(err.Error())
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
