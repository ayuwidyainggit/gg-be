package repository

import (
	"encoding/json"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type SalesTargetDistributorRepository interface {
	FindAllByCustId(dataFilter entity.SalesTargetDistributorQueryFilter, custId string) (data []model.SalesTargetDistributorYearly, total int, lastPage int, err error)
	FindOneByIdAndCustId(id int, custId string) (model.SalesTargetDistributorYearly, error)
	FindMonthlyDetailsByYearlyId(yearlyId int) ([]model.SalesTargetDistributorMonthly, error)
	FindChildCustIDByDistributorID(distributorID int) (string, error)
	FindAllocationSummaryByYearlyId(yearlyId int, custId string) (model.SalesTargetAllocationSummary, error)
	StoreYearly(tx *sqlx.Tx, data model.SalesTargetDistributorYearly) (int, error)
	StoreMonthly(tx *sqlx.Tx, data model.SalesTargetDistributorMonthly) (int, error)
	UpdateMonthlyTarget(tx *sqlx.Tx, monthlyId int, monthlyTarget int, updatedBy int64) error
	UpdateYearly(tx *sqlx.Tx, id int, request entity.UpdateSalesTargetDistributorRequest) error
	DeleteMonthlyByYearlyId(tx *sqlx.Tx, yearlyId int, deletedBy int64) error
	BeginTx() (*sqlx.Tx, error)
}

type salesTargetDistributorRepositoryImpl struct {
	*sqlx.DB
}

// NewSalesTargetDistributorRepository creates a new instance of SalesTargetDistributorRepository
func NewSalesTargetDistributorRepository(db *sqlx.DB) SalesTargetDistributorRepository {
	return &salesTargetDistributorRepositoryImpl{db}
}

// BeginTx starts a new database transaction
func (repository *salesTargetDistributorRepositoryImpl) BeginTx() (*sqlx.Tx, error) {
	return repository.Beginx()
}

func buildSalesTargetDistributorStatusClause(statuses []int, currentYear int) string {
	if len(statuses) == 0 {
		return ""
	}

	conditions := make([]string, 0, len(statuses))
	for _, s := range statuses {
		switch entity.SalesTargetStatus(s) {
		case entity.SALES_TARGET_STATUS_DRAFT:
			conditions = append(conditions, "std.status = 0")
		case entity.SALES_TARGET_STATUS_ACTIVE:
			conditions = append(conditions, fmt.Sprintf("(std.status != 0 AND std.year <= %d AND std.is_active = true)", currentYear))
		case entity.SALES_TARGET_STATUS_INACTIVE:
			conditions = append(conditions, fmt.Sprintf("(std.status != 0 AND (std.year > %d OR std.is_active = false))", currentYear))
		}
	}

	if len(conditions) == 0 {
		return ""
	}

	return strings.Join(conditions, " OR ")
}

// FindAllByCustId retrieves a paginated list of yearly sales targets for a specific customer ID
func (repository *salesTargetDistributorRepositoryImpl) FindAllByCustId(dataFilter entity.SalesTargetDistributorQueryFilter, custId string) (data []model.SalesTargetDistributorYearly, total int, lastPage int, err error) {
	data = []model.SalesTargetDistributorYearly{}
	selectCount := `COUNT(*) AS total `
	selectField := `
		std.cust_id, std.sales_target_distributor_yearly_id, std.area_id, std.region_id, std.distributor_id,
		std.year, std.yearly_target, std.status, std.is_active, std.user_inactive, std.inactive_at,
		std.created_by, std.created_at, std.updated_by, std.updated_at, std.is_del,
		COALESCE(d.distributor_code, '') AS distributor_code,
		COALESCE(d.distributor_name, '') AS distributor_name,
		COALESCE(a.area_code, '') AS area_code,
		COALESCE(a.area_name, '') AS area_name,
		COALESCE(r.region_code, '') AS region_code,
		COALESCE(r.region_name, '') AS region_name,
		CASE 
			WHEN std.updated_by IS NULL THEN COALESCE(u_created.user_fullname, '')
			ELSE COALESCE(u.user_fullname, '')
		END AS updated_by_name
	`

	qFrom := `
		FROM mst.m_sales_target_distributor_yearly std
		LEFT JOIN mst.m_distributor d ON d.distributor_id = std.distributor_id
		LEFT JOIN mst.m_area a ON a.area_id = std.area_id
		LEFT JOIN mst.m_region r ON r.region_id = std.region_id
		LEFT JOIN sys.m_user u ON u.user_id = std.updated_by
		LEFT JOIN sys.m_user u_created ON u_created.user_id = std.created_by
	`

	qWhere := `
		WHERE std.is_del = false
		AND std.cust_id = $1 `

	if dataFilter.Year != nil {
		qWhere += fmt.Sprintf(" AND std.year = %d ", *dataFilter.Year)
	}

	if dataFilter.Status != nil && len(*dataFilter.Status) > 0 {
		currentYear := time.Now().Year()
		statusClause := buildSalesTargetDistributorStatusClause(*dataFilter.Status, currentYear)
		if statusClause != "" {
			qWhere += fmt.Sprintf(" AND (%s) ", statusClause)
		}
	}

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	err = repository.QueryRow(queryCount, custId).Scan(&total)
	if err != nil {
		return data, 0, 0, err
	}

	sortBy := ``
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`std.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		querySelect += ` ORDER BY std.created_at DESC`
	}

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 20
	}

	if dataFilter.Limit > 9999 {
		dataFilter.Limit = 9999
	}

	page := dataFilter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit

	lastPage = int(math.Ceil(float64(float64(total) / float64(dataFilter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(dataFilter.Limit), strconv.Itoa(offset))

	err = repository.Select(&data, querySelect, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	return data, total, lastPage, nil
}

// FindOneByIdAndCustId retrieves a single yearly sales target by its ID and customer ID
func (repository *salesTargetDistributorRepositoryImpl) FindOneByIdAndCustId(id int, custId string) (model.SalesTargetDistributorYearly, error) {
	data := model.SalesTargetDistributorYearly{}
	query := `
		SELECT
			std.cust_id, std.sales_target_distributor_yearly_id, std.area_id, std.region_id, std.distributor_id,
			std.year, std.yearly_target, std.status, std.is_active, std.user_inactive, std.inactive_at,
			std.created_by, std.created_at, std.updated_by, std.updated_at, std.is_del,
			COALESCE(d.distributor_code, '') AS distributor_code,
			COALESCE(d.distributor_name, '') AS distributor_name,
			COALESCE(a.area_code, '') AS area_code,
			COALESCE(a.area_name, '') AS area_name,
			COALESCE(r.region_code, '') AS region_code,
			COALESCE(r.region_name, '') AS region_name,
			CASE 
				WHEN std.updated_by IS NULL THEN COALESCE(u_created.user_fullname, '')
				ELSE COALESCE(u.user_fullname, '')
			END AS updated_by_name
		FROM mst.m_sales_target_distributor_yearly std
		LEFT JOIN mst.m_distributor d ON d.distributor_id = std.distributor_id
		LEFT JOIN mst.m_area a ON a.area_id = std.area_id
		LEFT JOIN mst.m_region r ON r.region_id = std.region_id
		LEFT JOIN sys.m_user u ON u.user_id = std.updated_by
		LEFT JOIN sys.m_user u_created ON u_created.user_id = std.created_by
		WHERE std.is_del = false
		AND std.sales_target_distributor_yearly_id = $1
		AND std.cust_id = $2
	`
	err := repository.Get(&data, query, id, custId)
	if err != nil {
		return data, err
	}
	return data, nil
}

// FindMonthlyDetailsByYearlyId retrieves all monthly targets associated with a specific yearly target ID
func (repository *salesTargetDistributorRepositoryImpl) FindMonthlyDetailsByYearlyId(yearlyId int) ([]model.SalesTargetDistributorMonthly, error) {
	data := []model.SalesTargetDistributorMonthly{}
	query := `
		SELECT 
			cust_id, sales_target_distributor_monthly_id, sales_target_distributor_yearly_id,
			month, monthly_target, is_active
		FROM mst.m_sales_target_distributor_monthly
		WHERE is_del = false
		AND sales_target_distributor_yearly_id = $1
		ORDER BY month ASC
	`
	err := repository.Select(&data, query, yearlyId)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (repository *salesTargetDistributorRepositoryImpl) FindChildCustIDByDistributorID(distributorID int) (string, error) {
	var custID string
	query := `
		SELECT mc.cust_id
		FROM smc.m_customer mc
		WHERE mc.distributor_id = $1
		LIMIT 1
	`

	err := repository.Get(&custID, query, distributorID)
	if err != nil {
		return "", err
	}

	return custID, nil
}

// FindAllocationSummaryByYearlyId retrieves allocation summary for a yearly distributor target
func (repository *salesTargetDistributorRepositoryImpl) FindAllocationSummaryByYearlyId(yearlyId int, custId string) (model.SalesTargetAllocationSummary, error) {
	data := model.SalesTargetAllocationSummary{}
	query := `
		SELECT
			COALESCE(SUM(t.allocated_total), 0) AS allocated_total,
			CASE
				WHEN COUNT(a.salesman_id) > 0 OR COALESCE(SUM(t.allocated_total), 0) > 0 THEN true
				ELSE false
			END AS is_allocated
		FROM mst.m_sales_target_distributor_yearly y
		JOIN mst.m_sales_target_distributor_monthly m
			ON y.sales_target_distributor_yearly_id = m.sales_target_distributor_yearly_id
			AND m.is_del = false
		LEFT JOIN mst.m_sales_target t
			ON m.sales_target_distributor_monthly_id = t.sales_target_distributor_monthly_id
			AND t.is_del = false
		LEFT JOIN mst.m_sales_allocated a
			ON t.sales_target_id = a.sales_target_id
			AND a.is_del = false
		WHERE y.is_del = false
		AND y.sales_target_distributor_yearly_id = $1
		AND y.cust_id = $2
	`
	err := repository.Get(&data, query, yearlyId, custId)
	if err != nil {
		return data, err
	}
	return data, nil
}

// StoreYearly inserts a new yearly sales target record into the database within a transaction
func (repository *salesTargetDistributorRepositoryImpl) StoreYearly(tx *sqlx.Tx, data model.SalesTargetDistributorYearly) (int, error) {
	query := `
		INSERT INTO mst.m_sales_target_distributor_yearly (
			cust_id, area_id, region_id, distributor_id, year, yearly_target, status, is_active,
			created_by, created_at, is_del
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING sales_target_distributor_yearly_id
	`
	var id int
	err := tx.QueryRow(query,
		data.CustId, data.AreaId, data.RegionId, data.DistributorId, data.Year, data.YearlyTarget,
		data.Status, data.IsActive, data.CreatedBy, data.CreatedAt, data.IsDel,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// StoreMonthly inserts a new monthly sales target record into the database within a transaction
func (repository *salesTargetDistributorRepositoryImpl) StoreMonthly(tx *sqlx.Tx, data model.SalesTargetDistributorMonthly) (int, error) {
	query := `
		INSERT INTO mst.m_sales_target_distributor_monthly (
			cust_id, sales_target_distributor_yearly_id, month, monthly_target, is_active,
			created_by, created_at, is_del
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING sales_target_distributor_monthly_id
	`
	var monthlyID int
	err := tx.QueryRow(query,
		data.CustId, data.SalesTargetDistributorYearlyId, data.Month, data.MonthlyTarget,
		data.IsActive, data.CreatedBy, data.CreatedAt, data.IsDel,
	).Scan(&monthlyID)
	if err != nil {
		return 0, err
	}
	return monthlyID, nil
}

func (repository *salesTargetDistributorRepositoryImpl) UpdateMonthlyTarget(tx *sqlx.Tx, monthlyId int, monthlyTarget int, updatedBy int64) error {
	query := `UPDATE mst.m_sales_target_distributor_monthly
			  SET monthly_target = $1,
				  updated_by = $2,
				  updated_at = CURRENT_TIMESTAMP
			  WHERE sales_target_distributor_monthly_id = $3
				AND is_del = false`

	_, err := tx.Exec(query, monthlyTarget, updatedBy, monthlyId)
	if err != nil {
		return err
	}

	return nil
}

// UpdateYearly updates an existing yearly sales target record using patched fields
func (repository *salesTargetDistributorRepositoryImpl) UpdateYearly(tx *sqlx.Tx, id int, request entity.UpdateSalesTargetDistributorRequest) error {
	var (
		r            model.SalesTargetDistributorYearlyUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	// Manually set internal fields since they are tagged with json:"-" in entity
	r.UpdatedBy = &request.UpdatedBy
	if request.UserInactive != nil {
		r.UserInactive = request.UserInactive
	}
	if request.InactiveAt != nil {
		r.InactiveAt = request.InactiveAt
	}

	sqlPatch := sql_helper.SQLPatches(r)

	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sales_target_distributor_yearly
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false
			  AND cust_id = :cust_id
			  AND sales_target_distributor_yearly_id = :sales_target_distributor_yearly_id;`

	sqlPatch.Args["sales_target_distributor_yearly_id"] = id
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		return constant.ErrNoRowsAffected
	}
	if nRows == 0 {
		return constant.ErrNoRowsAffected
	}

	return nil
}

// DeleteMonthlyByYearlyId marks monthly targets as deleted for a specific yearly target ID within a transaction
func (repository *salesTargetDistributorRepositoryImpl) DeleteMonthlyByYearlyId(tx *sqlx.Tx, yearlyId int, deletedBy int64) error {
	query := `
		UPDATE mst.m_sales_target_distributor_monthly
		SET is_del = true,
			deleted_by = $1,
			deleted_at = CURRENT_TIMESTAMP
		WHERE sales_target_distributor_yearly_id = $2
	`
	_, err := tx.Exec(query, deletedBy, yearlyId)
	if err != nil {
		return err
	}
	return nil
}
