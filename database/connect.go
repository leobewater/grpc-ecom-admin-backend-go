package database

import (
	"github.com/leobewater/grpc-ecom-admin-backend-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect connects to database
func Connect() {
	database, err := gorm.Open(mysql.Open("root:secret@tcp(127.0.0.1:6333)/go_admin"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to the database")
	}

	DB = database

	// run migrations
	database.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Product{}, &models.Order{}, &models.OrderItem{})
}
