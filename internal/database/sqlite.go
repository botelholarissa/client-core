package database


import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connect() (*sql.DB, error) {
		db, err := sql.Open("sqlite3", "./clientcore.db") 
		if err != nil {
			return nil, err
		}

		err = createTables(db)
		if err != nil {
			return nil, err
		}

		return db, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			assets REAL NOT NULL,
			status TEXT NOT NULL,
			priority TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS processed_events (
    	event_id TEXT PRIMARY KEY,
		pipefy_card_id TEXT,
    	processed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}
