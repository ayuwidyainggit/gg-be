package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"master/entity"
	"master/model"
	"master/pkg/constant"
	"master/pkg/rabbitmq"
	"master/pkg/str"
	"master/pkg/structs"
	"master/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	mPriceDateLayout = "2006-01-02"
)

var errManagePricingNotAllowed = errors.New("you are not allowed to manage pricing")

type MPriceService interface {
	Detail(entity.DetailMPriceParams) (entity.MPriceResponse, error)
	List(entity.MPriceQueryFilter, string) (data []entity.MPriceResponse, total int, lastPage int, err error)
	Store(entity.CreateMPriceBody) (entity.MPriceResponse, error)
	Update(string, entity.UpdateMPriceRequest) error
	Publish(entity.PublishMPriceParams) error
	PublishByRMQ(entity.PublishByRmqMPriceReq) error
	Cancel(entity.CancelMPriceParams) error
	Delete(string, string, int64) error
	Template(string, string, string, int64) (*bytes.Buffer, string, string, error)
	Export(entity.MPriceQueryFilter, string, string) (*bytes.Buffer, string, string, error)
	Import(entity.MPriceImportRequest, string, string, int64, int64, string) (entity.MPriceImportResponse, error)
}

func NewMPriceService(
	priceRepository repository.MPriceRepository,
	distributorRepository repository.DistributorRepository,
	transPriceRepository repository.MTransactionPriceRepository,
) *MPriceServiceImpl {
	return &MPriceServiceImpl{
		MPriceRepository:            priceRepository,
		DistributorRepository:       distributorRepository,
		MTransactionPriceRepository: transPriceRepository,
	}
}

type MPriceServiceImpl struct {
	MPriceRepository            repository.MPriceRepository
	DistributorRepository       repository.DistributorRepository
	MTransactionPriceRepository repository.MTransactionPriceRepository
}

func (service *MPriceServiceImpl) Detail(detail entity.DetailMPriceParams) (response entity.MPriceResponse, err error) {
	price, err := service.MPriceRepository.FindOneByMPriceIDAndCustID(detail)
	if err != nil {
		return response, err
	}

	if err = structs.Automapper(price, &response); err != nil {
		return response, err
	}
	if price.EffectiveDate != nil {
		response.EffectiveDate = price.EffectiveDate.Format(mPriceDateLayout)
	}
	response.StatusDesc = response.GetPriceStatusDesc()

	response.Details = make([]entity.DistributorAreaRegionData, 0)
	distributorIDs, err := service.resolveDetailDistributorIDs(price, detail.ParentCustID)
	if err != nil {
		return response, err
	}
	if len(distributorIDs) > 0 {
		updatedDistributorIDs, err := service.resolveUpdatedDistributorIDs(price, detail.ParentCustID, distributorIDs)
		if err != nil {
			return response, err
		}

		distributors, err := service.DistributorRepository.FindAllAreaRegionByIDs(distributorIDs, detail.ParentCustID)
		if err != nil {
			return response, err
		}

		for _, row := range distributors {
			var distData entity.DistributorAreaRegionData
			if err = structs.Automapper(row, &distData); err != nil {
				return response, err
			}

			distData.IsUpdated = updatedDistributorIDs[row.DistributorID]
			distData.UpdateStatus = managePriceProductUpdateStatus(distData.IsUpdated)
			response.Details = append(response.Details, distData)
		}
	}

	return response, nil
}

func (service *MPriceServiceImpl) List(dataFilter entity.MPriceQueryFilter, custID string) (data []entity.MPriceResponse, total int, lastPage int, err error) {
	prices, total, lastPage, err := service.MPriceRepository.FindAllByCustID(dataFilter, custID)
	if err != nil {
		return data, total, lastPage, err
	}

	for _, row := range prices {
		var vResp entity.MPriceResponse
		if err = structs.Automapper(row, &vResp); err != nil {
			return data, total, lastPage, err
		}
		if row.EffectiveDate != nil {
			vResp.EffectiveDate = row.EffectiveDate.Format(mPriceDateLayout)
		}
		vResp.StatusDesc = vResp.GetPriceStatusDesc()
		data = append(data, vResp)
	}

	return data, total, lastPage, nil
}

func (service *MPriceServiceImpl) Store(request entity.CreateMPriceBody) (response entity.MPriceResponse, err error) {
	if err = service.validateManagePricingPermission(request.ParentCustID, request.DistributorID); err != nil {
		return response, err
	}

	if err = service.prepareCreateRequest(&request); err != nil {
		return response, err
	}

	effectiveDate, scheduleAt, isImmediate, err := parseManagePriceEffectiveDate(request.EffectiveDate)
	if err != nil {
		return response, err
	}

	request.EffectiveDate, err = str.DateStrToRfc3339String(request.EffectiveDate)
	if err != nil {
		return response, err
	}

	objectID := primitive.NewObjectID().Hex()
	nowUTC := time.Now().UTC()

	priceData := model.MPrice{}
	if err = structs.Automapper(request, &priceData); err != nil {
		return response, err
	}
	priceData.CustID = request.CustID
	priceData.PriceID = objectID
	priceData.CreatedAt = nowUTC
	priceData.CreatedByID = request.CreatedByID
	priceData.CreatedBy = request.CreatedBy
	priceData.UpdatedAt = nowUTC
	priceData.UpdatedByID = request.CreatedByID
	priceData.UpdatedBy = request.CreatedBy
	priceData.Status = 1

	if err = service.MPriceRepository.Store(&priceData); err != nil {
		return response, err
	}

	if err = structs.Automapper(priceData, &response); err != nil {
		return response, err
	}
	response.EffectiveDate = effectiveDate.Format(mPriceDateLayout)
	response.Status = 1

	payload := entity.PublishByRmqMPriceReq{
		CustID:        request.CustID,
		ParentCustID:  request.ParentCustID,
		DistributorID: request.DistributorID,
		PriceID:       response.PriceID,
		Status:        10,
		UpdatedBy:     request.CreatedBy,
		UpdatedByID:   request.CreatedByID,
	}

	if isImmediate {
		if err = service.PublishByRMQ(payload); err != nil {
			return response, err
		}
		response.Status = 10
		response.StatusDesc = response.GetPriceStatusDesc()
		return response, nil
	}

	service.enqueuePublish(payload, scheduleAt, request.ExpirationMs)
	response.StatusDesc = response.GetPriceStatusDesc()
	return response, nil
}

