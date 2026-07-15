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

type VehicleRepository interface {
	FindOneParentCustId(distCustId string) (model.MCustomer, error)
	FindOneByVehicleIdAndCustId(vehicleId int64, custId string) (model.Vehicle, error)
	FindOneByVehicleNoAndCustId(vehicleNo string, custId string) (model.Vehicle, error)
	FindOneByVehicleTypeAndCustId(vehicleTypeId int) (model.VehicleType, error)
	FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) (consPro []model.Vehicle, total int, lastPage int, err error)
	Store(vehicle model.Vehicle) (int64, error)
	Update(vehicleId int64, request entity.UpdateVehicleRequest) error
	Delete(custId string, vehicleId int64, deletedBy int64) error
	FindOneByDriverAndCustId(driverID int64, custId string) (model.Vehicle, error)
	FindOneByHelperAndCustId(helperID int64, custId string) (model.Vehicle, error)
}

func NewVehicleRepository(db *sqlx.DB) VehicleRepository {
	return &vehicleRepositoryImpl{db}
}

type vehicleRepositoryImpl struct {
	*sqlx.DB
}

func (repository *vehicleRepositoryImpl) FindOneParentCustId(distCustId string) (model.MCustomer, error) {
	mCustomer := model.MCustomer{}
	query := `SELECT 
				cust_id, cust_name, parent_cust_id
			  FROM smc.m_customer
			  WHERE cust_id = $1`
	err := repository.Get(&mCustomer, query, distCustId)
	if err != nil {
		log.Println("VehicleRepository, FindOneParentCustId, err:", err.Error())
		return mCustomer, err
	}

	return mCustomer, nil
}

func (repository *vehicleRepositoryImpl) FindOneByVehicleIdAndCustId(vehicleId int64, custId string) (model.Vehicle, error) {
	vehicle := model.Vehicle{}
	query := `SELECT 
				a.cust_id, a.vehicle_id, a.vehicle_no, 
				a.vehicle_desc, a.vehicle_type, 
				a.length, a.width, a.height, COALESCE(a.weight, 0) as weight, a.volume,
				a.driver_id, a.helper_id,
				COALESCE(vt.vehicle_type_name, '') as vehicle_type_name,COALESCE(ed.emp_name, '') as driver_name,COALESCE(eh.emp_name, '') as helper_name,
				a.is_active, a.created_by,
				a.created_at, a.updated_by, a.updated_at,
				a.is_del, a.deleted_by, a.deleted_at
			  FROM mst.m_vehicle a
				LEFT JOIN mst.m_vehicle_type vt ON a.vehicle_type = vt.vehicle_type_id
				LEFT JOIN mst.m_employee ed ON ed.emp_id = a.driver_id
				LEFT JOIN mst.m_employee eh ON eh.emp_id = a.helper_id
			  WHERE a.vehicle_id = $1 
			  AND a.cust_id = $2`
	err := repository.Get(&vehicle, query, vehicleId, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByVehicleIdAndCustId, err:", err.Error())
		return vehicle, err
	}

	return vehicle, nil
}

func (repository *vehicleRepositoryImpl) FindOneByVehicleNoAndCustId(vehicleNo string, custId string) (model.Vehicle, error) {
	vehicle := model.Vehicle{}
	query := `SELECT 
				cust_id, vehicle_id, vehicle_no, 
				vehicle_desc, vehicle_type, 
				length, width, height, COALESCE(weight, 0) as weight, volume,
				driver_id, helper_id,
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_vehicle 
			  WHERE vehicle_no = $1 
			  AND cust_id = $2`
	err := repository.Get(&vehicle, query, vehicleNo, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByVehicleIdAndCustId, err:", err.Error())
		return vehicle, err
	}

	return vehicle, nil
}

