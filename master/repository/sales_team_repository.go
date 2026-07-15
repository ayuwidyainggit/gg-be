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

type SalesTeamRepository interface {
	FindOneParentCustId(distCustId string) (model.MCustomer, error)
	FindOneCustomerByDistributorID(distributorID int64) (model.MCustomer, error)
	FindOneBySalesTeamIdAndCustId(salesTeamId int, custId string) (model.SalesTeam, error)
	FindOneBySalesTeamCodeAndCustId(salesTeamCode string, custId string) (model.SalesTeam, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.SalesTeam, total int, lastPage int, err error)
	FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.SalesTeam, total int, lastPage int, err error)
	Store(salesTeam model.SalesTeam) (int, error)
	Update(salesTeamId int, request entity.UpdateSalesTeamRequest) error
	Delete(custId string, salesTeamId int, deletedBy int64) error
}

func NewSalesTeamRepository(db *sqlx.DB) SalesTeamRepository {
	return &salesTeamRepositoryImpl{db}
}

type salesTeamRepositoryImpl struct {
	*sqlx.DB
}

func buildSalesTeamCustScopeCondition(distributorIDs []int, parentCustId, custId string) string {
	if len(distributorIDs) == 0 {
		return ` a.cust_id = '` + custId + `' `
	}

	scopeParentCustID := strings.TrimSpace(parentCustId)
	if scopeParentCustID == "" {
		scopeParentCustID = strings.TrimSpace(custId)
	}

	includePrincipalScope := false
	stringIDs := make([]string, 0, len(distributorIDs))
	seen := make(map[int]struct{})

	for _, distributorID := range distributorIDs {
		if _, exists := seen[distributorID]; exists {
			continue
		}
		seen[distributorID] = struct{}{}

		if distributorID == 0 {
			includePrincipalScope = true
			continue
		}

		if distributorID < 0 {
			continue
		}

		stringIDs = append(stringIDs, strconv.Itoa(distributorID))
	}

	conditions := make([]string, 0, 2)
	if includePrincipalScope {
		conditions = append(conditions, `a.cust_id = '`+scopeParentCustID+`'`)
	}
	if len(stringIDs) > 0 {
		conditions = append(conditions, `a.cust_id IN (
			SELECT mc.cust_id
			FROM smc.m_customer mc
			WHERE mc.parent_cust_id = '`+scopeParentCustID+`'
			AND mc.distributor_id IN (`+strings.Join(stringIDs, ",")+`)
		)`)
	}

	if len(conditions) == 0 {
		return ` a.cust_id = '` + custId + `' `
	}

	return `( ` + strings.Join(conditions, ` OR `) + ` )`
}

func extractSalesTeamDistributorIDs(dataFilter entity.GeneralQueryFilter) []int {
	if len(dataFilter.DistributorIDs) > 0 {
		return dataFilter.DistributorIDs
	}

	if dataFilter.DistributorID > 0 {
		return []int{int(dataFilter.DistributorID)}
	}

	return nil
}

func (repository *salesTeamRepositoryImpl) FindOneParentCustId(distCustId string) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE cust_id = $1`
	err := repository.Get(&mCustomer, query, distCustId)
	if err != nil {
		log.Println("salesTeamRepository, FindOneParentCustId, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func (repository *salesTeamRepositoryImpl) FindOneBySalesTeamIdAndCustId(salesTeamId int, custId string) (model.SalesTeam, error) {
	salesTeam := model.SalesTeam{}
	query := `SELECT 
				cust_id, sales_team_id, sales_team_code,
				sales_team_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sales_team 
			  WHERE sales_team_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&salesTeam, query, salesTeamId, custId)
	if err != nil {
		log.Println("salesTeamRepository, FindOneBySalesTeamCodeAndCustId, err:", err.Error())
		return salesTeam, err
	}

	return salesTeam, nil
}

