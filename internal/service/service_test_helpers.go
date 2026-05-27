package service

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func createInMemoryDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	_, err = db.Exec(`
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
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
