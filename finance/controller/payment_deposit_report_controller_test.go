package controller

import (
	"finance/entity"
	"testing"
)

func TestNormalizeAndValidatePaymentDepositFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  entity.PaymentDepositReportQueryFilter
		wantErr string
	}{
		{
			name: "missing deposit type",
			filter: entity.PaymentDepositReportQueryFilter{
				StartDate: "1767225600",
				EndDate:   "1772323199",
			},
			wantErr: "deposit_type is required",
		},
		{
			name: "invalid deposit type",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR,XX"},
				StartDate:   "1767225600",
				EndDate:     "1772323199",
			},
			wantErr: "invalid deposit_type: XX",
		},
		{
			name: "missing start_date",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR"},
				EndDate:     "1772323199",
			},
			wantErr: "start_date is required",
		},
		{
			name: "missing end_date",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR"},
				StartDate:   "1767225600",
			},
			wantErr: "end_date is required",
		},
		{
			name: "valid deposit type csv and dates",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR, AP"},
				EmpID:       []string{"381, 421"},
				DepositNo:   []string{"DP1, DP2"},
				Sort:        "created_date:desc",
				StartDate:   "1767225600",
				EndDate:     "1772323199",
			},
			wantErr: "",
		},
		{
			name: "invalid date range",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR"},
				StartDate:   "1772323199",
				EndDate:     "1767225600",
			},
			wantErr: "invalid date range: end_date must be greater than or equal to start_date",
		},
		{
			name: "invalid sort field",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR"},
				Sort:        "foo:desc",
				StartDate:   "1767225600",
				EndDate:     "1772323199",
			},
			wantErr: "invalid sort field: foo",
		},
		{
			name: "invalid sort direction",
			filter: entity.PaymentDepositReportQueryFilter{
				DepositType: []string{"AR"},
				Sort:        "deposit_date:drop",
				StartDate:   "1767225600",
				EndDate:     "1772323199",
			},
			wantErr: "invalid sort direction: drop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := normalizeAndValidatePaymentDepositFilter(&tt.filter)
			if tt.wantErr == "" && err != nil {
				t.Fatalf("normalizeAndValidatePaymentDepositFilter() unexpected error = %v", err)
			}
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("normalizeAndValidatePaymentDepositFilter() error = nil, want %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("normalizeAndValidatePaymentDepositFilter() error = %q, want %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestNormalizeAndValidatePaymentDepositDownloadFilter(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		StartDate: "1767225600",
		EndDate:   "1772323199",
	}

	if err := normalizeAndValidatePaymentDepositDownloadFilter(&filter); err != nil {
		t.Fatalf("expected default AR+AP download filter to be valid without salesman_id, got %v", err)
	}
	if got, want := len(filter.DepositType), 2; got != want || filter.DepositType[0] != "AP" || filter.DepositType[1] != "AR" {
		t.Fatalf("expected default deposit type AP+AR, got %#v", filter.DepositType)
	}
}

func TestNormalizeAndValidatePaymentDepositFilterCurlContract(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		Page:        1,
		Limit:       10,
		Sort:        "created_date:desc",
		StartDate:   "1775001600",
		EndDate:     "1780271999",
		EmpID:       []string{"421,415,381"},
		DepositType: []string{"AR,AP"},
	}

	if err := normalizeAndValidatePaymentDepositListFilter(&filter); err != nil {
		t.Fatalf("expected curl contract filter to be valid, got %v", err)
	}
	if got, want := filter.Sort, "deposit_date:desc"; got != want {
		t.Fatalf("Sort = %q, want %q", got, want)
	}
	if got, want := len(filter.EmpID), 3; got != want {
		t.Fatalf("EmpID length = %d, want %d", got, want)
	}
	if got, want := filter.DepositType, []string{"AP", "AR"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("DepositType = %#v, want %#v", got, want)
	}
}

func TestNormalizeAndValidatePaymentDepositFilterDoesNotRequireCollector(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AR,AP"},
		StartDate:   "1775001600",
		EndDate:     "1780271999",
	}

	if err := normalizeAndValidatePaymentDepositListFilter(&filter); err != nil {
		t.Fatalf("expected filter without salesman_id and emp_id to be valid, got %v", err)
	}
}

func TestNormalizeAndValidatePaymentDepositDownloadFilterAPOnlyDoesNotRequireSalesman(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AP"},
		StartDate:   "1767225600",
		EndDate:     "1772323199",
	}

	if err := normalizeAndValidatePaymentDepositDownloadFilter(&filter); err != nil {
		t.Fatalf("expected AP-only filter to be valid, got %v", err)
	}
}

func TestNormalizeAndValidatePaymentDepositFilter_AllNormalizesToARAP(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"All"},
		StartDate:   "1767225600",
		EndDate:     "1772323199",
	}
	if err := normalizeAndValidatePaymentDepositFilter(&filter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filter.DepositType) != 2 || filter.DepositType[0] != "AP" || filter.DepositType[1] != "AR" {
		t.Fatalf("DepositType = %#v, want [AP AR]", filter.DepositType)
	}
}

func TestNormalizePaymentDepositFilter(t *testing.T) {
	filter := entity.PaymentDepositReportQueryFilter{
		DepositType: []string{"AR, AP"},
		EmpID:       []string{"381, 421"},
		SalesmanID:  []string{"11,12"},
		DepositNo:   []string{"DP1, DP2"},
	}

	normalizePaymentDepositFilter(&filter)

	if got, want := len(filter.DepositType), 2; got != want {
		t.Fatalf("DepositType length = %d, want %d", got, want)
	}
	if got, want := len(filter.EmpID), 2; got != want {
		t.Fatalf("EmpID length = %d, want %d", got, want)
	}
	if got, want := len(filter.SalesmanID), 2; got != want {
		t.Fatalf("SalesmanID length = %d, want %d", got, want)
	}
	if got, want := len(filter.DepositNo), 2; got != want {
		t.Fatalf("DepositNo length = %d, want %d", got, want)
	}
}
