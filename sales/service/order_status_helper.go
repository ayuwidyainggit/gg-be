package service

import (
	"errors"
	"strings"

	"sales/entity"
	"sales/model"
)

const (
	LIMIT_ACTION_WARNING    = 1
	LIMIT_ACTION_RESTRICTED = 2

	ORDER_STATUS_REASON_VALIDATE1_FAILED          = "validate1_failed"
	ORDER_STATUS_REASON_RESTRICTED_BY_OUTLET_RULE = "restricted_by_outlet_rule"
	ORDER_STATUS_REASON_PROCESSED                 = "processed"
	ORDER_STATUS_REASON_VALIDATION_NOT_CONFIGURED = "validation_not_configured"
	ORDER_STATUS_VALIDATE1_FAILED_MESSAGE         = "order cannot be processed because stock validation failed"
)

type salesOrderOutletRules struct {
	CreditLimitAction       *int
	CreditLimitActionName   string
	SalesInvLimitAction     *int
	SalesInvLimitActionName string
	ObsLimitAction          *int
	ObsLimitActionName      string
}

type salesOrderStatusDecision struct {
	DataStatus int64
	Blocked    bool
	Reason     string
}

func outletRulesFromOutletRead(outlet model.OutletRead) salesOrderOutletRules {
	return salesOrderOutletRules{
		CreditLimitAction:       outlet.CreditLimitAction,
		CreditLimitActionName:   outlet.CreditLimitActionName,
		SalesInvLimitAction:     outlet.SalesInvLimitAction,
		SalesInvLimitActionName: outlet.SalesInvLimitActionName,
		ObsLimitAction:          outlet.ObsLimitAction,
		ObsLimitActionName:      outlet.ObsLimitActionName,
	}
}

func outletRulesFromOrderList(order model.OrderList) salesOrderOutletRules {
	return salesOrderOutletRules{
		CreditLimitAction:       order.CreditLimitAction,
		CreditLimitActionName:   order.CreditLimitActionName,
		SalesInvLimitAction:     order.SalesInvLimitAction,
		SalesInvLimitActionName: order.SalesInvLimitActionName,
		ObsLimitAction:          order.ObsLimitAction,
		ObsLimitActionName:      order.ObsLimitActionName,
	}
}

func normalizeLimitAction(action *int, actionName string) int {
	if action != nil {
		return *action
	}

	normalized := strings.ToLower(strings.TrimSpace(actionName))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	switch normalized {
	case "warning", "warn":
		return LIMIT_ACTION_WARNING
	case "restrict", "restricted", "restriction":
		return LIMIT_ACTION_RESTRICTED
	default:
		return 0
	}
}

func isRestrictedLimitAction(action *int, actionName string) bool {
	return normalizeLimitAction(action, actionName) == LIMIT_ACTION_RESTRICTED
}

func determineSalesOrderStatus(validation entity.ValidateResponse, outlet salesOrderOutletRules) salesOrderStatusDecision {
	if !validation.Validate1Success {
		return salesOrderStatusDecision{
			DataStatus: int64(entity.NEED_REVIEW),
			Blocked:    true,
			Reason:     ORDER_STATUS_REASON_VALIDATE1_FAILED,
		}
	}

	if !validation.Validate2Success && isRestrictedLimitAction(outlet.CreditLimitAction, outlet.CreditLimitActionName) {
		return salesOrderStatusDecision{
			DataStatus: int64(entity.NEED_REVIEW),
			Reason:     ORDER_STATUS_REASON_RESTRICTED_BY_OUTLET_RULE,
		}
	}

	if !validation.Validate3Success && isRestrictedLimitAction(outlet.SalesInvLimitAction, outlet.SalesInvLimitActionName) {
		return salesOrderStatusDecision{
			DataStatus: int64(entity.NEED_REVIEW),
			Reason:     ORDER_STATUS_REASON_RESTRICTED_BY_OUTLET_RULE,
		}
	}

	if !validation.Validate4Success && isRestrictedLimitAction(outlet.ObsLimitAction, outlet.ObsLimitActionName) {
		return salesOrderStatusDecision{
			DataStatus: int64(entity.NEED_REVIEW),
			Reason:     ORDER_STATUS_REASON_RESTRICTED_BY_OUTLET_RULE,
		}
	}

	return salesOrderStatusDecision{
		DataStatus: int64(entity.PROCESSED),
		Reason:     ORDER_STATUS_REASON_PROCESSED,
	}
}

