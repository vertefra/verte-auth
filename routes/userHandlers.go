package routes

import (
	"fmt"
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
			return err
		}

		// Checking if user is already in the database
		if user.IsPresent(db) {
			err := fmt.Errorf("%s is already in the database", user.Email)
			return err
		}

		// Ashing password
		hashedPass, err := utils.Encrypt(user.Password)
		if err != nil {
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
			return err
		}

		// Saving text password for later comparing
		p := user.Password

		result := db.Table("users").Where(&models.User{Email: user.Email}).First(user)
		if result.Error != nil {
			return result.Error
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p)); err != nil {
			return err
		}

		// Passwords match, generate auth token

		t, err := utils.GenerateToken(user.Email, 24, user.Key)
		if err != nil {
			return err
		}

		ctx.JSON(fiber.Map{
			"success": true,
			"token":   t,
			"key":     user.Key,
			"userId":  user.ID,
		})
		return nil
	})

	r.Get("/", func(ctx *fiber.Ctx) error {
		ctx.SendString("Hi user")
		return nil
	})
}
