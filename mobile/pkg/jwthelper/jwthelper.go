package jwthelper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"mobile/entity"
	"mobile/pkg/config/env"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims defines the structure for JWT claims with type safety.
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID               int64  `json:"user_id"`
	CustID               string `json:"cust_id"`
	ParentCustID         string `json:"parent_cust_id"`
	DistPriceGrpID       int    `json:"dist_price_grp_id"`
	UserRole             string `json:"user_role"`
	Email                string `json:"email"`
	EmpID                int64  `json:"emp_id"`
	EmpCode              string `json:"emp_code"`
	LangID               string `json:"lang_id"`
	EmpGrpID             int64  `json:"emp_grp_id"`
	MobileNo             string `json:"mobile_no"`
	Whatsapp             string `json:"whatsapp"`
	OprTypeCanvas        string `json:"opr_type_canvas"`
	OprTypeOrderTaking   string `json:"opr_type_order_taking"`
	AllowInputPrice      bool   `json:"allow_input_price"`
	TaxOption            string `json:"tax_option"`
	IsActiveGudangCanvas bool   `json:"is_active_gudang_canvas"`
	IsActiveGudangUtama  bool   `json:"is_active_gudang_utama"`
	DistributorID        int64  `json:"distributor_id"`
	Expires              int64  `json:"expires"`
	// Credentials stores dynamic permission flags
	Credentials  map[string]bool `json:"credentials,omitempty"`
	Username     string          `json:"user_name"`
	UserFullname string          `json:"user_fullname"`
	IsAdmin      bool            `json:"is_admin"`
}

// GenerateNewTokens func for generate a new Access & Refresh tokens.
func GenerateNewToken(user entity.UserData, credentials []string, additionalParam map[string]interface{}) (string, error) {
	// Generate JWT Access token.
	accessToken, err := generateNewAccessToken(user, credentials, additionalParam)
	if err != nil {
		// Return token generation error.
		return "", err
	}

	// Generate JWT Refresh token.
	// refreshToken, err := generateNewRefreshToken()
	// if err != nil {
	// 	// Return token generation error.
	// 	return nil, err
	// }

	return accessToken, err
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*entity.TokenMetadata, error) {
	tokenString := extractToken(c)
	if tokenString == "" {
		return nil, fmt.Errorf("missing authorization token")
	}

	// Parse with custom claims struct
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Access claims directly as struct fields
	return &entity.TokenMetadata{
		UserId:               claims.UserID,
		UserRole:             claims.UserRole,
		Email:                claims.Email,
		CustId:               claims.CustID,
		EmpId:                claims.EmpID,
		EmpCode:              claims.EmpCode,
		EmpGrpId:             claims.EmpGrpID,
		LangId:               claims.LangID,
		MobileNo:             claims.MobileNo,
		Whatsapp:             claims.Whatsapp,
		Expires:              claims.Expires,
		ParentCustId:         claims.ParentCustID,
		DistPriceGrpId:       claims.DistPriceGrpID,
		OprTypeCanvas:        claims.OprTypeCanvas,
		OprTypeOrderTaking:   claims.OprTypeOrderTaking,
		AllowInputPrice:      claims.AllowInputPrice,
		TaxOption:            claims.TaxOption,
		IsActiveGudangCanvas: claims.IsActiveGudangCanvas,
		IsActiveGudangUtama:  claims.IsActiveGudangUtama,
		DistributorID:        claims.DistributorID,
	}, nil
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}

func generateNewAccessToken(user entity.UserData, credentials []string, additionalParam map[string]interface{}) (string, error) {
	envCfg := env.NewCfgEnv()

	// Set secret key from .env file.
	secret := envCfg.Get("JWT_SECRET_KEY")

	// Set expires minutes count for secret key from .env file.
	minutesCount, _ := strconv.Atoi(envCfg.Get("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))
	expiresAt := time.Now().Add(time.Minute * time.Duration(minutesCount))

	// Build credentials map
	credentialsMap := make(map[string]bool)
	for _, credential := range credentials {
		credentialsMap[credential] = true
	}

	// Safely dereference pointer fields with defaults
	var empID int64
	if user.EmpId != nil {
		empID = *user.EmpId
	}
	var empCode string
	if user.EmpCode != nil {
		empCode = *user.EmpCode
	}
	var empGrpID int64
	if user.EmpGrpId != nil {
		empGrpID = *user.EmpGrpId
	}

	// Create claims using struct
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:               user.UserId,
		UserRole:             user.UserRole,
		Username:             user.Username,
		UserFullname:         user.UserFullname,
		IsAdmin:              user.IsAdmin,
		Email:                user.Email,
		MobileNo:             user.MobileNo,
		Whatsapp:             user.Whatsapp,
		LangID:               user.LangId,
		CustID:               user.CustId,
		ParentCustID:         user.ParentCustId,
		DistPriceGrpID:       user.DistPriceGrpId,
		EmpID:                empID,
		EmpCode:              empCode,
		EmpGrpID:             empGrpID,
		Expires:              expiresAt.Unix(),
		OprTypeCanvas:        user.OprTypeCanvas,
		OprTypeOrderTaking:   user.OprTypeOrderTaking,
		AllowInputPrice:      user.AllowInputPrice,
		TaxOption:            user.TaxOption,
		IsActiveGudangCanvas: user.IsActiveGudangCanvas,
		IsActiveGudangUtama:  user.IsActiveGudangUtama,
		DistributorID:        int64(user.DistributorID),
		Credentials:          credentialsMap,
	}

	claimsPrint, _ := json.MarshalIndent(claims, "", "\t")
	log.Println("### claimsPrint ###")
	log.Println(string(claimsPrint))
	log.Println("### End Of claimsPrint ###")

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

func generateNewRefreshToken() (string, error) {
	// Create a new SHA256 hash.
	hash := sha256.New()

	// Create a new now date and time string with salt.
	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()

	// See: https://pkg.go.dev/io#Writer.Write
	_, err := hash.Write([]byte(refresh))
	if err != nil {
		// Return error, it refresh token generation failed.
		return "", err
	}

	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	// Set expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hex.EncodeToString(hash.Sum(nil)) + "." + expireTime

	return t, nil
}

// ParseRefreshToken func for parse second argument from refresh token.
func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
