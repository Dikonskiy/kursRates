package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type DBItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Code  string `json:"code"`
	Value string `json:"value"`
	Date  string `json:"date"`
}

type Rates struct {
	XMLName xml.Name       `xml:"rates"`
	Items   []CurrencyItem `xml:"item"`
	Date    string         `xml:"date"`
}

type CurrencyItem struct {
	Title string `xml:"fullname"`
	Code  string `xml:"title"`
	Value string `xml:"description"`
}

func saveToDatabase(db *sql.DB, stmt *sql.Stmt, items []CurrencyItem, formattedDate string, done chan<- bool) {
	for _, item := range items {
		_, err := stmt.Exec(item.Title, item.Code, item.Value, formattedDate)
		if err != nil {
			log.Println("Error inserting data:", err)
		}
	}
	done <- true
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config.json:", err)
	}
	defer configFile.Close()

	var config struct {
		ListenPort            string `json:"listenPort"`
		MysqlConnectionString string `json:"mysqlConnectionString"`
	}

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config.json:", err)
	}

	db, err := sql.Open("mysql", config.MysqlConnectionString)
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

		parsedDate, err := time.Parse("02.01.2006", date)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		formattedDate := parsedDate.Format("2006-01-02")

		apiURL := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s", date)

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

		var rates Rates
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
		go saveToDatabase(db, stmt, rates.Items, formattedDate, done)
		<-done

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	})

	r.HandleFunc("/currency/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dateString := vars["date"]
		code := vars["code"]

		date, err := time.Parse("02.01.2006", dateString)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		formattedDate := date.Format("2006-01-02")

		query := "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params := []interface{}{formattedDate}
		params = append(params, code)

		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []DBItem
		for rows.Next() {
			var item DBItem
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
		dateStr := vars["date"]

		date, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		formattedDate := date.Format("2006-01-02")
		query := "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params := []interface{}{formattedDate}

		rows, err := db.Query(query, params...)
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []DBItem
		for rows.Next() {
			var item DBItem
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
		Addr:         ":" + config.ListenPort,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	log.Println("Listening on port", config.ListenPort, "...")
	log.Fatal(server.ListenAndServe())
}
