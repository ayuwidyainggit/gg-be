package repository

import (
	"math"
	"mobile/entity"
	"mobile/model"
	"time"

	"gorm.io/gorm"
)

type LeaveRequestRepository interface {
	ExistsOverlappingLeave(custID string, empID int64, start, end time.Time) (bool, error)
	Create(data *model.LeaveRequest) error
	Count(filter entity.LeaveRequestQuery) (int64, error)
	List(filter entity.LeaveRequestQuery) ([]model.LeaveRequestRead, int64, int, error)
	FindByCustIDAndEmpID(custID string, empID int64, startDate, endDate string) (*model.LeaveRequest, error)
}

type leaveRequestRepositoryImpl struct {
	*gorm.DB
}

func NewLeaveRequestRepository(db *gorm.DB) LeaveRequestRepository {
	return &leaveRequestRepositoryImpl{db}
}

func (repo *leaveRequestRepositoryImpl) ExistsOverlappingLeave(custID string, empID int64, start, end time.Time) (bool, error) {
	var leaveID int64
	err := repo.Table("mobile.leave_request lr").
		Select("lr.leave_id").
		Where("lr.cust_id = ? AND lr.emp_id = ?", custID, empID).
		Where("lr.start_date <= ? AND lr.end_date >= ?", end, start).
		Where("lr.approval != ?", "Canceled").
		Limit(1).
		Scan(&leaveID).Error
	if err != nil {
		return false, err
	}
	return leaveID > 0, nil
}

func (repo *leaveRequestRepositoryImpl) Create(data *model.LeaveRequest) error {
	return repo.DB.Create(data).Error
}

func (repo *leaveRequestRepositoryImpl) baseQuery(filter entity.LeaveRequestQuery) *gorm.DB {
	query := repo.Table("mobile.leave_request lr").
		Where("lr.cust_id = ? AND lr.emp_id = ?", filter.CustID, filter.EmpID)

	if filter.FilterEnd != "" {
		query = query.Where("lr.start_date <= ?", filter.FilterEnd)
	}
	if filter.FilterStart != "" {
		query = query.Where("lr.end_date >= ?", filter.FilterStart)
	}

	return query
}

func (repo *leaveRequestRepositoryImpl) Count(filter entity.LeaveRequestQuery) (int64, error) {
	var total int64
	err := repo.baseQuery(filter).Count(&total).Error
	return total, err
}

func (repo *leaveRequestRepositoryImpl) List(filter entity.LeaveRequestQuery) ([]model.LeaveRequestRead, int64, int, error) {
	var items []model.LeaveRequestRead
	var total int64

	countQuery := repo.baseQuery(filter)
	if err := countQuery.Count(&total).Error; err != nil {
		return items, 0, 0, err
	}

	limit, _, offset, lastPage := applyLeavePagination(filter.Limit, filter.Page, total, 10)

	err := repo.Table("mobile.leave_request lr").
		Select(`
			lr.cust_id,
			lr.emp_id,
			lr.start_date,
			lr.end_date,
			lr.reason,
			lr.file_url,
			lr.file_name,
			lr.approval,
			AGE(lr.end_date::date + INTERVAL '1 day', lr.start_date::date) AS duration,
			u_created.user_fullname AS created_by,
			lr.created_at,
			u_approved.user_fullname AS approved_by,
			lr.approved_at,
			u_canceled.user_fullname AS canceled_by,
			lr.canceled_at
		`).
		Joins("LEFT JOIN sys.m_user u_created ON u_created.user_id = lr.created_by").
		Joins("LEFT JOIN sys.m_user u_approved ON u_approved.user_id = lr.approved_by").
		Joins("LEFT JOIN sys.m_user u_canceled ON u_canceled.user_id = lr.canceled_by").
		Where("lr.cust_id = ? AND lr.emp_id = ?", filter.CustID, filter.EmpID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if filter.FilterEnd != "" {
				db = db.Where("lr.start_date <= ?", filter.FilterEnd)
			}
			if filter.FilterStart != "" {
				db = db.Where("lr.end_date >= ?", filter.FilterStart)
			}
			return db
		}).
		Order("lr.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&items).Error
	if err != nil {
		return items, 0, 0, err
	}

	return items, total, lastPage, nil
}

func applyLeavePagination(limit, page int, total int64, defaultLimit int) (int, int, int, int) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage < 1 {
		lastPage = 1
	}
	return limit, page, offset, lastPage
}

func (repo *leaveRequestRepositoryImpl) FindByCustIDAndEmpID(custID string, empID int64, startDate, endDate string) (*model.LeaveRequest, error) {
	var data model.LeaveRequest
	err := repo.Model(&model.LeaveRequest{}).Where("cust_id = ? AND emp_id = ?", custID, empID).
		Where("approval != ?", "canceled").
		Where("start_date <= ? AND end_date >= ?", startDate, endDate).
		Take(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
