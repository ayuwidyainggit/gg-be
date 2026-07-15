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

type MChannelRepository interface {
	FindOneByChannelCodeAndCustId(channelCode, custId string) (model.MChannel, error)
	Store(channel model.MChannel) (int, error)
	FindAllByCustId(dataFilter entity.ChannelQueryFilter, custId string) ([]model.MChannel, int, int, error)
	FindAllByCustIdLookupMode(dataFilter entity.ChannelQueryFilter, custId string) ([]model.MChannel, int, int, error)
	FindOneByChannelIdAndCustId(channelId int, custId string) (model.MChannel, error)
	Update(channelId int, request entity.ChannelUpdateRequest) error
	Delete(custId string, channelID int, deletedBy int64) error
}

type MChannelRepositoryImpl struct {
	*sqlx.DB
}

func NewMChannelRepository(db *sqlx.DB) *MChannelRepositoryImpl {
	return &MChannelRepositoryImpl{db}
}

func (repository *MChannelRepositoryImpl) FindOneByChannelCodeAndCustId(channelCode, custId string) (model.MChannel, error) {
	brand := model.MChannel{}
	query := `SELECT 
				*
			  FROM mst.m_channel 
			  WHERE is_del = false 
			  AND channel_code = $1 
			  AND cust_id = $2`
	err := repository.Get(&brand, query, channelCode, custId)
	if err != nil {
		log.Println("MChannelRepository, FindOneByChannelCodeAndCustId, err:", err.Error())
		return brand, err
	}

	return brand, nil
}

func (repository *MChannelRepositoryImpl) Store(channel model.MChannel) (int, error) {
	query :=
		`INSERT INTO mst.m_channel(
			cust_id, channel_code, channel_name, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10, $11
		) RETURNING channel_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		channel.CustID, channel.ChannelCode, channel.ChannelName,
		channel.IsActive, channel.CreatedBy, channel.CreatedAt, channel.UpdatedBy, channel.UpdatedAt,
		channel.IsDel, channel.DeletedBy, channel.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("channelRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *MChannelRepositoryImpl) FindAllByCustId(dataFilter entity.ChannelQueryFilter, custId string) ([]model.MChannel, int, int, error) {

	channels := []model.MChannel{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.*, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.channel_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.channel_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}

	qFrom := ` 	FROM mst.m_channel a
	LEFT JOIN sys.m_user u ON u.user_id = a.updated_by `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere
	fmt.Println(querySelect)
	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("MChannelRepository, count total, err:", err.Error())
		return channels, 0, 0, err
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
		sortBy := `a.channel_id`
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

	// log.Println("districtRepository, querySelect:", querySelect)
	err = repository.Select(&channels, querySelect)
	if err != nil {
		log.Println("MChannelRepository, FindAllByCustId, err:", err.Error())
		return channels, total, lastPage, err
	}

	return channels, total, lastPage, nil
}

func (repository *MChannelRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.ChannelQueryFilter, custId string) ([]model.MChannel, int, int, error) {

	districts := []model.MChannel{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.channel_id, a.channel_code, a.channel_name `
	qWhere := ` WHERE a.is_del = false AND a.is_active = true AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.district_code ILIKE '%` + dataFilter.Query + `%' 
					OR a.district_name ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` 	FROM mst.m_channel a 
				LEFT JOIN sys.m_user u ON u.user_id = a.updated_by   `

	queryCount := `SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect := `SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("districtRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("districtRepository, count total, err:", err.Error())
		return districts, 0, 0, err
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
		sortBy := `a.channel_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	err = repository.Select(&districts, querySelect)
	if err != nil {
		log.Println("districtRepository, FindAllByCustId, err:", err.Error())
		return districts, total, 1, err
	}

	return districts, total, 1, nil
}

func (repository *MChannelRepositoryImpl) FindOneByChannelIdAndCustId(channelId int, custId string) (model.MChannel, error) {
	channel := model.MChannel{}
	query := `SELECT 
				b.*, u.user_fullname AS updated_by_name
			  FROM mst.m_channel b
			  LEFT JOIN sys.m_user u ON u.user_id = b.updated_by 
			  WHERE b.is_del = false 
				AND b.channel_id = $1 
				AND b.cust_id = $2`
	err := repository.Get(&channel, query, channelId, custId)
	if err != nil {
		log.Println("MChannelRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return channel, err
	}

	return channel, nil
}

func (repository *MChannelRepositoryImpl) Update(channelId int, request entity.ChannelUpdateRequest) error {
	var (
		r            model.MChannelUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("districtRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_channel
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND channel_id = :channel_id_old;`

	log.Println("channelRepository, Update, query:", query)

	sqlPatch.Args["channel_id_old"] = channelId
	sqlPatch.Args["cust_id"] = request.CustID

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("channelRepository, Update, err:", err.Error())
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

func (repository *MChannelRepositoryImpl) Delete(custId string, channelID int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_channel
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND channel_id = :channel_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"channel_id": channelID,
		"deleted_by": deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("channelRepository, Delete, err:", err.Error())
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
