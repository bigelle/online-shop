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

type CartItem struct {
	ProductId uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
