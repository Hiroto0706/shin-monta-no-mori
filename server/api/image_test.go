package api_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	model "shin-monta-no-mori/server/internal/domains/models"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	_ "github.com/lib/pq"
)

func TestListImages(t *testing.T) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}
	setUp(t, config)
	// defer tearDown(t, config)

	type ListIllustrationsParams struct {
		Limit  int
		Offset int
	}

	tests := []struct {
		name    string
		arg     ListIllustrationsParams
		want    []model.Illustration
		wantErr bool
	}{
		{
			name: "正常系",
			arg: ListIllustrationsParams{
				Limit:  config.ImageFetchLimit,
				Offset: 0,
			},
			want: []model.Illustration{
				{
					Image: db.Image{
						ID:          10001,
						Title:       "test_image_title_10001",
						OriginalSrc: "test_image_original_src_10001.com",
						SimpleSrc: sql.NullString{
							String: "test_image_simple_src_10001.com",
							Valid:  true,
						},
					},
					Character: []db.Character{
						{
							ID:   10001,
							Name: "test_character_name_10001",
							Src:  "test_character_src_10001.com",
						},
					},
					ParentCategory: []db.ParentCategory{
						{
							ID:   10001,
							Name: "test_parent_category_name_10001",
							Src:  "test_parent_category_src_10001.com",
						},
					},
					ChildCategory: []db.ChildCategory{
						{
							ID:   10001,
							Name: "test_child_category_name_10001",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	log.Println(tests)
}

func createConn(config util.Config) *db.Store {
	conn, err := sql.Open(config.DBDriver, config.TestDBUrl)
	if err != nil {
		log.Fatal("cannot connect to test db: ", err)
	}

	return db.NewStore(conn)
}

func setUp(t *testing.T, config util.Config) {
	store := createConn(config)

	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(10001, 'test_image_title_10001', 'test_image_original_src_10001.com', 'test_image_simple_src_10001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(10001, 'test_character_name_10001', 'test_character_src_10001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_characters_relations (id, image_id, character_id)
		VALUES
		(10001, 10001, 10001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(10001, 'test_parent_category_name_10001', 'test_parent_category_src_10001.com');
		`),
		fmt.Sprintln(`
		INSERT INTO image_parent_categories_relations (id, image_id, parent_category_id)
		VALUES
		(10001, 10001, 10001);
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(10001, 'test_child_category_name_10001', 10001);
		`),
	}

	for _, query := range queries {
		if _, err := store.ExecQuery(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}
}

func tearDown(t *testing.T, config util.Config) {
	store := createConn(config)

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
