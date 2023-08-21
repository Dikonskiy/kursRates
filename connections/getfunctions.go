package connections

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/internal/models"
	"kursRates/util/dateutil"
	"kursRates/util/logutil"
	"net/http"

	"github.com/gorilla/mux"
)

func SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
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

	db, err := sql.Open("mysql", models.Config.MysqlConnectionString)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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
}

func GetCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]
	formattedDate := dateutil.DateFormat(date)

	db, err := sql.Open("mysql", models.Config.MysqlConnectionString)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var query string
	var params []interface{}

	if code == "" {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
		logutil.Info.Println("Currency by date accessed")
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
		logutil.Info.Println("Currency by date and code accessed")
	}

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

	if len(results) == 0 {
		http.Error(w, "No data found with these parameters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
