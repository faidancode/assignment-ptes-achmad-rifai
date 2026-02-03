package customer

import "errors"

var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
