package controller

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"sales/entity"
	"sales/pkg/validation"
	"sales/service"

	"github.com/gofiber/fiber/v2"
)

type mockReportServiceForController struct {
	publishSecondarySalesReportFn func(dataFilter entity.SecondarySalesReportQueryFilter) (entity.ReportList, error)
	sumByMonthFn                  func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error)
	groupSalesFn                  func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) ([]entity.SecondarySalesReportGroupResp, error)
	trendSalesFn                  func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.SumReportTrendSalesResp, error)
	activityTrendSalesFn          func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.ActivityReportTrendSalesResp, error)
	activityGeotagFn              func(authCustID, parentCustID string, req entity.ActivityReportGeotagPayload) (entity.ActivityReportGeotagResp, error)
}

func (m *mockReportServiceForController) PublishSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) (data entity.ReportList, err error) {
	if m.publishSecondarySalesReportFn != nil {
		return m.publishSecondarySalesReportFn(dataFilter)
	}
	return data, nil
}
func (m *mockReportServiceForController) List(dataFilter entity.ReportQueryFilter) (data []entity.ReportList, total int64, lastPage int, err error) {
	return nil, 0, 0, nil
}
func (m *mockReportServiceForController) SubscribeSecondarySalesReport(dataFilter entity.SecondarySalesReportQueryFilter) error {
	return nil
}
func (m *mockReportServiceForController) PublishActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (data entity.ReportList, err error) {
	return data, nil
}
func (m *mockReportServiceForController) SubscribeActivitySalesReport(dataFilter entity.ActivityReportQueryFilter) (err error) {
	return nil
}
func (m *mockReportServiceForController) PublishActivitySalesReportList(dataFilter entity.ActivityReportQueryFilterList) (results []entity.ActivityReportListResp, total int64, lastPage int, err error) {
	return nil, 0, 0, nil
}
func (m *mockReportServiceForController) ExtractReportSecondary(req entity.SecondarySalesReportDashboardExtractQueryFilter) (err error) {
	return nil
}
func (m *mockReportServiceForController) SecondarySalesReportSumReportByMonth(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error) {
	if m.sumByMonthFn != nil {
		return m.sumByMonthFn(authCustID, parentCustID, req)
	}
	return entity.SumReportByMonthModelResp{}, nil
}
func (m *mockReportServiceForController) SecondarySalesReportGroupSales(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) (datas []entity.SecondarySalesReportGroupResp, err error) {
	if m.groupSalesFn != nil {
		return m.groupSalesFn(authCustID, parentCustID, req)
	}
	return nil, nil
}
func (m *mockReportServiceForController) SecondarySalesReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.SumReportTrendSalesResp, err error) {
	if m.trendSalesFn != nil {
		return m.trendSalesFn(authCustID, parentCustID, year, requestedCustIDs)
	}
	return nil, nil
}
func (m *mockReportServiceForController) SalesmanActivityReportSumReportByMonth(authCustID, parentCustID string, req entity.SalesmanActivityReportDashboardSumPayload) (data entity.SalesmanActivityReportByMonthModelResp, err error) {
	return data, nil
}
func (m *mockReportServiceForController) SalesmanActivityReportTrendSales(authCustID, parentCustID string, year int, requestedCustIDs []string) (data []entity.ActivityReportTrendSalesResp, err error) {
	if m.activityTrendSalesFn != nil {
		return m.activityTrendSalesFn(authCustID, parentCustID, year, requestedCustIDs)
	}
	return nil, nil
}
func (m *mockReportServiceForController) SalesmanActivityReportGeotag(authCustID, parentCustID string, req entity.ActivityReportGeotagPayload) (entity.ActivityReportGeotagResp, error) {
	if m.activityGeotagFn != nil {
		return m.activityGeotagFn(authCustID, parentCustID, req)
	}
	return entity.ActivityReportGeotagResp{}, nil
}
func (m *mockReportServiceForController) SalesmanActivityReportGroupSales(authCustID, parentCustID string, req entity.SalesmanActivityReportDashboardGroupPayload) (datas []entity.SecondarySalesReportGroupResp, err error) {
	return nil, nil
}
func (m *mockReportServiceForController) SalesmanActivitySalesmanList(dataFilter entity.ActivityReportSalesmanListQueryFilter) (datas []entity.SalesmanActivityReportSalesmanListResp, err error) {
	return nil, nil
}

