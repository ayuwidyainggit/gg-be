package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type MPriceRepository interface {
	FindOneByMPriceIDAndCustID(entity.DetailMPriceParams) (model.MPriceDetail, error)
	FindAllByCustID(dataFilter entity.MPriceQueryFilter, custId string) (mPrice []model.MPrice, total int, lastPage int, err error)
	Store(mprice *model.MPrice) error
	Update(priceID string, request entity.UpdateMPriceRequest) error
	Delete(custId string, priceId string, deletedBy int64) error
	Cancel(entity.CancelMPriceParams) error
	FindOneProductByProID(proID int64, custID string) (model.Product, error)
	FindOneProductSnapshotByProID(proID int64, custID string) (model.MPriceProductSnapshot, error)
	FindOneProductSnapshotByCode(proCode string, custID string) (model.MPriceProductSnapshot, error)
	FindAffectedDistributorProductIDs(price model.MPriceDetail, parentCustID string) ([]int64, error)
	FindUpdatedDistributorProductIDs(price model.MPriceDetail, parentCustID string, distributorIDs []int64) ([]int64, error)
	UpdatePrincipalAssignedProductPrices(parentProID int64, distributorIDs []int64, detail model.MPriceDetail) error
	UpdateDistributorProductPrices(custID string, distributorID, proID int64, detail model.MPriceDetail) error
	PublishByRMQ(request entity.PublishByRmqMPriceReq) error
	FindBrokenDistributorChildLinks(parentProID int64, parentCustID string, distributorIDs []int64) ([]int64, error)
}

func NewMPriceRepository(db *sqlx.DB) MPriceRepository {
	return &MPriceRepositoryImpl{db}
}

type MPriceRepositoryImpl struct {
	*sqlx.DB
}

func (repository *MPriceRepositoryImpl) FindOneByMPriceIDAndCustID(detail entity.DetailMPriceParams) (model.MPriceDetail, error) {
	price := model.MPriceDetail{}
	query := `
		SELECT
			prc.cust_id,
			prc.price_id,
			prc.coverage,
			prc.effective_date,
			prc.pro_id,
			prc.unit_id1,
			prc.unit_id2,
			prc.unit_id3,
			prc.conv_unit2,
			prc.conv_unit3,
			prc.purch_price1,
			prc.purch_price2,
			prc.purch_price3,
			prc.sell_price1,
			prc.sell_price2,
			prc.sell_price3,
			prc.new_purch_price1,
			prc.new_purch_price2,
			prc.new_purch_price3,
			prc.new_sell_price1,
			prc.new_sell_price2,
			prc.new_sell_price3,
			prc.status,
			prc.created_by_id,
			prc.created_by,
			prc.created_at,
			prc.updated_by_id,
			prc.updated_by,
			prc.updated_at,
			prc.distributor_ids,
			prd.pro_code,
			prd.pro_name,
			COALESCE(un1.unit_name, '') AS unit_name1,
			COALESCE(un2.unit_name, '') AS unit_name2,
			COALESCE(un3.unit_name, '') AS unit_name3
		FROM mst.m_price prc
		LEFT JOIN mst.m_product prd ON prd.pro_id = prc.pro_id AND prd.cust_id = prc.cust_id
		LEFT JOIN mst.m_unit un1 ON un1.unit_id = prd.unit_id1 AND un1.cust_id = prd.cust_id
		LEFT JOIN mst.m_unit un2 ON un2.unit_id = prd.unit_id2 AND un2.cust_id = prd.cust_id
		LEFT JOIN mst.m_unit un3 ON un3.unit_id = prd.unit_id3 AND un3.cust_id = prd.cust_id
		WHERE prc.price_id = $1
		  AND prc.cust_id = $2
	`

	if err := repository.Get(&price, query, detail.PriceID, detail.CustID); err != nil {
		log.Error("priceRepository, FindOneByMPriceIDAndCustID, err:", err.Error())
		return price, err
	}

	return price, nil
}

