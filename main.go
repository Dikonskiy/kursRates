package main

import (
	"database/sql"
	"encoding/json"
	"kursRates/connections"
	"kursRates/internal/models"
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

	err = json.NewDecoder(configFile).Decode(&models.Config)
	if err != nil {
		log.Fatal("Error decoding config.json:", err)
	}
	db, err := sql.Open("mysql", models.Config.MysqlConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

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
	log.Fatal(server.ListenAndServe())
}
