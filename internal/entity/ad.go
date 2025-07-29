package entity

import (
	"time"
)

type Ad struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type AdWithAuthor struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorLogin string    `json:"author_login"`
}
