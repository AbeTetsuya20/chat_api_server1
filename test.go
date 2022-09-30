package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/api/:id")

		reqID := strings.TrimPrefix(r.URL.Path, "/api/")
		fmt.Println("reqID: ", reqID)

		// header を取得
		header_name := r.Header.Get("name")
		header_address := r.Header.Get("address")
		header_password := r.Header.Get("password")

		fmt.Println("header_name:", header_name)
		fmt.Println("header_address:", header_address)
		fmt.Println("header_password:", header_password)

		data := map[string]string{
			"message": "hello",
			"request": reqID,
		}
		render.JSON(w, r, data)
	})

	addr := os.Getenv("Addr")
	if addr == "" {
		addr = ":4444"
	}

	log.Printf("listen: %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("!! %+v", err)
	}
}
