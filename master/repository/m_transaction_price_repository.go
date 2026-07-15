package repository

import (
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"

	"github.com/jmoiron/sqlx"
)

type MTransactionPriceRepository interface {
	Store(mTransPrice *model.MTransactionPrice) error
	StoreBatch(transPrices []model.MTransactionPrice) error
	GetByCustPro(custID string, proID, distributorID int64, coverage string) (*model.MTransactionPrice, error)
	Delete(custID string, transPriceID string) error
	FindAllByCustID(filter entity.MTransactionPriceQueryFilter) ([]model.MTransactionPrice, int, int, error)
	DeleteByStartDate(custID string, proID int64, effectiveDate time.Time) error
	UpdateBatch(transPrices []model.MTransactionPrice) error
	StoreIfNotExists(mTransPrice *model.MTransactionPrice) error
}

func NewMTransactionPriceRepository(db *sqlx.DB) MTransactionPriceRepository {
	return &MTransactionPriceRepositoryImpl{db}
}

type MTransactionPriceRepositoryImpl struct {
	*sqlx.DB
}

func (repository *MTransactionPriceRepositoryImpl) Store(transPrice *model.MTransactionPrice) error {
	query :=
		`INSERT INTO mst.m_transaction_price(` +
			`cust_id, transaction_price_id, pro_id, ` +
			`purch_price1, purch_price2, purch_price3,` +
			`sell_price1, sell_price2, sell_price3,` +
			`source, created_by, created_at, ` +
			`start_date, end_date,` +
			`coverage, distributor_id, price_group_reff,` +
			`reference_id,` +
			`outlet_id)` +
			`VALUES ( ` +
			`:cust_id, :transaction_price_id, :pro_id, ` +
			`:purch_price1, :purch_price2, :purch_price3,` +
			`:sell_price1, :sell_price2, :sell_price3,` +
			`:source, :created_by, :created_at, ` +
			`:start_date, :end_date,` +
			`:coverage, :distributor_id, :price_group_reff,` +
			`:reference_id, ` +
			`:outlet_id)` +
			`RETURNING transaction_price_id;`

	_, err := repository.NamedExec(query, map[string]interface{}{
		"cust_id":              transPrice.CustID,
		"transaction_price_id": transPrice.TransactionPriceID,
		"pro_id":               transPrice.ProID,
		"purch_price1":         transPrice.PurchPrice1,
		"purch_price2":         transPrice.PurchPrice2,
		"purch_price3":         transPrice.PurchPrice3,
		"sell_price1":          transPrice.SellPrice1,
		"sell_price2":          transPrice.SellPrice2,
		"sell_price3":          transPrice.SellPrice3,
		"source":               transPrice.Source,
		"created_by":           transPrice.CreatedBy,
		"created_at":           transPrice.CreatedAt,
		"start_date":           transPrice.StartDate,
		"end_date":             transPrice.EndDate,
		"coverage":             transPrice.Coverage,
		"distributor_id":       transPrice.DistributorID,
		"price_group_reff":     transPrice.PriceGroupReff,
		"reference_id":         transPrice.ReferenceID,
		"outlet_id":            transPrice.OutletID,
	})
	if err != nil {
		log.Error("MTransactionPriceRepository, Store, err:", err.Error())
		return err
	}
	// log.Println("Store, result:", result)

	return nil
}

func (repository *MTransactionPriceRepositoryImpl) StoreBatch(transPrices []model.MTransactionPrice) error {
	query := `INSERT INTO mst.m_transaction_price(
		cust_id, transaction_price_id, pro_id, 
		purch_price1, purch_price2, purch_price3,
		sell_price1, sell_price2, sell_price3,
		source, created_by, created_at, 
		start_date, end_date,
		coverage, distributor_id, price_group_reff,
		reference_id)
		VALUES (
		:cust_id, :transaction_price_id, :pro_id, 
		:purch_price1, :purch_price2, :purch_price3,
		:sell_price1, :sell_price2, :sell_price3,
		:source, :created_by, :created_at, 
		:start_date, :end_date,
		:coverage, :distributor_id, :price_group_reff,
		:reference_id) RETURNING transaction_price_id;`

	var batchData []map[string]interface{}
	for _, transPrice := range transPrices {
		batchData = append(batchData, map[string]interface{}{
			"cust_id":              transPrice.CustID,
			"transaction_price_id": transPrice.TransactionPriceID,
			"pro_id":               transPrice.ProID,
			"purch_price1":         transPrice.PurchPrice1,
			"purch_price2":         transPrice.PurchPrice2,
			"purch_price3":         transPrice.PurchPrice3,
			"sell_price1":          transPrice.SellPrice1,
			"sell_price2":          transPrice.SellPrice2,
			"sell_price3":          transPrice.SellPrice3,
			"source":               transPrice.Source,
			"created_by":           transPrice.CreatedBy,
			"created_at":           transPrice.CreatedAt,
			"start_date":           transPrice.StartDate,
			"end_date":             transPrice.EndDate,
			"coverage":             transPrice.Coverage,
			"distributor_id":       transPrice.DistributorID,
			"price_group_reff":     transPrice.PriceGroupReff,
			"reference_id":         transPrice.ReferenceID,
		})
	}

	_, err := repository.NamedExec(query, batchData)
	if err != nil {
		log.Error("MTransactionPriceRepository, StoreBatch, err:", err.Error())
		return err
	}

	return nil
}

