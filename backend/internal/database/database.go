package database

import (
	"errors"
	"log"

	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	STATUS_PENDING = "pending"
	STATUS_PAYMENT = "payment"
)

func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)

}

func Close(db *gorm.DB) {
	if db == nil {
		log.Println("Warning: trying to close a nil database connection")
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting SQL DB from GORM: %s", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %s", err)
	}

}

func FindUser(db *gorm.DB, email string) (*models.User, error) {
	var usr models.User
	err := db.Model(&models.User{}).Where("email = ?", email).Select("*").First(&usr).Error
	return &usr, err
}

func AddUser(db *gorm.DB, l schemas.Login) error {
	return db.Create(&models.User{
		Username:       l.Username,
		Email:          l.Email,
		HashedPassword: l.Password,
	}).Error
}

func UpdateCart(db *gorm.DB, userID uint, updates map[uint]int) ([]models.CartItem, error) {
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
	return ViewCart(db, userID)
}

func RemoveFromCart(db *gorm.DB, userID uint, productIDs []uint) ([]models.CartItem, error) {
	if err := db.Where("user_id = ? AND product_id IN (?)", userID, productIDs).Delete(&models.CartItem{}).Error; err != nil {
		return nil, err
	}
	return ViewCart(db, userID)
}

func ClearCart(db *gorm.DB, userID uint) ([]models.CartItem, error) {
	if err := db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		return nil, err
	}
	return ViewCart(db, userID)
}

func ViewCart(db *gorm.DB, userID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func AddToOrders(db *gorm.DB, userID uint) (*models.Order, error) {
	var user models.User
	err := db.Preload("CartItems.Product").
		Preload("Orders.OrderItems.Product").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	cart := user.CartItems

	order := models.Order{
		UserID:     userID,
		TotalPrice: 0,
		Status:     STATUS_PAYMENT,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		var items []models.OrderItem
		total := 0

		for _, req := range cart {
			item := models.OrderItem{
				OrderID:         order.ID,
				ProductID:       req.ProductID,
				Quantity:        req.Quantity,
				PriceAtPurchase: req.Product.Price,
			}
			items = append(items, item)
			total += req.Quantity * req.Product.Price
		}

		if len(items) == 0 {
			return errors.New("no valid items to insert")
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		if err := tx.Model(&order).Update("TotalPrice", total).Error; err != nil {
			return err
		}
		order.TotalPrice = total
		order.OrderItems = items

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrderStatus(db *gorm.DB, orderID uint, status string) error {
	return db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func GetProduct(db *gorm.DB, productID uint) (*models.Product, error) {
	var product models.Product
	if err := db.First(&product, productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
