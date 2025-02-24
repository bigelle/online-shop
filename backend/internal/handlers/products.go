package handlers

import (
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (p *ProductHandler) GetAll(ctx *gin.Context) {
	var products []models.Product
	if err := p.DB.Find(&products).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (p *ProductHandler) GetById(ctx *gin.Context) {
	id := ctx.GetInt("id")
	var product models.Product
	if err := p.DB.First(&product, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