func ensureSalesOrderStatusDecisionAllowed(decision salesOrderStatusDecision) error {
	if decision.Blocked {
		return errors.New(ORDER_STATUS_VALIDATE1_FAILED_MESSAGE)
	}

	return nil
}

func applyValidationResultToOrderModel(orderModel *model.Order, validation entity.ValidateResponse) {
	orderModel.ValidateStok = &validation.Validate1Success
	orderModel.ValidateStokMessage = nullableValidationMessage(validation.Validate1)
	orderModel.ValidateCreditLimit = &validation.Validate2Success
	orderModel.ValidateCreditLimitMessage = validation.Validate2
	orderModel.ValidateCreditLimitValue = validation.Validate2value
	orderModel.ValidateOverdue = &validation.Validate3Success
	orderModel.ValidateOverdueMessage = validation.Validate3
	orderModel.ValidateOverdueValue = validation.Validate3Value
	orderModel.ValidateOutstanding = &validation.Validate4Success
	orderModel.ValidateOutstandingMessage = validation.Validate4
	orderModel.ValidateOutstandingValue = validation.Validate4Value
	orderModel.ValidateSummary = validation.IsSuccessValidate
}

func validationResultFromOrderList(order model.OrderList) entity.ValidateResponse {
	if !hasStoredValidationResult(order) {
		return entity.ValidateResponse{
			Validate1Success:  true,
			Validate2Success:  true,
			Validate3Success:  true,
			Validate4Success:  true,
			IsSuccessValidate: true,
		}
	}

	return entity.ValidateResponse{
		Validate1Success:  order.ValidateStok,
		Validate1:         stringFromPtr(order.ValidateStokMessage),
		Validate2Success:  order.ValidateCreditLimit,
		Validate2:         order.ValidateCreditLimitMessage,
		Validate2value:    order.ValidateCreditLimitValue,
		Validate3Success:  order.ValidateOverdue,
		Validate3:         order.ValidateOverdueMessage,
		Validate3Value:    order.ValidateOverdueValue,
		Validate4Success:  order.ValidateOutstanding,
		Validate4:         order.ValidateOutstandingMessage,
		Validate4Value:    order.ValidateOutstandingValue,
		IsSuccessValidate: order.ValidateSummary,
	}
}

func hasStoredValidationResult(order model.OrderList) bool {
	return order.ValidateSummary ||
		order.ValidateStok ||
		order.ValidateCreditLimit ||
		order.ValidateOverdue ||
		order.ValidateOutstanding ||
		order.ValidateStokMessage != nil ||
		order.ValidateCreditLimitMessage != "" ||
		order.ValidateCreditLimitValue != 0 ||
		order.ValidateOverdueMessage != "" ||
		order.ValidateOverdueValue != 0 ||
		order.ValidateOutstandingMessage != "" ||
		order.ValidateOutstandingValue != 0
}

func isProcessedDataStatus(dataStatus *int64) bool {
	return dataStatus != nil && *dataStatus == int64(entity.PROCESSED)
}

func isValidationResponseEmpty(validation entity.ValidateResponse) bool {
	return !validation.Validate1Success &&
		!validation.Validate2Success &&
		!validation.Validate3Success &&
		!validation.Validate4Success &&
		!validation.IsSuccessValidate &&
		validation.Validate1 == "" &&
		validation.Validate2 == "" &&
		validation.Validate3 == "" &&
		validation.Validate4 == "" &&
		validation.Validate2value == 0 &&
		validation.Validate3Value == 0 &&
		validation.Validate4Value == 0
}

