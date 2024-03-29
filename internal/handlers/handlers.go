package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"prx/internal/logger"
	"prx/internal/utils"
	"time"
)

func HandleRequests(w http.ResponseWriter, req *http.Request) {
        targetURL, ok := utils.RedirectRecords[req.Host]
        if !ok {
            logger.Log.Error("No redirect record found for host:", "host", req.Host)
            http.Error(w, "Not Found", http.StatusNotFound)
            return
        }

        logger.Log.Info("Proxying request", "host", req.Host, "target", targetURL)

        parsedURL, err := url.Parse(targetURL)
        if err != nil {
            logger.Log.Error("Failed to parse target URL", "target", targetURL, "error", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        proxy := httputil.NewSingleHostReverseProxy(parsedURL)
        proxy.ServeHTTP(w, req)
}

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	response := struct {
		Status string `json:"status"`
		Time   string `json:"time"`
	}{
		Status: "OK",
		Time:   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}