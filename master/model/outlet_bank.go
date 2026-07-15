package model

type MOutletBank struct {
	CustID       string `db:"cust_id" json:"cust_id"`
	OutletID     int64  `db:"outlet_id" json:"outlet_id"`
	BankId       int64  `db:"bank_id" json:"bank_id"`
	AccountNo    string `db:"account_no" json:"account_no"`
	AccountName  string `db:"account_name" json:"account_name"`
	OutletBankId *int64 `db:"outlet_bank_id" json:"outlet_bank_id"`
}

type OutletBankList struct {
	BankId       int64  `db:"bank_id" json:"bank_id"`
	AccountNo    string `db:"account_no" json:"account_no"`
	AccountName  string `db:"account_name" json:"account_name"`
	OutletBankId *int64 `db:"outlet_bank_id" json:"outlet_bank_id"`
}

type MOutletBankUpdate struct {
	BankId      *int64  `db:"bank_id" json:"bank_id"`
	AccountNo   *string `db:"account_no" json:"account_no"`
	AccountName *string `db:"account_name" json:"account_name"`
}

type MOutletBankRead struct {
	CustID       string  `db:"cust_id" json:"cust_id"`
	OutletID     int64   `db:"outlet_id" json:"outlet_id"`
	BankId       int64   `db:"bank_id" json:"bank_id"`
	BankCode     *string `db:"bank_code" json:"bank_code"`
	BankName     *string `db:"bank_name" json:"bank_name"`
	AccountNo    *string `db:"account_no" json:"account_no"`
	AccountName  *string `db:"account_name" json:"account_name"`
	OutletBankId *int64  `db:"outlet_bank_id" json:"outlet_bank_id"`
}
