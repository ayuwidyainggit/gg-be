package repository

import (
	"database/sql"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OutletCodeRepository interface {
	List(filter entity.OutletCodeListFilter, custIDs []string) ([]model.OutletCode, int, int, error)
	ExistsBySerialCodeAndYearCode(custId, serialCode string, yearCode int) (bool, error)
	Store(row model.OutletCode) error
	GetByID(id string) (*model.OutletCode, error)
	ExistsBySerialCodeAndYearCodeExceptID(custId, serialCode string, yearCode int, excludeID string) (bool, error)
	UpdateSerialCode(id, serialCode string, updatedBy *string) error
	UpdateStatus(id, status string, updatedBy *string) error
	GetActiveConfigForCustIdAndYear(custId string, year int) (*model.OutletCode, error)
	UpdateLastSequenceNo(id, newSequenceStr string, updatedBy *string) error
	GetActiveConfigAndIncrement(custId string, currentYear int, updatedBy *string) (*model.OutletCode, string, error)
	GetActiveConfigAndIncrementByCreatedBy(custId string, currentYear int, createdBy string, updatedBy *string) (*model.OutletCode, string, error)
	IncrementSequenceByID(id string, updatedBy *string) (string, error)
	// GetActiveConfigForUpdate returns Active config and next sequence (locks row, no update). Caller must UpdateLastSequenceNoWithTx after inserting outlet.
	GetActiveConfigForUpdate(tx *sqlx.Tx, custId string, currentYear int) (*model.OutletCode, string, error)
	GetActiveConfigForUpdateByCreatedBy(tx *sqlx.Tx, custId string, currentYear int, createdBy string) (*model.OutletCode, string, error)
	UpdateLastSequenceNoWithTx(tx *sqlx.Tx, id, newSequenceStr string, updatedBy *string) error
	FindOneByCustIdYearAndStatus(custId string, year int, statuses []string) (*model.OutletCode, error)
	FindOneByCustIdYearAndStatusAndCreatedBy(custId string, year int, statuses []string, createdBy string) (*model.OutletCode, error)
}

func NewOutletCodeRepository(db *sqlx.DB) OutletCodeRepository {
	return &outletCodeRepositoryImpl{db}
}

type outletCodeRepositoryImpl struct {
	*sqlx.DB
}

func (r *outletCodeRepositoryImpl) List(filter entity.OutletCodeListFilter, custIDs []string) ([]model.OutletCode, int, int, error) {
	list := []model.OutletCode{}
	qWhere := ` WHERE 1=1 `
	args := []interface{}{}
	argIdx := 1

	if filter.CustId != "" {
		qWhere += ` AND a.cust_id = $` + strconv.Itoa(argIdx)
		args = append(args, filter.CustId)
		argIdx++
	} else if len(custIDs) > 0 {
		placeholders := make([]string, len(custIDs))
		for i := range custIDs {
			placeholders[i] = `$` + strconv.Itoa(argIdx)
			args = append(args, custIDs[i])
			argIdx++
		}
		qWhere += ` AND a.cust_id IN (` + strings.Join(placeholders, ",") + `) `
	}

	if filter.Q != "" {
		esc := strings.ReplaceAll(filter.Q, "'", "''")
		qWhere += ` AND (a.serial_code ILIKE '%` + esc + `%' OR a.last_sequence_no ILIKE '%` + esc + `%') `
	}

	if len(filter.Status) > 0 {
		statusMap := map[string]string{
			"active":     "Active",
			"deactive":   "Deactivate",
			"non_active": "Non Active",
		}
		var in []string
		for _, s := range filter.Status {
			s = strings.TrimSpace(strings.ToLower(s))
			if v, ok := statusMap[s]; ok {
				in = append(in, "'"+strings.ReplaceAll(v, "'", "''")+"'")
			}
		}
		if len(in) > 0 {
			qWhere += ` AND a.status IN (` + strings.Join(in, ",") + `) `
		}
	}

	selectCount := ` COUNT(*) AS total `
	selectField := ` a.id, a.cust_id, a.serial_code, a.year_code, a.last_sequence_no, a.status,
		a.created_at, a.created_by, a.updated_at, a.updated_by `
	qFrom := ` FROM mst.m_outlet_code a `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	var total int
	if err := r.QueryRow(queryCount, args...).Scan(&total); err != nil {
		log.Println("outlet_code_repository List count err:", err.Error())
		return list, 0, 0, err
	}

	sortBy := `a.created_at DESC`
	if filter.Sort != "" {
		parts := strings.SplitN(filter.Sort, ":", 2)
		if len(parts) == 2 {
			col := strings.TrimSpace(parts[0])
			dir := strings.TrimSpace(strings.ToUpper(parts[1]))
			if dir != "ASC" && dir != "DESC" {
				dir = "DESC"
			}
			allowedCols := map[string]bool{
				"created_at": true, "updated_at": true, "serial_code": true, "year_code": true, "status": true,
			}
			if allowedCols[col] {
				sortBy = fmt.Sprintf(`a.%s %s`, col, dir)
			}
		}
	}
	querySelect := `SELECT ` + selectField + qFrom + qWhere + ` ORDER BY ` + sortBy

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, limit, offset)

	if err := r.Select(&list, querySelect, args...); err != nil {
		log.Println("outlet_code_repository List select err:", err.Error())
		return list, 0, 0, err
	}

	lastPage := 0
	if limit > 0 {
		lastPage = int(math.Ceil(float64(total) / float64(limit)))
	}
	return list, total, lastPage, nil
}