func (repository *MPriceRepositoryImpl) FindAllByCustID(dataFilter entity.MPriceQueryFilter, custID string) ([]model.MPrice, int, int, error) {
	mPrices := []model.MPrice{}
	selectCount := `COUNT(prc.price_id) AS total`
	selectField := `
		prc.cust_id,
		prc.price_id,
		prc.status,
		prc.pro_id,
		prd.pro_code,
		prd.pro_name,
		un1.unit_name AS unit_name1,
		un2.unit_name AS unit_name2,
		un3.unit_name AS unit_name3,
		prc.coverage,
		prc.effective_date,
		prc.unit_id1,
		prc.unit_id2,
		prc.unit_id3,
		prc.conv_unit2,
		prc.conv_unit3,
		prc.purch_price1,
		prc.purch_price2,
		prc.purch_price3,
		prc.sell_price1,
		prc.sell_price2,
		prc.sell_price3,
		prc.new_purch_price1,
		prc.new_purch_price2,
		prc.new_purch_price3,
		prc.new_sell_price1,
		prc.new_sell_price2,
		prc.new_sell_price3,
		prc.created_by_id,
		prc.created_by,
		prc.created_at,
		prc.updated_by_id,
		prc.updated_by,
		prc.updated_at,
		prc.distributor_ids
	`
	qFrom := `
		FROM mst.m_price prc
		LEFT JOIN mst.m_product prd ON prd.pro_id = prc.pro_id AND prd.cust_id = prc.cust_id
		LEFT JOIN mst.m_unit un1 ON un1.unit_id = prd.unit_id1 AND un1.cust_id = prd.cust_id
		LEFT JOIN mst.m_unit un2 ON un2.unit_id = prd.unit_id2 AND un2.cust_id = prd.cust_id
		LEFT JOIN mst.m_unit un3 ON un3.unit_id = prd.unit_id3 AND un3.cust_id = prd.cust_id
	`
	escapedCustID := strings.ReplaceAll(custID, "'", "''")
	conditions := []string{
		fmt.Sprintf("prc.cust_id = '%s'", escapedCustID),
		"prc.is_del IS FALSE",
	}

	if dataFilter.EffectiveDateStart != "" || dataFilter.EffectiveDateEnd != "" {
		start := dataFilter.EffectiveDateStart
		end := dataFilter.EffectiveDateEnd
		if start == "" {
			start = end
		}
		if end == "" {
			end = start
		}
		conditions = append(conditions, fmt.Sprintf(
			"prc.effective_date::date BETWEEN '%s' AND '%s'",
			start,
			end,
		))
	}

	if dataFilter.Query != "" {
		escaped := strings.ReplaceAll(dataFilter.Query, "'", "''")
		conditions = append(conditions, fmt.Sprintf(
			"(prd.pro_name ILIKE '%%%s%%' OR prd.pro_code ILIKE '%%%s%%')",
			escaped,
			escaped,
		))
	}

	if len(dataFilter.Status) > 0 {
		conditions = append(conditions, fmt.Sprintf("prc.status IN (%s)", joinInts(dataFilter.Status)))
	}

	if len(dataFilter.DistributorIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf(
			"prc.distributor_ids && ARRAY[%s]::bigint[]",
			joinInt64s(dataFilter.DistributorIDs),
		))
	}

	qWhere := "WHERE " + strings.Join(conditions, " AND ")

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	if err := repository.QueryRow(queryCount).Scan(&total); err != nil {
		log.Error("priceRepository, FindAllByCustID, count err:", err.Error())
		return mPrices, 0, 0, err
	}

	sortBy := "prc.created_at DESC"
	if dataFilter.Sort != "" {
		sortBy = buildSortClause(dataFilter.Sort)
		if sortBy == "" {
			sortBy = "prc.created_at DESC"
		}
	}
	querySelect += " ORDER BY " + sortBy

	if dataFilter.Limit == 0 {
		dataFilter.Limit = 5
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * dataFilter.Limit
	lastPage := int(math.Ceil(float64(total) / float64(dataFilter.Limit)))
	querySelect += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)

	if err := repository.Select(&mPrices, querySelect); err != nil {
		log.Error("priceRepository, FindAllByCustID, select err:", err.Error())
		return mPrices, total, lastPage, err
	}

	return mPrices, total, lastPage, nil
}

