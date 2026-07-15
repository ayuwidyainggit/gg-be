package entity

type ArSettlementQueryFilter struct {
	CustId          string
	ParentCustId    string
	EmpId           []int   `query:"emp_id"`
	DepositStatus   []int   `query:"deposit_status"`
	DepositDateFrom *int64  `query:"deposit_date_from" validate:"required_with=DepositDateTo,omitempty,gte=1000000000"`
	DepositDateTo   *int64  `query:"deposit_date_to" validate:"required_with=DepositDateFrom,omitempty,lte=9999999999,gtefield=DepositDateFrom"`
	Page            int     `query:"page"`
	Limit           int     `query:"limit" validate:"required"`
	Query           string  `query:"q"`
	Mode            string  `query:"mode"`
	Sort            string  `query:"sort"`
	IsActive        *int    `query:"is_active"`
	DistributorId   *string `query:"cust_id"`
}

type ArBranchSettlementQueryFilter struct {
	CustId *string `query:"cust_id"`
}

type ArSettlementResponse struct {
	CustId              *string                       `json:"cust_id"`
	CustName            *string                       `json:"cust_name"`
	DepositNo           *string                       `json:"deposit_no"`
	DepositDate         *string                       `json:"deposit_date"`
	CollectionNo        *string                       `json:"collection_no"`
	CollectionDate      *string                       `json:"collection_date"`
	OtGrpID             *int64                        `json:"ot_grp_id"`
	OtGrpCode           *string                       `json:"ot_grp_code"`
	OtGrpName           *string                       `json:"ot_grp_name"`
	EmpId               *int64                        `json:"emp_id"`
	EmpCode             *string                       `json:"emp_code"`
	EmpName             *string                       `json:"emp_name"`
	TotalDiscount       *float64                      `json:"total_discount"`
	TotalMaterai        *float64                      `json:"total_materai"`
	ExpenseTotal        *float64                      `json:"expense_total"`
	TotalPaymentBalance *float64                      `json:"total_payment_balance"`
	TotalPaymnet        *float64                      `json:"total_payment"`
	RemainingAmount     *float64                      `json:"remaining_amount"`
	DepositStatus       *int64                        `json:"deposit_status"`
	DepositStatusName   *string                       `json:"deposit_status_name"`
	Details             []ArSettlementPaymentResponse `json:"details"`
	Expense             []ArSettlementExpenseItem     `json:"expense"`
}

// ArSettlementExpenseItem is one expense item in AR Settlement Detail (from acf.deposit_expense + acf.expense).
type ArSettlementExpenseItem struct {
	DepositExpenseID int64   `json:"deposit_expense_id"`
	DocNo            string  `json:"doc_no"`
	Balance          float64 `json:"balance"`
	PaymentAmount    float64 `json:"payment_amount"`
}

type ArSettlementPaymentResponse struct {
	DepositPaymentID int64    `json:"deposit_payment_id"`
	InvoiceNo        *string  `json:"invoice_no"`
	InvoiceDate      *string  `json:"invoice_date"`
	PayType          *int64   `json:"pay_type"`
	PayTypeName      *string  `json:"pay_type_name"`
	DocumentNo       *string  `json:"document_no"`
	Balance          *float64 `json:"balance"`
	PaymentAmount    *float64 `json:"payment_amount"`
	SalesmanId       *int64   `json:"salesman_id"`
	SalesmanCode     *string  `json:"salesman_code"`
	SalesmanName     *string  `json:"salesman_name"`
	OutletId         *int64   `json:"outlet_id"`
	OutletCode       *string  `json:"outlet_code"`
	OutletName       *string  `json:"outlet_name"`
	Discount         *float64 `json:"discount"`
	Materai          *float64 `json:"materai"`
	PaymentBalance   *float64 `json:"payment_balance"`
	TotalPayment     *float64 `json:"total_payment"`
	RemainingPayment *float64 `json:"remaining_payment"`
}

type ArSettlementListResponse struct {
	CustId      *string `json:"cust_id"`
	DepositNo   *string `json:"deposit_no"`
	DepositDate *string `json:"deposit_date"`
	EmpId       *int64  `json:"emp_id"`
	EmpCode     *string `json:"emp_code"`
	EmpName     *string `json:"emp_name"`
	// InvoiceAmount     *float64 `json:"invoice_amount"`
	TotalPayment      *float64 `json:"total_payment"`
	RemainingAmount   *float64 `json:"remaining_amount"`
	DepositStatus     *int64   `json:"deposit_status"`
	DepositStatusName *string  `json:"deposit_status_name"`
	ApprovedBy        *int64   `json:"approved_by"`
	ApprovedByName    *string  `json:"approved_by_name"`
}

type DetailArSettlementParams struct {
	DepositNo string `params:"deposit_no" validate:"required" json:"deposit_no"`
}

type ApproveArSettlementParams struct {
	DepositNo string `params:"deposit_no" validate:"required" json:"deposit_no"`
}

