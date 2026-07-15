package controller

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestParseOutletClassIDs(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		expected    []int
		expectError bool
	}{
		{
			name:     "single value",
			raw:      "123",
			expected: []int{123},
		},
		{
			name:     "comma separated values",
			raw:      "123, 456",
			expected: []int{123, 456},
		},
		{
			name:     "empty value",
			raw:      "",
			expected: nil,
		},
		{
			name:        "invalid value",
			raw:         "abc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseOutletClassIDs(tt.raw)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(actual) != len(tt.expected) {
				t.Fatalf("expected %d items, got %d", len(tt.expected), len(actual))
			}

			for i := range actual {
				if actual[i] != tt.expected[i] {
					t.Fatalf("expected value %d at index %d, got %d", tt.expected[i], i, actual[i])
				}
			}
		})
	}
}

func TestParseIntSliceQuery_IgnoresZeroValues(t *testing.T) {
	var args fasthttp.Args
	args.Add("distributor_id", "0")

	actual, err := parseIntSliceQuery(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if actual != nil {
		t.Fatalf("expected nil values when distributor_id is zero, got %+v", actual)
	}
}

func TestParseIntSliceQuery_IgnoresZeroAndKeepsValidValues(t *testing.T) {
	var args fasthttp.Args
	args.Add("distributor_id", "0,102")

	actual, err := parseIntSliceQuery(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := []int{102}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected value %d at index %d, got %d", expected[i], i, actual[i])
		}
	}
}

func TestParseIntSliceQueryAllowZero_KeepsPrincipalDistributorScope(t *testing.T) {
	var args fasthttp.Args
	args.Add("distributor_id", "0,102")

	actual, err := parseIntSliceQueryAllowZero(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := []int{0, 102}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected value %d at index %d, got %d", expected[i], i, actual[i])
		}
	}
}

func TestContainsInt_DetectsOutletStatusAllValue(t *testing.T) {
	if !containsInt([]int{0, 1, 5}, 0) {
		t.Fatalf("expected zero to be detected")
	}
	if containsInt([]int{1, 5, 6, 7}, 0) {
		t.Fatalf("did not expect zero to be detected")
	}
}

func TestIncludeInactiveQueryParser_ParsesExplicitFlag(t *testing.T) {
	values, err := url.ParseQuery("include_inactive=1&is_active=1")
	if err != nil {
		t.Fatalf("unexpected parse query error: %v", err)
	}

	var args fasthttp.Args
	for k, vals := range values {
		for _, v := range vals {
			args.Add(k, v)
		}
	}

	includeInactive, err := parseIntSliceQueryAllowZero(&args, "include_inactive", "include_inactive")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	isActive, err := parseIntSliceQueryAllowZero(&args, "is_active", "is_active")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(includeInactive, []int{1}) {
		t.Fatalf("expected include_inactive [1], got %v", includeInactive)
	}
	if !reflect.DeepEqual(isActive, []int{1}) {
		t.Fatalf("expected is_active [1], got %v", isActive)
	}
}
