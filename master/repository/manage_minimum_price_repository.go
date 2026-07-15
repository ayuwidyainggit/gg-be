package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/sql_helper"
	"master/pkg/str"
	"math"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ManageMinimumPriceRepository interface {
	FindAllByCustId(dataFilter entity.ManageMinimumPriceQueryFilter, custId, parentCustId string) (manageMinimumPrice []model.ManageMinimumPriceRead, total int, lastPage int, err error)
	FindDetailById(manageMinimumPriceId int64, custId string) (model.ManageMinimumPriceRead, error)
	FindDetailByProId(ProId int64, custId string) (model.ManageMinimumPriceRead, error)
	Store(manageMinimumPrice model.ManageMinimumPrice) (int64, error)
	Update(manageMinimumPriceId int64, request entity.UpdateManageMinimumPrice) error
	Delete(custId string, manageMinimumPriceId int64, deletedBy int64) error
	UpdateStatus(manageMinimumPriceId int64, status int, custId string, userID int64) error
	TrxBegin()
	TrxCommit() error
	TrxRollback() error
}

func NewManageMinimumPriceRepository(db *sqlx.DB) ManageMinimumPriceRepository {
	return &manageMinumumPriceRepositoryImpl{db: db}
}

type manageMinumumPriceRepositoryImpl struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func (repository *manageMinumumPriceRepositoryImpl) TrxBegin() {
	repository.tx = repository.db.MustBegin()
}
func (repo *manageMinumumPriceRepositoryImpl) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *manageMinumumPriceRepositoryImpl) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *manageMinumumPriceRepositoryImpl) FindAllByCustId(dataFilter entity.ManageMinimumPriceQueryFilter, custId, parentCustId string) ([]model.ManageMinimumPriceRead, int, int, error) {

	manageMinimumPrice := []model.ManageMinimumPriceRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := `	mmp.cust_id,
						mmp.manage_minimum_price_id,
						mmp.base_price,
						mmp.limit_action,
						mmp.threshold,
						mmp.status_manage_minimum_price,
						mmp.pro_id,
						mmp.price1,
						mmp.price2,
						mmp.price3,
						mmp.price4,
						mmp.price5,
						mmp.price1_minimum,
						mmp.price2_minimum,
						mmp.price3_minimum,
						mmp.price4_minimum,
						mmp.price5_minimum,
						mmp.unit_id1,
						mmp.unit_id2,
						mmp.unit_id3,
						mmp.unit_id4,
						mmp.unit_id5,
						mmp.conv_unit2,
						mmp.conv_unit3,
						mmp.conv_unit4,
						mmp.conv_unit5,
						mmp.created_by,
						mmp.created_at,
						mmp.updated_by,
						mmp.updated_at,
						mp.pro_name,mp.pro_code`

	qWhere := ` WHERE mmp.cust_id = '` + custId + `' 
				AND mmp.is_del = false `

	if dataFilter.Query != "" {
		qWhere += ` AND mp.pro_code ILIKE '%` + dataFilter.Query + `%' `
	}

	if len(dataFilter.Status) > 0 {
		intArrStr := str.ArrayToString(dataFilter.Status, ",")
		qWhere += ` AND mmp.status_manage_minimum_price IN (` + intArrStr + `) `
	}

	qFrom := ` FROM mst.manage_minimum_price mmp
	LEFT JOIN sys.m_user u ON u.user_id = mmp.updated_by 
	LEFT JOIN mst.m_product mp ON mp.pro_id = mmp.pro_id`

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("manageMinimumPriceRepository, queryCount:", queryCount)
	var total int
	err := repository.db.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("manageMinimumPriceRepository, count total, err:", err.Error())
		return manageMinimumPrice, 0, 0, err
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
		sortBy := `mmp.manage_minimum_price_id`
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

	// fmt.Println("manageMinimumPriceRepository, querySelect:", querySelect)
	err = repository.db.Select(&manageMinimumPrice, querySelect)
	if err != nil {
		log.Println("manageMinimumPriceRepository, FindAllByCustId, err:", err.Error())
		return manageMinimumPrice, total, lastPage, err
	}

	return manageMinimumPrice, total, lastPage, nil
}

