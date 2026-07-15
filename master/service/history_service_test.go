package service

import "testing"

func TestTempTableByUploadType_OutletNew(t *testing.T) {
	table, includeCtid, err := tempTableByUploadType("outlet-new")
	if err != nil {
		t.Fatalf("tempTableByUploadType(outlet-new) error: %v", err)
	}
	if table != "import.outlet_temp" {
		t.Fatalf("table = %q, want import.outlet_temp", table)
	}
	if includeCtid {
		t.Fatalf("includeCtid = true, want false")
	}
}

func TestIsOutletNewUploadType(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"outlet-new", true},
		{"Outlet-New", true},
		{" outlet-new ", true},
		{"outlet", false},
		{"outlet-update", false},
	}
	for _, tc := range cases {
		if got := isOutletNewUploadType(tc.in); got != tc.want {
			t.Fatalf("isOutletNewUploadType(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
