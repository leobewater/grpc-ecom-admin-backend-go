package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/middlewares"
	"github.com/leobewater/udemy-orders-go-admin/models"
)

const routeRoles = "roles"

// AllRoles returns all roles from database
func AllRoles(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeRoles); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var roles []models.Role

	database.DB.Preload("Permissions").Find(&roles)

	return c.JSON(roles)
}

// CreateRole creates a new role
func CreateRole(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeRoles); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var roleDto fiber.Map

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	// get the permissions from roleDto and cast it as slice of interface
	list := roleDto["permissions"].([]interface{})

	// create a slice to store all the posted permissions
	permissions := make([]models.Permission, len(list))

	// loop each posted permission
	for i, permissionId := range list {
		//fmt.Printf("\n%T", permissionId)

		// case the permission id from interface to string then convert to int
		id, _ := strconv.Atoi(fmt.Sprintf("%v", permissionId)) //permissionId.(int)

		// fmt.Printf("\nParsed %d:", id)
		permissions[i] = models.Permission{
			Id: uint(id),
		}
	}

	role := models.Role{
		Name:        roleDto["name"].(string),
		Permissions: permissions,
	}

	database.DB.Create(&role)

	return c.JSON(role)
}

// GetRole returns role data by given role id
func GetRole(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeRoles); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	database.DB.Preload("Permissions").Find(&role)

	return c.JSON(role)
}

// UpdateRole updates role data
func UpdateRole(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeRoles); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	var roleDto fiber.Map

	if err := c.BodyParser(&roleDto); err != nil {
		return err
	}

	// get the permissions from roleDto and cast it as slice of interface
	list := roleDto["permissions"].([]interface{})

	// create a slice to store all the posted permissions
	permissions := make([]models.Permission, len(list))

	// loop each posted permission
	for i, permissionId := range list {
		// case the permission id as string then convert to int
		id, _ := permissionId.(float64) //strconv.Atoi(permissionId.(string))

		permissions[i] = models.Permission{
			Id: uint(id),
		}
	}

	// remove previous saved role_permissions
	var result interface{}
	database.DB.Table("role_permissions").Where("role_id", id).Delete(&result)

	// create the role again
	role := models.Role{
		Id:          uint(id),
		Name:        roleDto["name"].(string),
		Permissions: permissions,
	}

	database.DB.Model(&role).Updates(role)

	return c.JSON(role)
}

// DeleteRole deletes role from database
func DeleteRole(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeRoles); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	role := models.Role{
		Id: uint(id),
	}

	// remove previous saved role_permissions
	var result interface{}
	database.DB.Table("role_permissions").Where("role_id", id).Delete(&result)

	database.DB.Delete(&role)

	return nil
}
