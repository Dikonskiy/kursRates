// internal/app/server.go
package app

import (
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Repo *repository.Repository
}

func NewApplication(repo *repository.Repository) *Application {
	return &Application{
		Repo: repo,
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
		a.Repo.Logerr.Info("caught signal", map[string]string{
			"signal": s.String(),
		})
		os.Exit(0)
	}()

	log.Println("Listening on port", config.ListenPort, "...")
	a.Repo.Logerr.Info("Listening on port", config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
