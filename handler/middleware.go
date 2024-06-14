package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
)

type contextKey string

const (
	contextKeyUserID     contextKey = "user-id"
	contextKeyUserSecret contextKey = "user-secret"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessKey := r.Header.Get("Authorization")

		if accessKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if accessKey[:6] != "Basic " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessKey = accessKey[6:]

		decoded, err := base64.StdEncoding.DecodeString(accessKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(string(decoded), ":")
		if len(parts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := parts[0]
		userSecret := parts[1]

		r = r.WithContext(context.WithValue(r.Context(), contextKeyUserID, userID))
		r = r.WithContext(context.WithValue(r.Context(), contextKeyUserSecret, userSecret))

		next.ServeHTTP(w, r)
	})
}
