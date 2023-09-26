package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/internal/models"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	Client *http.Client
	Logger *slog.Logger
}

func NewService(Logger *slog.Logger) *Service {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Service{
		Client: client,
		Logger: Logger,
	}
}

func (s *Service) GetData(ctx context.Context, data string, Config *models.Config) *models.Rates {
	apiURL := fmt.Sprintf("%s?fdate=%s", Config.APIURL, data)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		s.Logger.Error("Failed to create request with context", err)
		return nil
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error("Failed to GET URL", err)
		return nil
	}
	defer resp.Body.Close()

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error("Failed to Read response Body", err)
		return nil
	}

	var rates *models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		s.Logger.Error("Failed to parse XML data", err)
		return nil
	}

	return rates
}
