package repository

import (
	"context"
	"errors"
	"fmt"
	"scyllax-pjp/data/response"
	"scyllax-pjp/helper"
	"scyllax-pjp/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RouteOutletRepositoryImpl struct {
	Db *gorm.DB
}

func NewRouteOutletRepositoryImpl(Db *gorm.DB) RouteOutletRepository {
	return &RouteOutletRepositoryImpl{Db: Db}
}

func (repo *RouteOutletRepositoryImpl) Insert(ctx context.Context, route model.RouteOutlet) {
	result := repo.Db.WithContext(ctx).Save(&route)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) Create(ctx context.Context, route model.RouteOutlet) {
	result := repo.Db.WithContext(ctx).Create(&route)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) CreateBulk(ctx context.Context, outlets []model.RouteOutlet) error {
	if len(outlets) == 0 {
		return nil
	}
	return repo.Db.WithContext(ctx).Create(&outlets).Error
}

func (repo *RouteOutletRepositoryImpl) CreateAdditionalRoute(ctx context.Context, route model.RouteOutletAdditional) {
	var existingRoute model.RouteOutletAdditional
	result := repo.Db.WithContext(ctx).
		Where("route_code = ? AND outlet_id = ? AND pjp_id = ?", route.RouteCode, route.OutletID, route.PjpID).
		First(&existingRoute)

	if result.RowsAffected > 0 {
		fmt.Println("Route already exists, skipping creation.")
		return
	}

	createResult := repo.Db.WithContext(ctx).Create(&route)
	helper.ErrorPanic(createResult.Error)
}

func (repo *RouteOutletRepositoryImpl) FindByRouteCode(ctx context.Context, code int) (model.RouteOutlet, error) {
	var route model.RouteOutlet

	result := repo.Db.WithContext(ctx).First(&route, "route_code = ?", code)

	if result.Error != nil {
		return route, errors.New("route code not found")
	}
	return route, nil
}

func (repo *RouteOutletRepositoryImpl) FindById(ctx context.Context, id int) (model.RouteOutlet, error) {
	var route model.RouteOutlet

	result := repo.Db.WithContext(ctx).First(&route, "id = ?", id)

	if result.Error != nil {
		return route, errors.New("route code not found")
	}
	return route, nil
}

func (repo *RouteOutletRepositoryImpl) FindAll(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.RouteOutlet {
	var route []model.RouteOutlet

	query := repo.Db.Model(&route).Preload("Pjp").Preload("PjpOld")

	for field, value := range filters {
		switch field {
		case "salesman_name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Joins("LEFT JOIN pjp.permanent_journey_plans ON pjp.route_outlet.pjp_id = pjp.permanent_journey_plans.id").
					Where("pjp.permanent_journey_plans.salesman_name = ?", v)
			}
		case "salesman_code":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}
				query = query.Joins("LEFT JOIN pjp.permanent_journey_plans ON pjp.route_outlet.pjp_id = pjp.permanent_journey_plans.id").
					Where("pjp.permanent_journey_plans.salesman_code IN ?", codes)
			}
		case "status":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}

				query = query.Where("status IN ?", codes)
			}
		case "start_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				query = query.Where("created_at >= ?", v)
			}
		case "end_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				v = v.AddDate(0, 0, 1)
				query = query.Where("created_at <= ?", v)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					codes := strings.Split(v, ",")
					for i := range codes {
						codes[i] = strings.TrimSpace(codes[i])
					}

					query = query.Where(fmt.Sprintf("%s IN ?", field), codes)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("%s = ?", field), v)
				}
			}
		}
	}

	result := query.WithContext(ctx).Where("pjp.route_outlet.cust_id = ?", currentCustomerId).Scopes(response.Scopes(page, limit)).Find(&route)
	helper.ErrorPanic(result.Error)
	return route
}
func (repo *RouteOutletRepositoryImpl) FindAllEnhance(ctx context.Context, page, limit int, filters map[string]interface{}, currentCustomerId string) []model.Pjp {
	var pjp []model.Pjp

	query := repo.Db.Model(&pjp).Preload("RouteOutlets").Where("approval_status != ?", "Draft")

	for field, value := range filters {
		switch field {
		case "salesman_name":
			if v, ok := value.(string); ok && v != "" {
				query = query.Where("salesman_name = ?", v)
			}
		case "salesman_code":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}
				query = query.Where("salesman_code IN ?", codes)
			}
		case "status":
			if v, ok := value.(string); ok && v != "" {
				codes := strings.Split(v, ",")
				for i := range codes {
					codes[i] = strings.TrimSpace(codes[i])
				}

				query = query.Where("approval_status IN ?", codes)
			}
		case "start_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				query = query.Where("created_at >= ?", v)
			}
		case "end_date":
			if v, ok := value.(time.Time); ok && !v.IsZero() {
				v = v.AddDate(0, 0, 1)
				query = query.Where("created_at <= ?", v)
			}
		default:
			switch v := value.(type) {
			case string:
				if v != "" {
					codes := strings.Split(v, ",")
					for i := range codes {
						codes[i] = strings.TrimSpace(codes[i])
					}

					query = query.Where(fmt.Sprintf("%s IN ?", field), codes)
				}
			case int:
				if v != 0 {
					query = query.Where(fmt.Sprintf("%s = ?", field), v)
				}
			}
		}
	}

	result := query.WithContext(ctx).Where("cust_id = ?", currentCustomerId).Scopes(response.Scopes(page, limit)).Find(&pjp)
	helper.ErrorPanic(result.Error)
	return pjp
}

