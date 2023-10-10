package healthcheck

import (
	"kursRates/internal/repository"
	"net/http"
	"sync"
)

type Health struct {
	mu     sync.RWMutex
	ready  bool
	live   bool
	Repo   repository.Repository
	APIURL string
}

func NewHealth(repo *repository.Repository, apiurl string) *Health {
	return &Health{
		ready:  true,
		live:   true,
		Repo:   *repo,
		APIURL: apiurl,
	}
}

func (h *Health) ReadyHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.Repo.Db == nil {
		h.ready = false
	}

	if !h.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.Repo.Logerr.Error("Readiness checker", "Status", "Not Ready")
		w.Write([]byte("Status: Not Ready"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Ready"))
	h.Repo.Logerr.Error("Readiness checker", "Status", "Ready")
}

func (h *Health) LiveHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if err := h.Repo.Db.Ping(); err != nil {
		h.live = false
	}

	if _, err := http.Get(h.APIURL); err != nil {
		h.live = false
	}

	if !h.live {
		w.WriteHeader(http.StatusServiceUnavailable)
		h.Repo.Logerr.Error("Liveness checker", "Status", "Not Live")
		w.Write([]byte("Status: Not Live"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status: Live"))
	h.Repo.Logerr.Error("Liveness checker", "Status", "Live")
}
