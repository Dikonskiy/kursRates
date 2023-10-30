package healthcheck

import (
	"kursRates/internal/repository"
	"net/http"
	"sync"
	"time"
)

type Health struct {
	mu     sync.RWMutex
	ready  bool
	live   bool
	Repo   *repository.Repository
	APIURL string
}

func NewHealth(repo *repository.Repository, apiurl string) *Health {
	return &Health{
		ready:  true,
		live:   true,
		Repo:   repo,
		APIURL: apiurl,
	}
}

func (h *Health) ReadinessCheck() bool {
	return h.Repo.Db != nil
}

func (h *Health) LivenessCheck() bool {
	if err := h.Repo.Db.Ping(); err != nil {
		return false
	} else if _, err := http.Get(h.APIURL); err != nil {
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
		h.Repo.Logerr.Error("Readiness checker", "Status", "Not Ready")
		w.Write([]byte("Status: Not Ready"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Ready"))
	h.Repo.Logerr.Info("Readiness checker", "Status", "Ready")
}

func (h *Health) LiveHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	h.live = h.LivenessCheck()

	if !h.live {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.Repo.Logerr.Error("Liveness checker", "Status", "Not Live")
		w.Write([]byte("Status: Not Live"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Live"))
	h.Repo.Logerr.Info("Liveness checker", "Status", "Live")
}

func (h *Health) PeriodicHealthCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		liveness := h.LivenessCheck()
		readiness := h.ReadinessCheck()

		if liveness && readiness {
			h.Repo.Logerr.Info("Readiness checker", "Status", "Ready")
		} else {
			h.Repo.Logerr.Error("Readiness checker", "Status", "Not Ready")
		}

		if liveness {
			h.Repo.Logerr.Info("Liveness checker", "Status", "Live")
		} else {
			h.Repo.Logerr.Error("Liveness checker", "Status", "Not Live")
		}
	}
}
