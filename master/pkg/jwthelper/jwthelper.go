package jwthelper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"master/entity"
	"master/model"
	"master/pkg/config/env"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateNewTokens func for generate a new Access & Refresh tokens.
func GenerateNewToken(user model.User, credentials []string, additionalParam map[string]interface{}) (*entity.Token, error) {
	// Generate JWT Access token.
	accessToken, err := generateNewAccessToken(user, credentials, additionalParam)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	// Generate JWT Refresh token.
	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	return &entity.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*entity.TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		var distPriceGrpId int
		var distributorId int64
		userId := int64(claims["user_id"].(float64))
		userName := claims["user_name"].(string)
		userFullName := claims["user_fullname"].(string)
		email := claims["email"].(string)
		isAdmin := claims["is_admin"].(bool)
		custId := claims["cust_id"].(string)
		parentCustId := claims["parent_cust_id"].(string)
		empId := int64(claims["emp_id"].(float64))
		if claims["distributor_id"] != nil {
			distributorId = int64(claims["distributor_id"].(float64))
		}
		if claims["dist_price_grp_id"] != nil {
			distPriceGrpId = int(claims["dist_price_grp_id"].(float64))
		}
		langId := claims["lang_id"].(string)
		mobileNo := claims["mobile_no"].(string)
		whatsapp := claims["whatsapp"].(string)

		// Expires time.
		expires := int64(claims["expires"].(float64))

		tokenMetadata := &entity.TokenMetadata{
			UserId:         userId,
			UserName:       userName,
			UserFullName:   userFullName,
			Email:          email,
			IsAdmin:        isAdmin,
			CustId:         custId,
			ParentCustId:   parentCustId,
			DistPriceGrpId: distPriceGrpId,
			EmpId:          empId,
			LangId:         langId,
			MobileNo:       mobileNo,
			Whatsapp:       whatsapp,
			Expires:        expires,
			DistributorId:  distributorId,
		}
		// log.Println("tokenMetadata:", tokenMetadata)
		return tokenMetadata, nil
	}

	return nil, err
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

func generateNewAccessToken(user model.User, credentials []string, additionalParam map[string]interface{}) (string, error) {
	envCfg := env.NewCfgEnv()

	// Set secret key from .env file.
	secret := envCfg.Get("JWT_SECRET_KEY")

	// Set expires minutes count for secret key from .env file.
	minutesCount, _ := strconv.Atoi(envCfg.Get("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	additionalParamPrint, _ := json.MarshalIndent(additionalParam, "", "\t")
	log.Println("### additionalParamPrint ###")
	log.Println(string(additionalParamPrint))
	log.Println("### End Of additionalParamPrint ###")

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["user_id"] = user.UserId
	claims["user_name"] = user.Username
	claims["user_fullname"] = user.Fullname
	claims["email"] = user.Email
	claims["mobile_no"] = user.MobileNo
	claims["whatsapp"] = user.Whatsapp
	claims["lang_id"] = user.LangId
	claims["cust_id"] = user.CustId
	claims["emp_id"] = user.EmpId
	claims["is_admin"] = user.IsAdmin
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	// claims["book:create"] = false
	// claims["book:update"] = false
	// claims["book:delete"] = false

	// Set private token credentials:
	for _, credential := range credentials {
		claims[credential] = true
	}

	if additionalParam != nil {
		for key, value := range additionalParam {
			claims[key] = value
		}
	}

	claimsPrint, _ := json.MarshalIndent(claims, "", "\t")
	log.Println("### claimsPrint ###")
	log.Println(string(claimsPrint))
	log.Println("### End Of claimsPrint ###")

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenPrint, _ := json.MarshalIndent(token, "", "\t")
	log.Println("### tokenPrint ###")
	log.Println(string(tokenPrint))
	log.Println("### End Of tokenPrint ###")

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
