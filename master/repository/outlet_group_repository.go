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

type OutletGroupRepository interface {
	FindOneByOutletGroupIdAndCustId(outletGroupId int64, custId string) (model.OutletGroup, error)
	FindOneByOutletGroupCodeAndCustId(outletGroupCode string, custId string) (model.OutletGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (outletGroup []model.OutletGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (outletGroup []model.OutletGroup, total int, lastPage int, err error)
	Store(outletGroup model.OutletGroup) (int, error)
	Update(outletGroupId int, request entity.UpdateOutletGroupRequest) error
	Delete(custId string, outletGroupId int, deletedBy int64) error
}

func NewOutletGroupRepository(db *sqlx.DB) OutletGroupRepository {
	return &outletGroupRepositoryImpl{db}
}

type outletGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *outletGroupRepositoryImpl) FindOneByOutletGroupIdAndCustId(outletGroupId int64, custId string) (model.OutletGroup, error) {
	outletGroup := model.OutletGroup{}
	query := `SELECT 
				cust_id, ot_grp_id, ot_grp_code,
				ot_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_group 
			  WHERE ot_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletGroup, query, outletGroupId, custId)
	if err != nil {
		log.Println("outletGroupRepository, FindOneByOutletGroupCodeAndCustId, err:", err.Error())
		return outletGroup, err
	}

	return outletGroup, nil
}

func (repository *outletGroupRepositoryImpl) FindOneByOutletGroupCodeAndCustId(outletGroupCode string, custId string) (model.OutletGroup, error) {
	outletGroup := model.OutletGroup{}
	query := `SELECT 
				cust_id, ot_grp_id, ot_grp_code,
				ot_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_group 
			  WHERE ot_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletGroup, query, outletGroupCode, custId)
	if err != nil {
		log.Println("outletGroupRepository, FindOneByOutletGroupCodeAndCustId, err:", err.Error())
		return outletGroup, err
	}

	return outletGroup, nil
}

func (repository *outletGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletGroup, int, int, error) {

	outletGroups := []model.OutletGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_grp_id, a.ot_grp_code, a.ot_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_outlet_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by  `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("outletGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletGroupRepository, count total, err:", err.Error())
		return outletGroups, 0, 0, err
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
		sortBy := `a.ot_grp_id`
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

	// log.Println("outletGroupRepository, querySelect:", querySelect)
	err = repository.Select(&outletGroups, querySelect)
	if err != nil {
		log.Println("outletGroupRepository, FindAllByCustId, err:", err.Error())
		return outletGroups, total, lastPage, err
	}

	return outletGroups, total, lastPage, nil
}

func (repository *outletGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletGroup, int, int, error) {

	outletGroups := []model.OutletGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_grp_id, a.ot_grp_code, a.ot_grp_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_outlet_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by  `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("outletGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletGroupRepository, count total, err:", err.Error())
		return outletGroups, 0, 0, err
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
		sortBy := `a.ot_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&outletGroups, querySelect)
	if err != nil {
		log.Println("outletGroupRepository, FindAllByCustId, err:", err.Error())
		return outletGroups, total, 1, err
	}

	return outletGroups, total, 1, nil
}

func (repository *outletGroupRepositoryImpl) Store(outletGroup model.OutletGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_outlet_group(
			cust_id, ot_grp_code, ot_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING ot_grp_id;`
	lastInsertId := outletGroup.OutletGroupId
	err := repository.QueryRow(query,
		outletGroup.CustId, outletGroup.OutletGroupCode, outletGroup.OutletGroupName,
		outletGroup.IsActive, outletGroup.CreatedBy, outletGroup.CreatedAt, outletGroup.UpdatedBy,
		outletGroup.UpdatedAt, outletGroup.IsDel, outletGroup.DeletedBy, outletGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("outletGroupRepository, Store, err:", err.Error())
		return outletGroup.OutletGroupId, err
	}
	return outletGroup.OutletGroupId, nil
}

func (repository *outletGroupRepositoryImpl) Update(outletGroupId int, request entity.UpdateOutletGroupRequest) error {
	var (
		r            model.OutletGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("outletGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_outlet_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND ot_grp_id = :ot_grp_id_old;`

	// log.Println("outletGroupRepository, Update, query:", query)

	sqlPatch.Args["ot_grp_id_old"] = outletGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("outletGroupRepository, Update, err:", err.Error())
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

func (repository *outletGroupRepositoryImpl) Delete(custId string, outletGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_outlet_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND ot_grp_id = :ot_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"ot_grp_id":  outletGroupId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("OutletGroupRepository, Delete, err:", err.Error())
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
