package main

import (
	"kursRates/connections"
	"kursRates/internal/app"
	"kursRates/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func init() {
	util.InitLogger()

	db, err := util.InitDB()
	if err != nil {
		util.Error.Println("Failed to initialize database:", err)
		return
	}
	defer db.Close()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", connections.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", connections.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", connections.GetCurrencyHandler)

	app.StartServer(r)
}
