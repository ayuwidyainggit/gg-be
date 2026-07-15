package repository

import (
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PickupReasonRepository interface {
	FindAllByCustIdLookupMode(dataFilter entity.PickupReasonQueryFilter, custId string) (pickupReason []model.PickupReason, total int64, lastPage int, err error)
	FindAllByCustId(dataFilter entity.PickupReasonQueryFilter, custId string) (pickupReason []model.PickupReason, total int64, lastPage int, err error)
}

// type (
// 	pickupReasonRepositoryImpl struct {
// 		*gorm.DB
// 	}
// )

type pickupReasonRepositoryImpl struct {
	*sqlx.DB
}

func NewPickupReasonRepository(db *sqlx.DB) PickupReasonRepository {
	return &pickupReasonRepositoryImpl{db}
}

func (repository *pickupReasonRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.PickupReasonQueryFilter, custId string) ([]model.PickupReason, int64, int, error) {

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
	var total int64
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

func (repository *pickupReasonRepositoryImpl) FindAllByCustId(dataFilter entity.PickupReasonQueryFilter, custId string) ([]model.PickupReason, int64, int, error) {

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
	var total int64
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