func (service *MPriceServiceImpl) Update(priceID string, request entity.UpdateMPriceRequest) error {
	detail := entity.DetailMPriceParams{
		CustID:       request.CustID,
		ParentCustID: request.ParentCustID,
		PriceID:      priceID,
	}
	price, err := service.MPriceRepository.FindOneByMPriceIDAndCustID(detail)
	if err != nil {
		return err
	}
	if price.Status != 1 {
		return errors.New("manage price status not allowed to edit")
	}

	if err = service.prepareUpdateRequest(&request); err != nil {
		return err
	}

	_, scheduleAt, isImmediate, err := parseManagePriceEffectiveDate(*request.EffectiveDate)
	if err != nil {
		return err
	}

	if effectiveDateRFC3339, convErr := str.DateStrToRfc3339String(*request.EffectiveDate); convErr == nil {
		request.EffectiveDate = &effectiveDateRFC3339
	} else {
		return convErr
	}

	if err = service.MPriceRepository.Update(priceID, request); err != nil {
		return err
	}

	payload := entity.PublishByRmqMPriceReq{
		CustID:        request.CustID,
		ParentCustID:  request.ParentCustID,
		DistributorID: request.DistributorID,
		PriceID:       priceID,
		Status:        10,
		UpdatedBy:     request.UpdatedBy,
		UpdatedByID:   request.UpdatedByID,
	}
	if isImmediate {
		return service.PublishByRMQ(payload)
	}

	service.enqueuePublish(payload, scheduleAt, 0)
	return nil
}

func (service *MPriceServiceImpl) Delete(custID string, priceID string, userID int64) error {
	return service.MPriceRepository.Delete(custID, priceID, userID)
}

func (service *MPriceServiceImpl) Cancel(detail entity.CancelMPriceParams) error {
	return service.MPriceRepository.Cancel(detail)
}

func (service *MPriceServiceImpl) Publish(params entity.PublishMPriceParams) error {
	price, err := service.MPriceRepository.FindOneByMPriceIDAndCustID(entity.DetailMPriceParams{
		CustID:       params.CustID,
		ParentCustID: params.ParentCustID,
		PriceID:      params.PriceID,
	})
	if err != nil {
		return err
	}
	if price.Status != 1 {
		return fmt.Errorf("only scheduled manage prices can be published; current status is %s", entity.MPriceStatusDesc[price.Status])
	}

	_, scheduleAt, _, err := parseManagePriceEffectiveDate(price.EffectiveDate.Format(mPriceDateLayout))
	if err != nil {
		return err
	}
	if time.Now().In(scheduleAt.Location()).Before(scheduleAt) {
		return errors.New("manage price effective date is not due yet")
	}

	return service.PublishByRMQ(entity.PublishByRmqMPriceReq{
		CustID:        params.CustID,
		ParentCustID:  params.ParentCustID,
		DistributorID: params.DistributorID,
		PriceID:       params.PriceID,
		Status:        10,
		UpdatedBy:     params.UpdatedBy,
		UpdatedByID:   params.UpdatedByID,
	})
}

func (service *MPriceServiceImpl) PublishByRMQ(request entity.PublishByRmqMPriceReq) error {
	detail := entity.DetailMPriceParams{
		CustID:       request.CustID,
		ParentCustID: request.ParentCustID,
		PriceID:      request.PriceID,
	}
	price, err := service.MPriceRepository.FindOneByMPriceIDAndCustID(detail)
	if err != nil {
		return err
	}
	if price.Status != 1 {
		return nil
	}

	_, scheduleAt, _, err := parseManagePriceEffectiveDate(price.EffectiveDate.Format(mPriceDateLayout))
	if err != nil {
		return err
	}
	nowJkt := time.Now().In(scheduleAt.Location())
	if nowJkt.Before(scheduleAt) {
		service.enqueuePublish(request, scheduleAt, 0)
		return nil
	}

	if err = service.applyPublishedProductPrices(request, price); err != nil {
		return err
	}
	if err = service.syncTransactionPrices(price); err != nil {
		return err
	}

	log.Info("Service, PublishByRMQ, request -> ", structs.StructToJson(request))
	if err = service.MPriceRepository.PublishByRMQ(request); err != nil {
		if err.Error() == "no rows affected" {
			return nil
		}
		return err
	}

	return nil
}

