// internal/logerr/logger.go
package logerr

import (
	"log"
	"log/slog"
	"os"
)

type Logerr struct {
	Logerr *slog.Logger
}

func NewLogerr(isProd bool) *Logerr {
	infoLogFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}
	var infoLog *slog.Logger
	if isProd {
		infoLog = slog.New(slog.NewJSONHandler(infoLogFile, nil))
	} else {
		infoLog = slog.New(slog.NewTextHandler(infoLogFile, nil))
	}
	return &Logerr{
		Logerr: infoLog,
	}
}
