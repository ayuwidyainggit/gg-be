package pjp

import (
	"context"
	"scyllax-pjp/data/response"

	"gorm.io/gorm"
)

func (repo *pjpRepository) IsPrincipalCustomer(ctx context.Context, tx *gorm.DB, custID string) (bool, error) {
	var distributorID *string

	err := tx.WithContext(ctx).Table("smc.m_customer").
		Select("distributor_id").
		Where("cust_id = ?", custID).
		Scan(&distributorID).Error
	if err != nil {
		return false, err
	}

	return len(custID) == 6 && distributorID == nil, nil
}

func (repo *pjpRepository) GetDestinationDetails(
	ctx context.Context,
	tx *gorm.DB,
	pjpID int,
	date string,
	limit int,
	page int,
	sortOrder string,
	custID string,
	isPrincipal bool,
) ([]response.DestinationDetailRow, int64, error) {
	if !isPrincipal {
		return repo.getDistributorDestinationDetails(ctx, tx, pjpID, date, limit, page, sortOrder, custID)
	}

	return repo.getPrincipalDestinationDetails(ctx, tx, pjpID, date, limit, page, sortOrder, custID)
}

func (repo *pjpRepository) getPrincipalDestinationDetails(
	ctx context.Context,
	tx *gorm.DB,
	pjpID int,
	date string,
	limit int,
	page int,
	sortOrder string,
	custID string,
) ([]response.DestinationDetailRow, int64, error) {
	var (
		rows  []response.DestinationDetailRow
		total int64
	)
	custScope := tx.Table("smc.m_customer").
		Select("cust_id").
		Where("cust_id = ? OR parent_cust_id = ?", custID, custID)

	baseQuery := tx.WithContext(ctx).
		Table("pjp_principles.destinations_history dh").
		Where("DATE(dh.date) = ?", date).
		Where("dh.pjp_id = ?", pjpID).
		Where("dh.cust_id IN (?)", custScope)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Select(`
			COALESCE(dh.route_code, 0) AS route_code,
			COALESCE(r.route_name, '') AS route_name,
			COALESCE(dh.week, 0) AS week,
			COALESCE(dh.year, 0) AS year,
			DATE(dh.date) AS date,
			COALESCE(dh.destination_id, 0) AS destination_id,
			COALESCE(dh.destination_code, '') AS destination_code,
			COALESCE(dh.destination_type, '') AS destination_type,
			COALESCE(dh.destination_name, '') AS destination_name,
			COALESCE(dh.longitude::text, '') AS longitude,
			COALESCE(dh.latitude::text, '') AS latitude,
			COALESCE(dh.destination_status::text, '') AS destination_status,
			COALESCE(dh.destination_address, '') AS destination_address
		`).
		Joins("LEFT JOIN pjp_principles.routes r ON r.route_code = dh.route_code").
		Order("dh.destination_id " + sortOrder).
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (repo *pjpRepository) getDistributorDestinationDetails(
	ctx context.Context,
	tx *gorm.DB,
	pjpID int,
	date string,
	limit int,
	page int,
	sortOrder string,
	custID string,
) ([]response.DestinationDetailRow, int64, error) {
	var (
		rows  []response.DestinationDetailRow
		total int64
	)

	baseQuery := tx.WithContext(ctx).
		Table("pjp.route_outlet_history roh").
		Where("roh.date = ?", date).
		Where("roh.pjp_id = ?", pjpID).
		Where("roh.cust_id = ?", custID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Select(`
			COALESCE(roh.route_code, 0) AS route_code,
			COALESCE(roh.route_name, '') AS route_name,
			COALESCE(roh.week, 0) AS week,
			COALESCE(roh.year, 0) AS year,
			roh.date AS date,
			COALESCE(roh.outlet_id, 0) AS destination_id,
			COALESCE(roh.outlet_code, '') AS destination_code,
			CASE WHEN COALESCE(roh.is_extra_call, FALSE) THEN 'distributor' ELSE 'outlet' END AS destination_type,
			COALESCE(roh.outlet_name, '') AS destination_name,
			COALESCE(roh.longitude::text, '') AS longitude,
			COALESCE(roh.latitude::text, '') AS latitude,
			COALESCE(roh.outlet_status::text, '') AS destination_status,
			COALESCE(roh.outlet_address, '') AS destination_address
		`).
		Order("roh.outlet_id " + sortOrder).
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
