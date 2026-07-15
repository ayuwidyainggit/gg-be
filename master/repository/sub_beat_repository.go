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

type SubBeatRepository interface {
	FindOneBySubBeatIdAndCustId(subBeatId int, custId string) (model.SubBeat, error)
	FindOneBySubBeatCodeAndCustId(subBeatCode string, custId string) (model.SubBeat, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (subBeat []model.SubBeat, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (subBeat []model.SubBeat, total int, lastPage int, err error)
	Store(subBeat model.SubBeat) (int, error)
	Update(subBeatId int, request entity.UpdateSubBeatRequest) error
	Delete(custId string, subBeatId int, deletedBy int64) error
}

func NewSubBeatRepository(db *sqlx.DB) SubBeatRepository {
	return &subBeatRepositoryImpl{db}
}

type subBeatRepositoryImpl struct {
	*sqlx.DB
}

func (repository *subBeatRepositoryImpl) FindOneBySubBeatIdAndCustId(subBeatId int, custId string) (model.SubBeat, error) {
	subBeat := model.SubBeat{}
	query := `SELECT 
				a.cust_id, a.sbeat_id, a.sbeat_code,
				a.sbeat_name, a.beat_id, a.district_id, a.is_active, a.created_by,
				a.created_at, a.updated_by, a.updated_at,
				a.is_del, a.deleted_by, a.deleted_at, b.district_code, b.district_name,
				c.beat_code, c.beat_name  
			  FROM mst.m_sub_beat a 
			  LEFT JOIN mst.m_district b on b.district_id = a.district_id
			  LEFT JOIN mst.m_beat c on c.beat_id = a.beat_id
			  WHERE a.sbeat_id = $1 
			  AND a.cust_id = $2`
	err := repository.Get(&subBeat, query, subBeatId, custId)
	if err != nil {
		log.Println("subBeatRepository, FindOneBySubBeatCodeAndCustId, err:", err.Error())
		return subBeat, err
	}

	return subBeat, nil
}

func (repository *subBeatRepositoryImpl) FindOneBySubBeatCodeAndCustId(subBeatCode string, custId string) (model.SubBeat, error) {
	subBeat := model.SubBeat{}
	query := `SELECT 
				cust_id, sbeat_id, sbeat_code,
				sbeat_name, beat_id, district_id, is_active, 
				created_by, created_at, updated_by, 
				updated_at, is_del, deleted_by, deleted_at
			  FROM mst.m_sub_beat 
			  WHERE sbeat_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&subBeat, query, subBeatCode, custId)
	if err != nil {
		log.Println("subBeatRepository, FindOneBySubBeatCodeAndCustId, err:", err.Error())
		return subBeat, err
	}

	return subBeat, nil
}

func (repository *subBeatRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SubBeat, int, int, error) {

	subBeats := []model.SubBeat{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sbeat_id, a.sbeat_code, a.sbeat_name, a.beat_id, 
					a.district_id, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, b.district_code, b.district_name,
					c.beat_code, c.beat_name, u.user_name AS updated_by_name`
	qWhere := ` WHERE a.is_del = false 
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sbeat_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sbeat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_sub_beat a 
				LEFT JOIN mst.m_district b on b.district_id = a.district_id
				LEFT JOIN mst.m_beat c on c.beat_id = a.beat_id 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("subBeatRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBeatRepository, count total, err:", err.Error())
		return subBeats, 0, 0, err
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
		sortBy := `a.sbeat_id`
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

	err = repository.Select(&subBeats, querySelect)
	if err != nil {
		log.Println("subBeatRepository, FindAllByCustId, err:", err.Error())
		return subBeats, total, lastPage, err
	}

	return subBeats, total, lastPage, nil
}

func (repository *subBeatRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.SubBeat, int, int, error) {

	subBeats := []model.SubBeat{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.sbeat_id, a.sbeat_code, a.sbeat_name, a.beat_id, 
					a.district_id, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, b.district_code, b.district_name,
					c.beat_code, c.beat_name, u.user_name AS updated_by_name`
	qWhere := ` WHERE a.is_del = false AND a.is_active = true
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.sbeat_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.sbeat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_sub_beat a 
				LEFT JOIN mst.m_district b on b.district_id = a.district_id
				LEFT JOIN mst.m_beat c on c.beat_id = a.beat_id 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by`
	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("subBeatRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("subBeatRepository, count total, err:", err.Error())
		return subBeats, 0, 0, err
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
		sortBy := `a.sbeat_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&subBeats, querySelect)
	if err != nil {
		log.Println("subBeatRepository, FindAllByCustId, err:", err.Error())
		return subBeats, total, 1, err
	}

	return subBeats, total, 1, nil
}

func (repository *subBeatRepositoryImpl) Store(subBeat model.SubBeat) (int, error) {
	query :=
		`INSERT INTO mst.m_sub_beat(
			cust_id, sbeat_code, sbeat_name, beat_id, 
			district_id, is_active, created_by, created_at, 
			updated_by, updated_at, is_del, deleted_by, 
			deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12,
			$13
		) RETURNING sbeat_id;`
	lastInsertId := subBeat.SbeatId
	err := repository.QueryRow(query,
		subBeat.CustId, subBeat.SbeatCode, subBeat.SbeatName, subBeat.BeatId,
		subBeat.DistrictId, subBeat.IsActive, subBeat.CreatedBy, subBeat.CreatedAt,
		subBeat.UpdatedBy, subBeat.UpdatedAt, subBeat.IsDel, subBeat.DeletedBy,
		subBeat.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("subBeatRepository, Store, err:", err.Error())
		return subBeat.SbeatId, err
	}
	return subBeat.SbeatId, nil
}

func (repository *subBeatRepositoryImpl) Update(subBeatId int, request entity.UpdateSubBeatRequest) error {
	var (
		r            model.SubBeatUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	data, _ := json.Marshal(sqlPatch)
	fmt.Printf("subBeatRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_sub_beat
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND sbeat_id = :sbeat_id_old;`

	sqlPatch.Args["sbeat_id_old"] = subBeatId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("subBeatRepository, Update, err:", err.Error())
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

func (repository *subBeatRepositoryImpl) Delete(custId string, subBeatId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_sub_beat
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND sbeat_id = :sbeat_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"sbeat_id":   subBeatId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("SubBeatRepository, Delete, err:", err.Error())
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
