// internal/httphandler/handlers.go
package httphandler

import (
	"context"
	"encoding/json"
	"kursRates/internal/models"
	"kursRates/internal/repository"
	"kursRates/internal/service"
	"kursRates/internal/metrics"
	"log/slog"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	log *slog.Logger
	r    *repository.Repository
	cnfg *models.Config
	metrics *metrics.Metrics
}

func NewHandler(log *slog.Logger, repo *repository.Repository, config *models.Config, metrics *metrics.Metrics) *Handler {
	if repo == nil {
		log.Error("Failed to initialize the repository")
	}

	return &Handler{
		log: log,
		r:    repo,
		cnfg: config,
		metrics: metrics,
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
	h.log.Error(errorMsg+": ", err)
}

// @Summary Save currency data
// @Description Save currency data for a specific date.
// @Tags currency
// @Accept json
// @Param date path string true "Date in DD.MM.YYYY format"
// @Router /currency/save/{date} [post]
func (h *Handler) SaveCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse and format the date", err)
		return
	}

	var service = service.NewService(h.log, h.metrics)

	go h.r.InsertData(*service.GetData(r.Context(), date, h.cnfg.APIURL), formattedDate)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	h.log.Info("Success: true")
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
func (h *Handler) GetCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.r.GetData(r.Context(), formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}
	h.log.Info("Data was showed")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) DeleteCurrencyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.r.DeleteData(r.Context(), formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve data", err)
		return
	}
	h.log.Info("Data was showed")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) StartScheduler(ctx context.Context) {
	date := time.Now().Format("02.01.2006")

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.log.Error("Cannot parse the Data")
	}

	h.r.HourTick(date, formattedDate, ctx, h.cnfg.APIURL)
}