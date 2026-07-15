package service

import "testing"

func TestOrderImportParserHelpers_SX2451(t *testing.T) {
	if v, ok := parseRequiredFloat(""); ok || v != 0 {
		t.Fatalf("parseRequiredFloat empty should be invalid; got v=%v ok=%v", v, ok)
	}
	if v, ok := parseRequiredFloat("12.5"); !ok || v != 12.5 {
		t.Fatalf("parseRequiredFloat 12.5 should be valid; got v=%v ok=%v", v, ok)
	}
	if _, ok := parseRequiredFloat("abc"); ok {
		t.Fatalf("parseRequiredFloat abc should be invalid")
	}
	if _, ok := parseOptionalFloat(""); !ok {
		t.Fatalf("parseOptionalFloat empty should be valid zero")
	}
	if _, ok := parseOptionalFloat("abc"); ok {
		t.Fatalf("parseOptionalFloat abc should be invalid")
	}
	if p, ok := intPtrFromStringStrict("0"); !ok || p == nil || *p != 0 {
		t.Fatalf("intPtrFromStringStrict 0 should be valid; got %v %v", p, ok)
	}
	if _, ok := intPtrFromStringStrict("3.5"); ok {
		t.Fatalf("intPtrFromStringStrict 3.5 should be invalid")
	}
	if p, ok := float64PtrFromStringStrict("99.5"); !ok || p == nil || *p != 99.5 {
		t.Fatalf("float64PtrFromStringStrict 99.5 should be valid; got %v %v", p, ok)
	}
	if _, ok := float64PtrFromStringStrict("oops"); ok {
		t.Fatalf("float64PtrFromStringStrict oops should be invalid")
	}
}
