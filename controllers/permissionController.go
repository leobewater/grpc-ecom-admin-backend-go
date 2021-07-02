package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/grpc-ecom-admin-backend-go/database"
	"github.com/leobewater/grpc-ecom-admin-backend-go/models"
)

const routePermissions = "permissions"

// AllPermissions returns all permissions from database
func AllPermissions(c *fiber.Ctx) error {
	// middleware
	// if err := middlewares.IsAuthorized(c, routePermissions); err != nil {
	// 	return c.SendStatus(fiber.StatusUnauthorized)
	// }

	var permissions []models.Permission

	database.DB.Find(&permissions)

	return c.JSON(permissions)
}
