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

type TopRepository interface {
	FindOneByTopAndCustId(top int, custId string) (model.Top, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Top, total int, lastPage int, err error)
	Store(top model.Top) (int, error)
	Update(topId int, request entity.UpdateTopRequest) error
	Delete(custId string, topId int, deletedBy int64) error
}

func NewTopRepository(db *sqlx.DB) TopRepository {
	return &topRepositoryImpl{db}
}

type topRepositoryImpl struct {
	*sqlx.DB
}

func (repository *topRepositoryImpl) FindOneByTopAndCustId(topPk int, custId string) (model.Top, error) {
	top := model.Top{}
	query := `SELECT 
				cust_id, top, 
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_top 
			  WHERE top = $1 
			  AND cust_id = $2`
	err := repository.Get(&top, query, topPk, custId)
	if err != nil {
		log.Println("topRepository, FindOneByTopAndCustId, err:", err.Error())
		return top, err
	}

	return top, nil
}

func (repository *topRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Top, int, int, error) {

	tops := []model.Top{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.top, a.is_active, a.created_by, a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.top ILIKE '%` + dataFilter.Query + `%')`
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
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_top a 
				   LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_top a 
					LEFT JOIN sys.m_user u ON u.user_id = a.updated_by` + qWhere

	// log.Println("topRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("topRepository, count total, err:", err.Error())
		return tops, 0, 0, err
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
		sortBy := `top`
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

	// log.Println("topRepository, querySelect:", querySelect)
	err = repository.Select(&tops, querySelect)
	if err != nil {
		log.Println("topRepository, FindAllByCustId, err:", err.Error())
		return tops, total, lastPage, err
	}

	return tops, total, lastPage, nil
}

func (repository *topRepositoryImpl) Store(top model.Top) (int, error) {
	query :=
		`INSERT INTO mst.m_top(
			cust_id, top, is_active, 
			created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, 
			deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, 
			$7, $8, $9, 
			$10
		) RETURNING top;`
	lastInsertId := top.Top
	err := repository.QueryRow(query,
		top.CustId, top.Top, top.IsActive,
		top.CreatedBy, top.CreatedAt, top.UpdatedBy,
		top.UpdatedAt, top.IsDel, top.DeletedBy,
		top.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("topRepository, Store, err:", err.Error())
		return top.Top, err
	}
	return top.Top, nil
}

func (repository *topRepositoryImpl) Update(topId int, request entity.UpdateTopRequest) error {
	var (
		r            model.TopUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("topRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_top
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND top = :top_id_old;`

	// log.Println("topRepository, Update, query:", query)

	sqlPatch.Args["top_id_old"] = topId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("topRepository, Update, err:", err.Error())
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

func (repository *topRepositoryImpl) Delete(custId string, topId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_top
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND top = :top_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"top_id":     topId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("TopRepository, Delete, err:", err.Error())
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
