package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"prx/internal/db"
	"prx/internal/handlers"
)

func main() {

	db.InitRedisClient()

	router := http.NewServeMux()
	router.HandleFunc("GET /", handlers.HandleRequest)
	router.HandleFunc("GET /api/status", handlers.StatusHandler)

	router.HandleFunc("GET /api/records", handlers.GetAllRedirectRecordsHandler)
	router.HandleFunc("POST /api/records", handlers.UpdateRedirectRecord)

	stack := handlers.CreateStack(
		handlers.Logging,
	)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:         ":443",
		Handler:      stack(router),
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.Fatal(srv.ListenAndServeTLS("./aproxynate.io.origin.pem.crt", "./aproxynate.io.origin.pem.key"))
}
