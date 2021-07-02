package middlewares

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/models"
	"github.com/leobewater/udemy-orders-go-admin/util"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	cookie := c.Cookies("jwt")

	Id, err := util.ParseJwt(cookie)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	userId, _ := strconv.Atoi(Id)

	user := models.User{
		Id: uint(userId),
	}

	// find user and role
	database.DB.Preload("Role").Find(&user)

	// check permission
	role := models.Role{
		Id: user.RoleId,
	}

	database.DB.Preload("Permissions").Find(&role)

	// bypass Admin
	// if role.Name == "Admin" {
	// 	return nil
	// }

	// fmt.Printf("%+v\n", user.Role)
	// fmt.Printf("%+v\n", role.Permissions)

	// check request method and permission
	if c.Method() == "GET" {
		for _, permission := range role.Permissions {
			if permission.Name == "view_"+page || permission.Name == "edit_"+page {
				return nil
			}
		}
	} else {
		for _, permission := range role.Permissions {
			if permission.Name == "edit_"+page {
				return nil
			}
		}
	}

	//c.Status(fiber.StatusUnauthorized)
	return errors.New("unauthorized")
}
