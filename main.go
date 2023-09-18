package main

import (
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/initconfig"
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"kursRates/internal/repository"

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

	app.StartServer(r, Logger, Cnfg)
}
