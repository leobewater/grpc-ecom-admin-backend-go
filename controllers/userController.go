package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/middlewares"
	"github.com/leobewater/udemy-orders-go-admin/models"
)

const tempPassword = "1234"

const routeUsers = "users"

// AllUsers returns all users from database
func AllUsers(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// get url param "page" and "page_size"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "2"))

	return c.JSON(models.Paginate(database.DB, &models.User{}, page, pageSize))
}

// CreateUser creates new user with auto generated password
func CreateUser(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// auto generate a password for new user
	user.SetPassword(tempPassword)

	database.DB.Create(&user)

	return c.JSON(user)
}

// GetUser returns user data by given user id
func GetUser(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	database.DB.Preload("Role").Find(&user)

	return c.JSON(user)
}

// UpdateUser updates user data
func UpdateUser(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

// DeleteUser deletes user from database
func DeleteUser(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	user := models.User{
		Id: uint(id),
	}

	database.DB.Delete(&user)

	return nil
}
