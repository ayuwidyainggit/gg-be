master/controller/product_controller.go: added GET /report route immediately after POST /report, before GET /:pro_id.
master/controller/product_report_controller_test.go: added portable runtime.Caller/path/filepath module-root setup via t.Chdir for router tests; updated both helper callers.
