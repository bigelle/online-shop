package handlers

import (
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/schemas"
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
	var payment schemas.PaymentRequest
	if err := ctx.BindJSON(&payment); err != nil {
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

	//making some stuff with this payment
	//imagine we redirecting it to some payment service and waiting for result
	//on success:
	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:   true,
			Code: http.StatusAccepted,
			Result: schemas.PaymentResponse{
				OrderID:   payment.OrderID,
				PaymentID: "tx_42", // from third-party payment service
				Status:    "success",
			},
		},
	)
}
