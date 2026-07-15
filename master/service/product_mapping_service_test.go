package service

import (
	"master/pkg/constant"
	"testing"
)

func TestValidateDistinctProductMappingUOM(t *testing.T) {
	tests := []struct {
		name    string
		largest string
		middle  string
		smallest string
		wantErr bool
	}{
		{
			name:    "valid unique uom",
			largest: "BOX",
			middle:  "PCS",
			smallest: "UNIT",
			wantErr: false,
		},
		{
			name:    "reject case-insensitive duplicate largest and middle",
			largest: "Box",
			middle:  "BOX",
			smallest: "UNIT",
			wantErr: true,
		},
		{
			name:    "reject case-insensitive duplicate largest and smallest",
			largest: "BOX",
			middle:  "PCS",
			smallest: "box",
			wantErr: true,
		},
		{
			name:    "duplicate largest and middle",
			largest: "BOX",
			middle:  "BOX",
			smallest: "UNIT",
			wantErr: true,
		},
		{
			name:    "duplicate largest and smallest",
			largest: "BOX",
			middle:  "",
			smallest: "BOX",
			wantErr: true,
		},
		{
			name:    "duplicate middle and smallest",
			largest: "BOX",
			middle:  "PCS",
			smallest: "PCS",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDistinctProductMappingUOM(tt.largest, tt.middle, tt.smallest)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateDistinctProductMappingUOM() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != constant.ProductMappingDuplicateUOMErrorMsg {
				t.Fatalf("validateDistinctProductMappingUOM() message = %q, want %q", err.Error(), constant.ProductMappingDuplicateUOMErrorMsg)
			}
		})
	}
}

func TestMapProductMappingUOMToProductColumns(t *testing.T) {
	unitID1, unitID2, unitID3 := mapProductMappingUOMToProductColumns("BOX", "PCS", "UNIT")
	if unitID1 != "UNIT" || unitID2 != "PCS" || unitID3 != "BOX" {
		t.Fatalf("unexpected mapping: unit_id1=%q unit_id2=%q unit_id3=%q", unitID1, unitID2, unitID3)
	}

	unitID1, unitID2, unitID3 = mapProductMappingUOMToProductColumns("BOX", "", "")
	if unitID1 != "" || unitID2 != "" || unitID3 != "BOX" {
		t.Fatalf("unexpected largest-only mapping: unit_id1=%q unit_id2=%q unit_id3=%q", unitID1, unitID2, unitID3)
	}

	unitID1, unitID2, unitID3 = mapProductMappingUOMToProductColumns("", "", "UNIT")
	if unitID1 != "UNIT" || unitID2 != "" || unitID3 != "" {
		t.Fatalf("unexpected smallest-only mapping: unit_id1=%q unit_id2=%q unit_id3=%q", unitID1, unitID2, unitID3)
	}
}

func TestBuildProductMappingImportCombinationKey(t *testing.T) {
	keyA := buildProductMappingImportCombinationKey(" TP-008 ", "PK1", "Wormie Hat")
	keyB := buildProductMappingImportCombinationKey("tp-008", " pk1 ", " wormie hat ")
	if keyA != keyB {
		t.Fatalf("expected normalized keys to match, got %q and %q", keyA, keyB)
	}

	keyC := buildProductMappingImportCombinationKey("TP-008", "PK2", "Wormie Hat")
	if keyA == keyC {
		t.Fatalf("expected different keys for different pro_code, got %q", keyC)
	}
}
