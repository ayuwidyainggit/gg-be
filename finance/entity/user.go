package entity

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type TokenMetadata struct {
	UserId         int64
	UserName       string
	UserFullName   string
	Email          string
	IsAdmin        bool
	CustId         string
	EmpId          int64
	LangId         string
	MobileNo       string
	Whatsapp       string
	Expires        int64
	ParentCustId   string
	DistPriceGrpId int
}