func (repository *MPriceRepositoryImpl) Store(price *model.MPrice) error {
	query := `
		INSERT INTO mst.m_price(
			cust_id,
			price_id,
			coverage,
			effective_date,
			pro_id,
			unit_id1,
			unit_id2,
			unit_id3,
			conv_unit2,
			conv_unit3,
			purch_price1,
			purch_price2,
			purch_price3,
			sell_price1,
			sell_price2,
			sell_price3,
			new_purch_price1,
			new_purch_price2,
			new_purch_price3,
			new_sell_price1,
			new_sell_price2,
			new_sell_price3,
			status,
			created_by_id,
			created_by,
			created_at,
			updated_by_id,
			updated_by,
			updated_at,
			distributor_ids
		) VALUES (
			:cust_id,
			:price_id,
			:coverage,
			:effective_date,
			:pro_id,
			:unit_id1,
			:unit_id2,
			:unit_id3,
			:conv_unit2,
			:conv_unit3,
			:purch_price1,
			:purch_price2,
			:purch_price3,
			:sell_price1,
			:sell_price2,
			:sell_price3,
			:new_purch_price1,
			:new_purch_price2,
			:new_purch_price3,
			:new_sell_price1,
			:new_sell_price2,
			:new_sell_price3,
			:status,
			:created_by_id,
			:created_by,
			:created_at,
			:updated_by_id,
			:updated_by,
			:updated_at,
			:distributor_ids
		)
	`

	args := map[string]interface{}{
		"cust_id":          price.CustID,
		"price_id":         price.PriceID,
		"coverage":         price.Coverage,
		"effective_date":   price.EffectiveDate,
		"pro_id":           price.ProID,
		"unit_id1":         price.UnitID1,
		"unit_id2":         price.UnitID2,
		"unit_id3":         price.UnitID3,
		"conv_unit2":       price.ConvUnit2,
		"conv_unit3":       price.ConvUnit3,
		"purch_price1":     price.PurchPrice1,
		"purch_price2":     price.PurchPrice2,
		"purch_price3":     price.PurchPrice3,
		"sell_price1":      price.SellPrice1,
		"sell_price2":      price.SellPrice2,
		"sell_price3":      price.SellPrice3,
		"new_purch_price1": price.NewPurchPrice1,
		"new_purch_price2": price.NewPurchPrice2,
		"new_purch_price3": price.NewPurchPrice3,
		"new_sell_price1":  price.NewSellPrice1,
		"new_sell_price2":  price.NewSellPrice2,
		"new_sell_price3":  price.NewSellPrice3,
		"status":           price.Status,
		"created_by_id":    price.CreatedByID,
		"created_by":       price.CreatedBy,
		"created_at":       price.CreatedAt,
		"updated_by_id":    price.UpdatedByID,
		"updated_by":       price.UpdatedBy,
		"updated_at":       price.UpdatedAt,
		"distributor_ids":  pq.Array(price.DistributorIDs),
	}

	if _, err := repository.NamedExec(query, args); err != nil {
		log.Error("priceRepository, Store, err:", err.Error())
		return err
	}

	return nil
}

