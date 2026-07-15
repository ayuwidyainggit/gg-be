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

type MissedPaymentReasonsRepository interface {
	FindOneByReasonAndCustId(reason, custId string) (model.MissedPaymentReasons, error)
	FindOneByMissedPaymentReasonsIdAndCustId(MissedPaymentReasonsId int, custId string) (model.MissedPaymentReasons, error)
	FindAllByCustId(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) (consPro []model.MissedPaymentReasons, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) (consPro []model.MissedPaymentReasons, total int, lastPage int, err error)
	Store(brand model.MissedPaymentReasons) (int, error)
	Update(MissedPaymentReasonsId int, request entity.UpdateMissedPaymentReasonsRequest) error
	Delete(custId string, MissedPaymentReasonsId int, deletedBy int64) error
}

func NewMissedPaymentReasonsRepository(db *sqlx.DB) MissedPaymentReasonsRepository {
	return &MissedPaymentReasonsRepositoryImpl{db}
}

type MissedPaymentReasonsRepositoryImpl struct {
	*sqlx.DB
}

func (repository *MissedPaymentReasonsRepositoryImpl) FindOneByReasonAndCustId(reason string, custId string) (model.MissedPaymentReasons, error) {
	MissedPaymentReasons := model.MissedPaymentReasons{}
	query := `SELECT 
				c.missed_payment_reasons_id, 
				c.missed_payment_reasons_code,
				c.missed_payment_reasons_name,
				COALESCE(c.image_url, '') as image_url,
				c.is_active, c.created_by,
				c.created_at, c.updated_by, c.updated_at,
				c.is_del, c.deleted_by, c.deleted_at
			  FROM mst.missed_payment_reasons c
			  WHERE c.is_del = false 
				AND c.missed_payment_reasons_code = $1 
				AND c.cust_id = $2`
	err := repository.Get(&MissedPaymentReasons, query, reason, custId)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return MissedPaymentReasons, err
	}

	return MissedPaymentReasons, nil
}

func (repository *MissedPaymentReasonsRepositoryImpl) FindOneByMissedPaymentReasonsIdAndCustId(MissedPaymentReasonsId int, custId string) (model.MissedPaymentReasons, error) {
	MissedPaymentReasons := model.MissedPaymentReasons{}
	query := `SELECT 
				c.missed_payment_reasons_id, 
				c.missed_payment_reasons_code,
				c.missed_payment_reasons_name,
				COALESCE(c.image_url, '') as image_url,
				c.is_active, c.created_by,
				c.created_at, c.updated_by, c.updated_at,
				c.is_del, c.deleted_by, c.deleted_at
			  FROM mst.missed_payment_reasons c
			  WHERE c.is_del = false 
				AND c.missed_payment_reasons_id = $1 
				AND c.cust_id = $2`
	err := repository.Get(&MissedPaymentReasons, query, MissedPaymentReasonsId, custId)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return MissedPaymentReasons, err
	}

	return MissedPaymentReasons, nil
}

func (repository *MissedPaymentReasonsRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) ([]model.MissedPaymentReasons, int, int, error) {

	MissedPaymentReasons := []model.MissedPaymentReasons{}
	selectCount := ` COUNT(c.*) AS total `
	selectField := `c.missed_payment_reasons_id,
					c.missed_payment_reasons_code ,c.missed_payment_reasons_name  `
	qWhere := ` WHERE c.is_del = false AND c.is_active = true 
				AND c.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (c.missed_payment_reasons_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.MissedPaymentReasonsId != 0 {
		qWhere += `AND c.missed_payment_reasons_id = ` + strconv.Itoa(dataFilter.MissedPaymentReasonsId) + ` `
	}

	qFrom := ` FROM mst.missed_payment_reasons c `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("MissedPaymentReasonsRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, count total, err:", err.Error())
		return MissedPaymentReasons, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`c.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `c.missed_payment_reasons_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("MissedPaymentReasonsRepository, querySelect:", querySelect)
	err = repository.Select(&MissedPaymentReasons, querySelect)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return MissedPaymentReasons, total, 1, err
	}

	return MissedPaymentReasons, total, 1, nil
}

func (repository *MissedPaymentReasonsRepositoryImpl) FindAllByCustId(dataFilter entity.MissedPaymentReasonsQueryFilter, custId string) ([]model.MissedPaymentReasons, int, int, error) {

	brands := []model.MissedPaymentReasons{}
	selectCount := ` COUNT(c.*) AS total `
	selectField := `c.missed_payment_reasons_id,
					c.missed_payment_reasons_code,
					c.missed_payment_reasons_name,
					COALESCE(c.image_url, '') as image_url,
					c.is_active, c.created_by,
					c.created_at, c.updated_by, c.updated_at,
					c.is_del, c.deleted_by, c.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE c.is_del = false 
				AND c.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (c.missed_payment_reasons_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND c.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND c.is_active = false `
		}
	}

	if dataFilter.MissedPaymentReasonsId != 0 {
		qWhere += `AND c.missed_payment_reasons_id = ` + strconv.Itoa(dataFilter.MissedPaymentReasonsId) + ` `
	}

	qFrom := ` FROM mst.missed_payment_reasons c
			   LEFT JOIN sys.m_user u ON u.user_id = c.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("MissedPaymentReasonsRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`c.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `c.missed_payment_reasons_id`
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

	// log.Println("MissedPaymentReasonsRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *MissedPaymentReasonsRepositoryImpl) Store(MissedPaymentReasons model.MissedPaymentReasons) (int, error) {
	query :=
		`INSERT INTO mst.missed_payment_reasons(
			cust_id, missed_payment_reasons_code, missed_payment_reasons_name, image_url, 
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11, $12 
		) RETURNING missed_payment_reasons_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		MissedPaymentReasons.CustId, MissedPaymentReasons.MissedPaymentReasonsCode, MissedPaymentReasons.MissedPaymentReasonsName, MissedPaymentReasons.ImageUrl, MissedPaymentReasons.IsActive,
		MissedPaymentReasons.CreatedBy, MissedPaymentReasons.CreatedAt, MissedPaymentReasons.UpdatedBy, MissedPaymentReasons.UpdatedAt,
		MissedPaymentReasons.IsDel, MissedPaymentReasons.DeletedBy, MissedPaymentReasons.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *MissedPaymentReasonsRepositoryImpl) Update(MissedPaymentReasonsId int, request entity.UpdateMissedPaymentReasonsRequest) error {
	var (
		r            model.MissedPaymentReasonsUpdate
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

	query := `UPDATE mst.missed_payment_reasons
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND missed_payment_reasons_id = :missed_payment_reasons_id;`

	sqlPatch.Args["missed_payment_reasons_id"] = MissedPaymentReasonsId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, Update, err:", err.Error())
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

func (repository *MissedPaymentReasonsRepositoryImpl) Delete(custId string, MissedPaymentReasonsId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.missed_payment_reasons
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND missed_payment_reasons_id = :missed_payment_reasons_id;`

	wMap := map[string]interface{}{
		"cust_id":                   custId,
		"missed_payment_reasons_id": MissedPaymentReasonsId,
		"deleted_by":                deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("MissedPaymentReasonsRepository, Delete, err:", err.Error())
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
