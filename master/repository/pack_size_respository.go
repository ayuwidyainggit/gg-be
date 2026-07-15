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

type PackSizeRepository interface {
	FindOneByPSizeIdAndCustId(pSizeId int, custId string) (model.PackSize, error)
	FindOneByPSizeCodeAndCustId(pSizeCode, custId string) (model.PackSize, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (pSize []model.PackSize, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (pSize []model.PackSize, total int, lastPage int, err error)
	Store(packSize model.PackSize) (int, error)
	Update(pSizeId int, request entity.UpdatePackSizeRequest) error
	Delete(custId string, pSizeId int, deletedBy int64) error
}

func NewPackSizeRepository(db *sqlx.DB) PackSizeRepository {
	return &packSizeRepositoryImpl{db}
}

type packSizeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *packSizeRepositoryImpl) FindOneByPSizeIdAndCustId(pSizeId int, custId string) (model.PackSize, error) {
	packSize := model.PackSize{}
	query := `SELECT 
				cust_id, psize_id, psize_code,
				psize_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_pack_size 
			  WHERE is_del = false 
			  AND psize_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&packSize, query, pSizeId, custId)
	if err != nil {
		log.Println("packSizeRepository, FindOneByPSizeIdAndCustId, err:", err.Error())
		return packSize, err
	}

	return packSize, nil
}

func (repository *packSizeRepositoryImpl) FindOneByPSizeCodeAndCustId(pSizeCode, custId string) (model.PackSize, error) {
	packSize := model.PackSize{}
	query := `SELECT 
				cust_id, psize_id, psize_code,
				psize_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_pack_size 
			  WHERE is_del = false 
			  AND psize_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&packSize, query, pSizeCode, custId)
	if err != nil {
		log.Println("packSizeRepository, FindOneByPSizeCode, err:", err.Error())
		return packSize, err
	}

	return packSize, nil
}

func (repository *packSizeRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PackSize, int, int, error) {

	packSizes := []model.PackSize{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.psize_id, p.psize_code,
					p.psize_name, p.is_active, p.created_by,
					p.created_at, p.updated_by, p.updated_at,
					p.is_del, p.deleted_by, p.deleted_at,
					u.user_fullname AS updated_by_name `

	qWhere := ` WHERE p.is_del = false 
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.psize_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.psize_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND p.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND p.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_pack_size p
				LEFT JOIN sys.m_user u ON u.user_id = p.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("packSizeRepository, count total, err:", err.Error())
		return packSizes, 0, 0, err
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
		sortBy := `p.psize_id`
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

	err = repository.Select(&packSizes, querySelect)
	if err != nil {
		log.Println("packSizeRepository, FindAllByCustId, err:", err.Error())
		return packSizes, total, lastPage, err
	}

	return packSizes, total, lastPage, nil
}

func (repository *packSizeRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PackSize, int, int, error) {

	packSizes := []model.PackSize{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` p.cust_id, p.psize_id, p.psize_code,
					p.psize_name `

	qWhere := ` WHERE p.is_del = false AND p.is_active = true
				AND p.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.psize_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.psize_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_pack_size p `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("packSizeRepository, count total, err:", err.Error())
		return packSizes, 0, 0, err
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
		sortBy := `p.psize_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&packSizes, querySelect)
	if err != nil {
		log.Println("packSizeRepository, FindAllByCustId, err:", err.Error())
		return packSizes, total, 1, err
	}

	return packSizes, total, 1, nil
}

func (repository *packSizeRepositoryImpl) Store(packSize model.PackSize) (int, error) {
	query :=
		`INSERT INTO mst.m_pack_size(
			cust_id, psize_code, psize_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING psize_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		packSize.CustId, packSize.PsizeCode, packSize.PsizeName, packSize.IsActive,
		packSize.CreatedBy, packSize.CreatedAt, packSize.UpdatedBy, packSize.UpdatedAt,
		packSize.IsDel, packSize.DeletedBy, packSize.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("packSizeRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *packSizeRepositoryImpl) Update(pSizeId int, request entity.UpdatePackSizeRequest) error {
	var (
		r            model.PackSizeUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("packSizeRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_pack_size
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND psize_id = :psize_id;`

	// log.Println("packSizeRepository, Update, query:", query)

	sqlPatch.Args["psize_id"] = pSizeId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("packSizeRepository, Update, err:", err.Error())
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

func (repository *packSizeRepositoryImpl) Delete(custId string, pSizeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_pack_size
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND psize_id = :psize_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"psize_id":   pSizeId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PackSizeRepository, Delete, err:", err.Error())
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
