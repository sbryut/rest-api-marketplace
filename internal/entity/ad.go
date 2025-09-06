// Package entity defines domain entities
package entity

import (
	"time"
)

// Ad represents an advertisement
type Ad struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// AdWithAuthor represents an ad along with author's login
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

// AdResponse represents ad response for API with ownership info
type AdResponse struct {
	AdWithAuthor AdWithAuthor
	IsOwner      *bool `json:"is_owner"`
}
