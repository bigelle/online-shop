package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrdersHandler struct {
	DB *gorm.DB
}

func NewOrdersHandler(db *gorm.DB) *OrdersHandler {
	return &OrdersHandler{DB: db}
}

func (h *OrdersHandler) GetAll(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}

func (h *OrdersHandler) Create(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}
