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

	// routing biasa (bisa pake middleware kalo mau)
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/my-articles", h.MyArticles)
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
