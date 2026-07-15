package service

import (
	"master/entity"
	"testing"
)

func TestResolveDistributorCustIDsForFilter_DoesNotAppendJWTCustID(t *testing.T) {
	actual := resolveDistributorCustIDsForFilter([]string{"C260020102"})

	expected := []string{"C260020102"}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d cust ids, got %d: %+v", len(expected), len(actual), actual)
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected cust id %s at index %d, got %s", expected[i], i, actual[i])
		}
	}
}

func TestResolveDistributorCustIDsForFilter_RemovesEmptyAndDuplicateCustIDs(t *testing.T) {
	actual := resolveDistributorCustIDsForFilter([]string{"C260020102", "", " C260020102 ", "C260020103"})

	expected := []string{"C260020102", "C260020103"}
	if len(actual) != len(expected) {
		t.Fatalf("expected %d cust ids, got %d: %+v", len(expected), len(actual), actual)
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Fatalf("expected cust id %s at index %d, got %s", expected[i], i, actual[i])
		}
	}
}

func TestResolveImportNewSheetLayout_LegacyTemplateUsesResponseFieldHeaders(t *testing.T) {
	rows := [][]string{
		{"Template Field", "Outlet Code", "Outlet Name", "Phone Number", "Created Date", "Outlet Status", "Outlet Address"},
		{"Response Field", "outlet_code", "outlet_name", "phone_no", "outlet_establishment_date", "outlet_status_name", "address1"},
		{"", "OS-001", "Toko 1", "0812345", "6/15/2026", "Active", "Jalan 1"},
		{"", "OS-002", "Toko 2", "0812346", "6/15/2026", "Registered", "Jalan 2"},
	}

	layout := resolveImportNewSheetLayout(rows)

	if layout.dataColumnOffset != 1 {
		t.Fatalf("expected column offset 1, got %d", layout.dataColumnOffset)
	}
	if len(layout.headers) != 6 || layout.headers[0] != "outlet_code" {
		t.Fatalf("expected technical headers from response row, got %+v", layout.headers)
	}
	if len(layout.dataRows) != 2 {
		t.Fatalf("expected 2 data rows, got %d", len(layout.dataRows))
	}
}

func TestBuildImportNewRowMap_LegacyTemplateMapsMandatoryFields(t *testing.T) {
	headers := []string{"outlet_code", "outlet_name", "phone_no", "outlet_establishment_date", "outlet_status_name", "address1"}
	row := []string{"", "OS-001", "Toko 1", "0812345", "6/15/2026", "Active", "Jalan 1"}

	rowMap := buildImportNewRowMap(headers, row, 1)

	if rowMap["outlet_code"] != "OS-001" {
		t.Fatalf("expected outlet_code OS-001, got %q", rowMap["outlet_code"])
	}
	if rowMap["phone_no"] != "0812345" {
		t.Fatalf("expected phone_no, got %q", rowMap["phone_no"])
	}
	if rowMap["outlet_establishment_date"] != "6/15/2026" {
		t.Fatalf("expected outlet_establishment_date, got %q", rowMap["outlet_establishment_date"])
	}
	if rowMap["address1"] != "Jalan 1" {
		t.Fatalf("expected address1, got %q", rowMap["address1"])
	}
}

func TestIsImportNewMetadataRow_DetectsResponseRowWithoutColumnA(t *testing.T) {
	row := []string{"outlet_code", "outlet_name", "phone_no", "outlet_establishment_date", "outlet_status_name", "address1", "outlet_province", "outlet_regency", "outlet_sub_district", "longitude", "latitude", "ot_type_name"}
	if !isImportNewMetadataRow(row, 0) {
		t.Fatal("expected response-field row without column A to be detected as metadata")
	}
}

func TestNormalizeImportNewDate_SupportsUSExcelFormat(t *testing.T) {
	got := normalizeImportNewDate("6/15/2026")
	if got != "2026-06-15" {
		t.Fatalf("expected 2026-06-15, got %q", got)
	}
}

func TestApplyImportNewDefaults_PopulatesIdentityAndTaxFields(t *testing.T) {
	data := entity.OutletTemp{
		OutletName: "Toko 7",
		Address1:   "Jalan 2",
		PhoneNo:    "0812345",
	}
	applyImportNewDefaults(&data)

	if data.IdentityType != "National ID" {
		t.Fatalf("expected identity type National ID, got %q", data.IdentityType)
	}
	if data.IdentityNo != "0000000000000000" {
		t.Fatalf("expected default identity no, got %q", data.IdentityNo)
	}
	if data.TaxIdentifierType != "National ID" {
		t.Fatalf("expected tax identifier type National ID, got %q", data.TaxIdentifierType)
	}
	if data.TaxName != "Toko 7" {
		t.Fatalf("expected tax name from outlet name, got %q", data.TaxName)
	}
	if data.AddressTax != "Jalan 2" {
		t.Fatalf("expected tax address from address1, got %q", data.AddressTax)
	}
}

