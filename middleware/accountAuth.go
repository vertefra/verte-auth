package middleware

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// AccessDetails access the id encoded in the token
type AccessDetails struct {
	ID string
}

// ExtractTokenMetadata verify the token and returns the encoded metadata
func ExtractTokenMetadata(ctx *fiber.Ctx) (*AccessDetails, error) {

	if err := tokenValid(ctx); err != nil {
		return nil, err
	}

	token, err := verifyToken(ctx)
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

func tokenValid(ctx *fiber.Ctx) error {
	token, err := verifyToken(ctx)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil

}

func verifyToken(ctx *fiber.Ctx) (*jwt.Token, error) {

	key := ctx.Get("key")

	if key == "" {
		err := fmt.Errorf("Error key not present")
		return nil, err
	}

	tokenString := extractToken(ctx)
	config.Err(tokenString)
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