func (repository *MPriceRepositoryImpl) Update(priceID string, request entity.UpdateMPriceRequest) error {
	var (
		r            model.MPriceUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	if sqlSetFields == "" {
		return errors.New("no fields to update")
	}

	distIDsUpdateValue := ""
	if r.DistributorIDs == nil {
		distIDsUpdateValue = "distributor_ids = NULL, "
	} else {
		sqlPatch.Args["distributor_ids"] = pq.Array(*r.DistributorIDs)
	}

	query := `UPDATE mst.m_price SET ` + sqlSetFields + `, ` + distIDsUpdateValue +
		`updated_at = CURRENT_TIMESTAMP WHERE cust_id = :cust_id AND price_id = :price_id AND is_del IS FALSE`

	sqlPatch.Args["price_id"] = priceID
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error("priceRepository, Update, err:", err.Error())
		return err
	}

	if nRows, err = result.RowsAffected(); err != nil {
		log.Error("priceRepository, Update, RowsAffected err:", err.Error())
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *MPriceRepositoryImpl) Delete(custID string, priceID string, deletedBy int64) error {
	query := `
		UPDATE mst.m_price
		SET is_del = true,
			deleted_at = CURRENT_TIMESTAMP,
			deleted_by = :deleted_by
		WHERE is_del IS FALSE
		  AND cust_id = :cust_id
		  AND price_id = :price_id
	`

	result, err := repository.NamedExec(query, map[string]interface{}{
		"cust_id":    custID,
		"price_id":   priceID,
		"deleted_by": deletedBy,
	})
	if err != nil {
		log.Error("MPriceRepository, Delete, err:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		log.Error("MPriceRepository, Delete, RowsAffected err:", err.Error())
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (repository *MPriceRepositoryImpl) Cancel(detail entity.CancelMPriceParams) error {
	query := `
		UPDATE mst.m_price
		SET status = 5,
			updated_by_id = :updated_by_id,
			updated_by = :updated_by,
			updated_at = CURRENT_TIMESTAMP
		WHERE cust_id = :cust_id
		  AND price_id = :price_id
		  AND status = 1
		  AND is_del IS FALSE
	`

	result, err := repository.NamedExec(query, map[string]interface{}{
		"cust_id":       detail.CustID,
		"price_id":      detail.PriceID,
		"updated_by_id": detail.UpdatedByID,
		"updated_by":    detail.UpdatedBy,
	})
	if err != nil {
		log.Error("MPriceRepository, Cancel, err:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		log.Error("MPriceRepository, Cancel, RowsAffected err:", err.Error())
		return errors.New("no row cancelled")
	}
	if nRows == 0 {
		return errors.New("no row cancelled")
	}

	return nil
}

func (repository *MPriceRepositoryImpl) FindOneProductByProID(proID int64, custID string) (model.Product, error) {
	productID := model.Product{}
	query := `SELECT pro_id FROM mst.m_product WHERE pro_id = $1 AND cust_id = $2 AND is_del = false`
	if err := repository.Get(&productID, query, proID, custID); err != nil {
		log.Error("priceRepository, FindOneProductByProID, err:", err.Error())
		return productID, err
	}
	return productID, nil
}

func (repository *MPriceRepositoryImpl) FindOneProductSnapshotByProID(proID int64, custID string) (model.MPriceProductSnapshot, error) {
	snapshot := model.MPriceProductSnapshot{}
	query := `
		SELECT
			pro_id,
			pro_code,
			pro_name,
			unit_id1,
			unit_id2,
			unit_id3,
			conv_unit2,
			conv_unit3,
			purch_price1,
			purch_price2,
			purch_price3,
			sell_price1,
			sell_price2,
			sell_price3,
			distributor_id,
			COALESCE(parent_pro_id, 0) AS parent_pro_id
		FROM mst.m_product
		WHERE pro_id = $1
		  AND cust_id = $2
		  AND is_del = false
	`
	if err := repository.Get(&snapshot, query, proID, custID); err != nil {
		log.Error("priceRepository, FindOneProductSnapshotByProID, err:", err.Error())
		return snapshot, err
	}
	return snapshot, nil
}

func (repository *MPriceRepositoryImpl) FindOneProductSnapshotByCode(proCode string, custID string) (model.MPriceProductSnapshot, error) {
	snapshot := model.MPriceProductSnapshot{}
	query := `
		SELECT
			pro_id,
			pro_code,
			pro_name,
			unit_id1,
			unit_id2,
			unit_id3,
			conv_unit2,
			conv_unit3,
			purch_price1,
			purch_price2,
			purch_price3,
			sell_price1,
			sell_price2,
			sell_price3,
			distributor_id,
			COALESCE(parent_pro_id, 0) AS parent_pro_id
		FROM mst.m_product
		WHERE pro_code = $1
		  AND cust_id = $2
		  AND is_del = false
	`
	if err := repository.Get(&snapshot, query, proCode, custID); err != nil {
		log.Error("priceRepository, FindOneProductSnapshotByCode, err:", err.Error())
		return snapshot, err
	}
	return snapshot, nil
}

func (repository *MPriceRepositoryImpl) FindAffectedDistributorProductIDs(price model.MPriceDetail, parentCustID string) ([]int64, error) {
	distributorIDs := make([]int64, 0)
	query := `
		SELECT DISTINCT distributor_id
		FROM mst.m_product
		WHERE is_del IS FALSE
		  AND distributor_id IS NOT NULL
	`
	args := []interface{}{}

	if price.CustID == parentCustID {
		query += ` AND parent_pro_id = ?`
		args = append(args, price.ProID)
	} else {
		query += ` AND cust_id = ? AND pro_id = ?`
		args = append(args, price.CustID, price.ProID)
	}

	query = repository.Rebind(query)
	if err := repository.Select(&distributorIDs, query, args...); err != nil {
		log.Error("priceRepository, FindAffectedDistributorProductIDs, err:", err.Error())
		return distributorIDs, err
	}
	return distributorIDs, nil
}

func (repository *MPriceRepositoryImpl) FindUpdatedDistributorProductIDs(price model.MPriceDetail, parentCustID string, distributorIDs []int64) ([]int64, error) {
	updatedDistributorIDs := make([]int64, 0)
	if len(distributorIDs) == 0 {
		return updatedDistributorIDs, nil
	}

	query := `
		SELECT DISTINCT distributor_id
		FROM mst.m_product
		WHERE is_del IS FALSE
		  AND distributor_id IN (?)
		  AND purch_price1 = ?
		  AND purch_price2 = ?
		  AND purch_price3 = ?
		  AND sell_price1 = ?
		  AND sell_price2 = ?
		  AND sell_price3 = ?
	`
	args := []interface{}{
		distributorIDs,
		price.NewPurchPrice1,
		price.NewPurchPrice2,
		price.NewPurchPrice3,
		price.NewSellPrice1,
		price.NewSellPrice2,
		price.NewSellPrice3,
	}

	if price.CustID == parentCustID {
		query += ` AND parent_pro_id = ?`
		args = append(args, price.ProID)
	} else {
		query += ` AND cust_id = ? AND pro_id = ?`
		args = append(args, price.CustID, price.ProID)
	}

	query, argsExpanded, err := sqlx.In(query, args...)
	if err != nil {
		log.Error("priceRepository, FindUpdatedDistributorProductIDs, sqlx.In err:", err.Error())
		return updatedDistributorIDs, err
	}
	query = repository.Rebind(query)

	if err := repository.Select(&updatedDistributorIDs, query, argsExpanded...); err != nil {
		log.Error("priceRepository, FindUpdatedDistributorProductIDs, err:", err.Error())
		return updatedDistributorIDs, err
	}
	return updatedDistributorIDs, nil
}

func (repository *MPriceRepositoryImpl) UpdatePrincipalAssignedProductPrices(parentProID int64, distributorIDs []int64, detail model.MPriceDetail) error {
	query := `
		UPDATE mst.m_product
		SET purch_price1 = ?,
			purch_price2 = ?,
			purch_price3 = ?,
			sell_price1 = ?,
			sell_price2 = ?,
			sell_price3 = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE parent_pro_id = ?
		  AND is_del IS FALSE
		  AND distributor_id IS NOT NULL
	`
	args := []interface{}{
		detail.NewPurchPrice1,
		detail.NewPurchPrice2,
		detail.NewPurchPrice3,
		detail.NewSellPrice1,
		detail.NewSellPrice2,
		detail.NewSellPrice3,
		parentProID,
	}

	if len(distributorIDs) > 0 {
		query += ` AND distributor_id IN (?)`
		args = append(args, distributorIDs)
	}

	query, argsExpanded, err := sqlx.In(query, args...)
	if err != nil {
		log.Error("priceRepository, UpdatePrincipalAssignedProductPrices, sqlx.In err:", err.Error())
		return err
	}
	query = repository.Rebind(query)

	if _, err := repository.Exec(query, argsExpanded...); err != nil {
		log.Error("priceRepository, UpdatePrincipalAssignedProductPrices, err:", err.Error())
		return err
	}
	return nil
}

func (repository *MPriceRepositoryImpl) UpdateDistributorProductPrices(custID string, distributorID, proID int64, detail model.MPriceDetail) error {
	query := `
		UPDATE mst.m_product
		SET purch_price1 = $1,
			purch_price2 = $2,
			purch_price3 = $3,
			sell_price1 = $4,
			sell_price2 = $5,
			sell_price3 = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE cust_id = $7
		  AND pro_id = $8
		  AND is_del IS FALSE
	`
	args := []interface{}{
		detail.NewPurchPrice1,
		detail.NewPurchPrice2,
		detail.NewPurchPrice3,
		detail.NewSellPrice1,
		detail.NewSellPrice2,
		detail.NewSellPrice3,
		custID,
		proID,
	}

	if distributorID > 0 {
		query += ` AND distributor_id = $9`
		args = append(args, distributorID)
	}

	if _, err := repository.Exec(query, args...); err != nil {
		log.Error("priceRepository, UpdateDistributorProductPrices, err:", err.Error())
		return err
	}
	return nil
}

func (repository *MPriceRepositoryImpl) PublishByRMQ(request entity.PublishByRmqMPriceReq) error {
	var (
		r            model.MPricePublish
		sqlSetFields string
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)
	for i := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")
	if sqlSetFields == "" {
		return errors.New("no fields to update")
	}

	query := `UPDATE mst.m_price SET ` + sqlSetFields + `, updated_at = CURRENT_TIMESTAMP ` +
		`WHERE cust_id = :cust_id AND price_id = :price_id AND status = 1 AND is_del IS FALSE`

	sqlPatch.Args["price_id"] = request.PriceID
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error("priceRepository, PublishByRMQ, err:", err.Error())
		return err
	}

	nRows, err := result.RowsAffected()
	if err != nil {
		log.Error("priceRepository, PublishByRMQ, RowsAffected err:", err.Error())
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (repository *MPriceRepositoryImpl) FindBrokenDistributorChildLinks(parentProID int64, parentCustID string, distributorIDs []int64) ([]int64, error) {
	result := make([]int64, 0)
	if len(distributorIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT DISTINCT p.distributor_id
		FROM mst.m_product p
		JOIN mst.m_distributor d
		  ON d.cust_id = p.cust_id
		 AND d.distributor_id = p.distributor_id
		 AND d.parent_cust_id = ?
		WHERE p.distributor_id IN (?)
		  AND p.is_del = false
		  AND COALESCE(p.parent_pro_id, 0) != ?
	`
	args := []interface{}{parentCustID, distributorIDs, parentProID}

	query, argsExpanded, err := sqlx.In(query, args...)
	if err != nil {
		log.Error("priceRepository, FindBrokenDistributorChildLinks, sqlx.In err:", err.Error())
		return result, err
	}
	query = repository.Rebind(query)

	if err := repository.Select(&result, query, argsExpanded...); err != nil {
		log.Error("priceRepository, FindBrokenDistributorChildLinks, err:", err.Error())
		return result, err
	}
	return result, nil
}

func buildSortClause(raw string) string {
	allowedColumns := map[string]string{
		"created_date":   "prc.created_at",
		"effective_date": "prc.effective_date",
		"pro_code":       "prd.pro_code",
		"pro_name":       "prd.pro_name",
		"status":         "prc.status",
	}

	parts := strings.Split(raw, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		colSort := strings.Split(part, ":")
		if len(colSort) != 2 {
			continue
		}
		column := strings.TrimSpace(colSort[0])
		direction := strings.ToUpper(strings.TrimSpace(colSort[1]))
		if column == "" {
			continue
		}
		mappedColumn, ok := allowedColumns[column]
		if !ok {
			continue
		}
		if direction != "ASC" && direction != "DESC" {
			direction = "DESC"
		}
		items = append(items, fmt.Sprintf("%s %s", mappedColumn, direction))
	}
	return strings.Join(items, ", ")
}

func joinInts(values []int) string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strconv.Itoa(value))
	}
	return strings.Join(result, ",")
}

func joinInt64s(values []int64) string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, strconv.FormatInt(value, 10))
	}
	return strings.Join(result, ",")
}