func (r *outletCodeRepositoryImpl) ExistsBySerialCodeAndYearCode(custId, serialCode string, yearCode int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM mst.m_outlet_code WHERE cust_id = $1 AND serial_code = $2 AND year_code = $3)`
	err := r.QueryRow(query, custId, serialCode, yearCode).Scan(&exists)
	return exists, err
}

func (r *outletCodeRepositoryImpl) Store(row model.OutletCode) error {
	query := `INSERT INTO mst.m_outlet_code (cust_id, serial_code, year_code, last_sequence_no, status, created_at, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, NOW(), $7)`
	_, err := r.Exec(query, row.CustId, row.SerialCode, row.YearCode, row.LastSequenceNo, row.Status, row.CreatedBy, row.UpdatedBy)
	return err
}

func (r *outletCodeRepositoryImpl) GetByID(id string) (*model.OutletCode, error) {
	var row model.OutletCode
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status,
		created_at, created_by, updated_at, updated_by
		FROM mst.m_outlet_code WHERE id = $1`
	if err := r.Get(&row, query, id); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *outletCodeRepositoryImpl) ExistsBySerialCodeAndYearCodeExceptID(custId, serialCode string, yearCode int, excludeID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(
		SELECT 1 FROM mst.m_outlet_code
		WHERE cust_id = $1 AND serial_code = $2 AND year_code = $3 AND id <> $4
	)`
	err := r.QueryRow(query, custId, serialCode, yearCode, excludeID).Scan(&exists)
	return exists, err
}

func (r *outletCodeRepositoryImpl) UpdateSerialCode(id, serialCode string, updatedBy *string) error {
	query := `UPDATE mst.m_outlet_code
		SET serial_code = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3`
	_, err := r.Exec(query, serialCode, updatedBy, id)
	return err
}

func (r *outletCodeRepositoryImpl) UpdateStatus(id, status string, updatedBy *string) error {
	query := `UPDATE mst.m_outlet_code
		SET status = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3`
	_, err := r.Exec(query, status, updatedBy, id)
	return err
}

func (r *outletCodeRepositoryImpl) GetActiveConfigForCustIdAndYear(custId string, year int) (*model.OutletCode, error) {
	var row model.OutletCode
	year2 := year % 100
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status,
		created_at, created_by, updated_at, updated_by
		FROM mst.m_outlet_code
		WHERE cust_id = $1 AND status = 'Active' AND (year_code = $2 OR year_code = $3)
		ORDER BY created_at DESC LIMIT 1`
	if err := r.Get(&row, query, custId, year, year2); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *outletCodeRepositoryImpl) UpdateLastSequenceNo(id, newSequenceStr string, updatedBy *string) error {
	query := `UPDATE mst.m_outlet_code
		SET last_sequence_no = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3`
	_, err := r.Exec(query, newSequenceStr, updatedBy, id)
	return err
}

func (r *outletCodeRepositoryImpl) GetActiveConfigAndIncrement(custId string, currentYear int, updatedBy *string) (*model.OutletCode, string, error) {
	tx, err := r.Beginx()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`SET LOCAL lock_timeout = '1000ms'`); err != nil {
		return nil, "", err
	}

	year2 := currentYear % 100
	var row model.OutletCode
	query := `WITH target AS (
			SELECT id
			FROM mst.m_outlet_code
			WHERE cust_id = $1 AND status = 'Active' AND (year_code = $2 OR year_code = $3)
			ORDER BY created_at DESC
			LIMIT 1
		),
		updated AS (
			UPDATE mst.m_outlet_code oc
			SET last_sequence_no = LPAD((COALESCE(NULLIF(TRIM(oc.last_sequence_no), ''), '0')::INT + 1)::TEXT, 4, '0'),
				updated_at = NOW(),
				updated_by = $4
			FROM target t
			WHERE oc.id = t.id
			RETURNING oc.id, oc.cust_id, oc.serial_code, oc.year_code, oc.last_sequence_no, oc.status
		)
		SELECT id, cust_id, serial_code, year_code, last_sequence_no, status
		FROM updated`
	if err := tx.Get(&row, query, custId, currentYear, year2, updatedBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		if strings.Contains(strings.ToLower(err.Error()), "lock timeout") {
			return nil, "", fmt.Errorf("outlet code sequence is busy, please retry")
		}
		return nil, "", err
	}

	nextSeq, _ := strconv.Atoi(strings.TrimSpace(row.LastSequenceNo))
	if nextSeq > 9999 {
		return nil, "", fmt.Errorf("outlet code sequence exceeded 9999 for year %d", currentYear)
	}
	if err := tx.Commit(); err != nil {
		return nil, "", err
	}
	return &row, row.LastSequenceNo, nil
}

