package service

import (
	"context"
	"errors"
	"finance/entity"
	"finance/model"
	"finance/pkg/structs"
	"finance/repository"
	"time"

	"gorm.io/gorm"
)

var (
	ErrExpenseTypeExists   = errors.New("data already exists")
	ErrExpenseTypeNotFound = errors.New("record not found")
)

type ExpenseService interface {
	List(dataFilter entity.ExpenseQueryFilter) (data []entity.ExpenseListResponse, total int64, lastPage int, err error)
	Store(request entity.CreateExpenseBody, userId int64) error
	Update(parentCustId string, expenseTypeId int, request entity.UpdateExpenseBody, userId int64) error
	Delete(parentCustId string, expenseTypeId int, deletedBy int64) error
}

type expenseServiceImpl struct {
	Repository  repository.ExpenseRepository
	Transaction repository.Dbtransaction
}

func NewExpenseService(expenseRepository repository.ExpenseRepository, transaction repository.Dbtransaction) *expenseServiceImpl {
	return &expenseServiceImpl{
		Repository:  expenseRepository,
		Transaction: transaction,
	}
}

func (service *expenseServiceImpl) List(dataFilter entity.ExpenseQueryFilter) ([]entity.ExpenseListResponse, int64, int, error) {
	ctx := context.Background()

	expenseTypes, total, lastPage, err := service.Repository.FindAll(ctx, dataFilter)
	if err != nil {
		return nil, 0, 0, err
	}

	response := make([]entity.ExpenseListResponse, 0, len(expenseTypes))
	for _, et := range expenseTypes {
		etResp := entity.ExpenseListResponse{
			ExpenseTypeID:   et.ExpenseTypeID,
			ExpenseTypeCode: et.ExpenseTypeCode,
			ExpenseTypeName: et.ExpenseTypeName,
			Status:          et.IsActive,
			UpdateBy:        0,
			UpdateByName:    "",
			UpdateDate:      "",
		}

		if et.UpdatedBy != nil {
			etResp.UpdateBy = *et.UpdatedBy
		}
		if et.UpdatedByName != nil {
			etResp.UpdateByName = *et.UpdatedByName
		}
		if et.UpdatedAt != nil {
			etResp.UpdateDate = et.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		response = append(response, etResp)
	}

	return response, total, lastPage, nil
}

func (service *expenseServiceImpl) Store(request entity.CreateExpenseBody, userId int64) error {
	ctx := context.Background()

	// Check duplicate: expense_type_code and expense_type_name cannot be the same
	_, err := service.Repository.FindByCodeAndName(ctx, request.ParentCustId, request.ExpenseTypeCode, request.ExpenseTypeName)
	if err == nil {
		return ErrExpenseTypeExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to save data, please try again")
	}

	var expenseTypeModel model.ExpenseType
	if err := structs.Automapper(request, &expenseTypeModel); err != nil {
		return errors.New("failed to save data, please try again")
	}

	expenseTypeModel.CreatedBy = int(userId)
	expenseTypeModel.CustID = request.ParentCustId
	expenseTypeModel.IsDel = false

	if err := service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		return service.Repository.Store(txCtx, &expenseTypeModel)
	}); err != nil {
		if errors.Is(err, ErrExpenseTypeExists) {
			return err
		}
		return errors.New("failed to save data, please try again")
	}

	return nil
}

func (service *expenseServiceImpl) Update(parentCustId string, expenseTypeId int, request entity.UpdateExpenseBody, userId int64) error {
	ctx := context.Background()

	// Check if expense type exists
	_, err := service.Repository.FindById(ctx, parentCustId, expenseTypeId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrExpenseTypeNotFound
	}
	if err != nil {
		return errors.New("failed to update data, please try again")
	}

	// Check duplicate: expense_type_code and expense_type_name cannot be the same (exclude current record)
	duplicate, err := service.Repository.FindByCodeAndName(ctx, parentCustId, request.ExpenseTypeCode, request.ExpenseTypeName)
	if err == nil && duplicate.ExpenseTypeID != expenseTypeId {
		return ErrExpenseTypeExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to update data, please try again")
	}

	now := time.Now()

	// Prepare update data
	updateData := map[string]interface{}{
		"expense_type_code": request.ExpenseTypeCode,
		"expense_type_name": request.ExpenseTypeName,
		"updated_by":        userId,
		"updated_at":        now,
	}

	// Handle IsActive: only update if provided in request
	if request.IsActive != nil {
		updateData["is_active"] = *request.IsActive
	}

	if err := service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		return service.Repository.Update(txCtx, parentCustId, expenseTypeId, updateData)
	}); err != nil {
		if errors.Is(err, ErrExpenseTypeExists) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrExpenseTypeNotFound
		}
		return errors.New("failed to update data, please try again")
	}

	return nil
}

func (service *expenseServiceImpl) Delete(parentCustId string, expenseTypeId int, deletedBy int64) error {
	ctx := context.Background()

	if err := service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		return service.Repository.Delete(txCtx, parentCustId, expenseTypeId, deletedBy)
	}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrExpenseTypeNotFound
		}
		return errors.New("failed to delete data, please try again")
	}

	return nil
}
