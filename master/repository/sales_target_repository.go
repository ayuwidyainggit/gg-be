package repository

import (
	"database/sql"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SalesTargetRepository interface {
	FindAll(filter entity.SalesTargetQueryFilter, custId string) (data []model.SalesTargetList, total int, lastPage int, err error)
	FindOneById(salesTargetId int64, custId string) (model.SalesTargetList, error)
	FindDetailsBySalesTargetId(salesTargetId int64, custId string) ([]model.SalesAllocatedDetail, error)
	FindMonthlyDistributor(query entity.SalesTargetMonthlyDistQuery) (model.SalesTargetMonthlyDist, error)
	Store(salesTarget model.SalesTarget) (int64, error)
	StoreAllocated(salesAllocated model.SalesAllocated) error
	Update(salesTargetId int64, salesTarget model.SalesTarget) error
	UpdatePartial(salesTargetId int64, custId string, updates map[string]interface{}) error
	DeleteAllocatedByTargetId(salesTargetId int64, custId string) error
	FindMonthlyAllocationByYearlyId(yearlyId int, custId string) ([]model.SalesTargetMonthlyAllocation, error)
	SyncTargetsToMonthly(tx *sqlx.Tx, custId string, yearlyId int, month int, monthlyId int, monthlyTarget int, updatedBy int64) error
	TrxBegin()
	TrxCommit() error
	TrxRollback() error
}

type salesTargetRepositoryImpl struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewSalesTargetRepository(db *sqlx.DB) SalesTargetRepository {
	return &salesTargetRepositoryImpl{db: db}
}

func (repo *salesTargetRepositoryImpl) TrxBegin() {
	repo.tx = repo.db.MustBegin()
}

func (repo *salesTargetRepositoryImpl) TrxCommit() error {
	if repo.tx == nil {
		return nil
	}
	err := repo.tx.Commit()
	repo.tx = nil
	return err
}

func (repo *salesTargetRepositoryImpl) TrxRollback() error {
	if repo.tx == nil {
		return nil
	}
	err := repo.tx.Rollback()
	repo.tx = nil
	return err
}

func (repo *salesTargetRepositoryImpl) FindAll(filter entity.SalesTargetQueryFilter, custId string) (data []model.SalesTargetList, total int, lastPage int, err error) {
	data = []model.SalesTargetList{}

	selectCount := `COUNT(*) AS total`
	selectField := `st.sales_target_id, st.month, st.year, st.allocated_total, st.monthly_target,
					st.remaining, st.status, st.created_at as created_at_raw,
					CASE
						WHEN st.updated_at IS NULL THEN u1.user_name
						ELSE u2.user_name
					END AS updated_by_name,
					CASE
						WHEN st.updated_at IS NULL THEN u1.user_name
						ELSE u2.user_name
					END AS created_by_name,
					CASE
						WHEN st.updated_at IS NULL THEN st.created_at
						ELSE st.updated_at
					END AS updated_at,
					st.created_at`

	qWhere := `WHERE st.cust_id = $1 AND st.is_del = false`

	if filter.Year > 0 {
		qWhere += fmt.Sprintf(` AND st.year = %d`, filter.Year)
	}

	qFrom := ` FROM mst.m_sales_target st
			   LEFT JOIN sys.m_user u1 ON u1.user_id = st.created_by
			   LEFT JOIN sys.m_user u2 ON u2.user_id = st.updated_by`

	queryCount := `SELECT ` + selectCount + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + qFrom + ` ` + qWhere

	err = repo.db.QueryRow(queryCount, custId).Scan(&total)
	if err != nil {
		log.Println("salesTargetRepository, FindAll, count error:", err.Error())
		return data, 0, 0, err
	}

	// Sorting
	sortBy := `st.created_at DESC`
	sortColumnMap := map[string]string{
		"created_date":    "st.created_at",
		"created_at":      "st.created_at",
		"month":           "st.month",
		"year":            "st.year",
		"allocated_total": "st.allocated_total",
		"monthly_target":  "st.monthly_target",
		"remaining":       "st.remaining",
	}

	if filter.Sort != "" {
		rawSort := strings.ToLower(strings.TrimSpace(filter.Sort))
		if rawSort == "asc" || rawSort == "desc" {
			sortBy = fmt.Sprintf("st.created_at %s", strings.ToUpper(rawSort))
		} else {
			mSortBy := strings.Split(filter.Sort, ",")
			sortClauses := []string{}
			for _, row := range mSortBy {
				colSort := strings.Split(strings.TrimSpace(row), ":")
				if len(colSort) != 2 {
					continue
				}

				columnKey := strings.ToLower(strings.TrimSpace(colSort[0]))
				direction := strings.ToLower(strings.TrimSpace(colSort[1]))
				if direction != "asc" && direction != "desc" {
					continue
				}

				columnName, exists := sortColumnMap[columnKey]
				if !exists {
					continue
				}

				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", columnName, strings.ToUpper(direction)))
			}

			if len(sortClauses) > 0 {
				sortBy = strings.Join(sortClauses, ", ")
			}
		}
	}
	querySelect += ` ORDER BY ` + sortBy

	// Pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * filter.Limit
	lastPage = int(math.Ceil(float64(total) / float64(filter.Limit)))

	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, filter.Limit, offset)

	err = repo.db.Select(&data, querySelect, custId)
	if err != nil {
		log.Println("salesTargetRepository, FindAll, select error:", err.Error())
		return data, total, lastPage, err
	}

	return data, total, lastPage, nil
}

