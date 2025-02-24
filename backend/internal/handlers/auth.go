package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	//TODO
	ctx.JSON(http.StatusOK, "pong")
}
