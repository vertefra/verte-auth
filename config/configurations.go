package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	port = os.Getenv("PORT")
)

func (c *config) InitConfig() {
	godotenv.Load()

	// development configuration
	if c.ENV == "development" {
		c.PORT = os.Getenv("PORT")
		c.KEY = os.Getenv("KEY")
		c.DBNAME = os.Getenv("DEV_DB")
		c.USER = os.Getenv("DEV_USER")
	}

	if c.ENV == "production" {
		c.PORT = os.Getenv("PORT")
		c.KEY = os.Getenv("KEY")
		c.DBNAME = os.Getenv("DATABASE_URL")
		c.USER = os.Getenv("PROD_USER")
	}

}
