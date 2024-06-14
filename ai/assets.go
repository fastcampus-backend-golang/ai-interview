package ai

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
)

type ChatAsset struct {
	SystemPrompt string
	ChatText     string
	ChatAudio    string
}

func GetChatAsset() (ChatAsset, error) {
	prompt, err := getSystemPrompt()
	if err != nil {
		return ChatAsset{}, err
	}

	text, err := getInitialText()
	if err != nil {
		return ChatAsset{}, err
	}

	audio, err := getInitialAudio()
	if err != nil {
		return ChatAsset{}, err
	}

	return ChatAsset{
		SystemPrompt: prompt,
		ChatText:     text,
		ChatAudio:    audio,
	}, nil
}

func getInitialAudio() (string, error) {
	basedir, err := os.Getwd()
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	assetPath := filepath.Join(basedir, "ai", "assets", "initial.mp3")

	audio, err := os.ReadFile(assetPath)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(audio)

	return encoded, nil
}

func getInitialText() (string, error) {
	basedir, err := os.Getwd()
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	assetPath := filepath.Join(basedir, "ai", "assets", "initial.txt")

	text, err := os.ReadFile(assetPath)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	return string(text), nil
}

func getSystemPrompt() (string, error) {
	basedir, err := os.Getwd()
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	assetPath := filepath.Join(basedir, "ai", "assets", "system.txt")

	text, err := os.ReadFile(assetPath)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}

	return string(text), nil
}
