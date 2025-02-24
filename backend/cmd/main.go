package main

import (
	"log"

	"github.com/bigelle/online-shop/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("can't load .env file: %s", err.Error())
	}

	r := gin.Default()
	r.POST("/api/auth/register", handlers.HandleAuthRegister)
	r.POST("/api/auth/login", handlers.HandleAuthLogin)
	r.POST("/api/auth/logout", handlers.HandleAuthLogout)
	r.GET("/api/products", handlers.HandleGetProducts)
	r.GET("/api/products/:id", handlers.HandleGetProductById)
	r.POST("/api/cart/add", handlers.HandleCartAdd)
	r.GET("/api/orders", handlers.HandleGetOrders)
	r.POST("/api/orders/create", handlers.HandleOrdersCreate)
	r.POST("/api/payment/checkout", handlers.HandlePaymentCheckout)

	r.Run()
}
