package model

type ChatResponse struct {
	Prompt Response `json:"prompt,omitempty"`
	Answer Response `json:"answer,omitempty"`
}

type Response struct {
	Audio string `json:"audio,omitempty"`
	Text  string `json:"text,omitempty"`
}
