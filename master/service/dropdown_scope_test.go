package service

import (
	"testing"
)

func TestNormalizeDropdownScope(t *testing.T) {
	cases := map[string]string{
		"Specific": "specific",
		"SPESIFIC": "specific",
		"selected": "specific",
		"ALL":      "all",
		"":         "all",
		"random":   "all",
	}

	for in, want := range cases {
		got := NormalizeScopeSet(in, in, in)
		if got.RegionScope != want || got.AreaScope != want || got.DistributorScope != want {
			t.Fatalf("input %s want %s got %+v", in, want, got)
		}
	}
}
