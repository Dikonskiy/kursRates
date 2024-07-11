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
	"kursRates/internal/repository"
	"kursRates/internal/service"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
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

func (a *Application) StartServer() {
	cnfg, err := initconfig.InitConfig("config.json")
	if err != nil {
		return
	}

	metrics := metrics.NewMetrics()
	logger := logerr.Logger(cnfg.IsProd)
	service := service.NewService(logger, metrics)
	repo := repository.NewRepository(cnfg.MysqlConnectionString, logger, metrics, service)
	health := healthcheck.NewHealth(logger, repo.GetDb(), cnfg.APIURL)
	hand := httphandler.NewHandler(logger, repo, cnfg, metrics, service)
	go hand.StartScheduler(context.TODO())

	r := mux.NewRouter()

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	r.HandleFunc("/currency/save/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		hand.SaveCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		hand.GetCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		hand.GetCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/delete/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		hand.DeleteCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/currency/delete/{date}", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(30*time.Second))
		defer cancel()

		hand.DeleteCurrencyHandler(w, r.WithContext(ctx))
	})

	r.HandleFunc("/live", health.LiveHealthCheckHandler)
	r.HandleFunc("/ready", health.ReadyHealthCheckHandler)
	go health.PeriodicHealthCheck()

	r.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			fmt.Println("Failed to start the metrics server:", err)
		}
	}()

	server := &http.Server{
		Addr:         ":" + cnfg.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	quit := make(chan os.Signal, 1)

	go shutdown(quit, *logger)

	log.Println("Listening on port", cnfg.ListenPort, "...")
	logger.Info("Listening on port", cnfg.ListenPort, "...")
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
