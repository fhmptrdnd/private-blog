package main

import (
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"

	"github.com/fhmptrdnd/private-blog/internal/handler"
	"github.com/fhmptrdnd/private-blog/internal/repository"
	"github.com/fhmptrdnd/private-blog/internal/service"
)

func main() {
	// initialize sqlite database
	repo, err := repository.NewSQLiteRepo("blog.db")
	if err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		return
	}
	
	// pake function types, bukan struct
	// clock sama idgen ini function yang bisa dipanggil
	clock := service.NewRealClock()  // function buat dapetin waktu
	idGen := service.NewRealIDGen()  // function buat generate id
	
	svc := service.NewArticleService(repo, clock, idGen)
	h := handler.NewHandler(svc)

	// WithLogging: log setiap request
	// WithPanicRecovery: tangkap panic biar server ga crash
	
	// routes yang cuma butuh GET
	http.HandleFunc("/", handler.Chain(h.Home, handler.WithLogging, handler.WithPanicRecovery))
	http.HandleFunc("/my-articles", handler.Chain(h.MyArticles, handler.WithLogging, handler.WithPanicRecovery))
	http.HandleFunc("/view/", handler.Chain(h.View, handler.WithLogging, handler.WithPanicRecovery))
	http.HandleFunc("/edit/", handler.Chain(h.Edit, handler.WithLogging, handler.WithPanicRecovery))
	
	// routes yang butuh POST (dengan method check)
	http.HandleFunc("/create", handler.Chain(h.Create, handler.WithLogging, handler.WithPanicRecovery))
	http.HandleFunc("/update/", handler.Chain(h.Update, handler.WithLogging, handler.WithPanicRecovery))
	http.HandleFunc("/delete/", handler.Chain(h.Delete, handler.WithLogging, handler.WithPanicRecovery))

	fmt.Println("Telegraph running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("server error:", err)
	}
}