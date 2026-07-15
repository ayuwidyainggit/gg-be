package service

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sales/adapter"
	"sales/entity"
	"sales/model"
	"sales/pkg/config/env"
	"sales/pkg/constant"
	"sales/pkg/conversion"
	"sales/pkg/rabbitmq"
	"sales/pkg/str"
	"sales/pkg/structs"
	"sales/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const activityReportGeotagMaxMeters = 200.0

var ErrUnauthorizedCustID = errors.New("cust_id is outside authorized scope")
var publishReportMessage = rabbitmq.PublishMessage

func (service *reportServiceImpl) publishExportMessage(reportID string, rmqConfig rabbitmq.RmqConfig) error {
	err := publishReportMessage(&rmqConfig)
	if err != nil {
		service.markReportFailed(reportID, err)
		return err
	}

	return nil
}

func buildReportObjectFileName(reportName, reportID string) string {
	baseName := reportID
	if reportName != "" {
		baseName = reportName
	}
	return filepath.ToSlash(filepath.Join("reports", baseName+".xlsx"))
}

func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371000.0
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

func effectiveGeotagStatus(st sql.NullInt32, masterLon, masterLat, actualLon, actualLat float64) sql.NullInt32 {
	if st.Valid {
		return st
	}
	if (actualLon == 0 && actualLat == 0) || (masterLon == 0 && masterLat == 0) {
		return sql.NullInt32{Valid: false}
	}
	if masterLat < -90 || masterLat > 90 || math.Abs(masterLon) > 180 ||
		actualLat < -90 || actualLat > 90 || math.Abs(actualLon) > 180 {
		return sql.NullInt32{Valid: false}
	}
	d := haversineMeters(masterLat, masterLon, actualLat, actualLon)
	if d <= activityReportGeotagMaxMeters {
		return sql.NullInt32{Int32: 1, Valid: true}
	}
	return sql.NullInt32{Int32: 0, Valid: true}
}

func geotagStatusDescFrom(st sql.NullInt32) string {
	if !st.Valid {
		return ""
	}
	switch st.Int32 {
	case 1:
		return "Match"
	case 0:
		return "Mismatch"
	default:
		return ""
	}
}

func geotagStatusPtr(st sql.NullInt32) *int {
	if !st.Valid {
		return nil
	}
	v := int(st.Int32)
	return &v
}

func activityReportLocationActual(lon, lat float64) string {
	if lon == 0 && lat == 0 {
		return ""
	}
	return fmt.Sprintf("%v,%v", lon, lat)
}

func activityReportLocationPair(lat, lon string) string {
	lat = strings.TrimSpace(lat)
	lon = strings.TrimSpace(lon)
	if lat == "" || lon == "" {
		return ""
	}
	return lat + ", " + lon
}

func activityReportParseCoord(s string) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return v
}

func activityReportPjpStatus(isPlanned bool) string {
	if isPlanned {
		return "Planned"
	}
	return "Unplanned"
}

func activityReportVisitStatus(skipAt sql.NullInt64, arriveAt int64) string {
	if skipAt.Valid && skipAt.Int64 != 0 {
		return "Skipped"
	}
	if arriveAt != 0 {
		return "Visited"
	}
	return "Pending"
}

func activityReportCompliance(isPlanned bool, arriveAt int64, skipAt sql.NullInt64) string {
	if isPlanned && arriveAt != 0 && (!skipAt.Valid || skipAt.Int64 == 0) {
		return "Yes"
	}
	return "No"
}

func activityReportGeotagFromLabel(label string) *int {
	switch strings.TrimSpace(label) {
	case "Match":
		v := 1
		return &v
	case "Mismatch":
		v := 0
		return &v
	default:
		return nil
	}
}

func activityReportDurationString(minutes int64) string {
	if minutes <= 0 {
		return ""
	}
	return strconv.FormatInt(minutes, 10)
}

func isActivityReportPrincipalUser(authCustID, parentCustID string) bool {
	authCustID = strings.TrimSpace(authCustID)
	parentCustID = strings.TrimSpace(parentCustID)
	if authCustID == "" {
		return false
	}
	if parentCustID != "" {
		return authCustID == parentCustID
	}
	return len(authCustID) <= 6
}

func formatPrincipalActivityReportColumns(row model.SalesActivityReportRow) (businessUnitCode, distributorCode, distributorName string) {
	businessUnitCode = row.BusinessUnitCode
	distributorCode = row.DistributorCode
	distributorName = row.DistributorName

	if strings.TrimSpace(row.BusinessUnitCode) != "" {
		return businessUnitCode, distributorCode, distributorName
	}

	businessUnitCode = "-"
	if strings.TrimSpace(row.OutletCode) != "" {
		distributorCode = "-"
		distributorName = "-"
	}
	return businessUnitCode, distributorCode, distributorName
}

func mapActivityReportRow(row model.SalesActivityReportRow, isPrincipalUser bool) entity.ActivityReportListResp {
	pjpCodeInt, _ := strconv.Atoi(strings.TrimSpace(row.PJPCode))
	pjpCodeNew := fmt.Sprintf("%04d", pjpCodeInt)

	pjpStatus := row.PjpStatus
	if pjpStatus == "" {
		pjpStatus = activityReportPjpStatus(row.IsPlanned)
	}
	visitStatus := row.VisitStatus
	if visitStatus == "" {
		arriveFlag := int64(0)
		if row.CheckinTime != "" {
			arriveFlag = 1
		}
		visitStatus = activityReportVisitStatus(row.SkipAt, arriveFlag)
	}
	compliance := row.Compliance
	if compliance == "" {
		arriveFlag := int64(0)
		if row.CheckinTime != "" {
			arriveFlag = 1
		}
		compliance = activityReportCompliance(row.IsPlanned, arriveFlag, row.SkipAt)
	}

	businessUnitCode := row.BusinessUnitCode
	distributorCode := row.DistributorCode
	distributorName := row.DistributorName
	if isPrincipalUser {
		businessUnitCode, distributorCode, distributorName = formatPrincipalActivityReportColumns(row)
	}

	return entity.ActivityReportListResp{
		BusinessUnitCode:    businessUnitCode,
		BusinessUnitName:    row.BusinessUnitName,
		DistributorCode:     distributorCode,
		DistributorName:     distributorName,
		PJPCode:             pjpCodeNew,
		EmployeeCode:        row.EmpCode,
		SalesmanName:        row.SalesmanName,
		OutletCode:          row.OutletCode,
		OutletPrincipalCode: row.OutletPrincipalCode,
		OutletName:          row.OutletName,
		Date:                str.FormatTimeToDateString(row.VisitDate),
		ClockInTime:         row.ClockInTime,
		ClockOutTime:        row.ClockOutTime,
		CheckinTime:         row.CheckinTime,
		CheckoutTime:        row.CheckoutTime,
		Duration:            activityReportDurationString(row.DurationMinutes),
		PjpStatus:           pjpStatus,
		VisitStatus:         visitStatus,
		Compliance:          compliance,
		SalesValue:          row.SalesValue,
		ReturnValue:         row.ReturnValue,
		PaymentCollected:    row.PaymentValue,
		LocationMaster:      row.LocationMaster,
		LocationActual:      row.LocationActual,
		GeotagStatus:        activityReportGeotagFromLabel(row.GeotagStatusLabel),
		GeotagStatusDesc:    row.GeotagStatusLabel,
		Remarks:             activityReportRemarksOrDefault(row.Remarks),
	}
}

