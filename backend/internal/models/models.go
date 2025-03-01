package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string
	Email          string
	HashedPassword string
	CsrfToken      string
	SessionToken   string
	CartItems      []CartItem `gorm:"foreignKey:UserID"`
	Orders         []Order    `gorm:"foreignKey:UserID"`
}

type CartItem struct {
	gorm.Model
	UserID    uint
	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int
}

type Product struct {
	gorm.Model
	Name        string
	Price       int
	Description string
	Tags        []string `gorm:"serializer:json"`
}

type Order struct {
	gorm.Model
	UserID     uint
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	TotalPrice int
	Status     string
}

type OrderItem struct {
	gorm.Model
	OrderID         uint
	ProductID       uint
	Product         Product `gorm:"foreignKey:ProductID"`
	Quantity        int
	PriceAtPurchase int
}
