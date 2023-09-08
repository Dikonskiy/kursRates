// internal/httphandler/handlers.go
package httphandler

import (
	"encoding/json"
	"kursRates/internal/repository"
	"kursRates/internal/service"

	"kursRates/internal/logerr"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	Repo   *repository.Repository
	logger = logerr.InitLogger()
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

	formattedDate, err := DateFormat(date)
	if err != nil {
		http.Error(w, "Failed to parse and format the date", http.StatusInternalServerError)
		return
	}

	rates, err := service.Service(date)
	if err != nil {
		logger.Error("Failed to parse: ", err)
	}

	db, err := repository.GetDB()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get the database connection", err)
		return
	}
	defer db.Close()

	go repository.InsertData(db, rates, formattedDate)

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

	data, err := repository.GetData(Repo.Db, formattedDate, code)
	if err != nil {
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
