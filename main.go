package main

import (
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/models"
	"kursRates/internal/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Repo *repository.Repository
var Hand *httphandler.Handler

func init() {
	err := models.InitConfig("config.json")
	if err != nil {
		Repo.Logerr.Error("Failed to initialize the configuration:", err)
		return
	}

	Repo = repository.NewRepository(models.Config.MysqlConnectionString)
	Hand = httphandler.NewHandler(models.Config.MysqlConnectionString)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", Hand.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", Hand.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", Hand.GetCurrencyHandler)

	app.StartServer(r)
}
