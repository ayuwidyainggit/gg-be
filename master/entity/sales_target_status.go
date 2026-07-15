package entity

type SalesTargetStatus int64

const (
	SALES_TARGET_STATUS_DRAFT    SalesTargetStatus = 0
	SALES_TARGET_STATUS_ACTIVE   SalesTargetStatus = 1
	SALES_TARGET_STATUS_INACTIVE SalesTargetStatus = 2
)

func (s SalesTargetStatus) String() string {
	switch s {
	case SALES_TARGET_STATUS_DRAFT:
		return "Draft"
	case SALES_TARGET_STATUS_ACTIVE:
		return "Active"
	case SALES_TARGET_STATUS_INACTIVE:
		return "Inactive"
	default:
		return "Unknown"
	}
}