func (repository *vehicleRepositoryImpl) FindOneByVehicleTypeAndCustId(vehicleTypeId int) (model.VehicleType, error) {
	vehicleType := model.VehicleType{}
	query := `SELECT 
				vt.vehicle_type_id, 
				COALESCE(vt.vehicle_type_name, '') as vehicle_type_name, 
				vt.is_active, vt.created_by,
				vt.created_at, vt.updated_by, vt.updated_at,
				vt.is_del, vt.deleted_by, vt.deleted_at
			FROM mst.m_vehicle_type vt
			WHERE vt.is_del = false 
				AND vt.vehicle_type_id = $1`
	err := repository.Get(&vehicleType, query, vehicleTypeId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByVehicleTypeAndCustId, err:", err.Error())
		return vehicleType, err
	}

	return vehicleType, nil
}

func (repository *vehicleRepositoryImpl) FindAllByCustId(dataFilter entity.GeneralQueryFilter, custId string) ([]model.Vehicle, int, int, error) {

	vehicles := []model.Vehicle{}
	selectCount := ` COUNT(*) AS total `
	selectField := ` a.cust_id, a.vehicle_id, a.vehicle_no, a.vehicle_desc, a.vehicle_type, 
					a.length, a.width, a.height, COALESCE(a.weight, 0) as weight, a.volume, a.driver_id, a.helper_id, 
					COALESCE(vt.vehicle_type_name, '') as vehicle_type_name,COALESCE(ed.emp_name, '') as driver_name,COALESCE(eh.emp_name, '') as helper_name,
					a.is_active, a.created_by, a.created_at, a.updated_by, a.updated_at,
					a.is_del, a.deleted_by, a.deleted_at, u.user_fullname AS updated_by_name `
	qWhere := ` WHERE a.is_del = false AND a.cust_id = '` + custId + `' `

	if dataFilter.Query != "" {
		qWhere += ` AND (a.vehicle_desc ILIKE '%` + dataFilter.Query + `%' 
					OR a.vehicle_no ILIKE '%` + dataFilter.Query + `%'
					OR a.vehicle_type ILIKE '%` + dataFilter.Query + `%' )`
	}
	if dataFilter.IsActive != nil {
		if *dataFilter.IsActive == 1 {
			qWhere += ` AND a.is_active = true `
		}
		if *dataFilter.IsActive == 2 {
			qWhere += ` AND a.is_active = false `
		}
	}
	queryJoin := `
					LEFT JOIN sys.m_user u ON u.user_id = a.updated_by
					LEFT JOIN mst.m_vehicle_type vt ON a.vehicle_type = vt.vehicle_type_id
					LEFT JOIN mst.m_employee ed ON ed.emp_id = a.driver_id
					LEFT JOIN mst.m_employee eh ON eh.emp_id = a.helper_id
				`
	queryCount := `SELECT ` + selectCount + ` FROM mst.m_vehicle a ` + queryJoin + qWhere
	querySelect := `SELECT ` + selectField + ` FROM mst.m_vehicle a ` + queryJoin + qWhere

	// log.Println("vehicleRepository, queryCount:", queryCount)
	var total int
	err := repository.QueryRow(queryCount).Scan(&total)
	if err != nil {
		log.Println("vehicleRepository, count total, err:", err.Error())
		return vehicles, 0, 0, err
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
		sortBy := `vehicle_id`
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

	// log.Println("vehicleRepository, querySelect:", querySelect)
	err = repository.Select(&vehicles, querySelect)
	if err != nil {
		log.Println("vehicleRepository, FindAllByCustId, err:", err.Error())
		return vehicles, total, lastPage, err
	}

	return vehicles, total, lastPage, nil
}

func (repository *vehicleRepositoryImpl) Store(vehicle model.Vehicle) (int64, error) {
	query :=
		`INSERT INTO mst.m_vehicle(
			cust_id, vehicle_no, vehicle_desc, vehicle_type, 
			length, width, height, volume,
			driver_id, helper_id, is_active, 
			created_by, created_at, updated_by, updated_at, 
			is_del, deleted_by, deleted_at,weight)
		VALUES ( 
			$1, $2, $3, $4, 
			$5, $6, $7, $8, 
			$9, $10, $11,
			$12, $13, $14, $15,
			$16, $17, $18, $19
		) RETURNING vehicle_id;`
	lastInsertId := vehicle.VehicleId
	err := repository.QueryRow(query,
		vehicle.CustId, vehicle.VehicleNo, vehicle.VehicleDesc, vehicle.VehicleType,
		vehicle.Length, vehicle.Width, vehicle.Height, vehicle.Volume,
		vehicle.DriverId, vehicle.HelperId, vehicle.IsActive,
		vehicle.CreatedBy, vehicle.CreatedAt, vehicle.UpdatedBy, vehicle.UpdatedAt,
		vehicle.IsDel, vehicle.DeletedBy, vehicle.DeletedAt, vehicle.Weight).Scan(&lastInsertId)
	if err != nil {
		log.Println("vehicleRepository, Store, err:", err.Error())
		return vehicle.VehicleId, err
	}
	return vehicle.VehicleId, nil
}

func (repository *vehicleRepositoryImpl) Update(vehicleId int64, request entity.UpdateVehicleRequest) error {
	var (
		r            model.VehicleUpdate
		sqlSetFields string
		nRows        int64
	)

	reqByte, _ := json.Marshal(request)
	_ = json.Unmarshal(reqByte, &r)
	sqlPatch := sql_helper.SQLPatches(r)

	// data, _ := json.Marshal(sqlPatch)
	// fmt.Printf("vehicleRepository, Update, Fields & Args: %s\n", data)

	for i, _ := range sqlPatch.Fields {
		sqlSetFields += sqlPatch.Fields[i] + ", "
		i++
	}
	sqlSetFields = strings.TrimSuffix(sqlSetFields, ", ")

	query := `UPDATE mst.m_vehicle
			  SET ` + sqlSetFields + `,
			  	updated_at = CURRENT_TIMESTAMP
			  WHERE is_del = false 
			  AND cust_id = :cust_id 
			  AND vehicle_id = :vehicle_id_old;`

	// log.Println("vehicleRepository, Update, query:", query)

	sqlPatch.Args["vehicle_id_old"] = vehicleId
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

func (repository *vehicleRepositoryImpl) Delete(custId string, vehicleId int64, deletedBy int64) error {
	var nRows int64
	query := `UPDATE mst.m_vehicle
			SET is_del = true,
				deleted_at = CURRENT_TIMESTAMP,
				deleted_by = :deleted_by 
			WHERE is_del = false
			ANd cust_id = :cust_id
			AND vehicle_id = :vehicle_id;`

	wMap := map[string]interface{}{
		"cust_id":    custId,
		"vehicle_id": vehicleId,
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

func (repository *vehicleRepositoryImpl) FindOneByDriverAndCustId(driverID int64, custId string) (model.Vehicle, error) {
	vehicle := model.Vehicle{}
	query := `SELECT 
				cust_id, vehicle_id, vehicle_no, 
				vehicle_desc, vehicle_type, 
				length, width, height, COALESCE(weight, 0) as weight, volume,
				driver_id, helper_id,
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_vehicle 
			  WHERE driver_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&vehicle, query, driverID, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByVehicleIdAndCustId, err:", err.Error())
		return vehicle, err
	}

	return vehicle, nil
}

func (repository *vehicleRepositoryImpl) FindOneByHelperAndCustId(helperID int64, custId string) (model.Vehicle, error) {
	vehicle := model.Vehicle{}
	query := `SELECT 
				cust_id, vehicle_id, vehicle_no, 
				vehicle_desc, vehicle_type, 
				length, width, height, COALESCE(weight, 0) as weight, volume,
				driver_id, helper_id,
				is_active, created_by,
				created_at, updated_by, updated_at,
				is_del, deleted_by, deleted_at
			  FROM mst.m_vehicle 
			  WHERE helper_id = $1 
			  AND cust_id = $2`
	err := repository.Get(&vehicle, query, helperID, custId)
	if err != nil {
		log.Println("vehicleRepository, FindOneByVehicleIdAndCustId, err:", err.Error())
		return vehicle, err
	}

	return vehicle, nil
}
