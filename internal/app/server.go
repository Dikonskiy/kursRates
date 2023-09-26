// internal/app/server.go
package app

import (
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"log"
	"log/slog"
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

	quit := make(chan os.Signal, 1)

	go shutdown(quit, *a.Logerr.Logerr)

	log.Println("Listening on port", config.ListenPort, "...")
	a.Logerr.Logerr.Info("Listening on port", config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}

func shutdown(quit chan os.Signal, logger slog.Logger) {
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info("caught signal", map[string]string{
		"signal": s.String(),
	})
	os.Exit(0)
}
