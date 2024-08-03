package util

import (
	"context"
	"fmt"
	"log"
	db "shin-monta-no-mori/internal/db/sqlc"

	"github.com/lib/pq"
)

func SeedingForDev(store *db.Store) {
	queries := []string{
		fmt.Sprintln(`
		INSERT INTO operators (name, hashed_password, email)
		VALUES ('test', '$2a$10$deT1QvqgziN0hZ15v/zh1O6kIGFf7xyJOg1w4i4Ukmmxshd41DhMa', 'test@test.com');
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			// Check if the error is a PostgreSQL duplicate key violation
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				log.Printf("Duplicate entry found, skipping: %v", err)
				continue
			} else {
				log.Fatalf("Failed to exec query: %v", err)
			}
		}
	}
}
