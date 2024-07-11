package initconfig

import (
	"encoding/json"
	"kursRates/internal/models"
	"os"
)

func InitConfig(filename string) (config *models.Config, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
