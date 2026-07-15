package live_monitoring

import (
	"context"
	"errors"
	"scyllax-pjp/constant"
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

// GetVisitInformationPrincipal retrieves visit information for principal users
func (r *liveMonitoringRepository) GetVisitInformationPrincipal(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) (*model.VisitInformationRow, error) {
	var result model.VisitInformationRow

	query := tx.WithContext(ctx).Table("pjp_principles.permanent_journey_plans pjp").
		Select(`
			me.emp_id,
			me.emp_code,
			me.emp_name,
			COUNT(d.destination_code) AS plan,
			COUNT(CASE WHEN ovl."start" IS NOT NULL THEN 1 END) AS on_going,
			COUNT(CASE WHEN ovl.finish IS NOT NULL THEN 1 END) AS visited,
			COUNT(CASE WHEN ovl.skip_at IS NOT NULL THEN 1 END) AS total_skip,
			COUNT(ovl.outlet_code) AS matched
		`).
		Joins("JOIN pjp_principles.route_pop_permanent rpp ON rpp.pjp_id = pjp.id").
		Joins("JOIN pjp_principles.routes r ON r.route_code = rpp.route_code").
		Joins("JOIN pjp_principles.destinations d ON d.route_code = r.route_code").
		Joins("JOIN mst.m_salesman ms2 ON pjp.salesman_id = ms2.emp_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms2.emp_id").
		Joins("JOIN smc.m_customer mc ON ms2.cust_id = mc.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Joins("LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.pjp_id = pjp.id AND ovl.outlet_code = d.destination_code AND DATE(ovl.date) = ?", date).
		Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms WHERE ms.cust_id IN ?)", custIDs).
		Where("DATE(rpp.date) = ?", date).
		Where("me.emp_id = ?", empID).
		Where("pjp.approval_status = ?", constant.ApprovalStatusApproved).
		Group("me.emp_id, me.emp_code, me.emp_name")

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// GetVisitInformationPrincipalFromHistory retrieves visit information for principal users from destinations history
func (r *liveMonitoringRepository) GetVisitInformationPrincipalFromHistory(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) (*model.VisitInformationRow, error) {
	var result model.VisitInformationRow

	query := tx.WithContext(ctx).Table("pjp_principles.destinations_history dh").
		Select(`
			pjp.salesman_id AS emp_id,
			pjp.salesman_code AS emp_code,
			COALESCE(ms.sales_name, '') AS emp_name,
			SUM(CASE WHEN dh.is_extra_call = false THEN 1 ELSE 0 END) AS plan,
			SUM(CASE WHEN ovl.arrive_at IS NOT NULL AND ovl.leave_at IS NULL THEN 1 ELSE 0 END) AS on_going,
			SUM(CASE WHEN dh.is_extra_call = true THEN 1 ELSE 0 END) AS extra_call,
			SUM(CASE WHEN ovl.arrive_at IS NOT NULL AND ovl.leave_at IS NOT NULL THEN 1 ELSE 0 END) AS visited,
			SUM(CASE WHEN ovl.skip_at IS NOT NULL THEN 1 ELSE 0 END) AS total_skip,
			COUNT(ovl.outlet_code) AS matched
		`).
		Joins("JOIN pjp_principles.permanent_journey_plans pjp ON pjp.id = dh.pjp_id AND pjp.cust_id = dh.cust_id").
		Joins("JOIN mst.m_outlet mo ON mo.outlet_id = dh.destination_id").
		Joins("LEFT JOIN pjp_principles.routes r ON r.route_code = dh.route_code").
		Joins("LEFT JOIN mst.m_distributor md ON md.distributor_id = dh.destination_id").
		Joins("LEFT JOIN pjp_principles.outlet_visit_list ovl ON ovl.outlet_id = dh.destination_id AND ovl.pjp_id = dh.pjp_id AND DATE(ovl.date) = DATE(dh.date)").
		Joins("LEFT JOIN mobile.visits v ON v.emp_code = pjp.salesman_code AND v.outlet_code = mo.outlet_code AND DATE(v.created_at) = DATE(dh.date)").
		Joins("LEFT JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id").
		Where("pjp.salesman_id IN (SELECT emp_id FROM mst.m_salesman ms2 WHERE ms2.cust_id IN ?)", custIDs).
		Where("pjp.salesman_id = ?", empID).
		Where("DATE(dh.date) = ?", date).
		Where("pjp.approval_status IN ?", []string{constant.ApprovalStatusApproved, "Need Review"}).
		Group("pjp.salesman_id, pjp.salesman_code, ms.sales_name")

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// CountTotalVisitsPrincipal counts total visits for a salesman on a specific date (Principal)
func (r *liveMonitoringRepository) CountTotalVisitsPrincipal(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) (int64, error) {
	var count int64

	// Join with PJP to get salesman_id, because ovl only has pjp_id
	query := tx.WithContext(ctx).Table("pjp_principles.outlet_visit_list ovl").
		Joins("JOIN pjp_principles.permanent_journey_plans pjp ON ovl.pjp_id = pjp.id").
		Joins("JOIN mst.m_salesman ms ON pjp.salesman_id = ms.emp_id").
		Where("ms.cust_id IN ?", custIDs).
		Where("pjp.salesman_id = ?", empID).
		Where("DATE(ovl.date) = ?", date)

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetVisitInformationDistributor retrieves visit information for distributor users
func (r *liveMonitoringRepository) GetVisitInformationDistributor(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (*model.VisitInformationRow, error) {
	var result model.VisitInformationRow

	query := tx.WithContext(ctx).Table("pjp.permanent_journey_plans pjp").
		Select(`
			me.emp_id,
			me.emp_code,
			me.emp_name,
			0 AS plan,
			0 AS extra_call,
			0 AS on_going,
			0 AS visited,
			0 AS total_skip,
			0 AS matched
		`).
		Joins("JOIN pjp.route_pop_permanent rpp ON rpp.pjp_id = pjp.id").
		Joins("JOIN pjp.routes r ON r.route_code = rpp.route_code").
		Joins("JOIN pjp.route_outlet_history roh ON roh.pjp_id = pjp.id AND roh.route_code = r.route_code AND DATE(roh.date) = ?", date).
		Joins("JOIN mst.m_salesman ms ON pjp.salesman_id = ms.emp_id").
		Joins("JOIN mst.m_employee me ON me.emp_id = ms.emp_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("pjp.salesman_id = ?", empID).
		Where("DATE(rpp.date) = ?", date).
		Where("md.distributor_id = ?", distributorID).
		Where("pjp.approval_status IN ?", []string{constant.ApprovalStatusApproved, "Need Review"}).
		Group("me.emp_id, me.emp_code, me.emp_name")

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *liveMonitoringRepository) CountDistributorPlannedVisits(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.route_outlet_history roh").
		Select("COUNT(DISTINCT roh.outlet_id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = roh.pjp_id AND pjp.cust_id = roh.cust_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("DATE(roh.date) = ?", date).
		Where("pjp.salesman_id = ?", empID).
		Where("md.distributor_id = ?", distributorID).
		Where("roh.is_extra_call = false").
		Where("pjp.approval_status IN ?", []string{constant.ApprovalStatusApproved, "Need Review"})

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *liveMonitoringRepository) CountDistributorExtraCalls(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.outlet_visit_list ovl").
		Select("COUNT(DISTINCT ovl.id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = ovl.pjp_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("DATE(ovl.date) = ?", date).
		Where("pjp.salesman_id = ?", empID).
		Where("md.distributor_id = ?", distributorID).
		Where("ovl.is_extra_call = true")

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *liveMonitoringRepository) CountDistributorOnGoingVisits(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.outlet_visit_list ovl").
		Select("COUNT(DISTINCT ovl.id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = ovl.pjp_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("DATE(ovl.date) = ?", date).
		Where("pjp.salesman_id = ?", empID).
		Where("md.distributor_id = ?", distributorID).
		Where("ovl.arrive_at IS NOT NULL").
		Where("ovl.leave_at IS NULL").
		Where("ovl.skip_at IS NULL").
		Where("ovl.finish IS NULL")

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *liveMonitoringRepository) CountDistributorVisitedVisits(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.outlet_visit_list ovl").
		Select("COUNT(DISTINCT ovl.id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = ovl.pjp_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("DATE(ovl.date) = ?", date).
		Where("pjp.salesman_id = ?", empID).
		Where("md.distributor_id = ?", distributorID).
		Where("ovl.leave_at IS NOT NULL").
		Where("ovl.skip_at IS NULL")

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *liveMonitoringRepository) CountDistributorSkippedVisits(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID, distributorID int,
) (int64, error) {
	var count int64

	query := tx.WithContext(ctx).Table("pjp.outlet_visit_list ovl").
		Select("COUNT(DISTINCT ovl.id)").
		Joins("JOIN pjp.permanent_journey_plans pjp ON pjp.id = ovl.pjp_id").
		Joins("JOIN mst.m_salesman ms ON ms.emp_id = pjp.salesman_id AND ms.cust_id = pjp.cust_id").
		Joins("JOIN smc.m_customer mc ON mc.cust_id = ms.cust_id").
		Joins("JOIN mst.m_distributor md ON md.distributor_id = mc.distributor_id").
		Where("DATE(ovl.date) = ?", date).
		Where("pjp.salesman_id = ?", empID).
		Where("md.distributor_id = ?", distributorID).
		Where("ovl.skip_at IS NOT NULL").
		Where("ovl.leave_at IS NULL")

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetSales retrieves sales data for a specific employee on a specific date
func (r *liveMonitoringRepository) GetSales(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) ([]model.SalesRow, error) {
	var results []model.SalesRow

	query := tx.WithContext(ctx).Table("sls.order ord").
		Select(`
			o.outlet_id,
			o.outlet_name,
			o.outlet_code,
			SUM(ord.total) AS sales_order
		`).
		Joins("INNER JOIN mst.m_outlet o ON ord.outlet_id = o.outlet_id AND ord.cust_id = o.cust_id").
		Where("ord.cust_id IN ?", custIDs).
		Where("DATE(ord.ro_date) = ?", date).
		Where("ord.salesman_id = ?", empID).
		Where("ord.is_del = false").
		Group("o.outlet_id, o.outlet_name, o.outlet_code").
		Order("sales_order DESC")

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetReturns retrieves return data for a specific employee on a specific date
func (r *liveMonitoringRepository) GetReturns(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) ([]model.ReturnRow, error) {
	var results []model.ReturnRow

	query := tx.WithContext(ctx).Table("sls.return r").
		Select(`
			o.outlet_id,
			o.outlet_name,
			o.outlet_code,
			SUM(r.total) AS return_total
		`).
		Joins("INNER JOIN mst.m_outlet o ON r.outlet_id = o.outlet_id AND r.cust_id = o.cust_id").
		Where("r.cust_id IN ?", custIDs).
		Where("r.emp_id = ?", empID).
		Where("DATE(r.return_date) = ?", date).
		Where("r.is_del = false").
		Group("o.outlet_id, o.outlet_name, o.outlet_code").
		Order("return_total DESC")

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetCollections retrieves paid collection data per outlet for a specific employee on a specific date.
func (r *liveMonitoringRepository) GetCollections(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) ([]model.CollectionRow, error) {
	var results []model.CollectionRow

	query := `
		WITH payment_per_invoice AS (
			SELECT
				deposit_no,
				cust_id,
				invoice_no,
				SUM(COALESCE(payment_amount, 0)) AS payment_amount
			FROM acf.deposit_payment
			GROUP BY deposit_no, cust_id, invoice_no
		)
		SELECT
			mo.outlet_id,
			mo.outlet_code,
			mo.outlet_name,
			SUM(COALESCE(ppi.payment_amount, 0)) AS collection_total
		FROM acf.deposit d
		JOIN acf.deposit_detail dd
			ON dd.deposit_no = d.deposit_no
			AND dd.cust_id = d.cust_id
		JOIN payment_per_invoice ppi
			ON ppi.deposit_no = dd.deposit_no
			AND ppi.cust_id = dd.cust_id
			AND ppi.invoice_no = dd.invoice_no
		JOIN sls."order" o
			ON o.invoice_no = dd.invoice_no
			AND o.cust_id = dd.cust_id
		JOIN mst.m_outlet mo
			ON mo.outlet_id = o.outlet_id
			AND mo.cust_id = o.cust_id
		WHERE d.cust_id IN ?
			AND DATE(d.deposit_date) = ?
			AND d.collection_no IS NOT NULL
			AND d.emp_id = ?
		GROUP BY mo.outlet_id, mo.outlet_code, mo.outlet_name
		ORDER BY mo.outlet_code
	`

	err := tx.WithContext(ctx).Raw(query, custIDs, date, empID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetExpenses retrieves expense data for a specific salesman and date.
func (r *liveMonitoringRepository) GetExpenses(
	ctx context.Context,
	tx *gorm.DB,
	custID string,
	empID int,
	date string,
) ([]model.ExpenseRow, error) {
	var results []model.ExpenseRow

	query := tx.WithContext(ctx).Table("acf.expense e").
		Select(`
			e.expense_type_id,
			et.expense_type_name,
			e.note,
			e.amount
		`).
		Joins("JOIN acf.expense_type et ON e.expense_type_id = et.expense_type_id").
		Where("e.cust_id = ?", custID).
		Where("e.collector_id = ?", empID).
		Where("DATE(e.date) = ?", date).
		Where("e.is_del = false")

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetShipments retrieves shipment data for a specific employee on a specific date
func (r *liveMonitoringRepository) GetShipments(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) ([]model.ShipmentRow, error) {
	var results []model.ShipmentRow

	query := tx.WithContext(ctx).Table("tms.shipment_invoices si").
		Select(`
			si.shipment_no,
			si.status,
			si.outlet_id,
			o.outlet_name,
			o.outlet_code,
			SUM(si.total_netto) AS total_netto
		`).
		Joins("INNER JOIN mst.m_outlet o ON si.outlet_id = o.outlet_id AND si.cust_id = o.cust_id").
		Where("si.cust_id IN ?", custIDs).
		Where("si.salesman_id = ?", empID).
		Where("DATE(si.delivery_date) = ?", date).
		Group("si.shipment_no, si.status, si.outlet_id, o.outlet_name, o.outlet_code")

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// GetActivityTime retrieves check-in time for a specific employee on a specific date
// Note: mobile.attendances uses emp_code, so we need to get it from m_employee first
func (r *liveMonitoringRepository) GetActivityTime(
	ctx context.Context,
	tx *gorm.DB,
	date string,
	empID int,
) (*string, error) {
	var result model.AttendanceRow

	// Get emp_code from emp_id first via subquery
	query := tx.WithContext(ctx).Table("mobile.attendances a").
		Select("a.created_at").
		Joins("JOIN mst.m_employee me ON me.emp_code = a.emp_code").
		Where("me.emp_id = ?", empID).
		Where("DATE(a.created_at) = ?", date).
		Where("a.type = 1"). // 1 = checkin
		Order("a.created_at ASC").
		Limit(1)

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result.CreatedAt, nil
}

// GetDistributorInfo retrieves distributor information
func (r *liveMonitoringRepository) GetDistributorInfo(
	ctx context.Context,
	tx *gorm.DB,
	distributorID int,
) (*model.DistributorInfoRow, error) {
	var result model.DistributorInfoRow

	query := tx.WithContext(ctx).Table("mst.m_distributor").
		Select("distributor_id, distributor_code, distributor_name").
		Where("distributor_id = ?", distributorID)

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

// GetUserFullname retrieves user fullname from sys.m_user for principal
func (r *liveMonitoringRepository) GetUserFullname(
	ctx context.Context,
	tx *gorm.DB,
	custID string,
) (*string, error) {
	var result model.UserFullnameRow

	query := tx.WithContext(ctx).Table("sys.m_user").
		Select("user_fullname").
		Where("cust_id = ?", custID).
		Limit(1)

	err := query.Take(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &result.UserFullname, nil
}

// GetChildCustIDs retrieves all child cust_ids for a principal
func (r *liveMonitoringRepository) GetChildCustIDs(
	ctx context.Context,
	tx *gorm.DB,
	parentCustID string,
) ([]string, error) {
	var custIDs []string
	// Get self + children
	err := tx.WithContext(ctx).Table("smc.m_customer").
		Where("cust_id = ? OR parent_cust_id = ?", parentCustID, parentCustID).
		Pluck("cust_id", &custIDs).Error
	return custIDs, err
}

// GetSalesmanCustID retrieves cust_id for a specific salesman
func (r *liveMonitoringRepository) GetSalesmanCustID(
	ctx context.Context,
	tx *gorm.DB,
	empID int,
) (string, error) {
	var custID string
	err := tx.WithContext(ctx).Table("mst.m_salesman").
		Select("cust_id").
		Where("emp_id = ?", empID).
		Take(&custID).Error
	return custID, err
}

// GetSubmittedSurveyData retrieves submitted survey aggregation per outlet
func (r *liveMonitoringRepository) GetSubmittedSurveyData(
	ctx context.Context,
	tx *gorm.DB,
	custIDs []string,
	date string,
	empID int,
) ([]model.SurveyDataRow, error) {
	var results []model.SurveyDataRow

	err := tx.WithContext(ctx).Table("mst.survey_answer sa").
		Select(`
			COUNT(sa.survey_answer_id) AS submission,
			ms.survey_title,
			mo.outlet_code,
			mo.outlet_name
		`).
		Joins("JOIN mst.m_survey ms ON ms.survey_id = sa.survey_id").
		Joins("JOIN mst.m_outlet mo ON mo.outlet_id = sa.outlet_id AND mo.cust_id = sa.cust_id").
		Where("sa.cust_id IN ?", custIDs).
		Where("DATE(sa.answer_date) = ?", date).
		Where("sa.emp_id = ?", empID).
		Where("sa.status = ?", "Submitted").
		Group("ms.survey_title, mo.outlet_code, mo.outlet_name").
		Order("ms.survey_title ASC, mo.outlet_code ASC, mo.outlet_name ASC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
