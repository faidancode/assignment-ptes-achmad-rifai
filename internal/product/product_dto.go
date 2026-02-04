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
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ListParams struct {
	Page     int      `form:"page"`
	PageSize int      `form:"page_size"`
	Name     *string  `form:"name"`
	Category *string  `form:"category"`
	MinPrice *float64 `form:"min_price"`
	MaxPrice *float64 `form:"max_price"`
	MinStock *int32   `form:"min_stock"`
	MaxStock *int32   `form:"max_stock"`
	Sort     *string  `form:"sort"`
}
