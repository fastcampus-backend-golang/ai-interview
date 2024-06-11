package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/madeindra/interview-ai/ai"
)

var client = ai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))

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

		r.Post("/answer", answerChat)
	})

	// health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"healthy\":true}"))
	})

	return r
}

func initialChat(w http.ResponseWriter, req *http.Request) {
	initialText, err := ai.GetInitialText()
	if err != nil {
		log.Printf("failed to get initial text: %v", err)
		http.Error(w, "failed to get initial text", http.StatusInternalServerError)
		return
	}

	initialAudio, err := ai.GetInitialAudio()
	if err != nil {
		log.Printf("failed to get initial audio: %v", err)
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
		log.Printf("failed to marshal response: %v", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func answerChat(w http.ResponseWriter, req *http.Request) {
	// read audio as multipart
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	if fileHeader == nil {
		log.Println("no file uploaded")
		http.Error(w, "no file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// transcribe the audio
	transcript, err := client.Transcribe(file, fileHeader.Filename)
	if err != nil {
		log.Printf("failed to transcribe audio: %v", err)
		http.Error(w, "failed to transcribe audio", http.StatusInternalServerError)
		return
	}

	// get chat completion
	chatCompletion, err := client.Chat(transcript.Text)
	if err != nil {
		log.Printf("failed to get chat completion: %v", err)
		http.Error(w, "failed to get chat completion", http.StatusInternalServerError)
		return
	}

	if len(chatCompletion.Choices) == 0 {
		log.Println("no chat completion")
		http.Error(w, "no chat completion", http.StatusInternalServerError)
		return
	}

	// create speech from the chat completion
	speech, err := client.TextToSpeech(chatCompletion.Choices[0].Message.Content)
	if err != nil {
		log.Printf("failed to create speech: %v", err)
		http.Error(w, "failed to create speech", http.StatusInternalServerError)
		return
	}

	// encode the speech to base64
	speechByte, err := io.ReadAll(speech)
	if err != nil {
		log.Printf("failed to read speech: %v", err)
		http.Error(w, "failed to read speech", http.StatusInternalServerError)
		return
	}
	speechBase64 := base64.StdEncoding.EncodeToString(speechByte)

	// send response
	response := ChatResponse{
		Prompt: Response{
			Text: transcript.Text,
		},
		Answer: Response{
			Text:  chatCompletion.Choices[0].Message.Content,
			Audio: speechBase64,
		},
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	// write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
