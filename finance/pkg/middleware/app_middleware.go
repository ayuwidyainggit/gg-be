package middleware

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AppMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func AppMiddleware(a *fiber.App) {
	var ConfigDefault = requestid.Config{
		Next:   nil,
		Header: fiber.HeaderXRequestID,
		Generator: func() string {
			objectID := primitive.NewObjectID() // Generate a new ObjectID
			objectIDString := objectID.Hex()    // Convert ObjectID to string
			return objectIDString
		},
		ContextKey: "requestid",
	}
	a.Use(requestid.New(ConfigDefault))

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
			`"queryParams":"${queryParams}",` +
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
			"resBodyV2": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
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
