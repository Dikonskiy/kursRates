package repository

import (
	"database/sql"
	"kursRates/internal/database"
	"kursRates/internal/logerr"
)

var (
	db, _  = database.InitDB()
	logger = logerr.InitLogger()
)

func AddData() (*sql.Stmt, error) {
	stmt, err := db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
	if err != nil {
		logger.Error("Failed to to prepare date", err)
	}
	return stmt, nil
}
