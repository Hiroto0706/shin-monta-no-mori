package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"shin-monta-no-mori/server/pkg/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../../")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	conn, err := sql.Open(config.DBDriver, config.TestDBUrl)
	if err != nil {
		log.Fatal("cannot connect to db :", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}

func SetUp(t *testing.T, db DBTX) {
	queries := []string{
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(10001, 'test_image_title_10001', 'test_image_original_src_10001', 'test_image_simple_src_10001');
		`),
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(20001, 'test_image_title_20001', 'test_image_original_src_20001', 'test_image_simple_src_20001'),
		(20002, 'test_image_title_20002', 'test_image_original_src_20002', '');
		`),
		// listの時は最後の値を取得したいので、IDを大きくする
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(99990, 'test_image_title_99990', 'test_image_original_src_99990', 'test_image_simple_src_99990'),
		(99991, 'test_image_title_99991', 'test_image_original_src_99991', 'test_image_simple_src_99991'),
		(99992, 'test_image_title_99992', 'test_image_original_src_99992', 'test_image_simple_src_99992');
		`),
		fmt.Sprintln(`
		INSERT INTO images (id, title, original_src, simple_src)
		VALUES
		(40001, 'test_image_title_40001', 'test_image_original_src_40001', 'test_image_simple_src_40001'),
		(40002, 'test_image_title_40002', 'test_image_original_src_40002', 'test_image_simple_src_40002'),
		(40003, 'test_image_title_40003', 'test_image_original_src_40003', 'test_image_simple_src_40003');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(10001, 'test_character_name_10001', 'test_character_src_10001');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(20001, 'test_character_name_20001', 'test_character_src_20001'),
		(20002, 'test_character_name_20002', '');
		`),
		// listの時は最後の値を取得したいので、IDを大きくする
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(99990, 'test_character_name_99990', 'test_character_src_99990'),
		(99991, 'test_character_name_99991', 'test_character_src_99991'),
		(99992, 'test_character_name_99992', '');
		`),
		fmt.Sprintln(`
		INSERT INTO characters (id, name, src)
		VALUES
		(40001, 'test_character_name_40001', 'test_character_src_40001'),
		(40002, 'test_character_name_40002', 'test_character_src_40002');
		`),
		fmt.Sprintln(`
		INSERT INTO operators (id, name, hashed_password, email)
		VALUES
		(10001, 'test_operator_name_10001', 'testtest', 'test_10001@test.com'),
		(10002, 'test_operator_name_10002', 'testtest', 'test_10002@test.com'),
		(10003, 'test_operator_name_10003', 'testtest', 'test_10003@test.com'),
		(10004, 'test_operator_name_10004', 'testtest', 'test_10004@test.com');
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(10001, 'test_parent_category_name_10001', 'test_parent_category_src_10001');
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(20001, 'test_parent_category_name_20001', 'test_parent_category_src_20001'),
		(20002, 'test_parent_category_name_20002', '');
		`),
		// listの時は最後の値を取得したいので、IDを大きくする
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(99990, 'test_parent_category_name_99990', 'test_parent_category_src_99990'),
		(99991, 'test_parent_category_name_99991', 'test_parent_category_src_99991'),
		(99992, 'test_parent_category_name_99992', 'test_parent_category_src_99992');
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(40001, 'test_parent_category_name_40001', 'test_parent_category_src_40001'),
		(40002, 'test_parent_category_name_40002', 'test_parent_category_src_40002'),
		(40003, 'test_parent_category_name_40003', 'test_parent_category_src_40003');
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(50001, 'test_parent_category_name_50001', 'test_parent_category_src_50001');
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(60001, 'test_parent_category_name_60001', 'test_parent_category_src_60001');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(10001, 'test_child_category_name_10001', 60001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(70001, 'test_parent_category_name_70001', 'test_parent_category_src_70001');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(20001, 'test_child_category_name_20001', 70001);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(80001, 'test_parent_category_name_80001', 'test_parent_category_src_80001'),
		(80002, 'test_parent_category_name_80002', 'test_parent_category_src_80002'),
		(80003, 'test_parent_category_name_80003', 'test_parent_category_src_80003');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(99991, 'test_child_category_name_99991', 80001),
		(99992, 'test_child_category_name_99992', 80002),
		(99993, 'test_child_category_name_99993', 80003);
		`),
		fmt.Sprintln(`
		INSERT INTO parent_categories (id, name, src)
		VALUES
		(90000, 'test_parent_category_name_90000', 'test_parent_category_src_90000'),
		(90001, 'test_parent_category_name_90001', 'test_parent_category_src_90001');
		`),
		fmt.Sprintln(`
		INSERT INTO child_categories (id, name, parent_id)
		VALUES
		(30001, 'test_child_category_name_30001', 90000),
		(30002, 'test_child_category_name_30002', 90000),
		(30003, 'test_child_category_name_30003', 90000);
		`),
	}
	for _, query := range queries {
		if _, err := db.ExecContext(context.Background(), query); err != nil {
			t.Fatalf("Failed to exec query: %v", err)
		}
	}
}

func TearDown(t *testing.T, db DBTX) {
	queries := []string{
		"TRUNCATE TABLE images RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE characters RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE operators RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE parent_categories RESTART IDENTITY CASCADE;",
		"TRUNCATE TABLE child_categories RESTART IDENTITY CASCADE;",
	}
	for _, query := range queries {
		if _, err := db.ExecContext(context.Background(), query); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}
	}
}
