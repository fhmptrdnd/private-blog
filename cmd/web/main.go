package main

import (
	"fmt"
	"net/http"

	"github.com/fhmptrdnd/weather-api-test-web-based/internal/handler"
	"github.com/fhmptrdnd/weather-api-test-web-based/internal/repository"
	"github.com/fhmptrdnd/weather-api-test-web-based/internal/service"
)

func main() {
    repo := repository.NewMemoryRepo()
    
    // pake function types, bukan struct
    // clock sama idgen ini function yang bisa dipanggil
    clock := service.NewRealClock()  // function buat dapetin waktu
    idGen := service.NewRealIDGen()  // function buat generate id
    
    svc := service.NewArticleService(repo, clock, idGen)
    h := handler.NewHandler(svc)

    // routing biasa (bisa pake middleware kalo mau)
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
