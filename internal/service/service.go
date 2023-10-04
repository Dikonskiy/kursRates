package service

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/internal/metrics"
	"kursRates/internal/models"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	Client  *http.Client
	Logger  *slog.Logger
	Metrics *metrics.Metrics
}

func NewService(Logger *slog.Logger, metrics *metrics.Metrics) *Service {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Service{
		Client:  client,
		Logger:  Logger,
		Metrics: metrics,
	}
}

func (s *Service) GetData(ctx context.Context, data string, APIURL string) *models.Rates {
	start := time.Now()
	apiURL := fmt.Sprintf("%s?fdate=%s", APIURL, data)

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

	duration := time.Since(start)
	statusCode := fmt.Sprintf("%d", resp.StatusCode)

	s.Metrics.IncRequestCount("URL", statusCode)
	s.Metrics.ObserveRequestDuration("URL", statusCode, duration.Seconds())

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
