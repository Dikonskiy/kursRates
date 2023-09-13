package initconfig

import (
	"encoding/json"
	"kursRates/internal/models"
	"os"
)

func InitConfig(configFilePath string) error {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&models.Config)
	if err != nil {
		return err
	}

	return nil
}
