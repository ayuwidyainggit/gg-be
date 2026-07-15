package request

type DestinationDetailsRequest struct {
	PjpID     int
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	SortOrder string `form:"sort_order"`
	Date      string `form:"date" binding:"required"`
}
