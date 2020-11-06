package middleware

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/models"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	fiber "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AccessDetails access the id encoded in the token
type AccessDetails struct {
	ID string
}

// UserAuth is the middleware that checks that the incoming request
// has the correct token signed with the secret key belonging to the
// system
func UserAuth() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		key := config.AppConfig.KEY

		config.Err("KEY found: ", key)
		data, err := ExtractTokenMetadata(ctx, key)
		if err != nil {
			config.Err("Error in UserAuth")
			config.Err(err)
			ctx.Status(404)
			return err
		}
		ctx.Locals("tokenData", data)
		return ctx.Next()
	}
}

// PrivateAccount Verify the token for all the user connected to a project. the key is extracted
// getting the owner of the project id from the headers and getting the key from the
// database
func PrivateAccount(db *gorm.DB) func(ctx *fiber.Ctx) error {

	return func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseUint(ctx.Get("ownerID"), 10, 64)
		if err != nil {
			ctx.Status(500)
			return err
		}

		key := models.Account{}.Key

		result := db.Table("accounts").Where(
			&models.Account{OwnerID: uint(id)}).Select(
			"key").Find(&key)
		if result.Error != nil {
			ctx.Status(404)
			return result.Error
		}

		// now that we have the key we can verify the token

		tokenData, err := ExtractTokenMetadata(ctx, key)
		if err != nil {
			// not authorized
			ctx.Status(404)
			return err
		}

		ctx.Locals("tokenData", tokenData)

		return ctx.Next()
	}
}

// Helper functions to manipulate token and extract token data

// ExtractTokenMetadata verify the token and returns the encoded metadata
func ExtractTokenMetadata(ctx *fiber.Ctx, key string) (*AccessDetails, error) {

	if err := tokenValid(ctx, key); err != nil {
		return nil, err
	}

	token, err := verifyToken(ctx, key)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, ok := claims["ID"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			ID: id,
		}, nil
	}
	return nil, err
}

func tokenValid(ctx *fiber.Ctx, key string) error {
	token, err := verifyToken(ctx, key)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil

}

func verifyToken(ctx *fiber.Ctx, key string) (*jwt.Token, error) {

	if key == "" {
		err := fmt.Errorf("Error key not present")
		return nil, err
	}

	tokenString := extractToken(ctx)
	token, err := jwt.Parse(
		tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signin method")
			}
			return []byte(key), nil
		})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// extractToken extract the token from the header and verify it
func extractToken(ctx *fiber.Ctx) string {
	bearerToken := ctx.Get("Authorization")
	// separates "Bearer" from token
	sliceToken := strings.Split(bearerToken, " ")
	if len(sliceToken) == 2 {
		return sliceToken[1]
	}
	return ""
}
