package main

import (
	"context"
	"fmt"
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/initconfig"
	"kursRates/internal/logerr"
	"kursRates/internal/metrics"
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"net/http"
	"time"

	_ "kursRates/docs"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	Metrics *metrics.Metrics
	Repo    *repository.Repository
	Hand    *httphandler.Handler
	Logger  *logerr.Logerr
	Cnfg    *models.Config
	App     *app.Application
)

func init() {
	var err error
	Cnfg, err = initconfig.InitConfig("config.json")
	if err != nil {
		Repo.Logerr.Error("Failed to initialize the configuration:", err)
		return
	}

	Metrics = metrics.NewMetrics()
	Logger = logerr.NewLogerr(Cnfg.IsProd)
	Repo = repository.NewRepository(Cnfg.MysqlConnectionString, Logger, Metrics)
	Hand = httphandler.NewHandler(Repo, Cnfg)
	App = app.NewApplication(Logger)
}

// @title Swagger kursRates API
// @version 0.1
// @description A web service that, upon request, collects data from the public API of the national bank and saves the data to the local TEST database
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
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

	r.HandleFunc("/health", Hand.HealthCheckHandler)

	r.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8081", r); err != nil {
			fmt.Println("Failed to start the metrics server:", err)
		}
	}()

	App.StartServer(r, Cnfg)

}
