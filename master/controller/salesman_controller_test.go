package controller

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestParseSalesmanDistributorIDs(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		expected    []int
		expectError bool
	}{
		{
			name:     "empty value",
			raw:      "",
			expected: nil,
		},
		{
			name:     "single value",
			raw:      "10",
			expected: []int{10},
		},
		{
			name:     "comma separated values",
			raw:      "10, 20, 30",
			expected: []int{10, 20, 30},
		},
		{
			name:     "duplicate values deduplicated",
			raw:      "10,20,10",
			expected: []int{10, 20},
		},
		{
			name:     "zero value preserved as principal scope",
			raw:      "0",
			expected: []int{0},
		},
		{
			name:     "zero value preserved with valid values",
			raw:      "0,20",
			expected: []int{0, 20},
		},
		{
			name:        "invalid alpha value",
			raw:         "abc",
			expectError: true,
		},
		{
			name:        "mixed invalid value",
			raw:         "1,abc,3",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseSalesmanDistributorIDs(tt.raw)
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

func TestParseSalesmanDistributorIDQuery_SupportsArrayStyle(t *testing.T) {
	var args fasthttp.Args
	args.Add("distributor_id[]", "10")
	args.Add("distributor_id[]", "20")

	actual, err := parseSalesmanDistributorIDQuery(&args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{10, 20}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected value %d at index %d, got %d", expected[i], i, actual[i])
		}
	}
}

func TestParseSalesmanDistributorIDQuery_SupportsRepeatedAndCommaSeparatedValues(t *testing.T) {
	var args fasthttp.Args
	args.Add("distributor_id", "0,10,20")
	args.Add("distributor_id", "30")

	actual, err := parseSalesmanDistributorIDQuery(&args)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{0, 10, 20, 30}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected value %d at index %d, got %d", expected[i], i, actual[i])
		}
	}
}
