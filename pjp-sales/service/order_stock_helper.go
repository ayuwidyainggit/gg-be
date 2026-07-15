package service

import "sales/pkg/conversion"

type apiStockBreakdown struct {
	Qty1 float64
	Qty2 float64
	Qty3 float64
}

func safeConvUnits(convUnit2, convUnit3 int) (int, int) {
	if convUnit2 <= 0 {
		return 0, 0
	}
	if convUnit3 <= 0 {
		return convUnit2, 0
	}
	return convUnit2, convUnit3
}

func toTotalSmallFromAPIUnits(qtyLarge, qtyMedium, qtySmall float64, convUnit2, convUnit3 int) int {
	conv2, conv3 := safeConvUnits(convUnit2, convUnit3)
	total := int(qtySmall)
	if conv2 <= 0 {
		return total
	}

	total += int(qtyMedium) * conv2
	if conv3 <= 0 {
		total += int(qtyLarge) * conv2
		return total
	}

	total += int(qtyLarge) * conv2 * conv3
	return total
}

func canonicalAPIStockBreakdown(totalSmall, convUnit2, convUnit3 int) apiStockBreakdown {
	if totalSmall == 0 {
		return apiStockBreakdown{}
	}

	conv2, conv3 := safeConvUnits(convUnit2, convUnit3)
	if conv2 <= 0 {
		return apiStockBreakdown{Qty3: float64(totalSmall)}
	}
	if conv3 <= 0 {
		qty1 := totalSmall / conv2
		qty3 := totalSmall % conv2
		return apiStockBreakdown{Qty1: float64(qty1), Qty2: 0, Qty3: float64(qty3)}
	}

	qty := &conversion.Qty{Qty: totalSmall, ConvUnit2: conv2, ConvUnit3: conv3}
	converted := qty.ConvToQtyConversion()
	return apiStockBreakdown{Qty1: float64(converted.Qty3), Qty2: float64(converted.Qty2), Qty3: float64(converted.Qty1)}
}

func computeDisplayedAvailableStockBreakdown(warehouseSmall int, qtyLarge, qtyMedium, qtySmall float64, includeOrder bool, convUnit2, convUnit3 int) apiStockBreakdown {
	totalSmall := warehouseSmall
	if includeOrder {
		totalSmall += toTotalSmallFromAPIUnits(qtyLarge, qtyMedium, qtySmall, convUnit2, convUnit3)
	}
	return canonicalAPIStockBreakdown(totalSmall, convUnit2, convUnit3)
}

func applyStockBreakdownToPointers(targetQty1, targetQty2, targetQty3 **float64, breakdown apiStockBreakdown) {
	qty1 := breakdown.Qty1
	qty2 := breakdown.Qty2
	qty3 := breakdown.Qty3
	*targetQty1 = &qty1
	*targetQty2 = &qty2
	*targetQty3 = &qty3
}