func activityReportRemarksOrDefault(remarks string) string {
	trimmed := strings.TrimSpace(remarks)
	if trimmed == "" {
		return "-"
	}
	return trimmed
}

func resolveSecondaryDashboardYear(year *int) int {
	if year == nil {
		return time.Now().Year()
	}

	return *year
}

func secondarySalesExcelHeaders() []interface{} {
	return []interface{}{
		"DistributorCode", "DisributorName", "trx type",
		"DocumentNo", "DocumentDate",
		"OutletCode", "OutletName",
		"SalesmanCode", "EmpName",
		"SupCode", "SupName",
		"ProCode", "ProName",
		"Price Unit 3", "Price Unit 2", "Price Unit 1",
		"UnitID3", "UnitID2", "UnitID1",
		"ConvUnit3", "ConvUnit2",
		"Qty3Final", "Qty2Final", "Qty1Final",
		"GrossSales", "SpecialDiscount", "Discount",
		"NetSalesExcPPN", "PPN", "NetSalesIncPPN",
	}
}

func secondarySalesExcelRow(row model.SecondarySalesReportUnion) []interface{} {
	documentDate := ""
	if !row.DocumentDate.IsZero() {
		documentDate = row.DocumentDate.Format(constant.DD_MM_YYYY)
	}

	return []interface{}{
		row.DistributorCode,
		row.DistributorName,
		row.TrxType,
		row.DocumentNo,
		documentDate,
		row.OutletCode,
		row.OutletName,
		row.EmpCode,
		row.EmpName,
		row.SupCode,
		row.SupName,
		row.ProCode,
		row.ProName,
		row.SellPrice3,
		row.SellPrice2,
		row.SellPrice1,
		row.UnitID3,
		row.UnitID2,
		row.UnitID1,
		row.ConvUnit3,
		row.ConvUnit2,
		row.Qty3,
		row.Qty2,
		row.Qty1,
		row.GrossSales,
		row.SpecialDiscount,
		row.Discount,
		row.NetSalesExcPPN,
		row.PPN,
		row.NetSalesIncPPN,
	}
}

func (service *reportServiceImpl) markReportFailed(reportID string, cause error) {
	if reportID == "" {
		return
	}

	if cause != nil {
		log.Error("Report export failed:", cause.Error())
	}

	failedUpdate := model.ReportList{ReportID: reportID, FileStatus: entity.FILE_STATUS_FAILED, FileURL: ""}
	if err := service.ReportRepository.UpdateReportByReportID(context.Background(), reportID, &failedUpdate); err != nil {
		log.Error("Report export mark failed update error:", err.Error())
	}
}

type ReportService interface {
	PublishSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) (data entity.ReportList, err error)
	List(dataFilter entity.ReportQueryFilter) (data []entity.ReportList, total int64, lastPage int, err error)
	SubscribeSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) error
	PublishActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (data entity.ReportList, err error)
	SubscribeActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (err error)
	PublishActivitySalesReportList(dataFilter entity.ActivityReportQueryFilterList) (results []entity.ActivityReportListResp, total int64, lastPage int, err error)
	ExtractReportSecondary(req entity.SecondarySalesReportDashboardExtractQueryFilter) (err error)
	SecondarySalesReportSumReportByMonth(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (data entity.SumReportByMonthModelResp, err error)
	SecondarySalesReportGroupSales(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) (datas []entity.SecondarySalesReportGroupResp, err error)
	SecondarySalesReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.SumReportTrendSalesResp, err error)
	SalesmanActivityReportSumReportByMonth(authCustID, parentCustID string, req entity.SalesmanActivityReportDashboardSumPayload) (data entity.SalesmanActivityReportByMonthModelResp, err error)
	SalesmanActivityReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.ActivityReportTrendSalesResp, err error)
	SalesmanActivityReportGeotag(authCustID, parentCustID string, req entity.ActivityReportGeotagPayload) (data entity.ActivityReportGeotagResp, err error)
	SalesmanActivityReportGroupSales(authCustID, parentCustID string, req entity.SalesmanActivityReportDashboardGroupPayload) (datas []entity.SecondarySalesReportGroupResp, err error)
	SalesmanActivitySalesmanList(dataFilter entity.ActivityReportSalesmanListQueryFilter) (datas []entity.SalesmanActivityReportSalesmanListResp, err error)
}

func NewReportService(
	config env.ConfigEnv,
	reportRepository repository.ReportRepository,
	transaction repository.Dbtransaction,
	obsAdapter adapter.ObsAdapter,
) *reportServiceImpl {
	return &reportServiceImpl{
		Config:           config,
		ReportRepository: reportRepository,
		Transaction:      transaction,
		ObsAdapter:       obsAdapter,
	}
}

type reportServiceImpl struct {
	Config           env.ConfigEnv
	ReportRepository repository.ReportRepository
	StockRepository  repository.StockRepository
	Transaction      repository.Dbtransaction
	ObsAdapter       adapter.ObsAdapter
}