func (service *MPriceServiceImpl) Template(format, custID, parentCustID string, distributorID int64) (*bytes.Buffer, string, string, error) {
	principalRows := buildManagePricePrincipalTemplateRows()
	distributorRows := buildManagePriceDistributorTemplateRows()
	templateRows := principalRows
	templateSheetName := "Template sebagai Principal"
	if isDistributorManagePriceTemplateScope(distributorID) {
		templateRows = distributorRows
		templateSheetName = "Template sebagai Distributor"
	}

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "", "xlsx":
		return exportManagePriceTemplateWorkbook(
			[]managePriceTemplateSheet{
				{Name: templateSheetName, Rows: templateRows},
			},
			"manage_price_template.xlsx",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		)
	case "xls":
		return exportManagePriceTemplateWorkbook(
			[]managePriceTemplateSheet{
				{Name: templateSheetName, Rows: templateRows},
			},
			"manage_price_template.xls",
			"application/vnd.ms-excel",
		)
	case "csv":
		if isDistributorManagePriceTemplateScope(distributorID) {
			return exportManagePriceTemplateCSV(distributorRows, "manage_price_template_distributor.csv")
		}
		return exportManagePriceTemplateCSV(principalRows, "manage_price_template_principal.csv")
	default:
		return nil, "", "", errors.New("invalid format")
	}
}

func (service *MPriceServiceImpl) Export(filter entity.MPriceQueryFilter, custID, parentCustID string) (*bytes.Buffer, string, string, error) {
	if filter.Limit == 0 {
		filter.Limit = 1000
	}

	rows, _, _, err := service.List(filter, custID)
	if err != nil {
		return nil, "", "", err
	}
	for i := range rows {
		detail, err := service.Detail(entity.DetailMPriceParams{
			CustID:       custID,
			ParentCustID: parentCustID,
			PriceID:      rows[i].PriceID,
		})
		if err != nil {
			return nil, "", "", err
		}
		rows[i].Details = detail.Details
	}

	switch strings.ToLower(filter.FileType) {
	case "csv":
		return exportManagePriceCSV(rows)
	case "xls":
		buf, _, _, err := exportManagePriceExcel(rows, "manage_price_export.xls")
		return buf, "application/vnd.ms-excel", "manage_price_export.xls", err
	default:
		return exportManagePriceExcel(rows, "manage_price_export.xlsx")
	}
}

func (service *MPriceServiceImpl) Import(req entity.MPriceImportRequest, custID, parentCustID string, userID int64, distributorID int64, userFullName string) (entity.MPriceImportResponse, error) {
	resp := entity.MPriceImportResponse{
		FileURL:     req.FileURL,
		ProcessedAt: time.Now().Format(time.RFC3339),
		FailedRows:  make([]string, 0),
	}

	if err := service.validateManagePricingPermission(parentCustID, distributorID); err != nil {
		return resp, err
	}

	rows, err := downloadManagePriceImportRows(req.FileURL)
	if err != nil {
		return resp, err
	}
	if len(rows) < 2 {
		return resp, errors.New("template does not contain data rows")
	}

	headerMap, dataStartRow := buildMPriceHeaderIndex(rows)
	requiredHeaders := []string{
		"effective_date",
		"pro_code",
		"new_purch_price1",
		"new_purch_price2",
		"new_purch_price3",
		"new_sell_price1",
		"new_sell_price2",
		"new_sell_price3",
	}
	for _, header := range requiredHeaders {
		if _, ok := headerMap[header]; !ok {
			return resp, fmt.Errorf("missing required header %s", header)
		}
	}

	for rowIndex, row := range rows[dataStartRow:] {
		if isEmptyManagePriceImportRow(row) {
			continue
		}
		actualRow := rowIndex + dataStartRow + 1
		resp.TotalRow++
		createReq, rowErr := service.buildCreateRequestFromImportRow(row, headerMap, custID, parentCustID, userID, distributorID, userFullName)
		if rowErr != nil {
			failedRow := fmt.Sprintf("row %d: %v", actualRow, rowErr)
			log.Info(fmt.Sprintf("MPriceService, Import, %s", failedRow))
			resp.FailedRows = append(resp.FailedRows, failedRow)
			resp.FailedRow++
			continue
		}
		if _, rowErr = service.Store(createReq); rowErr != nil {
			failedRow := fmt.Sprintf("row %d: %v", actualRow, rowErr)
			log.Info(fmt.Sprintf("MPriceService, Import, %s", failedRow))
			resp.FailedRows = append(resp.FailedRows, failedRow)
			resp.FailedRow++
			continue
		}
		resp.SuccessRow++
	}

	if resp.SuccessRow == 0 {
		return resp, errors.New("price data upload failed")
	}
	if resp.FailedRow > 0 {
		return resp, errors.New("price data upload processed with partial success")
	}

	return resp, nil
}

func (service *MPriceServiceImpl) validateManagePricingPermission(parentCustID string, distributorID int64) error {
	if distributorID <= 0 {
		return nil
	}
	if service.DistributorRepository == nil {
		return errManagePricingNotAllowed
	}

	allowed, err := service.DistributorRepository.FindAllowManagePricingByDistributorID(distributorID)
	if err != nil {
		log.Errorf("MPriceService, validateManagePricingPermission, distributor_id=%d parent_cust_id=%s, err=%v", distributorID, parentCustID, err)
		return errManagePricingNotAllowed
	}
	if !allowed {
		return errManagePricingNotAllowed
	}

	return nil
}

func (service *MPriceServiceImpl) prepareCreateRequest(request *entity.CreateMPriceBody) error {
	request.DistributorIDs = normalizeDistributorIDs(request.Coverage, request.DistributorIDs, request.DistributorID)
	if err := service.validateDistributors(request.Coverage, request.DistributorIDs, request.ParentCustID); err != nil {
		return err
	}

	snapshot, err := service.MPriceRepository.FindOneProductSnapshotByProID(request.ProID, request.CustID)
	if err != nil {
		log.Errorf("MPriceService, prepareCreateRequest, FindOneProductSnapshotByProID, pro_id=%d cust_id=%s, err=%v", request.ProID, request.CustID, err)
		return errors.New("pro_id not found")
	}
	applySnapshotToCreateRequest(request, snapshot)
	return nil
}