func (repository *manageMinumumPriceRepositoryImpl) Store(manageMinimumPrice model.ManageMinimumPrice) (int64, error) {
	query :=
		`INSERT INTO mst.manage_minimum_price(
			cust_id, base_price, limit_action,threshold, status_manage_minimum_price,
			pro_id,
			price1, price2, price3, price4,price5,
			price1_minimum, price2_minimum, price3_minimum, price4_minimum, price5_minimum,
			unit_id1, unit_id2,unit_id3,unit_id4, unit_id5,
			conv_unit2, conv_unit3, conv_unit4, conv_unit5,
			created_by, created_at, updated_by, updated_at, 
			is_del)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12, 
			$13, $14, $15, $16, 
			$17, $18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30
		) RETURNING manage_minimum_price_id;`
	lastInsertId := manageMinimumPrice.ManageMinimumPriceId
	err := repository.tx.QueryRow(query,
		manageMinimumPrice.CustId, manageMinimumPrice.BasePrice, manageMinimumPrice.LimitAction, manageMinimumPrice.Threshold,
		manageMinimumPrice.StatusManageMinimumPrice, manageMinimumPrice.ProId,
		manageMinimumPrice.Price1,
		manageMinimumPrice.Price2,
		manageMinimumPrice.Price3,
		manageMinimumPrice.Price4,
		manageMinimumPrice.Price5,
		manageMinimumPrice.PriceMinimum1,
		manageMinimumPrice.PriceMinimum2,
		manageMinimumPrice.PriceMinimum3,
		manageMinimumPrice.PriceMinimum4,
		manageMinimumPrice.PriceMinimum5,
		manageMinimumPrice.UnitId1,
		manageMinimumPrice.UnitId2,
		manageMinimumPrice.UnitId3,
		manageMinimumPrice.UnitId4,
		manageMinimumPrice.UnitId5,
		manageMinimumPrice.ConvUnit2,
		manageMinimumPrice.ConvUnit3,
		manageMinimumPrice.ConvUnit4,
		manageMinimumPrice.ConvUnit5,
		manageMinimumPrice.CreatedBy, manageMinimumPrice.CreatedAt, manageMinimumPrice.UpdatedBy, manageMinimumPrice.UpdatedAt, false).Scan(&lastInsertId)
	if err != nil {
		log.Println("manageMinimumPriceRepository, Store, err:", err.Error())
		return 0, err
	}
	return int64(*lastInsertId), nil
}

