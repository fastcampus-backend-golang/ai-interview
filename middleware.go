package main

import (
	"net/http"
)

// accessKeyMiddleware digunakan untuk memvalidasi access key dari header
func accessKeyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Access-Key") != accessKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
