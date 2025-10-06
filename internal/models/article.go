package models

import "time"

// Article represents a simple article entity.
type Article struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    Views     int       `json:"views"`
    OwnerID   string    `json:"owner_id"`
}
