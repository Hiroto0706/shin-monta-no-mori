package db

import (
	"database/sql"
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

	dbSource := util.MakeDBSource(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.TestDBName)
	conn, err := sql.Open(config.DBDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db :", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
