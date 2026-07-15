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

type PackTypeRepository interface {
	FindOneByPTypeIdAndCustId(pTypeId int, custId string) (model.PackType, error)
	FindOneByPTypeCodeAndCustId(pTypeCode, custId string) (model.PackType, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (packType []model.PackType, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (packType []model.PackType, total int, lastPage int, err error)
	Store(packType model.PackType) (int, error)
	Update(pTypeId int, request entity.UpdatePackTypeRequest) error
	Delete(custId string, pTypeId int, deletedBy int64) error
}

func NewPackTypeRepository(db *sqlx.DB) PackTypeRepository {
	return &packTypeRepositoryImpl{db}
}

type packTypeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *packTypeRepositoryImpl) FindOneByPTypeIdAndCustId(pTypeId int, custId string) (model.PackType, error) {
	packType := model.PackType{}
	query := `SELECT 
				cust_id, ptype_id, ptype_code,
				ptype_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_pack_type 
			  WHERE is_del = false 
			  AND ptype_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&packType, query, pTypeId, custId)
	if err != nil {
		log.Println("packTypeRepository, FindOneByPTypeIdAndCustId, err:", err.Error())
		return packType, err
	}

	return packType, nil
}

func (repository *packTypeRepositoryImpl) FindOneByPTypeCodeAndCustId(pTypeCode, custId string) (model.PackType, error) {
	packType := model.PackType{}
	query := `SELECT 
				cust_id, ptype_id, ptype_code,
				ptype_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_pack_type 
			  WHERE is_del = false 
			  AND ptype_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&packType, query, pTypeCode, custId)
	if err != nil {
		log.Println("packTypeRepository, FindOneByPTypeCode, err:", err.Error())
		return packType, err
	}

	return packType, nil
}

func (repository *packTypeRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PackType, int, int, error) {

	packTypes := []model.PackType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.ptype_id, p.ptype_code,
					p.ptype_name, p.is_active, p.created_by,
					p.created_at, p.updated_by, p.updated_at,
					p.is_del, p.deleted_by, p.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE p.is_del = false 
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.ptype_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.ptype_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_pack_type p
				LEFT JOIN sys.m_user u ON u.user_id = p.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("packTypeRepository, count total, err:", err.Error())
		return packTypes, 0, 0, err
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
		sortBy := `p.ptype_id`
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

	err = repository.Select(&packTypes, querySelect)
	if err != nil {
		log.Println("packTypeRepository, FindAllByCustId, err:", err.Error())
		return packTypes, total, lastPage, err
	}

	return packTypes, total, lastPage, nil
}

func (repository *packTypeRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PackType, int, int, error) {

	packTypes := []model.PackType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.ptype_id, p.ptype_code, p.ptype_name  `
	qWhere := ` WHERE p.is_del = false AND p.is_active = true
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.ptype_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.ptype_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_pack_type p `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("packTypeRepository, count total, err:", err.Error())
		return packTypes, 0, 0, err
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
		sortBy := `p.ptype_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&packTypes, querySelect)
	if err != nil {
		log.Println("packTypeRepository, FindAllByCustId, err:", err.Error())
		return packTypes, total, 1, err
	}

	return packTypes, total, 1, nil
}

func (repository *packTypeRepositoryImpl) Store(packType model.PackType) (int, error) {
	query :=
		`INSERT INTO mst.m_pack_type(
			cust_id, ptype_code, ptype_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING ptype_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		packType.CustId, packType.PtypeCode, packType.PtypeName, packType.IsActive,
		packType.CreatedBy, packType.CreatedAt, packType.UpdatedBy, packType.UpdatedAt,
		packType.IsDel, packType.DeletedBy, packType.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("packTypeRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *packTypeRepositoryImpl) Update(pTypeId int, request entity.UpdatePackTypeRequest) error {
	var (
		r            model.PackTypeUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("packTypeRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_pack_type
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND ptype_id = :ptype_id;`

	// log.Println("packTypeRepository, Update, query:", query)

	sqlPatch.Args["ptype_id"] = pTypeId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("packTypeRepository, Update, err:", err.Error())
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

func (repository *packTypeRepositoryImpl) Delete(custId string, pTypeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_pack_type
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND ptype_id = :ptype_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"ptype_id":   pTypeId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PackTypeRepository, Delete, err:", err.Error())
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
