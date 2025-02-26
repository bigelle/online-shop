package handlers

import (
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
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
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(
				http.StatusNotFound,
				schemas.Response{
					Ok:          false,
					Code:        http.StatusNotFound,
					Description: http.StatusText(http.StatusNotFound),
				},
			)
			return
		}
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	ctx.JSON(
		http.StatusOK,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusOK,
			Result: products,
		},
	)
}

func (p *ProductHandler) GetById(ctx *gin.Context) {
	id := ctx.GetInt("id")
	var product models.Product
	if err := p.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(
				http.StatusNotFound,
				schemas.Response{
					Ok:          false,
					Code:        http.StatusNotFound,
					Description: http.StatusText(http.StatusNotFound),
				},
			)
			return
		}
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: http.StatusText(http.StatusInternalServerError),
			},
		)
		return
	}
	ctx.JSON(
		http.StatusOK,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusOK,
			Result: product,
		},
	)
}
