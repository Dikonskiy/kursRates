package repository

import (
	"context"
	"database/sql"
	"kursRates/internal/logerr"
	"kursRates/internal/metrics"
	"kursRates/internal/models"
	"log/slog"
	"strconv"
	"time"
)

type Repository struct {
	Db      *sql.DB
	Logerr  *slog.Logger
	Metrics *metrics.Metrics
}

func NewRepository(MysqlConnectionString string, logerr *logerr.Logerr, metrics *metrics.Metrics) *Repository {
	db, err := sql.Open("mysql", MysqlConnectionString)
	if err != nil {
		logerr.Logerr.Error("Failed initialize database connection")
		return nil
	}

	db.SetMaxOpenConns(39)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil
	}

	return &Repository{
		Db:      db,
		Logerr:  logerr.Logerr,
		Metrics: metrics,
	}
}

func (r *Repository) InsertData(rates models.Rates, formattedDate string) {
	savedItemCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			r.Logerr.Error("Failed to convert float: %s", err)
			continue
		}

		startTime := time.Now()

		rows, err := r.Db.QueryContext(ctx, "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)", item.Title, item.Code, value, formattedDate)
		if err != nil {
			r.Logerr.Error("Failed to insert in the database:", err)
		} else {
			savedItemCount++
			r.Logerr.Info("Item saved",
				"count", savedItemCount,
			)
		}
		defer rows.Close()

		duration := time.Since(startTime).Seconds()
		go r.Metrics.ObserveInsertDuration("insert", "success", duration)
	}
	r.Logerr.Info("Items saved:",
		"All", savedItemCount,
	)
}

func (r *Repository) GetData(ctx context.Context, formattedDate, code string) ([]models.DBItem, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "SELECT ID, TITLE, CODE, VALUE, A_DATE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
	}

	startTime := time.Now()

	rows, err := r.Db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	duration := time.Since(startTime).Seconds()
	if code == "" {
		go r.Metrics.ObserveSelectDuration("select", "success", duration)
		go r.Metrics.IncSelectCount("select", "success")
	}

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

func (r *Repository) scheduler(ctx context.Context, formattedDate string, rates models.Rates) error {
	var count int
	err := r.Db.QueryRowContext(ctx, "SELECT COUNT(*) FROM R_CURRENCY WHERE A_DATE = ?", formattedDate).Scan(&count)
	if err != nil {
		return err
	}

	for _, item := range rates.Items {
		value, errr := strconv.ParseFloat(item.Value, 64)
		if errr != nil {
			r.Logerr.Error("Failed to convert float:", errr)
			continue
		}
		if count > 0 {
			_, err = r.Db.ExecContext(ctx, "UPDATE R_CURRENCY SET TITLE = ?, VALUE = ?, U_DATE = NOW() WHERE A_DATE = ? AND CODE = ?", item.Title, value, formattedDate, item.Code)
			if err != nil {
				return err
			}
		} else {
			_, err = r.Db.ExecContext(ctx, "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)", item.Title, item.Code, value, formattedDate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Repository) HourTick(ctx context.Context, formattedDate string, rates models.Rates) {
	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		err := r.scheduler(ctx, formattedDate, rates)
		if err != nil {
			r.Logerr.Error("Can't update the date:", err)
		}
	}

}
