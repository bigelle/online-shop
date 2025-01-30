package main

import (
	"github.com/bigelle/online-shop/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	//TODO: gorm open db, pass it to handlers struct

	r := gin.Default()
	r.POST("/auth/register", handlers.HandleAuthRegister)
	r.POST("/auth/login", handlers.HandleAuthLogin)
	r.GET("/products", handlers.HandleGetProducts)
	r.GET("/products/:id", handlers.HandleGetProductById)
	r.POST("/cart/add", handlers.HandleCartAdd)
	r.GET("/orders", handlers.HandleGetOrders)
	r.POST("/orders/create", handlers.HandleOrdersCreate)
	r.POST("/payment/checkout", handlers.HandlePaymentCheckout)

	r.Run()
}
