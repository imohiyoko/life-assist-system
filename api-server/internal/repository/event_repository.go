package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

// EventLogsテーブルを初期化する関数
func (r *Repository) InitTable() error {
	query := `
    	CREATE TABLE IF NOT EXISTS EventLogs (
    		event_id SERIAL PRIMARY KEY,
    		timestamp TIMESTAMPTZ NOT NULL,
    		event_type VARCHAR(50) NOT NULL,
    		details TEXT
    	);`
	_, err := r.DB.Exec(context.Background(), query)
	return err
}

// イベントをDBに保存する関数
func (r *Repository) SaveEvent(eventType string, details string) error {
	query := "INSERT INTO EventLogs (timestamp, event_type, details) VALUES ($1, $2, $3)"
	_, err := r.DB.Exec(context.Background(), query, time.Now(), eventType, details)
	if err != nil {
		return fmt.Errorf("unable to insert event: %w", err)
	}
	return nil
}
