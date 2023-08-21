package main

import (
	"kursRates/connections"
	"kursRates/dbutil"
	"kursRates/internal/models"
	"kursRates/logutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config.json:", err)
	}
	defer configFile.Close()

	logutil.InitLoggers()

	db, err := dbutil.InitDB()
	if err != nil {
		logutil.Error.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/currency/save/{date}", connections.SaveCurrencyHandler)
	r.HandleFunc("/currency/{date}/{code}", connections.GetCurrencyHandler)
	r.HandleFunc("/currency/{date}", connections.GetCurrencyHandler)

	server := &http.Server{
		Addr:         ":" + models.Config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	log.Println("Listening on port", models.Config.ListenPort, "...")
	logutil.Info.Println("Listening on port", models.Config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
