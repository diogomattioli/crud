package model

import "strings"

type Validator interface {
	IsValid() bool
}

func valid(str string) bool {
	return len(strings.TrimSpace(str)) != 0
}

func between(value int, min int, max int) bool {
	return value >= min && value <= max
}

func in(value int, conds ...int) bool {
	for i := 0; i < len(conds); i++ {
		if value == conds[i] {
			return true
		}
	}
	return false
}
