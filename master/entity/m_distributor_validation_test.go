package entity

import (
	"testing"

	"master/pkg/validation"
)

func strPtr(v string) *string { return &v }
func intPtr(v int) *int       { return &v }
func boolPtr(v bool) *bool    { return &v }
func int64Ptr(v int64) *int64 { return &v }

func validCreateDistributorBody() CreateDistributorBody {
	return CreateDistributorBody{
		CustId:                "CUST001",
		ParentCustId:          "CUST001",
		CreatedBy:             int64Ptr(1001),
		DistributorCode:       "DIST001",
		DistributorName:       "Distributor One",
		Barcode:               "1234567890123",
		RegionId:              1,
		AreaId:                2,
		ChannelId:             3,
		SubDistributorGroupId: 4,
		DistPriceGrpId:        5,
		Address:               "Jl Test 1",
		Latitude:              "-6.200000",
		Longitude:             "106.816666",
		IsActive:              boolPtr(false),
		Contacts: []DistributorContact{
			{
				ContactName:  "John Doe",
				JobTitle:     "Manager",
				PhoneNo:      "081234567890",
				IsWaNo:       true,
				WaNo:         "081234567890",
				Email:        "",
				IdentityNo:   "ID12345",
				IdentityType: "National ID",
			},
		},
		Tax: []DistributorTax{
			{
				TaxIdentifierNoType: "NPWP",
				TaxIdentifierNo:     "1234567890",
				Nitku:               "NITKU001",
				TaxName:             "PT Distributor",
				TaxAddress:          "Jl Pajak",
			},
		},
	}
}

func validUpdateDistributorRequest() UpdateDistributorRequest {
	return UpdateDistributorRequest{
		CustId:                "CUST001",
		UpdatedBy:             int64Ptr(1002),
		DistributorCode:       strPtr("DIST001"),
		DistributorName:       strPtr("Distributor Updated"),
		RegionId:              intPtr(1),
		AreaId:                intPtr(2),
		ChannelId:             intPtr(3),
		SubDistributorGroupId: intPtr(4),
		DistPriceGrpId:        intPtr(5),
		Address:               strPtr("Jl Update"),
		Latitude:              strPtr("-6.210000"),
		Longitude:             strPtr("106.826666"),
		IsActive:              boolPtr(true),
		Contacts: []DistributorContactUpdate{
			{
				DistributorContactId: nil,
				ContactName:          strPtr("Jane Doe"),
				JobTitle:             strPtr("Supervisor"),
				PhoneNo:              strPtr("081111111111"),
				IsWaNo:               boolPtr(true),
				WaNo:                 strPtr("081111111111"),
				Email:                strPtr(""),
				IdentityNo:           strPtr("ID67890"),
				IdentityType:         strPtr("Passport"),
			},
		},
	}
}

func TestCreateDistributorBodyValidation_OptionalAddressAndEmailAllowed(t *testing.T) {
	validator := validation.NewValidator()
	payload := validCreateDistributorBody()
	payload.ProvinceId = ""
	payload.RegencyId = ""
	payload.SubDistrictId = ""
	payload.WardId = ""
	payload.ZipCode = ""
	payload.Contacts[0].Email = ""

	if err := validator.Validator.Struct(payload); err != nil {
		t.Fatalf("expected payload valid, got error: %v", err)
	}
}

func TestCreateDistributorBodyValidation_BarcodeRules(t *testing.T) {
	validator := validation.NewValidator()

	t.Run("barcode alphanumeric accepted", func(t *testing.T) {
		payload := validCreateDistributorBody()
		payload.Barcode = "BAR750759"
		if err := validator.Validator.Struct(payload); err != nil {
			t.Fatalf("expected alphanumeric barcode valid, got error: %v", err)
		}
	})

	t.Run("barcode with dash accepted", func(t *testing.T) {
		payload := validCreateDistributorBody()
		payload.Barcode = "BAR-750759"
		if err := validator.Validator.Struct(payload); err != nil {
			t.Fatalf("expected barcode with dash valid, got error: %v", err)
		}
	})

	t.Run("barcode space rejected", func(t *testing.T) {
		payload := validCreateDistributorBody()
		payload.Barcode = "BAR 750759"
		if err := validator.Validator.Struct(payload); err == nil {
			t.Fatalf("expected validation error for barcode with space")
		}
	})

	t.Run("barcode length over 20 rejected", func(t *testing.T) {
		payload := validCreateDistributorBody()
		payload.Barcode = "123456789012345678901"
		if err := validator.Validator.Struct(payload); err == nil {
			t.Fatalf("expected validation error for barcode length")
		}
	})
}

