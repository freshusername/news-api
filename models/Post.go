package models

import (
	"time"
)

// Post model
// swagger:model Post
type Post struct {
	// example: 1
	ID int `json:"id"`
	// example: My First Post
	Title string `json:"title"`
	// example: This is the content of my first post.
	Content string `json:"content"`
	// example: 2024-02-015T00:00:00Z
	CreatedAt time.Time `json:"created_at"`
	// example: 2024-02-015T00:00:00Z
	UpdatedAt time.Time `json:"updated_at"`
}
