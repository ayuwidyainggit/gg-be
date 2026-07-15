package entity

// Request Entities

type CreateStockDisposalBody struct {
	CustID       string                           `json:"cust_id" validate:"required"`
	ParentCustID string                           `json:"parent_cust_id"`
	CreatedBy    int64                            `json:"created_by"`
	Date         string                           `json:"date" validate:"required"`
	SupID        int64                            `json:"sup_id" validate:"required"`
	WhID         int64                            `json:"wh_id" validate:"required"`
	GrNo         *string                          `json:"gr_no"`
	Note         string                           `json:"note"`
	Products     []CreateStockDisposalProductBody `json:"products" validate:"required,min=1"`
}

type CreateStockDisposalProductBody struct {
	ProID       int64                        `json:"pro_id" validate:"required"`
	UnitID1     string                       `json:"unit_id1" validate:"required"`
	UnitID2     string                       `json:"unit_id2" validate:"required"`
	UnitID3     string                       `json:"unit_id3" validate:"required"`
	Qty1        float64                      `json:"qty1" validate:"required,min=0"`
	Qty2        float64                      `json:"qty2" validate:"min=0"`
	Qty3        float64                      `json:"qty3" validate:"min=0"`
	PurchPrice1 float64                      `json:"purch_price1" validate:"required"`
	PurchPrice2 float64                      `json:"purch_price2" validate:"required"`
	PurchPrice3 float64                      `json:"purch_price3" validate:"required"`
	GrossPrice  float64                      `json:"gross_price" validate:"required"`
	Vat         float64                      `json:"vat" validate:"required"`
	VatValue    float64                      `json:"vat_value" validate:"required"`
	SubTotal    float64                      `json:"sub_total" validate:"required"`
	UploadFile  *CreateStockDisposalFileBody `json:"upload_file"`
}

type CreateStockDisposalFileBody struct {
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	MediaCategory string `json:"media_category" validate:"omitempty,oneof=image video"`
	FileUrl       string `json:"file_url"`
	FileSize      int64  `json:"file_size"`
}

type StockDisposalQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64  `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64  `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int     `query:"page"`
	Limit        int     `query:"limit" validate:"required"`
	Query        string  `query:"q"`
	Sort         string  `query:"sort"`
	TrCode       string  `query:"tr_code"`
	StockType    string  `query:"stock_type"`
	WhID         []int64 `query:"wh_id"`
	SupID        []int64 `query:"sup_id"`
}

type DetailStockDisposalParams struct {
	StockDisposalID int64 `params:"stock_disposal_id" validate:"required"`
}

type StockDisposalProductLookupQueryFilter struct {
	WhID          int64  `query:"wh_id"`
	DistributorID *int64 `query:"distributor_id"`
	SupID         *int64 `query:"sup_id"`
	Query         string `query:"q"`
	ZeroStock     *bool  `query:"zero_stock"`
	Page          int    `query:"page"`
	Limit         int    `query:"limit" validate:"required,min=1"`
}

type StockDisposalProductLookupResponse struct {
	ProID           int64   `json:"pro_id"`
	ProCode         string  `json:"pro_code"`
	ProName         string  `json:"pro_name"`
	Vat             int     `json:"vat"`
	ConvUnit2       int     `json:"conv_unit2"`
	ConvUnit3       int     `json:"conv_unit3"`
	UnitID1         string  `json:"unit_id1"`
	UnitID2         string  `json:"unit_id2"`
	UnitID3         string  `json:"unit_id3"`
	PurchPrice1     float64 `json:"purch_price1"`
	PurchPrice2     float64 `json:"purch_price2"`
	PurchPrice3     float64 `json:"purch_price3"`
	MinStockQty     float64 `json:"min_stock_qty"`
	SafStockQty     float64 `json:"saf_stock_qty"`
	Qty1            float64 `json:"qty1"`
	Qty2            float64 `json:"qty2"`
	Qty3            float64 `json:"qty3"`
	TotalQty        float64 `json:"total_qty"`
	InTransitStock1 float64 `json:"in_transit_stock1"`
	InTransitStock2 float64 `json:"in_transit_stock2"`
	InTransitStock3 float64 `json:"in_transit_stock3"`
}

// Response Entities

