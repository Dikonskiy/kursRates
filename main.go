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

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	r.HandleFunc("/currency/save/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.SaveCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.GetCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.GetCurrencyHandler(w, r.WithContext(ctx))
	})

	app.StartServer(r, Logger, Cnfg)
}
