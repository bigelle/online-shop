package handlers

import (
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/database"
	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{DB: db}
}

func (h *CartHandler) Update(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusUnauthorized,
				Description: "unauthorized",
			},
		)
		return
	}
	u := user.(models.User)

	var items []schemas.CartItemRequest
	if err := ctx.BindJSON(&items); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: "bad request",
			},
		)
		return
	}

	updates := make(map[uint]int, len(items))
	for _, item := range items {
		updates[item.ProductID] = item.Quantity
	}

	cart, err := database.UpdateCart(h.DB, u.ID, updates)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: "internal server error",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusAccepted,
			Result: wrapCart(cart),
		},
	)
}

func (h *CartHandler) Remove(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusUnauthorized,
				Description: "unauthorized",
			},
		)
		return
	}
	u := user.(models.User)

	var items []schemas.CartItemRequest
	if err := ctx.BindJSON(&items); err != nil {
		var singleItem schemas.CartItemRequest
		if err := ctx.BindJSON(&singleItem); err != nil {
			ctx.JSON(http.StatusBadRequest, schemas.Response{
				Ok:          false,
				Code:        http.StatusBadRequest,
				Description: "bad request",
			})
			return
		}
		items = []schemas.CartItemRequest{singleItem}
	}

	updates := make([]uint, len(items))
	for _, item := range items {
		updates = append(updates, item.ProductID)
	}

	cart, err := database.RemoveFromCart(h.DB, u.ID, updates)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: "internal server error",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusAccepted,
			Result: wrapCart(cart),
		},
	)
}

func (h *CartHandler) Clear(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusUnauthorized,
				Description: "unauthorized",
			},
		)
		return
	}
	u := user.(models.User)

	cart, err := database.ClearCart(h.DB, u.ID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: "internal server error",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusAccepted,
			Result: wrapCart(cart),
		},
	)
}

func (h *CartHandler) View(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(
			http.StatusUnauthorized,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusUnauthorized,
				Description: "unauthorized",
			},
		)
		return
	}
	u := user.(models.User)

	cart, err := database.ViewCart(h.DB, u.ID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: "internal server error",
			},
		)
		return
	}

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusAccepted,
			Result: wrapCart(cart),
		},
	)
}

func wrapCart(cartItems []models.CartItem) schemas.CartResponse {
	items := make([]schemas.CartItemResponse, len(cartItems))
	totalPrice := 0

	for i, item := range cartItems {
		items[i] = schemas.CartItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			Price:       item.Product.Price,
		}
		totalPrice += item.Quantity * item.Product.Price
	}

	return schemas.CartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}
}
