package handlers

import (
	"encoding/json"
	"net/http"
	"prx/internal/entities"
	"strings"
	"time"

	l "prx/internal/logger"

	"context"

	"prx/internal/utils"
)

func HandleRequest(w http.ResponseWriter, req *http.Request) {

	if utils.RedirectRecords[req.Host] == "" || RedirectRecords[req.Host] == nil {
		l.Log.Error("Redirect record does not exist", RedirectRecords[req.Host])
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	toURL := RedirectRecords[req.Host]
	http.Redirect(w, req, toURL, http.StatusFound)
}

func StatusHandler(w http.ResponseWriter, req *http.Request) {

	res := struct {
		Status string `json:"status"`
		Time   string `json:"time"`
	}{
		Status: "OK",
		Time:   time.Now().Format(time.RFC3339),
	}

	// Add any status checks you need here
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		l.Log.Error("Failed to encode records", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
