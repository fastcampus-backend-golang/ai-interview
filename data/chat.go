package data

import "github.com/fastcampus-backend-golang/ai-interview/ai"

type ChatEntry struct {
	ID      string `bson:"_id"`
	Secret  string
	History []ai.ChatMessage
}
