package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/dateutil"
	"kursRates/models"
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

	r.HandleFunc("/currency/save/{date}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		formattedDate := dateutil.DateFormat(date)

		apiURL := fmt.Sprintf("%s?fdate=%s", models.Config.APIURL, date)

		resp, err := http.Get(apiURL)
		if err != nil {
			http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		xmlData, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read XML response", http.StatusInternalServerError)
			return
		}

		var rates models.Rates
		if err := xml.Unmarshal(xmlData, &rates); err != nil {
			http.Error(w, "Failed to parse XML", http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		done := make(chan bool)
		go dateutil.SaveToDatabase(db, stmt, rates.Items, formattedDate, done)
		<-done

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	})

	r.HandleFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		code := vars["code"]

		formattedDate := dateutil.DateFormat(date)

		query := "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params := []interface{}{formattedDate}
		params = append(params, code)

		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []models.DBItem
		for rows.Next() {
			var item models.DBItem
			if err := rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			results = append(results, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	r.HandleFunc("/currency/{date}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]

		formattedDate := dateutil.DateFormat(date)
		query := "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params := []interface{}{formattedDate}

		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []models.DBItem
		for rows.Next() {
			var item models.DBItem
			if err := rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			results = append(results, item)
			if results == nil {
				http.Error(w, "we don't have a data with this parameters", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	server := &http.Server{
		Addr:         ":" + models.Config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	log.Println("Listening on port", models.Config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
