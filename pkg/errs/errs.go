package errs

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type (
	ErrHandler interface {
		HandleError(error) error
	}

	VoteErrorHandler struct{}
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate") ||
		strings.Contains(err.Error(), "duplicate key")
}

func (v VoteErrorHandler) HandleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	if IsDuplicateError(err) {
		return ErrAlreadyExists
	}
	if errors.Is(err, gorm.ErrInvalidData) {
		return ErrInvalidInput
	}
	return fmt.Errorf("internal server error")
}
