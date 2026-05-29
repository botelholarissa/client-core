package repository

import (
	"database/sql"
	"client-core/internal/models"
)

type ClientRepository struct {
	DB *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{DB: db}
}

func (r *ClientRepository) CreateClient(client models.Client) error {
	query := `
		INSERT INTO clients (
			name,
			email,
			assets,
			status,
			priority
		)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.DB.Exec(
		query,
		client.Name,
		client.Email,
		client.Assets,
		client.Status,
		client.Priority,
	)

	return err
}

func (r *ClientRepository) UpdateClient(client models.Client) error {
	query := `
		UPDATE clients
		SET
			status = ?,
			priority = ?
		WHERE email = ?
	`

	_, err := r.DB.Exec(
		query,
		client.Status,
		client.Priority,
		client.Email,
	)

	return err
}

func (r *ClientRepository) FindByEmail(email string) (*models.Client, error) {
	query := `
		SELECT
			id,
			name,
			email,
			assets,
			status,
			priority
		FROM clients
		WHERE email = ?
	`

	row := r.DB.QueryRow(query, email)

	var client models.Client

	err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.Assets,
		&client.Status,
		&client.Priority,
	)

	if err != nil {
		return nil, err
	}

	return &client, nil
}