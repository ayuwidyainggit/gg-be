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

type PluGroupRepository interface {
	FindOneByPluGroupIdAndCustId(pluGroupId int, custId string) (model.PluGroup, error)
	FindOneByPluGroupCodeAndCustId(pluGroupCode string, custId string) (model.PluGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (pluGroup []model.PluGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (pluGroup []model.PluGroup, total int, lastPage int, err error)
	Store(pluGroup model.PluGroup) (int, error)
	Update(pluGroupId int, request entity.UpdatePluGroupRequest) error
	Delete(custId string, pluGroupId int, deletedBy int64) error
}

func NewPluGroupRepository(db *sqlx.DB) PluGroupRepository {
	return &pluGroupRepositoryImpl{db}
}

type pluGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *pluGroupRepositoryImpl) FindOneByPluGroupIdAndCustId(pluGroupId int, custId string) (model.PluGroup, error) {
	pluGroup := model.PluGroup{}
	query := `SELECT 
				cust_id, plu_grp_id, plu_grp_code,
				plu_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_plu_group 
			  WHERE plu_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&pluGroup, query, pluGroupId, custId)
	if err != nil {
		log.Println("pluGroupRepository, FindOneByPluGroupCodeAndCustId, err:", err.Error())
		return pluGroup, err
	}

	return pluGroup, nil
}

func (repository *pluGroupRepositoryImpl) FindOneByPluGroupCodeAndCustId(pluGroupCode string, custId string) (model.PluGroup, error) {
	pluGroup := model.PluGroup{}
	query := `SELECT 
				cust_id, plu_grp_id, plu_grp_code,
				plu_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_plu_group 
			  WHERE plu_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&pluGroup, query, pluGroupCode, custId)
	if err != nil {
		log.Println("pluGroupRepository, FindOneByPluGroupCodeAndCustId, err:", err.Error())
		return pluGroup, err
	}

	return pluGroup, nil
}

func (repository *pluGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PluGroup, int, int, error) {

	pluGroups := []model.PluGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.plu_grp_id, a.plu_grp_code, a.plu_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.plu_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.plu_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_plu_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("pluGroupRepository, count total, err:", err.Error())
		return pluGroups, 0, 0, err
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
		sortBy := `a.plu_grp_id`
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

	err = repository.Select(&pluGroups, querySelect)
	if err != nil {
		log.Println("pluGroupRepository, FindAllByCustId, err:", err.Error())
		return pluGroups, total, lastPage, err
	}

	return pluGroups, total, lastPage, nil
}

func (repository *pluGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PluGroup, int, int, error) {

	pluGroups := []model.PluGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.plu_grp_id, a.plu_grp_code, a.plu_grp_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.plu_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.plu_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_plu_group a `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("pluGroupRepository, count total, err:", err.Error())
		return pluGroups, 0, 0, err
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
		sortBy := `a.plu_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&pluGroups, querySelect)
	if err != nil {
		log.Println("pluGroupRepository, FindAllByCustId, err:", err.Error())
		return pluGroups, total, 1, err
	}

	return pluGroups, total, 1, nil
}

func (repository *pluGroupRepositoryImpl) Store(pluGroup model.PluGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_plu_group(
			cust_id, plu_grp_code, plu_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING plu_grp_id;`
	lastInsertId := pluGroup.PluGrpId
	err := repository.QueryRow(query,
		pluGroup.CustId, pluGroup.PluGrpCode, pluGroup.PluGrpName,
		pluGroup.IsActive, pluGroup.CreatedBy, pluGroup.CreatedAt, pluGroup.UpdatedBy,
		pluGroup.UpdatedAt, pluGroup.IsDel, pluGroup.DeletedBy, pluGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("pluGroupRepository, Store, err:", err.Error())
		return pluGroup.PluGrpId, err
	}
	return pluGroup.PluGrpId, nil
}

func (repository *pluGroupRepositoryImpl) Update(pluGroupId int, request entity.UpdatePluGroupRequest) error {
	var (
		r            model.PluGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("pluGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_plu_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND plu_grp_id = :plu_grp_id_old;`

	// log.Println("pluGroupRepository, Update, query:", query)

	sqlPatch.Args["plu_grp_id_old"] = pluGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("pluGroupRepository, Update, err:", err.Error())
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

func (repository *pluGroupRepositoryImpl) Delete(custId string, pluGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_plu_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND plu_grp_id = :plu_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"plu_grp_id": pluGroupId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PluGroupRepository, Delete, err:", err.Error())
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
