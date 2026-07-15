package service

import (
	"context"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/model"
	"inventory/pkg/constant"
	"inventory/pkg/str"
	"inventory/repository"
	"math"
	"strings"
	"time"
)

const (
	TypeFinal         = "final"
	TypeReplenishment = "replenishment"
)

const (
	StatusNeedReview = 1 // Need Review - Replenishment Order sudah submit oleh Admin Distributor untuk diapprove oleh user di level Principal
	StatusApproved   = 2 // Approved - Replenishment Order sudah Di Approve oleh User yang memiliki akses ke menu Approval Replenishment Order
	StatusRejected   = 3 // Rejected - Replenishment Order direject di menu Replenishment Order Approval
	StatusOnDelivery = 4 // On Delivery - Status yang didapatkan ketika ERP/Sistem yang berelasi mem-post No PO, No SO, No DO, Delivery date, Vehicle Number, code-Product name, QTY, Price kemudian mengirim balikan status ke Scylla X
	StatusCompleted  = 5 // Completed - NO PO sudah berelasi dengan No dokumen Goods Receipt
	StatusProcessed  = 6 // Processed - Status yang didapatkan ketika ERP/Sistem yang berelasi mem-post NoPO, No SO, code-Product name, QTY, Price dan mengirim balikan status ke Scylla X
	StatusCancelled  = 7 // Cancelled - Status yang didapatkan ketika ERP/Sistem yang berelasi mem-post rejection
)

const (
	StatusNameOnDelivery = "On Delivery"
	StatusNameGoodsIssue = "Goods Issue"
	StatusNameSubmitted  = "submitted"
)

type ReplenishmentService interface {
	List(dataFilter entity.ReplenishmentQueryFilter, custId, parentCustId string, isPrincipal bool) (data []entity.ReplenishmentListResponse, total int64, lastPage int, err error)
	Store(request entity.CreateReplenishmentOrderBody) error
	Detail(replenishmentID, typeParam, statusParam string, custID, parentCustID string, isPrincipal bool) (*entity.ReplenishmentDetailResponse, error)
	ProductList(dataFilter entity.ReplenishmentProductQueryFilter, custId, parentCustId string) (data []entity.ReplenishmentProductListResponse, total int64, lastPage int, err error)
	ProductGrList(dataFilter entity.ProductGrListQueryFilter, custId, parentCustId string) (*entity.ProductGrListResponse, int64, int, error)
	PoList(dataFilter entity.PoListQueryFilter, custId, parentCustId string) (data []entity.PoListResponse, total int64, lastPage int, err error)
	ApprovalProductList(dataFilter entity.ReplenishmentApprovalProductQueryFilter, custId, parentCustId string, isPrincipal bool) (data []entity.ReplenishmentApprovalProductResponse, total int64, lastPage int, err error)
	ApprovalList(dataFilter entity.ReplenishmentApprovalListQueryFilter, custId, parentCustId string, isPrincipal bool) (data []entity.ReplenishmentApprovalListResponse, total int64, lastPage int, err error)
	UpdateApproval(replenishmentID int64, request entity.UpdateReplenishmentApprovalRequest, custID string, userID int64, employeeID int64) error
	BatchUpdateApproval(request entity.BatchReplenishmentApprovalRequest, custID string, userID int64, employeeID int64) ([]entity.BatchReplenishmentApprovalResultItem, error)
	IsAnyApprovalPICUser(userID int64, employeeID int64) (bool, error)
	CheckIsPrincipal(custID string) (bool, error)
	SummarizeReplanishment(replanishmentIDs []int64, custID, parentCustID string, isPrincipal bool, requestID string, userID int64) ([]entity.SummarizeReplanishmentResponse, error)
	SAPUpdateReplenishmentStatus(req entity.SAPReplStatusRequest, updatedBy int64) entity.SAPReplStatusResponse
	SAPGetReplenishmentExport(query entity.SAPReplExportQuery) ([]entity.SAPReplExportItem, error)
}

func NewReplenishmentService(
	replenishmentRepository repository.ReplenishmentRepository,
	transaction repository.Dbtransaction,
) *replenishmentServiceImpl {
	return &replenishmentServiceImpl{
		ReplenishmentRepository: replenishmentRepository,
		Transaction:             transaction,
	}
}

type replenishmentServiceImpl struct {
	ReplenishmentRepository repository.ReplenishmentRepository
	Transaction             repository.Dbtransaction
}

