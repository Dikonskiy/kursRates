package main

import (
	"kursRates/internal/app"
	"kursRates/internal/httphandler"
	"kursRates/internal/logerr"
	"kursRates/internal/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Repo *repository.Repository

func init() {
	logger := logerr.InitLogger()

	db, err := repository.InitDB()
	if err != nil {
		logger.Error("Failed to initialize database:", err)
		return
	}

	Repo = repository.NewRepository(db)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", httphandler.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", httphandler.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", httphandler.GetCurrencyHandler)

	app.StartServer(r)
}
