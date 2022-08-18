package model

import "github.com/diogomattioli/crud/pkg/data"

type Dummy struct {
	data.Validate[Dummy]
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (o Dummy) IsCreateValid() bool {
	return data.Valid(o.Title)
}