func (service *replenishmentServiceImpl) List(dataFilter entity.ReplenishmentQueryFilter, custId, parentCustId string, isPrincipal bool) ([]entity.ReplenishmentListResponse, int64, int, error) {
	ctx := context.Background()
	isPICUser, err := service.ReplenishmentRepository.IsUserApprovalPIC(ctx, dataFilter.UserID)
	if err != nil {
		return nil, 0, 0, err
	}
	if !isPICUser && dataFilter.EmpID > 0 {
		isPICUser, err = service.ReplenishmentRepository.IsUserApprovalPIC(ctx, dataFilter.EmpID)
		if err != nil {
			return nil, 0, 0, err
		}
	}
	dataFilter.IsPICUser = isPICUser

	replenishments, total, lastPage, err := service.ReplenishmentRepository.FindAllByCustId(ctx, dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []entity.ReplenishmentListResponse
	for _, replenishment := range replenishments {
		response := service.mapToResponse(&replenishment)
		responses = append(responses, response)
	}

	return responses, total, lastPage, nil
}

func (service *replenishmentServiceImpl) ApprovalList(dataFilter entity.ReplenishmentApprovalListQueryFilter, custId, parentCustId string, isPrincipal bool) ([]entity.ReplenishmentApprovalListResponse, int64, int, error) {
	ctx := context.Background()
	rows, total, lastPage, err := service.ReplenishmentRepository.FindApprovalOrderList(ctx, dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		return nil, 0, 0, err
	}
	out := make([]entity.ReplenishmentApprovalListResponse, 0, len(rows))
	for i := range rows {
		out = append(out, service.mapToApprovalListResponse(&rows[i]))
	}
	return out, total, lastPage, nil
}

func (service *replenishmentServiceImpl) mapToApprovalListResponse(row *model.ReplenishmentOrderList) entity.ReplenishmentApprovalListResponse {
	var r entity.ReplenishmentApprovalListResponse
	r.ReplenishmentID = row.ReplenishmentID
	r.ReplenishmentNo = row.ReplenishmentNo
	if row.Date != nil {
		r.Date = row.Date.In(constant.AsiaJakartaLocation).Format(constant.DATE_FORMAT_DISPLAY)
	}
	if row.DeliveryDate != nil {
		s := row.DeliveryDate.In(constant.AsiaJakartaLocation).Format(constant.DATE_FORMAT_DISPLAY)
		r.DeliveryDate = &s
	}
	if row.SupID != nil {
		r.SupID = *row.SupID
	}
	if row.SupCode != nil {
		r.SupCode = *row.SupCode
	}
	if row.SupName != nil {
		r.SupName = *row.SupName
	}
	r.DistributorID = row.DistributorID
	if row.DistributorCode != nil {
		r.DistributorCode = *row.DistributorCode
	}
	if row.DistributorName != nil {
		r.DistributorName = *row.DistributorName
	}
	if row.Address != nil {
		r.Address = *row.Address
	}
	if row.WhID != nil {
		r.WhID = *row.WhID
	}
	if row.WhCode != nil {
		r.WhCode = *row.WhCode
	}
	if row.WhName != nil {
		r.WhName = *row.WhName
	}
	r.CreatedBy = row.CreatedBy
	if row.CreatedByName != nil {
		r.CreatedByName = *row.CreatedByName
	}
	if !row.CreatedAt.IsZero() {
		r.CreatedAt = row.CreatedAt.UTC().Format(time.RFC3339)
	}
	if row.UpdatedBy != nil {
		r.UpdatedBy = row.UpdatedBy
	}
	if row.UpdatedByName != nil {
		r.UpdatedByName = *row.UpdatedByName
	}
	if row.UpdatedAt != nil {
		u := row.UpdatedAt.UTC().Format(time.RFC3339)
		r.UpdatedAt = &u
	}
	return r
}

func (service *replenishmentServiceImpl) mapStatusDisplayName(statusName string) string {
	if strings.EqualFold(strings.TrimSpace(statusName), StatusNameOnDelivery) {
		return StatusNameGoodsIssue
	}
	return statusName
}

func (service *replenishmentServiceImpl) mapToResponse(replenishment *model.ReplenishmentOrderList) entity.ReplenishmentListResponse {
	var response entity.ReplenishmentListResponse

	// Map basic fields
	response.ReplenishmentID = service.int64ToString(replenishment.ReplenishmentID)
	response.ReplenishmentNo = replenishment.ReplenishmentNo

	// Format date to DD/MM/YYYY with timezone Indonesia
	if replenishment.Date != nil {
		dateInJkt := replenishment.Date.In(constant.AsiaJakartaLocation)
		response.Date = dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
	}

	// Delivery date
	if replenishment.DeliveryDate != nil {
		dateInJkt := replenishment.DeliveryDate.In(constant.AsiaJakartaLocation)
		response.DeliveryDate = dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
	}

	// Supplier info
	if replenishment.SupID != nil {
		response.SupID = service.int64ToString(*replenishment.SupID)
	}
	if replenishment.SupCode != nil {
		response.SupCode = *replenishment.SupCode
	}
	if replenishment.SupName != nil {
		response.SupName = *replenishment.SupName
	}

	// Warehouse info
	if replenishment.WhID != nil {
		response.WhID = service.int64ToString(*replenishment.WhID)
	}
	if replenishment.WhCode != nil {
		response.WhCode = *replenishment.WhCode
	}
	if replenishment.WhName != nil {
		response.WhName = *replenishment.WhName
	}

	// Status
	if replenishment.StatusName != nil {
		response.Status = service.mapStatusDisplayName(*replenishment.StatusName)
	} else {
		response.Status = service.int64ToString(int64(replenishment.Status))
	}

	// Created by
	response.CreatedBy = entity.ReplenishmentUserInfo{
		UserID:   replenishment.CreatedBy,
		UserName: "",
	}
	if replenishment.CreatedByName != nil {
		response.CreatedBy.UserName = *replenishment.CreatedByName
	}
	response.DistributorID = replenishment.DistributorID
	if replenishment.DistributorCode != nil {
		response.DistributorCode = *replenishment.DistributorCode
	}
	if replenishment.DistributorName != nil {
		response.DistributorName = *replenishment.DistributorName
	}
	if replenishment.Address != nil {
		response.DistributorAddr = *replenishment.Address
	}
	if !replenishment.CreatedAt.IsZero() {
		response.CreatedBy.CreatedAt = replenishment.CreatedAt.In(constant.AsiaJakartaLocation).Format(time.RFC3339)
	}

	// Updated by
	if replenishment.UpdatedBy != nil {
		updatedBy := entity.ReplenishmentUserInfo{
			UserID:   *replenishment.UpdatedBy,
			UserName: "",
		}
		if replenishment.UpdatedByName != nil {
			updatedBy.UserName = *replenishment.UpdatedByName
		}
		if replenishment.UpdatedAt != nil {
			updatedBy.UpdatedAt = replenishment.UpdatedAt.In(constant.AsiaJakartaLocation).Format(time.RFC3339)
		}
		response.UpdatedBy = &updatedBy
	}

	return response
}

func (service *replenishmentServiceImpl) int64ToString(val int64) string {
	return fmt.Sprintf("%d", val)
}

func (service *replenishmentServiceImpl) Store(request entity.CreateReplenishmentOrderBody) error {
	ctx := context.Background()

	if len(request.Data) == 0 {
		return errors.New("data is required, minimum 1 product")
	}

	date, err := service.validateAndParseDate(request.CustID)
	if err != nil {
		return err
	}

	err = service.ReplenishmentRepository.FindSupplierByID(ctx, request.SupID, request.CustID)
	if err != nil {
		return fmt.Errorf("sup_id: %d, %s", request.SupID, err.Error())
	}

	distributorID, err := service.ReplenishmentRepository.GetDistributorIDByCustID(ctx, request.CustID)
	if err != nil {
		return fmt.Errorf("failed to get distributor_id: %w", err)
	}
	if request.DistributorID != nil && *request.DistributorID > 0 {
		distributorID = *request.DistributorID
	}

	err = service.ReplenishmentRepository.FindWarehouseByID(ctx, request.WhID, distributorID)
	if err != nil {
		return fmt.Errorf("wh_id: %d, %s", request.WhID, err.Error())
	}

	// Validate all products exist
	for _, product := range request.Data {
		err = service.ReplenishmentRepository.FindProductByID(ctx, product.ProID, request.CustID, request.ParentCustID)
		if err != nil {
			return fmt.Errorf("data[].pro_id: %d, %s", product.ProID, err.Error())
		}
	}

	// Parse optional dates (DD/MM/YYYY format)
	var soStartDate, soEndDate, deliveryDate *time.Time
	if request.SoStartDate != nil && *request.SoStartDate != "" {
		parsedTime, err := str.DateStrDdMmYyyyToTime(*request.SoStartDate)
		if err != nil {
			return fmt.Errorf("so_start_date: invalid date format, expected DD/MM/YYYY: %v", err)
		}
		soStartDate = &parsedTime
	}

	if request.SoEndDate != nil && *request.SoEndDate != "" {
		parsedTime, err := str.DateStrDdMmYyyyToTime(*request.SoEndDate)
		if err != nil {
			return fmt.Errorf("so_end_date: invalid date format, expected DD/MM/YYYY: %v", err)
		}
		soEndDate = &parsedTime
	}

	if request.DeliveryDate != nil && *request.DeliveryDate != "" {
		parsedTime, err := str.DateStrDdMmYyyyToTime(*request.DeliveryDate)
		if err != nil {
			return fmt.Errorf("delivery_date: invalid date format, expected DD/MM/YYYY: %v", err)
		}
		deliveryDate = &parsedTime
	}

	var distributorIDToStore *int64
	if request.DistributorID != nil && *request.DistributorID > 0 {
		distributorIDToStore = request.DistributorID
	} else if distributorID > 0 {
		distributorIDToStore = &distributorID
	}

	status := StatusNeedReview
	defaultFinalOrderToQtyTotal := false
	applyDetailApprovalOnCreate := false
	var isApproval *bool
	var approveBy *int64
	var approveAt *time.Time

	isPrincipal, err := service.ReplenishmentRepository.CheckIsPrincipal(ctx, request.CustID)
	if err != nil {
		return fmt.Errorf("failed to check principal status: %w", err)
	}
	isGlobalPIC, err := service.ReplenishmentRepository.IsUserApprovalPIC(ctx, request.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to check PIC role: %w", err)
	}
	if !isGlobalPIC && request.CreatedEmpID > 0 {
		isGlobalPIC, err = service.ReplenishmentRepository.IsUserApprovalPIC(ctx, request.CreatedEmpID)
		if err != nil {
			return fmt.Errorf("failed to check PIC role by employee_id: %w", err)
		}
	}

	creatingForChildDistributor := request.DistributorID != nil && *request.DistributorID > 0
	principalPICAutoApprove := (isPrincipal || isGlobalPIC) && creatingForChildDistributor

	if principalPICAutoApprove {
		status = StatusApproved
		applyDetailApprovalOnCreate = true
		approval := true
		now := time.Now()
		isApproval = &approval
		approveBy = &request.CreatedBy
		approveAt = &now
	} else if request.IsPic != nil && *request.IsPic {
		isPIC, err := service.ReplenishmentRepository.IsDistributorApprovalPIC(ctx, request.CreatedBy, request.SupID, distributorID)
		if err != nil {
			return fmt.Errorf("failed to check PIC setup: %w", err)
		}
		if !isPIC && request.CreatedEmpID > 0 {
			isPIC, err = service.ReplenishmentRepository.IsDistributorApprovalPIC(ctx, request.CreatedEmpID, request.SupID, distributorID)
			if err != nil {
				return fmt.Errorf("failed to check PIC setup by employee_id: %w", err)
			}
		}
		if isPIC {
			status = StatusApproved
			applyDetailApprovalOnCreate = true
			approval := true
			now := time.Now()
			isApproval = &approval
			approveBy = &request.CreatedBy
			approveAt = &now
		}
	} else {
		hasSetup, isApprovalRequired, err := service.ReplenishmentRepository.GetDistributorApprovalRequirement(ctx, request.SupID, distributorID)
		if err != nil {
			return fmt.Errorf("failed to check approval requirement: %w", err)
		}
		if !hasSetup {
			status = StatusOnDelivery
			defaultFinalOrderToQtyTotal = true
		} else if !isApprovalRequired {
			status = StatusApproved
		} else {
			status = StatusNeedReview
		}
	}

	deliveryType := strings.TrimSpace(request.DeliveryType)
	if deliveryType == "" {
		deliveryType = "Full"
	}

	replenishmentModel := model.ReplenishmentOrder{
		CustID:            request.CustID,
		ParentCustID:      request.ParentCustID,
		Date:              date,
		DistributorID:     distributorIDToStore,
		SupID:             request.SupID,
		WhID:              request.WhID,
		DeliveryType:      deliveryType,
		ReplenishmentType: request.ReplenishmentType,
		SoStartDate:       soStartDate,
		SoEndDate:         soEndDate,
		DeliveryDate:      deliveryDate,
		Note:              request.Note,
		Status:            status,
		IsApproval:        isApproval,
		ApproveBy:         approveBy,
		ApproveAt:         approveAt,
		IsAdditionFrom:    false,
		CreatedBy:         request.CreatedBy,
		UpdatedBy:         &request.CreatedBy,
	}

	// Store within transaction
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Store header (replenishment_no auto-generated; replenishment_id auto-increment)
		err := service.ReplenishmentRepository.Store(txCtx, &replenishmentModel)
		if err != nil {
			return err
		}

		// Store details
		for _, product := range request.Data {
			// Dereference pointers (validation already ensures they are not nil)
			orderBookingQty1 := float64(0)
			if product.OrderBookingQty1 != nil {
				orderBookingQty1 = *product.OrderBookingQty1
			}
			orderBookingQty2 := float64(0)
			if product.OrderBookingQty2 != nil {
				orderBookingQty2 = *product.OrderBookingQty2
			}
			orderBookingQty3 := float64(0)
			if product.OrderBookingQty3 != nil {
				orderBookingQty3 = *product.OrderBookingQty3
			}
			purchPrice1 := float64(0)
			if product.PurchPrice1 != nil {
				purchPrice1 = *product.PurchPrice1
			}
			purchPrice2 := float64(0)
			if product.PurchPrice2 != nil {
				purchPrice2 = *product.PurchPrice2
			}
			purchPrice3 := float64(0)
			if product.PurchPrice3 != nil {
				purchPrice3 = *product.PurchPrice3
			}
			estimatedPrice := float64(0)
			if product.EstimatedPrice != nil {
				estimatedPrice = *product.EstimatedPrice
			}

			detailModel := model.ReplenishmentOrderDetail{
				CustID:           request.CustID,
				ReplenishmentID:  replenishmentModel.ReplenishmentID,
				ProID:            product.ProID,
				OrderBookingQty1: orderBookingQty1,
				OrderBookingQty2: orderBookingQty2,
				OrderBookingQty3: orderBookingQty3,
				PurchPrice1:      purchPrice1,
				PurchPrice2:      purchPrice2,
				PurchPrice3:      purchPrice3,
				EstimatedPrice:   estimatedPrice,
				CreatedBy:        request.CreatedBy,
				UpdatedBy:        &request.CreatedBy,
			}
			if defaultFinalOrderToQtyTotal || applyDetailApprovalOnCreate {
				qtyApproval1 := orderBookingQty1
				qtyApproval2 := orderBookingQty2
				qtyApproval3 := orderBookingQty3
				allocation1 := float64(0)
				allocation2 := float64(0)
				allocation3 := float64(0)
				detailModel.QtyOrderApproval1 = &qtyApproval1
				detailModel.QtyOrderApproval2 = &qtyApproval2
				detailModel.QtyOrderApproval3 = &qtyApproval3
				detailModel.QtyOrderAllocation1 = &allocation1
				detailModel.QtyOrderAllocation2 = &allocation2
				detailModel.QtyOrderAllocation3 = &allocation3
			}

			err = service.ReplenishmentRepository.CreateDetail(txCtx, &detailModel)
			if err != nil {
				return err
			}
		}

		if status == StatusNeedReview {
			level, sequence, pic, err := service.ReplenishmentRepository.GetInitialDistributorApproval(txCtx, request.SupID, distributorID)
			if err != nil {
				return fmt.Errorf("failed to resolve initial approval PIC: %w", err)
			}
			err = service.ReplenishmentRepository.InsertReplenishmentOrderApproval(txCtx, request.CustID, replenishmentModel.ReplenishmentID, level, sequence, pic)
			if err != nil {
				return fmt.Errorf("failed to insert replenishment approval queue: %w", err)
			}
		}

		return nil
	})

	return err
}

