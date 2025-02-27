package database

import (
	"errors"

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
		&models.CartItem{},
	)

}

func Close(db *gorm.DB) {
	//TODO
}

func UpdateCart(db *gorm.DB, userID uint, updates map[uint]int) (*models.User, error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for productID, quantityChange := range updates {
		var item models.CartItem
		err := tx.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error

		switch {
		case err == nil:
			item.Quantity += quantityChange
			if item.Quantity > 0 {
				if err := tx.Save(&item).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			} else {
				if err := tx.Delete(&item).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		case errors.Is(err, gorm.ErrRecordNotFound) && quantityChange > 0:
			newItem := models.CartItem{UserID: userID, ProductID: productID, Quantity: quantityChange}
			if err := tx.Create(&newItem).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		default:
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return getUserWithCart(db, userID)
}

func RemoveFromCart(db *gorm.DB, userID uint, productIDs []uint) (*models.User, error) {
	if err := db.Where("user_id = ? AND product_id IN (?)", userID, productIDs).Delete(&models.CartItem{}).Error; err != nil {
		return nil, err
	}
	return getUserWithCart(db, userID)
}

func ClearCart(db *gorm.DB, userID uint) (*models.User, error) {
	if err := db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		return nil, err
	}
	return getUserWithCart(db, userID)
}

func getUserWithCart(db *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	err := db.Preload("CartItems.Product").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
