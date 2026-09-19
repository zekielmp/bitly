package dto

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartResponse struct {
	ID        uint                `json:"id"`
	UserID    uint                `json:"user_id"`
	CartItems []CartItemsResponse `json:"cart_items"`
	Total     float64             `json:"total"`
}

type CartItemsResponse struct {
	ID       uint            `json:"id"`
	Product  productResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
}

type OrderResponse struct {
	ID          uint                 `json:"id"`
	UserID      uint                 `json:"user_id"`
	Status      string               `json:"status"`
	TotalAmount float64              `json:"total_amount"`
	OrderItems  []OrderItemsResponse `json:"order_items"`
	CreatedAt   string               `json:"created_at"`
}

type OrderItemsResponse struct {
	ID       uint            `json:"id"`
	Product  productResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}
