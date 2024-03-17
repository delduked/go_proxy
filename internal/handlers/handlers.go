package handlers

import (
	"encoding/json"
	"net/http"
	"prx/internal/db"
	"prx/internal/entities"
	"strings"
	"time"

	l "prx/internal/logger"

	"context"

	"prx/internal/utils"

	"github.com/redis/go-redis/v9"
)

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	utils.Acquire()
	defer utils.Release()

	hashedDomain := db.HashValue(req.Host) // Use the same hash function

	value, err := db.Rdb.Get(ctx, hashedDomain).Result()
	if err == redis.Nil {
		http.NotFound(w, req)
		return
	} else if err != nil {
		l.Log.Error("Error retrieving redirect record", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Assuming you're storing the full 'uuid|~|from_value|~|to_value' string,
	// but you could adjust this to just retrieve and redirect to the 'to_value'.
	parts := strings.Split(value, "|~|")
	if len(parts) < 3 {
		l.Log.Error("Stored value has an unexpected format", "value", value)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	toURL := parts[2]
	http.Redirect(w, req, toURL, http.StatusFound)
}

func UpdateRedirectRecord(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	var entries []entities.RedirectEntry
	if err := json.NewDecoder(req.Body).Decode(&entries); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.UpdateRedirectRecord(ctx, entries)

	// Construct a detailed response based on the update results
	var response entities.UpdateRecordsResponse

	if len(result.Failures) > 0 {
		l.Log.Info("Some redirect records failed to update", "result", result)
		w.WriteHeader(http.StatusPartialContent) // Use a status code that indicates partial success/failure
	} else {
		l.Log.Info("All redirect records updated successfully")
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		l.Log.Error("Failed to encode response", "err", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func GetAllRedirectRecordsHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	entries, err := db.GetAllRedirectRecords(ctx) // Ensure db.Rdb is your initialized Redis client
	if err != nil {
		l.Log.Error("Failed to retrieve all redirect records", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		l.Log.Error("Failed to encode records", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
