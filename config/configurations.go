package config

import (
	"os"
)

var (
	port = os.Getenv("PORT")
)

type config struct {
	ENV        string
	PORT       string
	DBNAME     string
	USER       string
	KEY        string
	DROPDB     bool
	CREATEDB   bool
	MIGRATE    bool
	DROPTABLES bool
}

var devConfig = config{
	ENV:        "development",
	PORT:       "4999",
	DBNAME:     "dev_verte_auth_db",
	USER:       "verte",
	KEY:        os.Getenv("KEY"),
	DROPDB:     false,
	CREATEDB:   false,
	MIGRATE:    true,
	DROPTABLES: false,
}

var prodConfig = config{
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