func (service *orderServiceImpl) determineStatusForExistingOrder(roNo string, custID string) (salesOrderStatusDecision, error) {
	ro, err := service.OrderRepository.FindByNo(roNo, custID)
	if err != nil {
		return salesOrderStatusDecision{}, err
	}

	validation := validationResultFromOrderList(ro)
	if service.ValidateOrderRepository != nil && ro.OutletID != nil && ro.WhId != nil {
		details, err := service.OrderRepository.FindDetail(roNo, custID)
		if err != nil {
			return salesOrderStatusDecision{}, err
		}

		validationRequest := entity.ValidateOrderBody{
			CustID:   custID,
			OutletID: *ro.OutletID,
			WhID:     int(*ro.WhId),
			Total:    getValueOrDefault(ro.Total, 0),
		}
		for _, detail := range details {
			if detail.ItemType != 1 {
				continue
			}

			validationRequest.ProStok = append(validationRequest.ProStok, entity.ProductsValidate{
				ProductId: int64(detail.ProId),
				Qty1:      int64(getValueOrDefault(detail.Qty1Final, getValueOrDefault(detail.Qty1, 0))),
				Qty2:      int64(getValueOrDefault(detail.Qty2Final, getValueOrDefault(detail.Qty2, 0))),
				Qty3:      int64(getValueOrDefault(detail.Qty3Final, getValueOrDefault(detail.Qty3, 0))),
			})
		}

		validateOrderService := NewValidateOrderService(service.ValidateOrderRepository, service.Transaction)
		validation, _, _, err = validateOrderService.ValidateOrder(validationRequest)
		if err != nil {
			return salesOrderStatusDecision{}, err
		}
	}

	return determineSalesOrderStatus(validation, outletRulesFromOrderList(ro)), nil
}

func (service *orderServiceImpl) determineStatusForOrderProjection(ro model.OrderList, custID string, details []model.OrderDetailRead, headerUpdate model.Order) (salesOrderStatusDecision, error) {
	validation := validationResultFromOrderList(ro)
	if len(details) == 0 || !hasStoredValidationResult(ro) || service.ValidateOrderRepository == nil || ro.OutletID == nil || ro.WhId == nil {
		return determineSalesOrderStatus(validation, outletRulesFromOrderList(ro)), nil
	}

	validationRequest := entity.ValidateOrderBody{
		CustID:   custID,
		OutletID: *ro.OutletID,
		WhID:     int(*ro.WhId),
		Total:    getValueOrDefault(headerUpdate.TotalFinal, getValueOrDefault(headerUpdate.Total, getValueOrDefault(ro.TotalFinal, getValueOrDefault(ro.Total, 0)))),
	}

	for _, detail := range details {
		if detail.ItemType != 1 {
			continue
		}

		validationRequest.ProStok = append(validationRequest.ProStok, entity.ProductsValidate{
			ProductId: int64(detail.ProId),
			Qty1:      int64(getValueOrDefault(detail.Qty1Final, getValueOrDefault(detail.Qty1, 0))),
			Qty2:      int64(getValueOrDefault(detail.Qty2Final, getValueOrDefault(detail.Qty2, 0))),
			Qty3:      int64(getValueOrDefault(detail.Qty3Final, getValueOrDefault(detail.Qty3, 0))),
		})
	}

	validateOrderService := NewValidateOrderService(service.ValidateOrderRepository, service.Transaction)
	validation, _, _, err := validateOrderService.ValidateOrder(validationRequest)
	if err != nil {
		return salesOrderStatusDecision{}, err
	}

	return determineSalesOrderStatus(validation, outletRulesFromOrderList(ro)), nil
}
