package model

import "github.com/diogomattioli/crud/pkg/data"

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Salt   string `json:"salt"`
	Active bool   `json:"active"`
}

func (o *User) IsValidCreate() bool {
	return data.Valid(o.Name) && data.Valid(o.User) && data.Valid(o.Pass)
}

func (o *User) IsValidUpdate(old *User) bool {
	return data.Valid(o.Name) && data.Valid(o.User) && data.Valid(o.Pass) && o.User == old.User
}

func (o *User) IsValidDelete() bool {
	return o.ID > 1 // Admin
}
