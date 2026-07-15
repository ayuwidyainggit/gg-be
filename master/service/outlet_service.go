package service

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"

	// "math"
	"io"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"

	"master/entity"
	"master/model"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
)

var errUploadSecondaryOutletNotAllowed = errors.New("distributor tidak diizinkan untuk upload secondary outlet. Aktifkan Allow Upload Secondary Sales di Master Distributor")

// validation helpers for mandatory fields based on business steps
// AC-02: If Set as Outlet PKP = yes/1 then identity/tax fields are not required.
// AC-04: If payment_type = 1 or 2 (Cash On Delivery/Cash Before Delivery) then top is not required.
func validateCreateMandatory(row map[string]string) error {
	requiredFields := []string{
		"outlet_code",
		"outlet_name",
		"address1",
		"outlet_province",
		"outlet_regency",
		"outlet_sub_district",
		"outlet_ward",
		"zip_code",
		"ot_loc_name",
		"longitude",
		"latitude",
		"building_own",
		"outlet_establishment_date",
		"contact_name",
		"job_title",
		"contact_phone_no",
		"contact_is_wa_no",
		"contact_wa_no",
		"contact_email",
		"tax_invoice_form_name",
		"phone_no",
		"disc_grp_name",
		"ot_grp_name",
		"is_contra_bon",
		"price_grp_name",
		"district_name",
		"industry_name",
		"ot_class_name",
		"ot_type_name",
		"market_name",
		"delv_addr1",
		"delv_province",
		"delv_regency",
		"delv_sub_district",
		"delv_ward",
		"delv_longitude",
		"delv_latitude",
		"delv_zip_code",
		"delv_is_same_addr",
		"inv_addr1",
		"inv_province",
		"inv_regency",
		"inv_sub_district",
		"inv_ward",
		"inv_zip_code",
		"inv_is_same_addr",
		"payment_type_name",
		"ar_status_name",
		"bank_name",
		"account_no",
		"account_name",
	}

	// AC-02: If Set as Outlet PKP = yes/1, user does not need to fill identity_type, identity_no, tax fields
	isPkp := parseImportBool(row["is_pkp_outlet"], false) || parseImportBool(row["set_as_outlet_pkp"], false)
	if !isPkp {
		requiredFields = append(requiredFields, "identity_type", "identity_no", "tax_identifier_type", "tax_name", "address_tax")
	}

	// AC-04: If payment_type = 1 or 2, top is not required
	paymentTypeStr := strings.TrimSpace(strings.ToLower(row["payment_type_name"]))
	if paymentTypeStr != "cash on delivery" && paymentTypeStr != "cash before delivery" {
		pt := strings.TrimSpace(row["payment_type"])
		if pt != "1" && pt != "2" {
			requiredFields = append(requiredFields, "top")
		}
	}

	var missing []string
	for _, field := range requiredFields {
		if strings.TrimSpace(row[field]) == "" {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		displayNames := make([]string, 0, len(missing))
		for _, field := range missing {
			displayNames = append(displayNames, outletDisplayNameID(field))
		}
		return fmt.Errorf("kolom wajib belum diisi: %s", strings.Join(displayNames, ", "))
	}

	return nil
}

// validateImportMandatory: untuk import outlet, requirement hanya outlet_name (outlet_code bisa kosong, akan auto-generate).
func validateImportMandatory(row map[string]string) error {
	if strings.TrimSpace(row["outlet_name"]) == "" {
		return fmt.Errorf("kolom wajib belum diisi: %s", outletDisplayNameID("outlet_name"))
	}
	return nil
}

var outletImportNewTemplateHeaders = []string{
	"outlet_code",
	"outlet_name",
	"phone_no",
	"outlet_establishment_date",
	"outlet_status_name",
	"address1",
	"outlet_province",
	"outlet_regency",
	"outlet_sub_district",
	"longitude",
	"latitude",
	"ot_type_name",
}

var outletImportNewTemplateDisplayNames = map[string]string{
	"outlet_code":               "Outlet Code",
	"outlet_name":               "Outlet Name",
	"phone_no":                  "Phone Number",
	"outlet_establishment_date": "Created Date",
	"outlet_status_name":        "Outlet Status",
	"address1":                  "Outlet Address",
	"outlet_province":           "Province",
	"outlet_regency":            "City/Regency",
	"outlet_sub_district":       "Sub District",
	"longitude":                 "Longitude",
	"latitude":                  "Latitude",
	"ot_type_name":              "Outlet Type",
}

var outletImportNewInstructionByHeader = map[string]string{
	"outlet_code":               "alphanumeric with special character, maximum 10 char, Unique Value",
	"outlet_name":               "alphanumeric with special character, maximum 75 char",
	"phone_no":                  "Numeric, maximum 20 Char",
	"outlet_establishment_date": "format tanggal DD-MM-YYYY (02-01-2026)",
	"outlet_status_name":        "status outlet, status yang dicantumkan harus sesuai dengan settingan status yang ada di cust_id tersebut",
	"address1":                  "alphanumeric with special character, maximum 150 char",
	"outlet_province":           "Wajib ada di Province Name di Setup Parameter Web > Distributor > Province",
	"outlet_regency":            "Wajib ada di City/Regency Name di Setup Parameter Web > Distributor > City Regency",
	"outlet_sub_district":       "Wajib ada di Subdistrict Name di Setup Parameter Web > Distributor > Sub District",
	"longitude":                 "longitude toko",
	"latitude":                  "latitude toko",
	"ot_type_name":              "Wajib ada di Outlet Type Name di Setup Parameter Web > Outlet",
}

func importNewDisplayHeader(key string) string {
	if v, ok := outletImportNewTemplateDisplayNames[key]; ok {
		return v
	}
	return toDisplayHeader(key)
}

func importUploadType(req entity.ImportRequest) string {
	if req.IsImportNew {
		return "outlet-new"
	}
	return "outlet"
}

func stripLegacyImportNewMetadataRows(rows [][]string) [][]string {
	return filterImportNewDataRows(rows, 0)
}

type importNewSheetLayout struct {
	headers          []string
	dataRows         [][]string
	dataColumnOffset int
	headerRowNumber  int
}

func importNewFirstCell(row []string) string {
	if len(row) == 0 {
		return ""
	}
	return strings.TrimSpace(row[0])
}

func trimLeadingCells(row []string, offset int) []string {
	if offset <= 0 {
		out := make([]string, len(row))
		copy(out, row)
		return out
	}
	if len(row) <= offset {
		return nil
	}
	return row[offset:]
}

func importNewRowMatchesTemplateKeys(values []string) bool {
	if len(values) < len(outletImportNewTemplateHeaders) {
		return false
	}
	matchCount := 0
	for i, h := range outletImportNewTemplateHeaders {
		if strings.EqualFold(strings.TrimSpace(values[i]), h) {
			matchCount++
		}
	}
	return matchCount == len(outletImportNewTemplateHeaders)
}

func isImportNewMetadataRow(row []string, colOffset int) bool {
	first := importNewFirstCell(row)
	if strings.EqualFold(first, "Response Field") || strings.EqualFold(first, "Template Field") {
		return true
	}
	return importNewRowMatchesTemplateKeys(trimLeadingCells(row, colOffset))
}

func isImportNewEmptyDataRow(row []string, colOffset int) bool {
	for _, val := range trimLeadingCells(row, colOffset) {
		if strings.TrimSpace(val) != "" {
			return false
		}
	}
	return true
}

func filterImportNewDataRows(rows [][]string, colOffset int) [][]string {
	filtered := make([][]string, 0, len(rows))
	for _, row := range rows {
		if isImportNewMetadataRow(row, colOffset) {
			continue
		}
		if isImportNewEmptyDataRow(row, colOffset) {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func resolveImportNewSheetLayout(rows [][]string) importNewSheetLayout {
	layout := importNewSheetLayout{headerRowNumber: 1}
	if len(rows) == 0 {
		return layout
	}

	if strings.EqualFold(importNewFirstCell(rows[0]), "Template Field") {
		layout.dataColumnOffset = 1
		if len(rows) > 1 && strings.EqualFold(importNewFirstCell(rows[1]), "Response Field") {
			layout.headers = trimLeadingCells(rows[1], layout.dataColumnOffset)
			layout.headerRowNumber = 2
			layout.dataRows = filterImportNewDataRows(rows[2:], layout.dataColumnOffset)
			return layout
		}
		layout.headers = trimLeadingCells(rows[0], layout.dataColumnOffset)
		layout.dataRows = filterImportNewDataRows(rows[1:], layout.dataColumnOffset)
		return layout
	}

	layout.headers = rows[0]
	layout.dataRows = filterImportNewDataRows(rows[1:], 0)
	return layout
}

func buildImportNewRowMap(headers []string, rowValues []string, colOffset int) map[string]string {
	rowMap := make(map[string]string)
	values := trimLeadingCells(rowValues, colOffset)
	for j, val := range values {
		if j >= len(headers) {
			break
		}
		orig := strings.TrimSpace(headers[j])
		if orig == "" {
			continue
		}
		cleaned := canonicalizeHeader(orig)
		rowMap[strings.ToLower(cleaned)] = strings.TrimSpace(val)
	}
	return rowMap
}

func importTemplateHeaders(req entity.ImportRequest) []string {
	if req.IsImportNew {
		return outletImportNewTemplateHeaders
	}
	return outletExportHeaders
}

func validateImportRowMandatory(req entity.ImportRequest, row map[string]string) error {
	if req.IsImportNew {
		return validateImportNewMandatory(row)
	}
	return validateImportMandatory(row)
}

func importRowValue(row map[string]string, key string) string {
	keys := []string{key, strings.ToLower(key)}
	switch key {
	case "outlet_establishment_date":
		keys = append(keys, "created_date", "outlet_created_date")
	}
	for _, k := range keys {
		if v := strings.TrimSpace(row[k]); v != "" {
			return v
		}
	}
	return ""
}

func expandImportRowFromMap(req entity.ImportRequest, rowMap map[string]string) []string {
	headers := importTemplateHeaders(req)
	fullRow := make([]string, len(headers))
	for i, h := range headers {
		fullRow[i] = importRowValue(rowMap, canonicalizeHeader(h))
	}
	return fullRow
}

func validateImportNewMandatory(row map[string]string) error {
	required := []string{
		"outlet_name",
		"phone_no",
		"outlet_establishment_date",
		"address1",
		"outlet_province",
		"outlet_regency",
		"outlet_sub_district",
		"longitude",
		"latitude",
		"ot_type_name",
	}
	var missing []string
	for _, field := range required {
		if importRowValue(row, field) == "" {
			missing = append(missing, outletDisplayNameID(field))
		}
	}
	status := importRowValue(row, "outlet_status_name")
	if status == "" {
		status = importRowValue(row, "outlet_status")
	}
	if status == "" {
		missing = append(missing, "Outlet Status")
	}
	if len(missing) > 0 {
		return fmt.Errorf("kolom wajib belum diisi: %s", strings.Join(missing, ", "))
	}

	var invalid []string
	if _, err := parseImportCoordinate(importRowValue(row, "longitude"), "Longitude", -180, 180); err != nil {
		invalid = append(invalid, err.Error())
	}
	if _, err := parseImportCoordinate(importRowValue(row, "latitude"), "Latitude", -90, 90); err != nil {
		invalid = append(invalid, err.Error())
	}
	if len(invalid) > 0 {
		return fmt.Errorf("%s", strings.Join(invalid, "; "))
	}
	return nil
}

func setImportDefaultString(target *string, value string) {
	if strings.TrimSpace(*target) == "" {
		*target = value
	}
}

func importDateLayouts() []string {
	return []string{
		"02-01-2006", "2-1-2006",
		"2006-01-02",
		"02/01/2006", "2/1/2006",
		"2006/01/02",
		"1/2/2006", "01/02/2006",
		"1/2/06", "01/02/06",
	}
}

func tryParseExcelSerialDate(value string) (*time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false
	}
	matched, err := regexp.MatchString(`^\d+(\.\d+)?$`, value)
	if err != nil || !matched {
		return nil, false
	}
	serial, err := strconv.ParseFloat(value, 64)
	if err != nil || serial < 1 || serial > 2958465 {
		return nil, false
	}
	base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)
	wholeDays := int(serial)
	t := base.AddDate(0, 0, wholeDays)
	return &t, true
}

func normalizeImportNewDate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if t, ok := tryParseExcelSerialDate(value); ok && t != nil {
		return t.Format("2006-01-02")
	}
	for _, layout := range importDateLayouts() {
		if t, err := time.Parse(layout, value); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return value
}

func parseImportDateValue(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	if t, ok := tryParseExcelSerialDate(value); ok && t != nil {
		return t, nil
	}
	normalized := normalizeImportNewDate(value)
	if t, err := time.Parse("2006-01-02", normalized); err == nil {
		return &t, nil
	}
	for _, layout := range importDateLayouts() {
		if t, err := time.Parse(layout, value); err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("invalid date format: %s", value)
}

// normalizeCoordinateString menormalkan format koordinat dari Excel (locale ID/EU/US).
func normalizeCoordinateString(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "")

	commaCount := strings.Count(value, ",")
	dotCount := strings.Count(value, ".")

	switch {
	case commaCount > 0 && dotCount > 0:
		if strings.LastIndex(value, ",") > strings.LastIndex(value, ".") {
			value = strings.ReplaceAll(value, ".", "")
			value = strings.ReplaceAll(value, ",", ".")
		} else {
			value = strings.ReplaceAll(value, ",", "")
		}
	case commaCount > 1:
		if idx := strings.Index(value, ","); idx >= 0 {
			value = value[:idx] + "." + value[idx+1:]
			value = strings.ReplaceAll(value, ",", "")
		}
	case commaCount == 1:
		value = strings.ReplaceAll(value, ",", ".")
	case dotCount > 1:
		if idx := strings.LastIndex(value, "."); idx >= 0 {
			intPart := strings.ReplaceAll(value[:idx], ".", "")
			value = intPart + value[idx:]
		}
	}

	return value
}

func parseImportCoordinate(value, fieldLabel string, min, max float64) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s wajib diisi", fieldLabel)
	}

	normalized := normalizeCoordinateString(value)
	f, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return "", fmt.Errorf("%s format tidak valid", fieldLabel)
	}
	if f < min || f > max {
		return "", fmt.Errorf("%s di luar rentang valid (%.0f sampai %.0f)", fieldLabel, min, max)
	}
	return strconv.FormatFloat(f, 'f', -1, 64), nil
}

func resolveImportVerificationStatus(req entity.ImportRequest) int64 {
	if req.IsImportNew {
		return 1
	}
	return 2
}

func applyImportNewDefaults(data *entity.OutletTemp) {
	outletName := strings.TrimSpace(data.OutletName)
	address1 := strings.TrimSpace(data.Address1)

	data.OutletEstablishmentDate = normalizeImportNewDate(data.OutletEstablishmentDate)
	if v, err := parseImportCoordinate(data.Longitude, "Longitude", -180, 180); err == nil {
		data.Longitude = v
	}
	if v, err := parseImportCoordinate(data.Latitude, "Latitude", -90, 90); err == nil {
		data.Latitude = v
	}

	setImportDefaultString(&data.ContactName, outletName)
	setImportDefaultString(&data.IsEmbBail, "false")
	setImportDefaultString(&data.IdentityType, "National ID")
	setImportDefaultString(&data.IdentityNo, "0000000000000000")
	setImportDefaultString(&data.TaxIdentifierType, "National ID")
	setImportDefaultString(&data.TaxType, "National ID")
	setImportDefaultString(&data.TaxIdentifierNo, "0000000000000000")
	setImportDefaultString(&data.Nitku, "000000")
	setImportDefaultString(&data.TaxName, outletName)
	setImportDefaultString(&data.AddressTax, address1)

	setImportDefaultString(&data.BuildingOwn, "1")
	setImportDefaultString(&data.PaymentType, "1")
	setImportDefaultString(&data.ArStatus, "1")
	setImportDefaultString(&data.ArStatusName, "Normal")

	setImportDefaultString(&data.CreditLimitType, "1")
	setImportDefaultString(&data.CreditLimit, "1")
	setImportDefaultString(&data.CreditLimitActionName, "Warning")

	setImportDefaultString(&data.SalesInvLimitType, "1")
	setImportDefaultString(&data.SalesInvLimit, "1")
	setImportDefaultString(&data.SalesInvLimitActionName, "Warning")

	setImportDefaultString(&data.ObsType, "1")
	setImportDefaultString(&data.ObsTypeName, "1")
	setImportDefaultString(&data.Obs, "1")
	setImportDefaultString(&data.ObsLimitActionName, "Warning")

	setImportDefaultString(&data.BankCode, "000")
	setImportDefaultString(&data.BankName, "Non Bank")
	setImportDefaultString(&data.AccountNo, "000000")
	setImportDefaultString(&data.AccountName, "Non Bank")

	if strings.TrimSpace(data.ContactPhoneNo) == "" {
		data.ContactPhoneNo = strings.TrimSpace(data.PhoneNo)
	}
	if strings.TrimSpace(data.IsActive) == "" {
		data.IsActive = "true"
	}
}

func (service *outletServiceImpl) mapRowToOutletImportNewStruct(row []string) (entity.OutletTemp, error) {
	headers := outletImportNewTemplateHeaders
	headerIndex := make(map[string]int, len(headers))
	for i, key := range headers {
		headerIndex[key] = i
	}
	get := func(key string) string {
		if idx, ok := headerIndex[key]; ok && idx >= 0 && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}
	outletStatus := get("outlet_status_name")
	if outletStatus == "" {
		outletStatus = get("outlet_status")
	}
	data := entity.OutletTemp{
		OutletCode:              get("outlet_code"),
		OutletName:              get("outlet_name"),
		PhoneNo:                 get("phone_no"),
		OutletEstablishmentDate: get("outlet_establishment_date"),
		OutletStatus:            outletStatus,
		Address1:                get("address1"),
		OutletProvince:          get("outlet_province"),
		OutletRegency:           get("outlet_regency"),
		OutletSubDistrict:       get("outlet_sub_district"),
		Longitude:               get("longitude"),
		Latitude:                get("latitude"),
		OtTypeName:              get("ot_type_name"),
		IsActive:                "true",
	}
	applyImportNewDefaults(&data)
	return data, nil
}

func (service *outletServiceImpl) mapImportRowToOutletStruct(req entity.ImportRequest, row []string) (entity.OutletTemp, error) {
	if req.IsImportNew {
		return service.mapRowToOutletImportNewStruct(row)
	}
	return service.mapRowToOutletStruct(row)
}

func validateUpdateMandatory(req entity.UpdateOutletRequest) error {
	// Step 1 basics still required
	if strings.TrimSpace(req.OutletCode) == "" {
		return errors.New("outlet_code is required")
	}
	if req.OutletStatus != 0 && (req.OutletStatus < 1 || req.OutletStatus > 4) {
		return errors.New("outlet_status must be 1..4 when provided")
	}
	if strings.TrimSpace(req.Address1) == "" {
		return errors.New("address1 is required")
	}
	if req.BuldingOwn <= 0 {
		return errors.New("building_own is required")
	}

	// Main geo coordinates mandatory at update
	if strings.TrimSpace(req.Latitude) == "" {
		return errors.New("latitude is required")
	}
	if strings.TrimSpace(req.Longitude) == "" {
		return errors.New("longitude is required")
	}

	// Step 2 dropdowns
	if req.DiscGrpId <= 0 {
		return errors.New("disc_grp_id is required")
	}
	if req.OtGrpId <= 0 {
		return errors.New("ot_grp_id is required")
	}
	if req.PriceGrpId <= 0 {
		return errors.New("price_grp_id is required")
	}
	if req.DistrictId <= 0 {
		return errors.New("district_id is required")
	}
	if req.IndustryId <= 0 {
		return errors.New("industry_id is required")
	}
	if req.OtClassId <= 0 {
		return errors.New("ot_class_id is required")
	}
	if req.MarketId <= 0 {
		return errors.New("market_id is required")
	}
	if req.OtTypeId <= 0 {
		return errors.New("ot_type_id is required")
	}

	// Step 3: delivery/invoice toggles present, ensure zip codes if provided
	if req.InvIsSameAddress == nil {
		return errors.New("inv_is_same_addr is required")
	}
	if req.DelvIsSameAddress == nil {
		return errors.New("delv_is_same_addr is required")
	}

	// Step 4: payment
	if req.PaymentType < 1 || req.PaymentType > 3 {
		return errors.New("payment_type must be 1..3")
	}
	if req.ArStatus <= 0 {
		return errors.New("ar_status is required")
	}
	if req.PaymentType == 3 && req.Top <= 0 {
		return errors.New("top is required when payment_type is Credit")
	}

	// Bank
	if len(req.Details.OutletBank) == 0 {
		return errors.New("at least one bank is required")
	}
	b := req.Details.OutletBank[0]
	if b.BankId == nil || *b.BankId <= 0 {
		return errors.New("bank.bank_id is required")
	}
	if strings.TrimSpace(b.AccountNo) == "" {
		return errors.New("bank.account_no is required")
	}
	if strings.TrimSpace(b.AccountName) == "" {
		return errors.New("bank.account_name is required")
	}

	// Contact
	if len(req.Details.OutletContact) == 0 {
		return errors.New("at least one contact is required")
	}
	c := req.Details.OutletContact[0]
	if c.ContactName == nil || strings.TrimSpace(*c.ContactName) == "" {
		return errors.New("contact.contact_name is required")
	}
	if strings.TrimSpace(c.JobTitle) == "" {
		return errors.New("contact.job_title is required")
	}
	if strings.TrimSpace(c.IdentityType) == "" {
		return errors.New("contact.identity_type is required")
	}
	if strings.TrimSpace(c.IdentityNo) == "" {
		return errors.New("contact.identity_no is required")
	}
	if strings.TrimSpace(c.PhoneNo) == "" {
		return errors.New("contact.phone_no is required")
	}
	if strings.TrimSpace(c.WaNo) == "" {
		return errors.New("contact.wa_no is required")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("contact.email is required")
	}

	return nil
}

// normalizeHeader converts a header label from uploaded files into a stable key form
func normalizeHeader(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))

	// 1️⃣ Gabungkan teks + angka (contoh: "address 1" -> "address1")
	reJoin := regexp.MustCompile(`([a-z])\s+(\d)`)
	s = reJoin.ReplaceAllString(s, "${1}${2}")

	// 2️⃣ Ganti karakter pemisah umum dengan underscore (termasuk en-dash – dan em-dash —)
	repl := []struct{ old, new string }{
		{" ", "_"}, {"-", "_"}, {"–", "_"}, {"—", "_"}, {"/", "_"}, {"(", ""}, {")", ""}, {".", ""}, {",", ""},
	}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r.old, r.new)
	}

	// 3️⃣ Hilangkan underscore ganda
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	return s
}

// stripTrailingHeaderNotes removes descriptive suffix tokens (e.g. "_maksimal_10_karakter")
// from normalized headers so that core keys like "outlet_code" still match.
var headerNoteTokens = map[string]struct{}{
	"max": {}, "maximal": {}, "maximum": {},
	"min": {}, "minimal": {}, "minimum": {},
	"maks": {}, "maksimal": {}, "maksimum": {},
	"karakter": {}, "character": {}, "characters": {}, "char": {}, "chars": {},
	"digit": {}, "digits": {},
	"length": {}, "panjang": {},
	"wajib": {}, "harus": {}, "diisi": {}, "opsional": {}, "optional": {},
	"hanya": {}, "only": {}, "boleh": {}, "tidak": {}, "bisa": {},
	"lebih": {}, "melebihi": {}, "dari": {}, "hingga": {}, "sampai": {},
	"huruf": {}, "angka": {}, "kosong": {}, "null": {},
	"duplikat": {}, "duplicate": {}, "ganda": {}, "double": {},
	"diubah": {}, "ubah": {}, "perubahan": {}, "ganti": {}, "diganti": {},
	"isi": {}, "input": {}, "limit": {}, "maksimalnya": {}, "minimalnya": {},
	"catatan": {}, "notes": {},
}

func stripTrailingHeaderNotes(s string) (string, bool) {
	if s == "" {
		return s, false
	}

	tokens := strings.Split(s, "_")
	changed := false
	for len(tokens) > 0 {
		last := tokens[len(tokens)-1]
		if last == "" {
			tokens = tokens[:len(tokens)-1]
			changed = true
			continue
		}
		if _, ok := headerNoteTokens[last]; ok {
			tokens = tokens[:len(tokens)-1]
			changed = true
			continue
		}
		if _, err := strconv.Atoi(last); err == nil {
			// cek konteks: kalau sebelumnya termasuk token note (max/min/maksimal/dll) baru hapus
			if len(tokens) > 1 {
				prev := tokens[len(tokens)-2]
				if _, ok := headerNoteTokens[prev]; ok {
					tokens = tokens[:len(tokens)-1]
					changed = true
					continue
				}
			}
			break // angka valid → berhenti
		}
		break
	}

	if !changed || len(tokens) == 0 {
		return s, false
	}

	return strings.Join(tokens, "_"), true
}

// alias map from normalized display labels to canonical normalized keys
var aliasToCanonical = map[string]string{
	"status_toggle":                         "is_active",
	"contact_phone":                         "contact_phone_no",
	"set_as_whatsapp":                       "contact_is_wa_no",
	"set_as_whatsapp_toggle":                "contact_is_wa_no",
	"same_as_outlet_address_delivery":       "delv_is_same_addr",
	"same_as_outlet_address_invoice":        "inv_is_same_addr",
	"same_as_outlet_address_toggle_yesno":   "delv_is_same_addr",
	"same_as_outlet_address_invoice_toggle": "inv_is_same_addr",
	"delivery_address":                      "delv_addr1",
	"delivery_province":                     "delv_province",
	"delivery_city_regency":                 "delv_regency",
	"delivery_sub_district":                 "delv_sub_district",
	"delivery_village":                      "delv_ward",
	"delivery_postal_code":                  "delv_zip_code",
	"delivery_longitude":                    "delv_longitude",
	"delivery_latitude":                     "delv_latitude",
	"invoice_address":                       "inv_addr1",
	"invoice_province":                      "inv_province",
	"invoice_city_regency":                  "inv_regency",
	"invoice_sub_district":                  "inv_sub_district",
	"invoice_village":                       "inv_ward",
	"invoice_postal_code":                   "inv_zip_code",
	"ar_status":                             "ar_status_name",
	"credit_limit_idr":                      "credit_limit",
	"overdue_invoice_limit":                 "sales_inv_limit_type_name",
	"max_overdue_invoice":                   "sales_inv_limit",
	"overdue_limit_action":                  "sales_inv_limit_action_name",
	"outstanding_invoice_limit":             "obs_type_name",
	"max_outstanding_invoice":               "obs",
	"outstanding_limit_action":              "obs_limit_action_name",
	"phone_outlet":                          "phone_no",
	"outlet_whatsapp_number":                "wa_no",
	"outlet_whatsapp":                       "wa_no",
	"outlet_whatsapp_no":                    "wa_no",

	"status":                   "is_active",
	"set_as_pkp_outlet":        "is_pkp_outlet",
	"set_as_outlet_pkp":        "is_pkp_outlet",
	"is_pkp_outlet":            "is_pkp_outlet",
	"is_emb_bail":              "is_pkp_outlet", // legacy template column; historically labeled PKP
	"outlet_registration_date": "ot_reg_date",
	"outlet_start_date":        "ot_start_date",
	"addres1":                  "address1",
	"address":                  "address1",
	"province":                 "outlet_province",
	"province_id":              "outlet_province_id",
	"regency":                  "outlet_regency",
	"city_regency":             "outlet_regency",
	"regency_id":               "outlet_regency_id",
	"sub_district":             "outlet_sub_district",
	"sub_district_id":          "outlet_sub_district_id",
	"village":                  "outlet_ward",
	"ward":                     "outlet_ward",
	"ward_id":                  "outlet_ward_id",
	"postal_code":              "zip_code",
	// Contact
	"position":                "job_title",
	"identity_no":             "identity_no",
	"identity_type":           "identity_type",
	"whatsapp_number":         "contact_wa_no",
	"whatsapp_no":             "contact_wa_no",
	"wa_no":                   "contact_wa_no",
	"contact_whatsapp":        "contact_wa_no",
	"contact_whatsapp_no":     "contact_wa_no",
	"contact_whatsapp_number": "contact_wa_no",
	// Update-template contact display headers (spec: Contact Phone J=phone_no, Whatsapp L=wa_no, Email M=email; we use contact_* for contact table)
	"contact_name":            "contact_name",
	"phone":                   "phone_no",
	"whatsapp_number_contact": "contact_wa_no",
	"email":                   "contact_email",
	"identity_no.":            "identity_no",
	"contact_wa_active":       "contact_is_wa_no",
	"fax_number":              "fax_number",
	// Tax
	"tax_invoice_type": "tax_invoice_form_name",
	"tax_type":         "tax_identifier_type",
	"tax_number":       "tax_identifier_no",
	// Criteria
	"discount_group": "disc_grp_name",
	"outlet_group":   "ot_grp_name",
	"contra_bon":     "is_contra_bon",
	"control_bon":    "is_contra_bon",
	"price_group":    "price_grp_name",
	"district":       "district_name",
	"industry":       "industry_name",
	"outlet_class":   "ot_class_name",
	"outlet_type":    "ot_type_name",
	"agent_form":     "agent_from",
	"market":         "market_name",
	// Secondary outlet import-new template display headers
	"phone_number":          "phone_no",
	"created_date":          "outlet_establishment_date",
	"outlet_created_date":   "outlet_establishment_date",
	"outlet_address":        "address1",
	"outlet_status":  "outlet_status_name",
	// Payment/Validation
	"payment_type":          "payment_type_name",
	"terms_of_payment_top":  "top",
	"sales_weekly_average":  "avg_sales_week",
	"sales_monthly_average": "avg_sales_month",
	// Typos/variants observed
	"oby_type_name": "obs_type_name",
	// Update-template variants and synonyms
	"outlet_grp_code":  "outlet_grp_code", // keep as-is for update template
	"outlet_grp_name":  "outlet_grp_name",
	"outlet_grp_id":    "ot_grp_id",
	"outlet_loc_code":  "outlet_loc_code",
	"outlet_location":  "ot_loc_name",
	"outlet_loc_id":    "ot_loc_id",
	"outlet_type_code": "outlet_type_code",
	"outlet_type_name": "outlet_type_name",
	"outlet_type_id":   "ot_type_id",
	// accept *_code_new aliases
	"outlet_grp_code_new": "ot_grp_code_new",
	"outlet_loc_code_new": "ot_loc_code_new",
	"building_ownership":  "building_own",
	"bank":                "bank_name",
	"account_number":      "account_no",
	"tax_address":         "address_tax",
	"npwp":                "tax_identifier_no",
}

func cleanHeaderDisplayName(h string) string {
	if h == "" {
		return h
	}

	h = strings.TrimSpace(h)

	// hapus teks (maksimal XX karakter)
	re := regexp.MustCompile(`(?i)\s*\(.*maksimal.*karakter.*\)`)
	h = re.ReplaceAllString(h, "")

	// hapus tanda * di belakang
	h = strings.TrimSuffix(h, "*")
	h = strings.TrimSpace(h)

	return strings.TrimSpace(h)
}

// canonicalizeHeader returns the canonical normalized key used by code
func canonicalizeHeader(h string) string {
	h = cleanHeaderDisplayName(h)
	norm := normalizeHeader(h)
	logrus.Infof("Normalisasi header '%s' → '%s'", h, norm)
	if v, ok := aliasToCanonical[norm]; ok {
		return v
	}
	if trimmed, ok := stripTrailingHeaderNotes(norm); ok {
		if v, exists := aliasToCanonical[trimmed]; exists {
			return v
		}
		return trimmed
	}
	return norm
}

// toDisplayHeader converts an internal key (snake_case) into a display label used in templates
func toDisplayHeader(key string) string {
	// Specific overrides based on web labels
	overrides := map[string]string{
		"outlet_code":                 "Outlet Code",
		"outlet_name":                 "Outlet Name",
		"is_active":                   "Is Active",
		"outlet_status":               "Outlet Status",
		"address1":                    "Address",
		"outlet_province":             "Province Outlet Address",
		"outlet_regency":              "City/Regency Outlet Adddress",
		"outlet_sub_district":         "Subdistrict Outlet Address",
		"outlet_ward":                 "Village Outlet Address",
		"zip_code":                    "Postal Code Outlet Address",
		"ot_loc_name":                 "Outlet Location",
		"longitude":                   "Longitude",
		"latitude":                    "Latitude",
		"building_own":                "Building Ownership",
		"outlet_establishment_date":   "Outlet Created Date",
		"contact_name":                "Contact Name",
		"job_title":                   "Position",
		"is_pkp_outlet":               "Set as PKP Outlet",
		"identity_type":               "Identity Type",
		"identity_no":                 "Identity No",
		"contact_phone_no":            "Contact Phone",
		"contact_is_wa_no":            "Set As WhatsApp",
		"contact_wa_no":               "WhatsApp Number",
		"contact_email":               "Contact Email",
		"tax_invoice_form_name":       "Tax Invoice Type",
		"tax_identifier_type":         "Tax Type",
		"tax_identifier_no":           "Tax Identifier No",
		"nitku":                       "NITKU",
		"tax_name":                    "Tax Name",
		"address_tax":                 "Tax Address",
		"phone_no":                    "Outlet Phone",
		"fax_no":                      "Fax Number",
		"barcode":                     "Barcode",
		"disc_grp_name":               "Discount Group",
		"ot_grp_name":                 "Outlet Group",
		"is_contra_bon":               "Control Bon",
		"price_grp_name":              "Price Group",
		"district_name":               "District",
		"industry_name":               "Industry",
		"ot_class_name":               "Outlet Class",
		"ot_type_name":                "Outlet Type",
		"market_name":                 "Market",
		"agent_from":                  "Agent From",
		"delv_addr1":                  "Delivery Address",
		"delv_province":               "Delivery Province",
		"delv_regency":                "Delivery City/Regency",
		"delv_sub_district":           "Delivery Sub District",
		"delv_ward":                   "Delivery Village",
		"delv_longitude":              "Delivery Longgitude",
		"delv_latitude":               "Delivery Lattitude",
		"delv_zip_code":               "Delivery Postal Code",
		"delv_is_same_addr":           "Same As Outlet Address (Delivery)",
		"inv_addr1":                   "Invoice Address",
		"inv_province":                "Invoice Province",
		"inv_regency":                 "Invoice City/Regency",
		"inv_sub_district":            "Invoice Sub District",
		"inv_ward":                    "Invoice Village",
		"inv_zip_code":                "Invoice Postal Code",
		"inv_is_same_addr":            "Same As Outlet Address (Invoice)",
		"payment_type_name":           "Payment Type",
		"ar_status_name":              "AR Status",
		"bank_name":                   "Bank",
		"account_no":                  "Account Number",
		"account_name":                "Account Name",
		"top":                         "Terms of Payment (TOP)",
		"credit_limit_type_name":      "Credit Limit Type",
		"credit_limit":                "Credit Limit (IDR)",
		"credit_limit_action_name":    "Credit Limit Action",
		"sales_inv_limit_type_name":   "Overdue Invoice Limit",
		"sales_inv_limit":             "Max Overdue Invoice",
		"sales_inv_limit_action_name": "Overdue Limit Action",
		"obs_type_name":               "Outstanding Invoice Limit",
		"obs":                         "Max Outstanding Invoice",
		"obs_limit_action_name":       "Outstanding Invoice Action",
		"delv_addr2":                  "Delivery Address 2",
		"delv_province2":              "Delivery 2 Province",
		"delv_regency2":               "Delivery 2 City/Regency",
		"delv_sub_district2":          "Delivery 2 Sub District",
		"delv_ward2":                  "Delivery 2 Village",
		"delv_longitude2":             "Delivery 2 Longgitude",
		"delv_latitude2":              "Delivery 2 Lattitude",
		"delv_zip_code2":              "Delivery 2 Postal Code",
		"delv_city2":                  "Delivery 2 City/Regency",
		"principal_code":              "Outlet Principal Code",
		"outlet_principal_code":       "Outlet Principal Code",
		"ot_reg_date":                 "Outlet Registration Date",
		"closed_date":                 "Outlet Closed Date",
		"first_trans_date":            "First Transaction Date",
		"last_trans_date":             "Last Transaction Date",
		"ar_total":                    "AR Total",
	}
	if v, ok := overrides[key]; ok {
		return v
	}
	// Default Title Case
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	parts := strings.Split(strings.ReplaceAll(key, "_", " "), " ")
	for i, p := range parts {
		if p == "" {
			continue
		}
		switch strings.ToLower(p) {
		case "id":
			parts[i] = "ID"
		case "wa":
			parts[i] = "WA"
		default:
			r := []rune(p)
			for j := range r {
				if j == 0 {
					if r[j] >= 'a' && r[j] <= 'z' {
						r[j] = r[j] - ('a' - 'A')
					}
				} else {
					if r[j] >= 'A' && r[j] <= 'Z' {
						r[j] = r[j] + ('a' - 'A')
					}
				}
			}
			parts[i] = string(r)
		}
	}
	return strings.Join(parts, " ")
}

// ===== Enum label mappings (code <-> label) =====
// NOTE: These are used for Excel/CSV export (show labels) and import (accept labels).

// Normalization helper for human-entered labels
func normEnum(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// remove spaces, underscores and punctuation commonly used
	replacers := []string{" ", "", "_", "", "-", "", "/", "", ".", "", ",", ""}
	for i := 0; i < len(replacers); i += 2 {
		s = strings.ReplaceAll(s, replacers[i], replacers[i+1])
	}
	return s
}

// outlet_status
var outletStatusLabels = map[int]string{
	0: "Not set",
	1: "New Open Outlet",
	2: "Covered",
	3: "Non Active",
	4: "Closed",
	5: "Dormant",
	6: "Registered",
	7: "Active",
}
var outletStatusCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range outletStatusLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// payment_type
var paymentTypeLabels = map[int]string{1: "COD (Cash On Delivery)", 2: "CBD (Cash Before Delivery)", 3: "Credit"}
var paymentTypeCodes = func() map[string]int {
	m := map[string]int{"cod": 1, "cashondelivery": 1, "cbd": 2, "cashbeforedelivery": 2}
	for k, v := range paymentTypeLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// credit_limit_type
var creditLimitTypeLabels = map[int]string{1: "Unlimited", 2: "Limit By Total", 3: "Limit By Supplier"}
var creditLimitTypeCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range creditLimitTypeLabels {
		m[normEnum(v)] = k
	}
	m["unlinmited"] = 1
	return m
}()

// sales_inv_limit_type
var salesInvLimitTypeLabels = map[int]string{1: "Unlimited", 2: "Limit by Number of Invoice"}
var salesInvLimitTypeCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range salesInvLimitTypeLabels {
		m[normEnum(v)] = k
	}
	m["unlinmited"] = 1
	return m
}()

// obs_type
var obsTypeLabels = map[int]string{1: "Unlimited", 2: "Limit by Number of Invoice"}
var obsTypeCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range obsTypeLabels {
		m[normEnum(v)] = k
	}
	m["unlinmited"] = 1
	return m
}()

// Actions (credit_limit_action, sales_inv_limit_action, obs_limit_action)
var limitActionLabels = map[int]string{1: "Warning", 2: "Restricted"}
var limitActionCodes = func() map[string]int {
	m := map[string]int{"restrict": 2, "restricted": 2, "warning": 1}
	for k, v := range limitActionLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// verification_status
var verificationStatusLabels = map[int]string{1: "Approved", 2: "In Review", 3: "Reject"}
var verificationStatusCodes = func() map[string]int {
	m := map[string]int{"In Review": 2}
	for k, v := range verificationStatusLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// tax_invoice_form
var taxInvoiceFormLabels = map[int]string{1: "Faktur Pajak Standar", 2: "Faktur Pajak Gabungan"}
var taxInvoiceFormCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range taxInvoiceFormLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// ar_status
// Use entity constants for canonical labels to keep consistency across the app
var arStatusLabels = map[int]string{
	1: entity.AR_TYPE_NORMAL,
	2: entity.AR_TYPE_1X_GIRO_TOLAKAN,
	3: entity.AR_TYPE_2X_GIRO_TOLAKAN,
	4: entity.AR_TYPE_3X_GIRO_TOLAKAN,
}
var arStatusCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range arStatusLabels {
		m[normEnum(v)] = k
	}
	// friendly fallbacks
	m["normal"] = 1
	m["1xgirotolakan"] = 2
	m["2xgirotolakan"] = 3
	m["3xgirotolakan"] = 4
	return m
}()

var buildingOwnershipLabels = map[int]string{
	1: "Milik Sendiri",
	2: "Kontrak",
	3: "Sewa",
}

var buildingOwnershipCodes = func() map[string]int {
	m := map[string]int{}
	for k, v := range buildingOwnershipLabels {
		m[normEnum(v)] = k
	}
	return m
}()

// Helpers for formatting enum labels from pointers
func labelFromPtr(p *int, m map[int]string) string {
	if p == nil {
		return ""
	}
	if v, ok := m[*p]; ok {
		return v
	}
	return strconv.Itoa(*p)
}

func outletStatusExportLabel(p *int) string {
	if p == nil {
		return ""
	}
	if v, ok := outletStatusLabels[*p]; ok {
		return v
	}
	return fmt.Sprintf("Unknown (%d)", *p)
}

// Parse enum code from possibly-text input; if empty return empty string
func parseEnumCode(input string, nameMap map[string]int) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return ""
	}
	// Already numeric?
	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return s
	}
	if c, ok := nameMap[normEnum(s)]; ok {
		return strconv.Itoa(c)
	}
	return s // fallback; later parsing may error and be handled
}

// safeString returns a string from a *string, handling nil values
func safeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type OutletService interface {
	FindParentCustId(string) (entity.MCustomerResp, error)
	Detail(int64, string, string, string) (entity.OutletRespone, error)
	List(entity.OutletQueryFilter, string, string) (data []entity.OutletListRespone, total int, lastPage int, err error)
	ListByDistributor(entity.OutletQueryFilter) (data []entity.OutletListByDistributorRespone, total int, lastPage int, err error)
	OutletTypeList(entity.OutletQueryFilter, string, string) (data []entity.OutletTypeListRespone, total int, lastPage int, err error)
	OutletGroupList(entity.OutletQueryFilter, string, string) (data []entity.OutletGroupListRespone, total int, lastPage int, err error)
	Store(entity.CreateOutletBody) (entity.OutletRespone, error)
	Approve(entity.ApproveOutletBody) error
	Reject(entity.RejectOutletBody) error
	Update(int, entity.UpdateOutletRequest) error
	UpdateStatus(outletId int64, custId, parentCustId string, request entity.UpdateOutletStatusRequest, updatedBy int64) error
	UpdateStatuses() (int64, error)
	Delete(string, int, int64) error
	VerificationStatusList(entity.OutletQueryFilter, string, string) (data []entity.VerificationStatusListRespone, total int, lastPage int, err error)
	Export(filter entity.OutletQueryFilter) (*bytes.Buffer, string, string, error)
	ExportTemplate(format string, additional []string, fields []string) (*bytes.Buffer, string, string, error)
	ExportTemplateNew(format string) (*bytes.Buffer, string, string, error)
	ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error)
	ImportSecondaryCheck(jwtDistributorID int64, custId string) (*entity.OutletImportSecondaryCheckData, error)
	ValidateUploadSecondarySalesPermission(jwtDistributorID int64, custId string) error
	ImportOutletsFromXLSX(req entity.ImportRequest) error
	ImportOutletsFromCSV(req entity.ImportRequest) error
	ImportUpdateXLSX(req entity.ImportRequest) error
	ImportUpdateCSV(req entity.ImportRequest) error
	ImportUpdateXLS(req entity.ImportRequest) error
	ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error
	ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error
	OutletListApproval(dataFilter entity.OutletListApprovalQueryFilter, custId string) (data []entity.OutletListApprovalResponse, total int, lastPage int, err error)
	ApproveOutletList(request entity.OutletListApprovalRequest, custId string, userId int64) error
}

func NewOutletService(
	outletRepository repository.OutletRepository,
	outletConfigRepository repository.OutletConfigRepository,
	outletCodeRepository repository.OutletCodeRepository,
	distributorRepository repository.DistributorRepository,
) *outletServiceImpl {
	return &outletServiceImpl{
		OutletRepository:       outletRepository,
		OutletConfigRepository: outletConfigRepository,
		OutletCodeRepository:   outletCodeRepository,
		DistributorRepository:  distributorRepository,
	}
}

type outletServiceImpl struct {
	OutletRepository       repository.OutletRepository
	OutletConfigRepository repository.OutletConfigRepository
	OutletCodeRepository   repository.OutletCodeRepository
	DistributorRepository  repository.DistributorRepository
}

func isValidEmail(s string) bool {
	if s == "" {
		return true
	}
	return regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`).MatchString(strings.TrimSpace(s))
}

func (service *outletServiceImpl) resolveDistributorIDForImport(jwtDistributorID int64, custId string) (int64, error) {
	if jwtDistributorID > 0 {
		return jwtDistributorID, nil
	}
	if service.DistributorRepository == nil {
		return 0, errors.New("distributor repository not configured")
	}
	return service.DistributorRepository.FindDistributorIdByCustId(custId)
}

func (service *outletServiceImpl) ValidateUploadSecondarySalesPermission(jwtDistributorID int64, custId string) error {
	distributorID, err := service.resolveDistributorIDForImport(jwtDistributorID, custId)
	if err != nil || distributorID <= 0 {
		return errUploadSecondaryOutletNotAllowed
	}
	if service.DistributorRepository == nil {
		return errUploadSecondaryOutletNotAllowed
	}
	allowed, err := service.DistributorRepository.FindAllowUploadSecondarySalesByDistributorID(distributorID)
	if err != nil || !allowed {
		return errUploadSecondaryOutletNotAllowed
	}
	return nil
}

func (service *outletServiceImpl) ImportSecondaryCheck(jwtDistributorID int64, custId string) (*entity.OutletImportSecondaryCheckData, error) {
	distributorID, err := service.resolveDistributorIDForImport(jwtDistributorID, custId)
	if err != nil || distributorID <= 0 {
		return &entity.OutletImportSecondaryCheckData{
			DistributorId:             distributorID,
			AllowUploadSecondarySales: false,
		}, nil
	}
	if service.DistributorRepository == nil {
		return &entity.OutletImportSecondaryCheckData{
			DistributorId:             distributorID,
			AllowUploadSecondarySales: false,
		}, nil
	}
	allowed, err := service.DistributorRepository.FindAllowUploadSecondarySalesByDistributorID(distributorID)
	if err != nil {
		return nil, err
	}
	return &entity.OutletImportSecondaryCheckData{
		DistributorId:             distributorID,
		AllowUploadSecondarySales: allowed,
	}, nil
}

func (service *outletServiceImpl) FindParentCustId(custId string) (response entity.MCustomerResp, err error) {
	mCustomer, err := service.OutletRepository.FindOneParentCustId(custId)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(mCustomer, &response); err != nil {
		return response, err
	}

	return response, err
}

func (service *outletServiceImpl) Detail(outletId int64, custId string, parentCustId string, lang string) (response entity.OutletRespone, err error) {
	outlet, err := service.OutletRepository.FindOneByOutletIdAndCustId(outletId, custId, parentCustId)
	if err != nil {
		return response, err
	}
	err = structs.Automapper(outlet, &response)
	if err != nil {
		return response, err
	}
	if outlet.OutletPrincipalCode != nil {
		response.OutletPrincipalCode = *outlet.OutletPrincipalCode
	}
	if outlet.IsPkpOutlet != nil {
		response.IsPkpOutlet = *outlet.IsPkpOutlet
	} else if outlet.IsEmbBail != nil {
		response.IsPkpOutlet = *outlet.IsEmbBail
	}
	if outlet.OutletStatus != nil {
		response.OutletStatus = *outlet.OutletStatus
	}

	log.Println("outlet.ArStatus:", outlet.ArStatus)
	log.Println("AR_TYPE_NORMAL_ID:", entity.AR_TYPE_NORMAL_ID)
	log.Println("AR_TYPE_1X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_1X_GIRO_TOLAKAN_ID)
	log.Println("AR_TYPE_2X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_2X_GIRO_TOLAKAN_ID)
	log.Println("AR_TYPE_3X_GIRO_TOLAKAN_ID:", entity.AR_TYPE_3X_GIRO_TOLAKAN_ID)

	if *outlet.ArStatus == entity.AR_TYPE_NORMAL_ID {
		response.ArStatusName = entity.AR_TYPE_NORMAL
	} else if *outlet.ArStatus == entity.AR_TYPE_1X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_1X_GIRO_TOLAKAN
	} else if *outlet.ArStatus == entity.AR_TYPE_2X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_2X_GIRO_TOLAKAN
	} else if *outlet.ArStatus == entity.AR_TYPE_3X_GIRO_TOLAKAN_ID {
		response.ArStatusName = entity.AR_TYPE_3X_GIRO_TOLAKAN
	}

	if outlet.CloseDate != nil {
		closeDate := outlet.CloseDate.Format("2006-01-02")
		response.CloseDate = closeDate
	}

	outletSalesmanDetail, err := service.OutletRepository.GetDetailOutletSalesman(outletId, custId, lang)
	if err != nil {
		return response, err
	}
	for _, osDetail := range outletSalesmanDetail {
		var outletSalesResp entity.OutletSalesmanRead
		err = structs.Automapper(osDetail, &outletSalesResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletSalesman = append(response.Details.OutletSalesman, outletSalesResp)
	}

	outletBankDetail, err := service.OutletRepository.GetDetailOutletbank(outletId, custId)
	if err != nil {
		return response, err
	}
	for _, bankDetail := range outletBankDetail {
		var outletBankResp entity.OutletBankRead
		err = structs.Automapper(bankDetail, &outletBankResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletBank = append(response.Details.OutletBank, outletBankResp)
	}

	outletContactDetail, err := service.OutletRepository.GetDetailOutletContact(outletId, custId)
	if err != nil {
		return response, err
	}
	for _, contactDetail := range outletContactDetail {
		var outletContactResp entity.OutletContactRead
		err = structs.Automapper(contactDetail, &outletContactResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletContact = append(response.Details.OutletContact, outletContactResp)
	}

	outletTaxDetail, err := service.OutletRepository.GetDetailOutletTax(outletId, custId)
	if err != nil {
		return response, err
	}
	for _, TaxDetail := range outletTaxDetail {
		var outletTaxResp entity.OutletTaxRead
		err = structs.Automapper(TaxDetail, &outletTaxResp)
		if err != nil {
			return response, err
		}
		response.Details.OutletTax = append(response.Details.OutletTax, outletTaxResp)
	}

	// Ambil tanggal saat ini
	currentDate := time.Now()

	newDate := currentDate.AddDate(0, 0, response.Top)

	formattedDate := newDate.Format("2006-01-02")

	response.Duedate = formattedDate

	verificationStatusName := response.GenerateOutletVerificationStatusName()
	response.VerificationStatusName = &verificationStatusName

	response.TaxInvoiceFormName = response.GenerateTaxInvoiceFormName()

	return response, err
}

func mergeCustIdsForDistributorFilter(resolved []string, jwtCustId string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(resolved)+1)
	for _, id := range resolved {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	if cid := strings.TrimSpace(jwtCustId); cid != "" {
		if _, ok := seen[cid]; !ok {
			out = append(out, cid)
		}
	}
	return out
}

func resolveDistributorCustIDsForFilter(resolved []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0, len(resolved))
	for _, id := range resolved {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func distributorFilterIncludesPrincipalScope(distributorIDs []int) bool {
	for _, distributorID := range distributorIDs {
		if distributorID == 0 {
			return true
		}
	}
	return false
}

func positiveDistributorIDs(distributorIDs []int) []int {
	seen := make(map[int]struct{})
	ids := make([]int, 0, len(distributorIDs))
	for _, distributorID := range distributorIDs {
		if distributorID <= 0 {
			continue
		}
		if _, ok := seen[distributorID]; ok {
			continue
		}
		seen[distributorID] = struct{}{}
		ids = append(ids, distributorID)
	}
	return ids
}

func mergeOutletScopeCustIds(custIds ...[]string) []string {
	seen := make(map[string]struct{})
	merged := make([]string, 0)

	for _, group := range custIds {
		for _, custId := range group {
			custId = strings.TrimSpace(custId)
			if custId == "" {
				continue
			}
			if _, exists := seen[custId]; exists {
				continue
			}
			seen[custId] = struct{}{}
			merged = append(merged, custId)
		}
	}

	return merged
}

func flexStringTrimmed(p *entity.FlexString) string {
	if p == nil {
		return ""
	}
	return strings.TrimSpace(string(*p))
}

func (service *outletServiceImpl) upsertOutletGeoMastersForCreate(request entity.CreateOutletBody) error {
	geoCustID := strings.TrimSpace(request.ParentCustId)
	if geoCustID == "" {
		geoCustID = request.CustId
	}
	provID := flexStringTrimmed(request.OutletProvinceId)
	regID := flexStringTrimmed(request.OutletRegencyId)
	subID := flexStringTrimmed(request.OutletSubDistrictId)
	wardID := flexStringTrimmed(request.OutletWardId)

	if provID == "" && regID == "" && subID == "" && wardID == "" {
		return nil
	}

	provName := strings.TrimSpace(request.OutletProvince)
	if provName == "" && provID != "" {
		provName = provID
	}
	regName := strings.TrimSpace(request.OutletRegency)
	if regName == "" && regID != "" {
		regName = regID
	}
	subName := strings.TrimSpace(request.OutletSubDistrict)
	if subName == "" && subID != "" {
		subName = subID
	}
	wardName := strings.TrimSpace(request.OutletWard)
	if wardName == "" && wardID != "" {
		wardName = wardID
	}

	if provID != "" && provName != "" {
		if err := service.OutletRepository.UpsertProvince(geoCustID, provID, provName, request.CreatedBy); err != nil {
			return err
		}
	}
	if regID != "" && regName != "" {
		if err := service.OutletRepository.UpsertRegency(geoCustID, regID, regName, provID, request.CreatedBy); err != nil {
			return err
		}
	}
	if subID != "" && subName != "" {
		if err := service.OutletRepository.UpsertSubDistrict(geoCustID, subID, subName, provID, regID, request.CreatedBy); err != nil {
			return err
		}
	}
	if wardID != "" && wardName != "" {
		if err := service.OutletRepository.UpsertWard(geoCustID, wardID, wardName, provID, regID, subID, request.CreatedBy); err != nil {
			return err
		}
	}
	return nil
}

func (service *outletServiceImpl) upsertImportOutletGeoMasters(geoCustID string, data entity.OutletTemp, userId int64) error {
	if err := service.OutletRepository.UpsertProvince(geoCustID, data.OutletProvinceId, data.OutletProvince, userId); err != nil {
		return err
	}
	if err := service.OutletRepository.UpsertRegency(geoCustID, data.OutletRegencyId, data.OutletRegency, data.OutletProvinceId, userId); err != nil {
		return err
	}
	if err := service.OutletRepository.UpsertSubDistrict(geoCustID, data.OutletSubDistrictId, data.OutletSubDistrict, data.OutletProvinceId, data.OutletRegencyId, userId); err != nil {
		return err
	}
	if err := service.OutletRepository.UpsertWard(geoCustID, data.OutletWardId, data.OutletWard, data.OutletProvinceId, data.OutletRegencyId, data.OutletSubDistrictId, userId); err != nil {
		return err
	}
	return nil
}

func (service *outletServiceImpl) listOutletReadMultiCustForCitus(dataFilter entity.OutletQueryFilter, custId, parentCustId string, resolved []string) ([]model.OutletRead, int, int, error) {
	limit := dataFilter.Limit
	if limit == 0 {
		limit = 10
	}
	page := dataFilter.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	need := offset + limit

	total := 0
	for _, cid := range resolved {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		sub := dataFilter
		sub.ResolvedCustIdsForDistributor = []string{cid}
		_, t, _, err := service.OutletRepository.FindAllByCustId(sub, custId, parentCustId)
		if err != nil {
			return nil, 0, 0, err
		}
		total += t
	}

	seen := make(map[string]struct{})
	merged := make([]model.OutletRead, 0)
	for _, cid := range resolved {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		sub := dataFilter
		sub.ResolvedCustIdsForDistributor = []string{cid}
		sub.Page = 1
		sub.Limit = need
		rows, _, _, err := service.OutletRepository.FindAllByCustId(sub, custId, parentCustId)
		if err != nil {
			return nil, 0, 0, err
		}
		for _, row := range rows {
			key := row.CustId + ":" + strconv.Itoa(row.OutletId)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, row)
		}
	}
	sort.Slice(merged, func(i, j int) bool {
		return merged[i].OutletId > merged[j].OutletId
	})

	var outlets []model.OutletRead
	if offset < len(merged) {
		end := offset + limit
		if end > len(merged) {
			end = len(merged)
		}
		outlets = merged[offset:end]
	}
	lastPage := 0
	if total > 0 {
		lastPage = (total + limit - 1) / limit
	}
	return outlets, total, lastPage, nil
}

func (service *outletServiceImpl) List(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (data []entity.OutletListRespone, total int, lastPage int, err error) {
	// Filter status= (query param) → outlet_status di mst.m_outlet
	if dataFilter.OutletStatus == nil && strings.TrimSpace(dataFilter.Status) != "" {
		if s, parseErr := strconv.Atoi(strings.TrimSpace(dataFilter.Status)); parseErr == nil && s >= 0 && s <= 7 {
			dataFilter.OutletStatus = &s
		}
	}
	if len(dataFilter.DistributorID) > 0 {
		posIDs := positiveDistributorIDs(dataFilter.DistributorID)
		if len(posIDs) == 1 && posIDs[0] == 1 && !distributorFilterIncludesPrincipalScope(dataFilter.DistributorID) {
			scopeCust := strings.TrimSpace(custId)
			if scopeCust == "" {
				return data, 0, 0, nil
			}
			dataFilter.ResolvedCustIdsForDistributor = resolveDistributorCustIDsForFilter([]string{scopeCust})
		} else {
			resolvedGroups := make([][]string, 0, 2)
			if distributorFilterIncludesPrincipalScope(dataFilter.DistributorID) {
				scopeParentCustID := strings.TrimSpace(parentCustId)
				if scopeParentCustID == "" {
					scopeParentCustID = strings.TrimSpace(custId)
				}
				resolvedGroups = append(resolvedGroups, []string{scopeParentCustID})
			}
			if len(posIDs) > 0 {
				resolved, errResolve := service.OutletRepository.FindCustIdsByDistributorIds(parentCustId, posIDs)
				if errResolve != nil {
					return data, 0, 0, errResolve
				}
				resolvedGroups = append(resolvedGroups, resolved)
			}
			dataFilter.ResolvedCustIdsForDistributor = resolveDistributorCustIDsForFilter(mergeOutletScopeCustIds(resolvedGroups...))
			if len(dataFilter.ResolvedCustIdsForDistributor) == 0 {
				return data, 0, 0, nil
			}
		}
	}
	if len(dataFilter.DistributorID) == 0 && strings.TrimSpace(custId) == strings.TrimSpace(parentCustId) {
		resolved, errResolve := service.OutletRepository.FindCustIdsByParentCustId(parentCustId)
		if errResolve != nil {
			return data, 0, 0, errResolve
		}
		dataFilter.ResolvedCustIdsForDistributor = resolveDistributorCustIDsForFilter(resolved)
	}
	var outlets []model.OutletRead
	if len(dataFilter.ResolvedCustIdsForDistributor) > 1 {
		outlets, total, lastPage, err = service.listOutletReadMultiCustForCitus(dataFilter, custId, parentCustId, dataFilter.ResolvedCustIdsForDistributor)
	} else {
		outlets, total, lastPage, err = service.OutletRepository.FindAllByCustId(dataFilter, custId, parentCustId)
	}
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.OutletListRespone
			structs.Automapper(row, &vResp)
			if row.OutletStatusCode != nil {
				vResp.OutletStatus = *row.OutletStatusCode
			}
			if row.OutletStatusDesc != nil {
				vResp.OutletStatusName = *row.OutletStatusDesc
			}
			vResp.PreDormantStatus = row.PreDormantStatus
			if row.CreatedByName != nil {
				vResp.CreatedByName = *row.CreatedByName
			}
			if row.ContactName != nil {
				vResp.ContactName = *row.ContactName
			}
			if row.ContactPhoneNo != nil {
				vResp.ContactPhoneNo = *row.ContactPhoneNo
			}
			if strings.TrimSpace(vResp.PhoneNo) == "" && strings.TrimSpace(vResp.ContactPhoneNo) != "" {
				vResp.PhoneNo = vResp.ContactPhoneNo
			}
			if row.CloseDate != nil {
				closeDate := row.CloseDate.Format("2006-01-02")
				vResp.CloseDate = closeDate
			}

			// Ambil tanggal saat ini
			currentDate := time.Now()

			newDate := currentDate.AddDate(0, 0, *row.Top)

			formattedDate := newDate.Format("2006-01-02")

			vResp.Duedate = formattedDate

			verificationStatusName := vResp.GenerateOutletVerificationStatusName()
			vResp.VerificationStatusName = &verificationStatusName
			vResp.TaxInvoiceFormName = vResp.GenerateTaxInvoiceFormName()
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *outletServiceImpl) OutletTypeList(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (data []entity.OutletTypeListRespone, total int, lastPage int, err error) {
	outlets, total, lastPage, err := service.OutletRepository.FindAllOutletTypeByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.OutletTypeListRespone
			structs.Automapper(row, &vResp)
			if row.CloseDate != nil {
				// closeDate := row.CloseDate.Format("2006-01-02")
				// vResp.CloseDate = closeDate
			}
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *outletServiceImpl) OutletGroupList(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (data []entity.OutletGroupListRespone, total int, lastPage int, err error) {
	outlets, total, lastPage, err := service.OutletRepository.FindAllOutletGroupByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.OutletGroupListRespone
			structs.Automapper(row, &vResp)
			if row.CloseDate != nil {
				// closeDate := row.CloseDate.Format("2006-01-02")
				// vResp.CloseDate = closeDate
			}
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

func (service *outletServiceImpl) Store(request entity.CreateOutletBody) (response entity.OutletRespone, err error) {
	custIdForConfig := strings.TrimSpace(request.ParentCustId)
	if custIdForConfig == "" {
		custIdForConfig = request.CustId
	}
	if custIdForConfig == request.CustId {
		if parent, errParent := service.OutletRepository.FindOneParentCustId(request.CustId); errParent == nil && strings.TrimSpace(parent.ParentCustId) != "" {
			custIdForConfig = parent.ParentCustId
		}
	}
	cfgHeader, cfgDetails, cfgErr := service.OutletConfigRepository.FindActiveByCustId(custIdForConfig)
	if cfgErr != nil || cfgHeader == nil {
		return response, errors.New("Outlet config not found. Please setup outlet config first")
	}
	if len(cfgDetails) == 0 {
		return response, errors.New("Outlet config outlet status not found. Please setup outlet status in outlet config first")
	}

	updatedByStr := strconv.FormatInt(request.CreatedBy, 10)
	var trx repository.OutletTx
	var outletCodeConfig *model.OutletCode
	var nextSeqStr string
	var yearNow int
	var createdByName string

	outletCodeToUse := strings.TrimSpace(request.OutletCode)
	if outletCodeToUse == "0000000000" {
		outletCodeToUse = ""
	}
	// Hanya generate dari last sequence jika kode kosong (jika sudah diisi, pakai yang diisi — mis. dari import yang sudah auto-generate)
	if outletCodeToUse == "" {
		yearNow = time.Now().Year()
		createdByName = strings.TrimSpace(request.CreatedByName)
		statuses := []string{"active"}

		if createdByName != "" {
			outletCodeConfig, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(request.CustId, yearNow, statuses, createdByName)
			if (err != nil || outletCodeConfig == nil) && custIdForConfig != request.CustId {
				outletCodeConfig, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(custIdForConfig, yearNow, statuses, createdByName)
			}
		}

		if outletCodeConfig == nil && err == nil {
			outletCodeConfig, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatus(request.CustId, yearNow, statuses)
			if (err != nil || outletCodeConfig == nil) && custIdForConfig != request.CustId {
				outletCodeConfig, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatus(custIdForConfig, yearNow, statuses)
			}
		}

		if err != nil || outletCodeConfig == nil {
			return response, errors.New("Outlet code is required when outlet code config is not set up for current year (status Active)")
		}
	} else {
		// Kode sudah diisi (API atau dari import): jangan pakai config, jangan update last_sequence
		outletCodeConfig = nil
	}

	if len(outletCodeToUse) > 10 {
		return response, errors.New("Maksimal Outlet Code 10 Karakter")
	}
	if len(request.OutletName) > 75 {
		return response, errors.New("Maksimal Outlet Name 75 Karakter")
	}
	if request.Email != "" && !isValidEmail(request.Email) {
		return response, errors.New("Email Tidak Valid")
	}
	for _, c := range request.Details.OutletContact {
		if c.ContactName != nil && len(*c.ContactName) > 50 {
			return response, errors.New("Maksimal Contact Name 50 Karakter")
		}
		if len(c.JobTitle) > 20 {
			return response, errors.New("Maksimal position 20 Karakter")
		}
	}
	for _, c := range request.Details.OutletContact {
		identityType := strings.TrimSpace(c.IdentityType)
		identityNo := strings.TrimSpace(c.IdentityNo)
		// Identity dummy tetap boleh disimpan, tapi tidak ikut validasi unique antar outlet.
		if identityType != "" && identityNo != "" && identityNo != "0000000000000000" {
			exists, dupErr := service.OutletRepository.ExistsOutletContactIdentity(request.CustId, request.ParentCustId, c.IdentityType, c.IdentityNo, 0)
			if dupErr == nil && exists {
				return response, errors.New("Identity Type & Identity No already used by another outlet in this distributor")
			}
		}
	}

	if outletCodeToUse != "" {
		existing, err := service.OutletRepository.FindOneByOutletCodeAndCustId(outletCodeToUse, request.CustId, request.ParentCustId)
		if err == nil && existing.OutletId > 0 {
			return response, errors.New("Outlet Code has already taken: " + outletCodeToUse)
		}
	}

	var outletData model.Outlet

	if request.FirstTransDate != nil {
		if *request.FirstTransDate == "" {
			request.FirstTransDate = nil
		}
	}

	if request.LastTransDate != nil {
		if *request.LastTransDate == "" {
			request.LastTransDate = nil
		}
	}

	if request.OtStartDate != nil {
		if *request.OtStartDate == "" {
			request.OtStartDate = nil
		}
	}

	if request.OtRegDate != nil {
		if *request.OtRegDate == "" {
			request.OtRegDate = nil
		}
	}

	if request.Dob != nil {
		if *request.Dob == "" {
			request.Dob = nil
		}
	}

	if request.OutletEstablishmentDate != nil {
		if *request.OutletEstablishmentDate == "" {
			request.OutletEstablishmentDate = nil
		}
	}

	if request.CloseDate != nil {
		if *request.CloseDate == "" {
			request.CloseDate = nil
		} else {
			closedStart, err := str.DateStrToRfc3339String(*request.CloseDate)
			if err != nil {
				return response, err
			}
			request.CloseDate = &closedStart
		}
	}
	timeNow := time.Now().In(time.UTC)

	structs.Automapper(request, &outletData)
	isPkpOutlet := request.IsPkpOutlet
	outletData.IsPkpOutlet = &isPkpOutlet
	invCityStr := string(request.InvCity)
	outletData.InvCity = &invCityStr
	flexToStrPtr := func(p *entity.FlexString) *string {
		if p == nil {
			return nil
		}
		s := string(*p)
		return &s
	}
	flexToStrPtrOpt := func(p *entity.FlexString) *string {
		if p == nil {
			return nil
		}
		s := strings.TrimSpace(string(*p))
		if s == "" {
			return nil
		}
		return &s
	}
	strPtrOpt := func(s string) *string {
		v := strings.TrimSpace(s)
		if v == "" {
			return nil
		}
		return &v
	}
	outletData.OutletProvinceId = flexToStrPtrOpt(request.OutletProvinceId)
	outletData.OutletRegencyId = flexToStrPtrOpt(request.OutletRegencyId)
	outletData.OutletSubDistrictId = flexToStrPtrOpt(request.OutletSubDistrictId)
	outletData.OutletWardId = flexToStrPtrOpt(request.OutletWardId)
	outletData.DelvWardId = flexToStrPtr(request.DelvWardId)
	outletData.DelvZipCode = flexToStrPtr(request.DelvZipCode)
	outletData.DelvWardId2 = flexToStrPtr(request.DelvWardId2)
	outletData.DelvZipCode2 = flexToStrPtr(request.DelvZipCode2)
	outletData.InvWardId = flexToStrPtr(request.InvWardId)
	outletData.InvZipCode = flexToStrPtr(request.InvZipCode)

	creditLimitTypeName := strings.TrimSpace(request.CreditLimitTypeName)
	if creditLimitTypeName == "" {
		creditLimitTypeName = "Limit By Total"
	}
	creditLimitActionName := strings.TrimSpace(request.CreditLimitActionName)
	if creditLimitActionName == "" {
		creditLimitActionName = "Warning"
	}
	salesInvLimitTypeName := strings.TrimSpace(request.SalesInvLimitTypeName)
	if salesInvLimitTypeName == "" {
		salesInvLimitTypeName = "Limit by Invoice"
	}
	salesInvLimitActionName := strings.TrimSpace(request.SalesInvLimitActionName)
	if salesInvLimitActionName == "" {
		salesInvLimitActionName = "Warning"
	}
	outletData.CreditLimitTypeName = strPtrOpt(creditLimitTypeName)
	outletData.CreditLimitActionName = strPtrOpt(creditLimitActionName)
	outletData.SalesInvLimitTypeName = strPtrOpt(salesInvLimitTypeName)
	outletData.SalesInvLimitActionName = strPtrOpt(salesInvLimitActionName)

	statusForDb := request.Status
	if statusForDb == 0 {
		statusForDb = request.OutletStatus
	}
	if statusForDb >= 1 && statusForDb <= 7 {
		outletData.OutletStatus = &statusForDb
	}
	defaultOutletStatus := 0
	if len(cfgDetails) > 0 {
		defaultOutletStatus = cfgDetails[0].Status
	}
	if outletData.OutletStatus == nil || *outletData.OutletStatus == 0 {
		outletData.OutletStatus = &defaultOutletStatus
	}

	verificationStatus := 2
	outletData.VerificationStatus = &verificationStatus

	yyMMdd := timeNow.Format("060102") // YYMMDD
	regionCode, distributorCode, errRegion := service.OutletRepository.GetRegionAndDistributorCodeByCustId(request.CustId, request.ParentCustId)
	if errRegion != nil {
		return response, errRegion
	}
	var prefix string
	if regionCode != nil && distributorCode != nil {
		prefix = fmt.Sprintf("%s-%s-%s-", *regionCode, *distributorCode, yyMMdd)
	} else {
		prefix = fmt.Sprintf("%s-", yyMMdd)
	}
	trx, err = service.OutletRepository.TrxBegin()
	if err != nil {
		return response, err
	}
	seq, errSeq := service.OutletRepository.GetNextOutletPrincipalCodeSeqTx(trx.Tx(), prefix)
	if errSeq != nil {
		trx.TrxRollback()
		return response, errSeq
	}
	outletPrincipalCode := prefix + fmt.Sprintf("%04d", seq)
	outletData.OutletPrincipalCode = &outletPrincipalCode

	outletData.CreatedAt = &timeNow
	outletData.CreatedBy = &request.CreatedBy
	outletData.UpdatedBy = &request.UpdatedBy
	outletData.UpdatedAt = &timeNow

	if err = service.upsertOutletGeoMastersForCreate(request); err != nil {
		_ = trx.TrxRollback()
		return response, err
	}

	if outletCodeToUse == "" {
		if outletCodeConfig == nil {
			_ = trx.TrxRollback()
			return response, errors.New("Outlet code is required when outlet code config is not set up for current year (status Active)")
		}
		nextSeqStr, err = service.OutletCodeRepository.IncrementSequenceByID(outletCodeConfig.Id, &updatedByStr)
		if err != nil || strings.TrimSpace(nextSeqStr) == "" {
			_ = trx.TrxRollback()
			if err != nil {
				return response, err
			}
			return response, errors.New("Outlet code is required when outlet code config is not set up for current year (status Active)")
		}
		yearVal := outletCodeConfig.YearCode
		if yearVal >= 100 {
			yearVal = yearVal % 100
		}
		yearPart := fmt.Sprintf("%02d", yearVal)
		outletCodeToUse = outletCodeConfig.SerialCode + yearPart + nextSeqStr
		if len(outletCodeToUse) > 10 {
			_ = trx.TrxRollback()
			return response, errors.New("Kode Outlet Maksimal 10 Karakter: konfigurasi outlet code (serial + tahun + urutan) menghasilkan lebih dari 10 karakter. Persingkat Serial Code di setup outlet code")
		}
		request.OutletCode = outletCodeToUse
		outletData.OutletCode = &outletCodeToUse

		existing, errDup := service.OutletRepository.FindOneByOutletCodeAndCustId(outletCodeToUse, request.CustId, request.ParentCustId)
		if errDup == nil && existing.OutletId > 0 {
			_ = trx.TrxRollback()
			return response, errors.New("Outlet Code has already taken: " + outletCodeToUse)
		}
	}

	// outletId, err := service.OutletRepository.Store(outletData)
	// if err != nil {
	// 	return response, err
	// }

	if trx == nil {
		trx, err = service.OutletRepository.TrxBegin()
		if err != nil {
			return response, err
		}
	}
	err = trx.Store(&outletData)
	if err != nil {
		trx.TrxRollback()
		return response, err
	}

	response.OutletId = outletData.OutletId
	for _, detail := range request.Details.OutletSalesman {
		var outletSalesmanData model.MOutletSalesman
		err = structs.Automapper(detail, &outletSalesmanData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		outletSalesmanData.CustID = request.CustId
		outletSalesmanData.OutletID = int64(outletData.OutletId)
		err := trx.StoreDetailSalesman(&outletSalesmanData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	for _, detail := range request.Details.OutletBank {
		var outletBankData model.MOutletBank
		err = structs.Automapper(detail, &outletBankData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		outletBankData.CustID = request.CustId
		outletBankData.OutletID = int64(outletData.OutletId)
		err := trx.StoreDetailBank(&outletBankData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	for _, detail := range request.Details.OutletContact {
		var outletContactData model.MOutletContact
		err = structs.Automapper(detail, &outletContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		outletContactData.CustID = request.CustId
		outletContactData.OutletID = int64(outletData.OutletId)
		err := trx.StoreDetailContact(&outletContactData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	for _, detail := range request.Details.OutletTax {
		var outletTaxData model.MOutletTax
		err = structs.Automapper(detail, &outletTaxData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
		outletTaxData.CustID = request.CustId
		outletTaxData.OutletID = int64(outletData.OutletId)
		err := trx.StoreDetailTax(&outletTaxData)
		if err != nil {
			trx.TrxRollback()
			return response, err
		}
	}
	if outletCodeConfig != nil && strings.TrimSpace(nextSeqStr) != "" {
		if errUpd := service.OutletCodeRepository.UpdateLastSequenceNoWithTx(trx.Tx(), outletCodeConfig.Id, nextSeqStr, &updatedByStr); errUpd != nil {
			trx.TrxRollback()
			return response, errUpd
		}
	}
	trx.TrxCommit()

	return response, err
}

func (service *outletServiceImpl) Update(outletId int, request entity.UpdateOutletRequest) (err error) {

	// if err := validateUpdateMandatory(request); err != nil {
	// 	return err
	// }

	outletSalesIDs := []int64{}
	outletBankIDs := []int64{}
	outletContactIDs := []int64{}
	outletTaxIDs := []int64{}

	outlet, err := service.OutletRepository.FindOneByOutletCodeAndCustId(request.OutletCode, request.CustId, request.ParentCustId)
	if err == nil && outlet.OutletId != outletId {
		return errors.New("outlet_code: " + *outlet.OutletCode + " is already exists")
	}

	if request.CreditLimitType == nil {
		creditLimitType := 0
		request.CreditLimitType = &creditLimitType
	}

	if request.SalesInvLimitType == nil {
		salesInvLimitType := 0
		request.SalesInvLimitType = &salesInvLimitType
	}

	if request.ObsType == nil {
		obsType := 0
		request.ObsType = &obsType
	}

	if request.CreditLimitAction == nil {
		creditLimitAction := 0
		request.CreditLimitAction = &creditLimitAction
	}

	if request.SalesInvLimitAction == nil {
		salesInvLimitAction := 0
		request.SalesInvLimitAction = &salesInvLimitAction
	}

	if request.ObsLimitAction == nil {
		obsLimitAction := 0
		request.ObsLimitAction = &obsLimitAction
	}

	if request.Status != nil {
		request.OutletStatus = *request.Status
	}

	trx, err := service.OutletRepository.TrxBegin()
	if err != nil {
		return err
	}
	err = trx.Update(outletId, request)
	if err != nil {
		return err
	}

	/* Detail - Outlet Salesman */
	for _, detail := range request.Details.OutletSalesman {
		if detail.OutletSalesId != nil {
			outletSalesIDs = append(outletSalesIDs, *detail.OutletSalesId)
		}
	}

	if len(outletSalesIDs) > 0 {
		err := trx.DeleteDetailSalesNotIn(int64(outletId), outletSalesIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	} else {
		err := trx.DeleteDetailSalesByOutletId(int64(outletId))
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	for _, detail := range request.Details.OutletSalesman {
		if detail.OutletSalesId == nil || *detail.OutletSalesId == 0 {
			detail.OutletSalesId = nil
			var outletSalesData model.MOutletSalesman
			err = structs.Automapper(detail, &outletSalesData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			outletSalesData.CustID = request.CustId
			outletSalesData.OutletID = int64(outletId)

			err := trx.StoreDetail(&outletSalesData)
			if err != nil {
				trx.TrxRollback()
				return err
			}

		} else {
			var outletSalesData model.MOutletSalesman
			err = structs.Automapper(detail, &outletSalesData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			err := trx.UpdateDetailOutletSalesman(outletId, *detail.OutletSalesId, detail)
			if err != nil {
				log.Println("outletSalesService, UpdateDetail, err:", err.Error())
				trx.TrxRollback()
				return err
			}
		}
	}
	/* End of Detail - Outlet Salesman */

	/* Detail - Outlet Bank */
	for _, detail := range request.Details.OutletBank {
		if detail.OutletBankId != nil {
			outletBankIDs = append(outletBankIDs, *detail.OutletBankId)
		}
	}
	if len(outletBankIDs) > 0 {
		err := trx.DeleteDetailBankNotIn(int64(outletId), outletBankIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	} else {
		err := trx.DeleteDetailBankByOutletId(int64(outletId))
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	for _, detail := range request.Details.OutletBank {
		if detail.OutletBankId == nil || *detail.OutletBankId == 0 {
			detail.OutletBankId = nil
			var outletBankData model.MOutletBank
			err = structs.Automapper(detail, &outletBankData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			outletBankData.CustID = request.CustId
			outletBankData.OutletID = int64(outletId)

			err := trx.StoreDetailBank(&outletBankData)
			if err != nil {
				trx.TrxRollback()
				return err
			}

		} else {
			var outletBankData model.MOutletBank
			err = structs.Automapper(detail, &outletBankData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			err := trx.UpdateDetailOutletBank(outletId, *detail.OutletBankId, detail)
			if err != nil {
				log.Println("outletBankService, UpdateDetail, err:", err.Error())
				trx.TrxRollback()
				return err
			}
		}
	}
	/* End of Detail - Outlet Bank */

	/* Detail - Outlet Contact */
	for _, detail := range request.Details.OutletContact {
		if detail.OutletContactId != nil {
			outletContactIDs = append(outletContactIDs, *detail.OutletContactId)
		}
	}
	if len(outletContactIDs) > 0 {
		err := trx.DeleteDetailContactNotIn(int64(outletId), outletContactIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	} else {
		err := trx.DeleteDetailContactByOutletId(int64(outletId))
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	for _, detail := range request.Details.OutletContact {
		if detail.OutletContactId == nil || *detail.OutletContactId == 0 {
			detail.OutletContactId = nil
			var outletContactData model.MOutletContact
			err = structs.Automapper(detail, &outletContactData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			outletContactData.CustID = request.CustId
			outletContactData.OutletID = int64(outletId)

			err := trx.StoreDetailContact(&outletContactData)
			if err != nil {
				trx.TrxRollback()
				return err
			}

		} else {
			var outletContactData model.MOutletContact
			err = structs.Automapper(detail, &outletContactData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			err := trx.UpdateDetailOutletContact(outletId, *detail.OutletContactId, detail)
			if err != nil {
				log.Println("outletContactService, UpdateDetail, err:", err.Error())
				trx.TrxRollback()
				return err
			}
		}
	}
	/* End of Detail - Outlet Contact */

	/* Detail - Outlet Tax */
	for _, detail := range request.Details.OutletTax {
		if detail.OutletTaxId != nil {
			outletTaxIDs = append(outletTaxIDs, *detail.OutletTaxId)
		}
	}
	if len(outletTaxIDs) > 0 {
		err := trx.DeleteDetailTaxNotIn(int64(outletId), outletTaxIDs)
		if err != nil {
			trx.TrxRollback()
			return err
		}
	} else {
		err := trx.DeleteDetailTaxByOutletId(int64(outletId))
		if err != nil {
			trx.TrxRollback()
			return err
		}
	}
	for _, detail := range request.Details.OutletTax {
		if detail.OutletTaxId == nil || *detail.OutletTaxId == 0 {
			detail.OutletTaxId = nil
			var outletTaxData model.MOutletTax
			err = structs.Automapper(detail, &outletTaxData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			outletTaxData.CustID = request.CustId
			outletTaxData.OutletID = int64(outletId)

			err := trx.StoreDetailTax(&outletTaxData)
			if err != nil {
				trx.TrxRollback()
				return err
			}

		} else {
			var outletTaxData model.MOutletTax
			err = structs.Automapper(detail, &outletTaxData)
			if err != nil {
				trx.TrxRollback()
				return err
			}
			err := trx.UpdateDetailOutletTax(outletId, *detail.OutletTaxId, detail)
			if err != nil {
				log.Println("outletTaxService, UpdateDetail, err:", err.Error())
				trx.TrxRollback()
				return err
			}
		}
	}
	/* End of Detail - Outlet Tax */

	trx.TrxCommit()
	return err
}

func (service *outletServiceImpl) UpdateStatus(outletId int64, custId, parentCustId string, request entity.UpdateOutletStatusRequest, updatedBy int64) error {
	if request.Status == nil {
		return errors.New("status is required")
	}
	_, err := service.OutletRepository.FindOneByOutletIdAndCustId(outletId, custId, parentCustId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("outlet not found")
		}
		return err
	}
	return service.OutletRepository.UpdateOutletStatus(outletId, custId, *request.Status, updatedBy)
}

func (service *outletServiceImpl) UpdateStatuses() (int64, error) {
	ctx := context.Background()

	rows, err := service.OutletRepository.BulkUpdateStatuses(ctx)
	if err != nil {
		return 0, err
	}

	nooRows, err := service.OutletRepository.BulkPromoteRegisteredWithTransToNoo(ctx)
	if err != nil {
		return rows, err
	}

	return rows + nooRows, nil
}

func (service *outletServiceImpl) Delete(custId string, outletId int, userId int64) (err error) {

	err = service.OutletRepository.Delete(custId, outletId, userId)
	if err != nil {
		return err
	}
	return err
}

func (service *outletServiceImpl) Approve(request entity.ApproveOutletBody) (err error) {

	request.VerificationStatus = 1
	request.VerifiedAt = time.Now()

	trx, err := service.OutletRepository.TrxBegin()
	if err != nil {
		return err
	}

	err = trx.Approve(request)
	if err != nil {
		return err
	}

	trx.TrxCommit()

	return err
}

func (service *outletServiceImpl) Reject(request entity.RejectOutletBody) (err error) {

	request.VerificationStatus = 3
	request.VerifiedAt = time.Now()

	trx, err := service.OutletRepository.TrxBegin()
	if err != nil {
		return err
	}
	err = trx.Reject(request)
	if err != nil {
		return err
	}

	trx.TrxCommit()
	return err
}

func (service *outletServiceImpl) VerificationStatusList(dataFilter entity.OutletQueryFilter, custId, parentCustId string) (data []entity.VerificationStatusListRespone, total int, lastPage int, err error) {
	outlets, total, lastPage, err := service.OutletRepository.FindAllVerificationStatusByCustId(dataFilter, custId, parentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.VerificationStatusListRespone
			structs.Automapper(row, &vResp)
			verificationStatusName := vResp.GenerateOutletVerificationStatusName()
			vResp.VerificationStatusName = &verificationStatusName
			data = append(data, vResp)
		}
	}
	return data, total, lastPage, err
}

// exportOutletsFindAllPerCustID menjalankan FindAllExport sekali per cust_id agar aman di Citus
// (satu shard / satu nilai o.cust_id; IN multi-nilai memicu "complex joins are only supported...").
func (service *outletServiceImpl) exportOutletsFindAllPerCustID(base entity.OutletQueryFilter, custIDs []string) ([]model.OutletExport, error) {
	seen := make(map[string]struct{})
	outlets := make([]model.OutletExport, 0)
	for _, cid := range custIDs {
		cid = strings.TrimSpace(cid)
		if cid == "" {
			continue
		}
		f := base
		f.ResolvedCustIdsForDistributor = []string{cid}
		batch, _, err := service.OutletRepository.FindAllExport(f, base.CustId)
		if err != nil {
			return nil, err
		}
		for _, row := range batch {
			key := row.CustId + ":" + strconv.Itoa(row.OutletId)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			outlets = append(outlets, row)
		}
	}
	sort.Slice(outlets, func(i, j int) bool {
		return outlets[i].OutletId > outlets[j].OutletId
	})
	return outlets, nil
}

// fungsi export
func (service *outletServiceImpl) Export(filter entity.OutletQueryFilter) (*bytes.Buffer, string, string, error) {
	statusStr := strings.TrimSpace(filter.Status)
	multiStatus := strings.Contains(statusStr, ",")

	if !multiStatus && filter.OutletStatus == nil && statusStr != "" {
		if statusStr == "all" || statusStr == "0" {
		} else if s, err := strconv.Atoi(statusStr); err == nil && s >= 0 && s <= 7 {
			filter.OutletStatus = &s
		}
	}
	if !multiStatus && filter.OutletStatus == nil && filter.IsActive == nil && statusStr != "" {
		status := strings.ToLower(statusStr)
		switch status {
		case "active":
			val := 1
			filter.IsActive = &val
		case "deactive", "inactive":
			val := 2
			filter.IsActive = &val
		case "all", "0":
		default:
			if filter.OutletStatus == nil {
				return nil, "", "", fmt.Errorf("status must be active, deactive, all, or outlet_status 0-7")
			}
		}
	}

	var outlets []model.OutletExport
	var err error
	if len(filter.DistributorID) > 0 {
		resolved, errResolve := service.OutletRepository.FindCustIdsByDistributorIds(filter.ParentCustId, filter.DistributorID)
		if errResolve != nil {
			return nil, "", "", errResolve
		}
		merged := mergeCustIdsForDistributorFilter(resolved, filter.CustId)
		if len(merged) == 0 {
			outlets = []model.OutletExport{}
		} else if len(merged) == 1 {
			filter.ResolvedCustIdsForDistributor = merged
			outlets, _, err = service.OutletRepository.FindAllExport(filter, filter.CustId)
		} else {
			outlets, err = service.exportOutletsFindAllPerCustID(filter, merged)
		}
	} else {
		outlets, _, err = service.OutletRepository.FindAllExport(filter, filter.CustId)
	}
	if err != nil {
		logrus.Errorf("OutletService, Export, FindAllExport err: %v", err)
		return nil, "", "", err
	}

	if len(outlets) > 0 {
		parentCustId := ""
		if parent, errParent := service.OutletRepository.FindOneParentCustId(filter.CustId); errParent == nil && parent.ParentCustId != "" {
			parentCustId = parent.ParentCustId
		}
		if parentCustId != "" {
			missing := make(map[int64]struct{})
			for i := range outlets {
				if outlets[i].MarketId == nil {
					continue
				}
				if outlets[i].MarketName != nil && strings.TrimSpace(*outlets[i].MarketName) != "" {
					continue
				}
				id := int64(*outlets[i].MarketId)
				if id != 0 {
					missing[id] = struct{}{}
				}
			}
			if len(missing) > 0 {
				ids := make([]int64, 0, len(missing))
				for id := range missing {
					ids = append(ids, id)
				}
				if markets, errMarkets := service.OutletRepository.FindMarketsByIDs(parentCustId, ids); errMarkets == nil {
					for i := range outlets {
						if outlets[i].MarketId == nil {
							continue
						}
						if outlets[i].MarketName != nil && strings.TrimSpace(*outlets[i].MarketName) != "" {
							continue
						}
						id := int64(*outlets[i].MarketId)
						if market, ok := markets[id]; ok {
							name := market.MarketName
							code := market.MarketCode
							outlets[i].MarketName = &name
							outlets[i].MarketCode = &code
						}
					}
				} else {
					logrus.Warnf("OutletService, Export, fallback market lookup failed: %v", errMarkets)
				}
			}
		}
	}

	var buffer *bytes.Buffer
	var contentType, filename string

	switch filter.Format {
	case "csv":
		buffer, err = service.createOutletCSV(outlets)
		contentType = "text/csv"
		filename = "outlets.csv"
	case "xls":
		buffer, err = service.createOutletXLS(outlets)
		contentType = "application/vnd.ms-excel"
		filename = "outlets.xls"
	default:
		buffer, err = service.createOutletXLSX(outlets)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "outlets.xlsx"
	}

	if err != nil {
		return nil, "", "", err
	}
	return buffer, contentType, filename, nil
}

func (service *outletServiceImpl) createOutletXLSX(outlets []model.OutletExport) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Outlets"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := outletExportDataHeaders
	widths := make([]int, len(headers))
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		display := toDisplayHeader(header)
		f.SetCellValue(sheetName, cell, display)
		f.SetCellStyle(sheetName, cell, cell, style)
		widths[i] = utf8.RuneCountInString(display)
	}

	for r, rowData := range outlets {
		for c, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			val := formatOutletHeaderValue(rowData, header)
			f.SetCellValue(sheetName, cell, val)
			if l := utf8.RuneCountInString(val); l > widths[c] {
				widths[c] = l
			}
		}
	}

	for i := range headers {
		if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
			width := float64(widths[i] + 4)
			if width < 18 {
				width = 18
			}
			if width > 60 {
				width = 60
			}
			_ = f.SetColWidth(sheetName, colName, colName, width)
		}
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	return f.WriteToBuffer()
}

func (service *outletServiceImpl) createOutletXLS(outlets []model.OutletExport) (*bytes.Buffer, error) {
	return service.createOutletXLSX(outlets)
}

func (service *outletServiceImpl) createOutletCSV(outlets []model.OutletExport) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'

	headers := outletExportDataHeaders
	display := make([]string, len(headers))
	for i := range headers {
		display[i] = toDisplayHeader(headers[i])
	}
	if err := writer.Write(display); err != nil {
		return nil, err
	}

	for _, rowData := range outlets {
		record := make([]string, len(headers))
		for i, header := range headers {
			record[i] = formatOutletHeaderValue(rowData, header)
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer, nil
}

func safeIntPointer(value *int) string {
	if value != nil {
		return strconv.Itoa(*value)
	}
	return ""
}

func safeBoolPointer(value *bool) string {
	if value != nil {
		return strconv.FormatBool(*value)
	}
	return ""
}

func safeFloat64Pointer(value *float64) string {
	if value != nil {
		return strconv.FormatFloat(*value, 'f', 2, 64)
	}
	return ""
}

func boolToYesNo(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func boolToActiveDeactive(value bool) string {
	if value {
		return "Active"
	}
	return "Deactive"
}

func boolPtrToYesNo(value *bool) string {
	if value == nil {
		return ""
	}
	return boolToYesNo(*value)
}

func buildingOwnershipLabel(value *int) string {
	if value == nil {
		return ""
	}
	if lbl, ok := buildingOwnershipLabels[*value]; ok {
		return lbl
	}
	return strconv.Itoa(*value)
}

func safeInt64Pointer(value *int64) string {
	if value != nil {
		return strconv.FormatInt(*value, 10)
	}
	return ""
}

func formatOutletHeaderValue(row model.OutletExport, key string) string {
	switch key {
	case "outlet_code":
		return safeString(row.OutletCode)
	case "outlet_name":
		return safeString(row.OutletName)
	case "is_active":
		return boolToActiveDeactive(row.IsActive)
	case "outlet_status":
		return outletStatusExportLabel(row.OutletStatus)
	case "address1":
		return safeString(row.Address1)
	case "outlet_province":
		return safeString(row.OutletProvince)
	case "outlet_regency":
		return safeString(row.OutletRegency)
	case "outlet_sub_district":
		return safeString(row.OutletSubDistrict)
	case "outlet_ward":
		return safeString(row.OutletWard)
	case "zip_code":
		return safeString(row.ZipCode)
	case "ot_loc_name":
		return safeString(row.OtLocName)
	case "longitude":
		return safeString(row.Longitude)
	case "latitude":
		return safeString(row.Latitude)
	case "building_own":
		return buildingOwnershipLabel(row.BuldingOwn)
	case "outlet_establishment_date":
		// Export label "Outlet Created Date" — use record created_at (date only).
		if created := safeTimeDatePointer(row.CreatedAt); created != "" {
			return created
		}
		return safeString(row.OutletEstablishmentDate)
	case "contact_name":
		return safeString(row.ContactName)
	case "job_title":
		return safeString(row.JobTitle)
	case "identity_type":
		return safeString(row.IdentityType)
	case "identity_no":
		return safeString(row.IdentityNo)
	case "contact_phone_no":
		return safeString(row.ContactPhoneNo)
	case "contact_is_wa_no":
		return boolPtrToYesNo(row.ContactIsWaNo)
	case "contact_wa_no":
		return safeString(row.ContactWaNo)
	case "contact_email":
		return safeString(row.ContactEmail)
	case "tax_invoice_form_name":
		return labelFromPtr(row.TaxInvoiceForm, taxInvoiceFormLabels)
	case "tax_identifier_type":
		return safeString(row.TaxIdentifierType)
	case "tax_identifier_no":
		return safeString(row.TaxIdentifierNo)
	case "nitku":
		return safeString(row.Nitku)
	case "tax_name":
		return safeString(row.TaxName)
	case "address_tax":
		return safeString(row.AddressTax)
	case "phone_no":
		return safeString(row.PhoneNo)
	case "fax_no":
		return safeString(row.FaxNo)
	case "barcode":
		return safeString(row.Barcode)
	case "disc_grp_name":
		return safeString(row.DiscGrpName)
	case "ot_grp_name":
		return safeString(row.OtGrpName)
	case "is_contra_bon":
		return boolPtrToYesNo(row.IsContraBon)
	case "price_grp_name":
		return safeString(row.PriceGrpName)
	case "district_name":
		return safeString(row.DistrictName)
	case "industry_name":
		return safeString(row.Industryname)
	case "ot_class_name":
		return safeString(row.OtClassName)
	case "ot_type_name":
		return safeString(row.OtTypeName)
	case "market_name":
		return safeString(row.MarketName)
	case "agent_from":
		return safeString(row.AgentFrom)
	case "delv_addr1":
		return safeString(row.DelvAdd1)
	case "delv_province":
		return safeString(row.DelvProvince)
	case "delv_regency":
		return safeString(row.DelvRegency)
	case "delv_sub_district":
		return safeString(row.DelvSubDistrict)
	case "delv_ward":
		return safeString(row.DelvWard)
	case "delv_longitude":
		return safeString(row.DelvLongitude)
	case "delv_latitude":
		return safeString(row.DelvLatitude)
	case "delv_zip_code":
		return safeString(row.DelvZipCode)
	case "delv_is_same_addr":
		return boolPtrToYesNo(row.DelvIsSameAddress)
	case "inv_addr1":
		return safeString(row.InvAddr1)
	case "inv_province":
		return safeString(row.InvProvince)
	case "inv_regency":
		return safeString(row.InvRegency)
	case "inv_sub_district":
		return safeString(row.InvSubDistrict)
	case "inv_ward":
		return safeString(row.InvWard)
	case "inv_zip_code":
		return safeString(row.InvZipCode)
	case "inv_is_same_addr":
		return boolPtrToYesNo(row.InvIsSameAddress)
	case "payment_type_name":
		return labelFromPtr(row.PaymentType, paymentTypeLabels)
	case "ar_status_name":
		return labelFromPtr(row.ArStatus, arStatusLabels)
	case "bank_name":
		return safeString(row.BankName)
	case "account_no":
		return safeString(row.AccountNo)
	case "account_name":
		return safeString(row.AccountName)
	case "top":
		return safeIntPointer(row.Top)
	case "credit_limit_type_name":
		return labelFromPtr(row.CreditLimitType, creditLimitTypeLabels)
	case "credit_limit":
		return safeFloat64Pointer(row.CreditLimit)
	case "credit_limit_action_name":
		return labelFromPtr(row.CreditLimitAction, limitActionLabels)
	case "sales_inv_limit_type_name":
		return labelFromPtr(row.SalesInvLimitType, salesInvLimitTypeLabels)
	case "sales_inv_limit":
		return safeIntPointer(row.SalesInvLimit)
	case "sales_inv_limit_action_name":
		return labelFromPtr(row.SalesInvLimitAction, limitActionLabels)
	case "obs_type_name":
		return labelFromPtr(row.ObsType, obsTypeLabels)
	case "obs":
		return safeIntPointer(row.Obs)
	case "obs_limit_action_name":
		return labelFromPtr(row.ObsLimitAction, limitActionLabels)
	case "principal_code":
		return safeString(row.OutletPrincipalCode)
	case "ot_reg_date":
		return safeString(row.OtRegDate)
	case "closed_date":
		return safeTimePointer(row.CloseDate)
	case "first_trans_date":
		return safeString(row.FirstTransDate)
	case "last_trans_date":
		return safeString(row.LastTransDate)
	case "ar_total":
		return safeFloat64Pointer(row.ArTotal)
	default:
		return ""
	}
}

func safeTimePointer(value *time.Time) string {
	if value != nil {
		return value.Format("2006-01-02 15:04:05")
	}
	return ""
}

func safeTimeDatePointer(value *time.Time) string {
	if value != nil {
		return value.Format("2006-01-02")
	}
	return ""
}

func buildMandatoryOutletHeaderSet(instructions []entity.ImportInstruction) map[string]struct{} {
	mandatory := make(map[string]struct{})
	for _, ins := range instructions {
		if !ins.Mandatory {
			continue
		}
		key := canonicalizeHeader(ins.Kolom)
		if key != "" {
			mandatory[key] = struct{}{}
		}
	}

	if len(mandatory) == 0 {
		for key := range defaultMandatoryOutletHeaders {
			mandatory[key] = struct{}{}
		}
		return mandatory
	}

	for key := range defaultMandatoryOutletHeaders {
		if _, ok := mandatory[key]; !ok {
			mandatory[key] = struct{}{}
		}
	}

	return mandatory
}

// filterToMandatoryOutletHeaders returns only headers that are in the mandatory set (AC-01: template only mandatory fields).
func filterToMandatoryOutletHeaders(allHeaders []string, mandatory map[string]struct{}) []string {
	if len(mandatory) == 0 {
		return allHeaders
	}
	out := make([]string, 0, len(allHeaders))
	for _, h := range allHeaders {
		if _, ok := mandatory[canonicalizeHeader(h)]; ok {
			out = append(out, h)
		}
	}
	if len(out) == 0 {
		return allHeaders
	}
	return out
}

func isMandatoryOutletHeader(headerKey, displayName string, mandatory map[string]struct{}) bool {
	if len(mandatory) == 0 {
		return false
	}

	if _, ok := mandatory[canonicalizeHeader(headerKey)]; ok {
		return true
	}

	if displayName != "" {
		if _, ok := mandatory[canonicalizeHeader(displayName)]; ok {
			return true
		}
	}

	return false
}

// --- FUNGSI EKSPOR TEMPLATE ---

func (service *outletServiceImpl) ExportTemplate(format string, additional []string, fields []string) (*bytes.Buffer, string, string, error) {
	var buffer *bytes.Buffer
	var contentType, filename string
	var err error

	var optionalHeaders []string
	// Jika query "fields" dikirim (sama seperti export-template-update), gunakan logika yang sama.
	if len(fields) > 0 {
		optionalHeaders = resolveFieldsToOutletTemplateHeaders(fields)
	} else if len(additional) > 0 {
		// Aturan: (1) Jika hanya Sales Order Validation (Credit Limit / Overdue / Outstanding) dipilih → tampilkan 5 header wajib saja.
		// (2) Jika salah satu dari 3 validation tersebut dipilih dan ada additional lain → tampilkan 5 header wajib + header additional yang dipilih.
		hasValidationCredit := false
		hasValidationOverdue := false
		hasValidationOutstanding := false
		var otherAdditional []string
		for _, raw := range additional {
			token := normalizeHeader(raw)
			switch token {
			case "sales_order_validation_credit_limit":
				hasValidationCredit = true
			case "sales_order_validation_overdue", "sales_order_valiation_overdue":
				hasValidationOverdue = true
			case "sales_order_validation_outstanding":
				hasValidationOutstanding = true
			default:
				otherAdditional = append(otherAdditional, raw)
			}
		}
		var validationFive []string
		if hasValidationCredit {
			validationFive = []string{"outlet_code", "outlet_name", "credit_limit_type_name", "credit_limit", "credit_limit_action_name"}
		} else if hasValidationOverdue {
			validationFive = []string{"outlet_code", "outlet_name", "sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name"}
		} else if hasValidationOutstanding {
			validationFive = []string{"outlet_code", "outlet_name", "obs_type_name", "obs", "obs_limit_action_name"}
		}
		if len(validationFive) > 0 {
			// Hanya 5 header wajib untuk validation yang dipilih
			optionalHeaders = append([]string{}, validationFive...)
			// Jika ada additional lain, tambahkan header mereka setelah 5 header wajib
			if len(otherAdditional) > 0 {
				extra := resolveAdditionalOptions(otherAdditional)
				validationSet := make(map[string]struct{}, len(validationFive))
				for _, h := range validationFive {
					validationSet[h] = struct{}{}
				}
				otherSet := make(map[string]struct{})
				for _, h := range extra {
					if _, in := validationSet[h]; !in {
						otherSet[h] = struct{}{}
					}
				}
				otherOrdered := headersInOrder(outletTemplateColumnOrder, otherSet)
				optionalHeaders = append(optionalHeaders, otherOrdered...)
			}
		} else {
			extra := resolveAdditionalOptions(additional)
			if len(extra) > 0 {
				allowedSet := make(map[string]struct{}, len(outletTemplateMandatoryOnlyColumns)+len(extra))
				for _, h := range outletTemplateMandatoryOnlyColumns {
					allowedSet[h] = struct{}{}
				}
				for _, h := range extra {
					allowedSet[h] = struct{}{}
				}
				optionalHeaders = headersInOrder(outletTemplateColumnOrder, allowedSet)
			} else {
				optionalHeaders = nil
			}
		}
	}

	switch format {
	case "csv":
		// Package two CSV files into a ZIP: outlet_template.csv + instructions.csv
		templateBuf, err := service.createOutletTemplateCSV(optionalHeaders)
		if err != nil {
			return nil, "", "", err
		}
		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)
		// 1) outlet_template.csv
		tf, err := zw.Create("outlet_template.csv")
		if err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if _, err := tf.Write(templateBuf.Bytes()); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		instructionsForCsv, _ := service.OutletRepository.GetImportInstructions("outlet")
		iw, err := zw.Create("instructions.csv")
		if err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		cw := csv.NewWriter(iw)
		cw.Comma = ';'
		_ = cw.Write([]string{"Column Name (Template)", "Instruction"})
		for _, header := range outletInstructionColumnOrder {
			if header == "is_active" {
				continue
			}
			kolom := toDisplayHeader(header)
			keterangan := outletInstructionByHeader[header]
			if keterangan == "" {
				for _, ins := range instructionsForCsv {
					if canonicalizeHeader(ins.Kolom) == canonicalizeHeader(kolom) {
						keterangan = ins.Keterangan
						break
					}
				}
			}
			_ = cw.Write([]string{kolom, keterangan})
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if err := zw.Close(); err != nil {
			return nil, "", "", err
		}
		buffer = zipBuf
		contentType = "application/zip"
		filename = "outlet_template.zip"
	case "xls":
		buffer, err = service.createOutletTemplateXLS(optionalHeaders)
		contentType = "application/vnd.ms-excel"
		filename = "outlet_template.xls"
	default:
		buffer, err = service.createOutletTemplateXLSX(optionalHeaders)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "outlet_template.xlsx"
	}

	if err != nil {
		return nil, "", "", err
	}
	return buffer, contentType, filename, nil
}

func (service *outletServiceImpl) ExportTemplateNew(format string) (*bytes.Buffer, string, string, error) {
	var buffer *bytes.Buffer
	var contentType, filename string
	var err error

	switch format {
	case "csv":
		templateBuf, err := service.createOutletImportNewTemplateCSV()
		if err != nil {
			return nil, "", "", err
		}
		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)
		tf, err := zw.Create("outlet_template_new.csv")
		if err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if _, err := tf.Write(templateBuf.Bytes()); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		iw, err := zw.Create("instructions.csv")
		if err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		cw := csv.NewWriter(iw)
		cw.Comma = ';'
		_ = cw.Write([]string{"Column Name (Template)", "Instruction"})
		for _, header := range outletImportNewTemplateHeaders {
			kolom := importNewDisplayHeader(header)
			keterangan := outletImportNewInstructionByHeader[header]
			if keterangan == "" {
				keterangan = outletInstructionByHeader[header]
			}
			_ = cw.Write([]string{kolom, keterangan})
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if err := zw.Close(); err != nil {
			return nil, "", "", err
		}
		buffer = zipBuf
		contentType = "application/zip"
		filename = "outlet_template_new.zip"
	case "xls":
		buffer, err = service.createOutletImportNewTemplateXLSX()
		contentType = "application/vnd.ms-excel"
		filename = "outlet_template_new.xls"
	default:
		buffer, err = service.createOutletImportNewTemplateXLSX()
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "outlet_template_new.xlsx"
	}

	if err != nil {
		return nil, "", "", err
	}
	return buffer, contentType, filename, nil
}

func (service *outletServiceImpl) createOutletImportNewTemplateXLSX() (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Outlet Template"
	index, _ := f.NewSheet(sheetName)
	headers := outletImportNewTemplateHeaders
	widths := make([]int, len(headers))

	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for i, header := range headers {
		displayName := importNewDisplayHeader(header)
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, displayName)
		f.SetCellStyle(sheetName, cell, cell, boldStyle)
		widths[i] = utf8.RuneCountInString(displayName)
	}

	insSheet := "Instructions"
	f.NewSheet(insSheet)
	ih := []string{"Column Name (Template)", "Instruction"}
	for i, h := range ih {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(insSheet, cell, h)
		styleID, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}})
		f.SetCellStyle(insSheet, cell, cell, styleID)
	}
	rowIdx := 2
	for _, header := range outletImportNewTemplateHeaders {
		kolom := importNewDisplayHeader(header)
		keterangan := outletImportNewInstructionByHeader[header]
		if keterangan == "" {
			keterangan = outletInstructionByHeader[header]
		}
		f.SetCellValue(insSheet, "A"+strconv.Itoa(rowIdx), kolom)
		f.SetCellValue(insSheet, "B"+strconv.Itoa(rowIdx), keterangan)
		rowIdx++
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	for i := range widths {
		if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
			width := float64(widths[i] + 4)
			if width < 18 {
				width = 18
			}
			if width > 60 {
				width = 60
			}
			_ = f.SetColWidth(sheetName, colName, colName, width)
		}
	}
	return f.WriteToBuffer()
}

func (service *outletServiceImpl) createOutletImportNewTemplateCSV() (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'

	headers := outletImportNewTemplateHeaders
	templateRow := make([]string, len(headers))
	for i, header := range headers {
		templateRow[i] = importNewDisplayHeader(header)
	}
	_ = writer.Write(templateRow)

	writer.Flush()
	return buffer, writer.Error()
}

func (service *outletServiceImpl) createOutletTemplateXLSX(optionalHeaders []string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Outlet Template"
	index, _ := f.NewSheet(sheetName)

	instructions, _ := service.OutletRepository.GetImportInstructions("outlet")
	mandatoryHeaders := buildMandatoryOutletHeaderSet(instructions)

	var headers []string
	if len(optionalHeaders) > 0 {
		headers = optionalHeaders
	} else {
		headers = outletTemplateMandatoryOnlyColumns
		mandatoryHeaders = make(map[string]struct{}, len(outletTemplateMandatoryOnlyColumns))
		for _, h := range outletTemplateMandatoryOnlyColumns {
			mandatoryHeaders[h] = struct{}{}
		}
	}
	widths := make([]int, len(headers))

	styleCache := map[string]int{}
	for i, header := range headers {
		displayName := toDisplayHeader(header)
		if isMandatoryOutletHeader(header, displayName, mandatoryHeaders) {
			displayName = displayName + "*"
		}
		if maxLen, ok := outletMaxCharHeaders[canonicalizeHeader(header)]; ok && maxLen > 0 {
			displayName = fmt.Sprintf("%s (maksimal %d karakter)", displayName, maxLen)
		}
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, displayName)
		step := stepCategoryForHeader(header)
		styleID, err := headerStyleID(f, styleCache, step)
		if err != nil {
			return nil, err
		}
		f.SetCellStyle(sheetName, cell, cell, styleID)
		widths[i] = utf8.RuneCountInString(displayName)
	}

	// Sample data dihapus sesuai permintaan; hanya header yang disediakan.

	insSheet := "Instructions"
	f.NewSheet(insSheet)
	ih := []string{"Column Name (Template)", "Instruction"}
	for i, h := range ih {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(insSheet, cell, h)
		styleID, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}})
		f.SetCellStyle(insSheet, cell, cell, styleID)
	}
	rowIdx := 2
	for _, header := range outletInstructionColumnOrder {
		if header == "is_active" {
			continue
		}
		kolom := toDisplayHeader(header)
		keterangan := outletInstructionByHeader[header]
		if keterangan == "" {
			for _, ins := range instructions {
				if canonicalizeHeader(ins.Kolom) == canonicalizeHeader(kolom) {
					keterangan = ins.Keterangan
					break
				}
			}
		}
		f.SetCellValue(insSheet, "A"+strconv.Itoa(rowIdx), kolom)
		f.SetCellValue(insSheet, "B"+strconv.Itoa(rowIdx), keterangan)
		rowIdx++
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	for i := range headers {
		if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
			width := float64(widths[i] + 4)
			if width < 18 {
				width = 18
			}
			if width > 60 {
				width = 60
			}
			_ = f.SetColWidth(sheetName, colName, colName, width)
		}
	}
	return f.WriteToBuffer()
}

func (service *outletServiceImpl) createOutletTemplateXLS(optionalHeaders []string) (*bytes.Buffer, error) {
	// Logika identik dengan XLSX
	return service.createOutletTemplateXLSX(optionalHeaders)
}

func (service *outletServiceImpl) createOutletTemplateCSV(optionalHeaders []string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'

	var headers []string
	if len(optionalHeaders) > 0 {
		headers = optionalHeaders
	} else {
		headers = outletTemplateMandatoryOnlyColumns
	}

	// Write display headers matching XLSX labels
	disp := make([]string, len(headers))
	for i, header := range headers {
		displayName := toDisplayHeader(header)
		if maxLen, ok := outletMaxCharHeaders[canonicalizeHeader(header)]; ok && maxLen > 0 {
			displayName = fmt.Sprintf("%s (maksimal %d karakter)", displayName, maxLen)
		}
		disp[i] = displayName
	}
	_ = writer.Write(disp)

	// Sample data dihapus; hanya header.

	writer.Flush()
	return buffer, writer.Error()
}

// --- FUNGSI EKSPOR TEMPLATE UPDATE ---

// FieldGroupsOutlet maps each header/field to its logical group for ExportTemplateUpdate.
// Groups are aligned with the web step layout to make selection easier to reason about.
var FieldGroupsOutlet = map[string]string{
	"outlet_code":               "step1_basic",
	"outlet_name":               "step1_basic",
	"is_active":                 "step1_basic",
	"outlet_status":             "step1_basic",
	"address1":                  "step1_basic",
	"outlet_province":           "step1_basic",
	"outlet_regency":            "step1_basic",
	"outlet_sub_district":       "step1_basic",
	"outlet_ward":               "step1_basic",
	"zip_code":                  "step1_basic",
	"ot_loc_name":               "step1_basic",
	"longitude":                 "step1_basic",
	"latitude":                  "step1_basic",
	"building_own":              "step1_basic",
	"outlet_establishment_date": "step1_basic",
	"is_pkp_outlet":             "step1_basic",

	"phone_no": "step1_optional",
	"fax_no":   "step1_optional",
	"barcode":  "step1_optional",

	"contact_name":      "step1_contact",
	"job_title":         "step1_contact",
	"identity_type":     "step1_contact",
	"identity_no":       "step1_contact",
	"contact_phone_no":  "step1_contact",
	"contact_is_wa_no":  "step1_contact",
	"contact_wa_no":     "step1_contact",
	"contact_email":     "step1_contact",
	"outlet_contact_id": "step1_contact",

	"tax_invoice_form_name": "step1_tax",
	"tax_identifier_type":   "step1_tax",
	"tax_identifier_no":     "step1_tax",
	"nitku":                 "step1_tax",
	"tax_name":              "step1_tax",
	"address_tax":           "step1_tax",
	"outlet_tax_id":         "step1_tax",

	"disc_grp_name":  "step2",
	"ot_grp_name":    "step2",
	"is_contra_bon":  "step2",
	"price_grp_name": "step2",
	"district_name":  "step2",
	"industry_name":  "step2",
	"ot_class_name":  "step2",
	"ot_type_name":   "step2",
	"market_name":    "step2",
	"agent_from":     "step2_optional",

	"delv_addr1":        "step3_delivery",
	"delv_province":     "step3_delivery",
	"delv_regency":      "step3_delivery",
	"delv_sub_district": "step3_delivery",
	"delv_ward":         "step3_delivery",
	"delv_longitude":    "step3_delivery",
	"delv_latitude":     "step3_delivery",
	"delv_zip_code":     "step3_delivery",
	"delv_is_same_addr": "step3_delivery",

	"inv_addr1":        "step3_invoice",
	"inv_province":     "step3_invoice",
	"inv_regency":      "step3_invoice",
	"inv_sub_district": "step3_invoice",
	"inv_ward":         "step3_invoice",
	"inv_zip_code":     "step3_invoice",
	"inv_is_same_addr": "step3_invoice",

	"payment_type_name": "step4_payment",
	"ar_status_name":    "step4_payment",

	"bank_name":      "step4_bank",
	"account_no":     "step4_bank",
	"account_name":   "step4_bank",
	"outlet_bank_id": "step4_bank",

	"top":                         "step4_optional",
	"credit_limit_type_name":      "step4_optional",
	"credit_limit":                "step4_optional",
	"credit_limit_action_name":    "step4_optional",
	"sales_inv_limit_type_name":   "step4_optional",
	"sales_inv_limit":             "step4_optional",
	"sales_inv_limit_action_name": "step4_optional",
	"obs_type_name":               "step4_optional",
	"obs":                         "step4_optional",
	"obs_limit_action_name":       "step4_optional",

	"outlet_id":             "meta",
	"outlet_principal_code": "meta",
	"ot_reg_date":           "meta",
	"closed_date":           "meta",
	"first_trans_date":      "meta",
	"last_trans_date":       "meta",
	"ar_total":              "meta",
}

var (
	outletStep1BasicHeaders = []string{
		"outlet_code", "outlet_name", "is_active", "outlet_status",
		"address1", "outlet_province", "outlet_regency", "outlet_sub_district",
		"outlet_ward", "zip_code", "ot_loc_name", "longitude", "latitude",
		"building_own", "outlet_establishment_date",
	}
	outletStep1ContactHeaders = []string{
		"contact_name", "job_title", "is_pkp_outlet", "identity_type", "identity_no",
		"contact_phone_no", "contact_is_wa_no", "contact_wa_no", "contact_email",
	}
	outletStep1TaxHeaders = []string{
		"tax_invoice_form_name", "tax_identifier_type", "tax_identifier_no",
		"nitku", "tax_name", "address_tax",
	}
	outletStep1OptionalHeaders = []string{"phone_no", "fax_no", "barcode"}
	outletStep2Headers         = []string{"disc_grp_name", "ot_grp_name", "is_contra_bon", "price_grp_name", "district_name", "industry_name", "ot_class_name", "ot_type_name", "market_name"}
	outletStep2OptionalHeaders = []string{"agent_from"}
	outletDeliveryHeaders      = []string{"delv_addr1", "delv_province", "delv_regency", "delv_sub_district", "delv_ward", "delv_longitude", "delv_latitude", "delv_zip_code", "delv_is_same_addr"}
	outletInvoiceHeaders       = []string{"inv_addr1", "inv_province", "inv_regency", "inv_sub_district", "inv_ward", "inv_zip_code", "inv_is_same_addr"}
	outletStep4PaymentHeaders  = []string{"payment_type_name", "ar_status_name"}
	outletStep4BankHeaders     = []string{"bank_name", "account_no", "account_name"}
	outletStep4OptionalHeaders = []string{"top", "credit_limit_type_name", "credit_limit", "credit_limit_action_name", "sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name", "obs_type_name", "obs", "obs_limit_action_name"}

	outletExportHeaders = func() []string {
		headers := []string{}
		headers = append(headers, outletStep1BasicHeaders...)
		headers = append(headers, outletStep1ContactHeaders...)
		headers = append(headers, outletStep1TaxHeaders...)
		headers = append(headers, outletStep1OptionalHeaders...)
		headers = append(headers, outletStep2Headers...)
		headers = append(headers, outletStep2OptionalHeaders...)
		headers = append(headers, outletDeliveryHeaders...)
		headers = append(headers, outletInvoiceHeaders...)
		headers = append(headers, outletStep4PaymentHeaders...)
		headers = append(headers, outletStep4BankHeaders...)
		headers = append(headers, outletStep4OptionalHeaders...)
		return headers
	}()

	outletTemplateMandatoryOnlyColumns = []string{
		"outlet_code", "outlet_name", "outlet_status",
		"contact_name", "contact_phone_no", "is_pkp_outlet", "identity_type", "identity_no",
		"tax_identifier_type", "tax_identifier_no", "nitku", "tax_name", "address_tax",
		"disc_grp_name", "ot_grp_name", "price_grp_name", "ot_class_name", "ot_type_name",
		"address1", "longitude", "latitude",
		"payment_type_name", "top",
	}

	outletTemplateColumnOrder = []string{
		"outlet_code", "outlet_name", "barcode", "outlet_status",
		"contact_name", "job_title", "is_pkp_outlet", "identity_type", "identity_no",
		"contact_phone_no", "contact_wa_no", "contact_email",
		"tax_identifier_type", "tax_identifier_no", "nitku", "tax_name", "address_tax",
		"disc_grp_name", "ot_grp_name", "price_grp_name", "ot_class_name", "ot_type_name",
		"district_name", "agent_from", "industry_name", "market_name", "ot_loc_name",
		"address1", "outlet_province", "outlet_regency", "outlet_sub_district", "outlet_ward",
		"zip_code", "longitude", "latitude",
		"phone_no", "fax_no", "building_own", "outlet_establishment_date",
		"delv_addr1", "delv_province", "delv_regency", "delv_sub_district", "delv_ward",
		"delv_longitude", "delv_latitude", "delv_zip_code",
		"delv_addr2", "delv_province2", "delv_regency2", "delv_sub_district2", "delv_ward2",
		"delv_longitude2", "delv_latitude2", "delv_zip_code2",
		"inv_addr1", "inv_province", "inv_regency", "inv_sub_district", "inv_ward", "inv_zip_code",
		"payment_type_name", "top", "ar_status_name",
		"credit_limit_type_name", "credit_limit", "credit_limit_action_name",
		"sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name",
		"obs_type_name", "obs", "obs_limit_action_name",
		"bank_name", "account_no", "account_name",
	}

	outletExportDataHeaders = concatHeaders(excludeHeaders(outletExportHeaders, "is_active"), []string{"principal_code", "ot_reg_date", "closed_date", "first_trans_date", "last_trans_date", "ar_total"})

	outletInstructionColumnOrder = concatHeaders(outletTemplateColumnOrder, []string{"outlet_principal_code", "ot_reg_date", "closed_date", "first_trans_date", "last_trans_date", "ar_total"})

	defaultMandatoryOutletHeaders = map[string]struct{}{
		"outlet_code":               {},
		"outlet_name":               {},
		"address1":                  {},
		"outlet_province":           {},
		"outlet_regency":            {},
		"outlet_sub_district":       {},
		"outlet_ward":               {},
		"zip_code":                  {},
		"ot_loc_name":               {},
		"longitude":                 {},
		"latitude":                  {},
		"building_own":              {},
		"outlet_establishment_date": {},
		"contact_name":              {},
		"job_title":                 {},
		"identity_type":             {},
		"identity_no":               {},
		"is_pkp_outlet":             {},
		"contact_phone_no":          {},
		"contact_is_wa_no":          {},
		"contact_wa_no":             {},
		"contact_email":             {},
		"tax_invoice_form_name":     {},
		"tax_identifier_type":       {},
		"tax_name":                  {},
		"address_tax":               {},
		"phone_no":                  {},
		"disc_grp_name":             {},
		"ot_grp_name":               {},
		"is_contra_bon":             {},
		"price_grp_name":            {},
		"district_name":             {},
		"industry_name":             {},
		"ot_class_name":             {},
		"ot_type_name":              {},
		"market_name":               {},
		"delv_addr1":                {},
		"delv_province":             {},
		"delv_regency":              {},
		"delv_sub_district":         {},
		"delv_ward":                 {},
		"delv_longitude":            {},
		"delv_latitude":             {},
		"delv_zip_code":             {},
		"delv_is_same_addr":         {},
		"inv_addr1":                 {},
		"inv_province":              {},
		"inv_regency":               {},
		"inv_sub_district":          {},
		"inv_ward":                  {},
		"inv_zip_code":              {},
		"inv_is_same_addr":          {},
		"payment_type_name":         {},
		"ar_status_name":            {},
		"bank_name":                 {},
		"account_no":                {},
		"account_name":              {},
		"top":                       {},
	}

	outletMaxCharHeaders = map[string]int{
		"outlet_name":            75,
		"outlet_code":            10,
		"address1":               150,
		"delv_addr1":             50,
		"inv_addr1":              50,
		"zip_code":               6,
		"delv_zip_code":          6,
		"inv_zip_code":           6,
		"identity_no":            20,
		"barcode":                25,
		"address2":               150,
		"city":                   100,
		"phone_no":               20,
		"wa_no":                  20,
		"fax_no":                 20,
		"email":                  20,
		"agent_from":             50,
		"tax_name":               150,
		"tax_addr1":              150,
		"tax_addr2":              150,
		"tax_city":               100,
		"tax_no":                 30,
		"owner_name":             150,
		"owner_addr1":            150,
		"owner_addr2":            150,
		"owner_city":             100,
		"owner_phone_no":         20,
		"owner_id_no":            50,
		"delv_addr2":             150,
		"delv_city":              100,
		"delv_province2":         50,
		"delv_regency2":          50,
		"delv_sub_district2":     50,
		"delv_ward2":             50,
		"inv_addr2":              150,
		"inv_city":               100,
		"latitude":               50,
		"longitude":              50,
		"image_url":              255,
		"delv_city2":             50,
		"delv_latitude":          50,
		"delv_longitude":         50,
		"delv_latitude2":         50,
		"delv_longitude2":        50,
		"delv_zip_code2":         6,
		"contact_name":           150,
		"job_title":              100,
		"identity_type":          100,
		"contact_phone_no":       20,
		"contact_wa_no":          20,
		"contact_email":          100,
		"tax_invoice_form_name":  100,
		"tax_identifier_type":    100,
		"tax_identifier_no":      30,
		"nitku":                  50,
		"address_tax":            150,
		"disc_grp_name":          150,
		"ot_loc_name":            100,
		"ot_grp_name":            100,
		"price_grp_name":         150,
		"district_name":          150,
		"industry_name":          150,
		"ot_class_name":          150,
		"ot_type_name":           100,
		"market_name":            150,
		"bank_name":              150,
		"account_no":             50,
		"account_name":           150,
		"outlet_province":        50,
		"outlet_regency":         50,
		"outlet_sub_district":    50,
		"outlet_ward":            50,
		"delv_province":          50,
		"delv_regency":           50,
		"delv_sub_district":      50,
		"delv_ward":              50,
		"inv_province":           50,
		"inv_regency":            50,
		"inv_sub_district":       50,
		"inv_ward":               50,
		"outlet_province_id":     10,
		"outlet_regency_id":      10,
		"outlet_sub_district_id": 10,
		"outlet_ward_id":         10,
		"delv_province_id":       10,
		"delv_regency_id":        10,
		"delv_sub_district_id":   10,
		"delv_ward_id":           10,
		"inv_province_id":        10,
		"inv_regency_id":         10,
		"inv_sub_district_id":    10,
		"inv_ward_id":            10,
		"ot_type_code":           10,
		"ot_grp_code":            10,
		"price_grp_code":         10,
		"district_code":          10,
		"disc_grp_code":          10,
		"market_code":            10,
		"industry_code":          10,
	}

	outletFieldDisplayNameID = map[string]string{
		"outlet_code":               "Kode Outlet",
		"outlet_name":               "Nama Outlet",
		"address1":                  "Alamat Outlet",
		"outlet_province":           "Provinsi",
		"outlet_regency":            "Kota/Kabupaten",
		"outlet_sub_district":       "Kecamatan",
		"outlet_ward":               "Kelurahan/Desa",
		"delv_addr1":                "Alamat Pengiriman",
		"delv_province":             "Provinsi Pengiriman",
		"delv_regency":              "Kota/Kabupaten Pengiriman",
		"delv_sub_district":         "Kecamatan Pengiriman",
		"delv_ward":                 "Kelurahan/Desa Pengiriman",
		"inv_addr1":                 "Alamat Penagihan",
		"inv_province":              "Provinsi Penagihan",
		"inv_regency":               "Kota/Kabupaten Penagihan",
		"inv_sub_district":          "Kecamatan Penagihan",
		"inv_ward":                  "Kelurahan/Desa Penagihan",
		"zip_code":                  "Kode Pos",
		"delv_zip_code":             "Kode Pos Pengiriman",
		"inv_zip_code":              "Kode Pos Penagihan",
		"identity_no":               "Nomor Identitas",
		"ot_loc_name":               "Nama Lokasi Outlet",
		"building_own":              "Status Kepemilikan Gedung",
		"longitude":                 "Longitude",
		"latitude":                  "Latitude",
		"outlet_establishment_date": "Tanggal Berdiri Outlet",
		"disc_grp_name":             "Nama Grup Diskon",
		"ot_grp_name":               "Nama Grup Outlet",
		"ot_type_name":              "Jenis Outlet",
		"price_grp_name":            "Nama Grup Harga",
		"district_name":             "Nama Distrik",
		"industry_name":             "Nama Industri",
		"ot_class_name":             "Kelas Outlet",
		"market_name":               "Nama Pasar",
		"is_contra_bon":             "Control Bon",
		"payment_type_name":         "Jenis Pembayaran",
		"ar_status_name":            "Status AR",
		"top":                       "Termin Pembayaran (TOP)",
		"bank_name":                 "Nama Bank",
		"account_no":                "Nomor Rekening",
		"account_name":              "Nama Pemilik Rekening",
		"contact_name":              "Nama Kontak",
		"job_title":                 "Jabatan",
		"identity_type":             "Jenis Identitas",
		"contact_is_wa_no":          "Setel Sebagai WhatsApp",
		"contact_phone_no":          "Nomor Telepon Kontak",
		"contact_wa_no":             "Nomor WhatsApp Kontak",
		"phone_no":                  "Nomor Telepon",
		"wa_no":                     "Nomor WhatsApp",
		"contact_email":             "Email Kontak",
		"tax_invoice_form_name":     "Jenis Faktur Pajak",
		"tax_identifier_type":       "Jenis Pajak",
		"tax_name":                  "Nama Pajak",
		"address_tax":               "Alamat Pajak",
		"delv_longitude":            "Longitude Pengiriman",
		"delv_latitude":             "Latitude Pengiriman",
		"delv_is_same_addr":         "Alamat Pengiriman Sama",
		"inv_is_same_addr":          "Alamat Penagihan Sama",
	}

	outletTemplateGroupSelectors = map[string][]string{
		"step1_basic":    outletStep1BasicHeaders,
		"step1_contact":  outletStep1ContactHeaders,
		"step1_tax":      outletStep1TaxHeaders,
		"step1_optional": outletStep1OptionalHeaders,
		"contact":        outletStep1ContactHeaders,
		"tax":            outletStep1TaxHeaders,
		"step1": concatHeaders(
			outletStep1BasicHeaders,
			outletStep1ContactHeaders,
			outletStep1TaxHeaders,
			outletStep1OptionalHeaders,
		),
		"step2":               concatHeaders(outletStep2Headers, outletStep2OptionalHeaders),
		"step2_mandatory":     outletStep2Headers,
		"step2_optional":      outletStep2OptionalHeaders,
		"step2_optional_only": outletStep2OptionalHeaders,
		"step3_delivery":      outletDeliveryHeaders,
		"step3_invoice":       outletInvoiceHeaders,
		"step3":               concatHeaders(outletDeliveryHeaders, outletInvoiceHeaders),
		"delivery":            outletDeliveryHeaders,
		"invoice":             outletInvoiceHeaders,
		"step4_payment":       outletStep4PaymentHeaders,
		"step4_bank":          outletStep4BankHeaders,
		"step4_optional":      outletStep4OptionalHeaders,
		"payment":             outletStep4PaymentHeaders,
		"bank":                outletStep4BankHeaders,
		"optional": concatHeaders(
			outletStep1OptionalHeaders,
			outletStep2OptionalHeaders,
			outletStep4OptionalHeaders,
		),
		"step4": concatHeaders(
			outletStep4PaymentHeaders,
			outletStep4BankHeaders,
			outletStep4OptionalHeaders,
		),
		"all": outletExportHeaders,
	}

	outletReferenceSelectors = map[string][]string{
		"location":     {"cust_id", "outlet_loc_code", "outlet_loc_name"},
		"type":         {"cust_id", "outlet_type_code", "outlet_type_name"},
		"group":        {"cust_id", "outlet_grp_code", "outlet_grp_name"},
		"district":     {"cust_id", "district_code", "district_name"},
		"province":     {"cust_id", "province_id", "province"},
		"regency":      {"cust_id", "regency_id", "regency", "province_id"},
		"sub_district": {"cust_id", "sub_district_id", "sub_district", "province_id", "regency_id"},
		"ward":         {"cust_id", "ward_id", "ward", "province_id", "regency_id", "sub_district_id"},
		"bank":         {"cust_id", "bank_code", "bank_name"},
		"price":        {"cust_id", "price_grp_code", "price_grp_name"},
		"class":        {"cust_id", "ot_class_code", "ot_class_name"},
		"discount":     {"cust_id", "disc_grp_code", "disc_grp_name"},
		"market":       {"cust_id", "market_code", "market_name"},
		"industry":     {"cust_id", "industry_code", "industry_name"},
	}

	outletReferenceSample = map[string][]string{
		"location":     {"C001", "000", "Non Lokasi"},
		"type":         {"C001", "TK", "TK RETAIL"},
		"group":        {"C001", "000", "Non Group"},
		"district":     {"C001", "360101", "Not Set"},
		"province":     {"C001", "11", "DKI Jakarta"},
		"regency":      {"C001", "1101", "Kota Jakarta Pusat", "11"},
		"sub_district": {"C001", "110101", "Gambir", "11", "1101"},
		"ward":         {"C001", "11010101", "Cideng", "11", "1101", "110101"},
		"bank":         {"C001", "000000", "Nama Bank"},
		"price":        {"C001", "000", "Non Price"},
		"class":        {"C001", "000", "Non Class"},
		"discount":     {"C001", "000", "Non Disc"},
		"market":       {"C001", "000", "Non Market"},
		"industry":     {"C001", "000", "Non Industry"},
	}
	outletReferenceColumns = func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, cols := range outletReferenceSelectors {
			for _, col := range cols {
				m[normalizeHeader(col)] = struct{}{}
			}
		}
		return m
	}()

	additionalOptionToInternalHeaders = map[string][]string{
		"barcode":                             {"barcode"},
		"position":                            {"job_title"},
		"set_as_pkp_outlet":                   {"is_pkp_outlet"},
		"whatsapp_number":                     {"contact_wa_no"},
		"contact_email":                       {"contact_email"},
		"district":                            {"district_name"},
		"agent_from":                          {"agent_from"},
		"industry":                            {"industry_name"},
		"market":                              {"market_name"},
		"outlet_location":                     {"ot_loc_name"},
		"province_outlet_address":             {"outlet_province"},
		"cityregency_outlet_adddress":         {"outlet_regency"},
		"subdistrict_outlet_address":          {"outlet_sub_district"},
		"village_outlet_address":              {"outlet_ward"},
		"postal_code_outlet_address":          {"zip_code"},
		"outlet_phone":                        {"phone_no"},
		"fax_number":                          {"fax_no"},
		"building_ownership":                  {"building_own"},
		"outlet_created_date":                 {"outlet_establishment_date"},
		"delivery_address":                    {"delv_addr1"},
		"delivery_province":                   {"delv_province"},
		"delivery_cityregency":                {"delv_regency"},
		"delivery_subdistrict":                {"delv_sub_district"},
		"delivery_village":                    {"delv_ward"},
		"delivery_longgitude":                 {"delv_longitude"},
		"delivery_lattitude":                  {"delv_latitude"},
		"delivery_postal_code":                {"delv_zip_code"},
		"delivery_2_address":                  {"delv_addr2"},
		"delivery_2_province":                 {"delv_province2"},
		"delivery_2_cityregency":              {"delv_regency2"},
		"delivery_2_subdistrict":              {"delv_sub_district2"},
		"delivery_2_village":                  {"delv_ward2"},
		"delivery_2_longgitude":               {"delv_longitude2"},
		"delivery_2_lattitude":                {"delv_latitude2"},
		"delivery_2_postal_code":              {"delv_zip_code2"},
		"invoice_address":                     {"inv_addr1"},
		"invoice_province":                    {"inv_province"},
		"invoice_cityregency":                 {"inv_regency"},
		"invoice_subdistrict":                 {"inv_sub_district"},
		"invoice_village":                     {"inv_ward"},
		"invoice_postal_code":                 {"inv_zip_code"},
		"ar_status":                           {"ar_status_name"},
		"sales_order_validation_credit_limit": {"credit_limit_type_name", "credit_limit", "credit_limit_action_name"},
		"sales_order_validation_overdue":      {"sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name"},
		"sales_order_valiation_overdue":       {"sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name"},
		"sales_order_validation_outstanding":  {"obs_type_name", "obs", "obs_limit_action_name"},
		"bank":                                {"bank_name", "account_no", "account_name"},
	}

	outletTemplateAllowedHeaders = func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, h := range outletExportHeaders {
			m[h] = struct{}{}
		}
		for _, h := range []string{"is_pkp_outlet", "delv_addr2", "delv_city2", "delv_province2", "delv_regency2", "delv_sub_district2", "delv_ward2", "delv_latitude2", "delv_longitude2", "delv_zip_code2"} {
			m[h] = struct{}{}
		}
		return m
	}()

	headerStepFillColors = map[string]string{
		"step1":     "C6EFCE", // soft green
		"step2":     "FFEB9C", // soft yellow
		"step3":     "BDD7EE", // soft blue
		"step4":     "F7C6D9", // soft Red
		"reference": "D9E1F2", // soft violet-blue for lookup sheets
		"meta":      "E6E6E6", // neutral grey
		"default":   "CCCCCC", // fallback
	}
	stepColorLabels = map[string]string{
		"step1":     "Hijau Muda",
		"step2":     "Kuning Muda",
		"step3":     "Biru Muda",
		"step4":     "Merah Muda",
		"reference": "Ungu Muda",
		"meta":      "Abu-abu Muda",
		"default":   "",
	}
)

const (
	instructionMandatoryYes = "Kolom Wajib Di Isi"
	instructionMandatoryNo  = "Kolom Tidak Wajib Di Isi"
)

var outletInstructionByHeader = map[string]string{
	"outlet_code":                 "alphanumeric with special character, maximum 10 char, Unique Value",
	"outlet_name":                 "alphanumeric with special character, maximum 75 char",
	"barcode":                     "Alphanumeric, maximum 25 char",
	"outlet_status":               "1 = New Open Outlet (noo), 4 = Close, 5 = Dormant, 6 = Registered, 7 = Active",
	"contact_name":                "Alphanumeric, maximum 150 Char",
	"job_title":                   "Alphanumeric, maximum 100 Char",
	"is_pkp_outlet":               "true/false",
	"identity_type":               "NATIONAL ID / PASSPORT / OTHER ID",
	"identity_no":                 "Alphanumeric, maximum 20 char",
	"contact_phone_no":            "Numeric, maximum 20 Char",
	"contact_wa_no":               "Numeric, maximum 20 Char",
	"contact_email":               "Alphanumeric with special character, maximum 100 Char",
	"tax_identifier_type":         "NATIONAL ID (NIK) / PASSPORT / TIN / OTHER ID",
	"tax_identifier_no":           "Numeric, maximum 30 Char",
	"nitku":                       "Numeric, maximum 50 Char",
	"tax_name":                    "Alphanumeric, maximum 150 Char",
	"address_tax":                 "Alphanumeric, maximum 150 Char",
	"disc_grp_name":               "Wajib ada di Discount Group Name di Setup Parameter Web > Outlet",
	"ot_grp_name":                 "Wajib ada di Outlet Group Name di Setup Parameter Web > Outlet",
	"price_grp_name":              "Wajib ada di Price Group Name di Setup Parameter Web > Outlet",
	"ot_class_name":               "Wajib ada di Outlet Classification Name di Setup Parameter Web > Outlet",
	"ot_type_name":                "Wajib ada di Outlet Type Name di Setup Parameter Web > Outlet",
	"district_name":               "Wajib ada di District Name di Setup Parameter Web > Outlet",
	"agent_from":                  "Alphanumeric, maximum 50 Char",
	"industry_name":               "Wajib ada di Industry Name di Setup Parameter Web > Outlet",
	"market_name":                 "Wajib ada di Market Name di Setup Parameter Web > Outlet",
	"ot_loc_name":                 "Wajib ada di Outlet Location Name di Setup Parameter Web > Outlet",
	"address1":                    "alphanumeric with special character, maximum 150 char",
	"outlet_province":             "Wajib ada di Province Name di Setup Parameter Web > Distributor > Province",
	"outlet_regency":              "Wajib ada di City/Regency Name di Setup Parameter Web > Distributor > City/Regency",
	"outlet_sub_district":         "Wajib ada di Subdistrict Name di Setup Parameter Web > Distributor > Subdistrict",
	"outlet_ward":                 "Wajib ada di Village Name di Setup Parameter Web > Distributor > Village",
	"zip_code":                    "Numeric, maximum 6 Char",
	"longitude":                   "Alphanumeric with special Char, maximum 50 char",
	"latitude":                    "Alphanumeric with special Char, maximum 50 char",
	"phone_no":                    "Numeric, maximum 20 Char",
	"fax_no":                      "Numeric, maximum 20 Char",
	"building_own":                "Milik Sendiri / Kontrak / Sewa",
	"outlet_establishment_date":   "format tanggal DD-MM-YYYY (02-01-2026)",
	"delv_addr1":                  "alphanumeric with special character, maximum 150 char",
	"delv_province":               "Wajib ada di Province Name di Setup Parameter Web > Distributor > Province",
	"delv_regency":                "Wajib ada di City/Regency Name di Setup Parameter Web > Distributor > City/Regency",
	"delv_sub_district":           "Wajib ada di Subdistrict Name di Setup Parameter Web > Distributor > Subdistrict",
	"delv_ward":                   "Wajib ada di Village Name di Setup Parameter Web > Distributor > Village",
	"delv_longitude":              "Alphanumeric with special Char, maximum 50 char",
	"delv_latitude":               "Alphanumeric with special Char, maximum 50 char",
	"delv_zip_code":               "Numeric, maximum 6 Char",
	"delv_addr2":                  "alphanumeric with special character, maximum 150 char",
	"delv_province2":              "Wajib ada di Province Name di Setup Parameter Web > Distributor > Province",
	"delv_regency2":               "Wajib ada di City/Regency Name di Setup Parameter Web > Distributor > City/Regency",
	"delv_sub_district2":          "Wajib ada di Subdistrict Name di Setup Parameter Web > Distributor > Subdistrict",
	"delv_ward2":                  "Wajib ada di Village Name di Setup Parameter Web > Distributor > Village",
	"delv_longitude2":             "Alphanumeric with special Char, maximum 50 char",
	"delv_latitude2":              "Alphanumeric with special Char, maximum 50 char",
	"delv_zip_code2":              "Numeric, maximum 6 Char",
	"inv_addr1":                   "alphanumeric with special character, maximum 150 char",
	"inv_province":                "Wajib ada di Province Name di Setup Parameter Web > Distributor > Province",
	"inv_regency":                 "Wajib ada di City/Regency Name di Setup Parameter Web > Distributor > City/Regency",
	"inv_sub_district":            "Wajib ada di Subdistrict Name di Setup Parameter Web > Distributor > Subdistrict",
	"inv_ward":                    "Wajib ada di Village Name di Setup Parameter Web > Distributor > Village",
	"inv_zip_code":                "Numeric, maximum 6 Char",
	"payment_type_name":           "- 1 = Cash On Delivery\n- 2 = Cash Before Delivery\n- 3 = Credit",
	"top":                         "Numeric, Maximum 3 Char",
	"ar_status_name":              "- 1 = Normal\n- 2 = 1x Giro Tolakan\n- 3 = 2x Giro Tolakan\n- 4 = 3x Giro Tolakan",
	"credit_limit_type_name":      "- Null = Unlimited\n- 2 = Limited by Total",
	"credit_limit":                "Numeric",
	"credit_limit_action_name":    "- 1 = Warning\n- 2 = Restricted",
	"sales_inv_limit_type_name":   "- Null = Unlimited\n- 2 = Limited by Invoice",
	"sales_inv_limit":             "Numeric",
	"sales_inv_limit_action_name": "- 1 = Warning\n- 2 = Restricted",
	"obs_type_name":               "- Null = Unlimited\n- 2 = Limited by Total",
	"obs":                         "Numeric",
	"obs_limit_action_name":       "- 1 = Warning\n- 2 = Restricted",
	"bank_name":                   "Wajib ada di Setup Parameter Web > Bank",
	"account_no":                  "Alphanumeric with special char, maximum 30 char",
	"account_name":                "Alphanumeric with special char, maximum 30 char",
	"outlet_principal_code":       "Hanya tampil ketika melakukan eksport data outlet",
	"ot_reg_date":                 "Hanya tampil ketika melakukan eksport data outlet",
	"closed_date":                 "Hanya tampil ketika melakukan eksport data outlet",
	"first_trans_date":            "Hanya tampil ketika melakukan eksport data outlet",
	"last_trans_date":             "Hanya tampil ketika melakukan eksport data outlet",
	"ar_total":                    "Hanya tampil ketika melakukan eksport data outlet",
}

func concatHeaders(groups ...[]string) []string {
	total := 0
	for _, g := range groups {
		total += len(g)
	}
	result := make([]string, 0, total)
	for _, g := range groups {
		result = append(result, g...)
	}
	return result
}

func excludeHeaders(headers []string, excluded ...string) []string {
	if len(headers) == 0 || len(excluded) == 0 {
		return headers
	}
	skip := make(map[string]struct{}, len(excluded))
	for _, h := range excluded {
		skip[h] = struct{}{}
	}
	result := make([]string, 0, len(headers))
	for _, h := range headers {
		if _, ok := skip[h]; ok {
			continue
		}
		result = append(result, h)
	}
	return result
}

func headersInOrder(orderedList []string, allowed map[string]struct{}) []string {
	out := make([]string, 0, len(orderedList))
	for _, h := range orderedList {
		if _, ok := allowed[h]; ok {
			out = append(out, h)
		}
	}
	return out
}

func resolveAdditionalOptions(additional []string) []string {
	seen := make(map[string]struct{})
	out := make([]string, 0)
	for _, opt := range additional {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		norm := normalizeHeader(opt)
		headers, ok := additionalOptionToInternalHeaders[norm]
		if !ok {
			continue
		}
		for _, h := range headers {
			if _, allowed := outletTemplateAllowedHeaders[h]; !allowed {
				continue
			}
			if _, ok := seen[h]; ok {
				continue
			}
			seen[h] = struct{}{}
			out = append(out, h)
		}
	}
	return out
}

func stepCategoryForHeader(header string) string {
	key := normalizeHeader(header)
	if grp, ok := FieldGroupsOutlet[key]; ok {
		switch grp {
		case "step1_basic", "step1_contact", "step1_tax", "step1_optional":
			return "step1"
		case "step2", "step2_optional":
			return "step2"
		case "step3_delivery", "step3_invoice":
			return "step3"
		case "step4_payment", "step4_bank", "step4_optional":
			return "step4"
		case "meta":
			return "meta"
		}
	}
	if _, ok := outletReferenceColumns[key]; ok {
		return "reference"
	}
	return "default"
}

func headerStyleID(f *excelize.File, cache map[string]int, step string) (int, error) {
	if id, ok := cache[step]; ok {
		return id, nil
	}
	color, ok := headerStepFillColors[step]
	if !ok {
		color = headerStepFillColors["default"]
	}
	id, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
	})
	if err != nil {
		return 0, err
	}
	cache[step] = id
	return id, nil
}

func reorderSelectedHeaders(selected []string, appearance []string) []string {
	if len(selected) <= 1 {
		return selected
	}
	remaining := make(map[string]struct{}, len(selected))
	for _, h := range selected {
		remaining[h] = struct{}{}
	}

	ordered := make([]string, 0, len(selected))
	add := func(list []string) {
		for _, h := range list {
			if _, ok := remaining[h]; ok {
				ordered = append(ordered, h)
				delete(remaining, h)
			}
		}
	}

	add(outletStep1BasicHeaders)
	add(outletStep1ContactHeaders)
	add(outletStep1TaxHeaders)
	add(outletStep1OptionalHeaders)
	add(outletStep2Headers)
	add(outletStep2OptionalHeaders)
	add(outletDeliveryHeaders)
	add(outletInvoiceHeaders)
	add(outletStep4PaymentHeaders)
	add(outletStep4BankHeaders)
	add(outletStep4OptionalHeaders)
	add([]string{"outlet_id"})

	if len(remaining) > 0 {
		firstSeen := make(map[string]int, len(appearance))
		for idx, h := range appearance {
			if _, ok := firstSeen[h]; !ok {
				firstSeen[h] = idx
			}
		}
		remainder := make([]string, 0, len(remaining))
		for h := range remaining {
			remainder = append(remainder, h)
		}
		sort.SliceStable(remainder, func(i, j int) bool {
			return firstSeen[remainder[i]] < firstSeen[remainder[j]]
		})
		ordered = append(ordered, remainder...)
	}

	return ordered
}

func instructionStepAndColor(kolom string, step string) (string, string) {
	step = strings.TrimSpace(step)
	stepKey := normalizeStepKey(step)
	stepLabel := step
	if stepLabel == "" {
		canonical := canonicalizeHeader(kolom)
		stepKey = stepCategoryForHeader(canonical)
		stepLabel = stepDisplayName(stepKey)
	} else if stepKey == "" {
		stepKey = normalizeStepKey(stepLabel)
	}

	colorLabel := ""
	if stepKey != "" {
		if label, ok := stepColorLabels[stepKey]; ok && label != "" {
			colorLabel = label
		} else if hex, ok := headerStepFillColors[stepKey]; ok && hex != "" {
			colorLabel = "#" + strings.ToUpper(hex)
		}
	}

	return stepLabel, colorLabel
}

func stepDisplayName(stepKey string) string {
	switch stepKey {
	case "step1":
		return "Step 1"
	case "step2":
		return "Step 2"
	case "step3":
		return "Step 3"
	case "step4":
		return "Step 4"
	case "reference":
		return "Reference"
	case "meta":
		return "Metadata"
	default:
		return ""
	}
}

func normalizeStepKey(step string) string {
	if step == "" {
		return ""
	}
	processed := strings.ToLower(strings.TrimSpace(step))
	replacer := strings.NewReplacer(" ", "", "-", "", "_", "")
	processed = replacer.Replace(processed)
	switch processed {
	case "step1", "step01", "langkah1":
		return "step1"
	case "step2", "step02", "langkah2":
		return "step2"
	case "step3", "step03", "langkah3":
		return "step3"
	case "step4", "step04", "langkah4":
		return "step4"
	case "reference", "referensi":
		return "reference"
	case "meta", "metadata":
		return "meta"
	default:
		if strings.HasPrefix(processed, "step") {
			digits := strings.TrimPrefix(processed, "step")
			switch digits {
			case "1", "01":
				return "step1"
			case "2", "02":
				return "step2"
			case "3", "03":
				return "step3"
			case "4", "04":
				return "step4"
			}
		}
	}
	return ""
}

// resolveFieldsToOutletTemplateHeaders mengubah fields (nama tampilan/alias) menjadi daftar header internal
// dengan logika sama seperti ExportTemplateUpdate: validation five + field/group/reference.
func resolveFieldsToOutletTemplateHeaders(fields []string) []string {
	headers := []string{}
	appearance := []string{}

	headerAlias := map[string]string{}
	for key := range FieldGroupsOutlet {
		headerAlias[normalizeHeader(key)] = key
		headerAlias[normalizeHeader(toDisplayHeader(key))] = key
	}
	groupAlias := map[string]string{}
	for key := range outletTemplateGroupSelectors {
		norm := normalizeHeader(key)
		groupAlias[norm] = key
		groupAlias[strings.ReplaceAll(norm, "_", "")] = key
	}
	legacyGroupMap := map[string]string{
		"outlet":         "all",
		"outlet_contact": "step1_contact",
		"outlet_bank":    "step4_bank",
		"outlet_tax":     "step1_tax",
	}
	for legacy, target := range legacyGroupMap {
		norm := normalizeHeader(legacy)
		groupAlias[norm] = target
		groupAlias[strings.ReplaceAll(norm, "_", "")] = target
	}
	referenceAlias := map[string]string{}
	for key := range outletReferenceSelectors {
		norm := normalizeHeader(key)
		referenceAlias[norm] = key
		referenceAlias[strings.ReplaceAll(norm, "_", "")] = key
	}

	hasValidationCredit := false
	hasValidationOverdue := false
	hasValidationOutstanding := false

	for _, raw := range fields {
		token := normalizeHeader(raw)
		if token == "" {
			continue
		}
		switch token {
		case "sales_order_validation_credit_limit":
			hasValidationCredit = true
		case "sales_order_validation_overdue", "sales_order_valiation_overdue":
			hasValidationOverdue = true
		case "sales_order_validation_outstanding":
			hasValidationOutstanding = true
		}
		if canonical, ok := headerAlias[token]; ok {
			headers = append(headers, canonical)
			appearance = append(appearance, canonical)
			continue
		}
		if key, ok := groupAlias[token]; ok {
			groupHeaders := outletTemplateGroupSelectors[key]
			headers = append(headers, groupHeaders...)
			appearance = append(appearance, groupHeaders...)
			continue
		}
		if key, ok := groupAlias[strings.ReplaceAll(token, "_", "")]; ok {
			groupHeaders := outletTemplateGroupSelectors[key]
			headers = append(headers, groupHeaders...)
			appearance = append(appearance, groupHeaders...)
			continue
		}
		if key, ok := referenceAlias[token]; ok {
			referenceHeaders := outletReferenceSelectors[key]
			headers = append(headers, referenceHeaders...)
			appearance = append(appearance, referenceHeaders...)
			continue
		}
		if key, ok := referenceAlias[strings.ReplaceAll(token, "_", "")]; ok {
			referenceHeaders := outletReferenceSelectors[key]
			headers = append(headers, referenceHeaders...)
			appearance = append(appearance, referenceHeaders...)
			continue
		}
	}

	if hasValidationCredit || hasValidationOverdue || hasValidationOutstanding {
		var validationFive []string
		switch {
		case hasValidationCredit:
			validationFive = []string{"outlet_code", "outlet_name", "credit_limit_type_name", "credit_limit", "credit_limit_action_name"}
		case hasValidationOverdue:
			validationFive = []string{"outlet_code", "outlet_name", "sales_inv_limit_type_name", "sales_inv_limit", "sales_inv_limit_action_name"}
		case hasValidationOutstanding:
			validationFive = []string{"outlet_code", "outlet_name", "obs_type_name", "obs", "obs_limit_action_name"}
		}
		validationSet := make(map[string]struct{}, len(validationFive))
		for _, h := range validationFive {
			validationSet[h] = struct{}{}
		}
		otherHeaders := make([]string, 0, len(headers))
		for _, h := range headers {
			if _, in := validationSet[h]; !in {
				otherHeaders = append(otherHeaders, h)
			}
		}
		headers = append(append([]string{}, validationFive...), otherHeaders...)
		otherAppearance := make([]string, 0, len(appearance))
		for _, a := range appearance {
			if _, in := validationSet[a]; !in {
				otherAppearance = append(otherAppearance, a)
			}
		}
		appearance = append(append([]string{}, validationFive...), otherAppearance...)
	}

	// Jika TIDAK ada Sales Order Validation yang dipilih,
	// selalu sertakan header wajib dasar (outletTemplateMandatoryOnlyColumns)
	// lalu tambahkan header dari fields (tanpa duplikasi).
	if !hasValidationCredit && !hasValidationOverdue && !hasValidationOutstanding {
		if len(headers) == 0 {
			return outletTemplateMandatoryOnlyColumns
		}
		combined := make([]string, 0, len(outletTemplateMandatoryOnlyColumns)+len(headers))
		combined = append(combined, outletTemplateMandatoryOnlyColumns...)
		combined = append(combined, headers...)
		seenBase := map[string]struct{}{}
		dedup := make([]string, 0, len(combined))
		for _, h := range combined {
			if _, ok := seenBase[h]; ok {
				continue
			}
			seenBase[h] = struct{}{}
			dedup = append(dedup, h)
		}
		headers = dedup
	}

	if len(headers) == 0 {
		return outletTemplateMandatoryOnlyColumns
	}
	if len(headers) > 1 {
		seen := map[string]struct{}{}
		uniq := headers[:0]
		for _, h := range headers {
			if _, ok := seen[h]; ok {
				continue
			}
			seen[h] = struct{}{}
			uniq = append(uniq, h)
		}
		headers = uniq
	}
	return reorderSelectedHeaders(headers, appearance)
}

func (service *outletServiceImpl) ExportTemplateUpdate(custId string, format string, fields []string) (*bytes.Buffer, string, string, error) {
	headers := []string{}
	appearance := []string{}
	sampleRow := []string{}

	headerAlias := map[string]string{}
	for key := range FieldGroupsOutlet {
		headerAlias[normalizeHeader(key)] = key
		headerAlias[normalizeHeader(toDisplayHeader(key))] = key
	}

	groupAlias := map[string]string{}
	for key := range outletTemplateGroupSelectors {
		norm := normalizeHeader(key)
		groupAlias[norm] = key
		groupAlias[strings.ReplaceAll(norm, "_", "")] = key
	}
	legacyGroupMap := map[string]string{
		"outlet":         "all",
		"outlet_contact": "step1_contact",
		"outlet_bank":    "step4_bank",
		"outlet_tax":     "step1_tax",
	}
	for legacy, target := range legacyGroupMap {
		norm := normalizeHeader(legacy)
		groupAlias[norm] = target
		groupAlias[strings.ReplaceAll(norm, "_", "")] = target
	}

	referenceAlias := map[string]string{}
	for key := range outletReferenceSelectors {
		norm := normalizeHeader(key)
		referenceAlias[norm] = key
		referenceAlias[strings.ReplaceAll(norm, "_", "")] = key
	}

	repoFieldsSet := map[string]struct{}{"outlet": {}}
	referenceTokens := map[string]struct{}{}

	hasValidationCredit := false
	hasValidationOverdue := false
	hasValidationOutstanding := false

	for _, raw := range fields {
		// logrus.Info("Processing field token: %s", raw)
		token := normalizeHeader(raw)
		// logrus.Info("Normalized token: %s", token)
		if token == "" {
			continue
		}

		switch token {
		case "sales_order_validation_credit_limit":
			hasValidationCredit = true
		case "sales_order_validation_overdue", "sales_order_valiation_overdue":
			hasValidationOverdue = true
		case "sales_order_validation_outstanding":
			hasValidationOutstanding = true
		}

		if canonical, ok := headerAlias[token]; ok {
			headers = append(headers, canonical)
			logrus.Infof("Added header: %s", canonical)
			appearance = append(appearance, canonical)
			continue
		}

		if key, ok := groupAlias[token]; ok {
			groupHeaders := outletTemplateGroupSelectors[key]
			headers = append(headers, groupHeaders...)
			appearance = append(appearance, groupHeaders...)
			continue
		}
		if key, ok := groupAlias[strings.ReplaceAll(token, "_", "")]; ok {
			groupHeaders := outletTemplateGroupSelectors[key]
			headers = append(headers, groupHeaders...)
			appearance = append(appearance, groupHeaders...)
			continue
		}

		if key, ok := referenceAlias[token]; ok {
			referenceHeaders := outletReferenceSelectors[key]
			headers = append(headers, referenceHeaders...)
			appearance = append(appearance, referenceHeaders...)
			if sample := outletReferenceSample[key]; sample != nil {
				sampleRow = append(sampleRow, sample...)
			}
			referenceTokens[key] = struct{}{}
			continue
		}
		if key, ok := referenceAlias[strings.ReplaceAll(token, "_", "")]; ok {
			referenceHeaders := outletReferenceSelectors[key]
			headers = append(headers, referenceHeaders...)
			appearance = append(appearance, referenceHeaders...)
			if sample := outletReferenceSample[key]; sample != nil {
				sampleRow = append(sampleRow, sample...)
			}
			referenceTokens[key] = struct{}{}
			continue
		}

	}

	if hasValidationCredit || hasValidationOverdue || hasValidationOutstanding {
		var validationFive []string
		switch {
		case hasValidationCredit:
			validationFive = []string{
				"outlet_code",
				"outlet_name",
				"credit_limit_type_name",
				"credit_limit",
				"credit_limit_action_name",
			}
		case hasValidationOverdue:
			validationFive = []string{
				"outlet_code",
				"outlet_name",
				"sales_inv_limit_type_name",
				"sales_inv_limit",
				"sales_inv_limit_action_name",
			}
		case hasValidationOutstanding:
			validationFive = []string{
				"outlet_code",
				"outlet_name",
				"obs_type_name",
				"obs",
				"obs_limit_action_name",
			}
		}
		validationSet := make(map[string]struct{}, len(validationFive))
		for _, h := range validationFive {
			validationSet[h] = struct{}{}
		}
		otherHeaders := make([]string, 0, len(headers))
		for _, h := range headers {
			if _, in := validationSet[h]; !in {
				otherHeaders = append(otherHeaders, h)
			}
		}
		headers = append(append([]string{}, validationFive...), otherHeaders...)
		otherAppearance := make([]string, 0, len(appearance))
		for _, a := range appearance {
			if _, in := validationSet[a]; !in {
				otherAppearance = append(otherAppearance, a)
			}
		}
		appearance = append(append([]string{}, validationFive...), otherAppearance...)
	}

	if len(headers) == 0 {
		headers = append(headers, outletExportHeaders...)
		appearance = append(appearance, outletExportHeaders...)
	}

	if len(headers) > 1 {
		seen := map[string]struct{}{}
		uniq := headers[:0]
		for _, h := range headers {
			if _, ok := seen[h]; ok {
				continue
			}
			seen[h] = struct{}{}
			uniq = append(uniq, h)
		}
		headers = uniq
	}
	headers = reorderSelectedHeaders(headers, appearance)

	hasHeader := func(target string) bool {
		for _, h := range headers {
			if strings.EqualFold(h, target) {
				return true
			}
		}
		return false
	}

	needContact := false
	needBank := false
	needTax := false
	needOutletID := false

	for _, h := range headers {
		switch FieldGroupsOutlet[h] {
		case "step1_contact":
			needContact = true
			needOutletID = true
		case "step4_bank":
			needBank = true
			needOutletID = true
		case "step1_tax":
			if h == "tax_identifier_type" || h == "tax_identifier_no" || h == "nitku" || h == "address_tax" {
				needTax = true
			}
			needOutletID = true
		case "step1_basic", "step2", "step2_optional", "step3_delivery", "step3_invoice", "step4_payment", "step4_optional":
			needOutletID = true
		}
	}

	// Explicitly remove outlet_id if it was requested or added earlier
	if hasHeader("outlet_id") {
		filtered := headers[:0]
		for _, h := range headers {
			if !strings.EqualFold(h, "outlet_id") {
				filtered = append(filtered, h)
			}
		}
		headers = filtered
	}

	if needOutletID {
		if !hasHeader("outlet_code") {
			headers = append(headers, "outlet_code")
		}
		if !hasHeader("outlet_name") {
			headers = append(headers, "outlet_name")
		}
	}

	if len(headers) > 0 {
		var lead, rest []string
		for _, h := range headers {
			if strings.EqualFold(h, "outlet_code") || strings.EqualFold(h, "outlet_name") {
				lead = append(lead, h)
			} else {
				rest = append(rest, h)
			}
		}
		var ordered []string
		for _, w := range []string{"outlet_code", "outlet_name"} {
			for _, l := range lead {
				if strings.EqualFold(l, w) {
					ordered = append(ordered, l)
					break
				}
			}
		}
		headers = append(ordered, rest...)
	}

	if needContact {
		repoFieldsSet["outlet_contact"] = struct{}{}
	}
	if needBank {
		repoFieldsSet["outlet_bank"] = struct{}{}
	}
	if needTax {
		repoFieldsSet["outlet_tax"] = struct{}{}
	}

	for token := range referenceTokens {
		repoFieldsSet[token] = struct{}{}
	}

	/*
		repoFields := make([]string, 0, len(repoFieldsSet))
		repoFields = append(repoFields, "outlet")
		for token := range repoFieldsSet {
			if token == "outlet" {
				continue
			}
			repoFields = append(repoFields, token)
		}

		data, err := service.OutletRepository.GetOutletDataForTemplateUpdate(custId, repoFields)
		if err != nil {
			return nil, "", "", err
		}

		firstStr := func(row []string) string {
			if len(row) > 0 {
				return row[0]
			}
			return ""
		}
		codeToLabel := func(code string, labels map[int]string) string {
			code = strings.TrimSpace(code)
			if code == "" {
				return ""
			}
			if iv, err := strconv.Atoi(code); err == nil {
				if lbl, ok := labels[iv]; ok {
					return lbl
				}
			}
			return code
		}

		fillPairs := []struct {
			codeKey string
			nameKey string
			labels  map[int]string
		}{
			{"payment_type", "payment_type_name", paymentTypeLabels},
			{"credit_limit_type", "credit_limit_type_name", creditLimitTypeLabels},
			{"sales_inv_limit_type", "sales_inv_limit_type_name", salesInvLimitTypeLabels},
			{"obs_type", "obs_type_name", obsTypeLabels},
			{"credit_limit_action", "credit_limit_action_name", limitActionLabels},
			{"sales_inv_limit_action", "sales_inv_limit_action_name", limitActionLabels},
			{"obs_limit_action", "obs_limit_action_name", limitActionLabels},
			{"tax_invoice_form", "tax_invoice_form_name", taxInvoiceFormLabels},
			{"ar_status", "ar_status_name", arStatusLabels},
		}
		for _, p := range fillPairs {
			if col, ok := data[p.codeKey]; ok {
				out := make([][]string, len(col))
				for i := range col {
					out[i] = []string{codeToLabel(firstStr(col[i]), p.labels)}
				}
				data[p.nameKey] = out
			}
		}
	*/
	data := map[string][][]string{}

	var buffer *bytes.Buffer
	var err error
	var contentType, filename string

	switch strings.ToLower(format) {
	case "csv":
		// Build CSV content first
		csvBuf, cerr := service.createTemplateUpdateCSV(custId, headers, data, sampleRow)
		if cerr != nil {
			return nil, "", "", cerr
		}
		// Package as ZIP together with instructions.csv
		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)
		// 1) template csv
		cf, zerr := zw.Create("outlet_template_update.csv")
		if zerr != nil {
			_ = zw.Close()
			return nil, "", "", zerr
		}
		if _, zerr = cf.Write(csvBuf.Bytes()); zerr != nil {
			_ = zw.Close()
			return nil, "", "", zerr
		}
		inf, zerr := zw.Create("instructions.csv")
		if zerr != nil {
			_ = zw.Close()
			return nil, "", "", zerr
		}
		cw := csv.NewWriter(inf)
		cw.Comma = ';'
		_ = cw.Write([]string{"Column Name (Template)", "Instruction"})
		for _, header := range outletInstructionColumnOrder {
			if header == "is_active" {
				continue
			}
			kolom := toDisplayHeader(header)
			keterangan := outletInstructionByHeader[header]
			_ = cw.Write([]string{kolom, keterangan})
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			_ = zw.Close()
			return nil, "", "", err
		}
		if zerr := zw.Close(); zerr != nil {
			return nil, "", "", zerr
		}
		buffer = zipBuf
		contentType = "application/zip"
		filename = "outlet_template_update.zip"
	case "xls":
		buffer, err = service.createTemplateUpdateXLSX(custId, headers, data, sampleRow)
		contentType = "application/vnd.ms-excel"
		filename = "outlet_template_update.xls"
	default:
		buffer, err = service.createTemplateUpdateXLSX(custId, headers, data, sampleRow)
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		filename = "outlet_template_update.xlsx"
	}
	if err != nil {
		return nil, "", "", err
	}
	return buffer, contentType, filename, nil
}

func (service *outletServiceImpl) createTemplateUpdateXLSX(custId string, headers []string, data map[string][][]string, sampleRow []string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheetName := "Template Update"
	index, _ := f.NewSheet(sheetName)

	// Apply per-step colored styles similar to initial export template
	styleCache := map[string]int{}
	widths := make([]int, len(headers))
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		displayName := toDisplayHeader(header)
		if maxLen, ok := outletMaxCharHeaders[canonicalizeHeader(header)]; ok && maxLen > 0 {
			displayName = fmt.Sprintf("%s (maksimal %d karakter)", displayName, maxLen)
		}
		f.SetCellValue(sheetName, cell, displayName)
		step := stepCategoryForHeader(header)
		styleID, err := headerStyleID(f, styleCache, step)
		if err == nil {
			f.SetCellStyle(sheetName, cell, cell, styleID)
		}
		widths[i] = utf8.RuneCountInString(displayName)
	}

	// Hide internal columns in the template (Excel only)
	hiddenColumns := map[string]struct{}{
		"cust_id":           {},
		"outlet_id":         {},
		"outlet_contact_id": {},
		"outlet_bank_id":    {},
		"outlet_tax_id":     {},
	}
	for i, h := range headers {
		if _, hide := hiddenColumns[strings.ToLower(h)]; hide {
			if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
				// Best-effort: hide the column; ignore error to avoid breaking export
				_ = f.SetColVisible(sheetName, colName, false)
			}
		}
	}

	// Tulis contoh data pada baris ke-2 bila tersedia
	if len(sampleRow) == len(headers) {
		for i, v := range sampleRow {
			cell, _ := excelize.CoordinatesToCellName(i+1, 2)
			f.SetCellValue(sheetName, cell, v)
			if l := utf8.RuneCountInString(v); l > widths[i] {
				widths[i] = l
			}
		}
	}

	// Hitung jumlah baris data maksimum dari map `data`
	maxRows := 0
	for _, col := range data {
		if len(col) > maxRows {
			maxRows = len(col)
		}
	}

	// Tulis data dari repository dimulai baris ke-3 (setelah sample)
	for r := 0; r < maxRows; r++ {
		for c, h := range headers {
			var val string
			if h == "cust_id" {
				val = custId
			} else if col, ok := data[h]; ok {
				if r < len(col) && len(col[r]) > 0 {
					val = col[r][0]
				}
			}

			cell, _ := excelize.CoordinatesToCellName(c+1, r+3)
			f.SetCellValue(sheetName, cell, val)
			if l := utf8.RuneCountInString(val); l > widths[c] {
				widths[c] = l
			}
		}
	}

	insSheet := "Instructions"
	f.NewSheet(insSheet)
	ih := []string{"Column Name (Template)", "Instruction"}
	styleHeaderID, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1}})
	for i, hdr := range ih {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(insSheet, cell, hdr)
		f.SetCellStyle(insSheet, cell, cell, styleHeaderID)
	}
	rowIdx := 2
	for _, header := range outletInstructionColumnOrder {
		if header == "is_active" {
			continue
		}
		kolom := toDisplayHeader(header)
		keterangan := outletInstructionByHeader[header]
		f.SetCellValue(insSheet, "A"+strconv.Itoa(rowIdx), kolom)
		f.SetCellValue(insSheet, "B"+strconv.Itoa(rowIdx), keterangan)
		rowIdx++
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")
	for i := range headers {
		if colName, err := excelize.ColumnNumberToName(i + 1); err == nil {
			width := float64(widths[i] + 4)
			if width < 18 {
				width = 18
			}
			if width > 60 {
				width = 60
			}
			_ = f.SetColWidth(sheetName, colName, colName, width)
		}
	}
	return f.WriteToBuffer()
}

func (service *outletServiceImpl) createTemplateUpdateCSV(custId string, headers []string, data map[string][][]string, sampleRow []string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'

	// Write display headers
	disp := make([]string, len(headers))
	for i := range headers {
		displayName := toDisplayHeader(headers[i])
		if maxLen, ok := outletMaxCharHeaders[canonicalizeHeader(headers[i])]; ok && maxLen > 0 {
			displayName = fmt.Sprintf("%s (maksimal %d karakter)", displayName, maxLen)
		}
		disp[i] = displayName
	}
	writer.Write(disp)

	// Tulis contoh data jika panjangnya sesuai header
	if len(sampleRow) == len(headers) {
		writer.Write(sampleRow)
	}

	maxRows := 0
	for _, fieldData := range data {
		if len(fieldData) > maxRows {
			maxRows = len(fieldData)
		}
	}

	for r := 0; r < maxRows; r++ {
		record := make([]string, len(headers))
		for c, h := range headers {
			if h == "cust_id" {
				record[c] = custId
				continue
			}
			if strings.EqualFold(h, "outlet_id") {
				// never output outlet_id in template csv
				continue
			}
			if col, ok := data[h]; ok {
				if r < len(col) && len(col[r]) > 0 {
					record[c] = col[r][0]
				}
			}
		}
		writer.Write(record)
	}

	writer.Flush()
	return buffer, writer.Error()
}

// --- FUNGSI IMPOR ---
func outletDisplayNameID(field string) string {
	key := strings.ToLower(strings.TrimSpace(field))
	if name, ok := outletFieldDisplayNameID[key]; ok {
		return name
	}
	return toDisplayHeader(field)
}

func validateMaxLength(row map[string]string) error {
	for field, maxLen := range outletMaxCharHeaders {
		if val, ok := row[field]; ok && strings.TrimSpace(val) != "" {
			if len([]rune(val)) > maxLen {
				return fmt.Errorf("%s melebihi panjang maksimal %d karakter (panjang saat ini %d)", outletDisplayNameID(field), maxLen, len([]rune(val)))
			}
		}
	}
	return nil
}

type outletRowValidationMeta struct {
	code      string
	rowNumber int
}

func formatOutletValidationPrefix(outletCode string, rowNumber int) string {
	code := strings.TrimSpace(outletCode)
	if code != "" {
		return fmt.Sprintf("Kode outlet %s", code)
	}
	return fmt.Sprintf("Baris ke-%d", rowNumber)
}

// precheckActiveOutletCodeConfigForImport hanya memastikan konfigurasi outlet code ada (tanpa mengubah last_sequence).
// Generate kode sebenarnya dilakukan di dalam transaksi insert (GetActiveConfigForUpdate + UpdateLastSequenceNoWithTx).
func (service *outletServiceImpl) precheckActiveOutletCodeConfigForImport(custId, parentCustId, createdByName string) error {
	custIdForConfig := strings.TrimSpace(parentCustId)
	if custIdForConfig == "" {
		custIdForConfig = custId
	}
	if custIdForConfig == custId {
		if parent, errParent := service.OutletRepository.FindOneParentCustId(custId); errParent == nil && strings.TrimSpace(parent.ParentCustId) != "" {
			custIdForConfig = parent.ParentCustId
		}
	}
	yearNow := time.Now().Year()
	statuses := []string{"active"}
	createdByName = strings.TrimSpace(createdByName)

	if createdByName != "" {
		cfg, err := service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(custId, yearNow, statuses, createdByName)
		if err != nil {
			return fmt.Errorf("failed to verify setup outlet code: %w", err)
		}
		if cfg != nil {
			return nil
		}
		if custIdForConfig != custId {
			cfg, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(custIdForConfig, yearNow, statuses, createdByName)
			if err != nil {
				return fmt.Errorf("failed to verify setup outlet code: %w", err)
			}
			if cfg != nil {
				return nil
			}
		}
	}

	cfg, err := service.OutletCodeRepository.FindOneByCustIdYearAndStatus(custId, yearNow, statuses)
	if err != nil {
		return fmt.Errorf("failed to verify setup outlet code: %w", err)
	}
	if cfg != nil {
		return nil
	}
	if custIdForConfig != custId {
		cfg, err = service.OutletCodeRepository.FindOneByCustIdYearAndStatus(custIdForConfig, yearNow, statuses)
		if err != nil {
			return fmt.Errorf("failed to verify setup outlet code: %w", err)
		}
		if cfg != nil {
			return nil
		}
	}
	return errors.New("outlet code is required or setup outlet code (status Active) for current year is not available")
}

func (service *outletServiceImpl) ImportOutletsFromXLSX(req entity.ImportRequest) error {
	f, err := excelize.OpenReader(req.File)
	if err != nil {
		return errors.New("gagal membuka file Excel (XLS/XLSX)")
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return errors.New("gagal mengambil data baris dari sheet Excel")
	}

	if len(rows) < 1 {
		return errors.New("file Excel (XLS/XLSX) harus memiliki header dan minimal satu baris data")
	}

	var (
		headers          []string
		dataRows         [][]string
		dataColumnOffset int
		headerRowNumber  int
	)
	if req.IsImportNew {
		layout := resolveImportNewSheetLayout(rows)
		headers = layout.headers
		dataRows = layout.dataRows
		dataColumnOffset = layout.dataColumnOffset
		headerRowNumber = layout.headerRowNumber
	} else {
		headers = rows[0]
		dataRows = rows[1:]
		headerRowNumber = 1
	}

	if len(headers) == 0 {
		return errors.New("file Excel (XLS/XLSX) harus memiliki header")
	}
	if len(dataRows) == 0 {
		return errors.New("file Excel (XLS/XLSX) harus memiliki minimal satu baris data")
	}
	logrus.Infof("XLSX Header: %+v", headers)

	rowValidationErrors := make(map[int][]string)
	rowOrder := make([]int, 0)
	rowMeta := make(map[int]outletRowValidationMeta)
	resolvedRowMaps := make([]map[string]string, 0, len(dataRows))
	var importOutletCodePrecheck sync.Once
	var importOutletCodePrecheckErr error
	for i, rowValues := range dataRows {
		var rowMap map[string]string
		if req.IsImportNew {
			rowMap = buildImportNewRowMap(headers, rowValues, dataColumnOffset)
		} else {
			rowMap = make(map[string]string)
			for j, val := range rowValues {
				if j < len(headers) {
					orig := headers[j]
					cleaned := canonicalizeHeader(orig)
					logrus.Infof("Header real: '%s' -> cleaned: '%s'", headers[j], cleaned)
					rowMap[strings.ToLower(cleaned)] = strings.TrimSpace(val)
				}
			}
		}
		outletCode := strings.TrimSpace(rowMap["outlet_code"])
		if outletCode == "0000000000" {
			outletCode = ""
		}
		if outletCode == "" {
			importOutletCodePrecheck.Do(func() {
				importOutletCodePrecheckErr = service.precheckActiveOutletCodeConfigForImport(req.CustId, req.ParentCustId, req.CreatedByName)
			})
			if importOutletCodePrecheckErr != nil {
				rowIdx := i + 1
				if _, exists := rowValidationErrors[rowIdx]; !exists {
					rowOrder = append(rowOrder, rowIdx)
				}
				rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], importOutletCodePrecheckErr.Error())
			}
		}
		resolvedRowMaps = append(resolvedRowMaps, rowMap)
		rowIdx := i + 1
		rowNumber := headerRowNumber + 1 + i
		outletCode = strings.TrimSpace(rowMap["outlet_code"])
		rowMeta[rowIdx] = outletRowValidationMeta{code: outletCode, rowNumber: rowNumber}
		if err := validateImportRowMandatory(req, rowMap); err != nil {
			if _, exists := rowValidationErrors[rowIdx]; !exists {
				rowOrder = append(rowOrder, rowIdx)
			}
			rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
		}
		if err := validateMaxLength(rowMap); err != nil {
			if _, exists := rowValidationErrors[rowIdx]; !exists {
				rowOrder = append(rowOrder, rowIdx)
			}
			rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
		}
	}

	if len(rowValidationErrors) > 0 {
		historyId, err := service.OutletRepository.CreateImportHistory(importUploadType(req), req.Filename, req.CustId, req.UserId, len(dataRows))
		if err != nil {
			return err
		}
		validationSummaries := make([]string, 0, len(rowOrder))
		for _, rowIdx := range rowOrder {
			msgs := rowValidationErrors[rowIdx]
			meta := rowMeta[rowIdx]
			prefix := formatOutletValidationPrefix(meta.code, meta.rowNumber)
			validationSummaries = append(validationSummaries, fmt.Sprintf("%s: %s", prefix, strings.Join(msgs, "; ")))
			rowValues := dataRows[rowIdx-1]
			var rowMap map[string]string
			if req.IsImportNew {
				rowMap = buildImportNewRowMap(headers, rowValues, dataColumnOffset)
			} else {
				rowMap = make(map[string]string)
				for j, val := range rowValues {
					if j < len(headers) {
						cleaned := canonicalizeHeader(headers[j])
						rowMap[strings.ToLower(cleaned)] = strings.TrimSpace(val)
					}
				}
			}
			fullRow := expandImportRowFromMap(req, rowMap)
			outletData, mapErr := service.mapImportRowToOutletStruct(req, fullRow)
			if mapErr != nil {
				if len(rowValues) > 0 {
					outletData.OutletCode = rowValues[0]
				}
				if len(rowValues) > 1 {
					outletData.OutletName = rowValues[1]
				}
			}
			outletData.ErrorMessage = fmt.Sprintf("%s: %s", prefix, strings.Join(msgs, "; "))
			if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, outletData); tempErr != nil {
				logrus.Errorf("gagal menyimpan data outlet_temp pada baris %d: %v", rowIdx+1, tempErr)
			}
		}
		_ = service.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rowValidationErrors), true)
		return fmt.Errorf("validasi gagal (%d baris): %s", len(rowValidationErrors), strings.Join(validationSummaries, "; "))
	}

	historyId, err := service.OutletRepository.CreateImportHistory(importUploadType(req), req.Filename, req.CustId, req.UserId, len(dataRows))
	if err != nil {
		return err
	}

	// Expand rows from resolvedRowMaps (termasuk outlet_code yang di-auto-generate) agar fullRow dipakai processOutletRows
	expandedRows := make([][]string, 0, len(resolvedRowMaps))
	for _, rowMap := range resolvedRowMaps {
		expandedRows = append(expandedRows, expandImportRowFromMap(req, rowMap))
	}

	// Async: proses jalan di background
	go func() {
		if err := service.processOutletRows(req, expandedRows, req.Filename, historyId); err != nil {
			log.Printf("[ASYNC XLSX IMPORT] error: %v", err)
		}
	}()

	return nil
}

func (service *outletServiceImpl) ImportOutletsFromCSV(req entity.ImportRequest) error {
	buf, err := io.ReadAll(req.File)
	if err != nil {
		return errors.New("gagal membaca file")
	}

	content := string(buf)
	lines := strings.Split(content, "\n")
	// Find the first non-empty, non-comment line as header
	headerIdx := -1
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		headerIdx = i
		break
	}
	if headerIdx == -1 {
		return errors.New("file CSV harus memiliki baris header")
	}

	header := strings.Split(strings.TrimSpace(lines[headerIdx]), ";")
	dataColumnOffset := 0
	dataStartIdx := headerIdx + 1
	if req.IsImportNew && len(header) > 0 && strings.EqualFold(strings.TrimSpace(header[0]), "Template Field") {
		dataColumnOffset = 1
		header = trimLeadingCells(header, dataColumnOffset)
		for nextIdx := headerIdx + 1; nextIdx < len(lines); nextIdx++ {
			line := strings.TrimSpace(lines[nextIdx])
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			fields := strings.Split(line, ";")
			if len(fields) > 0 && strings.EqualFold(strings.TrimSpace(fields[0]), "Response Field") {
				header = trimLeadingCells(fields, dataColumnOffset)
				dataStartIdx = nextIdx + 1
				break
			}
			break
		}
	}
	log.Println("CSV Header:", header)

	records := [][]string{}
	for i, line := range lines[dataStartIdx:] {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ";")

		// FIX: kalau header nempel dengan row pertama
		if strings.Contains(fields[0], "tax_identifier_no") {
			parts := strings.SplitN(fields[0], "tax_identifier_no", 2)
			header[len(header)-1] = "tax_identifier_no"
			if len(parts) == 2 && parts[1] != "" {
				fields[0] = parts[1]
			} else {
				continue
			}
		}

		if len(fields) < len(header) {
			log.Printf("Row %d skipped, expected %d columns, got %d", i+1, len(header), len(fields))
			continue
		}

		records = append(records, fields)
		log.Printf("Row %d parsed: %+v", i+1, fields[:5])
	}

	if req.IsImportNew {
		records = filterImportNewDataRows(records, dataColumnOffset)
	}

	// VALIDASI MANDATORY (sama seperti versi XLSX) dengan normalisasi header display -> internal
	rowValidationErrors := make(map[int][]string)
	rowOrder := make([]int, 0)
	rowMeta := make(map[int]outletRowValidationMeta)
	resolvedRowMapsCSV := make([]map[string]string, 0, len(records))
	var importOutletCodePrecheckCSV sync.Once
	var importOutletCodePrecheckErrCSV error
	if len(records) > 0 {
		// fungsi normalisasi lokal (duplikasi dari ImportOutletsFromXLSX untuk konsistensi cepat)
		for i, rowValues := range records { // i mulai 0, baris data asli adalah i+2
			var rowMap map[string]string
			if req.IsImportNew {
				rowMap = buildImportNewRowMap(header, rowValues, dataColumnOffset)
			} else {
				rowMap = make(map[string]string)
				for j, val := range rowValues {
					if j < len(header) {
						cleaned := canonicalizeHeader(header[j])
						rowMap[cleaned] = strings.TrimSpace(val)
					}
				}
			}
			outletCode := strings.TrimSpace(rowMap["outlet_code"])
			if outletCode == "0000000000" {
				outletCode = ""
			}
			if outletCode == "" {
				importOutletCodePrecheckCSV.Do(func() {
					importOutletCodePrecheckErrCSV = service.precheckActiveOutletCodeConfigForImport(req.CustId, req.ParentCustId, req.CreatedByName)
				})
				if importOutletCodePrecheckErrCSV != nil {
					rowIdx := i
					if _, exists := rowValidationErrors[rowIdx]; !exists {
						rowOrder = append(rowOrder, rowIdx)
					}
					rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], importOutletCodePrecheckErrCSV.Error())
				}
			}
			resolvedRowMapsCSV = append(resolvedRowMapsCSV, rowMap)
			rowIdx := i // index dalam records slice
			rowNumber := i + 2
			outletCode = strings.TrimSpace(rowMap["outlet_code"])
			rowMeta[rowIdx] = outletRowValidationMeta{code: outletCode, rowNumber: rowNumber}
			if err := validateImportRowMandatory(req, rowMap); err != nil {
				if _, exists := rowValidationErrors[rowIdx]; !exists {
					rowOrder = append(rowOrder, rowIdx)
				}
				rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
			}
			if err := validateMaxLength(rowMap); err != nil {
				if _, exists := rowValidationErrors[rowIdx]; !exists {
					rowOrder = append(rowOrder, rowIdx)
				}
				rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
			}
		}
	}

	if len(records) == 0 {
		return errors.New("no valid data rows found")
	}

	if len(rowValidationErrors) > 0 {
		historyId, err := service.OutletRepository.CreateImportHistory(importUploadType(req), req.Filename, req.CustId, req.UserId, len(records))
		if err != nil {
			return err
		}
		validationSummaries := make([]string, 0, len(rowOrder))
		for _, rowIdx := range rowOrder {
			msgs := rowValidationErrors[rowIdx]
			meta := rowMeta[rowIdx]
			prefix := formatOutletValidationPrefix(meta.code, meta.rowNumber)
			validationSummaries = append(validationSummaries, fmt.Sprintf("%s: %s", prefix, strings.Join(msgs, "; ")))
			rowValues := records[rowIdx]
			rowMap := make(map[string]string)
			for j, val := range rowValues {
				if j < len(header) {
					cleaned := canonicalizeHeader(header[j])
					rowMap[cleaned] = strings.TrimSpace(val)
				}
			}
			fullRow := expandImportRowFromMap(req, rowMap)
			outletData, mapErr := service.mapImportRowToOutletStruct(req, fullRow)
			if mapErr != nil {
				if len(rowValues) > 0 {
					outletData.OutletCode = rowValues[0]
				}
				if len(rowValues) > 1 {
					outletData.OutletName = rowValues[1]
				}
			}
			outletData.ErrorMessage = fmt.Sprintf("%s: %s", prefix, strings.Join(msgs, "; "))
			if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, outletData); tempErr != nil {
				logrus.Errorf("gagal menyimpan data outlet_temp pada baris %d: %v", meta.rowNumber, tempErr)
			}
		}
		_ = service.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rowValidationErrors), true)
		return fmt.Errorf("validasi gagal (%d baris): %s", len(rowValidationErrors), strings.Join(validationSummaries, "; "))
	}

	historyId, err := service.OutletRepository.CreateImportHistory(importUploadType(req), req.Filename, req.CustId, req.UserId, len(records))
	if err != nil {
		return err
	}

	// Expand rows from resolvedRowMapsCSV (termasuk outlet_code yang di-auto-generate)
	expandedRows := make([][]string, 0, len(resolvedRowMapsCSV))
	for _, rowMap := range resolvedRowMapsCSV {
		expandedRows = append(expandedRows, expandImportRowFromMap(req, rowMap))
	}

	// Async: proses jalan di background
	go func() {
		if err := service.processOutletRows(req, expandedRows, req.Filename, historyId); err != nil {
			log.Printf("[ASYNC CSV IMPORT] error: %v", err)
		}
	}()

	return nil
}

// --- Lanjut proses ---

// func (service *outletServiceImpl) validateOutletHeader(header []string) error {
// 	requiredHeaders := map[string]bool{
// 		"outlet_code": false,
// 		"outlet_name": false,
// 		// "area_code":   true,
// 		// "area_name":   true,
// 	}

// 	headerMap := make(map[string]bool)
// 	for _, h := range header {
// 		headerMap[strings.ToLower(h)] = true
// 	}

// 	var missingHeaders []string
// 	for reqHeader := range requiredHeaders {
// 		if !headerMap[reqHeader] {
// 			missingHeaders = append(missingHeaders, reqHeader)
// 		}
// 	}

// 	if len(missingHeaders) > 0 {
// 		return fmt.Errorf("missing required header columns: %s", strings.Join(missingHeaders, ", "))
// 	}
// 	return nil
// }

type outletImportJob struct {
	index int
	row   []string
}

type outletImportResult struct {
	success  bool
	nonFatal bool
	err      error
}

func (service *outletServiceImpl) resolveBankForImport(req entity.ImportRequest, data *entity.OutletTemp) (*model.MOutletBank, string, error) {
	bankName := strings.TrimSpace(data.BankName)
	accountNo := strings.TrimSpace(data.AccountNo)
	accountName := strings.TrimSpace(data.AccountName)

	if bankName == "" && accountNo == "" && accountName == "" {
		return nil, "", nil
	}

	if bankName == "" || accountNo == "" || accountName == "" {
		return nil, "", nil
	}

	bankCode := strings.TrimSpace(data.BankCode)
	var bankID int64
	if bankCode != "" {
		id, err := service.OutletRepository.FindBankIdByCode(req.CustId, bankCode)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// fallback by name
			} else {
				return nil, "", fmt.Errorf("gagal mengambil master bank (kode %s): %w", bankCode, err)
			}
		} else if id > 0 {
			bankID = id
		}
	}

	if bankID == 0 {
		id, err := service.OutletRepository.FindBankIdByName(req.CustId, bankName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, "", nil
			}
			return nil, "", fmt.Errorf("gagal mengambil master bank \"%s\": %w", bankName, err)
		}
		bankID = id
	}

	if bankID == 0 {
		return nil, "", nil
	}

	return &model.MOutletBank{
		CustID:      req.CustId,
		BankId:      bankID,
		AccountNo:   accountNo,
		AccountName: accountName,
	}, "", nil
}

func (service *outletServiceImpl) processOutletRows(req entity.ImportRequest, rows [][]string, filename string, historyId int64) error {
	mappingCustId := req.CustId
	if req.ParentCustId != "" {
		mappingCustId = req.ParentCustId
	} else {
		if parent, err := service.OutletRepository.FindOneParentCustId(req.CustId); err == nil && parent.ParentCustId != "" {
			mappingCustId = parent.ParentCustId
		}
	}

	total := len(rows)
	if total == 0 {
		if err := service.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, 0, false); err != nil {
			return err
		}
		return nil
	}

	workerCount := runtime.NumCPU()
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > total {
		workerCount = total
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan outletImportJob)
	results := make(chan outletImportResult)
	var workersWg sync.WaitGroup

	for w := 0; w < workerCount; w++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					res := service.processSingleOutletRow(req, mappingCustId, job.row, job.index+2, historyId)
					isFatal := res.err != nil && !res.nonFatal
					select {
					case results <- res:
					case <-ctx.Done():
						return
					}
					if isFatal {
						cancel()
						return
					}
				}
			}
		}()
	}

	go func() {
		workersWg.Wait()
		close(results)
	}()

	var successCount atomic.Int64
	var nonFatalCount atomic.Int64
	var fatalErr error
	var fatalOnce sync.Once

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		defer resultsWg.Done()
		for res := range results {
			if res.success {
				successCount.Add(1)
			}
			if res.nonFatal {
				nonFatalCount.Add(1)
			}
			if res.err != nil && !res.nonFatal {
				fatalOnce.Do(func() {
					fatalErr = res.err
				})
			}
		}
	}()

sendLoop:
	for idx, row := range rows {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- outletImportJob{index: idx, row: row}:
		}
	}
	close(jobs)

	resultsWg.Wait()

	if fatalErr != nil {
		successes := int(successCount.Load())
		failed := total - successes
		if other := int(nonFatalCount.Load()); other > failed {
			failed = other
		}
		if failed < 0 {
			failed = 0
		}
		if err := service.OutletRepository.UpdateImportHistoryOutlet(historyId, successes, failed, true); err != nil {
			logrus.Errorf("gagal memperbarui import_history pada kegagalan fatal: %v", err)
		}
		return fatalErr
	}

	successes := int(successCount.Load())
	failed := total - successes
	if other := int(nonFatalCount.Load()); other > failed {
		failed = other
	}

	if err := service.OutletRepository.UpdateImportHistoryOutlet(historyId, successes, failed, false); err != nil {
		return err
	}

	return nil
}

func (service *outletServiceImpl) processSingleOutletRow(req entity.ImportRequest, mappingCustId string, row []string, rowNumber int, historyId int64) outletImportResult {
	importData, err := service.mapImportRowToOutletStruct(req, row)
	logrus.Info(err)
	if err != nil {
		if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, importData); tempErr != nil {
			logrus.Errorf("gagal menyimpan data outlet_temp pada baris %d: %v", rowNumber, tempErr)
		}
		return outletImportResult{nonFatal: true}
	}

	code := strings.TrimSpace(importData.OutletCode)
	if code == "0000000000" {
		code = ""
		importData.OutletCode = ""
	}
	if code != "" {
		exists, err := service.OutletRepository.CheckOutletCodeExists(req.CustId, code)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal memeriksa kode outlet pada baris %d: %w", rowNumber, err)}
		}
		if exists {
			msg := fmt.Sprintf("Kode outlet %s sudah terdaftar", code)
			importData.ErrorMessage = msg
			if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, importData); tempErr != nil {
				logrus.Errorf("gagal mencatat duplikasi outlet pada baris %d: %v", rowNumber, tempErr)
			}
			return outletImportResult{nonFatal: true}
		}
	}

	// Map human labels to enum codes before parsing
	importData.OutletStatus = parseEnumCode(importData.OutletStatus, outletStatusCodes)
	if strings.TrimSpace(importData.PaymentType) == "" {
		importData.PaymentType = parseEnumCode(importData.PaymentTypeName, paymentTypeCodes)
	} else {
		importData.PaymentType = parseEnumCode(importData.PaymentType, paymentTypeCodes)
	}
	if strings.TrimSpace(importData.CreditLimitType) == "" {
		importData.CreditLimitType = parseEnumCode(importData.CreditLimitTypeName, creditLimitTypeCodes)
	} else {
		importData.CreditLimitType = parseEnumCode(importData.CreditLimitType, creditLimitTypeCodes)
	}
	if strings.TrimSpace(importData.SalesInvLimitType) == "" {
		importData.SalesInvLimitType = parseEnumCode(importData.SalesInvLimitTypeName, salesInvLimitTypeCodes)
	} else {
		importData.SalesInvLimitType = parseEnumCode(importData.SalesInvLimitType, salesInvLimitTypeCodes)
	}
	if arStatusCode := parseEnumCode(importData.ArStatus, arStatusCodes); arStatusCode != "" {
		importData.ArStatus = arStatusCode
	} else {
		importData.ArStatus = parseEnumCode(importData.ArStatusName, arStatusCodes)
	}
	if strings.TrimSpace(importData.ObsType) == "" {
		importData.ObsType = parseEnumCode(importData.ObyTypeName, obsTypeCodes)
	} else {
		importData.ObsType = parseEnumCode(importData.ObsType, obsTypeCodes)
	}
	if strings.TrimSpace(importData.CreditLimitAction) == "" {
		importData.CreditLimitAction = parseEnumCode(importData.CreditLimitActionName, limitActionCodes)
	} else {
		importData.CreditLimitAction = parseEnumCode(importData.CreditLimitAction, limitActionCodes)
	}
	if strings.TrimSpace(importData.SalesInvLimitAction) == "" {
		importData.SalesInvLimitAction = parseEnumCode(importData.SalesInvLimitActionName, limitActionCodes)
	} else {
		importData.SalesInvLimitAction = parseEnumCode(importData.SalesInvLimitAction, limitActionCodes)
	}
	if strings.TrimSpace(importData.ObsLimitAction) == "" {
		importData.ObsLimitAction = parseEnumCode(importData.ObsLimitActionName, limitActionCodes)
	} else {
		importData.ObsLimitAction = parseEnumCode(importData.ObsLimitAction, limitActionCodes)
	}
	if strings.TrimSpace(importData.VerificationStatus) == "" {
		importData.VerificationStatus = parseEnumCode(importData.VerificationStatus, verificationStatusCodes)
	} else {
		importData.VerificationStatus = parseEnumCode(importData.VerificationStatus, verificationStatusCodes)
	}
	if strings.TrimSpace(importData.TaxInvoiceForm) == "" {
		importData.TaxInvoiceForm = parseEnumCode(importData.TaxInvoiceFormName, taxInvoiceFormCodes)
	} else {
		importData.TaxInvoiceForm = parseEnumCode(importData.TaxInvoiceForm, taxInvoiceFormCodes)
	}

	// u := req.UserId

	missingMasters := make([]string, 0, 20)
	missingMasterAdded := make(map[string]struct{})
	addMissingMaster := func(msg string) {
		msg = strings.TrimSpace(msg)
		if msg == "" {
			return
		}
		if _, exists := missingMasterAdded[msg]; exists {
			return
		}
		missingMasterAdded[msg] = struct{}{}
		missingMasters = append(missingMasters, msg)
	}

	if name := strings.TrimSpace(importData.OutletProvince); name != "" {
		provinceId, err := service.OutletRepository.GetProvinceIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master outlet province pada baris %d: %w", rowNumber, err)}
		}
		importData.OutletProvinceId = provinceId
		if provinceId == "" {
			addMissingMaster("Master outlet province tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.OutletRegency); name != "" {
		regencyId, err := service.OutletRepository.GetRegencyIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master outlet regency pada baris %d: %w", rowNumber, err)}
		}
		importData.OutletRegencyId = regencyId
		if regencyId == "" {
			addMissingMaster("Master outlet regency tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.OutletSubDistrict); name != "" {
		subDistrictId, err := service.OutletRepository.GetSubDistrictIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master outlet sub_district pada baris %d: %w", rowNumber, err)}
		}
		importData.OutletSubDistrictId = subDistrictId
		if subDistrictId == "" {
			addMissingMaster("Master outlet sub_district tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.OutletWard); name != "" {
		wardId, err := service.OutletRepository.GetWardIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master outlet ward pada baris %d: %w", rowNumber, err)}
		}
		importData.OutletWardId = wardId
		if wardId == "" {
			addMissingMaster("Master outlet ward tidak ditemukan")
		}
	}

	if name := strings.TrimSpace(importData.DelvProvince); name != "" {
		provinceId, err := service.OutletRepository.GetProvinceIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master delivery province pada baris %d: %w", rowNumber, err)}
		}
		importData.DelvProvinceId = provinceId
		if provinceId == "" {
			addMissingMaster("Master delivery province tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.DelvRegency); name != "" {
		regencyId, err := service.OutletRepository.GetRegencyIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master delivery regency pada baris %d: %w", rowNumber, err)}
		}
		importData.DelvRegencyId = regencyId
		if regencyId == "" {
			addMissingMaster("Master delivery regency tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.DelvSubDistrict); name != "" {
		subDistrictId, err := service.OutletRepository.GetSubDistrictIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master delivery sub district pada baris %d: %w", rowNumber, err)}
		}
		importData.DelvSubDistrictId = subDistrictId
		if subDistrictId == "" {
			addMissingMaster("Master delivery sub district tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.DelvWard); name != "" {
		wardId, err := service.OutletRepository.GetWardIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master delivery ward pada baris %d: %w", rowNumber, err)}
		}
		importData.DelvWardId = wardId
		if wardId == "" {
			addMissingMaster("Master delivery ward tidak ditemukan")
		}
	}

	if name := strings.TrimSpace(importData.InvProvince); name != "" {
		provinceId, err := service.OutletRepository.GetProvinceIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master invoice province pada baris %d: %w", rowNumber, err)}
		}
		importData.InvProvinceId = provinceId
		if provinceId == "" {
			addMissingMaster("Master invoice province tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.InvRegency); name != "" {
		regencyId, err := service.OutletRepository.GetRegencyIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master invoice regency pada baris %d: %w", rowNumber, err)}
		}
		importData.InvRegencyId = regencyId
		if regencyId == "" {
			addMissingMaster("Master invoice regency tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.InvSubDistrict); name != "" {
		subDistrictId, err := service.OutletRepository.GetSubDistrictIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master invoice sub_district pada baris %d: %w", rowNumber, err)}
		}
		importData.InvSubDistrictId = subDistrictId
		if subDistrictId == "" {
			addMissingMaster("Master invoice sub_district tidak ditemukan")
		}
	}
	if name := strings.TrimSpace(importData.InvWard); name != "" {
		wardId, err := service.OutletRepository.GetWardIdByName(mappingCustId, name)
		if err != nil {
			return outletImportResult{err: fmt.Errorf("gagal mengambil master invoice ward pada baris %d: %w", rowNumber, err)}
		}
		importData.InvWardId = wardId
		if wardId == "" {
			addMissingMaster("Master invoice ward tidak ditemukan")
		}
	}

	resolveMaster := func(label, errorLabel, code, name string, findByCode func(string) (int64, error), findByName func(string) (int64, error)) (int64, error) {
		code = strings.TrimSpace(code)
		name = strings.TrimSpace(name)
		if code == "" && name == "" {
			return 0, nil
		}

		recordFatal := func(detail error) (int64, error) {
			prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
			message := fmt.Sprintf("%s: gagal mengambil master %s: %v", prefix, errorLabel, detail)
			importData.ErrorMessage = message
			if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, importData); tempErr != nil {
				logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, tempErr)
			}
			fullErr := fmt.Errorf("gagal mengambil master %s pada baris %d: %w", label, rowNumber, detail)
			return 0, fullErr
		}

		if code != "" && findByCode != nil {
			id, errLookup := findByCode(code)
			if errLookup == nil && id > 0 {
				return id, nil
			}
			if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
				return recordFatal(errLookup)
			}
		}
		if code == "" && name != "" && findByCode != nil {
			id, errLookup := findByCode(name)
			if errLookup == nil && id > 0 {
				return id, nil
			}
			if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
				return recordFatal(errLookup)
			}
		}

		if name != "" && findByName != nil {
			id, errLookup := findByName(name)
			if errLookup == nil && id > 0 {
				return id, nil
			}
			if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
				return recordFatal(errLookup)
			}
		}

		addMissingMaster(fmt.Sprintf("Master %s tidak ditemukan", errorLabel))
		return 0, nil
	}

	var (
		discGrpId  int64
		otLocId    int64
		otGrpId    int64
		priceGrpId int64
		districtId int64
		otClassId  int64
		industryId int64
		marketId   int64
		otTypeId   int64
	)

	if id, err := resolveMaster(
		"disc group",
		"disc group",
		importData.DiscGrpCode,
		importData.DiscGrpName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindDiscGroupIdByCode(req.CustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindDiscGroupIdByName(req.CustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		discGrpId = id
	}

	if id, err := resolveMaster(
		"outlet location",
		"outlet location",
		importData.OtLocCode,
		importData.OtLocName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindOutletLocIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindOutletLocIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		otLocId = id
	}

	if id, err := resolveMaster(
		"outlet group",
		"outlet group",
		importData.OtGrpCode,
		importData.OtGrpName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindOutletGroupIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindOutletGroupIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		otGrpId = id
	}

	if id, err := resolveMaster(
		"price group",
		"price group",
		importData.PriceGrpCode,
		importData.PriceGrpName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindPriceGroupIdByCode(req.CustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindPriceGroupIdByName(req.CustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		priceGrpId = id
	}

	if id, err := resolveMaster(
		"district",
		"district",
		importData.DistrictCode,
		importData.DistrictName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindDistrictIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindDistrictIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		districtId = id
	}

	if id, err := resolveMaster(
		"outlet class",
		"outlet class",
		importData.OtClassCode,
		importData.OtClassName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindOutletClassIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindOutletClassIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		otClassId = id
	}

	if id, err := resolveMaster(
		"industry",
		"industry",
		importData.IndustryCode,
		importData.IndustryName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindIndustryIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindIndustryIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		industryId = id
	}

	if id, err := resolveMaster(
		"market",
		"market",
		importData.MarketCode,
		importData.MarketName,
		func(code string) (int64, error) { return service.OutletRepository.FindMarketIdByCode(req.CustId, code) },
		func(name string) (int64, error) { return service.OutletRepository.FindMarketIdByName(req.CustId, name) },
	); err != nil {
		return outletImportResult{err: err}
	} else {
		marketId = id
	}

	if id, err := resolveMaster(
		"outlet type",
		"outlet type",
		importData.OtTypeCode,
		importData.OtTypeName,
		func(code string) (int64, error) {
			return service.OutletRepository.FindOutletTypeIdByCode(mappingCustId, code)
		},
		func(name string) (int64, error) {
			return service.OutletRepository.FindOutletTypeIdByName(mappingCustId, name)
		},
	); err != nil {
		return outletImportResult{err: err}
	} else {
		otTypeId = id
	}

	var validationMessages []string

	if len(missingMasters) > 0 {
		msgBody := strings.Join(missingMasters, ", ")
		validationMessages = append(validationMessages, msgBody)
	}

	bankModel, bankMsg, bankErr := service.resolveBankForImport(req, &importData)
	if bankErr != nil {
		return outletImportResult{err: fmt.Errorf("gagal memeriksa bank pada baris %d: %w", rowNumber, bankErr)}
	}
	if bankMsg != "" {
		validationMessages = append(validationMessages, bankMsg)
	}

	if len(validationMessages) > 0 {
		prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
		msg := fmt.Sprintf("%s: %s", prefix, strings.Join(validationMessages, "; "))
		if existing := strings.TrimSpace(importData.ErrorMessage); existing != "" {
			msg = existing + "; " + msg
		}
		logrus.Warnf("import outlet: %s", msg)
		importData.ErrorMessage = msg
		if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, importData); tempErr != nil {
			logrus.Errorf("gagal menyimpan data outlet_temp pada baris %d: %v", rowNumber, tempErr)
		}
		return outletImportResult{nonFatal: true}
	}

	outletStatus, _ := strconv.ParseInt(importData.OutletStatus, 10, 16)
	top, _ := strconv.ParseInt(importData.Top, 10, 32)
	paymentType, _ := strconv.ParseInt(importData.PaymentType, 10, 16)
	// AC-04/AC-05: payment_type 1 or 2 (Cash On Delivery/Cash Before Delivery) -> top not required; if user fills TOP, set top=0
	if paymentType == 1 || paymentType == 2 {
		top = 0
	}
	isContraBon, _ := strconv.ParseBool(importData.IsContraBon)
	creditLimitType, _ := strconv.ParseInt(importData.CreditLimitType, 10, 16)
	creditLimit, _ := strconv.ParseFloat(importData.CreditLimit, 64)
	salesInvLimitType, _ := strconv.ParseInt(importData.SalesInvLimitType, 10, 16)
	salesInvLimit, _ := strconv.ParseInt(importData.SalesInvLimit, 10, 16)
	buildingOwnStr := parseEnumCode(importData.BuildingOwn, buildingOwnershipCodes)
	buildingOwn, _ := strconv.ParseInt(buildingOwnStr, 10, 16)
	arStatus, _ := strconv.ParseInt(importData.ArStatus, 10, 16)
	arTotal, _ := strconv.ParseFloat(importData.ArTotal, 64)
	isEmbBail, _ := strconv.ParseBool(importData.IsEmbBail)
	isActive := true
	if val := strings.TrimSpace(importData.IsActive); val != "" {
		if parsed, errBool := strconv.ParseBool(val); errBool == nil {
			isActive = parsed
		}
	}
	isWaNo, _ := strconv.ParseBool(importData.IsWaNo)
	isObs, _ := strconv.ParseBool(importData.IsObs)
	delvIsSameAddr, _ := strconv.ParseBool(importData.DelvIsSameAddr)
	invIsSameAddr, _ := strconv.ParseBool(importData.InvIsSameAddr)
	verificationStatus := resolveImportVerificationStatus(req)
	taxInvoiceForm, _ := strconv.ParseInt(importData.TaxInvoiceForm, 10, 16)
	obsType, _ := strconv.ParseInt(importData.ObsTypeName, 10, 64)
	creditLimitAction, _ := strconv.ParseInt(importData.CreditLimitAction, 10, 64)
	salesInvLimitAction, _ := strconv.ParseInt(importData.SalesInvLimitAction, 10, 64)
	obsLimitAction, _ := strconv.ParseInt(importData.ObsLimitAction, 10, 64)
	obs, _ := strconv.ParseInt(importData.Obs, 10, 64)
	beatId, _ := strconv.ParseInt(importData.BeatId, 10, 64)
	sbeatId, _ := strconv.ParseInt(importData.SbeatId, 10, 64)
	pluGrpId, _ := strconv.ParseInt(importData.PluGrpId, 10, 64)
	convGrpId, _ := strconv.ParseInt(importData.ConvGrpId, 10, 64)
	discInvId, _ := strconv.ParseInt(importData.DiscInvId, 10, 64)
	avgSalesWeek, _ := strconv.ParseFloat(importData.AvgSalesWeek, 64)
	avgSalesMonth, _ := strconv.ParseFloat(importData.AvgSalesMonth, 64)
	firstWeekNo64, _ := strconv.ParseInt(importData.FirstWeekNo, 10, 16)
	firstWeekNo := int16(firstWeekNo64)

	var firstTransDate, lastTransDate, otStartDate, otRegDate, dob, closedDate, outletEstablishmentDate *time.Time
	if importData.FirstTransDate != "" {
		if t, err := time.Parse("2006-01-02", importData.FirstTransDate); err == nil {
			firstTransDate = &t
		}
	}
	if importData.LastTransDate != "" {
		if t, err := time.Parse("2006-01-02", importData.LastTransDate); err == nil {
			lastTransDate = &t
		}
	}
	if importData.OtStartDate != "" {
		if t, err := time.Parse("2006-01-02", importData.OtStartDate); err == nil {
			otStartDate = &t
		}
	}
	if importData.OtRegDate != "" {
		if t, err := time.Parse("2006-01-02", importData.OtRegDate); err == nil {
			otRegDate = &t
		}
	}
	if importData.Dob != "" {
		if t, err := time.Parse("2006-01-02", importData.Dob); err == nil {
			dob = &t
		}
	}
	if importData.ClosedDate != "" {
		if t, err := time.Parse("2006-01-02", importData.ClosedDate); err == nil {
			closedDate = &t
		}
	}
	if importData.OutletEstablishmentDate != "" {
		if req.IsImportNew {
			if t, err := parseImportDateValue(importData.OutletEstablishmentDate); err == nil {
				outletEstablishmentDate = t
			}
		} else if t, err := time.Parse("2006-01-02", importData.OutletEstablishmentDate); err == nil {
			outletEstablishmentDate = &t
		}
	}

	now := time.Now()
	createdBy := int64(req.UserId)
	updatedBy := int64(req.UserId)

	parentCustId := req.ParentCustId
	if parentCustId == "" {
		parentCustId = mappingCustId
	}
	yyMMdd := now.Format("060102")
	regionCode, distributorCode, errRegion := service.OutletRepository.GetRegionAndDistributorCodeByCustId(req.CustId, parentCustId)
	if errRegion != nil {
		return outletImportResult{err: fmt.Errorf("gagal mengambil region/distributor pada baris %d: %w", rowNumber, errRegion)}
	}
	var prefix string
	if regionCode != nil && distributorCode != nil {
		prefix = fmt.Sprintf("%s-%s-%s-", *regionCode, *distributorCode, yyMMdd)
	} else {
		prefix = fmt.Sprintf("%s-", yyMMdd)
	}

	updatedByStr := strconv.FormatInt(req.UserId, 10)
	custIdForOutletCode := strings.TrimSpace(parentCustId)
	if custIdForOutletCode == "" {
		custIdForOutletCode = req.CustId
	}
	if custIdForOutletCode == req.CustId {
		if parent, errParent := service.OutletRepository.FindOneParentCustId(req.CustId); errParent == nil && strings.TrimSpace(parent.ParentCustId) != "" {
			custIdForOutletCode = parent.ParentCustId
		}
	}
	yearNow := now.Year()
	createdByNameForCode := strings.TrimSpace(req.CreatedByName)

	trx, err := service.OutletRepository.TrxBegin()
	if err != nil {
		return outletImportResult{err: fmt.Errorf("failed to open outlet transaction at row %d: %w", rowNumber, err)}
	}

	var outletCodeConfig *model.OutletCode
	var nextOutletCodeSeqStr string
	needAutoOutletCode := strings.TrimSpace(importData.OutletCode) == ""
	if needAutoOutletCode {
		statuses := []string{"active"}
		var errCode error
		if createdByNameForCode != "" {
			outletCodeConfig, errCode = service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(req.CustId, yearNow, statuses, createdByNameForCode)
			if (errCode != nil || outletCodeConfig == nil) && custIdForOutletCode != req.CustId {
				outletCodeConfig, errCode = service.OutletCodeRepository.FindOneByCustIdYearAndStatusAndCreatedBy(custIdForOutletCode, yearNow, statuses, createdByNameForCode)
			}
		}
		if outletCodeConfig == nil && errCode == nil {
			outletCodeConfig, errCode = service.OutletCodeRepository.FindOneByCustIdYearAndStatus(req.CustId, yearNow, statuses)
			if (errCode != nil || outletCodeConfig == nil) && custIdForOutletCode != req.CustId {
				outletCodeConfig, errCode = service.OutletCodeRepository.FindOneByCustIdYearAndStatus(custIdForOutletCode, yearNow, statuses)
			}
		}
		if errCode != nil {
			_ = trx.TrxRollback()
			return outletImportResult{err: fmt.Errorf("failed to reserve outlet code sequence at row %d: %w", rowNumber, errCode)}
		}
		if outletCodeConfig == nil {
			_ = trx.TrxRollback()
			return outletImportResult{err: fmt.Errorf("setup outlet code (Active) not found for generate code at row %d", rowNumber)}
		}
		nextOutletCodeSeqStr, errCode = service.OutletCodeRepository.IncrementSequenceByID(outletCodeConfig.Id, &updatedByStr)
		if errCode != nil || strings.TrimSpace(nextOutletCodeSeqStr) == "" {
			_ = trx.TrxRollback()
			if errCode != nil {
				return outletImportResult{err: fmt.Errorf("failed to reserve outlet code sequence at row %d: %w", rowNumber, errCode)}
			}
			return outletImportResult{err: fmt.Errorf("failed to reserve outlet code sequence at row %d", rowNumber)}
		}
		yearVal := outletCodeConfig.YearCode
		if yearVal >= 100 {
			yearVal = yearVal % 100
		}
		yearPart := fmt.Sprintf("%02d", yearVal)
		generated := strings.TrimSpace(outletCodeConfig.SerialCode) + yearPart + nextOutletCodeSeqStr
		if len(generated) > 10 {
			_ = trx.TrxRollback()
			return outletImportResult{err: fmt.Errorf("generated outlet code exceeds 10 characters at row %d", rowNumber)}
		}
		importData.OutletCode = generated
		exists, errDup := service.OutletRepository.CheckOutletCodeExists(req.CustId, generated)
		if errDup != nil {
			_ = trx.TrxRollback()
			return outletImportResult{err: fmt.Errorf("failed to check outlet code at row %d: %w", rowNumber, errDup)}
		}
		if exists {
			_ = trx.TrxRollback()
			importData.ErrorMessage = fmt.Sprintf("outlet code %s already registered", generated)
			if tempErr := service.OutletRepository.CreateOutletTemp(historyId, "failed", req.CustId, importData); tempErr != nil {
				logrus.Errorf("failed to record duplicate outlet at row %d: %v", rowNumber, tempErr)
			}
			return outletImportResult{nonFatal: true}
		}
	}

	seq, errSeq := service.OutletRepository.GetNextOutletPrincipalCodeSeqTx(trx.Tx(), prefix)
	if errSeq != nil {
		_ = trx.TrxRollback()
		return outletImportResult{err: fmt.Errorf("failed to generate outlet_principal_code at row %d: %w", rowNumber, errSeq)}
	}
	outletPrincipalCode := prefix + fmt.Sprintf("%04d", seq)

	if req.IsImportNew {
		if err := service.upsertImportOutletGeoMasters(mappingCustId, importData, req.UserId); err != nil {
			return outletImportResult{err: fmt.Errorf("failed to upsert outlet geo masters at row %d: %w", rowNumber, err)}
		}
	}

	importGeoProvinceID := ""
	importGeoRegencyID := ""
	importGeoSubDistrictID := ""
	if req.IsImportNew {
		importGeoProvinceID = importData.OutletProvinceId
		importGeoRegencyID = importData.OutletRegencyId
		importGeoSubDistrictID = importData.OutletSubDistrictId
	}

	processedData := entity.ProcessedOutlet{
		CustId:                  req.CustId,
		OutletPrincipalCode:     outletPrincipalCode,
		OutletCode:              importData.OutletCode,
		Barcode:                 importData.Barcode,
		OutletName:              importData.OutletName,
		OutletStatus:            int16(outletStatus),
		Address1:                importData.Address1,
		Address2:                importData.Address2,
		City:                    importData.City,
		ZipCode:                 importData.ZipCode,
		PhoneNo:                 outletPhoneOrContactPhone(importData.PhoneNo, importData.ContactPhoneNo),
		WaNo:                    importData.WaNo,
		FaxNo:                   importData.FaxNo,
		Email:                   importData.Email,
		DiscGrpId:               discGrpId,
		OtLocId:                 otLocId,
		OtGrpId:                 otGrpId,
		PriceGrpId:              priceGrpId,
		DistrictId:              districtId,
		BeatId:                  beatId,
		SbeatId:                 sbeatId,
		OtClassId:               otClassId,
		IndustryId:              industryId,
		MarketId:                marketId,
		Top:                     int32(top),
		PaymentType:             int16(paymentType),
		IsContraBon:             isContraBon,
		PluGrpId:                pluGrpId,
		ConvGrpId:               convGrpId,
		DiscInvId:               discInvId,
		AgentFrom:               importData.AgentFrom,
		CreditLimitType:         int16(creditLimitType),
		CreditLimitTypeName:     importData.CreditLimitTypeName,
		CreditLimit:             creditLimit,
		SalesInvLimitType:       int16(salesInvLimitType),
		SalesInvLimitTypeName:   importData.SalesInvLimitTypeName,
		SalesInvLimit:           int16(salesInvLimit),
		AvgSalesWeek:            avgSalesWeek,
		AvgSalesMonth:           avgSalesMonth,
		FirstTransDate:          firstTransDate,
		LastTransDate:           lastTransDate,
		FirstWeekNo:             firstWeekNo,
		OtStartDate:             otStartDate,
		OtRegDate:               otRegDate,
		BuildingOwn:             int16(buildingOwn),
		Dob:                     dob,
		ArStatus:                int16(arStatus),
		ArTotal:                 arTotal,
		ClosedDate:              closedDate,
		IsEmbBail:               isEmbBail,
		IsPkpOutlet:             isEmbBail,
		TaxName:                 importData.TaxName,
		TaxAddr1:                importData.TaxAddr1,
		TaxAddr2:                importData.TaxAddr2,
		TaxCity:                 importData.TaxCity,
		TaxNo:                   importData.TaxNo,
		OwnerName:               importData.OwnerName,
		OwnerAddr1:              importData.OwnerAddr1,
		OwnerAddr2:              importData.OwnerAddr2,
		OwnerCity:               importData.OwnerCity,
		OwnerPhoneNo:            importData.OwnerPhoneNo,
		OwnerIdNo:               importData.OwnerIdNo,
		DelvAddr1:               importData.DelvAddr1,
		DelvAddr2:               importData.DelvAddr2,
		DelvCity:                importData.DelvCity,
		InvAddr1:                importData.InvAddr1,
		InvAddr2:                importData.InvAddr2,
		InvCity:                 importData.InvCity,
		IsActive:                isActive,
		CreatedBy:               &createdBy,
		CreatedAt:               &now,
		UpdatedBy:               &updatedBy,
		UpdatedAt:               &now,
		IsDel:                   false,
		Latitude:                importData.Latitude,
		Longitude:               importData.Longitude,
		ImageUrl:                importData.ImageUrl,
		OtTypeId:                otTypeId,
		IsObs:                   isObs,
		Obs:                     obs,
		OutletProvinceId:        importGeoProvinceID,
		OutletRegencyId:         importGeoRegencyID,
		OutletSubDistrictId:     importGeoSubDistrictID,
		OutletWardId:            importData.OutletWardId,
		IsWaNo:                  isWaNo,
		DelvWardId:              importData.DelvWardId,
		DelvZipCode:             importData.DelvZipCode,
		DelvIsSameAddr:          delvIsSameAddr,
		InvWardId:               importData.InvWardId,
		InvZipCode:              importData.InvZipCode,
		InvIsSameAddr:           invIsSameAddr,
		VerificationStatus:      int16(verificationStatus),
		TaxInvoiceForm:          int16(taxInvoiceForm),
		ObsType:                 obsType,
		CreditLimitAction:       creditLimitAction,
		CreditLimitActionName:   importData.CreditLimitActionName,
		SalesInvLimitAction:     salesInvLimitAction,
		SalesInvLimitActionName: importData.SalesInvLimitActionName,
		ObsLimitAction:          obsLimitAction,
		OutletEstablishmentDate: outletEstablishmentDate,
		DelvCity2:               importData.DelvCity2,
		DelvLatitude:            importData.DelvLatitude,
		DelvLongitude:           importData.DelvLongitude,
		DelvLatitude2:           importData.DelvLatitude2,
		DelvLongitude2:          importData.DelvLongitude2,
		DelvWardId2:             importData.DelvWardId2,
		DelvZipCode2:            importData.DelvZipCode2,
	}
	logrus.Info(processedData.DelvZipCode, processedData.InvZipCode, processedData.ZipCode)

	outletID, err := service.OutletRepository.CreateOutletTx(trx.Tx(), processedData)
	if err != nil {
		_ = trx.TrxRollback()
		return outletImportResult{err: fmt.Errorf("failed to save outlet (code %s): %w", importData.OutletCode, err)}
	}
	if needAutoOutletCode && outletCodeConfig != nil && strings.TrimSpace(nextOutletCodeSeqStr) != "" {
		if errUpd := service.OutletCodeRepository.UpdateLastSequenceNoWithTx(trx.Tx(), outletCodeConfig.Id, nextOutletCodeSeqStr, &updatedByStr); errUpd != nil {
			_ = trx.TrxRollback()
			return outletImportResult{err: fmt.Errorf("failed to update outlet code last sequence at row %d: %w", rowNumber, errUpd)}
		}
	}
	if err := trx.TrxCommit(); err != nil {
		return outletImportResult{err: fmt.Errorf("failed to commit outlet at row %d: %w", rowNumber, err)}
	}

	if err := service.insertImportOutletDetails(req, outletID, importData, int64(taxInvoiceForm), isEmbBail, bankModel); err != nil {
		return outletImportResult{err: fmt.Errorf("failed to save outlet detail (code %s): %w", importData.OutletCode, err)}
	}

	return outletImportResult{success: true}
}

// outletPhoneOrContactPhone: untuk import, jika kolom Outlet Phone kosong tapi Contact Phone terisi, pakai Contact Phone sebagai phone_no di m_outlet agar tampil di list.
func outletPhoneOrContactPhone(outletPhone, contactPhone string) string {
	if strings.TrimSpace(outletPhone) != "" {
		return outletPhone
	}
	return strings.TrimSpace(contactPhone)
}

func parseImportBool(val string, defaultValue bool) bool {
	trimmed := strings.TrimSpace(strings.ToLower(val))
	switch trimmed {
	case "1", "true", "t", "ya", "yes", "y":
		return true
	case "0", "false", "f", "tidak", "no", "n":
		return false
	case "":
		return defaultValue
	default:
		return defaultValue
	}
}

func (service *outletServiceImpl) insertImportOutletDetails(req entity.ImportRequest, outletID int64, data entity.OutletTemp, taxInvoiceForm int64, isEmbBail bool, bankModel *model.MOutletBank) error {
	if bankModel != nil && bankModel.BankId != 0 {
		bankModel.CustID = req.CustId
		bankModel.OutletID = outletID
		if _, err := service.OutletRepository.CreateOutletBank(*bankModel); err != nil {
			return fmt.Errorf("failed to save outlet bank: %w", err)
		}
	} else {
		defaultBank := model.MOutletBank{
			CustID:      req.CustId,
			OutletID:    outletID,
			BankId:      0,
			AccountNo:   "000000",
			AccountName: "Non Bank",
		}
		if _, err := service.OutletRepository.CreateOutletBank(defaultBank); err != nil {
			return fmt.Errorf("failed to save outlet bank default: %w", err)
		}
	}

	// Step 1 Contact (mandatory)
	contactName := strings.TrimSpace(data.ContactName)
	jobTitle := strings.TrimSpace(data.JobTitle)
	identityType := strings.TrimSpace(data.IdentityType)
	identityNo := strings.TrimSpace(data.IdentityNo)
	if !isEmbBail && !req.IsImportNew {
		identityType = ""
		identityNo = ""
	}
	phone := strings.TrimSpace(data.ContactPhoneNo)
	wa := strings.TrimSpace(data.ContactWaNo)
	email := strings.TrimSpace(data.ContactEmail)
	if contactName != "" || jobTitle != "" || identityType != "" || identityNo != "" || phone != "" || wa != "" || email != "" {
		contactComplete := contactName != "" && phone != ""
		if contactComplete {
			isWa := parseImportBool(data.ContactIsWaNo, false)
			contactModel := model.MOutletContact{
				CustID:       req.CustId,
				OutletID:     outletID,
				ContactName:  contactName,
				JobTitle:     jobTitle,
				PhoneNo:      phone,
				WaNo:         wa,
				IdentityNo:   identityNo,
				IsWaNo:       isWa,
				Email:        email,
				IdentityType: identityType,
				FaxNumber:    strings.TrimSpace(data.FaxNumber),
			}
			if _, err := service.OutletRepository.CreateOutletContact(contactModel); err != nil {
				return fmt.Errorf("failed to save outlet contact: %w", err)
			}
		}
	}

	taxName := strings.TrimSpace(data.TaxName)
	addressTax := strings.TrimSpace(data.AddressTax)
	taxIdentifierType := strings.TrimSpace(data.TaxIdentifierType)
	if !isEmbBail && !req.IsImportNew {
		taxName = ""
		addressTax = ""
		taxIdentifierType = ""
	}
	if taxName != "" || addressTax != "" || taxIdentifierType != "" {
		if taxName != "" && addressTax != "" && taxIdentifierType != "" {
			taxModel := model.MOutletTax{
				CustID:            req.CustId,
				OutletID:          outletID,
				IsEmbBail:         isEmbBail,
				TaxName:           taxName,
				TaxAddr1:          addressTax,
				TaxAddr2:          "",
				TaxCity:           "",
				TaxNo:             strings.TrimSpace(data.TaxIdentifierNo),
				TaxInvoiceId:      taxInvoiceForm,
				TaxType:           strings.TrimSpace(data.TaxType),
				Nitku:             strings.TrimSpace(data.Nitku),
				AdddressTax:       addressTax,
				TaxIdentifierType: taxIdentifierType,
				TaxIdentifierNo:   strings.TrimSpace(data.TaxIdentifierNo),
			}
			if _, err := service.OutletRepository.CreateOutletTax(taxModel); err != nil {
				return fmt.Errorf("failed to save outlet tax: %w", err)
			}
		}
	}

	return nil
}

func (service *outletServiceImpl) mapRowToOutletStruct(row []string) (entity.OutletTemp, error) {
	headerIndex := make(map[string]int, len(outletExportHeaders))
	for i, key := range outletExportHeaders {
		headerIndex[key] = i
	}

	get := func(key string) string {
		if idx, ok := headerIndex[key]; ok {
			if idx >= 0 && idx < len(row) {
				return row[idx]
			}
		}
		return ""
	}

	data := entity.OutletTemp{
		OutletCode:              get("outlet_code"),
		OutletName:              get("outlet_name"),
		IsActive:                get("is_active"),
		OutletStatus:            get("outlet_status"),
		IsEmbBail:               get("is_pkp_outlet"),
		Address1:                get("address1"),
		OutletProvince:          get("outlet_province"),
		OutletRegency:           get("outlet_regency"),
		OutletSubDistrict:       get("outlet_sub_district"),
		OutletWard:              get("outlet_ward"),
		ZipCode:                 get("zip_code"),
		OtLocName:               get("ot_loc_name"),
		Longitude:               get("longitude"),
		Latitude:                get("latitude"),
		BuildingOwn:             get("building_own"),
		OutletEstablishmentDate: get("outlet_establishment_date"),
		ContactName:             get("contact_name"),
		JobTitle:                get("job_title"),
		IdentityType:            get("identity_type"),
		IdentityNo:              get("identity_no"),
		ContactPhoneNo:          get("contact_phone_no"),
		ContactIsWaNo:           get("contact_is_wa_no"),
		ContactWaNo:             get("contact_wa_no"),
		ContactEmail:            get("contact_email"),
		TaxInvoiceFormName:      get("tax_invoice_form_name"),
		TaxType:                 get("tax_identifier_type"),
		TaxIdentifierType:       get("tax_identifier_type"),
		TaxIdentifierNo:         get("tax_identifier_no"),
		Nitku:                   get("nitku"),
		TaxName:                 get("tax_name"),
		AddressTax:              get("address_tax"),
		PhoneNo:                 get("phone_no"),
		FaxNo:                   get("fax_no"),
		Barcode:                 get("barcode"),
		DiscGrpName:             get("disc_grp_name"),
		OtGrpName:               get("ot_grp_name"),
		IsContraBon:             get("is_contra_bon"),
		PriceGrpName:            get("price_grp_name"),
		DistrictName:            get("district_name"),
		IndustryName:            get("industry_name"),
		OtClassName:             get("ot_class_name"),
		OtTypeName:              get("ot_type_name"),
		MarketName:              get("market_name"),
		AgentFrom:               get("agent_from"),
		DelvAddr1:               get("delv_addr1"),
		DelvProvince:            get("delv_province"),
		DelvRegency:             get("delv_regency"),
		DelvSubDistrict:         get("delv_sub_district"),
		DelvWard:                get("delv_ward"),
		DelvLongitude:           get("delv_longitude"),
		DelvLatitude:            get("delv_latitude"),
		DelvZipCode:             get("delv_zip_code"),
		DelvIsSameAddr:          get("delv_is_same_addr"),
		InvAddr1:                get("inv_addr1"),
		InvProvince:             get("inv_province"),
		InvRegency:              get("inv_regency"),
		InvSubDistrict:          get("inv_sub_district"),
		InvWard:                 get("inv_ward"),
		InvZipCode:              get("inv_zip_code"),
		InvIsSameAddr:           get("inv_is_same_addr"),
		PaymentTypeName:         get("payment_type_name"),
		ArStatusName:            get("ar_status_name"),
		BankName:                get("bank_name"),
		AccountNo:               get("account_no"),
		AccountName:             get("account_name"),
		Top:                     get("top"),
		CreditLimitTypeName:     get("credit_limit_type_name"),
		CreditLimit:             get("credit_limit"),
		CreditLimitActionName:   get("credit_limit_action_name"),
		SalesInvLimitTypeName:   get("sales_inv_limit_type_name"),
		SalesInvLimit:           get("sales_inv_limit"),
		SalesInvLimitActionName: get("sales_inv_limit_action_name"),
		ObsTypeName:             get("obs_type_name"),
		Obs:                     get("obs"),
		ObsLimitActionName:      get("obs_limit_action_name"),
	}

	if strings.TrimSpace(data.IsActive) == "" {
		data.IsActive = "true"
	}

	return data, nil
}

func (s *outletServiceImpl) mapRowToOutletStructReupload(row []string, headerMap map[string]int) (entity.OutletTemp, error) {
	logrus.Info(row, "\n==========================================================\n")
	get := func(key string) string {
		if idx, ok := headerMap[canonicalizeHeader(key)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	data := entity.OutletTemp{
		OutletCode:   get("outlet_code"),
		OutletName:   get("outlet_name"),
		Barcode:      get("barcode"),
		OutletStatus: get("outlet_status"),
		Address1:     get("address1"),
		Address2:     get("address2"),
		City:         get("city"),
		ZipCode:      get("zip_code"),
		PhoneNo:      get("phone_no"),
		WaNo:         get("wa_no"),
		FaxNo:        get("fax_no"),
		Email:        get("email"),

		DiscGrpId:   get("disc_grp_id"),
		DiscGrpCode: get("disc_grp_code"),
		DiscGrpName: get("disc_grp_name"),

		OtLocId:   get("ot_loc_id"),
		OtLocCode: get("ot_loc_code"),
		OtLocName: get("ot_loc_name"),

		OtGrpId:   get("ot_grp_id"),
		OtGrpCode: get("ot_grp_code"),
		OtGrpName: get("ot_grp_name"),

		PriceGrpId:   get("price_grp_id"),
		PriceGrpCode: get("price_grp_code"),
		PriceGrpName: get("price_grp_name"),

		DistrictId:   get("district_id"),
		DistrictCode: get("district_code"),
		DistrictName: get("district_name"),

		OutletProvinceId:    get("outlet_province_id"),
		OutletProvince:      get("outlet_province"),
		OutletRegencyId:     get("outlet_regency_id"),
		OutletRegency:       get("outlet_regency"),
		OutletSubDistrictId: get("outlet_sub_district_id"),
		OutletSubDistrict:   get("outlet_sub_district"),
		OutletWardId:        get("outlet_ward_id"),
		OutletWard:          get("outlet_ward"),

		BeatId:  get("beat_id"),
		SbeatId: get("sbeat_id"),

		OtClassId:   get("ot_class_id"),
		OtClassCode: get("ot_class_code"),
		OtClassName: get("ot_class_name"),

		IndustryId:   get("industry_id"),
		IndustryCode: get("industry_code"),
		IndustryName: get("industry_name"),

		MarketId:   get("market_id"),
		MarketCode: get("market_code"),
		MarketName: get("market_name"),

		Top:             get("top"),
		PaymentType:     get("payment_type"),
		PaymentTypeName: get("payment_type_name"),
		IsContraBon:     get("is_contra_bon"),

		// dst isi semua kolom lain pakai get("nama_kolom")
		// ...
		BankId:   get("bank_id"),
		BankCode: get("bank_code"),
		BankName: get("bank_name"),

		AccountNo:      get("account_no"),
		AccountName:    get("account_name"),
		ContactName:    get("contact_name"),
		JobTitle:       get("job_title"),
		ContactPhoneNo: get("contact_phone_no"),
		ContactWaNo:    get("contact_wa_no"),
		ContactEmail:   get("contact_email"),
		IdentityNo:     get("identity_no"),
		ContactIsWaNo:  get("contact_is_wa_no"),
		IdentityType:   get("identity_type"),
		FaxNumber:      get("fax_number"),

		TaxInvoiceId:      get("tax_invoice_id"),
		IsEmbBail2:        get("is_emb_bail2"),
		TaxNo2:            get("tax_no2"),
		TaxName2:          get("tax_name2"),
		TaxCity2:          get("tax_city2"),
		TaxAddr1_2:        get("tax_addr1_2"),
		TaxAddr2_2:        get("tax_addr2_2"),
		TaxType:           get("tax_type"),
		Nitku:             get("nitku"),
		AddressTax:        get("address_tax"),
		TaxIdentifierType: get("tax_identifier_type"),
		TaxIdentifierNo:   get("tax_identifier_no"),
	}
	if strings.TrimSpace(data.IsActive) == "" {
		data.IsActive = "true"
	}
	logrus.Info(data)

	// validasi minimal
	// if data.OutletCode == "" || data.OutletName == "" {
	// 	return data, errors.New("outlet_code and outlet_name cannot be empty")
	// }
	// if data.OtLocCode == "" || data.OtLocName == "" {
	// 	return data, errors.New("ot_loc_code and ot_loc_name cannot be empty")
	// }
	// if data.OtTypeCode == "" || data.OtTypeName == "" {
	// 	return data, errors.New("ot_type_code and ot_type_name cannot be empty")
	// }
	// if data.OtGrpCode == "" || data.OtGrpName == "" {
	// 	return data, errors.New("ot_grp_code and ot_grp_name cannot be empty")
	// }
	// if data.DistrictCode == "" || data.DistrictName == "" {
	// 	return data, errors.New("district_code and district_name cannot be empty")
	// }
	// if data.OtClassCode == "" || data.OtClassName == "" {
	// 	return data, errors.New("ot_class_code and ot_class_name cannot be empty")
	// }
	// if data.DiscGrpCode == "" || data.DiscGrpName == "" {
	// 	return data, errors.New("disc_grp_code and disc_grp_name cannot be empty")
	// }
	// if data.MarketCode == "" || data.MarketName == "" {
	// 	return data, errors.New("market_code and market_name cannot be empty")
	// }
	// if data.IndustryCode == "" || data.IndustryName == "" {
	// 	return data, errors.New("industry_code and industry_name cannot be empty")
	// }
	// if data.BankCode == "" || data.BankName == "" {
	// 	return data, errors.New("bank_code and bank_name cannot be empty")
	// }
	// if data.PriceGrpCode == "" || data.PriceGrpName == "" {
	// 	return data, errors.New("price_grp_code and price_grp_name cannot be empty")
	// }

	return data, nil
}

// --- HELPER FUNCTIONS (GENERIC) ---

// (Fungsi-fungsi ini bisa dipindah ke package terpisah)
// func (service *outletServiceImpl) createTemplateUpdateXLSX(headers []string, data map[string][][]string) (*bytes.Buffer, error) {
// 	f := excelize.NewFile()
// 	sheetName := "Template Update"
// 	index, _ := f.NewSheet(sheetName)

// 	style, _ := f.NewStyle(&excelize.Style{
// 		Font: &excelize.Font{Bold: true},
// 		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
// 	})

// 	for i, header := range headers {
// 		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
// 		f.SetCellValue(sheetName, cell, header)
// 		f.SetCellStyle(sheetName, cell, cell, style)
// 	}

// 	colOffset := 0
// 	for _, fieldGroup := range headers {
// 		// Asumsi group header (e.g., "area_id", "area_code", "area_name")
// 		// kita hanya perlu key utamanya ("area_id") untuk ambil data.
// 		mainFieldKey := fieldGroup
// 		if rows, ok := data[mainFieldKey]; ok {
// 			for rowIdx, rowData := range rows {
// 				for valIdx, value := range rowData {
// 					cell, _ := excelize.CoordinatesToCellName(colOffset+valIdx+1, rowIdx+2)
// 					f.SetCellValue(sheetName, cell, value)
// 				}
// 			}
// 		}
// 		// Ini perlu disesuaikan tergantung struktur `data` yang dikembalikan repository
// 		// Contoh sederhana ini mengasumsikan setiap field di `headers` adalah key di `data`.
// 		colOffset++
// 	}

// 	f.SetActiveSheet(index)
// 	f.DeleteSheet("Sheet1")
// 	return f.WriteToBuffer()
// }

// func (service *outletServiceImpl) createTemplateUpdateCSV(headers []string, data map[string][][]string) (*bytes.Buffer, error) {
// 	buffer := new(bytes.Buffer)
// 	writer := csv.NewWriter(buffer)
// 	writer.Comma = ';'
// 	writer.Write(headers)

// 	maxRows := 0
// 	for _, fieldData := range data {
// 		if len(fieldData) > maxRows {
// 			maxRows = len(fieldData)
// 		}
// 	}

// 	for i := 0; i < maxRows; i++ {
// 		record := make([]string, len(headers))
// 		// Logika ini perlu penyesuaian besar tergantung bagaimana data dari repo distrukturkan
// 		// dan bagaimana `headers` berhubungan dengan `data`.
// 		// Ini adalah implementasi yang disederhanakan.
// 		writer.Write(record)
// 	}

// 	writer.Flush()
// 	return buffer, writer.Error()
// }

func (s *outletServiceImpl) ImportUpdateXLSX(req entity.ImportRequest) error {
	f, err := excelize.OpenReader(req.File)
	if err != nil {
		return errors.New("gagal membuka file Excel (XLS/XLSX)")
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return errors.New("gagal membaca baris data")
	}

	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	headers := rows[0]
	headerMap := make(map[string]int, len(headers))
	for idx, h := range headers {
		headerMap[canonicalizeHeader(h)] = idx
	}

	// ✅ Validasi panjang karakter dulu (sebelum masuk ke DB)
	rowValidationErrors := make(map[int][]string)
	var validationSummaries []string
	for i, rowValues := range rows[1:] {
		rowMap := make(map[string]string)
		for j, val := range rowValues {
			if j < len(headers) {
				cleaned := canonicalizeHeader(headers[j])
				rowMap[cleaned] = strings.TrimSpace(val)
			}
		}
		rowIdx := i + 1
		rowNumber := i + 2
		if err := validateMaxLength(rowMap); err != nil {
			rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
			validationSummaries = append(validationSummaries, fmt.Sprintf("baris %d: %v", rowNumber, err))
		}
	}

	if len(rowValidationErrors) > 0 {
		historyId, err := s.OutletRepository.CreateImportHistory("outlet-update", req.Filename, req.CustId, req.UserId, len(rows)-1)
		if err != nil {
			return err
		}
		for rowIdx, msgs := range rowValidationErrors {
			rowData := rows[rowIdx]
			temp := s.mapRowToOutletUpdateTemp(historyId, req.CustId, rowData, headerMap, strings.Join(msgs, "; "))
			if err := s.OutletRepository.InsertOutletUpdateTemp(temp); err != nil {
				logrus.Errorf("gagal simpan ke outlet_update_temp baris %d: %v", rowIdx+1, err)
			}
		}
		_ = s.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rowValidationErrors), true)
		return fmt.Errorf("validasi gagal (%d baris): %s", len(rowValidationErrors), strings.Join(validationSummaries, "; "))
	}

	historyId, err := s.OutletRepository.CreateImportHistory("outlet-update", req.Filename, req.CustId, req.UserId, len(rows)-1)
	if err != nil {
		return err
	}

	go func() {
		if err := s.processUpdateRows(req, rows, historyId); err != nil {
			_ = s.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rows)-1, true)
			logrus.Errorf("Import Update XLSX async gagal (historyId=%d): %v", historyId, err)
			return
		}
		logrus.Infof("Import Update XLSX async success, historyId=%d", historyId)
	}()

	return nil
}

func (s *outletServiceImpl) hasHeader(headers []string, key string) bool {
	for _, h := range headers {
		if h == key {
			return true
		}
	}
	return false
}

func (s *outletServiceImpl) ImportUpdateCSV(req entity.ImportRequest) error {
	reader := csv.NewReader(req.File)
	reader.TrimLeadingSpace = true
	reader.Comma = ';'

	rows, err := reader.ReadAll()
	if err != nil {
		return errors.New("gagal membaca file CSV")
	}

	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	headers := rows[0]
	headerMap := make(map[string]int, len(headers))
	for idx, h := range headers {
		headerMap[canonicalizeHeader(h)] = idx
	}

	// ✅ Validasi panjang karakter sebelum masuk ke DB
	rowValidationErrors := make(map[int][]string)
	var validationSummaries []string
	for i, rowValues := range rows[1:] {
		rowMap := make(map[string]string)
		for j, val := range rowValues {
			if j < len(headers) {
				cleaned := canonicalizeHeader(headers[j])
				rowMap[cleaned] = strings.TrimSpace(val)
			}
		}
		rowIdx := i + 1
		rowNumber := i + 2
		if err := validateMaxLength(rowMap); err != nil {
			rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
			validationSummaries = append(validationSummaries, fmt.Sprintf("baris %d: %v", rowNumber, err))
		}
	}

	if len(rowValidationErrors) > 0 {
		historyId, err := s.OutletRepository.CreateImportHistory("outlet-update", req.Filename, req.CustId, req.UserId, len(rows)-1)
		if err != nil {
			return err
		}
		for rowIdx, msgs := range rowValidationErrors {
			rowData := rows[rowIdx]
			temp := s.mapRowToOutletUpdateTemp(historyId, req.CustId, rowData, headerMap, strings.Join(msgs, "; "))
			if err := s.OutletRepository.InsertOutletUpdateTemp(temp); err != nil {
				logrus.Errorf("gagal simpan ke outlet_update_temp baris %d: %v", rowIdx+1, err)
			}
		}
		_ = s.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rowValidationErrors), true)
		return fmt.Errorf("validasi gagal (%d baris): %s", len(rowValidationErrors), strings.Join(validationSummaries, "; "))
	}

	historyId, err := s.OutletRepository.CreateImportHistory("outlet-update", req.Filename, req.CustId, req.UserId, len(rows)-1)
	if err != nil {
		return err
	}

	go func() {
		if err := s.processUpdateRows(req, rows, historyId); err != nil {
			_ = s.OutletRepository.UpdateImportHistoryOutlet(historyId, 0, len(rows)-1, true)
			logrus.Errorf("Import Update CSV async gagal (historyId=%d): %v", historyId, err)
			return
		}
		logrus.Infof("Import Update CSV async success, historyId=%d", historyId)
	}()

	return nil
}

func (s *outletServiceImpl) ImportUpdateXLS(req entity.ImportRequest) error {
	// For XLS files, we'll treat them the same as XLSX since we're using excelize
	// which handles both formats
	return s.ImportUpdateXLSX(req)
}

// ReuploadImportUpdateFile processes a reupload file for outlet-update failed rows
func (s *outletServiceImpl) ReuploadImportUpdateFile(custId string, historyId int64, req entity.ImportRequest) error {
	var rows [][]string
	var err error
	success := 0
	switch strings.ToLower(req.Format) {
	case "csv":
		reader := csv.NewReader(req.File)
		reader.TrimLeadingSpace = true
		reader.Comma = ';'
		rows, err = reader.ReadAll()
	default:
		f, e := excelize.OpenReader(req.File)
		if e != nil {
			return errors.New("gagal membuka file Excel (XLS/XLSX)")
		}
		defer f.Close()
		rows, err = f.GetRows(f.GetSheetName(0))
	}
	if err != nil {
		return errors.New("gagal membaca file")
	}
	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	header := rows[0]
	// Build headerMap with canonical + normalized display variants
	headerMap := map[string]int{}
	normalizeHeaderKey := func(h string) string {
		k := strings.TrimSpace(strings.ToLower(h))
		replacer := strings.NewReplacer(" ", "_", "-", "_", "/", "_")
		k = replacer.Replace(k)
		for strings.Contains(k, "__") {
			k = strings.ReplaceAll(k, "__", "_")
		}
		return k
	}
	for i, h := range header {
		canon := canonicalizeHeader(h)
		headerMap[canon] = i
		headerMap[normalizeHeaderKey(h)] = i // allow display header variants
	}
	// Synchronous validation: collect all mandatory errors to return in one response
	rowValidationErrors := make(map[int][]string)
	var validationErrs []string
	for i, row := range rows[1:] {
		rowMap := make(map[string]string)
		for j, val := range row {
			if j < len(header) {
				orig := header[j]
				rowMap[strings.ToLower(orig)] = val
				rowMap[normalizeHeaderKey(orig)] = val
				rowMap[canonicalizeHeader(orig)] = val
			}
		}
		if err := validateCreateMandatory(rowMap); err != nil {
			rowIdx := i + 1
			rowValidationErrors[rowIdx] = append(rowValidationErrors[rowIdx], err.Error())
			validationErrs = append(validationErrs, fmt.Sprintf("baris %d: %v", i+2, err))
		}
	}
	if len(rowValidationErrors) > 0 {
		newHistoryId, err := s.OutletRepository.CreateImportHistory("outlet_update", req.Filename, req.CustId, req.UserId, len(rows)-1)
		if err != nil {
			return err
		}
		for rowIdx, msgs := range rowValidationErrors {
			rowData := rows[rowIdx]
			temp := s.mapRowToOutletUpdateTemp(newHistoryId, req.CustId, rowData, headerMap, strings.Join(msgs, "; "))
			if err := s.OutletRepository.InsertOutletUpdateTemp(temp); err != nil {
				logrus.Errorf("gagal simpan ke outlet_update_temp baris %d: %v", rowIdx+1, err)
			}
		}
		_ = s.OutletRepository.UpdateImportHistoryOutlet(newHistoryId, 0, len(rowValidationErrors), true)
		return fmt.Errorf("validasi gagal (%d baris): %s", len(rowValidationErrors), strings.Join(validationErrs, "; "))
	}
	get := func(row []string, key string) string {
		if idx, ok := headerMap[canonicalizeHeader(key)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	total, err := s.OutletRepository.CountOutletUpdateTemp(historyId)
	if err != nil {
		return err
	}
	newHistoryId, err := s.OutletRepository.CreateImportHistory("outlet_update", req.Filename, req.CustId, req.UserId, total)

	// local helpers to mirror processUpdateRows parsing
	parseInt64 := func(val string) int64 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, e := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if e != nil {
			return 0
		}
		return v
	}
	parseInt32 := func(val string) int32 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, e := strconv.ParseInt(strings.TrimSpace(val), 10, 32)
		if e != nil {
			return 0
		}
		return int32(v)
	}
	parseInt32WithCodes := func(val string, codes map[string]int) int32 {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return 0
		}
		if v, e := strconv.ParseInt(trimmed, 10, 32); e == nil {
			return int32(v)
		}
		if codes != nil {
			if code, ok := codes[normEnum(trimmed)]; ok {
				return int32(code)
			}
		}
		return 0
	}
	parseInt64WithCodes := func(val string, codes map[string]int) int64 {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return 0
		}
		if v, e := strconv.ParseInt(trimmed, 10, 64); e == nil {
			return v
		}
		if codes != nil {
			if code, ok := codes[normEnum(trimmed)]; ok {
				return int64(code)
			}
		}
		return 0
	}
	parseFloat64 := func(val string) float64 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, e := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if e != nil {
			return 0
		}
		return v
	}
	parseBoolWithDefault := func(val string, def bool) bool {
		t := strings.ToLower(strings.TrimSpace(val))
		switch t {
		case "1", "true", "t", "ya", "yes", "y":
			return true
		case "0", "false", "f", "tidak", "no", "n":
			return false
		case "":
			return def
		}
		return def
	}
	parseActiveFlag := func(val string) *bool {
		t := strings.ToLower(strings.TrimSpace(val))
		if t == "" {
			return nil
		}
		switch t {
		case "1", "true", "t", "ya", "yes", "y", "active", "aktif":
			v := true
			return &v
		case "0", "false", "f", "tidak", "no", "n", "inactive", "deactive", "nonactive", "non-active", "non active", "deactivated", "nonaktif":
			v := false
			return &v
		default:
			return nil
		}
	}
	parseDate := func(val string) (time.Time, bool) {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return time.Time{}, false
		}
		layouts := []string{"2006-01-02", time.RFC3339, "2006-01-02 15:04:05", "02/01/2006", "02-01-2006"}
		for _, layout := range layouts {
			if t, e := time.Parse(layout, trimmed); e == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}

	go func() {
		for i, row := range rows[1:] {
			_ = i
			ctid := get(row, "ctid")
			if ctid == "" {
				continue
			}
			// bank
			if code := get(row, "bank_code"); code != "" {
				if id, err := s.OutletRepository.FindBankIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportBank(custId, code, get(row, "bank_name"), id)
				}
				// ====== update outlet main fields (mirror processUpdateRows) ======
				{
					// resolve outlet id from identifiers or fallback to code
					resolveOutletId := func() (int64, error) {
						if v := strings.TrimSpace(get(row, "outlet_id")); v != "" {
							if oid, perr := strconv.ParseInt(v, 10, 64); perr == nil {
								out, e := s.OutletRepository.FindOneByOutletIdAndCustId(oid, req.CustId, req.ParentCustId)
								if e == nil && out.OutletId > 0 {
									return int64(out.OutletId), nil
								}
								return 0, fmt.Errorf("outlet_id %d tidak ditemukan", oid)
							}
						}
						if code := strings.TrimSpace(get(row, "outlet_code")); code != "" {
							out, e := s.OutletRepository.FindOneByOutletCodeAndCustId(code, req.CustId, req.ParentCustId)
							if e == nil && out.OutletId > 0 {
								return int64(out.OutletId), nil
							}
							return 0, fmt.Errorf("outlet_code %s tidak ditemukan", code)
						}
						return 0, errors.New("wajib isi outlet_id atau outlet_code")
					}

					otGrpId, _ := s.OutletRepository.FindOutletGroupIdByName(req.ParentCustId, get(row, "ot_grp_name"))
					discGrpId, _ := s.OutletRepository.FindDiscGroupIdByName(req.ParentCustId, get(row, "disc_grp_name"))
					otLocId, _ := s.OutletRepository.FindOutletLocIdByName(req.ParentCustId, get(row, "ot_loc_name"))
					priceGrpId, _ := s.OutletRepository.FindPriceGroupIdByName(req.ParentCustId, get(row, "price_grp_name"))
					districtId, _ := s.OutletRepository.FindDistrictIdByName(req.ParentCustId, get(row, "district_name"))
					industryId, _ := s.OutletRepository.FindIndustryIdByName(req.ParentCustId, get(row, "industry_name"))
					otClassId, _ := s.OutletRepository.FindOutletClassIdByName(req.ParentCustId, get(row, "ot_class_name"))
					marketId, _ := s.OutletRepository.FindMarketIdByName(req.ParentCustId, get(row, "market_name"))
					otTypeId, _ := s.OutletRepository.FindOutletTypeIdByName(req.ParentCustId, get(row, "ot_type_name"))
					wardId, _ := s.OutletRepository.GetWardIdByName(req.ParentCustId, get(row, "outlet_ward"))

					outletId, rerr := resolveOutletId()
					if rerr == nil {
						activeFlag := parseActiveFlag(get(row, "is_active"))
						if activeFlag == nil {
							activeFlag = parseActiveFlag(get(row, "outlet_status"))
						}
						updateReq := entity.ProcessedUpdateOutlet{
							CustId:              req.CustId,
							OutletId:            outletId,
							OutletCode:          get(row, "outlet_code"),
							Barcode:             get(row, "barcode"),
							OutletName:          get(row, "outlet_name"),
							OutletStatus:        parseInt32WithCodes(get(row, "outlet_status"), outletStatusCodes),
							Address1:            get(row, "address1"),
							Address2:            get(row, "address2"),
							City:                get(row, "city"),
							ZipCode:             get(row, "zip_code"),
							PhoneNo:             get(row, "phone_no"),
							WaNo:                get(row, "wa_no"),
							FaxNo:               get(row, "fax_no"),
							Email:               get(row, "email"),
							DiscGrpId:           discGrpId,
							OtLocId:             otLocId,
							OtGrpId:             otGrpId,
							PriceGrpId:          priceGrpId,
							DistrictId:          districtId,
							BeatId:              parseInt64(get(row, "beat_id")),
							SbeatId:             parseInt64(get(row, "sbeat_id")),
							OtClassId:           otClassId,
							IndustryId:          industryId,
							MarketId:            marketId,
							OtTypeId:            otTypeId,
							OutletWardId:        wardId, //parseInt64(get(row, "outlet_ward_id")),
							Top:                 parseInt32(get(row, "top")),
							PaymentType:         parseInt32WithCodes(get(row, "payment_type"), paymentTypeCodes),
							IsContraBon:         parseBoolWithDefault(get(row, "is_contra_bon"), true),
							CreditLimitType:     parseInt32WithCodes(get(row, "credit_limit_type"), creditLimitTypeCodes),
							CreditLimit:         parseFloat64(get(row, "credit_limit")),
							SalesInvLimitType:   parseInt32WithCodes(get(row, "sales_inv_limit_type"), salesInvLimitTypeCodes),
							SalesInvLimit:       parseInt32(get(row, "sales_inv_limit")),
							AvgSalesWeek:        parseFloat64(get(row, "avg_sales_week")),
							AvgSalesMonth:       parseFloat64(get(row, "avg_sales_month")),
							CreditLimitAction:   parseInt64WithCodes(get(row, "credit_limit_action"), limitActionCodes),
							SalesInvLimitAction: parseInt64WithCodes(get(row, "sales_inv_limit_action"), limitActionCodes),
							PluGrpId:            parseInt64(get(row, "plu_grp_id")),
							ConvGrpId:           parseInt64(get(row, "conv_grp_id")),
							DiscInvId:           parseInt64(get(row, "disc_inv_id")),
							AgentFrom:           get(row, "agent_from"),
							FirstWeekNo:         parseInt32(get(row, "first_week_no")),
							BuildingOwn:         parseInt32(get(row, "building_own")),
							ArStatus:            parseInt32WithCodes(get(row, "ar_status"), arStatusCodes),
							ArTotal:             parseFloat64(get(row, "ar_total")),
							IsEmbBail:           parseBoolWithDefault(get(row, "is_emb_bail"), false),
							IsActive:            activeFlag,
							TaxName:             get(row, "tax_name"),
							TaxAddr1:            get(row, "tax_addr1"),
							TaxAddr2:            get(row, "tax_addr2"),
							TaxCity:             get(row, "tax_city"),
							TaxNo:               get(row, "tax_no"),
							TaxInvoiceForm:      parseInt32WithCodes(get(row, "tax_invoice_form"), taxInvoiceFormCodes),
							OwnerName:           get(row, "owner_name"),
							OwnerAddr1:          get(row, "owner_addr1"),
							OwnerAddr2:          get(row, "owner_addr2"),
							OwnerCity:           get(row, "owner_city"),
							OwnerPhoneNo:        get(row, "owner_phone_no"),
							OwnerIdNo:           get(row, "owner_id_no"),
							DelvAddr1:           get(row, "delv_addr1"),
							DelvAddr2:           get(row, "delv_addr2"),
							DelvCity:            get(row, "delv_city"),
							DelvCity2:           get(row, "delv_city2"),
							DelvWardId:          get(row, "delv_ward_id"),
							DelvWardId2:         get(row, "delv_ward_id2"),
							DelvZipCode:         get(row, "delv_zip_code"),
							DelvZipCode2:        get(row, "delv_zip_code2"),
							DelvIsSameAddr:      parseBoolWithDefault(get(row, "delv_is_same_addr"), false),
							DelvLatitude:        get(row, "delv_latitude"),
							DelvLongitude:       get(row, "delv_longitude"),
							DelvLatitude2:       get(row, "delv_latitude2"),
							DelvLongitude2:      get(row, "delv_longitude2"),
							InvAddr1:            get(row, "inv_addr1"),
							InvAddr2:            get(row, "inv_addr2"),
							InvCity:             get(row, "inv_city"),
							InvWardId:           get(row, "inv_ward_id"),
							InvZipCode:          get(row, "inv_zip_code"),
							InvIsSameAddr:       parseBoolWithDefault(get(row, "inv_is_same_addr"), false),
							Latitude:            get(row, "latitude"),
							Longitude:           get(row, "longitude"),
							ImageUrl:            get(row, "image_url"),
							IsObs:               parseBoolWithDefault(get(row, "is_obs"), false),
							Obs:                 parseInt64(get(row, "obs")),
							ObsType:             parseInt64WithCodes(get(row, "obs_type"), obsTypeCodes),
							ObsLimitAction:      parseInt64WithCodes(get(row, "obs_limit_action"), limitActionCodes),
							IsWaNo:              parseBoolWithDefault(get(row, "is_wa_no"), false),
							VerificationStatus:  parseInt32WithCodes(get(row, "verification_status"), verificationStatusCodes),
							VerifiedBy:          parseInt64(get(row, "verified_by")),
							UpdatedBy:           req.UserId,
						}
						if v := strings.TrimSpace(get(row, "is_pkp_outlet")); v != "" {
							b := parseBoolWithDefault(v, false)
							updateReq.IsPkpOutlet = &b
						}
						if t, ok := parseDate(get(row, "first_trans_date")); ok {
							updateReq.FirstTransDate = t
						}
						if t, ok := parseDate(get(row, "last_trans_date")); ok {
							updateReq.LastTransDate = t
						}
						if t, ok := parseDate(get(row, "ot_start_date")); ok {
							updateReq.OtStartDate = t
						}
						if t, ok := parseDate(get(row, "ot_reg_date")); ok {
							updateReq.OtRegDate = t
						}
						if t, ok := parseDate(get(row, "dob")); ok {
							updateReq.Dob = t
						}
						if t, ok := parseDate(get(row, "closed_date")); ok {
							updateReq.ClosedDate = t
						}
						if t, ok := parseDate(get(row, "outlet_establishment_date")); ok {
							updateReq.OutletEstablishmentDate = t
						}
						if t, ok := parseDate(get(row, "verified_at")); ok {
							updateReq.VerifiedAt = t
						}
						if uerr := s.OutletRepository.UpdateOutletImport(updateReq); uerr != nil {
							logrus.Errorf("baris %d: gagal memperbarui data outlet: %v", i+2, uerr)
						}
					}
				}
			}
			// outlet class
			if code := get(row, "ot_class_code"); code != "" {
				if id, err := s.OutletRepository.FindOutletClassIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportOutletClass(custId, code, get(row, "ot_class_name"), id)
				}
			}
			// outlet group
			if code := get(row, "outlet_grp_code"); code != "" {
				if id, err := s.OutletRepository.FindOutletGroupIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportOutletGroup(custId, code, get(row, "outlet_grp_name"), id)
				}
			}
			// outlet location
			if code := get(row, "outlet_loc_code"); code != "" {
				if id, err := s.OutletRepository.FindOutletLocIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportOutletLoc(custId, code, get(row, "outlet_loc_name"), id)
				}
			}
			// outlet type
			if code := get(row, "outlet_type_code"); code != "" {
				if id, err := s.OutletRepository.FindOutletTypeIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportOutletType(custId, code, get(row, "outlet_type_name"), id)
				}
			}
			// district
			if code := get(row, "district_code"); code != "" {
				if id, err := s.OutletRepository.FindDistrictIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportDistrict(custId, code, get(row, "district_name"), id)
				}
			}
			// province
			if code := get(row, "outlet_province_id"); code != "" {
				_ = s.OutletRepository.UpdateImportProvince(custId, code, get(row, "outlet_province"))
			}
			// regency
			if code := get(row, "outlet_regency_id"); code != "" {
				_ = s.OutletRepository.UpdateImportRegency(custId, code, get(row, "outlet_regency"), get(row, "outlet_province_id"))
			}
			// sub_district
			if code := get(row, "outlet_sub_district_id"); code != "" {
				_ = s.OutletRepository.UpdateImportSubDistrict(custId, code, get(row, "outlet_sub_district"), get(row, "outlet_province_id"), get(row, "outlet_regency_id"))
			}
			// ward
			if code := get(row, "outlet_ward_id"); code != "" {
				_ = s.OutletRepository.UpdateImportWard(custId, code, get(row, "outlet_ward"), get(row, "outlet_province_id"), get(row, "outlet_regency_id"), get(row, "outlet_sub_district_id"))
			}

			// disc group
			if code := get(row, "disc_grp_code"); code != "" {
				if id, err := s.OutletRepository.FindDiscGroupIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportDiscGroup(custId, code, get(row, "disc_grp_name"), id)
				}
			}
			// market
			if code := get(row, "market_code"); code != "" {
				if id, err := s.OutletRepository.FindMarketIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportMarket(custId, code, get(row, "market_name"), id)
				}
			}
			// industry
			if code := get(row, "industry_code"); code != "" {
				if id, err := s.OutletRepository.FindIndustryIdByCode(custId, code); err == nil && id > 0 {
					_ = s.OutletRepository.UpdateImportIndustry(custId, code, get(row, "industry_name"), id)
				}
			}
			// price group
			if code := get(row, "price_grp_code"); code != "" {
				if id, err := s.OutletRepository.FindPriceGroupIdByCode(custId, code); err == nil && id > 0 {
					_, _ = s.OutletRepository.UpdateImportPriceGroup(custId, id, code, get(row, "price_grp_name"))
				}
			}
			success++

			// delete temp row on success basis (best effort here)
			_ = s.OutletRepository.DeleteOutletUpdateTempByCtid(ctid)
		}

		// failed, err := s.OutletRepository.CountOutletUpdateTemp(historyId)
		// if err != nil { return err }
		// total, err := s.OutletRepository.GetImportTotalData(historyId)
		// if err != nil { return err }
		newFailed := total - success
		if success < 0 {
			success = 0
		}
		err = s.OutletRepository.UpdateImportHistoryOutlet(newHistoryId, success, newFailed, newFailed > 0)
		if err != nil {
			logrus.Info("error update import history")
		}
	}()
	return nil
	// return nil
}

func (s *outletServiceImpl) ReuploadImportInsertFile(custId string, historyId int64, req entity.ImportRequest) error {
	var rows [][]string
	var err error
	switch strings.ToLower(req.Format) {
	case "csv":
		reader := csv.NewReader(req.File)
		reader.TrimLeadingSpace = true
		reader.Comma = ';'
		rows, err = reader.ReadAll()
	default:
		f, e := excelize.OpenReader(req.File)
		if e != nil {
			return errors.New("gagal membuka file Excel (XLS/XLSX)")
		}
		defer f.Close()
		rows, err = f.GetRows(f.GetSheetName(0))
	}
	if err != nil {
		return errors.New("gagal membaca file")
	}
	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	header := rows[0]
	headerMap := map[string]int{}
	for i, h := range header {
		headerMap[canonicalizeHeader(h)] = i
	}

	// Synchronous pre-validation of mandatory fields so controller can return failure
	normalizeHeaderKey := func(h string) string {
		k := strings.TrimSpace(strings.ToLower(h))
		replacer := strings.NewReplacer(" ", "_", "-", "_", "/", "_")
		k = replacer.Replace(k)
		for strings.Contains(k, "__") {
			k = strings.ReplaceAll(k, "__", "_")
		}
		return k
	}
	var validationErrors []string
	for i, row := range rows[1:] { // i=0 => sheet row 2
		rowMap := make(map[string]string)
		for j, val := range row {
			if j < len(header) {
				orig := header[j]
				rowMap[strings.ToLower(orig)] = val
				rowMap[normalizeHeaderKey(orig)] = val
				rowMap[canonicalizeHeader(orig)] = val
			}
		}
		var errVal error
		if req.IsImportNew {
			errVal = validateImportNewMandatory(rowMap)
		} else {
			errVal = validateCreateMandatory(rowMap)
		}
		if errVal != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("baris %d: %v", i+2, errVal))
		}
	}
	if len(validationErrors) > 0 {
		return fmt.Errorf("validasi gagal (%d baris): %s", len(validationErrors), strings.Join(validationErrors, "; "))
	}
	// get := func(row []string, key string) string {
	//     if idx, ok := headerMap[strings.ToLower(key)]; ok && idx < len(row) { return strings.TrimSpace(row[idx]) }
	//     return ""
	// }
	// create new import history for this reupload session (mirror update flow)
	total, err := s.OutletRepository.CountOutletTemp(historyId)
	if err != nil {
		return err
	}
	newHistoryId, err := s.OutletRepository.CreateImportHistory(importUploadType(req), req.Filename, req.CustId, req.UserId, total)
	if err != nil {
		return err
	}

	go func() {
		success := 0
		for i, row := range rows[1:] {
			rowNumber := i + 2
			_ = i
			// Rows already validated; just parse
			importData, err := s.mapRowToOutletStructReupload(row, headerMap)
			logrus.Info(err)
			if err != nil {
				// Record into outlet_temp failed rows to mirror ImportInsert
				if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
					logrus.Errorf("gagal mencatat hasil validasi reupload pada baris %d: %v", rowNumber, e)
				}
				continue
			}

			code := strings.TrimSpace(importData.OutletCode)
			if code != "" {
				exists, err := s.OutletRepository.CheckOutletCodeExists(req.CustId, code)
				if err != nil {
					logrus.Errorf("gagal memeriksa kode outlet pada baris %d: %v", rowNumber, err)
					importData.ErrorMessage = fmt.Sprintf("Gagal memeriksa kode outlet: %v", err)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat hasil pemeriksaan duplikasi pada baris %d: %v", rowNumber, e)
					}
					continue
				}
				if exists {
					importData.ErrorMessage = fmt.Sprintf("Kode outlet %s sudah terdaftar", code)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet duplikat pada baris %d: %v", rowNumber, e)
					}
					continue
				}
			}

			missingMasters := make([]string, 0, 20)
			missingMasterAdded := make(map[string]struct{})
			addMissingMaster := func(msg string) {
				msg = strings.TrimSpace(msg)
				if msg == "" {
					return
				}
				if _, exists := missingMasterAdded[msg]; exists {
					return
				}
				missingMasterAdded[msg] = struct{}{}
				missingMasters = append(missingMasters, msg)
			}

			if name := strings.TrimSpace(importData.OutletProvince); name != "" {
				provinceId, provinceErr := s.OutletRepository.GetProvinceIdByName(req.CustId, name)
				if provinceErr != nil {
					logrus.Errorf("gagal mengambil master outlet province pada baris %d: %v", rowNumber, provinceErr)
					prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
					importData.ErrorMessage = fmt.Sprintf("%s: gagal mengambil master outlet province: %v", prefix, provinceErr)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
					}
					continue
				}
				importData.OutletProvinceId = provinceId
				if provinceId == "" {
					addMissingMaster("Master outlet province tidak ditemukan")
				}
			}
			if name := strings.TrimSpace(importData.OutletRegency); name != "" {
				regencyId, regencyErr := s.OutletRepository.GetRegencyIdByName(req.CustId, name)
				if regencyErr != nil {
					logrus.Errorf("gagal mengambil master outlet regency pada baris %d: %v", rowNumber, regencyErr)
					prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
					importData.ErrorMessage = fmt.Sprintf("%s: gagal mengambil master outlet regency: %v", prefix, regencyErr)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
					}
					continue
				}
				importData.OutletRegencyId = regencyId
				if regencyId == "" {
					addMissingMaster("Master outlet regency tidak ditemukan")
				}
			}
			if name := strings.TrimSpace(importData.OutletSubDistrict); name != "" {
				subDistrictId, subErr := s.OutletRepository.GetSubDistrictIdByName(req.CustId, name)
				if subErr != nil {
					logrus.Errorf("gagal mengambil master outlet sub district pada baris %d: %v", rowNumber, subErr)
					prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
					importData.ErrorMessage = fmt.Sprintf("%s: gagal mengambil master outlet sub district: %v", prefix, subErr)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
					}
					continue
				}
				importData.OutletSubDistrictId = subDistrictId
				if subDistrictId == "" {
					addMissingMaster("Master outlet sub district tidak ditemukan")
				}
			}
			if name := strings.TrimSpace(importData.OutletWard); name != "" {
				wardId, wardErr := s.OutletRepository.GetWardIdByName(req.CustId, name)
				if wardErr != nil {
					logrus.Errorf("gagal mengambil master outlet ward pada baris %d: %v", rowNumber, wardErr)
					prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
					importData.ErrorMessage = fmt.Sprintf("%s: gagal mengambil master outlet ward: %v", prefix, wardErr)
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
					}
					continue
				}
				importData.OutletWardId = wardId
				if wardId == "" {
					addMissingMaster("Master outlet ward tidak ditemukan")
				}
			}

			resolveMaster := func(label, errorLabel, code, name string, findByCode func(string) (int64, error), findByName func(string) (int64, error)) (int64, error) {
				code = strings.TrimSpace(code)
				name = strings.TrimSpace(name)
				if code == "" && name == "" {
					return 0, nil
				}

				recordFatal := func(detail error) (int64, error) {
					prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
					message := fmt.Sprintf("%s: gagal mengambil master %s: %v", prefix, errorLabel, detail)
					importData.ErrorMessage = message
					if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
						logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
					}
					fullErr := fmt.Errorf("gagal mengambil master %s pada baris %d: %w", label, rowNumber, detail)
					return 0, fullErr
				}

				if code != "" && findByCode != nil {
					id, errLookup := findByCode(code)
					if errLookup == nil && id > 0 {
						return id, nil
					}
					if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
						return recordFatal(errLookup)
					}
				}

				if code == "" && name != "" && findByCode != nil {
					id, errLookup := findByCode(name)
					if errLookup == nil && id > 0 {
						return id, nil
					}
					if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
						return recordFatal(errLookup)
					}
				}

				if name != "" && findByName != nil {
					id, errLookup := findByName(name)
					if errLookup == nil && id > 0 {
						return id, nil
					}
					if errLookup != nil && !errors.Is(errLookup, sql.ErrNoRows) {
						return recordFatal(errLookup)
					}
				}

				addMissingMaster(fmt.Sprintf("Master %s tidak ditemukan", errorLabel))
				return 0, nil
			}

			var (
				discGrpId  int64
				otLocId    int64
				otGrpId    int64
				priceGrpId int64
				districtId int64
				otClassId  int64
				industryId int64
				marketId   int64
				otTypeId   int64
			)

			if id, err := resolveMaster(
				"disc group",
				"disc group",
				importData.DiscGrpCode,
				importData.DiscGrpName,
				func(code string) (int64, error) { return s.OutletRepository.FindDiscGroupIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindDiscGroupIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				discGrpId = id
			}

			if id, err := resolveMaster(
				"outlet location",
				"outlet location",
				importData.OtLocCode,
				importData.OtLocName,
				func(code string) (int64, error) { return s.OutletRepository.FindOutletLocIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindOutletLocIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				otLocId = id
			}

			if id, err := resolveMaster(
				"outlet group",
				"outlet group",
				importData.OtGrpCode,
				importData.OtGrpName,
				func(code string) (int64, error) { return s.OutletRepository.FindOutletGroupIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindOutletGroupIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				otGrpId = id
			}

			if id, err := resolveMaster(
				"price group",
				"price group",
				importData.PriceGrpCode,
				importData.PriceGrpName,
				func(code string) (int64, error) { return s.OutletRepository.FindPriceGroupIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindPriceGroupIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				priceGrpId = id
			}

			if id, err := resolveMaster(
				"district",
				"district",
				importData.DistrictCode,
				importData.DistrictName,
				func(code string) (int64, error) { return s.OutletRepository.FindDistrictIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindDistrictIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				districtId = id
			}

			if id, err := resolveMaster(
				"outlet class",
				"outlet class",
				importData.OtClassCode,
				importData.OtClassName,
				func(code string) (int64, error) { return s.OutletRepository.FindOutletClassIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindOutletClassIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				otClassId = id
			}

			if id, err := resolveMaster(
				"industry",
				"industry",
				importData.IndustryCode,
				importData.IndustryName,
				func(code string) (int64, error) { return s.OutletRepository.FindIndustryIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindIndustryIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				industryId = id
			}

			if id, err := resolveMaster(
				"market",
				"market",
				importData.MarketCode,
				importData.MarketName,
				func(code string) (int64, error) { return s.OutletRepository.FindMarketIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindMarketIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				marketId = id
			}

			if id, err := resolveMaster(
				"outlet type",
				"outlet type",
				importData.OtTypeCode,
				importData.OtTypeName,
				func(code string) (int64, error) { return s.OutletRepository.FindOutletTypeIdByCode(req.CustId, code) },
				func(name string) (int64, error) { return s.OutletRepository.FindOutletTypeIdByName(req.CustId, name) },
			); err != nil {
				logrus.Errorf(err.Error())
				continue
			} else {
				otTypeId = id
			}

			var validationMessages []string

			if len(missingMasters) > 0 {
				msgBody := strings.Join(missingMasters, ", ")
				validationMessages = append(validationMessages, msgBody)
			}

			bankModel, bankMsg, bankErr := s.resolveBankForImport(req, &importData)
			if bankErr != nil {
				logrus.Errorf("gagal memeriksa bank pada baris %d: %v", rowNumber, bankErr)
				continue
			}
			if bankMsg != "" {
				validationMessages = append(validationMessages, bankMsg)
			}

			if len(validationMessages) > 0 {
				prefix := formatOutletValidationPrefix(importData.OutletCode, rowNumber)
				msg := fmt.Sprintf("%s: %s", prefix, strings.Join(validationMessages, "; "))
				if existing := strings.TrimSpace(importData.ErrorMessage); existing != "" {
					msg = existing + "; " + msg
				}
				logrus.Warnf("import outlet reupload: %s", msg)
				importData.ErrorMessage = msg
				if e := s.OutletRepository.CreateOutletTemp(newHistoryId, "failed", req.CustId, importData); e != nil {
					logrus.Errorf("gagal mencatat data outlet_temp pada baris %d: %v", rowNumber, e)
				}
				continue
			}

			logrus.Info(importData.OutletWardId, importData.OutletWard, importData.OutletProvinceId, importData.OutletRegencyId, importData.OutletSubDistrictId)
			err = s.OutletRepository.UpsertProvince(req.CustId, importData.OutletProvinceId, importData.OutletProvince, req.UserId)
			if err != nil {
				fmt.Errorf("gagal memproses outlet_type pada baris %d: %w", rowNumber, err)
			}
			err = s.OutletRepository.UpsertRegency(req.CustId, importData.OutletRegencyId, importData.OutletRegency, importData.OutletProvinceId, req.UserId)
			if err != nil {
				fmt.Errorf("gagal memproses outlet_type pada baris %d: %w", rowNumber, err)
			}
			err = s.OutletRepository.UpsertSubDistrict(req.CustId, importData.OutletSubDistrictId, importData.OutletSubDistrict, importData.OutletProvinceId, importData.OutletRegencyId, req.UserId)
			if err != nil {
				fmt.Errorf("gagal memproses outlet_type pada baris %d: %w", rowNumber, err)
			}
			err = s.OutletRepository.UpsertWard(req.CustId, importData.OutletWardId, importData.OutletWard, importData.OutletProvinceId, importData.OutletRegencyId, importData.OutletSubDistrictId, req.UserId)
			if err != nil {
				fmt.Errorf("gagal memproses outlet_type pada baris %d: %w", rowNumber, err)
			}

			// Map human labels to enum codes before parsing
			importData.OutletStatus = parseEnumCode(importData.OutletStatus, outletStatusCodes)
			if strings.TrimSpace(importData.PaymentType) == "" {
				importData.PaymentType = parseEnumCode(importData.PaymentTypeName, paymentTypeCodes)
			} else {
				importData.PaymentType = parseEnumCode(importData.PaymentType, paymentTypeCodes)
			}
			if strings.TrimSpace(importData.CreditLimitType) == "" {
				importData.CreditLimitType = parseEnumCode(importData.CreditLimitTypeName, creditLimitTypeCodes)
			} else {
				importData.CreditLimitType = parseEnumCode(importData.CreditLimitType, creditLimitTypeCodes)
			}
			if strings.TrimSpace(importData.SalesInvLimitType) == "" {
				importData.SalesInvLimitType = parseEnumCode(importData.SalesInvLimitTypeName, salesInvLimitTypeCodes)
			} else {
				importData.SalesInvLimitType = parseEnumCode(importData.SalesInvLimitType, salesInvLimitTypeCodes)
			}
			if strings.TrimSpace(importData.ObsType) == "" {
				importData.ObsType = parseEnumCode(importData.ObyTypeName, obsTypeCodes)
			} else {
				importData.ObsType = parseEnumCode(importData.ObsType, obsTypeCodes)
			}
			if strings.TrimSpace(importData.CreditLimitAction) == "" {
				importData.CreditLimitAction = parseEnumCode(importData.CreditLimitActionName, limitActionCodes)
			} else {
				importData.CreditLimitAction = parseEnumCode(importData.CreditLimitAction, limitActionCodes)
			}
			if strings.TrimSpace(importData.SalesInvLimitAction) == "" {
				importData.SalesInvLimitAction = parseEnumCode(importData.SalesInvLimitActionName, limitActionCodes)
			} else {
				importData.SalesInvLimitAction = parseEnumCode(importData.SalesInvLimitAction, limitActionCodes)
			}
			if strings.TrimSpace(importData.ObsLimitAction) == "" {
				importData.ObsLimitAction = parseEnumCode(importData.ObsLimitActionName, limitActionCodes)
			} else {
				importData.ObsLimitAction = parseEnumCode(importData.ObsLimitAction, limitActionCodes)
			}
			if strings.TrimSpace(importData.VerificationStatus) == "" {
				importData.VerificationStatus = parseEnumCode(importData.VerificationStatus, verificationStatusCodes)
			} else {
				importData.VerificationStatus = parseEnumCode(importData.VerificationStatus, verificationStatusCodes)
			}
			if strings.TrimSpace(importData.TaxInvoiceForm) == "" {
				importData.TaxInvoiceForm = parseEnumCode(importData.TaxInvoiceFormName, taxInvoiceFormCodes)
			} else {
				importData.TaxInvoiceForm = parseEnumCode(importData.TaxInvoiceForm, taxInvoiceFormCodes)
			}

			// Konversi tipe data
			outletStatus, _ := strconv.ParseInt(importData.OutletStatus, 10, 16)
			top, _ := strconv.ParseInt(importData.Top, 10, 32)
			paymentType, _ := strconv.ParseInt(importData.PaymentType, 10, 16)
			isContraBon, _ := strconv.ParseBool(importData.IsContraBon)
			creditLimitType, _ := strconv.ParseInt(importData.CreditLimitType, 10, 16)
			creditLimit, _ := strconv.ParseFloat(importData.CreditLimit, 64)
			salesInvLimitType, _ := strconv.ParseInt(importData.SalesInvLimitType, 10, 16)
			salesInvLimit, _ := strconv.ParseInt(importData.SalesInvLimit, 10, 16)
			buildingOwnStr := parseEnumCode(importData.BuildingOwn, buildingOwnershipCodes)
			buildingOwn, _ := strconv.ParseInt(buildingOwnStr, 10, 16)
			arStatus, _ := strconv.ParseInt(importData.ArStatus, 10, 16)
			arTotal, _ := strconv.ParseFloat(importData.ArTotal, 64)
			isEmbBail, _ := strconv.ParseBool(importData.IsEmbBail)
			isActive := true
			if val := strings.TrimSpace(importData.IsActive); val != "" {
				if parsed, errBool := strconv.ParseBool(val); errBool == nil {
					isActive = parsed
				}
			}
			isWaNo, _ := strconv.ParseBool(importData.IsWaNo)
			isObs, _ := strconv.ParseBool(importData.IsObs)
			delvIsSameAddr, _ := strconv.ParseBool(importData.DelvIsSameAddr)
			invIsSameAddr, _ := strconv.ParseBool(importData.InvIsSameAddr)
			verificationStatus := resolveImportVerificationStatus(req)
			taxInvoiceForm, _ := strconv.ParseInt(importData.TaxInvoiceForm, 10, 16)
			obsType, _ := strconv.ParseInt(importData.ObsType, 10, 64)
			creditLimitAction, _ := strconv.ParseInt(importData.CreditLimitAction, 10, 64)
			salesInvLimitAction, _ := strconv.ParseInt(importData.SalesInvLimitAction, 10, 64)
			obsLimitAction, _ := strconv.ParseInt(importData.ObsLimitAction, 10, 64)
			obs, _ := strconv.ParseInt(importData.Obs, 10, 64)
			beatId, _ := strconv.ParseInt(importData.BeatId, 10, 64)
			sbeatId, _ := strconv.ParseInt(importData.SbeatId, 10, 64)
			pluGrpId, _ := strconv.ParseInt(importData.PluGrpId, 10, 64)
			convGrpId, _ := strconv.ParseInt(importData.ConvGrpId, 10, 64)
			discInvId, _ := strconv.ParseInt(importData.DiscInvId, 10, 64)
			avgSalesWeek, _ := strconv.ParseFloat(importData.AvgSalesWeek, 64)
			avgSalesMonth, _ := strconv.ParseFloat(importData.AvgSalesMonth, 64)
			firstWeekNo64, _ := strconv.ParseInt(importData.FirstWeekNo, 10, 16)
			firstWeekNo := int16(firstWeekNo64)

			// Siapkan time.Time fields jika tidak kosong
			var firstTransDate, lastTransDate, otStartDate, otRegDate, dob, closedDate, outletEstablishmentDate *time.Time
			if importData.FirstTransDate != "" {
				t, err := time.Parse("2006-01-02", importData.FirstTransDate)
				if err == nil {
					firstTransDate = &t
				}
			}
			if importData.LastTransDate != "" {
				t, err := time.Parse("2006-01-02", importData.LastTransDate)
				if err == nil {
					lastTransDate = &t
				}
			}
			if importData.OtStartDate != "" {
				t, err := time.Parse("2006-01-02", importData.OtStartDate)
				if err == nil {
					otStartDate = &t
				}
			}
			if importData.OtRegDate != "" {
				t, err := time.Parse("2006-01-02", importData.OtRegDate)
				if err == nil {
					otRegDate = &t
				}
			}
			if importData.Dob != "" {
				t, err := time.Parse("2006-01-02", importData.Dob)
				if err == nil {
					dob = &t
				}
			}
			if importData.ClosedDate != "" {
				t, err := time.Parse("2006-01-02", importData.ClosedDate)
				if err == nil {
					closedDate = &t
				}
			}
			if importData.OutletEstablishmentDate != "" {
				if req.IsImportNew {
					if t, err := parseImportDateValue(importData.OutletEstablishmentDate); err == nil {
						outletEstablishmentDate = t
					}
				} else if t, err := time.Parse("2006-01-02", importData.OutletEstablishmentDate); err == nil {
					outletEstablishmentDate = &t
				}
			}
			// if importData.MOutletEstablishmentDate != "" {
			// 	t, err := time.Parse("2006-01-02", importData.MOutletEstablishmentDate)
			// 	if err == nil {
			// 		mOutletEstablishmentDate = &t
			// 	}
			// }

			// Set default timestamps
			now := time.Now()
			createdBy := int64(req.UserId)
			updatedBy := int64(req.UserId)

			reuploadGeoProvinceID := ""
			reuploadGeoRegencyID := ""
			reuploadGeoSubDistrictID := ""
			if req.IsImportNew {
				reuploadGeoProvinceID = importData.OutletProvinceId
				reuploadGeoRegencyID = importData.OutletRegencyId
				reuploadGeoSubDistrictID = importData.OutletSubDistrictId
			}

			// Siapkan data untuk disimpan
			processedData := entity.ProcessedOutlet{
				CustId:                  req.CustId,
				OutletCode:              importData.OutletCode,
				Barcode:                 importData.Barcode,
				OutletName:              importData.OutletName,
				OutletStatus:            int16(outletStatus),
				Address1:                importData.Address1,
				Address2:                importData.Address2,
				City:                    importData.City,
				ZipCode:                 importData.ZipCode,
				PhoneNo:                 importData.PhoneNo,
				WaNo:                    importData.WaNo,
				FaxNo:                   importData.FaxNo,
				Email:                   importData.Email,
				DiscGrpId:               discGrpId,
				OtLocId:                 otLocId,
				OtGrpId:                 otGrpId,
				PriceGrpId:              priceGrpId,
				DistrictId:              districtId,
				BeatId:                  beatId,
				SbeatId:                 sbeatId,
				OtClassId:               otClassId,
				IndustryId:              industryId,
				MarketId:                marketId,
				Top:                     int32(top),
				PaymentType:             int16(paymentType),
				IsContraBon:             isContraBon,
				PluGrpId:                pluGrpId,
				ConvGrpId:               convGrpId,
				DiscInvId:               discInvId,
				AgentFrom:               importData.AgentFrom,
				CreditLimitType:         int16(creditLimitType),
				CreditLimitTypeName:     importData.CreditLimitTypeName,
				CreditLimit:             creditLimit,
				SalesInvLimitType:       int16(salesInvLimitType),
				SalesInvLimitTypeName:   importData.SalesInvLimitTypeName,
				SalesInvLimit:           int16(salesInvLimit),
				AvgSalesWeek:            avgSalesWeek,
				AvgSalesMonth:           avgSalesMonth,
				FirstTransDate:          firstTransDate,
				LastTransDate:           lastTransDate,
				FirstWeekNo:             firstWeekNo,
				OtStartDate:             otStartDate,
				OtRegDate:               otRegDate,
				BuildingOwn:             int16(buildingOwn),
				Dob:                     dob,
				ArStatus:                int16(arStatus),
				ArTotal:                 arTotal,
				ClosedDate:              closedDate,
				IsEmbBail:               isEmbBail,
				IsPkpOutlet:             isEmbBail,
				TaxName:                 importData.TaxName,
				TaxAddr1:                importData.TaxAddr1,
				TaxAddr2:                importData.TaxAddr2,
				TaxCity:                 importData.TaxCity,
				TaxNo:                   importData.TaxNo,
				OwnerName:               importData.OwnerName,
				OwnerAddr1:              importData.OwnerAddr1,
				OwnerAddr2:              importData.OwnerAddr2,
				OwnerCity:               importData.OwnerCity,
				OwnerPhoneNo:            importData.OwnerPhoneNo,
				OwnerIdNo:               importData.OwnerIdNo,
				DelvAddr1:               importData.DelvAddr1,
				DelvAddr2:               importData.DelvAddr2,
				DelvCity:                importData.DelvCity,
				InvAddr1:                importData.InvAddr1,
				InvAddr2:                importData.InvAddr2,
				InvCity:                 importData.InvCity,
				IsActive:                isActive,
				CreatedBy:               &createdBy,
				CreatedAt:               &now,
				UpdatedBy:               &updatedBy,
				UpdatedAt:               &now,
				IsDel:                   false,
				Latitude:                importData.Latitude,
				Longitude:               importData.Longitude,
				ImageUrl:                importData.ImageUrl,
				OtTypeId:                otTypeId,
				IsObs:                   isObs,
				Obs:                     obs,
				OutletProvinceId:        reuploadGeoProvinceID,
				OutletRegencyId:         reuploadGeoRegencyID,
				OutletSubDistrictId:     reuploadGeoSubDistrictID,
				OutletWardId:            importData.OutletWardId,
				IsWaNo:                  isWaNo,
				DelvWardId:              importData.DelvWardId,
				DelvZipCode:             importData.DelvZipCode,
				DelvIsSameAddr:          delvIsSameAddr,
				InvWardId:               importData.InvWardId,
				InvZipCode:              importData.InvZipCode,
				InvIsSameAddr:           invIsSameAddr,
				VerificationStatus:      int16(verificationStatus),
				TaxInvoiceForm:          int16(taxInvoiceForm),
				ObsType:                 obsType,
				CreditLimitAction:       creditLimitAction,
				CreditLimitActionName:   importData.CreditLimitActionName,
				SalesInvLimitAction:     salesInvLimitAction,
				SalesInvLimitActionName: importData.SalesInvLimitActionName,
				ObsLimitAction:          obsLimitAction,
				OutletEstablishmentDate: outletEstablishmentDate,
				DelvCity2:               importData.DelvCity2,
				DelvLatitude:            importData.DelvLatitude,
				DelvLongitude:           importData.DelvLongitude,
				DelvLatitude2:           importData.DelvLatitude2,
				DelvLongitude2:          importData.DelvLongitude2,
				DelvWardId2:             importData.DelvWardId2,
				DelvZipCode2:            importData.DelvZipCode2,
			}

			outletID, err := s.OutletRepository.CreateOutlet(processedData)
			if err != nil {
				logrus.Errorf("gagal menyimpan outlet (kode %s): %v", importData.OutletCode, err)
				continue
			}
			if err := s.insertImportOutletDetails(req, outletID, importData, int64(taxInvoiceForm), isEmbBail, bankModel); err != nil {
				logrus.Errorf("gagal menyimpan detail outlet (kode %s): %v", importData.OutletCode, err)
				continue
			}
			success++

			// delete temp row on success basis (best effort here)
			// _ = s.OutletRepository.DeleteOutletUpdateTempByCtid(ctid)
		}

		newFailed := total - success
		if success < 0 {
			success = 0
		}
		err = s.OutletRepository.UpdateImportHistoryOutlet(newHistoryId, success, newFailed, newFailed > 0)
		if err != nil {
			logrus.Info("Error update history")
		}
	}()
	return nil
}

func minOutletContactID(rows []model.MOutletContactRead) int64 {
	var min int64
	for _, c := range rows {
		if c.OutletContactId == nil {
			continue
		}
		id := *c.OutletContactId
		if min == 0 || id < min {
			min = id
		}
	}
	return min
}

func minOutletTaxID(rows []model.MOutletTaxRead) int64 {
	var min int64
	for _, t := range rows {
		if t.OutletTaxId == nil {
			continue
		}
		id := *t.OutletTaxId
		if min == 0 || id < min {
			min = id
		}
	}
	return min
}

func (s *outletServiceImpl) processUpdateRows(req entity.ImportRequest, rows [][]string, historyId int64) error {
	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	headers := rows[0]
	headerMap := map[string]int{}
	for i, h := range headers {
		headerMap[canonicalizeHeader(h)] = i
	}
	logrus.Info("Header Map:", headerMap)

	get := func(row []string, key string) string {
		if idx, ok := headerMap[canonicalizeHeader(key)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// getSyn tries multiple header keys (canonical and common synonyms) for a logical field
	getSyn := func(row []string, key string, fallbacks ...string) string {
		// try primary
		if v := get(row, key); v != "" {
			return v
		}
		// try fallbacks
		for _, fb := range fallbacks {
			if v := get(row, fb); v != "" {
				return v
			}
		}
		return ""
	}

	parseInt64 := func(val string) int64 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
		if err != nil {
			return 0
		}
		return v
	}

	parseInt32 := func(val string) int32 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, err := strconv.ParseInt(strings.TrimSpace(val), 10, 32)
		if err != nil {
			return 0
		}
		return int32(v)
	}

	parseInt32WithCodes := func(val string, codes map[string]int) int32 {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return 0
		}
		if v, err := strconv.ParseInt(trimmed, 10, 32); err == nil {
			return int32(v)
		}
		if codes != nil {
			if code, ok := codes[normEnum(trimmed)]; ok {
				return int32(code)
			}
		}
		return 0
	}

	parseInt64WithCodes := func(val string, codes map[string]int) int64 {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return 0
		}
		if v, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
			return v
		}
		if codes != nil {
			if code, ok := codes[normEnum(trimmed)]; ok {
				return int64(code)
			}
		}
		return 0
	}

	parseFloat64 := func(val string) float64 {
		if strings.TrimSpace(val) == "" {
			return 0
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return 0
		}
		return v
	}

	parseBoolWithDefault := func(val string, defaultValue bool) bool {
		trimmed := strings.TrimSpace(strings.ToLower(val))
		if trimmed == "" {
			return defaultValue
		}
		switch trimmed {
		case "1", "true", "t", "ya", "yes", "y":
			return true
		case "0", "false", "f", "tidak", "no", "n":
			return false
		default:
			return defaultValue
		}
	}
	parseActiveFlag := func(val string) *bool {
		trimmed := strings.ToLower(strings.TrimSpace(val))
		if trimmed == "" {
			return nil
		}
		switch trimmed {
		case "1", "true", "t", "ya", "yes", "y", "active", "aktif":
			v := true
			return &v
		case "0", "false", "f", "tidak", "no", "n", "inactive", "deactive", "nonactive", "non-active", "non active", "deactivated", "nonaktif":
			v := false
			return &v
		default:
			return nil
		}
	}

	// parseDate attempts to parse common date/datetime layouts; returns (zeroTime,false) when parsing fails
	parseDate := func(val string) (time.Time, bool) {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return time.Time{}, false
		}
		layouts := []string{
			"2006-01-02",
			time.RFC3339,
			"2006-01-02 15:04:05",
			"02/01/2006",
			"02-01-2006",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, trimmed); err == nil {
				return t, true
			}
		}
		return time.Time{}, false
	}

	total := len(rows) - 1
	success := 0
	start := time.Now()

	var (
		// provinceId     string
		// regencyId      string
		// subDistrictId  string
		wardId       string
		otGrpId      int64
		otLocId      int64
		industryId   int64
		marketId     int64
		outletTypeId int64
		// bankId    int64
		priceGrpId int64
		otClassId  int64
		discGrpId  int64
		districtId int64
	)

	for i, row := range rows[1:] {
		rowErrors := []string{}
		missingMasters := make([]string, 0, 16)
		missingMasterAdded := make(map[string]struct{})
		addMissingMaster := func(msg string) {
			msg = strings.TrimSpace(msg)
			if msg == "" {
				return
			}
			if _, exists := missingMasterAdded[msg]; exists {
				return
			}
			missingMasterAdded[msg] = struct{}{}
			missingMasters = append(missingMasters, msg)
		}
		addRowError := func(format string, args ...interface{}) {
			rowErrors = append(rowErrors, fmt.Sprintf(format, args...))
		}

		// Helper: resolve outlet_id by provided identifiers: outlet_id, outlet_bank_id, outlet_contact_id, outlet_tax_id, then fallback to outlet_code
		resolveOutletId := func() (int64, error) {
			// 1) outlet_id direct
			if v := strings.TrimSpace(get(row, "outlet_id")); v != "" {
				if oid, perr := strconv.ParseInt(v, 10, 64); perr == nil {
					out, err := s.OutletRepository.FindOneByOutletIdAndCustId(oid, req.CustId, req.ParentCustId)
					if err == nil && out.OutletId > 0 {
						return int64(out.OutletId), nil
					}
					return 0, fmt.Errorf("outlet_id %d tidak ditemukan", oid)
				}
			}
			// 2) by outlet_contact_id
			if v := strings.TrimSpace(get(row, "outlet_contact_id")); v != "" {
				if id, perr := strconv.ParseInt(v, 10, 64); perr == nil {
					if oid, err := s.OutletRepository.FindOutletIdByOutletContactId(req.CustId, id); err == nil && oid > 0 {
						return oid, nil
					}
					return 0, fmt.Errorf("outlet_contact_id %d tidak ditemukan", id)
				}
			}
			// 3) by outlet_bank_id
			if v := strings.TrimSpace(get(row, "outlet_bank_id")); v != "" {
				if id, perr := strconv.ParseInt(v, 10, 64); perr == nil {
					if oid, err := s.OutletRepository.FindOutletIdByOutletBankId(req.CustId, id); err == nil && oid > 0 {
						return oid, nil
					}
					return 0, fmt.Errorf("outlet_bank_id %d tidak ditemukan", id)
				}
			}
			// 4) by outlet_tax_id
			if v := strings.TrimSpace(get(row, "outlet_tax_id")); v != "" {
				if id, perr := strconv.ParseInt(v, 10, 64); perr == nil {
					if oid, err := s.OutletRepository.FindOutletIdByOutletTaxId(req.CustId, id); err == nil && oid > 0 {
						return oid, nil
					}
					return 0, fmt.Errorf("outlet_tax_id %d tidak ditemukan", id)
				}
			}
			// 5) fallback to outlet_code
			if code := strings.TrimSpace(get(row, "outlet_code")); code != "" {
				out, err := s.OutletRepository.FindOneByOutletCodeAndCustId(code, req.CustId, req.ParentCustId)
				if err == nil && out.OutletId > 0 {
					return int64(out.OutletId), nil
				}
				return 0, fmt.Errorf("outlet_code %s tidak ditemukan", code)
			}
			return 0, errors.New("wajib isi salah satu: outlet_id, outlet_contact_id, outlet_bank_id, outlet_tax_id, atau outlet_code")
		}

		// ========= OUTLET BANK DETAIL (Prefer match by outlet_bank_id; else by bank_id/bank_code + account_no) =========
		if get(row, "account_no") != "" || get(row, "account_name") != "" || get(row, "bank_id") != "" || get(row, "bank_code") != "" || get(row, "outlet_bank_id") != "" {
			outletId, err := resolveOutletId()
			if err != nil {
				addRowError("%v", err)
			} else {
				// derive bankId from bank_id or bank_code
				var bankId int64
				if v := strings.TrimSpace(get(row, "bank_id")); v != "" {
					bankId, _ = strconv.ParseInt(v, 10, 64)
				}
				if bankId == 0 {
					if code := strings.TrimSpace(get(row, "bank_name")); code != "" {
						if bid, e := s.OutletRepository.FindBankIdByName(req.CustId, code); e == nil {
							bankId = bid
						}
					}
				}
				accNo := strings.TrimSpace(get(row, "account_no"))
				accName := strings.TrimSpace(get(row, "account_name"))

				// fetch existing bank details for outlet
				existing, e := s.OutletRepository.GetDetailOutletbank(outletId, req.CustId)
				if e != nil {
					addRowError("gagal mengambil data outlet_bank: %v", e)
				} else {
					// Prefer match by outlet_bank_id if provided
					var targetId int64
					if idStr := strings.TrimSpace(get(row, "outlet_bank_id")); idStr != "" {
						if id, pe := strconv.ParseInt(idStr, 10, 64); pe == nil {
							for _, b := range existing {
								if b.OutletBankId != nil && *b.OutletBankId == id {
									targetId = id
									break
								}
							}
							if targetId == 0 {
								addRowError("outlet_bank_id %d tidak ditemukan untuk outlet", id)
							}
						}
					}
					// Fallback match by bank_id+account_no
					if targetId == 0 && bankId > 0 && accNo != "" {
						for _, b := range existing {
							if b.BankId == bankId {
								if b.AccountNo != nil && strings.EqualFold(strings.TrimSpace(*b.AccountNo), accNo) {
									if b.OutletBankId != nil {
										targetId = *b.OutletBankId
									}
									break
								}
							}
						}
					}

					trx, e := s.OutletRepository.TrxBegin()
					if e != nil {
						addRowError("gagal mulai transaksi bank: %v", e)
					} else {
						if targetId > 0 {
							// update requires at least new values; if bankId still 0, keep existing via repository logic
							if bankId == 0 && accNo == "" && accName == "" {
								_ = trx.TrxRollback()
								addRowError("tidak ada data bank baru untuk diupdate")
							} else {
								reqBank := entity.OutletBank{BankId: nil, AccountNo: accNo, AccountName: accName}
								if bankId > 0 {
									reqBank.BankId = &bankId
								}
								if e := trx.UpdateDetailOutletBank(int(outletId), targetId, reqBank); e != nil {
									_ = trx.TrxRollback()
									addRowError("pembaruan outlet_bank gagal: %v", e)
								} else {
									_ = trx.TrxCommit()
								}
							}
						} else {
							// insert requires bankId, account_no, account_name
							if bankId == 0 || accNo == "" || accName == "" {
								_ = trx.TrxRollback()
								addRowError("data bank tidak lengkap (bank_id/bank_code, account_no, account_name)")
							} else {
								m := model.MOutletBank{CustID: req.CustId, OutletID: outletId, BankId: bankId, AccountNo: accNo, AccountName: accName}
								if e := trx.StoreDetailBank(&m); e != nil {
									_ = trx.TrxRollback()
									addRowError("penyimpanan outlet_bank gagal: %v", e)
								} else {
									_ = trx.TrxCommit()
								}
							}
						}
					}
				}
			}
		}

		if get(row, "contact_name") != "" || get(row, "job_title") != "" || get(row, "contact_phone_no") != "" || get(row, "contact_email") != "" || get(row, "outlet_contact_id") != "" || get(row, "contact_wa_no") != "" || get(row, "identity_type") != "" || get(row, "identity_no") != "" || get(row, "contact_is_wa_no") != "" {
			outletId, err := resolveOutletId()
			if err != nil {
				addRowError("%v", err)
			} else {
				name := strings.TrimSpace(getSyn(row, "contact_name", "name", "contact"))
				existing, e := s.OutletRepository.GetDetailOutletContact(outletId, req.CustId)
				if e != nil {
					addRowError("gagal mengambil data outlet_contact: %v", e)
				} else {
					var targetId int64
					if idStr := strings.TrimSpace(get(row, "outlet_contact_id")); idStr != "" {
						if id, pe := strconv.ParseInt(idStr, 10, 64); pe == nil {
							for _, c := range existing {
								if c.OutletContactId != nil && *c.OutletContactId == id {
									targetId = id
									break
								}
							}
							if targetId == 0 {
								addRowError("outlet_contact_id %d tidak ditemukan untuk outlet", id)
							}
						}
					}
					if targetId == 0 && name != "" {
						for _, c := range existing {
							if c.ContactName != nil && strings.EqualFold(strings.TrimSpace(*c.ContactName), name) {
								if c.OutletContactId != nil {
									targetId = *c.OutletContactId
								}
								break
							}
						}
					}
					if targetId == 0 {
						targetId = minOutletContactID(existing)
					}

					phone := strings.TrimSpace(getSyn(row, "contact_phone_no", "phone_no", "phone", "contact_phone", "contact_phone_number"))
					wa := strings.TrimSpace(getSyn(row, "contact_wa_no", "whatsapp_no", "wa_no", "whatsapp_number", "whatsapp", "wa"))
					email := strings.TrimSpace(getSyn(row, "contact_email", "email"))
					idNo := strings.TrimSpace(get(row, "identity_no"))
					isWaStr := strings.ToLower(strings.TrimSpace(getSyn(row, "contact_is_wa_no", "is_wa_no", "contact_wa_active", "contact_wa", "wa_active")))
					idType := strings.TrimSpace(get(row, "identity_type"))
					fax := strings.TrimSpace(getSyn(row, "fax_number", "fax"))
					job := strings.TrimSpace(getSyn(row, "job_title", "position"))

					trx, e := s.OutletRepository.TrxBegin()
					if e != nil {
						addRowError("gagal mulai transaksi contact: %v", e)
					} else {
						if targetId > 0 {
							var patch model.MOutletContactUpdate
							if name != "" {
								patch.ContactName = &name
							}
							if job != "" {
								patch.JobTitle = &job
							}
							if phone != "" {
								patch.PhoneNo = &phone
							}
							if wa != "" {
								patch.WaNo = &wa
							}
							if email != "" {
								patch.Email = &email
							}
							if idNo != "" {
								patch.IdentityNo = &idNo
							}
							if isWaStr != "" {
								isWa := isWaStr == "1" || isWaStr == "true" || isWaStr == "ya" || isWaStr == "yes" || isWaStr == "y"
								patch.IsWaNo = &isWa
							}
							if idType != "" {
								patch.IdentityType = &idType
							}
							if fax != "" {
								patch.FaxNumber = &fax
							}
							if e := trx.UpdateDetailOutletContactPartial(int(outletId), targetId, patch); e != nil {
								_ = trx.TrxRollback()
								addRowError("pembaruan outlet_contact gagal: %v", e)
							} else {
								_ = trx.TrxCommit()
							}
						} else {
							if name == "" {
								_ = trx.TrxRollback()
								addRowError("contact_name wajib diisi untuk menambah contact")
							} else {
								isWa := isWaStr == "1" || isWaStr == "true" || isWaStr == "ya" || isWaStr == "yes" || isWaStr == "y"
								m := model.MOutletContact{
									CustID:       req.CustId,
									OutletID:     outletId,
									ContactName:  name,
									JobTitle:     job,
									PhoneNo:      phone,
									WaNo:         wa,
									Email:        email,
									IdentityNo:   idNo,
									IsWaNo:       isWa,
									IdentityType: idType,
									FaxNumber:    fax,
								}
								if e := trx.StoreDetailContact(&m); e != nil {
									_ = trx.TrxRollback()
									addRowError("penyimpanan outlet_contact gagal: %v", e)
								} else {
									_ = trx.TrxCommit()
								}
							}
						}
					}
				}
			}
		}

		if get(row, "outlet_tax_id") != "" || get(row, "tax_invoice_id") != "" || get(row, "tax_identifier_no") != "" || get(row, "address_tax") != "" || get(row, "tax_type") != "" || get(row, "nitku") != "" || get(row, "tax_identifier_type") != "" || get(row, "tax_name") != "" {
			outletId, err := resolveOutletId()
			if err != nil {
				addRowError("%v", err)
			} else {
				idStr := strings.TrimSpace(get(row, "outlet_tax_id"))
				taxInvStr := strings.TrimSpace(get(row, "tax_invoice_id"))
				var taxInvId int64
				if taxInvStr != "" {
					taxInvId, _ = strconv.ParseInt(taxInvStr, 10, 64)
				}
				taxType := strings.TrimSpace(get(row, "tax_type"))
				nitku := strings.TrimSpace(get(row, "nitku"))
				addrTax := strings.TrimSpace(get(row, "address_tax"))
				idType := strings.TrimSpace(get(row, "tax_identifier_type"))
				idNo := strings.TrimSpace(get(row, "tax_identifier_no"))

				existing, e := s.OutletRepository.GetDetailOutletTax(outletId, req.CustId)
				if e != nil {
					addRowError("gagal mengambil data outlet_tax: %v", e)
				} else {
					var targetId int64
					if idStr != "" {
						if id, pe := strconv.ParseInt(idStr, 10, 64); pe == nil {
							for _, t := range existing {
								if t.OutletTaxId != nil && *t.OutletTaxId == id {
									targetId = id
									break
								}
							}
							if targetId == 0 {
								addRowError("outlet_tax_id %d tidak ditemukan untuk outlet", id)
							}
						}
					}
					if targetId == 0 && taxInvStr != "" {
						for _, t := range existing {
							if t.TaxInvoiceId != nil && *t.TaxInvoiceId == taxInvId {
								if t.OutletTaxId != nil {
									targetId = *t.OutletTaxId
								}
								break
							}
						}
					}
					if targetId == 0 {
						targetId = minOutletTaxID(existing)
					}

					trx, e := s.OutletRepository.TrxBegin()
					if e != nil {
						addRowError("gagal mulai transaksi tax: %v", e)
					} else {
						if targetId > 0 {
							var patch model.MOutletTaxUpdate
							if taxInvStr != "" {
								v := taxInvId
								patch.TaxInvoiceId = &v
							}
							if taxType != "" {
								patch.TaxType = &taxType
							}
							if nitku != "" {
								patch.Nitku = &nitku
							}
							if addrTax != "" {
								patch.AdddressTax = &addrTax
							}
							if idType != "" {
								patch.TaxIdentifierType = &idType
							}
							if idNo != "" {
								patch.TaxIdentifierNo = &idNo
							}
							if emb := strings.TrimSpace(get(row, "is_emb_bail")); emb != "" {
								v := parseBoolWithDefault(emb, false)
								patch.IsEmbBail = &v
							}
							if v := strings.TrimSpace(get(row, "tax_no")); v != "" {
								patch.TaxNo = &v
							}
							if v := strings.TrimSpace(get(row, "tax_name")); v != "" {
								patch.TaxName = &v
							}
							if v := strings.TrimSpace(get(row, "tax_addr1")); v != "" {
								patch.TaxAddr1 = &v
							}
							if v := strings.TrimSpace(get(row, "tax_addr2")); v != "" {
								patch.TaxAddr2 = &v
							}
							if v := strings.TrimSpace(get(row, "tax_city")); v != "" {
								patch.TaxCity = &v
							}
							if e := trx.UpdateDetailOutletTaxPartial(int(outletId), targetId, patch); e != nil {
								_ = trx.TrxRollback()
								addRowError("pembaruan outlet_tax gagal: %v", e)
							} else {
								_ = trx.TrxCommit()
							}
						} else {
							m := model.MOutletTax{
								CustID:            req.CustId,
								OutletID:          outletId,
								TaxInvoiceId:      taxInvId,
								TaxType:           taxType,
								Nitku:             nitku,
								AdddressTax:       addrTax,
								TaxIdentifierType: idType,
								TaxIdentifierNo:   idNo,
							}
							if e := trx.StoreDetailTax(&m); e != nil {
								_ = trx.TrxRollback()
								addRowError("penyimpanan outlet_tax gagal: %v", e)
							} else {
								_ = trx.TrxCommit()
							}
						}
					}
				}
			}
		}

		// BANK
		{
			raw := get(row, "bank_id")
			codeKey := get(row, "bank_name")
			codeNew := get(row, "bank_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindBankIdByName(req.CustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master bank tidak ditemukan")
				} else {
					id = v
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckBankExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan bank gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: bank_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckBankCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi bank_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: bank_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportBank(req.CustId, codeNew, get(row, "bank_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan data bank gagal: %v", i+2, err))
			// 	}
			// }
		}

		// OUTLET CLASS
		{
			raw := get(row, "ot_class_id")
			codeKey := get(row, "ot_class_name")
			codeNew := get(row, "ot_class_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindOutletClassIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master outlet class tidak ditemukan")
				} else {
					id = v
					otClassId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckOutletClassExists(req.ParentCustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan outlet_class gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: ot_class_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckOutletClassCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi outlet_class_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: ot_class_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportOutletClass(req.CustId, codeNew, get(row, "ot_class_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan outlet_class gagal: %v", i+2, err))
			// 	}
			// }
		}

		// OUTLET GROUP
		{
			raw := get(row, "outlet_grp_id")
			codeKey := get(row, "ot_grp_name")
			codeNew := get(row, "ot_grp_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if codeKey != "" {
				v, err := s.OutletRepository.FindOutletGroupIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master outlet group tidak ditemukan")
				} else {
					id = v
					otGrpId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckOutletGroupExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan outlet_group gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_grp_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckOutletGroupCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi outlet_group_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_grp_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportOutletGroup(req.CustId, codeNew, get(row, "outlet_grp_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan outlet_group gagal: %v", i+2, err))
			// 	}
			// }
		}

		// OUTLET LOCATION
		{
			raw := get(row, "outlet_loc_id")
			codeKey := get(row, "ot_loc_name")
			codeNew := get(row, "ot_loc_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindOutletLocIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master outlet location tidak ditemukan")
				} else {
					id = v
					otLocId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckOutletLocExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan outlet_location gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_loc_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckOutletLocCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi outlet_loc_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_loc_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportOutletLoc(req.CustId, codeNew, get(row, "outlet_loc_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan outlet_location gagal: %v", i+2, err))
			// 	}
			// }
		}

		// OUTLET TYPE
		{
			raw := get(row, "outlet_type_id")
			codeKey := get(row, "ot_type_name")
			codeNew := get(row, "outlet_type_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindOutletTypeIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master outlet type tidak ditemukan")
				} else {
					id = v
					outletTypeId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckOutletTypeExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan outlet_type gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_type_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckOutletTypeCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi outlet_type_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: outlet_type_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportOutletType(req.CustId, codeNew, get(row, "outlet_type_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan outlet_type gagal: %v", i+2, err))
			// 	}
			// }
		}

		// DISTRICT
		{
			raw := get(row, "district_id")
			codeKey := get(row, "district_name")
			codeNew := get(row, "district_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindDistrictIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master district tidak ditemukan")
				} else {
					id = v
					districtId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckDistrictExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan district gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: district_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckDistrictCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi district_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: district_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportDistrict(req.CustId, codeNew, get(row, "district_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan district gagal: %v", i+2, err))
			// 	}
			// }
		}

		// DISC GROUP
		{
			raw := get(row, "disc_grp_id")
			codeKey := get(row, "disc_grp_name")
			codeNew := get(row, "disc_grp_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindDiscGroupIdByName(req.CustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master disc group tidak ditemukan")
				} else {
					id = v
					discGrpId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckDiscGroupExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan disc_group gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: disc_grp_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckDiscGroupCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi disc_grp_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: disc_grp_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportDiscGroup(req.CustId, codeNew, get(row, "disc_grp_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan disc_group gagal: %v", i+2, err))
			// 	}
			// }
		}

		// MARKET
		{
			raw := get(row, "market_id")
			codeKey := get(row, "market_name")
			codeNew := get(row, "market_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindMarketIdByName(req.CustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master market tidak ditemukan")
				} else {
					id = v
					marketId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckMarketExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan market gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: market_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckMarketCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi market_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: market_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportMarket(req.CustId, codeNew, get(row, "market_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan market gagal: %v", i+2, err))
			// 	}
			// }
		}

		// INDUSTRY
		{
			raw := get(row, "industry_id")
			codeKey := get(row, "industry_name")
			codeNew := get(row, "industry_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindIndustryIdByName(req.ParentCustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master industry tidak ditemukan")
				} else {
					id = v
					industryId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckIndustryExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan industry gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: industry_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckIndustryCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi industry_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: industry_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if err := s.OutletRepository.UpdateImportIndustry(req.CustId, codeNew, get(row, "industry_name"), id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan industry gagal: %v", i+2, err))
			// 	}
			// }
		}

		// PRICE GROUP
		{
			raw := get(row, "price_grp_id")
			codeKey := get(row, "price_grp_name")
			codeNew := get(row, "price_grp_code_new")
			if codeNew == "" {
				codeNew = codeKey
			}
			var id int64
			if raw != "" {
				id, _ = strconv.ParseInt(raw, 10, 64)
			}
			if id == 0 && codeKey != "" {
				v, err := s.OutletRepository.FindPriceGroupIdByName(req.CustId, codeKey)
				if err != nil || v == 0 {
					addMissingMaster("Master price group tidak ditemukan")
				} else {
					id = v
					priceGrpId = id
				}
			}
			// if id > 0 {
			// 	if exists, err := s.OutletRepository.CheckPriceGroupExists(req.CustId, id); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pengecekan price_group gagal: %v", i+2, err))
			// 	} else if !exists {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: price_grp_id %d tidak ditemukan", i+2, id))
			// 	}
			// 	if codeNew != "" {
			// 		if dup, err := s.OutletRepository.CheckPriceGroupCodeDuplicate(req.CustId, id, codeNew); err != nil {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: validasi price_grp_code gagal: %v", i+2, err))
			// 		} else if dup {
			// 			rowErrors = append(rowErrors, fmt.Sprintf("baris %d: price_grp_code %s sudah digunakan", i+2, codeNew))
			// 		}
			// 	}
			// 	if _, err := s.OutletRepository.UpdateImportPriceGroup(req.CustId, id, codeNew, get(row, "price_grp_name")); err != nil {
			// 		rowErrors = append(rowErrors, fmt.Sprintf("baris %d: pembaruan price_group gagal: %v", i+2, err))
			// 	}
			// }
		}

		// PROVINCE
		{
			provinceId := get(row, "outlet_province_id")
			province := get(row, "outlet_province")
			if provinceId != "" || province != "" {
				provinceId, _ = s.OutletRepository.GetProvinceIdByName(req.ParentCustId, province)
				if strings.TrimSpace(provinceId) == "" {
					addMissingMaster("Master outlet province tidak ditemukan")
				} else if exists, err := s.OutletRepository.CheckProvinceExists(req.ParentCustId, provinceId); err != nil {
					addRowError("pengecekan master outlet province gagal: %v", err)
				} else if !exists {
					addMissingMaster("Master outlet province tidak ditemukan")
				}
			}
		}

		// REGENCY
		{
			regencyId := get(row, "outlet_regency_id")
			regency := get(row, "outlet_regency")
			if regencyId != "" || regency != "" {
				regencyId, _ = s.OutletRepository.GetRegencyIdByName(req.ParentCustId, regency)
				if strings.TrimSpace(regencyId) == "" {
					addMissingMaster("Master outlet regency tidak ditemukan")
				} else if exists, err := s.OutletRepository.CheckRegencyExists(req.ParentCustId, regencyId); err != nil {
					addRowError("pengecekan master outlet regency gagal: %v", err)
				} else if !exists {
					addMissingMaster("Master outlet regency tidak ditemukan")
				}
			}
		}

		// SUB-DISTRICT
		{
			subDistrictId := get(row, "outlet_sub_district_id")
			subDistrict := get(row, "outlet_sub_district")
			if subDistrictId != "" || subDistrict != "" {
				subDistrictId, _ = s.OutletRepository.GetSubDistrictIdByName(req.ParentCustId, subDistrict)
				if strings.TrimSpace(subDistrictId) == "" {
					addMissingMaster("Master outlet sub_district tidak ditemukan")
				} else if exists, err := s.OutletRepository.CheckSubDistrictExists(req.ParentCustId, subDistrictId); err != nil {
					addRowError("pengecekan master outlet sub_district gagal: %v", err)
				} else if !exists {
					addMissingMaster("Master outlet sub_district tidak ditemukan")
				}
			}
		}

		// WARD
		{
			ward := get(row, "outlet_ward")
			if ward != "" {
				var err error
				wardId, err = s.OutletRepository.GetWardIdByName(req.ParentCustId, ward)
				if err != nil {
					addRowError("pengecekan master outlet ward gagal: %v", err)
					wardId = ""
				} else if strings.TrimSpace(wardId) == "" {
					addMissingMaster("Master outlet ward tidak ditemukan")
				} else if exists, err := s.OutletRepository.CheckWardExists(req.ParentCustId, wardId); err != nil {
					addRowError("pengecekan master outlet ward gagal: %v", err)
				} else if !exists {
					addMissingMaster("Master outlet ward tidak ditemukan")
				}
			} else {
				wardId = ""
			}
		}

		// update outlet diakhir saja, sehingga kalau error outlet tidak terupdate
		{
			outletIdRaw := get(row, "outlet_id")
			outletCode := get(row, "outlet_code")

			if outletIdRaw != "" || outletCode != "" {
				resolvedOutletId, err := resolveOutletId()
				if err != nil {
					addRowError("%v", err)
				} else {
					activeFlag := parseActiveFlag(get(row, "is_active"))
					if activeFlag == nil {
						activeFlag = parseActiveFlag(get(row, "outlet_status"))
					}
					updateReq := entity.ProcessedUpdateOutlet{
						CustId:              req.CustId,
						OutletId:            resolvedOutletId,
						OutletCode:          outletCode,
						Barcode:             get(row, "barcode"),
						OutletName:          get(row, "outlet_name"),
						OutletStatus:        parseInt32WithCodes(get(row, "outlet_status"), outletStatusCodes),
						Address1:            get(row, "address1"),
						Address2:            get(row, "address2"),
						City:                get(row, "city"),
						ZipCode:             get(row, "zip_code"),
						PhoneNo:             get(row, "phone_no"),
						WaNo:                get(row, "wa_no"),
						FaxNo:               get(row, "fax_no"),
						Email:               get(row, "email"),
						DiscGrpId:           discGrpId,
						OtLocId:             otLocId,
						OtGrpId:             otGrpId,
						PriceGrpId:          priceGrpId,
						DistrictId:          districtId,
						BeatId:              parseInt64(get(row, "beat_id")),
						SbeatId:             parseInt64(get(row, "sbeat_id")),
						OtClassId:           otClassId,
						IndustryId:          industryId,
						MarketId:            marketId,
						OtTypeId:            outletTypeId,
						OutletWardId:        wardId,
						Top:                 parseInt32(get(row, "top")),
						PaymentType:         parseInt32WithCodes(get(row, "payment_type"), paymentTypeCodes),
						IsContraBon:         parseBoolWithDefault(get(row, "is_contra_bon"), true),
						CreditLimitType:     parseInt32WithCodes(get(row, "credit_limit_type"), creditLimitTypeCodes),
						CreditLimit:         parseFloat64(get(row, "credit_limit")),
						SalesInvLimitType:   parseInt32WithCodes(get(row, "sales_inv_limit_type"), salesInvLimitTypeCodes),
						SalesInvLimit:       parseInt32(get(row, "sales_inv_limit")),
						AvgSalesWeek:        parseFloat64(get(row, "avg_sales_week")),
						AvgSalesMonth:       parseFloat64(get(row, "avg_sales_month")),
						CreditLimitAction:   parseInt64WithCodes(get(row, "credit_limit_action"), limitActionCodes),
						SalesInvLimitAction: parseInt64WithCodes(get(row, "sales_inv_limit_action"), limitActionCodes),
						PluGrpId:            parseInt64(get(row, "plu_grp_id")),
						ConvGrpId:           parseInt64(get(row, "conv_grp_id")),
						DiscInvId:           parseInt64(get(row, "disc_inv_id")),
						AgentFrom:           get(row, "agent_from"),
						FirstWeekNo:         parseInt32(get(row, "first_week_no")),
						BuildingOwn:         parseInt32(get(row, "building_own")),
						ArStatus:            parseInt32WithCodes(get(row, "ar_status"), arStatusCodes),
						ArTotal:             parseFloat64(get(row, "ar_total")),
						IsEmbBail:           parseBoolWithDefault(get(row, "is_emb_bail"), false),
						IsActive:            activeFlag,
						TaxName:             get(row, "tax_name"),
						TaxAddr1:            get(row, "tax_addr1"),
						TaxAddr2:            get(row, "tax_addr2"),
						TaxCity:             get(row, "tax_city"),
						TaxNo:               get(row, "tax_no"),
						TaxInvoiceForm:      parseInt32WithCodes(get(row, "tax_invoice_form"), taxInvoiceFormCodes),
						OwnerName:           get(row, "owner_name"),
						OwnerAddr1:          get(row, "owner_addr1"),
						OwnerAddr2:          get(row, "owner_addr2"),
						OwnerCity:           get(row, "owner_city"),
						OwnerPhoneNo:        get(row, "owner_phone_no"),
						OwnerIdNo:           get(row, "owner_id_no"),
						DelvAddr1:           get(row, "delv_addr1"),
						DelvAddr2:           get(row, "delv_addr2"),
						DelvCity:            get(row, "delv_city"),
						DelvCity2:           get(row, "delv_city2"),
						DelvWardId:          get(row, "delv_ward_id"),
						DelvWardId2:         get(row, "delv_ward_id2"),
						DelvZipCode:         get(row, "delv_zip_code"),
						DelvZipCode2:        get(row, "delv_zip_code2"),
						DelvIsSameAddr:      parseBoolWithDefault(get(row, "delv_is_same_addr"), false),
						DelvLatitude:        get(row, "delv_latitude"),
						DelvLongitude:       get(row, "delv_longitude"),
						DelvLatitude2:       get(row, "delv_latitude2"),
						DelvLongitude2:      get(row, "delv_longitude2"),
						InvAddr1:            get(row, "inv_addr1"),
						InvAddr2:            get(row, "inv_addr2"),
						InvCity:             get(row, "inv_city"),
						InvWardId:           get(row, "inv_ward_id"),
						InvZipCode:          get(row, "inv_zip_code"),
						InvIsSameAddr:       parseBoolWithDefault(get(row, "inv_is_same_addr"), false),
						Latitude:            get(row, "latitude"),
						Longitude:           get(row, "longitude"),
						ImageUrl:            get(row, "image_url"),
						IsObs:               parseBoolWithDefault(get(row, "is_obs"), false),
						Obs:                 parseInt64(get(row, "obs")),
						ObsType:             parseInt64WithCodes(get(row, "obs_type"), obsTypeCodes),
						ObsLimitAction:      parseInt64WithCodes(get(row, "obs_limit_action"), limitActionCodes),
						IsWaNo:              parseBoolWithDefault(get(row, "is_wa_no"), false),
						VerificationStatus:  parseInt32WithCodes(get(row, "verification_status"), verificationStatusCodes),
						VerifiedBy:          parseInt64(get(row, "verified_by")),
						UpdatedBy:           req.UserId,
					}
					if v := strings.TrimSpace(get(row, "is_pkp_outlet")); v != "" {
						b := parseBoolWithDefault(v, false)
						updateReq.IsPkpOutlet = &b
					}

					if t, ok := parseDate(get(row, "first_trans_date")); ok {
						updateReq.FirstTransDate = t
					}
					if t, ok := parseDate(get(row, "last_trans_date")); ok {
						updateReq.LastTransDate = t
					}
					if t, ok := parseDate(get(row, "ot_start_date")); ok {
						updateReq.OtStartDate = t
					}
					if t, ok := parseDate(get(row, "ot_reg_date")); ok {
						updateReq.OtRegDate = t
					}
					if t, ok := parseDate(get(row, "dob")); ok {
						updateReq.Dob = t
					}
					if t, ok := parseDate(get(row, "closed_date")); ok {
						updateReq.ClosedDate = t
					}
					if t, ok := parseDate(get(row, "outlet_establishment_date")); ok {
						updateReq.OutletEstablishmentDate = t
					}
					if t, ok := parseDate(get(row, "verified_at")); ok {
						updateReq.VerifiedAt = t
					}
					if err := s.OutletRepository.UpdateOutletImport(updateReq); err != nil {
						addRowError("gagal memperbarui data outlet: %v", err)
					}
				}
			} else {
				addRowError("outlet_id / outlet_code wajib diisi")
			}
		}

		if len(missingMasters) > 0 {
			rowErrors = append(rowErrors, strings.Join(missingMasters, ", "))
		}

		if len(rowErrors) > 0 {
			prefix := formatOutletValidationPrefix(get(row, "outlet_code"), i+2)
			errMsg := fmt.Sprintf("%s: %s", prefix, strings.Join(rowErrors, "; "))
			temp := entity.ImportOutletUpdateTemp{
				HistoryID:      historyId,
				CustID:         req.CustId,
				OutletLocCode:  get(row, "outlet_loc_code"),
				OutletLocName:  get(row, "outlet_loc_name"),
				OutletTypeCode: get(row, "outlet_type_code"),
				OutletTypeName: get(row, "outlet_type_name"),
				OutletGrpCode:  get(row, "outlet_grp_code"),
				OutletGrpName:  get(row, "outlet_grp_name"),
				DistrictCode:   get(row, "district_code"),
				DistrictName:   get(row, "district_name"),
				OtClassCode:    get(row, "ot_class_code"),
				OtClassName:    get(row, "ot_class_name"),
				DiscGrpCode:    get(row, "disc_grp_code"),
				DiscGrpName:    get(row, "disc_grp_name"),
				MarketCode:     get(row, "market_code"),
				MarketName:     get(row, "market_name"),
				IndustryCode:   get(row, "industry_code"),
				IndustryName:   get(row, "industry_name"),
				PriceGrpCode:   get(row, "price_grp_code"),
				PriceGrpName:   get(row, "price_grp_name"),
				BankCode:       get(row, "bank_code"),
				BankName:       get(row, "bank_name"),

				// tambahan kolom baru dari DDL
				OutletCode:              get(row, "outlet_code"),
				OutletName:              get(row, "outlet_name"),
				Address1:                get(row, "address1"),
				OutletProvince:          get(row, "outlet_province"),
				OutletRegency:           get(row, "outlet_regency"),
				OutletSubDistrict:       get(row, "outlet_sub_district"),
				OutletWard:              get(row, "outlet_ward"),
				ZipCode:                 get(row, "zip_code"),
				Longitude:               get(row, "longitude"),
				Latitude:                get(row, "latitude"),
				BuildingOwn:             get(row, "building_own"),
				OutletEstablishmentDate: get(row, "outlet_establishment_date"),

				PhoneNo: get(row, "phone_no"),
				FaxNo:   get(row, "fax_no"),
				Barcode: get(row, "barcode"),

				ContactName:     get(row, "contact_name"),
				JobTitle:        get(row, "job_title"),
				IdentityType:    get(row, "identity_type"),
				IdentityNo:      get(row, "identity_no"),
				ContactPhoneNo:  get(row, "contact_phone_no"),
				ContactIsWaNo:   get(row, "contact_is_wa_no"),
				ContactWaNo:     get(row, "contact_wa_no"),
				ContactEmail:    get(row, "contact_email"),
				OutletContactID: get(row, "outlet_contact_id"),

				TaxInvoiceFormName: get(row, "tax_invoice_form_name"),
				TaxIdentifierType:  get(row, "tax_identifier_type"),
				TaxIdentifierNo:    get(row, "tax_identifier_no"),
				Nitku:              get(row, "nitku"),
				TaxName:            get(row, "tax_name"),
				AddressTax:         get(row, "address_tax"),
				OutletTaxID:        get(row, "outlet_tax_id"),

				IsContraBon: get(row, "is_contra_bon"),
				AgentFrom:   get(row, "agent_from"),

				DelvAddr1:       get(row, "delv_addr1"),
				DelvProvince:    get(row, "delv_province"),
				DelvRegency:     get(row, "delv_regency"),
				DelvSubDistrict: get(row, "delv_sub_district"),
				DelvWard:        get(row, "delv_ward"),
				DelvLongitude:   get(row, "delv_longitude"),
				DelvLatitude:    get(row, "delv_latitude"),
				DelvZipCode:     get(row, "delv_zip_code"),
				DelvIsSameAddr:  get(row, "delv_is_same_addr"),

				InvAddr1:       get(row, "inv_addr1"),
				InvProvince:    get(row, "inv_province"),
				InvRegency:     get(row, "inv_regency"),
				InvSubDistrict: get(row, "inv_sub_district"),
				InvWard:        get(row, "inv_ward"),
				InvZipCode:     get(row, "inv_zip_code"),
				InvIsSameAddr:  get(row, "inv_is_same_addr"),

				PaymentTypeName: get(row, "payment_type_name"),
				ArStatusName:    get(row, "ar_status_name"),

				AccountNo:    get(row, "account_no"),
				AccountName:  get(row, "account_name"),
				OutletBankID: get(row, "outlet_bank_id"),

				Top:                     get(row, "top"),
				CreditLimitTypeName:     get(row, "credit_limit_type_name"),
				CreditLimit:             get(row, "credit_limit"),
				CreditLimitActionName:   get(row, "credit_limit_action_name"),
				SalesInvLimitTypeName:   get(row, "sales_inv_limit_type_name"),
				SalesInvLimit:           get(row, "sales_inv_limit"),
				SalesInvLimitActionName: get(row, "sales_inv_limit_action_name"),
				ObsTypeName:             get(row, "obs_type_name"),
				Obs:                     get(row, "obs"),
				ObsLimitActionName:      get(row, "obs_limit_action_name"),

				OutletID: get(row, "outlet_id"),

				StatusInsert: "failed",
				ErrorMessage: errMsg,
			}
			if err := s.OutletRepository.InsertOutletUpdateTemp(temp); err != nil {
				logrus.Errorf("gagal simpan ke outlet_update_temp baris %d: %v", i+2, err)
			}
			continue
		}

		success++

		if time.Since(start) > 10*time.Minute {
			logrus.Warnf("Timeout tercapai, sisa data dialihkan ke outlet_update_temp mulai baris %d", i+2)
			for j, r := range rows[i+2:] {
				temp := s.mapRowToOutletUpdateTemp(historyId, req.CustId, r, headerMap, "timeout: proses lebih dari 1 menit")
				if err := s.OutletRepository.InsertOutletUpdateTemp(temp); err != nil {
					logrus.Errorf("gagal simpan ke outlet_update_temp baris %d: %v", j+i+3, err)
				}
			}
			break
		}
	}

	failed := total - success
	if err := s.OutletRepository.UpdateImportHistoryOutlet(historyId, success, failed, false); err != nil {
		return err
	}
	return nil
}

func (s *outletServiceImpl) mapRowToOutletUpdateTemp(historyId int64, custId string, row []string, headerMap map[string]int, errMsg string) entity.ImportOutletUpdateTemp {
	get := func(key string) string {
		if idx, ok := headerMap[canonicalizeHeader(key)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}
	return entity.ImportOutletUpdateTemp{
		HistoryID:      historyId,
		CustID:         custId,
		OutletLocCode:  get("outlet_loc_code"),
		OutletLocName:  get("outlet_loc_name"),
		OutletTypeCode: get("outlet_type_code"),
		OutletTypeName: get("outlet_type_name"),
		OutletGrpCode:  get("outlet_grp_code"),
		OutletGrpName:  get("outlet_grp_name"),
		DistrictCode:   get("district_code"),
		DistrictName:   get("district_name"),
		OtClassCode:    get("ot_class_code"),
		OtClassName:    get("ot_class_name"),
		DiscGrpCode:    get("disc_grp_code"),
		DiscGrpName:    get("disc_grp_name"),
		MarketCode:     get("market_code"),
		MarketName:     get("market_name"),
		IndustryCode:   get("industry_code"),
		IndustryName:   get("industry_name"),
		PriceGrpCode:   get("price_grp_code"),
		PriceGrpName:   get("price_grp_name"),
		BankCode:       get("bank_code"),
		BankName:       get("bank_name"),
		StatusInsert:   "failed",
		ErrorMessage:   errMsg,
	}
}

func (s *outletServiceImpl) processUpdateRowsOld(req entity.ImportRequest, rows [][]string) error {
	if len(rows) < 2 {
		return errors.New("file harus berisi header dan data")
	}

	// baca header
	headers := rows[0]
	headerMap := map[string]int{}
	for i, h := range headers {
		headerMap[canonicalizeHeader(h)] = i
	}
	logrus.Info("Header Map:", headerMap)

	// helper function untuk ambil value berdasarkan header
	get := func(row []string, key string) string {
		if idx, ok := headerMap[canonicalizeHeader(key)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// loop baris data
	for i, row := range rows[1:] {
		// --- BANK ---
		// Mendukung 2 skenario kunci:
		// 1) bank_id ada -> pakai seperti semula
		// 2) bank_id kosong, tapi bank_code ada -> cari bank_id berdasarkan bank_code
		{
			rawBankId := strings.TrimSpace(get(row, "bank_id"))
			bankCodeKey := strings.TrimSpace(get(row, "bank_code"))
			bankCodeNew := strings.TrimSpace(get(row, "bank_code_new"))
			if bankCodeNew == "" {
				bankCodeNew = bankCodeKey
			}

			var bankId int64
			var err error
			if rawBankId != "" {
				bankId, _ = strconv.ParseInt(rawBankId, 10, 64)
			}

			if bankId == 0 && bankCodeKey != "" {
				// Cari ID berdasarkan code, agar bisa update dengan patokan bank_code
				bankId, err = s.OutletRepository.FindBankIdByCode(req.CustId, bankCodeKey)
				if err != nil || bankId == 0 {
					return fmt.Errorf("baris %d: bank_code %s tidak ditemukan", i+2, bankCodeKey)
				}
			}

			if bankId > 0 {
				exists, err := s.OutletRepository.CheckBankExists(req.CustId, bankId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan bank gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: bank_id %d tidak ditemukan", i+2, bankId)
				}

				// Jika user mengirim bank_code baru, cek duplikat (exclude dirinya sendiri)
				if code := bankCodeNew; code != "" {
					dup, err := s.OutletRepository.CheckBankCodeDuplicate(req.CustId, bankId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi bank_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: bank_code %s sudah digunakan bank lain", i+2, code)
					}
				}

				if err := s.OutletRepository.UpdateImportBank(req.CustId, bankCodeNew, get(row, "bank_name"), bankId); err != nil {
					return fmt.Errorf("baris %d: pembaruan data bank gagal: %w", i+2, err)
				}
			}
		}

		// --- OUTLET CLASS --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "ot_class_id"))
			codeKey := strings.TrimSpace(get(row, "ot_class_code"))
			codeNew := strings.TrimSpace(get(row, "ot_class_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var otClassId int64
			if rawId != "" {
				otClassId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if otClassId == 0 && codeKey != "" {
				var err error
				otClassId, err = s.OutletRepository.FindOutletClassIdByCode(req.CustId, codeKey)
				if err != nil || otClassId == 0 {
					return fmt.Errorf("baris %d: ot_class_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if otClassId > 0 {
				exists, err := s.OutletRepository.CheckOutletClassExists(req.CustId, otClassId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan outlet_class gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: ot_class_id %d tidak ditemukan", i+2, otClassId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckOutletClassCodeDuplicate(req.CustId, otClassId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi outlet_class_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: ot_class_code %s sudah digunakan outlet class lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportOutletClass(req.CustId, codeNew, get(row, "ot_class_name"), otClassId); err != nil {
					return fmt.Errorf("baris %d: pembaruan outlet_class gagal: %w", i+2, err)
				}
			}
		}

		// --- OUTLET GROUP --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "ot_grp_id"))
			codeKey := strings.TrimSpace(get(row, "ot_grp_code"))
			codeNew := strings.TrimSpace(get(row, "ot_grp_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var otGrpId int64
			if rawId != "" {
				otGrpId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if otGrpId == 0 && codeKey != "" {
				var err error
				otGrpId, err = s.OutletRepository.FindOutletGroupIdByCode(req.CustId, codeKey)
				if err != nil || otGrpId == 0 {
					return fmt.Errorf("baris %d: ot_grp_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if otGrpId > 0 {
				exists, err := s.OutletRepository.CheckOutletGroupExists(req.CustId, otGrpId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan outlet_group gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: ot_grp_id %d tidak ditemukan", i+2, otGrpId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckOutletGroupCodeDuplicate(req.CustId, otGrpId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi outlet_group_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: ot_grp_code %s sudah digunakan outlet group lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportOutletGroup(req.CustId, codeNew, get(row, "ot_grp_name"), otGrpId); err != nil {
					return fmt.Errorf("baris %d: pembaruan outlet_group gagal: %w", i+2, err)
				}
			}
		}

		// --- OUTLET LOC --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "ot_loc_id"))
			codeKey := strings.TrimSpace(get(row, "ot_loc_code"))
			codeNew := strings.TrimSpace(get(row, "ot_loc_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var otLocId int64
			if rawId != "" {
				otLocId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if otLocId == 0 && codeKey != "" {
				var err error
				otLocId, err = s.OutletRepository.FindOutletLocIdByCode(req.CustId, codeKey)
				if err != nil || otLocId == 0 {
					return fmt.Errorf("baris %d: ot_loc_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if otLocId > 0 {
				exists, err := s.OutletRepository.CheckOutletLocExists(req.CustId, otLocId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan outlet_location gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: ot_loc_id %d tidak ditemukan", i+2, otLocId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckOutletLocCodeDuplicate(req.CustId, otLocId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi outlet_loc_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: ot_loc_code %s sudah digunakan outlet loc lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportOutletLoc(req.CustId, codeNew, get(row, "ot_loc_name"), otLocId); err != nil {
					return fmt.Errorf("baris %d: pembaruan outlet_location gagal: %w", i+2, err)
				}
			}
		}

		// --- OUTLET TYPE --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "ot_type_id"))
			codeKey := strings.TrimSpace(get(row, "ot_type_code"))
			codeNew := strings.TrimSpace(get(row, "ot_type_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var otTypeId int64
			if rawId != "" {
				otTypeId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if otTypeId == 0 && codeKey != "" {
				var err error
				otTypeId, err = s.OutletRepository.FindOutletTypeIdByCode(req.CustId, codeKey)
				if err != nil || otTypeId == 0 {
					return fmt.Errorf("baris %d: ot_type_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if otTypeId > 0 {
				exists, err := s.OutletRepository.CheckOutletTypeExists(req.CustId, otTypeId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan outlet_type gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: ot_type_id %d tidak ditemukan", i+2, otTypeId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckOutletTypeCodeDuplicate(req.CustId, otTypeId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi outlet_type_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: ot_type_code %s sudah digunakan outlet type lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportOutletType(req.CustId, codeNew, get(row, "ot_type_name"), otTypeId); err != nil {
					return fmt.Errorf("baris %d: pembaruan outlet_type gagal: %w", i+2, err)
				}
			}
		}

		// --- DISTRICT --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "district_id"))
			codeKey := strings.TrimSpace(get(row, "district_code"))
			codeNew := strings.TrimSpace(get(row, "district_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var districtId int64
			if rawId != "" {
				districtId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if districtId == 0 && codeKey != "" {
				var err error
				districtId, err = s.OutletRepository.FindDistrictIdByCode(req.CustId, codeKey)
				if err != nil || districtId == 0 {
					return fmt.Errorf("baris %d: district_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if districtId > 0 {
				exists, err := s.OutletRepository.CheckDistrictExists(req.CustId, districtId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan district gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: district_id %d tidak ditemukan", i+2, districtId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckDistrictCodeDuplicate(req.CustId, districtId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi district_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: district_code %s sudah digunakan district lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportDistrict(req.CustId, codeNew, get(row, "district_name"), districtId); err != nil {
					return fmt.Errorf("baris %d: pembaruan district gagal: %w", i+2, err)
				}
			}
		}

		// --- DISC GROUP --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "disc_grp_id"))
			codeKey := strings.TrimSpace(get(row, "disc_grp_code"))
			codeNew := strings.TrimSpace(get(row, "disc_grp_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var discGrpId int64
			if rawId != "" {
				discGrpId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if discGrpId == 0 && codeKey != "" {
				var err error
				discGrpId, err = s.OutletRepository.FindDiscGroupIdByCode(req.CustId, codeKey)
				if err != nil || discGrpId == 0 {
					return fmt.Errorf("baris %d: disc_grp_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if discGrpId > 0 {
				exists, err := s.OutletRepository.CheckDiscGroupExists(req.CustId, discGrpId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan disc_group gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: disc_grp_id %d tidak ditemukan", i+2, discGrpId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckDiscGroupCodeDuplicate(req.CustId, discGrpId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi disc_group_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: disc_grp_code %s sudah digunakan disc group lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportDiscGroup(req.CustId, codeNew, get(row, "disc_grp_name"), discGrpId); err != nil {
					return fmt.Errorf("baris %d: pembaruan disc_group gagal: %w", i+2, err)
				}
			}
		}

		// --- MARKET --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "market_id"))
			codeKey := strings.TrimSpace(get(row, "market_code"))
			codeNew := strings.TrimSpace(get(row, "market_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var marketId int64
			if rawId != "" {
				marketId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if marketId == 0 && codeKey != "" {
				var err error
				marketId, err = s.OutletRepository.FindMarketIdByCode(req.CustId, codeKey)
				if err != nil || marketId == 0 {
					return fmt.Errorf("baris %d: market_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if marketId > 0 {
				exists, err := s.OutletRepository.CheckMarketExists(req.CustId, marketId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan market gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: market_id %d tidak ditemukan", i+2, marketId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckMarketCodeDuplicate(req.CustId, marketId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi market_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: market_code %s sudah digunakan market lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportMarket(req.CustId, codeNew, get(row, "market_name"), marketId); err != nil {
					return fmt.Errorf("baris %d: pembaruan market gagal: %w", i+2, err)
				}
			}
		}

		// --- INDUSTRY --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "industry_id"))
			codeKey := strings.TrimSpace(get(row, "industry_code"))
			codeNew := strings.TrimSpace(get(row, "industry_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var industryId int64
			if rawId != "" {
				industryId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if industryId == 0 && codeKey != "" {
				var err error
				industryId, err = s.OutletRepository.FindIndustryIdByCode(req.CustId, codeKey)
				if err != nil || industryId == 0 {
					return fmt.Errorf("baris %d: industry_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if industryId > 0 {
				exists, err := s.OutletRepository.CheckIndustryExists(req.CustId, industryId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan industry gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: industry_id %d tidak ditemukan", i+2, industryId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckIndustryCodeDuplicate(req.CustId, industryId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi industry_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: industry_code %s sudah digunakan industry lain", i+2, code)
					}
				}
				if err := s.OutletRepository.UpdateImportIndustry(req.CustId, codeNew, get(row, "industry_name"), industryId); err != nil {
					return fmt.Errorf("baris %d: pembaruan industry gagal: %w", i+2, err)
				}
			}
		}

		// --- PRICE GROUP --- (support code fallback)
		{
			rawId := strings.TrimSpace(get(row, "price_grp_id"))
			codeKey := strings.TrimSpace(get(row, "price_grp_code"))
			codeNew := strings.TrimSpace(get(row, "price_grp_code_new"))
			if codeNew == "" {
				codeNew = codeKey
			}
			var priceGrpId int64
			if rawId != "" {
				priceGrpId, _ = strconv.ParseInt(rawId, 10, 64)
			}
			if priceGrpId == 0 && codeKey != "" {
				var err error
				priceGrpId, err = s.OutletRepository.FindPriceGroupIdByCode(req.CustId, codeKey)
				if err != nil || priceGrpId == 0 {
					return fmt.Errorf("baris %d: price_grp_code %s tidak ditemukan", i+2, codeKey)
				}
			}
			if priceGrpId > 0 {
				exists, err := s.OutletRepository.CheckPriceGroupExists(req.CustId, priceGrpId)
				if err != nil {
					return fmt.Errorf("baris %d: pengecekan price_group gagal: %w", i+2, err)
				}
				if !exists {
					return fmt.Errorf("baris %d: price_grp_id %d tidak ditemukan", i+2, priceGrpId)
				}
				if code := codeNew; code != "" {
					dup, err := s.OutletRepository.CheckPriceGroupCodeDuplicate(req.CustId, priceGrpId, code)
					if err != nil {
						return fmt.Errorf("baris %d: validasi price_grp_code gagal: %w", i+2, err)
					}
					if dup {
						return fmt.Errorf("baris %d: price_grp_code %s sudah digunakan price group lain", i+2, code)
					}
				}
				_, err = s.OutletRepository.UpdateImportPriceGroup(req.CustId, priceGrpId, codeNew, get(row, "price_grp_name"))
				if err != nil {
					return fmt.Errorf("baris %d: pembaruan price_group gagal: %w", i+2, err)
				}
			}
		}
	}

	return nil
}

// func strPtr(v *string) string {
// 	if v == nil {
// 		return ""
// 	}
// 	return *v
// }

func floatPtr(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func intPtr(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func timePtr(v *time.Time) string {
	if v == nil {
		return ""
	}
	return v.Format("2006-01-02 15:04:05")
}

func floatPtrStr(v *float64) string {
	if v == nil {
		return "0"
	}
	return strconv.FormatFloat(*v, 'f', -1, 64)
}

func intPtrStr(v *int) string {
	if v == nil {
		return "0"
	}
	return strconv.Itoa(*v)
}
func (service *outletServiceImpl) ListByDistributor(dataFilter entity.OutletQueryFilter) (data []entity.OutletListByDistributorRespone, total int, lastPage int, err error) {

	if len(dataFilter.DistributorID) == 0 {
		return data, 0, 0, fmt.Errorf("distributor_id is required for list-by-distributor")
	}
	mCustomer, err := service.OutletRepository.FindOneCustomerByDistributorID(int64(dataFilter.DistributorID[0]))
	if err != nil {
		return data, total, lastPage, err
	}
	dataFilter.CustId = mCustomer.CustId

	outlets, total, lastPage, err := service.OutletRepository.FindAllByCustId(dataFilter, mCustomer.CustId, mCustomer.ParentCustId)
	if err != nil {
		return data, total, lastPage, err
	}

	if len(outlets) > 0 {
		for _, row := range outlets {
			var vResp entity.OutletListByDistributorRespone
			structs.Automapper(row, &vResp)
			data = append(data, vResp)
		}
	}

	return data, total, lastPage, err
}

// OutletListApproval retrieves a paginated list of outlet change requests filtered by status and customer ID.
// Sets default values for pagination (page=1, limit=5) and sorting (created_date:desc) if not provided.
// Converts repository models to response entities and returns the list with pagination metadata.
func (service *outletServiceImpl) OutletListApproval(dataFilter entity.OutletListApprovalQueryFilter, custId string) (data []entity.OutletListApprovalResponse, total int, lastPage int, err error) {
	if dataFilter.Page < 1 {
		dataFilter.Page = 1
	}
	if dataFilter.Limit < 1 {
		dataFilter.Limit = 5
	}
	if dataFilter.Sort == "" {
		dataFilter.Sort = "created_date:desc"
	}

	outlets, total, lastPage, err := service.OutletRepository.FindAllOutletCrByStatus(dataFilter, custId)
	if err != nil {
		return nil, 0, 0, err
	}

	if len(outlets) == 0 {
		return []entity.OutletListApprovalResponse{}, total, lastPage, nil
	}

	data = make([]entity.OutletListApprovalResponse, 0, len(outlets))
	for _, row := range outlets {
		var vResp entity.OutletListApprovalResponse
		structs.Automapper(row, &vResp)
		data = append(data, vResp)
	}

	return data, total, lastPage, nil
}

// ApproveOutletList processes approval or rejection of outlet change requests.
// Validates that all requests belong to the current customer and are in pending status (status=1).
// If approved (status=2), updates the outlet location with new latitude/longitude from change request details.
// Uses transaction with defer rollback to ensure data consistency on error.
func (service *outletServiceImpl) ApproveOutletList(request entity.OutletListApprovalRequest, custId string, userId int64) error {
	trx, err := service.OutletRepository.TrxBegin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			trx.TrxRollback()
		}
	}()

	outletCrs, err := trx.GetOutletCrByOutletCrIds(request.OutletCrId)
	if err != nil {
		return err
	}

	if len(outletCrs) == 0 {
		return errors.New("outlet_cr_id not found")
	}

	for _, ocr := range outletCrs {
		if ocr.CustId != custId {
			return errors.New("outlet_cr_id does not belong to current customer")
		}
		if ocr.Status != 1 {
			return errors.New("only pending status can be approved/rejected")
		}
	}

	err = trx.UpdateOutletCrStatus(request.OutletCrId, request.Status, userId)
	if err != nil {
		return err
	}

	if request.Status == 2 {
		details, err := trx.GetOutletCrDetByOutletCrIds(request.OutletCrId)
		if err != nil {
			return err
		}

		outletCrMap := make(map[int64]map[string]string)
		for _, detail := range details {
			if outletCrMap[detail.OutletCrId] == nil {
				outletCrMap[detail.OutletCrId] = make(map[string]string)
			}
			if detail.FieldName == "latitude" && detail.NewValue != nil {
				outletCrMap[detail.OutletCrId]["latitude"] = *detail.NewValue
			}
			if detail.FieldName == "longitude" && detail.NewValue != nil {
				outletCrMap[detail.OutletCrId]["longitude"] = *detail.NewValue
			}
		}

		for _, ocr := range outletCrs {
			if coords, exists := outletCrMap[ocr.OutletCrId]; exists {
				latitude := coords["latitude"]
				longitude := coords["longitude"]
				if latitude != "" && longitude != "" {
					err = trx.UpdateOutletLocationFromCr(ocr.OutletId, latitude, longitude)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	err = trx.TrxCommit()
	return err
}
