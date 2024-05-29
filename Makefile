DB_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori?sslmode=disable
TEST_DB_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori-test?sslmode=disable

.PHONY: createdb
createdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test

.PHONY: dropdb
dropdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres  shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres shin-monta-no-mori-test

.PHONY: new_migration
new_migration:
	migrate create -ext sql -dir server/internal/db/migration -seq $(name)

.PHONY: new_seed
new_seed:
	migrate create -ext sql -dir server/internal/db/seed -seq $(name)

.PHONY: migrateup
migrateup:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose up 1
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up 1

.PHONY: migratedown
migratedown:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path server/internal/db/migration -database "$(DB_URL)" -verbose down 1
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose down 1

.PHONY: seedup
seedup:
	migrate -path server/internal/db/seed -database "$(DB_URL)" -verbose up

.PHONY: dc-up
dc-up:
	docker compose up --build

.PHONY: dc-down
dc-down:
	docker compose down

.PHONY: serve
serve:
	cd ./server && air -c .air.toml

.PHONY: front
front:
	cd ./client && npm run dev

.PHONY: sqlc
sqlc:
	cd server/ && sqlc generate

.PHONY: test
test:
	# テスト実行環境の構築
	docker exec shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up
	mkdir -p coverage

	# 各サブディレクトリのテストを実行し、個別のカバレッジファイルを生成
	go test ./server/api/admin/... -coverprofile=./coverage/api_admin.out
	go test ./server/api/user/... -coverprofile=./coverage/api_user.out
	go test ./server/api/middleware/... -coverprofile=./coverage/api_middleware.out
	go test ./server/pkg/... -coverprofile=./coverage/pkg.out
	go test ./server/internal/db/... -coverprofile=./coverage/db.out
	go test ./server/internal/domains/... -coverprofile=./coverage/domains.out

	# カバレッジファイルの結合
	echo "mode: set" > ./coverage/coverage.out
	tail -n +2 ./coverage/api_admin.out >> ./coverage/coverage.out
	tail -n +2 ./coverage/api_user.out >> ./coverage/coverage.out
	tail -n +2 ./coverage/api_middleware.out >> ./coverage/coverage.out
	tail -n +2 ./coverage/pkg.out >> ./coverage/coverage.out
	tail -n +2 ./coverage/db.out >> ./coverage/coverage.out
	tail -n +2 ./coverage/domains.out >> ./coverage/coverage.out

	# テスト結果の集計・出力
	go tool cover -func=./coverage/coverage.out > ./coverage/report.txt
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	./tools/aggregate_coverage.sh ./coverage/report.txt

# テストが途中で失敗したなどの理由でテスト環境が汚れてしまった時に使う
.PHONY: test-reset
test-reset:
	docker exec shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test
	migrate -path server/internal/db/migration -database "$(TEST_DB_URL)" -verbose up