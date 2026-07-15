package service

import (
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/pkg/structs"
	"master/repository"
)

func deriveSalesTargetStatus(status int, year int, month int) string {
	if status == entity.StatusDraft {
		return "Draft"
	}

	now := currentTimeUTC()
	if year != now.Year() || month != int(now.Month()) {
		return "Inactive"
	}

	switch status {
	case entity.StatusActive:
		return "Active"
	case entity.StatusInactive:
		return "Inactive"
	default:
		return "Unknown"
	}
}

type SalesTargetService interface {
	List(filter entity.SalesTargetQueryFilter, custId string) (data []entity.SalesTargetListResponse, total int, lastPage int, err error)
	Detail(salesTargetId int64, custId string) (entity.SalesTargetDetailResponse, error)
	MonthlyDistributor(query entity.SalesTargetMonthlyDistQuery) (entity.SalesTargetMonthlyDistResp, error)
	Store(request entity.CreateSalesTargetRequest) error
	Update(salesTargetId int64, request entity.UpdateSalesTargetRequest) error
}

type salesTargetServiceImpl struct {
	SalesTargetRepository repository.SalesTargetRepository
}

func NewSalesTargetService(salesTargetRepository repository.SalesTargetRepository) SalesTargetService {
	return &salesTargetServiceImpl{
		SalesTargetRepository: salesTargetRepository,
	}
}

