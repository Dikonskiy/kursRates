package dateutil

import (
	"database/sql"
	"kursRates/internal/models"
	"log"
	"time"
)

func SaveToDatabase(db *sql.DB, stmt *sql.Stmt, items []models.CurrencyItem, formattedDate string, done chan<- bool) {
	for _, item := range items {
		_, err := stmt.Exec(item.Title, item.Code, item.Value, formattedDate)
		if err != nil {
			log.Println("Error inserting data:", err)
		}
	}
	done <- true
}

func DateFormat(date string) string {
	parsedDate, _ := time.Parse("02.01.2006", date)
	formattedDate := parsedDate.Format("2006-01-02")
	return formattedDate
}
