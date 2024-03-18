package main

import (
	"crypto/tls"
	"net/http"
	"prx/internal/db"
	"prx/internal/handlers"

	l "prx/internal/logger"

	"github.com/lucas-clemente/quic-go/http3"
)

func main() {
	db.InitRedisClient()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HandleRequest)
	mux.HandleFunc("/api/update-redirect", handlers.ValidateUpdatePostRequest(handlers.UpdateRedirectRecord))
	mux.HandleFunc("/api/get-redirects", handlers.ValidateUpdateGetRequest(handlers.GetAllRedirectRecordsHandler))
	mux.HandleFunc("/api/status", handlers.StatusHandler)

	cert, err := tls.LoadX509KeyPair("./tls.crt", "./tls.key")
	if err != nil {
		panic(err)
	}

	server := &http3.Server{
		Addr: "0.0.0.0:443",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: mux,
	}

	l.Log.Info("Starting HTTP/3 server...")
	if err := server.ListenAndServe(); err != nil {
    l.Log.Fatal("Failed to start server", "err", err)
	}
}
