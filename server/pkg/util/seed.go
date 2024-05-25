package util

import (
	"context"
	"fmt"
	"log"
	db "shin-monta-no-mori/server/internal/db/sqlc"
)

func Seeding(store *db.Store) {
	queries := []string{
		fmt.Sprintln(`
		INSERT INTO operators (name, hashed_password, email)
		VALUES ('test1', '$2a$10$7tG1j8/S4E3M7O/oiyuHO.P7xx6cK3G3Y6kP8k1HT8wT7J5P8zD36', 'test1@test.com');
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			log.Fatalf("Failed to exec query: %v", err)
		}
	}
}
