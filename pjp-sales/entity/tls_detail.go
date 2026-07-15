package entity

type CreateTlsDetBody struct {
	TlsId    int64   `json:"tls_id"`
	OutletId int64   `json:"outlet_id"`
	Notes    *string `json:"notes"`
}

type TlsDetResponse struct {
	TlsId      int     `json:"tls_id"`
	TlsDetId   int64   `json:"tls_det_id"`
	OutletId   int64   `json:"outlet_id"`
	OutletCode string  `json:"outlet_code"`
	OutletName string  `json:"outlet_name"`
	Notes      *string `json:"notes"`
}

type UpdateTlsDetBody struct {
	TlsDetId *int64  `json:"tls_det_id"`
	OutletId int64   `json:"outlet_id"`
	Notes    *string `json:"notes"`
}
