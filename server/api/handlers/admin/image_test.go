package handlers_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	_ "github.com/lib/pq"
)

func TestListImages(t *testing.T) {
	setUp(t)
	defer tearDown(t)

	log.Println("test")
}

func createConn() *db.Store {
	config, err := util.LoadConfig("../../../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	conn, err := sql.Open(config.DBDriver, config.TestDBUrl)
	if err != nil {
		log.Fatal("cannot connect to test db: ", err)
	}

	return db.NewStore(conn)
}

func setUp(t *testing.T) {
	store := createConn()

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(10001, 'test_image_title_10001', 'test_image_original_src_10001', 'test_image_simple_src_10001'),
		(10002, 'test_image_title_10002', 'test_image_original_src_10002', 'test_image_simple_src_10002');
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}
}

func tearDown(t *testing.T) {
	store := createConn()

	queries := []string{
		"TRUNCATE TABLE images RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE operators RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE parent_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE child_categories RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
