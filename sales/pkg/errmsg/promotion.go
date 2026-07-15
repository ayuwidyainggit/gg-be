package errmsg

const (
	ERROR_PROMO_ID_ALREADY_EXISTS = "The promo_id already exists"
)

func PromoInsufficientStockMessage(promoID string) string {
	return "Promo " + promoID + " tidak tersedia karena stok tidak mencukupi."
}
