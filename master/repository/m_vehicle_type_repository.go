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

type VehicleTypeRepository interface {
	FindOneByVehicleTypeIdAndCustId(vehicleTypeId int, custId string) (model.VehicleType, error)
	FindOneByVehicleTypeNameAndCustId(vehicleTypeName string, custId string) (model.VehicleType, error)
	FindAllByCustId(dataFilter entity.VehicletypeQueryFilter, custId string) (consPro []model.VehicleType, total int, lastPage int, err error)
	FindAllByCustIdLookupMode(dataFilter entity.VehicletypeQueryFilter, custId string) (consPro []model.VehicleType, total int, lastPage int, err error)
	Store(brand model.VehicleType) (int, error)
	Update(vehicleTypeId int, request entity.UpdateVehicletypeRequest) error
	Delete(custId string, vehicleTypeId int, deletedBy int64) error
}

func NewVehicleTypeRepository(db *sqlx.DB) VehicleTypeRepository {
	return &VehicleTypeRepositoryImpl{db}
}

type VehicleTypeRepositoryImpl struct {
	*sqlx.DB
}

func (repository *VehicleTypeRepositoryImpl) FindOneByVehicleTypeIdAndCustId(vehicleTypeId int, custId string) (model.VehicleType, error) {
	vehicleType := model.VehicleType{}
	query := `SELECT 
				vt.vehicle_type_id, 
				vt.vehicle_type_name, 
				vt.is_active, vt.created_by,
				vt.created_at, vt.updated_by, vt.updated_at,
				vt.is_del, vt.deleted_by, vt.deleted_at
			  FROM mst.m_vehicle_type vt
			  WHERE vt.is_del = false 
				AND vt.vehicle_type_id = $1 
				AND vt.cust_id = $2`
	err := repository.Get(&vehicleType, query, vehicleTypeId, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByBrandIdAndCustId, err:", err.Error())
		return vehicleType, err
	}

	return vehicleType, nil
}

func (repository *VehicleTypeRepositoryImpl) FindOneByVehicleTypeNameAndCustId(vehicleTypeName string, custId string) (model.VehicleType, error) {
	vehicleType := model.VehicleType{}
	query := `SELECT 
					vt.vehicle_type_id, 
					vt.vehicle_type_name, 
					vt.is_active, vt.created_by,
					vt.created_at, vt.updated_by, vt.updated_at,
					vt.is_del, vt.deleted_by, vt.deleted_at
				FROM mst.m_vehicle_type vt
				WHERE vt.is_del = false 
					AND vt.vehicle_type_name = $1 
					AND vt.cust_id = $2`
	err := repository.Get(&vehicleType, query, vehicleTypeName, custId)
	if err != nil {
		log.Println("bankRepository, FindOneByvehicleTypeNameAndCustId, err:", err.Error())
		return vehicleType, err
	}

	return vehicleType, nil
}

