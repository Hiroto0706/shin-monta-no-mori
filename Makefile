DB_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori?sslmode=disable

createdb:
	docker exec -it shin-monta_no_mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori

dropdb:
	docker exec -it shin-monta-no-mori-db dropdb shin-monta-no-mori

new_migration:
	migrate create -ext sql -dir server/internal/db/migration -seq $(name)

migrateup:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down 1

dc-up:
	docker compose up --build

dc-down:
	docker compose down