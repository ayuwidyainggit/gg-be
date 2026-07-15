package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/errmsg"
	"master/pkg/sql_helper"
	"master/pkg/structs"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type SpPriceRepository interface {
	TrxBegin() (*spPriceTransaction, error)
	FindOneBySpPriceIdAndCustID(params entity.MSpPriceParams) (model.MSpPriceWithDetails, error)
	FindDetailSpPriceIdAndCustId(params entity.MSpPriceParams) ([]model.MSpPriceDetView, error)
	FindDetailSpPriceIdAndCustIdPublish(params entity.MSpPriceParams) ([]model.MSpPriceDetPublish, error)
	FindAllByCustID(dataFilter entity.MSpPriceQueryFilter, custId string) ([]model.MSpPriceDetail, int, int, error)
	FindOneProductByProID(proID int64, custID string) (model.Product, error)
	UpdateStatusByRMQ(request entity.PublishUnpublishSPriceReq) error
}

func NewSpPriceRepository(db *sqlx.DB) SpPriceRepository {
	return &spPriceRepositoryImpl{db}
}

type spPriceRepositoryImpl struct {
	*sqlx.DB
}

type spPriceTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

func NewTransactionSpPrice(db *sqlx.DB) (trxObj *spPriceTransaction, err error) {
	trx := db.MustBegin()

	return &spPriceTransaction{tx: trx, db: db}, nil
}

