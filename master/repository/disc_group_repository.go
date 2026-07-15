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

type DiscGroupRepository interface {
	FindOneByDiscGroupIdAndCustId(discGroupId int, custId string) (model.DiscGroup, error)
	FindOneByDiscGroupCodeAndCustId(discGroupCode string, custId string) (model.DiscGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.DiscGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.DiscGroup, total int, lastPage int, err error)
	Store(discGroup model.DiscGroup) (int, error)
	Update(discGroupId int, request entity.UpdateDiscGroupRequest) error
	Delete(custId string, discGroupId int, deletedBy int64) error
}

func NewDiscGroupRepository(db *sqlx.DB) DiscGroupRepository {
	return &discGroupRepositoryImpl{db}
}

type discGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *discGroupRepositoryImpl) FindOneByDiscGroupIdAndCustId(discGroupId int, custId string) (model.DiscGroup, error) {
	discGroup := model.DiscGroup{}
	query := `SELECT 
				cust_id, disc_grp_id, disc_grp_code,
				disc_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_disc_group 
			  WHERE disc_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&discGroup, query, discGroupId, custId)
	if err != nil {
		log.Println("discGroupRepository, FindOneByDiscGroupCodeAndCustId, err:", err.Error())
		return discGroup, err
	}

	return discGroup, nil
}

func (repository *discGroupRepositoryImpl) FindOneByDiscGroupCodeAndCustId(discGroupCode string, custId string) (model.DiscGroup, error) {
	discGroup := model.DiscGroup{}
	query := `SELECT 
				cust_id, disc_grp_id, disc_grp_code,
				disc_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_disc_group 
			  WHERE disc_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&discGroup, query, discGroupCode, custId)
	if err != nil {
		log.Println("discGroupRepository, FindOneByDiscGroupCodeAndCustId, err:", err.Error())
		return discGroup, err
	}

	return discGroup, nil
}

func (repository *discGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.DiscGroup, int, int, error) {

	discGroups := []model.DiscGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.disc_grp_id, a.disc_grp_code, a.disc_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.disc_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.disc_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_disc_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("discGroupRepository, count total, err:", err.Error())
		return discGroups, 0, 0, err
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
		sortBy := `a.disc_grp_id`
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

	err = repository.Select(&discGroups, querySelect)
	if err != nil {
		log.Println("discGroupRepository, FindAllByCustId, err:", err.Error())
		return discGroups, total, lastPage, err
	}

	return discGroups, total, lastPage, nil
}

func (repository *discGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.DiscGroup, int, int, error) {

	discGroups := []model.DiscGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.disc_grp_id, a.disc_grp_code, a.disc_grp_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.disc_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.disc_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_disc_group a  `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("discGroupRepository, count total, err:", err.Error())
		return discGroups, 0, 0, err
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
		sortBy := `a.disc_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&discGroups, querySelect)
	if err != nil {
		log.Println("discGroupRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return discGroups, total, 1, err
	}

	return discGroups, total, 1, nil
}

func (repository *discGroupRepositoryImpl) Store(discGroup model.DiscGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_disc_group(
			cust_id, disc_grp_code, disc_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING disc_grp_id;`
	lastInsertId := discGroup.DiscGrpId
	err := repository.QueryRow(query,
		discGroup.CustId, discGroup.DiscGrpCode, discGroup.DiscGrpName,
		discGroup.IsActive, discGroup.CreatedBy, discGroup.CreatedAt, discGroup.UpdatedBy,
		discGroup.UpdatedAt, discGroup.IsDel, discGroup.DeletedBy, discGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("discGroupRepository, Store, err:", err.Error())
		return discGroup.DiscGrpId, err
	}
	return discGroup.DiscGrpId, nil
}

func (repository *discGroupRepositoryImpl) Update(discGroupId int, request entity.UpdateDiscGroupRequest) error {
	var (
		r            model.DiscGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("discGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_disc_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND disc_grp_id = :disc_grp_id_old;`

	// log.Println("discGroupRepository, Update, query:", query)

	sqlPatch.Args["disc_grp_id_old"] = discGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("discGroupRepository, Update, err:", err.Error())
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

func (repository *discGroupRepositoryImpl) Delete(custId string, discGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_disc_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND disc_grp_id = :disc_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"disc_grp_id": discGroupId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("DiscGroupRepository, Delete, err:", err.Error())
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
