package main

type ChatResponse struct {
	Prompt Response `json:"prompt"`
	Answer Response `json:"answer"`
}

type Response struct {
	Audio string `json:"audio"`
	Text  string `json:"text"`
}
