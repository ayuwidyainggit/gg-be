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

type DistPriceGroupRepository interface {
	FindOneByDistPriceGroupIdAndCustId(distPriceGroupId int, custId string) (model.DistPriceGroup, error)
	FindOneByDistPriceGroupCodeAndCustId(distPriceGroupCode string, custId string) (model.DistPriceGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.DistPriceGroup, total int, lastPage int, err error)
	Store(DistPriceGroup model.DistPriceGroup) (int, error)
	Update(distPriceGroupId int, request entity.UpdateDistPriceGroupRequest) error
	Delete(custId string, DistPriceGroupId int, deletedBy int64) error
}

func NewDistPriceGroupRepository(db *sqlx.DB) DistPriceGroupRepository {
	return &DistPriceGroupRepositoryImpl{db}
}

type DistPriceGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *DistPriceGroupRepositoryImpl) FindOneByDistPriceGroupIdAndCustId(distPriceGroupId int, custId string) (model.DistPriceGroup, error) {
	DistPriceGroup := model.DistPriceGroup{}
	query := `SELECT 
				cust_id, dist_price_grp_id, dist_price_grp_code,
				dist_price_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_dist_price_group 
			  WHERE dist_price_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&DistPriceGroup, query, distPriceGroupId, custId)
	if err != nil {
		log.Println("DistPriceGroupRepository, FindOneByDistPriceGroupIdAndCustId, err:", err.Error())
		return DistPriceGroup, err
	}

	return DistPriceGroup, nil
}

func (repository *DistPriceGroupRepositoryImpl) FindOneByDistPriceGroupCodeAndCustId(distPriceGroupCode string, custId string) (model.DistPriceGroup, error) {
	DistPriceGroup := model.DistPriceGroup{}
	query := `SELECT 
				cust_id, dist_price_grp_id, dist_price_grp_code,
				dist_price_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_dist_price_group 
			  WHERE dist_price_grp_code = $1 
			  AND cust_id = $2 AND is_del = false`
	err := repository.Get(&DistPriceGroup, query, distPriceGroupCode, custId)
	if err != nil {
		log.Println("DistPriceGroupRepository, FindOneByDistPriceGroupCodeAndCustId, err:", err.Error())
		return DistPriceGroup, err
	}

	return DistPriceGroup, nil
}

func (repository *DistPriceGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.DistPriceGroup, int, int, error) {

	DistPriceGroups := []model.DistPriceGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.dist_price_grp_id, a.dist_price_grp_code, a.dist_price_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` LEFT JOIN sys.m_user u ON u.user_id = a.updated_by 
				WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.dist_price_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.dist_price_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	queryFrom := ` FROM mst.m_dist_price_group a `
	queryCount := `SELECT ` + selectCount + ` ` + queryFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + queryFrom + ` ` + qWhere

	// log.Println("DistPriceGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("DistPriceGroupRepository, count total, err:", err.Error())
		return DistPriceGroups, 0, 0, err
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
		sortBy := `a.dist_price_grp_id`
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

	// log.Println("DistPriceGroupRepository, querySelect:", querySelect)
	err = repository.Select(&DistPriceGroups, querySelect)
	if err != nil {
		log.Println("DistPriceGroupRepository, FindAllByCustId, err:", err.Error())
		return DistPriceGroups, total, lastPage, err
	}

	return DistPriceGroups, total, lastPage, nil
}

func (repository *DistPriceGroupRepositoryImpl) Store(DistPriceGroup model.DistPriceGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_dist_price_group(
			cust_id, dist_price_grp_code, dist_price_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING dist_price_grp_id;`
	lastInsertId := DistPriceGroup.DistPriceGrpId
	err := repository.QueryRow(query,
		DistPriceGroup.CustId, DistPriceGroup.DistPriceGrpCode, DistPriceGroup.DistPriceGrpName,
		DistPriceGroup.IsActive, DistPriceGroup.CreatedBy, DistPriceGroup.CreatedAt, DistPriceGroup.UpdatedBy,
		DistPriceGroup.UpdatedAt, DistPriceGroup.IsDel, DistPriceGroup.DeletedBy, DistPriceGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("DistPriceGroupRepository, Store, err:", err.Error())
		return DistPriceGroup.DistPriceGrpId, err
	}
	return DistPriceGroup.DistPriceGrpId, nil
}

func (repository *DistPriceGroupRepositoryImpl) Update(distPriceGroupId int, request entity.UpdateDistPriceGroupRequest) error {
	var (
		r            model.DistPriceGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("DistPriceGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_dist_price_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND dist_price_grp_id = :dist_price_grp_id_old;`

	// log.Println("DistPriceGroupRepository, Update, query:", query)

	sqlPatch.Args["dist_price_grp_id_old"] = distPriceGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("DistPriceGroupRepository, Update, err:", err.Error())
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

func (repository *DistPriceGroupRepositoryImpl) Delete(custId string, distPriceGrpId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_dist_price_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND dist_price_grp_id = :dist_price_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":           custId,
		"dist_price_grp_id": distPriceGrpId,
		"deleted_by":        deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("DistPriceGroupRepository, Delete, err:", err.Error())
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
