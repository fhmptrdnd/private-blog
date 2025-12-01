// Package service implements business logic with Functional Programming principles.
//
// PRINSIP FUNCTIONAL PROGRAMMING YANG DITERAPKAN:
//
// 1. IMMUTABILITY (Ketidakberubahan Data)
//    - Semua method menerima dan mengembalikan VALUE (models.Article), bukan pointer
//    - Saat update, kita membuat COPY baru dari struct, bukan mengubah yang lama
//    - Contoh: updated := a; updated.Views++ (copy-on-write pattern)
//
// 2. PURE FUNCTIONS (Fungsi Murni)
//    - sanitizeHTML() adalah pure function: input sama -> output sama
//    - Tidak ada side effects tersembunyi dalam fungsi bisnis logic
//    - Deterministik dan mudah di-test
//
// 3. DEPENDENCY INJECTION (Isolasi Side Effects)
//    - Clock dan IDGenerator di-inject sebagai interface
//    - Side effects (time.Now, random) diisolasi dari business logic
//    - Memudahkan testing dengan mock objects
//
// 4. COMMAND-QUERY SEPARATION (CQS)
//    - Get() adalah QUERY: hanya membaca, tidak mengubah state
//    - IncrementViews() adalah COMMAND: mengubah state secara eksplisit
//    - Pemisahan jelas antara operasi baca dan tulis
//
// 5. VALUE SEMANTICS
//    - Repository bekerja dengan values, bukan pointers
//    - Menghindari shared mutable state
//    - Lebih aman dari race conditions
package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/fhmptrdnd/private-blog/internal/models"
	"github.com/fhmptrdnd/private-blog/internal/repository"
)

// Clock abstracts time generation.
// FP Principle: Dependency Injection untuk isolasi side effects.
// Dengan interface ini, kita bisa inject mock clock saat testing.
type Clock interface {
	Now() time.Time
}

// IDGenerator abstracts ID generation.
// FP Principle: Dependency Injection untuk isolasi side effects.
// Random ID generation adalah side effect yang perlu diisolasi.
type IDGenerator interface {
	Generate() string
}

// RealClock implements Clock using system time.
type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

// RealIDGenerator implements IDGenerator using crypto/rand.
type RealIDGenerator struct{}

func (RealIDGenerator) Generate() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ArticleService manages articles.
type ArticleService struct {
	repo  repository.Repository
	clock Clock
	idGen IDGenerator
}

// NewArticleService creates a new service with dependencies.
func NewArticleService(r repository.Repository, c Clock, g IDGenerator) *ArticleService {
	return &ArticleService{
		repo:  r,
		clock: c,
		idGen: g,
	}
}

// sanitizeHTML is a pure function for content sanitization.
// FP Principle: Pure Function - input yang sama selalu menghasilkan output yang sama.
// Tidak ada side effects, tidak mengubah state global, deterministik.
func sanitizeHTML(content string) string {
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "<br>")
	return content
}

// Create creates a new article. Returns the created article value.
// FP Principle: Mengembalikan VALUE baru, bukan pointer (Immutability).
// Side effects (ID, Time) sudah diisolasi melalui injected dependencies.
func (s *ArticleService) Create(title, author, content, ownerID string) (models.Article, error) {
	a := models.Article{
		ID:        s.idGen.Generate(),
		Title:     title,
		Author:    author,
		Content:   sanitizeHTML(content),
		CreatedAt: s.clock.Now(),
		Views:     0,
		OwnerID:   ownerID,
	}
	if err := s.repo.Create(a); err != nil {
		return models.Article{}, err
	}
	return a, nil
}

// Get retrieves an article by ID. It does NOT increment views (Query).
// FP Principle: Command-Query Separation - ini adalah QUERY murni tanpa side effects.
func (s *ArticleService) Get(id string) (models.Article, error) {
	return s.repo.Get(id)
}

// IncrementViews increments the view count for an article (Command).
// FP Principle: Command-Query Separation - ini adalah COMMAND yang mengubah state.
// FP Principle: Immutability - kita buat copy baru (updated), bukan mutasi langsung.
func (s *ArticleService) IncrementViews(id string) error {
	a, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	// FP: Copy-on-write pattern - buat salinan baru, jangan ubah yang asli
	updated := a
	updated.Views++
	return s.repo.Update(updated)
}

// Update updates an existing article. Returns the updated article value.
// FP Principle: Mengembalikan VALUE baru hasil update, bukan mengubah pointer.
func (s *ArticleService) Update(id, title, author, content, ownerID string) (models.Article, error) {
	a, err := s.repo.Get(id)
	if err != nil {
		return models.Article{}, err
	}
	if a.OwnerID != ownerID {
		return models.Article{}, repository.ErrNotFound
	}

	// FP: Copy-on-write - buat versi baru dari article, jangan mutasi yang lama
	updated := a
	updated.Title = title
	updated.Author = author
	updated.Content = sanitizeHTML(content)
	
	if err := s.repo.Update(updated); err != nil {
		return models.Article{}, err
	}
	return updated, nil
}

// Delete removes an article.
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
