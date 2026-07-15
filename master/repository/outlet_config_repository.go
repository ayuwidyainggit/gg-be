package repository

import (
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type OutletConfigRepository interface {
	FindAll(filter entity.OutletConfigListFilter, custIDs []string) ([]model.OutletConfig, int, int, error)
	FindByOutletConfigID(outletConfigID int, custID string) (*model.OutletConfigHeader, []model.OutletConfigDet, error)
	FindHeaderByID(outletConfigID int) (*model.OutletConfigHeader, error)
	GetCustIdsByParentCustId(parentCustId string) ([]string, error)
	FindActiveByCustId(custID string) (*model.OutletConfigHeader, []model.OutletConfigDet, error)
	Store(custID string, body entity.CreateOutletConfigBody, createdBy int64) error
	Update(outletConfigID int, custID string, body entity.CreateOutletConfigBody, updatedBy int64) error
	Delete(outletConfigID int, custID string, updatedBy int64) error
}

type OutletConfigStatusRepository interface {
	FindAll(filter entity.OutletConfigStatusListFilter) ([]model.OutletConfigStatus, int, int, error)
}

func NewOutletConfigStatusRepository(db *sqlx.DB) OutletConfigStatusRepository {
	return &outletConfigStatusRepositoryImpl{db}
}

type outletConfigStatusRepositoryImpl struct {
	*sqlx.DB
}

func NewOutletConfigRepository(db *sqlx.DB) OutletConfigRepository {
	return &outletConfigRepositoryImpl{db}
}

type outletConfigRepositoryImpl struct {
	*sqlx.DB
}

func (r *outletConfigRepositoryImpl) FindAll(filter entity.OutletConfigListFilter, custIDs []string) ([]model.OutletConfig, int, int, error) {
	list := []model.OutletConfig{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.outlet_config_id, a.verification_status, a.rules_type,
		a.created_by, a.created_at, a.updated_by, a.updated_at, a.status,
		u1.user_fullname AS created_by_name, u2.user_fullname AS updated_by_name `
	qWhere := ` WHERE 1=1 `
	args := []interface{}{}
	argIdx := 1

	if len(custIDs) > 0 {
		qWhere += ` AND a.cust_id IN (`
		for i, cid := range custIDs {
			if i > 0 {
				qWhere += `,`
			}
			qWhere += `$` + strconv.Itoa(argIdx)
			args = append(args, cid)
			argIdx++
		}
		qWhere += `)`
	}
	if filter.Q != "" {
		qWhere += ` AND (a.verification_status ILIKE $` + strconv.Itoa(argIdx) + ` OR a.rules_type ILIKE $` + strconv.Itoa(argIdx) + `) `
		args = append(args, "%"+filter.Q+"%")
		argIdx++
	}
	if filter.RulesType != "" {
		qWhere += ` AND a.rules_type = $` + strconv.Itoa(argIdx)
		args = append(args, filter.RulesType)
		argIdx++
	}
	if filter.Status != nil {
		if *filter.Status == 1 {
			qWhere += ` AND a.status = 1 `
		} else if *filter.Status == 0 {
			qWhere += ` AND a.status = 0 `
		}
	}

	qWhere += ` AND a.is_del = false `
	qFrom := ` FROM mst.m_outlet_config a
		LEFT JOIN sys.m_user u1 ON u1.user_id = a.created_by
		LEFT JOIN sys.m_user u2 ON u2.user_id = a.updated_by `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	var total int
	row := r.QueryRowx(queryCount, args...)
	if err := row.Scan(&total); err != nil {
		log.Println("outletConfigRepository FindAll count err:", err)
		return list, 0, 0, err
	}

	sortBy := `a.created_at DESC`
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ",")
		var sortParts []string
		for _, p := range parts {
			colOrder := strings.Split(strings.TrimSpace(p), ":")
			if len(colOrder) >= 2 {
				col := strings.TrimSpace(colOrder[0])
				order := strings.TrimSpace(strings.ToUpper(colOrder[1]))
				if order != "ASC" && order != "DESC" {
					order = "DESC"
				}
				if col == "created_date" {
					col = "created_at"
				}
				if col == "updated_date" {
					col = "updated_at"
				}
				sortParts = append(sortParts, fmt.Sprintf("a.%s %s", col, order))
			}
		}
		if len(sortParts) > 0 {
			sortBy = strings.Join(sortParts, ", ")
		}
	}
	querySelect := `SELECT ` + selectField + qFrom + qWhere + ` ORDER BY ` + sortBy

	limit := filter.Limit
	if limit <= 0 {
		limit = 5
	}
	page := filter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if total == 0 {
		lastPage = 0
	}

	querySelect += ` LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	if err := r.Select(&list, querySelect, args...); err != nil {
		log.Println("outletConfigRepository FindAll select err:", err)
		return list, 0, 0, err
	}
	return list, total, lastPage, nil
}

func (r *outletConfigRepositoryImpl) FindByOutletConfigID(outletConfigID int, custID string) (*model.OutletConfigHeader, []model.OutletConfigDet, error) {
	header := &model.OutletConfigHeader{}
	queryHeader := `SELECT outlet_config_id, cust_id, verification_status, rules_type
		FROM mst.m_outlet_config
		WHERE outlet_config_id = $1 AND cust_id = $2 AND is_del = false`
	if err := r.Get(header, queryHeader, outletConfigID, custID); err != nil {
		log.Println("outletConfigRepository FindByOutletConfigID header err:", err)
		return nil, nil, err
	}
	details := []model.OutletConfigDet{}
	queryDet := `SELECT d.outlet_config_det_id, d.outlet_config_id, d.status, d.validate_trx, d.counting_period,
		s.status_description
		FROM mst.m_outlet_config_det d
		LEFT JOIN mst.m_outlet_config_status s 
			ON s.status_code = d.status::varchar AND s.is_del = false
		WHERE d.outlet_config_id = $1 AND d.is_del = false
		ORDER BY d.outlet_config_det_id`
	if err := r.Select(&details, queryDet, outletConfigID); err != nil {
		log.Println("outletConfigRepository FindByOutletConfigID details err:", err)
		return nil, nil, err
	}
	return header, details, nil
}

func (r *outletConfigRepositoryImpl) FindHeaderByID(outletConfigID int) (*model.OutletConfigHeader, error) {
	header := &model.OutletConfigHeader{}
	query := `SELECT outlet_config_id, cust_id, verification_status, rules_type
		FROM mst.m_outlet_config
		WHERE outlet_config_id = $1 AND is_del = false`
	if err := r.Get(header, query, outletConfigID); err != nil {
		return nil, err
	}
	return header, nil
}

func (r *outletConfigRepositoryImpl) GetCustIdsByParentCustId(parentCustId string) ([]string, error) {
	if parentCustId == "" {
		return nil, nil
	}
	var rows []struct {
		CustId string `db:"cust_id"`
	}
	query := `SELECT cust_id FROM smc.m_customer WHERE parent_cust_id = $1`
	if err := r.Select(&rows, query, parentCustId); err != nil {
		log.Println("outletConfigRepository GetCustIdsByParentCustId err:", err)
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.CustId)
	}
	return ids, nil
}

func (r *outletConfigRepositoryImpl) FindActiveByCustId(custID string) (*model.OutletConfigHeader, []model.OutletConfigDet, error) {
	header := &model.OutletConfigHeader{}
	q := `SELECT outlet_config_id, verification_status, rules_type
		FROM mst.m_outlet_config
		WHERE cust_id = $1 AND status = 1 AND is_del = false
		ORDER BY outlet_config_id ASC LIMIT 1`
	if err := r.Get(header, q, custID); err != nil {
		return nil, nil, err
	}
	details := []model.OutletConfigDet{}
	qDet := `SELECT d.outlet_config_det_id, d.outlet_config_id, d.status, d.validate_trx, d.counting_period,
		s.status_description
		FROM mst.m_outlet_config_det d
		LEFT JOIN mst.m_outlet_config_status s 
			ON s.status_code = d.status::varchar AND s.is_del = false
		WHERE d.outlet_config_id = $1 AND d.is_del = false
		ORDER BY d.outlet_config_det_id ASC`
	if err := r.Select(&details, qDet, header.OutletConfigId); err != nil {
		return nil, nil, err
	}
	return header, details, nil
}

func (r *outletConfigRepositoryImpl) Store(custID string, body entity.CreateOutletConfigBody, createdBy int64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	insertHeader := `INSERT INTO mst.m_outlet_config (
		cust_id, verification_status, rules_type,
		created_by, created_at, updated_by, updated_at, status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, 1)
	RETURNING outlet_config_id`
	var outletConfigID int
	if err = tx.QueryRowx(insertHeader,
		custID, body.VerificationStatus, body.RulesType,
		createdBy, now, createdBy, now,
	).Scan(&outletConfigID); err != nil {
		log.Println("outletConfigRepository Store header err:", err)
		return err
	}

	insertDet := `INSERT INTO mst.m_outlet_config_det (
		cust_id, outlet_config_id, status, validate_trx, counting_period,
		created_by, created_at, updated_by, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, d := range body.OutletStatus {
		if _, err = tx.Exec(insertDet, custID, outletConfigID, d.Status, d.ValidateTrx, d.CountingPeriod,
			createdBy, now, createdBy, now); err != nil {
			log.Println("outletConfigRepository Store detail err:", err)
			return err
		}
	}
	return tx.Commit()
}

func (r *outletConfigRepositoryImpl) Update(outletConfigID int, custID string, body entity.CreateOutletConfigBody, updatedBy int64) error {
	tx, err := r.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()
	updateHeader := `UPDATE mst.m_outlet_config SET
		verification_status = $1, rules_type = $2, updated_by = $3, updated_at = $4
		WHERE outlet_config_id = $5 AND cust_id = $6`
	res, err := tx.Exec(updateHeader, body.VerificationStatus, body.RulesType, updatedBy, now, outletConfigID, custID)
	if err != nil {
		log.Println("outletConfigRepository Update header err:", err)
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("sql: no rows in result set")
	}

	_, err = tx.Exec(`DELETE FROM mst.m_outlet_config_det WHERE outlet_config_id = $1`, outletConfigID)
	if err != nil {
		log.Println("outletConfigRepository Update delete details err:", err)
		return err
	}

	insertDet := `INSERT INTO mst.m_outlet_config_det (
		cust_id, outlet_config_id, status, validate_trx, counting_period,
		created_by, created_at, updated_by, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	for _, d := range body.OutletStatus {
		if _, err = tx.Exec(insertDet, custID, outletConfigID, d.Status, d.ValidateTrx, d.CountingPeriod,
			updatedBy, now, updatedBy, now); err != nil {
			log.Println("outletConfigRepository Update insert detail err:", err)
			return err
		}
	}
	return tx.Commit()
}

func (r *outletConfigRepositoryImpl) Delete(outletConfigID int, custID string, updatedBy int64) error {
	now := time.Now().UTC()
	res, err := r.Exec(`UPDATE mst.m_outlet_config SET is_del = true, updated_by = $1, updated_at = $2 WHERE outlet_config_id = $3 AND cust_id = $4`,
		updatedBy, now, outletConfigID, custID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("sql: no rows in result set")
	}
	return nil
}

func (r *outletConfigStatusRepositoryImpl) FindAll(filter entity.OutletConfigStatusListFilter) ([]model.OutletConfigStatus, int, int, error) {
	list := []model.OutletConfigStatus{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` outlet_config_status_id, status_code, status_description, is_trx, sort_order, is_active `
	qWhere := ` WHERE 1=1 `
	args := []interface{}{}
	argIdx := 1

	if filter.Q != "" {
		qWhere += ` AND (status_code ILIKE $` + strconv.Itoa(argIdx) + ` OR status_description ILIKE $` + strconv.Itoa(argIdx) + `) `
		args = append(args, "%"+filter.Q+"%")
		argIdx++
	}

	if filter.IsActive != nil {
		if *filter.IsActive {
			qWhere += ` AND is_active = true `
		} else {
			qWhere += ` AND is_active = false `
		}
	}

	qFrom := ` FROM mst.m_outlet_config_status `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	var total int
	row := r.QueryRowx(queryCount, args...)
	if err := row.Scan(&total); err != nil {
		log.Println("outletConfigStatusRepository FindAll count err:", err)
		return list, 0, 0, err
	}

	sortBy := `sort_order ASC`
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ",")
		var sortParts []string
		for _, p := range parts {
			colOrder := strings.Split(strings.TrimSpace(p), ":")
			if len(colOrder) >= 2 {
				col := strings.TrimSpace(colOrder[0])
				order := strings.TrimSpace(strings.ToUpper(colOrder[1]))
				if order != "ASC" && order != "DESC" {
					order = "ASC"
				}
				if col == "updated_date" {
					col = "updated_at"
				}
				if col == "created_date" {
					col = "created_at"
				}
				sortParts = append(sortParts, fmt.Sprintf("%s %s", col, order))
			}
		}
		if len(sortParts) > 0 {
			sortBy = strings.Join(sortParts, ", ")
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
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if total == 0 {
		lastPage = 0
	}

	querySelect += ` LIMIT $` + strconv.Itoa(argIdx) + ` OFFSET $` + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	if err := r.Select(&list, querySelect, args...); err != nil {
		log.Println("outletConfigStatusRepository FindAll select err:", err)
		return list, 0, 0, err
	}
	return list, total, lastPage, nil
}
