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

type InvoiceDiscRepository interface {
	FindOneByInvDiscIdAndCustId(InvDiscId int, custId string) (model.InvoiceDisc, error)
	FindOneByInvDiscCodeAndCustId(InvDiscCode, custId string) (model.InvoiceDisc, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.InvoiceDisc, total int, lastPage int, err error)
	Store(invoiceDisc model.InvoiceDisc) (int, error)
	Update(InvDiscId int, request entity.UpdateInvoiceDiscRequest) error
	Delete(custId string, invDiscId int, deletedBy int64) error
}

func NewInvoiceDiscRepository(db *sqlx.DB) InvoiceDiscRepository {
	return &InvoiceDiscRepositoryImpl{db}
}

type InvoiceDiscRepositoryImpl struct {
	*sqlx.DB
}

func (repository *InvoiceDiscRepositoryImpl) FindOneByInvDiscIdAndCustId(InvDiscId int, custId string) (model.InvoiceDisc, error) {
	InvoiceDisc := model.InvoiceDisc{}
	query := `SELECT 
				cust_id, inv_disc_id, inv_disc_code,
				inv_disc_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_invoice_disc 
			  WHERE is_del = false 
			  AND inv_disc_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&InvoiceDisc, query, InvDiscId, custId)
	if err != nil {
		log.Println("InvoiceDiscRepository, FindOneByInvDiscIdAndCustId, err:", err.Error())
		return InvoiceDisc, err
	}

	return InvoiceDisc, nil
}

func (repository *InvoiceDiscRepositoryImpl) FindOneByInvDiscCodeAndCustId(invDiscCode, custId string) (model.InvoiceDisc, error) {
	invoiceDisc := model.InvoiceDisc{}
	query := `SELECT 
				cust_id, inv_disc_id, inv_disc_code,
				inv_disc_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_invoice_disc 
			  WHERE is_del = false 
			  AND inv_disc_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&invoiceDisc, query, invDiscCode, custId)
	if err != nil {
		log.Println("InvoiceDiscRepository, FindOneByInvDiscCode, err:", err.Error())
		return invoiceDisc, err
	}

	return invoiceDisc, nil
}

func (repository *InvoiceDiscRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.InvoiceDisc, int, int, error) {

	invoiceDiscs := []model.InvoiceDisc{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` id.cust_id, id.inv_disc_id, id.inv_disc_code,
	id.inv_disc_name, id.is_active, id.created_by,
	id.created_at, id.updated_by, id.updated_at,
	id.is_del, id.deleted_by, id.deleted_at,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE id.is_del = false 
				AND id.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (id.inv_disc_code ILIKE '%` + dataFilter.Query + `%' 
					OR id.inv_disc_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND id.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND id.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_invoice_disc id
	LEFT JOIN sys.m_user u ON u.user_id = id.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("InvoiceDiscRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("InvoiceDiscRepository, count total, err:", err.Error())
		return invoiceDiscs, 0, 0, err
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
		sortBy := `inv_disc_id`
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

	// log.Println("InvoiceDiscRepository, querySelect:", querySelect)
	err = repository.Select(&invoiceDiscs, querySelect)
	if err != nil {
		log.Println("InvoiceDiscRepository, FindAllByCustId, err:", err.Error())
		return invoiceDiscs, total, lastPage, err
	}

	// log.Println("lastPage:", lastPage)
	return invoiceDiscs, total, lastPage, nil
}

func (repository *InvoiceDiscRepositoryImpl) Store(invoiceDisc model.InvoiceDisc) (int, error) {
	query :=
		`INSERT INTO mst.m_invoice_disc(
			cust_id, inv_disc_code, inv_disc_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING inv_disc_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		invoiceDisc.CustId, invoiceDisc.InvDiscCode, invoiceDisc.InvDiscName, invoiceDisc.IsActive,
		invoiceDisc.CreatedBy, invoiceDisc.CreatedAt, invoiceDisc.UpdatedBy, invoiceDisc.UpdatedAt,
		invoiceDisc.IsDel, invoiceDisc.DeletedBy, invoiceDisc.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("InvoiceDiscRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *InvoiceDiscRepositoryImpl) Update(invDiscId int, request entity.UpdateInvoiceDiscRequest) error {
	var (
		r            model.InvoiceDiscUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("InvoiceDiscRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_invoice_disc
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND inv_disc_id = :inv_disc_id;`

	// log.Println("InvoiceDiscRepository, Update, query:", query)

	sqlPatch.Args["inv_disc_id"] = invDiscId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("InvoiceDiscRepository, Update, err:", err.Error())
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

func (repository *InvoiceDiscRepositoryImpl) Delete(custId string, invDiscId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_invoice_disc
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND inv_disc_id = :inv_disc_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"inv_disc_id": invDiscId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("InvoiceDiscRepository, Delete, err:", err.Error())
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
