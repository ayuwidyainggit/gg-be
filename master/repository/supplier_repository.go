package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/structs"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type SupplierRepository interface {
	FindOneBySupplierIdAndCustId(supplierId int, custId string) (model.SupplierRead, error)
	FindOneBySupplierCodeAndCustId(supplierCode string, custId string) (model.Supplier, error)
	FindAllByCustId(dataFilter entity.SupplierQueryFilter, custId string) (sup []model.Supplier, total int, lastPage int, err error)
	FindAllLookupByCustId(dataFilter entity.SupplierQueryFilter, scope entity.SupplierLookupScope) (sup []model.Supplier, total int, lastPage int, err error)
	Store(supplier model.Supplier) (int, error)
	Update(supplierId int, request entity.UpdateSupplierRequest) error
	Delete(custId string, supplierId int, deletedBy int64) error
}

func NewSupplierRepository(db *sqlx.DB) SupplierRepository {
	return &SupplierRepositoryImpl{db}
}

type SupplierRepositoryImpl struct {
	*sqlx.DB
}

func (repository *SupplierRepositoryImpl) FindOneBySupplierIdAndCustId(supplierId int, custId string) (model.SupplierRead, error) {
	supplier := model.SupplierRead{}
	query := `SELECT 
				mst.m_supplier.cust_id, sup_id, sup_code, sup_name, 
				address1, address2, city, phone_no,
				fax_no, sup_type, contact_name, tax_name,
				tax_no, pay_term, is_credit_limit, credit_limit,
				mst.m_supplier.province_id,province,mst.m_supplier.regency_id,regency,mst.m_supplier.sub_district_id,sub_district,mst.m_supplier.ward_id,ward,zip_code,ot_loc_id,latitude,longitude,email,is_wa_no,wa_no,contact_type,credit_limit_type,
				mst.m_supplier.is_active, mst.m_supplier.created_by,
				mst.m_supplier.created_at, mst.m_supplier.updated_by, mst.m_supplier.updated_at,
				mst.m_supplier.is_del, mst.m_supplier.deleted_by, mst.m_supplier.deleted_at,
				mst.m_supplier.phone, mst.m_supplier.fax_number, mst.m_supplier.tax_identifier_no, mst.m_supplier.nitku, mst.m_supplier.tax_address
			  FROM mst.m_supplier
			  left join mst.m_province pro on mst.m_supplier.province_id = pro.province_id
			  left join mst.m_regency reg on mst.m_supplier.regency_id = reg.regency_id
			  left join mst.m_sub_district dist on mst.m_supplier.sub_district_id = dist.sub_district_id
			  left join mst.m_ward wrd on mst.m_supplier.ward_id = wrd.ward_id
			  WHERE mst.m_supplier.sup_id = $1 
			  AND mst.m_supplier.cust_id = $2`
	err := repository.Get(&supplier, query, supplierId, custId)
	if err != nil {
		log.Println("SupplierRepository, FindOneBySupplierIdAndCustId, err:", err.Error())
		return supplier, err
	}

	return supplier, nil
}

