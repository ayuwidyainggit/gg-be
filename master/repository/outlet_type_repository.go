package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OutletTypeRepository interface {
	FindOneByOutletTypeIdAndCustId(outletTypeId int64, custId string) (model.OutletType, error)
	FindOneByOutletTypeCodeAndCustId(outletTypeCode string, custId string) (model.OutletType, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.OutletType, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.OutletType, total int, lastPage int, err error)
	Store(outletType model.OutletType) (int, error)
	Update(outletTypeId int, request entity.UpdateOutletTypeRequest) error
	Delete(custId string, outletTypeId int, deletedBy int64) error
	FindAllByOutletTypeIDsAndCustID(outletTypeIDs []int, custID string) ([]model.OutletType, error)
}

func NewOutletTypeRepository(db *sqlx.DB) OutletTypeRepository {
	return &outletTypeRepositoryImpl{db}
}

type outletTypeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *outletTypeRepositoryImpl) FindOneByOutletTypeIdAndCustId(outletTypeId int64, custId string) (model.OutletType, error) {
	log.Info("outletTypeID:", outletTypeId)
	log.Info("custId:", custId)
	outletType := model.OutletType{}
	query := `SELECT 
				cust_id, ot_type_id, ot_type_code,
				ot_type_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_type 
			  WHERE ot_type_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletType, query, outletTypeId, custId)
	if err != nil {
		log.Error("outletTypeRepository, FindOneByOutletTypeCodeAndCustId, err:", err.Error())
		return outletType, err
	}

	return outletType, nil
}

func (repository *outletTypeRepositoryImpl) FindOneByOutletTypeCodeAndCustId(outletTypeCode string, custId string) (model.OutletType, error) {
	outletType := model.OutletType{}
	query := `SELECT 
				cust_id, ot_type_id, ot_type_code,
				ot_type_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_type 
			  WHERE ot_type_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&outletType, query, outletTypeCode, custId)
	if err != nil {
		log.Error("outletTypeRepository, FindOneByOutletTypeCodeAndCustId, err:", err.Error())
		return outletType, err
	}

	return outletType, nil
}

func (repository *outletTypeRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletType, int, int, error) {

	outletTypes := []model.OutletType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_type_id, a.ot_type_code, a.ot_type_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_type_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_outlet_type a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Error("outletTypeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("outletTypeRepository, count total, err:", err.Error())
		return outletTypes, 0, 0, err
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
		sortBy := `a.ot_type_id`
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

	// log.Error("outletTypeRepository, querySelect:", querySelect)
	err = repository.Select(&outletTypes, querySelect)
	if err != nil {
		log.Error("outletTypeRepository, FindAllByCustId, err:", err.Error())
		return outletTypes, total, lastPage, err
	}

	return outletTypes, total, lastPage, nil
}

func (repository *outletTypeRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.OutletType, int, int, error) {

	outletTypes := []model.OutletType{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.ot_type_id, a.ot_type_code, a.ot_type_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.ot_type_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.ot_type_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_outlet_type a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Error("outletTypeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("outletTypeRepository, count total, err:", err.Error())
		return outletTypes, 0, 0, err
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
		sortBy := `a.ot_type_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&outletTypes, querySelect)
	if err != nil {
		log.Error("outletTypeRepository, FindAllByCustId, err:", err.Error())
		return outletTypes, total, 1, err
	}

	return outletTypes, total, 1, nil
}

func (repository *outletTypeRepositoryImpl) Store(outletType model.OutletType) (int, error) {
	query :=
		`INSERT INTO mst.m_outlet_type(
			cust_id, ot_type_code, ot_type_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING ot_type_id;`
	lastInsertId := outletType.OtTypeId
	err := repository.QueryRow(query,
		outletType.CustId, outletType.OtTypeCode, outletType.OtTypeName,
		outletType.IsActive, outletType.CreatedBy, outletType.CreatedAt, outletType.UpdatedBy,
		outletType.UpdatedAt, outletType.IsDel, outletType.DeletedBy, outletType.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Error("outletTypeRepository, Store, err:", err.Error())
		return outletType.OtTypeId, err
	}
	return outletType.OtTypeId, nil
}

func (repository *outletTypeRepositoryImpl) Update(outletTypeId int, request entity.UpdateOutletTypeRequest) error {
	var (
		r            model.OutletTypeUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("outletTypeRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_outlet_type
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND ot_type_id = :ot_type_id_old;`

	// log.Error("outletTypeRepository, Update, query:", query)

	sqlPatch.Args["ot_type_id_old"] = outletTypeId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error("outletTypeRepository, Update, err:", err.Error())
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

func (repository *outletTypeRepositoryImpl) Delete(custId string, outletTypeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_outlet_type
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND ot_type_id = :ot_type_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"ot_type_id": outletTypeId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Error("OutletTypeRepository, Delete, err:", err.Error())
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

func (repository *outletTypeRepositoryImpl) FindAllByOutletTypeIDsAndCustID(outletTypeIDs []int, custID string) ([]model.OutletType, error) {
	outletTypes := []model.OutletType{}
	query := `SELECT 
				cust_id, ot_type_id, ot_type_code,
				ot_type_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_outlet_type 
			  WHERE ot_type_id IN $1 
			  AND cust_id = $2`
	err := repository.Get(&outletTypes, query, outletTypeIDs, custID)
	if err != nil {
		log.Error(err.Error())
		return outletTypes, err
	}

	return outletTypes, nil
}