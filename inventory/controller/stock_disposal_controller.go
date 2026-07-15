package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"inventory/entity"
	"inventory/pkg/constant"
	"inventory/pkg/middleware"
	"inventory/pkg/responsebuild"
	"inventory/pkg/validation"
	"inventory/service"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
)

var (
	// Regex patterns compiled once for efficiency
	productsArrayRegex = regexp.MustCompile(`^products\[(\d+)\]\[(.+)\]$`)
)

type StockDisposalController struct {
	StockDisposalService service.StockDisposalService
	validator            *validation.Validate
}

func NewStockDisposalController(stockDisposalService service.StockDisposalService, validator *validation.Validate) *StockDisposalController {
	return &StockDisposalController{
		StockDisposalService: stockDisposalService,
		validator:            validator,
	}
}

func (controller *StockDisposalController) Route(app *fiber.App) {
	sdRouteV1 := app.Group("/v1/stock-disposal", middleware.JWTProtected())
	sdRouteV1.Post("/", controller.Store)
	sdRouteV1.Get("/products", controller.ProductLookup)
	sdRouteV1.Get("/:stock_disposal_id", controller.Detail)
	sdRouteV1.Get("/", controller.List)
}

// getAcceptLanguage extracts accept language from request headers
func (controller *StockDisposalController) getAcceptLanguage(c *fiber.Ctx) string {
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		return c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	return ""
}

// buildResponse builds response payload with accept language header
func (controller *StockDisposalController) buildResponse(c *fiber.Ctx) *responsebuild.DataRespReq {
	return responsebuild.BuildResponse(c.Locals("requestid").(string), controller.getAcceptLanguage(c))
}

// handleServiceError handles service errors and maps them to appropriate HTTP status codes
func (controller *StockDisposalController) handleServiceError(c *fiber.Ctx, err error, responsePayload *responsebuild.DataRespReq) error {
	statusCode := fiber.StatusBadRequest
	errMsg := err.Error()
	if errors.Is(err, sql.ErrNoRows) {
		statusCode = fiber.StatusNotFound
		errMsg = "Not found"
	}
	responsePayload.Setmsg(errMsg)
	return c.Status(statusCode).JSON(responsePayload.GetRespPayload())
}

// setPagination sets pagination information in response payload
func (controller *StockDisposalController) setPagination(responsePayload *responsebuild.DataRespReq, total int64, page, limit, lastPage int) {
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: page,
		PageLimit:   limit,
		PageTotal:   lastPage,
	})
}

