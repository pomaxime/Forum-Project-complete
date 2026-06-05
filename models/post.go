package models

import "time"

type Post struct {
	ID        int
	UserID    int
	Username  string
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
}
