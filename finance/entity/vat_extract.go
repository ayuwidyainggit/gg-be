package entity

import (
	"strconv"
	"time"
	"unicode"
)

const (
	VAT_TYPE_IN  = 1
	VAT_TYPE_OUT = 2
)

type VatExtractListResponse struct {
	TransactionID    uint    `json:"transaction_id"`
	InvoiceNo        string  `json:"invoice_no"`
	InvoiceDate      string  `json:"invoice_date"`
	InvoiceType      string  `json:"invoice_type"`
	NPWP             string  `json:"npwp"`
	SupplierCode     string  `json:"supplier_code"`
	SupplierName     string  `json:"supplier_name"`
	Address          string  `json:"address"`
	DPP              float64 `json:"dpp"`
	PPN              float64 `json:"ppn"`
	PPNBM            float64 `json:"ppnbm"`
	TaxNo            string  `json:"tax_no"`
	TaxDate          string  `json:"tax_date"`
	TaxExtractDate   string  `json:"tax_extract_date"`
	ReturnDocumentNo string  `json:"return_no,omitempty"`
	ReturnDate       string  `json:"return_date,omitempty"`
}
type VatExtractParams struct {
	VatExtractID int64 `params:"vat_extract_id" validate:"required"`
}
type VatExtractQueryFilter struct {
	TransactionID    []uint `query:"transaction_id"`
	VatType          int    `query:"vat_type" validate:"required,oneof=1 2"`
	InvoiceType      string `query:"invoice_type" validate:"required,oneof=I R"`
	ExtractionStatus string `query:"extraction_status"` // E, NE
	CustId           string
	ParentCustId     string
	From             *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To               *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page             int    `query:"page"`
	Limit            int    `query:"limit"`
	Query            string `query:"q"`
	Mode             string `query:"mode"`
	Sort             string `query:"sort"`
	IsActive         *int   `query:"is_active"`
	Taxes            bool   `query:"taxes"`
}
type VatExtractReq struct {
	TransactionID []int64 `json:"transaction_id"`
	VatType       int     `json:"vat_type" validate:"required,oneof=1 2"`
	InvoiceType   string  `json:"invoice_type" validate:"required,oneof=I R"`
	CustID        string
	CreatedBy     int64 `json:"created_by"`
}

type VatExtractResp struct {
	ID int64 `json:"id"`
}

type VatExtractResultQueryFilter struct {
	CustId       string
	ParentCustId string
	From         *int64 `query:"from" validate:"required_with=To,omitempty,gte=1000000000"`
	To           *int64 `query:"to" validate:"required_with=From,omitempty,lte=9999999999,gtefield=From"`
	Page         int    `query:"page"`
	Limit        int    `query:"limit" validate:""`
	Sort         string `query:"sort"`
	VatTypes     []int  `query:"vat_type"`
}
type VatExtractList struct {
	VatExtractID   *int64     `json:"vat_extract_id"`
	VatExtractType int        `json:"vat_extract_type"`
	InvoiceType    string     `json:"invoice_type"`
	ExtractTotal   int        `json:"extract_total"`
	CreatedBy      int64      `json:"created_by"`
	CreatedAt      *time.Time `json:"created_at"`
}

type VatExtractDowloadVatInInvoice struct {
	KdJenisTransaksi string  `json:"kd_jenis_transaksi"`
	FgPengganti      string  `json:"fg_pengganti"`
	NomorFaktur      string  `json:"nomor_faktur"`
	MasaPajak        int     `json:"masa_pajak"`
	TahunPajak       int     `json:"tahun_pajak"`
	TanggalFaktur    string  `json:"tanggal_faktur"`
	NPWP             string  `json:"npwp"`
	Nama             string  `json:"nama"`
	AlamatLengkap    string  `json:"alamat_lengkap"`
	JumlahDPP        float64 `json:"jumlah_dpp"`
	JumlahPPN        float64 `json:"jumlah_ppn"`
	JumlahPPNBM      float64 `json:"jumlah_ppnbm"`
	IsCreditable     int     `json:"is_creditable"`
}

type VatExtractDowloadVatInReturn struct {
	Npwp              string  `json:"npwp"`
	Nama              string  `json:"nama"`
	KdJenisTransaksi  string  `json:"kd_jenis_transaksi"`
	FgPengganti       string  `json:"fg_pengganti"`
	NomorFaktur       string  `json:"nomor_faktur"`
	TanggalFaktur     string  `json:"tanggal_faktur"`
	IsCreditable      int     `json:"is_creditable"`
	NomorDokumenRetur string  `json:"nomor_dokumen_retur"`
	TanggalRetur      string  `json:"tanggal_retur"`
	MasaPajakRetur    int     `json:"masa_pajak_retur"`
	TahunPajak        int     `json:"tahun_pajak"`
	NilaiReturDPP     float64 `json:"nilai_retur_dpp"`
	NilaiReturPPN     float64 `json:"nilai_retur_ppn"`
	NilaiReturPPNBM   float64 `json:"nilai_retur_ppnbm"`
}

type Faktur struct {
	KdJjenisTansaksi string
	FgPengganti      string
	NomorFaktur      string
	TahunPajak       int
}

func GenerateFakturComponents(input string) Faktur {
	kdJenisTransaksi := input[:2]

	fgPengganti := input[2:3]

	remaining := input[3:]
	nomorFaktur := ""
	for _, char := range remaining {
		if unicode.IsDigit(char) {
			nomorFaktur += string(char)
		}
	}

	tahunStr := input[8:10] // Karakter ke-9 dan ke-10 (index 8-9)

	// Mengonversi angka tahun menjadi integer
	tahunInt, _ := strconv.Atoi(tahunStr)

	// Menambahkan 2000 untuk mendapatkan tahun penuh
	tahun := 2000 + tahunInt

	return Faktur{
		KdJjenisTansaksi: kdJenisTransaksi,
		FgPengganti:      fgPengganti,
		NomorFaktur:      nomorFaktur,
		TahunPajak:       tahun,
	}
}
