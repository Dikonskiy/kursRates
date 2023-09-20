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

func (r *Repository) InsertData(rates models.Rates, formattedDate string) {
	savedItemCount := 0

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			r.Logerr.Error("Failed to convert float: %s", err.Error())
			continue
		}

		rows, err := r.Db.Query("INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)", item.Title, item.Code, value, formattedDate)
		if err != nil {
			r.Logerr.Error("Failed to insert in the database:", err.Error())
		} else {
			savedItemCount++
			r.Logerr.Info("Item saved",
				"count", savedItemCount,
			)
		}
		defer rows.Close()
	}
	r.Logerr.Info("Items saved:",
		"All", savedItemCount,
	)
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
		r.Logerr.Error("No data found with these parameters")
	}

	return results, nil
}
