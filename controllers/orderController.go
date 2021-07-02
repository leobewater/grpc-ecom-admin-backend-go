package controllers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/leobewater/udemy-orders-go-admin/database"
	"github.com/leobewater/udemy-orders-go-admin/middlewares"
	"github.com/leobewater/udemy-orders-go-admin/models"
)

const routeOrders = "orders"

// AllOrders returns all orders from database
func AllOrders(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeOrders); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// get url param "page" and "page_size"
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "15"))

	return c.JSON(models.Paginate(database.DB, &models.Order{}, page, pageSize))
}

// Export exports orders into a CSV file
func Export(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeOrders); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	filePath := "./csv/orders.csv"
	if err := CreateFile(filePath); err != nil {
		return err
	}
	return c.Download(filePath)
}

func CreateFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a csv writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// query all orders
	var orders []models.Order
	database.DB.Preload("OrderItems").Find(&orders)

	// csv header
	writer.Write([]string{
		"ID", "Name", "Email", "Product Title", "Price", "Quantity",
	})

	for _, order := range orders {
		data := []string{
			strconv.Itoa(int(order.Id)),
			order.FirstName + " " + order.LastName,
			order.Email,
			"",
			"",
			"",
		}

		if err := writer.Write(data); err != nil {
			return err
		}

		for _, orderItem := range order.OrderItems {
			data := []string{
				"",
				"",
				"",
				orderItem.ProductTitle,
				fmt.Sprintf("%f", orderItem.Price),
				fmt.Sprintf("%d", orderItem.Quantity),
			}

			if err := writer.Write(data); err != nil {
				return err
			}
		}
	}

	return nil
}

// Sales holds the chart return data
type Sales struct {
	Date string `json:"date"`
	Sum  string `json:"sum"`
}

func Chart(c *fiber.Ctx) error {
	// middleware
	if err := middlewares.IsAuthorized(c, routeOrders); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var sales []Sales

	database.DB.Raw(`
		SELECT DATE_FORMAT(o.created_at, '%Y-%m-%d') as date,
		SUM(oi.price * oi.quantity) as sum
		FROM orders o
		JOIN order_items oi on o.id = oi.order_id
		GROUP BY date
	`).Scan(&sales)

	return c.JSON(sales)
}