// validateAndParseDate validates and parses date (use current date if not provided)
func (service *replenishmentServiceImpl) validateAndParseDate(custID string) (time.Time, error) {
	now := time.Now().In(constant.AsiaJakartaLocation)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, constant.AsiaJakartaLocation), nil
}

func (service *replenishmentServiceImpl) Detail(replenishmentID, typeParam, statusParam string, custID, parentCustID string, isPrincipal bool) (*entity.ReplenishmentDetailResponse, error) {
	ctx := context.Background()

	// Get header by replenishment_id
	replenishment, err := service.ReplenishmentRepository.FindByReplenishmentID(ctx, replenishmentID, custID, parentCustID, isPrincipal)
	if err != nil {
		return nil, err
	}

	// Validate type and status logic based on requirement
	// IF type == "final" AND (status == "On Delivery" OR status == 3) THEN RETURN EMPTY RESULT
	if typeParam == TypeFinal {
		statusInt := replenishment.Status
		var statusName string
		if replenishment.StatusName != nil {
			statusName = *replenishment.StatusName
		}
		if strings.EqualFold(strings.TrimSpace(statusParam), StatusNameOnDelivery) ||
			strings.EqualFold(strings.TrimSpace(statusParam), StatusNameGoodsIssue) ||
			strings.EqualFold(strings.TrimSpace(statusName), StatusNameOnDelivery) ||
			strings.EqualFold(strings.TrimSpace(statusName), StatusNameGoodsIssue) ||
			statusInt == StatusRejected {
			return nil, nil
		}
	}

	// IF type == "final" AND status == 7 THEN show final data
	// IF type == "replenishment" AND status == "submitted" THEN show replenishment data
	// Note: If type is empty, show default behavior (replenishment type)
	if typeParam != "" {
		if typeParam == TypeFinal {
			if replenishment.Status != StatusCompleted {
				return nil, fmt.Errorf("invalid type or status: type 'final' requires status %d (Completed)", StatusCompleted)
			}
		} else if typeParam == TypeReplenishment {
			// Type replenishment: status == "submitted" means status code 1 (Need Review)
			if statusParam == StatusNameSubmitted {
				if replenishment.Status != StatusNeedReview {
					return nil, fmt.Errorf("invalid type or status: type 'replenishment' with status 'submitted' requires status code %d (Need Review)", StatusNeedReview)
				}
			}
		}
	}

	var response entity.ReplenishmentDetailResponse

	// Map header to response
	response.ReplenishmentID = service.int64ToString(replenishment.ReplenishmentID)
	response.ReplenishmentNo = replenishment.ReplenishmentNo

	// Format date
	if !replenishment.Date.IsZero() {
		dateInJkt := replenishment.Date.In(constant.AsiaJakartaLocation)
		response.Date = dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
	}

	// Format optional dates
	if replenishment.SoStartDate != nil {
		dateInJkt := replenishment.SoStartDate.In(constant.AsiaJakartaLocation)
		formatted := dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
		response.SoStartDate = &formatted
	}
	if replenishment.SoEndDate != nil {
		dateInJkt := replenishment.SoEndDate.In(constant.AsiaJakartaLocation)
		formatted := dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
		response.SoEndDate = &formatted
	}
	if replenishment.DeliveryDate != nil {
		dateInJkt := replenishment.DeliveryDate.In(constant.AsiaJakartaLocation)
		formatted := dateInJkt.Format(constant.DATE_FORMAT_DETAIL)
		response.DeliveryDate = &formatted
	}

	response.ReplenishmentNo = replenishment.ReplenishmentNo
	response.ReplenishmentType = replenishment.ReplenishmentType
	response.DeliveryType = replenishment.DeliveryType
	response.SupID = service.int64ToString(replenishment.SupID)
	if replenishment.SupCode != nil {
		response.SupCode = *replenishment.SupCode
	}
	if replenishment.SupName != nil {
		response.SupName = *replenishment.SupName
	}
	response.WhID = service.int64ToString(replenishment.WhID)
	if replenishment.WhCode != nil {
		response.WhCode = *replenishment.WhCode
	}
	if replenishment.WhName != nil {
		response.WhName = *replenishment.WhName
	}
	if replenishment.StatusName != nil {
		response.Status = service.mapStatusDisplayName(*replenishment.StatusName)
	} else {
		response.Status = service.int64ToString(int64(replenishment.Status))
	}

	// Delivery fee from GR
	if replenishment.DeliveryFee != nil {
		response.DeliveryFee = int64(*replenishment.DeliveryFee)
	}
	response.Notes = replenishment.Note

	// Distributor fields
	response.DistributorID = replenishment.DistributorID
	response.DistributorCode = replenishment.DistributorCode
	response.DistributorName = replenishment.DistributorName
	response.Address = replenishment.Address

	// Get details with additional fields
	details, err := service.ReplenishmentRepository.FindDetailByReplenishmentID(ctx, replenishmentID, custID, parentCustID, isPrincipal)
	if err != nil {
		return nil, err
	}

	// Map details to response based on type
	var productData []entity.ReplenishmentDetailProductResponse
	for _, detail := range details {
		product := entity.ReplenishmentDetailProductResponse{
			ReplenishmentDetID:  service.int64ToString(detail.ReplenishmentDetailID),
			ProductID:           detail.ProID,
			PurchPrice1:         detail.PurchPrice1,
			PurchPrice2:         detail.PurchPrice2,
			PurchPrice3:         detail.PurchPrice3,
			EstimatedPrice:      detail.EstimatedPrice,
			Vat:                 0,
			SafStockQty:         0,
			MinStockQty:         0,
			Qty1:                0,
			Qty2:                0,
			Qty3:                0,
			QtyOrderAllocation1: 0,
			QtyOrderAllocation2: 0,
			QtyOrderAllocation3: 0,
			QtyOrderApproval1:   0,
			QtyOrderApproval2:   0,
			QtyOrderApproval3:   0,
		}

		// Set product info
		if detail.ProCode != nil {
			product.ProductCode = *detail.ProCode
		}
		if detail.ProName != nil {
			product.ProductName = *detail.ProName
		}
		if detail.UnitID1 != nil {
			product.UnitID1 = *detail.UnitID1
		}
		if detail.UnitID2 != nil {
			product.UnitID2 = *detail.UnitID2
		}
		if detail.UnitID3 != nil {
			product.UnitID3 = *detail.UnitID3
		}
		if detail.Vat != nil {
			product.Vat = *detail.Vat
		}

		if typeParam == "final" && replenishment.Status == 7 {
			// Type: final, Status: 7 (Completed)
			// Show: qty_final, stock_received, sub_total
			if detail.QtyFinal1 != nil {
				product.QtyFinal1 = detail.QtyFinal1
			}
			if detail.QtyFinal2 != nil {
				product.QtyFinal2 = detail.QtyFinal2
			}
			if detail.QtyFinal3 != nil {
				product.QtyFinal3 = detail.QtyFinal3
			}
			if detail.StockReceived1 != nil {
				product.StockReceived1 = detail.StockReceived1
			}
			if detail.StockReceived2 != nil {
				product.StockReceived2 = detail.StockReceived2
			}
			if detail.StockReceived3 != nil {
				product.StockReceived3 = detail.StockReceived3
			}

			// Calculate sub_total from qty_final
			qty1 := float64(0)
			qty2 := float64(0)
			qty3 := float64(0)
			if detail.QtyFinal1 != nil {
				qty1 = *detail.QtyFinal1
			}
			if detail.QtyFinal2 != nil {
				qty2 = *detail.QtyFinal2
			}
			if detail.QtyFinal3 != nil {
				qty3 = *detail.QtyFinal3
			}
			product.SubTotal = (qty1 * product.PurchPrice1) +
				(qty2 * product.PurchPrice2) +
				(qty3 * product.PurchPrice3)

		} else {
			// Type: replenishment or default
			// Show: qty_ro, qty_order_allocation, qty_total, sub_total
			product.QtyRo1 = detail.OrderBookingQty1
			product.QtyRo2 = detail.OrderBookingQty2
			product.QtyRo3 = detail.OrderBookingQty3

			var qtyOrderAlloc1, qtyOrderAlloc2, qtyOrderAlloc3 float64
			if detail.QtyOrderAllocation1 != nil {
				qtyOrderAlloc1 = *detail.QtyOrderAllocation1
				product.QtyOrderAllocation1 = *detail.QtyOrderAllocation1
			} else {
				product.QtyOrderAllocation1 = 0
			}
			if detail.QtyOrderAllocation2 != nil {
				qtyOrderAlloc2 = *detail.QtyOrderAllocation2
				product.QtyOrderAllocation2 = *detail.QtyOrderAllocation2
			} else {
				product.QtyOrderAllocation2 = 0
			}
			if detail.QtyOrderAllocation3 != nil {
				qtyOrderAlloc3 = *detail.QtyOrderAllocation3
				product.QtyOrderAllocation3 = *detail.QtyOrderAllocation3
			} else {
				product.QtyOrderAllocation3 = 0
			}

			// Get qty_order_approval from detail
			if detail.QtyOrderApproval1 != nil {
				product.QtyOrderApproval1 = *detail.QtyOrderApproval1
			} else {
				product.QtyOrderApproval1 = 0
			}
			if detail.QtyOrderApproval2 != nil {
				product.QtyOrderApproval2 = *detail.QtyOrderApproval2
			} else {
				product.QtyOrderApproval2 = 0
			}
			if detail.QtyOrderApproval3 != nil {
				product.QtyOrderApproval3 = *detail.QtyOrderApproval3
			} else {
				product.QtyOrderApproval3 = 0
			}

			// Calculate qty_total based on status
			// Status 1 (Need Review): qty_total = qty_ro
			// Other statuses: qty_total = qty_ro + qty_order_allocation (follow minus/plus value)
			if replenishment.Status == 1 {
				product.QtyTotal1 = detail.OrderBookingQty1
				product.QtyTotal2 = detail.OrderBookingQty2
				product.QtyTotal3 = detail.OrderBookingQty3
			} else {
				// Other statuses: Total Order = Replenishment order + Order Allocation (follow minus/plus value)
				product.QtyTotal1 = detail.OrderBookingQty1 + qtyOrderAlloc1
				product.QtyTotal2 = detail.OrderBookingQty2 + qtyOrderAlloc2
				product.QtyTotal3 = detail.OrderBookingQty3 + qtyOrderAlloc3
			}

			// Calculate sub_total according to requirement:
			// sub_total = (qty_order_allocation1 * purch_price1) + (qty_order_allocation2 * purch_price2) + (qty_order_allocation3 * purch_price3)
			// If qty_order_allocation is nil or 0, use qty_order_approval as fallback
			var subTotalQty1, subTotalQty2, subTotalQty3 float64
			if qtyOrderAlloc1 != 0 {
				subTotalQty1 = qtyOrderAlloc1
			} else if detail.QtyOrderApproval1 != nil && *detail.QtyOrderApproval1 != 0 {
				subTotalQty1 = *detail.QtyOrderApproval1
			}
			if qtyOrderAlloc2 != 0 {
				subTotalQty2 = qtyOrderAlloc2
			} else if detail.QtyOrderApproval2 != nil && *detail.QtyOrderApproval2 != 0 {
				subTotalQty2 = *detail.QtyOrderApproval2
			}
			if qtyOrderAlloc3 != 0 {
				subTotalQty3 = qtyOrderAlloc3
			} else if detail.QtyOrderApproval3 != nil && *detail.QtyOrderApproval3 != 0 {
				subTotalQty3 = *detail.QtyOrderApproval3
			}
			product.SubTotal = (subTotalQty1 * product.PurchPrice1) +
				(subTotalQty2 * product.PurchPrice2) +
				(subTotalQty3 * product.PurchPrice3)

			// Calculate estimated_price from qty_order_approval to ensure consistency with minus/plus values
			// estimated_price = (qty_order_approval1 * purch_price1) + (qty_order_approval2 * purch_price2) + (qty_order_approval3 * purch_price3)
			var estQty1, estQty2, estQty3 float64
			if detail.QtyOrderApproval1 != nil {
				estQty1 = *detail.QtyOrderApproval1
			}
			if detail.QtyOrderApproval2 != nil {
				estQty2 = *detail.QtyOrderApproval2
			}
			if detail.QtyOrderApproval3 != nil {
				estQty3 = *detail.QtyOrderApproval3
			}
			product.EstimatedPrice = (estQty1 * product.PurchPrice1) +
				(estQty2 * product.PurchPrice2) +
				(estQty3 * product.PurchPrice3)
		}

		// In transit stock (sum of order_booking_qty from all replenishment orders with status On Delivery)
		product.InTransitStock1 = detail.InTransitStock1
		product.InTransitStock2 = detail.InTransitStock2
		product.InTransitStock3 = detail.InTransitStock3

		// Saf stock qty and min stock qty
		if detail.SafStockQty != nil {
			product.SafStockQty = *detail.SafStockQty
		}
		if detail.MinStockQty != nil {
			product.MinStockQty = *detail.MinStockQty
		}

		// Warehouse stock qty1/2/3
		if detail.Qty1 != nil {
			product.Qty1 = *detail.Qty1
		}
		if detail.Qty2 != nil {
			product.Qty2 = *detail.Qty2
		}
		if detail.Qty3 != nil {
			product.Qty3 = *detail.Qty3
		}

		productData = append(productData, product)
	}

	response.ProductData = productData

	// Create map of product_data by product_id for quick lookup
	productDataMap := make(map[int64]entity.ReplenishmentDetailProductResponse)
	for _, product := range productData {
		productDataMap[product.ProductID] = product
	}

	if replenishment.Status >= StatusOnDelivery {
		finalData, err := service.ReplenishmentRepository.FindFinalByReplenishmentID(ctx, replenishment.ReplenishmentNo, custID, parentCustID, isPrincipal)
		if err != nil {
			response.FinalReplanishment = []entity.ReplenishmentFinalResponse{}
		} else {
			var finalArray []entity.ReplenishmentFinalResponse
			for _, final := range finalData {
				// Get purch_price from product_data based on product_id
				var purchPrice1, purchPrice2, purchPrice3 float64
				if product, exists := productDataMap[final.ProID]; exists {
					purchPrice1 = product.PurchPrice1
					purchPrice2 = product.PurchPrice2
					purchPrice3 = product.PurchPrice3
				}

				finalItem := entity.ReplenishmentFinalResponse{
					ReplenishmentDetID:  service.int64ToString(final.ReplenishmentDetailID),
					ProductID:           final.ProID,
					PurchPriceDelivery1: &purchPrice1,
					PurchPriceDelivery2: &purchPrice2,
					PurchPriceDelivery3: &purchPrice3,
					FinalOrder1:         final.FinalOrder1,
					FinalOrder2:         final.FinalOrder2,
					FinalOrder3:         final.FinalOrder3,
					GrPrice1:            final.GrPrice1,
					GrPrice2:            final.GrPrice2,
					GrPrice3:            final.GrPrice3,
					StockReceived1:      final.StockReceived1,
					StockReceived2:      final.StockReceived2,
					StockReceived3:      final.StockReceived3,
				}

				if final.ProCode != nil {
					finalItem.ProductCode = *final.ProCode
				}
				if final.ProName != nil {
					finalItem.ProductName = *final.ProName
				}
				if final.UnitID1 != nil {
					finalItem.UnitID1 = final.UnitID1
				}
				if final.UnitID2 != nil {
					finalItem.UnitID2 = final.UnitID2
				}
				if final.UnitID3 != nil {
					finalItem.UnitID3 = final.UnitID3
				}

				// Vat from product
				if final.Vat != nil {
					finalItem.Vat = final.Vat
				}

				// purch_price_delivery already set from rod.purch_price in query (same as product_data.purch_price)
				// No need to override with gr_price

				// Calculate sub_total: SUM(gr_price * stock_receipt)
				var grPrice1, grPrice2, grPrice3 float64
				var stock1, stock2, stock3 float64

				if final.GrPrice1 != nil {
					grPrice1 = *final.GrPrice1
				}
				if final.GrPrice2 != nil {
					grPrice2 = *final.GrPrice2
				}
				if final.GrPrice3 != nil {
					grPrice3 = *final.GrPrice3
				}
				if final.StockReceived1 != nil {
					stock1 = *final.StockReceived1
				}
				if final.StockReceived2 != nil {
					stock2 = *final.StockReceived2
				}
				if final.StockReceived3 != nil {
					stock3 = *final.StockReceived3
				}

				finalItem.SubTotal = (grPrice1 * stock1) + (grPrice2 * stock2) + (grPrice3 * stock3)
				finalArray = append(finalArray, finalItem)
			}
			response.FinalReplanishment = finalArray
		}
	} else {
		response.FinalReplanishment = []entity.ReplenishmentFinalResponse{}
	}

	// Get good receipt data
	goodReceiptData, err := service.ReplenishmentRepository.FindGoodReceiptByReplenishmentNo(ctx, replenishment.ReplenishmentNo, custID, parentCustID, isPrincipal)
	if err != nil {
		response.GoodReceipt = []entity.ReplenishmentGoodReceiptResponse{}
		return &response, nil
	}

	goodReceiptArray := make([]entity.ReplenishmentGoodReceiptResponse, 0)
	for _, gr := range goodReceiptData {
		grItem := entity.ReplenishmentGoodReceiptResponse{
			PoNo:           gr.PoNo,
			ProID:          gr.ProID,
			UnitPrice1:     gr.UnitPrice1,
			UnitPrice2:     gr.UnitPrice2,
			UnitPrice3:     gr.UnitPrice3,
			QtyReceived1:   gr.QtyReceived1,
			QtyReceived2:   gr.QtyReceived2,
			QtyReceived3:   gr.QtyReceived3,
			EstimatedPrice: gr.EstimatedPrice,
			Vat:            0,
		}

		if gr.ProCode != nil {
			grItem.ProCode = *gr.ProCode
		}
		if gr.ProName != nil {
			grItem.ProName = *gr.ProName
		}
		if gr.Vat != nil {
			grItem.Vat = *gr.Vat
		}

		goodReceiptArray = append(goodReceiptArray, grItem)
	}
	response.GoodReceipt = goodReceiptArray

	return &response, nil
}