func (repository *VehicleTypeRepositoryImpl) FindAllByCustIdLookupMode(dataFilter entity.VehicletypeQueryFilter, custId string) ([]model.VehicleType, int, int, error) {

	vehicleType := []model.VehicleType{}
	selectCount := ` COUNT(vt.*) AS total `
	selectField := `vt.vehicle_type_id,
					vt.vehicle_type_name  `
	qWhere := ` WHERE vt.is_del = false AND vt.is_active = true 
				AND vt.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (vt.vehicle_type_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.VehicleTypeId != 0 {
		qWhere += `AND vt.vehicle_type_id = ` + strconv.Itoa(dataFilter.VehicleTypeId) + ` `
	}

	qFrom := ` FROM mst.m_vehicle_type vt `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + `  ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("vehicleTypeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("vehicleTypeRepository, count total, err:", err.Error())
		return vehicleType, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`vt.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `vt.vehicle_type_id`
		querySelect += fmt.Sprintf(`ORDER BY %s DESC`, sortBy)
	}

	// log.Println("vehicleTypeRepository, querySelect:", querySelect)
	err = repository.Select(&vehicleType, querySelect)
	if err != nil {
		log.Println("vehicleTypeRepository, FindAllByCustIdLookupMode, err:", err.Error())
		return vehicleType, total, 1, err
	}

	return vehicleType, total, 1, nil
}

func (repository *VehicleTypeRepositoryImpl) FindAllByCustId(dataFilter entity.VehicletypeQueryFilter, custId string) ([]model.VehicleType, int, int, error) {

	brands := []model.VehicleType{}
	selectCount := ` COUNT(vt.*) AS total `
	selectField := `vt.vehicle_type_id,
					vt.vehicle_type_name, 
					vt.is_active, vt.created_by,
					vt.created_at, vt.updated_by, vt.updated_at,
					vt.is_del, vt.deleted_by, vt.deleted_at,
					u.user_fullname AS updated_by_name `
	qWhere := ` WHERE vt.is_del = false 
				AND vt.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (vt.vehicle_type_name ILIKE '%` + dataFilter.Query + `%')`
	}

	if dataFilter.IsActive != nil {
		fmt.Println("dataFilter.IsActive:", *dataFilter.IsActive)
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND vt.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND vt.is_active = false `
		}
	}

	if dataFilter.VehicleTypeId != 0 {
		qWhere += `AND vt.vehicle_type_id = ` + strconv.Itoa(dataFilter.VehicleTypeId) + ` `
	}

	qFrom := ` FROM mst.m_vehicle_type vt
			   LEFT JOIN sys.m_user u ON u.user_id = vt.updated_by `

	queryCount :=
		`SELECT ` + selectCount + ` ` + qFrom + ` ` + qWhere
	querySelect :=
		`SELECT ` + selectField + ` ` + qFrom + ` ` + qWhere

	// log.Println("vehicleTypeRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("vehicleTypeRepository, count total, err:", err.Error())
		return brands, 0, 0, err
	}

	sortBy := `` // default sort by
	if dataFilter.Sort != "" {
		mSortBy := strings.Split(dataFilter.Sort, ",")
		for _, row := range mSortBy {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortBy += fmt.Sprintf(`vt.%s %s, `, colSort[0], colSort[1])
			}
		}
		sortBy = strings.TrimSuffix(sortBy, ", ")
		querySelect += fmt.Sprintf(`ORDER BY %s`, sortBy)
	} else {
		sortBy := `vt.vehicle_type_id`
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

	// log.Println("vehicleTypeRepository, querySelect:", querySelect)
	err = repository.Select(&brands, querySelect)
	if err != nil {
		log.Println("vehicleTypeRepository, FindAllByCustId, err:", err.Error())
		return brands, total, lastPage, err
	}

	return brands, total, lastPage, nil
}

func (repository *VehicleTypeRepositoryImpl) Store(vehicleType model.VehicleType) (int, error) {
	query :=
		`INSERT INTO mst.m_vehicle_type(
			cust_id, vehicle_type_name, 
			is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at)
		VALUES ( 
			$1, $2, $3, 
			$4, $5, $6, $7, 
			$8,	$9, $10 
		) RETURNING vehicle_type_id;`
	lastInsertId := 0
	err := repository.QueryRow(query,
		vehicleType.CustId, vehicleType.VehicleTypeName, vehicleType.IsActive,
		vehicleType.CreatedBy, vehicleType.CreatedAt, vehicleType.UpdatedBy, vehicleType.UpdatedAt,
		vehicleType.IsDel, vehicleType.DeletedBy, vehicleType.DeletedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("vehicleTypeRepository, Store, err:", err.Error())
		return lastInsertId, err
	}
	return lastInsertId, nil
}

func (repository *VehicleTypeRepositoryImpl) Update(vehicleTypeId int, request entity.UpdateVehicletypeRequest) error {
	var (
		r            model.VehicleTypeUpdate
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

	query := `UPDATE mst.m_vehicle_type
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND vehicle_type_id = :vehicle_type_id;`

	sqlPatch.Args["vehicle_type_id"] = vehicleTypeId
	sqlPatch.Args["cust_id"] = request.CustId

	result, err := repository.NamedExec(query, sqlPatch.Args)
	if err != nil {
		log.Println("vehicleTypeRepository, Update, err:", err.Error())
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

func (repository *VehicleTypeRepositoryImpl) Delete(custId string, vehicleTypeId int, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_vehicle_type
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND vehicle_type_id = :vehicle_type_id;`

	wMap := map[string]interface{}{
		"cust_id":         custId,
		"vehicle_type_id": vehicleTypeId,
		"deleted_by":      deletedBy,
	}

	result, err := repository.NamedExec(query, wMap)
	if err != nil {
		log.Println("VehicleTypeRepository, Delete, err:", err.Error())
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
