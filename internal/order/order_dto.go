package order

import "time"

type OrderItemRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id" binding:"required"`
	Items      []OrderItemRequest `json:"items" binding:"required,gt=0,dive"` //gt=0 slice validation
}

type OrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	CustomerID    string              `json:"customer_id"`
	TotalQuantity int                 `json:"total_quantity"`
	TotalPrice    float64             `json:"total_price"`
	CreatedAt     time.Time           `json:"created_at"`
	Items         []OrderItemResponse `json:"items,omitempty"`
}
