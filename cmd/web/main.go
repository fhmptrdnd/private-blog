package main

import (
	"fmt"
	"net/http"

	"github.com/fhmptrdnd/private-blog/internal/handler"
	"github.com/fhmptrdnd/private-blog/internal/repository"
	"github.com/fhmptrdnd/private-blog/internal/service"
)

func main() {
    repo := repository.NewMemoryRepo()
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
