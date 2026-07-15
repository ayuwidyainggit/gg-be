package entity

// ExpenseQueryFilter represents query parameters for expense type list
type ExpenseQueryFilter struct {
	Q            string `query:"q"`
	Page         int    `query:"page" validate:"required"`
	Limit        int    `query:"limit" validate:"required"`
	Sort         string `query:"sort"`
	CustId       string
	ParentCustId string
}

// ExpenseListResponse represents expense type data in list response
type ExpenseListResponse struct {
	ExpenseTypeID   int    `json:"expense_type_id"`
	ExpenseTypeCode string `json:"expense_type_code"`
	ExpenseTypeName string `json:"expense_type_name"`
	Status          bool   `json:"status"`
	UpdateBy        int64  `json:"update_by"`
	UpdateByName    string `json:"update_by_name"`
	UpdateDate      string `json:"update_date"`
}

// CreateExpenseBody represents request body for creating expense type
type CreateExpenseBody struct {
	ExpenseTypeCode string `json:"expense_type_code" validate:"required,max=20"`
	ExpenseTypeName string `json:"expense_type_name" validate:"required,max=50"`
	IsActive        bool   `json:"is_active"`
	CustId          string `json:"-"`
	ParentCustId    string `json:"-"`
}

// UpdateExpenseBody represents request body for updating expense type
type UpdateExpenseBody struct {
	ExpenseTypeCode string `json:"expense_type_code" validate:"required,max=20"`
	ExpenseTypeName string `json:"expense_type_name" validate:"required,max=50"`
	IsActive        *bool  `json:"is_active"`
}

// DetailExpenseParams represents path parameters for expense type detail
type DetailExpenseParams struct {
	ExpenseTypeID int `params:"expense_type_id" validate:"required"`
}

// UpdateExpenseParams represents path parameters for expense type update
type UpdateExpenseParams struct {
	ExpenseTypeID int `params:"expense_type_id" validate:"required"`
}

// DeleteExpenseParams represents path parameters for expense type delete
type DeleteExpenseParams struct {
	ExpenseTypeID int `params:"expense_type_id" validate:"required"`
}

type DepositExpense struct {
	ExpenseID   int64   `json:"expense_id,omitempty"`
	DocNo       string  `json:"doc_no,omitempty"`
	ExpenseName string  `json:"expense_name,omitempty"`
	Amount      float64 `json:"amount,omitempty"`

	// Deposit detail required fields.
	DepositExpenseID int64   `json:"deposit_expense_id,omitempty"`
	DocumentNo       string  `json:"document_no,omitempty"`
	Balance          float64 `json:"balance"`
	PaymentAmount    float64 `json:"payment_amount"`
}
