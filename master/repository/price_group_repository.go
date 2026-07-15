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

type PriceGroupRepository interface {
	FindOneByPriceGroupIdAndCustId(priceGroupId int, custId string) (model.PriceGroup, error)
	FindOneByPriceGroupCodeAndCustId(priceGroupCode string, custId string) (model.PriceGroup, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (priceGroup []model.PriceGroup, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (priceGroup []model.PriceGroup, total int, lastPage int, err error)
	Store(priceGroup model.PriceGroup) (int, error)
	Update(priceGroupId int, request entity.UpdatePriceGroupRequest) error
	Delete(custId string, priceGroupId int, deletedBy int64) error
}

func NewPriceGroupRepository(db *sqlx.DB) PriceGroupRepository {
	return &priceGroupRepositoryImpl{db}
}

type priceGroupRepositoryImpl struct {
	*sqlx.DB
}

func (repository *priceGroupRepositoryImpl) FindOneByPriceGroupIdAndCustId(priceGroupId int, custId string) (model.PriceGroup, error) {
	priceGroup := model.PriceGroup{}
	query := `SELECT 
				cust_id, price_grp_id, price_grp_code,
				price_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_price_group 
			  WHERE price_grp_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&priceGroup, query, priceGroupId, custId)
	if err != nil {
		log.Println("priceGroupRepository, FindOneByPriceGroupCodeAndCustId, err:", err.Error())
		return priceGroup, err
	}

	return priceGroup, nil
}

func (repository *priceGroupRepositoryImpl) FindOneByPriceGroupCodeAndCustId(priceGroupCode string, custId string) (model.PriceGroup, error) {
	priceGroup := model.PriceGroup{}
	query := `SELECT 
				cust_id, price_grp_id, price_grp_code,
				price_grp_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_price_group 
			  WHERE price_grp_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&priceGroup, query, priceGroupCode, custId)
	if err != nil {
		log.Println("priceGroupRepository, FindOneByPriceGroupCodeAndCustId, err:", err.Error())
		return priceGroup, err
	}

	return priceGroup, nil
}

func (repository *priceGroupRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PriceGroup, int, int, error) {

	priceGroups := []model.PriceGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.price_grp_id, a.price_grp_code, a.price_grp_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.price_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.price_grp_name ILIKE '%` + dataFilter.Query + `%' )`
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

	qFrom := ` 	FROM mst.m_price_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("priceGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("priceGroupRepository, count total, err:", err.Error())
		return priceGroups, 0, 0, err
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
		sortBy := `a.price_grp_id`
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

	err = repository.Select(&priceGroups, querySelect)
	if err != nil {
		log.Println("priceGroupRepository, FindAllByCustId, err:", err.Error())
		return priceGroups, total, lastPage, err
	}

	return priceGroups, total, lastPage, nil
}

func (repository *priceGroupRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.PriceGroup, int, int, error) {

	priceGroups := []model.PriceGroup{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.price_grp_id, a.price_grp_code, a.price_grp_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.price_grp_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.price_grp_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_price_group a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("priceGroupRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("priceGroupRepository, count total, err:", err.Error())
		return priceGroups, 0, 0, err
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
		sortBy := `a.price_grp_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&priceGroups, querySelect)
	if err != nil {
		log.Println("priceGroupRepository, FindAllByCustId, err:", err.Error())
		return priceGroups, total, 1, err
	}

	return priceGroups, total, 1, nil
}

func (repository *priceGroupRepositoryImpl) Store(priceGroup model.PriceGroup) (int, error) {
	query :=
		`INSERT INTO mst.m_price_group(
			cust_id, price_grp_code, price_grp_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING price_grp_id;`
	lastInsertId := priceGroup.PriceGrpId
	err := repository.QueryRow(query,
		priceGroup.CustId, priceGroup.PriceGrpCode, priceGroup.PriceGrpName,
		priceGroup.IsActive, priceGroup.CreatedBy, priceGroup.CreatedAt, priceGroup.UpdatedBy,
		priceGroup.UpdatedAt, priceGroup.IsDel, priceGroup.DeletedBy, priceGroup.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("priceGroupRepository, Store, err:", err.Error())
		return priceGroup.PriceGrpId, err
	}
	return priceGroup.PriceGrpId, nil
}

func (repository *priceGroupRepositoryImpl) Update(priceGroupId int, request entity.UpdatePriceGroupRequest) error {
	var (
		r            model.PriceGroupUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("priceGroupRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_price_group
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND price_grp_id = :price_grp_id_old;`

	// log.Println("priceGroupRepository, Update, query:", query)

	sqlPatch.Args["price_grp_id_old"] = priceGroupId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("priceGroupRepository, Update, err:", err.Error())
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

func (repository *priceGroupRepositoryImpl) Delete(custId string, priceGroupId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_price_group
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND price_grp_id = :price_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":      custId,
		"price_grp_id": priceGroupId,
		"deleted_by":   deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("PriceGroupRepository, Delete, err:", err.Error())
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
