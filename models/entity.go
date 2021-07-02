package models

import "gorm.io/gorm"

type Entity interface {
	// count the total of records and return int64
	Count(db *gorm.DB) int64

	// using interface all different type of structs
	Take(db *gorm.DB, pageSize int, offset int) interface{}
}
