package entity

type Product struct {
	Id int64 `json:"pro_id" validate:"required"`
}

type ProductId struct {
	Product Product
}

type DetailsNormalProductId struct {
	ProductIds []ProductId `json:"details.normal[].pro_id" validate:"unique,dive"`
}

type DetailsPromoProductId struct {
	ProductIds []ProductId `json:"details.promo[].pro_id" validate:"unique,dive"`
}

type DetailsProductId struct {
	ProductIds []ProductId `json:"details.pro_id" validate:"unique,dive"`
}
