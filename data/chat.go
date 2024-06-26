package data

import "github.com/fastcampus-backend-golang/ai-interview/ai"

type ChatEntry struct {
	ID      string `bson:"_id"`
	Secret  string
	History []ai.ChatMessage
}

type InitialChatData struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`

	Chat
}

type ChatData struct {
	Prompt Chat `json:"prompt,omitempty"`
	Answer Chat `json:"answer,omitempty"`
}

type Chat struct {
	Audio string `json:"audio,omitempty"`
	Text  string `json:"text,omitempty"`
}