func (r *outletCodeRepositoryImpl) GetActiveConfigAndIncrementByCreatedBy(custId string, currentYear int, createdBy string, updatedBy *string) (*model.OutletCode, string, error) {
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		return nil, "", nil
	}

	tx, err := r.Beginx()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`SET LOCAL lock_timeout = '1000ms'`); err != nil {
		return nil, "", err
	}

	year2 := currentYear % 100
	var row model.OutletCode
	query := `WITH target AS (
			SELECT id
			FROM mst.m_outlet_code
			WHERE cust_id = $1 AND status = 'Active' AND (year_code = $2 OR year_code = $3) AND created_by = $4
			ORDER BY created_at DESC
			LIMIT 1
		),
		updated AS (
			UPDATE mst.m_outlet_code oc
			SET last_sequence_no = LPAD((COALESCE(NULLIF(TRIM(oc.last_sequence_no), ''), '0')::INT + 1)::TEXT, 4, '0'),
				updated_at = NOW(),
				updated_by = $5
			FROM target t
			WHERE oc.id = t.id
			RETURNING oc.id, oc.cust_id, oc.serial_code, oc.year_code, oc.last_sequence_no, oc.status
		)
		SELECT id, cust_id, serial_code, year_code, last_sequence_no, status
		FROM updated`
	if err := tx.Get(&row, query, custId, currentYear, year2, createdBy, updatedBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		if strings.Contains(strings.ToLower(err.Error()), "lock timeout") {
			return nil, "", fmt.Errorf("outlet code sequence is busy, please retry")
		}
		return nil, "", err
	}

	nextSeq, _ := strconv.Atoi(strings.TrimSpace(row.LastSequenceNo))
	if nextSeq > 9999 {
		return nil, "", fmt.Errorf("outlet code sequence exceeded 9999 for year %d", currentYear)
	}
	if err := tx.Commit(); err != nil {
		return nil, "", err
	}
	return &row, row.LastSequenceNo, nil
}

func (r *outletCodeRepositoryImpl) IncrementSequenceByID(id string, updatedBy *string) (string, error) {
	var nextSeq string
	query := `WITH base AS (
			SELECT COALESCE(NULLIF(TRIM(last_sequence_no), ''), '0')::INT AS base_seq
			FROM mst.m_outlet_code
			WHERE id = $1
		),
		upsert AS (
			INSERT INTO mst.m_outlet_code_seq(outlet_code_id, last_sequence_no, updated_by, created_at, updated_at)
			SELECT $1::uuid, base_seq + 1, $2, NOW(), NOW()
			FROM base
			ON CONFLICT (outlet_code_id)
			DO UPDATE SET
				last_sequence_no = mst.m_outlet_code_seq.last_sequence_no + 1,
				updated_by = EXCLUDED.updated_by,
				updated_at = NOW()
			RETURNING last_sequence_no
		)
		SELECT LPAD(last_sequence_no::TEXT, 4, '0') FROM upsert`
	if err := r.QueryRow(query, id, updatedBy).Scan(&nextSeq); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	seqInt, _ := strconv.Atoi(strings.TrimSpace(nextSeq))
	if seqInt > 9999 {
		return "", fmt.Errorf("outlet code sequence exceeded 9999")
	}
	return nextSeq, nil
}

