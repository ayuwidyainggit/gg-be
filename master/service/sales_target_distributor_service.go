package service

import (
	"fmt"
	"master/entity"
	"master/model"
	"master/pkg/structs"
	"master/repository"
	"time"

	"github.com/jmoiron/sqlx"
)

type SalesTargetDistributorService interface {
	List(entity.SalesTargetDistributorQueryFilter, string) (data []entity.SalesTargetDistributorListResponse, total int, lastPage int, err error)
	Detail(int, string) (entity.SalesTargetDistributorDetailResponse, error)
	Store(entity.CreateSalesTargetDistributorBody) (int, error)
	Update(int, entity.UpdateSalesTargetDistributorRequest) error
}

type salesTargetDistributorServiceImpl struct {
	SalesTargetDistributorRepository repository.SalesTargetDistributorRepository
	SalesTargetRepository            repository.SalesTargetRepository
}

// NewSalesTargetDistributorService creates a new instance of salesTargetDistributorServiceImpl
func NewSalesTargetDistributorService(salesTargetDistributorRepository repository.SalesTargetDistributorRepository, salesTargetRepository repository.SalesTargetRepository) *salesTargetDistributorServiceImpl {
	return &salesTargetDistributorServiceImpl{
		SalesTargetDistributorRepository: salesTargetDistributorRepository,
		SalesTargetRepository:            salesTargetRepository,
	}
}

func currentTimeUTC() time.Time {
	return time.Now().UTC()
}

func modelTimeNowUTC() time.Time {
	return time.Now().UTC()
}

func buildMonthlyAllocationMap(rows []model.SalesTargetMonthlyAllocation) map[int]model.SalesTargetMonthlyAllocation {
	result := make(map[int]model.SalesTargetMonthlyAllocation, len(rows))
	for _, row := range rows {
		result[row.Month] = row
	}

	return result
}

func buildMonthlyDetailMap(rows []model.SalesTargetDistributorMonthly) map[int]model.SalesTargetDistributorMonthly {
	result := make(map[int]model.SalesTargetDistributorMonthly, len(rows))
	for _, row := range rows {
		result[row.Month] = row
	}

	return result
}

func computeMonthlyFlags(year int, month int, currentTime time.Time, allocatedTotal int64) (bool, bool) {
	isPastMonth := year < currentTime.Year() || (year == currentTime.Year() && month < int(currentTime.Month()))
	isAllocated := allocatedTotal > 0

	return isPastMonth, !isPastMonth && !isAllocated
}

func ensureMonthlyTargetsConsistent(monthlyTarget int, allocatedTotal int64, month int) error {
	if allocatedTotal > int64(monthlyTarget) {
		return fmt.Errorf("month %d cannot be updated because allocated total (%d) exceeds monthly target (%d)", month, allocatedTotal, monthlyTarget)
	}

	return nil
}

func (service *salesTargetDistributorServiceImpl) syncMonthlyTargets(tx *sqlx.Tx, yearlyId int, syncCustID string, request entity.UpdateSalesTargetDistributorRequest, existingMonthly map[int]model.SalesTargetDistributorMonthly, allocations map[int]model.SalesTargetMonthlyAllocation) error {
	timeNow := currentTimeUTC()

	for _, monthly := range request.Data {
		allocation := allocations[monthly.Month]
		if err := ensureMonthlyTargetsConsistent(monthly.MonthlyTarget, allocation.AllocatedTotal, monthly.Month); err != nil {
			return err
		}

		monthlyRow, exists := existingMonthly[monthly.Month]
		if !exists {
			monthlyRow = model.SalesTargetDistributorMonthly{
				CustId:                         request.CustId,
				SalesTargetDistributorYearlyId: yearlyId,
				Month:                          monthly.Month,
				MonthlyTarget:                  monthly.MonthlyTarget,
				IsActive:                       true,
				CreatedBy:                      request.UpdatedBy,
				CreatedAt:                      timeNow,
				IsDel:                          false,
			}

			monthlyID, err := service.SalesTargetDistributorRepository.StoreMonthly(tx, monthlyRow)
			if err != nil {
				return err
			}

			monthlyRow.SalesTargetDistributorMonthlyId = monthlyID
		} else {
			if err := service.SalesTargetDistributorRepository.UpdateMonthlyTarget(tx, monthlyRow.SalesTargetDistributorMonthlyId, monthly.MonthlyTarget, request.UpdatedBy); err != nil {
				return err
			}
		}

		if err := service.SalesTargetRepository.SyncTargetsToMonthly(tx, syncCustID, yearlyId, monthly.Month, monthlyRow.SalesTargetDistributorMonthlyId, monthly.MonthlyTarget, request.UpdatedBy); err != nil {
			return err
		}
	}

	return nil
}

func deriveSalesTargetDistributorStatus(rawStatus int, year int, isActive bool, currentYear int) string {
	if entity.SalesTargetStatus(rawStatus) == entity.SALES_TARGET_STATUS_DRAFT {
		return entity.SALES_TARGET_STATUS_DRAFT.String()
	}

	if year > currentYear || !isActive {
		return entity.SALES_TARGET_STATUS_INACTIVE.String()
	}

	return entity.SALES_TARGET_STATUS_ACTIVE.String()
}

