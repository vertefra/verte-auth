package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler Intercepts all the error error and format them in json
func ErrorHandler() func(ctx *fiber.Ctx, err error) error {

	return func(ctx *fiber.Ctx, err error) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		ctx.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
		return nil
	}
}

// currently not used
