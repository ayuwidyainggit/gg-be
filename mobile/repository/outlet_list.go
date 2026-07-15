package repository

import (
	"errors"
	"fmt"
	"log"
	"math"
	"mobile/entity"
	"mobile/model"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type OutletListRepository interface {
	FindAll(filter entity.OutletListQueryFilter, custId string) ([]model.MOutletRead, int64, int, error)
	SoftDelete(custId string, outletId int64, deletedBy int64) error
	Update(custId string, outletId int64, data map[string]interface{}) error
	UpdateContact(custId string, outletId int64, contact entity.UpdateOutletContact) error
}

type outletListRepositoryImpl struct {
	*sqlx.DB
}

func NewOutletListRepository(db *sqlx.DB) OutletListRepository {
	return &outletListRepositoryImpl{db}
}

// FindAll - GET /v1/outlet-list
func (r *outletListRepositoryImpl) FindAll(filter entity.OutletListQueryFilter, custId string) ([]model.MOutletRead, int64, int, error) {
	outlets := []model.MOutletRead{}

	selectCount := ` COUNT(*) AS total `
	selectField := ` 
        o.outlet_id, 
        o.outlet_code, 
        o.outlet_name, 
        o.address1, 
        o.longitude, 
        o.latitude, 
        COALESCE(o.outlet_status, 0) AS outlet_status
    `

	qFrom := ` FROM mst.m_outlet o `
	qWhere := ` WHERE o.is_del = false AND o.cust_id = $1 `
	args := []interface{}{custId}
	argIdx := 2

	// Filter by search query
	if filter.Query != "" {
		qWhere += fmt.Sprintf(` AND (o.outlet_code ILIKE $%d OR o.outlet_name ILIKE $%d) `, argIdx, argIdx)
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	// Filter by outlet_status
	if len(filter.OutletStatus) > 0 {
		placeholders := make([]string, len(filter.OutletStatus))
		for i, status := range filter.OutletStatus {
			placeholders[i] = fmt.Sprintf("$%d", argIdx)
			args = append(args, status)
			argIdx++
		}
		qWhere += ` AND o.outlet_status IN (` + strings.Join(placeholders, ",") + `) `
	}

	// Filter by is_active
	if filter.IsActive != nil {
		if *filter.IsActive == 1 {
			qWhere += ` AND o.is_active = true `
		} else if *filter.IsActive == 0 {
			qWhere += ` AND o.is_active = false `
		}
	}

	// Count total
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	var total int64
	err := r.QueryRow(queryCount, args...).Scan(&total)
	if err != nil {
		log.Println("OutletListRepository, FindAll, count error:", err.Error())
		return outlets, 0, 0, err
	}

	// Build select query
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// Sorting
	if filter.Sort != "" {
		sortParts := strings.Split(filter.Sort, ":")
		if len(sortParts) == 2 {
			col := sortParts[0]
			dir := strings.ToUpper(sortParts[1])
			if dir != "ASC" && dir != "DESC" {
				dir = "ASC"
			}
			querySelect += fmt.Sprintf(` ORDER BY %s %s `, col, dir)
		}
	} else {
		querySelect += ` ORDER BY outlet_code ASC `
	}

	// Pagination
	if filter.Limit <= 0 || filter.Limit > 999 {
		filter.Limit = 5
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(float64(filter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %d OFFSET %d`, filter.Limit, offset)

	// Execute query
	err = r.Select(&outlets, querySelect, args...)
	if err != nil {
		log.Println("OutletListRepository, FindAll, select error:", err.Error())
		return outlets, total, lastPage, err
	}

	return outlets, total, lastPage, nil
}

// SoftDelete - DELETE /v1/outlet-list/:outlet_id
func (r *outletListRepositoryImpl) SoftDelete(custId string, outletId int64, deletedBy int64) error {
	query := `
        UPDATE mst.m_outlet
        SET 
            is_del = true,
            deleted_at = $1,
            deleted_by = $2
        WHERE is_del = false
        AND cust_id = $3
        AND outlet_id = $4
    `

	now := time.Now().UTC()
	result, err := r.Exec(query, now, deletedBy, custId, outletId)
	if err != nil {
		log.Println("OutletListRepository, SoftDelete, error:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to get affected rows")
	}
	if nRows == 0 {
		return errors.New("record not found")
	}

	return nil
}

// Update - PATCH /v1/m-outlets/:outlet_id
func (r *outletListRepositoryImpl) Update(custId string, outletId int64, data map[string]interface{}) error {
	// Build dynamic update query
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	// Always update updated_at and updated_by or verified_by
	data["updated_at"] = time.Now().UTC()

	for key, value := range data {
		if value != nil {
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, argIdx))
			args = append(args, value)
			argIdx++
		}
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf(`
        UPDATE mst.m_outlet
        SET %s
        WHERE is_del = false
        AND cust_id = $%d
        AND outlet_id = $%d
    `, strings.Join(setClauses, ", "), argIdx, argIdx+1)

	args = append(args, custId, outletId)

	result, err := r.Exec(query, args...)
	if err != nil {
		log.Println("OutletListRepository, Update, error:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to get affected rows")
	}
	if nRows == 0 {
		return errors.New("record not found")
	}

	return nil
}

// UpdateContact - Update outlet contact for PATCH /v1/m-outlets/:outlet_id
func (r *outletListRepositoryImpl) UpdateContact(custId string, outletId int64, contact entity.UpdateOutletContact) error {
	// If outlet_contact_id exists, update by outlet_contact_id
	if contact.OutletContactId != nil && *contact.OutletContactId > 0 {
		// Update existing contact by outlet_contact_id
		setClauses := []string{}
		args := []interface{}{}
		argIdx := 1

		if contact.ContactName != nil {
			setClauses = append(setClauses, fmt.Sprintf("contact_name = $%d", argIdx))
			args = append(args, *contact.ContactName)
			argIdx++
		}
		if contact.JobTitle != nil {
			setClauses = append(setClauses, fmt.Sprintf("job_title = $%d", argIdx))
			args = append(args, *contact.JobTitle)
			argIdx++
		}
		if contact.PhoneNo != nil {
			setClauses = append(setClauses, fmt.Sprintf("phone_no = $%d", argIdx))
			args = append(args, *contact.PhoneNo)
			argIdx++
		}
		if contact.WaNo != nil {
			setClauses = append(setClauses, fmt.Sprintf("wa_no = $%d", argIdx))
			args = append(args, *contact.WaNo)
			argIdx++
		}
		if contact.Email != nil {
			setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
			args = append(args, *contact.Email)
			argIdx++
		}
		if contact.IdentityNo != nil {
			setClauses = append(setClauses, fmt.Sprintf("identity_no = $%d", argIdx))
			args = append(args, *contact.IdentityNo)
			argIdx++
		}

		if len(setClauses) == 0 {
			return nil // No fields to update
		}

		query := fmt.Sprintf(`
			UPDATE mst.m_outlet_contact
			SET %s
			WHERE cust_id = $%d
			AND outlet_id = $%d
			AND outlet_contact_id = $%d
		`, strings.Join(setClauses, ", "), argIdx, argIdx+1, argIdx+2)

		args = append(args, custId, outletId, *contact.OutletContactId)

		_, err := r.Exec(query, args...)
		if err != nil {
			log.Println("OutletListRepository, UpdateContact, update error:", err.Error())
			return err
		}
	} else if contact.ContactName != nil {
		// UPSERT: Insert or update by primary key (cust_id, outlet_id, contact_name)
		jobTitle := ""
		if contact.JobTitle != nil {
			jobTitle = *contact.JobTitle
		}
		phoneNo := ""
		if contact.PhoneNo != nil {
			phoneNo = *contact.PhoneNo
		}
		waNo := ""
		if contact.WaNo != nil {
			waNo = *contact.WaNo
		}
		email := ""
		if contact.Email != nil {
			email = *contact.Email
		}
		identityNo := ""
		if contact.IdentityNo != nil {
			identityNo = *contact.IdentityNo
		}

		query := `
			INSERT INTO mst.m_outlet_contact (
				cust_id, outlet_id, contact_name, job_title, phone_no, wa_no, email, identity_no
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (cust_id, outlet_id, contact_name) DO UPDATE SET
				job_title = EXCLUDED.job_title,
				phone_no = EXCLUDED.phone_no,
				wa_no = EXCLUDED.wa_no,
				email = EXCLUDED.email,
				identity_no = EXCLUDED.identity_no
		`

		_, err := r.Exec(query, custId, outletId, *contact.ContactName, jobTitle, phoneNo, waNo, email, identityNo)
		if err != nil {
			log.Println("OutletListRepository, UpdateContact, upsert error:", err.Error())
			return err
		}
	}

	return nil
}
