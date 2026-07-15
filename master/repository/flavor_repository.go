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

type FlavorRepository interface {
	FindOneByFlavorIdAndCustId(flavorId int, custId string) (model.Flavor, error)
	FindOneByFlavorCodeAndCustId(flavorCode, custId string) (model.Flavor, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (flavor []model.Flavor, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (flavor []model.Flavor, total int, lastPage int, err error)
	Store(flavor model.Flavor) (int, error)
	Update(flavorId int, request entity.UpdateFlavorRequest) error
	Delete(custId string, flavorId int, deletedBy int64) error
}

func NewFlavorRepository(db *sqlx.DB) FlavorRepository {
	return &flavorRepositoryImpl{db}
}

type flavorRepositoryImpl struct {
	*sqlx.DB
}

func (repository *flavorRepositoryImpl) FindOneByFlavorIdAndCustId(flavorId int, custId string) (model.Flavor, error) {
	flavor := model.Flavor{}
	query := `SELECT 
				cust_id, flavor_id, flavor_code,
				flavor_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_flavor 
			  WHERE is_del = false 
			  AND flavor_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&flavor, query, flavorId, custId)
	if err != nil {
		log.Println("flavorRepository, FindOneByFlavorIdAndCustId, err:", err.Error())
		return flavor, err
	}

	return flavor, nil
}

func (repository *flavorRepositoryImpl) FindOneByFlavorCodeAndCustId(flavorCode, custId string) (model.Flavor, error) {
	flavor := model.Flavor{}
	query := `SELECT 
				cust_id, flavor_id, flavor_code,
				flavor_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_flavor 
			  WHERE is_del = false 
			  AND flavor_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&flavor, query, flavorCode, custId)
	if err != nil {
		log.Println("flavorRepository, FindOneByFlavorCode, err:", err.Error())
		return flavor, err
	}

	return flavor, nil
}

func (repository *flavorRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Flavor, int, int, error) {

	flavors := []model.Flavor{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` f.cust_id, f.flavor_id, f.flavor_code,
					f.flavor_name, f.is_active, f.created_by,
					f.created_at, f.updated_by, f.updated_at,
					f.is_del, f.deleted_by, f.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE f.is_del = false 
				AND f.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (f.flavor_code ILIKE '%` + dataFilter.Query + `%' 
					OR f.flavor_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND f.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND f.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_flavor f
				LEFT JOIN sys.m_user u ON u.user_id = f.updated_by `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("flavorRepository, count total, err:", err.Error())
		return flavors, 0, 0, err
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
		sortBy := `f.flavor_id`
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

	err = repository.Select(&flavors, querySelect)
	if err != nil {
		log.Println("flavorRepository, FindAllByCustId, err:", err.Error())
		return flavors, total, lastPage, err
	}

	return flavors, total, lastPage, nil
}

func (repository *flavorRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Flavor, int, int, error) {

	flavors := []model.Flavor{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` f.cust_id, f.flavor_id, f.flavor_code,
					f.flavor_name  `
	qWhere := ` WHERE f.is_del = false AND f.is_active = true
				AND f.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (f.flavor_code ILIKE '%` + dataFilter.Query + `%' 
					OR f.flavor_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_flavor f`

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("flavorRepository, count total, err:", err.Error())
		return flavors, 0, 0, err
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
		sortBy := `f.flavor_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&flavors, querySelect)
	if err != nil {
		log.Println("flavorRepository, FindAllByCustId, err:", err.Error())
		return flavors, total, 1, err
	}

	return flavors, total, 1, nil
}

func (repository *flavorRepositoryImpl) Store(flavor model.Flavor) (int, error) {
	query :=
		`INSERT INTO mst.m_flavor(
			cust_id, flavor_code, flavor_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING flavor_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		flavor.CustId, flavor.FlavorCode, flavor.FlavorName, flavor.IsActive,
		flavor.CreatedBy, flavor.CreatedAt, flavor.UpdatedBy, flavor.UpdatedAt,
		flavor.IsDel, flavor.DeletedBy, flavor.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("flavorRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *flavorRepositoryImpl) Update(flavorId int, request entity.UpdateFlavorRequest) error {
	var (
		r            model.FlavorUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("flavorRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_flavor
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND flavor_id = :flavor_id;`

	// log.Println("flavorRepository, Update, query:", query)

	sqlPatch.Args["flavor_id"] = flavorId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("flavorRepository, Update, err:", err.Error())
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

func (repository *flavorRepositoryImpl) Delete(custId string, flavorId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_flavor
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND flavor_id = :flavor_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"flavor_id":  flavorId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("FlavorRepository, Delete, err:", err.Error())
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
