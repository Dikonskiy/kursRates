package main

import (
	"kursRates/connections"
	"kursRates/internal/app"
	"kursRates/util"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config.json:", err)
	}
	defer configFile.Close()

	util.InitLogger()

	util.InitDB()

	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", connections.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", connections.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", connections.GetCurrencyHandler)

	app.StartServer(r)
}
