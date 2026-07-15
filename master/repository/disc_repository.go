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

func NewDiscRepository(db *sqlx.DB) *discRepositoryImpl {
	return &discRepositoryImpl{db}
}

type discRepositoryImpl struct {
	*sqlx.DB
}

type discTransaction struct {
	tx *sqlx.Tx
	db *sqlx.DB
}

type DiscRepository interface {
	TrxBegin() (*discTransaction, error)
	FindOneByDiscIdAndCustId(discId int64, custId string) (model.Disc, error)
	FindOneByDiscCodeAndCustId(discCode, custId string) (model.Disc, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (disc []model.Disc, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (disc []model.Disc, total int, lastPage int, err error)
	//Store(disc model.Disc) (int, error)
	//Update(discId int64, request entity.UpdateDiscRequest) error
	Delete(custId string, discId int, deletedBy int64) error
	FindDetailByDiscIdAndCustId(discId int64, custId string) ([]model.DiscDet, error)
}

func NewdiscRepository(db *sqlx.DB) *discRepositoryImpl {
	return &discRepositoryImpl{db}
}

func NewDiscTransaction(db *sqlx.DB) (trxObj *discTransaction, err error) {
	trx := db.MustBegin()

	return &discTransaction{tx: trx, db: db}, nil
}

func (repo *discRepositoryImpl) TrxBegin() (*discTransaction, error) {
	trxObj, err := NewDiscTransaction(repo.DB)
	if err != nil {
		return nil, err
	}
	return trxObj, nil
}

func (repo *discTransaction) TrxCommit() error {
	return repo.tx.Commit()
}

func (repo *discTransaction) TrxRollback() error {
	return repo.tx.Rollback()
}

func (repository *discRepositoryImpl) FindOneByDiscIdAndCustId(discId int64, custId string) (model.Disc, error) {
	disc := model.Disc{}
	query := `SELECT 
				cust_id, disc_id, disc_code, disc_name,
				start_date, end_date, range_type, is_multiple,
				purchase_limit, disc_type, disc_perc, disc_value,
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_disc 
			  WHERE is_del = false 
			  AND disc_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&disc, query, discId, custId)
	if err != nil {
		log.Println("discRepository, FindOneBySBrand1IdAndCustId, err:", err.Error())
		return disc, err
	}

	return disc, nil
}
func (repo *discRepositoryImpl) FindDetailByDiscIdAndCustId(discId int64, custId string) ([]model.DiscDet, error) {
	discDet := []model.DiscDet{}
	query := `SELECT 
	cust_id, disc_id, row_no, min_value,
	max_value, disc_type, disc_perc, disc_value,
	disc_det_id
  FROM mst.m_disc_det 
  WHERE disc_id = $1 
  AND cust_id = $2`
	err := repo.Select(&discDet, query, discId, custId)
	if err != nil {
		log.Println("discRepository, FindOneBySBrand1IdAndCustId, err:", err.Error())
		return discDet, err
	}

	return discDet, nil
}

func (repository *discRepositoryImpl) FindOneByDiscCodeAndCustId(discCode, custId string) (model.Disc, error) {
	disc := model.Disc{}
	query := `SELECT 
				cust_id, disc_id, disc_code, disc_name, 
				start_date, end_date, range_type, is_multiple,
				purchase_limit, disc_type, disc_perc, disc_value,
				is_active, created_by, created_at, updated_by, 
				updated_at, is_del, deleted_by, deleted_at
			  FROM mst.m_disc 
			  WHERE is_del = false 
			  AND disc_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&disc, query, discCode, custId)
	if err != nil {
		log.Println("discRepository, FindOneByDiscCodeAndCustId, err:", err.Error())
		return disc, err
	}

	return disc, nil
}

func (repository *discRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Disc, int, int, error) {

	discs := []model.Disc{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` d.cust_id, d.disc_id, d.disc_code, d.disc_name,
	d.start_date, d.end_date, d.range_type, d.is_multiple,
	d.purchase_limit, d.disc_type, d.disc_perc, d.disc_value,
	d.is_active, d.created_by, d.created_at, d.updated_by, 
	d.updated_at, d.is_del, d.deleted_by, d.deleted_at,
	u.user_fullname AS updated_by_name `
	qWhere := ` WHERE d.is_del = false 
				AND d.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (d.disc_code ILIKE '%` + dataFilter.Query + `%' 
					OR d.disc_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND d.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND d.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_disc d
	LEFT JOIN sys.m_user u ON u.user_id = d.updated_by `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("discRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("discRepository, count total, err:", err.Error())
		return discs, 0, 0, err
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
		sortBy := `disc_id`
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

	// log.Println("discRepository, querySelect:", querySelect)
	err = repository.Select(&discs, querySelect)
	if err != nil {
		log.Println("discRepository, FindAllByCustId, err:", err.Error())
		return discs, total, lastPage, err
	}

	return discs, total, lastPage, nil
}

func (repository *discRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Disc, int, int, error) {

	discs := []model.Disc{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` disc_id, disc_code, disc_name, start_date, end_date,
					range_type, is_multiple, purchase_limit, disc_type, disc_perc, disc_value `
	qWhere := ` WHERE is_del = false AND is_active = true 
				AND cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (disc_code ILIKE '%` + dataFilter.Query + `%' 
					OR disc_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	queryCount := `SELECT ` + selectCount + ` FROM mst.m_disc ` + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_disc ` + qWhere

	// log.Println("discRepository, FindAllByCustIdLookupMode, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("discRepository, FindAllByCustIdLookupMode, count total, err:", err.Error())
		return discs, 0, 0, err
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
		sortBy := `disc_code`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("discRepository, FindAllByCustIdLookupMode, querySelect:", querySelect)
	err = repository.Select(&discs, querySelect)
	if err != nil {
		log.Println("discRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return discs, total, 1, err
	}

	return discs, total, 1, nil
}

func (repository *discTransaction) Store(disc *model.Disc) error {
	query :=
		`INSERT INTO mst.m_disc(
			cust_id, disc_code, disc_name, 
			start_date, end_date, range_type, 
			is_multiple, purchase_limit, disc_type, 
			disc_perc, disc_value, is_active, 
			created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, 
			deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8, $9, 
			$10, $11, $12, 
			$13, $14, $15, 
			$16, $17, $18, 
			$19
		) RETURNING disc_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		disc.CustId, disc.DiscCode, disc.DiscName,
		disc.StartDate, disc.EndDate, disc.RangeType,
		disc.IsMultiple, disc.PurchaseLimit, disc.DiscType,
		disc.DiscPerc, disc.DiscValue, disc.IsActive,
		disc.CreatedBy, disc.CreatedAt, disc.UpdatedBy,
		disc.UpdatedAt, disc.IsDel, disc.DeletedBy,
		disc.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("discRepository, Store, err:", err.Error())
		return err
	}
	disc.DiscId = lastInsertId
	return nil
}

func (repository *discTransaction) StoreDetail(discDet *model.DiscDet) error {
	query :=
		`INSERT INTO mst.m_disc_det(
			cust_id, disc_id, row_no, min_value,
	max_value, disc_type, disc_perc, disc_value)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8
		) RETURNING disc_id;`
	var lastInsertId int64
	err := repository.tx.QueryRow(query,
		discDet.CustID, discDet.DiscID, discDet.RowNo, discDet.MinValue,
		discDet.MaxValue, discDet.DiscType, discDet.DiscPerc, discDet.DiscValue).Scan(&lastInsertId)
	if err != nil {
		log.Println("discRepository, Store, err:", err.Error())
		return err
	}
	discDet.DiscDetID = &lastInsertId
	return nil
}

func (repository *discTransaction) Update(discId int64, request entity.UpdateDiscRequest) error {
	var (
		r            model.DiscUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)

	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("discRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_disc
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND disc_id = :disc_id;`

	log.Println("discRepository, Update, query:", query)

	sqlPatch.Args["disc_id"] = discId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("discRepository, Update, err:", err.Error())
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

func (repository *discTransaction) UpdateDetail(discId int64, discDetID int64, request entity.UpdateDiscDetBody) error {
	var (
		r            model.DiscDetUpdate
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
	fmt.Println(sqlSetFields)
	query := `UPDATE mst.m_disc_det
			  SET ` + sqlSetFields + `
			  WHERE  disc_det_id = :disc_det_id
			  AND disc_id=:disc_id;`

	// log.Println("tprRepository, Update, query:", query)
	sqlPatch.Args["disc_det_id"] = discDetID
	sqlPatch.Args["disc_id"] = discId

	result, err := repository.tx.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("discRepository, Update, err:", err.Error())
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

func (repository *discTransaction) DeleteDetailNotIn(discId int64, discDetId []int64) error {
	query, args, err := sqlx.In("DELETE FROM mst.m_disc_det	WHERE disc_id = ? AND disc_det_id not in (?);", discId, discDetId)
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
func (repository *discRepositoryImpl) Delete(custId string, discId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_disc
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND disc_id = :disc_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"disc_id":    discId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("DiscRepository, Delete, err:", err.Error())
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

func (repo *discTransaction) InsertDiscount(disc *model.Disc) (int, error) {
	query :=
		`INSERT INTO mst.m_disc(
			cust_id, disc_code, disc_name, 
			start_date, end_date, range_type, 
			is_multiple, purchase_limit, disc_type, 
			disc_perc, disc_value, is_active, 
			created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, 
			deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6,
			$7, $8, $9, 
			$10, $11, $12, 
			$13, $14, $15, 
			$16, $17, $18, 
			$19
		) RETURNING disc_id;`
	lastInsertId := 0
	err := repo.tx.QueryRow(query,
		disc.CustId, disc.DiscCode, disc.DiscName,
		disc.StartDate, disc.EndDate, disc.RangeType,
		disc.IsMultiple, disc.PurchaseLimit, disc.DiscType,
		disc.DiscPerc, disc.DiscValue, disc.IsActive,
		disc.CreatedBy, disc.CreatedAt, disc.UpdatedBy,
		disc.UpdatedAt, disc.IsDel, disc.DeletedBy,
		disc.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("discRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}
