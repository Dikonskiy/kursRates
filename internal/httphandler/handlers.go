// internal/httphandler/handlers.go
package httphandler

import (
	"context"
	"encoding/json"
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"kursRates/internal/service"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	R    *repository.Repository
	Cnfg *models.Config
}

func NewHandler(repo *repository.Repository, config *models.Config) *Handler {
	if repo == nil {
		repo.Logerr.Error("Failed to initialize the repository")
	}
	return &Handler{
		R:    repo,
		Cnfg: config,
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

// @Summary Save currency data
// @Description Save currency data for a specific date.
// @Tags currency
// @Accept json
// @Param date path string true "Date in DD.MM.YYYY format"
// @Router /currency/save/{date} [post]
func (h *Handler) SaveCurrencyHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	date := vars["date"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse and format the date", err)
		return
	}

	var service = service.NewService(h.R.Logerr, h.R.Metrics)

	go h.R.InsertData(*service.GetData(ctx, date, h.Cnfg.APIURL), formattedDate)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	h.R.Logerr.Info("Success: true")
}

// @Summary Get currency data by date
// @Description Get currency data for a specific date.
// @Tags currency
// @Accept json
// @Param date path string true "Date in DD.MM.YYYY format"
// @Router /currency/{date} [get]
// @Summary Get currency data by date and code
// @Description Get currency data for a specific date and currency code.
// @Tags currency
// @Accept json
// @Param code path string true "Currency code (e.g., USD)"
// @Router /currency/{date}/{code} [get]
func (h *Handler) GetCurrencyHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.R.GetData(ctx, formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}
	h.R.Logerr.Info("Data was showed")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// @Summary Check the health status of the application
// @Description Returns the health status of the application, including the database availability.
// @ID health-check
// @Produce  json
// @Success 200 {string} string "Status: Available"
// @Failure 503 {string} string "Status: Not available"
// @Router /health [get]
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if h.R.Db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.R.Logerr.Info("Health checker",
			"Status", "Not available")
		w.Write([]byte("Status: Not available"))
		return
	}

	if err := h.R.Db.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Status: Not available"))
		h.R.Logerr.Info("Health checker",
			"Status", "Not available")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Available"))
	h.R.Logerr.Info("Health checker",
		"Status", "Available")
}
