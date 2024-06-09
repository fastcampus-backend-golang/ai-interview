package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewHandler() *chi.Mux {
	r := chi.NewRouter()

	// gunakan middleware CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Access-Key"},
	}))

	// gunakan middleware access key pada api-api berikut
	r.Group(func(r chi.Router) {
		r.Use(accessKeyMiddleware)
	})

	// health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"healthy\":true}"))
	})

	return r
}
