package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Number struct {
	ID     int
	Number string
}

// Open opens (or creates) a SQLite database at the given path.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("pinging db: %w", err)
	}
	return db, nil
}

// Migrate creates the phone_numbers table if it does not exist.
func Migrate(db *sql.DB) error {
	const stmt = `
CREATE TABLE IF NOT EXISTS phone_numbers (
	id INTEGER PRIMARY KEY,
	value TEXT NOT NULL
);`
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("migrating: %w", err)
	}
	return nil
}

// Seed inserts the given phone numbers into the table.
func Seed(db *sql.DB, numbers []string) error {
	stmt, err := db.Prepare(`INSERT INTO phone_numbers(value) VALUES(?)`)
	if err != nil {
		return fmt.Errorf("preparing insert: %w", err)
	}
	defer stmt.Close()

	for _, n := range numbers {
		if _, err := stmt.Exec(n); err != nil {
			return fmt.Errorf("inserting %q: %w", n, err)
		}
	}
	return nil
}

// Reset clears all rows from phone_numbers.
func Reset(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM phone_numbers`)
	if err != nil {
		return fmt.Errorf("resetting: %w", err)
	}
	return nil
}

// AllNumbers returns every phone number row.
func AllNumbers(db *sql.DB) ([]Number, error) {
	rows, err := db.Query(`SELECT id, value FROM phone_numbers ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("querying numbers: %w", err)
	}
	defer rows.Close()

	var numbers []Number
	for rows.Next() {
		var n Number
		if err := rows.Scan(&n.ID, &n.Number); err != nil {
			return nil, fmt.Errorf("scanning number: %w", err)
		}
		numbers = append(numbers, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating numbers: %w", err)
	}
	return numbers, nil
}

// FindNumber returns the first row with the given phone number value, or nil if none.
func FindNumber(db *sql.DB, number string) (*Number, error) {
	var n Number
	err := db.QueryRow(`SELECT id, value FROM phone_numbers WHERE value = ?`, number).
		Scan(&n.ID, &n.Number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("finding number: %w", err)
	}
	return &n, nil
}

// UpdateNumber sets the value for the given row id.
func UpdateNumber(db *sql.DB, id int, number string) error {
	_, err := db.Exec(`UPDATE phone_numbers SET value = ? WHERE id = ?`, number, id)
	if err != nil {
		return fmt.Errorf("updating number %d: %w", id, err)
	}
	return nil
}

// DeleteNumber removes the row with the given id.
func DeleteNumber(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM phone_numbers WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting number %d: %w", id, err)
	}
	return nil
}
