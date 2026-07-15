package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DistributorReplenishmentSetupRepository interface {
	FindAll(dataFilter entity.DistributorReplenishmentSetupQueryFilter) ([]model.DistributorReplenishmentSetup, int, int, error)
	FindAllByPicUserID(dataFilter entity.DistributorReplenishmentSetupQueryFilter, picUserID int) ([]model.DistributorReplenishmentSetup, int, int, error)
	FindDistributorsForPic(f entity.DistributorReplenishmentDistributorQueryFilter) ([]model.DistributorReplenishmentDistributorRow, int, int, error)
	FindByIDAndCustID(id int, custID, parentCustID string) (model.DistributorReplenishmentSetup, error)
	FindApprovalsBySetupIDAndCustID(setupID int, custID, parentCustID string) ([]model.DistributorReplenishmentApproval, error)
	FindSuppliersForPic(f entity.DistributorReplenishmentSupplierQueryFilter) ([]model.DistributorReplenishmentSupplierRow, int, int, error)
	Create(custID string, userID int64, payload entity.DistributorReplenishmentSetupCreatePayload) (int, error)
	Update(id int, custID string, userID int64, payload entity.DistributorReplenishmentSetupCreatePayload) error
	Delete(id int, custID string, userID int64) error
}

func NewDistributorReplenishmentSetupRepository(db *sqlx.DB) DistributorReplenishmentSetupRepository {
	return &distributorReplenishmentSetupRepositoryImpl{db}
}

type distributorReplenishmentSetupRepositoryImpl struct {
	*sqlx.DB
}

