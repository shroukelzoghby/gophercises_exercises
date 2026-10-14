package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/shroukelzoghby/gophercises_exercises/phone/db"
	"github.com/shroukelzoghby/gophercises_exercises/phone/phone"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dbPath := "phones.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	sqlDB, err := db.Open(dbPath)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := db.Migrate(sqlDB); err != nil {
		return err
	}

	// Start from a clean slate each run so the demo is repeatable.
	if err := db.Reset(sqlDB); err != nil {
		return err
	}

	rawNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	if err := db.Seed(sqlDB, rawNumbers); err != nil {
		return err
	}

	fmt.Println("Before normalization:")
	if err := printNumbers(sqlDB); err != nil {
		return err
	}

	if err := normalizeAndDedupe(sqlDB); err != nil {
		return err
	}

	fmt.Println("\nAfter normalization:")
	return printNumbers(sqlDB)
}

func normalizeAndDedupe(sqlDB *sql.DB) error {
	numbers, err := db.AllNumbers(sqlDB)
	if err != nil {
		return err
	}

	for _, n := range numbers {
		normalized := phone.Normalize(n.Number)
		existing, err := db.FindNumber(sqlDB, normalized)
		if err != nil {
			return err
		}

		// Another row already holds this normalized value — drop the duplicate.
		if existing != nil && existing.ID != n.ID {
			if err := db.DeleteNumber(sqlDB, n.ID); err != nil {
				return err
			}
			continue
		}

		if n.Number != normalized {
			if err := db.UpdateNumber(sqlDB, n.ID, normalized); err != nil {
				return err
			}
		}
	}
	return nil
}

func printNumbers(sqlDB *sql.DB) error {
	numbers, err := db.AllNumbers(sqlDB)
	if err != nil {
		return err
	}
	for _, n := range numbers {
		fmt.Printf("  %d: %s\n", n.ID, n.Number)
	}
	return nil
}
