// internal/app/server.go
package app

import (
	"kursRates/internal/models"
	"kursRates/util"
	"log"
	"net/http"
	"time"
)

func StartServer(router http.Handler) {
	server := &http.Server{
		Addr:         ":" + models.Config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	log.Println("Listening on port", models.Config.ListenPort, "...")
	util.Info.Println("Listening on port", models.Config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