func (r *distributorReplenishmentSetupRepositoryImpl) FindAll(dataFilter entity.DistributorReplenishmentSetupQueryFilter) ([]model.DistributorReplenishmentSetup, int, int, error) {
	fromClause := r.buildFromClause()
	conditions, args := r.buildConditionsWithBase(dataFilter, []interface{}{dataFilter.CustId, dataFilter.ParentCustID}, 3)
	queryBody := fromClause + conditions

	countQuery := "SELECT COUNT(*) " + queryBody
	selectQuery := r.buildSelectColumns() + queryBody
	selectQuery += r.buildOrderBy(dataFilter.Sort)

	if dataFilter.Limit > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	total, err := r.executeCountQuery(countQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := r.executeSelectQuery(selectQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)
	return rows, total, lastPage, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) FindAllByPicUserID(dataFilter entity.DistributorReplenishmentSetupQueryFilter, picUserID int) ([]model.DistributorReplenishmentSetup, int, int, error) {
	fromClause := r.buildFromClause() + `
	AND EXISTS (
		SELECT 1 FROM mst.distributor_replenishment_approval appr
		WHERE appr.dist_replenishment_setup_id = drs.id
			AND appr.is_del = false
			AND appr.pic = $3
	)`
	conditions, args := r.buildConditionsWithBase(dataFilter, []interface{}{dataFilter.CustId, dataFilter.ParentCustID, picUserID}, 4)
	queryBody := fromClause + conditions

	countQuery := "SELECT COUNT(*) " + queryBody
	selectQuery := r.buildSelectColumns() + queryBody
	selectQuery += r.buildOrderBy(dataFilter.Sort)

	if dataFilter.Limit > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	total, err := r.executeCountQuery(countQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	rows, err := r.executeSelectQuery(selectQuery, args)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := sql_helper.CalculateLastPage(total, dataFilter.Limit)
	return rows, total, lastPage, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) buildSelectColumns() string {
	return `SELECT
		drs.id,
		drs.sup_id,
		COALESCE(sup.sup_code, '') AS sup_code,
		COALESCE(sup.sup_name, '') AS sup_name,
		drs.distributor_id,
		COALESCE(d.distributor_code, '') AS distributor_code,
		COALESCE(d.distributor_name, '') AS distributor_name,
		drs.distributor_type,
		COALESCE(drs.wh_limit_action, '') AS wh_limit_action,
		drs.wh_capacity,
		drs.wh_volume,
		drs.credit_limit_action,
		drs.plafon_credit,
		drs.lead_time_days,
		drs.is_approval_required,
		drs.created_by,
		uc.user_name AS created_by_name,
		drs.created_at,
		drs.updated_by,
		uu.user_name AS updated_by_name,
		drs.updated_at
	`
}

func (r *distributorReplenishmentSetupRepositoryImpl) buildFromClause() string {
	return ` FROM mst.distributor_replenishment_setup drs
	INNER JOIN mst.m_distributor d ON d.distributor_id = drs.distributor_id AND d.is_del = false
	LEFT JOIN LATERAL (
		SELECT
			s.sup_code,
			s.sup_name
		FROM mst.m_supplier s
		WHERE s.sup_id = drs.sup_id
			AND s.is_del = false
			AND s.cust_id IN ($1, $2)
		ORDER BY
			CASE WHEN COALESCE(s.sup_code, '') <> '' OR COALESCE(s.sup_name, '') <> '' THEN 1 ELSE 0 END DESC,
			CASE WHEN s.cust_id = drs.cust_id THEN 1 ELSE 0 END DESC
		LIMIT 1
	) sup ON true
	LEFT JOIN sys.m_user uc ON uc.user_id = drs.created_by
	LEFT JOIN sys.m_user uu ON uu.user_id = drs.updated_by
	WHERE drs.is_del = false
	AND drs.cust_id = $1`
}

func (r *distributorReplenishmentSetupRepositoryImpl) buildConditionsWithBase(dataFilter entity.DistributorReplenishmentSetupQueryFilter, baseArgs []interface{}, startIdx int) (string, []interface{}) {
	var b strings.Builder
	args := append([]interface{}{}, baseArgs...)
	idx := startIdx

	if len(dataFilter.DistributorIDs) > 0 {
		placeholders := r.placeholders(len(dataFilter.DistributorIDs), idx)
		b.WriteString(fmt.Sprintf(" AND drs.distributor_id IN (%s)", placeholders))
		for _, id := range dataFilter.DistributorIDs {
			args = append(args, id)
		}
		idx += len(dataFilter.DistributorIDs)
	}

	if len(dataFilter.SupplierIDs) > 0 {
		placeholders := r.placeholders(len(dataFilter.SupplierIDs), idx)
		b.WriteString(fmt.Sprintf(" AND drs.sup_id IN (%s)", placeholders))
		for _, id := range dataFilter.SupplierIDs {
			args = append(args, id)
		}
		idx += len(dataFilter.SupplierIDs)
	}

	if dataFilter.Q != "" {
		term := "%" + dataFilter.Q + "%"
		b.WriteString(fmt.Sprintf(` AND (
			sup.sup_code ILIKE $%d OR sup.sup_name ILIKE $%d
			OR d.distributor_code ILIKE $%d OR d.distributor_name ILIKE $%d
		)`, idx, idx+1, idx+2, idx+3))
		args = append(args, term, term, term, term)
	}

	return b.String(), args
}

func (r *distributorReplenishmentSetupRepositoryImpl) placeholders(n, start int) string {
	p := make([]string, n)
	for i := 0; i < n; i++ {
		p[i] = fmt.Sprintf("$%d", start+i)
	}
	return strings.Join(p, ", ")
}

func (r *distributorReplenishmentSetupRepositoryImpl) buildOrderBy(sort string) string {
	columnMap := map[string]string{
		"id":                "drs.id",
		"created_at":        "drs.created_at",
		"created_date":      "drs.created_at",
		"updated_at":        "drs.updated_at",
		"sup_code":          "sup.sup_code",
		"sup_name":          "sup.sup_name",
		"distributor_code":  "d.distributor_code",
		"distributor_name":  "d.distributor_name",
		"wh_capacity":       "drs.wh_capacity",
		"lead_time_days":    "drs.lead_time_days",
	}

	if sort == "" {
		return " ORDER BY drs.created_at DESC"
	}

	parts := strings.Split(sort, ",")
	var orders []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		segs := strings.SplitN(p, ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "DESC"
		if len(segs) > 1 {
			d := strings.ToUpper(strings.TrimSpace(segs[1]))
			if d == "ASC" || d == "DESC" {
				dir = d
			}
		}
		col, ok := columnMap[key]
		if !ok {
			continue
		}
		orders = append(orders, col+" "+dir)
	}
	if len(orders) == 0 {
		return " ORDER BY drs.created_at DESC"
	}
	return " ORDER BY " + strings.Join(orders, ", ")
}

func (r *distributorReplenishmentSetupRepositoryImpl) executeCountQuery(query string, args []interface{}) (int, error) {
	var total int
	if err := r.Get(&total, query, args...); err != nil {
		log.Error(err)
		return 0, err
	}
	return total, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) executeSelectQuery(query string, args []interface{}) ([]model.DistributorReplenishmentSetup, error) {
	var rows []model.DistributorReplenishmentSetup
	if err := r.Select(&rows, query, args...); err != nil {
		log.Error(err)
		return nil, err
	}
	return rows, nil
}

// FindByIDAndCustID returns one setup row scoped by tenant
func (r *distributorReplenishmentSetupRepositoryImpl) FindByIDAndCustID(id int, custID, parentCustID string) (model.DistributorReplenishmentSetup, error) {
	var row model.DistributorReplenishmentSetup
	q := r.buildSelectColumns() + r.buildFromClause() + " AND drs.id = $3"
	err := r.Get(&row, q, custID, parentCustID, id)
	if err != nil {
		log.Error(err)
	}
	return row, err
}

func (r *distributorReplenishmentSetupRepositoryImpl) FindApprovalsBySetupIDAndCustID(setupID int, custID, parentCustID string) ([]model.DistributorReplenishmentApproval, error) {
	q := `SELECT
		a.id,
		a.dist_replenishment_setup_id,
		a.level,
		a.sequence,
		a.business_unit,
		COALESCE(bu.business_unit_name, '') AS business_unit_name,
		a.pic,
		'' AS pic_name,
		a.is_active
	FROM mst.distributor_replenishment_approval a
	INNER JOIN mst.distributor_replenishment_setup drs ON drs.id = a.dist_replenishment_setup_id AND drs.is_del = false
	LEFT JOIN LATERAL (
		SELECT
			dbu.distributor_name AS business_unit_name
		FROM mst.m_distributor dbu
		WHERE dbu.distributor_id = a.business_unit
			AND dbu.is_del = false
		ORDER BY
			CASE WHEN dbu.distributor_name IS NOT NULL AND dbu.distributor_name <> '' THEN 1 ELSE 0 END DESC
		LIMIT 1
	) bu ON true
	WHERE drs.id = $1 AND drs.cust_id = $2
		AND COALESCE(a.is_del, false) = false
	ORDER BY a.level ASC, a.sequence ASC`
	var rows []model.DistributorReplenishmentApproval
	err := r.Select(&rows, q, setupID, custID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err := r.enrichApprovalPicNames(rows, custID, parentCustID); err != nil {
		log.Error(err)
		return nil, err
	}
	return rows, nil
}

type approvalPicEmployeeRow struct {
	EmpID   int    `db:"emp_id"`
	CustID  string `db:"cust_id"`
	EmpName string `db:"emp_name"`
	EmpCode string `db:"emp_code"`
}

func (r *distributorReplenishmentSetupRepositoryImpl) enrichApprovalPicNames(rows []model.DistributorReplenishmentApproval, custID, parentCustID string) error {
	seen := make(map[int]struct{})
	var empIDs []int
	for _, row := range rows {
		if row.Pic <= 0 {
			continue
		}
		if _, ok := seen[row.Pic]; ok {
			continue
		}
		seen[row.Pic] = struct{}{}
		empIDs = append(empIDs, row.Pic)
	}
	if len(empIDs) == 0 {
		return nil
	}
	q := `SELECT emp_id, cust_id, emp_name, emp_code
		FROM mst.m_employee
		WHERE emp_id = ANY($1)
			AND is_del = false`
	var emps []approvalPicEmployeeRow
	if err := r.Select(&emps, q, pq.Array(empIDs)); err != nil {
		return err
	}
	byEmp := make(map[int][]approvalPicEmployeeRow)
	for _, e := range emps {
		byEmp[e.EmpID] = append(byEmp[e.EmpID], e)
	}
	for i := range rows {
		list := byEmp[rows[i].Pic]
		if len(list) == 0 {
			continue
		}
		best := list[0]
		bestScore := custMatchScore(list[0].CustID, custID, parentCustID)
		for j := 1; j < len(list); j++ {
			s := custMatchScore(list[j].CustID, custID, parentCustID)
			if s > bestScore {
				bestScore = s
				best = list[j]
			}
		}
		name := strings.TrimSpace(best.EmpName)
		if name == "" {
			name = strings.TrimSpace(best.EmpCode)
		}
		if name != "" {
			rows[i].PicName = sql.NullString{String: name, Valid: true}
		}
	}
	return nil
}

func custMatchScore(empCustID, custID, parentCustID string) int {
	if empCustID == custID {
		return 2
	}
	if parentCustID != "" && empCustID == parentCustID {
		return 1
	}
	return 0
}

const picSupplierInnerFrom = `
FROM mst.distributor_replenishment_setup drs
INNER JOIN mst.distributor_replenishment_approval dra
	ON dra.dist_replenishment_setup_id = drs.id
	AND dra.is_del = false
	AND dra.cust_id = drs.cust_id
	AND dra.pic = $3
LEFT JOIN LATERAL (
	SELECT s.sup_code, s.sup_name
	FROM mst.m_supplier s
	WHERE s.sup_id = drs.sup_id
		AND s.is_del = false
		AND (s.is_active = true OR s.is_active IS NULL)
		AND s.cust_id IN ($1, $2)
	ORDER BY
		CASE WHEN COALESCE(s.sup_code, '') <> '' OR COALESCE(s.sup_name, '') <> '' THEN 1 ELSE 0 END DESC,
		CASE WHEN s.cust_id = drs.cust_id THEN 1 ELSE 0 END DESC
	LIMIT 1
) ms ON true
WHERE drs.is_del = false
	AND drs.cust_id = $1
	AND ($4::int IS NULL OR drs.distributor_id = $4)
	AND ($5 = '' OR COALESCE(ms.sup_code, '') ILIKE '%' || $5 || '%' OR COALESCE(ms.sup_name, '') ILIKE '%' || $5 || '%')
GROUP BY drs.sup_id`

const picDistributorInnerFrom = `
FROM mst.distributor_replenishment_setup drs
INNER JOIN mst.distributor_replenishment_approval dra
	ON dra.dist_replenishment_setup_id = drs.id
	AND dra.is_del = false
	AND dra.cust_id = drs.cust_id
	AND dra.pic = $2
INNER JOIN mst.m_distributor md
	ON md.distributor_id = drs.distributor_id
	AND md.is_del = false
WHERE drs.is_del = false
	AND drs.cust_id = $1
	AND ($3 = '' OR COALESCE(md.distributor_code, '') ILIKE '%' || $3 || '%' OR COALESCE(md.distributor_name, '') ILIKE '%' || $3 || '%')
GROUP BY drs.distributor_id, md.distributor_code, md.distributor_name`

func buildPicSupplierOrderBy(sort string) string {
	columnMap := map[string]string{
		"created_date": "t.agg_created_at",
		"created_at":   "t.agg_created_at",
		"sup_code":     "t.sup_code",
		"sup_name":     "t.sup_name",
	}
	if sort == "" {
		return " ORDER BY t.agg_created_at DESC"
	}
	parts := strings.Split(sort, ",")
	var orders []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		segs := strings.SplitN(p, ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "DESC"
		if len(segs) > 1 {
			d := strings.ToUpper(strings.TrimSpace(segs[1]))
			if d == "ASC" || d == "DESC" {
				dir = d
			}
		}
		col, ok := columnMap[key]
		if !ok {
			continue
		}
		orders = append(orders, col+" "+dir)
	}
	if len(orders) == 0 {
		return " ORDER BY t.agg_created_at DESC"
	}
	return " ORDER BY " + strings.Join(orders, ", ")
}

func buildPicDistributorOrderBy(sort string) string {
	columnMap := map[string]string{
		"id":               "t.id",
		"created_date":     "t.agg_created_at",
		"created_at":       "t.agg_created_at",
		"distributor_code": "t.distributor_code",
		"distributor_name": "t.distributor_name",
	}
	if sort == "" {
		return " ORDER BY t.agg_created_at DESC"
	}
	parts := strings.Split(sort, ",")
	var orders []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		segs := strings.SplitN(p, ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "DESC"
		if len(segs) > 1 {
			d := strings.ToUpper(strings.TrimSpace(segs[1]))
			if d == "ASC" || d == "DESC" {
				dir = d
			}
		}
		col, ok := columnMap[key]
		if !ok {
			continue
		}
		orders = append(orders, col+" "+dir)
	}
	if len(orders) == 0 {
		return " ORDER BY t.agg_created_at DESC"
	}
	return " ORDER BY " + strings.Join(orders, ", ")
}

func (r *distributorReplenishmentSetupRepositoryImpl) FindSuppliersForPic(f entity.DistributorReplenishmentSupplierQueryFilter) ([]model.DistributorReplenishmentSupplierRow, int, int, error) {
	qPattern := strings.TrimSpace(f.Q)
	var distArg interface{}
	if f.DistributorID != nil {
		distArg = *f.DistributorID
	}
	args := []interface{}{f.CustId, f.ParentCustID, f.Pic, distArg, qPattern}

	countQuery := `SELECT COUNT(*) FROM (SELECT drs.sup_id ` + picSupplierInnerFrom + `) cnt`
	total, err := r.executeCountQuery(countQuery, args)
	if err != nil {
		log.Error(err)
		return nil, 0, 0, err
	}

	page := f.Page
	limit := f.Limit
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}
	offset := (page - 1) * limit

	selectQuery := `SELECT t.sup_id, t.sup_code, t.sup_name FROM (
		SELECT drs.sup_id AS sup_id,
			MAX(COALESCE(ms.sup_code, '')) AS sup_code,
			MAX(COALESCE(ms.sup_name, '')) AS sup_name,
			MAX(drs.created_at) AS agg_created_at` + picSupplierInnerFrom + `
	) t ` + buildPicSupplierOrderBy(f.Sort) + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var rows []model.DistributorReplenishmentSupplierRow
	if err := r.Select(&rows, selectQuery, args...); err != nil {
		log.Error(err)
		return nil, 0, 0, err
	}
	lastPage := sql_helper.CalculateLastPage(total, limit)
	return rows, total, lastPage, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) FindDistributorsForPic(f entity.DistributorReplenishmentDistributorQueryFilter) ([]model.DistributorReplenishmentDistributorRow, int, int, error) {
	qPattern := strings.TrimSpace(f.Q)
	args := []interface{}{f.CustId, f.Pic, qPattern}

	countQuery := `SELECT COUNT(*) FROM (SELECT drs.distributor_id ` + picDistributorInnerFrom + `) cnt`
	total, err := r.executeCountQuery(countQuery, args)
	if err != nil {
		log.Error(err)
		return nil, 0, 0, err
	}

	page := f.Page
	limit := f.Limit
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 5
	}
	offset := (page - 1) * limit

	selectQuery := `SELECT t.id, t.distributor_id, t.distributor_code, t.distributor_name FROM (
		SELECT
			MIN(drs.id) AS id,
			drs.distributor_id,
			MAX(COALESCE(md.distributor_code, '')) AS distributor_code,
			MAX(COALESCE(md.distributor_name, '')) AS distributor_name,
			MAX(drs.created_at) AS agg_created_at ` + picDistributorInnerFrom + `
	) t ` + buildPicDistributorOrderBy(f.Sort) + fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var rows []model.DistributorReplenishmentDistributorRow
	if err := r.Select(&rows, selectQuery, args...); err != nil {
		log.Error(err)
		return nil, 0, 0, err
	}
	lastPage := sql_helper.CalculateLastPage(total, limit)
	return rows, total, lastPage, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) Create(custID string, userID int64, payload entity.DistributorReplenishmentSetupCreatePayload) (int, error) {
	tx, err := r.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var existingID int
	duplicateCheckQuery := `SELECT id
		FROM mst.distributor_replenishment_setup
		WHERE cust_id = $1
			AND distributor_id = $2
			AND sup_id = $3
			AND is_del = false
		LIMIT 1`
	err = tx.Get(&existingID, duplicateCheckQuery, custID, payload.DistributorID, payload.SupID)
	if err == nil {
		return 0, errors.New("combination of distributor_id and sup_id already exists")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	now := time.Now().UTC()
	var setupID int
	createSetupQuery := `INSERT INTO mst.distributor_replenishment_setup (
		cust_id, sup_id, distributor_id, distributor_type, wh_limit_action, wh_capacity, wh_volume,
		credit_limit_action, plafon_credit, lead_time_days, is_approval_required, is_active,
		created_by, created_at, updated_by, updated_at, is_del
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true, $12, $13, $14, $15, false
	) RETURNING id`

	err = tx.QueryRow(
		createSetupQuery,
		custID,
		payload.SupID,
		payload.DistributorID,
		payload.DistributorType,
		whLimitActionForStorage(payload.WhLimitAction),
		nullableIntPtr(payload.WhCapacity),
		nullableIntPtr(payload.WhVolume),
		int(payload.CreditLimitAction),
		nullableIntPtr(payload.PlafonCredit),
		payload.LeadTimeDays,
		ptrBoolOr(payload.IsApprovalRequired, false),
		userID,
		now,
		userID,
		now,
	).Scan(&setupID)
	if err != nil {
		return 0, err
	}

	if err := r.replaceApprovals(tx, setupID, custID, userID, payload.ApprovalData); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return setupID, nil
}

func (r *distributorReplenishmentSetupRepositoryImpl) Update(id int, custID string, userID int64, payload entity.DistributorReplenishmentSetupCreatePayload) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	updateSetupQuery := `UPDATE mst.distributor_replenishment_setup SET
		sup_id = $1,
		distributor_id = $2,
		distributor_type = $3,
		wh_limit_action = $4,
		wh_capacity = $5,
		wh_volume = $6,
		credit_limit_action = $7,
		plafon_credit = $8,
		lead_time_days = $9,
		is_approval_required = $10,
		updated_by = $11,
		updated_at = $12
	WHERE id = $13 AND cust_id = $14 AND is_del = false`

	res, err := tx.Exec(
		updateSetupQuery,
		payload.SupID,
		payload.DistributorID,
		payload.DistributorType,
		whLimitActionForStorage(payload.WhLimitAction),
		nullableIntPtr(payload.WhCapacity),
		nullableIntPtr(payload.WhVolume),
		int(payload.CreditLimitAction),
		nullableIntPtr(payload.PlafonCredit),
		payload.LeadTimeDays,
		ptrBoolOr(payload.IsApprovalRequired, false),
		userID,
		now,
		id,
		custID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("sql: no rows in result set")
	}

	if err := r.replaceApprovals(tx, id, custID, userID, payload.ApprovalData); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *distributorReplenishmentSetupRepositoryImpl) Delete(id int, custID string, userID int64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	_, err = tx.Exec(
		`UPDATE mst.distributor_replenishment_approval
		SET is_del = true, deleted_by = $1, deleted_at = $2
		WHERE dist_replenishment_setup_id = $3 AND cust_id = $4 AND is_del = false`,
		userID, now, id, custID,
	)
	if err != nil {
		return err
	}

	res, err := tx.Exec(
		`UPDATE mst.distributor_replenishment_setup
		SET is_del = true, deleted_by = $1, deleted_at = $2, updated_by = $3, updated_at = $4
		WHERE id = $5 AND cust_id = $6 AND is_del = false`,
		userID, now, userID, now, id, custID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("sql: no rows in result set")
	}

	return tx.Commit()
}

func (r *distributorReplenishmentSetupRepositoryImpl) replaceApprovals(tx *sqlx.Tx, setupID int, custID string, userID int64, approvals []entity.DistributorReplenishmentSetupApprovalPayload) error {
	now := time.Now().UTC()
	_, err := tx.Exec(
		`UPDATE mst.distributor_replenishment_approval
		SET is_del = true, deleted_by = $1, deleted_at = $2
		WHERE dist_replenishment_setup_id = $3 AND cust_id = $4 AND is_del = false`,
		userID, now, setupID, custID,
	)
	if err != nil {
		return err
	}

	for _, item := range approvals {
		_, err = tx.Exec(
			`INSERT INTO mst.distributor_replenishment_approval (
				cust_id, dist_replenishment_setup_id, level, sequence, business_unit, pic,
				is_active, created_by, created_at, updated_by, updated_at, is_del
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, false
			)`,
			custID,
			setupID,
			item.Level,
			item.Sequence,
			item.BusinessUnit,
			item.Pic,
			ptrBoolOr(item.IsActive, true),
			userID,
			now,
			userID,
			now,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func ptrBoolOr(v *bool, defaultVal bool) bool {
	if v == nil {
		return defaultVal
	}
	return *v
}

func nullableStrPtr(p *string) interface{} {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return nil
	}
	return s
}

func whLimitActionForStorage(p *string) string {
	if p == nil {
		return "Unrestricted"
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return "Unrestricted"
	}
	return s
}

func nullableIntPtr(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
