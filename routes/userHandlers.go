package routes

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/middleware"
	"github/vertefra/verte_auth_server/models"
	"github/vertefra/verte_auth_server/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserHandler defines all the routes that start with /api/users
func UserHandler(r fiber.Router, db *gorm.DB) {

	// REMEMBER!! I'm passing gorm.DB as pointer because it's a struct and
	// structs pass their value that is copied in the function in a
	// new memory position, while fiber.Router is an interface that pass
	// just the reference to the values so you don't need to pass a pointer
	// beacause it's already a pointer

	r.Post("/signup", func(ctx *fiber.Ctx) error {

		user, err := models.ParseUser(ctx)
		if err != nil {
			ctx.Status(404)
			return err
		}

		// Checking if user is already in the database
		if user.IsPresent(db) {
			err := fmt.Errorf("%s is already in the database", user.Email)
			ctx.Status(404)
			return err
		}

		// Ashing password
		hashedPass, err := utils.Encrypt(user.Password)
		if err != nil {
			ctx.Status(500)
			return err
		}

		// Creating API key
		key := utils.GenerateRanKey(8)

		// Assigning crypted password and api Key
		user.Password = hashedPass
		user.Key = key

		// Creating a new user. An evenual error is stored inside
		// result
		result := db.Create(user)
		if result.Error != nil {
			ctx.Status(500)
			return result.Error
		}

		// Deleting password from the result that will be returned
		user.Password = ""
		ctx.Status(201).JSON(fiber.Map{
			"success":     true,
			"createdUser": user,
		})

		return nil
	})

	r.Post("/login", func(ctx *fiber.Ctx) error {

		user, err := models.ParseUser(ctx)
		if err != nil {
			ctx.Status(404)
			return err
		}

		// Saving text password for later comparing
		p := user.Password

		result := db.Table("users").Where(&models.User{Email: user.Email}).First(user)
		if result.Error != nil {
			ctx.Status(500)
			return result.Error
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p)); err != nil {
			ctx.Status(404)
			return err
		}

		// Passwords match, generate auth token

		// For the user that subscribe to the service the key used to
		// sign the toke is the one stored in the .env file, while for the
		// authentication service the key is the randomly generated one
		// that will be sent in the header requests

		tokenKey := config.AppConfig.KEY

		t, err := utils.GenerateToken(user.Email, 24, tokenKey)
		if err != nil {
			ctx.Status(500)
			return err
		}

		ctx.JSON(fiber.Map{
			"success":     true,
			"token":       t,
			"key":         user.Key,
			"ID":          user.ID,
			"email":       user.Email,
			"redirectURL": user.RedirectURL,
		})
		return nil
	})

	// @desc	update user info
	// @route	PUT	/api/users/:id
	// @private	token key is in .env

	r.Put("/:userID", middleware.UserAuth(), func(ctx *fiber.Ctx) error {
		user := new(models.User)
		if err := ctx.BodyParser(user); err != nil {
			ctx.Status(404)
			return err
		}

		if ok, err := models.UpdateUser(db, user); !ok {
			config.Err("Error while updating user")
			ctx.Status(500)
			return err
		}

		return ctx.JSON(fiber.Map{
			"success": true,
			"user":    user,
		})

	})
}
