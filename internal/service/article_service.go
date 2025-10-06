package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/fhmptrdnd/weather-api-test-web-based/internal/models"
	"github.com/fhmptrdnd/weather-api-test-web-based/internal/repository"
)

// ArticleService manages articles.
type ArticleService struct {
    repo repository.Repository
}

// NewArticleService creates a new service with the provided repository.
func NewArticleService(r repository.Repository) *ArticleService {
    return &ArticleService{repo: r}
}

func generateID() string {
    b := make([]byte, 8)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// sanitizeHTML is a simple sanitizer for demonstration.
func sanitizeHTML(content string) string {
    content = strings.ReplaceAll(content, "\r", "")
    content = strings.ReplaceAll(content, "\n", "<br>")
    return content
}

func (s *ArticleService) Create(title, author, content, ownerID string) (*models.Article, error) {
    a := &models.Article{
        ID:        generateID(),
        Title:     title,
        Author:    author,
        Content:   sanitizeHTML(content),
        CreatedAt: time.Now(),
        Views:     0,
        OwnerID:   ownerID,
    }
    if err := s.repo.Create(a); err != nil {
        return nil, err
    }
    return a, nil
}

func (s *ArticleService) Get(id string) (*models.Article, error) {
    a, err := s.repo.Get(id)
    if err != nil {
        return nil, err
    }
    // increment view count
    a.Views++
    _ = s.repo.Update(a)
    return a, nil
}

// GetNoIncrement returns the article without changing the view count.
func (s *ArticleService) GetNoIncrement(id string) (*models.Article, error) {
    return s.repo.Get(id)
}

func (s *ArticleService) Update(id, title, author, content, ownerID string) (*models.Article, error) {
    a, err := s.repo.Get(id)
    if err != nil {
        return nil, err
    }
    if a.OwnerID != ownerID {
        return nil, repository.ErrNotFound
    }
    a.Title = title
    a.Author = author
    a.Content = sanitizeHTML(content)
    if err := s.repo.Update(a); err != nil {
        return nil, err
    }
    return a, nil
}

func (s *ArticleService) Delete(id, ownerID string) error {
    a, err := s.repo.Get(id)
    if err != nil {
        return err
    }
    if a.OwnerID != ownerID {
        return repository.ErrNotFound
    }
    return s.repo.Delete(id)
}
