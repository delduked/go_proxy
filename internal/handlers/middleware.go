package handlers

import (
	"net/http"
)

// validateUpdateRequest ensures that the request is a POST and Content-Type is application/json.
func ValidateUpdatePostRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
			http.Error(w, "Invalid Content-Type. Only application/json is allowed", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, req)
	}
}

// validateUpdateRequest ensures that the request is a POST and Content-Type is application/json.
func ValidateUpdateGetRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	}
}
