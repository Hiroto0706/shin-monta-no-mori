package main

import (
	"database/sql"
	"log"
	"shin-monta-no-mori/server/api"
	"shin-monta-no-mori/server/internal/app"
	db "shin-monta-no-mori/server/internal/db/sqlc"
	"shin-monta-no-mori/server/pkg/token"
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

	// DBのマイグレーションを実行
	runDBMigration(config.MigrationURL, config.DBUrl)

	store := db.NewStore(conn)
	// DBのシードファイルを実行
	if config.Environment == "dev" {
		util.Seeding(store)
	}
	token, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker : %w", err)
	}
	server := app.NewServer(config, store, token)
	server.Router.Use(app.CORSMiddleware(config))

	// Userサイドのルート設定
	api.SetUserRouters(server)
	// Adminサイドのルート設定
	api.SetAdminRouters(server)

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
