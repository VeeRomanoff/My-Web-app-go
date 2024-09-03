package models

// article model definition
type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Author  int64  `json:"author"`
	Content string `json:"content"`
}
