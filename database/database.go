package database

import (
	"context"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/models"
	"strings"

	"github.com/jackc/pgx/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//ConnectDB provide connection with the database specified in the config varibales
func ConnectDB(automigrate bool) (*gorm.DB, error) {

	// Retrieving the db name from config file
	dbName := config.AppConfig.DBNAME

	// creating url connection string
	url := "postgresql:///" + dbName

	// Opening database
	DB, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		config.Err("\nError while opening database: ", dbName)
		config.Err("Error: ", err)
		return nil, err
	}

	config.Msg("database " + dbName + "connected!")

	if automigrate == true {
		config.Msg("\nRunning auto migration")

		// need to make this dynamic
		if err := DB.AutoMigrate(&models.Account{}, &models.User{}); err != nil {
			config.Err("\nError while executing automigration")
			config.Err("Error: ", err)
			return nil, err
		}
	}

	return DB, nil

}

func CreateDBIfNotExists(confirm bool) error {

	if confirm {
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
			config.Msg("Database %s created ", dbName)
		} else if exist == true {
			config.Msg("\nFound a db with name: ", dbName)
		}

		conn.Close(context.Background())

		return nil
	} else {
		return nil
	}
}

// DropDb executes the DROP DATABASE sql command with the db name given
func DropDb(confirm bool) error {

	if confirm {
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
		return nil
	} else {
		return nil
	}
}
