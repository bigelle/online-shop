package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Price       int
	Description string
	Tags        []string `gorm:"serializer:json"`
}

type User struct {
	gorm.Model
	Username       string
	Email          string
	HashedPassword string
	CsrfToken      string
	SessionToken   string
	CartItems      []CartItem `gorm:"foreignKey:UserID"`
}

type CartItem struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int
}
