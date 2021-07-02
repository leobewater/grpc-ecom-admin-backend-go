package models

import "gorm.io/gorm"

type Order struct {
	Id         uint        `json:"id"`
	FirstName  string      `json:"-"`              // hide this
	LastName   string      `json:"-"`              // hide this
	Name       string      `json:"name" gorm:"-"`  // not adding column to database
	Email      string      `json:"email"`
	Total      float64     `json:"total" gorm:"-"` // not adding column to database
	UpdatedAt  string      `json:"updated_at"`
	CreatedAt  string      `json:"created_at"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderId"`
}

type OrderItem struct {
	Id           uint    `json:"id"`
	OrderId      uint    `json:"order_id"`
	ProductTitle string  `json:"product_title"`
	Price        float64 `json:"price"`
	Quantity     uint    `json:"quantity"`
}

// Count() implicit implements the Entity Interface and returns the total of orders
func (order *Order) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(order).Count(&total)
	return total
}

// Take() implicit implements the Entity Interface and returns the quered orders
func (order *Order) Take(db *gorm.DB, pageSize int, offset int) interface{} {
	var orders []Order
	db.Preload("OrderItems").Offset(offset).Limit(pageSize).Find(&orders)

	for i := range orders {

		// calculate order total
		var total float64 = 0
		for _, orderItem := range orders[i].OrderItems {
			total += orderItem.Price * float64(orderItem.Quantity)
		}

		// concat first name and last name to name
		orders[i].Name = orders[i].FirstName + " " + orders[i].LastName
		orders[i].Total = total
	}

	return orders
}
