package service

import (
	"strings"

	"sales/entity"
	"sales/model"
)

const orderTypeTakingOrder = "O"

func NormalizedOrderType(orderType *string) string {
	if orderType == nil {
		return ""
	}
	return strings.TrimSpace(*orderType)
}

func IsTakingOrder(orderType *string) bool {
	return NormalizedOrderType(orderType) == orderTypeTakingOrder
}

func ShouldValidateStockOnCreate(orderType *string) bool {
	return !IsTakingOrder(orderType)
}

func ShouldMutateInventoryOnCreate(orderType *string) bool {
	return !IsTakingOrder(orderType)
}

func resolveCreateOrderDataStatus(orderType *string, statusDecision salesOrderStatusDecision) int64 {
	if IsTakingOrder(orderType) {
		return int64(entity.NEED_REVIEW)
	}

	return statusDecision.DataStatus
}

func BuildCreateOrderValidationBypassResponse(orderType *string) entity.ValidateResponse {
	if !IsTakingOrder(orderType) {
		return entity.ValidateResponse{}
	}

	return entity.ValidateResponse{
		Validate1Success: true,
		Validate1:        "Sufficient Stock",
		Validate2Success: true,
		Validate2:        "Within Limit",
		Validate3Success: true,
		Validate3:        "Allowed",
		Validate4Success: true,
		Validate4:        "Allowed",
	}
}

func takingOrderQtySource(detail entity.CreateOrderDetBody) (float64, float64, float64) {
	qty1 := getValueOrDefault(detail.Qty1, getValueOrDefault(detail.QtyPo1, 0))
	qty2 := getValueOrDefault(detail.Qty2, getValueOrDefault(detail.QtyPo2, 0))
	qty3 := getValueOrDefault(detail.Qty3, getValueOrDefault(detail.QtyPo3, 0))

	return qty1, qty2, qty3
}

func nullableValidationMessage(message string) *string {
	if strings.TrimSpace(message) == "" {
		return nil
	}
	return stringPtr(message)
}

func applyTakingOrderValidationSnapshot(orderModel *model.Order) {
	validateStok := false
	orderModel.ValidateStok = &validateStok
	orderModel.ValidateStokMessage = nil
}

func applyTakingOrderDetailFields(detail entity.CreateOrderDetBody, target *model.OrderDetail, totalQty float64) {
	qtyPo1, qtyPo2, qtyPo3 := takingOrderQtySource(detail)

	target.QtyPo = totalQty
	target.QtyPo1 = float64Ptr(qtyPo1)
	target.QtyPo2 = float64Ptr(qtyPo2)
	target.QtyPo3 = float64Ptr(qtyPo3)
	target.OriginalQtyPo1 = float64Ptr(qtyPo1)
	target.OriginalQtyPo2 = float64Ptr(qtyPo2)
	target.OriginalQtyPo3 = float64Ptr(qtyPo3)
	target.SellPricePo1 = detail.SellPrice1
	target.SellPricePo2 = detail.SellPrice2
	target.SellPricePo3 = detail.SellPrice3

	if detail.DiscPo != nil {
		target.DiscPo = detail.DiscPo
	} else {
		target.DiscPo = detail.DiscValue
	}

	if detail.VatValuePo != nil {
		target.VatValuePo = detail.VatValuePo
	} else {
		target.VatValuePo = detail.VatValue
	}

	target.PromoPo1 = target.PromoSo1
	target.PromoPo2 = target.PromoSo2
	target.PromoPo3 = target.PromoSo3
	target.PromoPo4 = target.PromoSo4
	target.PromoPo5 = target.PromoSo5
	target.PromoRemarksPo = model.JSONStringArray(append([]string{}, target.PromoRemarksSo...))
	target.IsProductPromotionPo = target.IsProductPromotionSo
}
