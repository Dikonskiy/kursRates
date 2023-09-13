// internal/httphandler/handlers.go
package httphandler

import (
	"encoding/json"
	"kursRates/internal/repository"
	"kursRates/internal/service"
	"log"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	R *repository.Repository
}

func NewHandler(mysqlConnectionString string) *Handler {
	repo := repository.NewRepository(mysqlConnectionString)
	if repo == nil {
		log.Fatal("Failed to initialize the repository")
	}
	return &Handler{
		R: repo,
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

func (h *Handler) RespondWithError(w http.ResponseWriter, status int, errorMsg string, err error) {
	http.Error(w, errorMsg, status)
	h.R.Logerr.Error(errorMsg+": ", err)
}

func (h *Handler) SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse and format the date", err)
		return
	}

	rates, err := service.Service(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse", err)
	}

	go h.R.InsertData(rates, formattedDate)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	h.R.Logerr.Info("Success: true")
}

func (h *Handler) GetCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]
	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.R.GetData(formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}
	h.R.Logerr.Info("Data was showed")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
