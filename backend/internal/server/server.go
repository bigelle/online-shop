package server

import (
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/handlers"
	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
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
	authGroup.POST("/logout", authHandler.Logout).Use(authorize(authHandler.DB))

	// /cart
	cartHandler := handlers.NewCartHandler(db)
	cartGroup := api.Group("/cart").Use(authorize(authHandler.DB))
	cartGroup.POST("/add", cartHandler.Update)
	cartGroup.POST("/remove", cartHandler.Remove)
	cartGroup.POST("/clear", cartHandler.Clear)

	// //orders
	ordersHandler := handlers.NewOrdersHandler(db)
	ordersGroup := api.Group("/orders").Use(authorize(authHandler.DB))
	ordersGroup.GET("/", ordersHandler.GetAll)
	ordersGroup.POST("/create", ordersHandler.Create)

	// /payment
	paymentHandler := handlers.NewPaymentHandler(db)
	paymentGroup := api.Group("/payment").Use(authorize(authHandler.DB))
	paymentGroup.POST("/checkout", paymentHandler.Checkout)
}

func authorize(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		st, err := ctx.Request.Cookie("session_token")
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				schemas.Response{
					Ok:          false,
					Code:        http.StatusUnauthorized,
					Description: "login is required for this path",
				},
			)
		}

		var usr models.User
		err = db.Model(&models.User{}).Where("session_token = ?", st.Value).Select("*").First(&usr).Error
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				schemas.Response{
					Ok:          false,
					Code:        http.StatusUnauthorized,
					Description: "invalid session",
				},
			)
		}

		ctx.Set("user", usr)
		ctx.Next()
	}
}