func (service *replenishmentServiceImpl) ProductList(dataFilter entity.ReplenishmentProductQueryFilter, custId, parentCustId string) ([]entity.ReplenishmentProductListResponse, int64, int, error) {
	ctx := context.Background()

	products, total, lastPage, err := service.ReplenishmentRepository.FindProductList(ctx, dataFilter, custId, parentCustId)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []entity.ReplenishmentProductListResponse
	for _, product := range products {
		response := service.mapToProductResponse(&product)
		responses = append(responses, response)
	}

	return responses, total, lastPage, nil
}

func (service *replenishmentServiceImpl) mapToProductResponse(product *model.ReplenishmentProductList) entity.ReplenishmentProductListResponse {
	var response entity.ReplenishmentProductListResponse

	response.ProID = product.ProID
	response.ProCode = product.ProCode
	response.ProName = product.ProName
	response.PurchPrice1 = product.PurchPrice1
	response.PurchPrice2 = product.PurchPrice2
	response.PurchPrice3 = product.PurchPrice3
	response.UnitID1 = product.UnitID1
	response.UnitID2 = product.UnitID2
	response.UnitID3 = product.UnitID3
	response.Vat = product.Vat
	response.Qty1 = product.Qty1
	response.Qty2 = product.Qty2
	response.Qty3 = product.Qty3
	response.InTransitStock1 = product.InTransitStock1
	response.InTransitStock2 = product.InTransitStock2
	response.InTransitStock3 = product.InTransitStock3

	// Calculate estimated_price: (qty1 * purch_price1) + (qty2 * purch_price2) + (qty3 * purch_price3)
	response.EstimatedPrice = (product.Qty1 * product.PurchPrice1) + (product.Qty2 * product.PurchPrice2) + (product.Qty3 * product.PurchPrice3)

	return response
}

