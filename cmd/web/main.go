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
    svc := service.NewArticleService(repo)
    h := handler.NewHandler(svc)

    http.HandleFunc("/", h.Home)
    http.HandleFunc("/create", h.Create)
    http.HandleFunc("/view/", h.View)
    http.HandleFunc("/edit/", h.Edit)
    http.HandleFunc("/update/", h.Update)
    http.HandleFunc("/delete/", h.Delete)

    fmt.Println("📝 Telegraph Clone berjalan di http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("server error:", err)
    }
}
