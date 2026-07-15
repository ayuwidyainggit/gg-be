# D6 Coverage Remediation

Command:

```text
cd sales && rtk go test -coverprofile=/tmp/sx2516.cov -coverpkg=./service/... -run 'TestImportSecondarySales|TestParseImportOrders|TestValidateImportDate' ./service/...
cd sales && rtk go tool cover -func=/tmp/sx2516.cov
```

Result:

```text
Go test: 18 passed in 1 packages
sales/service/order_service.go:6592:             validateImportDate                         100.0%
sales/service/order_service.go:6604:             importSecondarySales                       91.5%
```

Added direct tests for invalid scope date, nil date scope skip, optional date parse errors, StoreDetail failure, and successful optional-date/order-detail mapping. No production files changed.
