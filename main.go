package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fastcampus-backend-golang/ai-interview/handler"
)

var (
	port   = os.Getenv("PORT")
	apiKey = os.Getenv("OPENAI_API_KEY")
	dbURI  = os.Getenv("DB_URI")
)

func main() {
	// pastikan semua variabel yang dibutuhkan tersedia
	if err := validateEnv(); err != nil {
		log.Fatal(err)
	}

	// buat handler
	router := handler.NewHandler(apiKey, dbURI)

	// buat server
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
		Handler: router,
	}

	// jalankan server
	log.Printf("Server listening on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// validateEnv digunakan untuk memvalidasi variabel yang dibutuhkan
func validateEnv() error {
	if port == "" {
		port = "8080"
	}

	if dbURI == "" {
		return errors.New("DB_URI is required")
	}

	if apiKey == "" {
		return errors.New("OPENAI_API_KEY is required")
	}

	return nil
}
