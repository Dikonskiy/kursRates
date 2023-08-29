package httphandler

import (
	"database/sql"
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

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = util.InitDB()
	if err != nil {
		logerr.Error.Fatalf("Failed to initialize database: %v", err)
	}
}

func DateFormat(date string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return "", err
	}
	formattedDate := parsedDate.Format("2006-01-02")
	return formattedDate, nil
}

func respondWithError(w http.ResponseWriter, status int, errorMsg string, err error) {
	http.Error(w, errorMsg, status)
	logerr.Error.Printf("%s: %v\n", errorMsg, err)
}

func SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	formattedDate, err := DateFormat(date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

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
	itemsChan := make(chan models.CurrencyItem, len(rates.Items))
	done := make(chan bool)

	for _, item := range rates.Items {
		go func(item models.CurrencyItem) {
			itemsChan <- item
		}(item)
	}

	go func() {
		for i := 0; i < len(rates.Items); i++ {
			<-done
		}
		close(itemsChan)
	}()

	stmt, err := db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	savedItemCount := 0

	for item := range itemsChan {
		go func(item models.CurrencyItem) {
			value, err := strconv.ParseFloat(item.Value, 64)
			if err != nil {
				http.Error(w, "Failed to convert float", http.StatusInternalServerError)
				done <- true
				return
			}

			_, err = stmt.Exec(item.Title, item.Code, value, formattedDate)
			if err != nil {
				http.Error(w, "Failed insert in database: data item.Title, value:", http.StatusInternalServerError)
			} else {
				logerr.Info.Printf("%d Item saved", savedItemCount)
				savedItemCount++
			}
			done <- true
		}(item)
	}
	logerr.Info.Printf("%d Items saved\n", savedItemCount)

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
		respondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

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
		respondWithError(w, http.StatusNotFound, "No data found with these parameters", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
