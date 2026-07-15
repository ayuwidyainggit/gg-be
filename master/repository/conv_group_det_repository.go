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

type ConvGroupDetRepository interface {
	FindOneByConvGroupIdAndCustId(convGrpId int, custId string, proId int) (model.ConvGroupDetRead, error)
	FindOneByProIdAndCustId(convGrpId int, proId int, custId string) (model.ConvGroupDetRead, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.ConvGroupDetRead, total int, lastPage int, err error)
	Store(convGroupDet model.ConvGroupDet) (int, error)
	Update(convGroupId int, proId int, request entity.UpdateConvGroupDetRequest) error
	Delete(custId string, convGroupId int, proId int, deletedBy int64) error
}

func NewConvGroupDetRepository(db *sqlx.DB) ConvGroupDetRepository {
	return &convGroupDetRepositoryImpl{db}
}

type convGroupDetRepositoryImpl struct {
	*sqlx.DB
}

func (repository *convGroupDetRepositoryImpl) FindOneByConvGroupIdAndCustId(convGrpId int, custId string, proId int) (model.ConvGroupDetRead, error) {
	convGroupDet := model.ConvGroupDetRead{}
	query :=
		`SELECT cg.*,p.pro_code, p.pro_name 
	FROM mst.m_conv_group_det cg 
	LEFT JOIN mst.m_product p ON p.pro_id = cg.pro_id 
	WHERE conv_grp_id = $1
	AND cg.cust_id = $2 
	AND cg.pro_id = $3 `
	err := repository.Get(&convGroupDet, query, convGrpId, custId, proId)

	if err != nil {
		log.Println("convGroupDetRepository, FindOneByConvGroupIdAndCustId, err:", err.Error())
		return convGroupDet, err
	}

	return convGroupDet, nil
}

func (repository *convGroupDetRepositoryImpl) FindOneByProIdAndCustId(convGrpId int, proId int, custId string) (model.ConvGroupDetRead, error) {
	convGroup := model.ConvGroupDetRead{}
	query :=
		`SELECT cg.*,p.pro_code, p.pro_name
	FROM mst.m_conv_group_det cg 
	LEFT JOIN mst.m_product p ON p.pro_id = cg.pro_id
	WHERE conv_grp_id = $1 
	AND cg.pro_id = $2 
	AND cg.cust_id = $3`
	err := repository.Get(&convGroup, query, convGrpId, proId, custId)
	if err != nil {
		log.Println("ConvGroupDetRepository, FindOneByProIdAndCustId, err:", err.Error())
		return convGroup, err
	}

	return convGroup, nil
}

func (repository *convGroupDetRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.ConvGroupDetRead, int, int, error) {

	convGroups := []model.ConvGroupDetRead{}
	selectCount := ` COUNT(*) AS total `
	selectField := `cg.*,
					u.user_fullname AS updated_by_name, p.pro_code, p.pro_name `
	qWhere := ` WHERE cg.is_del = false 
				AND cg.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (cg.unit_id1 ILIKE '%` + dataFilter.Query + `%' 
					OR cg.unit_id2 ILIKE '%` + dataFilter.Query + `%' 
					OR cg.unit_id3 ILIKE '%` + dataFilter.Query + `%' )`
	}

	qFrom := ` FROM mst.m_conv_group_det cg
	LEFT JOIN sys.m_user u ON u.user_id = cg.updated_by
	LEFT JOIN mst.m_product p ON p.pro_id = cg.pro_id `

	queryCount := `SELECT ` + selectCount + qFrom + qWhere
	querySelect := `SELECT ` + selectField + qFrom + qWhere

	// log.Println("convGroupDetRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("convGroupDetRepository, count total, err:", err.Error())
		return convGroups, 0, 0, err
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
		sortBy := `cg.pro_id`
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

	// log.Println("convGroupDetRepository, querySelect:", querySelect)
	err = repository.Select(&convGroups, querySelect)
	if err != nil {
		log.Println("convGroupDetRepository, FindAllByCustId, err:", err.Error())
		return convGroups, total, lastPage, err
	}

	return convGroups, total, lastPage, nil
}

func (repository *convGroupDetRepositoryImpl) Store(convGroupDet model.ConvGroupDet) (int, error) {
	query :=
		`INSERT INTO mst.m_conv_group_det(
			cust_id, conv_grp_id, pro_id, 
			unit_id1, unit_id2, unit_id3, unit_id4, unit_id5, 
			conv_unit2, conv_unit3, conv_unit4, conv_unit5, 
			new_unit_id1, new_unit_id2, new_unit_id3, new_unit_id4, new_unit_id5, 
			new_conv_unit2, new_conv_unit3, new_conv_unit4, new_conv_unit5,
			created_by, created_at, updated_by, updated_at, is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, $8, 
			$9, $10, $11, $12, 
			$13, $14, $15, $16, $17, 
			$18, $19, $20, $21,
			$22, $23, $24, $25, $26, $27, $28
		) RETURNING conv_grp_id;`
	lastInsertId := convGroupDet.ConvGrpId
	err := repository.QueryRow(query,
		convGroupDet.CustId, convGroupDet.ConvGrpId, convGroupDet.ProId,
		convGroupDet.UnitId1, convGroupDet.UnitId2, convGroupDet.UnitId3, convGroupDet.UnitId4, convGroupDet.UnitId5,
		convGroupDet.ConvUnit2, convGroupDet.ConvUnit3, convGroupDet.ConvUnit4, convGroupDet.ConvUnit5,
		convGroupDet.NewUnitId1, convGroupDet.NewUnitId2, convGroupDet.NewUnitId3, convGroupDet.NewUnitId4, convGroupDet.NewUnitId5,
		convGroupDet.NewConvUnit2, convGroupDet.NewConvUnit3, convGroupDet.NewConvUnit4, convGroupDet.NewConvUnit5,
		convGroupDet.CreatedBy, convGroupDet.CreatedAt, convGroupDet.UpdatedBy, convGroupDet.UpdatedAt, convGroupDet.IsDel, convGroupDet.DeletedBy, convGroupDet.DeletedAt).Scan(&lastInsertId)

	if err != nil {
		log.Println("convGroupDetRepository, Store, err:", err.Error())
		return convGroupDet.ConvGrpId, err
	}
	return convGroupDet.ConvGrpId, nil
}

func (repository *convGroupDetRepositoryImpl) Update(convGroupId int, proId int, request entity.UpdateConvGroupDetRequest) error {
	var (
		r            model.ConvGroupDetUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("convGroupDetRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_conv_group_det
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND conv_grp_id = :conv_grp_id_old AND pro_id =:pro_id_old ;`

	// log.Println("convGroupDetRepository, Update, query:", query)

	sqlPatch.Args["conv_grp_id_old"] = convGroupId
	sqlPatch.Args["pro_id_old"] = proId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("convGroupDetRepository, Update, err:", err.Error())
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

func (repository *convGroupDetRepositoryImpl) Delete(custId string, convGroupId int, proId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_conv_group_det
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			AND cust_id = :cust_id
			AND conv_grp_id = :conv_grp_id;`

	wMap := map[string]interface{}{
		"cust_id":     custId,
		"conv_grp_id": convGroupId,
		"pro_id":      proId,
		"deleted_by":  deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("ConvGroupDetRepository, Delete, err:", err.Error())
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
