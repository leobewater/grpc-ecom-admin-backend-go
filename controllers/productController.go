package controllers

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/middlewares"
	"github.com/leobewater/udemy-orders-go-admin/models"
)

const routeProducts = "products"

// AllProducts returns all products from database
func AllProducts(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeProducts); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// get url param "page" and "page_size"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "2"))

	return c.JSON(models.Paginate(database.DB, &models.Product{}, page, pageSize))
}

// CreateProduct creates new product
func CreateProduct(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeProducts); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		// TODO - convert string to float64 err

		// fmt.Println(err)
		// c.Status(fiber.StatusBadRequest)
		// return c.JSON(fiber.Map{
		// 	"error": err
		// })
		return err
	}

	// convert request product price from string to float
	fmt.Println(reflect.TypeOf(product.Price))
	//if reflect.TypeOf(product.Price) == "string"
	//product.Price, = strconv.ParseFloat(product.Price, 64)

	database.DB.Create(&product)

	return c.JSON(product)
}

// GetProduct returns product data by given product id
func GetProduct(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeProducts); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{
		Id: uint(id),
	}

	database.DB.Find(&product)

	return c.JSON(product)
}

// UpdateProduct updates product data
func UpdateProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{
		Id: uint(id),
	}

	if err := c.BodyParser(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(product)

	return c.JSON(product)
}

// DeleteProduct deletes product from database
func DeleteProduct(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeProducts); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, _ := strconv.Atoi(c.Params("id"))

	product := models.Product{
		Id: uint(id),
	}

	database.DB.Delete(&product)

	return nil
}
