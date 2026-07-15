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

type BeatRepository interface {
	FindOneByBeatIdAndCustId(beatId int, custId string) (model.Beat, error)
	FindOneByBeatCodeAndCustId(beatCode string, custId string) (model.Beat, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Beat, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Beat, total int, lastPage int, err error)
	Store(beat model.Beat) (int, error)
	Update(beatId int, request entity.UpdateBeatRequest) error
	Delete(custId string, beatId int, deletedBy int64) error
}

func NewBeatRepository(db *sqlx.DB) BeatRepository {
	return &beatRepositoryImpl{db}
}

type beatRepositoryImpl struct {
	*sqlx.DB
}

func (repository *beatRepositoryImpl) FindOneByBeatIdAndCustId(beatId int, custId string) (model.Beat, error) {
	beat := model.Beat{}
	query := `SELECT 
				a.cust_id, a.beat_id, a.beat_code,
				a.beat_name, a.district_id, a.is_active, a.created_by,
				a.created_at, a.updated_by, a.updated_at,
				a.is_del, a.deleted_by, a.deleted_at, b.district_code, b.district_name
			  FROM mst.m_beat a
			  LEFT JOIN mst.m_district b on b.district_id = a.district_id AND b.cust_id = '`+ custId +`'
			  WHERE a.beat_id = $1 
			  AND a.cust_id = $2`
	err := repository.Get(&beat, query, beatId, custId)
	if err != nil {
		log.Println("beatRepository, FindOneByBeatCodeAndCustId, err:", err.Error())
		return beat, err
	}

	return beat, nil
}

func (repository *beatRepositoryImpl) FindOneByBeatCodeAndCustId(beatCode string, custId string) (model.Beat, error) {
	beat := model.Beat{}
	query := `SELECT 
				cust_id, beat_id, beat_code,
				beat_name, district_id, is_active, 
				created_by, created_at, updated_by, 
				updated_at, is_del, deleted_by, deleted_at
			  FROM mst.m_beat 
			  WHERE beat_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&beat, query, beatCode, custId)
	if err != nil {
		log.Println("beatRepository, FindOneByBeatCodeAndCustId, err:", err.Error())
		return beat, err
	}

	return beat, nil
}

func (repository *beatRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Beat, int, int, error) {

	beats := []model.Beat{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.beat_id, a.beat_code, a.beat_name, 
					a.district_id, a.is_active, a.created_by,
					a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, 
					b.district_code, b.district_name, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false 
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.beat_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.beat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_beat a 
				LEFT JOIN mst.m_district b on b.district_id = a.district_id AND b.cust_id = '`+ custId +`' 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("beatRepository, count total, err:", err.Error())
		return beats, 0, 0, err
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
		sortBy := `a.beat_id`
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

	err = repository.Select(&beats, querySelect)
	if err != nil {
		log.Println("beatRepository, FindAllByCustId, err:", err.Error())
		return beats, total, lastPage, err
	}

	return beats, total, lastPage, nil
}

func (repository *beatRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Beat, int, int, error) {

	beats := []model.Beat{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.beat_id, a.beat_code, a.beat_name, 
					a.district_id, b.district_code, b.district_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true 
				AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.beat_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.beat_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_beat a 
				LEFT JOIN mst.m_district b on b.district_id = a.district_id AND b.cust_id = '`+ custId +`'
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("beatRepository, count total, err:", err.Error())
		return beats, 0, 0, err
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
		sortBy := `a.beat_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&beats, querySelect)
	if err != nil {
		log.Println("beatRepository, FindAllByCustId, err:", err.Error())
		return beats, total, 1, err
	}

	return beats, total, 1, nil
}

func (repository *beatRepositoryImpl) Store(beat model.Beat) (int, error) {
	query :=
		`INSERT INTO mst.m_beat(
			cust_id, beat_code, beat_name, district_id, 
			is_active, created_by, created_at, updated_by, 
			updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11, $12
		) RETURNING beat_id;`
	lastInsertId := beat.BeatId
	err := repository.QueryRow(query,
		beat.CustId, beat.BeatCode, beat.BeatName, beat.DistrictId,
		beat.IsActive, beat.CreatedBy, beat.CreatedAt, beat.UpdatedBy,
		beat.UpdatedAt, beat.IsDel, beat.DeletedBy, beat.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("beatRepository, Store, err:", err.Error())
		return beat.BeatId, err
	}
	return beat.BeatId, nil
}

func (repository *beatRepositoryImpl) Update(beatId int, request entity.UpdateBeatRequest) error {
	var (
		r            model.BeatUpdate
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

	query := `UPDATE mst.m_beat
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND beat_id = :beat_id_old;`

	log.Println("beatRepository, Update, query:", query)

	sqlPatch.Args["beat_id_old"] = beatId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("beatRepository, Update, err:", err.Error())
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

func (repository *beatRepositoryImpl) Delete(custId string, beatId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_beat
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND beat_id = :beat_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"beat_id":    beatId,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("BeatRepository, Delete, err:", err.Error())
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