type StockDisposalResponse struct {
	SdNumber      string                     `json:"sd_number"`
	DisposalDate  string                     `json:"disposal_date"`
	TrCode        string                     `json:"tr_code"`
	SupID         *int64                     `json:"sup_id"`
	SupCode       string                     `json:"sup_code"`
	SupName       string                     `json:"sup_name"`
	WhID          int64                      `json:"wh_id"`
	WhCode        string                     `json:"wh_code"`
	WhName        string                     `json:"wh_name"`
	StockType     string                     `json:"stock_type"`
	StockTypeDesc string                     `json:"stock_type_desc"`
	GrNo          *string                    `json:"gr_no"`
	Note          *string                    `json:"note"`
	SubTotal      float64                    `json:"sub_total"`
	VatValue      float64                    `json:"vat_value"`
	Total         float64                    `json:"total"`
	CreatedBy     int64                      `json:"created_by"`
	CreatedByName string                     `json:"created_by_name"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedBy     *int64                     `json:"updated_by"`
	UpdatedByName string                     `json:"updated_by_name"`
	UpdatedAt     string                     `json:"updated_at"`
	Details       []StockDisposalDetResponse `json:"details"`
}

type StockDisposalDetResponse struct {
	SdDetailID    int64   `json:"sd_detail_id"`
	ProID         int64   `json:"pro_id"`
	ProCode       string  `json:"pro_code"`
	ProName       string  `json:"pro_name"`
	UnitID1       string  `json:"unit_id1"`
	UnitID2       string  `json:"unit_id2"`
	UnitID3       string  `json:"unit_id3"`
	Qty1          float64 `json:"qty1"`
	Qty2          float64 `json:"qty2"`
	Qty3          float64 `json:"qty3"`
	PurchPrice1   float64 `json:"purch_price1"`
	PurchPrice2   float64 `json:"purch_price2"`
	PurchPrice3   float64 `json:"purch_price3"`
	GrossPrice    float64 `json:"gross_price"`
	Vat           float64 `json:"vat"`
	VatValue      float64 `json:"vat_value"`
	SubTotal      float64 `json:"sub_total"`
	FileName      *string `json:"file_name"`
	FileType      *string `json:"file_type"`
	MediaCategory *string `json:"media_category"`
	FileBase64    *string `json:"file_base64"`
	FileSize      *int64  `json:"file_size"`
}

type StockDisposalListResponse struct {
	SdID     int64   `json:"sd_id"`
	Date     string  `json:"date"`
	SdNumber string  `json:"sd_number"`
	WhID     int64   `json:"wh_id"`
	WhCode   string  `json:"wh_code"`
	WhName   string  `json:"wh_name"`
	SupID    int64   `json:"sup_id"`
	SupCode  string  `json:"sup_code"`
	SupName  string  `json:"sup_name"`
	Subtotal float64 `json:"subtotal"`
	Vat      float64 `json:"vat"`
	VatValue float64 `json:"vat_value"`
	Total    float64 `json:"total"`
}

type StockDisposalDetailResponse struct {
	SdDetailID   int64                          `json:"sd_detail_id"`
	SdID         int64                          `json:"sd_id"`
	Date         string                         `json:"date"`
	SdNumber     string                         `json:"sd_number"`
	SupID        int64                          `json:"sup_id"`
	SupCode      string                         `json:"sup_code"`
	SupName      string                         `json:"sup_name"`
	WhID         int64                          `json:"wh_id"`
	WhCode       string                         `json:"wh_code"`
	WhName       string                         `json:"wh_name"`
	StockType    string                         `json:"stock_type"`
	GrNo         *string                        `json:"gr_no"`
	Note         string                         `json:"note"`
	DataProducts []StockDisposalProductResponse `json:"data_products"`
	Subtotal     float64                        `json:"subtotal"`
	Vat          float64                        `json:"vat"`
	Total        float64                        `json:"total"`
}

type StockDisposalProductResponse struct {
	ProID       int64                      `json:"pro_id"`
	ProCode     string                     `json:"pro_code"`
	ProName     string                     `json:"pro_name"`
	UnitID1     string                     `json:"unit_id1"`
	UnitID2     string                     `json:"unit_id2"`
	UnitID3     string                     `json:"unit_id3"`
	Qty1        float64                    `json:"qty1"`
	Qty2        float64                    `json:"qty2"`
	Qty3        float64                    `json:"qty3"`
	PurchPrice1 float64                    `json:"purch_price1"`
	PurchPrice2 float64                    `json:"purch_price2"`
	PurchPrice3 float64                    `json:"purch_price3"`
	GrossPrice  float64                    `json:"gross_price"`
	Vat         float64                    `json:"vat"`
	VatValue    float64                    `json:"vat_value"`
	SubTotal    float64                    `json:"sub_total"`
	UploadFile  *StockDisposalFileResponse `json:"upload_file"`
}

type StockDisposalFileResponse struct {
	FileName      string `json:"file_name"`
	FileType      string `json:"file_type"`
	MediaCategory string `json:"media_category"`
	FileUrl       string `json:"file_url"`
	FileSize      string `json:"file_size"`
}
