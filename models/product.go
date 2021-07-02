package models

import "gorm.io/gorm"

type Product struct {
	Id          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

// Count() implicit implements the Entity Interface and returns the total of products
func (product *Product) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(product).Count(&total)
	return total
}

// Take() implicit implements the Entity Interface and returns the quered products
func (product *Product) Take(db *gorm.DB, pageSize int, offset int) interface{} {
	var products []Product
	db.Offset(offset).Limit(pageSize).Find(&products)
	return products
}
