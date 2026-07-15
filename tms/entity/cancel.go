package entity

type CancelRequest struct {
	ProductID []int `validate:"required,notEmptyIntSlice" json:"product_id"`
}
