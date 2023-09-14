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

var Repo *repository.Repository
var Hand *httphandler.Handler
var Logger *logerr.Logerr

func init() {
	err := initconfig.InitConfig("config.json")
	if err != nil {
		Repo.Logerr.Error("Failed to initialize the configuration:", err)
		return
	}
	Repo = repository.NewRepository(models.Config.MysqlConnectionString)
	Hand = httphandler.NewHandler(models.Config.MysqlConnectionString)
	Logger = logerr.NewLogerr(models.Config.IsProd)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", Hand.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", Hand.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", Hand.GetCurrencyHandler)

	app.StartServer(r)
}	
