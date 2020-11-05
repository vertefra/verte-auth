package main

import (
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/database"
	"github/vertefra/verte_auth_server/middleware"
	"github/vertefra/verte_auth_server/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	// templating engine for authentication views
	engine := django.New("./views", ".html")

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		Views:        engine,
	})

	// CHANGE TO true TO DROP THE DATABASE.
	if err := database.DropDb(false); err != nil {
		config.Err("was not able to drop the database")
		os.Exit(1)
	}

	// CHANGE TO true to CREATE A NEW DATABASE

	if err := database.CreateDBIfNotExists(false); err != nil {
		config.Err("was not able to create a new database ")
		os.Exit(1)
	}

	// initialize the automigration for the database
	db, err := database.ConnectDB(true)
	if err != nil {
		config.Err("\nFailed connecting database ")
		config.Err("Error: ", err)
		os.Exit(1)
	}

	// test route

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Api  running")
	})

	// static folder

	app.Static("/", "./public")

	// Routes groups

	// public
	usersAPI := app.Group("/api/users")

	accountAPI := app.Group("/api/users/:userID/accounts")

	privateAPI := app.Group("/private/accounts/:userID")

	// Routes Handlers
	routes.UserHandler(usersAPI, db)
	routes.AccountHandler(accountAPI, db)
	routes.PrivateHandler(privateAPI, db)

	config.Msg("\nServer listening on port: ", config.AppConfig.PORT)

	err = app.Listen(":" + config.AppConfig.PORT)

	if err != nil {
		config.Err("\nError while trying to listen on PORT: ", config.AppConfig.PORT)
		config.Err("Err: ", err)
		os.Exit(1)
	}

}
