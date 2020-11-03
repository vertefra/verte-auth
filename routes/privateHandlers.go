package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PrivateHandler defines all the private routes that will be used to make
// login and signup requests, and token verification from accounts
//
func PrivateHandler(r fiber.Router, db *gorm.DB) {

	r.Get("/", func(ctx *fiber.Ctx) error {
		data := ctx.Locals("tokenData")
		fmt.Println(data)
		ctx.SendString("received: ")
		return nil
	})
}
