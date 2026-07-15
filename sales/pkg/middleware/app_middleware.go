package middleware

import (
	"encoding/json"
	"net/url"
	"regexp"
	"sales/pkg/config/env"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/utils"
)

// AppMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func AppMiddleware(a *fiber.App, envCfg env.ConfigEnv) {
	var ConfigDefault = requestid.Config{
		Next:   nil,
		Header: fiber.HeaderXRequestID,
		Generator: func() string {
			return utils.UUID()
		},
		ContextKey: "requestid",
	}
	a.Use(requestid.New(ConfigDefault))

	// Enable response compression (gzip) for all routes
	// This reduces network transfer size significantly for JSON responses
	a.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Parse included routes from environment variable (comma-separated)
	// Only routes in this list will log response bodies and request headers; all others will return empty string
	// Example: LOGGER_INCLUDED_ROUTES="/v1/products,/api/reports/export,/health"
	// If empty or not set, all routes will return empty string for response body and request headers logging
	includedRoutesStr := envCfg.Get("LOGGER_INCLUDED_ROUTES")
	var includedRoutes []string

	if includedRoutesStr != "" {
		// Split comma-separated routes and trim whitespace
		routes := strings.Split(includedRoutesStr, ",")
		for _, route := range routes {
			trimmedRoute := strings.TrimSpace(route)
			if trimmedRoute != "" {
				includedRoutes = append(includedRoutes, trimmedRoute)
			}
		}
	}

	// Helper function to check if a path should be included
	shouldIncludePath := func(currentPath string) bool {
		// If no included routes specified, return false (don't include)
		if len(includedRoutes) == 0 {
			return false
		}

		// Check if current path matches any included route
		for _, includedRoute := range includedRoutes {
			if strings.HasPrefix(currentPath, includedRoute) || strings.Contains(currentPath, includedRoute) {
				return true
			}
		}
		return false
	}

	// setup custom logger
	var configDefault = logger.Config{
		Format: `{` +
			`"pid":"${pid}",` +
			`"time":"${time}",` +
			`"ip":"${ip}:${port}",` +
			`"host":"${host}",` +
			`"method":"${method}",` +
			`"path":"${path}",` +
			`"url":"${url}",` +
			`"ua":"${ua}",` +
			`"latency":"${latency}",` +
			`"status":${status},` +
			`"resBody":${resBodyV2},` +
			`"reqHeaders":"${reqHeadersV2}",` +
			`"queryParams":"${queryParamsV2}",` +
			`"reqBody":${reqBody},` +
			`"bytesSent":${bytesSent},` +
			`"bytesReceived":${bytesReceived},` +
			`"route":"${route}",` +
			`"lastErr":"${error}",` +
			`"requestId":"${locals:requestid}"` +
			`}` + "\n",
		TimeFormat:    time.RFC3339Nano,
		TimeZone:      "UTC",
		DisableColors: true,
		CustomTags: map[string]logger.LogFunc{
			"reqHeadersV2": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				// Check if current path should be included for request headers logging
				currentPath := c.Path()
				if !shouldIncludePath(currentPath) {
					return output.Write([]byte(""))
				}

				// Route is included, log the request headers
				reqHeaders := make([]string, 0)
				for k, v := range c.GetReqHeaders() {
					keyValueStr := k + "=" + strings.Join(v, ",")
					regEx := regexp.MustCompile(`["]`) // Use regex to remove all \n, \t
					reqHeadersCleaned := regEx.ReplaceAllString(string(keyValueStr), ``)
					reqHeaders = append(reqHeaders, reqHeadersCleaned)
				}
				return output.Write([]byte(strings.Join(reqHeaders, "&")))
			},
			"reqBody": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				reqBody := c.Body()
				regEx := regexp.MustCompile(`[\n\t]`) // Use regex to remove all \n, \t
				reqBodyCleaned := regEx.ReplaceAllString(string(reqBody), ``)
				reqBodyJsonStr, err := json.Marshal(reqBodyCleaned)
				if err != nil {
					return output.Write(reqBody)
				}
				return output.Write(reqBodyJsonStr)
			},
			"queryParamsV2": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				// Get raw query string and URL decode it
				rawQuery := string(c.Request().URI().QueryString())

				if rawQuery == "" {
					return output.Write([]byte(""))
				}

				// URL decode the query string
				decodedQuery, err := url.QueryUnescape(rawQuery)
				if err != nil {
					// If decoding fails, use original query string
					decodedQuery = rawQuery
				}
				return output.Write([]byte(decodedQuery))
			},
			"resBodyV2": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				// Check if current path should be included for response body logging
				currentPath := c.Path()
				if !shouldIncludePath(currentPath) {
					emptyStr, _ := json.Marshal("")
					return output.Write(emptyStr)
				}

				// Route is included, log the response body
				resBody := c.Response().Body()
				resBodyJsonStr, err := json.Marshal(string(resBody))
				if err != nil {
					return output.Write(resBody)
				}
				return output.Write(resBodyJsonStr)
			},
		},
	}

	a.Use(
		// Add CORS to each route.
		cors.New(),
		logger.New(configDefault),
		recover.New(),
	)

}