func (service *replenishmentServiceImpl) UpdateApproval(replenishmentID int64, request entity.UpdateReplenishmentApprovalRequest, custID string, userID int64, employeeID int64) error {
	ctx := context.Background()
	replenishment, err := service.ReplenishmentRepository.GetReplenishmentOrderByID(ctx, replenishmentID, "")
	if err != nil {
		return fmt.Errorf("replenishment order not found: %w", err)
	}
	return service.updateApprovalForOrder(ctx, replenishment, request, custID, userID, employeeID)
}

func (service *replenishmentServiceImpl) updateApprovalForOrder(ctx context.Context, replenishment *model.ReplenishmentOrder, request entity.UpdateReplenishmentApprovalRequest, custID string, userID int64, employeeID int64) error {
	var err error
	distributorID := int64(0)
	if replenishment.DistributorID != nil && *replenishment.DistributorID > 0 {
		distributorID = *replenishment.DistributorID
	}
	if distributorID == 0 {
		distributorID, err = service.ReplenishmentRepository.GetDistributorIDByCustID(ctx, replenishment.CustID)
		if err != nil {
			return fmt.Errorf("failed to get distributor_id: %w", err)
		}
	}

	authorizedPIC := userID
	isPIC, err := service.ReplenishmentRepository.IsDistributorApprovalPIC(ctx, userID, replenishment.SupID, distributorID)
	if err != nil {
		return fmt.Errorf("failed to check replenishment PIC setup: %w", err)
	}
	if !isPIC && employeeID > 0 {
		isPIC, err = service.ReplenishmentRepository.IsDistributorApprovalPIC(ctx, employeeID, replenishment.SupID, distributorID)
		if err != nil {
			return fmt.Errorf("failed to check replenishment PIC setup by employee_id: %w", err)
		}
		if isPIC {
			authorizedPIC = employeeID
		}
	}
	if !isPIC {
		return fmt.Errorf("only setup replenishment PIC is allowed to approve replenishment data")
	}

	now := time.Now()

	// Validate approval field
	if request.Approval == nil {
		return fmt.Errorf("approval field is required")
	}

	// Update header based on approval status
	if *request.Approval {
		// Approve: status = 2 (Approved)
		approval := true
		replenishment.Status = StatusApproved
		replenishment.IsApproval = &approval
		replenishment.ApproveBy = &userID
		replenishment.ApproveAt = &now
		replenishment.UpdatedBy = &userID
	} else {
		// Reject: status = 3 (Rejected)
		approval := false
		replenishment.Status = StatusRejected
		replenishment.IsApproval = &approval
		replenishment.ApproveBy = &userID
		replenishment.ApproveAt = &now
		replenishment.UpdatedBy = &userID
	}

	// Process in transaction
	err = service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Update header first
		if err := service.ReplenishmentRepository.UpdateApproval(txCtx, replenishment); err != nil {
			return err
		}
		approvalQueueStatus := StatusApproved
		var approvalRemarks *string
		if !*request.Approval {
			approvalQueueStatus = StatusRejected
			approvalRemarks = request.Remarks
		}
		if err := service.ReplenishmentRepository.UpdateReplenishmentOrderApprovalStatus(txCtx, replenishment.CustID, replenishment.ReplenishmentID, authorizedPIC, approvalQueueStatus, approvalRemarks); err != nil {
			return err
		}

		// If approval is true and has data changes, process detail changes
		if *request.Approval && (len(request.Data) > 0 || len(request.DeletedData) > 0) {
			// Validate all items before processing
			for _, item := range request.Data {
				if item.ReplenishmentDetID != nil {
					// EDIT produk: replenishment_det_id must be NOT NULL
					if item.OrderBookingQty1 == nil || item.OrderBookingQty2 == nil || item.OrderBookingQty3 == nil ||
						item.PurchPrice1 == nil || item.PurchPrice2 == nil || item.PurchPrice3 == nil {
						return fmt.Errorf("when replenishment_det_id is provided (edit product), all fields (order_booking_qty1/2/3, purch_price1/2/3) are required")
					}
				} else if item.ProID != nil {
					// ADD new produk: replenishment_det_id must be NULL
					if item.OrderBookingQty1 == nil || item.OrderBookingQty2 == nil || item.OrderBookingQty3 == nil ||
						item.PurchPrice1 == nil || item.PurchPrice2 == nil || item.PurchPrice3 == nil {
						return fmt.Errorf("when replenishment_det_id is NULL (new product), pro_id and all fields (order_booking_qty1/2/3, purch_price1/2/3) are required")
					}
				} else {
					return fmt.Errorf("either replenishment_det_id (for edit) or pro_id (for add) must be provided")
				}
				if len(item.ReturnReasonID) == 0 || item.ReturnReasonID[0] <= 0 {
					return fmt.Errorf("return_reason_id is required for each changed product")
				}
				returnReasonID := item.ReturnReasonID[0]
				exists, err := service.ReplenishmentRepository.IsReturnReasonDistributorExists(txCtx, returnReasonID)
				if err != nil {
					return fmt.Errorf("failed to validate return_reason_id %d: %w", returnReasonID, err)
				}
				if !exists {
					return fmt.Errorf("return_reason_id %d not found", returnReasonID)
				}
			}

			// Get existing details - MUST use txCtx to be within transaction
			// Since only principal can reach here (checked above), don't filter by cust_id for query
			existingDetails, err := service.ReplenishmentRepository.FindDetailByReplenishmentIDForUpdate(txCtx, replenishment.ReplenishmentID, "")
			if err != nil {
				return fmt.Errorf("failed to get existing details: %w", err)
			}

			// Create map of existing details by replenishment_detail_id
			existingDetailMap := make(map[int64]*model.ReplenishmentOrderDetail)
			existingDetailMapByProID := make(map[int64]*model.ReplenishmentOrderDetail)
			for i := range existingDetails {
				existingDetailMap[existingDetails[i].ReplenishmentDetailID] = &existingDetails[i]
				existingDetailMapByProID[existingDetails[i].ProID] = &existingDetails[i]
			}

			// Process each item in request data
			for _, item := range request.Data {
				if item.ReplenishmentDetID != nil {
					// UPDATE existing detail by replenishment_det_id
					existingDetail, exists := existingDetailMap[*item.ReplenishmentDetID]
					if !exists {
						return fmt.Errorf("replenishment_detail_id %d not found", *item.ReplenishmentDetID)
					}

					// Update approval quantities (required fields when editing)
					// Use order_booking_qty from request as qty_order_approval
					existingDetail.QtyOrderApproval1 = item.OrderBookingQty1
					existingDetail.QtyOrderApproval2 = item.OrderBookingQty2
					existingDetail.QtyOrderApproval3 = item.OrderBookingQty3

					// Calculate allocation (approval - booking)
					allocation1 := *item.OrderBookingQty1 - existingDetail.OrderBookingQty1
					allocation2 := *item.OrderBookingQty2 - existingDetail.OrderBookingQty2
					allocation3 := *item.OrderBookingQty3 - existingDetail.OrderBookingQty3
					existingDetail.QtyOrderAllocation1 = &allocation1
					existingDetail.QtyOrderAllocation2 = &allocation2
					existingDetail.QtyOrderAllocation3 = &allocation3

					// Update prices
					existingDetail.PurchPrice1 = *item.PurchPrice1
					existingDetail.PurchPrice2 = *item.PurchPrice2
					existingDetail.PurchPrice3 = *item.PurchPrice3

					// Calculate estimated_price automatically: (order_booking_qty1 * purch_price1) + (order_booking_qty2 * purch_price2) + (order_booking_qty3 * purch_price3)
					existingDetail.EstimatedPrice = (*item.OrderBookingQty1 * *item.PurchPrice1) +
						(*item.OrderBookingQty2 * *item.PurchPrice2) +
						(*item.OrderBookingQty3 * *item.PurchPrice3)

					existingDetail.UpdatedBy = &userID
					returnReasonID := item.ReturnReasonID[0]
					existingDetail.ReturnReasonID = &returnReasonID

					if err := service.ReplenishmentRepository.UpdateDetail(txCtx, existingDetail); err != nil {
						return err
					}
				} else if item.ProID != nil {
					// INSERT new detail (add product) - replenishment_det_id is NULL means new product
					// Calculate estimated_price automatically: (order_booking_qty1 * purch_price1) + (order_booking_qty2 * purch_price2) + (order_booking_qty3 * purch_price3)
					estimatedPrice := (*item.OrderBookingQty1 * *item.PurchPrice1) +
						(*item.OrderBookingQty2 * *item.PurchPrice2) +
						(*item.OrderBookingQty3 * *item.PurchPrice3)

					returnReasonID := item.ReturnReasonID[0]
					newDetail := &model.ReplenishmentOrderDetail{
						CustID:              replenishment.CustID, // Use replenishment.CustID instead of custID parameter to match the replenishment order's cust_id
						ReplenishmentID:     replenishment.ReplenishmentID,
						ProID:               *item.ProID,
						OrderBookingQty1:    0, // New product, booking qty is 0
						OrderBookingQty2:    0,
						OrderBookingQty3:    0,
						QtyOrderApproval1:   item.OrderBookingQty1,
						QtyOrderApproval2:   item.OrderBookingQty2,
						QtyOrderApproval3:   item.OrderBookingQty3,
						QtyOrderAllocation1: item.OrderBookingQty1, // For new product, allocation = approval (because booking = 0)
						QtyOrderAllocation2: item.OrderBookingQty2,
						QtyOrderAllocation3: item.OrderBookingQty3,
						PurchPrice1:         *item.PurchPrice1,
						PurchPrice2:         *item.PurchPrice2,
						PurchPrice3:         *item.PurchPrice3,
						EstimatedPrice:      estimatedPrice,
						ReturnReasonID:      &returnReasonID,
						CreatedBy:           userID,
						IsDel:               false,
					}

					if err := service.ReplenishmentRepository.CreateDetailForApproval(txCtx, newDetail); err != nil {
						return err
					}

					// Add newly created detail to maps for deletion lookup
					existingDetailMap[newDetail.ReplenishmentDetailID] = newDetail
					existingDetailMapByProID[newDetail.ProID] = newDetail
				}
			}

			// DELETE details from deleted_data array (soft delete)
			// Since only principal can reach here (checked above), don't filter by cust_id for query
			// Before soft delete, update qty_order_allocation for deleted items
			for _, detID := range request.DeletedData {
				// Try to find detail by replenishment_detail_id first
				deletedDetail, exists := existingDetailMap[detID]
				if !exists {
					// If not found, try to find by pro_id (in case deleted_data contains pro_id instead of replenishment_detail_id)
					deletedDetail, exists = existingDetailMapByProID[detID]
					if !exists {
						return fmt.Errorf("replenishment_detail_id or pro_id %d not found for deletion", detID)
					}
				}

				// Update qty_order_allocation for deleted items
				// qty_total is calculated as: qty_ro + qty_order_allocation (for status != 1)
				// User requested:
				//   1. qty_order_allocation = qty_ro (value from qty_ro)
				//   2. qty_total = 0
				// To satisfy both requirements, we need: qty_order_allocation = -qty_ro (negative of qty_ro)
				// This will make: qty_total = qty_ro + (-qty_ro) = 0
				// Setting allocation with negative value of qty_ro to achieve qty_total = 0
				allocation1 := -deletedDetail.OrderBookingQty1
				allocation2 := -deletedDetail.OrderBookingQty2
				allocation3 := -deletedDetail.OrderBookingQty3
				deletedDetail.QtyOrderAllocation1 = &allocation1
				deletedDetail.QtyOrderAllocation2 = &allocation2
				deletedDetail.QtyOrderAllocation3 = &allocation3

				deletedDetail.UpdatedBy = &userID

				// Update the detail before soft delete
				// Use replenishment_detail_id for update
				if err := service.ReplenishmentRepository.UpdateDetail(txCtx, deletedDetail); err != nil {
					return fmt.Errorf("failed to update detail %d before deletion: %w", deletedDetail.ReplenishmentDetailID, err)
				}

				// Then soft delete using replenishment_detail_id
				if err := service.ReplenishmentRepository.SoftDeleteDetail(txCtx, deletedDetail.ReplenishmentDetailID, "", userID); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return err
}

func (service *replenishmentServiceImpl) ProductGrList(dataFilter entity.ProductGrListQueryFilter, custId, parentCustId string) (*entity.ProductGrListResponse, int64, int, error) {
	ctx := context.Background()

	header, details, total, lastPage, err := service.ReplenishmentRepository.FindProductGrList(ctx, dataFilter, custId, parentCustId)
	if err != nil {
		return nil, 0, 0, err
	}

	response := &entity.ProductGrListResponse{
		ReplenishmentNo:   header.ReplenishmentNo,
		ReplenishmentType: header.ReplenishmentType,
		Details:           []entity.ProductGrListDetailResponse{},
	}

	// Map details
	for _, detail := range details {
		detailResp := entity.ProductGrListDetailResponse{
			ReplenishmentDetID: detail.ReplenishmentDetailID,
			ProID:              detail.ProID,
			StockReceived1:     0, // temporary set 0
			StockReceived2:     0, // temporary set 0
			StockReceived3:     0, // temporary set 0
			PurchPrice1:        detail.PurchPrice1,
			PurchPrice2:        detail.PurchPrice2,
			PurchPrice3:        detail.PurchPrice3,
		}

		if detail.ProCode != nil {
			detailResp.ProCode = *detail.ProCode
		}
		if detail.ProName != nil {
			detailResp.ProName = *detail.ProName
		}
		if detail.Vat != nil {
			detailResp.Vat = detail.Vat
		}
		// stock_shipment from qty_order_approval (default 0 if nil)
		if detail.QtyOrderApproval1 != nil {
			detailResp.StockShipment1 = *detail.QtyOrderApproval1
		}
		if detail.QtyOrderApproval2 != nil {
			detailResp.StockShipment2 = *detail.QtyOrderApproval2
		}
		if detail.QtyOrderApproval3 != nil {
			detailResp.StockShipment3 = *detail.QtyOrderApproval3
		}
		if detail.UnitID1 != nil {
			detailResp.UnitID1 = *detail.UnitID1
		}
		if detail.UnitID2 != nil {
			detailResp.UnitID2 = *detail.UnitID2
		}
		if detail.UnitID3 != nil {
			detailResp.UnitID3 = *detail.UnitID3
		}

		detailResp.SubTotal = (detailResp.StockReceived1 * detailResp.PurchPrice1) +
			(detailResp.StockReceived2 * detailResp.PurchPrice2) +
			(detailResp.StockReceived3 * detailResp.PurchPrice3)

		response.Details = append(response.Details, detailResp)
	}

	return response, total, lastPage, nil
}

func (service *replenishmentServiceImpl) PoList(dataFilter entity.PoListQueryFilter, custId, parentCustId string) ([]entity.PoListResponse, int64, int, error) {
	ctx := context.Background()

	poList, total, lastPage, err := service.ReplenishmentRepository.FindPoList(ctx, dataFilter, custId, parentCustId)
	if err != nil {
		return nil, 0, 0, err
	}

	var responses []entity.PoListResponse
	for _, po := range poList {
		response := entity.PoListResponse{
			ReplenishmentNo:   po.ReplenishmentNo,
			ReplenishmentType: po.ReplenishmentType,
			WhID:              po.WhID,
			SupID:             po.SupID,
		}

		if po.WhCode != nil {
			response.WhCode = *po.WhCode
		}
		if po.WhName != nil {
			response.WhName = *po.WhName
		}
		if po.SupCode != nil {
			response.SupCode = *po.SupCode
		}
		if po.SupName != nil {
			response.SupName = *po.SupName
		}

		responses = append(responses, response)
	}

	return responses, total, lastPage, nil
}

func (service *replenishmentServiceImpl) ApprovalProductList(dataFilter entity.ReplenishmentApprovalProductQueryFilter, custId, parentCustId string, isPrincipal bool) ([]entity.ReplenishmentApprovalProductResponse, int64, int, error) {
	ctx := context.Background()
	products, total, lastPage, err := service.ReplenishmentRepository.FindApprovalProducts(ctx, dataFilter, custId, parentCustId, isPrincipal)
	if err != nil {
		return nil, 0, 0, err
	}

	var resp []entity.ReplenishmentApprovalProductResponse
	for i := range products {
		p := products[i]
		resp = append(resp, entity.ReplenishmentApprovalProductResponse{
			ProID:           p.ProID,
			ProCode:         p.ProCode,
			ProName:         p.ProName,
			Ripening:        p.Ripening,
			InTransitStock1: p.InTransitStock1,
			InTransitStock2: p.InTransitStock2,
			InTransitStock3: p.InTransitStock3,
			SafStockQty:     p.SafStockQty,
			MinStockQty:     p.MinStockQty,
			Vat:             p.Vat,
			ConvUnit2:       p.ConvUnit2,
			ConvUnit3:       p.ConvUnit3,
			UnitID1:         p.UnitID1,
			UnitID2:         p.UnitID2,
			UnitID3:         p.UnitID3,
			PurchPrice1:     p.PurchPrice1,
			PurchPrice2:     p.PurchPrice2,
			PurchPrice3:     p.PurchPrice3,
			Qty1:            p.Qty1,
			Qty2:            p.Qty2,
			Qty3:            p.Qty3,
			TotalQty:        p.TotalQty,
		})
	}

	return resp, total, lastPage, nil
}

func (service *replenishmentServiceImpl) CheckIsPrincipal(custID string) (bool, error) {
	ctx := context.Background()
	return service.ReplenishmentRepository.CheckIsPrincipal(ctx, custID)
}

func (service *replenishmentServiceImpl) IsAnyApprovalPICUser(userID int64, employeeID int64) (bool, error) {
	ctx := context.Background()
	ok, err := service.ReplenishmentRepository.IsUserApprovalPIC(ctx, userID)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	if employeeID > 0 {
		return service.ReplenishmentRepository.IsUserApprovalPIC(ctx, employeeID)
	}
	return false, nil
}

func (service *replenishmentServiceImpl) BatchUpdateApproval(request entity.BatchReplenishmentApprovalRequest, custID string, userID int64, employeeID int64) ([]entity.BatchReplenishmentApprovalResultItem, error) {
	if len(request.Data) == 0 {
		return nil, fmt.Errorf("data must contain at least one replenishment order")
	}

	results := make([]entity.BatchReplenishmentApprovalResultItem, 0, len(request.Data))
	for _, order := range request.Data {
		if order.ReplenishmentID <= 0 {
			results = append(results, entity.BatchReplenishmentApprovalResultItem{
				ReplenishmentID: order.ReplenishmentID,
				Status:          "failed",
				Message:         "replenishment_id is required and must be greater than zero",
			})
			continue
		}

		effectiveApproval := order.Approval
		if effectiveApproval == nil {
			results = append(results, entity.BatchReplenishmentApprovalResultItem{
				ReplenishmentID: order.ReplenishmentID,
				Status:          "failed",
				Message:         "approval is required on each data item (with replenishment_id)",
			})
			continue
		}
		details := make([]entity.UpdateReplenishmentApprovalDetail, 0, len(order.Details))
		for i := range order.Details {
			details = append(details, order.Details[i].ToUpdateReplenishmentApprovalDetail())
		}

		var singleRemarks *string
		if order.Remarks != nil && strings.TrimSpace(*order.Remarks) != "" {
			r := strings.TrimSpace(*order.Remarks)
			singleRemarks = &r
		} else {
			singleRemarks = request.Remarks
		}

		single := entity.UpdateReplenishmentApprovalRequest{
			Approval:    effectiveApproval,
			Remarks:     singleRemarks,
			Data:        details,
			DeletedData: order.DeletedData,
		}

		ctx := context.Background()
		ro, findErr := service.ReplenishmentRepository.GetReplenishmentOrderByID(ctx, order.ReplenishmentID, "")
		if findErr != nil {
			results = append(results, entity.BatchReplenishmentApprovalResultItem{
				ReplenishmentID: order.ReplenishmentID,
				Status:          "failed",
				Message:         findErr.Error(),
			})
			continue
		}

		err := service.updateApprovalForOrder(ctx, ro, single, custID, userID, employeeID)
		status := "approved"
		if !*effectiveApproval {
			status = "rejected"
		}
		if err != nil {
			results = append(results, entity.BatchReplenishmentApprovalResultItem{
				ReplenishmentID: order.ReplenishmentID,
				ReplenishmentNo: ro.ReplenishmentNo,
				Status:          "failed",
				Message:         err.Error(),
			})
			continue
		}
		results = append(results, entity.BatchReplenishmentApprovalResultItem{
			ReplenishmentID: order.ReplenishmentID,
			ReplenishmentNo: ro.ReplenishmentNo,
			Status:          status,
		})
	}
	return results, nil
}

func (service *replenishmentServiceImpl) SummarizeReplanishment(replanishmentIDs []int64, custID, parentCustID string, isPrincipal bool, requestID string, userID int64) ([]entity.SummarizeReplanishmentResponse, error) {
	ctx := context.Background()
	rows, err := service.ReplenishmentRepository.FindSummarizeReplanishment(ctx, replanishmentIDs, custID, parentCustID, isPrincipal, userID)
	if err != nil {
		return nil, err
	}

	headerMap := make(map[int64]*entity.SummarizeReplanishmentResponse)
	ordered := make([]entity.SummarizeReplanishmentResponse, 0)
	for i := range rows {
		row := rows[i]
		if _, exists := headerMap[row.ReplanishmentID]; !exists {
			header := entity.SummarizeReplanishmentResponse{
				ReplanishmentID: row.ReplanishmentID,
				ReplanishmentNo: row.ReplanishmentNo,
				DisributorID:    row.DisributorID,
				SupID:           row.SupID,
				Details:         make([]entity.SummarizeReplanishmentDetailItem, 0),
			}
			if row.DistributorCode != nil {
				header.DistributorCode = *row.DistributorCode
			}
			if row.DistributorName != nil {
				header.DistributorName = *row.DistributorName
			}
			if row.SupCode != nil {
				header.SupCode = *row.SupCode
			}
			if row.SupName != nil {
				header.SupName = *row.SupName
			}
			header.WhID = row.WhID
			if row.WhCode != nil {
				header.WhCode = *row.WhCode
			}
			if row.WhName != nil {
				header.WhName = *row.WhName
			}
			ordered = append(ordered, header)
			headerMap[row.ReplanishmentID] = &ordered[len(ordered)-1]
		}

		detail := entity.SummarizeReplanishmentDetailItem{
			ReplanishmentDetailID: row.ReplanishmentDetailID,
			ProID:                 row.ProID,
			PurchPrice1:           row.PurchPrice1,
			PurchPrice2:           row.PurchPrice2,
			PurchPrice3:           row.PurchPrice3,
			QtyRo1:                row.QtyRo1,
			QtyRo2:                row.QtyRo2,
			QtyRo3:                row.QtyRo3,
			EstimatedPrice:        row.EstimatedPrice,
			RequestID:             requestID,
			ReturnReasonID:        row.ReturnReasonID,
		}
		if row.ProCode != nil {
			detail.ProCode = *row.ProCode
		}
		if row.ProName != nil {
			detail.ProName = *row.ProName
		}
		if row.WhStockLarge != nil {
			detail.WhStockLarge = int64(math.Round(*row.WhStockLarge))
		}
		if row.WhStockMedium != nil {
			detail.WhStockMedium = int64(math.Round(*row.WhStockMedium))
		}
		if row.WhStockSmall != nil {
			detail.WhStockSmall = int64(math.Round(*row.WhStockSmall))
		}
		if row.OptimumQty != nil {
			detail.OptimumQty = int64(math.Round(*row.OptimumQty))
		}
		if row.Ripening != nil {
			detail.Ripening = int64(math.Round(*row.Ripening))
		}
		if row.UnitID1 != nil {
			detail.UnitID1 = *row.UnitID1
		}
		if row.UnitID2 != nil {
			detail.UnitID2 = *row.UnitID2
		}
		if row.UnitID3 != nil {
			detail.UnitID3 = *row.UnitID3
		}

		headerMap[row.ReplanishmentID].Details = append(headerMap[row.ReplanishmentID].Details, detail)
	}

	return ordered, nil
}
