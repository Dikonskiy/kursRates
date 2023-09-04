package main

import (
	"kursRates/internal/app"
	"kursRates/internal/database"
	"kursRates/internal/httphandler"
	"kursRates/internal/logerr"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func init() {
	log := logerr.InitLogger()

	db, err := database.InitDB()
	if err != nil {
		log.Error("Failed to initialize database:", err)
		return
	}
	defer db.Close()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", httphandler.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", httphandler.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", httphandler.GetCurrencyHandler)

	app.StartServer(r)
}