func TestParseImportDateValue_SupportsUSExcelFormat(t *testing.T) {
	got, err := parseImportDateValue("6/15/2026")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected parsed date")
	}
	if got.Format("2006-01-02") != "2026-06-15" {
		t.Fatalf("expected 2026-06-15, got %s", got.Format("2006-01-02"))
	}
}

func TestResolveImportVerificationStatus_ImportNewUsesApproved(t *testing.T) {
	if got := resolveImportVerificationStatus(entity.ImportRequest{IsImportNew: true}); got != 1 {
		t.Fatalf("expected verification_status 1 for import-new, got %d", got)
	}
}

func TestResolveImportVerificationStatus_StandardImportUsesInReview(t *testing.T) {
	if got := resolveImportVerificationStatus(entity.ImportRequest{IsImportNew: false}); got != 2 {
		t.Fatalf("expected verification_status 2 for standard import, got %d", got)
	}
}

func TestBuildImportNewRowMap_DisplayHeaderCreatedDate(t *testing.T) {
	headers := []string{"Outlet Code", "Outlet Name", "Phone Number", "Created Date", "Outlet Status", "Outlet Address"}
	row := []string{"OS-001", "Toko 1", "0812345", "6/15/2026", "Active", "Jalan 1"}

	rowMap := buildImportNewRowMap(headers, row, 0)

	if rowMap["outlet_establishment_date"] != "6/15/2026" {
		t.Fatalf("expected outlet_establishment_date from Created Date header, got %q", rowMap["outlet_establishment_date"])
	}
}

func TestExpandImportRowFromMap_CreatedDateMapsToEstablishmentDate(t *testing.T) {
	headers := make([]string, len(outletImportNewTemplateHeaders))
	for i, key := range outletImportNewTemplateHeaders {
		headers[i] = importNewDisplayHeader(key)
	}
	row := []string{
		"OS-001", "Toko 1", "0812345", "02-01-2026", "Active", "Jalan 1",
		"Jawa Barat", "Bandung", "Coblong", "107.6", "-6.9", "Retail",
	}
	rowMap := buildImportNewRowMap(headers, row, 0)
	fullRow := expandImportRowFromMap(entity.ImportRequest{IsImportNew: true}, rowMap)

	if fullRow[3] != "02-01-2026" {
		t.Fatalf("expected establishment date at index 3, got %q (fullRow=%v)", fullRow[3], fullRow)
	}

	data := entity.OutletTemp{
		OutletEstablishmentDate: fullRow[3],
	}
	applyImportNewDefaults(&data)
	if data.OutletEstablishmentDate != "2026-01-02" {
		t.Fatalf("expected normalized date 2026-01-02, got %q", data.OutletEstablishmentDate)
	}

	parsed, err := parseImportDateValue(data.OutletEstablishmentDate)
	if err != nil || parsed == nil {
		t.Fatalf("expected parsed date, err=%v", err)
	}
	if parsed.Format("2006-01-02") != "2026-01-02" {
		t.Fatalf("expected 2026-01-02, got %s", parsed.Format("2006-01-02"))
	}
}

func TestParseImportDateValue_SupportsExcelSerialNumber(t *testing.T) {
	got, err := parseImportDateValue("46188")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected parsed date from Excel serial")
	}
	if got.Format("2006-01-02") != "2026-06-15" {
		t.Fatalf("expected 2026-06-15 from Excel serial 46188, got %s", got.Format("2006-01-02"))
	}
}

func TestNormalizeCoordinateString_ExcelLocaleFormats(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"107.6", "107.6"},
		{"-6.9", "-6.9"},
		{"-6,9", "-6.9"},
		{"106,852371564894", "106.852371564894"},
		{"-6,255529494521080", "-6.255529494521080"},
		{"-6,255,529,494,521,080", "-6.255529494521080"},
		{"10,685,237,156,489,400", "10.685237156489400"},
		{"1.234.567,89", "1234567.89"},
	}
	for _, tc := range cases {
		got := normalizeCoordinateString(tc.in)
		if got != tc.want {
			t.Fatalf("normalizeCoordinateString(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseImportCoordinate_ValidatesRange(t *testing.T) {
	if _, err := parseImportCoordinate("106.852371", "Longitude", -180, 180); err != nil {
		t.Fatalf("expected valid longitude, got %v", err)
	}
	if _, err := parseImportCoordinate("-6.255529", "Latitude", -90, 90); err != nil {
		t.Fatalf("expected valid latitude, got %v", err)
	}
	if _, err := parseImportCoordinate("200", "Longitude", -180, 180); err == nil {
		t.Fatal("expected longitude out of range error")
	}
}

func TestApplyImportNewDefaults_NormalizesCoordinates(t *testing.T) {
	data := entity.OutletTemp{
		Longitude: "106,852,371,564,894",
		Latitude:  "-6,255,529,494,521,080",
	}
	applyImportNewDefaults(&data)
	if data.Latitude != "-6.25552949452108" {
		t.Fatalf("expected normalized latitude, got %q", data.Latitude)
	}
	if data.Longitude != "106.852371564894" {
		t.Fatalf("expected normalized longitude, got %q", data.Longitude)
	}
}