func (service *reportServiceImpl) PublishSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) (data entity.ReportList, err error) {
	// Resolve effective cust for transactional data filtering.
	// authCustID is preserved separately so reportList.CustID stays as the auth user.
	authCustID := dataFilter.CustID
	requestedCustIDs := dataFilter.CustIDs
	if len(requestedCustIDs) == 0 && len(dataFilter.RequestedCustIDs) > 0 {
		requestedCustIDs = []string(dataFilter.RequestedCustIDs)
	}
	if len(requestedCustIDs) == 0 && dataFilter.RequestedCustID != "" {
		requestedCustIDs = []string{dataFilter.RequestedCustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, dataFilter.ParentCustID, requestedCustIDs)
	if err != nil {
		return data, err
	}
	dataFilter.CustIDs = effectiveCustIDs
	if len(effectiveCustIDs) == 1 {
		dataFilter.CustID = effectiveCustIDs[0]
		dataFilter.RequestedCustID = effectiveCustIDs[0]
	} else {
		dataFilter.CustID = authCustID
		dataFilter.RequestedCustID = ""
	}
	dataFilter.RequestedCustIDs = effectiveCustIDs

	objectID := primitive.NewObjectID() // Generate a new ObjectID
	objectIDString := objectID.Hex()    // Convert ObjectID to string

	loc, _ := time.LoadLocation("Asia/Jakarta")
	dataFilter.ExportDate = time.Now().In(loc).Format("020106") // ddmmyy
	sequence := service.ReportRepository.CountSecondarySalesReportByDate(dataFilter)

	reportName := entity.REPORT_NAME_SECONDARY_SALES + "-" + dataFilter.ExportDate + "-" + fmt.Sprintf("%03d", sequence)
	log.Info("Generated file name:", reportName)

	var reportList model.ReportList
	reportList.ReportID = objectIDString
	// report.list.cust_id must remain the auth user so the principal sees the row in GET /v1/reports.
	reportList.CustID = authCustID
	reportList.CreatedBy = dataFilter.ExportBy
	reportList.StartDate = str.UnixTimestampToAsiaJakartaTime(*dataFilter.From)
	reportList.EndDate = str.UnixTimestampToAsiaJakartaTime(*dataFilter.To)
	reportList.ReportName = reportName
	reportList.FileStatus = entity.FILE_STATUS_PROCESSING

	err = service.ReportRepository.StoreReportList(context.Background(), &reportList)
	if err != nil {
		log.Error("SecondarySales, StoreReportList:", err.Error())
		return data, err
	}

	if err = structs.Automapper(reportList, &data); err != nil {
		log.Error("SecondarySales, Automapper:", err.Error())
		return data, err
	}

	data.StartDate = reportList.StartDate.Format(constant.YYYY_MM_DD)
	data.EndDate = reportList.EndDate.Format(constant.YYYY_MM_DD)

	dataFilter.ReportID = objectIDString
	rmqConfig := rabbitmq.RmqConfig{
		MessageID:      objectIDString,
		ExchangeName:   constant.RMQ_DEFAULT_EXCHANGE,
		RoutingKey:     constant.RMQ_SECONDARY_SALES_EXPORT,
		QueueName:      constant.RMQ_SECONDARY_SALES_EXPORT,
		DelayQueueName: constant.RMQ_SECONDARY_SALES_EXPORT + constant.RMQ_DEFAULT_DELAY_SUFFIX,
		MessageTTL:     service.Config.Get("REPORT_DELAY_SECONDS"),
		Message:        structs.StructToJson(dataFilter),
	}
	if err = service.publishExportMessage(objectIDString, rmqConfig); err != nil {
		log.Error("PublishSecondarySalesReport, publishExportMessage:", err.Error())
		return data, err
	}

	return data, err
}

func (service *reportServiceImpl) List(dataFilter entity.ReportQueryFilter) (data []entity.ReportList, total int64, lastPage int, err error) {
	reports, total, lastPage, err := service.ReportRepository.FindAllByCustID(dataFilter)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range reports {
		var vResp entity.ReportList
		structs.Automapper(row, &vResp)
		vResp.FileStatusName = vResp.GetFileStatusName()
		data = append(data, vResp)
	}

	return data, total, lastPage, err
}

func (service *reportServiceImpl) SubscribeSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) (err error) {
	defer func() {
		if err != nil {
			service.markReportFailed(dataFilter.ReportID, err)
		}
	}()

	requestedCustIDs := dataFilter.CustIDs
	if len(requestedCustIDs) == 0 && len(dataFilter.RequestedCustIDs) > 0 {
		requestedCustIDs = []string(dataFilter.RequestedCustIDs)
	}
	if len(requestedCustIDs) == 0 && dataFilter.RequestedCustID != "" {
		requestedCustIDs = []string{dataFilter.RequestedCustID}
	}
	if len(requestedCustIDs) == 0 && dataFilter.CustID != "" {
		requestedCustIDs = []string{dataFilter.CustID}
	}
	if normalized, err := entity.NormalizeStringList(requestedCustIDs); err != nil {
		return err
	} else {
		dataFilter.CustIDs = normalized
	}
	if len(dataFilter.CustIDs) == 1 {
		dataFilter.CustID = dataFilter.CustIDs[0]
	}

	report, err := service.ReportRepository.GetReportByReportID(dataFilter.ReportID)
	if err != nil {
		return err
	}
	secondaryReports, err := service.ReportRepository.SecondarySalesUnion(dataFilter)
	if err != nil {
		return err
	}

	// 1) Siapkan file Excel
	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName(f.GetSheetName(f.GetActiveSheetIndex()), sheet)
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			log.Error("SecondarySales, Close file:", closeErr.Error())
		}
	}()

	if err := sw.SetRow("A1", secondarySalesExcelHeaders()); err != nil {
		return err
	}

	rowIdx := 2
	for _, row := range secondaryReports {
		cell, _ := excelize.CoordinatesToCellName(1, rowIdx) // kolom A, baris rowIdx
		if err := sw.SetRow(cell, secondarySalesExcelRow(row)); err != nil {
			return err
		}
		rowIdx++
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	// (opsional) atur lebar kolom biar rapi
	_ = f.SetColWidth(sheet, "A", "Z", 16)

	// 3) Serialize ke buffer
	buf, err := f.WriteToBuffer() // menghasilkan *bytes.Buffer
	if err != nil {
		return err
	}

	// 4) Penamaan file + upload ke storage seperti sebelumnya (ubah ekstensi dan content-type)
	var reportList model.ReportList

	// loc, _ := time.LoadLocation("Asia/Jakarta")
	// dataFilter.ExportDate = time.Now().In(loc).Format("020106") // ddmmyy
	// sequence := service.ReportRepository.CountSecondarySalesReportByDate(dataFilter)
	//fileName := "SecondarySales-" + dataFilter.ExportDate + "-" + fmt.Sprintf("%03d", sequence) + ".xlsx"
	fileName := buildReportObjectFileName(report.ReportName, dataFilter.ReportID)
	log.Info("Generated file name:", fileName)

	upload := &model.UploadCsv{
		FileName:    fileName,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		FileData:    bytes.NewReader(buf.Bytes()),
	}

	downloadUrl, err := service.ObsAdapter.UploadFileCsv(upload)
	if err != nil {
		log.Error("SecondarySales, UploadFileCsv:", err.Error())
		return err
	}
	if downloadUrl == "" {
		return fmt.Errorf("secondary sales export upload returned empty file url")
	}
	log.Info("Uploaded Excel URL:", downloadUrl)

	reportList.ReportID = dataFilter.ReportID
	reportList.FileURL = downloadUrl
	reportList.FileStatus = entity.FILE_STATUS_READY

	if err = service.ReportRepository.UpdateReportByReportID(context.Background(), dataFilter.ReportID, &reportList); err != nil {
		log.Error("SecondarySales, UpdateReportByReportID:", err.Error())
		return err
	}

	return err
}

func (service *reportServiceImpl) applyActivityReportListCustIDs(dataFilter *entity.ActivityReportQueryFilterList) error {
	authCustID := strings.TrimSpace(dataFilter.AuthCustID)
	if authCustID == "" {
		authCustID = strings.TrimSpace(dataFilter.CustID)
	}
	if authCustID == "" {
		return nil
	}
	dataFilter.AuthCustID = authCustID

	requestedCustIDs := dataFilter.CustIDs
	if len(requestedCustIDs) == 0 && strings.TrimSpace(dataFilter.RequestedCustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(dataFilter.RequestedCustID)}
	}

	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, dataFilter.ParentCustID, requestedCustIDs)
	if err != nil {
		return err
	}
	dataFilter.CustIDs = effectiveCustIDs
	if len(effectiveCustIDs) == 1 {
		dataFilter.CustID = effectiveCustIDs[0]
	} else {
		dataFilter.CustID = authCustID
	}
	return nil
}

