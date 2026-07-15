package middleware

import (
	"inventory/entity"
	"inventory/pkg/config/env"
	"inventory/pkg/jwthelper"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/jwt/v2"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(ctx *fiber.Ctx) error {
	envCfg := env.NewCfgEnv()

	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:     []byte(envCfg.Get("JWT_SECRET_KEY")),
		ContextKey:     "jwt", // used in private routes
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccess,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(entity.ApiResponse{
			Message: err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(entity.ApiResponse{
		Message: fiber.ErrUnauthorized.Message,
	})
}

func jwtSuccess(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := jwthelper.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(entity.ApiResponse{
			Message: err.Error(),
		})
	}

	c.Locals("cust_id", claims.CustId)
	c.Locals("parent_cust_id", claims.ParentCustId)
	c.Locals("user_id", claims.UserId)
	c.Locals("user_name", claims.UserName)
	c.Locals("user_fullname", claims.UserFullName)
	c.Locals("user_email", claims.Email)
	c.Locals("is_admin", claims.IsAdmin)
	c.Locals("employee_id", claims.EmpId)
	c.Locals("mobile_no", claims.MobileNo)
	c.Locals("user_lang", claims.LangId)
	c.Locals("dist_price_grp_id", claims.DistPriceGrpId)
	c.Locals("distributor_id", claims.DistributorID)

	if int64(claims.Expires) < time.Now().UTC().Unix() {
		// return nil, errors.New("jwt is expired")
		return c.Status(fiber.StatusUnauthorized).JSON(entity.ApiResponse{
			Message: "jwt is expired",
		})
	}

	return c.Next()
}
