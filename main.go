package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/madeindra/interview-ai/handler"
)

var (
	port   = os.Getenv("PORT")
	apiKey = os.Getenv("OPENAI_API_KEY")
)

func main() {
	// pastikan semua variabel yang dibutuhkan tersedia
	if err := validateEnv(); err != nil {
		log.Fatal(err)
	}

	// buat handler
	router := handler.NewHandler(apiKey)

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

	if apiKey == "" {
		return errors.New("OPENAI_API_KEY is required")
	}

	return nil
}
