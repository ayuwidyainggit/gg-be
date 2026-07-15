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
	"time"

	"github.com/jmoiron/sqlx"
)

type SalesmanRepository interface {
	FindOneParentCustId(distCustId string) (model.MCustomer, error)
	FindOneByEmpIdAndCustId(params entity.DetailSalesmanParams) (model.SalesmanList, error)
	FindSalesmanNameByEmpID(empID int64) (string, error)
	FindAllByCustId(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) (consPro []model.SalesmanList, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) (consPro []model.SalesmanList, total int, lastPage int, err error)
	Store(salesman model.Salesman) (int64, error)
	StoreWarehouse(warehouse model.Warehouse) (int, error)
	StoreSalesmanCanvas(salesman model.SalesmanCanvas) (int64, error)
	FindOneByWarehouseCodeAndCustId(warehouseCode string, custId string) (model.Warehouse, error)
	FindOneSalesmanCanvasByEmpIdAndCustId(empId int64, custId string) (model.Salesman, error)

	Update(salesmanId int64, request entity.UpdateSalesmanRequest) error
	UpdateCanvas(salesmanId int64, request entity.UpdateSalesmanRequest) error
	UpdateIsTakingOrder(salesmanId int64, empId string) error
	Delete(custId string, salesmanId int64, deletedBy int64) error
	FindDetailByIdAndCustId(params entity.DetailSalesmanParams) ([]model.SalesmanDetailRead, error)
	FindDetailById(empId int64, custId string) ([]model.SalesmanDetailRead, error)
	StoreDetail(salesmanDetail model.SalesmanDetail) error
	DeleteDetailNotInIDs(empId int64, MSalesmanDetID []int64) error
	DeleteDetails(empId int64, custId string) error
	// DeleteSalesmanCanvas(empId int64, custId string) error
	UpdateDetail(empId int64, request model.SalesmanDetail) error
	UpdateIsActive(empId int64, custId string) error
	UpdateDeActive(empId int64, custId string) error

	CheckDate()

	TrxBegin()
	TrxCommit() error
	TrxRollback() error
}

func NewSalesmanRepository(db *sqlx.DB) SalesmanRepository {
	return &salesmanRepositoryImpl{db: db}
}

type salesmanRepositoryImpl struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func buildSalesmanEmployeeJoin() string {
	return `LEFT JOIN mst.m_employee emp ON emp.emp_id = s.emp_id AND emp.cust_id = s.cust_id`
}

func buildSalesmanSalesTeamJoin() string {
	return `LEFT JOIN mst.m_sales_team st ON st.sales_team_id = s.sales_team_id AND st.cust_id = s.cust_id`
}

func buildSalesmanWarehouseJoin(alias, warehouseSource string) string {
	return `LEFT JOIN mst.m_warehouse ` + alias + ` ON ` + alias + `.wh_id = ` + warehouseSource + ` AND ` + alias + `.cust_id = s.cust_id`
}

func buildSalesmanCanvasJoin() string {
	return `LEFT JOIN mst.m_salesman_canvas msc on msc.emp_id = s.emp_id and msc.cust_id = s.cust_id`
}

func buildSalesmanVehicleJoin() string {
	return `LEFT JOIN mst.m_vehicle mv on mv.vehicle_id = msc.vehicle_id and mv.cust_id = s.cust_id`
}

func buildSalesmanDriverJoin() string {
	return `LEFT JOIN mst.m_employee me on me.emp_id = mv.driver_id and me.cust_id = s.cust_id`
}

func buildSalesmanCustIDInCondition(custIds []string) string {
	quoted := make([]string, 0, len(custIds))
	seen := make(map[string]struct{})

	for _, id := range custIds {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		quoted = append(quoted, "'"+strings.ReplaceAll(id, "'", "''")+"'")
	}

	if len(quoted) == 0 {
		return ` s.cust_id = '' `
	}
	if len(quoted) == 1 {
		return ` s.cust_id = ` + quoted[0] + ` `
	}

	return ` s.cust_id IN (` + strings.Join(quoted, ",") + `) `
}