func (repo *salesTargetRepositoryImpl) FindOneById(salesTargetId int64, custId string) (model.SalesTargetList, error) {
	var data model.SalesTargetList

	query := `SELECT st.sales_target_id, st.month, st.year, st.allocated_total, st.monthly_target,
			  st.remaining, st.status, st.created_at as created_at_raw,
			  CASE
				  WHEN st.updated_at IS NULL THEN u1.user_name
				  ELSE u2.user_name
			  END AS updated_by_name,
			  CASE
				  WHEN st.updated_at IS NULL THEN u1.user_name
				  ELSE u2.user_name
			  END AS created_by_name,
			  CASE
				  WHEN st.updated_at IS NULL THEN st.created_at
				  ELSE st.updated_at
			  END AS updated_at,
			  st.created_at
			  FROM mst.m_sales_target st
			  LEFT JOIN sys.m_user u1 ON u1.user_id = st.created_by
			  LEFT JOIN sys.m_user u2 ON u2.user_id = st.updated_by
			  WHERE st.sales_target_id = $1 AND st.cust_id = $2 AND st.is_del = false`

	err := repo.db.Get(&data, query, salesTargetId, custId)
	if err != nil {
		log.Println("salesTargetRepository, FindOneById, error:", err.Error())
		return data, err
	}

	return data, nil
}

func (repo *salesTargetRepositoryImpl) FindDetailsBySalesTargetId(salesTargetId int64, custId string) ([]model.SalesAllocatedDetail, error) {
	details := []model.SalesAllocatedDetail{}

	query := `SELECT sa.sales_allocated_id, sa.sales_target_id, sa.salesman_id,
			  COALESCE(s.sales_name, '') as sales_name,
			  COALESCE(CASE
				WHEN s.opr_type = 'C' THEN 'Canvas'
				WHEN s.opr_type = 'O' THEN 'Taking Order'
				WHEN s.opr_type = 'S' THEN 'Shop Sales'
				ELSE s.opr_type
			  END, '') as opr_type,
			  stdy.distributor_id,
			  COALESCE(d.distributor_code, '') as distributor_code,
			  COALESCE(d.distributor_name, '') as distributor_name,
			  d.channel_id,
			  COALESCE(mc.channel_code, '') as channel_code,
			  COALESCE(mc.channel_name, '') as channel_name,
			  sa.sales_team_id,
			  COALESCE(st.sales_team_code, '') as sales_team_code,
			  COALESCE(st.sales_team_name, '') as sales_team_name,
			  sa.allocated, sa.is_active
			  FROM mst.m_sales_allocated sa
			  JOIN mst.m_sales_target mst ON mst.sales_target_id = sa.sales_target_id
			  JOIN mst.m_sales_target_distributor_yearly stdy ON stdy.sales_target_distributor_yearly_id = mst.sales_target_distributor_yearly_id
			  LEFT JOIN mst.m_distributor d ON d.distributor_id = stdy.distributor_id
			  LEFT JOIN mst.m_channel mc ON mc.channel_id = d.channel_id
			  LEFT JOIN mst.m_salesman s ON s.emp_id = sa.salesman_id AND s.cust_id = sa.cust_id
			  LEFT JOIN mst.m_sales_team st ON st.sales_team_id = sa.sales_team_id
			  WHERE sa.sales_target_id = $1 AND sa.cust_id = $2 AND sa.is_del = false
			  ORDER BY sa.sales_allocated_id`

	err := repo.db.Select(&details, query, salesTargetId, custId)
	if err != nil {
		log.Println("salesTargetRepository, FindDetailsBySalesTargetId, error:", err.Error())
		return details, err
	}

	return details, nil
}