type RejectArSettlementParams struct {
	DepositNo string `params:"deposit_no" validate:"required" json:"deposit_no"`
}

// BulkApproveArSettlementItem is one item in the bulk approve request body.
type BulkApproveArSettlementItem struct {
	DepositNo string `json:"deposit_no" validate:"required"`
	CustId    string `json:"cust_id" validate:"required"`
}

const (
	DEPOSIT_STATUS_IN_REVIEW   = 1
	DEPOSIT_STATUS_IN_APPROVED = 2
	DEPOSIT_STATUS_IN_REJECTED = 3
)

var dataDepositStatusName = map[int64]string{
	DEPOSIT_STATUS_IN_REVIEW:   "In Review",
	DEPOSIT_STATUS_IN_APPROVED: "Approved",
	DEPOSIT_STATUS_IN_REJECTED: "Rejected",
}

func (arSettlement ArSettlementListResponse) GenerateDataDepositStatusName() string {
	if arSettlement.DepositStatus != nil {
		return dataDepositStatusName[*arSettlement.DepositStatus]
	}
	return ""
}

func (arSettlement ArSettlementResponse) GenerateDataDepositStatusName() string {
	if arSettlement.DepositStatus != nil {
		return dataDepositStatusName[*arSettlement.DepositStatus]
	}
	return ""
}

type DepositStatusLookupResponse struct {
	DepositStatus     *int64  `json:"deposit_status"`
	DepositStatusName *string `json:"deposit_status_name"`
}

func (depositStatus DepositStatusLookupResponse) GenerateDataDepositStatusName() string {
	if depositStatus.DepositStatus != nil {
		return dataDepositStatusName[*depositStatus.DepositStatus]
	}
	return ""
}

var dataPayTypeName = map[int64]string{
	1: "Cash",
	2: "Cheque",
	3: "Transfer",
	4: "Return",
	5: "Debit/Credit",
}

func (arSettlementPayment ArSettlementPaymentResponse) GenerateDataPayTypeName() string {
	if arSettlementPayment.PayType != nil {
		return dataPayTypeName[*arSettlementPayment.PayType]
	}
	return ""
}

const (
	AR_BRANCH_VERIFICATION_STATUS_NEED_REVIEW = 1
	AR_BRANCH_VERIFICATION_STATUS_APPROVED    = 2
	AR_BRANCH_VERIFICATION_STATUS_REJECTED    = 3
)

var arBranchVerificationStatusNameList = map[int64]string{
	AR_BRANCH_VERIFICATION_STATUS_NEED_REVIEW: "Need Review",
	AR_BRANCH_VERIFICATION_STATUS_APPROVED:    "Approved",
	AR_BRANCH_VERIFICATION_STATUS_REJECTED:    "Rejected",
}

type RejectVerifyReport struct {
	DepositNo                 string                          `json:"deposit_no"`
	CustId                    string                          `json:"cust_id"`
	Expense                   RejectVerifyExpense              `json:"expense"`
	ChequeGiro                []RejectVerifyDocAmount          `json:"cheque_giro"`
	BankTransfer              []RejectVerifyDocAmount          `json:"bank_transfer"`
	Return                    []RejectVerifyDocAmount          `json:"return"`
	Cndn                      []RejectVerifyCndnAmount         `json:"cndn"`
	DepositPaymentValidation  []RejectVerifyDepositPaymentItem `json:"deposit_payment_validation"`
}

type RejectVerifyExpense struct {
	TotalPaymentAmount float64                   `json:"total_payment_amount"`
	Items              []RejectVerifyExpenseItem `json:"items"`
}

type RejectVerifyExpenseItem struct {
	ExpenseId      int64   `json:"expense_id"`
	PaymentAmount  float64 `json:"payment_amount"`
	CurrentBalance float64 `json:"current_balance"`
	ExpenseExists  bool    `json:"expense_exists"`
}

type RejectVerifyDocAmount struct {
	DocumentNo             string  `json:"document_no"`
	SumAmount              float64 `json:"sum_amount"`
	CurrentPaidAmount      float64 `json:"current_paid_amount"`
	CurrentRemainingAmount float64 `json:"current_remaining_amount"`
	RowExists              bool    `json:"row_exists"`
}

type RejectVerifyCndnAmount struct {
	DocumentNo            string  `json:"document_no"`
	SumAmount             float64 `json:"sum_amount"`
	CurrentUsedAmount     float64 `json:"current_used_amount"`
	CurrentRemaningAmount float64 `json:"current_remaning_amount"`
	RowExists             bool    `json:"row_exists"`
}

type RejectVerifyDepositPaymentItem struct {
	DocumentNo    string  `json:"document_no"`
	PayType       int     `json:"pay_type"`
	PayTypeName   string  `json:"pay_type_name"`
	PaymentAmount float64 `json:"payment_amount"`
	RowExists     bool    `json:"row_exists"`
	PayTypeOK     bool    `json:"pay_type_ok"`
}
