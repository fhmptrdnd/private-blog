package models

import "time"

// article, struktur data buat artikel
type Article struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Views     int       `json:"views"`
    OwnerID   string    `json:"owner_id"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"` // nullable, buat soft delete
}
