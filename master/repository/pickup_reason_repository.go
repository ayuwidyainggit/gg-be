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

type PickupReasonRepository interface {
	FindOneByPickupReasonIdAndCustId(pickupReasonId int, custId string) (model.PickupReason, error)
	FindOneByPickupReasonCodeAndCustId(pickupReasonCode string, custId string) (model.PickupReason, error)
	FindAllByCustIdLookupMode(dataFilter entity.PickupReasonQueryFilter, custId string) (pickupReason []model.PickupReason, total int, lastPage int, err error)
	FindAllByCustId(dataFilter entity.PickupReasonQueryFilter, custId string) (pickupReason []model.PickupReason, total int, lastPage int, err error)
	Store(pickupReason model.PickupReason) (int, error)
	Update(pickupReasonId int, request entity.UpdatePickupReasonRequest) error
	Delete(custId string, pickupReasonId int, deletedBy int64) error
}

func NewPickupReasonRepository(db *sqlx.DB) PickupReasonRepository {
	return &pickupReasonRepositoryImpl{db}
}

type pickupReasonRepositoryImpl struct {
	*sqlx.DB
}

func (repository *pickupReasonRepositoryImpl) FindOneByPickupReasonIdAndCustId(pickupReasonId int, custId string) (model.PickupReason, error) {
	pickupReason := model.PickupReason{}
	query := `SELECT 
				cust_id, pickup_reason_id, pickup_reason_code,
				pickup_reason_name, is_active, created_by,
				created_at, updated_by, updated_at, deleted_by, deleted_at
			  FROM mst.pickup_reasons 
			  WHERE pickup_reason_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&pickupReason, query, pickupReasonId, custId)
	if err != nil {
		log.Println("pickupReasonRepository, FindOneByPickupReasonIdAndCustId, err:", err.Error())
		return pickupReason, err
	}

	return pickupReason, nil
}

func (repository *pickupReasonRepositoryImpl) FindOneByPickupReasonCodeAndCustId(pickupReasonCode string, custId string) (model.PickupReason, error) {
	pickupReason := model.PickupReason{}
	query := `SELECT 
				cust_id, pickup_reason_id, pickup_reason_code,
				pickup_reason_name, is_active, created_by,
				created_at, updated_by, updated_at,
				deleted_by, deleted_at
			  FROM mst.pickup_reasons 
			  WHERE pickup_reason_code = $1 
			  AND deleted_at IS NULL
			  AND cust_id = $2`
	err := repository.Get(&pickupReason, query, pickupReasonCode, custId)
	if err != nil {
		log.Println("pickupReasonRepository, FindOneByPickupReasonCodeAndCustId, err:", err.Error())
		return pickupReason, err
	}

	return pickupReason, nil
}

func (repository *pickupReasonRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.PickupReasonQueryFilter, custId string) ([]model.PickupReason, int, int, error) {

	pickupReasons := []model.PickupReason{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` cust_id, 
					pickup_reason_id, 
					pickup_reason_code,
					pickup_reason_name `
	qWhere := ` WHERE is_active = true AND deleted_at IS NULL
				AND cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (pickup_reason_code ILIKE '%` + dataFilter.Query + `%' 
					OR pickup_reason_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := `FROM mst.pickup_reasons`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("pickupReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("pickupReasonRepository, count total, err:", err.Error())
		return pickupReasons, 0, 0, err
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
		sortBy := `pickup_reason_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("pickupReasonRepository, querySelect:", querySelect)
	err = repository.Select(&pickupReasons, querySelect)
	if err != nil {
		log.Println("pickupReasonRepository, FindAllByCustId, err:", err.Error())
		return pickupReasons, total, 1, err
	}

	return pickupReasons, total, 1, nil
}

func (repository *pickupReasonRepositoryImpl) FindAllByCustId(dataFilter entity.PickupReasonQueryFilter, custId string) ([]model.PickupReason, int, int, error) {

	pickupReasons := []model.PickupReason{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` t.cust_id, t.pickup_reason_id, t.pickup_reason_code,
					t.pickup_reason_name, t.is_active, t.created_by,
					t.created_at, t.updated_by, t.updated_at,
					t.deleted_by, t.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE t.deleted_at IS NULL
				AND t.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (pickup_reason_code ILIKE '%` + dataFilter.Query + `%' 
					OR pickup_reason_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND t.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND t.is_active = false `
		}
	}
	queryCount := `SELECT ` + selectCount + ` FROM mst.pickup_reasons t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.pickup_reasons t LEFT JOIN sys.m_user u ON u.user_id = t.updated_by ` + qWhere

	// log.Println("pickupReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("pickupReasonRepository, count total, err:", err.Error())
		return pickupReasons, 0, 0, err
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
		sortBy := `pickup_reason_id`
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

	// log.Println("pickupReasonRepository, querySelect:", querySelect)
	err = repository.Select(&pickupReasons, querySelect)
	if err != nil {
		log.Println("pickupReasonRepository, FindAllByCustId, err:", err.Error())
		return pickupReasons, total, lastPage, err
	}

	return pickupReasons, total, lastPage, nil
}

func (repository *pickupReasonRepositoryImpl) Store(pickupReason model.PickupReason) (int, error) {
	query :=
		`INSERT INTO mst.pickup_reasons(
			cust_id, pickup_reason_code, pickup_reason_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10
		) RETURNING pickup_reason_id;`
	lastInsertId := pickupReason.PickupReasonId
	err := repository.QueryRow(query,
		pickupReason.CustId, pickupReason.PickupReasonCode, pickupReason.PickupReasonName, pickupReason.IsActive,
		pickupReason.CreatedBy, pickupReason.CreatedAt, pickupReason.UpdatedBy, pickupReason.UpdatedAt,
		pickupReason.DeletedBy, pickupReason.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("pickupReasonRepository, Store, err:", err.Error())
		return pickupReason.PickupReasonId, err
	}
	return pickupReason.PickupReasonId, nil
}

func (repository *pickupReasonRepositoryImpl) Update(pickupReasonId int, request entity.UpdatePickupReasonRequest) error {
	var (
		r            model.PickupReasonUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("pickupReasonRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.pickup_reasons
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE deleted_at IS NULL
			  AND cust_id = :cust_id 
			  AND pickup_reason_id = :pickup_reason_id_old;`

	// log.Println("pickupReasonRepository, Update, query:", query)

	sqlPatch.Args["pickup_reason_id_old"] = pickupReasonId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("pickupReasonRepository, Update, err:", err.Error())
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

func (repository *pickupReasonRepositoryImpl) Delete(custId string, pickupReasonId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.pickup_reasons
			SET 
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE deleted_at IS NULL
			AND cust_id = :cust_id
			AND pickup_reason_id = :pickup_reason_id;`

	wMap := map[string]interface{}{
		"cust_id":          custId,
		"pickup_reason_id": pickupReasonId,
		"deleted_by":       deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PickupReasonRepository, Delete, err:", err.Error())
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
