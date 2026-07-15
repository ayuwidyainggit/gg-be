package entity

// PromoRewardType represents the type of reward calculation for promotions.
// Used in promo_criteria table.
type PromoRewardType int

const (
	// PromoRewardTypeQuantity represents a quantity-based reward (free product)
	// Reward is in product quantity
	PromoRewardTypeQuantity PromoRewardType = 1

	// PromoRewardTypeFixedValue represents a fixed value reward
	// Reward is amount in currency
	PromoRewardTypeFixedValue PromoRewardType = 2

	// PromoRewardTypePercent represents percentage reward
	// Calculated as percentage of subtotal
	PromoRewardTypePercent PromoRewardType = 3
)

// IsQuantity returns true if the reward type is quantity-based
func (p PromoRewardType) IsQuantity() bool {
	return p == PromoRewardTypeQuantity
}

// IsFixedValue returns true if the reward type is a fixed value
func (p PromoRewardType) IsFixedValue() bool {
	return p == PromoRewardTypeFixedValue
}

// IsPercent returns true if the reward type is percentage-based
func (p PromoRewardType) IsPercent() bool {
	return p == PromoRewardTypePercent
}

// String returns the display name of the PromoRewardType
func (p PromoRewardType) String() string {
	switch p {
	case PromoRewardTypeQuantity:
		return "Quantity"
	case PromoRewardTypeFixedValue:
		return "Fixed Value"
	case PromoRewardTypePercent:
		return "Percentage"
	default:
		return ""
	}
}

// DisplayName returns the display name (alias for String)
func (p PromoRewardType) DisplayName() string {
	return p.String()
}
