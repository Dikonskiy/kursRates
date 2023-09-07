package models

import (
	"encoding/json"
	"os"
)

var Config struct {
	ListenPort            string `json:"listenPort"`
	MysqlConnectionString string `json:"mysqlConnectionString"`
	APIURL                string `json:"apiURL"`
}

func InitConfig(configFilePath string) error {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}
