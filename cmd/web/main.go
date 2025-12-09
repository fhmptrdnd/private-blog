package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"

	"github.com/fhmptrdnd/private-blog/internal/handler"
	"github.com/fhmptrdnd/private-blog/internal/repository"
	"github.com/fhmptrdnd/private-blog/internal/service"
)

func main() {
	// buka database file "blog.db"
	db, err := sql.Open("sqlite", "blog.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// cek koneksi
	if err := db.Ping(); err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	fmt.Println("Sukses terkoneksi ke SQLite (blog.db)")

	// automigration
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
		log.Fatal("Gagal membuat tabel:", err)
	}

	repo := repository.NewSQLiteRepo(db)

	// pake function types, bukan struct
	// clock sama idgen ini function yang bisa dipanggil
	clock := service.RealClock{} // function buat dapetin waktu
	idGen := service.RealIDGenerator{} // function buat generate id

	svc := service.NewArticleService(repo, clock, idGen)
	h := handler.NewHandler(svc)

	// WithLogging: log setiap request
	// WithPanicRecovery: tangkap panic biar server ga crash

	// routes yang cuma butuh GET
	http.HandleFunc("/", Chain(h.Home, WithLogging, WithPanicRecovery))
	http.HandleFunc("/create", Chain(h.Create, WithLogging, WithPanicRecovery))
	http.HandleFunc("/view/", Chain(h.View, WithLogging, WithPanicRecovery))
	http.HandleFunc("/edit/", Chain(h.Edit, WithLogging, WithPanicRecovery))

	// routes yang buth POST (pakai method check)
	http.HandleFunc("/update/", Chain(h.Update, WithLogging, WithPanicRecovery))
	http.HandleFunc("/delete/", Chain(h.Delete, WithLogging, WithPanicRecovery))

	fmt.Println("Telegraph running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Server error:", err)
	}
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// setiap request yg masuk ke terminal nanti dicatat di sini
func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// next handler
		next(w, r)
		
		// log kalau selesai
		log.Printf("[%s] %s took %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func WithPanicRecovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC RECOVERED: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}