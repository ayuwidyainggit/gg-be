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

type MParkingFundRepository interface {
	FindOneByOutletAndProIdIdAndCustId(OutletId int, ProId int, custId string) (model.MParkingFundList, error)
	FindOneBymParkingFundIdAndCustId(ParkingFundId int, custId string) (model.MParkingFundList, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.MParkingFundList, total int, lastPage int, err error)
	Store(mParkingFund model.MParkingFund) (int, error)
	Update(mParkingFundId int, request entity.UpdateMParkingFundRequest) error
	Delete(custId string, parkingFundId int, deletedBy int64) error
}

func NewMParkingFundRepository(db *sqlx.DB) MParkingFundRepository {
	return &mParkingFundRepositoryImpl{db}
}

type mParkingFundRepositoryImpl struct {
	*sqlx.DB
}

func (repository *mParkingFundRepositoryImpl) FindOneByOutletAndProIdIdAndCustId(OutletId int, ProId int, custId string) (model.MParkingFundList, error) {
	mParkingFund := model.MParkingFundList{}
	query := `SELECT 
		pf.*,o.outlet_code, o.outlet_name,
		p.pro_code, p.pro_name FROM mst.m_parking_fund pf
	LEFT JOIN mst.m_outlet o ON o.outlet_id = pf.outlet_id AND o.cust_id = '` + custId + `'
	LEFT JOIN mst.m_product_dist p ON p.pro_id = pf.pro_id 
	WHERE pf.outlet_id = $1 AND pf.pro_id = $2 AND cust_id = $3 and pf.is_del = FALSE`
	err := repository.Get(&mParkingFund, query, OutletId, ProId, custId)
	if err != nil {
		log.Println("mParkingFundRepository, FindOneByOutletAndProIdIdAndCustId, err:", err.Error())
		return mParkingFund, err
	}

	return mParkingFund, nil
}

func (repository *mParkingFundRepositoryImpl) FindOneBymParkingFundIdAndCustId(ParkingFundId int, custId string) (model.MParkingFundList, error) {
	mParkingFund := model.MParkingFundList{}
	query := `SELECT 
		pf.*,o.outlet_code, o.outlet_name,
		p.pro_code, p.pro_name FROM mst.m_parking_fund pf
	LEFT JOIN mst.m_outlet o ON o.outlet_id = pf.outlet_id AND o.cust_id = '` + custId + `'
	LEFT JOIN mst.m_product p ON p.pro_id = pf.pro_id
	WHERE pf.parking_fund_id = $1 AND pf.cust_id = $2 and pf.is_del = FALSE`
	err := repository.Get(&mParkingFund, query, ParkingFundId, custId)
	if err != nil {
		log.Println("mParkingFundRepository, FindOneByOutletAndProIdIdAndCustId, err:", err.Error())
		return mParkingFund, err
	}

	return mParkingFund, nil
}

func (repository *mParkingFundRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.MParkingFundList, int, int, error) {

	mParkingFund := []model.MParkingFundList{}
	selectCount := ` COUNT(pf.*) AS total `
	selectField := ` pf.cust_id, pf.parking_fund_id, pf.outlet_id, pf.pro_id, 
					pf.updated_at, pf.p_disc,
					o.outlet_code, o.outlet_name,
					p.pro_code, p.pro_name,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE pf.is_del = false and pf.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (o.outlet_code ILIKE '%` + dataFilter.Query + `%' OR 
						 p.pro_code ILIKE '%` + dataFilter.Query + `%' OR 
						 o.outlet_name ILIKE '%` + dataFilter.Query + `%' OR 
						 p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` FROM mst.m_parking_fund pf
			   LEFT JOIN sys.m_user u ON u.user_id = pf.updated_by 
			   LEFT JOIN mst.m_outlet o ON o.outlet_id = pf.outlet_id AND o.cust_id = '` + custId + `'
			   LEFT JOIN mst.m_product p ON p.pro_id = pf.pro_id `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("mParkingFundRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("mParkingFundRepository, count total, err:", err.Error())
		return mParkingFund, 0, 0, err
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
		sortBy := `parking_fund_id`
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

	// log.Println("mParkingFundRepository, querySelect:", querySelect)
	err = repository.Select(&mParkingFund, querySelect)
	if err != nil {
		log.Println("mParkingFundRepository, FindAllByCustId, err:", err.Error())
		return mParkingFund, total, lastPage, err
	}

	return mParkingFund, total, lastPage, nil
}

func (repository *mParkingFundRepositoryImpl) Store(mParkingFund model.MParkingFund) (int, error) {
	query :=
		`INSERT INTO mst.m_parking_fund(cust_id, outlet_id, pro_id, p_disc, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING parking_fund_id;`
	lastInsertId := mParkingFund.ParkingFundId
	err := repository.QueryRow(query,
		mParkingFund.CustId, mParkingFund.OutletId, mParkingFund.ProId, mParkingFund.PDisc,
		mParkingFund.CreatedBy, mParkingFund.CreatedAt, mParkingFund.UpdatedBy, mParkingFund.UpdatedAt, mParkingFund.IsDel, mParkingFund.DeletedBy, mParkingFund.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("mParkingFundRepository, Store, err:", err.Error())
		return mParkingFund.ParkingFundId, err
	}
	return mParkingFund.ParkingFundId, nil
}

func (repository *mParkingFundRepositoryImpl) Update(parkingFundId int, request entity.UpdateMParkingFundRequest) error {
	var (
		r            model.MParkingFundUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("mParkingFundRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_parking_fund
			  SET ` + sqlSetFields + `, updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false AND cust_id = :cust_id AND parking_fund_id = :parking_fund_id_old;`

	log.Println("mParkingFundRepository, Update, query:", query)

	sqlPatch.Args["parking_fund_id_old"] = parkingFundId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("mParkingFundRepository, Update, err:", err.Error())
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

func (repository *mParkingFundRepositoryImpl) Delete(custId string, parkingFundId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_parking_fund SET is_del = true, deleted_at = CURRENT_TIMESTAMP, deleted_by = :deleted_by 
			WHERE is_del = false AND cust_id = :cust_id AND parking_fund_id = :parking_fund_id;`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"parking_fund_id": parkingFundId,
		"deleted_by":      deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("mParkingFundRepository, Delete, err:", err.Error())
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
