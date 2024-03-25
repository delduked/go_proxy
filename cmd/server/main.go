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
	  parts := strings.Split(record, "=")
		if len(parts) != 2 {
			l.Log.Info("Invalid record format:", record)
			continue
    	}
    l.Log.Info(record)
   
    parts[0] = strings.ReplaceAll(parts[0], "=", "")
  	parts[1] = strings.ReplaceAll(parts[1], "=", "")

		utils.RedirectRecords[parts[0]] = parts[1]
	}

	router := http.NewServeMux()
	router.HandleFunc("/", handlers.HandleRequest)
	router.HandleFunc("/api/status", handlers.StatusHandler)

  stack := handlers.CreateStack(
    handlers.Logging,
  )

	srv := &http.Server{
		Addr:         ":80",
		Handler:      stack(router),
	}

  l.Log.Info("Server started...")
  if err := srv.ListenAndServe(); err != nil {
    l.Log.Fatal("server failed: ", err)
  }

}
