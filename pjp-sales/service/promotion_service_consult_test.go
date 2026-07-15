package service

import (
	"testing"

	"sales/entity"
	"sales/model"
	"sales/pkg/errmsg"
)

func TestConsultPromoRespHasRewardIgnoresZeroQtyProductRewards(t *testing.T) {
	resp := entity.ConsultPromoResp{
		RewardProduct: []entity.PromoRewardProductDet{
			{ProID: 10807, Qty1: 0, Qty2: 0, Qty3: 0},
		},
	}
	if consultPromoRespHasReward(resp) {
		t.Fatal("expected zero-qty product reward to be treated as no reward")
	}
}

func TestConsultPromoRespHasRewardAcceptsPositiveQtyProductRewards(t *testing.T) {
	resp := entity.ConsultPromoResp{
		RewardProduct: []entity.PromoRewardProductDet{
			{ProID: 10807, Qty1: 1},
		},
	}
	if !consultPromoRespHasReward(resp) {
		t.Fatal("expected positive product reward qty to count as reward")
	}
}

func TestOrderQtyInRewardUomSumsMatchingProducts(t *testing.T) {
	smallest := model.UomTypeSmallest
	details := []entity.ConPromoV2Det{
		{ProID: 10836, Qty1: 3, Qty2: 0, Qty3: 0},
		{ProID: 10836, Qty1: 1, Qty2: 0, Qty3: 0},
		{ProID: 9999, Qty1: 5, Qty2: 0, Qty3: 0},
	}
	got := orderQtyInRewardUom(details, 10836, &smallest)
	if got != 4 {
		t.Fatalf("expected order qty 4 smallest, got %v", got)
	}
}

func TestResolveConsultPromoStockErrorWithoutRequestedPromoIDs(t *testing.T) {
	err := resolveConsultPromoStockError(nil, []string{"STRATA003-003-003"}, nil)
	if err == nil {
		t.Fatal("expected insufficient stock error")
	}
	want := errmsg.PromoInsufficientStockMessage("STRATA003-003-003")
	if err.Error() != want {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}

func TestResolveConsultPromoStockErrorKeepsOtherPromosWhenSomeSucceed(t *testing.T) {
	err := resolveConsultPromoStockError(nil, []string{"PROMO-A"}, []entity.ConsultPromoResp{{PromoID: "PROMO-B"}})
	if err != nil {
		t.Fatalf("unexpected error when another promo succeeded: %v", err)
	}
}

func TestResolveConsultPromoStockErrorForRequestedPromoID(t *testing.T) {
	err := resolveConsultPromoStockError([]string{"STRATA003-003-003"}, []string{"STRATA003-003-003"}, []entity.ConsultPromoResp{{PromoID: "OTHER"}})
	if err == nil {
		t.Fatal("expected error for requested promo with insufficient stock")
	}
	if err.Error() != errmsg.PromoInsufficientStockMessage("STRATA003-003-003") {
		t.Fatalf("unexpected message: %s", err.Error())
	}
}

func TestPromoIDInList(t *testing.T) {
	if !promoIDInList("A", []string{"A", "B"}) {
		t.Fatal("expected promo A to be found")
	}
	if promoIDInList("C", []string{"A", "B"}) {
		t.Fatal("expected promo C to be missing")
	}
}

func TestNormalizePromoIDListAllowsHyphen(t *testing.T) {
	got, err := entity.NormalizePromoIDList([]string{"STRATA003-003-003"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "STRATA003-003-003" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestCalculateStrataRuleValueUsesFullOrderGrossForValueRule(t *testing.T) {
	strata := model.PromotionV2Strata{RuleType: model.RuleTypeValue}
	criteriaDetails := map[int]*entity.ConPromoV2Det{
		10836: {ProID: 10836, GrossValue: 400000},
	}
	orderDetails := []entity.ConPromoV2Det{
		{ProID: 10807, GrossValue: 3600000},
		{ProID: 8436, GrossValue: 10000000},
		{ProID: 10836, GrossValue: 400000},
	}
	got := calculateStrataRuleValue(strata, criteriaDetails, orderDetails)
	want := 14000000.0
	if got != want {
		t.Fatalf("expected full order gross %v, got %v", want, got)
	}
}

func TestResolveNonSequentialStrataOrdinalCapsAtHighestTier(t *testing.T) {
	strata := []model.PromotionV2Strata{
		{Ordinal: 1, RangeFrom: 1_000_000, RangeTo: 2_000_000},
		{Ordinal: 2, RangeFrom: 2_000_001, RangeTo: 3_000_000},
		{Ordinal: 3, RangeFrom: 3_000_001, RangeTo: 9_999_999},
	}
	got := resolveNonSequentialStrataOrdinal(strata, 14_000_000)
	if got != 3 {
		t.Fatalf("expected ordinal 3 for value above top tier, got %d", got)
	}
}

func TestResolveNonSequentialStrataOrdinalUsesInRangeTier(t *testing.T) {
	strata := []model.PromotionV2Strata{
		{Ordinal: 1, RangeFrom: 1_000_000, RangeTo: 2_000_000},
		{Ordinal: 2, RangeFrom: 2_000_001, RangeTo: 3_000_000},
		{Ordinal: 3, RangeFrom: 3_000_001, RangeTo: 9_999_999},
	}
	got := resolveNonSequentialStrataOrdinal(strata, 5_000_000)
	if got != 3 {
		t.Fatalf("expected ordinal 3, got %d", got)
	}
}
