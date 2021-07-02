package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"-"` // hide password
	RoleId    uint   `json:"role_id"`
	Role      Role   `json:"role" gorm:"foreignKey:RoleId"` // specify the foreign key
}

// SetPassword hash the given password and attached to the user model
func (user *User) SetPassword(password string) {
	// convert password string into byte and set the cost to 14
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

// Count() implicit implements the Entity Interface and returns the total of users
func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(user).Count(&total)
	return total
}

// Take() implicit implements the Entity Interface and returns the quered users
func (user *User) Take(db *gorm.DB, pageSize int, offset int) interface{} {
	var users []User
	db.Preload("Role").Offset(offset).Limit(pageSize).Find(&users)
	return users
}
