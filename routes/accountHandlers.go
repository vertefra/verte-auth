package routes

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/models"
	"github/vertefra/verte_auth_server/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AccountHandler defines all the routes that start with /api/users/:userID/accounts
func AccountHandler(r fiber.Router, db *gorm.DB) {

	// @route	POST /api/users/:userID/accounts/signup
	// @desc	create a new authentication Account belonging to a user with userID
	// @key     only an key in the request body matching the one in stored with the user id will allow to use
	// 			this route

	r.Post("/signup", func(ctx *fiber.Ctx) error {

		queryID, err := strconv.ParseUint(ctx.Params("userID"), 10, 64)
		if err != nil {
			config.Err("Error: ", err)
			return err
		}

		// Conversion of userID to uint from uint64
		userID := uint(queryID)

		// retrieving api KEY from headers
		key := ctx.Get("key")

		// Retrieving account info from request body
		account, err := models.ParseAccount(ctx)
		if err != nil {
			return err
		}

		// Checking if username is a valid email
		if utils.IsEmailValid(account.Username) == false {
			err := fmt.Errorf("%v is not a valid email", account.Username)
			return err
		}

		// Ashing password for account
		hashedPass, err := utils.Encrypt(account.Password)
		if err != nil {
			return err
		}

		// Finding the User that owns the project for the account
		user, err := models.FindUserByID(db, userID)
		if err != nil {
			return err
		}

		// If the user.Key and the account.Key received in the
		// request body don't match exit with error
		if user.Key != key {
			err := fmt.Errorf("Api key does not match the one stored for user")
			return err
		}

		// creating a token for user authentication. exp time is set to 24 hrs for now
		t, err := utils.GenerateToken(account.Username, 24, user.Key)
		if err != nil {
			return err
		}

		// Assigning ashed password api key ownerID and token to the account
		account.Password = hashedPass
		// account.Key = user.Key
		account.OwnerID = user.ID
		account.Token = t

		// Checking if the user already has a account with same name
		unique, err := models.IsUnique(db, account)
		if err != nil {
			return err
		}

		if unique {

			if err := account.AddAccountToUser(db, userID); err != nil {
				return err
			}

			// Deleting the password
			account.Password = ""

			ctx.Status(201).JSON(fiber.Map{
				"sucess":         true,
				"createdAccount": account,
			})
			return nil
		}

		return err

	})

	// @route	POST /api/users/:userID/accounts/login
	// @desc	verify the Account for user with userID, checking username and password
	// @public

	r.Post("/login", func(ctx *fiber.Ctx) error {

		// Retrieving account info from request body
		account, err := models.ParseAccount(ctx)
		if err != nil {
			return err
		}

		// Saving Text Password to compare later
		p := account.Password

		// Get access key from headers

		// Checking if username is a valid email
		if utils.IsEmailValid(account.Username) == false {
			err := fmt.Errorf("%v is not a valid email", account.Username)
			return err
		}

		// Retrieving the account from the database and comparing the ashed passwords
		result := db.Table("accounts").Where(&models.Account{Username: account.Username}).First(&account)
		if result.Error != nil {
			err := fmt.Errorf("Cant find an account associated with %v", account.Username)
			result.AddError(err)
			return result.Error
		}

		if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(p)); err != nil {
			return err
		}

		// Generate a token
		// will contain in the payload the username of the account holder
		t, err := utils.GenerateToken(account.Username, 24, account.Key)
		if err != nil {
			return err
		}

		// should pull a redirect address from the database that is where you want
		// the user to go if succesfully authenticated

		ctx.JSON(fiber.Map{
			"success": true,
			"token":   t,
		})

		return nil
	})

	// @route	GET /api/users/:userID/accounts
	// @desc	gets all the credentials stored under the user with id userID
	// @token	You need to send the auth token in the request header

	r.Get("/", func(ctx *fiber.Ctx) error {

		userID, err := strconv.ParseUint(ctx.Params("userID"), 10, 64)
		if err != nil {
			ctx.Status(404)
			err = fmt.Errorf("ID sent in url probably incorrect - %v", err)
			return err
		}

		accounts := new([]models.Account)

		result := db.Table("accounts").Where(
			&models.Account{OwnerID: uint(userID)}).Find(accounts)

		if result.Error != nil {
			err := fmt.Errorf("Error Retrieving data from database - %v", result.Error)
			ctx.Status(404)
			return err
		}

		ctx.Status(201).JSON(fiber.Map{
			"success":  true,
			"accounts": accounts,
		})

		return nil
	})
}
