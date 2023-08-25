package httphandler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	util "kursRates/internal/database"

	logerr "kursRates/internal/logerr"
	"kursRates/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func DateFormat(date string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return "", err
	}
	formattedDate := parsedDate.Format("2006-01-02")
	return formattedDate, nil
}

func SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	formattedDate, err := DateFormat(date)
	if err != nil {
		logerr.Error.Println("Failed to parse the date:", err)
	}

	apiURL := fmt.Sprintf("%s?fdate=%s", models.Config.APIURL, date)

	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		logerr.Error.Println("Failed to fetch data from API:", err)
		return
	}
	defer resp.Body.Close()

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read XML response", http.StatusInternalServerError)
		logerr.Error.Println("Failed to read XML response", http.StatusInternalServerError)
		return
	}

	var rates models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		http.Error(w, "Failed to parse XML", http.StatusInternalServerError)
		logerr.Error.Println("Failed to parse XML", http.StatusInternalServerError)
		return
	}

	db, err := util.InitDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		logerr.Error.Println("Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		logerr.Error.Println("Failed to prepare statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			logerr.Error.Println("[ ERR ] convert to float", err.Error())
			continue
		}

		_, err = stmt.Exec(item.Title, item.Code, value, formattedDate)
		if err != nil {
			logerr.Error.Println("[ ERR ] insert in database: data item.Title, value:", item.Title, err.Error())
			continue
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	logerr.Info.Println("date was saved")
}

func GetCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]
	formattedDate, err := DateFormat(date)
	if err != nil {
		logerr.Error.Println("Failed to parse the date:", err)
	}

	db, err := util.GetDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		logerr.Error.Println("Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var query string
	var params []interface{}

	if code == "" {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
		logerr.Info.Println("Currency by date accessed")
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
		logerr.Info.Println("Currency by date and code accessed")
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		logerr.Error.Println("Failed to fetch data from API:", err)
		return
	}
	defer rows.Close()

	var results []models.DBItem
	for rows.Next() {
		var item models.DBItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			logerr.Error.Println("Failed to scan row", http.StatusInternalServerError)
			return
		}
		results = append(results, item)
	}

	if len(results) == 0 {
		http.Error(w, "No data found with these parameters", http.StatusInternalServerError)
		logerr.Error.Println("No data found with these parameters", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
