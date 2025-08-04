package entity

import (
	"time"
)

type Ad struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type AdWithAuthor struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	AuthorLogin string    `json:"author_login"`
}

type AdResponse struct {
	AdWithAuthor AdWithAuthor
	IsOwner      *bool `json:"is_owner"`
}