func (service *reportServiceImpl) applyActivityReportExportCustIDs(dataFilter *entity.ActivityReportQueryFilter) error {
	authCustID := strings.TrimSpace(dataFilter.AuthCustID)
	if authCustID == "" {
		authCustID = strings.TrimSpace(dataFilter.CustID)
	}
	if authCustID == "" {
		return nil
	}
	dataFilter.AuthCustID = authCustID

	requestedCustIDs := dataFilter.CustIDs
	if len(requestedCustIDs) == 0 && len(dataFilter.RequestedCustIDs) > 0 {
		requestedCustIDs = dataFilter.RequestedCustIDs
	}
	if len(requestedCustIDs) == 0 && strings.TrimSpace(dataFilter.RequestedCustID) != "" {
		requestedCustIDs = []string{strings.TrimSpace(dataFilter.RequestedCustID)}
	}

	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, dataFilter.ParentCustID, requestedCustIDs)
	if err != nil {
		return err
	}
	dataFilter.CustIDs = effectiveCustIDs
	dataFilter.RequestedCustIDs = effectiveCustIDs
	if len(effectiveCustIDs) == 1 {
		dataFilter.CustID = effectiveCustIDs[0]
		dataFilter.RequestedCustID = effectiveCustIDs[0]
	} else {
		dataFilter.CustID = authCustID
		dataFilter.RequestedCustID = ""
	}
	return nil
}

func (service *reportServiceImpl) PublishActivitySalesReportList(dataFilter entity.ActivityReportQueryFilterList) (results []entity.ActivityReportListResp, total int64, lastPage int, err error) {
	if err = service.applyActivityReportListCustIDs(&dataFilter); err != nil {
		return
	}
	reports, total, lastPage, err := service.ReportRepository.ActivitySalesReportList(dataFilter)
	if err != nil {
		return
	}

	results = make([]entity.ActivityReportListResp, 0, len(reports))
	isPrincipalUser := isActivityReportPrincipalUser(dataFilter.AuthCustID, dataFilter.ParentCustID)
	for _, report := range reports {
		results = append(results, mapActivityReportRow(report, isPrincipalUser))
	}
	return
}

func (service *reportServiceImpl) PublishActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (data entity.ReportList, err error) {
	authCustID := dataFilter.AuthCustID
	if authCustID == "" {
		authCustID = dataFilter.CustID
	}
	if err = service.applyActivityReportExportCustIDs(&dataFilter); err != nil {
		return data, err
	}

	objectID := primitive.NewObjectID() // Generate a new ObjectID
	objectIDString := objectID.Hex()    // Convert ObjectID to string

	sequence := service.ReportRepository.CountReportByDateAndReportName(authCustID, dataFilter.ExportDate, entity.TYPE_REPORT_SALESMAN_ACTIVITY_REPORT)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	dataFilter.ExportDate = time.Now().In(loc).Format("020106") // ddmmyy

	reportName := entity.TYPE_REPORT_SALESMAN_ACTIVITY_REPORT + "-" + dataFilter.ExportDate + "-" + fmt.Sprintf("%03d", sequence)
	log.Info("Generated file name:", reportName)

	fromDate, err := str.ConvertStringDateToTimeObject(dataFilter.FromDate)
	if err != nil {
		log.Error("activityReport, PublishActivitySalesReport:", err.Error())
		return data, err
	}

	toDate, err := str.ConvertStringDateToTimeObject(dataFilter.ToDate)
	if err != nil {
		log.Error("activityReport, PublishActivitySalesReport:", err.Error())
		return data, err
	}

	var reportList model.ReportList
	reportList.ReportID = objectIDString
	reportList.CustID = authCustID
	reportList.CreatedBy = dataFilter.ExportBy
	reportList.StartDate = fromDate
	reportList.EndDate = toDate
	reportList.ReportName = reportName
	reportList.FileStatus = entity.FILE_STATUS_PROCESSING

	err = service.ReportRepository.StoreReportList(context.Background(), &reportList)
	if err != nil {
		log.Error("activityReport, PublishActivitySalesReport:", err.Error())
		return data, err
	}

	if err = structs.Automapper(reportList, &data); err != nil {
		log.Error("activityReport, Automapper:", err.Error())
		return data, err
	}

	data.StartDate = reportList.StartDate.Format(constant.YYYY_MM_DD)
	data.EndDate = reportList.EndDate.Format(constant.YYYY_MM_DD)

	dataFilter.ReportID = objectIDString
	rmqConfig := rabbitmq.RmqConfig{
		MessageID:      objectIDString,
		ExchangeName:   constant.RMQ_DEFAULT_EXCHANGE,
		RoutingKey:     constant.RMQ_SALESMAN_ACTIVITY_REPORT_SALES_EXPORT,
		QueueName:      constant.RMQ_SALESMAN_ACTIVITY_REPORT_SALES_EXPORT,
		DelayQueueName: constant.RMQ_SALESMAN_ACTIVITY_REPORT_SALES_EXPORT + constant.RMQ_DEFAULT_DELAY_SUFFIX,
		MessageTTL:     service.Config.Get("REPORT_DELAY_SECONDS"),
		Message:        structs.StructToJson(dataFilter),
	}
	if err = service.publishExportMessage(objectIDString, rmqConfig); err != nil {
		log.Error("PublishActivitySalesReport, publishExportMessage:", err.Error())
		return data, err
	}

	return data, err
}

