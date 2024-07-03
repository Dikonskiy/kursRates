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
	client  *http.Client
	logger  *slog.Logger
	metrics *metrics.Metrics
}

func NewService(logger *slog.Logger, metrics *metrics.Metrics) *Service {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Service{
		client:  client,
		logger:  logger,
		metrics: metrics,
	}
}

func (s *Service) GetData(ctx context.Context, data string, APIURL string) models.Rates {
	start := time.Now()
	apiURL := fmt.Sprintf("%s?fdate=%s", APIURL, data)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		s.logger.Error("Failed to create request with context", err)
		return models.Rates{}
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to GET URL", err)
		return models.Rates{}
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	statusCode := fmt.Sprintf("%d", resp.StatusCode)

	go s.metrics.IncRequestCount("URL", statusCode)
	go s.metrics.ObserveRequestDuration("URL", statusCode, duration.Seconds())

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to Read response Body", err)
		return models.Rates{}
	}

	var rates models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		s.logger.Error("Failed to parse XML data", err)
		return models.Rates{}
	}

	return rates
}
