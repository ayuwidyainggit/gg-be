package entity

type OperationType string

const (
	OperationTypeTakingOrder = "O"
	OperationTypeCanvas      = "C"
)

func (op OperationType) IsCanvas() bool {
	return op == OperationTypeCanvas
}

func (op OperationType) IsTakingOrder() bool {
	return op == OperationTypeTakingOrder
}

func (op OperationType) String() string {
	return string(op)
}

func (op OperationType) IsValid() bool {
	switch op {
	case OperationTypeTakingOrder, OperationTypeCanvas:
		return true
	default:
		return false
	}
}
