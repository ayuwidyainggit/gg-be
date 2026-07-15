package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"master/entity"
	"master/pkg/constant"
	"master/pkg/middleware"
	"master/pkg/responsebuild"
	"master/pkg/validation"
	"master/service"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type BusinessUnitController struct {
	Service   service.BusinessUnitService
	Validator *validation.Validate
}

type businessUnitListQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Query string `query:"q"`
	Sort  string `query:"sort"`
}

func NewBusinessUnitController(svc service.BusinessUnitService, validator *validation.Validate) *BusinessUnitController {
	return &BusinessUnitController{
		Service:   svc,
		Validator: validator,
	}
}

func (controller *BusinessUnitController) Route(app *fiber.App) {
	route := app.Group("/v1/business-unit", middleware.JWTProtected())
	route.Get("", controller.List)
}

func (controller *BusinessUnitController) List(c *fiber.Ctx) error {
	var dataFilter entity.BusinessUnitQueryFilter
	var query businessUnitListQuery
	var headerAcceptLang string
	if len(c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG]) > 0 {
		headerAcceptLang = c.GetReqHeaders()[constant.HEADER_ACCEPT_LANG][0]
	}
	responsePayload := responsebuild.BuildResponse(c.Locals("requestid").(string), headerAcceptLang)

	if err := c.QueryParser(&query); err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusUnprocessableEntity).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.Page = query.Page
	dataFilter.Limit = query.Limit
	dataFilter.Query = query.Query
	dataFilter.Sort = query.Sort

	dataFilter.CustId = c.Locals("cust_id").(string)
	dataFilter.ParentCustId = c.Locals("parent_cust_id").(string)
	dataFilter.UserName = c.Locals("user_name").(string)
	dataFilter.EmployeeId = localIntValue(c.Locals("employee_id"))
	distributorID := localIntValue(c.Locals("distributor_id"))
	dataFilter.DistributorId = &distributorID

	regionIDs, err := normalizeIntArrayQueryStrict(c.Context().QueryArgs(), "region_id[]", "region_id")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.RegionId = regionIDs

	areaIDs, err := normalizeIntArrayQueryStrict(c.Context().QueryArgs(), "area_id[]", "area_id")
	if err != nil {
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}
	dataFilter.AreaId = areaIDs

	if active := normalizeIntArrayQuery(c.Context().QueryArgs(), "is_active[]", "is_active"); len(active) > 0 {
		dataFilter.IsActive = active
	}

	data, total, lastPage, err := controller.Service.GetBusinessUnit(dataFilter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responsePayload.Setmsg(constant.RECORD_NOT_FOUND)
			return c.Status(fiber.StatusNotFound).JSON(responsePayload.GetRespPayload())
		}
		responsePayload.Setmsg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(responsePayload.GetRespPayload())
	}

	responsePayload.Setdata(data)
	if total == 0 {
		responsePayload.Setmsg(constant.NO_DATA)
	} else {
		responsePayload.Setmsg(constant.SUCCESS_NO_DATA_DISPLAYED)
	}
	responsePayload.Setpaging(entity.Pagination{
		TotalRecord: total,
		PageCurrent: dataFilter.Page,
		PageLimit:   dataFilter.Limit,
		PageTotal:   lastPage,
	})
	return c.Status(fiber.StatusOK).JSON(responsePayload.GetRespPayload())
}

func localIntValue(v interface{}) int {
	switch value := v.(type) {
	case int:
		return value
	case int8:
		return int(value)
	case int16:
		return int(value)
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		return 0
	}
}

func normalizeIntArrayQuery(args *fasthttp.Args, keys ...string) []int {
	values, _ := normalizeIntArrayQueryStrict(args, keys...)
	return values
}

func normalizeIntArrayQueryStrict(args *fasthttp.Args, keys ...string) ([]int, error) {
	result := make([]int, 0)
	seen := make(map[int]struct{})

	for _, key := range keys {
		values := args.PeekMulti(key)
		for _, value := range values {
			tokens := strings.Split(string(value), ",")
			for _, token := range tokens {
				cleaned := strings.TrimSpace(token)
				if cleaned == "" {
					continue
				}

				parsedValue, err := strconv.Atoi(cleaned)
				if err != nil {
					return nil, fmt.Errorf("invalid %s value %q", normalizeArrayQueryKey(keys[0]), cleaned)
				}

				if _, exists := seen[parsedValue]; exists {
					continue
				}

				seen[parsedValue] = struct{}{}
				result = append(result, parsedValue)
			}
		}
	}

	return result, nil
}

func normalizeArrayQueryKey(key string) string {
	return strings.TrimSuffix(key, "[]")
}