func (repo *RouteOutletRepositoryImpl) FindByRouteCodeAndPjpCode(ctx context.Context, routeCode int, pjpCode int) ([]model.RouteOutlet, error) {
	var route []model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where("route_code = ? AND pjp_code = ?", routeCode, pjpCode).Find(&route)

	if result.Error != nil {
		return route, result.Error
	}

	if result.RowsAffected == 0 {
		return route, errors.New("route outlet not found")
	}

	return route, nil
}

func (repo *RouteOutletRepositoryImpl) Save(ctx context.Context, route model.RouteOutlet) error {
	var err error

	dataset := model.RouteOutlet{
		RouteCode:  route.RouteCode,
		RouteName:  route.RouteName,
		PjpID:      route.PjpID,
		PjpCode:    route.PjpCode,
		OldPjpID:   route.PjpID,
		OldPjpCode: route.PjpCode,
		CustID:     route.CustID,
	}
	tx := repo.Db.Begin()
	err = tx.Model(&route).
		WithContext(ctx).
		Where("route_code = ?", dataset.RouteCode).
		Updates(dataset).Error
	if err != nil {
		tx.Rollback()
		helper.ErrorPanic(err)
	}

	return tx.Commit().Error
}

func (repo *RouteOutletRepositoryImpl) UpdatePivot(ctx context.Context, route model.RouteOutlet) {
	dataset := model.RouteOutlet{
		ID: route.ID,
		// RouteCode:    route.RouteCode,
		// OutletCode:   route.OutletCode,
		Status:       route.Status,
		VerifiedDate: route.VerifiedDate,
	}
	result := repo.Db.Model(&route).WithContext(ctx).
		Where("id = ?", dataset.ID).
		// Where("route_code = ?", dataset.RouteCode).
		// Where("outlet_code = ?", dataset.OutletCode).
		Updates(dataset)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) UpdatePjp(ctx context.Context, route model.RouteOutlet) error {
	tx := repo.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Model(&route).WithContext(ctx).
		Where("route_code = ?", route.RouteCode).
		Where("pjp_code = ?", route.PjpCode).
		Where("pjp_id = ?", route.PjpID).
		Delete(&route)

	data := tx.Table("pjp.routes").
		WithContext(ctx).
		Where("route_code = ?", route.RouteCode).
		Updates(map[string]interface{}{
			"is_assign": false,
		})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if data.Error != nil {
		tx.Rollback()
		return data.Error
	}

	tx.Commit()
	return nil

}

func (repo *RouteOutletRepositoryImpl) UpdatePjpRouteOutlet(ctx context.Context, route model.RouteOutlet) error {
	tx := repo.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Model(&route).WithContext(ctx).
		Where("pjp_code = ?", route.PjpCode).
		Where("pjp_id = ?", route.PjpID).
		Updates(map[string]interface{}{
			"pjp_code": nil,
			"pjp_id":   nil,
		})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil

}

