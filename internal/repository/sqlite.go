package repository

import (
	"database/sql"

	"github.com/fhmptrdnd/private-blog/internal/models"
	_ "modernc.org/sqlite"
)

// return repository struct yang isinya function-function
// state (db connection) disimpan dalam closure
func NewSQLiteRepo(dbPath string) (Repository, error) {
	// open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return Repository{}, err
	}

	// test connection
	if err := db.Ping(); err != nil {
		return Repository{}, err
	}

	// create table kalo belum ada
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		views INTEGER DEFAULT 0,
		owner_id TEXT NOT NULL,
		deleted_at DATETIME
	);
	`
	if _, err := db.Exec(createTableSQL); err != nil {
		return Repository{}, err
	}

	// return repository dengan closures yang capture db connection
	return Repository{
		// create, insert artikel baru
		Create: func(a models.Article) error {
			query := `
				INSERT INTO articles (id, title, author, content, created_at, updated_at, views, owner_id)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`
			_, err := db.Exec(query, a.ID, a.Title, a.Author, a.Content, a.CreatedAt, a.UpdatedAt, a.Views, a.OwnerID)
			return err
		},

		// get, ambil artikel by id
		Get: func(id string) (models.Article, error) {
			query := `
				SELECT id, title, author, content, created_at, updated_at, views, owner_id, deleted_at
				FROM articles
				WHERE id = ? AND deleted_at IS NULL
			`
			var a models.Article
			err := db.QueryRow(query, id).Scan(
				&a.ID, &a.Title, &a.Author, &a.Content,
				&a.CreatedAt, &a.UpdatedAt, &a.Views, &a.OwnerID, &a.DeletedAt,
			)
			if err == sql.ErrNoRows {
				return models.Article{}, ErrNotFound
			}
			if err != nil {
				return models.Article{}, err
			}
			return a, nil
		},

		// update, update artikel yang ada
		Update: func(a models.Article) error {
			query := `
				UPDATE articles
				SET title = ?, author = ?, content = ?, updated_at = ?, views = ?
				WHERE id = ? AND owner_id = ?
			`
			result, err := db.Exec(query, a.Title, a.Author, a.Content, a.UpdatedAt, a.Views, a.ID, a.OwnerID)
			if err != nil {
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return ErrNotFound
			}
			return nil
		},

		// delete, soft delete artikel (set deleted_at)
		Delete: func(id string) error {
			query := `UPDATE articles SET deleted_at = datetime('now') WHERE id = ? AND deleted_at IS NULL`
			result, err := db.Exec(query, id)
			if err != nil {
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return ErrNotFound
			}
			return nil
		},

		// listbyowner, ambil semua artikel milik user tertentu
		ListByOwner: func(ownerID string) ([]models.Article, error) {
			query := `
				SELECT id, title, author, content, created_at, updated_at, views, owner_id, deleted_at
				FROM articles
				WHERE owner_id = ? AND deleted_at IS NULL
				ORDER BY created_at DESC
			`
			rows, err := db.Query(query, ownerID)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var articles []models.Article
			for rows.Next() {
				var a models.Article
				err := rows.Scan(
					&a.ID, &a.Title, &a.Author, &a.Content,
					&a.CreatedAt, &a.UpdatedAt, &a.Views, &a.OwnerID, &a.DeletedAt,
				)
				if err != nil {
					return nil, err
				}
				articles = append(articles, a)
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}
			return articles, nil
		},
	}, nil
}
