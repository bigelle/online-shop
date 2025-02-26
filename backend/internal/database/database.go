package database

import (
	"github.com/bigelle/online-shop/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Product{},
		&models.User{},
	)
}

func Close(db *gorm.DB) {
	//TODO
}
