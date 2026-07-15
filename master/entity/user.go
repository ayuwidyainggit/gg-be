package entity

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type TokenMetadata struct {
	UserId         int64
	DistributorId  int64
	UserName       string
	UserFullName   string
	Email          string
	IsAdmin        bool
	CustId         string
	ParentCustId   string
	DistPriceGrpId int
	EmpId          int64
	LangId         string
	MobileNo       string
	Whatsapp       string
	Expires        int64
}
