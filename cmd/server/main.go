package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"prx/internal/handlers"
	"prx/internal/logger"
	"prx/internal/utils"
)

func main() {
	var redirectPairs utils.RedirectFlag

	flag.Var(&redirectPairs, "record", "Specify redirect records in the format FROM=TO. Multiple records can be specified.")
	flag.Parse()

	if err := utils.ParseRedirectRecords(redirectPairs); err != nil {
		logger.Log.Error("Failed to parse redirect records:", "error", err)
		return
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", handlers.HandleRequests)
	router.HandleFunc("GET /api/status", handlers.StatusHandler)

	middlewareStack := handlers.MiddlewareStack(handlers.LoggingMiddleware)

	server := &http.Server{
		Addr:    ":443",
		Handler: middlewareStack(router),
		TLSConfig: &tls.Config{
			GetCertificate: func(req *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return utils.GetOrCreateCertificate(req.ServerName)
			},
		},
	}

	// Listen and serve with TLS
	log.Println("Starting HTTPS server on port 443")
	err := server.ListenAndServeTLS("", "") // Cert and key are provided by the TLSConfig
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
