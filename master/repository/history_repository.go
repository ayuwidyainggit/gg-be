package repository

import (
	"fmt"
	"master/entity"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type HistoryRepository interface {
	ListImportHistory(uploadType, search *string, page, limit int, custId string) ([]entity.ImportHistoryRow, int, int, error)
	GetImportHistoryById(historyId int64) (entity.ImportHistoryRow, error)
	FetchTempRows(table string, historyId int64, custId string, includeCtid bool) ([]string, [][]interface{}, error)
	DeleteTempByCtid(table string, ctid string) error
	CountTemp(table string, historyId int64) (int, error)
	GetImportTotalData(historyId int64) (int, error)
	UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error
	// Fetch import instructions for a given type (e.g., "outlet") to include in reupload templates
	GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error)
	GetDistinctTempValues(table, column string, historyId int64) ([]string, error)
}

type historyRepositoryImpl struct{ *sqlx.DB }

func NewHistoryRepository(db *sqlx.DB) HistoryRepository { return &historyRepositoryImpl{db} }

func (r *historyRepositoryImpl) ListImportHistory(uploadType, search *string, page, limit int, custId string) ([]entity.ImportHistoryRow, int, int, error) {
	rows := []entity.ImportHistoryRow{}
	var total int
	args := []interface{}{}
	whereParts := []string{}
	if custId != "" {
		whereParts = append(whereParts, fmt.Sprintf("ih.cust_id = $%d", len(args)+1))
		args = append(args, custId)
	}
	if uploadType != nil && *uploadType != "" {
		whereParts = append(whereParts, fmt.Sprintf("ih.upload_type = $%d", len(args)+1))
		args = append(args, *uploadType)
	}
	if search != nil && strings.TrimSpace(*search) != "" {
		pattern := "%" + strings.TrimSpace(*search) + "%"
		whereParts = append(whereParts, fmt.Sprintf("(ih.file_name ILIKE $%d OR CAST(ih.history_id AS TEXT) ILIKE $%d OR u.user_name ILIKE $%d OR u.user_fullname ILIKE $%d OR u_by.user_name ILIKE $%d OR u_by.user_fullname ILIKE $%d)", len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5, len(args)+6))
		args = append(args, pattern, pattern, pattern, pattern, pattern, pattern)
	}
	where := ""
	if len(whereParts) > 0 {
		where = "WHERE " + strings.Join(whereParts, " AND ")
	}

	countQuery := "SELECT COUNT(1) FROM import.import_history ih LEFT JOIN sys.m_user u ON u.cust_id = ih.cust_id AND u.user_id = ih.uploaded_by LEFT JOIN sys.m_user u_by ON u_by.user_id = ih.uploaded_by " + where
	if err := r.DB.Get(&total, countQuery, args...); err != nil {
		return nil, 0, 0, err
	}
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	listQuery := "SELECT ih.history_id, ih.file_name, ih.uploaded_by, COALESCE(u.user_fullname, u.user_name, u_by.user_fullname, u_by.user_name) AS uploaded_by_name, ih.upload_date, ih.successful_data, ih.failed_data, ih.total_data, ih.status, ih.status_reupload, ih.upload_type FROM import.import_history ih LEFT JOIN sys.m_user u ON u.cust_id = ih.cust_id AND u.user_id = ih.uploaded_by LEFT JOIN sys.m_user u_by ON u_by.user_id = ih.uploaded_by " + where + " ORDER BY ih.history_id DESC OFFSET $" + fmt.Sprint(len(args)+1) + " LIMIT $" + fmt.Sprint(len(args)+2)
	args = append(args, offset, limit)
	if err := r.DB.Select(&rows, listQuery, args...); err != nil {
		return nil, 0, 0, err
	}
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return rows, total, lastPage, nil
}

func (r *historyRepositoryImpl) GetImportHistoryById(historyId int64) (entity.ImportHistoryRow, error) {
	var row entity.ImportHistoryRow
	err := r.DB.Get(&row, `SELECT ih.history_id, ih.file_name, ih.uploaded_by, COALESCE(u.user_fullname, u.user_name, u_by.user_fullname, u_by.user_name) AS uploaded_by_name, ih.upload_date, ih.successful_data, ih.failed_data, ih.total_data, ih.status, ih.status_reupload, ih.upload_type FROM import.import_history ih LEFT JOIN sys.m_user u ON u.cust_id = ih.cust_id AND u.user_id = ih.uploaded_by LEFT JOIN sys.m_user u_by ON u_by.user_id = ih.uploaded_by WHERE ih.history_id = $1`, historyId)
	return row, err
}

func (r *historyRepositoryImpl) FetchTempRows(table string, historyId int64, custId string, includeCtid bool) ([]string, [][]interface{}, error) {
	selectClause := "*"
	if includeCtid {
		selectClause = "ctid, *"
	}
	query := fmt.Sprintf("SELECT %s FROM %s WHERE history_id = $1 AND cust_id = $2", selectClause, table)
	rows, err := r.DB.Queryx(query, historyId, custId)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	var data [][]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		data = append(data, vals)
	}
	return cols, data, rows.Err()
}

func (r *historyRepositoryImpl) DeleteTempByCtid(table string, ctid string) error {
	_, err := r.DB.Exec(fmt.Sprintf("DELETE FROM %s WHERE ctid = $1", table), ctid)
	return err
}

func (r *historyRepositoryImpl) CountTemp(table string, historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE history_id = $1", table), historyId)
	return n, err
}

func (r *historyRepositoryImpl) GetImportTotalData(historyId int64) (int, error) {
	var n int
	err := r.DB.Get(&n, `SELECT total_data FROM import.import_history WHERE history_id = $1`, historyId)
	return n, err
}

func (r *historyRepositoryImpl) UpdateImportHistory(historyId int64, success, failed int, statusReupload bool) error {
	_, err := r.DB.Exec(`UPDATE import.import_history SET successful_data = $1, failed_data = $2, status_reupload = $3 WHERE history_id = $4`, success, failed, statusReupload, historyId)
	return err
}

func (r *historyRepositoryImpl) GetDistinctTempValues(table, column string, historyId int64) ([]string, error) {
	if strings.TrimSpace(table) == "" || strings.TrimSpace(column) == "" {
		return nil, nil
	}
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE history_id = $1 AND %s IS NOT NULL AND btrim(%s::text) <> ''", column, table, column, column)
	values := []string{}
	if err := r.DB.Select(&values, query, historyId); err != nil {
		return nil, err
	}
	return values, nil
}

func (r *historyRepositoryImpl) GetImportInstructions(instructionType string) ([]entity.ImportInstruction, error) {
	rows := []entity.ImportInstruction{}

	query := `
		SELECT 
			instruction_id, 
			instruction_type, 
			kolom, 
			mandatory, 
			keterangan, 
			step
		FROM import.import_instructions
		WHERE instruction_type = $1
		ORDER BY 
			CASE 
				WHEN step ILIKE 'Step 1%' THEN 1
				WHEN step ILIKE 'Step 2%' THEN 2
				WHEN step ILIKE 'Step 3%' THEN 3
				WHEN step ILIKE 'Step 4%' THEN 4
				ELSE 5
			END,
			instruction_id;
	`

	if err := r.DB.Select(&rows, query, instructionType); err != nil {
		return nil, err
	}
	return rows, nil
}
