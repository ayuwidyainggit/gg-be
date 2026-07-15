package repository

import (
	"fmt"
	"master/entity"
	"master/model"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BusinessUnitRepository interface {
	FindUserByUsername(username string) (model.UserInfo, error)
	FindCustomerNameByCustId(custId string) (string, error)
	FindDistributorsByCustId(dataFilter entity.BusinessUnitQueryFilter) ([]model.BusinessUnitDistributor, int, int, error)
	FindDistributorByDistributorId(distributorId int, custId string) (model.BusinessUnitDistributor, error)
}

type businessUnitRepositoryImpl struct {
	*sqlx.DB
}

func NewBusinessUnitRepository(db *sqlx.DB) *businessUnitRepositoryImpl {
	return &businessUnitRepositoryImpl{db}
}

func (repo *businessUnitRepositoryImpl) FindUserByUsername(username string) (model.UserInfo, error) {
	user := model.UserInfo{}
	query := `SELECT user_id, user_fullname FROM sys.m_user WHERE user_name = $1 AND is_del = false`
	err := repo.Get(&user, query, username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *businessUnitRepositoryImpl) FindCustomerNameByCustId(custId string) (string, error) {
	var custName string
	query := `SELECT cust_name FROM smc.m_customer WHERE cust_id = $1 LIMIT 1`
	err := repo.Get(&custName, query, custId)
	if err != nil {
		return "", err
	}
	return custName, nil
}

func (repo *businessUnitRepositoryImpl) FindDistributorsByCustId(dataFilter entity.BusinessUnitQueryFilter) ([]model.BusinessUnitDistributor, int, int, error) {
	distributors := []model.BusinessUnitDistributor{}

	queryCount, countArgs, querySelect, selectArgs := buildFindDistributorsByCustIDQuery(dataFilter)

	queryCount, countArgs, err := sqlx.In(queryCount, countArgs...)
	if err != nil {
		return distributors, 0, 0, err
	}
	queryCount = repo.Rebind(queryCount)

	var total int
	err = repo.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		return distributors, 0, 0, err
	}

	limit := dataFilter.Limit
	if limit == 0 {
		limit = 9999
	}
	if limit > 9999 {
		limit = 9999
	}
	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	querySelect, selectArgs, err = sqlx.In(querySelect, selectArgs...)
	if err != nil {
		return distributors, total, lastPage, err
	}
	querySelect = repo.Rebind(querySelect)

	err = repo.Select(&distributors, querySelect, selectArgs...)
	if err != nil {
		return distributors, total, lastPage, err
	}

	return distributors, total, lastPage, nil
}

func buildFindDistributorsByCustIDQuery(dataFilter entity.BusinessUnitQueryFilter) (string, []interface{}, string, []interface{}) {
	selectCount := `COUNT(DISTINCT md.distributor_id) AS total`
	selectField := `DISTINCT md.cust_id, md.distributor_id, md.distributor_code, md.distributor_name,
                     md.area_id, ma.area_code, ma.area_name,
                     md.region_id, mr.region_code, mr.region_name`

	qFrom := ` FROM mst.m_distributor md
               LEFT JOIN mst.m_area ma ON ma.area_id = md.area_id
               LEFT JOIN mst.m_region mr ON mr.region_id = md.region_id`

	whereClauses := []string{"md.is_del = false"}
	queryArgs := make([]interface{}, 0)

	if dataFilter.Scope.DistributorScope == "specific" {
		qFrom += ` INNER JOIN mst.m_employee_distributor_mapping edm
			ON edm.distributor_id = md.distributor_id
			AND edm.cust_id = ?
			AND edm.emp_id = ?
			AND edm.is_del = false`
		queryArgs = append(queryArgs, dataFilter.CustId, dataFilter.EmployeeId)
	} else {
		if dataFilter.Scope.AreaScope == "specific" {
			qFrom += ` INNER JOIN mst.m_employee_area_mapping eam
				ON eam.area_id = md.area_id
				AND eam.cust_id = ?
				AND eam.emp_id = ?
				AND eam.is_del = false`
			queryArgs = append(queryArgs, dataFilter.CustId, dataFilter.EmployeeId)
		}
		if dataFilter.Scope.RegionScope == "specific" {
			qFrom += ` INNER JOIN mst.m_employee_region_mapping erm
				ON erm.region_id = md.region_id
				AND erm.cust_id = ?
				AND erm.emp_id = ?
				AND erm.is_del = false`
			queryArgs = append(queryArgs, dataFilter.CustId, dataFilter.EmployeeId)
		}
		if dataFilter.Scope.AreaScope != "specific" && dataFilter.Scope.RegionScope != "specific" {
			scopeParent := strings.TrimSpace(dataFilter.ParentCustId)
			if scopeParent == "" {
				scopeParent = dataFilter.CustId
			}
			whereClauses = append(whereClauses, "md.parent_cust_id = ?")
			queryArgs = append(queryArgs, scopeParent)
		}
	}

	if len(dataFilter.RegionId) > 0 {
		whereClauses = append(whereClauses, "md.region_id IN (?)")
		queryArgs = append(queryArgs, dataFilter.RegionId)
	}

	if len(dataFilter.AreaId) > 0 {
		whereClauses = append(whereClauses, "md.area_id IN (?)")
		queryArgs = append(queryArgs, dataFilter.AreaId)
	}

	if len(dataFilter.IsActive) > 0 {
		for _, status := range dataFilter.IsActive {
			if status == 1 {
				whereClauses = append(whereClauses, "md.is_active = ?")
				queryArgs = append(queryArgs, true)
				break
			}
			if status == 0 {
				whereClauses = append(whereClauses, "md.is_active = ?")
				queryArgs = append(queryArgs, false)
				break
			}
		}
	}

	if dataFilter.Query != "" {
		whereClauses = append(whereClauses, "(md.distributor_code ILIKE ? OR md.distributor_name ILIKE ?)")
		queryLike := "%" + dataFilter.Query + "%"
		queryArgs = append(queryArgs, queryLike, queryLike)
	}

	qWhere := ` WHERE ` + strings.Join(whereClauses, " AND ")
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		sortBy := ""
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		querySelect += ` ORDER BY md.distributor_id ASC`
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 9999
	}
	if dataFilter.Limit > 9999 {
		dataFilter.Limit = 9999
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit
	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, dataFilter.Limit, offset)

	return queryCount, queryArgs, querySelect, queryArgs
}

func (repo *businessUnitRepositoryImpl) FindDistributorByDistributorId(distributorId int, custId string) (model.BusinessUnitDistributor, error) {
	distributor := model.BusinessUnitDistributor{}
	query := `SELECT md.cust_id, md.distributor_id, md.distributor_code, md.distributor_name,
	                    md.area_id, ma.area_code, ma.area_name,
	                    md.region_id, mr.region_code, mr.region_name
	             FROM mst.m_distributor md
	             LEFT JOIN mst.m_area ma ON ma.area_id = md.area_id
	             LEFT JOIN mst.m_region mr ON mr.region_id = md.region_id
	             WHERE md.distributor_id = $1 AND md.is_del = false`
	err := repo.Get(&distributor, query, distributorId)
	if err != nil {
		return distributor, err
	}
	return distributor, nil
}
