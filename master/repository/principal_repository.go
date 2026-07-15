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

type PrincipalRepository interface {
	FindOneByPrincipalIdAndCustId(custId string, principalId int) (model.Principal, error)
	FindOneByPrincipalCodeAndCustId(principalCode, custId string) (model.Principal, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Principal, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Principal, total int, lastPage int, err error)
	Store(principal model.Principal) (int, error)
	Update(PrincipalId int, request entity.UpdatePrincipalRequest) error
	Delete(custId string, PrincipalId int, deletedBy int64) error
}

func NewPrincipalRepository(db *sqlx.DB) PrincipalRepository {
	return &principalRepositoryImpl{db}
}

type principalRepositoryImpl struct {
	*sqlx.DB
}

func (repository *principalRepositoryImpl) FindOneByPrincipalIdAndCustId(custId string, principalId int) (model.Principal, error) {
	principal := model.Principal{}
	query := `SELECT 
				pc.cust_id, pc.principal_id, pc.principal_code,
				pc.principal_name, pc.is_active, 
				us.user_fullname AS updated_by_name, 
				pc.updated_at
			  FROM mst.m_principal pc
			  LEFT JOIN sys.m_user us ON us.user_id = pc.updated_by
			  WHERE pc.is_del = false 
			  AND pc.cust_id = $1 
			  AND pc.principal_id = $2`
	err := repository.Get(&principal, query, custId, principalId)
	if err != nil {
		log.Println("principalRepository, FindOneByPrincipalIdAndCustId, err:", err.Error())
		return principal, err
	}

	return principal, nil
}

func (repository *principalRepositoryImpl) FindOneByPrincipalCodeAndCustId(principalCode, custId string) (model.Principal, error) {
	principal := model.Principal{}
	query := `SELECT 
				cust_id, principal_id, principal_code,
				principal_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_principal 
			  WHERE is_del = false 
			  AND cust_id = $1 
			  AND principal_code = $2`
	err := repository.Get(&principal, query, custId, principalCode)
	if err != nil {
		log.Println("principalRepository, FindOneByPrincipalCodeAndCustId, err:", err.Error())
		return principal, err
	}

	return principal, nil
}

func (repository *principalRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, CustId string) ([]model.Principal, int, int, error) {

	Principals := []model.Principal{}
	selectCount := ` COUNT(*) AS total `
	selectField := `p.principal_id, p.principal_code,
	p.principal_name, p.is_active, p.created_by,
	p.created_at, p.updated_by, p.updated_at,
	p.is_del, p.deleted_by, p.deleted_at,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE p.is_del = false 
				AND p.cust_id = '` + CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.principal_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.principal_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	qFrom := ` FROM mst.m_principal p
	LEFT JOIN sys.m_user u ON u.user_id = p.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("principalRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("err count:", err.Error())
	}
	// log.Println("total:", total)

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
		sortBy := `principal_id`
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

	// log.Println("principalRepository, List, querySelect:", querySelect)
	err = repository.Select(&Principals, querySelect)
	if err != nil {
		log.Println("principalRepository, List, err:", err.Error())
		return Principals, total, lastPage, err
	}

	return Principals, total, lastPage, nil
}

func (repository *principalRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, CustId string) ([]model.Principal, int, int, error) {

	Principals := []model.Principal{}
	selectCount := ` COUNT(*) AS total `
	selectField := `p.principal_id, p.principal_code,
					p.principal_name, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE p.is_del = false AND p.is_active = true
				AND p.cust_id = '` + CustId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.principal_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.principal_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` FROM mst.m_principal p
			   LEFT JOIN sys.m_user u ON u.user_id = p.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("principalRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("err count:", err.Error())
	}
	// log.Println("total:", total)

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
		sortBy := `principal_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	lastPage := 1
	// log.Println("principalRepository, List, querySelect:", querySelect)
	err = repository.Select(&Principals, querySelect)
	if err != nil {
		log.Println("principalRepository, List, err:", err.Error())
		return Principals, total, lastPage, err
	}

	return Principals, total, lastPage, nil
}

func (repository *principalRepositoryImpl) Store(principal model.Principal) (int, error) {
	query :=
		`INSERT INTO mst.m_principal(
			cust_id, principal_code, principal_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING principal_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		principal.CustId, principal.PrincipalCode, principal.PrincipalName, principal.IsActive,
		principal.CreatedBy, principal.CreatedAt, principal.UpdatedBy, principal.UpdatedAt,
		principal.IsDel, principal.DeletedBy, principal.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("principalRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *principalRepositoryImpl) Update(PrincipalId int, request entity.UpdatePrincipalRequest) error {
	var (
		r            model.PrincipalUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("principalRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_principal
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND principal_id = :principal_id;`

	// log.Println("principalRepository, Update, query:", query)

	sqlPatch.Args["principal_id"] = PrincipalId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("principalRepository, Update, err.Error():", err.Error())
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

func (repository *principalRepositoryImpl) Delete(custId string, PrincipalId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_principal
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE cust_id = :cust_id 
			AND is_del = false
			AND principal_id = :principal_id;`

	wMap := map[string]interface{}{
		"cust_id":      custId,
		"principal_id": PrincipalId,
		"deleted_by":   deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("principalRepository, Delete, err:", err.Error())
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
