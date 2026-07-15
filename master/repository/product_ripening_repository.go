package repository

import (
	"database/sql"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type ProductRipeningRepository interface {
	HasAssignedDistributor(custID, parentCustID string, picUserID int64) (bool, error)
	FindAssignedDistributorByCode(custID, parentCustID string, picUserID int64, distributorCode string) (model.ProductRipeningAssignedDistributor, error)
	FindDistributorByCode(parentCustID, distributorCode string) (model.ProductRipeningAssignedDistributor, error)
	FindWeekByYearAndWeekID(parentCustID, distributorCustID string, perYear, weekID int, today time.Time) (model.ProductRipeningWeek, error)
	FindRowsByPlan(custID, parentCustID string, picUserID, distributorID int64, perYear, weekID int) ([]model.ProductRipening, error)
	ReplacePlanRows(custID string, distributorID int64, perYear, perID, weekID int, rows []model.ProductRipening, preserveBefore time.Time, userID int64) error
	ListPlans(filter entity.ProductRipeningQueryFilter, picUserID int64, today time.Time) ([]model.ProductRipeningPlanListRow, int, int, error)
	ExportRows(filter entity.ProductRipeningQueryFilter, picUserID int64) ([]model.ProductRipeningPlanListRow, error)
}

func NewProductRipeningRepository(db *sqlx.DB) ProductRipeningRepository {
	return &productRipeningRepositoryImpl{db}
}

type productRipeningRepositoryImpl struct {
	*sqlx.DB
}

func (r *productRipeningRepositoryImpl) HasAssignedDistributor(custID, parentCustID string, picUserID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM mst.distributor_replenishment_setup drs
			INNER JOIN mst.distributor_replenishment_approval dra
				ON dra.dist_replenishment_setup_id = drs.id
				AND dra.cust_id = drs.cust_id
				AND dra.is_del = false
				AND dra.pic = $3
			INNER JOIN mst.m_distributor md
				ON md.distributor_id = drs.distributor_id
				AND md.is_del = false
			WHERE drs.cust_id = $1
				AND drs.is_del = false
				AND md.parent_cust_id = $2
		)`
	err := r.Get(&exists, query, custID, parentCustID, picUserID)
	return exists, err
}

func (r *productRipeningRepositoryImpl) FindAssignedDistributorByCode(custID, parentCustID string, picUserID int64, distributorCode string) (model.ProductRipeningAssignedDistributor, error) {
	row := model.ProductRipeningAssignedDistributor{}
	query := `
		SELECT DISTINCT
			md.distributor_id,
			md.distributor_code,
			md.distributor_name,
			md.cust_id AS dist_cust_id
		FROM mst.distributor_replenishment_setup drs
		INNER JOIN mst.distributor_replenishment_approval dra
			ON dra.dist_replenishment_setup_id = drs.id
			AND dra.cust_id = drs.cust_id
			AND dra.is_del = false
			AND dra.pic = $3
		INNER JOIN mst.m_distributor md
			ON md.distributor_id = drs.distributor_id
			AND md.is_del = false
		WHERE drs.cust_id = $1
			AND drs.is_del = false
			AND md.parent_cust_id = $2
			AND UPPER(md.distributor_code) = UPPER($4)
		LIMIT 1`
	err := r.Get(&row, query, custID, parentCustID, picUserID, distributorCode)
	return row, err
}

func (r *productRipeningRepositoryImpl) FindDistributorByCode(parentCustID, distributorCode string) (model.ProductRipeningAssignedDistributor, error) {
	row := model.ProductRipeningAssignedDistributor{}
	query := `
		SELECT
			md.distributor_id,
			md.distributor_code,
			md.distributor_name,
			md.cust_id AS dist_cust_id
		FROM mst.m_distributor md
		WHERE md.parent_cust_id = $1
			AND md.is_del = false
			AND UPPER(md.distributor_code) = UPPER($2)
		LIMIT 1`
	err := r.Get(&row, query, parentCustID, distributorCode)
	return row, err
}

func (r *productRipeningRepositoryImpl) FindWeekByYearAndWeekID(parentCustID, distributorCustID string, perYear, weekID int, today time.Time) (model.ProductRipeningWeek, error) {
	row := model.ProductRipeningWeek{}
	query := `
		SELECT cust_id, per_year, per_id, week_id, week_start::text, week_end::text
		FROM mst.m_week
		WHERE (
				(cust_id = $1 AND working_day_calendar_id IS NOT NULL)
				OR (cust_id = $2 AND working_day_calendar_id IS NULL)
			)
			AND per_year = $3
			AND week_id = $4
			AND is_closed = false
			AND week_end::date >= $5::date
		ORDER BY
			CASE
				WHEN cust_id = $2 AND working_day_calendar_id IS NULL THEN 0
				WHEN cust_id = $1 AND working_day_calendar_id IS NOT NULL THEN 1
				ELSE 2
			END,
			per_id DESC
		LIMIT 1`
	err := r.Get(&row, query, parentCustID, distributorCustID, perYear, weekID, today.Format("2006-01-02"))
	return row, err
}

func (r *productRipeningRepositoryImpl) FindRowsByPlan(custID, parentCustID string, picUserID, distributorID int64, perYear, weekID int) ([]model.ProductRipening, error) {
	rows := []model.ProductRipening{}
	query := `
		SELECT
			pr.id,
			pr.cust_id,
			pr.distributor_id,
			pr.pro_id,
			pr.per_year,
			pr.per_id,
			pr.week_id,
			pr.sunday_qty,
			pr.monday_qty,
			pr.tuesday_qty,
			pr.wednesday_qty,
			pr.thursday_qty,
			pr.friday_qty,
			pr.saturday_qty,
			pr.created_by,
			pr.created_at,
			pr.updated_by,
			pr.updated_at,
			pr.deleted_by,
			pr.deleted_at,
			pr.is_del,
			COALESCE(md.distributor_code, '') AS distributor_code,
			COALESCE(md.distributor_name, '') AS distributor_name,
			COALESCE(mp.pro_code, '') AS product_code,
			COALESCE(mp.pro_name, '') AS product_name,
			COALESCE(mw.week_start::text, '') AS week_start,
			COALESCE(mw.week_end::text, '') AS week_end,
			COALESCE(uc.user_name, '') AS created_by_name,
			uu.user_name AS updated_by_name
		FROM mst.product_ripening pr
		INNER JOIN mst.m_distributor md
			ON md.distributor_id = pr.distributor_id
			AND md.parent_cust_id = $2
			AND md.is_del = false
		INNER JOIN mst.m_product mp
			ON mp.pro_id = pr.pro_id
			AND mp.cust_id = md.cust_id
			AND mp.deleted_at IS NULL
		LEFT JOIN LATERAL (
			SELECT week_start, week_end
			FROM mst.m_week mw_lookup
			WHERE mw_lookup.per_year = pr.per_year
				AND mw_lookup.week_id = pr.week_id
				AND mw_lookup.is_closed = false
				AND (
					(mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL)
					OR (mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL)
				)
			ORDER BY CASE
				WHEN mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL THEN 0
				WHEN mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL THEN 1
				ELSE 2
			END
			LIMIT 1
		) mw ON true
		LEFT JOIN sys.m_user uc ON uc.user_id = pr.created_by
		LEFT JOIN sys.m_user uu ON uu.user_id = pr.updated_by
		WHERE pr.cust_id = $1
			AND pr.distributor_id = $3
			AND pr.per_year = $4
			AND pr.week_id = $5
			AND pr.is_del = false
		ORDER BY mp.pro_code ASC`
	err := r.Select(&rows, query, custID, parentCustID, distributorID, perYear, weekID)
	_ = picUserID
	return rows, err
}

func (r *productRipeningRepositoryImpl) ReplacePlanRows(custID string, distributorID int64, perYear, perID, weekID int, rows []model.ProductRipening, preserveBefore time.Time, userID int64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	existing, err := r.findRowsByPlanTx(tx, custID, distributorID, perYear, perID, weekID)
	if err != nil {
		return err
	}

	existingByProID := make(map[int64]model.ProductRipening, len(existing))
	for _, row := range existing {
		existingByProID[row.ProID] = row
	}

	now := time.Now().UTC()
	seen := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		seen[row.ProID] = struct{}{}
		if old, ok := existingByProID[row.ProID]; ok {
			finalRow := row
			finalRow.ID = old.ID
			finalRow.CreatedBy = old.CreatedBy
			finalRow.CreatedAt = old.CreatedAt
			finalRow.UpdatedBy = &userID
			finalRow.UpdatedAt = &now
			if !preserveBefore.IsZero() {
				preservePastRipeningValues(&finalRow, old, preserveBefore)
			}
			if err := r.updateRowTx(tx, finalRow); err != nil {
				return err
			}
			continue
		}

		row.CreatedBy = userID
		row.CreatedAt = now
		row.UpdatedBy = &userID
		row.UpdatedAt = &now
		if err := r.insertRowTx(tx, row); err != nil {
			return err
		}
	}

	for _, old := range existing {
		if _, ok := seen[old.ProID]; ok {
			continue
		}
		if !preserveBefore.IsZero() {
			updated := old
			zeroFutureRipeningValues(&updated, preserveBefore)
			updated.UpdatedBy = &userID
			updated.UpdatedAt = &now
			if err := r.updateRowTx(tx, updated); err != nil {
				return err
			}
			continue
		}

		if _, err := tx.Exec(`
			UPDATE mst.product_ripening
			SET is_del = true, deleted_by = $1, deleted_at = $2, updated_by = $1, updated_at = $2
			WHERE id = $3`, userID, now, old.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *productRipeningRepositoryImpl) ListPlans(filter entity.ProductRipeningQueryFilter, picUserID int64, today time.Time) ([]model.ProductRipeningPlanListRow, int, int, error) {
	args := []interface{}{filter.CustId, filter.ParentCustId}
	where, args := r.buildPlanWhere(filter, args, 3)
	fromClause := `
		FROM mst.product_ripening pr
		INNER JOIN mst.m_distributor md
			ON md.distributor_id = pr.distributor_id
			AND md.parent_cust_id = $2
			AND md.is_del = false
		INNER JOIN LATERAL (
			SELECT week_start, week_end
			FROM mst.m_week mw_lookup
			WHERE mw_lookup.per_year = pr.per_year
				AND mw_lookup.week_id = pr.week_id
				AND mw_lookup.is_closed = false
				AND (
					(mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL)
					OR (mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL)
				)
			ORDER BY CASE
				WHEN mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL THEN 0
				WHEN mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL THEN 1
				ELSE 2
			END
			LIMIT 1
		) mw ON true
		LEFT JOIN sys.m_user uc ON uc.user_id = pr.created_by
		LEFT JOIN sys.m_user uu ON uu.user_id = pr.updated_by
		WHERE pr.cust_id = $1
			AND pr.is_del = false` + where

	countQuery := `
		SELECT COUNT(*)
		FROM (
			SELECT 1 ` + fromClause + `
			GROUP BY pr.cust_id, pr.distributor_id, pr.per_year, pr.week_id, md.distributor_code, md.distributor_name
		) grouped_plans`
	var total int
	if err := r.Get(&total, countQuery, args...); err != nil {
		return nil, 0, 0, err
	}

	query := `
		SELECT
			pr.cust_id,
			pr.distributor_id,
			pr.per_year,
			MIN(pr.per_id) AS per_id,
			pr.week_id,
			COALESCE(md.distributor_code, '') AS distributor_code,
			COALESCE(md.distributor_name, '') AS distributor_name,
			COALESCE(MIN(mw.week_start)::text, '') AS week_start,
			COALESCE(MIN(mw.week_end)::text, '') AS week_end,
			(ARRAY_AGG(pr.id ORDER BY pr.created_at ASC, pr.id ASC))[1] AS id,
			COUNT(DISTINCT pr.pro_id) AS total_product,
			(ARRAY_AGG(pr.created_by ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_by,
			(ARRAY_AGG(COALESCE(uc.user_name, '') ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_by_name,
			(ARRAY_AGG(pr.created_at ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_at,
			(ARRAY_AGG(pr.updated_by ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_by,
			(ARRAY_AGG(uu.user_name ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_by_name,
			(ARRAY_AGG(pr.updated_at ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_at ` + fromClause + `
		GROUP BY pr.cust_id, pr.distributor_id, pr.per_year, pr.week_id, md.distributor_code, md.distributor_name`
	query += buildRipeningPlanOrderBy(filter.Sort)
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)
	}

	rows := []model.ProductRipeningPlanListRow{}
	if err := r.Select(&rows, query, args...); err != nil {
		return nil, 0, 0, err
	}
	lastPage := sql_helper.CalculateLastPage(total, filter.Limit)
	_ = picUserID
	_ = today
	return rows, total, lastPage, nil
}

func (r *productRipeningRepositoryImpl) ExportRows(filter entity.ProductRipeningQueryFilter, picUserID int64) ([]model.ProductRipeningPlanListRow, error) {
	args := []interface{}{filter.CustId, filter.ParentCustId}
	where, args := r.buildPlanWhere(filter, args, 3)
	query := `
		SELECT
			pr.cust_id,
			pr.distributor_id,
			pr.per_year,
			MIN(pr.per_id) AS per_id,
			pr.week_id,
			COALESCE(md.distributor_code, '') AS distributor_code,
			COALESCE(md.distributor_name, '') AS distributor_name,
			COALESCE(MIN(mw.week_start)::text, '') AS week_start,
			COALESCE(MIN(mw.week_end)::text, '') AS week_end,
			(ARRAY_AGG(pr.id ORDER BY pr.created_at ASC, pr.id ASC))[1] AS id,
			COUNT(DISTINCT pr.pro_id) AS total_product,
			(ARRAY_AGG(pr.created_by ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_by,
			(ARRAY_AGG(COALESCE(uc.user_name, '') ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_by_name,
			(ARRAY_AGG(pr.created_at ORDER BY pr.created_at ASC, pr.id ASC))[1] AS created_at,
			(ARRAY_AGG(pr.updated_by ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_by,
			(ARRAY_AGG(uu.user_name ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_by_name,
			(ARRAY_AGG(pr.updated_at ORDER BY pr.created_at ASC, pr.id ASC))[1] AS updated_at
		FROM mst.product_ripening pr
		INNER JOIN mst.m_distributor md
			ON md.distributor_id = pr.distributor_id
			AND md.parent_cust_id = $2
			AND md.is_del = false
		INNER JOIN LATERAL (
			SELECT week_start, week_end
			FROM mst.m_week mw_lookup
			WHERE mw_lookup.per_year = pr.per_year
				AND mw_lookup.week_id = pr.week_id
				AND mw_lookup.is_closed = false
				AND (
					(mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL)
					OR (mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL)
				)
			ORDER BY CASE
				WHEN mw_lookup.cust_id = pr.cust_id AND mw_lookup.working_day_calendar_id IS NULL THEN 0
				WHEN mw_lookup.cust_id = $2 AND mw_lookup.working_day_calendar_id IS NOT NULL THEN 1
				ELSE 2
			END
			LIMIT 1
		) mw ON true
		LEFT JOIN sys.m_user uc ON uc.user_id = pr.created_by
		LEFT JOIN sys.m_user uu ON uu.user_id = pr.updated_by
		WHERE pr.cust_id = $1
			AND pr.is_del = false` + where + `
		GROUP BY pr.cust_id, pr.distributor_id, pr.per_year, pr.week_id, md.distributor_code, md.distributor_name`
	query += buildRipeningPlanOrderBy(filter.Sort)
	rows := []model.ProductRipeningPlanListRow{}
	err := r.Select(&rows, query, args...)
	_ = picUserID
	return rows, err
}

func (r *productRipeningRepositoryImpl) buildPlanWhere(filter entity.ProductRipeningQueryFilter, args []interface{}, startIdx int) (string, []interface{}) {
	var b strings.Builder
	idx := startIdx
	if filter.Query != "" {
		term := "%" + strings.TrimSpace(filter.Query) + "%"
		b.WriteString(fmt.Sprintf(` AND (
			md.distributor_code ILIKE $%d OR md.distributor_name ILIKE $%d
			OR CAST(pr.per_year AS TEXT) ILIKE $%d OR CAST(pr.week_id AS TEXT) ILIKE $%d
		)`, idx, idx+1, idx+2, idx+3))
		args = append(args, term, term, term, term)
		idx += 4
	}
	if len(filter.DistributorID) > 0 {
		holders := make([]string, 0, len(filter.DistributorID))
		for _, id := range filter.DistributorID {
			holders = append(holders, fmt.Sprintf("$%d", idx))
			args = append(args, id)
			idx++
		}
		b.WriteString(" AND pr.distributor_id IN (" + strings.Join(holders, ", ") + ")")
	}
	if filter.JwtDistributorId > 0 {
		b.WriteString(fmt.Sprintf(" AND pr.distributor_id = $%d", idx))
		args = append(args, filter.JwtDistributorId)
		idx++
	}
	if filter.WeekID != nil {
		b.WriteString(fmt.Sprintf(" AND pr.week_id = $%d", idx))
		args = append(args, *filter.WeekID)
		idx++
	}
	if filter.PerYear != nil {
		b.WriteString(fmt.Sprintf(" AND pr.per_year = $%d", idx))
		args = append(args, *filter.PerYear)
		idx++
	}
	if filter.PerID != nil {
		b.WriteString(fmt.Sprintf(" AND pr.per_id = $%d", idx))
		args = append(args, *filter.PerID)
	}
	return b.String(), args
}

func buildRipeningPlanOrderBy(sort string) string {
	columnMap := map[string]string{
		"id":               "MIN(pr.id)",
		"distributor_code": "md.distributor_code",
		"distributor_name": "md.distributor_name",
		"week_id":          "pr.week_id",
		"per_year":         "pr.per_year",
		"per_id":           "MIN(pr.per_id)",
		"total_product":    "COUNT(DISTINCT pr.pro_id)",
	}
	if sort == "" {
		return " ORDER BY pr.per_year DESC, MIN(pr.per_id) DESC, pr.week_id DESC, md.distributor_code ASC"
	}
	parts := strings.Split(sort, ",")
	orders := make([]string, 0, len(parts))
	for _, part := range parts {
		segs := strings.SplitN(strings.TrimSpace(part), ":", 2)
		key := strings.TrimSpace(segs[0])
		dir := "DESC"
		if len(segs) > 1 {
			upper := strings.ToUpper(strings.TrimSpace(segs[1]))
			if upper == "ASC" || upper == "DESC" {
				dir = upper
			}
		}
		col, ok := columnMap[key]
		if ok {
			orders = append(orders, col+" "+dir)
		}
	}
	if len(orders) == 0 {
		return " ORDER BY pr.per_year DESC, MIN(pr.per_id) DESC, pr.week_id DESC, md.distributor_code ASC"
	}
	return " ORDER BY " + strings.Join(orders, ", ")
}

func (r *productRipeningRepositoryImpl) findRowsByPlanTx(tx *sqlx.Tx, custID string, distributorID int64, perYear, perID, weekID int) ([]model.ProductRipening, error) {
	rows := []model.ProductRipening{}
	err := tx.Select(&rows, `
		SELECT id, cust_id, distributor_id, pro_id, per_year, per_id, week_id,
			sunday_qty, monday_qty, tuesday_qty, wednesday_qty, thursday_qty, friday_qty, saturday_qty,
			created_by, created_at, updated_by, updated_at, deleted_by, deleted_at, is_del
		FROM mst.product_ripening
		WHERE cust_id = $1 AND distributor_id = $2 AND per_year = $3 AND per_id = $4 AND week_id = $5 AND is_del = false`,
		custID, distributorID, perYear, perID, weekID)
	return rows, err
}

func (r *productRipeningRepositoryImpl) insertRowTx(tx *sqlx.Tx, row model.ProductRipening) error {
	_, err := tx.Exec(`
		INSERT INTO mst.product_ripening (
			cust_id, distributor_id, pro_id, per_year, per_id, week_id,
			sunday_qty, monday_qty, tuesday_qty, wednesday_qty, thursday_qty, friday_qty, saturday_qty,
			created_by, created_at, updated_by, updated_at, is_del
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16, $17, false
		)`,
		row.CustID, row.DistributorID, row.ProID, row.PerYear, row.PerID, row.WeekID,
		row.SundayQty, row.MondayQty, row.TuesdayQty, row.WednesdayQty, row.ThursdayQty, row.FridayQty, row.SaturdayQty,
		row.CreatedBy, row.CreatedAt, row.UpdatedBy, row.UpdatedAt)
	return err
}

func (r *productRipeningRepositoryImpl) updateRowTx(tx *sqlx.Tx, row model.ProductRipening) error {
	_, err := tx.Exec(`
		UPDATE mst.product_ripening
		SET sunday_qty = $1,
			monday_qty = $2,
			tuesday_qty = $3,
			wednesday_qty = $4,
			thursday_qty = $5,
			friday_qty = $6,
			saturday_qty = $7,
			updated_by = $8,
			updated_at = $9,
			is_del = false,
			deleted_by = NULL,
			deleted_at = NULL
		WHERE id = $10`,
		row.SundayQty, row.MondayQty, row.TuesdayQty, row.WednesdayQty, row.ThursdayQty, row.FridayQty, row.SaturdayQty,
		row.UpdatedBy, row.UpdatedAt, row.ID)
	return err
}

func preservePastRipeningValues(target *model.ProductRipening, existing model.ProductRipening, preserveBefore time.Time) {
	idx := weekdayIndexFromDate(preserveBefore)
	past := []int{
		existing.SundayQty,
		existing.MondayQty,
		existing.TuesdayQty,
		existing.WednesdayQty,
		existing.ThursdayQty,
		existing.FridayQty,
		existing.SaturdayQty,
	}
	next := []int{
		target.SundayQty,
		target.MondayQty,
		target.TuesdayQty,
		target.WednesdayQty,
		target.ThursdayQty,
		target.FridayQty,
		target.SaturdayQty,
	}
	for i := 0; i < idx; i++ {
		next[i] = past[i]
	}
	assignRipeningWeekdays(target, next)
}

func zeroFutureRipeningValues(target *model.ProductRipening, preserveBefore time.Time) {
	values := []int{
		target.SundayQty,
		target.MondayQty,
		target.TuesdayQty,
		target.WednesdayQty,
		target.ThursdayQty,
		target.FridayQty,
		target.SaturdayQty,
	}
	idx := weekdayIndexFromDate(preserveBefore)
	for i := idx; i < len(values); i++ {
		values[i] = 0
	}
	assignRipeningWeekdays(target, values)
}

func assignRipeningWeekdays(target *model.ProductRipening, values []int) {
	target.SundayQty = values[0]
	target.MondayQty = values[1]
	target.TuesdayQty = values[2]
	target.WednesdayQty = values[3]
	target.ThursdayQty = values[4]
	target.FridayQty = values[5]
	target.SaturdayQty = values[6]
}

func weekdayIndexFromDate(t time.Time) int {
	switch t.Weekday() {
	case time.Sunday:
		return 0
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	default:
		return 6
	}
}

var _ ProductRipeningRepository = (*productRipeningRepositoryImpl)(nil)
var _ = sql.ErrNoRows
