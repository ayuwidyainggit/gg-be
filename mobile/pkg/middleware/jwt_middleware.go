package middleware

import (
	"mobile/pkg/config/env"
	"mobile/pkg/jwthelper"
	"time"

	"mobile/entity"

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
	c.Locals("user_id", claims.UserId)
	c.Locals("email", claims.Email)
	c.Locals("emp_id", claims.EmpId)
	c.Locals("emp_code", claims.EmpCode)
	c.Locals("mobile_no", claims.MobileNo)
	c.Locals("user_lang", claims.LangId)
	c.Locals("parent_cust_id", claims.ParentCustId)
	c.Locals("dist_price_grp_id", claims.DistPriceGrpId)
	c.Locals("emp_grp_id", claims.EmpGrpId)
	c.Locals("is_active_gudang_utama", claims.IsActiveGudangUtama)
	c.Locals("is_active_gudang_canvas", claims.IsActiveGudangCanvas)
	c.Locals("distributor_id", claims.DistributorID)
	c.Locals("distributor_code", claims.DistributorCode)
	c.Locals("is_distributor", func() bool { return claims.DistributorID > 0 }())

	if int64(claims.Expires) < time.Now().UTC().Unix() {
		// return nil, errors.New("jwt is expired")
		return c.Status(fiber.StatusUnauthorized).JSON(entity.ApiResponse{
			Message: "jwt is expired",
		})
	}

	return c.Next()
}
