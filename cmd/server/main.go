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

	server := &http3.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{ /* Your TLS certificates */ },
			NextProtos:   []string{"quic"},
		},
		Handler: mux,
	}

	l.Log.Info("Starting HTTP/3 server...")
	if err := server.ListenAndServeTLS("./nated.site.crt", "./nated.site.key"); err != nil {
		l.Log.Fatal("Failed to start server", "err", err)
	}
}
