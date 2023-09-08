package repository

import (
	"database/sql"
	"kursRates/internal/models"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Db: db,
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
