package handler

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"path"

	"github.com/fastcampus-backend-golang/ai-interview/ai"
	"github.com/fastcampus-backend-golang/ai-interview/data"
	"github.com/fastcampus-backend-golang/ai-interview/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// sajikan direktori static ke /public
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))

	// rute untuk homepage
	r.Get("/", h.Homepage)

	// rute untuk chat
	r.Get("/chat/start", h.StartChat)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/chat/answer", h.AnswerChat)
	})

	return r
}

func (h *handler) Homepage(w http.ResponseWriter, req *http.Request) {
	// sajikan halaman index.html
	pagePath := path.Join("page", "index.html")
	http.ServeFile(w, req, pagePath)
}

func (h *handler) StartChat(w http.ResponseWriter, req *http.Request) {
	// ambil teks awal dari AI
	asset, err := ai.GetChatAsset()
	if err != nil {
		log.Printf("failed to get initial text: %v", err)
		sendResponse(w, nil, "failed to get initial text", http.StatusInternalServerError)

		return
	}

	// buat kata sandi
	plainSecret := generateRandom()
	hashed, err := createHash(plainSecret)
	if err != nil {
		log.Printf("failed to create hash: %v", err)
		sendResponse(w, nil, "failed to create hash", http.StatusInternalServerError)

		return
	}

	// buat chat baru
	entry := data.ChatEntry{
		Secret: hashed,
		History: []ai.ChatMessage{
			{
				Role:    ai.ROLE_SYSTEM,
				Content: asset.SystemPrompt,
			},
			{
				Role:    ai.ROLE_ASSISTANT,
				Content: asset.ChatText,
			},
		},
	}

	newID, err := h.db.InsertChat(entry)
	if err != nil {
		log.Printf("failed to create new chat: %v", err)
		sendResponse(w, nil, "failed to create new chat", http.StatusInternalServerError)

		return
	}

	// kirim respons awal
	initialChat := model.StartChatResponse{
		ID:     newID,
		Secret: plainSecret,
		Chat: model.Chat{
			Text:  asset.ChatText,
			Audio: asset.ChatAudio,
		},
	}

	sendResponse(w, initialChat, "a new chat created", http.StatusOK)
}

func (h *handler) AnswerChat(w http.ResponseWriter, req *http.Request) {
	// ambil user ID dan kata sandi dari konteks (diatur oleh middleware)
	userID := req.Context().Value(contextKeyUserID).(string)
	userSecret := req.Context().Value(contextKeyUserSecret).(string)

	// pastikan user ID dan kata sandi tidak kosong
	if userID == "" || userSecret == "" {
		log.Println("user ID or secret is missing")
		sendResponse(w, nil, "missing required authentication", http.StatusUnauthorized)

		return
	}

	// ambil chat entry berdasarkan user ID
	entry, err := h.db.GetChat(userID)
	if err != nil {
		log.Printf("failed to get chat: %v", err)
		sendResponse(w, nil, "failed to get chat", http.StatusInternalServerError)

		return
	}

	// bandingkan kata sandi
	if err := compareHash(userSecret, entry.Secret); err != nil {
		log.Println("invalid user secret")
		sendResponse(w, nil, "invalid user secret", http.StatusUnauthorized)

		return
	}

	// baca file audio dari form
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

	// ubah audio menjadi teks
	transcript, err := h.ai.Transcribe(file, fileHeader.Filename)
	if err != nil {
		log.Printf("failed to transcribe audio: %v", err)
		sendResponse(w, nil, "failed to transcribe audio", http.StatusInternalServerError)

		return
	}

	// pastikan teks tidak kosong
	if transcript.Text == "" {
		log.Println("cannot complete audio transcription: no transcript")
		sendResponse(w, nil, "cannot complete audio transcription", http.StatusInternalServerError)

		return
	}

	// gabungkan teks ke chat history
	chatHistory := append(entry.History, ai.ChatMessage{
		Role:    ai.ROLE_USER,
		Content: transcript.Text,
	})

	// kirim history ke AI
	chatCompletion, err := h.ai.Chat(chatHistory)
	if err != nil {
		log.Printf("failed to get chat completion: %v", err)
		sendResponse(w, nil, "failed to get chat completion", http.StatusInternalServerError)

		return
	}

	// pastikan chat completion tidak kosong
	if len(chatCompletion.Choices) == 0 {
		log.Println("cannot complete chat completion: no chat completion")
		sendResponse(w, nil, "cannot complete chat completion", http.StatusInternalServerError)

		return
	}

	speechInput := sanitizeString(chatCompletion.Choices[0].Message.Content)

	// buat audio dari teks AI
	speech, err := h.ai.TextToSpeech(speechInput)
	if err != nil {
		log.Printf("failed to create speech: %v", err)
		sendResponse(w, nil, "failed to create speech", http.StatusInternalServerError)

		return
	}

	// ubah audio menjadi base64
	speechByte, err := io.ReadAll(speech)
	if err != nil {
		log.Printf("failed to read speech: %v", err)
		sendResponse(w, nil, "failed to read speech", http.StatusInternalServerError)

		return
	}
	speechBase64 := base64.StdEncoding.EncodeToString(speechByte)

	// gabungkan teks AI ke chat history
	chatHistory = append(chatHistory, ai.ChatMessage{
		Role:    ai.ROLE_ASSISTANT,
		Content: chatCompletion.Choices[0].Message.Content,
	})

	// update chat entry
	entry.History = chatHistory
	if err := h.db.UpdateChat(userID, entry); err != nil {
		log.Printf("failed to update chat: %v", err)
		sendResponse(w, nil, "failed to update chat", http.StatusInternalServerError)

		return
	}

	// kirim respons
	response := model.AnswerChatResponse{
		Prompt: model.Chat{
			Text: transcript.Text,
		},
		Answer: model.Chat{
			Text:  chatCompletion.Choices[0].Message.Content,
			Audio: speechBase64,
		},
	}

	sendResponse(w, response, "success", http.StatusOK)
}
