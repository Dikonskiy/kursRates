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
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var (
	Repo   *repository.Repository
	Hand   *httphandler.Handler
	Logger *logerr.Logerr
	Cnfg   *models.Config
	App    *app.Application
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
	App = app.NewApplication(Logger)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.SaveCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	r.HandleFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.GetCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	r.HandleFunc("/currency/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		Hand.GetCurrencyHandler(w, r.WithContext(ctx), ctx)
	})

	App.StartServer(r, Cnfg)
}
