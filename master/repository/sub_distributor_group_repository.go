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

type SubDistributorGroupRepository interface {
	FindOneBySubDistributorGroupCodeAndCustId(subDistributorGroupCode, custId string) (model.SubDistributorGroup, error)
	Store(subDistributorGroup model.SubDistributorGroup) (int, error)
	FindAllByCustId(dataFilter entity.SubDistributorGroupQueryFilter, custId string) ([]model.SubDistributorGroup, int, int, error)
	FindAllByCustIdLookupMode(dataFilter entity.SubDistributorGroupQueryFilter, custId string) ([]model.SubDistributorGroup, int, int, error)
	FindOneBySubDistributorGroupIdAndCustId(subDistributorGroupId int, custId string) (model.SubDistributorGroup, error)
	Update(channelId int, request entity.SubDistributorGroupUpdateRequest) error
	Delete(custId string, ubDistributorGroupID int, deletedBy int64) error
}

type SubDistributorGroupRepositoryImpl struct {
	*sqlx.DB
}

func NewSubDistributorGroupRepository(db *sqlx.DB) *SubDistributorGroupRepositoryImpl {
	return &SubDistributorGroupRepositoryImpl{db}
}

func (repository *SubDistributorGroupRepositoryImpl) FindOneBySubDistributorGroupCodeAndCustId(subDistributorGroupCode, custId string) (model.SubDistributorGroup, error) {
	subDistributorGroup := model.SubDistributorGroup{}
	query := `SELECT 
				*
			  FROM mst.m_sub_distributor_group 
			  WHERE is_del = false 
			  AND sub_distributor_group_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&subDistributorGroup, query, subDistributorGroupCode, custId)
	if err != nil {
		log.Println("SubDistributorGroupRepository, FindOneBySubDistributorGroupCodeAndCustId, err:", err.Error())
		return subDistributorGroup, err
	}

	return subDistributorGroup, nil
}

func (repository *SubDistributorGroupRepositoryImpl) Store(subDistributorGroup model.SubDistributorGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_sub_distributor_group(
			cust_id, sub_distributor_group_code, sub_distributor_group_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11
		) RETURNING sub_distributor_group_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		subDistributorGroup.CustID, subDistributorGroup.SubDistributorGroupCode, subDistributorGroup.SubDistributorGroupName,
		subDistributorGroup.IsActive, subDistributorGroup.CreatedBy, subDistributorGroup.CreatedAt, subDistributorGroup.UpdatedBy, subDistributorGroup.UpdatedAt,
		subDistributorGroup.IsDel, subDistributorGroup.DeletedBy, subDistributorGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("subDistributorGroupRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *SubDistributorGroupRepositoryImpl) FindAllByCustId(dataFilter entity.SubDistributorGroupQueryFilter, custId string) ([]model.SubDistributorGroup, int, int, error) {

	subDistributorGroups := []model.SubDistributorGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.*, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sub_distributor_group_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sub_distributor_group_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_sub_distributor_group a
	LEFT JOIN sys.m_user u ON u.user_id = a.updated_by `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere
	fmt.Println(querySelect)
	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("SubDistributorGroupRepository, count total, err:", err.Error())
		return subDistributorGroups, 0, 0, err
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
		sortBy := `a.sub_distributor_group_id`
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

	// log.Println("districtRepository, querySelect:", querySelect)
	err = repository.Select(&subDistributorGroups, querySelect)
	if err != nil {
		log.Println("SubDistributorGroupRepository, FindAllByCustId, err:", err.Error())
		return subDistributorGroups, total, lastPage, err
	}

	return subDistributorGroups, total, lastPage, nil
}

func (repository *SubDistributorGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.SubDistributorGroupQueryFilter, custId string) ([]model.SubDistributorGroup, int, int, error) {

	districts := []model.SubDistributorGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.sub_distributor_group_id, a.sub_distributor_group_code, a.sub_distributor_group_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.district_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.district_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_sub_distributor_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("districtRepository, count total, err:", err.Error())
		return districts, 0, 0, err
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
		sortBy := `a.sub_distributor_group_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&districts, querySelect)
	if err != nil {
		log.Println("districtRepository, FindAllByCustId, err:", err.Error())
		return districts, total, 1, err
	}

	return districts, total, 1, nil
}

func (repository *SubDistributorGroupRepositoryImpl) FindOneBySubDistributorGroupIdAndCustId(subDistributorGroupId int, custId string) (model.SubDistributorGroup, error) {
	subDistributorGroup := model.SubDistributorGroup{}
	query := `SELECT 
				b.*, u.user_fullname AS updated_by_name
			  FROM mst.m_sub_distributor_group b
			  LEFT JOIN sys.m_user u ON u.user_id = b.updated_by 
			  WHERE b.is_del = false 
				AND b.sub_distributor_group_id = $1 
				AND b.cust_id = $2`
	err := repository.Get(&subDistributorGroup, query, subDistributorGroupId, custId)
	if err != nil {
		log.Println("SubDistributorGroupRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return subDistributorGroup, err
	}

	return subDistributorGroup, nil
}

func (repository *SubDistributorGroupRepositoryImpl) Update(subDistributorGroupId int, request entity.SubDistributorGroupUpdateRequest) error {
	var (
		r            model.SubDistributorGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("districtRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sub_distributor_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sub_distributor_group_id = :sub_distributor_group_id_old;`

	log.Println("subDistributorGroupRepository, Update, query:", query)

	sqlPatch.Args["sub_distributor_group_id_old"] = subDistributorGroupId
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("subDistributorGroupRepository, Update, err:", err.Error())
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

func (repository *SubDistributorGroupRepositoryImpl) Delete(custId string, subDistributorGroupID int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sub_distributor_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND sub_distributor_group_id = :sub_distributor_group_id;`

	wMap := map[string]interface{}{
		"cust_id":                  custId,
		"sub_distributor_group_id": subDistributorGroupID,
		"deleted_by":               deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("subDistributorGroupRepository, Delete, err:", err.Error())
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