func (service *MPriceServiceImpl) prepareUpdateRequest(request *entity.UpdateMPriceRequest) error {
	coverage := ""
	if request.Coverage != nil {
		coverage = *request.Coverage
	}
	distIDs := make([]int64, 0)
	if request.DistributorIDs != nil {
		distIDs = *request.DistributorIDs
	}
	distIDs = normalizeDistributorIDs(coverage, distIDs, request.DistributorID)
	if err := service.validateDistributors(coverage, distIDs, request.ParentCustID); err != nil {
		return err
	}
	if coverage == "N" {
		request.DistributorIDs = nil
	} else {
		request.DistributorIDs = &distIDs
	}

	if request.ProID == nil {
		return errors.New("pro_id not found")
	}
	snapshot, err := service.MPriceRepository.FindOneProductSnapshotByProID(*request.ProID, request.CustID)
	if err != nil {
		log.Errorf("MPriceService, prepareUpdateRequest, FindOneProductSnapshotByProID, pro_id=%d cust_id=%s, err=%v", *request.ProID, request.CustID, err)
		return errors.New("pro_id not found")
	}
	applySnapshotToUpdateRequest(request, snapshot)
	return nil
}

func (service *MPriceServiceImpl) validateDistributors(coverage string, distributorIDs []int64, parentCustID string) error {
	if coverage == "N" {
		return nil
	}
	if len(distributorIDs) == 0 {
		return errors.New("distributor_ids is required")
	}
	distributors, err := service.DistributorRepository.FindAllAreaRegionByIDs(distributorIDs, parentCustID)
	if err != nil {
		log.Errorf("MPriceService, validateDistributors, FindAllAreaRegionByIDs, distributor_ids=%v parent_cust_id=%s, err=%v", distributorIDs, parentCustID, err)
		return errors.New("distributor_ids not found")
	}
	if len(distributorIDs) != len(distributors) {
		return errors.New("distributor_ids not found")
	}
	return nil
}

func (service *MPriceServiceImpl) applyPublishedProductPrices(request entity.PublishByRmqMPriceReq, price model.MPriceDetail) error {
	if request.CustID == request.ParentCustID {
		distributorIDs := []int64(price.DistributorIDs)
		brokenIDs, err := service.MPriceRepository.FindBrokenDistributorChildLinks(price.ProID, request.ParentCustID, distributorIDs)
		if err != nil {
			return err
		}
		if len(brokenIDs) > 0 {
			log.Warnf("applyPublishedProductPrices: broken distributor child links detected price_id=%s parent_pro_id=%d broken_distributor_ids=%v",
				price.PriceID, price.ProID, brokenIDs)
		}
		return service.MPriceRepository.UpdatePrincipalAssignedProductPrices(price.ProID, distributorIDs, price)
	}

	targetDistributorID := request.DistributorID
	if targetDistributorID == 0 && len(price.DistributorIDs) == 1 {
		targetDistributorID = price.DistributorIDs[0]
	}
	return service.MPriceRepository.UpdateDistributorProductPrices(request.CustID, targetDistributorID, price.ProID, price)
}

func (service *MPriceServiceImpl) syncTransactionPrices(price model.MPriceDetail) error {
	if price.Coverage != "N" {
		var transPricesInsert, transPricesUpdate []model.MTransactionPrice
		for i := range price.DistributorIDs {
			transPriceExist, err := service.MTransactionPriceRepository.GetByCustPro(price.CustID, price.ProID, price.DistributorIDs[i], "D")
			if err != nil {
				transPrice := setupTransactionPriceData(price, price.DistributorIDs[i])
				transPricesInsert = append(transPricesInsert, transPrice)
				continue
			}
			transPrice := setupTransactionPriceDataUpdate(price, price.DistributorIDs[i], transPriceExist.TransactionPriceID)
			transPricesUpdate = append(transPricesUpdate, transPrice)
		}

		if len(transPricesInsert) > 0 {
			if err := service.MTransactionPriceRepository.StoreBatch(transPricesInsert); err != nil {
				return err
			}
		}
		if len(transPricesUpdate) > 0 {
			if err := service.MTransactionPriceRepository.UpdateBatch(transPricesUpdate); err != nil {
				return err
			}
		}
		return nil
	}

	service.MTransactionPriceRepository.DeleteByStartDate(price.CustID, price.ProID, *price.EffectiveDate)
	transPrice := setupTransactionPriceData(price, 0)
	return service.MTransactionPriceRepository.Store(&transPrice)
}

func (service *MPriceServiceImpl) resolveDetailDistributorIDs(price model.MPriceDetail, parentCustID string) ([]int64, error) {
	if len(price.DistributorIDs) > 0 {
		return []int64(price.DistributorIDs), nil
	}
	return service.MPriceRepository.FindAffectedDistributorProductIDs(price, parentCustID)
}

func (service *MPriceServiceImpl) resolveUpdatedDistributorIDs(price model.MPriceDetail, parentCustID string, distributorIDs []int64) (map[int64]bool, error) {
	updated := make(map[int64]bool, len(distributorIDs))
	if price.Status != 10 {
		return updated, nil
	}

	updatedDistributorIDs, err := service.MPriceRepository.FindUpdatedDistributorProductIDs(price, parentCustID, distributorIDs)
	if err != nil {
		return updated, err
	}
	for _, distributorID := range updatedDistributorIDs {
		updated[distributorID] = true
	}
	return updated, nil
}

