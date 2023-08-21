package util

import (
	"database/sql"
	"encoding/json"
	"os"

	"kursRates/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&models.Config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", models.Config.MysqlConnectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