func (service *reportServiceImpl) SubscribeActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (err error) {
	if err := service.applyActivityReportExportCustIDs(&dataFilter); err != nil {
		return err
	}

	report, err := service.ReportRepository.GetReportByReportID(dataFilter.ReportID)
	if err != nil {
		return err
	}

	activitySalesReports, _, _, err := service.ReportRepository.ActivitySalesReport(dataFilter)
	if err != nil {
		return err
	}
	// 1) Siapkan file Excel
	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName(f.GetSheetName(f.GetActiveSheetIndex()), sheet)
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return err
	}

	headers := []interface{}{
		"Business Unit Code",
		"Business Unit Name",
		"PJP Code",
		"Employee Code",
		"Salesman Name",
		"Distributor Code",
		"Distributor Name",
		"Outlet Code",
		"Outlet Principal Code",
		"Outlet Name",
		"Date",
		"Clock-in Time",
		"Clock-out Time",
		"Check-in Time",
		"Check-out Time",
		"Duration (In minutes)",
		"PJP Status",
		"Visit Status",
		"Compliance",
		"Sales Value",
		"Return Value",
		"Payment Collected",
		"Location Master",
		"Location Actual",
		"Geotag Status",
		"Remarks",
	}
	// tulis header di A1
	if err := sw.SetRow("A1", headers); err != nil {
		return err
	}
	// 2) Tulis data baris-per-baris (mulai dari row 2)
	rowIdx := 2
	isPrincipalUser := isActivityReportPrincipalUser(dataFilter.AuthCustID, dataFilter.ParentCustID)
	for _, row := range activitySalesReports {
		mapped := mapActivityReportRow(row, isPrincipalUser)
		rec := []interface{}{
			mapped.BusinessUnitCode,
			mapped.BusinessUnitName,
			mapped.PJPCode,
			mapped.EmployeeCode,
			mapped.SalesmanName,
			mapped.DistributorCode,
			mapped.DistributorName,
			mapped.OutletCode,
			mapped.OutletPrincipalCode,
			mapped.OutletName,
			mapped.Date,
			mapped.ClockInTime,
			mapped.ClockOutTime,
			mapped.CheckinTime,
			mapped.CheckoutTime,
			mapped.Duration,
			mapped.PjpStatus,
			mapped.VisitStatus,
			mapped.Compliance,
			formatDownloadAmountValue(mapped.SalesValue),
			mapped.ReturnValue,
			formatDownloadAmountValue(mapped.PaymentCollected),
			mapped.LocationMaster,
			mapped.LocationActual,
			mapped.GeotagStatusDesc,
			mapped.Remarks,
		}

		cell, _ := excelize.CoordinatesToCellName(1, rowIdx) // kolom A, baris rowIdx
		if err := sw.SetRow(cell, rec); err != nil {
			return err
		}
		rowIdx++
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	// (opsional) atur lebar kolom biar rapi
	_ = f.SetColWidth(sheet, "A", "Z", 16)

	// 3) Serialize ke buffer
	buf, err := f.WriteToBuffer() // menghasilkan *bytes.Buffer
	if err != nil {
		return err
	}

	// 4) Penamaan file + upload ke storage seperti sebelumnya (ubah ekstensi dan content-type)
	var reportList model.ReportList

	fileName := buildReportObjectFileName(report.ReportName, dataFilter.ReportID)
	log.Info("Generated file name:", fileName)

	upload := &model.UploadCsv{
		FileName:    fileName,
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		FileData:    bytes.NewReader(buf.Bytes()),
	}

	downloadUrl, err := service.ObsAdapter.UploadFileCsv(upload)
	if err != nil {
		log.Error("SalesmanActivityReport, UploadFileCsv:", err.Error())
		return err
	}
	log.Info("Uploaded Excel URL:", downloadUrl)

	reportList.ReportID = dataFilter.ReportID
	reportList.FileURL = downloadUrl
	reportList.FileStatus = entity.FILE_STATUS_READY

	if err = service.ReportRepository.UpdateReportByReportID(context.Background(), dataFilter.ReportID, &reportList); err != nil {
		log.Error("SalesmanActivityReport, UpdateReportByReportID:", err.Error())
		return err
	}

	return err
}