func buildSalesmanCustScopeCondition(custIds []string, distributorIDs []int, parentCustId, custId string) string {
	if len(custIds) > 0 {
		return buildSalesmanCustIDInCondition(custIds)
	}

	if len(distributorIDs) == 0 {
		return ` s.cust_id = '` + custId + `' `
	}

	scopeParentCustID := strings.TrimSpace(parentCustId)
	if scopeParentCustID == "" {
		scopeParentCustID = strings.TrimSpace(custId)
	}

	includePrincipalScope := false
	stringIDs := make([]string, 0, len(distributorIDs))
	seen := make(map[int]struct{})
	for _, distributorID := range distributorIDs {
		if _, exists := seen[distributorID]; exists {
			continue
		}
		seen[distributorID] = struct{}{}

		if distributorID == 0 {
			includePrincipalScope = true
			continue
		}

		if distributorID <= 0 {
			continue
		}

		stringIDs = append(stringIDs, strconv.Itoa(distributorID))
	}

	conditions := make([]string, 0, 2)
	if includePrincipalScope {
		conditions = append(conditions, `s.cust_id = '`+scopeParentCustID+`'`)
	}
	if len(stringIDs) > 0 {
		conditions = append(conditions, `s.cust_id IN (
		SELECT mc.cust_id
		FROM smc.m_customer mc
		WHERE mc.parent_cust_id = '`+scopeParentCustID+`'
		AND mc.distributor_id IN (`+strings.Join(stringIDs, ",")+`)
	)`) 
	}

	if len(conditions) == 0 {
		return ` s.cust_id = '` + custId + `' `
	}

	return `( ` + strings.Join(conditions, ` OR `) + ` )`
}

