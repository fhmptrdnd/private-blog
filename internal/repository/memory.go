package repository

import (
	"errors"
	"sync"

	"github.com/fhmptrdnd/weather-api-test-web-based/internal/models"
)

var ErrNotFound = errors.New("not found")

// MemoryRepo is a simple in-memory repository (not for production).
type MemoryRepo struct {
    mu       sync.RWMutex
    articles map[string]*models.Article
}

// NewMemoryRepo creates a new MemoryRepo.
func NewMemoryRepo() *MemoryRepo {
    return &MemoryRepo{
        articles: make(map[string]*models.Article),
    }
}

func (m *MemoryRepo) Create(a *models.Article) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.articles[a.ID] = a
    return nil
}

func (m *MemoryRepo) Get(id string) (*models.Article, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    a, ok := m.articles[id]
    if !ok {
        return nil, ErrNotFound
    }
    return a, nil
}

func (m *MemoryRepo) Update(a *models.Article) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, ok := m.articles[a.ID]; !ok {
        return ErrNotFound
    }
    m.articles[a.ID] = a
    return nil
}

func (m *MemoryRepo) Delete(id string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, ok := m.articles[id]; !ok {
        return ErrNotFound
    }
    delete(m.articles, id)
    return nil
}
