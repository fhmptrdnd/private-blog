package repository

import (
	"database/sql"
	"errors"
	"github.com/fhmptrdnd/private-blog/internal/models"
)

var ErrNotFound = errors.New("not found")

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(db *sql.DB) *SQLiteRepo {
	return &SQLiteRepo{db: db}
}

func (s *SQLiteRepo) Create(a models.Article) error {
	query := `INSERT INTO articles (id, title, author, content, created_at, views, owner_id) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := s.db.Exec(query, a.ID, a.Title, a.Author, a.Content, a.CreatedAt, a.Views, a.OwnerID)
	return err
}

func (s *SQLiteRepo) Get(id string) (models.Article, error) {
	query := `SELECT id, title, author, content, created_at, views, owner_id FROM articles WHERE id = ?`
	
	row := s.db.QueryRow(query, id)

	var a models.Article

	err := row.Scan(&a.ID, &a.Title, &a.Author, &a.Content, &a.CreatedAt, &a.Views, &a.OwnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Article{}, ErrNotFound
		}
		return models.Article{}, err
	}
	return a, nil
}

func (s *SQLiteRepo) Update(a models.Article) error {
	query := `UPDATE articles SET title=?, author=?, content=?, views=? WHERE id=?`
	
	res, err := s.db.Exec(query, a.Title, a.Author, a.Content, a.Views, a.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteRepo) Delete(id string) error {
	query := `DELETE FROM articles WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}