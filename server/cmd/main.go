package main

import (
	"database/sql"
	"log"
	"shin-monta-no-mori/server/api"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	log.Println(config.DBUrl, config.MigrationURL)

	// DBのマイグレーションを実行
	runDBMigration(config.MigrationURL, config.DBUrl)

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot connect to server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance :", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to ru migrate up : ", err)
	}

	log.Println("db migrated successfully")
}
