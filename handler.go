package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/madeindra/interview-ai/ai"
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

		r.Get("/initial", initialChat)

		r.Post("/answer", func(w http.ResponseWriter, r *http.Request) {
			// read audio as multipart

			// transcribe the audio

			// get chat completion

			// create speech from the chat completion

			// send response
		})
	})

	// health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"healthy\":true}"))
	})

	return r
}

func initialChat(w http.ResponseWriter, r *http.Request) {
	initialText, err := ai.GetInitialText()
	if err != nil {
		http.Error(w, "failed to get initial text", http.StatusInternalServerError)
		return
	}

	initialAudio, err := ai.GetInitialAudio()
	if err != nil {
		http.Error(w, "failed to get initial audio", http.StatusInternalServerError)
		return
	}

	response := ChatResponse{
		Answer: Response{
			Text:  initialText,
			Audio: initialAudio,
		},
	}

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