func (repo *salesTargetRepositoryImpl) FindMonthlyDistributor(query entity.SalesTargetMonthlyDistQuery) (model.SalesTargetMonthlyDist, error) {
	var data model.SalesTargetMonthlyDist

	querySQL := `SELECT stdy.distributor_id, stdy.year, stdm.month, stdm.monthly_target
				 FROM mst.m_sales_target_distributor_monthly stdm
				 JOIN mst.m_sales_target_distributor_yearly stdy 
					ON stdm.sales_target_distributor_yearly_id = stdy.sales_target_distributor_yearly_id
				 WHERE stdy.distributor_id = $1 AND stdy.year = $2 AND stdm.month = $3
				 LIMIT 1`

	err := repo.db.Get(&data, querySQL, query.DistributorId, query.Year, query.Month)
	if err != nil {
		log.Println("salesTargetRepository, FindMonthlyDistributor, error:", err.Error())
		return data, err
	}

	return data, nil
}

func (repo *salesTargetRepositoryImpl) Store(salesTarget model.SalesTarget) (int64, error) {
	query := `INSERT INTO mst.m_sales_target (
				cust_id, sales_target_distributor_yearly_id, sales_target_distributor_monthly_id,
				month, year, allocated_total, monthly_target, remaining, status,
				created_by, created_at, updated_by, updated_at, is_del
			  ) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
			  ) RETURNING sales_target_id`

	var salesTargetId int64
	err := repo.tx.QueryRow(query,
		salesTarget.CustId, salesTarget.SalesTargetDistributorYearlyId, salesTarget.SalesTargetDistributorMonthlyId,
		salesTarget.Month, salesTarget.Year, salesTarget.AllocatedTotal, salesTarget.MonthlyTarget,
		salesTarget.Remaining, salesTarget.Status, salesTarget.CreatedBy, salesTarget.CreatedAt,
		salesTarget.UpdatedBy, salesTarget.UpdatedAt, salesTarget.IsDel,
	).Scan(&salesTargetId)

	if err != nil {
		log.Println("salesTargetRepository, Store, error:", err.Error())
		return 0, err
	}

	return salesTargetId, nil
}

func (repo *salesTargetRepositoryImpl) StoreAllocated(salesAllocated model.SalesAllocated) error {
	query := `INSERT INTO mst.m_sales_allocated (
				cust_id, sales_target_id, salesman_id, sales_team_id, allocated, is_active,
				created_by, created_at, updated_by, updated_at, is_del
			  ) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			  )`

	_, err := repo.tx.Exec(query,
		salesAllocated.CustId, salesAllocated.SalesTargetId, salesAllocated.SalesmanId,
		salesAllocated.SalesTeamId, salesAllocated.Allocated, salesAllocated.IsActive,
		salesAllocated.CreatedBy, salesAllocated.CreatedAt, salesAllocated.UpdatedBy,
		salesAllocated.UpdatedAt, salesAllocated.IsDel,
	)

	if err != nil {
		log.Println("salesTargetRepository, StoreAllocated, error:", err.Error())
		return err
	}

	return nil
}

