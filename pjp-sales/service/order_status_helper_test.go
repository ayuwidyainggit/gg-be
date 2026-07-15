package service

import (
	"testing"

	"sales/entity"
	"sales/model"
)

func TestDetermineSalesOrderStatus(t *testing.T) {
	warning := LIMIT_ACTION_WARNING
	restricted := LIMIT_ACTION_RESTRICTED

	tests := []struct {
		name        string
		validation  entity.ValidateResponse
		outlet      salesOrderOutletRules
		wantStatus  int64
		wantBlocked bool
	}{
		{
			name: "all validations true should be processed",
			validation: entity.ValidateResponse{
				Validate1Success:  true,
				Validate2Success:  true,
				Validate3Success:  true,
				Validate4Success:  true,
				IsSuccessValidate: true,
			},
			wantStatus: int64(entity.PROCESSED),
		},
		{
			name: "credit limit warning should be processed",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: false,
				Validate3Success: true,
				Validate4Success: true,
			},
			outlet:     salesOrderOutletRules{CreditLimitAction: &warning},
			wantStatus: int64(entity.PROCESSED),
		},
		{
			name: "credit limit restrict should need review",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: false,
				Validate3Success: true,
				Validate4Success: true,
			},
			outlet:     salesOrderOutletRules{CreditLimitAction: &restricted},
			wantStatus: int64(entity.NEED_REVIEW),
		},
		{
			name: "overdue warning should be processed",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: true,
				Validate3Success: false,
				Validate4Success: true,
			},
			outlet:     salesOrderOutletRules{SalesInvLimitAction: &warning},
			wantStatus: int64(entity.PROCESSED),
		},
		{
			name: "overdue restrict should need review",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: true,
				Validate3Success: false,
				Validate4Success: true,
			},
			outlet:     salesOrderOutletRules{SalesInvLimitActionName: "Restricted"},
			wantStatus: int64(entity.NEED_REVIEW),
		},
		{
			name: "outstanding warning should be processed",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: true,
				Validate3Success: true,
				Validate4Success: false,
			},
			outlet:     salesOrderOutletRules{ObsLimitActionName: "Warning"},
			wantStatus: int64(entity.PROCESSED),
		},
		{
			name: "outstanding restrict should need review",
			validation: entity.ValidateResponse{
				Validate1Success: true,
				Validate2Success: true,
				Validate3Success: true,
				Validate4Success: false,
			},
			outlet:     salesOrderOutletRules{ObsLimitAction: &restricted},
			wantStatus: int64(entity.NEED_REVIEW),
		},
		{
			name: "stock validation failure should be blocked",
			validation: entity.ValidateResponse{
				Validate1Success: false,
				Validate2Success: true,
				Validate3Success: true,
				Validate4Success: true,
			},
			wantStatus:  int64(entity.NEED_REVIEW),
			wantBlocked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineSalesOrderStatus(tt.validation, tt.outlet)
			if got.DataStatus != tt.wantStatus {
				t.Fatalf("unexpected data status: got=%d want=%d", got.DataStatus, tt.wantStatus)
			}
			if got.Blocked != tt.wantBlocked {
				t.Fatalf("unexpected blocked flag: got=%v want=%v", got.Blocked, tt.wantBlocked)
			}
		})
	}
}

func TestValidationResultFromOrderList_NilValidateStokMessage(t *testing.T) {
	validation := validationResultFromOrderList(model.OrderList{
		ValidateStok:               false,
		ValidateStokMessage:        nil,
		ValidateCreditLimit:        true,
		ValidateCreditLimitMessage: "Within Limit",
		ValidateOverdue:            true,
		ValidateOverdueMessage:     "Allowed",
		ValidateOutstanding:        true,
		ValidateOutstandingMessage: "Allowed",
	})

	if validation.Validate1 != "" {
		t.Fatalf("expected empty validation message from nil pointer, got %q", validation.Validate1)
	}
	if validation.Validate1Success {
		t.Fatalf("expected stock validation bool to remain false, got %+v", validation)
	}
}

func TestNormalizeLimitAction(t *testing.T) {
	if got := normalizeLimitAction(nil, "restricted"); got != LIMIT_ACTION_RESTRICTED {
		t.Fatalf("expected restricted action, got %d", got)
	}
	if got := normalizeLimitAction(nil, "Warning"); got != LIMIT_ACTION_WARNING {
		t.Fatalf("expected warning action, got %d", got)
	}
}
