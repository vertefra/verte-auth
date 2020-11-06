package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Err for displaying error in Red underlined color
var Err = color.New(color.FgHiRed).Add(color.Underline).Println

// Msg for displaying comunication messages
var Msg = color.New(color.FgHiGreen).Println

// Start define the environment variables of the application
func start() *config {

	var c *config

	loadEnv()

	// flags

	env := flag.String("env", "dev", "can be set 'dev' or 'prod' default 'dev'")
	db := flag.String("db", "automigrate", "can be set 'drop', 'create', 'automigrate'")
	dropTables := flag.Bool("dropTables", false, "set drop table to drop the existing tables")
	connProdDB := flag.Bool("connProdDB", false, "connect to production db from local")
	flag.Parse()

	Msg("\n--env 			= ", *env)

	// Evaluating environment
	if *env == "dev" {
		c = initConfig("development")

	} else if *env == "prod" {
		c = initConfig("production")

	} else if *env == "testDB" {
		c = initConfig("testProdDatabase")

	} else {
		Err("Environment not valid")
	}

	// Evaluating DB options
	if *db == "automigrate" {

		c.MIGRATE = true

	} else if *db == "drop" {

		if c.ENV == "production" {
			Err("Cannot drop database in production")
			os.Exit(1)
		}
		c.DROPDB = true

	} else if *db == "create" {

		if c.ENV == "production" {
			Err("Cannot create database in production")
			os.Exit(1)
		}

		c.CREATEDB = true
	}

	// evaluating table dropping

	if *dropTables == true {
		c.DROPTABLES = true
		c.MIGRATE = false
	}

	// evaluating remote connection with prodution db

	if *connProdDB == true {
		c.ENV = "testProdDatabase"
		c.DATABASEURL = os.Getenv("DATABASE_URL")
	}

	// if len(cmd) == 1 {
	// 	Err("No environment found, running development")
	// 	c.ENV = "development"
	// } else if cmd[1] == "dev" {
	// 	Msg("-- Requested environment: ", cmd[1])
	// 	c.ENV = "development"
	// 	c.InitConfig()
	// } else if cmd[1] == "prod" {
	// 	Msg("-- Requested environment: ", cmd[1])
	// 	c.ENV = "production"
	// 	c.InitConfig()
	// }
	// if env == "prod" {
	// 	Msg("Running production mode")
	// 	c.ENV = "production"
	// 	c.InitConfig()
	// }

	fmt.Println("config ", c)

	return c
}

// AppConfig is the config struct that will be shared by all the packages in the app
var AppConfig = start()
