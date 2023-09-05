// internal/httphandler/handlers.go
package httphandler

import (
	"database/sql"
	"encoding/json"
	"kursRates/internal/repository"
	"kursRates/internal/service"

	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	logger = logerr.InitLogger()
	db     *sql.DB
	err    error
)

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
	logger.Error("%s: %v", errorMsg, err)
}

func SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	rates, err := service.Service(date)
	if err != nil {
		logger.Error("Failed to parse: ", err)
	}

	formattedDate, err := DateFormat(rates.Date)
	if err != nil {
		http.Error(w, "Failed to parse and format the date", http.StatusInternalServerError)
		return
	}

	stmt, err := repository.AddData()
	if err != nil {
		logger.Error("Failed to prepare data", err)
	}

	savedItemCount := 0

	go func(rates models.Rates, stmt *sql.Stmt) {
		defer stmt.Close()
		for _, item := range rates.Items {
			value, err := strconv.ParseFloat(item.Value, 64)
			if err != nil {
				logger.Error("Failed to convert float: %s", err.Error())
				continue
			}

			_, err = stmt.Exec(item.Title, item.Code, value, formattedDate)
			if err != nil {
				logger.Error("Failed insert in database:", err.Error())
			} else {
				logger.Info("Item saved", savedItemCount)
				savedItemCount++
			}
		}
		logger.Info("Items saved", savedItemCount)
	}(rates, stmt)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	logger.Info("date was saved")
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
		logger.Info("Currency by date accessed")
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
		logger.Info("Currency by date and code accessed")
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
		respondWithError(w, http.StatusNotFound, "No data found with these parameters", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
