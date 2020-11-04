package middleware

import (
	"github/vertefra/verte_auth_server/config"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler Intercepts all the error error and format them in json
func ErrorHandler(ctx fiber.Ctx, err error) error {

	code := 500
	if e, ok := err.(*fiber.Error); ok {
		config.Msg(e)
		code = e.Code
	}
	return ctx.Status(code).JSON(fiber.Map{
		"error": "this is the error",
	})
}

// currently not used
