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

type MarketRepository interface {
	FindOneByMarketIdAndCustId(marketId int, custId string) (model.Market, error)
	FindOneByMarketCodeAndCustId(marketCode string, custId string) (model.Market, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Market, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Market, total int, lastPage int, err error)
	Store(market model.Market) (int, error)
	Update(marketId int, request entity.UpdateMarketRequest) error
	Delete(custId string, marketId int, deletedBy int64) error
}

func NewMarketRepository(db *sqlx.DB) MarketRepository {
	return &marketRepositoryImpl{db}
}

type marketRepositoryImpl struct {
	*sqlx.DB
}

func (repository *marketRepositoryImpl) FindOneByMarketIdAndCustId(marketId int, custId string) (model.Market, error) {
	market := model.Market{}
	query := `SELECT 
				cust_id, market_id, market_code,
				market_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_market 
			  WHERE market_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&market, query, marketId, custId)
	if err != nil {
		log.Println("marketRepository, FindOneByMarketCodeAndCustId, err:", err.Error())
		return market, err
	}

	return market, nil
}

func (repository *marketRepositoryImpl) FindOneByMarketCodeAndCustId(marketCode string, custId string) (model.Market, error) {
	market := model.Market{}
	query := `SELECT 
				cust_id, market_id, market_code,
				market_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_market 
			  WHERE market_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&market, query, marketCode, custId)
	if err != nil {
		log.Println("marketRepository, FindOneByMarketCodeAndCustId, err:", err.Error())
		return market, err
	}

	return market, nil
}

func (repository *marketRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Market, int, int, error) {

	markets := []model.Market{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.market_id, a.market_code, a.market_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.market_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.market_name ILIKE '%` + dataFilter.Query + `%' )`
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

	qFrom := ` 	FROM mst.m_market a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + `  ` + qWhere

	// log.Println("marketRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("marketRepository, count total, err:", err.Error())
		return markets, 0, 0, err
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
		sortBy := `a.market_id`
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

	// log.Println("marketRepository, querySelect:", querySelect)
	err = repository.Select(&markets, querySelect)
	if err != nil {
		log.Println("marketRepository, FindAllByCustId, err:", err.Error())
		return markets, total, lastPage, err
	}

	return markets, total, lastPage, nil
}

func (repository *marketRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Market, int, int, error) {

	markets := []model.Market{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.market_id, a.market_code, a.market_name `
	qWhere := ` WHERE a.is_del = false AND 
					a.is_active = true AND 
					a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.market_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.market_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_market a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + `  ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("FindAllByCustIdLookupMode, count total, err:", err.Error())
		return markets, 0, 0, err
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
		sortBy := `a.market_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 10
	}

	err = repository.Select(&markets, querySelect)
	if err != nil {
		log.Println("FindAllByCustIdLookupMode, err:", err.Error())
		return markets, total, 1, err
	}

	return markets, total, 1, nil
}

func (repository *marketRepositoryImpl) Store(market model.Market) (int, error) {
	query :=
		`INSERT INTO mst.m_market(
			cust_id, market_code, market_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING market_id;`
	lastInsertId := market.MarketId
	err := repository.QueryRow(query,
		market.CustId, market.MarketCode, market.MarketName,
		market.IsActive, market.CreatedBy, market.CreatedAt, market.UpdatedBy,
		market.UpdatedAt, market.IsDel, market.DeletedBy, market.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("marketRepository, Store, err:", err.Error())
		return market.MarketId, err
	}
	return market.MarketId, nil
}

func (repository *marketRepositoryImpl) Update(marketId int, request entity.UpdateMarketRequest) error {
	var (
		r            model.MarketUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("marketRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_market
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND market_id = :market_id_old;`

	// log.Println("marketRepository, Update, query:", query)

	sqlPatch.Args["market_id_old"] = marketId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("marketRepository, Update, err:", err.Error())
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

func (repository *marketRepositoryImpl) Delete(custId string, marketId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_market
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND market_id = :market_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"market_id":  marketId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("MarketRepository, Delete, err:", err.Error())
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