// GetActiveConfigForUpdate locks row (FOR UPDATE) and returns config + next sequence string. Caller must call UpdateLastSequenceNoWithTx after inserting outlet.
func (r *outletCodeRepositoryImpl) GetActiveConfigForUpdate(tx *sqlx.Tx, custId string, currentYear int) (*model.OutletCode, string, error) {
	year2 := currentYear % 100
	var row model.OutletCode
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status
		FROM mst.m_outlet_code
		WHERE cust_id = $1 AND status = 'Active' AND (year_code = $2 OR year_code = $3)
		ORDER BY created_at DESC LIMIT 1 FOR UPDATE`
	if err := tx.Get(&row, query, custId, currentYear, year2); err != nil {
		return nil, "", err
	}

	currentSeq, _ := strconv.Atoi(strings.TrimSpace(row.LastSequenceNo))
	if currentSeq < 0 {
		currentSeq = 0
	}
	nextSeq := currentSeq + 1
	nextSeqStr := fmt.Sprintf("%04d", nextSeq)
	if nextSeq > 9999 {
		return nil, "", fmt.Errorf("outlet code sequence exceeded 9999 for year %d", currentYear)
	}
	return &row, nextSeqStr, nil
}

func (r *outletCodeRepositoryImpl) GetActiveConfigForUpdateByCreatedBy(tx *sqlx.Tx, custId string, currentYear int, createdBy string) (*model.OutletCode, string, error) {
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		return nil, "", nil
	}
	year2 := currentYear % 100
	var row model.OutletCode
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status
		FROM mst.m_outlet_code
		WHERE cust_id = $1 AND status = 'Active' AND (year_code = $2 OR year_code = $3) AND created_by = $4
		ORDER BY created_at DESC LIMIT 1 FOR UPDATE`
	if err := tx.Get(&row, query, custId, currentYear, year2, createdBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}

	currentSeq, _ := strconv.Atoi(strings.TrimSpace(row.LastSequenceNo))
	if currentSeq < 0 {
		currentSeq = 0
	}
	nextSeq := currentSeq + 1
	nextSeqStr := fmt.Sprintf("%04d", nextSeq)
	if nextSeq > 9999 {
		return nil, "", fmt.Errorf("outlet code sequence exceeded 9999 for year %d", currentYear)
	}
	return &row, nextSeqStr, nil
}

func (r *outletCodeRepositoryImpl) UpdateLastSequenceNoWithTx(tx *sqlx.Tx, id, newSequenceStr string, updatedBy *string) error {
	query := `UPDATE mst.m_outlet_code SET last_sequence_no = $1, updated_at = NOW(), updated_by = $2 WHERE id = $3`
	_, err := tx.Exec(query, newSequenceStr, updatedBy, id)
	return err
}

func (r *outletCodeRepositoryImpl) FindOneByCustIdYearAndStatus(custId string, year int, statuses []string) (*model.OutletCode, error) {
	if len(statuses) == 0 {
		return nil, nil
	}
	statusMap := map[string]string{
		"active":     "Active",
		"deactive":   "Deactivate",
		"non_active": "Non Active",
	}
	var in []string
	for _, s := range statuses {
		s = strings.TrimSpace(strings.ToLower(s))
		if v, ok := statusMap[s]; ok {
			in = append(in, "'"+strings.ReplaceAll(v, "'", "''")+"'")
		} else {
			in = append(in, "'"+strings.ReplaceAll(strings.TrimSpace(s), "'", "''")+"'")
		}
	}
	if len(in) == 0 {
		return nil, nil
	}
	yearShort := year % 100
	var row model.OutletCode
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status,
		created_at, created_by, updated_at, updated_by
		FROM mst.m_outlet_code
		WHERE cust_id = $1 AND (year_code = $2 OR year_code = $3) AND status IN (` + strings.Join(in, ",") + `)
		ORDER BY created_at DESC LIMIT 1`
	if err := r.Get(&row, query, custId, year, yearShort); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

func (r *outletCodeRepositoryImpl) FindOneByCustIdYearAndStatusAndCreatedBy(custId string, year int, statuses []string, createdBy string) (*model.OutletCode, error) {
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		return nil, nil
	}
	if len(statuses) == 0 {
		return nil, nil
	}
	statusMap := map[string]string{
		"active":     "Active",
		"deactive":   "Deactivate",
		"non_active": "Non Active",
	}
	var in []string
	for _, s := range statuses {
		s = strings.TrimSpace(strings.ToLower(s))
		if v, ok := statusMap[s]; ok {
			in = append(in, "'"+strings.ReplaceAll(v, "'", "''")+"'")
		} else {
			in = append(in, "'"+strings.ReplaceAll(strings.TrimSpace(s), "'", "''")+"'")
		}
	}
	if len(in) == 0 {
		return nil, nil
	}
	yearShort := year % 100
	var row model.OutletCode
	query := `SELECT id, cust_id, serial_code, year_code, last_sequence_no, status,
		created_at, created_by, updated_at, updated_by
		FROM mst.m_outlet_code
		WHERE cust_id = $1 AND (year_code = $2 OR year_code = $3) AND status IN (` + strings.Join(in, ",") + `)
		AND created_by = $4
		ORDER BY created_at DESC LIMIT 1`
	if err := r.Get(&row, query, custId, year, yearShort, createdBy); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
