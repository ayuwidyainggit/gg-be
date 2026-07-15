package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"inventory/entity"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	sapRespSuccess        = "success"
	sapRespPartialSuccess = "partial_success"
	sapRespError          = "error"
)

func newSAPReplRequestID() string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("REQ-%s-%d", time.Now().UTC().Format("20060102"), time.Now().UnixNano()%100000)
	}
	return fmt.Sprintf("REQ-%s-%s", time.Now().UTC().Format("20060102"), strings.ToUpper(hex.EncodeToString(b)))
}

func isKnownReplenishmentStatus(code int) bool {
	switch code {
	case StatusNeedReview, StatusApproved, StatusRejected, StatusOnDelivery,
		StatusCompleted, StatusProcessed, StatusCancelled:
		return true
	default:
		return false
	}
}

func sapLifecycleRank(st int) int {
	switch st {
	case StatusNeedReview:
		return 10
	case StatusApproved:
		return 20
	case StatusOnDelivery:
		return 30
	case StatusProcessed:
		return 35
	case StatusCompleted:
		return 40
	case StatusRejected:
		return 100
	case StatusCancelled:
		return 101
	default:
		return 0
	}
}

func assertSAPAllowedTransition(from, to int) error {
	if from == StatusRejected || from == StatusCancelled {
		return fmt.Errorf("replenishment is in terminal status %d and cannot be updated", from)
	}
	if to == StatusRejected || to == StatusCancelled {
		return nil
	}
	rFrom, rTo := sapLifecycleRank(from), sapLifecycleRank(to)
	if rFrom == 0 || rTo == 0 {
		return nil
	}
	if rTo < rFrom {
		return errors.New("status cannot move backward")
	}
	return nil
}

func ptrEqualFoldUnit(a *string, b string) bool {
	b = strings.TrimSpace(b)
	if a == nil {
		return b == ""
	}
	return strings.EqualFold(strings.TrimSpace(*a), b)
}

func (service *replenishmentServiceImpl) resolveSAPReplCustID(
	ctx context.Context,
	req entity.SAPReplStatusRequest,
) (string, []entity.SAPReplFieldError) {
	custID := strings.TrimSpace(req.CustID)
	distributorCode := req.DistributorCode

	if custID == "" && distributorCode <= 0 {
		return "", []entity.SAPReplFieldError{
			{Field: "cust_id", Message: "cust_id or distributor_code is required"},
			{Field: "distributor_code", Message: "cust_id or distributor_code is required"},
		}
	}

	if distributorCode <= 0 {
		return custID, nil
	}

	resolvedCustID, err := service.ReplenishmentRepository.FindCustIDByDistributorCode(ctx, distributorCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", []entity.SAPReplFieldError{
				{Field: "distributor_code", Message: "distributor not found"},
			}
		}
		return "", []entity.SAPReplFieldError{
			{Field: "distributor_code", Message: err.Error()},
		}
	}

	if custID == "" {
		return resolvedCustID, nil
	}

	if !strings.EqualFold(custID, resolvedCustID) {
		return "", []entity.SAPReplFieldError{
			{Field: "cust_id", Message: "cust_id does not match distributor_code"},
			{Field: "distributor_code", Message: "cust_id does not match distributor_code"},
		}
	}

	return custID, nil
}

