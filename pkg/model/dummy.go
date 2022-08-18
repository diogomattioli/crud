package model

type Dummy struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (o *Dummy) IsValid() bool {
	return valid(o.Title)
}
