package repository

import (
	"database/sql"
	"errors"
	"kursRates/internal/models"
)

func GetData(db *sql.DB, formattedDate, code string) ([]models.DBItem, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
	}

	rows, err := db.Query(query, params...)
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
		return nil, errors.New("no data found with these parameters")
	}

	return results, nil
}
