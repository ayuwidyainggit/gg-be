package repository

import (
	"context"
	"mobile/model"

	"gorm.io/gorm"
)

type (
	RepositoryMSalesmanImpl struct {
		*gorm.DB
	}
)
type MSalesmanRepository interface {
	FindOneByEmpId(custId string, empId int64) (salesman model.MSalesmanRead, err error)
	FindPjpSalesmanDetail(empId int64) (detail model.PjpSalesmanDetail, err error)
	FindPjpSalesmanWarehouses(empId int64, custId string) (warehouses []model.PjpSalesmanWarehouse, err error)
	// Update(c context.Context, ConfigId string, data model.MSalesman) error
}

func NewMSalesmanRepository(db *gorm.DB) *RepositoryMSalesmanImpl {
	return &RepositoryMSalesmanImpl{db}
}

func (repo *RepositoryMSalesmanImpl) model(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return repo.WithContext(ctx)
}

func (repository *RepositoryMSalesmanImpl) FindOneByEmpId(custId string, empId int64) (salesman model.MSalesmanRead, err error) {
	err = repository.
		Select("mst.m_salesman.*,st.sales_team_name, msc.opr_type as opr_type_canvas, mst.m_salesman.opr_type as opr_type_order_taking, CASE WHEN mst.m_salesman.is_taking_order IS NULL THEN false ELSE mst.m_salesman.is_taking_order END AS is_active_gudang_utama, msc.is_active as is_active_gudang_canvas").
		Joins("LEFT JOIN mst.m_sales_team st ON st.sales_team_id = mst.m_salesman.sales_team_id ").
		Joins("LEFT JOIN mst.m_salesman_canvas msc ON msc.emp_id = mst.m_salesman.emp_id ").
		Where("mst.m_salesman.emp_id = ? AND mst.m_salesman.is_active = true AND mst.m_salesman.cust_id = ?", empId, custId).
		Take(&salesman).Error
	return salesman, err
}

// FindPjpSalesmanDetail finds PJP salesman detail with customer, distributor, and sales team data
func (repository *RepositoryMSalesmanImpl) FindPjpSalesmanDetail(empId int64) (detail model.PjpSalesmanDetail, err error) {
	err = repository.
		Select(`
			ms.cust_id,
			mc.cust_name,
			mc.distributor_id,
			md.distributor_code,
			md.distributor_name,
			ms.emp_id,
			ms.sales_name,
			ms.sales_team_id,
			st.sales_team_code,
			st.sales_team_name,
	  	ms.is_taking_order,
		  ms.opr_type,
			ms.wh_id,
			msc.is_active as is_active_salesman_canvas,
			msc.wh_id as wh_id_canvas,
			msc.opr_type as opr_type_canvas
		`).
		Table("mst.m_salesman AS ms").
		Joins("LEFT JOIN smc.m_customer AS mc ON mc.cust_id = ms.cust_id AND mc.is_del = false").
		Joins("LEFT JOIN mst.m_distributor AS md ON md.distributor_id = mc.distributor_id").
		Joins("LEFT JOIN mst.m_sales_team AS st ON st.sales_team_id = ms.sales_team_id").
		Joins("LEFT JOIN mst.m_salesman_canvas AS msc ON msc.emp_id = ms.emp_id ").
		Where("ms.emp_id = ? AND ms.is_active = true AND ms.is_del = false", empId).
		Take(&detail).Error
	return detail, err
}

// FindPjpSalesmanWarehouses finds warehouse data for PJP salesman based on opr_type
func (repository *RepositoryMSalesmanImpl) FindPjpSalesmanWarehouses(empId int64, custId string) (warehouses []model.PjpSalesmanWarehouse, err error) {
	// Query warehouse from m_salesman (order taking - opr_type from m_salesman)
	var orderTakingWarehouses []model.PjpSalesmanWarehouse
	errOrderTaking := repository.
		Select(`
			ms.opr_type,
			ms.wh_id,
			mw.wh_code,
			mw.wh_name
		`).
		Table("mst.m_salesman AS ms").
		Joins("LEFT JOIN mst.m_warehouse AS mw ON mw.wh_id = ms.wh_id AND mw.cust_id = ?", custId).
		Where("ms.emp_id = ? AND ms.is_active = true AND ms.is_del = false AND ms.wh_id IS NOT NULL AND ms.opr_type IS NOT NULL", empId).
		Find(&orderTakingWarehouses).Error

	if errOrderTaking != nil {
		return warehouses, errOrderTaking
	}

	if len(orderTakingWarehouses) > 0 {
		warehouses = append(warehouses, orderTakingWarehouses...)
	}

	// Query warehouse from m_salesman_canvas (canvas - opr_type from m_salesman_canvas)
	var canvasWarehouses []model.PjpSalesmanWarehouse
	errCanvas := repository.
		Select(`
			msc.opr_type,
			msc.wh_id,
			mw.wh_code,
			mw.wh_name
		`).
		Table("mst.m_salesman_canvas AS msc").
		Joins("LEFT JOIN mst.m_warehouse AS mw ON mw.wh_id = msc.wh_id AND mw.cust_id = ?", custId).
		Where("msc.emp_id = ? AND msc.is_active = true AND msc.wh_id IS NOT NULL AND msc.opr_type IS NOT NULL", empId).
		Find(&canvasWarehouses).Error

	if errCanvas != nil {
		return warehouses, errCanvas
	}

	if len(canvasWarehouses) > 0 {
		warehouses = append(warehouses, canvasWarehouses...)
	}

	if len(warehouses) == 0 {
		return []model.PjpSalesmanWarehouse{}, nil
	}

	return warehouses, nil
}

// func (repository *RepositoryMSalesmanImpl) Update(c context.Context, ConfigId string, data model.MSalesman) error {
// 	result := repository.model(c).Model(&data).Where("config_id=?", ConfigId).Updates(&data)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	if result.RowsAffected == 0 {
// 		return errors.New("no rows affected")
// 	}
// 	return nil
// }
