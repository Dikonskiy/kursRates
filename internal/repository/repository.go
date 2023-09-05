package repository

import (
	"database/sql"
	"encoding/json"
	"kursRates/internal/models"
	"os"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(Db *sql.DB) *Repository {
	return &Repository{
		Db: Db,
	}
}

func GetDB() (*sql.DB, error) {
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

func InitDB() (*sql.DB, error) {
	Db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return Db, nil
}
