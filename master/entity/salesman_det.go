package entity

type SalesmanDetGroup struct {
	ProductLine []SalesmanDetGroupProductLine `json:"product_line"`
	Brand       []SalesmanDetGroupBrand       `json:"brand"`
	SubBrand    []SalesmanDetGroupSubBrand    `json:"sub_brand"`
}

type SalesmanDetGroupProductLine struct {
	PlId                   int    `json:"pl_id"`
	MSalesmanProductTypeID int64  `json:"m_salesman_product_type_id"`
	RefID                  int64  `json:"ref_id"`
	PlCode                 string `json:"pl_code"`
	PlName                 string `json:"pl_name"`
}

type SalesmanDetGroupBrand struct {
	PlId                   int    `json:"pl_id"`
	MSalesmanProductTypeID int64  `json:"m_salesman_product_type_id"`
	RefID                  int64  `json:"ref_id"`
	BrandCode              string `json:"brand_code"`
	BrandName              string `json:"brand_name"`
}

type SalesmanDetGroupSubBrand struct {
	PlId                   int    `json:"pl_id"`
	MSalesmanProductTypeID int64  `json:"m_salesman_product_type_id"`
	RefID                  int64  `json:"ref_id"`
	SBrand1Code            string `json:"sbrand1_code"`
	SBrand1Name            string `json:"sbrand1_name"`
}

type SalesmanDetCreateDetGroup struct {
	ProductLine []SalesmanDetGroupProductLine `json:"product_line"`
	Brand       []SalesmanDetGroupBrand       `json:"brand"`
	SubBrand    []SalesmanDetGroupSubBrand    `json:"sub_brand"`
}

type SalesmanDetCreateProductLine struct {
	RefID  int64  `json:"ref_id"`
	PlCode string `json:"pl_code"`
	PlName string `json:"pl_name"`
}

type SalesmanDetCreateGroupBrand struct {
	RefID     int64  `json:"ref_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type SalesmanDetCreateGroupSubBrand struct {
	RefID       int64  `json:"ref_id"`
	SBrand1Code string `json:"sbrand1_code"`
	SBrand1Name string `json:"sbrand1_name"`
}
type SalesmanDetGroupUpdate struct {
	ProductLine []SalesmanDetGroupProductLineUpdate `json:"product_line"`
	Brand       []SalesmanDetGroupBrandUpdate       `json:"brand"`
	SubBrand    []SalesmanDetGroupSubBrandUpdate    `json:"sub_brand"`
}

type SalesmanDetGroupProductLineUpdate struct {
	MSalesmanProductTypeID *int64 `json:"m_salesman_product_type_id"`
	RefID                  *int64 `json:"ref_id"`
}

type SalesmanDetGroupBrandUpdate struct {
	MSalesmanProductTypeID *int64 `json:"m_salesman_product_type_id"`
	RefID                  *int64 `json:"ref_id"`
}

type SalesmanDetGroupSubBrandUpdate struct {
	MSalesmanProductTypeID *int64 `json:"m_salesman_product_type_id"`
	RefID                  *int64 `json:"ref_id"`
}
