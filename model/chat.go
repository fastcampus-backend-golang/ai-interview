package model

import "github.com/madeindra/interview-ai/ai"

type ChatEntry struct {
	ID      string
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
