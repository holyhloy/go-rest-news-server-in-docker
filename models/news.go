package models

type News struct {
	Id         int    `db:"id" json:"Id"`
	Title      string `db:"title" json:"Title" validate:"required,min=3,max=100"`
	Content    string `db:"content" json:"Content" validate:"required,min=10"`
	Categories []int  `json:"Categories"`
}
