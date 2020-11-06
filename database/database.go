package database

import (
	"context"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/models"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//ConnectDB provide connection with the database specified in the config varibales
func ConnectDB(automigrate bool) (*gorm.DB, error) {

	// creating url connection string
	var url string

	if config.AppConfig.ENV == "development" {
		url = "postgresql:///" + config.AppConfig.DBNAME
	} else if config.AppConfig.ENV == "testProdDatabase" {
		url = config.AppConfig.DATABASEURL
	} else {
		url = config.AppConfig.DBNAME
	}

	config.Msg("connecting to => ", url)
	config.Msg("KEY ", config.AppConfig.KEY)

	// url = config.AppConfig.DBNAME

	// Opening database
	DB, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		config.Err("\nError while opening database: ", config.AppConfig.DBNAME)
		config.Err("Error: ", err)
		return nil, err
	}

	config.Msg("-- database connected!")

	if automigrate == true {
		config.Msg("\n-- Running auto migration")

		// need to make this dynamic
		if err := DB.AutoMigrate(&models.Account{}, &models.User{}); err != nil {
			config.Err("\nError while executing automigration")
			config.Err("Error: ", err)
			return nil, err
		}
	}

	return DB, nil

}

// DropTables drops the edxisting tables
func DropTables() error {

	var url string

	if config.AppConfig.ENV == "development" {
		url = "postgresql:///" + config.AppConfig.DBNAME
	} else if config.AppConfig.ENV == "testProdDatabase" {
		url = config.AppConfig.DATABASEURL

	} else {
		url = config.AppConfig.DBNAME

	}

	DB, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		config.Err("\nError while opening database: ", config.AppConfig.DBNAME)
		config.Err("Error: ", err)
		return err
	}

	if err := DB.Migrator().DropTable(&models.User{}, &models.Account{}); err != nil {
		config.Err("\nError dropping tables")
		config.Err("Error: ", err)
		os.Exit(1)
	}
	config.Msg("\nTables dropped!\n")
	os.Exit(0)
	return nil
}

// CreateDBIfNotExists will create a new databse if not already present
// DEVELOPMENT ONLY!!
func CreateDBIfNotExists() error {

	dbName := config.AppConfig.DBNAME
	user := config.AppConfig.USER

	conn, err := pgx.Connect(context.Background(), "postgresql:///"+user)

	if err != nil {
		config.Err("\nCan't connect to ", dbName)
		config.Err("Error: ", err)
		return err
	}

	check := strings.TrimSpace("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '" + dbName + "');")

	row := conn.QueryRow(context.Background(), check)

	var exist bool
	err = row.Scan(&exist)

	if err != nil {
		config.Err("error during scanning: ", err)
		return err
	}

	if exist == false {
		sqlCmd := "CREATE DATABASE " + dbName + ";"
		_, err := conn.Exec(context.Background(), sqlCmd)
		if err != nil {
			config.Err("Error creating DB: ", err)
			return err
		}
		config.Msg("Database created: ", dbName)
	} else if exist == true {
		config.Msg("\nFound a db with name: ", dbName)
	}

	conn.Close(context.Background())
	config.Msg("Database Created")
	os.Exit(0)
	return nil
}

// DropDb executes the DROP DATABASE sql command with the db name given
// ONLY DEVELOPMENT!!
func DropDb() error {

	dbName := config.AppConfig.DBNAME
	user := config.AppConfig.USER

	conn, err := pgx.Connect(context.Background(), "postgresql:///"+user)

	if err != nil {
		config.Err("\nCan't connect to ", dbName)
		config.Err("Error: ", err)
		return err
	}
	sqlCmd := "DROP DATABASE " + dbName + ";"

	_, err = conn.Exec(context.Background(), sqlCmd)

	if err != nil {
		config.Err("Error dropping database ", dbName)
		config.Err("Error: ", err)
		return err
	}
	config.Err("Your database has been erased")
	os.Exit(0)
	return nil
}
