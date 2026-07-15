package arrival_report

import (
	"scyllax-pjp/model"

	"gorm.io/gorm"
)

// ArrivalReportRepository defines the interface for arrival report repository
type ArrivalReportRepository interface {
	Create(db *gorm.DB, arrivalReport *model.ArrivalReport) error
	FindById(db *gorm.DB, arrivalId int64) (*model.ArrivalReport, error)
	FindByOutletIdAndUserId(db *gorm.DB, outletId int64, userId int64, activity string) (*model.ArrivalReport, error)
	Update(db *gorm.DB, arrivalReport *model.ArrivalReport) error
}

type arrivalReportRepositoryImpl struct{}

// NewArrivalReportRepository creates a new instance of ArrivalReportRepository
func NewArrivalReportRepository() ArrivalReportRepository {
	return &arrivalReportRepositoryImpl{}
}

// Create inserts a new arrival report record
func (r *arrivalReportRepositoryImpl) Create(db *gorm.DB, arrivalReport *model.ArrivalReport) error {
	return db.Create(arrivalReport).Error
}

// FindById retrieves an arrival report by its ID
func (r *arrivalReportRepositoryImpl) FindById(db *gorm.DB, arrivalId int64) (*model.ArrivalReport, error) {
	var arrivalReport model.ArrivalReport
	err := db.Where("arrival_id = ? AND (is_del = false OR is_del IS NULL)", arrivalId).First(&arrivalReport).Error
	if err != nil {
		return nil, err
	}
	return &arrivalReport, nil
}

// FindByOutletIdAndUserId retrieves an arrival report by outlet ID, user ID, and activity
func (r *arrivalReportRepositoryImpl) FindByOutletIdAndUserId(db *gorm.DB, outletId int64, userId int64, activity string) (*model.ArrivalReport, error) {
	var arrivalReport model.ArrivalReport
	err := db.Where("outlet_id = ? AND user_id = ? AND activity = ? AND (is_del = false OR is_del IS NULL)", outletId, userId, activity).
		Order("created_at DESC").
		First(&arrivalReport).Error
	if err != nil {
		return nil, err
	}
	return &arrivalReport, nil
}

// Update updates an existing arrival report record
func (r *arrivalReportRepositoryImpl) Update(db *gorm.DB, arrivalReport *model.ArrivalReport) error {
	return db.Save(arrivalReport).Error
}