func (repository *MTransactionPriceRepositoryImpl) UpdateBatch(transPrices []model.MTransactionPrice) error {
	query := `UPDATE mst.m_transaction_price SET
		purch_price1 = :purch_price1,
		purch_price2 = :purch_price2,
		purch_price3 = :purch_price3,
		sell_price1 = :sell_price1,
		sell_price2 = :sell_price2,
		sell_price3 = :sell_price3,
		source = :source,
		created_by = :created_by,
		created_at = :created_at,
		start_date = :start_date,
		end_date = :end_date,
		coverage = :coverage,
		distributor_id = :distributor_id,
		price_group_reff = :price_group_reff,
		reference_id = :reference_id
	WHERE transaction_price_id = :transaction_price_id;`

	var batchData []map[string]interface{}
	for _, transPrice := range transPrices {
		batchData = append(batchData, map[string]interface{}{
			"transaction_price_id": transPrice.TransactionPriceID,
			"purch_price1":         transPrice.PurchPrice1,
			"purch_price2":         transPrice.PurchPrice2,
			"purch_price3":         transPrice.PurchPrice3,
			"sell_price1":          transPrice.SellPrice1,
			"sell_price2":          transPrice.SellPrice2,
			"sell_price3":          transPrice.SellPrice3,
			"source":               transPrice.Source,
			"created_by":           transPrice.CreatedBy,
			"created_at":           transPrice.CreatedAt,
			"start_date":           transPrice.StartDate,
			"end_date":             transPrice.EndDate,
			"coverage":             transPrice.Coverage,
			"distributor_id":       transPrice.DistributorID,
			"price_group_reff":     transPrice.PriceGroupReff,
			"reference_id":         transPrice.ReferenceID,
		})
	}

	_, err := repository.NamedExec(query, batchData)
	if err != nil {
		log.Error("MTransactionPriceRepository, UpdateBatch, err:", err.Error())
		return err
	}

	return nil
}

