package handler

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/madeindra/interview-ai/ai"
	"github.com/madeindra/interview-ai/data"
	"github.com/madeindra/interview-ai/model"
)

type handler struct {
	ai ai.Client
	db data.Client
}

func NewHandler(apiKey string, dbURI string) *chi.Mux {
	h := &handler{
		ai: ai.NewOpenAI(apiKey),
		db: data.NewMongo(dbURI),
	}

	r := chi.NewRouter()

	// gunakan middleware CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Access-Key"},
	}))

	// sajikan direktori static ke /public
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))

	// rute untuk homepage
	r.Get("/", h.Homepage)

	// rute untuk chat
	r.Get("/chat/start", h.StartChat)
	r.Post("/chat/answer", h.AnswerChat)

	return r
}

func sendResponse(w http.ResponseWriter, data any, message string, status int) {
	resp, err := json.Marshal(model.Response{
		Message: message,
		Data:    data,
	})
	if err != nil {
		log.Printf("failed to marshal response: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "an error occured while processing the request"}`))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func (h *handler) Homepage(w http.ResponseWriter, req *http.Request) {
	pagePath := path.Join("page", "index.html")
	http.ServeFile(w, req, pagePath)
}

func (h *handler) StartChat(w http.ResponseWriter, req *http.Request) {
	initialText, err := ai.GetInitialText()
	if err != nil {
		log.Printf("failed to get initial text: %v", err)
		sendResponse(w, nil, "failed to get initial text", http.StatusInternalServerError)

		return
	}

	entry := data.ChatEntry{
		ID:     "", // TODO: generate new
		Secret: "", // TODO: generate new
		History: []ai.ChatMessage{
			{
				Role:    ai.ROLE_SYSTEM,
				Content: initialText,
			},
		},
	}

	// TODO: store to database

	initialAudio, err := ai.GetInitialAudio()
	if err != nil {
		log.Printf("failed to get initial audio: %v", err)
		sendResponse(w, nil, "failed to get initial audio", http.StatusInternalServerError)

		return
	}

	initialChat := data.InitialChatData{
		ID:     entry.ID,
		Secret: entry.Secret,
		Chat: data.Chat{
			Text:  entry.History[0].Content,
			Audio: initialAudio,
		},
	}

	sendResponse(w, initialChat, "a new chat created", http.StatusOK)
}

func (h *handler) AnswerChat(w http.ResponseWriter, req *http.Request) {
	// TODO: validate id and secret

	// TODO: fetch data from database

	// TODO: match secret

	// read audio as multipart
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		log.Printf("failed to read file: %v", err)
		sendResponse(w, nil, "failed to read file", http.StatusInternalServerError)

		return
	}
	if fileHeader == nil {
		log.Println("required file is missing")
		sendResponse(w, nil, "required file is missing", http.StatusBadRequest)

		return
	}
	defer file.Close()

	// transcribe the audio
	transcript, err := h.ai.Transcribe(file, fileHeader.Filename)
	if err != nil {
		log.Printf("failed to transcribe audio: %v", err)
		sendResponse(w, nil, "failed to transcribe audio", http.StatusInternalServerError)

		return
	}

	if transcript.Text == "" {
		log.Println("cannot complete audio transcription: no transcript")
		sendResponse(w, nil, "cannot complete audio transcription", http.StatusInternalServerError)

		return
	}

	// TODO: append the transcript to the chat history
	chatHistory := []ai.ChatMessage{}

	// get chat completion
	chatCompletion, err := h.ai.Chat(chatHistory)
	if err != nil {
		log.Printf("failed to get chat completion: %v", err)
		sendResponse(w, nil, "failed to get chat completion", http.StatusInternalServerError)

		return
	}

	if len(chatCompletion.Choices) == 0 {
		log.Println("cannot complete chat completion: no chat completion")
		sendResponse(w, nil, "cannot complete chat completion", http.StatusInternalServerError)

		return
	}

	// create speech from the chat completion
	speech, err := h.ai.TextToSpeech(chatCompletion.Choices[0].Message.Content)
	if err != nil {
		log.Printf("failed to create speech: %v", err)
		sendResponse(w, nil, "failed to create speech", http.StatusInternalServerError)

		return
	}

	// encode the speech to base64
	speechByte, err := io.ReadAll(speech)
	if err != nil {
		log.Printf("failed to read speech: %v", err)
		sendResponse(w, nil, "failed to read speech", http.StatusInternalServerError)

		return
	}
	speechBase64 := base64.StdEncoding.EncodeToString(speechByte)

	// TODO: update chat history in database

	// send response
	response := data.ChatData{
		Prompt: data.Chat{
			Text: transcript.Text,
		},
		Answer: data.Chat{
			Text:  chatCompletion.Choices[0].Message.Content,
			Audio: speechBase64,
		},
	}

	sendResponse(w, response, "success", http.StatusOK)
}