func (repository *salesmanRepositoryImpl) TrxBegin() {
	repository.tx = repository.db.MustBegin()
}
func (repo *salesmanRepositoryImpl) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *salesmanRepositoryImpl) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *salesmanRepositoryImpl) FindOneParentCustId(distCustId string) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE cust_id = $1`
	err := repository.db.Get(&mCustomer, query, distCustId)
	if err != nil {
		log.Println("salesTeamRepository, FindOneParentCustId, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func (repository *salesmanRepositoryImpl) FindOneByWarehouseCodeAndCustId(warehouseCode string, custId string) (model.Warehouse, error) {
	warehouse := model.Warehouse{}
	query := `SELECT 
				cust_id, wh_id, wh_code, stock_type,
				wh_name, is_active, 
				created_by, created_at, updated_by, 
				updated_at, is_del, deleted_by, deleted_at
			  FROM mst.m_warehouse 
			  WHERE cust_id = $2  
			  AND wh_code = $1
			  AND is_del = false`
	err := repository.db.Get(&warehouse, query, warehouseCode, custId)
	if err != nil {
		log.Println("warehouseRepository, FindOneByWarehouseCodeAndCustId, err:", err.Error())
		return warehouse, err
	}

	return warehouse, nil
}

func (repository *salesmanRepositoryImpl) FindOneSalesmanCanvasByEmpIdAndCustId(empId int64, custId string) (model.Salesman, error) {

	salesman := model.Salesman{}
	query := `SELECT 
				cust_id, emp_id
			  FROM mst.m_salesman_canvas 
			  WHERE cust_id = $2  
			  AND emp_id = $1`
	err := repository.db.Get(&salesman, query, empId, custId)
	if err != nil {
		log.Println("salesmanRepositoryImpl, FindOneSalesmanByEmpIdAndCustId, err:", err.Error())
		return salesman, err
	}

	return salesman, nil
}

func (repository *salesmanRepositoryImpl) FindOneByEmpIdAndCustId(params entity.DetailSalesmanParams) (model.SalesmanList, error) {
	salesman := model.SalesmanList{}
	query := `SELECT
			s.cust_id, s.emp_id, s.sales_name, s.sales_team_id, s.opr_type, s.is_bonus_rep, s.trans_date, s.wh_id, s.inc_grp_id, ig.inc_grp_name,
			s.official_id, of.official_name, s.sale_system, s.sm_is_transfer, s.sm_valid_route, s.sm_geoloc_valid, s.sm_radius, emp.image_url,
			s.sm_is_barcode, s.sm_is_photo_profile, COALESCE(s.sm_password, '') AS sm_password, s.is_active, s.updated_at, u.user_fullname AS updated_by_name,
			COALESCE(emp.emp_code, '') AS emp_code, COALESCE(emp.emp_name, '') AS emp_name, emp.last_education, emp.phone_no, emp.address, emp.email, wh.wh_code, wh.wh_name,
			of.official_type, ofh.hierarchy_code, st.sales_team_code, st.sales_team_name, msc.opr_type as opr_type_canvas, s.job_type,s.tax_option,s.start_date,s.end_date,s.allow_input_price,
			CASE WHEN msc.is_active IS NULL THEN false ELSE msc.is_active END AS is_active_canvas, 
			wh2.wh_name as wh_name_canvas,  mv.vehicle_id, mv.vehicle_desc as vehicle_name, me.emp_name as driver_name, 
			CASE WHEN s.is_taking_order IS NULL THEN false ELSE s.is_taking_order END AS is_taking_order
			FROM mst.m_salesman s
			LEFT JOIN sys.m_user u ON u.user_id = s.updated_by 
			` + buildSalesmanEmployeeJoin() + `
			` + buildSalesmanWarehouseJoin("wh", "s.wh_id") + `
			LEFT JOIN mst.m_official of ON of.official_id = s.official_id AND of.cust_id = '` + params.ParentCustId + `'
			LEFT JOIN mst.m_inc_group ig ON ig.inc_grp_id = s.inc_grp_id AND ig.cust_id = '` + params.ParentCustId + `'
			LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type AND ofh.cust_id = '` + params.ParentCustId + `'
			` + buildSalesmanSalesTeamJoin() + `
			` + buildSalesmanCanvasJoin() + `
			` + buildSalesmanWarehouseJoin("wh2", "msc.wh_id") + `
			` + buildSalesmanVehicleJoin() + `
			` + buildSalesmanDriverJoin() + `
			WHERE s.emp_id = $1 AND s.is_del = false AND s.deleted_at IS NULL
			AND s.cust_id = $2`
	err := repository.db.Get(&salesman, query, params.EmpId, params.CustId)
	if err != nil {
		log.Println("salesmanRepository, FindOneByEmpIdAndCustId, err:", err.Error())
		return salesman, err
	}

	return salesman, nil
}

func (repository *salesmanRepositoryImpl) FindSalesmanNameByEmpID(empID int64) (string, error) {
	var salesmanName string
	query := `SELECT sales_name
		FROM mst.m_salesman
		WHERE emp_id = $1 AND is_del = false AND deleted_at IS NULL
		ORDER BY updated_at DESC NULLS LAST, created_at DESC NULLS LAST
		LIMIT 1`
	err := repository.db.Get(&salesmanName, query, empID)
	if err != nil {
		return "", err
	}

	return salesmanName, nil
}

func (repo *salesmanRepositoryImpl) FindDetailByIdAndCustId(params entity.DetailSalesmanParams) ([]model.SalesmanDetailRead, error) {
	Details := []model.SalesmanDetailRead{}
	query := `SELECT 
	 CASE 
        WHEN pl.pl_id IS NOT NULL THEN pl.pl_id
        WHEN b.pl_id IS NOT NULL THEN b.pl_id
        WHEN sb.sbrand1_id IS NOT NULL THEN b2.pl_id
    END AS pl_id,
	CASE 
		WHEN pl.pl_id is not null THEN pl.pl_id
		WHEN b.brand_id is not null THEN b.brand_id
		WHEN sb.sbrand1_id is not null THEN sb.sbrand1_id 
	END as ref_id,
	CASE 
		WHEN pl.pl_code is not null THEN pl.pl_code
		WHEN b.brand_code is not null THEN b.brand_code
		WHEN sb.sbrand1_code is not null THEN sb.sbrand1_code 
	END as ref_code,
	CASE 
		WHEN pl.pl_name is not null THEN pl.pl_name
		WHEN b.brand_name is not null THEN b.brand_name
		WHEN sb.sbrand1_name is not null THEN sb.sbrand1_name
	END as ref_name,
	spt.group_type,
	spt.m_salesman_product_type_id
	FROM mst.m_salesman_product_type spt
	LEFT JOIN mst.m_product_line pl ON pl.pl_id = spt.ref_id AND spt.group_type=1
	LEFT JOIN mst.m_brand b ON b.brand_id = spt.ref_id AND spt.group_type=2
	LEFT JOIN mst.m_sub_brand1 sb ON sb.sbrand1_id = spt.ref_id AND spt.group_type=3
	LEFT JOIN mst.m_brand b2 ON b2.brand_id = sb.brand_id
	WHERE spt.emp_id=$1 AND spt.cust_id=$2 order by pl_id`
	err := repo.db.Select(&Details, query, params.EmpId, params.CustId)
	if err != nil {
		log.Println("spPriceRepository, FindDetailSpPriceIdAndCustId, err:", err.Error())
		return Details, err
	}

	return Details, nil
}

func (repository *salesmanRepositoryImpl) FindAllByCustId(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) ([]model.SalesmanList, int, int, error) {

	salesmans := []model.SalesmanList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` s.cust_id, s.emp_id, s.sales_name, s.sales_team_id, s.opr_type, s.is_bonus_rep, s.trans_date,
					s.wh_id, s.inc_grp_id, s.official_id, s.sale_system, s.sm_is_transfer, s.sm_valid_route,
					s.sm_geoloc_valid, s.sm_radius, s.image_url, s.sm_is_barcode, s.sm_is_photo_profile, COALESCE(s.sm_password, '') AS sm_password,
					s.is_active, s.updated_at, u.user_fullname AS updated_by_name, COALESCE(emp.emp_code, '') AS emp_code, COALESCE(emp.emp_name, '') AS emp_name, 
					CASE WHEN wh.wh_code IS NULL THEN wh2.wh_code ELSE wh.wh_code END AS wh_code, 
    				CASE WHEN wh.wh_name IS NULL THEN wh2.wh_name ELSE wh.wh_name END AS wh_name, 
					of.official_type, ofh.hierarchy_code, st.sales_team_code, st.sales_team_name, msc.opr_type as opr_type_canvas,
					CASE WHEN msc.is_active IS NULL THEN false ELSE msc.is_active END AS is_active_canvas, 
					CASE WHEN msc.is_active IS NOT NULL THEN wh2.wh_name ELSE NULL END AS wh_name_canvas, 
					CASE WHEN msc.is_active IS NOT NULL THEN wh2.wh_id ELSE NULL END AS wh_canvas_id, 

					CASE WHEN s.is_taking_order IS NULL THEN false ELSE s.is_taking_order END AS is_taking_order`

	qWhere := ` WHERE ` + buildSalesmanCustScopeCondition(dataFilter.CustIds, dataFilter.DistributorID, parentCustId, custId) + `
				AND s.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND s.sales_name ILIKE '%` + dataFilter.Query + `%' `
	}

	if dataFilter.SalesTeamId != "" {
		qWhere += ` AND s.sales_team_id IN (` + dataFilter.SalesTeamId + `) `
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND s.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND s.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_salesman s
	LEFT JOIN sys.m_user u ON u.user_id = s.updated_by 
	` + buildSalesmanEmployeeJoin() + `
	` + buildSalesmanWarehouseJoin("wh", "s.wh_id") + `
	LEFT JOIN mst.m_official of ON of.official_id = s.official_id AND of.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type AND ofh.cust_id = '` + parentCustId + `'
	` + buildSalesmanSalesTeamJoin() + `
	` + buildSalesmanCanvasJoin() + `
	` + buildSalesmanWarehouseJoin("wh2", "msc.wh_id") + ` `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("salesmanRepository, queryCount:", queryCount)
	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesmanRepository, count total, err:", err.Error())
		return salesmans, 0, 0, err
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
		sortBy := `s.emp_id`
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

	// log.Println("salesmanRepository, querySelect:", querySelect)
	err = repository.db.Select(&salesmans, querySelect)
	if err != nil {
		log.Println("salesmanRepository, FindAllByCustId, err:", err.Error())
		return salesmans, total, lastPage, err
	}

	return salesmans, total, lastPage, nil
}