func (repo *salesTargetRepositoryImpl) Update(salesTargetId int64, salesTarget model.SalesTarget) error {
	query := `UPDATE mst.m_sales_target 
			  SET sales_target_distributor_yearly_id = $1, 
				  sales_target_distributor_monthly_id = $2,
				  month = $3, year = $4, allocated_total = $5, 
				  monthly_target = $6, remaining = $7, status = $8,
				  updated_by = $9, updated_at = $10
			  WHERE sales_target_id = $11 AND cust_id = $12 AND is_del = false`

	result, err := repo.tx.Exec(query,
		salesTarget.SalesTargetDistributorYearlyId, salesTarget.SalesTargetDistributorMonthlyId,
		salesTarget.Month, salesTarget.Year, salesTarget.AllocatedTotal,
		salesTarget.MonthlyTarget, salesTarget.Remaining, salesTarget.Status,
		salesTarget.UpdatedBy, salesTarget.UpdatedAt,
		salesTargetId, salesTarget.CustId,
	)

	if err != nil {
		log.Println("salesTargetRepository, Update, error:", err.Error())
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (repo *salesTargetRepositoryImpl) UpdatePartial(salesTargetId int64, custId string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil // nothing to update
	}

	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf(`UPDATE mst.m_sales_target SET %s WHERE sales_target_id = $%d AND cust_id = $%d AND is_del = false`,
		strings.Join(setClauses, ", "), argIndex, argIndex+1)
	args = append(args, salesTargetId, custId)

	execTx := false
	if repo.tx != nil {
		execTx = true
	}

	var (
		result sql.Result
		err    error
	)
	if execTx {
		result, err = repo.tx.Exec(query, args...)
	} else {
		result, err = repo.db.Exec(query, args...)
	}
	if err != nil {
		log.Println("salesTargetRepository, UpdatePartial, error:", err.Error())
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (repo *salesTargetRepositoryImpl) DeleteAllocatedByTargetId(salesTargetId int64, custId string) error {
	query := `DELETE FROM mst.m_sales_allocated 
			  WHERE sales_target_id = $1 AND cust_id = $2`

	_, err := repo.tx.Exec(query, salesTargetId, custId)
	if err != nil {
		log.Println("salesTargetRepository, DeleteAllocatedByTargetId, error:", err.Error())
		return err
	}

	return nil
}

func (repo *salesTargetRepositoryImpl) FindMonthlyAllocationByYearlyId(yearlyId int, custId string) ([]model.SalesTargetMonthlyAllocation, error) {
	data := []model.SalesTargetMonthlyAllocation{}

	query := `SELECT
				month,
				COALESCE(SUM(allocated_total), 0) AS allocated_total,
				COUNT(*) AS target_count
			  FROM mst.m_sales_target
			  WHERE sales_target_distributor_yearly_id = $1
				AND cust_id = $2
				AND is_del = false
			  GROUP BY month
			  ORDER BY month ASC`

	err := repo.db.Select(&data, query, yearlyId, custId)
	if err != nil {
		log.Println("salesTargetRepository, FindMonthlyAllocationByYearlyId, error:", err.Error())
		return data, err
	}

	return data, nil
}

func (repo *salesTargetRepositoryImpl) SyncTargetsToMonthly(tx *sqlx.Tx, custId string, yearlyId int, month int, monthlyId int, monthlyTarget int, updatedBy int64) error {
	query := `UPDATE mst.m_sales_target
			  SET sales_target_distributor_monthly_id = $1,
				  monthly_target = $2,
				  remaining = $2 - allocated_total,
				  updated_by = $3,
				  updated_at = CURRENT_TIMESTAMP
			  WHERE sales_target_distributor_yearly_id = $4
				AND month = $5
				AND cust_id = $6
				AND is_del = false`

	_, err := tx.Exec(query, monthlyId, monthlyTarget, updatedBy, yearlyId, month, custId)
	if err != nil {
		log.Println("salesTargetRepository, SyncTargetsToMonthly, error:", err.Error())
		return err
	}

	return nil
}
