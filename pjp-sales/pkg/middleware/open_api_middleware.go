package middleware

import (
	"strings"

	"sales/entity"
	"sales/pkg/constant"
	"sales/service"

	"github.com/gofiber/fiber/v2"
)

func OpenAPIProtected(openAPIService service.OpenAPIService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientID := strings.TrimSpace(c.Get(constant.HeaderOpenAPIClientID))
		clientSecret := strings.TrimSpace(c.Get(constant.HeaderOpenAPIClientSecret))
		if clientID == "" || clientSecret == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(entity.ApiResponse{
				Message: constant.MsgOpenAPIMissingCredentials,
			})
		}

		custID := strings.TrimSpace(c.Get(constant.HeaderOpenAPICustID))
		if custID == "" {
			custID = strings.TrimSpace(c.Query("cust_id"))
		}

		authCtx, err := openAPIService.Authenticate(
			clientID,
			clientSecret,
			c.Method(),
			c.Path(),
			custID,
		)
		if err != nil {
			status := fiber.StatusUnauthorized
			msg := constant.MsgOpenAPIUnauthorized
			switch err {
			case service.ErrOpenAPIForbidden:
				status = fiber.StatusForbidden
				msg = constant.MsgOpenAPIForbidden
			case service.ErrOpenAPIEndpointNotFound:
				status = fiber.StatusNotFound
				msg = constant.MsgOpenAPIEndpointNotAllowed
			}
			return c.Status(status).JSON(entity.ApiResponse{Message: msg})
		}

		c.Locals("open_api_config_id", authCtx.ConfigID)
		c.Locals("open_api_system_integration", authCtx.SystemIntegration)
		c.Locals("open_api_cust_id", authCtx.CustID)
		c.Locals("parent_cust_id", authCtx.CustID)
		c.Locals("cust_id", authCtx.CustID)
		c.Locals("user_fullname", authCtx.SystemIntegration)
		c.Locals("user_name", authCtx.SystemIntegration)

		return c.Next()
	}
}
