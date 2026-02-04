package category

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type CategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListParams struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}
