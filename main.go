package main

import (
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"kursRates/internal/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Repo *repository.Repository

func init() {
	logger := logerr.InitLogger()

	// Initialize the configuration
	err := models.InitConfig("config.json")
	if err != nil {
		logger.Error("Failed to initialize the configuration:", err)
		return
	}

	// Initialize the database connection
	db, err := repository.GetDB()
	if err != nil {
		logger.Error("Failed to initialize the database:", err)
		return
	}
	defer db.Close()

	Repo = repository.NewRepository(db)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", httphandler.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", httphandler.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", httphandler.GetCurrencyHandler)

	app.StartServer(r)
}
