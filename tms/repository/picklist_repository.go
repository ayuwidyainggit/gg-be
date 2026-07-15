package repository

import (
	"context"
	"fmt"
	"log"
	"scyllax-tms/entity"
	"scyllax-tms/helper"
	"scyllax-tms/model"

	"gorm.io/gorm"
)

type PicklistRepository interface {
	FindAll(ctx context.Context, dataFilter entity.GeneralQueryFilter, picklistFilter entity.PicklistFilter, customerId string) ([]model.Picklist, error)
	FindByID(ctx context.Context, id string) (model.Picklist, error)
	InsertPicklist(ctx context.Context, tx *gorm.DB, data model.Picklist) error
	UpdatePicklist(ctx context.Context, data model.Picklist) error
	DeletePicklist(ctx context.Context, id string) error
	InsertOrderPicklist(ctx context.Context, tx *gorm.DB, data model.OrderPicklist) (resp int, err error)
	InsertOrderProduct(ctx context.Context, tx *gorm.DB, data model.OrderProduct) error
	GetOrdersByPicklistNo(ctx context.Context, picklistNo string) ([]model.OrderPicklist, error)
	GetProductsByOrderID(ctx context.Context, orderID uint) ([]model.OrderProduct, error)
	GetNextSequence(ctx context.Context, yy, mm, dd string) (int, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	FindAllByInvoiceNo(ctx context.Context, customerId string) []string
	GetTotalCount(ctx context.Context, picklistFilter entity.PicklistFilter) (int, error)
	CountAll(ctx context.Context, filter entity.PicklistFilter, customerId string) (int64, error)
}

type PicklistRepositoryImpl struct {
	db *gorm.DB
}

func NewPicklistRepositoryImpl(db *gorm.DB) PicklistRepository {
	return &PicklistRepositoryImpl{db: db}
}

func (repo *PicklistRepositoryImpl) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (repo *PicklistRepositoryImpl) FindAll(ctx context.Context, dataFilter entity.GeneralQueryFilter, picklistFilter entity.PicklistFilter, customerId string) ([]model.Picklist, error) {
	var data []model.Picklist
	db := repo.db.Debug().WithContext(ctx).
		Model(&model.Picklist{}).
		Where("cust_id = ?", customerId).
		Limit(dataFilter.Limit).
		Offset((dataFilter.Page - 1) * dataFilter.Limit)

	if picklistFilter.Driver != "" {
		db = db.Where("driver = ?", picklistFilter.Driver)
	}

	if picklistFilter.Vehicle != "" {
		db = db.Where("vehicle = ?", picklistFilter.Vehicle)
	}

	if picklistFilter.StartDate != "" {
		if t := helper.ParseWIBDateOnly(picklistFilter.StartDate, true); t != nil {
			db = db.Where("created_at >= ?", *t)
		}
	}

	if picklistFilter.EndDate != "" {
		if t := helper.ParseWIBDateOnly(picklistFilter.EndDate, false); t != nil {
			db = db.Where("created_at <= ?", *t)
		}
	}

	db = db.Order("created_at DESC").Find(&data)

	if db.Error != nil {
		log.Println(db.Error)
		return nil, db.Error // Return nil and the error if an error occurs
	}
	return data, nil // Return the data and no error
}

func (repo *PicklistRepositoryImpl) FindByID(ctx context.Context, id string) (model.Picklist, error) {
	var data model.Picklist
	db := repo.db.WithContext(ctx)
	db.Model(&model.Picklist{}).Where("picklist_no = ?", id).First(&data)

	if db.Error != nil {
		return data, db.Error
	}
	return data, nil
}

func (repo *PicklistRepositoryImpl) InsertPicklist(ctx context.Context, tx *gorm.DB, data model.Picklist) error {
	if err := tx.Model(&data).Create(&data).Error; err != nil {
		return fmt.Errorf("failed to insert picklist: %w", err)
	}
	return nil
}

func (repo *PicklistRepositoryImpl) UpdatePicklist(ctx context.Context, data model.Picklist) error {
	db := repo.db.WithContext(ctx)
	db.Model(&data).Updates(&data)

	if db.Error != nil {
		return fmt.Errorf("failed to update picklist: %w", db.Error)
	}

	return nil
}

func (repo *PicklistRepositoryImpl) DeletePicklist(ctx context.Context, id string) error {
	db := repo.db.WithContext(ctx)
	db.Model(&model.Picklist{}).Where("picklist_no = ?", id).Delete(&model.Picklist{})

	if db.Error != nil {
		return fmt.Errorf("failed to delete picklist: %w", db.Error)
	}

	return nil
}

func (repo *PicklistRepositoryImpl) InsertOrderPicklist(ctx context.Context, tx *gorm.DB, data model.OrderPicklist) (resp int, err error) {
	if err := tx.Model(&data).Create(&data).Error; err != nil {
		return resp, fmt.Errorf("failed to insert order picklist: %w", err)
	}
	return int(data.ID), nil
}

func (repo *PicklistRepositoryImpl) InsertOrderProduct(ctx context.Context, tx *gorm.DB, data model.OrderProduct) error {
	if err := tx.Model(&data).Create(&data).Error; err != nil {
		return fmt.Errorf("failed to insert order product: %w", err)
	}
	return nil
}

func (repo *PicklistRepositoryImpl) GetOrdersByPicklistNo(ctx context.Context, picklistNo string) ([]model.OrderPicklist, error) {
	var orders []model.OrderPicklist
	db := repo.db.WithContext(ctx)
	db.Model(&model.OrderPicklist{}).Where("picklist_no = ?", picklistNo).Find(&orders)

	if db.Error != nil {
		return nil, db.Error
	}
	return orders, nil
}

func (repo *PicklistRepositoryImpl) GetProductsByOrderID(ctx context.Context, orderID uint) ([]model.OrderProduct, error) {
	var products []model.OrderProduct
	db := repo.db.WithContext(ctx)
	db.Model(&model.OrderProduct{}).Where("order_id = ?", orderID).Find(&products)

	if db.Error != nil {
		return nil, db.Error
	}
	return products, nil
}

func (repo *PicklistRepositoryImpl) GetNextSequence(ctx context.Context, yy, mm, dd string) (int, error) {
	var count int64
	datePrefix := fmt.Sprintf("PL%s%s%s", yy, mm, dd)
	db := repo.db.WithContext(ctx)
	db.Model(&model.Picklist{}).Where("picklist_no LIKE ?", datePrefix+"%").Count(&count)

	if db.Error != nil {
		return 0, db.Error
	}

	return int(count) + 1, nil
}

func (repo *PicklistRepositoryImpl) FindAllByInvoiceNo(ctx context.Context, customerId string) []string {
	var order_no []string

	result := repo.db.WithContext(ctx).
		Table("picklist.order_picklist").
		Select("order_no").
		Where("cust_id = ?", customerId). // 👈 filter by customerId
		Debug().
		Find(&order_no)

	helper.ErrorPanic(result.Error)
	return order_no

}

func (repo *PicklistRepositoryImpl) GetTotalCount(ctx context.Context, picklistFilter entity.PicklistFilter) (int, error) {
	var count int64
	db := repo.db.Debug().WithContext(ctx).Model(&model.Picklist{})

	if picklistFilter.Driver != "" {
		db = db.Where("driver = ?", picklistFilter.Driver)
	}

	if picklistFilter.Vehicle != "" {
		db = db.Where("vehicle = ?", picklistFilter.Vehicle)
	}

	if picklistFilter.StartDate != "" {
		db = db.Where("created_at >= ?", picklistFilter.StartDate)
	}

	if picklistFilter.EndDate != "" {
		db = db.Where("created_at <= ?", picklistFilter.EndDate)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *PicklistRepositoryImpl) CountAll(ctx context.Context, filter entity.PicklistFilter, customerId string) (int64, error) {
	var count int64
	query := r.db.Model(&model.Picklist{}).Where("cust_id = ?", customerId)

	// Apply filters if they exist
	if filter.Driver != "" {
		query = query.Where("driver LIKE ?", "%"+filter.Driver+"%")
	}
	if filter.Vehicle != "" {
		query = query.Where("vehicle LIKE ?", "%"+filter.Vehicle+"%")
	}

	if filter.StartDate != "" {
		if t := helper.ParseWIBDateOnly(filter.StartDate, true); t != nil {
			query = query.Where("created_at >= ?", *t)
		}
	}

	if filter.EndDate != "" {
		if t := helper.ParseWIBDateOnly(filter.EndDate, false); t != nil {
			query = query.Where("created_at <= ?", *t)
		}
	}
	// Add other filters as needed

	err := query.Count(&count).Error
	return count, err
}
