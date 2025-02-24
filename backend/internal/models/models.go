package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Price       int
	Description string
	Tags        []string `gorm:"type:text"`
}
