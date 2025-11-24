package models

import "time"

// article, struktur data buat artikel
type Article struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    Views     int       `json:"views"`
    OwnerID   string    `json:"owner_id"`
}
