package schemas

type Response struct {
	//true on success
	Ok bool `json:"ok"`
	//response code
	Code int `json:"code"`
	//if ok == false, there would be an error description
	Description string `json:"description,omitzero"`
	// if ok == true, there would be a result body (only if applicable)
	Result any `json:"result,omitempty"`
}

type Login struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice int                `json:"total_price"`
}

type CartItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type CartItemResponse struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}

type OrderRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type OrderResponse struct {
	ID         uint        `json:"id"`
	TotalPrice int         `json:"total_price"`
	Status     string      `json:"status"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID       uint   `json:"product_id"`
	ProductName     string `json:"product_name"`
	Quantity        int    `json:"quantity"`
	PriceAtPurchase int    `json:"price_at_purchase"`
}

type PaymentRequest struct {
	OrderID uint `json:"order_id"`
}

type PaymentResponse struct {
	OrderID   uint   `json:"order_id"`
	PaymentID string `json:"payment_id,omitzero,"`
	Status    string `json:"status"`
}
