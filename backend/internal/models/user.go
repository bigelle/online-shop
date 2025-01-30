package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Cart []Product
}

type Product struct {
	gorm.Model
	Name        string
	Description string
	Rating      float32
	Price       int
	Shop        Shop
}

type Shop struct {
	gorm.Model
	Name        string
	Description string
	Rating      float32
}
