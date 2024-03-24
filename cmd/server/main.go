package main

import (
	l "prx/internal/logger"
	"net/http"
	"prx/internal/handlers"
	"prx/internal/utils"
	"strings"
  "flag"
)

func main() {

  var records utils.RecordsFlag

  flag.Var(&records, "record", "Add a record in the format FROM=TO. Multiple values can be specified.")
  flag.Parse()

  for _, record := range records {
	parts := strings.SplitN(record, "=", 2)
		if len(parts) != 2 {
			l.Log.Info("Invalid record format:", record)
			continue
    	}
		utils.RedirectRecords[parts[0]] = parts[1]
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /", handlers.HandleRequest)
	router.HandleFunc("GET /api/status", handlers.StatusHandler)

	stack := handlers.CreateStack(
		handlers.Logging,
	)

	srv := &http.Server{
		Addr:         ":80",
		Handler:      stack(router),
	}

	l.Log.Fatal(srv.ListenAndServe())
}
