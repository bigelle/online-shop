package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bigelle/online-shop/backend/internal/database"
	"github.com/bigelle/online-shop/backend/internal/models"
	"github.com/bigelle/online-shop/backend/internal/schemas"
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

	resp := make([]schemas.OrderResponse, len(u.Orders))
	for i, order := range u.Orders {
		resp[i] = wrapOrder(order)
	}

	ctx.JSON(
		http.StatusOK,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusOK,
			Result: resp,
		},
	)
}

func (h *OrdersHandler) Create(ctx *gin.Context) {
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

	order, err := database.AddToOrders(h.DB, u.ID)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: fmt.Sprintf("internal server error: %s", err.Error()),
			},
		)
		return
	}

	paymentReq := schemas.PaymentRequest{
		OrderID: order.ID,
	}
	paymentResp, err := h.redirectToPaymentService(paymentReq, ctx)
	if err != nil || paymentResp.Status != "success" {
		log.Println("RIGHT HERE")
		ctx.JSON(
			http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: "payment failed",
			},
		)
		return
	}

	if err := database.UpdateOrderStatus(h.DB, order.ID, database.STATUS_PENDING); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			schemas.Response{
				Ok:          false,
				Code:        http.StatusInternalServerError,
				Description: fmt.Sprintf("failed to update order status: %s", err.Error()),
			},
		)
		return
	}
	order.Status = database.STATUS_PENDING

	ctx.JSON(
		http.StatusAccepted,
		schemas.Response{
			Ok:     true,
			Code:   http.StatusAccepted,
			Result: wrapOrder(*order),
		},
	)
}

func (h *OrdersHandler) redirectToPaymentService(req schemas.PaymentRequest, ctx *gin.Context) (*schemas.PaymentResponse, error) {
	client := &http.Client{}
	jsonData, _ := json.Marshal(req)

	paymentURL := "http://localhost:8080/api/payment/checkout"
	request, err := http.NewRequest("POST", paymentURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	if csrfToken := ctx.GetHeader("X-CSRF-Token"); csrfToken != "" {
		request.Header.Set("X-CSRF-Token", csrfToken)
	}
	for _, cookie := range ctx.Request.Cookies() {
		request.AddCookie(cookie)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response schemas.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	paymentData, ok := response.Result.(map[string]any)
	if !ok {
		return nil, err
	}

	paymentResp := schemas.PaymentResponse{
		OrderID:   uint(paymentData["order_id"].(float64)),
		PaymentID: paymentData["payment_id"].(string),
		Status:    paymentData["status"].(string),
	}

	return &paymentResp, nil
}

func wrapOrder(order models.Order) schemas.OrderResponse {
	items := make([]schemas.OrderItem, len(order.OrderItems))
	for j, item := range order.OrderItems {
		items[j] = schemas.OrderItem{
			ProductID:       item.ProductID,
			ProductName:     item.Product.Name,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
		}
	}

	return schemas.OrderResponse{
		ID:         order.ID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		Items:      items,
	}
}
