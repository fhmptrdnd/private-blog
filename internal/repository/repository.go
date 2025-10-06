package repository

import "github.com/fhmptrdnd/weather-api-test-web-based/internal/models"

// Repository defines storage operations for articles.
type Repository interface {
    Create(a *models.Article) error
    Get(id string) (*models.Article, error)
    Update(a *models.Article) error
    Delete(id string) error
}