func (repo *spPriceRepositoryImpl) TrxBegin() (*spPriceTransaction, error) {
	trxObj, err := NewTransactionSpPrice(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *spPriceTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *spPriceTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repo *spPriceTransaction) InsertPrice(spPrice *model.MSpPrice) error {

	log.Info("spPrice:", structs.StructToJson(spPrice))
	query :=
		`INSERT INTO mst.m_sp_price(
			cust_id, sp_price_id, start_date, end_date, price_grp_id, pro_id, 
			unit_id1, unit_id2, unit_id3,
			sell_price1, sell_price2, sell_price3, 
			new_sell_price1, new_sell_price2, new_sell_price3, 
			conv_unit2, conv_unit3, status, 
			created_by, created_at, updated_by, updated_at
	)
	VALUES ( 
		$1, $2, $3, $4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12, 
		$13, $14, $15, 
		$16, $17, $18, 
		$19, $20, $21, $22
	) RETURNING sp_price_id;`
	lastInsertId := spPrice.SpPriceID
	err := repo.tx.QueryRow(query,
		spPrice.CustID, spPrice.SpPriceID, spPrice.StartDate, spPrice.EndDate, spPrice.PriceGrpID, spPrice.ProID,
		spPrice.UnitId1, spPrice.UnitId2, spPrice.UnitId3,
		spPrice.SellPrice1, spPrice.SellPrice2, spPrice.SellPrice3,
		spPrice.NewSellPrice1, spPrice.NewSellPrice2, spPrice.NewSellPrice3,
		spPrice.ConvUnit2, spPrice.ConvUnit3, spPrice.Status,
		spPrice.CreatedBy, spPrice.CreatedAt, spPrice.UpdatedBy, spPrice.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Error("spPriceRepository, Store, err:", err.Error())
		return err
	}
	spPrice.SpPriceID = lastInsertId
	return nil
}

func (repo *spPriceTransaction) InsertPriceDetail(spPriceID string, spPriceDetail *model.MSpPriceDet) error {
	query := `INSERT INTO mst.m_sp_price_det(cust_id, ref_type, sp_price_id, sp_price_det_id, ` +
		`ref_id, sell_price1, sell_price2, sell_price3, ` +
		`new_sell_price1, new_sell_price2, new_sell_price3, ` +
		`created_by, created_at, updated_by, updated_at) ` +
		`VALUES ($1, $2, $3, $4, ` +
		`$5, $6, $7, $8, ` +
		`$9, $10, $11, ` +
		`$12, $13, $14, $15) ` +
		`RETURNING sp_price_det_id;`

	lastInsertId := spPriceDetail.SpPriceDetID
	err := repo.tx.QueryRow(query, spPriceDetail.CustID, spPriceDetail.RefType, spPriceID, spPriceDetail.SpPriceDetID,
		spPriceDetail.RefID, spPriceDetail.SellPrice1, spPriceDetail.SellPrice2, spPriceDetail.SellPrice3,
		spPriceDetail.NewSellPrice1, spPriceDetail.NewSellPrice2, spPriceDetail.NewSellPrice3,
		spPriceDetail.CreatedBy, spPriceDetail.CreatedAt, spPriceDetail.UpdatedBy, spPriceDetail.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Error("spPriceRepository, Store, err:", err.Error())
		return err
	}
	spPriceDetail.SpPriceDetID = lastInsertId
	return nil
}

func (repository *spPriceTransaction) Update(spPriceID string, request entity.UpdateMSpPriceBody) error {
	var (
		r            model.MSpPriceUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	if err := json.Unmarshal(reqByte, &r); err != nil {
		return err
	}

	sqlPatch := sql_helper.SQLPatches(r)
	_, _ = json.Marshal(sqlPatch)
	// log.Info("spPriceTransaction, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sp_price ` +
		`SET ` + sqlSetFields + `,` +
		`updated_at = CURRENT_TIMESTAMP ` +
		`WHERE cust_id = :cust_id ` +
		`AND sp_price_id = :sp_price_id;`

	sqlPatch.Args["sp_price_id"] = spPriceID
	sqlPatch.Args["cust_id"] = request.CustID
	// log.Info("spPriceRepository, Update, query:", query)
	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error("spPriceRepository, Update, err:", err.Error())
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

func (repository *spPriceTransaction) UpdateDetail(spPriceID, spPriceDetID string, request entity.MSpPriceDetUpdate) error {
	var (
		r            model.MSpPriceUpdateDet
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("tprRepository, Update, Fields & Args: %s\n", data)
	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sp_price_det
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE sp_price_det_id = :sp_price_det_id
			  AND sp_price_id=:sp_price_id;`

	// log.Error("tprRepository, Update, query:", query)
	sqlPatch.Args["sp_price_id"] = spPriceID
	sqlPatch.Args["sp_price_det_id"] = spPriceDetID

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error("spPriceRepository, Update, err:", err.Error())
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

func (repo *spPriceRepositoryImpl) FindOneBySpPriceIdAndCustID(params entity.MSpPriceParams) (model.MSpPriceWithDetails, error) {
	spPrice := model.MSpPriceWithDetails{}
	query := `SELECT sp.*, ` +
		`p.pro_code, p.pro_name, p.purch_price1, p.purch_price2, p.purch_price3, ` +
		`pg.sp_price_grp_code AS price_grp_code, pg.sp_price_grp_name AS price_grp_name, ` +
		`un1.unit_name AS unit_name1, un2.unit_name AS unit_name2, un3.unit_name AS unit_name3 ` +
		`FROM mst.m_sp_price sp ` +
		`LEFT JOIN mst.m_sp_price_group pg ON pg.sp_price_grp_id = sp.price_grp_id AND pg.cust_id = '` + params.CustID + `' ` +
		`LEFT JOIN mst.m_product p ON p.pro_id = sp.pro_id AND p.cust_id = '` + params.ParentCustID + `' ` +
		`LEFT JOIN mst.m_unit un1 ON un1.unit_id = sp.unit_id1 AND un1.cust_id = '` + params.ParentCustID + `' ` +
		`LEFT JOIN mst.m_unit un2 ON un2.unit_id = sp.unit_id2 AND un2.cust_id = '` + params.ParentCustID + `' ` +
		`LEFT JOIN mst.m_unit un3 ON un3.unit_id = sp.unit_id3 AND un3.cust_id = '` + params.ParentCustID + `' ` +
		`WHERE sp_price_id = $1 AND sp.cust_id = $2;`
	err := repo.Get(&spPrice, query, params.SpPriceID, params.CustID)
	if err != nil {
		log.Error(err.Error())
		return spPrice, err
	}

	return spPrice, nil
}

func (repo *spPriceRepositoryImpl) FindDetailSpPriceIdAndCustId(params entity.MSpPriceParams) ([]model.MSpPriceDetView, error) {
	spPriceDetail := []model.MSpPriceDetView{}
	query := `SELECT *,
	CASE 
		WHEN spd.ref_type = 1 THEN 
			(SELECT sales_team_code FROM mst.m_sales_team st WHERE st.cust_id = '` + params.ParentCustID + `' AND sales_team_id = spd.ref_id)
		WHEN spd.ref_type = 5 THEN 
			(SELECT ot_type_code FROM mst.m_outlet_type ot WHERE ot.cust_id = '` + params.ParentCustID + `' AND ot_type_id = spd.ref_id)
		WHEN spd.ref_type = 10 THEN 
			(SELECT ot_grp_code FROM mst.m_outlet_group og WHERE og.cust_id = '` + params.ParentCustID + `' AND ot_grp_id = spd.ref_id)
		WHEN spd.ref_type = 15 THEN 
			(SELECT outlet_code FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id)
		WHEN spd.ref_type = 20 THEN 
			(SELECT outlet_code FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id)
		ELSE ''
	END AS ref_code,
	CASE 
		WHEN spd.ref_type = 1 THEN 
			(SELECT sales_team_name FROM mst.m_sales_team st WHERE st.cust_id = '` + params.ParentCustID + `' AND sales_team_id = spd.ref_id)
		WHEN spd.ref_type = 5 THEN 
			(SELECT ot_type_name FROM mst.m_outlet_type ot WHERE ot.cust_id = '` + params.ParentCustID + `' AND ot_type_id = spd.ref_id)
		WHEN spd.ref_type = 10 THEN 
			(SELECT ot_grp_name FROM mst.m_outlet_group og WHERE og.cust_id = '` + params.ParentCustID + `' AND ot_grp_id = spd.ref_id)
		WHEN spd.ref_type = 15 THEN 
			(SELECT outlet_name FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id)
		WHEN spd.ref_type = 20 THEN 
			(SELECT outlet_name FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id)
		ELSE ''
	END AS ref_name
	FROM mst.m_sp_price_det spd
	WHERE spd.sp_price_id = $1 
	AND spd.cust_id = $2`
	err := repo.Select(&spPriceDetail, query, params.SpPriceID, params.CustID)
	if err != nil {
		log.Error("spPriceRepository, FindDetailSpPriceIdAndCustId, err:", err.Error())
		return spPriceDetail, err
	}

	return spPriceDetail, nil
}

func (repository *spPriceRepositoryImpl) FindAllByCustID(dataFilter entity.MSpPriceQueryFilter, custId string) ([]model.MSpPriceDetail, int, int, error) {
	mspPrice := []model.MSpPriceDetail{}

	selectCount := ` COUNT(*) AS total `
	selectField := ` sp.*,
					p.pro_code, p.pro_name,
					u1.unit_name AS unit_name1, u2.unit_name AS unit_name2, u3.unit_name AS unit_name3,
					pg.sp_price_grp_code AS price_grp_code, pg.sp_price_grp_name AS price_grp_name  `
	qWhere := ` WHERE sp.cust_id = '` + custId + `' `

	if dataFilter.Status > 0 {
		qWhere += `AND sp.status = ` + strconv.Itoa(dataFilter.Status) + ` `
	}
	qFrom := ` 	FROM mst.m_sp_price sp
				LEFT JOIN mst.m_sp_price_group pg ON pg.sp_price_grp_id = sp.price_grp_id AND pg.cust_id = '` + dataFilter.CustID + `'
				LEFT JOIN mst.m_product p ON p.pro_id = sp.pro_id AND p.cust_id = '` + dataFilter.ParentCustID + `'
				LEFT JOIN mst.m_unit u1 ON u1.unit_id = p.unit_id1 AND u1.cust_id = '` + dataFilter.ParentCustID + `'
				LEFT JOIN mst.m_unit u2 ON u2.unit_id = p.unit_id2 AND u2.cust_id = '` + dataFilter.ParentCustID + `'
				LEFT JOIN mst.m_unit u3 ON u3.unit_id = p.unit_id3 AND u3.cust_id = '` + dataFilter.ParentCustID + `'`
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Error("spPriceRepository, count total, err:", err.Error())
		return mspPrice, 0, 0, err
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
		sortBy := `sp.sp_price_id`
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

	err = repository.Select(&mspPrice, querySelect)
	if err != nil {
		log.Error("spPriceRepository, FindAllByCustID, err:", err.Error())
		return mspPrice, total, lastPage, err
	}

	return mspPrice, total, lastPage, nil
}

func (repository *spPriceTransaction) Delete(custId, sppriceId string) error {
	var nRows int64
	query := `DELETE FROM mst.m_sp_price ` +
		`WHERE cust_id = :cust_id AND ` +
		`sp_price_id = :sp_price_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"sp_price_id": sppriceId,
	}
	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Error("\nspPriceRepository, Delete, err:", err.Error())
		return err
	}
	if nRows, err = result.RowsAffected(); err != nil {
		log.Error(err)
		return err
	}
	if nRows == 0 {
		return errors.New(errmsg.ERROR_NO_ROWS_AFFECTED)
	}

	return nil
}

func (repository *spPriceTransaction) DeleteDetailNotIn(sppriceId string, sppriceDetId []string) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_sp_price_det	WHERE sp_price_id = ? AND sp_price_det_id NOT IN (?);", sppriceId, sppriceDetId)
	if err != nil {
		return err
	}
	query = repository.db.Rebind(query)
	_, err = repository.tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *spPriceRepositoryImpl) FindOneProductByProID(proID int64, custID string) (model.Product, error) {
	productID := model.Product{}

	query := `SELECT pro_id, pro_code, pro_name, unit_id1, unit_id2, unit_id3, conv_unit2, conv_unit3, sell_price1, sell_price2, sell_price3
			  FROM mst.m_product 
			  WHERE pro_id = $1 
			  AND cust_id = $2`

	err := repository.Get(&productID, query, proID, custID)
	if err != nil {
		log.Error("priceRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return productID, err
	}

	return productID, nil
}

func (repository *spPriceTransaction) DeleteDetails(custId, sppriceId string) error {
	var nRows int64
	query := `DELETE FROM mst.m_sp_price_det ` +
		`WHERE cust_id = :cust_id AND ` +
		`sp_price_id = :sp_price_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"sp_price_id": sppriceId,
	}
	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if nRows, err = result.RowsAffected(); err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("sp price det, total deleted: ", nRows)

	return nil
}

func (repo *spPriceRepositoryImpl) FindDetailSpPriceIdAndCustIdPublish(params entity.MSpPriceParams) ([]model.MSpPriceDetPublish, error) {
	spPriceDetail := []model.MSpPriceDetPublish{}
	query := `SELECT spd.*, sp.pro_id, sp.start_date, sp.end_date, ` +
		`CASE ` +
		`WHEN spd.ref_type = 1 THEN ` +
		`(SELECT sales_team_code FROM mst.m_sales_team st WHERE st.cust_id = '` + params.ParentCustID + `' AND sales_team_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 5 THEN ` +
		`(SELECT ot_type_code FROM mst.m_outlet_type ot WHERE ot.cust_id = '` + params.ParentCustID + `' AND ot_type_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 10 THEN ` +
		`(SELECT ot_grp_code FROM mst.m_outlet_group og WHERE og.cust_id = '` + params.ParentCustID + `' AND ot_grp_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 15 THEN ` +
		`(SELECT outlet_code FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 20 THEN ` +
		`(SELECT outlet_code FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id) ` +
		`ELSE '' ` +
		`END AS ref_code,` +
		`CASE ` +
		`WHEN spd.ref_type = 1 THEN ` +
		`(SELECT sales_team_name FROM mst.m_sales_team st WHERE st.cust_id = '` + params.ParentCustID + `' AND sales_team_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 5 THEN ` +
		`(SELECT ot_type_name FROM mst.m_outlet_type ot WHERE ot.cust_id = '` + params.ParentCustID + `' AND ot_type_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 10 THEN ` +
		`(SELECT ot_grp_name FROM mst.m_outlet_group og WHERE og.cust_id = '` + params.ParentCustID + `' AND ot_grp_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 15 THEN ` +
		`(SELECT outlet_name FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id) ` +
		`WHEN spd.ref_type = 20 THEN ` +
		`(SELECT outlet_name FROM mst.m_outlet o WHERE o.cust_id = '` + params.CustID + `' AND outlet_id = spd.ref_id) ` +
		`ELSE '' ` +
		`END AS ref_name ` +
		`FROM mst.m_sp_price_det spd ` +
		`LEFT JOIN mst.m_sp_price sp ON sp.sp_price_id = spd.sp_price_id AND sp.cust_id = spd.cust_id ` +
		`WHERE spd.sp_price_id = $1 ` +
		`AND spd.cust_id = $2`
	err := repo.Select(&spPriceDetail, query, params.SpPriceID, params.CustID)
	if err != nil {
		log.Error(err.Error())
		return spPriceDetail, err
	}

	return spPriceDetail, nil
}

func (repository *spPriceRepositoryImpl) UpdateStatusByRMQ(request entity.PublishUnpublishSPriceReq) error {
	var (
		r            model.MSpPricePublish
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	oldStatus := "1"
	if request.Status == 7 {
		oldStatus = "10"
	}
	query := `UPDATE mst.m_sp_price ` +
		`SET ` + sqlSetFields + `, ` +
		`updated_at = CURRENT_TIMESTAMP ` +
		`WHERE cust_id = :cust_id ` +
		`AND sp_price_id = :sp_price_id ` +
		`AND status = ` + oldStatus + `;`

	sqlPatch.Args["sp_price_id"] = request.SpPriceID
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// log.Info("PublishByRMQ, result -> ", structs.StructToJson(result))
	if nRows, err = result.RowsAffected(); err != nil {
		return errors.New("no rows affected")
	}
	if nRows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
