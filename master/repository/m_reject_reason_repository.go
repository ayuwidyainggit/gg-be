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

type RejectReasonRepository interface {
	FindOneParentCustId(distCustId string) (model.MCustomer, error)
	FindOneByRejectReasonIdAndCustId(rejectReasonId int, custId string) (model.RejectReason, error)
	FindOneByRejectReasonCodeAndCustId(rejectReasonCode, custId string) (model.RejectReason, error)
	FindAllByCustId(dataFilter entity.RejectReasonQueryFilter, custId string) (consPro []model.RejectReason, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.RejectReasonQueryFilter, custId string) (consPro []model.RejectReason, total int, lastPage int, err error)
	Store(brand model.RejectReason) (int, error)
	Update(rejectReasonId int, request entity.UpdateRejectReasonRequest) error
	Delete(custId string, rejectReasonId int, deletedBy int64) error
}

func NewRejectReasonRepository(db *sqlx.DB) RejectReasonRepository {
	return &RejectReasonRepositoryImpl{db}
}

type RejectReasonRepositoryImpl struct {
	*sqlx.DB
}

func (repository *RejectReasonRepositoryImpl) FindOneParentCustId(distCustId string) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE cust_id = $1`
	err := repository.Get(&mCustomer, query, distCustId)
	if err != nil {
		log.Println("RejectReasonRepositoryImpl, FindOneParentCustId, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func (repository *RejectReasonRepositoryImpl) FindOneByRejectReasonIdAndCustId(rejectReasonId int, custId string) (model.RejectReason, error) {
	rejectReason := model.RejectReason{}
	query := `SELECT 
				rr.reject_reason_id, 
				rr.reject_reason_code,
				rr.reject_reason_name, 
				rr.is_active, rr.created_by,
				rr.created_at, rr.updated_by, rr.updated_at,
				rr.is_del, rr.deleted_by, rr.deleted_at
			  FROM mst.m_reject_reason rr
			  WHERE rr.is_del = false 
				AND rr.reject_reason_id = $1 
				AND rr.cust_id = $2`
	err := repository.Get(&rejectReason, query, rejectReasonId, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return rejectReason, err
	}

	return rejectReason, nil
}

func (repository *RejectReasonRepositoryImpl) FindOneByRejectReasonCodeAndCustId(rejectReasonCode, custId string) (model.RejectReason, error) {
	rejectReason := model.RejectReason{}
	query := `SELECT 
				rr.reject_reason_id, 
				rr.reject_reason_code,
				rr.reject_reason_name, 
				rr.is_active, rr.created_by,
				rr.created_at, rr.updated_by, rr.updated_at,
				rr.is_del, rr.deleted_by, rr.deleted_at
			  FROM mst.m_reject_reason rr
			  WHERE rr.is_del = false 
				AND rr.reject_reason_code = $1 
				AND rr.cust_id = $2`
	err := repository.Get(&rejectReason, query, rejectReasonCode, custId)
	if err != nil {
		log.Println("RejectReasonRepositoryImpl, FindOneByRejectReasonIdAndCustId, err:", err.Error())
		return rejectReason, err
	}

	return rejectReason, nil
}

func (repository *RejectReasonRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.RejectReasonQueryFilter, custId string) ([]model.RejectReason, int, int, error) {

	rejectReason := []model.RejectReason{}
	selectCount := ` COUNT(rr.*) AS total `
	selectField := `rr.reject_reason_id,
					rr.reject_reason_name,
					rr.reject_reason_code  `
	qWhere := ` WHERE rr.is_del = false AND rr.is_active = true 
				AND rr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (rr.reject_reason_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.RejectReasonId != 0 {
		qWhere += `AND rr.reject_reason_id = ` + strconv.Itoa(dataFilter.RejectReasonId) + ` `
	}

	qFrom := ` FROM mst.m_reject_reason rr `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("rejectReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("rejectReasonRepository, count total, err:", err.Error())
		return rejectReason, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`rr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `rr.reject_reason_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("rejectReasonRepository, querySelect:", querySelect)
	err = repository.Select(&rejectReason, querySelect)
	if err != nil {
		log.Println("rejectReasonRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return rejectReason, total, 1, err
	}

	return rejectReason, total, 1, nil
}

func (repository *RejectReasonRepositoryImpl) FindAllByCustId(dataFilter entity.RejectReasonQueryFilter, custId string) ([]model.RejectReason, int, int, error) {

	brands := []model.RejectReason{}
	selectCount := ` COUNT(rr.*) AS total `
	selectField := `rr.reject_reason_id,
					rr.reject_reason_code,
					rr.reject_reason_name, 
					rr.is_active, rr.created_by,
					rr.created_at, rr.updated_by, rr.updated_at,
					rr.is_del, rr.deleted_by, rr.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE rr.is_del = false 
				AND rr.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (rr.reject_reason_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND rr.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND rr.is_active = false `
		}
	}

	if dataFilter.RejectReasonId != 0 {
		qWhere += `AND rr.reject_reason_id = ` + strconv.Itoa(dataFilter.RejectReasonId) + ` `
	}

	qFrom := ` FROM mst.m_reject_reason rr
			   LEFT JOIN sys.m_user u ON u.user_id = rr.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("rejectReasonRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("rejectReasonRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`rr.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `rr.reject_reason_id`
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

	// log.Println("rejectReasonRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("rejectReasonRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *RejectReasonRepositoryImpl) Store(rejectReason model.RejectReason) (int, error) {
	query :=
		`INSERT INTO mst.m_reject_reason(
			cust_id, reject_reason_name, reject_reason_code,
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11
		) RETURNING reject_reason_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		rejectReason.CustId, rejectReason.RejectReasonName, rejectReason.RejectReasonCode, rejectReason.IsActive,
		rejectReason.CreatedBy, rejectReason.CreatedAt, rejectReason.UpdatedBy, rejectReason.UpdatedAt,
		rejectReason.IsDel, rejectReason.DeletedBy, rejectReason.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("rejectReasonRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *RejectReasonRepositoryImpl) Update(rejectReasonId int, request entity.UpdateRejectReasonRequest) error {
	var (
		r            model.RejectReasonUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_reject_reason
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND reject_reason_id = :reject_reason_id;`

	sqlPatch.Args["reject_reason_id"] = rejectReasonId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("rejectReasonRepository, Update, err:", err.Error())
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

func (repository *RejectReasonRepositoryImpl) Delete(custId string, rejectReasonId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_reject_reason
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND reject_reason_id = :reject_reason_id;`

	wMap := map[string]interface{}{
		"cust_id":          custId,
		"reject_reason_id": rejectReasonId,
		"deleted_by":       deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("RejectReasonRepository, Delete, err:", err.Error())
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
