package main

import (
	"database/sql"
	"fmt"
	"log"
	"shin-monta-no-mori/api/handlers"
	"shin-monta-no-mori/pkg/util"

	"github.com/gin-gonic/gin"

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

	conn, err := sql.Open(config.DBDriver, util.MakeDBSource(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	fmt.Println(conn)

	router := gin.Default()
	router.Use(CORSMiddleware())
	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "new Hello World from server.",
	// 	})
	// })
	router.GET("/", handlers.Greet)

	router.Run()
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: local, stg, prdで値を変更する
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:3030"}
		origin := c.GetHeader("Origin")

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

				break
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
