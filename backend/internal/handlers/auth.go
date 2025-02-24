package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username" form:"username"`
	Password int    `json:"password" form:"password"`
}

func HandleAuthRegister(ctx *gin.Context) {
	var l Login
	if err := ctx.ShouldBind(&l); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"bad request": http.StatusBadRequest,
		})
		return
	}

}

func HandleAuthLogin(ctx *gin.Context) {
	return
}

func HandleAuthLogout(ctx *gin.Context) {
	return
}

func HandleGetProducts(ctx *gin.Context) {
	return
}

func HandleGetProductById(ctx *gin.Context) {
	return
}

func HandleCartAdd(ctx *gin.Context) {
	return
}

func HandleGetOrders(ctx *gin.Context) {
	return
}

func HandleOrdersCreate(ctx *gin.Context) {
	return
}

func HandlePaymentCheckout(ctx *gin.Context) {
	return
}
