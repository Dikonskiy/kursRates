package healthcheck

import (
	"kursRates/internal/repository"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Health struct {
	log    *slog.Logger
	mu     sync.RWMutex
	ready  bool
	live   bool
	repo   *repository.Repository
	apiUrl string
}

func NewHealth(log *slog.Logger, repo *repository.Repository, apiurl string) *Health {
	return &Health{
		log:    log,
		ready:  true,
		live:   true,
		repo:   repo,
		apiUrl: apiurl,
	}
}

func (h *Health) ReadinessCheck() bool {
	return h.repo.GetDb() != nil
}

func (h *Health) LivenessCheck() bool {
	if err := h.repo.GetDb().Ping(); err != nil {
		return false
	} else if _, err := http.Get(h.apiUrl); err != nil {
		return false
	}
	return true
}

func (h *Health) ReadyHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	h.ready = h.ReadinessCheck()
	h.ready = h.LivenessCheck()

	if !h.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.log.Error("Readiness checker", "Status", "Not Ready")
		w.Write([]byte("Status: Not Ready"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Ready"))
	h.log.Info("Readiness checker", "Status", "Ready")
}

func (h *Health) LiveHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	h.live = h.LivenessCheck()

	if !h.live {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.log.Error("Liveness checker", "Status", "Not Live")
		w.Write([]byte("Status: Not Live"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Live"))
	h.log.Info("Liveness checker", "Status", "Live")
}

func (h *Health) PeriodicHealthCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		liveness := h.LivenessCheck()
		readiness := h.ReadinessCheck()

		if liveness && readiness {
			h.log.Info("Readiness checker", "Status", "Ready")
		} else {
			h.log.Error("Readiness checker", "Status", "Not Ready")
		}

		if liveness {
			h.log.Info("Liveness checker", "Status", "Live")
		} else {
			h.log.Error("Liveness checker", "Status", "Not Live")
		}
	}
}
