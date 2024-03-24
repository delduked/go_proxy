package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"

	"prx/internal/entities"
	l "prx/internal/logger"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
)

func InitRedisClient() {

	var REDIS_PASS string
	if ok := os.Getenv("REDIS_PASS"); ok == "" {
		l.Log.Error("REDIS_PASS is not set. Using default value")
		REDIS_PASS = "your_password"
	} else {
		REDIS_PASS = os.Getenv("REDIS_PASS")
	}

	var REDIS_PORT string
	if ok := os.Getenv("REDIS_PORT"); ok == "" {
		l.Log.Error("REDIS_PORT is not set. Using default value")
		REDIS_PORT = "6379"
	} else {
		REDIS_PORT = os.Getenv("REDIS_PORT")
	}

	var REDIS_ADDR string
	if ok := os.Getenv("REDIS_ADDR"); ok == "" {
		l.Log.Error("REDIS_ADDR is not set. Using default value")
		REDIS_ADDR = "localhost"
	} else {
		REDIS_ADDR = os.Getenv("REDIS_ADDR")
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR + ":" + REDIS_PORT,
		Password: REDIS_PASS,
		DB:       0,
	})

	l.Log.Info("* Redis client initialized * ")
	l.Log.Info("*")
	l.Log.Info("*")
	l.Log.Info("* Pass: ", REDIS_PASS)
	l.Log.Info("* Addr: ", REDIS_ADDR)
	l.Log.Info("* Port: ", REDIS_PORT)
	l.Log.Info("*")
	l.Log.Info("*")
}

// UpdateResult holds the results of the update operations, including successes and failures.
type UpdateResult struct {
	Successes []string
	Failures  []string
}

func UpdateRedirectRecord(ctx context.Context, entries []entities.RedirectEntry) UpdateResult {
	var result UpdateResult

	for _, entry := range entries {
		hashedFrom := HashValue(entry.From)

		value, err := json.Marshal(entry)
		if err != nil {
			l.Log.Error("Failed to marshal redirect entry", "from", entry.From, "err", err)
			result.Failures = append(result.Failures, entry.From) // Include more details as needed
			continue
		}

		err = Rdb.Set(ctx, hashedFrom, string(value), 0).Err()
		if err != nil {
			l.Log.Error("Failed to save redirect entry in the database", "hashedFrom", hashedFrom, "err", err)
			result.Failures = append(result.Failures, entry.From) // Include more details as needed
			continue
		}

		result.Successes = append(result.Successes, entry.From)
	}

	return result
}

func GetAllRedirectRecords(ctx context.Context) ([]entities.RedirectEntry, error) {
	var entries []entities.RedirectEntry
	var cursor uint64
	var err error

	for {
		var keys []string
		keys, cursor, err = Rdb.Scan(ctx, cursor, "redirect:*", 0).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			val, err := Rdb.Get(ctx, key).Result()
			if err != nil {
				l.Log.Error("Failed to fetch redirect entry", "key", key, "err", err)
				continue // Optionally, handle this more gracefully
			}

			var entry entities.RedirectEntry
			if err := json.Unmarshal([]byte(val), &entry); err != nil {
				l.Log.Error("Failed to unmarshal redirect entry", "val", val, "err", err)
				continue // Optionally, handle this more gracefully
			}

			entries = append(entries, entry)
		}

		if cursor == 0 {
			break
		}
	}

	return entries, nil
}

// HashValue takes a string value and returns a SHA-256 hash as a hex string.
func HashValue(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed)
}
