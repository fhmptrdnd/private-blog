// package service, ini tempat logic bisnis aplikasi
package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/fhmptrdnd/private-blog/internal/models"
	"github.com/fhmptrdnd/private-blog/internal/repository"
)

// clockfunc, function buat dapetin waktu sekarang
// ini contoh "function as first-class citizen", function bisa jadi tipe data
type ClockFunc func() time.Time

// idgenfunc, function buat generate id random
type IDGenFunc func() string

// newrealclock, bikin function yang return waktu sekarang
// ini contoh function yang return function (closure)
func NewRealClock() ClockFunc {
	return func() time.Time { return time.Now() }
}

// newrealidgen, bikin function yang generate id random
func NewRealIDGen() IDGenFunc {
	return func() string {
		b := make([]byte, 8)
		rand.Read(b)
		return hex.EncodeToString(b)
	}
}

// articleservice, struct utama buat manage artikel
type ArticleService struct {
	repo  repository.Repository
	clock ClockFunc // ini function, bukan interface!
	idGen IDGenFunc
}

// newarticleservice, bikin service baru
// parameter clock sama idgen itu function, bukan struct
func NewArticleService(r repository.Repository, clock ClockFunc, idGen IDGenFunc) *ArticleService {
	return &ArticleService{
		repo:  r,
		clock: clock,
		idGen: idGen,
	}
}

// sanitizehtml, bersihin html, ganti newline jadi <br>
// pure function: input sama = output sama, ga ada efek samping
func sanitizeHTML(content string) string {
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "<br>")
	return content
}

// create, bikin artikel baru
// return value (bukan pointer) biar immutable
func (s *ArticleService) Create(title, author, content, ownerID string) (models.Article, error) {
	now := s.clock()
	a := models.Article{
		ID:        s.idGen(),
		Title:     title,
		Author:    author,
		Content:   sanitizeHTML(content),
		CreatedAt: now,
		UpdatedAt: now,  // set updatedat = createdat saat create
		Views:     0,
		OwnerID:   ownerID,
	}
	if err := s.repo.Create(a); err != nil {
		return models.Article{}, err
	}
	return a, nil
}

// get, ambil artikel berdasarkan id
// cuma baca aja, ga ngubah apapun (pure query)
func (s *ArticleService) Get(id string) (models.Article, error) {
	return s.repo.Get(id)
}

// incrementviews, nambah jumlah views artikel
// ini ngubah state (command), beda sama get yang cuma baca
func (s *ArticleService) IncrementViews(id string) error {
	a, err := s.repo.Get(id)
	if err != nil {
		return err
	}
	// bikin copy dulu, baru ubah (biar immutable)
	updated := a
	updated.Views++
	return s.repo.Update(updated)
}

// update, update artikel yang udah ada
func (s *ArticleService) Update(id, title, author, content, ownerID string) (models.Article, error) {
	a, err := s.repo.Get(id)
	if err != nil {
		return models.Article{}, err
	}
	if a.OwnerID != ownerID {
		return models.Article{}, repository.ErrNotFound
	}

	// copy dulu, baru update field-nya
	updated := a
	updated.Title = title
	updated.Author = author
	updated.Content = sanitizeHTML(content)
	updated.UpdatedAt = s.clock()  // update timestamp saat update
	
	if err := s.repo.Update(updated); err != nil {
		return models.Article{}, err
	}
	return updated, nil
}

// delete, hapus artikel
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

// listmyarticles, ambil semua artikel milik user
// ini query, ga ngubah state
func (s *ArticleService) ListMyArticles(ownerID string) ([]models.Article, error) {
	return s.repo.ListByOwner(ownerID)
}
