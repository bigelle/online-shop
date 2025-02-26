package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Price       int
	Description string
	Tags        []string `gorm:"type:text"`
}

type User struct {
	gorm.Model
	Username       string
	Email          string
	HashedPassword string
	CsrfToken      string
	SessionToken   string
}
