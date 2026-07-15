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

type tprRepositoryImpl struct {
	*sqlx.DB
}

func NewTprRepository(db *sqlx.DB) *tprRepositoryImpl {
	return &tprRepositoryImpl{db}
}

type TprRepository interface {
	FindOneByTprIdAndCustId(tprid int64, custId string) (model.MTpr, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MTpr, int, int, error)
	Store(tpr model.MTpr) (int64, error)
	FindOneByTprCodeAndCustId(code string, custId string) (model.MTpr, error)
	Update(tprID int64, request entity.UpdateTprRequest) error
	Delete(custId string, tprId int64, deletedBy int64) error
}

func (repository *tprRepositoryImpl) FindOneByTprIdAndCustId(tprid int64, custId string) (model.MTpr, error) {
	tpr := model.MTpr{}
	query := `SELECT 
				*
			  FROM mst.m_tpr 
			  WHERE tpr_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&tpr, query, tprid, custId)
	if err != nil {
		log.Println("tprRepository, FindOneByTprIdAndCustId, err:", err.Error())
		return tpr, err
	}

	return tpr, nil
}
func (repository *tprRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MTpr, int, int, error) {
	tprs := []model.MTpr{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` t.*,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE t.is_del = false 
				AND t.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (t.tpr_code ILIKE '%` + dataFilter.Query + `%' 
					OR t.tpr_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND t.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND t.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_tpr t
	LEFT JOIN sys.m_user u ON u.user_id = t.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("vehicleRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("tprRepository, count total, err:", err.Error())
		return tprs, 0, 0, err
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
		sortBy := `tpr_id`
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

	// log.Println("tprRepository, querySelect:", querySelect)
	err = repository.Select(&tprs, querySelect)
	if err != nil {
		log.Println("tprRepository, FindAllByCustId, err:", err.Error())
		return tprs, total, lastPage, err
	}

	return tprs, total, lastPage, nil
}

func (repository *tprRepositoryImpl) FindOneByTprCodeAndCustId(code string, custId string) (model.MTpr, error) {
	tpr := model.MTpr{}
	query := `SELECT 
				*
			  FROM mst.m_tpr 
			  WHERE tpr_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&tpr, query, code, custId)
	if err != nil {
		log.Println("tprRepository, FindOneByTprIdAndCustId, err:", err.Error())
		return tpr, err
	}

	return tpr, nil
}

func (repository *tprRepositoryImpl) Store(tpr model.MTpr) (int64, error) {
	query :=
		`INSERT INTO mst.m_tpr(
			cust_id, tpr_code, tpr_name,date_start,date_end,range_type,promo_item_type,
			is_multiple_promo,is_all_ot_type,is_all_ot_grp,is_all_sales,is_all_ot,
			is_all_sales_team,is_all_industry,is_max,is_max_value,deduction,min_invoice_value,
			is_active,created_by,created_at,updated_by,updated_at,is_del,deleted_by,deleted_at
		)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		) RETURNING tpr_id;`
	lastInsertId := tpr.TprID

	err := repository.QueryRow(query,
		tpr.CustID, tpr.TprCode, tpr.TprName, tpr.DateStart, tpr.DateEnd, tpr.RangeType, tpr.PromoItemType, tpr.IsMultiplePromo, tpr.IsAllOtType, tpr.IsAllOtGrp, tpr.IsAllSales, tpr.IsAllOt, tpr.IsAllSalesTeam, tpr.IsAllIndustry, tpr.IsMax, tpr.IsMaxValue, tpr.Deduction, tpr.MinInvoiceValue, tpr.IsActive, tpr.CreatedBy, tpr.CreatedAt, tpr.UpdatedBy, tpr.UpdatedAt, tpr.IsDel, tpr.DeletedBy, tpr.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("tprRepository, Store, err:", err.Error())
		return tpr.TprID, err
	}
	return tpr.TprID, nil
}

func (repository *tprRepositoryImpl) Update(tprID int64, request entity.UpdateTprRequest) error {
	var (
		r            model.MTprUpdate
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

	query := `UPDATE mst.m_tpr
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND tpr_id = :tpr_id_old;`

	// log.Println("tprRepository, Update, query:", query)

	sqlPatch.Args["tpr_id_old"] = tprID
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("vehicleRepository, Update, err:", err.Error())
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

func (repository *tprRepositoryImpl) Delete(custId string, tprId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_tpr
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND tpr_id = :tpr_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"tpr_id":     tprId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("VehicleRepository, Delete, err:", err.Error())
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
