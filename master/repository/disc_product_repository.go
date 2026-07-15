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

type DiscProductRepository interface {
	FindOneByDiscIdAndCustId(discId int, custId string) (model.DiscProductRead, error)
	FindOneByDiscProductIdAndCustId(discId int, productId int, custId string) (model.DiscProductRead, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.DiscProductRead, total int, lastPage int, err error)
	Store(discProduct model.DiscProduct) (int, error)
	Update(discId int, proId int, request entity.UpdateDiscProductRequest) error
	Delete(custId string, discId int, proId int, deletedBy int64) error
}

func NewDiscProductRepository(db *sqlx.DB) DiscProductRepository {
	return &discProductRepositoryImpl{db}
}

type discProductRepositoryImpl struct {
	*sqlx.DB
}

func (repository *discProductRepositoryImpl) FindOneByDiscIdAndCustId(discId int, custId string) (model.DiscProductRead, error) {
	disProduct := model.DiscProductRead{}
	query := `SELECT dp.*,p.pro_code, p.pro_name, d.disc_code, d.disc_name
			FROM mst.m_disc_product dp 
			LEFT JOIN mst.m_disc d ON d.disc_id = dp.disc_id AND d.cust_id = $2
			LEFT JOIN mst.m_product p ON p.pro_id = dp.pro_id 
			  WHERE dp.outlet_id = $1 
			  AND dp.cust_id = $2`
	err := repository.Get(&disProduct, query, discId, custId)
	if err != nil {
		log.Println("discProductRepository, FindOneByDiscIdAndCustId, err:", err.Error())
		return disProduct, err
	}

	return disProduct, nil
}

func (repository *discProductRepositoryImpl) FindOneByDiscProductIdAndCustId(discId int, productId int, custId string) (model.DiscProductRead, error) {
	discProduct := model.DiscProductRead{}
	query := `SELECT dp.*,p.pro_code, p.pro_name, d.disc_code, d.disc_name
	FROM mst.m_disc_product dp 
	LEFT JOIN mst.m_disc d ON d.disc_id = dp.disc_id AND d.cust_id = $3
	LEFT JOIN mst.m_product p ON p.pro_id = dp.pro_id
	WHERE dp.disc_id = $1 AND dp.pro_id = $2 AND dp.cust_id = $3`
	err := repository.Get(&discProduct, query, discId, productId, custId)
	if err != nil {
		log.Println("discProductRepository, FindOneByDiscProductIdAndCustId, err:", err.Error())
		return discProduct, err
	}

	return discProduct, nil
}

func (repository *discProductRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.DiscProductRead, int, int, error) {
	discProducts := []model.DiscProductRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` dp.*, p.pro_code, p.pro_name, d.disc_code, d.disc_name,
					u.user_name AS updated_by_name, p.pro_code, p.pro_name `
	qWhere := ` WHERE dp.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (p.pro_code ILIKE '%` + dataFilter.Query + `%' 
					OR p.pro_name ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND dp.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND dp.is_active = false `
		}
	}
	qFrom := ` FROM mst.m_disc_product dp
	LEFT JOIN sys.m_user u ON u.user_id = dp.updated_by
	LEFT JOIN mst.m_disc d ON d.disc_id = dp.disc_id AND d.cust_id = '` + custId + `' 
	LEFT JOIN mst.m_product p ON p.pro_id = dp.pro_id `
	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("discProductRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("discProductRepository, count total, err:", err.Error())
		return discProducts, 0, 0, err
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

	// log.Println("discProductRepository, querySelect:", querySelect)
	err = repository.Select(&discProducts, querySelect)
	if err != nil {
		log.Println("discProductRepository, FindAllByCustId, err:", err.Error())
		return discProducts, total, lastPage, err
	}

	return discProducts, total, lastPage, nil
}

func (repository *discProductRepositoryImpl) Store(discProduct model.DiscProduct) (int, error) {
	query :=
		`INSERT INTO mst.m_disc_product(cust_id, disc_id, pro_id, min_qty, min_qty_str, max_qty, max_qty_str, disc_perc, 
				is_active, created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING disc_id;`
	lastInsertId := discProduct.DiscId
	err := repository.QueryRow(query,
		discProduct.CustId, discProduct.DiscId, discProduct.ProId,
		discProduct.MinQty, discProduct.MinQtyStr, discProduct.MaxQty, discProduct.MaxQtyStr, discProduct.DiscPerc, discProduct.IsActive,
		discProduct.CreatedBy, discProduct.CreatedAt, discProduct.UpdatedBy, discProduct.UpdatedAt, discProduct.IsDel, discProduct.DeletedBy,
		discProduct.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("discProductRepository, Store, err:", err.Error())
		return discProduct.DiscId, err
	}
	return discProduct.DiscId, nil
}

func (repository *discProductRepositoryImpl) Update(discId int, proId int, request entity.UpdateDiscProductRequest) error {
	var (
		r            model.DiscProductUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("discProductRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_disc_product
			  SET ` + sqlSetFields + `
			  WHERE is_del = false AND cust_id = :cust_id AND disc_id = :disc_id_old AND pro_id = :pro_id_old;`

	log.Println("discProductRepository, Update, query:", query)

	sqlPatch.Args["disc_id_old"] = discId
	sqlPatch.Args["pro_id_old"] = proId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("discProductRepository, Update, err:", err.Error())
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

func (repository *discProductRepositoryImpl) Delete(custId string, discId int, proId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_disc_product
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND disc_id = :disc_id
			AND pro_id = :pro_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"disc_id":    discId,
		"pro_id":     proId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("DiscProductRepository, Delete, err:", err.Error())
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