func managePriceProductUpdateStatus(isUpdated bool) string {
	if isUpdated {
		return "updated"
	}
	return "not_updated"
}

func (service *MPriceServiceImpl) enqueuePublish(payload entity.PublishByRmqMPriceReq, scheduleAt time.Time, expirationMs int) {
	delta := scheduleAt.Sub(time.Now().In(scheduleAt.Location())).Milliseconds()
	if delta < 0 {
		delta = 0
	}
	if expirationMs > 0 {
		delta = int64(expirationMs)
	}

	rmqConfig := rabbitmq.RmqConfig{
		MessageID:      primitive.NewObjectID().Hex(),
		ExchangeName:   constant.RMQ_DEFAULT_EXCHANGE,
		RoutingKey:     constant.RMQ_MANAGE_PRICE_CREATE_EVENT,
		QueueName:      constant.RMQ_MANAGE_PRICE_CREATE_EVENT,
		DelayQueueName: constant.RMQ_MANAGE_PRICE_CREATE_EVENT + constant.RMQ_DEFAULT_DELAY_SUFFIX,
		MessageTTL:     strconv.FormatInt(delta, 10),
		Message:        structs.StructToJson(payload),
	}
	go func(cfg rabbitmq.RmqConfig) {
		if err := rabbitmq.PublishMessage(&cfg); err != nil {
			log.Errorf("MPriceService, enqueuePublish, failed to publish manage price event: %v", err)
		}
	}(rmqConfig)
}

func (service *MPriceServiceImpl) buildCreateRequestFromImportRow(row []string, headerMap map[string]int, custID, parentCustID string, userID int64, distributorID int64, userFullName string) (entity.CreateMPriceBody, error) {
	req := entity.CreateMPriceBody{
		CustID:        custID,
		ParentCustID:  parentCustID,
		CreatedBy:     userFullName,
		CreatedByID:   &userID,
		DistributorID: distributorID,
		EffectiveDate: getCell(row, headerMap, "effective_date"),
	}
	var err error
	req.EffectiveDate, err = normalizeManagePriceImportDate(req.EffectiveDate)
	if err != nil {
		return req, err
	}

	proCode := getCell(row, headerMap, "pro_code")
	if proCode == "" {
		return req, errors.New("pro_code is required")
	}
	snapshot, err := service.MPriceRepository.FindOneProductSnapshotByCode(proCode, custID)
	if err != nil {
		log.Errorf("MPriceService, buildCreateRequestFromImportRow, FindOneProductSnapshotByCode, pro_code=%s cust_id=%s, err=%v", proCode, custID, err)
		return req, fmt.Errorf("product %s not found", proCode)
	}
	req.ProID = snapshot.ProID

	distributorCodesRaw := getCell(row, headerMap, "distributor_code")
	req.Coverage = defaultManagePriceImportCoverage(custID, parentCustID, distributorID)
	if distributorCodesRaw != "" {
		req.Coverage = "D"
		distCodes := splitCommaSeparated(distributorCodesRaw)
		distributors, err := service.DistributorRepository.FindAllAreaRegionByCodes(distCodes, parentCustID)
		if err != nil {
			log.Errorf("MPriceService, buildCreateRequestFromImportRow, FindAllAreaRegionByCodes, distributor_codes=%v parent_cust_id=%s, err=%v", distCodes, parentCustID, err)
			return req, errors.New("distributor_ids not found")
		}
		if len(distributors) != len(distCodes) {
			return req, errors.New("distributor_ids not found")
		}
		req.DistributorIDs = make([]int64, 0, len(distributors))
		for _, distributor := range distributors {
			req.DistributorIDs = append(req.DistributorIDs, distributor.DistributorID)
		}
	}

	var parseErr error
	req.NewPurchPrice1, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_purch_price1"), snapshot.PurchPrice1)
	if parseErr != nil {
		return req, parseErr
	}
	req.NewPurchPrice2, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_purch_price2"), snapshot.PurchPrice2)
	if parseErr != nil {
		return req, parseErr
	}
	req.NewPurchPrice3, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_purch_price3"), snapshot.PurchPrice3)
	if parseErr != nil {
		return req, parseErr
	}
	req.NewSellPrice1, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_sell_price1"), snapshot.SellPrice1)
	if parseErr != nil {
		return req, parseErr
	}
	req.NewSellPrice2, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_sell_price2"), snapshot.SellPrice2)
	if parseErr != nil {
		return req, parseErr
	}
	req.NewSellPrice3, parseErr = parseManagePriceImportPriceOrCurrent(getCell(row, headerMap, "new_sell_price3"), snapshot.SellPrice3)
	if parseErr != nil {
		return req, parseErr
	}

	return req, nil
}

func parseManagePriceEffectiveDate(raw string) (effectiveDate time.Time, scheduleAt time.Time, immediate bool, err error) {
	asiaJkt, _ := time.LoadLocation("Asia/Jakarta")
	effectiveDate, err = time.ParseInLocation(mPriceDateLayout, raw, asiaJkt)
	if err != nil {
		return effectiveDate, scheduleAt, false, err
	}

	nowJkt := time.Now().In(asiaJkt)
	currentDate := time.Date(nowJkt.Year(), nowJkt.Month(), nowJkt.Day(), 0, 0, 0, 0, asiaJkt)
	if effectiveDate.Before(currentDate) {
		return effectiveDate, scheduleAt, false, errors.New("effective date must be current date or later")
	}

	scheduleAt = time.Date(
		effectiveDate.Year(),
		effectiveDate.Month(),
		effectiveDate.Day(),
		0,
		1,
		0,
		0,
		asiaJkt,
	)
	immediate = effectiveDate.Equal(currentDate)
	return effectiveDate, scheduleAt, immediate, nil
}

