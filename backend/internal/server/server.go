package server

import (
	"github.com/bigelle/online-shop/backend/internal/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	// /products
	productsHandler := handlers.NewProductHandler(db)
	productsGroup := api.Group("/products")
	productsGroup.GET("/:id", productsHandler.GetById)
	productsGroup.GET("/", productsHandler.GetAll)

	// /auth
	authHandler := handlers.NewAuthHandler(db)
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	// /cart
	cartHandler := handlers.NewCartHandler(db)
	cartGroup := api.Group("/cart")
	cartGroup.POST("/add", cartHandler.Add)

	// //orders
	ordersHandler := handlers.NewOrdersHandler(db)
	ordersGroup := api.Group("/orders")
	ordersGroup.GET("/", ordersHandler.GetAll)
	ordersGroup.POST("/create", ordersHandler.Create)

	// /payment
	paymentHandler := handlers.NewPaymentHandler(db)
	paymentGroup := api.Group("/payment")
	paymentGroup.POST("/checkout", paymentHandler.Checkout)
}
