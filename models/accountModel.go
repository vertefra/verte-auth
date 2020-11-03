package models

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"

	fiber "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Account is thje struct with all the authorization fields for a User
type Account struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	OwnerID  uint   `json:"ownerID"`
	Key      string `json:"key"`
	Token    string `json:"token"`
}

// ParseAccount returns the parsed body with Account struct
func ParseAccount(ctx *fiber.Ctx) (*Account, error) {
	Account := new(Account)
	if err := ctx.BodyParser(Account); err != nil {
		config.Err("Error while parsing the request body")
		config.Err("Error: ", err)
		return nil, err
	}

	if Account.Username == "" {
		err := fmt.Errorf("At model.ParseError: Account cannot have empty username or password")
		config.Err(err)
		return nil, err
	}
	return Account, nil
}

// IsUnique check if in the database there is already a Account with same
// username belonging to the same owner
func IsUnique(db *gorm.DB, a *Account) (bool, error) {
	foundAccount := new(Account)
	result := db.Table("accounts").Where(
		&Account{Username: a.Username, OwnerID: a.OwnerID}).First(&foundAccount)
	if result.Error != nil {
		// No matching Account has been found the Account is unique
		config.Err("isUnique => ", result.Error)
		return true, nil
	}
	err := fmt.Errorf("Found another Account belonging to OwnerID: %v with username: %s .\n An owner cannot have two Accounts with same username", a.OwnerID, a.Username)
	return false, err
}

// AddAccountToUser get the validate Account structs coming from a ParseAccount function and create
// a new entry in the database adding UserID from the user id in the url
func (a *Account) AddAccountToUser(db *gorm.DB, userID uint) error {

	// Check if the userID exists in database
	user, err := FindUserByID(db, userID)
	if err != nil {
		config.Err("Error In: AddAccountToUser")
		config.Err("Error: ", err)
		return err
	}
	// Assigning the userID and the api key to the Account
	a.OwnerID = user.ID
	a.Key = user.Key

	// Creating a new Account
	result := db.Create(a)
	if result.Error != nil {
		config.Err("Error in: CreateAccount")
		config.Err("Error: ", result.Error)
		return result.Error
	}

	return nil
}