func normalizeDistributorIDs(coverage string, distributorIDs []int64, fallbackDistributorID int64) []int64 {
	if coverage == "N" {
		return nil
	}
	if len(distributorIDs) == 0 && fallbackDistributorID > 0 {
		return []int64{fallbackDistributorID}
	}
	return distributorIDs
}

func applySnapshotToCreateRequest(request *entity.CreateMPriceBody, snapshot model.MPriceProductSnapshot) {
	request.ProID = snapshot.ProID
	request.UnitID1 = snapshot.UnitID1
	request.UnitID2 = snapshot.UnitID2
	request.UnitID3 = snapshot.UnitID3
	request.ConvUnit2 = snapshot.ConvUnit2
	request.ConvUnit3 = snapshot.ConvUnit3
	request.PurchPrice1 = snapshot.PurchPrice1
	request.PurchPrice2 = snapshot.PurchPrice2
	request.PurchPrice3 = snapshot.PurchPrice3
	request.SellPrice1 = snapshot.SellPrice1
	request.SellPrice2 = snapshot.SellPrice2
	request.SellPrice3 = snapshot.SellPrice3
}

func applySnapshotToUpdateRequest(request *entity.UpdateMPriceRequest, snapshot model.MPriceProductSnapshot) {
	request.ProID = &snapshot.ProID
	request.UnitID1 = mPriceStringPtr(snapshot.UnitID1)
	request.UnitID2 = mPriceStringPtr(snapshot.UnitID2)
	request.UnitID3 = mPriceStringPtr(snapshot.UnitID3)
	request.ConvUnit2 = mPriceIntPtr(snapshot.ConvUnit2)
	request.ConvUnit3 = mPriceIntPtr(snapshot.ConvUnit3)
	request.PurchPrice1 = mPriceFloat64Ptr(snapshot.PurchPrice1)
	request.PurchPrice2 = mPriceFloat64Ptr(snapshot.PurchPrice2)
	request.PurchPrice3 = mPriceFloat64Ptr(snapshot.PurchPrice3)
	request.SellPrice1 = mPriceFloat64Ptr(snapshot.SellPrice1)
	request.SellPrice2 = mPriceFloat64Ptr(snapshot.SellPrice2)
	request.SellPrice3 = mPriceFloat64Ptr(snapshot.SellPrice3)
}

func exportManagePriceCSV(rows []entity.MPriceResponse) (*bytes.Buffer, string, string, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	for _, row := range buildManagePriceExportRows(rows) {
		if err := writer.Write(row); err != nil {
			return nil, "", "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", "", err
	}
	return buf, "text/csv", "manage_price_export.csv", nil
}

func buildManagePriceExportRows(rows []entity.MPriceResponse) [][]string {
	exportRows := [][]string{
		{
			"Effective Date", "Product Code", "Product Name", "Area Code", "Area Name",
			"Previous Purchase Price", "", "", "New Purchase Price", "", "",
			"Previous Selling Price", "", "", "New Selling Price", "", "",
			"Updated", "", "Not Updated", "",
		},
		{
			"", "", "", "", "",
			"Largest", "Middle", "Smallest",
			"Largest", "Middle", "Smallest",
			"Largest", "Middle", "Smallest",
			"Largest", "Middle", "Smallest",
			"Distributor Code", "Distributor Name",
			"Distributor Code", "Distributor Name",
		},
	}

	for _, row := range rows {
		details := row.Details
		if len(details) == 0 {
			details = []entity.DistributorAreaRegionData{{UpdateStatus: managePriceProductUpdateStatus(false)}}
		}
		prevPurchasePrices, newPurchasePrices := managePriceExportPriceGroups(
			[]float64{row.PurchPrice1, row.PurchPrice2, row.PurchPrice3},
			[]float64{row.NewPurchPrice1, row.NewPurchPrice2, row.NewPurchPrice3},
		)
		prevSellingPrices, newSellingPrices := managePriceExportPriceGroups(
			[]float64{row.SellPrice1, row.SellPrice2, row.SellPrice3},
			[]float64{row.NewSellPrice1, row.NewSellPrice2, row.NewSellPrice3},
		)

		for _, detail := range details {
			updatedDistributorCode, updatedDistributorName := "-", "-"
			notUpdatedDistributorCode, notUpdatedDistributorName := "-", "-"
			if detail.IsUpdated {
				updatedDistributorCode = detail.DistributorCode
				updatedDistributorName = detail.DistributorName
			} else {
				notUpdatedDistributorCode = detail.DistributorCode
				notUpdatedDistributorName = detail.DistributorName
			}

			exportRows = append(exportRows, []string{
				row.EffectiveDate,
				row.ProCode,
				row.ProName,
				strconv.Itoa(detail.AreaID),
				detail.AreaName,
				prevPurchasePrices[0],
				prevPurchasePrices[1],
				prevPurchasePrices[2],
				newPurchasePrices[0],
				newPurchasePrices[1],
				newPurchasePrices[2],
				prevSellingPrices[0],
				prevSellingPrices[1],
				prevSellingPrices[2],
				newSellingPrices[0],
				newSellingPrices[1],
				newSellingPrices[2],
				updatedDistributorCode,
				updatedDistributorName,
				notUpdatedDistributorCode,
				notUpdatedDistributorName,
			})
		}
	}

	return exportRows
}

func managePriceExportPriceGroups(previous, next []float64) ([]string, []string) {
	changed := false
	for i := range previous {
		if previous[i] != next[i] {
			changed = true
			break
		}
	}
	if !changed {
		return []string{"N.A", "N.A", "N.A"}, []string{"N.A", "N.A", "N.A"}
	}

	previousValues := make([]string, len(previous))
	nextValues := make([]string, len(next))
	for i := range previous {
		previousValues[i] = floatToString(previous[i])
		nextValues[i] = floatToString(next[i])
	}
	return previousValues, nextValues
}

func exportManagePriceTemplateCSV(rows [][]string, filename string) (*bytes.Buffer, string, string, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return nil, "", "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", "", err
	}
	return buf, "text/csv", filename, nil
}

