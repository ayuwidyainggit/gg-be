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

type BankRepository interface {
	FindOneByBankIdAndCustId(bankId int, custId string) (model.Bank, error)
	FindOneByBankCodeAndCustId(bankCode string, custId string) (model.Bank, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (bank []model.Bank, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (bank []model.Bank, total int, lastPage int, err error)
	Store(bank model.Bank) (int, error)
	Update(bankId int, request entity.UpdateBankRequest) error
	Delete(custId string, bankId int, deletedBy int64) error

	FindDistictOutletBankByCustIdLookupMode(dataFilter entity.QueryFilterOutletBank, custId string) ([]model.BankLookup, int, int, error)
	FindBankOutletByBankIdAndCustId(dataFilter entity.QueryFilterOutletBank, custId string) ([]model.OutletBankList, int, int, error)
}

func NewBankRepository(db *sqlx.DB) BankRepository {
	return &bankRepositoryImpl{db}
}

type bankRepositoryImpl struct {
	*sqlx.DB
}

func (repository *bankRepositoryImpl) FindOneByBankIdAndCustId(bankId int, custId string) (model.Bank, error) {
	bank := model.Bank{}
	query := `SELECT 
				cust_id, bank_id, bank_code,
				bank_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_bank 
			  WHERE bank_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&bank, query, bankId, custId)
	if err != nil {
		log.Println("bankRepository, FindOneByBankCodeAndCustId, err:", err.Error())
		return bank, err
	}

	return bank, nil
}

func (repository *bankRepositoryImpl) FindOneByBankCodeAndCustId(bankCode string, custId string) (model.Bank, error) {
	bank := model.Bank{}
	query := `SELECT 
				cust_id, bank_id, bank_code,
				bank_name, is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_bank 
			  WHERE bank_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&bank, query, bankCode, custId)
	if err != nil {
		log.Println("bankRepository, FindOneByBankCodeAndCustId, err:", err.Error())
		return bank, err
	}

	return bank, nil
}

func (repository *bankRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Bank, int, int, error) {

	banks := []model.Bank{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.bank_id, a.bank_code,
					a.bank_name, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.bank_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.bank_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_bank a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("bankRepository, count total, err:", err.Error())
		return banks, 0, 0, err
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
		sortBy := `a.bank_id`
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

	// log.Println("bankRepository, querySelect:", querySelect)
	err = repository.Select(&banks, querySelect)
	if err != nil {
		log.Println("bankRepository, FindAllByCustId, err:", err.Error())
		return banks, total, lastPage, err
	}

	return banks, total, lastPage, nil
}

func (repository *bankRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Bank, int, int, error) {

	banks := []model.Bank{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.bank_id, a.bank_code,
					a.bank_name, a.is_active `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.bank_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.bank_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_bank a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("bankRepository, count total, err:", err.Error())
		return banks, 0, 0, err
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
		sortBy := `a.bank_id`
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

	// log.Println("bankRepository, querySelect:", querySelect)
	err = repository.Select(&banks, querySelect)
	if err != nil {
		log.Println("bankRepository, FindAllByCustId, err:", err.Error())
		return banks, total, lastPage, err
	}

	return banks, total, lastPage, nil
}

func (repository *bankRepositoryImpl) Store(bank model.Bank) (int, error) {
	query :=
		`INSERT INTO mst.m_bank(
			cust_id, bank_code, bank_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11
		) RETURNING bank_id;`
	lastInsertId := bank.BankId
	err := repository.QueryRow(query,
		bank.CustId, bank.BankCode, bank.BankName,
		bank.IsActive, bank.CreatedBy, bank.CreatedAt, bank.UpdatedBy,
		bank.UpdatedAt, bank.IsDel, bank.DeletedBy, bank.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("bankRepository, Store, err:", err.Error())
		return bank.BankId, err
	}
	return bank.BankId, nil
}

func (repository *bankRepositoryImpl) Update(bankId int, request entity.UpdateBankRequest) error {
	var (
		r            model.BankUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("bankRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_bank
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND bank_id = :bank_id_old;`

	log.Println("bankRepository, Update, query:", query)

	sqlPatch.Args["bank_id_old"] = bankId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("bankRepository, Update, err:", err.Error())
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

func (repository *bankRepositoryImpl) Delete(custId string, bankId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_bank
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND bank_id = :bank_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"bank_id":    bankId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("BankRepository, Delete, err:", err.Error())
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

func (repository *bankRepositoryImpl) FindDistictOutletBankByCustIdLookupMode(dataFilter entity.QueryFilterOutletBank, custId string) ([]model.BankLookup, int, int, error) {

	banks := []model.BankLookup{}
	selectCount := ` COUNT(distinct(mob.bank_id)) AS total `
	selectField := ` distinct(mob.bank_id),mb.bank_name,mb.bank_code `
	qWhere := ` WHERE mob.cust_id = '` + custId + `' 
				AND (mob.bank_id = 0 
					OR (mb.bank_id IS NOT NULL AND COALESCE(mb.is_del, false) = false AND mb.is_active = true AND mb.cust_id = '` + custId + `')) `

	if dataFilter.Query != "" {
		qWhere += ` AND (mb.bank_code ILIKE '%` + dataFilter.Query + `%' 
					OR mb.bank_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.OutletID != nil && *dataFilter.OutletID != 0 {
		qWhere += fmt.Sprintf(` AND mob.outlet_id = %d`, *dataFilter.OutletID)
	}

	qFrom := ` 	FROM mst.m_outlet_bank mob 
				LEFT JOIN mst.m_bank mb ON mb.bank_id = mob.bank_id   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletBankRepository, count total, err:", err.Error())
		return banks, 0, 0, err
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
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		sortBy := `mb.bank_name`
		querySelect += fmt.Sprintf(` ORDER BY %s ASC`, sortBy)
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

	// log.Println("outletBankRepository, querySelect:", querySelect)
	err = repository.Select(&banks, querySelect)
	if err != nil {
		log.Println("outletBankRepository, FindDistictOutletBankByCustIdLookupMode, err:", err.Error())
		return banks, total, lastPage, err
	}

	return banks, total, lastPage, nil
}

func (repository *bankRepositoryImpl) FindBankOutletByBankIdAndCustId(dataFilter entity.QueryFilterOutletBank, custId string) ([]model.OutletBankList, int, int, error) {

	banks := []model.OutletBankList{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` mob.bank_id,mob.account_no,mob.account_name,mob.outlet_bank_id `
	qWhere := ` WHERE mob.cust_id = '` + custId + `' 
				AND (mob.bank_id = 0 
					OR (mb.bank_id IS NOT NULL AND COALESCE(mb.is_del, false) = false AND mb.is_active = true AND mb.cust_id = '` + custId + `')) `

	if dataFilter.Query != "" {
		qWhere += ` AND (mob.account_no ILIKE '%` + dataFilter.Query + `%' 
					OR mob.account_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.BankID != nil && *dataFilter.BankID != 0 {
		qWhere += fmt.Sprintf(` AND mob.bank_id = %d`, *dataFilter.BankID)
	}

	if dataFilter.OutletID != nil && *dataFilter.OutletID != 0 {
		qWhere += fmt.Sprintf(` AND mob.outlet_id = %d`, *dataFilter.OutletID)
	}

	qFrom := ` FROM mst.m_outlet_bank mob 
				LEFT JOIN mst.m_bank mb ON mb.bank_id = mob.bank_id `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("outletBankRepository, count total, err:", err.Error())
		return banks, 0, 0, err
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
		querySelect += fmt.Sprintf(` ORDER BY %s`, sortBy)
	} else {
		sortBy := `mb.bank_name`
		querySelect += fmt.Sprintf(` ORDER BY %s ASC`, sortBy)
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

	// log.Println("outletBankRepository, querySelect:", querySelect)
	err = repository.Select(&banks, querySelect)
	if err != nil {
		log.Println("outletBankRepository, FindBankOutletByBankIdAndCustId, err:", err.Error())
		return banks, total, lastPage, err
	}

	return banks, total, lastPage, nil
}