func (repository *manageMinumumPriceRepositoryImpl) Update(manageMinimumPriceId int64, request entity.UpdateManageMinimumPrice) error {
	var (
		r            model.ManageMinimumPriceUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("manageMinimumPriceRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.manage_minimum_price
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP,
			  	status_manage_minimum_price = 1
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND manage_minimum_price_id = :manage_minimum_price_id_old;`

	log.Println("manageMinimumPriceRepository, Update, query:", query)

	sqlPatch.Args["manage_minimum_price_id_old"] = manageMinimumPriceId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("manageMinimumPriceRepository, Update, err:", err.Error())
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

func (repository *manageMinumumPriceRepositoryImpl) Delete(custId string, manageMinimumPriceId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.manage_minimum_price
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND manage_minimum_price_id = :manage_minimum_price_id;`

	wMap := map[string]interface{}{
		"cust_id":                 custId,
		"manage_minimum_price_id": manageMinimumPriceId,
		"deleted_by":              deletedBy,
	}

	result, err := repository.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("ManageMinimumPriceRepository, Delete, err:", err.Error())
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

func (repo *manageMinumumPriceRepositoryImpl) FindDetailById(manageMinimumPriceId int64, custId string) (model.ManageMinimumPriceRead, error) {
	Details := model.ManageMinimumPriceRead{}
	query := `SELECT mmp.cust_id,
						mmp.manage_minimum_price_id,
						mmp.base_price,
						mmp.limit_action,
						mmp.threshold,
						mmp.status_manage_minimum_price,
						mmp.pro_id,
						mmp.price1,
						mmp.price2,
						mmp.price3,
						mmp.price4,
						mmp.price5,
						mmp.price1_minimum,
						mmp.price2_minimum,
						mmp.price3_minimum,
						mmp.price4_minimum,
						mmp.price5_minimum,
						mmp.unit_id1,
						mmp.unit_id2,
						mmp.unit_id3,
						mmp.unit_id4,
						mmp.unit_id5,
						mmp.conv_unit2,
						mmp.conv_unit3,
						mmp.conv_unit4,
						mmp.conv_unit5,
						mmp.created_by,
						mmp.created_at,
						mmp.updated_by,
						mmp.updated_at,mp.pro_name,mp.pro_code FROM mst.manage_minimum_price mmp
	LEFT JOIN sys.m_user u ON u.user_id = mmp.updated_by 
	LEFT JOIN mst.m_product mp ON mp.pro_id = mmp.pro_id AND mmp.cust_id = $2
	WHERE mmp.manage_minimum_price_id=$1 AND mmp.cust_id=$2 AND mmp.is_del=false`
	err := repo.db.Get(&Details, query, manageMinimumPriceId, custId)
	if err != nil {
		log.Println("ManageMinimumPriceRepository, FindDetailManageMinimumPriceAndCustId, err:", err.Error())
		return Details, err
	}

	return Details, nil
}

func (repo *manageMinumumPriceRepositoryImpl) FindDetailByProId(ProId int64, custId string) (model.ManageMinimumPriceRead, error) {
	Details := model.ManageMinimumPriceRead{}
	query := `SELECT mmp.cust_id,
						mmp.manage_minimum_price_id,
						mmp.base_price,
						mmp.limit_action,
						mmp.threshold,
						mmp.status_manage_minimum_price,
						mmp.pro_id,
						mmp.price1,
						mmp.price2,
						mmp.price3,
						mmp.price4,
						mmp.price5,
						mmp.price1_minimum,
						mmp.price2_minimum,
						mmp.price3_minimum,
						mmp.price4_minimum,
						mmp.price5_minimum,
						mmp.unit_id1,
						mmp.unit_id2,
						mmp.unit_id3,
						mmp.unit_id4,
						mmp.unit_id5,
						mmp.conv_unit2,
						mmp.conv_unit3,
						mmp.conv_unit4,
						mmp.conv_unit5,
						mmp.created_by,
						mmp.created_at,
						mmp.updated_by,
						mmp.updated_at,mp.pro_name,mp.pro_code FROM mst.manage_minimum_price mmp
	LEFT JOIN sys.m_user u ON u.user_id = mmp.updated_by 
	LEFT JOIN mst.m_product mp ON mp.pro_id = mmp.pro_id AND mmp.cust_id = $2
	WHERE mmp.pro_id=$1 AND mmp.cust_id=$2 AND mmp.is_del=false`
	err := repo.db.Get(&Details, query, ProId, custId)
	if err != nil {
		log.Println("ManageMinimumPriceRepository, FindDetailManageMinimumPriceAndCustId, err:", err.Error())
		return Details, err
	}

	return Details, nil
}

func (repo *manageMinumumPriceRepositoryImpl) UpdateStatus(manageMinimumPriceId int64, status int, custId string, userID int64) error {
	var nRows int64
	query := `UPDATE mst.manage_minimum_price
			SET status_manage_minimum_price = :status,
				updated_at = CURRENT_TIMESTAMP,
				updated_by = :updated_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND manage_minimum_price_id = :manage_minimum_price_id;`

	wMap := map[string]interface{}{
		"cust_id":                 custId,
		"manage_minimum_price_id": manageMinimumPriceId,
		"status":                  status,
		"updated_by":              userID,
	}

	result, err := repo.db.NamedExec(query, wMap)
	if err != nil {
		log.Println("ManageMinimumPriceRepository, Update Is Active, err:", err.Error())
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
