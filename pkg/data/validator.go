package data

import "strings"

type CreateValidator interface {
	IsValidCreate() bool
}

type UpdateValidator[T any] interface {
	IsValidUpdate(T) bool
}

type DeleteValidator interface {
	IsValidDelete() bool
}

type Validate[T any] struct {
}

func (*Validate[T]) IsValidCreate() bool {
	return true
}

func (v *Validate[T]) IsValidUpdate(T) bool {
	return v.IsValidCreate()
}

func (*Validate[T]) IsValidDelete() bool {
	return true
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
