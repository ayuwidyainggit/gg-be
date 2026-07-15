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

type ConvGroupRepository interface {
	FindOneByConvGroupIdAndCustId(convGroupId int, custId string) (model.ConvGroup, error)
	FindOneByConvGroupCodeAndCustId(convGroupCode string, custId string) (model.ConvGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ConvGroup, total int, lastPage int, err error)
	Store(convGroup model.ConvGroup) (int, error)
	Update(convGroupId int, request entity.UpdateConvGroupRequest) error
	Delete(custId string, convGroupId int, deletedBy int64) error
}

func NewConvGroupRepository(db *sqlx.DB) ConvGroupRepository {
	return &convGroupRepositoryImpl{db}
}

type convGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *convGroupRepositoryImpl) FindOneByConvGroupIdAndCustId(convGroupId int, custId string) (model.ConvGroup, error) {
	convGroup := model.ConvGroup{}
	query := `SELECT 
				cust_id, conv_grp_id, conv_grp_code,
				conv_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_conv_group 
			  WHERE conv_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&convGroup, query, convGroupId, custId)
	if err != nil {
		log.Println("convGroupRepository, FindOneByConvGroupCodeAndCustId, err:", err.Error())
		return convGroup, err
	}

	return convGroup, nil
}

func (repository *convGroupRepositoryImpl) FindOneByConvGroupCodeAndCustId(convGroupCode string, custId string) (model.ConvGroup, error) {
	convGroup := model.ConvGroup{}
	query := `SELECT 
				cust_id, conv_grp_id, conv_grp_code,
				conv_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_conv_group 
			  WHERE conv_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&convGroup, query, convGroupCode, custId)
	if err != nil {
		log.Println("convGroupRepository, FindOneByConvGroupCodeAndCustId, err:", err.Error())
		return convGroup, err
	}

	return convGroup, nil
}

func (repository *convGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ConvGroup, int, int, error) {

	convGroups := []model.ConvGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.conv_grp_id, a.conv_grp_code, a.conv_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.conv_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.conv_grp_name ILIKE '%` + dataFilter.Query + `%' )`
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
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_conv_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_conv_group a LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere

	// log.Println("convGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("convGroupRepository, count total, err:", err.Error())
		return convGroups, 0, 0, err
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
		sortBy := `conv_grp_id`
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

	// log.Println("convGroupRepository, querySelect:", querySelect)
	err = repository.Select(&convGroups, querySelect)
	if err != nil {
		log.Println("convGroupRepository, FindAllByCustId, err:", err.Error())
		return convGroups, total, lastPage, err
	}

	return convGroups, total, lastPage, nil
}

func (repository *convGroupRepositoryImpl) Store(convGroup model.ConvGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_conv_group(
			cust_id, conv_grp_code, conv_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING conv_grp_id;`
	lastInsertId := convGroup.ConvGrpId
	err := repository.QueryRow(query,
		convGroup.CustId, convGroup.ConvGrpCode, convGroup.ConvGrpName,
		convGroup.IsActive, convGroup.CreatedBy, convGroup.CreatedAt, convGroup.UpdatedBy,
		convGroup.UpdatedAt, convGroup.IsDel, convGroup.DeletedBy, convGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("convGroupRepository, Store, err:", err.Error())
		return convGroup.ConvGrpId, err
	}
	return convGroup.ConvGrpId, nil
}

func (repository *convGroupRepositoryImpl) Update(convGroupId int, request entity.UpdateConvGroupRequest) error {
	var (
		r            model.ConvGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("convGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_conv_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND conv_grp_id = :conv_grp_id_old;`

	// log.Println("convGroupRepository, Update, query:", query)

	sqlPatch.Args["conv_grp_id_old"] = convGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("convGroupRepository, Update, err:", err.Error())
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

func (repository *convGroupRepositoryImpl) Delete(custId string, convGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_conv_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND conv_grp_id = :conv_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"conv_grp_id": convGroupId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ConvGroupRepository, Delete, err:", err.Error())
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
