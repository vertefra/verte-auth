package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Project is thje struct with all the authorization fields for a user
type Project struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   uint   `json:"userID" gom:"foreignKey:UserRefer"`
}

// ParseProject returns the parsed body with Project struct
func ParseProject(ctx *fiber.Ctx) (*Project, error) {
	project := new(Project)
	if err := ctx.BodyParser(project); err != nil {
		return nil, err
	}
	return project, nil
}