func (controller *StockDisposalController) Store(c *fiber.Ctx) error {
	responsePayload := controller.buildResponse(c)

	// Parse multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		log.Error("StockDisposalController, Store, MultipartForm:", err.Error())
		responsePayload.Setmsg("Invalid multipart form data")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	var request entity.CreateStockDisposalBody

	// Parse form fields directly from form.Value (more reliable than c.FormValue after MultipartForm)
	if dateValues := form.Value["date"]; len(dateValues) > 0 && dateValues[0] != "" {
		request.Date = dateValues[0]
	}
	if supIDValues := form.Value["sup_id"]; len(supIDValues) > 0 && supIDValues[0] != "" {
		supID, err := strconv.ParseInt(supIDValues[0], 10, 64)
		if err != nil {
			log.Error("StockDisposalController, Store, ParseInt sup_id:", err.Error())
			responsePayload.Setmsg("Invalid sup_id format")
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		request.SupID = supID
	}
	if whIDValues := form.Value["wh_id"]; len(whIDValues) > 0 && whIDValues[0] != "" {
		whID, err := strconv.ParseInt(whIDValues[0], 10, 64)
		if err != nil {
			log.Error("StockDisposalController, Store, ParseInt wh_id:", err.Error())
			responsePayload.Setmsg("Invalid wh_id format")
			return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
		}
		request.WhID = whID
	}
	if grNoValues := form.Value["gr_no"]; len(grNoValues) > 0 && grNoValues[0] != "" {
		request.GrNo = &grNoValues[0]
	}
	if noteValues := form.Value["note"]; len(noteValues) > 0 {
		request.Note = noteValues[0]
	}

	// Parse products array from form fields like products[0][pro_id], products[0][unit_id1], etc.
	request.Products, err = controller.parseProductsFromForm(form)
	if err != nil {
		log.Error("StockDisposalController, Store, ParseProducts:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(request.Products) == 0 {
		log.Error("StockDisposalController, Store, No products found")
		responsePayload.Setmsg("At least one product is required")
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	userId := c.Locals("user_id").(int64)
	parentCustId := c.Locals("parent_cust_id").(string)

	request.CustID = custId
	request.ParentCustID = parentCustId
	request.CreatedBy = userId

	errs := controller.validator.ValidateStruct(request, controller.getAcceptLanguage(c))
	if errs != nil {
		log.Error("StockDisposalController, Store, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	// Call service Store method directly (no file upload needed, URL already in request)
	data, err := controller.StockDisposalService.Store(request)
	if err != nil {
		log.Error("StockDisposalController, Store, Service.Store, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	responsePayload.Setdata(data)
	responsePayload.Setmsg("Created Successfully")
	return c.Status(fiber.StatusCreated).JSON(responsePayload.GetRespPayload())
}

func (controller *StockDisposalController) Detail(c *fiber.Ctx) error {
	var params entity.DetailStockDisposalParams
	responsePayload := controller.buildResponse(c)

	if err := c.ParamsParser(&params); err != nil {
		log.Error("StockDisposalController, Detail, ParamsParser:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	errs := controller.validator.ValidateStruct(params, controller.getAcceptLanguage(c))
	if errs != nil {
		log.Error("StockDisposalController, Detail, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, err := controller.StockDisposalService.Detail(params.StockDisposalID, custId, parentCustId)
	if err != nil {
		log.Error("StockDisposalController, Detail, err:", err.Error())
		return controller.handleServiceError(c, err, responsePayload)
	}

	responsePayload.Setdata(data)
	return c.JSON(responsePayload.GetRespPayload())
}

func (controller *StockDisposalController) List(c *fiber.Ctx) error {
	var dataFilter entity.StockDisposalQueryFilter
	responsePayload := controller.buildResponse(c)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockDisposalController, List, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockDisposalController, List, ValidateStruct(params), errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.StockDisposalService.List(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("StockDisposalController, List, err:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg("")
		responsePayload.Data = nil
	} else {
		responsePayload.Setdata(data)
	}

	controller.setPagination(responsePayload, total, dataFilter.Page, dataFilter.Limit, lastPage)

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func (controller *StockDisposalController) ProductLookup(c *fiber.Ctx) error {
	var dataFilter entity.StockDisposalProductLookupQueryFilter
	responsePayload := controller.buildResponse(c)

	if err := c.QueryParser(&dataFilter); err != nil {
		log.Error("StockDisposalController, ProductLookup, query parser filter:", err.Error())
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}

	headerAcceptLang := ""
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	errs := controller.validator.ValidateStruct(dataFilter, headerAcceptLang)
	if errs != nil {
		log.Error("StockDisposalController, ProductLookup, ValidateStruct, errs:", errs)
		responsePayload.Setmsg(fiber.ErrBadRequest.Message)
		responsePayload.Seterrors(errs)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	custId := c.Locals("cust_id").(string)
	parentCustId := c.Locals("parent_cust_id").(string)

	data, total, lastPage, err := controller.StockDisposalService.ProductLookup(dataFilter, custId, parentCustId)
	if err != nil {
		log.Error("StockDisposalController, ProductLookup, err:", err.Error())
		responsePayload.Setmsg(constant.STOCK_DISPOSAL_PRODUCT_LOOKUP_FAILED)
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	if len(data) == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
		responsePayload.Data = nil
	} else {
		responsePayload.Setdata(data)
		responsePayload.Setmsg(constant.STOCK_DISPOSAL_PRODUCT_LOOKUP_SUCCESS)
	}

	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})

	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

// parseProductsFromForm parses products array from form fields like products[0][pro_id]
func (controller *StockDisposalController) parseProductsFromForm(form *multipart.Form) ([]entity.CreateStockDisposalProductBody, error) {
	var products []entity.CreateStockDisposalProductBody
	productMap := make(map[int]*entity.CreateStockDisposalProductBody)

	// Parse all form values
	for key, values := range form.Value {
		matches := productsArrayRegex.FindStringSubmatch(key)
		if len(matches) != 3 {
			continue
		}

		index, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid product index in field '%s': %w", key, err)
		}
		fieldName := matches[2]
		if len(values) == 0 {
			continue
		}
		value := values[0]

		// Initialize product if not exists
		if productMap[index] == nil {
			productMap[index] = &entity.CreateStockDisposalProductBody{}
		}
		product := productMap[index]

		// Parse field value based on field name
		switch fieldName {
		case "pro_id":
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid pro_id for product[%d]: %w", index, err)
			}
			product.ProID = val
		case "unit_id1":
			product.UnitID1 = value
		case "unit_id2":
			product.UnitID2 = value
		case "unit_id3":
			product.UnitID3 = value
		case "qty1":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid qty1 for product[%d]: %w", index, err)
			}
			product.Qty1 = val
		case "qty2":
			if value != "" {
				val, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid qty2 for product[%d]: %w", index, err)
				}
				product.Qty2 = val
			}
		case "qty3":
			if value != "" {
				val, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid qty3 for product[%d]: %w", index, err)
				}
				product.Qty3 = val
			}
		case "purch_price1":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid purch_price1 for product[%d]: %w", index, err)
			}
			product.PurchPrice1 = val
		case "purch_price2":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid purch_price2 for product[%d]: %w", index, err)
			}
			product.PurchPrice2 = val
		case "purch_price3":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid purch_price3 for product[%d]: %w", index, err)
			}
			product.PurchPrice3 = val
		case "gross_price":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid gross_price for product[%d]: %w", index, err)
			}
			product.GrossPrice = val
		case "vat":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid vat for product[%d]: %w", index, err)
			}
			product.Vat = val
		case "vat_value":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid vat_value for product[%d]: %w", index, err)
			}
			product.VatValue = val
		case "sub_total":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid sub_total for product[%d]: %w", index, err)
			}
			product.SubTotal = val
		case "file_url":
			// Parse file URL and extract metadata
			if value != "" {
				if product.UploadFile == nil {
					product.UploadFile = &entity.CreateStockDisposalFileBody{}
				}
				product.UploadFile.FileUrl = value
				// Extract metadata from URL
				fileName, fileType, mediaCategory := extractFileMetadataFromURL(value)
				product.UploadFile.FileName = fileName
				product.UploadFile.FileType = fileType
				product.UploadFile.MediaCategory = mediaCategory
				// FileSize cannot be extracted from URL, set to 0
				product.UploadFile.FileSize = 0
			}
		}
	}

	// Convert map to slice maintaining order
	maxIndex := -1
	for idx := range productMap {
		if idx > maxIndex {
			maxIndex = idx
		}
	}

	for i := 0; i <= maxIndex; i++ {
		if product, exists := productMap[i]; exists {
			products = append(products, *product)
		}
	}

	return products, nil
}

// extractFileMetadataFromURL extracts file metadata (filename, type, category) from URL
func extractFileMetadataFromURL(fileURL string) (fileName, fileType, mediaCategory string) {
	// Parse URL to get path
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		// If URL parsing fails, try to extract filename from string directly
		fileName = filepath.Base(fileURL)
	} else {
		// Extract filename from URL path
		fileName = filepath.Base(parsedURL.Path)
	}

	// If filename is empty or just "/", use a default
	if fileName == "" || fileName == "/" {
		fileName = "file"
	}

	// Extract file extension and determine type
	ext := strings.ToLower(filepath.Ext(fileName))

	// Determine file type and media category based on extension
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		fileType = "image"
		mediaCategory = "image"
	case ".mp4", ".avi", ".mov", ".mkv":
		fileType = "video"
		mediaCategory = "video"
	default:
		// Default to image if cannot determine
		fileType = "image"
		mediaCategory = "image"
	}

	return fileName, fileType, mediaCategory
}
