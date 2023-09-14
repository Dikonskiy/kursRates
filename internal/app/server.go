// internal/app/server.go
package app

import (
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"log"
	"net/http"
	"time"
)

func StartServer(router http.Handler, Logerr *logerr.Logerr, config *models.Config) {
	server := &http.Server{
		Addr:         ":" + config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	log.Println("Listening on port", config.ListenPort, "...")
	Logerr.Logerr.Info("Listening on port", config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
