// internal/logerr/logger.go
package logerr

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	Log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return Log
}
