package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/internal/models"
	"log/slog"
	"net/http"
)

type Service struct {
	Rates  *models.Rates
	Logger *slog.Logger
}

func NewService(data string, Config *models.Config, logerr *slog.Logger) (*Service, error) {
	apiURL := fmt.Sprintf("%s?fdate=%s", Config.APIURL, data)

	resp, err := http.Get(apiURL)
	if err != nil {
		logerr.Error("Failed to GET URL", err)
		return nil, err
	}
	defer resp.Body.Close()

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		logerr.Error("Failed to Read response Body", err)
		return nil, err
	}

	var rates *models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		logerr.Error("Failed to parse XML data", err)
		return nil, err
	}

	return &Service{
		Rates:  rates,
		Logger: logerr,
	}, nil
}
