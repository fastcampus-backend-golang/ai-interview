package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"github.com/fastcampus-backend-golang/ai-interview/model"
	"golang.org/x/crypto/bcrypt"
)

func sendResponse(w http.ResponseWriter, data any, message string, status int) {
	// marshal sebagai JSON
	resp, err := json.Marshal(model.Response{
		Message: message,
		Data:    data,
	})
	// jika error, kirim response error
	if err != nil {
		log.Printf("failed to marshal response: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "an error occured while processing the request"}`))

		return
	}

	// kirim response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func generateRandom() string {
	// karakter & panjang yang digunakan
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	// buat random string
	random := make([]byte, length)
	for i := range random {
		random[i] = charset[rand.Intn(len(charset))]
	}

	return string(random)
}

func createHash(plain string) (string, error) {
	// generate hash dari password
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHash(plain, hash string) error {
	// bandingkan hash dengan password
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

func sanitizeString(text string) string {
	// hapus tebal
	reStrong := regexp.MustCompile(`\*\*([^*]+)\*\*`)
	text = reStrong.ReplaceAllString(text, "$1")

	// hapus miring
	reItalic := regexp.MustCompile(`\*([^*]+)\*`)
	text = reItalic.ReplaceAllString(text, "$1")

	// Remove links
	reLink := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	text = reLink.ReplaceAllString(text, "$1")

	// hapus new line
	reBullet := regexp.MustCompile(`\n- `)
	text = reBullet.ReplaceAllString(text, ", ")

	// ganti new line dengan spasi
	text = strings.Replace(text, "\n", " ", -1)

	return text
}