func TestParseDownloadSalesmanIDs(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []int64
	}{
		{
			name:     "repeated salesman_id values",
			query:    "/?salesman_id=62&salesman_id=204&salesman_id=206",
			expected: []int64{62, 204, 206},
		},
		{
			name:     "repeated salesman_id bracket values",
			query:    "/?salesman_id[]=62&salesman_id[]=204&salesman_id[]=206",
			expected: []int64{62, 204, 206},
		},
		{
			name:     "comma separated values",
			query:    "/?salesman_id=62,204,206",
			expected: []int64{62, 204, 206},
		},
		{
			name:     "mixed repeated and comma separated values",
			query:    "/?salesman_id=62,204&salesman_id=206&salesman_id[]=207,208&salesman_id[]=209",
			expected: []int64{62, 204, 206, 207, 208, 209},
		},
		{
			name:     "empty and invalid values become no filter",
			query:    "/?salesman_id=&salesman_id[]=abc&salesman_id[]=%20",
			expected: []int64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/", func(c *fiber.Ctx) error {
				parsed := parseDownloadSalesmanIDs(c)
				return c.JSON(parsed)
			})

			req := httptest.NewRequest("GET", tt.query, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app test request failed: %v", err)
			}

			defer resp.Body.Close()

			var got []int64
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("unexpected parsed ids: got=%v expected=%v", got, tt.expected)
			}
		})
	}
}

