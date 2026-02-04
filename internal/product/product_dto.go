package product

type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	CategoryID    string  `json:"category_id" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"gte=0"`
	IsActive      *bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	CategoryID    string  `json:"category_id" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"gte=0"`
	IsActive      *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	IsActive      bool    `json:"is_active"`
	TotalSold     int     `json:"total_sold"`

	Category CategoryResponse `json:"category"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