func (service *replenishmentServiceImpl) SAPUpdateReplenishmentStatus(
	req entity.SAPReplStatusRequest,
	updatedBy int64,
) entity.SAPReplStatusResponse {
	requestID := newSAPReplRequestID()
	ctx := context.Background()

	custID, custErrs := service.resolveSAPReplCustID(ctx, req)
	if len(custErrs) > 0 {
		return entity.SAPReplStatusResponse{
			RequestID: requestID,
			Status:    sapRespError,
			Message:   "Validation error",
			Errors: []entity.SAPReplStatusReplErrWrap{
				{ReplenishmentNo: "", Errors: custErrs},
			},
		}
	}

	if len(req.Replenishments) == 0 {
		return entity.SAPReplStatusResponse{
			RequestID: requestID,
			Status:    sapRespError,
			Message:   "Validation error",
			Errors: []entity.SAPReplStatusReplErrWrap{
				{ReplenishmentNo: "", Errors: []entity.SAPReplFieldError{{Field: "replenishments", Message: "replenishments must contain at least one item"}}},
			},
		}
	}

	var successItems []entity.SAPReplStatusOKItem
	var failedItems []entity.SAPReplStatusFailedItem

	for _, rep := range req.Replenishments {
		if fieldErrs := validateSAPReplenishmentPayload(rep); len(fieldErrs) > 0 {
			failedItems = append(failedItems, entity.SAPReplStatusFailedItem{
				ReplenishmentNo: rep.ReplenishmentNo,
				Status:          sapRespError,
				Errors:          fieldErrs,
			})
			continue
		}

		err := service.Transaction.WithinTransaction(ctx, func(txCtx context.Context) error {
			order, err := service.ReplenishmentRepository.FindByReplenishmentNo(txCtx, rep.ReplenishmentNo, custID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("replenishment not found")
				}
				return err
			}
			if order.IsDel {
				return fmt.Errorf("replenishment not found")
			}

			if err := assertSAPAllowedTransition(order.Status, rep.Status); err != nil {
				return err
			}

			if err := service.ReplenishmentRepository.UpdateStatusByReplenishmentNo(txCtx, rep.ReplenishmentNo, custID, rep.Status, updatedBy); err != nil {
				return err
			}

			for i, d := range rep.Details {
				prefix := fmt.Sprintf("details[%d].", i)
				if strings.TrimSpace(d.ProCode.String()) == "" {
					return fmt.Errorf("%spro_code is required", prefix)
				}
				if strings.TrimSpace(d.UnitID3) == "" {
					return fmt.Errorf("%sunit_id3 is required", prefix)
				}

				proID, unitMst, err := service.ReplenishmentRepository.LookupProductByProCode(txCtx, custID, d.ProCode.String())
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("%spro_code: product not found", prefix)
					}
					return err
				}
				if !ptrEqualFoldUnit(unitMst, d.UnitID3) {
					return fmt.Errorf("%sunit_id3 does not match product master", prefix)
				}

				det, err := service.ReplenishmentRepository.FindActiveReplenishmentDetail(txCtx, custID, order.ReplenishmentID, proID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("%spro_code: line not found on replenishment", prefix)
					}
					return err
				}

				if err := service.ReplenishmentRepository.UpdateReplenishmentDetailSAPFields(txCtx, custID, det.ReplenishmentDetailID, d.Qty3, d.PurchPrice3, updatedBy); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			failedItems = append(failedItems, entity.SAPReplStatusFailedItem{
				ReplenishmentNo: rep.ReplenishmentNo,
				Status:          sapRespError,
				Errors: []entity.SAPReplFieldError{
					{Field: "replenishment_no", Message: err.Error()},
				},
			})
			continue
		}
		successItems = append(successItems, entity.SAPReplStatusOKItem{
			ReplenishmentNo: rep.ReplenishmentNo,
			Status:          sapRespSuccess,
		})
	}

	if len(failedItems) == 0 {
		return entity.SAPReplStatusResponse{
			RequestID: requestID,
			Status:    sapRespSuccess,
			Message:   "All replenishment status updated successfully",
			Data:      successItems,
		}
	}

	if len(successItems) == 0 {
		return entity.SAPReplStatusResponse{
			RequestID: requestID,
			Status:    sapRespError,
			Message:   "Validation error",
			Errors:    failedAsValidationWrap(failedItems),
		}
	}

	return entity.SAPReplStatusResponse{
		RequestID: requestID,
		Status:    sapRespPartialSuccess,
		Message:   "Some replenishments failed to update",
		Data: entity.SAPReplStatusPartialData{
			Success: successItems,
			Failed:  failedItems,
		},
	}
}

func failedAsValidationWrap(failed []entity.SAPReplStatusFailedItem) []entity.SAPReplStatusReplErrWrap {
	out := make([]entity.SAPReplStatusReplErrWrap, 0, len(failed))
	for _, f := range failed {
		out = append(out, entity.SAPReplStatusReplErrWrap{
			ReplenishmentNo: f.ReplenishmentNo,
			Errors:          f.Errors,
		})
	}
	return out
}

func validateSAPReplenishmentPayload(rep entity.SAPReplStatusReplenishment) []entity.SAPReplFieldError {
	var errs []entity.SAPReplFieldError
	if strings.TrimSpace(rep.ReplenishmentNo) == "" {
		errs = append(errs, entity.SAPReplFieldError{Field: "replenishment_no", Message: "replenishment_no is required"})
	}
	if !isKnownReplenishmentStatus(rep.Status) {
		errs = append(errs, entity.SAPReplFieldError{Field: "status", Message: "status must be integer and a valid replenishment status code"})
	}
	if len(rep.Details) == 0 {
		errs = append(errs, entity.SAPReplFieldError{Field: "details", Message: "details must contain at least one line"})
	}
	for i, d := range rep.Details {
		pref := fmt.Sprintf("details[%d].", i)
		if strings.TrimSpace(d.ProCode.String()) == "" {
			errs = append(errs, entity.SAPReplFieldError{Field: pref + "pro_code", Message: "pro_code is required"})
		}
		if strings.TrimSpace(d.UnitID3) == "" {
			errs = append(errs, entity.SAPReplFieldError{Field: pref + "unit_id3", Message: "unit_id3 is required"})
		}
	}
	return errs
}
