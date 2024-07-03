package repository

import (
	"context"
	"database/sql"
	"kursRates/internal/metrics"
	"kursRates/internal/models"
	"kursRates/internal/service"
	"log/slog"
	"strconv"
	"time"
)

type Repository struct {
	db      *sql.DB
	logerr  *slog.Logger
	metrics *metrics.Metrics
}

func NewRepository(MysqlConnectionString string, logerr *slog.Logger, metrics *metrics.Metrics) *Repository {
	db, err := sql.Open("mysql", MysqlConnectionString)
	if err != nil {
		logerr.Error("Failed initialize database connection")
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
		db:      db,
		logerr:  logerr,
		metrics: metrics,
	}
}

func (r *Repository) InsertData(rates models.Rates, formattedDate string) {
	savedItemCount := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
	defer cancel()

	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			r.logerr.Error("Failed to convert float: %s", err)
			continue
		}

		startTime := time.Now()

		rows, err := r.db.QueryContext(ctx, "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)", item.Title, item.Code, value, formattedDate)
		if err != nil {
			r.logerr.Error("Failed to insert in the database:", err)
		} else {
			savedItemCount++
			r.logerr.Info("Item saved",
				"count", savedItemCount,
			)
		}
		defer rows.Close()

		duration := time.Since(startTime).Seconds()
		go r.metrics.ObserveInsertDuration("insert", "success", duration)
	}
	r.logerr.Info("Items saved:",
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

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	duration := time.Since(startTime).Seconds()
	if code == "" {
		go r.metrics.ObserveSelectDuration("select", "success", duration)
		go r.metrics.IncSelectCount("select", "success")
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
		r.logerr.Error("No data found with these parameters")
	}

	return results, nil
}

func (r *Repository) scheduler(ctx context.Context, formattedDate string, rates models.Rates) error {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM R_CURRENCY WHERE A_DATE = ?", formattedDate).Scan(&count)
	if err != nil {
		return err
	}

	for _, item := range rates.Items {
		value, errr := strconv.ParseFloat(item.Value, 64)
		if errr != nil {
			r.logerr.Error("Failed to convert float:", errr)
			continue
		}
		if count > 0 {
			_, err = r.db.ExecContext(ctx, "UPDATE R_CURRENCY SET TITLE = ?, VALUE = ?, U_DATE = NOW() WHERE A_DATE = ? AND CODE = ?", item.Title, value, formattedDate, item.Code)
			if err != nil {
				return err
			}
		} else {
			_, err = r.db.ExecContext(ctx, "INSERT INTO R_CURRENCY (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, ?)", item.Title, item.Code, value, formattedDate)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (r *Repository) HourTick(date, formattedDate string, ctx context.Context, APIURL string) {

	var service = service.NewService(r.logerr, r.metrics)

	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		err := r.scheduler(ctx, formattedDate, *service.GetData(ctx, date, APIURL))
		if err != nil {
			r.logerr.Error("Can't update the date:", err)
		}
	}
}

func (r *Repository) DeleteData(ctx context.Context, formattedDate, code string) ([]models.DBItem, error) {
	var query string
	var params []interface{}

	if code == "" {
		query = "DELETE FROM R_CURRENCY WHERE A_DATE = ?"
		params = []interface{}{formattedDate}
	} else {
		query = "DELETE FROM R_CURRENCY WHERE A_DATE = ? AND CODE = ?"
		params = []interface{}{formattedDate, code}
	}

	startTime := time.Now()

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	duration := time.Since(startTime).Seconds()
	if code == "" {
		go r.metrics.ObserveSelectDuration("select", "success", duration)
		go r.metrics.IncSelectCount("select", "success")
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
		r.logerr.Error("No data found with these parameters")
	}

	return results, nil
}

func (r *Repository) GetDb() *sql.DB {
	return r.db
}
