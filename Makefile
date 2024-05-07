DB_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori?sslmode=disable
TEST_DB_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori-test?sslmode=disable

createdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test

dropdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres  shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres shin-monta-no-mori-test

new_migration:
	migrate create -ext sql -dir server/internal/db/migration -seq $(name)

migrateup:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up

migrateup1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up 1
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up 1

migratedown:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose down

migratedown1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down 1
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose down 1

dc-up:
	docker compose up --build

dc-down:
	docker compose down

serve:
	cd ./server && air -c .air.toml

sqlc:
	cd server/ && sqlc generate

test:
	docker exec shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up
	mkdir -p coverage
	go test ./server/... -coverprofile=./coverage/coverage.out
	go tool cover -func=./coverage/coverage.out > coverage/report.txt
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	./tools/aggregate_coverage.sh ./coverage/report.txt
