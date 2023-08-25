package main

import (
	"kursRates/internal/app"
	util "kursRates/internal/database"
	"kursRates/internal/httphandler"
	logerr "kursRates/internal/logerr"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func init() {
	logerr.InitLogger()

	db, err := util.InitDB()
	if err != nil {
		logerr.Error.Println("Failed to initialize database:", err)
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
