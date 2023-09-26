// internal/app/server.go
package app

import (
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Logerr *logerr.Logerr
}

func NewApplication(logerr *logerr.Logerr) *Application {
	return &Application{
		Logerr: logerr,
	}
}

func (a *Application) StartServer(router http.Handler, config *models.Config) {
	server := &http.Server{
		Addr:         ":" + config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		a.Logerr.Logerr.Info("caught signal", map[string]string{
			"signal": s.String(),
		})
		os.Exit(0)
	}()

	log.Println("Listening on port", config.ListenPort, "...")
	a.Logerr.Logerr.Info("Listening on port", config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
