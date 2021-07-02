package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/middlewares"
	"github.com/leobewater/udemy-orders-go-admin/models"
	"github.com/leobewater/udemy-orders-go-admin/util"
)

// Register user registration
func Register(c *fiber.Ctx) error {
	// create map with string, string to store request data
	var data map[string]string

	// get the request data
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// check password
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// map request data to model
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    1, // assign Admin role when register
	}

	// set password for user
	user.SetPassword(data["password"])

	// create user
	database.DB.Create(&user)

	return c.JSON(user)
}

// Login user login
func Login(c *fiber.Ctx) error {
	// create map with string, string to store request data
	var data map[string]string

	// get the request data
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// find user
	var user models.User
	// if user is found, map db data to &user
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "not found",
		})
	}

	// compare password
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// generate jwt and uses user.ID as the issuer
	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// store token into cookie expires in 24 hours with http only
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Claims
// type Claims struct {
// 	jwt.StandardClaims
// }

// User returns user information
func User(c *fiber.Ctx) error {
	// get the jwt cookie
	cookie := c.Cookies("jwt")

	// decode the jwt and check it validation
	id, _ := util.ParseJwt(cookie)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(user)
}

// Logout creates an expired jwt cookie
func Logout(c *fiber.Ctx) error {
	// remove cookie by creating a new cookie with a past expiration date
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// UpdateInfo updates user first, last name and email only
func UpdateInfo(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// get the jwt cookie
	cookie := c.Cookies("jwt")

	// decode the jwt and check it validation
	id, _ := util.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}

// UpdatePassword update user password
func UpdatePassword(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeUsers); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// check password
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	// get the jwt cookie
	cookie := c.Cookies("jwt")

	// decode the jwt and check it validation
	id, _ := util.ParseJwt(cookie)
	userId, _ := strconv.Atoi(id)

	user := models.User{
		Id: uint(userId),
	}

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(user)

	return c.JSON(user)
}
