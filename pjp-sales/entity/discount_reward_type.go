package entity

// DiscountRewardType represents the type of reward calculation for discounts.
// Used in discount_criteria table.
type DiscountRewardType int

const (
	// DiscountRewardTypeValue represents a fixed value reward (amount in currency)
	// Direct amount deduction
	DiscountRewardTypeValue DiscountRewardType = 1

	// DiscountRewardTypePercentage represents a percentage-based reward
	// Calculated as (SlabReward * SubTotal) / 100
	DiscountRewardTypePercentage DiscountRewardType = 2
)

// IsPercentage returns true if the reward type is percentage-based
func (d DiscountRewardType) IsPercentage() bool {
	return d == DiscountRewardTypePercentage
}

// IsValue returns true if the reward type is a fixed value
func (d DiscountRewardType) IsValue() bool {
	return d == DiscountRewardTypeValue
}

// String returns the display name of the DiscountRewardType
func (d DiscountRewardType) String() string {
	switch d {
	case DiscountRewardTypeValue:
		return "Value"
	case DiscountRewardTypePercentage:
		return "Percentage"
	default:
		return ""
	}
}

// DisplayName returns the display name (alias for String)
func (d DiscountRewardType) DisplayName() string {
	return d.String()
}