func (repository *MTransactionPriceRepositoryImpl) Delete(custID, transPriceID string) error {
	var nRows int64
	query := `DELETE mst.m_transaction_price WHERE cust_id = :cust_id AND transaction_price_id = :transaction_price_id;`

	wMap := map[string]interface{}{
		"cust_id":              custID,
		"transaction_price_id": transPriceID,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Error("MTransactionPriceRepository, Delete, err:", err.Error())
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

func (repository *MTransactionPriceRepositoryImpl) GetByCustPro(custID string, proID, distributorID int64, coverage string) (*model.MTransactionPrice, error) {
	query := `SELECT * FROM mst.m_transaction_price
	WHERE cust_id = $1 AND pro_id = $2 AND distributor_id = $3 AND coverage = $4 LIMIT 1;`

	var transPrice model.MTransactionPrice

	err := repository.Get(&transPrice, query, custID, proID, distributorID, coverage)
	if err != nil {
		log.Error("MTransactionPriceRepository, GetByCustPro, err:", err.Error())
		return nil, err
	}

	return &transPrice, nil
}

func (repository *MTransactionPriceRepositoryImpl) FindAllByCustID(filter entity.MTransactionPriceQueryFilter) ([]model.MTransactionPrice, int, int, error) {

	transPrices := []model.MTransactionPrice{}
	selectCount := ` COUNT(prc.*) AS total `
	selectField := `prc.cust_id, prc.price_id, prc.status,
					prc.pro_id, prd.pro_code, prd.pro_name,
					un1.unit_name AS unit_name1,
					un2.unit_name AS unit_name2,
					un3.unit_name AS unit_name3,
					prc.coverage, prc.effective_date,
					prc.unit_id1, prc.unit_id2, prc.unit_id3,
					prc.conv_unit2, prc.conv_unit3,
					prc.purch_price1, prc.purch_price2, prc.purch_price3,
					prc.sell_price1, prc.sell_price2, prc.sell_price3,
					prc.new_purch_price1, prc.new_purch_price2, prc.new_purch_price3,
					prc.new_sell_price1, prc.new_sell_price2, prc.new_sell_price3,
					prc.created_by, prc.created_at, prc.updated_by, prc.updated_at `
	qWhere := ` WHERE prc.cust_id = '` + filter.CustID + `' `

	// if filter.Query != "" {
	// 	qWhere += ` AND (prd.price_name ILIKE '%` + dataFilter.Query + `%')`
	// }

	// if dataFilter.PriceId != 0 {
	// 	qWhere += `AND prc.price_id = ` + strconv.Itoa(dataFilter.PriceId) + ` `
	// }

	qFrom := ` FROM mst.m_transaction_price tp
	LEFT JOIN mst.m_product prd ON prd.pro_id = tp.pro_id
	LEFT JOIN mst.m_unit un1 ON prd.unit_id1 = un1.unit_id
	LEFT JOIN mst.m_unit un2 ON prd.unit_id2 = un2.unit_id
	LEFT JOIN mst.m_unit un3 ON prd.unit_id3 = un3.unit_id `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("MTransactionPriceRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("MTransactionPriceRepository, count total, err:", err.Error())
		return transPrices, 0, 0, err
	}

	sortBy := `` // default sort by
	if filter.Sort != "" {
		mSortBy := strings.Split(filter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `prc.price_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	if filter.Limit == 0 {
		filter.Limit = 10
	}

	page := filter.Page
	if page-1 < 1 {
		page = 1
	}
	offset := (page - 1) * filter.Limit

	lastPage := int(math.Ceil(float64(float64(total) / float64(filter.Limit))))

	querySelect += fmt.Sprintf(` LIMIT %s OFFSET %s`, strconv.Itoa(filter.Limit), strconv.Itoa(offset))

	// log.Println("MTransactionPriceRepository, querySelect:", querySelect)
	err = repository.Select(&transPrices, querySelect)
	if err != nil {
		log.Error("MTransactionPriceRepository, FindAllByCustId, err:", err.Error())
		return transPrices, total, lastPage, err
	}

	return transPrices, total, lastPage, nil
}

func (repository *MTransactionPriceRepositoryImpl) DeleteByStartDate(custID string, proID int64, effectiveDate time.Time) error {
	var nRows int64
	query := `DELETE FROM mst.m_transaction_price tp
			  WHERE cust_id = :cust_id
			  AND pro_id = :pro_id AND start_date >= :start_date`

	wMap := map[string]interface{}{
		"cust_id":    custID,
		"pro_id":     proID,
		"start_date": effectiveDate,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Error("InvoiceDiscDetRepository, Delete, err:", err.Error())
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

func (repository *MTransactionPriceRepositoryImpl) StoreIfNotExists(transPrice *model.MTransactionPrice) error {
	// First, check if a record already exists with the same key fields
	queryCount := `SELECT COUNT(*) AS total FROM mst.m_transaction_price 
		WHERE cust_id = $1 
		AND pro_id = $2 
		AND source = $3 
		AND coverage = $4 
		AND start_date = $5
		AND reference_id = $6`
	var total int
	err := repository.QueryRow(queryCount, 
		transPrice.CustID, transPrice.ProID, 
		transPrice.Source, transPrice.Coverage, 
		transPrice.StartDate, transPrice.ReferenceID).
		Scan(&total)
	if err != nil {
		log.Error("MTransactionPriceRepository, StoreIfNotExists, check existing, err:", err.Error())
		return err
	}

	if total > 0 {
		log.Error("MTransactionPriceRepository, StoreIfNotExists, record already exists, skipping insert")
		return nil
	}

	// Record doesn't exist, proceed with insert
	query :=
		`INSERT INTO mst.m_transaction_price(` +
			`cust_id, transaction_price_id, pro_id, ` +
			`purch_price1, purch_price2, purch_price3,` +
			`sell_price1, sell_price2, sell_price3,` +
			`source, created_by, created_at, ` +
			`start_date, end_date,` +
			`coverage, distributor_id, price_group_reff,` +
			`reference_id,` +
			`outlet_id)` +
			`VALUES ( ` +
			`:cust_id, :transaction_price_id, :pro_id, ` +
			`:purch_price1, :purch_price2, :purch_price3,` +
			`:sell_price1, :sell_price2, :sell_price3,` +
			`:source, :created_by, :created_at, ` +
			`:start_date, :end_date,` +
			`:coverage, :distributor_id, :price_group_reff,` +
			`:reference_id, ` +
			`:outlet_id)` +
			`RETURNING transaction_price_id;`

	_, err = repository.NamedExec(query, map[string]interface{}{
		"cust_id":              transPrice.CustID,
		"transaction_price_id": transPrice.TransactionPriceID,
		"pro_id":               transPrice.ProID,
		"purch_price1":         transPrice.PurchPrice1,
		"purch_price2":         transPrice.PurchPrice2,
		"purch_price3":         transPrice.PurchPrice3,
		"sell_price1":          transPrice.SellPrice1,
		"sell_price2":          transPrice.SellPrice2,
		"sell_price3":          transPrice.SellPrice3,
		"source":               transPrice.Source,
		"created_by":           transPrice.CreatedBy,
		"created_at":           transPrice.CreatedAt,
		"start_date":           transPrice.StartDate,
		"end_date":             transPrice.EndDate,
		"coverage":             transPrice.Coverage,
		"distributor_id":       transPrice.DistributorID,
		"price_group_reff":     transPrice.PriceGroupReff,
		"reference_id":         transPrice.ReferenceID,
		"outlet_id":            transPrice.OutletID,
	})

	if err != nil {
		log.Error("MTransactionPriceRepository, StoreIfNotExists, insert, err:", err.Error())
		return err
	}

	log.Info("MTransactionPriceRepository, StoreIfNotExists, record inserted successfully")
	return nil
}
