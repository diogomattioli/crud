package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CreateValidator interface {
	GetID() int
	ValidateCreate() error
}

type UpdateValidator[T any] interface {
	ValidateUpdate(T) error
}

type DeleteValidator interface {
	ValidateDelete() error
}

type Validate[T any] struct {
}

func (*Validate[T]) ValidateCreate() error {
	return nil
}

func (v *Validate[T]) ValidateUpdate(T) error {
	return v.ValidateCreate()
}

func (*Validate[T]) ValidateDelete() error {
	return nil
}

type ValidationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {

	bytes, err := json.Marshal(&e)
	if err == nil {
		return string(bytes)
	}

	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func ValidationErrorNew(code int, message string) ValidationError {
	return ValidationError{code, message}
}

func Valid(str string) bool {
	return len(strings.TrimSpace(str)) != 0
}

func Between(value int, min int, max int) bool {
	return value >= min && value <= max
}

func In(value int, conds ...int) bool {
	for i := 0; i < len(conds); i++ {
		if value == conds[i] {
			return true
		}
	}
	return false
}