func (repo *RouteOutletRepositoryImpl) DeleteByOutletCode(ctx context.Context, route model.RouteOutlet) error {
	results := repo.Db.WithContext(ctx).Model(&route).Where("route_code = ? AND outlet_code = ?", route.RouteCode, route.OutletCode).Delete(&route)
	if results.Error != nil {
		return results.Error
	}
	if results.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (repo *RouteOutletRepositoryImpl) DeleteByOutletCodeAdditional(ctx context.Context, route model.RouteOutletAdditional) error {
	results := repo.Db.WithContext(ctx).Model(&route).Where("route_code = ? AND outlet_code = ?", route.RouteCode, route.OutletCode).Delete(&route)
	if results.Error != nil {
		return results.Error
	}
	if results.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (repo *RouteOutletRepositoryImpl) FindByPjpCode(ctx context.Context, pjpCode int) (model.RouteOutlet, error) {
	var route model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where("pjp_code = ?", pjpCode).Find(&route)

	if result.Error != nil {
		return route, result.Error
	}

	if result.RowsAffected == 0 {
		return route, errors.New("pjp code not found")
	}

	return route, nil
}

func (repo *RouteOutletRepositoryImpl) FindByPjpCodeEnhance(ctx context.Context, pjpCode int) ([]model.RouteOutlet, error) {
	var routes []model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where("pjp_code = ?", pjpCode).Find(&routes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("pjp code not found")
	}

	return routes, nil
}

func (repo *RouteOutletRepositoryImpl) FindByPjpId(ctx context.Context, pjpId int) ([]model.RouteOutlet, error) {
	var routes []model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where("pjp_id = ?", pjpId).Find(&routes)

	if result.Error != nil {
		return nil, result.Error
	}

	return routes, nil
}

func (repo *RouteOutletRepositoryImpl) FindByRouteCodes(ctx context.Context, routeCode, pjpCode int) []model.RouteOutlet {
	var route []model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where("route_code = ?", routeCode).Where("pjp_code = ?", pjpCode).Find(&route)
	helper.ErrorPanic(result.Error)

	return route
}

func (repo *RouteOutletRepositoryImpl) Update(ctx context.Context, code int, name string) {
	dataset := map[string]interface{}{
		"route_name": name,
	}

	result := repo.Db.Model(model.RouteOutlet{}).
		Where("route_code = ?", code).WithContext(ctx).Updates(dataset)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) UpdateOrCreate(ctx context.Context, route model.RouteOutlet) {
	var existingRoute model.RouteOutlet
	result := repo.Db.WithContext(ctx).Where(&model.RouteOutlet{OutletCode: route.OutletCode}).Assign(&route).FirstOrCreate(&existingRoute)
	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) Count(ctx context.Context, currentCustomerId string) (int64, error) {
	var data []model.RouteOutlet
	var totalRows int64
	result := repo.Db.WithContext(ctx).Where("cust_id = ?", currentCustomerId).Find(&data).Count(&totalRows)
	helper.ErrorPanic(result.Error)
	return totalRows, nil
}

func (repo *RouteOutletRepositoryImpl) CountAllEnhance(ctx context.Context, currentCustomerId string) (int64, error) {
	var count int64
	query := repo.Db.Model(&model.Pjp{}).Where("approval_status != ?", "Draft").Where("status = ?", "true")

	// Tambahkan filter customerId jika dibutuhkan
	if currentCustomerId != "" {
		query = query.Where("cust_id = ?", currentCustomerId)
	}

	// Hitung total
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *RouteOutletRepositoryImpl) UpdateNewRoute(ctx context.Context, route model.RouteOutlet) {

	// now := time.Now()

	dataset := model.RouteOutlet{
		OutletID:     route.OutletID,
		OutletCode:   route.OutletCode,
		RouteCode:    route.RouteCode,
		RouteName:    route.RouteName,
		PjpID:        route.PjpID,
		PjpCode:      route.PjpCode,
		OldPjpID:     route.OldPjpID,
		OldPjpCode:   route.OldPjpCode,
		OldRouteCode: route.OldRouteCode,
		OldRouteName: route.OldRouteName,
		// VerifiedDate: &now,
	}

	//fmt.Println("old_pjp_code", dataset.OldPjpCode)

	result := repo.Db.Model(&route).WithContext(ctx).
		Where("old_route_code = ?", dataset.OldRouteCode).
		Where("outlet_code = ?", dataset.OutletCode).
		Where("old_pjp_code = ?", dataset.OldPjpCode).
		Updates(dataset)

	helper.ErrorPanic(result.Error)
}

func (repo *RouteOutletRepositoryImpl) GetAllOutletBySalesCode(ctx context.Context, salesCode, custId, date, routeCode string) (data []model.RouteOutlet) {
	result := repo.Db.Raw(`
		WITH selected_pjp AS (
			SELECT id
			FROM pjp.permanent_journey_plans
			WHERE salesman_code = ? AND cust_id = ?
		)
		SELECT ro.outlet_code, ro.outlet_id, ro.outlet_name, ro.outlet_address, ro.outlet_status, ro.status, ro.route_code, ro.pjp_id, ro.pjp_code, ro.longitude, ro.latitude
		FROM pjp.route_outlet ro
		LEFT JOIN (
			SELECT DISTINCT ON (pjp_id) *
			FROM pjp.route_pop_daily
			WHERE date = ?
		) rpd ON ro.pjp_id = rpd.pjp_id
		WHERE ro.pjp_id IN (SELECT id FROM selected_pjp)
		AND (ro.status = 'Approved' OR ro.status = 'Approved With Propose')
		AND ro.route_code = ?
	`, salesCode, custId, date, routeCode).Scan(&data)

	helper.ErrorPanic(result.Error)
	return data
}

func (repo *RouteOutletRepositoryImpl) FindByRouteCodeAndOutletIDAndPjpNull(ctx context.Context, routeCode int, outletID int) (*model.RouteOutlet, error) {
	var routeOutlet model.RouteOutlet
	err := repo.Db.Where(
		"route_code = ? AND outlet_id = ? AND pjp_id IS NULL AND pjp_code IS NULL AND old_pjp_id IS NULL AND old_pjp_code IS NULL", routeCode, outletID).
		First(&routeOutlet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &routeOutlet, nil
}

func (repo *RouteOutletRepositoryImpl) FindAllOutletIdByPjpId(ctx context.Context, pjpIds []int) []int {
	var outletIDs []int
	queryResult := repo.Db.WithContext(ctx).
		Table("(SELECT DISTINCT ro.outlet_id FROM pjp.route_outlet ro "+
			"JOIN pjp.route_pop_daily rpd ON ro.route_code = rpd.route_code "+
			"WHERE ro.status LIKE ? AND rpd.pjp_id IN ? "+
			"UNION "+
			"SELECT DISTINCT roa.outlet_id FROM pjp.route_outlet_additional roa "+
			"JOIN pjp.route_pop_daily rpd ON roa.route_code = rpd.parent_route "+
			"WHERE roa.status LIKE ? AND rpd.pjp_id IN ?) as combined_outlets",
			"%Approved%", pjpIds, "%Approved%", pjpIds).
		Pluck("outlet_id", &outletIDs)

	helper.ErrorPanic(queryResult.Error)
	return outletIDs
}

func (repo *RouteOutletRepositoryImpl) FindAllOutletIdByPjpIdToday(ctx context.Context, pjpIds []int) []int {
	var outletIDs []int
	timeNow := time.Now().Format("2006-01-02")

	queryResult := repo.Db.WithContext(ctx).
		Table("(SELECT DISTINCT ro.outlet_id FROM pjp.route_outlet ro "+
			"JOIN pjp.route_pop_daily rpd ON ro.route_code = rpd.route_code "+
			"WHERE ro.status LIKE ? AND rpd.pjp_id IN ? AND rpd.date = ? "+
			"UNION "+
			"SELECT DISTINCT roa.outlet_id FROM pjp.route_outlet_additional roa "+
			"JOIN pjp.route_pop_daily rpd ON roa.route_code = rpd.parent_route "+
			"WHERE roa.status LIKE ? AND rpd.pjp_id IN ? AND rpd.date = ?) as combined_outlets",
			"%Approved%", pjpIds, timeNow, "%Approved%", pjpIds, timeNow).
		Pluck("outlet_id", &outletIDs)

	helper.ErrorPanic(queryResult.Error)
	return outletIDs
}

func (repo *RouteOutletRepositoryImpl) SearchOutletIdByPjpId(ctx context.Context, pjpIds []int, search string) []int {
	var outletIDs []int
	queryResult := repo.Db.WithContext(ctx).
		Table("(SELECT DISTINCT ro.outlet_id FROM pjp.route_outlet ro "+
			"JOIN pjp.route_pop_daily rpd ON ro.route_code = rpd.route_code "+
			"WHERE (ro.status LIKE ? AND rpd.pjp_id IN ? AND LOWER(ro.outlet_name) LIKE LOWER(?)) "+
			"UNION "+
			"SELECT DISTINCT roa.outlet_id FROM pjp.route_outlet_additional roa "+
			"JOIN pjp.route_pop_daily rpd ON roa.route_code = rpd.parent_route "+
			"WHERE (roa.status LIKE ? AND rpd.pjp_id IN ? AND LOWER(roa.outlet_name) LIKE LOWER(?))) as combined_outlets",
			"%Approved%", pjpIds, "%"+search+"%", "%Approved%", pjpIds, "%"+search+"%").
		Pluck("outlet_id", &outletIDs)

	helper.ErrorPanic(queryResult.Error)
	return outletIDs
}

func (repo *RouteOutletRepositoryImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.Db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (repo *RouteOutletRepositoryImpl) MobileCancelAddOutletToRoute(ctx context.Context, route model.RouteOutletAdditional, tx ...*gorm.DB) error {
	var db *gorm.DB
	if len(tx) > 0 && tx[0] != nil { // Jika tx diberikan dan tidak nil, gunakan tx
		db = tx[0]
	} else { // Jika tidak, gunakan koneksi default repo.Db
		db = repo.Db
	}

	today := time.Now().Format("2006-01-02")

	results := db.WithContext(ctx).
		Model(&route).
		Where("route_code = ? AND outlet_code = ? AND pjp_id = ? AND DATE(date) = ? AND is_planned = ?",
			route.RouteCode, route.OutletCode, route.PjpID, today, false).
		Delete(&route)

	if results.Error != nil {
		return results.Error
	}
	if results.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