func (repository *SupplierRepositoryImpl) FindOneBySupplierCodeAndCustId(supplierCode string, custId string) (model.Supplier, error) {
	supplier := model.Supplier{}
	query := `SELECT 
				cust_id, sup_id, sup_code, sup_name, 
				address1, address2, city, phone_no,
				fax_no, sup_type, contact_name, tax_name,
				tax_no, pay_term, is_credit_limit, credit_limit,
				province_id,regency_id,sub_district_id,ward_id,zip_code,ot_loc_id,latitude,longitude,email,is_wa_no,wa_no,contact_type,credit_limit_type,
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_supplier 
			  WHERE sup_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&supplier, query, supplierCode, custId)
	if err != nil {
		log.Println("SupplierRepository, FindOneBySupplierCodeAndCustId, err:", err.Error())
		return supplier, err
	}

	return supplier, nil
}

func (repository *SupplierRepositoryImpl) FindAllByCustId(dataFilter entity.SupplierQueryFilter, custId string) ([]model.Supplier, int, int, error) {

	suppliers := []model.Supplier{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sup_id, a.sup_code, a.sup_name, 
					a.address1, a.address2, a.city, a.phone_no,
					a.fax_no, a.sup_type, a.contact_name, a.tax_name,
					a.tax_no, a.pay_term, a.is_credit_limit, a.credit_limit,
					a.province_id, a.regency_id, a.sub_district_id, a.ward_id, a.zip_code, a.ot_loc_id, a.latitude, a.longitude, a.email, a.is_wa_no, a.wa_no, a.contact_type, a.credit_limit_type,
					a.is_active, a.created_by, a.created_at, a.updated_by, 
					a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sup_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sup_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		// fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_supplier a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_supplier a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere

	// log.Println("SupplierRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("SupplierRepository, count total, err:", err.Error())
		return suppliers, 0, 0, err
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
		sortBy := `sup_id`
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

	// log.Println("SupplierRepository, querySelect:", querySelect)
	err = repository.Select(&suppliers, querySelect)
	if err != nil {
		log.Println("SupplierRepository, FindAllByCustId, err:", err.Error())
		return suppliers, total, lastPage, err
	}

	// log.Println("data:", structs.StructToJson(suppliers))

	return suppliers, total, lastPage, nil
}

func supplierLookupOrderSQL(sort string) string {
	if sort != "" {
		var sortBy strings.Builder
		mSortBy := strings.Split(sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy.WriteString(fmt.Sprintf(`%s %s, `, colSort[0], colSort[1]))
			}
		}
		s := strings.TrimSuffix(sortBy.String(), ", ")
		if s != "" {
			return ` ORDER BY ` + s
		}
	}
	return ` ORDER BY sup_id DESC `
}

func (repository *SupplierRepositoryImpl) FindAllLookupByCustId(dataFilter entity.SupplierQueryFilter, scope entity.SupplierLookupScope) ([]model.Supplier, int, int, error) {
	useUnion := scope.IncludeParentID && len(scope.DistributorIDs) > 0
	if useUnion {
		return repository.findAllLookupUnion(dataFilter, scope)
	}
	return repository.findAllLookupSingle(dataFilter, scope.ParentCustID)
}

func (repository *SupplierRepositoryImpl) findAllLookupSingle(dataFilter entity.SupplierQueryFilter, scopeCustID string) ([]model.Supplier, int, int, error) {
	suppliers := []model.Supplier{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` cust_id, sup_id, sup_code, sup_name,
					address1,city, phone_no, sup_type `
	args := []interface{}{scopeCustID}
	qWhere := ` WHERE is_del = false AND is_active = true AND cust_id = $1 `

	if dataFilter.Query != "" {
		qWhere += ` AND (sup_code ILIKE $2 OR sup_name ILIKE $2) `
		args = append(args, "%"+dataFilter.Query+"%")
	}
	if dataFilter.SupType != "" {
		n := len(args) + 1
		qWhere += fmt.Sprintf(` AND sup_type = $%d `, n)
		args = append(args, dataFilter.SupType)
	}

	qFrom := ` FROM mst.m_supplier `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount, args...).Scan(&total)
	if err != nil {
		log.Println("SupplierRepository, count total, err:", err.Error())
		return suppliers, 0, 0, err
	}

	querySelect += supplierLookupOrderSQL(dataFilter.Sort)

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := sql_helper.CalculateLastPage(total, limit)
	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	err = repository.Select(&suppliers, querySelect, args...)
	if err != nil {
		log.Println("SupplierRepository, findAllLookupSingle, err:", err.Error())
		return suppliers, total, lastPage, err
	}

	return suppliers, total, lastPage, nil
}

func (repository *SupplierRepositoryImpl) findAllLookupUnion(dataFilter entity.SupplierQueryFilter, scope entity.SupplierLookupScope) ([]model.Supplier, int, int, error) {
	suppliers := []model.Supplier{}
	args := []interface{}{scope.CustID, pq.Array(scope.DistributorIDs)}
	nextArg := 3
	qFilter := ""
	if dataFilter.Query != "" {
		pat := "%" + dataFilter.Query + "%"
		qFilter = fmt.Sprintf(` AND (ms.sup_code ILIKE $%d OR ms.sup_name ILIKE $%d)`, nextArg, nextArg)
		args = append(args, pat)
		nextArg++
	}
	supTypeFilter := ""
	if dataFilter.SupType != "" {
		supTypeFilter = fmt.Sprintf(` AND ms.sup_type = $%d`, nextArg)
		args = append(args, dataFilter.SupType)
		nextArg++
	}

	unionInner := fmt.Sprintf(`
(
	SELECT ms.cust_id, ms.sup_id, ms.sup_code, ms.sup_name, ms.address1, ms.city, ms.phone_no, ms.sup_type
	FROM mst.m_supplier ms
	WHERE ms.cust_id = $1 AND ms.is_del = false AND ms.is_active = true
	%s%s
)
UNION
(
	SELECT ms.cust_id, ms.sup_id, ms.sup_code, ms.sup_name, ms.address1, ms.city, ms.phone_no, ms.sup_type
	FROM mst.m_supplier ms
	INNER JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id
	WHERE mc.distributor_id = ANY($2::int[])
		AND ms.is_del = false AND ms.is_active = true
	%s%s
)`, qFilter, supTypeFilter, qFilter, supTypeFilter)

	countQuery := `SELECT COUNT(*) FROM (` + unionInner + `) AS u`
	var total int
	if err := repository.Get(&total, countQuery, args...); err != nil {
		log.Println("SupplierRepository, findAllLookupUnion count, err:", err.Error())
		return suppliers, 0, 0, err
	}

	limit := dataFilter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := sql_helper.CalculateLastPage(total, limit)

	orderSQL := supplierLookupOrderSQL(dataFilter.Sort)
	// sort diterapkan pada hasil gabungan; kolom disamakan dengan lookup tunggal
	selectQuery := `SELECT cust_id, sup_id, sup_code, sup_name, address1, city, phone_no, sup_type FROM (` + unionInner + `) AS u ` +
		orderSQL + fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	if err := repository.Select(&suppliers, selectQuery, args...); err != nil {
		log.Println("SupplierRepository, findAllLookupUnion, err:", err.Error())
		return suppliers, total, lastPage, err
	}

	return suppliers, total, lastPage, nil
}

func (repository *SupplierRepositoryImpl) Store(supplier model.Supplier) (int, error) {

	log.Println("supplier:", structs.StructToJson(supplier))
	query :=
		`INSERT INTO mst.m_supplier(
			cust_id, distributor_id, sup_code, sup_name, address1, 
			province_id,regency_id,sub_district_id,ward_id,zip_code,ot_loc_id,latitude,longitude,email,is_wa_no,wa_no,contact_type,credit_limit_type,
			address2, city, phone_no, fax_no, 
			sup_type, contact_name, tax_name, tax_no, 
			pay_term, is_credit_limit, credit_limit, phone, fax_number, tax_identifier_no, nitku, tax_address, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			:cust_id, :distributor_id, :sup_code, :sup_name, :address1, 
			:province_id, :regency_id, :sub_district_id, :ward_id, 
			:zip_code, :ot_loc_id, :latitude, :longitude,
			:email, :is_wa_no, :wa_no, :contact_type, :credit_limit_type,
			:address2, :city, :phone_no, :fax_no,
			:sup_type, :contact_name, :tax_name, :tax_no,
			:pay_term, :is_credit_limit, :credit_limit, :phone,
			:fax_number, :tax_identifier_no, :nitku, :tax_address,
			:is_active, :created_by, :created_at, :updated_by, :updated_at,
			:is_del, :deleted_by, :deleted_at
		) RETURNING sup_id;`
	rows, err := repository.NamedQuery(query, supplier)
	if err != nil {
		log.Println("SupplierRepository, Store, err:", err.Error())
		return supplier.SupplierId, err
	}
	defer rows.Close()

	var lastInsertId int
	if rows.Next() {
		if err := rows.Scan(&lastInsertId); err != nil {
			log.Println("SupplierRepository, Store, scan err:", err.Error())
			return supplier.SupplierId, err
		}
		return lastInsertId, nil
	}

	return supplier.SupplierId, errors.New("no rows affected")
}

func (repository *SupplierRepositoryImpl) Update(supplierId int, request entity.UpdateSupplierRequest) error {
	var (
		r            model.SupplierUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("SupplierRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_supplier
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sup_id = :sup_id_old;`

	// log.Println("SupplierRepository, Update, query:", query)

	sqlPatch.Args["sup_id_old"] = supplierId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("SupplierRepository, Update, err:", err.Error())
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

func (repository *SupplierRepositoryImpl) Delete(custId string, supplierId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_supplier
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND sup_id = :sup_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"sup_id":     supplierId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SupplierRepository, Delete, err:", err.Error())
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
