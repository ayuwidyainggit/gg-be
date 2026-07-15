package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"strings"
	"time"
)

type VehicleRepoImpl struct {
	db *gorm.DB
}

type VehicleRepo interface {
	GetVehicle(ctx context.Context, dataFilter entity.VehicleQueryFilter) (data []entity.VehicleResponse, total int)
}

func NewVehicleRepoImpl(db *gorm.DB) VehicleRepo {
	return &VehicleRepoImpl{db: db}
}

func (repo *VehicleRepoImpl) GetVehicle(ctx context.Context, dataFilter entity.VehicleQueryFilter) (data []entity.VehicleResponse, total int) {
	selectQuery := `a.cust_id, a.vehicle_id, a.vehicle_no, a.vehicle_desc, a.vehicle_type, 
			a.length, a.width, a.height, a.volume, a.driver_id, a.helper_id, 
			d.emp_name as driver_name, h.emp_name as helper_name`
	countQuery := `COUNT(*)`

	rawQuery := `
		WITH availableDrivers AS (
			SELECT
				cust_id,
				vehicle_id,
				vehicle_no,
				vehicle_desc,
				vehicle_type,
				length,
				width,
				height,
				volume,
				driver_id,
				helper_id
			FROM
				mst.m_vehicle 
		),
		getEmployee AS (
			SELECT
				emp_id,
				emp_grp_id,
				emp_name
			FROM
				mst.m_employee
		)
		SELECT 
			%s
		FROM 
			availableDrivers a
		JOIN getEmployee d ON a.driver_id = d.emp_id
		JOIN getEmployee h ON a.helper_id = h.emp_id
		WHERE a.cust_id = 'C220010001'
	`

	baseQuery := fmt.Sprintf(rawQuery, selectQuery)

	if dataFilter.DeliveryDate != "" {
		dateFilterQuery := ` AND NOT EXISTS (SELECT 1 FROM tms.shipments s WHERE a.vehicle_id = s.vehicle_id AND s.delivery_date = CURRENT_DATE) `
		if dataFilter.DeliveryDate != time.Now().Format("2006-01-02") {
			dateFilterQuery = " AND tms.shipments.delivery_date = '" + dataFilter.DeliveryDate + "'"
		}
		baseQuery += dateFilterQuery
	}

	sortBy := "a.vehicle_id DESC"
	if dataFilter.Sort != "" {
		var sortClauses []string
		for _, row := range strings.Split(dataFilter.Sort, ",") {
			colSort := strings.Split(row, ":")
			if len(colSort) > 1 {
				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", colSort[0], colSort[1]))
			}
		}
		if len(sortClauses) > 0 {
			sortBy = strings.Join(sortClauses, ", ")
		}
	}
	baseQuery += " ORDER BY " + sortBy

	if dataFilter.Limit > 0 && dataFilter.Page > 0 {
		offset := (dataFilter.Page - 1) * dataFilter.Limit
		baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", dataFilter.Limit, offset)
	}

	result := repo.db.Raw(baseQuery).WithContext(ctx).Scan(&data)
	helper.ErrorPanic(result.Error)

	countQuery = fmt.Sprintf(rawQuery, countQuery)

	if dataFilter.DeliveryDate != "" {
		if dataFilter.DeliveryDate == time.Now().Format("2006-01-02") {
			countQuery += `AND NOT EXISTS (SELECT 1 FROM tms.shipments s WHERE a.vehicle_id = s.vehicle_id AND s.delivery_date = CURRENT_DATE) `
		} else {
			countQuery += " AND tms.shipments.delivery_date = '" + dataFilter.DeliveryDate + "'"
		}
	}

	totalRaw := repo.db.Raw(countQuery).WithContext(ctx).Scan(&total)
	if totalRaw.Error != nil {
		return nil, 0
	}

	return data, total
}
