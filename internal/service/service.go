package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"kursRates/internal/models"
	"net/http"
)

func Service(date string) (models.Rates, error) {
	rates, err := FetchRatesFromAPI(date)
	if err != nil {
		return models.Rates{}, err
	}
	return rates, nil
}

func FetchRatesFromAPI(date string) (models.Rates, error) {
	apiURL := fmt.Sprintf("%s?fdate=%s", models.Config.APIURL, date)

	resp, err := http.Get(apiURL)
	if err != nil {
		return models.Rates{}, err
	}
	defer resp.Body.Close()

	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Rates{}, err
	}

	var rates models.Rates
	if err := xml.Unmarshal(xmlData, &rates); err != nil {
		return models.Rates{}, err
	}

	return rates, nil
}
