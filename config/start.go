package config

import (
	"os"

	"github.com/fatih/color"
)

type config struct {
	ENV    string
	PORT   string
	DBNAME string
	USER   string
	KEY    string
}

// Err for displaying error in Red underlined color
var Err = color.New(color.FgHiRed).Add(color.Underline).Println

// Msg for displaying comunication messages
var Msg = color.New(color.FgHiGreen).Println

// Start define the environment variables of the application
func start(env string) config {
	var c config
	cmd := os.Args

	if len(cmd) == 1 {
		Err("No environment found, running development")
		c.ENV = "development"
	} else if cmd[1] == "dev" {
		Msg("-- Requested environment: ", cmd[1])
		c.ENV = "development"
		c.InitConfig()
	} else if cmd[1] == "prod" {
		Msg("-- Requested environment: ", cmd[1])
		c.ENV = "production"
		c.InitConfig()
	} else if env == "prod" {
		Msg("Running production mode")
		c.ENV = "production"
		c.InitConfig()
	}

	Msg("-- Starting env    => ", c.ENV)
	Msg("-- Setting port to => ", c.PORT)

	return c
}

// AppConfig is the config struct that will be shared by all the packages in the app
var AppConfig = start("prod")
