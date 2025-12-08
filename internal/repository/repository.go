package repository

import "github.com/fhmptrdnd/private-blog/internal/models"

type Repository interface {
    Create(a models.Article) error
    Get(id string) (models.Article, error)
    Update(a models.Article) error
    Delete(id string) error
}
