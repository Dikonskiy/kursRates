package repository

import (
	"database/sql"
	"kursRates/internal/models"
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

func InitDB(repo *Repository) error {
	Db, err := GetDB()
	if err != nil {
		return err
	}
	repo.Db = Db
	return nil
}
