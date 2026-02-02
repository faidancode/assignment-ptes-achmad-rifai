package category

import "errors"

var (
	ErrInvalidCategoryName = errors.New("invalid category name")
	ErrCategoryNotFound    = errors.New("category not found")
)
