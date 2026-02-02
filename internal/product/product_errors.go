package product

import "errors"

var (
	ErrInvalidProductName = errors.New("invalid product name")
	ErrProductNotFound    = errors.New("product not found")
)
