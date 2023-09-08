package repository

import (
	"database/sql"
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"strconv"
)

var (
	Repo   *Repository
	logger = logerr.InitLogger()
)

func AddData(db *sql.DB) (*sql.Stmt, error) {
	stmt, err := Repo.Db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
	if err != nil {
		logger.Error("Failed to prepare statement", err)
		return nil, err
	}
	return stmt, nil
}

func InsertData(db *sql.DB, rates models.Rates, formattedDate string) {
	stmt, err := AddData(Repo.Db)
	if err != nil {
		logger.Error("Failed to prepare data", err)
		return
	}
	defer stmt.Close()

	savedItemCount := 0

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			logger.Error("Failed to convert float: %s", err.Error())
			continue
		}

		_, err = stmt.Exec(item.Title, item.Code, value, formattedDate)
		if err != nil {
			logger.Error("Failed to insert in the database:", err.Error())
		} else {
			savedItemCount++
			logger.Info("Item saved", savedItemCount)
		}
	}
	logger.Info("Items saved:", savedItemCount)
}
