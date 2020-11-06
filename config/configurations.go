package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	port = os.Getenv("PORT")
)

type config struct {
	ENV         string
	PORT        string
	DBNAME      string
	USER        string
	KEY         string
	DROPDB      bool
	CREATEDB    bool
	MIGRATE     bool
	DROPTABLES  bool
	DATABASEURL string
}

func initConfig(c string) *config {

	loadEnv()

	var devConfig = &config{
		ENV:         "development",
		PORT:        "4999",
		DBNAME:      "dev_verte_auth_db",
		USER:        "verte",
		KEY:         os.Getenv("KEY"),
		DROPDB:      false,
		CREATEDB:    false,
		MIGRATE:     true,
		DROPTABLES:  false,
		DATABASEURL: "",
	}

	var prodConfig = &config{
		ENV:        "production",
		PORT:       os.Getenv("PORT"),
		DBNAME:     os.Getenv("DATABASE_URL"),
		USER:       os.Getenv("PROD_USER"),
		KEY:        os.Getenv("KEY"),
		DROPDB:     false,
		CREATEDB:   false,
		MIGRATE:    true,
		DROPTABLES: false,
	}

	var prodDatabaseTest = &config{
		ENV:         "testProdDatabase",
		PORT:        "4999",
		DBNAME:      "testing real db",
		USER:        os.Getenv("PROD_USER"),
		KEY:         os.Getenv("KEY"),
		DROPDB:      false,
		CREATEDB:    false,
		MIGRATE:     false,
		DROPTABLES:  false,
		DATABASEURL: os.Getenv("DATABASE_URL"),
	}

	if c == "development" {
		return devConfig
	}

	if c == "production" {
		return prodConfig
	}

	if c == "testProdDatabase" {
		return prodDatabaseTest
	}

	return nil
}

func loadEnv() {
	godotenv.Load()
}
