package repository

import (
	"database/sql"
)

type WebhookRepository struct {
	DB *sql.DB
}

func NewWebhookRepository(db *sql.DB) *WebhookRepository {
	return &WebhookRepository{DB: db}
}

func (r *WebhookRepository) IsProcessed(eventID string) (bool, error) {
	query := `SELECT event_id FROM processed_events WHERE event_id = ?`
	row := r.DB.QueryRow(query, eventID)
	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *WebhookRepository) MarkProcessed(eventID string, cardID string) error {
	query := `INSERT OR IGNORE INTO processed_events (event_id, pipefy_card_id) VALUES (?, ?)`
	_, err := r.DB.Exec(query, eventID, cardID)
	return err
}
