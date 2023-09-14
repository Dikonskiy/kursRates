package repository

import (
	"database/sql"
	"kursRates/internal/logerr"
	"kursRates/internal/models"
	"log/slog"
	"strconv"
)

type Repository struct {
	Db     *sql.DB
	Logerr *slog.Logger
}

func NewRepository(MysqlConnectionString string, logerr *logerr.Logerr) *Repository {
	db, err := sql.Open("mysql", MysqlConnectionString)
	if err != nil {
		logerr.Logerr.Error("Failed initialize database connection")
		return nil
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil
	}

	return &Repository{
		Db:     db,
		Logerr: logerr.Logerr,
	}
}

func (r *Repository) AddData() (*sql.Stmt, error) {
	stmt, err := r.Db.Prepare("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)")
	if err != nil {
		r.Logerr.Error("Failed to prepare statement", err)
		return nil, err
	}
	return stmt, nil
}

func (r *Repository) InsertData(rates models.Rates, formattedDate string) {
	stmt, err := r.AddData()
	if err != nil {
		r.Logerr.Error("Failed to prepare data", err)
		return
	}
	defer stmt.Close()

	savedItemCount := 0

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			r.Logerr.Error("Failed to convert float: %s", err.Error())
			continue
		}

		_, err = stmt.Exec(item.Title, item.Code, value, formattedDate)
		if err != nil {
			r.Logerr.Error("Failed to insert in the database:", err.Error())
		} else {
			savedItemCount++
			r.Logerr.Info("Item saved", savedItemCount)
		}
	}
	r.Logerr.Info("Items saved:", savedItemCount)
}

func (r *Repository) GetData(formattedDate, code string) ([]models.DBItem, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
	}

	rows, err := r.Db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DBItem
	for rows.Next() {
		var item models.DBItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Code, &item.Value, &item.Date); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if len(results) == 0 {
		r.Logerr.Error("no data found with these parameters")
	}

	return results, nil
}