type managePriceTemplateSheet struct {
	Name string
	Rows [][]string
}

func exportManagePriceTemplateWorkbook(sheets []managePriceTemplateSheet, filename, contentType string) (*bytes.Buffer, string, string, error) {
	f := excelize.NewFile()
	defaultSheet := f.GetSheetName(0)

	styleID, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
	})

	for sheetIndex, sheet := range sheets {
		sheetName := sheet.Name
		if sheetIndex == 0 {
			f.SetSheetName(defaultSheet, sheetName)
		} else {
			f.NewSheet(sheetName)
		}
		for rowIndex, row := range sheet.Rows {
			for colIndex, value := range row {
				cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
				f.SetCellValue(sheetName, cell, value)
				_ = f.SetCellStyle(sheetName, cell, cell, styleID)
			}
		}
	}
	f.SetActiveSheet(0)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}
	return buf, contentType, filename, nil
}

func exportManagePriceExcel(rows []entity.MPriceResponse, filename string) (*bytes.Buffer, string, string, error) {
	f := excelize.NewFile()
	sheet := "Manage Price"
	f.SetSheetName(f.GetSheetName(0), sheet)

	styleID, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	exportRows := buildManagePriceExportRows(rows)
	for rowIndex, row := range exportRows {
		for colIndex, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+1)
			f.SetCellValue(sheet, cell, value)
			if rowIndex < 2 {
				_ = f.SetCellStyle(sheet, cell, cell, styleID)
			}
		}
	}

	mergeRanges := [][2]string{
		{"A1", "A2"}, {"B1", "B2"}, {"C1", "C2"}, {"D1", "D2"}, {"E1", "E2"},
		{"F1", "H1"}, {"I1", "K1"}, {"L1", "N1"}, {"O1", "Q1"}, {"R1", "S1"}, {"T1", "U1"},
	}
	for _, mergeRange := range mergeRanges {
		_ = f.MergeCell(sheet, mergeRange[0], mergeRange[1])
	}
	for colIndex := 1; colIndex <= len(exportRows[0]); colIndex++ {
		colName, _ := excelize.ColumnNumberToName(colIndex)
		_ = f.SetColWidth(sheet, colName, colName, 16)
	}
	_ = f.SetColWidth(sheet, "C", "C", 32)
	_ = f.SetColWidth(sheet, "E", "E", 24)
	_ = f.SetColWidth(sheet, "S", "S", 28)
	_ = f.SetColWidth(sheet, "U", "U", 28)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}
	return buf, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", filename, nil
}

func downloadManagePriceImportRows(fileURL string) ([][]string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err == nil {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, errors.New("excel file has no sheets")
		}
		var fallbackRows [][]string
		for _, sheet := range sheets {
			rows, rowsErr := f.GetRows(sheet)
			if rowsErr != nil {
				continue
			}
			if fallbackRows == nil {
				fallbackRows = rows
			}
			if hasManagePriceImportDataRows(rows) {
				return rows, nil
			}
		}
		if fallbackRows != nil {
			return fallbackRows, nil
		}
		return nil, errors.New("excel file has no sheets")
	}

	reader := csv.NewReader(bytes.NewReader(fileBytes))
	rows, csvErr := reader.ReadAll()
	if csvErr == nil && len(rows) > 0 {
		return rows, nil
	}

	return nil, err
}

func buildManagePricePrincipalTemplateRows() [][]string {
	return [][]string{
		{
			"Effective date",
			"Product Code",
			"New Purchase Price",
			"",
			"",
			"New Selling Price",
			"",
			"",
			"Distributor Code",
		},
		{
			"",
			"",
			"Largest Unit",
			"Middle Unit",
			"Smallest Unit",
			"Largest Unit",
			"Middle Unit",
			"Smallest Unit",
			"",
		},
	}
}

func buildManagePriceDistributorTemplateRows() [][]string {
	return [][]string{
		{
			"Effective date",
			"Product Code",
			"New Purchase Price",
			"",
			"",
			"New Selling Price",
			"",
			"",
		},
		{
			"",
			"",
			"Largest Unit",
			"Middle Unit",
			"Smallest Unit",
			"Largest Unit",
			"Middle Unit",
			"Smallest Unit",
		},
	}
}

func hasManagePriceImportDataRows(rows [][]string) bool {
	if len(rows) < 3 {
		return false
	}
	headerMap, dataStartRow := buildMPriceHeaderIndex(rows)
	requiredHeaders := []string{
		"effective_date",
		"pro_code",
		"new_purch_price1",
		"new_purch_price2",
		"new_purch_price3",
		"new_sell_price1",
		"new_sell_price2",
		"new_sell_price3",
	}
	for _, header := range requiredHeaders {
		if _, ok := headerMap[header]; !ok {
			return false
		}
	}
	for _, row := range rows[dataStartRow:] {
		if !isEmptyManagePriceImportRow(row) {
			return true
		}
	}
	return false
}

func isDistributorManagePriceScope(custID, parentCustID string, distributorID int64) bool {
	return custID != "" && parentCustID != "" && custID != parentCustID && distributorID > 0
}