func applyStatusTransitionMetadata(request *entity.UpdateSalesTargetDistributorRequest, timeNow time.Time) {
	if request.Status == nil {
		return
	}

	status := entity.SalesTargetStatus(*request.Status)
	if status == entity.SALES_TARGET_STATUS_INACTIVE {
		userInactive := request.UpdatedBy
		request.UserInactive = &userInactive
		inactiveAt := timeNow
		request.InactiveAt = &inactiveAt

		isActive := false
		request.IsActive = &isActive
		return
	}

	if status == entity.SALES_TARGET_STATUS_ACTIVE {
		request.UserInactive = nil
		request.InactiveAt = nil

		isActive := true
		request.IsActive = &isActive
	}
}

// List returns a list of yearly sales targets with status calculations
func (service *salesTargetDistributorServiceImpl) List(dataFilter entity.SalesTargetDistributorQueryFilter, custId string) (data []entity.SalesTargetDistributorListResponse, total int, lastPage int, err error) {
	data = []entity.SalesTargetDistributorListResponse{}
	list, total, lastPage, err := service.SalesTargetDistributorRepository.FindAllByCustId(dataFilter, custId)
	if err != nil {
		return data, total, lastPage, err
	}

	currentYear := currentTimeUTC().Year()

	for _, row := range list {
		var vResp entity.SalesTargetDistributorListResponse

		// Manual mapping to avoid type mismatch on Status field (int in model vs string in entity)
		vResp.SalesTargetDistributorYearlyId = row.SalesTargetDistributorYearlyId
		vResp.DistributorId = row.DistributorId
		vResp.DistributorCode = row.DistributorCode
		vResp.DistributorName = row.DistributorName
		vResp.Year = row.Year
		vResp.YearlyTarget = row.YearlyTarget
		vResp.UserInactive = row.UserInactive
		vResp.InactiveAt = row.InactiveAt

		vResp.Status = deriveSalesTargetDistributorStatus(row.Status, row.Year, row.IsActive, currentYear)

		// Audit column fallback logic:
		// If data has not been edited (updated_at is NULL), use created_by/created_at
		if row.UpdatedAt == nil {
			vResp.UpdatedBy = row.CreatedBy
			vResp.UpdatedAt = &row.CreatedAt
		} else {
			if row.UpdatedBy != nil {
				vResp.UpdatedBy = *row.UpdatedBy
			}
			vResp.UpdatedAt = row.UpdatedAt
		}
		// UpdatedByName is already handled by CASE WHEN in the repository query
		// It will automatically use created_by's name when updated_by is NULL
		vResp.UpdatedByName = row.UpdatedByName

		data = append(data, vResp)
	}

	return data, total, lastPage, nil
}

// Detail returns the full details of a yearly sales target, including its monthly components
func (service *salesTargetDistributorServiceImpl) Detail(id int, custId string) (response entity.SalesTargetDistributorDetailResponse, err error) {
	if service.SalesTargetRepository == nil {
		return response, fmt.Errorf("sales target repository is required")
	}

	yearly, err := service.SalesTargetDistributorRepository.FindOneByIdAndCustId(id, custId)
	if err != nil {
		return response, err
	}

	// Manual mapping to avoid type mismatch on Status field
	response.SalesTargetDistributorYearlyId = yearly.SalesTargetDistributorYearlyId
	response.Year = yearly.Year
	response.YearlyTarget = yearly.YearlyTarget
	response.AreaId = yearly.AreaId
	response.AreaCode = yearly.AreaCode
	response.AreaName = yearly.AreaName
	response.RegionId = yearly.RegionId
	response.RegionCode = yearly.RegionCode
	response.RegionName = yearly.RegionName
	response.DistributorId = yearly.DistributorId
	response.DistributorCode = yearly.DistributorCode
	response.DistributorName = yearly.DistributorName
	response.UserInactive = yearly.UserInactive
	response.InactiveAt = yearly.InactiveAt

	currentYear := currentTimeUTC().Year()
	response.Status = deriveSalesTargetDistributorStatus(yearly.Status, yearly.Year, yearly.IsActive, currentYear)

	// Audit column fallback logic:
	// If data has not been edited (updated_at is NULL), use created_by/created_at
	if yearly.UpdatedAt == nil {
		response.UpdatedBy = yearly.CreatedBy
		response.UpdatedAt = &yearly.CreatedAt
	} else {
		if yearly.UpdatedBy != nil {
			response.UpdatedBy = *yearly.UpdatedBy
		}
		response.UpdatedAt = yearly.UpdatedAt
	}
	// UpdatedByName is already handled by CASE WHEN in the repository query
	// It will automatically use created_by's name when updated_by is NULL
	response.UpdatedByName = yearly.UpdatedByName

	// Fetch details
	monthlyList, err := service.SalesTargetDistributorRepository.FindMonthlyDetailsByYearlyId(id)
	if err != nil {
		return response, err
	}

	distributorChildCustID, err := service.SalesTargetDistributorRepository.FindChildCustIDByDistributorID(yearly.DistributorId)
	if err != nil {
		return response, err
	}

	monthlyAllocations, err := service.SalesTargetRepository.FindMonthlyAllocationByYearlyId(id, distributorChildCustID)
	if err != nil {
		return response, err
	}

	allocationMap := buildMonthlyAllocationMap(monthlyAllocations)
	currentTime := currentTimeUTC()

	for _, row := range monthlyList {
		var detail entity.SalesTargetDistributorMonthlyDetail
		err = structs.Automapper(row, &detail)
		if err != nil {
			return response, err
		}

		allocation := allocationMap[row.Month]
		detail.AllocatedTotal = allocation.AllocatedTotal
		detail.Remaining = int64(row.MonthlyTarget) - allocation.AllocatedTotal
		detail.IsAllocated = allocation.TargetCount > 0 || allocation.AllocatedTotal > 0
		detail.IsPastMonth, detail.IsEditable = computeMonthlyFlags(response.Year, row.Month, currentTime, allocation.AllocatedTotal)
		if detail.IsPastMonth {
			detail.DisableReason = "past_month"
		} else if detail.IsAllocated {
			detail.DisableReason = "allocated"
		}
		response.Details = append(response.Details, detail)
	}

	if response.Details == nil {
		response.Details = []entity.SalesTargetDistributorMonthlyDetail{}
	}

	allocationSummary, err := service.SalesTargetDistributorRepository.FindAllocationSummaryByYearlyId(id, custId)
	if err != nil {
		// Non-fatal for detail endpoint: fallback to safe default
		response.IsAllocated = false
		response.AllocationTotal = 0
	} else {
		response.IsAllocated = allocationSummary.IsAllocated
		response.AllocationTotal = allocationSummary.AllocatedTotal
	}

	return response, nil
}

