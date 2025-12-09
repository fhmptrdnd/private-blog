package repository

import (
	"errors"
	"github.com/fhmptrdnd/weather-api-test-web-based/internal/models"
)

// errnotfound, error kalo data ga ketemu
var ErrNotFound = errors.New("not found")

// createfunc, function type buat create artikel
type CreateFunc func(models.Article) error

// getfunc, function type buat get artikel
type GetFunc func(string) (models.Article, error)

// updatefunc, function type buat update artikel
type UpdateFunc func(models.Article) error

// deletefunc, function type buat delete artikel
type DeleteFunc func(string) error

// listbyownerfunc, function type buat list artikel by owner
type ListByOwnerFunc func(ownerID string) ([]models.Article, error)

// repository, struct yang isinya function-function (bukan interface!)
// ini penerapan "functions as first-class citizens" di layer data
type Repository struct {
	Create      CreateFunc
	Get         GetFunc
	Update      UpdateFunc
	Delete      DeleteFunc
	ListByOwner ListByOwnerFunc
}