func (service *salesTargetServiceImpl) List(filter entity.SalesTargetQueryFilter, custId string) (data []entity.SalesTargetListResponse, total int, lastPage int, err error) {
	salesTargets, total, lastPage, err := service.SalesTargetRepository.FindAll(filter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	data = make([]entity.SalesTargetListResponse, 0, len(salesTargets))
	for _, row := range salesTargets {
		resp := entity.SalesTargetListResponse{
			SalesTargetId:  row.SalesTargetId,
			Year:           row.Year,
			Month:          row.Month,
			AllocatedTotal: row.AllocatedTotal,
			MonthlyTarget:  row.MonthlyTarget,
			Remaining:      row.Remaining,
			Status:         deriveSalesTargetStatus(row.Status, row.Year, row.Month),
		}

		// Handle updated_by logic: fallback to creator if not yet updated
		resp.UpdatedAt = row.CreatedAt
		if row.CreatedBy != nil {
			resp.UpdatedBy = *row.CreatedBy
		}

		if row.UpdatedAt != nil {
			resp.UpdatedAt = row.UpdatedAt
			if row.UpdatedBy != nil {
				resp.UpdatedBy = *row.UpdatedBy
			}
		}

		data = append(data, resp)
	}

	return data, total, lastPage, nil
}

func (service *salesTargetServiceImpl) Detail(salesTargetId int64, custId string) (entity.SalesTargetDetailResponse, error) {
	var response entity.SalesTargetDetailResponse

	// Get sales target header
	salesTarget, err := service.SalesTargetRepository.FindOneById(salesTargetId, custId)
	if err != nil {
		return response, err
	}

	response = entity.SalesTargetDetailResponse{
		SalesTargetId:  salesTarget.SalesTargetId,
		Year:           salesTarget.Year,
		Month:          salesTarget.Month,
		AllocatedTotal: salesTarget.AllocatedTotal,
		MonthlyTarget:  salesTarget.MonthlyTarget,
		Remaining:      salesTarget.Remaining,
		Status:         deriveSalesTargetStatus(salesTarget.Status, salesTarget.Year, salesTarget.Month),
	}

	// Handle updated_by logic: fallback to creator if not yet updated
	response.UpdatedAt = salesTarget.CreatedAt
	if salesTarget.CreatedBy != nil {
		response.UpdatedBy = *salesTarget.CreatedBy
	}

	if salesTarget.UpdatedAt != nil {
		response.UpdatedAt = salesTarget.UpdatedAt
		if salesTarget.UpdatedBy != nil {
			response.UpdatedBy = *salesTarget.UpdatedBy
		}
	}

	// Get sales allocated details
	details, err := service.SalesTargetRepository.FindDetailsBySalesTargetId(salesTargetId, custId)
	if err != nil {
		return response, err
	}

	for _, detail := range details {
		var detailResp entity.SalesAllocatedDetailResp
		err = structs.Automapper(detail, &detailResp)
		if err != nil {
			return response, err
		}

		// Handle nullable fields
		if detail.DistributorCode != nil {
			detailResp.DistributorCode = *detail.DistributorCode
		}
		if detail.DistributorName != nil {
			detailResp.DistributorName = *detail.DistributorName
		}
		if detail.SalesTeamCode != nil {
			detailResp.SalesTeamCode = *detail.SalesTeamCode
		}
		if detail.SalesTeamName != nil {
			detailResp.SalesTeamName = *detail.SalesTeamName
		}
		if detail.ChannelCode != nil {
			detailResp.ChannelCode = *detail.ChannelCode
		}
		if detail.ChannelName != nil {
			detailResp.ChannelName = *detail.ChannelName
		}

		response.Details = append(response.Details, detailResp)
	}

	// Initialize empty array if no details
	if response.Details == nil {
		response.Details = []entity.SalesAllocatedDetailResp{}
	}

	return response, nil
}

func (service *salesTargetServiceImpl) MonthlyDistributor(query entity.SalesTargetMonthlyDistQuery) (entity.SalesTargetMonthlyDistResp, error) {
	var response entity.SalesTargetMonthlyDistResp

	data, err := service.SalesTargetRepository.FindMonthlyDistributor(query)
	if err != nil {
		return response, err
	}

	err = structs.Automapper(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (service *salesTargetServiceImpl) Store(request entity.CreateSalesTargetRequest) error {
	if request.Status == nil {
		return fmt.Errorf("status is required")
	}

	if request.Remaining == nil {
		return fmt.Errorf("remaining is required")
	}

	if len(request.Data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	// Validate allocated_total matches sum of allocated
	var sumAllocated int64
	for _, item := range request.Data {
		sumAllocated += item.Allocated
	}
	if sumAllocated != request.AllocatedTotal {
		return fmt.Errorf("allocated_total (%d) must equal sum of allocated (%d)", request.AllocatedTotal, sumAllocated)
	}

	service.SalesTargetRepository.TrxBegin()

	defer func() {
		if r := recover(); r != nil {
			service.SalesTargetRepository.TrxRollback()
			panic(r)
		}
	}()

	timeNow := modelTimeNowUTC()

	// Prepare sales target model
	salesTarget := model.SalesTarget{
		CustId:                          request.CustId,
		SalesTargetDistributorYearlyId:  request.SalesTargetDistributorYearlyId,
		SalesTargetDistributorMonthlyId: request.SalesTargetDistributorMonthlyId,
		Month:                           request.Month,
		Year:                            request.Year,
		AllocatedTotal:                  request.AllocatedTotal,
		MonthlyTarget:                   request.MonthlyTarget,
		Remaining:                       *request.Remaining,
		Status:                          *request.Status,
		CreatedBy:                       request.CreatedBy,
		CreatedAt:                       timeNow,
		IsDel:                           false,
	}

	// Store sales target
	salesTargetId, err := service.SalesTargetRepository.Store(salesTarget)
	if err != nil {
		service.SalesTargetRepository.TrxRollback()
		return err
	}

	// Store sales allocated details
	for _, item := range request.Data {
		salesAllocated := model.SalesAllocated{
			CustId:        request.CustId,
			SalesTargetId: salesTargetId,
			SalesmanId:    item.SalesmanId,
			SalesTeamId:   &item.SalesTeamId,
			Allocated:     item.Allocated,
			IsActive:      true,
			CreatedBy:     request.CreatedBy,
			CreatedAt:     timeNow,
			IsDel:         false,
		}

		err = service.SalesTargetRepository.StoreAllocated(salesAllocated)
		if err != nil {
			service.SalesTargetRepository.TrxRollback()
			return err
		}
	}

	err = service.SalesTargetRepository.TrxCommit()
	if err != nil {
		service.SalesTargetRepository.TrxRollback()
		return err
	}

	return nil
}

func (service *salesTargetServiceImpl) Update(salesTargetId int64, request entity.UpdateSalesTargetRequest) error {
	// Check exists
	_, err := service.SalesTargetRepository.FindOneById(salesTargetId, request.CustId)
	if err != nil {
		return constant.ErrSalesTargetNotFound
	}

	updates := make(map[string]interface{})
	timeNow := modelTimeNowUTC()

	// Build updates map for non-nil fields only
	if request.SalesTargetDistributorYearlyId != nil {
		updates["sales_target_distributor_yearly_id"] = *request.SalesTargetDistributorYearlyId
	}
	if request.SalesTargetDistributorMonthlyId != nil {
		updates["sales_target_distributor_monthly_id"] = *request.SalesTargetDistributorMonthlyId
	}
	if request.Month != nil {
		updates["month"] = *request.Month
	}
	if request.Year != nil {
		updates["year"] = *request.Year
	}
	if request.AllocatedTotal != nil {
		updates["allocated_total"] = *request.AllocatedTotal
	}
	if request.MonthlyTarget != nil {
		updates["monthly_target"] = *request.MonthlyTarget
	}
	if request.Remaining != nil {
		updates["remaining"] = *request.Remaining
	}
	if request.Status != nil {
		updates["status"] = *request.Status
	}

	// Always set updated_by and updated_at
	updates["updated_by"] = request.UpdatedBy
	updates["updated_at"] = timeNow

	// Start transaction only if we have Data to update
	if request.Data != nil {
		// Validate allocated sum if both Data and AllocatedTotal provided
		if request.AllocatedTotal != nil {
			var sumAllocated int64
			for _, item := range request.Data {
				sumAllocated += item.Allocated
			}
			if sumAllocated != *request.AllocatedTotal {
				return fmt.Errorf("allocated_total (%d) must equal sum of allocated (%d)", *request.AllocatedTotal, sumAllocated)
			}
		}

		service.SalesTargetRepository.TrxBegin()
		defer func() {
			if r := recover(); r != nil {
				service.SalesTargetRepository.TrxRollback()
				panic(r)
			}
		}()

		// Update main record
		err = service.SalesTargetRepository.UpdatePartial(salesTargetId, request.CustId, updates)
		if err != nil {
			service.SalesTargetRepository.TrxRollback()
			return err
		}

		// Delete and recreate allocated
		err = service.SalesTargetRepository.DeleteAllocatedByTargetId(salesTargetId, request.CustId)
		if err != nil {
			service.SalesTargetRepository.TrxRollback()
			return err
		}

		// Store new allocated details
		for _, item := range request.Data {
			salesAllocated := model.SalesAllocated{
				CustId:        request.CustId,
				SalesTargetId: salesTargetId,
				SalesmanId:    item.SalesmanId,
				SalesTeamId:   &item.SalesTeamId,
				Allocated:     item.Allocated,
				IsActive:      true,
				CreatedBy:     request.UpdatedBy,
				CreatedAt:     timeNow,
				IsDel:         false,
			}

			err = service.SalesTargetRepository.StoreAllocated(salesAllocated)
			if err != nil {
				service.SalesTargetRepository.TrxRollback()
				return err
			}
		}

		return service.SalesTargetRepository.TrxCommit()
	}

	// Simple case: no Data, just update fields (no transaction needed)
	return service.SalesTargetRepository.UpdatePartial(salesTargetId, request.CustId, updates)
}