func (service *reportServiceImpl) ExtractReportSecondary(req entity.SecondarySalesReportDashboardExtractQueryFilter) (err error) {
	custIDsOrder, err := service.ReportRepository.ListCustIDReportSecondarySalesReportOrder(req.Date)
	if err != nil {
		return err
	}

	for _, custIDOrder := range custIDsOrder {
		err = service.ExtractReportSalesOrder(custIDOrder.CustID, req)
		if err != nil {
			return err
		}
	}

	custIDsReturns, err := service.ReportRepository.ListCustIDReportSecondarySalesReportReturn(req.Date)
	if err != nil {
		return err
	}
	for _, custIDReturn := range custIDsReturns {
		err = service.ExtractReportSalesReturn(custIDReturn.CustID, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *reportServiceImpl) ExtractReportSalesOrder(custID string, req entity.SecondarySalesReportDashboardExtractQueryFilter) (err error) {
	c := context.Background()

	context.Background()

	offset := 0

	for {
		var batchSize = 1000

		orders, err := service.ReportRepository.GetReportSecondarySalesReportOrder(custID, req.Date, batchSize, offset)
		if err != nil {
			return fmt.Errorf("get orders batch %d: %w", offset/batchSize+1, err)
		}

		if len(orders) == 0 {
			break // tidak ada data lagi
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			var productCategoriesModel []model.DimProductCategory
			var productModel []model.DimProduct
			var outletModel []model.DimOutlet
			var salesmanModel []model.DimSalesman
			var orderModel []model.FactOrder

			uniqueProductCategories := make(map[int64]struct{}) // atau map[uint]struct{} tergantung tipe ID
			uniqueProduct := make(map[int64]struct{})           // atau map[uint]struct{} tergantung tipe ID
			uniqueOutlet := make(map[int64]struct{})            // atau map[uint]struct{} tergantung tipe ID
			uniqueSalesman := make(map[int64]struct{})          // atau map[uint]struct{} tergantung tipe ID

			// Kumpulkan semua tanggal unik
			dateSet := make(map[string]time.Time)
			for _, order := range orders {
				t, err := time.Parse(time.RFC3339, order.RoDate)
				if err != nil {
					return fmt.Errorf("invalid RoDate %s: %w", order.RoDate, err)
				}
				key := fmt.Sprintf("%d-%d-%d", t.Day(), int(t.Month()), t.Year())
				dateSet[key] = t
			}

			// Siapkan slice tanggal untuk batch insert dim_date
			dates := make([]time.Time, 0, len(dateSet))
			for _, t := range dateSet {
				dates = append(dates, t)
			}

			// Dapatkan semua ID tanggal (existing/new)
			dateMap, err := service.ReportRepository.GetOrCreateBatchDimDate(txCtx, dates)
			if err != nil {
				return err
			}

			for _, order := range orders {
				now := time.Now()
				// product categories
				if _, exists := uniqueProductCategories[order.PcatID]; !exists {
					uniqueProductCategories[order.PcatID] = struct{}{}
					productCategoriesModel = append(productCategoriesModel, model.DimProductCategory{
						ID:   order.PcatID,
						Code: order.PcatCode,
						Name: order.PcatName,
					})
				}

				// end product categories

				// product
				if _, exists := uniqueProduct[order.ProductID]; !exists {
					uniqueProduct[order.ProductID] = struct{}{}
					productModel = append(productModel, model.DimProduct{
						ID:         order.ProductID,
						CategoryID: order.PcatID,
						Code:       order.ProCode,
						Name:       order.ProName,
						UnitID1:    order.UnitID1,
						UnitID2:    order.UnitID2,
						UnitID3:    order.UnitID3,
						ConvUnit2:  order.ConvUnit2,
						ConvUnit3:  order.ConvUnit3,
					})

				}

				// outlet
				if _, exists := uniqueOutlet[order.OutletID]; !exists {
					uniqueOutlet[order.OutletID] = struct{}{}
					outletModel = append(outletModel, model.DimOutlet{
						ID:   order.OutletID,
						Code: order.OutletCode,
						Name: order.OutletName,
					})
				}

				// salesman
				if _, exists := uniqueSalesman[order.SalesmanID]; !exists {
					uniqueSalesman[order.SalesmanID] = struct{}{}
					salesmanModel = append(salesmanModel, model.DimSalesman{
						ID:   order.SalesmanID,
						Code: order.EmpCode,
						Name: order.EmpName,
					})
				}

				// order
				t, _ := time.Parse(time.RFC3339, order.RoDate)
				key := fmt.Sprintf("%d-%d-%d", t.Day(), int(t.Month()), t.Year())

				dateID, ok := dateMap[key]
				if !ok {
					return fmt.Errorf("DateID not found for RoDate: %s", order.RoDate)
				}

				orderModel = append(orderModel, model.FactOrder{
					CustID:          order.CustID,
					RoNo:            order.RoNo,
					InvoiceNo:       order.InvoiceNo,
					DateID:          dateID,
					SalesmanID:      order.SalesmanID,
					OutletID:        order.OutletID,
					ProID:           order.ProductID,
					Qty1:            order.Qty1,
					Qty2:            order.Qty2,
					Qty3:            order.Qty3,
					Qty:             order.Qty,
					ItemType:        order.ItemType,
					GrossSales:      order.GrossSales,
					SpecialDiscount: order.SpecialDiscount,
					Discount:        order.Discount,
					NetSalesExcPPN:  order.NetSalesExcPPN,
					PPN:             order.PPN,
					NetSalesIncPPN:  order.NetSalesIncPPN,
					SellPrice1:      order.SellPrice1,
					SellPrice2:      order.SellPrice2,
					SellPrice3:      order.SellPrice3,
					ExtractedAt:     now,
				})

			}
			// Save batch kategori
			if err := service.saveBatchProductCategories(txCtx, productCategoriesModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchProducts(txCtx, productModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchOutlets(txCtx, outletModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchSalesman(txCtx, salesmanModel); err != nil {
				return err
			}

			// Save batch order
			if err := service.saveBatchOrder(txCtx, orderModel); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("transaction batch %d: %w", offset/batchSize+1, err)
		}

		fmt.Println(fmt.Sprintf("[INFO] Batch %d processed (%d rows)", offset/batchSize+1, len(orders)))
		offset += batchSize
	}

	return nil
}

func (service *reportServiceImpl) ExtractReportSalesReturn(custID string, req entity.SecondarySalesReportDashboardExtractQueryFilter) (err error) {
	c := context.Background()

	context.Background()

	offset := 0

	for {
		var batchSize = 1000

		returns, err := service.ReportRepository.GetReportSecondarySalesReportReturn(custID, req.Date, batchSize, offset)
		if err != nil {
			return fmt.Errorf("get orders batch %d: %w", offset/batchSize+1, err)
		}

		if len(returns) == 0 {
			break // tidak ada data lagi
		}

		err = service.Transaction.WithinTransaction(c, func(txCtx context.Context) error {
			var productCategoriesModel []model.DimProductCategory
			var productModel []model.DimProduct
			var outletModel []model.DimOutlet
			var salesmanModel []model.DimSalesman
			var returnModel []model.FactReturn

			uniqueProductCategories := make(map[int64]struct{}) // atau map[uint]struct{} tergantung tipe ID
			uniqueProduct := make(map[int64]struct{})           // atau map[uint]struct{} tergantung tipe ID
			uniqueOutlet := make(map[int64]struct{})            // atau map[uint]struct{} tergantung tipe ID
			uniqueSalesman := make(map[int64]struct{})          // atau map[uint]struct{} tergantung tipe ID

			// Kumpulkan semua tanggal unik
			dateSet := make(map[string]time.Time)
			for _, rtn := range returns {
				t, err := time.Parse(time.RFC3339, rtn.ReturnDate)
				if err != nil {
					return fmt.Errorf("invalid RoDate %s: %w", rtn.ReturnDate, err)
				}
				key := fmt.Sprintf("%d-%d-%d", t.Day(), int(t.Month()), t.Year())
				dateSet[key] = t
			}

			// Siapkan slice tanggal untuk batch insert dim_date
			dates := make([]time.Time, 0, len(dateSet))
			for _, t := range dateSet {
				dates = append(dates, t)
			}

			// Dapatkan semua ID tanggal (existing/new)
			dateMap, err := service.ReportRepository.GetOrCreateBatchDimDate(txCtx, dates)
			if err != nil {
				return err
			}

			for _, order := range returns {
				now := time.Now()
				// product categories
				if _, exists := uniqueProductCategories[order.PcatID]; !exists {
					uniqueProductCategories[order.PcatID] = struct{}{}
					productCategoriesModel = append(productCategoriesModel, model.DimProductCategory{
						ID:   order.PcatID,
						Code: order.PcatCode,
						Name: order.PcatName,
					})
				}

				// end product categories

				// product
				if _, exists := uniqueProduct[order.ProductID]; !exists {
					uniqueProduct[order.ProductID] = struct{}{}
					productModel = append(productModel, model.DimProduct{
						ID:         order.ProductID,
						CategoryID: order.PcatID,
						Code:       order.ProCode,
						Name:       order.ProName,
						UnitID1:    order.UnitID1,
						UnitID2:    order.UnitID2,
						UnitID3:    order.UnitID3,
						ConvUnit2:  order.ConvUnit2,
						ConvUnit3:  order.ConvUnit3,
					})

				}

				// outlet
				if _, exists := uniqueOutlet[order.OutletID]; !exists {
					uniqueOutlet[order.OutletID] = struct{}{}
					outletModel = append(outletModel, model.DimOutlet{
						ID:   order.OutletID,
						Code: order.OutletCode,
						Name: order.OutletName,
					})
				}

				// salesman
				if _, exists := uniqueSalesman[order.SalesmanID]; !exists {
					uniqueSalesman[order.SalesmanID] = struct{}{}
					salesmanModel = append(salesmanModel, model.DimSalesman{
						ID:   order.SalesmanID,
						Code: order.EmpCode,
						Name: order.EmpName,
					})
				}

				// order
				t, _ := time.Parse(time.RFC3339, order.ReturnDate)
				key := fmt.Sprintf("%d-%d-%d", t.Day(), int(t.Month()), t.Year())

				dateID, ok := dateMap[key]
				if !ok {
					return fmt.Errorf("DateID not found for RoDate: %s", order.ReturnDate)
				}

				QtyUnit := &conversion.QtyUnit{
					Qty1:      int(order.Qty1),
					Qty2:      int(order.Qty2),
					Qty3:      int(order.Qty3),
					ConvUnit2: int(order.ConvUnit2),
					ConvUnit3: int(order.ConvUnit3),
				}

				totalQty, err := QtyUnit.ToTotalQuantity()
				if err != nil {
					return err
				}

				fmt.Println(fmt.Sprintf("cust id = %v", order.CustID))
				returnModel = append(returnModel, model.FactReturn{
					CustID:          order.CustID,
					ReturnNo:        order.ReturnNo,
					InvoiceNo:       order.InvoiceNo,
					DateID:          dateID,
					SalesmanID:      order.SalesmanID,
					OutletID:        order.OutletID,
					ProID:           order.ProductID,
					Qty1:            order.Qty1,
					Qty2:            order.Qty2,
					Qty3:            order.Qty3,
					Qty:             float64(totalQty),
					ItemType:        order.ItemType,
					GrossSales:      order.GrossSales,
					SpecialDiscount: order.SpecialDiscount,
					Discount:        order.Discount,
					NetSalesExcPPN:  order.NetSalesExcPPN,
					PPN:             order.PPN,
					NetSalesIncPPN:  order.NetSalesIncPPN,
					SellPrice1:      order.SellPrice1,
					SellPrice2:      order.SellPrice2,
					SellPrice3:      order.SellPrice3,
					ExtractedAt:     now,
				})

			}
			// Save batch kategori
			if err := service.saveBatchProductCategories(txCtx, productCategoriesModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchProducts(txCtx, productModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchOutlets(txCtx, outletModel); err != nil {
				return err
			}

			// Save batch produk
			if err := service.saveBatchSalesman(txCtx, salesmanModel); err != nil {
				return err
			}

			// Save batch order
			if err := service.saveBatchReturn(txCtx, returnModel); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("transaction batch %d: %w", offset/batchSize+1, err)
		}

		fmt.Println(fmt.Sprintf("[INFO] Batch %d processed (%d rows)", offset/batchSize+1, len(returns)))
		offset += batchSize
	}

	return nil
}

func (service *reportServiceImpl) saveBatchProductCategories(ctx context.Context, data []model.DimProductCategory) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveProductCategoriesDim(ctx, batch); err != nil {
			return fmt.Errorf("save product categories batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) saveBatchProducts(ctx context.Context, data []model.DimProduct) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveProductDim(ctx, batch); err != nil {
			return fmt.Errorf("save product batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) saveBatchOutlets(ctx context.Context, data []model.DimOutlet) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveOutletsDim(ctx, batch); err != nil {
			return fmt.Errorf("save outlet batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) saveBatchSalesman(ctx context.Context, data []model.DimSalesman) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveSalemanDim(ctx, batch); err != nil {
			return fmt.Errorf("save salesman batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) saveBatchOrder(ctx context.Context, data []model.FactOrder) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveOrderfact(ctx, batch); err != nil {
			return fmt.Errorf("save salesman batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) saveBatchReturn(ctx context.Context, data []model.FactReturn) error {
	var batchSize = 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]

		if err := service.ReportRepository.SaveReturnfact(ctx, batch); err != nil {
			return fmt.Errorf("save salesman batch %d-%d: %w", i, end, err)
		}
	}
	return nil
}

func (service *reportServiceImpl) resolveSecondaryDashboardCustIDs(authCustID, parentCustID string, requested []string) ([]string, error) {
	normalized, err := entity.NormalizeStringList(requested)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return []string{authCustID}, nil
	}

	if authCustID != parentCustID {
		if len(normalized) != 1 || normalized[0] != authCustID {
			return nil, ErrUnauthorizedCustID
		}
		return normalized, nil
	}

	for _, custID := range normalized {
		if custID == authCustID {
			continue
		}
		allowed, err := service.ReportRepository.ExistsCustomerInParentScope(custID, parentCustID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, ErrUnauthorizedCustID
		}
	}

	return normalized, nil
}

func (service *reportServiceImpl) resolveSecondaryDashboardCustID(authCustID, parentCustID, requestedCustID string) (string, error) {
	custIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, []string{requestedCustID})
	if err != nil {
		return "", err
	}
	if len(custIDs) == 0 {
		return authCustID, nil
	}
	return custIDs[0], nil
}

func (service *reportServiceImpl) SecondarySalesReportSumReportByMonth(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (data entity.SumReportByMonthModelResp, err error) {
	requestedCustIDs := req.CustIDs
	if len(requestedCustIDs) == 0 && req.CustID != "" {
		requestedCustIDs = []string{req.CustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return
	}

	effectiveYear := resolveSecondaryDashboardYear(req.Year)

	sumReportModel, err := service.ReportRepository.SecondarySalesReportSumReportByMonth(effectiveCustIDs, req, effectiveYear)
	if err != nil {
		return
	}

	sumReportReturnModel, err := service.ReportRepository.SecondarySalesReportReturnSumReportByMonth(effectiveCustIDs, req.Month, effectiveYear)
	if err != nil {
		return
	}

	lastUpdate := sumReportModel.LastUpdate
	if lastUpdate == nil && sumReportReturnModel.LastUpdate != nil {
		lastUpdate = sumReportReturnModel.LastUpdate
	}

	data.NetSales = sumReportModel.NetSales
	data.TotalPPN = sumReportModel.TotalPPN
	data.NetSalesExcPPN = sumReportModel.NetSalesExcPPN
	data.Qty = sumReportModel.Qty
	data.TotalDiscountPromo = sumReportModel.TotalDiscountPromo
	data.TotalGrossSales = sumReportModel.TotalGrossSales
	data.TotalOutlet = sumReportModel.TotalOutlet
	data.TotalProduct = sumReportModel.TotalProduct
	data.TotalSalesman = sumReportModel.TotalSalesman
	data.QtyReturn = sumReportModel.QtyReturn
	data.ReturnRate = sumReportModel.ReturnRate
	data.NetSalesReturn = sumReportModel.NetSalesReturn
	data.LastUpdate = lastUpdate
	return
}

func (service *reportServiceImpl) SecondarySalesReportGroupSales(
	authCustID, parentCustID string,
	req entity.SecondarySalesReportDashboardGroupPayload,
) (datas []entity.SecondarySalesReportGroupResp, err error) {
	var results []model.SecondarySalesReportGroup

	requestedCustIDs := req.CustIDs
	if len(requestedCustIDs) == 0 && req.CustID != "" {
		requestedCustIDs = []string{req.CustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return nil, err
	}

	effectiveYear := resolveSecondaryDashboardYear(req.Year)

	switch req.GroupBy {
	case entity.SECONDARY_SALES_GROUP_OUTLET:
		results, err = service.ReportRepository.SecondarySalesReportGroupOutlet(effectiveCustIDs, req.Month, effectiveYear)
	case entity.SECONDARY_SALES_GROUP_SALESMAN:
		results, err = service.ReportRepository.SecondarySalesReportGroupSalesman(effectiveCustIDs, req.Month, effectiveYear)
	case entity.SECONDARY_SALES_GROUP_PRODUCT_CATEGORY:
		results, err = service.ReportRepository.SecondarySalesReportProductCategory(effectiveCustIDs, req.Month, effectiveYear)
	default:
		results, err = service.ReportRepository.SecondarySalesReportProduct(effectiveCustIDs, req.Month, effectiveYear)
	}

	if err != nil {
		return nil, err
	}

	datas = make([]entity.SecondarySalesReportGroupResp, len(results))
	for i, r := range results {
		datas[i] = entity.SecondarySalesReportGroupResp{
			ID:       r.ID,
			Code:     r.Code,
			Name:     r.Name,
			NetSales: r.NetSales,
		}
	}

	return datas, nil
}

func (service *reportServiceImpl) SecondarySalesReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.SumReportTrendSalesResp, err error) {
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return
	}

	trendSalesModels, err := service.ReportRepository.SecondarySalesReportTrendSales(effectiveCustIDs, year)
	if err != nil {
		return
	}

	for _, trendSalesModel := range trendSalesModels {
		data = append(data, entity.SumReportTrendSalesResp{
			Month:              trendSalesModel.Month,
			TotalGrossSales:    trendSalesModel.TotalGrossSales,
			TotalDiscountPromo: trendSalesModel.TotalDiscountPromo,
			NetSales:           trendSalesModel.NetSales,
		})
	}

	return
}

func (service *reportServiceImpl) SalesmanActivityReportSumReportByMonth(authCustID, parentCustID string, req entity.SalesmanActivityReportDashboardSumPayload) (data entity.SalesmanActivityReportByMonthModelResp, err error) {
	requestedCustIDs := req.CustIDs
	if len(requestedCustIDs) == 0 && req.CustID != "" {
		requestedCustIDs = []string{req.CustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return
	}

	effectiveYear := resolveSecondaryDashboardYear(req.Year)

	sumReportModel, err := service.ReportRepository.SalesmanActivityReportSumByMonth(effectiveCustIDs, req.Month, effectiveYear)
	if err != nil {
		return
	}

	data.TotalSales = sumReportModel.TotalSales
	data.TotalReturn = sumReportModel.TotalReturn
	data.SalesmanTotal = sumReportModel.TotalSalesman
	data.LastUpdate = sumReportModel.LastUpdate

	return
}

func (service *reportServiceImpl) SalesmanActivityReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.ActivityReportTrendSalesResp, err error) {
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return
	}

	trendSalesModels, err := service.ReportRepository.SalesmanActivityReportTrendSales(effectiveCustIDs, year)
	if err != nil {
		return
	}

	for _, trendSalesModel := range trendSalesModels {
		data = append(data, entity.ActivityReportTrendSalesResp{
			Month:        trendSalesModel.Month,
			TotalInvoice: trendSalesModel.TotalInvoice,
			TotalReturn:  trendSalesModel.TotalReturn,
			NetSales:     trendSalesModel.NetSales,
		})
	}

	return
}

func (service *reportServiceImpl) SalesmanActivityReportGroupSales(
	authCustID, parentCustID string,
	req entity.SalesmanActivityReportDashboardGroupPayload,
) (datas []entity.SecondarySalesReportGroupResp, err error) {
	requestedCustIDs := req.CustIDs
	if len(requestedCustIDs) == 0 && req.CustID != "" {
		requestedCustIDs = []string{req.CustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return nil, err
	}

	effectiveYear := resolveSecondaryDashboardYear(req.Year)

	var results []model.SecondarySalesReportGroup
	var resultsReturn []model.ReturnReportGroup

	if req.ActivityType == entity.ACTIVITY_SALESMAN_GROUP_SALES {
		results, err = service.ReportRepository.ActivitySalesmanReportGroupSalesman(effectiveCustIDs, req.Month, effectiveYear)
		if err != nil {
			return nil, err
		}
		datas = make([]entity.SecondarySalesReportGroupResp, len(results))
		for i, r := range results {
			datas[i] = entity.SecondarySalesReportGroupResp{
				ID:       r.ID,
				Code:     r.Code,
				Name:     r.Name,
				NetSales: r.NetSales,
			}
		}
		return datas, nil
	}

	resultsReturn, err = service.ReportRepository.ActivitySalesmanReturnReportGroupSalesman(effectiveCustIDs, req.Month, effectiveYear)
	if err != nil {
		return nil, err
	}

	datas = make([]entity.SecondarySalesReportGroupResp, len(resultsReturn))
	for i, r := range resultsReturn {
		datas[i] = entity.SecondarySalesReportGroupResp{
			ID:       r.ID,
			Code:     r.Code,
			Name:     r.Name,
			NetSales: r.NetSales,
		}
	}

	return datas, nil
}

func (service *reportServiceImpl) SalesmanActivityReportGeotag(
	authCustID, parentCustID string,
	req entity.ActivityReportGeotagPayload,
) (data entity.ActivityReportGeotagResp, err error) {
	requestedCustIDs := req.CustIDs
	if len(requestedCustIDs) == 0 && req.CustID != "" {
		requestedCustIDs = []string{req.CustID}
	}
	effectiveCustIDs, err := service.resolveSecondaryDashboardCustIDs(authCustID, parentCustID, requestedCustIDs)
	if err != nil {
		return data, err
	}

	rows, err := service.ReportRepository.ActivityReportGeotag(parentCustID, effectiveCustIDs, req.Year, req.EmpID)
	if err != nil {
		return data, err
	}

	var totalVisit, totalMatch, totalUnmatch int64
	details := make([]entity.ActivityReportGeotagDetailResp, 0, len(rows))
	for _, row := range rows {
		totalVisit += row.TotalVisit
		totalMatch += row.GeotagMatchCount
		totalUnmatch += row.GeotagUnmatchCount
		details = append(details, entity.ActivityReportGeotagDetailResp{
			SalesmanCode:            strconv.FormatInt(row.SalesmanCode, 10),
			SalesmanName:            row.SalesmanName,
			TotalVisit:              row.TotalVisit,
			GeotagMatchCount:        row.GeotagMatchCount,
			GeotagUnmatchCount:      row.GeotagUnmatchCount,
			GeotagMatchPercentage:   row.GeotagMatchPct,
			GeotagUnmatchPercentage: row.GeotagUnmatchPct,
		})
	}

	data.Details = details
	if totalVisit > 0 {
		data.TotalGeotagMatchPercentage = roundActivityReportGeotagServicePct(float64(totalMatch), float64(totalVisit))
		data.TotalGeotagUnmatchPercentage = roundActivityReportGeotagServicePct(float64(totalUnmatch), float64(totalVisit))
	}

	return data, nil
}

func roundActivityReportGeotagServicePct(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round(100.0*part/total*100) / 100
}

func (service *reportServiceImpl) SalesmanActivitySalesmanList(dataFilter entity.ActivityReportSalesmanListQueryFilter) (datas []entity.SalesmanActivityReportSalesmanListResp, err error) {
	salesmans, err := service.ReportRepository.ActivitySalesReportSalesmanList(dataFilter)
	if err != nil {
		return
	}

	for _, salesman := range salesmans {
		datas = append(datas, entity.SalesmanActivityReportSalesmanListResp{
			SalesmanID:   salesman.SalesmanID,
			SalesmanCode: salesman.SalesmanCode,
			SalesmanName: salesman.SalesmanName,
		})
	}
	return
}
