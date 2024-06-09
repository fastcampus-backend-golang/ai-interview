package ai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Client interface {
	Chat(input string) error
	TextToSpeech(input string) error
	Transcribe(audio io.Reader) error
}

type OpenAI struct {
	APIKey          string
	BaseURL         string
	ChatModel       string
	TranscriptModel string
	TTSModel        string
	TTSVoice        string
}

const (
	baseURL         = "https://api.openai.com/v1"
	chatModel       = "gpt-4o"
	transcriptModel = "whisper-1"
	ttsModel        = "tts-1"
	ttsVoice        = "nova"
)

// NewOpenAI digunakan untuk membuat instance client OpenAI
func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		APIKey:          apiKey,
		BaseURL:         baseURL,
		ChatModel:       chatModel,
		TranscriptModel: transcriptModel,
		TTSModel:        ttsModel,
		TTSVoice:        ttsVoice,
	}
}

// Chat digunakan untuk melakukan chat
func (c *OpenAI) Chat(input string) error {
	url, err := url.JoinPath(c.BaseURL, "/engines/chat/completions")
	if err != nil {
		return err
	}

	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"model": "%s",
		"messages": [
			{
				"role": "system",
				"content": "You are an interviewer."
			},
			{
				"role": "user",
				"content": "%s"
			}
		]
	}`, c.ChatModel, input)))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// TextToSpeech digunakan untuk mengubah teks menjadi suara
func (c *OpenAI) TextToSpeech(input string) error {
	url, err := url.JoinPath(c.BaseURL, "/audio/speech")
	if err != nil {
		return err
	}

	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
		"model": "%s",
		"voice": "%s",
		"input": "%s"
	}`, c.TTSModel, c.TTSVoice, input)))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// SpeechToText digunakan untuk mengubah suara menjadi teks
func (c *OpenAI) Transcribe(audio io.Reader) error {
	url, err := url.JoinPath(c.BaseURL, "/audio/transcriptions")
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the audio file to the multipart message
	part, err := writer.CreateFormFile("file", "audio.wav")
	if err != nil {
		panic(err)
	}

	// Copy the audio file to the part
	_, err = io.Copy(part, audio)
	if err != nil {
		panic(err)
	}

	// Add the model field
	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		panic(err)
	}

	// Close the writer to finalize the multipart message
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Create a new HTTP request with the appropriate headers
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
