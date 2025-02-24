package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{DB: db}
}

func (h *PaymentHandler) Checkout(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}
