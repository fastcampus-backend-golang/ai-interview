package ai

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
)

func GetInitialAudio() (string, error) {
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

func GetInitialText() (string, error) {
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
