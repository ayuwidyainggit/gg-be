package service

import "testing"

func TestCanonicalAPIStockBreakdown_TotalSmall13_ReturnsLargeMediumSmall(t *testing.T) {
	got := canonicalAPIStockBreakdown(13, 5, 1)
	if got.Qty1 != 2 || got.Qty2 != 0 || got.Qty3 != 3 {
		t.Fatalf("expected 13 small => 2,0,3 got %+v", got)
	}
}

func TestCanonicalAPIStockBreakdown_ExactLarge(t *testing.T) {
	got := canonicalAPIStockBreakdown(10, 5, 2)
	if got.Qty1 != 1 || got.Qty2 != 0 || got.Qty3 != 0 {
		t.Fatalf("expected 10 small => 1,0,0 got %+v", got)
	}
}

func TestCanonicalAPIStockBreakdown_WithMediumRemainder(t *testing.T) {
	got := canonicalAPIStockBreakdown(14, 2, 5)
	if got.Qty1 != 1 || got.Qty2 != 2 || got.Qty3 != 0 {
		t.Fatalf("expected 14 small => 1,2,0 got %+v", got)
	}
}

func TestCanonicalAPIStockBreakdown_WithoutLargeConversion(t *testing.T) {
	got := canonicalAPIStockBreakdown(7, 5, 0)
	if got.Qty1 != 1 || got.Qty2 != 0 || got.Qty3 != 2 {
		t.Fatalf("expected 7 small => 1,0,2 got %+v", got)
	}
}

func TestCanonicalAPIStockBreakdown_WithoutConversionFallsBackToSmallOnly(t *testing.T) {
	got := canonicalAPIStockBreakdown(7, 0, 0)
	if got.Qty1 != 0 || got.Qty2 != 0 || got.Qty3 != 7 {
		t.Fatalf("expected 7 small => 0,0,7 got %+v", got)
	}
}

func TestComputeDisplayedAvailableStockBreakdown_RenormalizesCombinedStock(t *testing.T) {
	got := computeDisplayedAvailableStockBreakdown(11, 0, 0, 2, true, 5, 1)
	if got.Qty1 != 2 || got.Qty2 != 0 || got.Qty3 != 3 {
		t.Fatalf("expected combined stock => 2,0,3 got %+v", got)
	}
}

func TestComputeDisplayedAvailableStockBreakdown_SX2508_ReturnsLargeMediumSmall(t *testing.T) {
	// 13 small units with conv2=5, conv3=2 => 1 large (10), 0 medium, 3 small
	got := computeDisplayedAvailableStockBreakdown(13, 0, 0, 0, true, 5, 2)
	if got.Qty1 != 1 || got.Qty2 != 0 || got.Qty3 != 3 {
		t.Fatalf("SX-2508: expected 13 small => Qty1=1 (Large), Qty2=0 (Medium), Qty3=3 (Small), got %+v", got)
	}
}