func (repository *salesTeamRepositoryImpl) FindOneBySalesTeamCodeAndCustId(salesTeamCode string, custId string) (model.SalesTeam, error) {
	salesTeam := model.SalesTeam{}
	query := `SELECT 
				cust_id, sales_team_id, sales_team_code,
				sales_team_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_sales_team 
			  WHERE sales_team_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&salesTeam, query, salesTeamCode, custId)
	if err != nil {
		log.Println("salesTeamRepository, FindOneBySalesTeamCodeAndCustId, err:", err.Error())
		return salesTeam, err
	}

	return salesTeam, nil
}

func (repository *salesTeamRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SalesTeam, int, int, error) {

	salesTeams := []model.SalesTeam{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sales_team_id, a.sales_team_code, a.sales_team_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at, a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND ` + buildSalesTeamCustScopeCondition(extractSalesTeamDistributorIDs(dataFilter), dataFilter.ParentCustId, custId)

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sales_team_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sales_team_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := `FROM mst.m_sales_team a 
			  LEFT JOIN sys.m_user u ON u.user_id = a.updated_by `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("salesTeamRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesTeamRepository, count total, err:", err.Error())
		return salesTeams, 0, 0, err
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
		sortBy := `sales_team_id`
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

	// log.Println("salesTeamRepository, querySelect:", querySelect)
	err = repository.Select(&salesTeams, querySelect)
	if err != nil {
		log.Println("salesTeamRepository, FindAllByCustId, err:", err.Error())
		return salesTeams, total, lastPage, err
	}

	return salesTeams, total, lastPage, nil
}

func (repository *salesTeamRepositoryImpl) FindAllByCustIdLookup(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SalesTeam, int, int, error) {

	salesTeams := []model.SalesTeam{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sales_team_id, a.sales_team_code, a.sales_team_name  `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND ` + buildSalesTeamCustScopeCondition(extractSalesTeamDistributorIDs(dataFilter), dataFilter.ParentCustId, custId)

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sales_team_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sales_team_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := `FROM mst.m_sales_team a `
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("salesTeamRepository, count total, err:", err.Error())
		return salesTeams, 0, 0, err
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
		sortBy := `a.sales_team_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&salesTeams, querySelect)
	if err != nil {
		log.Println("salesTeamRepository, FindAllByCustId, err:", err.Error())
		return salesTeams, total, 1, err
	}

	return salesTeams, total, 1, nil
}

func (repository *salesTeamRepositoryImpl) Store(salesTeam model.SalesTeam) (int, error) {
	query :=
		`INSERT INTO mst.m_sales_team(
			cust_id, sales_team_code, sales_team_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING sales_team_id;`
	lastInsertId := salesTeam.SalesTeamId
	err := repository.QueryRow(query,
		salesTeam.CustId, salesTeam.SalesTeamCode, salesTeam.SalesTeamName,
		salesTeam.IsActive, salesTeam.CreatedBy, salesTeam.CreatedAt, salesTeam.UpdatedBy,
		salesTeam.UpdatedAt, salesTeam.IsDel, salesTeam.DeletedBy, salesTeam.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("salesTeamRepository, Store, err:", err.Error())
		return salesTeam.SalesTeamId, err
	}
	return salesTeam.SalesTeamId, nil
}

func (repository *salesTeamRepositoryImpl) Update(salesTeamId int, request entity.UpdateSalesTeamRequest) error {
	var (
		r            model.SalesTeamUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("salesTeamRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sales_team
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sales_team_id = :sales_team_id_old;`

	// log.Println("salesTeamRepository, Update, query:", query)

	sqlPatch.Args["sales_team_id_old"] = salesTeamId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("salesTeamRepository, Update, err:", err.Error())
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

func (repository *salesTeamRepositoryImpl) Delete(custId string, salesTeamId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sales_team
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND sales_team_id = :sales_team_id;`

	wMap := map[string]interface{}{
		"cust_id":       custId,
		"sales_team_id": salesTeamId,
		"deleted_by":    deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SalesTeamRepository, Delete, err:", err.Error())
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

func (repository *salesTeamRepositoryImpl) FindOneCustomerByDistributorID(distributorID int64) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE distributor_id = $1`
	err := repository.Get(&mCustomer, query, distributorID)
	if err != nil {
		log.Println("salesTeamRepository, FindOneCustomerByDistributorID, err:", err.Error())
		return mCustomer, err
	}
	return mCustomer, nil
}
