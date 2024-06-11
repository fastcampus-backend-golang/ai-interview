package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Client interface {
	Chat(string) (ChatResponse, error)
	TextToSpeech(string) (io.ReadCloser, error)
	Transcribe(*os.File) (TranscriptResponse, error)
}

type OpenAI struct {
	APIKey             string
	BaseURL            string
	ChatModel          string
	TranscriptModel    string
	TranscriptLanguage string
	TTSModel           string
	TTSVoice           string
}

const (
	baseURL            = "https://api.openai.com/v1"
	chatModel          = "gpt-4o"
	transcriptModel    = "whisper-1"
	transcriptLanguage = "en"
	ttsModel           = "tts-1"
	ttsVoice           = "nova"

	defaultSystemContent = "You are an interviewer for a backend engineer role focusing on Golang. In this session of interview, focus on exploring the interviewee's professional experience and how they can fit in as a backend engineer. Ask common interview questions like introduction, experience with Golang, experience in backend development, personal weaknesses & strengths, motivation to join the company, leadership experience, problem-solving, conflict resolution, and what they are looking for in their next role."
)

// NewOpenAI digunakan untuk membuat instance client OpenAI
func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		APIKey:             apiKey,
		BaseURL:            baseURL,
		ChatModel:          chatModel,
		TranscriptModel:    transcriptModel,
		TTSModel:           ttsModel,
		TTSVoice:           ttsVoice,
		TranscriptLanguage: transcriptLanguage,
	}
}

// Chat digunakan untuk melakukan chat
func (c *OpenAI) Chat(message string) (ChatResponse, error) {
	url, err := url.JoinPath(c.BaseURL, "/chat/completions")
	if err != nil {
		return ChatResponse{}, err
	}

	chatReq := ChatRequest{
		Model: c.ChatModel,
		Messages: []ChatMessage{
			{
				Role:    ROLE_SYSTEM,
				Content: defaultSystemContent,
			},
			{
				Role:    ROLE_USER,
				Content: message,
			},
		},
	}

	body, err := json.Marshal(chatReq)
	if err != nil {
		return ChatResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return ChatResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChatResponse{}, err
	}

	var chatResp ChatResponse
	err = unmarshalJSONResponse(resp, &chatResp)
	if err != nil {
		return ChatResponse{}, err
	}

	return chatResp, nil
}

// TextToSpeech digunakan untuk mengubah teks menjadi suara
func (c *OpenAI) TextToSpeech(input string) (io.ReadCloser, error) {
	url, err := url.JoinPath(c.BaseURL, "/audio/speech")
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"model": "%s",
		"voice": "%s",
		"input": "%s"
	}`, c.TTSModel, c.TTSVoice, input)))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// SpeechToText digunakan untuk mengubah suara menjadi teks
func (c *OpenAI) Transcribe(file *os.File) (TranscriptResponse, error) {
	if file == nil {
		return TranscriptResponse{}, fmt.Errorf("audio is nil")
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return TranscriptResponse{}, err
	}

	url, err := url.JoinPath(c.BaseURL, "/audio/transcriptions")
	if err != nil {
		return TranscriptResponse{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", info.Name())
	if err != nil {
		return TranscriptResponse{}, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return TranscriptResponse{}, err
	}

	err = writer.WriteField("model", c.TranscriptModel)
	if err != nil {
		return TranscriptResponse{}, err
	}

	err = writer.WriteField("language", c.TranscriptLanguage)
	if err != nil {
		return TranscriptResponse{}, err
	}

	err = writer.Close()
	if err != nil {
		return TranscriptResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return TranscriptResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return TranscriptResponse{}, err
	}

	if resp == nil || resp.Body == nil {
		return TranscriptResponse{}, fmt.Errorf("response is nil")
	}

	var transcriptResp TranscriptResponse
	err = unmarshalJSONResponse(resp, &transcriptResp)
	if err != nil {
		return TranscriptResponse{}, err
	}

	return transcriptResp, nil
}

// getResponseBody digunakan untuk mendapatkan response body dari http.Response
func getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// unmarshalJSONResponse digunakan untuk mengubah response dari byte menjadi struct
func unmarshalJSONResponse(resp *http.Response, v interface{}) error {
	respBody, err := getResponseBody(resp)
	if err != nil {
		return err
	}
	if respBody == nil {
		return fmt.Errorf("response body is nil")
	}
	defer respBody.Close()

	respByte, err := io.ReadAll(respBody)
	if err != nil {
		return err
	}

	return json.Unmarshal(respByte, v)
}
