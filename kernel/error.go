package kernel

import (
	"errors"
	"fmt"
)

var (
	EntityNotFoundError = errors.New("entity not found")
	EntityAlreadyExists = func(duplicateKey string) error {
		return fmt.Errorf("entity already exists: %v", duplicateKey)
	}
)

type Error interface {
	Code() int
	Error() string
}

type customError struct {
	code    int
	message string
}

func NewError(code int, message string) Error {
	return &customError{code, message}
}

func (e customError) Error() string {
	return e.message
}

func (e customError) Code() int {
	return e.code
}