func isDistributorManagePriceTemplateScope(distributorID int64) bool {
	return distributorID > 0
}

func defaultManagePriceImportCoverage(custID, parentCustID string, distributorID int64) string {
	if isDistributorManagePriceScope(custID, parentCustID, distributorID) {
		return "D"
	}
	return "N"
}

func normalizeManagePriceImportDate(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", errors.New("effective_date is required")
	}

	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"2/1/2006",
		"02-01-2006",
		"2-1-2006",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed.Format(mPriceDateLayout), nil
		}
	}

	if numericValue, err := strconv.ParseFloat(value, 64); err == nil {
		baseDate := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		parsed := baseDate.AddDate(0, 0, int(numericValue))
		return parsed.Format(mPriceDateLayout), nil
	}

	return "", fmt.Errorf("invalid effective_date %s", raw)
}

func buildMPriceHeaderIndex(rows [][]string) (map[string]int, int) {
	result := make(map[string]int)
	if len(rows) == 0 {
		return result, 0
	}

	firstRow := rows[0]
	if len(rows) > 1 {
		secondRow := rows[1]
		matchedGroupedTemplate := false
		lastHeader := ""
		for index, header := range firstRow {
			if strings.TrimSpace(header) != "" {
				lastHeader = header
			} else {
				header = lastHeader
			}
			subHeader := ""
			if index < len(secondRow) {
				subHeader = secondRow[index]
			}
			key := normalizeManagePriceHeader(header, subHeader)
			if key == "" {
				continue
			}
			matchedGroupedTemplate = true
			result[key] = index
		}
		if matchedGroupedTemplate {
			return result, 2
		}
	}

	for index, header := range firstRow {
		key := normalizeManagePriceHeader(header, "")
		if key == "" {
			continue
		}
		result[key] = index
	}
	return result, 1
}

func getCell(row []string, headerMap map[string]int, key string) string {
	index, ok := headerMap[key]
	if !ok || index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func parseManagePriceImportPriceOrCurrent(raw string, currentPrice float64) (float64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return currentPrice, nil
	}

	price, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	if price == 0 {
		return currentPrice, nil
	}
	return price, nil
}

func splitCommaSeparated(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func normalizeManagePriceHeader(header string, subHeader string) string {
	normalizedHeader := normalizeManagePriceHeaderCell(header)
	normalizedSubHeader := normalizeManagePriceHeaderCell(subHeader)

	switch normalizedHeader {
	case "effective_date":
		return "effective_date"
	case "product_code", "pro_code":
		return "pro_code"
	case "coverage":
		return "coverage"
	case "distributor_code":
		return "distributor_code"
	case "new_purch_price1", "new_purch_price2", "new_purch_price3",
		"new_sell_price1", "new_sell_price2", "new_sell_price3":
		return normalizedHeader
	case "new_purchase_price":
		switch normalizedSubHeader {
		case "largest_unit":
			return "new_purch_price3"
		case "middle_unit":
			return "new_purch_price2"
		case "smallest_unit":
			return "new_purch_price1"
		}
	case "new_selling_price":
		switch normalizedSubHeader {
		case "largest_unit":
			return "new_sell_price3"
		case "middle_unit":
			return "new_sell_price2"
		case "smallest_unit":
			return "new_sell_price1"
		}
	}

	return ""
}

func normalizeManagePriceHeaderCell(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\r", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	return normalized
}

func isEmptyManagePriceImportRow(row []string) bool {
	for _, value := range row {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}

func floatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func mPriceStringPtr(value string) *string {
	return &value
}

func mPriceIntPtr(value int) *int {
	return &value
}

func mPriceFloat64Ptr(value float64) *float64 {
	return &value
}

func setupTransactionPriceData(price model.MPriceDetail, distributorID int64) (transPrice model.MTransactionPrice) {
	objectIDString := primitive.NewObjectID().Hex()
	transPrice = model.MTransactionPrice{
		CustID:             price.CustID,
		TransactionPriceID: objectIDString,
		ProID:              price.ProID,
		PurchPrice1:        price.NewPurchPrice1,
		PurchPrice2:        price.NewPurchPrice2,
		PurchPrice3:        price.NewPurchPrice3,
		SellPrice1:         price.NewSellPrice1,
		SellPrice2:         price.NewSellPrice2,
		SellPrice3:         price.NewSellPrice3,
		Source:             10,
		CreatedBy:          price.CreatedBy,
		CreatedAt:          time.Now().UTC(),
		StartDate:          price.EffectiveDate,
		EndDate:            nil,
		Coverage:           price.Coverage,
		DistributorID:      distributorID,
		ReferenceID:        price.PriceID,
	}
	return transPrice
}

func setupTransactionPriceDataUpdate(price model.MPriceDetail, distributorID int64, transactionPriceID string) (transPrice model.MTransactionPrice) {
	transPrice = model.MTransactionPrice{
		CustID:             price.CustID,
		TransactionPriceID: transactionPriceID,
		ProID:              price.ProID,
		PurchPrice1:        price.NewPurchPrice1,
		PurchPrice2:        price.NewPurchPrice2,
		PurchPrice3:        price.NewPurchPrice3,
		SellPrice1:         price.NewSellPrice1,
		SellPrice2:         price.NewSellPrice2,
		SellPrice3:         price.NewSellPrice3,
		Source:             10,
		CreatedBy:          price.CreatedBy,
		CreatedAt:          time.Now().UTC(),
		StartDate:          price.EffectiveDate,
		EndDate:            nil,
		Coverage:           price.Coverage,
		DistributorID:      distributorID,
		ReferenceID:        price.PriceID,
	}
	return transPrice
}
