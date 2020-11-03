package models

import (
	"fmt"
	"github/vertefra/verte_auth_server/config"
	"github/vertefra/verte_auth_server/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// User is the structure associated with all the user that use the provided service
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

// ParseUser returns the parsed body with user Struct
func ParseUser(ctx *fiber.Ctx) (*User, error) {
	user := new(User)
	if err := ctx.BodyParser(user); err != nil {
		config.Err("Error while parsing the request body")
		config.Err("Error: ", err)
		return nil, err
	}

	if user.Password == "" {
		err := fmt.Errorf("Password field is empty")
		return nil, err
	}

	if utils.IsEmailValid(user.Email) {
		return user, nil
	}

	err := fmt.Errorf("Email is not valid")

	return nil, err
}

// FindUserByID returns a pointer to a User Struct with ID matching the ID of the user
func FindUserByID(db *gorm.DB, ID uint) (*User, error) {
	foundUser := new(User)
	result := db.Table("users").First(&foundUser, ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return foundUser, nil
}

// IsPresent check if an user of type User is already in the database
func (u *User) IsPresent(db *gorm.DB) bool {
	foundUser := new(User)
	db.Table("users").Where(&User{Email: u.Email}).Find(&foundUser)
	return u.Email == foundUser.Email
}
