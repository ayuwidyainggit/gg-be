package entity

type PaymentType int

const (
	PayTypeCash PaymentType = iota + 1
	PayTypeTransfer
	PayTypeGiro
	PayTypeReturn
	PayTypeCredit
)

func (c PaymentType) Value() int16 {
	return int16(c)
}

func (c PaymentType) IsValid() bool {
	switch c {
	case PayTypeCash, PayTypeTransfer, PayTypeGiro, PayTypeReturn, PayTypeCredit:
		return true
	default:
		return false
	}
}

func (c PaymentType) IsCash() bool {
	return c == PayTypeCash
}

func (c PaymentType) IsTransfer() bool {
	return c == PayTypeTransfer
}

func (c PaymentType) IsGiro() bool {
	return c == PayTypeGiro
}

func (c PaymentType) IsCredit() bool {
	return c == PayTypeCredit
}

func (c PaymentType) String() string {
	switch c {
	case PayTypeCash:
		return "Cash"
	case PayTypeTransfer:
		return "Transfer"
	case PayTypeGiro:
		return "Cheque/Bilyet Giro"
	case PayTypeReturn:
		return "Return"
	case PayTypeCredit:
		return "Credit/Debit"
	default:
		return ""
	}
}
