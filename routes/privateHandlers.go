package routes

import (
	"fmt"
	"github/vertefra/verte_auth_server/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PrivateHandler defines all the private routes that will be used to make
// login and signup requests, and token verification from accounts
//
func PrivateHandler(r fiber.Router, db *gorm.DB) {

	r.Get("/", middleware.PrivateAccount(db), func(ctx *fiber.Ctx) error {

		// for now access details does not contains any meainginful information
		// just the key passed from the middleware
		data := ctx.Locals("tokenData").(*middleware.AccessDetails)
		fmt.Println(data)
		ctx.SendString("received: " + string(data.ID))
		return nil
	})
}
