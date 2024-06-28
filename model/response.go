package model

type Response struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type StartChatResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`

	Chat
}

type Chat struct {
	Audio string `json:"audio,omitempty"`
	Text  string `json:"text,omitempty"`
}

type AnswerChatResponse struct {
	Prompt Chat `json:"prompt,omitempty"`
	Answer Chat `json:"answer,omitempty"`
}
