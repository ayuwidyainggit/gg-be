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

type IncentiveGroupRepository interface {
	FindOneByIncentiveGroupIdAndCustId(incentiveGroupId int, custId string) (model.IncentiveGroupList, error)
	FindOneByIncentiveGroupCodeAndCustId(incentiveGroupCode string, custId string) (model.IncentiveGroup, error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.IncentiveGroupList, total int, lastPage int, err error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.IncentiveGroupList, total int, lastPage int, err error)
	Store(incentiveGroup model.IncentiveGroup) (int, error)
	Update(incentiveGroupId int, request entity.UpdateIncentiveGroupRequest) error
	Delete(custId string, incentiveGroupId int, deletedBy int64) error
}

func NewIncentiveGroupRepository(db *sqlx.DB) IncentiveGroupRepository {
	return &incentiveGroupRepositoryImpl{db}
}

type incentiveGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *incentiveGroupRepositoryImpl) FindOneByIncentiveGroupIdAndCustId(incentiveGroupId int, custId string) (model.IncentiveGroupList, error) {
	incentiveGroup := model.IncentiveGroupList{}
	query := `SELECT 
				a.cust_id, a.inc_grp_id, a.inc_grp_code, 
				a.inc_grp_name, a.is_active, 
				u.user_fullname AS updated_by_name, 
				a.updated_at
			  FROM mst.m_inc_group a
			  LEFT JOIN sys.m_user u ON u.user_id = a.updated_by AND u.cust_id = a.cust_id
			  WHERE a.inc_grp_id = $1 
			  AND a.cust_id = $2`
	err := repository.Get(&incentiveGroup, query, incentiveGroupId, custId)
	if err != nil {
		log.Println("incentiveGroupRepository, FindOneByIncentiveGroupCodeAndCustId, err:", err.Error())
		return incentiveGroup, err
	}

	return incentiveGroup, nil
}

func (repository *incentiveGroupRepositoryImpl) FindOneByIncentiveGroupCodeAndCustId(incentiveGroupCode string, custId string) (model.IncentiveGroup, error) {
	incentiveGroup := model.IncentiveGroup{}
	query := `SELECT 
				a.cust_id, a.inc_grp_id, a.inc_grp_code,
				a.inc_grp_name, a.is_active, a.updated_by, a.updated_at
			  FROM mst.m_inc_group a
			  LEFT JOIN sys.m_user u ON u.user_id = a.updated_by
			  WHERE a.inc_grp_code = $1 
			  AND a.cust_id = $2`
	err := repository.Get(&incentiveGroup, query, incentiveGroupCode, custId)
	if err != nil {
		log.Println("incentiveGroupRepository, FindOneByIncentiveGroupCodeAndCustId, err:", err.Error())
		return incentiveGroup, err
	}

	return incentiveGroup, nil
}

func (repository *incentiveGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.IncentiveGroupList, int, int, error) {

	incentiveGroups := []model.IncentiveGroupList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.inc_grp_id, a.inc_grp_code, a.inc_grp_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' AND a.is_active = true `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.inc_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.inc_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_inc_group a `
	querySelect := `SELECT ` + selectField + ` FROM mst.m_inc_group a `

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("FindAllByCustIdLookupMode, count total, err:", err.Error())
		return incentiveGroups, 0, 0, err
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
		sortBy := `inc_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}

	lastPage := 1

	err = repository.Select(&incentiveGroups, querySelect)
	if err != nil {
		log.Println("incentiveGroupRepository, FindAllByCustId, err:", err.Error())
		return incentiveGroups, total, lastPage, err
	}

	return incentiveGroups, total, lastPage, nil
}

func (repository *incentiveGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.IncentiveGroupList, int, int, error) {

	incentiveGroups := []model.IncentiveGroupList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.inc_grp_id, a.inc_grp_code, a.inc_grp_name, a.is_active, 
					a.updated_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.inc_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.inc_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_inc_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_inc_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere

	// log.Println("incentiveGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("incentiveGroupRepository, count total, err:", err.Error())
		return incentiveGroups, 0, 0, err
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
		sortBy := `inc_grp_id`
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

	// log.Println("incentiveGroupRepository, querySelect:", querySelect)
	err = repository.Select(&incentiveGroups, querySelect)
	if err != nil {
		log.Println("incentiveGroupRepository, FindAllByCustId, err:", err.Error())
		return incentiveGroups, total, lastPage, err
	}

	return incentiveGroups, total, lastPage, nil
}

func (repository *incentiveGroupRepositoryImpl) Store(incentiveGroup model.IncentiveGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_inc_group(
			cust_id, inc_grp_code, inc_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING inc_grp_id;`
	lastInsertId := incentiveGroup.IncGrpID
	err := repository.QueryRow(query,
		incentiveGroup.CustID, incentiveGroup.IncGrpCode, incentiveGroup.IncGrpName,
		incentiveGroup.IsActive, incentiveGroup.CreatedBy, incentiveGroup.CreatedAt, incentiveGroup.UpdatedBy,
		incentiveGroup.UpdatedAt, incentiveGroup.IsDel, incentiveGroup.DeletedBy, incentiveGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("incentiveGroupRepository, Store, err:", err.Error())
		return incentiveGroup.IncGrpID, err
	}
	return incentiveGroup.IncGrpID, nil
}

func (repository *incentiveGroupRepositoryImpl) Update(incentiveGroupId int, request entity.UpdateIncentiveGroupRequest) error {
	var (
		r            model.IncentiveGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("incentiveGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_inc_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND inc_grp_id = :inc_grp_id_old;`

	// log.Println("incentiveGroupRepository, Update, query:", query)

	sqlPatch.Args["inc_grp_id_old"] = incentiveGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("incentiveGroupRepository, Update, err:", err.Error())
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

func (repository *incentiveGroupRepositoryImpl) Delete(custId string, incentiveGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_inc_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND inc_grp_id = :inc_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"inc_grp_id": incentiveGroupId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("IncentiveGroupRepository, Delete, err:", err.Error())
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
