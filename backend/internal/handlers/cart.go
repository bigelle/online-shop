package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	//TODO
	return &CartHandler{DB: db}
}

func (h *CartHandler) Add(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}
