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

type OfficialHierarchyRepository interface {
	TrxBegin() (*officialHierarchyTransaction, error)
	FindOneByOfficialTypeAndCustId(officialType int, custId string) (model.OfficialHierarchy, error)
	FindAllByCustId(dataFilter entity.OfficialHierarchyQueryFilter, custId string) (official []model.OfficialHierarchy, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.OfficialHierarchyQueryFilter, custId string) (official []model.OfficialHierarchy, total int, lastPage int, err error)
	Store(official model.OfficialHierarchy) (int, error)
	Upsert(request entity.CreateOfficialHierarchyBody) error
	Update(officialId int, request entity.UpdateOfficialHierarchyRequest) error
	Delete(custId string, officialType int, deletedBy int64) error
}

func NewOfficialHierarchyRepository(db *sqlx.DB) OfficialHierarchyRepository {
	return &officialHierarchyRepositoryImpl{db}
}

func NewOfficialHierarchyTransaction(db *sqlx.DB) (trxObj *officialHierarchyTransaction, err error) {
	trx := db.MustBegin()

	return &officialHierarchyTransaction{tx: trx, db: db}, nil
}

type officialHierarchyRepositoryImpl struct {
	*sqlx.DB
}

type officialHierarchyTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func (repo *officialHierarchyRepositoryImpl) TrxBegin() (*officialHierarchyTransaction, error) {
	trxObj, err := NewOfficialHierarchyTransaction(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *officialHierarchyTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *officialHierarchyTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *officialHierarchyRepositoryImpl) FindOneByOfficialTypeAndCustId(officialType int, custId string) (model.OfficialHierarchy, error) {
	official := model.OfficialHierarchy{}
	query := `SELECT 
				of.cust_id, 
				of.official_type, of.hierarchy_code, of.is_active, 
				of.created_by, of.created_at, of.updated_by, of.updated_at
			  FROM mst.m_official_hierarchy of
			  WHERE of.official_type = $1 
			  AND of.cust_id = $2`
	err := repository.Get(&official, query, officialType, custId)
	if err != nil {
		log.Println("OfficialHierarchyRepository, FindOneByOfficialTypeAndCustId, err:", err.Error())
		return official, err
	}

	return official, nil
}

func (repository *officialHierarchyRepositoryImpl) FindAllByCustId(dataFilter entity.OfficialHierarchyQueryFilter, custId string) ([]model.OfficialHierarchy, int, int, error) {

	officials := []model.OfficialHierarchy{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` 	of.cust_id, 
						of.official_type, of.hierarchy_code, of.is_active, 
						of.created_by, of.created_at, of.updated_by, of.updated_at `
	qWhere := ` WHERE of.is_del = false 
				AND of.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (of.hierarchy_code = '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OfficialType != 0 {
		qWhere += ` AND of.official_type = ` + strconv.Itoa(dataFilter.OfficialType) + ` `
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND of.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND of.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_official_hierarchy of `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("OfficialHierarchyRepository, count total, err:", err.Error())
		return officials, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`of.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `of.official_type`
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

	err = repository.Select(&officials, querySelect)
	if err != nil {
		log.Println("OfficialHierarchyRepository, FindAllByCustId, err:", err.Error())
		return officials, total, lastPage, err
	}

	return officials, total, lastPage, nil
}

func (repository *officialHierarchyRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.OfficialHierarchyQueryFilter, custId string) ([]model.OfficialHierarchy, int, int, error) {

	officials := []model.OfficialHierarchy{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` 	of.cust_id, 
						of.official_type, of.hierarchy_code, of.is_active, 
						of.created_by, of.created_at, of.updated_by, of.updated_at,
						of.is_del, of.deleted_by, of.deleted_at  `
	qWhere := ` WHERE of.is_del = false AND of.is_active = true
				AND of.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (of.hierarchy_code ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OfficialType != 0 {
		qWhere += ` AND of.official_type = ` + strconv.Itoa(dataFilter.OfficialType) + ` `
	}

	qFrom := ` FROM mst.m_official_hierarchy of `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("OfficialHierarchyRepository, count total, err:", err.Error())
		return officials, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`of.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `of.official_type`
		querySelect += fmt.Sprintf(`ORDER BY %s ASC`, sortBy)
	}

	lastPage := 1

	// log.Println("officialRepository, querySelect:", querySelect)
	err = repository.Select(&officials, querySelect)
	if err != nil {
		log.Println("OfficialHierarchyRepository, FindAllByCustIdLookup, err:", err.Error())
		return officials, total, lastPage, err
	}

	return officials, total, lastPage, nil
}

func (repository *officialHierarchyRepositoryImpl) Store(official model.OfficialHierarchy) (int, error) {
	query :=
		`INSERT INTO mst.m_official_hierarchy(
			cust_id, official_type, hierarchy_code, is_active, 
			created_by, created_at, updated_by, updated_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8
		) RETURNING official_type;`
	lastInsertOfficialType := official.OfficialType
	err := repository.QueryRow(query,
		official.CustId, official.OfficialType, official.HierarchyCode, official.IsActive,
		official.CreatedBy, official.CreatedAt, official.UpdatedBy, official.UpdatedAt).Scan(&lastInsertOfficialType)
	if err != nil {
		log.Println("OfficialHierarchyRepository, Store, err:", err.Error())
		return official.OfficialType, err
	}
	return official.OfficialType, nil
}

func (repository *officialHierarchyRepositoryImpl) Upsert(request entity.CreateOfficialHierarchyBody) error {
	var (
		r            model.OfficialHierarchyUpsert
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("OfficialHierarchyRepository, Upsert, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query :=
		`
		INSERT INTO mst.m_official_hierarchy (
			cust_id, official_type, hierarchy_code, is_active, 
			created_by, created_at, updated_by, updated_at
		) VALUES (:cust_id, :official_type, :hierarchy_code, :is_active, 
			:created_by, :created_at, :updated_by, :updated_at
		)
		ON CONFLICT (cust_id, official_type)
		DO UPDATE SET ` + sqlSetFields + `; `

	sqlPatch.Args["cust_id"] = request.CustId
	sqlPatch.Args["official_type"] = request.OfficialType

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("OfficialHierarchyRepository, Upsert, err:", err.Error())
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

func (repository *officialHierarchyRepositoryImpl) Update(officialType int, request entity.UpdateOfficialHierarchyRequest) error {
	var (
		r            model.OfficialHierarchyUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("OfficialHierarchyRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_official_hierarchy
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND official_type = :official_type_old;`

	// log.Println("OfficialHierarchyRepository, Update, query:", query)

	sqlPatch.Args["official_type_old"] = officialType
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("OfficialHierarchyRepository, Update, err:", err.Error())
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

func (repository *officialHierarchyRepositoryImpl) Delete(custId string, officialType int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_official_hierarchy
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND official_type = :official_type;`

	wMap := map[string]interface{}{
		"cust_id":       custId,
		"official_type": officialType,
		"deleted_by":    deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OfficialHierarchyRepository, Delete, err:", err.Error())
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

func (repository *officialHierarchyTransaction) UpsertWithTrx(request entity.CreateOfficialHierarchyBody) error {
	var (
		r            model.OfficialHierarchyUpsert
		sqlSetFields string
		nRows        int64
	)

	log.Println("UpsertWithTrx - REPOSITORY")
	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("OfficialHierarchyRepository, Upsert, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query :=
		`
		INSERT INTO mst.m_official_hierarchy (
			cust_id, official_type, hierarchy_code, is_active, 
			created_by, created_at, updated_by, updated_at
		) VALUES (:cust_id, :official_type, :hierarchy_code, :is_active, 
			:created_by, :created_at, :updated_by, :updated_at
		)
		ON CONFLICT (cust_id, official_type)
		DO UPDATE SET ` + sqlSetFields + `; `

	sqlPatch.Args["cust_id"] = request.CustId
	sqlPatch.Args["official_type"] = request.OfficialType

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("OfficialHierarchyRepository, Upsert, err:", err.Error())
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