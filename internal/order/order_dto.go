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

type ListParams struct {
	Page     int
	PageSize int
}

type OrderItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	Quantity     int     `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	CategoryName string  `json:"category_name,omitempty"`
}

type OrderResponse struct {
	ID            string              `json:"id"`
	CustomerID    string              `json:"customer_id"`
	CustomerName  string              `json:"customer_name,omitempty"`
	CustomerEmail string              `json:"customer_email,omitempty"`
	TotalQuantity int32               `json:"total_quantity"`
	TotalPrice    float64             `json:"total_price"`
	CreatedAt     time.Time           `json:"created_at"`
	Items         []OrderItemResponse `json:"items,omitempty"`
}
