// internal/app/server.go
package app

import (
	"context"
	"fmt"
	"kursRates/internal/healthcheck"
	"kursRates/internal/httphandler"
	"kursRates/internal/initconfig"
	"kursRates/internal/logerr"
	"kursRates/internal/metrics"
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

var (
	Metrics *metrics.Metrics
	Repo    *repository.Repository
	Hand    *httphandler.Handler
	Logger  *logerr.Logerr
	Cnfg    *models.Config
	Health  *healthcheck.Health
)

func init() {
	var err error
	Cnfg, err = initconfig.InitConfig("config.json")
	if err != nil {
		Logger.Logerr.Error("Failed to initialize the configuration:", err)
		return
	}

	Metrics = metrics.NewMetrics()
	Logger = logerr.NewLogerr(Cnfg.IsProd)
	Repo = repository.NewRepository(Cnfg.MysqlConnectionString, Logger, Metrics)
	Health = healthcheck.NewHealth(Repo, Cnfg.APIURL)
	Hand = httphandler.NewHandler(Repo, Cnfg)
	go Hand.StartScheduler(context.TODO())
}

func (a *Application) StartServer() {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

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

	r.HandleFunc("/live", Health.LiveHealthCheckHandler)
	r.HandleFunc("/ready", Health.ReadyHealthCheckHandler)
	go Health.PeriodicHealthCheck()

	r.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			fmt.Println("Failed to start the metrics server:", err)
		}
	}()

	server := &http.Server{
		Addr:         ":" + Cnfg.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	quit := make(chan os.Signal, 1)

	go shutdown(quit, *Logger.Logerr)

	log.Println("Listening on port", Cnfg.ListenPort, "...")
	Logger.Logerr.Info("Listening on port", Cnfg.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}

func shutdown(quit chan os.Signal, logger slog.Logger) {
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	logger.Info("caught signal",
		"signal", s.String(),
	)
	os.Exit(0)
}
