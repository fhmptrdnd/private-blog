package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/fhmptrdnd/private-blog/internal/handler"
	"github.com/fhmptrdnd/private-blog/internal/repository"
	"github.com/fhmptrdnd/private-blog/internal/service"
)

func main() {
	db, err := sql.Open("sqlite", "blog.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Sukses terkoneksi ke SQLite (file: blog.db)")

	queryBuatTabel := `
    CREATE TABLE IF NOT EXISTS articles (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        author TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME NOT NULL,
        views INTEGER DEFAULT 0,
        owner_id TEXT
    );`

	if _, err := db.Exec(queryBuatTabel); err != nil {
		panic(fmt.Sprintf("Gagal membuat tabel: %v", err))
	}

	repo := repository.NewSQLiteRepo(db)
	clock := service.RealClock{}
	idGen := service.RealIDGenerator{}
	svc := service.NewArticleService(repo, clock, idGen)
	h := handler.NewHandler(svc)

	http.HandleFunc("/", h.Home)
	http.HandleFunc("/create", h.Create)
	http.HandleFunc("/view/", h.View)
	http.HandleFunc("/edit/", h.Edit)
	http.HandleFunc("/update/", h.Update)
	http.HandleFunc("/delete/", h.Delete)

	fmt.Println("Telegraph running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}