func TestUpdateDistributorRequestValidation_OptionalAddressAndEmailAllowed(t *testing.T) {
	validator := validation.NewValidator()
	payload := validUpdateDistributorRequest()
	payload.ProvinceId = nil
	payload.RegencyId = nil
	payload.SubDistrictId = nil
	payload.WardId = nil
	payload.ZipCode = nil
	payload.Barcode = nil
	payload.Contacts[0].Email = strPtr("")

	if err := validator.Validator.Struct(payload); err != nil {
		t.Fatalf("expected payload valid, got error: %v", err)
	}
}

func TestUpdateDistributorRequestValidation_IsActiveOptional(t *testing.T) {
	validator := validation.NewValidator()
	payload := validUpdateDistributorRequest()
	payload.IsActive = nil

	if err := validator.Validator.Struct(payload); err != nil {
		t.Fatalf("expected payload valid without is_active, got error: %v", err)
	}
}

func TestUpdateDistributorRequestValidation_InvalidZipCodeRejected(t *testing.T) {
	validator := validation.NewValidator()
	payload := validUpdateDistributorRequest()
	payload.ZipCode = strPtr("12-345")

	if err := validator.Validator.Struct(payload); err == nil {
		t.Fatalf("expected validation error for invalid zip_code")
	}
}

func TestUpdateDistributorRequestValidation_LocationOptionalScenarios(t *testing.T) {
	validator := validation.NewValidator()

	t.Run("update without location fields should be valid", func(t *testing.T) {
		payload := validUpdateDistributorRequest()
		payload.ProvinceId = nil
		payload.RegencyId = nil
		payload.SubDistrictId = nil
		payload.WardId = nil

		if err := validator.Validator.Struct(payload); err != nil {
			t.Fatalf("expected payload valid without location fields, got error: %v", err)
		}
	})

	t.Run("update with partial location fields should be valid", func(t *testing.T) {
		payload := validUpdateDistributorRequest()
		payload.RegencyId = nil
		payload.SubDistrictId = nil
		payload.WardId = nil
		payload.ProvinceId = strPtr("11")

		if err := validator.Validator.Struct(payload); err != nil {
			t.Fatalf("expected payload valid with partial location fields, got error: %v", err)
		}
	})

	t.Run("update with invalid location format should fail", func(t *testing.T) {
		payload := validUpdateDistributorRequest()
		payload.ProvinceId = strPtr("abc")

		if err := validator.Validator.Struct(payload); err == nil {
			t.Fatalf("expected validation error when province_id is invalid")
		}
	})
}

func TestCreateDistributorBodyValidation_DistributorCodeRules(t *testing.T) {
	validator := validation.NewValidator()

	tests := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{name: "accepts dash", code: "DIST-NEW1", wantErr: false},
		{name: "accepts underscore", code: "DIST_001", wantErr: false},
		{name: "backward compatible numeric", code: "12345", wantErr: false},
		{name: "rejects empty", code: "", wantErr: true},
		{name: "rejects 21 chars", code: "DIST12345678901234567", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := validCreateDistributorBody()
			payload.DistributorCode = tt.code
			err := validator.Validator.Struct(payload)
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error for code %q", tt.code)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected code %q valid, got error: %v", tt.code, err)
			}
		})
	}
}

func TestUpdateDistributorRequestValidation_DistributorCodeRules(t *testing.T) {
	validator := validation.NewValidator()

	tests := []struct {
		name    string
		code    *string
		wantErr bool
	}{
		{name: "accepts alphanumeric with dash", code: strPtr("DIST128-AB"), wantErr: false},
		{name: "rejects space", code: strPtr("DIST 001"), wantErr: true},
		{name: "rejects dot", code: strPtr("DIST.001"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := validUpdateDistributorRequest()
			payload.DistributorCode = tt.code
			err := validator.Validator.Struct(payload)
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected valid distributor code, got error: %v", err)
			}
		})
	}
}

func TestDistributorContactValidation_IdentityTypeOneOf(t *testing.T) {
	validator := validation.NewValidator()
	contact := DistributorContact{
		ContactName:  "John Doe",
		JobTitle:     "Manager",
		PhoneNo:      "081234567890",
		WaNo:         "081234567890",
		Email:        "",
		IdentityNo:   "ID12345",
		IdentityType: "KTP",
	}

	if err := validator.Validator.Struct(contact); err == nil {
		t.Fatalf("expected validation error for invalid identity_type")
	}
}
