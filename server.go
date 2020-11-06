package main

import (
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/database"
	"github/vertefra/verte_auth_server/middleware"
	"github/vertefra/verte_auth_server/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/django"
)

func main() {

	// templating engine for authentication views
	engine := django.New("./views", ".html")

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		Views:        engine,
	})

	config.Msg("\nverifing database options: \n")
	config.Msg("- createDb:		", config.AppConfig.CREATEDB)
	config.Msg("- dropDb: 		", config.AppConfig.DROPDB)
	config.Msg("- migrateTables:", config.AppConfig.MIGRATE)
	config.Msg("- dropTables: 	", config.AppConfig.DROPTABLES)
	config.Msg("- system Key	", config.AppConfig.KEY)

	if config.AppConfig.DROPDB {
		if err := database.DropDb(); err != nil {
			config.Err("was not able to drop the database")
			os.Exit(1)
		}
	}

	// Create new database
	if config.AppConfig.CREATEDB {
		if err := database.CreateDBIfNotExists(); err != nil {
			config.Err("\nwas not able to create a new database ")
			os.Exit(1)
		}
	}

	// Drop tables if flag dropTables is called

	if config.AppConfig.DROPTABLES {
		if err := database.DropTables(); err != nil {
			config.Err("Cannot drop tables")
			os.Exit(1)
		}
	}

	// Initialize the automigration for the database
	db, err := database.ConnectDB(config.AppConfig.MIGRATE)
	if err != nil {
		config.Err("\nFailed connecting database ")
		config.Err("Error: ", err)
		os.Exit(1)
	}

	// test route

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Api running")
	})

	// REMEMBER!!! SET REAL CORSE WHEN READY FOR PRODUCTION

	app.Use(cors.New())

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
