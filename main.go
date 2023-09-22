package main

import (
	"context"
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/initconfig"
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	Repo   *repository.Repository
	Hand   *httphandler.Handler
	Logger *logerr.Logerr
	Cnfg   *models.Config
)

func init() {
	var err error
	Cnfg, err = initconfig.InitConfig("config.json")
	if err != nil {
		Repo.Logerr.Error("Failed to initialize the configuration:", err)
		return
	}
	Logger = logerr.NewLogerr(Cnfg.IsProd)
	Repo = repository.NewRepository(Cnfg.MysqlConnectionString, Logger)
	Hand = httphandler.NewHandler(Repo, Cnfg)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", Hand.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", Hand.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", Hand.GetCurrencyHandler)

	server := &http.Server{
		Addr:    Cnfg.ListenPort,
		Handler: r,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			Logger.Logerr.Error("Server shutdown error:", err)
		} else {
			Logger.Logerr.Info("Server gracefully stopped")
		}
	}()

	app.StartServer(r, Logger, Cnfg)
}