func (repository *salesmanRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.SalesmanQueryFilter, custId, parentCustId string) ([]model.SalesmanList, int, int, error) {

	salesmans := []model.SalesmanList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` s.cust_id, s.emp_id, s.sales_name, s.sales_team_id, s.opr_type, s.is_bonus_rep, s.trans_date, s.wh_id,
					s.inc_grp_id, s.official_id, s.sale_system, s.sm_is_transfer, s.sm_valid_route,
					s.sm_geoloc_valid, s.sm_radius, s.image_url, s.sm_is_barcode, s.sm_is_photo_profile, COALESCE(s.sm_password, '') AS sm_password,
					s.is_active, s.updated_at, u.user_fullname AS updated_by_name, COALESCE(emp.emp_code, '') AS emp_code, COALESCE(emp.emp_name, '') AS emp_name, wh.wh_code, wh.wh_name,
					of.official_type, ofh.hierarchy_code, st.sales_team_code, st.sales_team_name, 
					CASE WHEN msc.is_active IS NULL THEN false ELSE msc.is_active END AS is_active_canvas, 
					CASE WHEN s.is_taking_order IS NULL THEN false ELSE s.is_taking_order END AS is_taking_order`
	qWhere := ` WHERE s.is_del = false 
				AND ` + buildSalesmanCustScopeCondition(dataFilter.CustIds, dataFilter.DistributorID, parentCustId, custId)

	if dataFilter.Query != "" {
		qWhere += ` AND s.sales_name ILIKE '%` + dataFilter.Query + `%' `
	}

	if dataFilter.SalesTeamId != "" {
		qWhere += ` AND s.sales_team_id IN (` + dataFilter.SalesTeamId + `) `
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND s.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND s.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_salesman s
	LEFT JOIN sys.m_user u ON u.user_id = s.updated_by 
	` + buildSalesmanEmployeeJoin() + `
	` + buildSalesmanWarehouseJoin("wh", "s.wh_id") + `
	LEFT JOIN mst.m_official of ON of.official_id = s.official_id AND of.cust_id = '` + parentCustId + `'
	LEFT JOIN mst.m_official_hierarchy ofh ON ofh.official_type = of.official_type AND ofh.cust_id = '` + parentCustId + `'
	` + buildSalesmanSalesTeamJoin() + `
	` + buildSalesmanCanvasJoin() + ` `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("salesmanRepository, queryCount:", queryCount)
	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesmanRepository, count total, err:", err.Error())
		return salesmans, 0, 0, err
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
		sortBy := `s.emp_id`
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

	// log.Println("salesmanRepository, querySelect:", querySelect)
	err = repository.db.Select(&salesmans, querySelect)
	if err != nil {
		log.Println("salesmanRepository, FindAllByCustId, err:", err.Error())
		return salesmans, total, lastPage, err
	}

	return salesmans, total, lastPage, nil
}

func (repository *salesmanRepositoryImpl) StoreDetail(salesmanDetail model.SalesmanDetail) error {
	query :=
		`INSERT INTO mst.m_salesman_product_type(
			cust_id, emp_id, group_type, ref_id)
		VALUES ( 
			$1, $2, $3, $4
		) RETURNING m_salesman_product_type_id;`
	lastInsertId := salesmanDetail.MSalesmanProductTypeID
	err := repository.tx.QueryRow(query,
		salesmanDetail.CustId, salesmanDetail.EmpId, salesmanDetail.GroupType, salesmanDetail.RefID,
	).Scan(&lastInsertId)
	if err != nil {
		log.Println("salesmanRepository, Store, err:", err.Error())
		return err
	}
	return nil
}

func (repository *salesmanRepositoryImpl) Store(salesman model.Salesman) (int64, error) {
	query :=
		`INSERT INTO mst.m_salesman(
			cust_id, emp_id, sales_name, sales_team_id,
			opr_type, is_bonus_rep, trans_date, wh_id, 
			inc_grp_id, official_id, sale_system,
			sm_is_transfer, sm_valid_route, sm_geoloc_valid, sm_radius, job_type, tax_option,start_date, end_date,allow_input_price,
			sm_password, is_active, created_by, created_at, 
			updated_by, updated_at, is_del, deleted_by, 
			deleted_at, image_url,sm_is_barcode,sm_is_photo_profile,is_taking_order)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12, 
			$13, $14, $15, $16, 
			$17, $18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30, $31, $32, $33
		) RETURNING emp_id;`
	lastInsertId := salesman.EmpId
	err := repository.tx.QueryRow(query,
		salesman.CustId, salesman.EmpId, salesman.SalesName, salesman.SalesTeamId,
		salesman.OprType, salesman.IsBonusRep, salesman.TransDate, salesman.WhIdTackingOrder,
		salesman.IncGrpId, salesman.OfficialId, salesman.SaleSystem,
		salesman.SmIsTransfer, salesman.SmValidRoute, salesman.SmGeolocValid, salesman.SmRadius, salesman.JobType, salesman.TaxOption, salesman.StartDate, salesman.EndDate, salesman.AllowInputPrice,
		salesman.SmPassword, salesman.IsActive, salesman.CreatedBy, salesman.CreatedAt,
		salesman.UpdatedBy, salesman.UpdatedAt, salesman.IsDel, salesman.DeletedBy,
		salesman.DeletedAt, salesman.ImageUrl, salesman.SmIsBarcode, salesman.SmIsPhotoProfile, salesman.IsTakingOrder).Scan(&lastInsertId)
	if err != nil {
		log.Println("salesmanRepository, Store, err:", err.Error())
		return salesman.EmpId, err
	}
	return salesman.EmpId, nil
}

func (repository *salesmanRepositoryImpl) StoreSalesmanCanvas(salesmanCanvas model.SalesmanCanvas) (int64, error) {
	query :=
		`INSERT INTO mst.m_salesman_canvas(
			cust_id, emp_id, wh_id, vehicle_id, is_active, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at, opr_type)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) RETURNING emp_id;`
	lastInsertId := salesmanCanvas.EmpId
	err := repository.tx.QueryRow(query,
		salesmanCanvas.CustId, salesmanCanvas.EmpId, salesmanCanvas.WhId, salesmanCanvas.VehicleId, salesmanCanvas.IsActive, salesmanCanvas.CreatedBy, salesmanCanvas.CreatedAt,
		salesmanCanvas.UpdatedBy, salesmanCanvas.UpdatedAt, salesmanCanvas.DeletedBy,
		salesmanCanvas.DeletedAt, salesmanCanvas.OprTypeCanvas).Scan(&lastInsertId)
	if err != nil {
		log.Println("salesmanCanvasRepository, Store, err:", err.Error())
		return salesmanCanvas.EmpId, err
	}
	return salesmanCanvas.EmpId, nil
}

func (repository *salesmanRepositoryImpl) StoreWarehouse(warehouse model.Warehouse) (int, error) {
	query :=
		`INSERT INTO mst.m_warehouse(
			cust_id, wh_code, wh_name, 
			is_active, created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, deleted_at, stock_type)
		VALUES ( 
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12
		) RETURNING wh_id;`
	lastInsertId := warehouse.WarehouseId
	err := repository.tx.QueryRow(query,
		warehouse.CustId, warehouse.WarehouseCode, warehouse.WarehouseName, warehouse.IsActive, warehouse.CreatedBy, warehouse.CreatedAt, warehouse.UpdatedBy, warehouse.UpdatedAt,
		warehouse.IsDel, warehouse.DeletedBy, warehouse.DeletedAt, warehouse.StockType).Scan(&lastInsertId)
	if err != nil {
		log.Println("warehouseRepository, Store, err:", err.Error())
		return warehouse.WarehouseId, err
	}
	warehouse.WarehouseId = lastInsertId
	fmt.Println("RESPONE WH >>>", warehouse)
	return warehouse.WarehouseId, nil
}

func (repository *salesmanRepositoryImpl) Update(empId int64, request entity.UpdateSalesmanRequest) error {
	var (
		r            model.SalesmanUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("salesmanRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_salesman
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND emp_id = :emp_id_old;`

	// log.Println("salesmanRepository, Update, query:", query)

	sqlPatch.Args["emp_id_old"] = empId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("salesmanRepository, Update, err:", err.Error())
		return err
	}

	if request.EndDate == nil {
		query = `UPDATE mst.m_salesman
			  SET end_date = NULL,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false
			  AND cust_id = :cust_id
			  AND emp_id = :emp_id_old;`
		result, err = repository.tx.NamedExec(query, map[string]interface{}{
			"cust_id":    request.CustId,
			"emp_id_old": empId,
		})
		if err != nil {
			log.Println("salesmanRepository, Update, err:", err.Error())
			return err
		}
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *salesmanRepositoryImpl) UpdateCanvas(empId int64, request entity.UpdateSalesmanRequest) error {
	fmt.Println("empId >>>>", empId)

	var (
		r            model.SalesmanCanvasUpdate
		sqlSetFields string
		nRows        int64
	)

	fmt.Println("request >>>>", request)
	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	fmt.Println("sqlPatch >>>>", sqlPatch)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("salesmanRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_salesman_canvas
			  SET ` + sqlSetFields + `, updated_at = CURRENT_TIMESTAMP
			  WHERE cust_id = :cust_id 
			  AND emp_id = :emp_id_old;`

	// log.Println("salesmanRepository, Update, query:", query)

	sqlPatch.Args["emp_id_old"] = empId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("salesmanRepository, Update, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		// return errors.New("no rows affected")
	}

	return nil
}

func (repository *salesmanRepositoryImpl) UpdateIsTakingOrder(empId int64, custId string) error {

	var nRows int64
	query := `UPDATE mst.m_salesman
			SET is_taking_order = false,wh_id = 0,
				updated_at = CURRENT_TIMESTAMP
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND emp_id = :emp_id;`

	wMap := map[string]interface{}{
		"cust_id": custId,
		"emp_id":  empId,
	}

	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesmanRepository, Delete, err:", err.Error())
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

func (repository *salesmanRepositoryImpl) Delete(custId string, empId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_salesman
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND emp_id = :emp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"emp_id":     empId,
		"deleted_by": deletedBy,
	}

	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesmanRepository, Delete, err:", err.Error())
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

func (repository *salesmanRepositoryImpl) DeleteDetailNotInIDs(empId int64, MSalesmanDetID []int64) error {

	query, args, err := sqlx.In("DELETE from mst.m_salesman_product_type WHERE emp_id =? AND m_salesman_product_type_id NOT IN(?) ;", empId, MSalesmanDetID)
	if err != nil {
		return err
	}
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = repository.tx.Rebind(query)
	qr, err := repository.tx.Query(query, args...)
	if err != nil {
		return err
	}
	qr.Next()
	return nil
}

func (repository *salesmanRepositoryImpl) DeleteDetails(empId int64, custId string) error {

	query, args, err := sqlx.In(
		`DELETE from mst.m_salesman_product_type 
		 WHERE cust_id = ? AND emp_id = ? ;`, custId, empId)
	if err != nil {
		return err
	}
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = repository.tx.Rebind(query)
	qr, err := repository.tx.Query(query, args...)
	if err != nil {
		return err
	}
	qr.Next()
	return nil
}

// func (repository *salesmanRepositoryImpl) DeleteSalesmanCanvas(empId int64, custId string) error {

// 	query, args, err := sqlx.In(
// 		`DELETE from mst.m_salesman_canvas
// 		 WHERE cust_id = ? AND emp_id = ? ;`, custId, empId)
// 	if err != nil {
// 		return err
// 	}
// 	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
// 	query = repository.tx.Rebind(query)
// 	qr, err := repository.tx.Query(query, args...)
// 	if err != nil {
// 		return err
// 	}
// 	qr.Next()
// 	return nil
// }

func (repository *salesmanRepositoryImpl) UpdateDetail(empId int64, request model.SalesmanDetail) error {
	var (
		sqlSetFields string
		nRows        int64
	)
	sqlPatch := sql_helper.SQLPatches(request)

	data, err := json.Marshal(sqlPatch)
	if err != nil {
		return err
	}
	fmt.Printf("salesmanRepository, Update det, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_salesman_product_type
			  SET ` + sqlSetFields + `
			  WHERE
			  emp_id = :emp_id_old  
			  AND m_salesman_product_type_id = :m_salesman_product_type_id;`

	log.Println("salesmanRepository, Update detail, query:", query)

	sqlPatch.Args["emp_id_old"] = empId
	sqlPatch.Args["m_salesman_product_type_id"] = request.MSalesmanProductTypeID
	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("salesmanRepository, Update, err:", err.Error())
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

func (repo *salesmanRepositoryImpl) FindDetailById(empId int64, custId string) ([]model.SalesmanDetailRead, error) {
	Details := []model.SalesmanDetailRead{}
	query := `SELECT 
	 CASE 
        WHEN pl.pl_id IS NOT NULL THEN pl.pl_id
        WHEN b.pl_id IS NOT NULL THEN b.pl_id
        WHEN sb.sbrand1_id IS NOT NULL THEN b2.pl_id
    END AS pl_id,
	CASE 
		WHEN pl.pl_id is not null THEN pl.pl_id
		WHEN b.brand_id is not null THEN b.brand_id
		WHEN sb.sbrand1_id is not null THEN sb.sbrand1_id 
	END as ref_id,
	CASE 
		WHEN pl.pl_code is not null THEN pl.pl_code
		WHEN b.brand_code is not null THEN b.brand_code
		WHEN sb.sbrand1_code is not null THEN sb.sbrand1_code 
	END as ref_code,
	CASE 
		WHEN pl.pl_name is not null THEN pl.pl_name
		WHEN b.brand_name is not null THEN b.brand_name
		WHEN sb.sbrand1_name is not null THEN sb.sbrand1_name
	END as ref_name,
	spt.group_type,
	spt.m_salesman_product_type_id
	FROM mst.m_salesman_product_type spt
	LEFT JOIN mst.m_product_line pl ON pl.pl_id = spt.ref_id AND spt.group_type=1
	LEFT JOIN mst.m_brand b ON b.brand_id = spt.ref_id AND spt.group_type=2
	LEFT JOIN mst.m_sub_brand1 sb ON sb.sbrand1_id = spt.ref_id AND spt.group_type=3
	LEFT JOIN mst.m_brand b2 ON b2.brand_id = sb.brand_id
	WHERE spt.emp_id=$1 AND spt.cust_id=$2 order by pl_id`
	err := repo.db.Select(&Details, query, empId, custId)
	if err != nil {
		log.Println("spPriceRepository, FindDetailSpPriceIdAndCustId, err:", err.Error())
		return Details, err
	}

	return Details, nil
}

func (repository *salesmanRepositoryImpl) UpdateDeActive(empId int64, custId string) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	var nRows int64
	query := `UPDATE mst.m_salesman
			SET is_active = false,
				updated_at = CURRENT_TIMESTAMP
			WHERE is_del = false
			AND cust_id = :cust_id
			AND end_date IS NOT NULL
			AND end_date <= :end_date
			AND emp_id = :emp_id;`

	wMap := map[string]interface{}{
		"cust_id":  custId,
		"emp_id":   empId,
		"end_date": now.Format("2006-01-02"),
	}

	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesmanRepository, Update De Active, err:", err.Error())
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

func (repository *salesmanRepositoryImpl) UpdateIsActive(empId int64, custId string) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	var nRows int64
	query := `UPDATE mst.m_salesman
			SET is_active = true,
				updated_at = CURRENT_TIMESTAMP
			WHERE is_del = false
			AND cust_id = :cust_id
			AND start_date IS NOT NULL
			AND start_date <= :start_date
			AND emp_id = :emp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"emp_id":     empId,
		"start_date": now.Format("2006-01-02"),
	}

	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesmanRepository, Update Is Active, err:", err.Error())
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

func (repository *salesmanRepositoryImpl) CheckDate() {
	var result struct {
		CurrentTimestamp time.Time `db:"current_timestamp"`
		CurrentDate      time.Time `db:"current_date"`
	}

	query := `SELECT CURRENT_TIMESTAMP, CURRENT_DATE`
	err := repository.db.Get(&result, query) // Gunakan .Get untuk single row
	if err != nil {
		log.Println("SalesmanRepository, cekDate error:", err.Error())
		return
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	fmt.Println("Current Timestamp Go:", now)
	fmt.Println("Current Timestamp Go Format:", now.Format("2006-01-02"))

	log.Println("Current Timestamp:", result.CurrentTimestamp)
	log.Println("Current Date:", result.CurrentDate.Format("2006-01-02"))
}