// Store saves a new yearly sales target and its associated monthly targets in a single transaction
func (service *salesTargetDistributorServiceImpl) Store(request entity.CreateSalesTargetDistributorBody) (int, error) {
	tx, err := service.SalesTargetDistributorRepository.BeginTx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	timeNow := currentTimeUTC()

	yearlyData := model.SalesTargetDistributorYearly{
		CustId:        request.CustId,
		AreaId:        request.AreaId,
		RegionId:      request.RegionId,
		DistributorId: request.DistributorId,
		Year:          request.Year,
		YearlyTarget:  request.YearlyTarget,
		Status:        *request.Status,
		IsActive:      true,
		CreatedBy:     request.CreatedBy,
		CreatedAt:     timeNow,
		IsDel:         false,
	}

	yearlyId, err := service.SalesTargetDistributorRepository.StoreYearly(tx, yearlyData)
	if err != nil {
		return 0, err
	}

	for _, monthly := range request.Data {
		monthlyData := model.SalesTargetDistributorMonthly{
			CustId:                         request.CustId,
			SalesTargetDistributorYearlyId: yearlyId,
			Month:                          monthly.Month,
			MonthlyTarget:                  monthly.MonthlyTarget,
			IsActive:                       true,
			CreatedBy:                      request.CreatedBy,
			CreatedAt:                      timeNow,
			IsDel:                          false,
		}
		_, err = service.SalesTargetDistributorRepository.StoreMonthly(tx, monthlyData)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return yearlyId, nil
}

// Update modifies a yearly sales target and replaces its monthly details if provided, within a transaction
func (service *salesTargetDistributorServiceImpl) Update(id int, request entity.UpdateSalesTargetDistributorRequest) error {
	tx, err := service.SalesTargetDistributorRepository.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	timeNow := currentTimeUTC()

	applyStatusTransitionMetadata(&request, timeNow)

	// Update Yearly
	// TODO: Refactor to use Transaction.WithinTransaction pattern when available
	err = service.SalesTargetDistributorRepository.UpdateYearly(tx, id, request)
	if err != nil {
		return err
	}

	if len(request.Data) > 0 {
		yearly, err := service.SalesTargetDistributorRepository.FindOneByIdAndCustId(id, request.CustId)
		if err != nil {
			return err
		}

		existingMonthly, err := service.SalesTargetDistributorRepository.FindMonthlyDetailsByYearlyId(id)
		if err != nil {
			return err
		}

		resolvedDistributorID := yearly.DistributorId
		if request.DistributorId != nil {
			resolvedDistributorID = *request.DistributorId
		}

		syncCustID, err := service.SalesTargetDistributorRepository.FindChildCustIDByDistributorID(resolvedDistributorID)
		if err != nil {
			return err
		}

		monthlyAllocations, err := service.SalesTargetRepository.FindMonthlyAllocationByYearlyId(id, syncCustID)
		if err != nil {
			return err
		}

		err = service.syncMonthlyTargets(tx, id, syncCustID, request, buildMonthlyDetailMap(existingMonthly), buildMonthlyAllocationMap(monthlyAllocations))
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
