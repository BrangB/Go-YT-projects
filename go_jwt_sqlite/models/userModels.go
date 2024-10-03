package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Rooms    []Room `gorm:"foreignKey:OwnerID"` // Establish foreign key relationship
}

type Room struct {
	gorm.Model
	OwnerID     uint // Foreign key for User
	Owner       User `gorm:"foreignKey:OwnerID"` // Reference User model by OwnerID
	Title       string
	Description string
}