func TestSecondaryReportSalesSumMonthReturnsForbiddenForUnauthorizedCustID(t *testing.T) {
	validator := validation.NewValiditor()
	controller := NewReportController(&mockReportServiceForController{
		sumByMonthFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error) {
			return entity.SumReportByMonthModelResp{}, service.ErrUnauthorizedCustID
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-1")
		c.Locals("cust_id", "DIST1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/sum-date", controller.SecondaryReportSalesSumMonth)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/sum-date?month=5&year=2026&cust_id=SIBLING1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["request_id"] != "req-1" {
		t.Fatalf("expected request_id req-1, got %v", body["request_id"])
	}
}

func TestSecondaryReportSalesGroupReturnsForbiddenForUnauthorizedCustID(t *testing.T) {
	validator := validation.NewValiditor()
	controller := NewReportController(&mockReportServiceForController{
		groupSalesFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) ([]entity.SecondarySalesReportGroupResp, error) {
			return nil, service.ErrUnauthorizedCustID
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-2")
		c.Locals("cust_id", "DIST1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/group", controller.SecondaryReportSalesGroup)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/group?month=5&year=2026&group_by=outlet&cust_id=SIBLING1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["request_id"] != "req-2" {
		t.Fatalf("expected request_id req-2, got %v", body["request_id"])
	}
}

func TestSecondaryReportSalesSumMonthParsesQueryAndAuthLocals(t *testing.T) {
	validator := validation.NewValiditor()
	var gotAuth, gotParent string
	var gotReq entity.SecondarySalesReportDashboardSumPayload
	ctrl := NewReportController(&mockReportServiceForController{
		sumByMonthFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error) {
			gotAuth = authCustID
			gotParent = parentCustID
			gotReq = req
			return entity.SumReportByMonthModelResp{}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-sum-success")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/sum-date", ctrl.SecondaryReportSalesSumMonth)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotAuth != "AUTH1" || gotParent != "PARENT1" {
		t.Fatalf("unexpected auth locals auth=%s parent=%s", gotAuth, gotParent)
	}
	if gotReq.Month != 6 {
		t.Fatalf("expected month 6, got %d", gotReq.Month)
	}
	if gotReq.Year == nil || *gotReq.Year != 2026 {
		t.Fatalf("expected year 2026, got %#v", gotReq.Year)
	}
	if gotReq.CustID != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %q", gotReq.CustID)
	}
	if !reflect.DeepEqual(gotReq.CustIDs, []string{"C260020001"}) {
		t.Fatalf("expected cust_ids [C260020001], got %#v", gotReq.CustIDs)
	}
}

func TestSecondaryReportSalesSumMonthAllowsMissingYear(t *testing.T) {
	validator := validation.NewValiditor()
	var gotReq entity.SecondarySalesReportDashboardSumPayload
	ctrl := NewReportController(&mockReportServiceForController{
		sumByMonthFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error) {
			gotReq = req
			if authCustID != "AUTH1" || parentCustID != "PARENT1" {
				t.Fatalf("unexpected auth locals auth=%s parent=%s", authCustID, parentCustID)
			}
			return entity.SumReportByMonthModelResp{}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-sum-no-year")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/sum-date", ctrl.SecondaryReportSalesSumMonth)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/sum-date?month=6&cust_id=C260020001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotReq.Month != 6 {
		t.Fatalf("expected month 6, got %d", gotReq.Month)
	}
	if gotReq.Year != nil {
		t.Fatalf("expected nil year, got %#v", gotReq.Year)
	}
	if gotReq.CustID != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %q", gotReq.CustID)
	}
	if !reflect.DeepEqual(gotReq.CustIDs, []string{"C260020001"}) {
		t.Fatalf("expected cust_ids [C260020001], got %#v", gotReq.CustIDs)
	}
}

func TestSecondaryReportSalesSumMonthInvalidMonthYearReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-sum-invalid")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/sum-date", ctrl.SecondaryReportSalesSumMonth)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/sum-date?month=13&year=10000&cust_id=C260020001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSecondaryReportSalesSumMonthReturnsSubtractValuesFromService(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{
		sumByMonthFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardSumPayload) (entity.SumReportByMonthModelResp, error) {
			return entity.SumReportByMonthModelResp{
				Qty:                134,
				TotalDiscountPromo: 1_238_740,
				QtyReturn:          16,
				NetSalesReturn:     200,
				ReturnRate:         11.94,
			}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-sum-values")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/sum-date", ctrl.SecondaryReportSalesSumMonth)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/sum-date?month=6&year=2026&cust_id=C260020001", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
	if got := data["qty"]; got != float64(134) {
		t.Fatalf("expected qty 134, got %#v", got)
	}
	if got := data["total_discount_promo"]; got != float64(1238740) {
		t.Fatalf("expected total_discount_promo 1238740, got %#v", got)
	}
}

func TestSecondaryReportSalesGroupParsesQueryAndAuthLocals(t *testing.T) {
	validator := validation.NewValiditor()
	var gotAuth, gotParent string
	var gotReq entity.SecondarySalesReportDashboardGroupPayload
	ctrl := NewReportController(&mockReportServiceForController{
		groupSalesFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) ([]entity.SecondarySalesReportGroupResp, error) {
			gotAuth = authCustID
			gotParent = parentCustID
			gotReq = req
			return nil, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-group-success")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/group", ctrl.SecondaryReportSalesGroup)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/group?month=6&year=2026&cust_id=C260020001&group_by=outlet", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotAuth != "AUTH1" || gotParent != "PARENT1" {
		t.Fatalf("unexpected auth locals auth=%s parent=%s", gotAuth, gotParent)
	}
	if gotReq.Month != 6 {
		t.Fatalf("expected month 6, got %d", gotReq.Month)
	}
	if gotReq.Year == nil || *gotReq.Year != 2026 {
		t.Fatalf("expected year 2026, got %#v", gotReq.Year)
	}
	if gotReq.CustID != "C260020001" {
		t.Fatalf("expected cust_id C260020001, got %q", gotReq.CustID)
	}
	if !reflect.DeepEqual(gotReq.CustIDs, []string{"C260020001"}) {
		t.Fatalf("expected cust_ids [C260020001], got %#v", gotReq.CustIDs)
	}
	if gotReq.GroupBy != "outlet" {
		t.Fatalf("expected group_by outlet, got %q", gotReq.GroupBy)
	}
}

func TestSecondaryReportSalesGroupInvalidMonthYearReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-group-invalid")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/group", ctrl.SecondaryReportSalesGroup)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/group?month=13&year=10000&cust_id=C260020001&group_by=outlet", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSecondarySalesExportReturnsForbiddenForDistributorSiblingCust(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{
		publishSecondarySalesReportFn: func(dataFilter entity.SecondarySalesReportQueryFilter) (entity.ReportList, error) {
			return entity.ReportList{}, service.ErrUnauthorizedCustID
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-403")
		c.Locals("cust_id", "DIST1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_fullname", "Tester")
		return c.Next()
	})
	app.Post("/reports/secondary-sales", ctrl.SecondarySales)

	from := int64(1777568400)
	to := int64(1779123599)
	bodyJSON, _ := json.Marshal(map[string]interface{}{
		"from":    from,
		"to":      to,
		"cust_id": "SIBLING1",
	})
	req := httptest.NewRequest("POST", "/reports/secondary-sales", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestSecondarySalesExportAuthCustNotOverwrittenByBody(t *testing.T) {
	validator := validation.NewValiditor()
	var gotFilter entity.SecondarySalesReportQueryFilter
	ctrl := NewReportController(&mockReportServiceForController{
		publishSecondarySalesReportFn: func(dataFilter entity.SecondarySalesReportQueryFilter) (entity.ReportList, error) {
			gotFilter = dataFilter
			return entity.ReportList{}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-auth")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_fullname", "Tester")
		return c.Next()
	})
	app.Post("/reports/secondary-sales", ctrl.SecondarySales)

	from := int64(1777568400)
	to := int64(1779123599)
	// Body tries to spoof _cust_id and _parent_cust_id — must be ignored.
	bodyJSON, _ := json.Marshal(map[string]interface{}{
		"from":            from,
		"to":              to,
		"_cust_id":        "SPOOFED",
		"_parent_cust_id": "SPOOFED",
	})
	req := httptest.NewRequest("POST", "/reports/secondary-sales", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotFilter.CustID != "AUTH1" {
		t.Fatalf("expected auth cust_id AUTH1, got %s", gotFilter.CustID)
	}
	if gotFilter.ParentCustID != "PARENT1" {
		t.Fatalf("expected parent_cust_id PARENT1, got %s", gotFilter.ParentCustID)
	}
}

func TestSecondarySalesBodyParserSupportsStringAndArrayCustID(t *testing.T) {
	bodyString := []byte(`{"cust_id":"CHILD1"}`)
	var stringBody rawSecondarySalesExportBody
	if err := json.Unmarshal(bodyString, &stringBody); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	var stringList entity.StringListOrScalar
	if err := json.Unmarshal(stringBody.RequestedCustIDRaw, &stringList); err != nil {
		t.Fatalf("unexpected string cust_id decode error: %v", err)
	}
	if !reflect.DeepEqual([]string(stringList), []string{"CHILD1"}) {
		t.Fatalf("unexpected string cust_id decode: %#v", []string(stringList))
	}

	bodyArray := []byte(`{"cust_id":["CHILD1","CHILD2"]}`)
	var arrayBody rawSecondarySalesExportBody
	if err := json.Unmarshal(bodyArray, &arrayBody); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	var arrayList entity.StringListOrScalar
	if err := json.Unmarshal(arrayBody.RequestedCustIDRaw, &arrayList); err != nil {
		t.Fatalf("unexpected array cust_id decode error: %v", err)
	}
	if !reflect.DeepEqual([]string(arrayList), []string{"CHILD1", "CHILD2"}) {
		t.Fatalf("unexpected array cust_id decode: %#v", []string(arrayList))
	}
}

func TestSecondaryReportSalesTrendSalesReturnsForbiddenForDistributorSibling(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{
		trendSalesFn: func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.SumReportTrendSalesResp, error) {
			return nil, service.ErrUnauthorizedCustID
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-trend-403")
		c.Locals("cust_id", "DIST1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/trend-sales", ctrl.SecondaryReportSalesTrendSales)

	bodyJSON := []byte(`{"cust_id":"SIBLING1"}`)
	req := httptest.NewRequest("GET", "/reports/secondary-sales/trend-sales?year=2026", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestSecondaryReportSalesTrendSalesPassesCustIDFromBody(t *testing.T) {
	validator := validation.NewValiditor()
	var gotAuth, gotParent, gotRequested string
	var gotYear int
	ctrl := NewReportController(&mockReportServiceForController{
		trendSalesFn: func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.SumReportTrendSalesResp, error) {
			gotAuth = authCustID
			gotParent = parentCustID
			gotYear = year
			gotRequested = strings.Join(requestedCustIDs, ",")
			return nil, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-trend-body")
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/trend-sales", ctrl.SecondaryReportSalesTrendSales)

	bodyJSON := []byte(`{"cust_id":"CHILD1"}`)
	req := httptest.NewRequest("GET", "/reports/secondary-sales/trend-sales?year=2026", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotAuth != "PARENT1" || gotParent != "PARENT1" || gotYear != 2026 || gotRequested != "CHILD1" {
		t.Fatalf("unexpected service args auth=%s parent=%s year=%d requested=%s", gotAuth, gotParent, gotYear, gotRequested)
	}
}

func TestSecondaryReportSalesTrendSalesPassesCustIDQueryMulti(t *testing.T) {
	validator := validation.NewValiditor()
	var gotRequested []string
	ctrl := NewReportController(&mockReportServiceForController{
		trendSalesFn: func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.SumReportTrendSalesResp, error) {
			gotRequested = append([]string(nil), requestedCustIDs...)
			return nil, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-trend-querymulti")
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/trend-sales", ctrl.SecondaryReportSalesTrendSales)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/trend-sales?year=2026&cust_id=CHILD1,CHILD2&cust_id=CHILD3", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if !reflect.DeepEqual(gotRequested, []string{"CHILD1", "CHILD2", "CHILD3"}) {
		t.Fatalf("unexpected requested cust ids: %#v", gotRequested)
	}
}

func TestSecondaryReportSalesTrendSalesYearRequiredReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-trend-noyear")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/trend-sales", ctrl.SecondaryReportSalesTrendSales)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/trend-sales", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 when year missing, got %d", resp.StatusCode)
	}
}

// TestSecondarySalesExportRequestedCustIDPassedToService verifies the controller
// correctly passes RequestedCustID from body to the service filter.
func TestSecondarySalesExportRequestedCustIDPassedToService(t *testing.T) {
	validator := validation.NewValiditor()
	var gotFilter entity.SecondarySalesReportQueryFilter
	ctrl := NewReportController(&mockReportServiceForController{
		publishSecondarySalesReportFn: func(dataFilter entity.SecondarySalesReportQueryFilter) (entity.ReportList, error) {
			gotFilter = dataFilter
			return entity.ReportList{}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-child")
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_fullname", "Tester")
		return c.Next()
	})
	app.Post("/reports/secondary-sales", ctrl.SecondarySales)

	from := int64(1777568400)
	to := int64(1779123599)
	bodyJSON, _ := json.Marshal(map[string]interface{}{
		"from":    from,
		"to":      to,
		"cust_id": "CHILD1",
	})
	req := httptest.NewRequest("POST", "/reports/secondary-sales", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if gotFilter.RequestedCustID != "CHILD1" {
		t.Fatalf("expected RequestedCustID CHILD1, got %s", gotFilter.RequestedCustID)
	}
	if !reflect.DeepEqual([]string(gotFilter.RequestedCustIDs), []string{"CHILD1"}) {
		t.Fatalf("expected RequestedCustIDs [CHILD1], got %#v", []string(gotFilter.RequestedCustIDs))
	}
	if gotFilter.CustID != "PARENT1" {
		t.Fatalf("expected CustID PARENT1 from JWT, got %s", gotFilter.CustID)
	}
}

func TestSecondarySalesExportRequestedCustIDArrayPassedToService(t *testing.T) {
	validator := validation.NewValiditor()
	var gotFilter entity.SecondarySalesReportQueryFilter
	ctrl := NewReportController(&mockReportServiceForController{
		publishSecondarySalesReportFn: func(dataFilter entity.SecondarySalesReportQueryFilter) (entity.ReportList, error) {
			gotFilter = dataFilter
			return entity.ReportList{}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-array")
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_fullname", "Tester")
		return c.Next()
	})
	app.Post("/reports/secondary-sales", ctrl.SecondarySales)

	from := int64(1777568400)
	to := int64(1779123599)
	bodyJSON, _ := json.Marshal(map[string]interface{}{
		"from":    from,
		"to":      to,
		"cust_id": []string{"CHILD1", "CHILD2"},
	})
	req := httptest.NewRequest("POST", "/reports/secondary-sales", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if !reflect.DeepEqual([]string(gotFilter.RequestedCustIDs), []string{"CHILD1", "CHILD2"}) {
		t.Fatalf("expected RequestedCustIDs [CHILD1 CHILD2], got %#v", []string(gotFilter.RequestedCustIDs))
	}
}

func TestSecondaryReportSalesGroupPassesCodeResponse(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{
		groupSalesFn: func(authCustID, parentCustID string, req entity.SecondarySalesReportDashboardGroupPayload) ([]entity.SecondarySalesReportGroupResp, error) {
			return []entity.SecondarySalesReportGroupResp{{ID: 1, Code: "OUT-1", Name: "Outlet A", NetSales: 100}}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-group-code")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/group", ctrl.SecondaryReportSalesGroup)

	req := httptest.NewRequest("GET", "/reports/secondary-sales/group?month=5&group_by=outlet", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	data, ok := body["data"].([]interface{})
	if !ok || len(data) != 1 {
		t.Fatalf("unexpected response data: %#v", body["data"])
	}
	item := data[0].(map[string]interface{})
	if item["code"] != "OUT-1" {
		t.Fatalf("expected code OUT-1, got %#v", item["code"])
	}
}

func TestSecondarySalesExportInvalidCustIDReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-export-invalid-cust")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("user_fullname", "Tester")
		return c.Next()
	})
	app.Post("/reports/secondary-sales", ctrl.SecondarySales)

	from := int64(1777568400)
	to := int64(1779123599)
	bodyJSON, _ := json.Marshal(map[string]interface{}{
		"from":    from,
		"to":      to,
		"cust_id": "BAD CUST!",
	})
	req := httptest.NewRequest("POST", "/reports/secondary-sales", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSecondaryReportSalesTrendSalesInvalidCustIDReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("requestid", "req-trend-invalid-cust")
		c.Locals("cust_id", "AUTH1")
		c.Locals("parent_cust_id", "PARENT1")
		return c.Next()
	})
	app.Get("/reports/secondary-sales/trend-sales", ctrl.SecondaryReportSalesTrendSales)

	bodyJSON := []byte(`{"cust_id":"BAD CUST!"}`)
	req := httptest.NewRequest("GET", "/reports/secondary-sales/trend-sales?year=2026", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app test request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestSalesmanActivityTrendSalesReturnsForbiddenForDistributorSibling(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{
		activityTrendSalesFn: func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.ActivityReportTrendSalesResp, error) {
			return nil, service.ErrUnauthorizedCustID
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("cust_id", "DIST1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("requestid", "req-activity-trend-403")
		return c.Next()
	})
	app.Get("/reports/activity-report-sales/trend-sales", ctrl.SalesmanActivityTrendSales)

	req := httptest.NewRequest("GET", "/reports/activity-report-sales/trend-sales?year=2026&cust_id=SIBLING1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestSalesmanActivityTrendSalesPassesCustIDFromQuery(t *testing.T) {
	validator := validation.NewValiditor()
	var gotCustIDs []string
	ctrl := NewReportController(&mockReportServiceForController{
		activityTrendSalesFn: func(authCustID, parentCustID string, year int, requestedCustIDs []string) ([]entity.ActivityReportTrendSalesResp, error) {
			gotCustIDs = requestedCustIDs
			return []entity.ActivityReportTrendSalesResp{{Month: 1, TotalInvoice: 100, TotalReturn: 10, NetSales: 90}}, nil
		},
	}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("requestid", "req-activity-trend-query")
		return c.Next()
	})
	app.Get("/reports/activity-report-sales/trend-sales", ctrl.SalesmanActivityTrendSales)

	req := httptest.NewRequest("GET", "/reports/activity-report-sales/trend-sales?year=2026&cust_id=CHILD1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(gotCustIDs) != 1 || gotCustIDs[0] != "CHILD1" {
		t.Fatalf("expected CHILD1 cust id, got %#v", gotCustIDs)
	}
}

func TestSalesmanActivityTrendSalesYearRequiredReturns400(t *testing.T) {
	validator := validation.NewValiditor()
	ctrl := NewReportController(&mockReportServiceForController{}, validator)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("cust_id", "PARENT1")
		c.Locals("parent_cust_id", "PARENT1")
		c.Locals("requestid", "req-activity-trend-noyear")
		return c.Next()
	})
	app.Get("/reports/activity-report-sales/trend-sales", ctrl.SalesmanActivityTrendSales)

	req := httptest.NewRequest("GET", "/reports/activity-report-sales/trend-sales", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

// ensure strings import is used
var _ = strings.Contains
